package command

import (
	"net"
	"net/url"
	"os"

	"github.com/docker/docker/client"
	"github.com/frantjc/forge/internal/dind"
	"github.com/spf13/cobra"
)

// NewSleep returns the command which acts as
// the entrypoint for `shim sleep`.
func NewSleep() *cobra.Command {
	var (
		forgeSock string
		mounts    map[string]string
		cmd       = &cobra.Command{
			Use:           "sleep",
			Short:         "Sleep until signalled",
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(cmd *cobra.Command, _ []string) error {
				ctx := cmd.Context()

				if forgeSock != "" {
					if lis, err := net.Listen("unix", forgeSock); err == nil {
						defer os.Remove(forgeSock)

						if dockerHost := os.Getenv(client.EnvOverrideHost); dockerHost != "" {
							if dockerSock, err := url.Parse(dockerHost); err == nil {
								return dind.ServeDockerdProxy(ctx, mounts, lis, dockerSock)
							}
						}
					}
				}

				<-ctx.Done()

				return ctx.Err()
			},
		}
	)

	cmd.Flags().StringToStringVar(&mounts, "mount", nil, "mounts for forge")
	cmd.Flags().StringVar(&forgeSock, "sock", "", "unix socket for forge")
	_ = cmd.MarkFlagFilename("sock", ".sock")

	return cmd
}
