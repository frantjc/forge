package command

import (
	"github.com/frantjc/forge/concourse"
	"github.com/spf13/cobra"
)

// NewPut returns the command which acts as
// the entrypoint for `forge put`.
func NewPut() *cobra.Command {
	return newResource(concourse.MethodPut, false)
}
