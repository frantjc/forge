package command

import (
	"runtime"

	"github.com/frantjc/forge"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var (
		workdir   string
		verbosity int
		cmd       = &cobra.Command{
			Use:     "forge",
			Version: forge.Semver(),
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				cmd.SetContext(
					WithWorkdir(
						forge.WithLogger(cmd.Context(), forge.NewLogger().V(verbosity)),
						workdir,
					),
				)
			},
			SilenceErrors: true,
			SilenceUsage:  true,
		}
	)

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.PersistentFlags().CountVarP(&verbosity, "verbose", "v", "verbosity for forge")
	cmd.PersistentFlags().StringVarP(&workdir, "workdir", "d", "", "working directory for forge")
	_ = cmd.MarkFlagDirname("workdir")
	cmd.AddCommand(NewUse(), NewGet(), NewPut(), NewCheck(), NewPrune())

	return cmd
}
