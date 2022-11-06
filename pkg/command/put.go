package command

import (
	"github.com/frantjc/forge/pkg/concourse"
	"github.com/spf13/cobra"
)

func NewPut() *cobra.Command {
	return newNonCheckResource(concourse.MethodPut)
}
