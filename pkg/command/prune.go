package command

import (
	"os"

	"github.com/frantjc/forge/internal/hostfs"
	"github.com/spf13/cobra"
)

func NewPrune() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Prune the Forge cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			return os.RemoveAll(hostfs.ActionsCache)
		},
	}

	return cmd
}
