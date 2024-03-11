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

	"github.com/docker/docker/api/types/container"
	xslice "github.com/frantjc/x/slice"
)

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

				moby, err := net.Dial(network, address)
				if err != nil {
					return err
				}

				// Copy responses from the Docker daemon straight back to the client.
				go func() {
					errC <- func() error {
						if _, err = io.Copy(cli, moby); err != nil {
							return err
						}

						return nil
					}()
				}()

				go func() {
					buf := bufio.NewReader(cli)

					errC <- func() error {
						for {
							req, err := http.ReadRequest(buf)
							if errors.Is(err, io.EOF) {
								break
							} else if err != nil {
								return err
							}

							if req.Method == http.MethodPost && strings.HasSuffix(req.URL.Path, "/containers/create") {
								body := &struct {
									HostConfig *container.HostConfig
									Rest       json.RawMessage
								}{
									HostConfig: &container.HostConfig{},
								}

								if err := json.NewDecoder(req.Body).Decode(body); err != nil {
									return err
								}

								if err = req.Body.Close(); err != nil {
									return err
								}

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
									}

									return bind
								})
								buf := new(bytes.Buffer)

								if err = json.NewEncoder(buf).Encode(body); err != nil {
									return err
								}

								req.Body = io.NopCloser(buf)
							}

							if err := req.Write(moby); err != nil {
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
