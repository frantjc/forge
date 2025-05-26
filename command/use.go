package command

import (
	"encoding/json"
	"os"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/githubactions"
	"github.com/frantjc/forge/internal/hostfs"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// NewUse returns the command which acts as
// the entrypoint for `forge use`.
func NewUse() *cobra.Command {
	var (
		attach, outputs, envVars bool
		env, with                map[string]string
		debug                    bool
		cmd                      = setCommon(&cobra.Command{
			Use:     "use [flags] (action)",
			Aliases: []string{"github", "action", "act", "gh"},
			Short:   "Use a GitHub Action",
			Args:    cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx = cmd.Context()
					a   = &forge.Action{
						ID:   uuid.NewString(),
						Uses: args[0],
						With: with,
						Env:  env,
					}
				)

				wd, err := os.Getwd()
				if err != nil {
					return err
				}

				a.GlobalContext, err = githubactions.NewGlobalContextFromPath(wd)
				if err != nil {
					a.GlobalContext = githubactions.NewGlobalContextFromEnv()
				}

				if debug {
					a.GlobalContext.EnableDebug()
				}

				for _, dir := range []string{hostfs.RunnerTmp, hostfs.RunnerToolCache} {
					if err = os.MkdirAll(dir, 0o755); err != nil {
						return err
					}
				}

				cr, opts, err := runOptsAndContainerRuntime(cmd, envVars, outputs)
				if err != nil {
					return err
				}

				opts.Mounts = []forge.Mount{
					{
						Source:      wd,
						Destination: forge.GitHubWorkspace(opts.WorkingDir),
					},
					{
						Source:      hostfs.RunnerTmp,
						Destination: forge.GitHubRunnerTmp(opts.WorkingDir),
					},
					{
						Source:      hostfs.RunnerToolCache,
						Destination: forge.GitHubRunnerToolCache(opts.WorkingDir),
					},
				}

				if attach {
					forge.HookContainerStarted.Listen(hookAttach(cmd, opts.WorkingDir, envVars, outputs))
				}

				if envVars {
					defer func() {
						_ = json.NewEncoder(cmd.OutOrStdout()).Encode(a.GlobalContext.EnvContext)
					}()
				}

				var ()

				if outputs {
					defer func() {
						_ = json.NewEncoder(cmd.OutOrStdout()).Encode(a.GlobalContext.StepsContext[a.ID].Outputs)
					}()
				}

				return a.Run(ctx, cr, opts)
			},
		})
	)

	cmd.Flags().BoolVarP(&attach, "attach", "a", false, "Attach to containers before executing action")
	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "Print debug logs")
	cmd.Flags().BoolVar(&outputs, "outputs", false, "Print step outputs")
	cmd.Flags().BoolVar(&envVars, "env-vars", false, "Print step environment variables")
	cmd.Flags().StringToStringVarP(&env, "env", "e", nil, "Env values for use")
	cmd.Flags().StringToStringVarP(&with, "with", "w", nil, "With values for use")
	cmd.Flags().StringVar(&forge.Node12ImageReference, "node12-image", forge.DefaultNode12ImageReference, "The node12 container image for use")
	cmd.Flags().StringVar(&forge.Node16ImageReference, "node16-image", forge.DefaultNode16ImageReference, "The node16 container image for use")
	cmd.Flags().StringVar(&forge.Node20ImageReference, "node20-image", forge.DefaultNode20ImageReference, "The node20 container image for use")

	return cmd
}
