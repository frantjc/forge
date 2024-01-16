package command

import (
	"os"

	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/forgeazure"
	"github.com/frantjc/forge/internal/containerfs"
	"github.com/frantjc/forge/internal/hooks"
	"github.com/frantjc/forge/ore"
	"github.com/frantjc/forge/runtime/docker"
	"github.com/spf13/cobra"
)

// NewTask returns the command which acts as
// the entrypoint for `forge task`.
func NewTask() *cobra.Command {
	var (
		attach, cache      bool
		inputs             map[string]string
		execution, workdir string
		cmd                = &cobra.Command{
			Use:           "task",
			Short:         "Run an Azure DevOps Task",
			Args:          cobra.ExactArgs(1),
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()

				c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
				if err != nil {
					return err
				}

				if attach {
					hooks.ContainerStarted.Listen(hookAttach(cmd, containerfs.WorkingDir))
				}

				t := &ore.Task{
					Task:      args[0],
					Inputs:    inputs,
					Execution: execution,
				}

				var o forge.Ore = t
				if cache {
					o = &ore.Cache{Ore: o}
				}

				return forge.NewFoundry(docker.New(c)).Process(
					ctx,
					o,
					commandDrains(cmd),
				)
			},
		}
	)

	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	cmd.Flags().BoolVarP(&attach, "attach", "a", false, "attach to containers")
	cmd.Flags().BoolVar(&cache, "cache", false, "use cache")
	cmd.Flags().StringToStringVarP(&inputs, "input", "i", nil, "inputs")
	cmd.Flags().StringVar(&forgeazure.NodeImageReference, "node-image", forgeazure.DefaultNodeImageReference, "Node image")
	cmd.Flags().StringVar(&forgeazure.Node10ImageReference, "node10-image", forgeazure.DefaultNode10ImageReference, "Node10 image")
	cmd.Flags().StringVar(&forgeazure.Node16ImageReference, "node16-image", forgeazure.DefaultNode16ImageReference, "Node16 image")
	cmd.Flags().StringVarP(&execution, "exec", "e", "Node", "task execution")
	cmd.Flags().StringVar(&workdir, "workdir", wd, "working directory for use")
	_ = cmd.MarkFlagDirname("workdir")

	return cmd
}
