package command

import (
	"runtime"

	"github.com/frantjc/forge/internal/containerutil"
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
	cmd.PersistentFlags().BoolVar(&containerutil.UseForgeSock, "use-sock", false, "enable use of forge.sock")
	cmd.PersistentFlags().Bool("no-dind", false, "disable Docker in Docker")
	cmd.MarkFlagsMutuallyExclusive("no-dind", "use-sock")
	cmd.AddCommand(NewUse(), NewGet(), NewPut(), NewCheck(), NewTask(), NewCloudBuild(), NewCache())

	return cmd
}
