package actions_test

import (
	"testing"

	"github.com/frantjc/forge/pkg/github/actions"
)

func TestWorkflowCommandToString(t *testing.T) {
	var (
		wc = &actions.WorkflowCommand{
			Command: "set-output",
			Parameters: map[string]string{
				"name":       "var",
				"otherParam": "param",
			},
			Value: "value",
		}
		expected = "::set-output name=var,otherParam=param::value"
		actual   = wc.CommandString()
	)

	if actual != expected {
		t.Error("actual", actual, "does not match expected", expected)
		t.FailNow()
	}
}
