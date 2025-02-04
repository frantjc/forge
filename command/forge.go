package command

import (
	"runtime"

	"github.com/spf13/cobra"
)

// New returns the "root" command for `forge`
// which acts as Forge's CLI entrypoint.
func NewForge() *cobra.Command {
	var (
		verbosity int
		cmd       = &cobra.Command{
			Use:           "forge",
			SilenceErrors: true,
			SilenceUsage:  true,
		}
	)

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.PersistentFlags().CountVarP(&verbosity, "verbose", "V", "verbosity for forge")
	cmd.PersistentFlags().Bool("fix-dind", false, "intercept and fix traffic to docker.sock")
	cmd.PersistentFlags().Bool("no-dind", false, "disable Docker in Docker")
	cmd.MarkFlagsMutuallyExclusive("no-dind", "intercept-sock")
	cmd.AddCommand(NewUse(), NewGet(), NewPut(), NewCheck(), NewTask(), NewCloudBuild(), NewCache())

	return cmd
}
