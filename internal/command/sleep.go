package command

import (
	"net"
	"net/url"
	"os"
	"path/filepath"

	"github.com/frantjc/forge/internal/dind"
	"github.com/spf13/cobra"
)

// NewSleep returns the command which acts as
// the entrypoint for `shim sleep`.
func NewSleep() *cobra.Command {
	var (
		workingdir string
		mounts     map[string]string
		cmd        = &cobra.Command{
			Use:           "sleep",
			Short:         "Sleep until signalled",
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(cmd *cobra.Command, _ []string) error {
				ctx := cmd.Context()

				if len(mounts) > 0 && workingdir != "" {
					forgeSock := filepath.Join(workingdir, "forge.sock")

					if lis, err := net.Listen("unix", forgeSock); err == nil {
						defer os.Remove(forgeSock)

						if dockerHost := os.Getenv("DOCKER_HOST"); dockerHost != "" {
							if dockerSock, err := url.Parse(dockerHost); err == nil {
								return dind.NewProxy(ctx, mounts, lis, dockerSock)
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
	cmd.Flags().StringVar(&workingdir, "wd", "", "working directory for forge")
	_ = cmd.MarkFlagDirname("wd")

	return cmd
}

