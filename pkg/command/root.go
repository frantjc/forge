package command

import (
	"context"
	"os"
	"path"
	"runtime"

	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/hooks"
	"github.com/frantjc/forge/pkg/basin/bucket"
	"github.com/frantjc/forge/pkg/concourse"
	"github.com/frantjc/forge/pkg/github/actions"
	"github.com/frantjc/forge/pkg/ore"
	"github.com/frantjc/forge/pkg/runtime/container/docker"
	"github.com/spf13/cobra"
)

func NewRoot() Command {
	var (
		verbosity int
		cmd       = &cobra.Command{
			Use:     "4ge",
			Version: forge.Semver(),
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				cmd.SetContext(forge.WithLogger(cmd.Context(), forge.NewLogger().V(verbosity)))
			},
			Run: func(cmd *cobra.Command, args []string) {
				ctx := cmd.Context()

				wd, err := os.Getwd()
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				globalContext, err := actions.NewGlobalContextFromPath(ctx, wd)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				cache, err := os.UserCacheDir()
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				cache = path.Join(cache, "forge")
				if err = os.MkdirAll(cache, 0777); err != nil {
					cmd.PrintErrln(err)
					return
				}

				basin, err := bucket.New(ctx, "file://"+cache)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				var (
					_                  = globalContext
					_ *concourse.Input = nil
				)

				hooks.ContainerStarted.Listen(func(ctx context.Context, container forge.Container) {
					// streams, restore := forge.StdTerminalStreams()
					// defer func() {
					// 	if err = restore(); err != nil {
					//    cmd.PrintErrln(err)
					//    return
					// 	}
					// }()

					// if _, err = container.Exec(ctx, &forge.ContainerConfig{
					// 	Entrypoint: []string{"sh"},
					// }, streams); err != nil {
					//   cmd.PrintErrln(err)
					//   return
					// }
				})

				foundry := &forge.Foundry{ContainerRuntime: docker.New(c), Basin: basin}
				if _, err = foundry.Process(
					ctx,
					&ore.Alloy{
						Ores: []forge.Ore{
							// &ore.Action{
							// 	Uses: "actions/setup-go@v3",
							// 	With: map[string]string{
							// 		"go-version": "1.19",
							// 	},
							// 	GlobalContext: globalContext,
							// },
							// &ore.Resource{
							// 	Method: "get",
							// 	Resource: &concourse.Resource{
							// 		Name: "github.com/frantjc/forge",
							// 		Type: "git",
							// 		Source: map[string]string{
							// 			"uri":    "https://github.com/frantjc/forge",
							// 			"branch": "main",
							// 		},
							// 	},
							// 	ResourceType: &concourse.ResourceType{
							// 		Name: "git",
							// 		Source: &concourse.Source{
							// 			Repository: "docker.io/concourse/git-resource",
							// 			Tag:        "alpine",
							// 		},
							// 	},
							// },
							&ore.Pure{
								Image:      "alpine",
								Entrypoint: []string{"ls", "-al"},
							},
							&ore.Lava{
								From: &ore.Pure{
									Image:      "alpine",
									Entrypoint: []string{"echo", "hello"},
								},
								To: &ore.Pure{
									Image:      "alpine",
									Entrypoint: []string{"base64"},
								},
							},
						},
					},
					forge.StdDrains(),
				); err != nil {
					cmd.PrintErrln(err)
					return
				}
			},
		}
	)

	cmd.PersistentFlags().CountVarP(&verbosity, "verbose", "v", "verbosity")
	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.AddCommand()

	return cmd
}
