package command

import (
	"os"

	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/frantjc/forge/internal/hostfs"
	"github.com/frantjc/forge/pkg/forgeactions"
	"github.com/frantjc/forge/pkg/github/actions"
	"github.com/frantjc/forge/pkg/ore"
	"github.com/frantjc/forge/pkg/runtime/docker"
	"github.com/spf13/cobra"
)

func NewUse() *cobra.Command {
	var (
		env, with map[string]string
		cmd       = &cobra.Command{
			Use:           "use",
			Short:         "Use a GitHub Action",
			Args:          cobra.ExactArgs(1),
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx = cmd.Context()
					wd  = WorkdirFrom(ctx)
				)

				globalContext, err := actions.NewGlobalContextFromPath(ctx, wd)
				if err != nil {
					return err
				}

				c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
				if err != nil {
					return err
				}

				for _, dir := range []string{hostfs.RunnerTmp, hostfs.RunnerToolcache} {
					if err = os.MkdirAll(dir, 0755); err != nil {
						return err
					}
				}

				_, err = forge.NewFoundry(docker.New(c)).Process(
					contaminate.WithMounts(ctx, []*forge.Mount{
						{
							Source:      wd,
							Destination: forgeactions.DefaultWorkspace,
						},
						{
							Source:      hostfs.RunnerTmp,
							Destination: forgeactions.DefaultRunnerTemp,
						},
						{
							Source:      hostfs.RunnerToolcache,
							Destination: forgeactions.DefaultRunnerToolCache,
						},
					}...),
					&ore.Action{
						Uses:          args[0],
						With:          with,
						Env:           env,
						GlobalContext: globalContext,
					},
					forge.StdDrains(),
				)
				return err
			},
		}
	)

	cmd.Flags().StringToStringVarP(&with, "env", "e", nil, "env values")
	cmd.Flags().StringToStringVarP(&with, "with", "w", nil, "with values")

	return cmd
}
