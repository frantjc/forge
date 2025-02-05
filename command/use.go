package command

import (
	"encoding/json"
	"os"
	"strconv"

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
		workDir                  string
		env, with                map[string]string
		cmd                      = &cobra.Command{
			Use:           "use [flags] (action)",
			Aliases:       []string{"github", "action", "act", "gh"},
			Short:         "Use a GitHub Action",
			Args:          cobra.ExactArgs(1),
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(cmd *cobra.Command, args []string) error {
				globalContext, err := githubactions.NewGlobalContextFromPath(workDir)
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

				cr, opts, err := runOptsAndContainerRuntime(cmd)
				if err != nil {
					return err
				}

				opts.Mounts = []forge.Mount{
					{
						Source:      workDir,
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
						_ = json.NewEncoder(cmd.OutOrStdout()).Encode(globalContext.EnvContext)
					}()
				}

				var (
					ctx = cmd.Context()
					a   = &forge.Action{
						ID:            uuid.NewString(),
						Uses:          args[0],
						With:          with,
						Env:           env,
						GlobalContext: globalContext,
					}
				)

				if outputs {
					defer func() {
						_ = json.NewEncoder(cmd.OutOrStdout()).Encode(globalContext.StepsContext[a.ID].Outputs)
					}()
				}

				return a.Run(ctx, cr, opts)
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
	cmd.Flags().StringVar(&forge.Node12ImageReference, "node12-image", forge.DefaultNode12ImageReference, "node12 image")
	cmd.Flags().StringVar(&forge.Node16ImageReference, "node16-image", forge.DefaultNode16ImageReference, "node16 image")
	cmd.Flags().StringVar(&forge.Node20ImageReference, "node20-image", forge.DefaultNode20ImageReference, "node20 image")
	cmd.Flags().StringVar(&workDir, "workdir", wd, "working directory for use")
	_ = cmd.MarkFlagDirname("workdir")

	return cmd
}
