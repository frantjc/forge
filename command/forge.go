package command

import (
	"runtime"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/containerutil"
	"github.com/spf13/cobra"
)

// New returns the "root" command for `forge`
// which acts as forge's CLI entrypoint.
func NewForge() *cobra.Command {
	var (
		verbosity int
		cmd       = &cobra.Command{
			Use:           "forge",
			Version:       forge.SemVer(),
			SilenceErrors: true,
			SilenceUsage:  true,
		}
	)

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.PersistentFlags().CountVarP(&verbosity, "verbose", "V", "verbosity for forge")
	cmd.PersistentFlags().BoolVar(&containerutil.NoUseForgeSock, "no-sock", false, "disable use of forge.sock")
	cmd.AddCommand(NewUse(), NewGet(), NewPut(), NewCheck(), NewTask(), NewCloudBuild(), NewCache())

	return cmd
}
