package command

import (
	"runtime"

	"github.com/frantjc/forge"
	"github.com/spf13/cobra"
)

// New returns the "root" command for `forge`
// which acts as forge's CLI entrypoint.
func NewForge() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "forge",
		Version:       forge.GetSemver(),
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.AddCommand(NewUse(), NewGet(), NewPut(), NewCheck(), NewPrune())

	return cmd
}
