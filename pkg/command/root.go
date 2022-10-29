package command

import (
	"runtime"

	"github.com/frantjc/forge"
	"github.com/spf13/cobra"
)

func NewRoot() *cobra.Command {
	var (
		verbosity int
		cmd       = &cobra.Command{
			Use:     "4ge",
			Version: forge.Semver(),
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				cmd.SetContext(forge.WithLogger(cmd.Context(), forge.NewLogger().V(verbosity)))
			},
			SilenceErrors: true,
			SilenceUsage:  true,
		}
	)

	cmd.PersistentFlags().CountVarP(&verbosity, "verbose", "v", "verbosity")
	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.AddCommand(NewUse(), NewGet(), NewPut(), NewPrune())

	return cmd
}
