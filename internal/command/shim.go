package command

import (
	"runtime"

	"github.com/frantjc/forge"
	"github.com/spf13/cobra"
)

// NewShim returns the command which acts as
// the entrypoint for `shim`.
func NewShim() *cobra.Command {
	var (
		verbosity int
		cmd       = &cobra.Command{
			Use:           "shim",
			Version:       forge.SemVer(),
			SilenceErrors: true,
			SilenceUsage:  true,
			PersistentPreRun: func(cmd *cobra.Command, _ []string) {
				cmd.SetContext(
					forge.WithLogger(cmd.Context(), forge.NewLogger().V(2-verbosity)),
				)
			},
		}
	)

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.PersistentFlags().CountVarP(&verbosity, "verbose", "V", "verbosity for forge")
	cmd.AddCommand(NewSleep(), NewExec())

	return cmd
}
