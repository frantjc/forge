package githubactions

import (
	"fmt"

	"github.com/frantjc/forge/pkg/rangemap"
)

func (c *WorkflowCommand) CommandString() string {
	s := fmt.Sprintf("::%s", c.Command)

	paramSpl := " "
	numParams := len(c.Parameters)
	paramsAdded := 0
	rangemap.Ascending(c.Parameters, func(k, v string) {
		s = fmt.Sprintf("%s%s%s=%s", s, paramSpl, k, v)
		paramSpl = ","
		paramsAdded++
		if paramsAdded == numParams {
			paramSpl = ""
		}
	})

	return fmt.Sprintf("%s::%s", s, c.Value)
}

func (c *WorkflowCommand) GoString() string {
	return "&WorkflowCommand{" + c.CommandString() + "}"
}

func (c *WorkflowCommand) GetName() string {
	if c.Parameters != nil {
		if name, ok := c.Parameters["name"]; ok {
			return name
		}
	}

	return ""
}
