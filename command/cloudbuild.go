package command

import (
	"io"
	"os"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/cloudbuild"
	"github.com/frantjc/forge/envconv"
	"github.com/frantjc/forge/internal/hostfs"
	"github.com/spf13/cobra"
)

// NewCloudBuild returns the command which acts as
// the entrypoint for `forge cloudbuild`.
func NewCloudBuild() *cobra.Command {
	var (
		attach        bool
		workDir       string
		script        string
		substitutions map[string]string
		cb            = &forge.CloudBuild{}
		cmd           = &cobra.Command{
			Use:           "cloudbuild [flags] (builder) [--] [args]",
			Aliases:       []string{"cb"},
			Short:         "Run a Google Cloud Build step",
			Args:          cobra.MinimumNArgs(1),
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx   = cmd.Context()
					iArgs = 1
				)

				cb.Name = args[0]

				if len(args) > 1 {
					if args[1] == "--" {
						iArgs++
					}
				}

				cb.Args = args[iArgs:]

				if script != "" {
					var r io.Reader
					if script == "-" {
						r = cmd.InOrStdin()
					} else {
						f, err := os.Open(script)
						if err != nil {
							return err
						}
						defer f.Close()

						r = f
					}

					b, err := io.ReadAll(r)
					if err != nil {
						return err
					}

					cb.Script = string(b)
				}

				for _, dir := range []string{hostfs.CloudBuildWorkspace} {
					if err := os.MkdirAll(dir, 0o755); err != nil {
						return err
					}
				}

				subs, err := cloudbuild.NewSubstituionsFromPath(workDir, substitutions)
				if err != nil {
					if subs, err = cloudbuild.NewSubstitutionsFromEnv(substitutions); err != nil {
						return err
					}
				}

				cb.Substitutions = envconv.ArrToMap(subs.Env())

				cr, opts, err := oreOptsAndContainerRuntime(cmd)
				if err != nil {
					return err
				}

				opts.Mounts = []forge.Mount{
					{
						Source:      workDir,
						Destination: forge.CloudBuildWorkingDir(opts.WorkingDir),
					},
					{
						Source:      hostfs.CloudBuildWorkspace,
						Destination: cloudbuild.WorkspacePath,
					},
				}

				if attach {
					forge.HookContainerStarted.Listen(hookAttach(cmd, forge.CloudBuildWorkingDir(opts.WorkingDir)))
				}

				return cb.Liquify(ctx, cr, opts)
			},
		}
	)

	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	cmd.Flags().BoolVarP(&attach, "attach", "a", false, "attach to containers")
	cmd.Flags().StringVar(&workDir, "workdir", wd, "working directory for cloudbuild")
	_ = cmd.MarkFlagDirname("workdir")

	cmd.Flags().StringVar(&cb.Entrypoint, "entrypoint", "", "entrypoint for cloudbuild")
	cmd.Flags().StringVar(&script, "script", "", "script for cloudbuild")
	_ = cmd.MarkFlagFilename("script")
	cmd.Flags().StringArrayVarP(&cb.Env, "env", "e", nil, "env for cloudbuild")
	cmd.Flags().StringToStringVarP(&substitutions, "sub", "s", nil, "substitutions for cloudbuild")
	cmd.Flags().BoolVar(&cb.AutomapSubstitutions, "automap-substitutions", false, "automap substitutions for cloudbuild")

	return cmd
}
