package command

import (
	"github.com/frantjc/forge/pkg/concourse"
	"github.com/spf13/cobra"
)

func NewGet() *cobra.Command {
	return newResource(concourse.MethodGet, false)
}
