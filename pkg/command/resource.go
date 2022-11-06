package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/frantjc/forge/pkg/forgeconcourse"
	"github.com/frantjc/forge/pkg/ore"
	"github.com/frantjc/forge/pkg/runtime/docker"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func newNonCheckResource(method string) *cobra.Command {
	var (
		params map[string]string
		cmd    = newResource(method)
	)

	cmd.Flags().StringToStringVarP(&params, "param", "p", nil, "params for resource")

	return cmd
}

func newResource(method string) *cobra.Command {
	var (
		conf            string
		params, version map[string]string
		cmd             = &cobra.Command{
			Use:           method,
			Short:         fmt.Sprintf("%s a Concourse resource", cases.Title(language.English).String(method)),
			Args:          cobra.ExactArgs(1),
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx    = cmd.Context()
					_      = forge.LoggerFrom(ctx)
					name   = args[0]
					config = &forgeconcourse.Config{}
					wd     = WorkdirFrom(ctx)
					file   io.Reader
					err    error
				)

				if filepath.IsAbs(conf) {
					if file, err = os.Open(conf); err != nil {
						return err
					}
				} else {
					if file, err = os.Open(filepath.Join(wd, conf)); err != nil {
						return err
					}
				}

				if err = yaml.NewDecoder(file).Decode(config); err != nil {
					return err
				}

				o := &ore.Resource{
					Method:  method,
					Version: version,
					Params:  params,
				}
				for _, r := range config.GetResources() {
					if r.GetName() == name {
						o.Resource = r
					}
				}
				if o.GetResource() == nil {
					return fmt.Errorf("resource not found: %s", name)
				}

				for _, t := range config.GetResourceTypes() {
					if t.GetName() == o.GetResource().GetType() {
						o.ResourceType = t
					}
				}
				if o.GetResourceType() == nil {
					return fmt.Errorf("resource type not found: %s", o.GetResource().GetType())
				}

				c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
				if err != nil {
					return err
				}

				_, err = forge.NewFoundry(docker.New(c)).Process(
					contaminate.WithMounts(ctx, &forge.Mount{
						Source:      wd,
						Destination: filepath.Join(forgeconcourse.DefaultRootPath, o.GetResource().GetName()),
					}), o, forge.StdDrains(),
				)
				return err
			},
		}
	)

	cmd.Flags().StringToStringVarP(&version, "version", "i", nil, "version for resource")
	cmd.Flags().StringVarP(&conf, "conf", "c", "forge.yml", "config file for resource")
	_ = cmd.MarkFlagFilename("conf")

	return cmd
}
