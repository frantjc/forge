package command

import (
	"github.com/frantjc/forge/pkg/concourse"
	"github.com/spf13/cobra"
)

// NewCheck returns the command which acts as
// the entrypoint for `forge check`.
func NewCheck() *cobra.Command {
	return newResource(concourse.MethodCheck, true)
}
