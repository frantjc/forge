package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/concourse"
	"github.com/frantjc/forge/forgeconcourse"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/frantjc/forge/internal/hooks"
	"github.com/frantjc/forge/internal/yaml"
	"github.com/frantjc/forge/ore"
	"github.com/frantjc/forge/runtime/docker"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func newResource(method string) *cobra.Command {
	var (
		attach          bool
		conf, workdir   string
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
					cr       = &ore.Resource{
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
					if file, err = os.Open(filepath.Join(workdir, conf)); err != nil {
						return err
					}
				}

				if err = yaml.NewDecoder(file).Decode(pipeline); err != nil {
					return err
				}

				for _, r := range pipeline.Resources {
					if r.Name == name {
						resource := r
						cr.Resource = &resource
					}
				}
				if cr.Resource == nil {
					return fmt.Errorf("resource not found: %s", name)
				}

				for _, t := range append(pipeline.ResourceTypes, concourse.BuiltinResourceTypes...) {
					if t.Name == cr.Resource.Type {
						resourceType := t
						cr.ResourceType = &resourceType
					}
				}
				if cr.ResourceType == nil {
					return fmt.Errorf("resource type not found: %s", cr.Resource.Type)
				}

				c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
				if err != nil {
					return err
				}

				destination := filepath.Join(forgeconcourse.DefaultRootPath, cr.Resource.Name)

				if attach {
					hooks.ContainerStarted.Listen(hookAttach(cmd, destination))
				}

				return forge.NewFoundry(docker.New(c, !cmd.Flag("no-dind").Changed)).Process(
					contaminate.WithMounts(ctx, forge.Mount{
						Source:      workdir,
						Destination: destination,
					}),
					cr,
					commandDrains(cmd),
				)
			},
		}
	)

	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	if method != "check" {
		cmd.Flags().VarP(newStringToPrimitive(nil, &params), "param", "p", "params for resource")
	}
	cmd.Flags().BoolVarP(&attach, "attach", "a", false, "attach to containers")
	cmd.Flags().VarP(newStringToPrimitive(nil, &version), "version", "v", "version for resource")
	cmd.Flags().StringVarP(&conf, "conf", "c", ".forge.yml", "config file for resource")
	_ = cmd.MarkFlagFilename("conf", ".yaml", ".yml", ".json")
	cmd.Flags().StringVar(&workdir, "workdir", wd, "working directory for resource")
	_ = cmd.MarkFlagDirname("workdir")

	return cmd
}
