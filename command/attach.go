package command

import (
	"github.com/frantjc/forge"
	"github.com/spf13/cobra"
)

func commandStreams(cmd *cobra.Command) *forge.Streams {
	return commandDrains(cmd).ToStreams(cmd.InOrStdin())
}
