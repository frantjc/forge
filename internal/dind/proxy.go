package dind

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/volume"
)

// ServeDockerdProxy listens on the given listener for traffic intended for `dockerd`.
// It proxies responses back from `dockerd` directly, but modifies requests
// to `dockerd` to try to translate them for the host `dockerd` when the
// client is calling from inside of a container.
//
// As of writing, it only does this by modifying CreateContainer HostConfig.Binds
// to use the host path equivalent of the client's mount source. Note that this
// only works if the client's mount source is mounted from the host, which, in
// `forge`, is often the case. Unfortunately, it can't support _everything_.
//
// It always returns an error and doesn't exit until the given context.Context
// is done or an error is encountered, similar to http.Serve.
func ServeDockerdProxy(ctx context.Context, mounts map[string]string, lis net.Listener, dockerSock *url.URL) error { //nolint: gocyclo
	var (
		network = "tcp"
		address = dockerSock.Host
		errC    = make(chan error, 1)
	)

	if strings.EqualFold("unix", dockerSock.Scheme) {
		network = "unix"
		address = dockerSock.Path
	}

	var (
		director = func(r *http.Request) {
			if r.URL.Scheme != "https" {
				r.URL.Scheme = "http"
			}

			if r.URL.Host == "" {
				r.URL.Host = "api.moby.localhost"
			}
		}
		// Return errors to the client instead of logging them server-side.
		errorLog     = log.New(io.Discard, "", 0)
		errorHandler = func(w http.ResponseWriter, _ *http.Request, err error) {
			w.WriteHeader(http.StatusBadGateway)

			_ = json.NewEncoder(w).Encode(&types.ErrorResponse{
				Message: err.Error(),
			})
		}
		dialer = &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
	)

	go func() {
		errC <- (&http.Server{
			ReadHeaderTimeout: time.Second * 5,
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				errStatusCode := http.StatusInternalServerError

				if err := func() error {
					if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/containers/create") {
						body := &struct {
							HostConfig       container.HostConfig
							NetworkingConfig map[string]any
							container.Config `json:",inline"`
						}{}

						if err := json.NewDecoder(r.Body).Decode(body); err != nil {
							errStatusCode = http.StatusBadRequest
							return err
						}

						// Replace requested mount sources with their
						// host path equivalents. Error if impossible.
						//
						// For example, if the client is running in container1 which
						// has mount `/host/path:/container1/path` and requests mount
						// `/container1/path/subpath:/container2/path`, then we modify the
						// request to be for the mount `/host/path/subpath:/container2/path`.
						for i, bind := range body.HostConfig.Binds {
							parts := strings.SplitN(bind, ":", 2)

							if len(parts) != 2 {
								errStatusCode = http.StatusBadRequest
								return fmt.Errorf("invalid volume `%s`", bind)
							}

							var (
								src       = parts[0]
								dst       = parts[1]
								satisfied bool
							)

							for k, v := range mounts {
								if strings.HasPrefix(src, v) {
									body.HostConfig.Binds[i] = fmt.Sprintf("%s:%s",
										filepath.Join(
											k, strings.TrimPrefix(src, v),
										),
										dst,
									)
									satisfied = true
									break
								}
							}

							if !satisfied {
								errStatusCode = http.StatusBadRequest
								return fmt.Errorf("volume `%s` cannot be satisfied by Forge because it does not exist on the host where the Docker daemon is running", bind)
							}
						}

						buf := new(bytes.Buffer)

						if err := json.NewEncoder(buf).Encode(body); err != nil {
							return err
						}

						// Since we possibly modified the request body,
						// the Content-Length has possibly changed.
						lenBuf := buf.Len()
						r.Body = io.NopCloser(buf)
						r.Header.Set("Content-Length", fmt.Sprint(lenBuf))
						r.ContentLength = int64(lenBuf)
					} else if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/volumes/create") {
						body := &volume.CreateOptions{}

						if err := json.NewDecoder(r.Body).Decode(body); err != nil {
							errStatusCode = http.StatusBadRequest
							return err
						}

						if strings.EqualFold("local", body.Driver) && strings.EqualFold("none", body.DriverOpts["type"]) {
							if device, ok := body.DriverOpts["device"]; ok {
								var satisfied bool

								for k, v := range mounts {
									if strings.HasPrefix(device, v) {
										body.DriverOpts["device"] = filepath.Join(
											k, strings.TrimPrefix(device, v),
										)
										satisfied = true
										break
									}
								}

								if !satisfied {
									errStatusCode = http.StatusBadRequest
									return fmt.Errorf("volume `%s` cannot be satisfied by Forge because it does not exist on the host where the Docker daemon is running", device)
								}
							}
						}

						buf := new(bytes.Buffer)

						if err := json.NewEncoder(buf).Encode(body); err != nil {
							return err
						}

						// Since we possibly modified the request body,
						// the Content-Length has possibly changed.
						lenBuf := buf.Len()
						r.Body = io.NopCloser(buf)
						r.Header.Set("Content-Length", fmt.Sprint(lenBuf))
						r.ContentLength = int64(lenBuf)
					}

					daemon, err := dialer.DialContext(r.Context(), network, address)
					if err != nil {
						return err
					}
					closeDaemon := sync.OnceFunc(func() {
						_ = daemon.Close()
					})
					defer closeDaemon()

					if strings.HasSuffix(r.URL.Path, "/attach") {
						h, ok := w.(http.Hijacker)
						if !ok {
							return fmt.Errorf("not a hijacker")
						}

						cli, buf, err := h.Hijack()
						if err != nil {
							return err
						}
						closeCli := sync.OnceFunc(func() {
							_ = buf.Flush()
							_ = cli.Close()
						})
						defer closeCli()

						_errC := make(chan error, 1)

						// Copy responses from `dockerd` straight back to the client.
						go func() {
							// Close the client connection once we
							// reach EOF on the response.
							defer closeCli()
							defer closeDaemon()

							bufD := bufio.NewReader(daemon)

							_errC <- func() error {
								res, err := http.ReadResponse(bufD, r)
								if err != nil {
									return err
								}

								if err = res.Write(cli); err != nil {
									return err
								}

								if _, err = io.Copy(buf, daemon); err != nil {
									return err
								}

								return nil
							}()
						}()

						if err = r.Write(daemon); err != nil {
							return err
						}

						if err := <-_errC; err != nil {
							return err
						}
					} else {
						transport, ok := http.DefaultTransport.(*http.Transport)
						if !ok {
							return fmt.Errorf("default transport is not a transport")
						}

						// We want to use http.DefaultTransport, but it won't dial non-http schemes by default,
						// so convince it that it's not and override the dialer to always connect to `dockerd`.
						transport.DialContext = func(_ context.Context, _, _ string) (net.Conn, error) {
							return daemon, nil
						}

						(&httputil.ReverseProxy{
							// This director makes http.DefaultTransport happy
							// just long enough for it to use our transport.
							Director:  director,
							Transport: transport,
							// Return errors to the client instead of logging them server-side.
							ErrorLog:     errorLog,
							ErrorHandler: errorHandler,
						}).ServeHTTP(w, r)
					}

					return nil
				}(); err != nil {
					w.WriteHeader(errStatusCode)

					_ = json.NewEncoder(w).Encode(&types.ErrorResponse{
						Message: err.Error(),
					})
				}
			}),
		}).Serve(lis)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errC:
		return err
	}
}
