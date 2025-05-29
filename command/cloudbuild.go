package command

import (
	"io"
	"os"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/cloudbuild"
	"github.com/frantjc/forge/internal/envconv"
	"github.com/frantjc/forge/internal/hostfs"
	"github.com/spf13/cobra"
)

// NewCloudBuild returns the command which acts as
// the entrypoint for `forge cloudbuild`.
func NewCloudBuild() *cobra.Command {
	var (
		attach        bool
		script        string
		substitutions map[string]string
		cb            = &forge.CloudBuild{}
		cmd           = setCommon(&cobra.Command{
			Use:     "cloudbuild [flags] (builder) [--] [args]",
			Aliases: []string{"cb"},
			Short:   "Run a Google Cloud Build step",
			Args:    cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx   = cmd.Context()
					iArgs = 1
				)

				wd, err := os.Getwd()
				if err != nil {
					return err
				}

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

				subs, err := cloudbuild.NewSubstituionsFromPath(wd, substitutions)
				if err != nil {
					if subs, err = cloudbuild.NewSubstitutionsFromEnv(substitutions); err != nil {
						return err
					}
				}

				cb.Substitutions = envconv.ArrToMap(subs.Env())

				cr, opts, err := runOptsAndContainerRuntime(cmd)
				if err != nil {
					return err
				}

				opts.Mounts = []forge.Mount{
					{
						Source:      wd,
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

				return cb.Run(ctx, cr, opts)
			},
		})
	)

	cmd.Flags().BoolVarP(&attach, "attach", "a", false, "Attach to container before executing cloudbuild")
	cmd.Flags().StringVar(&cb.Entrypoint, "entrypoint", "", "Entrypoint for cloudbuild")
	cmd.Flags().StringVar(&script, "script", "", "Script for cloudbuild")
	_ = cmd.MarkFlagFilename("script")
	cmd.Flags().StringArrayVarP(&cb.Env, "env", "e", nil, "Env for cloudbuild")
	cmd.Flags().StringToStringVarP(&substitutions, "sub", "s", nil, "Substitutions for cloudbuild")
	cmd.Flags().BoolVar(&cb.AutomapSubstitutions, "automap-substitutions", false, "Automap substitutions for cloudbuild")

	return cmd
}
