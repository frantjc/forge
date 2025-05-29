package githubactions

import (
	"fmt"

	"github.com/frantjc/forge/internal/rangemap"
)

type WorkflowCommand struct {
	Command    string
	Parameters map[string]string
	Value      string
}

func (c *WorkflowCommand) String() string {
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
	return "&WorkflowCommand{" + c.String() + "}"
}

// GetName returns the value of the name parameter from the workflow command.
// Useful for set-env workflow commands as they require it to specify the name
// of the environment variable.
func (c *WorkflowCommand) GetName() string {
	if c.Parameters != nil {
		if name, ok := c.Parameters["name"]; ok {
			return name
		}
	}

	return ""
}
