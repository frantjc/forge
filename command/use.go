package command

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/forgeactions"
	"github.com/frantjc/forge/githubactions"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/frantjc/forge/internal/hooks"
	"github.com/frantjc/forge/internal/hostfs"
	"github.com/frantjc/forge/ore"
	"github.com/frantjc/forge/runtime/docker"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// NewUse returns the command which acts as
// the entrypoint for `forge use`.
func NewUse() *cobra.Command {
	var (
		attach, outputs, envVars bool
		workdir                  string
		env, with                map[string]string
		cmd                      = &cobra.Command{
			Use:           "use [flags] (action)",
			Aliases:       []string{"github", "action", "act", "gh"},
			Short:         "Use a GitHub Action",
			Args:          cobra.ExactArgs(1),
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()

				globalContext, err := githubactions.NewGlobalContextFromPath(workdir)
				if err != nil {
					globalContext = githubactions.NewGlobalContextFromEnv()
				}

				if verbosity, _ := strconv.Atoi(cmd.Flag("verbose").Value.String()); verbosity > 0 {
					globalContext.EnableDebug()
				}

				for _, dir := range []string{hostfs.RunnerTmp, hostfs.RunnerToolCache} {
					if err = os.MkdirAll(dir, 0o755); err != nil {
						return err
					}
				}

				c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
				if err != nil {
					return err
				}

				if attach {
					hooks.ContainerStarted.Listen(hookAttach(cmd, forgeactions.DefaultMapping.Workspace))
				}

				a := &ore.Action{
					ID:            uuid.NewString(),
					Uses:          args[0],
					With:          with,
					Env:           env,
					GlobalContext: globalContext,
				}

				if outputs {
					defer func() {
						_ = json.NewEncoder(cmd.OutOrStdout()).Encode(globalContext.StepsContext[a.ID].Outputs)
					}()
				}

				if envVars {
					defer func() {
						_ = json.NewEncoder(cmd.OutOrStdout()).Encode(globalContext.EnvContext)
					}()
				}

				return forge.NewFoundry(docker.New(c, !cmd.Flag("no-dind").Changed)).Process(
					contaminate.WithMounts(ctx, []forge.Mount{
						{
							Source:      workdir,
							Destination: forgeactions.DefaultWorkspace,
						},
						{
							Source:      hostfs.RunnerTmp,
							Destination: forgeactions.DefaultRunnerTemp,
						},
						{
							Source:      hostfs.RunnerToolCache,
							Destination: forgeactions.DefaultRunnerToolCache,
						},
					}...),
					a,
					commandDrains(cmd, outputs, envVars),
				)
			},
		}
	)

	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	cmd.Flags().BoolVarP(&attach, "attach", "a", false, "attach to containers")
	cmd.Flags().BoolVar(&outputs, "outputs", false, "print step outputs")
	cmd.Flags().BoolVar(&envVars, "env-vars", false, "print step environment variables")
	cmd.Flags().StringToStringVarP(&env, "env", "e", nil, "env values")
	cmd.Flags().StringToStringVarP(&with, "with", "w", nil, "with values")
	cmd.Flags().StringVar(&forgeactions.Node12ImageReference, "node12-image", forgeactions.DefaultNode12ImageReference, "node12 image")
	cmd.Flags().StringVar(&forgeactions.Node16ImageReference, "node16-image", forgeactions.DefaultNode16ImageReference, "node16 image")
	cmd.Flags().StringVar(&forgeactions.Node20ImageReference, "node20-image", forgeactions.DefaultNode20ImageReference, "node20 image")
	cmd.Flags().StringVar(&workdir, "workdir", wd, "working directory for use")
	_ = cmd.MarkFlagDirname("workdir")

	return cmd
}
