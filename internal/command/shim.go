package command

import (
	"runtime"

	"github.com/spf13/cobra"
)

// NewShim returns the command which acts as
// the entrypoint for `shim`.
func NewShim() *cobra.Command {
	var (
		verbosity int
		cmd       = &cobra.Command{
			Use:           "shim",
			SilenceErrors: true,
			SilenceUsage:  true,
		}
	)

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.PersistentFlags().CountVarP(&verbosity, "verbose", "V", "verbosity for shim")
	cmd.AddCommand(NewSleep(), NewExec())

	return cmd
}
