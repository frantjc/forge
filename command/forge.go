package command

import (
	"runtime"

	"github.com/frantjc/forge"
	"github.com/spf13/cobra"
)

// New returns the "root" command for `forge`
// which acts as forge's CLI entrypoint.
func NewForge() *cobra.Command {
	var (
		verbosity int
		cmd       = &cobra.Command{
			Use:           "forge",
			Version:       forge.GetSemver(),
			SilenceErrors: true,
			SilenceUsage:  true,
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				cmd.SetContext(
					forge.WithLogger(cmd.Context(), forge.NewLogger().V(2-verbosity)),
				)
			},
		}
	)

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.PersistentFlags().CountVarP(&verbosity, "verbose", "V", "verbosity for forge")
	cmd.AddCommand(NewUse(), NewGet(), NewPut(), NewCheck(), NewCache())

	return cmd
}
