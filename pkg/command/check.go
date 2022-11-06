package command

import (
	"github.com/frantjc/forge/pkg/concourse"
	"github.com/spf13/cobra"
)

func NewCheck() *cobra.Command {
	return newResource(concourse.MethodCheck)
}
