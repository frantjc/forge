package command

import (
	"github.com/frantjc/forge"
	"github.com/spf13/cobra"
)

func commandDrains(cmd *cobra.Command) *forge.Drains {
	return &forge.Drains{
		Out: cmd.OutOrStdout(),
		Err: cmd.ErrOrStderr(),
	}
}
