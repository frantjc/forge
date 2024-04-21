package command

import (
	"io"

	"github.com/frantjc/forge"
	xslice "github.com/frantjc/x/slice"
	"github.com/spf13/cobra"
)

func commandDrains(cmd *cobra.Command, stdoutUsed ...bool) *forge.Drains {
	return &forge.Drains{
		Out: func() io.Writer {
			if xslice.Some(stdoutUsed, func(b bool, _ int) bool {
				return b
			}) {
				return cmd.ErrOrStderr()
			}

			return cmd.OutOrStdout()
		}(),
		Err: cmd.ErrOrStderr(),
	}
}
