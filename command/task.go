package command

import (
	"os"

	"github.com/frantjc/forge"
	"github.com/spf13/cobra"
)

// NewTask returns the command which acts as
// the entrypoint for `forge task`.
func NewTask() *cobra.Command {
	var (
		attach             bool
		inputs             map[string]string
		execution, workdir string
		cmd                = &cobra.Command{
			Use:           "task",
			Aliases:       []string{"azure", "ado", "az"},
			Short:         "Run an Azure DevOps Task",
			Args:          cobra.ExactArgs(1),
			Hidden:        true,
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx = cmd.Context()
					t   = &forge.Task{
						Task:      args[0],
						Inputs:    inputs,
						Execution: execution,
					}
				)

				cr, opts, err := runOptsAndContainerRuntime(cmd)
				if err != nil {
					return err
				}

				if attach {
					forge.HookContainerStarted.Listen(hookAttach(cmd, forge.AzureDevOpsTaskWorkingDir(opts.WorkingDir)))
				}

				return t.Run(ctx, cr, opts)
			},
		}
	)

	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	cmd.Flags().BoolVarP(&attach, "attach", "a", false, "Attach to containers")
	cmd.Flags().StringToStringVarP(&inputs, "input", "i", nil, "inputs")
	cmd.Flags().StringVar(&forge.NodeImageReference, "node-image", forge.DefaultNodeImageReference, "Node image for task")
	cmd.Flags().StringVar(&forge.Node10ImageReference, "node10-image", forge.DefaultNode10ImageReference, "Node10 image for task")
	cmd.Flags().StringVar(&forge.Node16ImageReference, "node16-image", forge.DefaultNode16ImageReference, "Node16 image for task")
	cmd.Flags().StringVarP(&execution, "exec", "e", "Node", "Task execution")
	cmd.Flags().StringVar(&workdir, "workdir", wd, "Working directory for task")
	_ = cmd.MarkFlagDirname("workdir")

	return cmd
}
