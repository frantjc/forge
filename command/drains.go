package command

import (
	"github.com/frantjc/forge"
	"github.com/frantjc/go-fn"
	"github.com/spf13/cobra"
)

func commandDrains(cmd *cobra.Command, stdoutUsed ...bool) *forge.Drains {
	return &forge.Drains{
		Out: fn.Ternary(
			fn.Some(stdoutUsed, func(b bool, _ int) bool {
				return b
			}),
			cmd.ErrOrStderr(),
			cmd.OutOrStdout(),
		),
		Err: cmd.ErrOrStderr(),
	}
}
