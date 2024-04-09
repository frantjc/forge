package dind

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
				var (
					statusCode = http.StatusInternalServerError
					errC2      = make(chan error, 1)
				)

				conn, cli, err := func() (net.Conn, *bufio.ReadWriter, error) {
					h, ok := w.(http.Hijacker)
					if !ok {
						return nil, nil, fmt.Errorf("not a hijacker")
					}

					return h.Hijack()
				}()
				if err != nil {
					_ = json.NewEncoder(w).Encode(&types.ErrorResponse{
						Message: err.Error(),
					})
					return
				}

				closeConn := sync.OnceFunc(func() {
					_ = cli.Flush()
					_ = conn.Close()
				})
				defer closeConn()

				if err = func() error {
					daemon, err := net.Dial(network, address)
					if err != nil {
						return err
					}

					// Copy responses from `dockerd` straight back to the client.
					go func() {
						// Close the client connection once we
						// reach EOF on the response.
						defer closeConn()
						defer daemon.Close()

						buf := bufio.NewReader(daemon)

						errC2 <- func() error {
							res, err := http.ReadResponse(buf, r)
							if err != nil {
								return err
							}

							if err = res.Write(io.MultiWriter(os.Stdout, cli)); err != nil {
								return err
							}

							return nil
						}()
					}()

					// Intercept, inspect and potentially modify requests to
					// `dockerd` from the client.
					// TODO: Presumably there are other requests that
					// need intercepted and modified to work properly.
					if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/containers/create") {
						body := &struct {
							HostConfig       container.HostConfig
							NetworkingConfig map[string]any
							container.Config `json:",inline"`
						}{}

						if err := json.NewDecoder(cli).Decode(body); err != nil {
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
								parts = strings.SplitN(bind, ":", 2)
								src   = parts[0]
								dst   = parts[1]
							)

							for k, v := range mounts {
								if strings.HasPrefix(src, v) {
									body.HostConfig.Binds[i] = fmt.Sprintf("%s:%s",
										filepath.Join(
											k, strings.TrimPrefix(src, v),
										),
										dst,
									)
									continue
								}
							}

							return fmt.Errorf("one or more requested mounts cannot be satisfied by Forge because it exists inside of the container that Forge is running your process inside of, but not on the host where the Docker daemon is running")
						}

						buf := new(bytes.Buffer)

						if err = json.NewEncoder(buf).Encode(body); err != nil {
							return err
						}

						// Since we possibly modified the request body,
						// the Content-Length has possibly changed.
						lenBuf := buf.Len()
						r.Body = io.NopCloser(buf)
						r.Header.Set("Content-Length", fmt.Sprint(lenBuf))
						r.ContentLength = int64(lenBuf)
					}

					if err := r.WriteProxy(io.MultiWriter(os.Stdout, daemon)); err != nil {
						return err
					}

					return <-errC2
				}(); err != nil {
					buf := new(bytes.Buffer)

					_ = json.NewEncoder(buf).Encode(&types.ErrorResponse{
						Message: err.Error(),
					})

					_ = (&http.Response{
						Status:     http.StatusText(statusCode),
						StatusCode: statusCode,
						Proto:      r.Proto,
						ProtoMajor: r.ProtoMajor,
						ProtoMinor: r.ProtoMajor,
						Request:    r,
						Body:       io.NopCloser(buf),
					}).Write(cli)
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
