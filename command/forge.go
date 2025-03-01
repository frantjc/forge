package command

import (
	"runtime"

	"github.com/spf13/cobra"
)

// New returns the "root" command for `forge`
// which acts as Forge's CLI entrypoint.
func NewForge() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "forge",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.PersistentFlags().CountP("verbose", "V", "Verbosity for forge")
	cmd.PersistentFlags().Bool("fix-dind", false, "Intercept and fix traffic to docker.sock,")
	cmd.PersistentFlags().Bool("no-dind", false, "Disable Docker in Docker,")
	cmd.MarkFlagsMutuallyExclusive("no-dind", "fix-dind")
	cmd.AddCommand(NewUse(), NewGet(), NewPut(), NewCheck(), NewTask(), NewCloudBuild(), NewCache())

	return cmd
}
