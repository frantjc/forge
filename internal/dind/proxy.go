package dind

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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
func ServeDockerdProxy(ctx context.Context, mounts map[string]string, lis net.Listener, dockerSock *url.URL) error {
	var (
		network = "tcp"
		address = dockerSock.Host
		errC    = make(chan error, 1)
	)

	if strings.EqualFold("unix", dockerSock.Scheme) {
		network = "unix"
		address = dockerSock.Path
	}

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
							var (
								parts     = strings.SplitN(bind, ":", 2)
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
								return fmt.Errorf("volume `%s` cannot be satisfied by Forge because it exists inside of the container that Forge is running your process inside of, but not on the host where the Docker daemon is running", bind)
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

					var (
						dialer = &net.Dialer{
							Timeout:   30 * time.Second,
							KeepAlive: 30 * time.Second,
						}
						transport, ok = http.DefaultTransport.(*http.Transport)
					)
					if !ok {
						return fmt.Errorf("default transport is not a transport")
					}

					transport.DialContext = func(ctx context.Context, _, _ string) (net.Conn, error) {
						return dialer.DialContext(ctx, network, address)
					}

					(&httputil.ReverseProxy{
						Director: func(r *http.Request) {
							fmt.Println(r.Header)
							r.URL.Scheme = "http"

							if r.Host == "" {
								r.Host = "api.moby.localhost"
							}

							r.URL.Host = r.Host
							r.Header.Set("Host", r.Host)
						},
						Transport: transport,
					}).ServeHTTP(w, r)

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
