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
			Run: func(cmd *cobra.Command, args []string) {
				if err := os.RemoveAll(hfs.ActionCache); err != nil {
					cmd.PrintErrln(err)
				}
			},
		}
	)

	return cmd
}
