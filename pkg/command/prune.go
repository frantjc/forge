package command

import (
	"os"

	hfs "github.com/frantjc/forge/internal/hostfs"
	"github.com/spf13/cobra"
)

func NewPrune() *cobra.Command {
	var (
		cmd = &cobra.Command{
			Use: "prune",
			RunE: func(cmd *cobra.Command, args []string) error {
				return os.RemoveAll(hfs.ActionsCache)
			},
		}
	)

	return cmd
}
