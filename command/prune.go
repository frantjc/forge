package command

import (
	"os"

	"github.com/frantjc/forge/internal/hostfs"
	"github.com/spf13/cobra"
)

// NewPrune returns the command which acts as
// the entrypoint for `forge prune`.
func NewPrune() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "prune",
		Short:         "Prune the Forge cache",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return os.RemoveAll(hostfs.ActionsCache)
		},
	}

	return cmd
}
