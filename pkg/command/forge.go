package command

import (
	"runtime"

	"github.com/frantjc/forge"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var (
		verbosity int
		cmd       = &cobra.Command{
			Use:     "forge",
			Version: forge.GetSemver(),
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				cmd.SetContext(
					forge.WithLogger(cmd.Context(), forge.NewLogger().V(verbosity)),
				)
			},
			SilenceErrors: true,
			SilenceUsage:  true,
		}
	)

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.PersistentFlags().CountVarP(&verbosity, "verbose", "v", "verbosity for forge")
	cmd.AddCommand(NewUse(), NewGet(), NewPut(), NewCheck(), NewPrune())

	return cmd
}
