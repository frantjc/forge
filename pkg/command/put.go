package command

import "github.com/spf13/cobra"

func NewPut() *cobra.Command {
	var (
		params = map[string]string{}
		cmd    = &cobra.Command{
			Use:  "put",
			Args: cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				if err := processResource(cmd.Context(), cmd.Use, args[0], params); err != nil {
					cmd.PrintErrln(err)
					return
				}
			},
		}
	)

	cmd.Flags().StringToStringVarP(&params, "param", "p", make(map[string]string), "params")

	return cmd
}
