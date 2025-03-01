package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/concourse"
	"github.com/frantjc/forge/internal/yaml"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func newResource(method string) *cobra.Command {
	var (
		attach          bool
		conf, workDir   string
		version, params map[string]any
		cmd             = &cobra.Command{
			Use:           fmt.Sprintf("%s [flags] (resource)", method),
			Short:         fmt.Sprintf("%s a Concourse resource", cases.Title(language.English).String(method)),
			Args:          cobra.ExactArgs(1),
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx      = cmd.Context()
					name     = args[0]
					pipeline = &concourse.Pipeline{}
					file     io.Reader
					err      error
					r        = &forge.Resource{
						Method:  method,
						Version: version,
						Params:  params,
					}
				)

				if cmd.Flag("conf").Changed {
					if file, err = os.Open(conf); err != nil {
						return err
					}
				} else {
					if file, err = os.Open(filepath.Join(workDir, conf)); err != nil {
						return err
					}
				}

				if err = yaml.NewDecoder(file).Decode(pipeline); err != nil {
					return err
				}

				for _, p := range pipeline.Resources {
					if p.Name == name {
						r.Resource = &p
					}
				}
				if r.Resource == nil {
					return fmt.Errorf("resource not found: %s", name)
				}

				for _, t := range append(pipeline.ResourceTypes, concourse.BuiltinResourceTypes...) {
					if t.Name == r.Resource.Type {
						resourceType := t
						r.ResourceType = &resourceType
					}
				}
				if r.ResourceType == nil {
					return fmt.Errorf("resource type not found: %s", r.Resource.Type)
				}

				cr, opts, err := runOptsAndContainerRuntime(cmd)
				if err != nil {
					return err
				}

				opts.Mounts = []forge.Mount{
					{
						Source:      workDir,
						Destination: opts.WorkingDir,
					},
				}

				if attach {
					forge.HookContainerStarted.Listen(hookAttach(cmd, opts.WorkingDir))
				}

				return r.Run(ctx, cr, opts)
			},
		}
	)

	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	if method != concourse.MethodCheck {
		cmd.Flags().VarP(newStringToPrimitive(nil, &params), "param", "p", "Params for resource")
	}
	cmd.Flags().BoolVarP(&attach, "attach", "a", false, "Attach to container")
	cmd.Flags().VarP(newStringToPrimitive(nil, &version), "version", "v", "Version for resource")
	cmd.Flags().StringVarP(&conf, "conf", "c", ".forge.yml", "Config file for resource")
	_ = cmd.MarkFlagFilename("conf", ".yaml", ".yml", ".json")
	cmd.Flags().StringVar(&workDir, "workdir", wd, "Working directory for resource")
	_ = cmd.MarkFlagDirname("workdir")

	return cmd
}
