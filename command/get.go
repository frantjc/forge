package command

import (
	"github.com/frantjc/forge/concourse"
	"github.com/spf13/cobra"
)

// NewGet returns the command which acts as
// the entrypoint for `forge get`.
func NewGet() *cobra.Command {
	return newResource(concourse.MethodGet, false)
}
