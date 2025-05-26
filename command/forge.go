package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newVersion returns the command which acts as
// the entrypoint for `forge version`.
func newVersion(version string) *cobra.Command {
	return setCommon(&cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), version)
			return err
		},
	})
}

// New returns the "root" command for `forge`
// which acts as Forge's CLI entrypoint.
func NewForge(version string) *cobra.Command {
	cmd := setCommon(&cobra.Command{Use: "forge"})

	cmd.AddCommand(NewUse(), NewGet(), NewPut(), NewCheck(), NewTask(), NewCloudBuild(), NewCache())

	if version != "" {
		cmd.AddCommand(newVersion(version))
	}

	return cmd
}
