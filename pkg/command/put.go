package command

import (
	"github.com/frantjc/forge/pkg/concourse"
	"github.com/spf13/cobra"
)

func NewPut() *cobra.Command {
	var (
		params  = map[string]string{}
		version = map[string]string{}
		cmd     = &cobra.Command{
			Use:   concourse.MethodPut,
			Short: "Put a Concourse Resource",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return processResource(cmd.Context(), cmd.Use, args[0], params, version)
			},
		}
	)

	cmd.Flags().StringToStringVarP(&params, "param", "p", make(map[string]string), "params")
	cmd.Flags().StringToStringVarP(&version, "version", "i", make(map[string]string), "version")

	return cmd
}
