package actions_test

import (
	"strings"
	"testing"

	"github.com/frantjc/forge/pkg/github/actions"
)

func TestWorkflowCommandToString(t *testing.T) {
	var (
		command = &actions.WorkflowCommand{
			Command: "set-output",
			Parameters: map[string]string{
				"name":       "var",
				"otherParam": "param",
			},
			Value: "value",
		}
		expected = []string{"::set-output name=var,otherParam=param::value", "::set-output otherParam=param,name=var::value"}
		actual   = command.String()
	)

	match := false

	for _, e := range expected {
		match = match || strings.EqualFold(actual, e)
	}
	if !match {
		t.Error("actual", actual, "does not match", expected[0], "or", expected[1])
		t.FailNow()
	}
}
