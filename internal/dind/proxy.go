package dind

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	xslice "github.com/frantjc/x/slice"
)

// NewProxy listens on the given listener for traffic intended for `dockerd`.
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
func NewProxy(ctx context.Context, mounts map[string]string, lis net.Listener, dockerSock *url.URL) error {
	var (
		errC    = make(chan error)
		network = "tcp"
		address = dockerSock.Host
	)

	if strings.EqualFold("unix", dockerSock.Scheme) {
		network = "unix"
		address = dockerSock.Path
	}

	go func() {
		errC <- func() error {
			for {
				cli, err := lis.Accept()
				if err != nil {
					return err
				}

				dockerd, err := net.Dial(network, address)
				if err != nil {
					return err
				}

				// Copy responses from `dockerd` straight back to the client.
				go func() {
					// Close the client connection once we
					// reach EOF on the response.
					defer dockerd.Close()
					defer cli.Close()

					errC <- func() error {
						if _, err = io.Copy(cli, dockerd); err != nil {
							return err
						}

						return nil
					}()
				}()

				// Intercept, inspect and potentially modify requests to
				// `dockerd` from the client.
				go func() {
					errC <- func() error {
						buf := bufio.NewReader(cli)

						for {
							req, err := http.ReadRequest(buf)
							if errors.Is(err, io.EOF) {
								break
							} else if err != nil {
								return err
							}

							// TODO: Presumably there are other requests that
							// need intercepted and modified to work properly.
							if req.Method == http.MethodPost && strings.HasSuffix(req.URL.Path, "/containers/create") {
								body := &struct {
									HostConfig       container.HostConfig
									NetworkingConfig map[string]any
									container.Config `json:",inline"`
								}{}

								if err := json.NewDecoder(req.Body).Decode(body); err != nil {
									return err
								}

								if err = req.Body.Close(); err != nil {
									return err
								}

								// Replace requested mount sources with their
								// host path equivalents. Error if impossible.
								//
								// For example, if the client is running in container1 which
								// has mount `/host/path:/container1/path` and requests mount
								// `/container1/path/subpath:/container2/path`, then we modify the
								// request to be for the mount `/host/path/subpath:/container2/path`.
								mountsOk := true
								body.HostConfig.Binds = xslice.Map(body.HostConfig.Binds, func(bind string, _ int) string {
									var (
										parts = strings.SplitN(bind, ":", 2)
										src   = parts[0]
										dst   = parts[1]
									)

									for k, v := range mounts {
										if strings.HasPrefix(src, v) {
											return fmt.Sprintf("%s:%s",
												filepath.Join(
													k, strings.TrimPrefix(src, v),
												),
												dst,
											)
										}

										mountsOk = false
									}

									return bind
								})

								buf := new(bytes.Buffer)

								if !mountsOk {
									if err = json.NewEncoder(buf).Encode(&types.ErrorResponse{
										Message: "one or more requested mounts cannot be satisfied by Forge because it exists inside of the container that Forge is running your process inside of, but not on the host where the Docker daemon is running",
									}); err != nil {
										return err
									}

									if err = (&http.Response{
										Status:        http.StatusText(http.StatusInternalServerError),
										StatusCode:    http.StatusInternalServerError,
										Proto:         req.Proto,
										ProtoMajor:    req.ProtoMajor,
										ProtoMinor:    req.ProtoMinor,
										Body:          io.NopCloser(buf),
										ContentLength: int64(buf.Len()),
										Request:       req,
									}).Write(cli); err != nil {
										return err
									}

									return nil
								}

								if err = json.NewEncoder(buf).Encode(body); err != nil {
									return err
								}

								// Since we possibly modified the request body,
								// the Content-Length has possibly changed.
								req.Body = io.NopCloser(buf)
								req.Header.Set("Content-Length", fmt.Sprint(buf.Len()))
								req.ContentLength = int64(buf.Len())
							}

							if err := req.WriteProxy(dockerd); err != nil {
								return err
							}
						}

						return nil
					}()
				}()
			}
		}()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errC:
			if err != nil {
				return err
			}
		}
	}
}
