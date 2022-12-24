package githubactions_test

import (
	"testing"

	"github.com/frantjc/forge/githubactions"
)

func TestParseCommandNoParams(t *testing.T) {
	var (
		command  = "::debug::hello there"
		expected = &githubactions.WorkflowCommand{
			Command:    "debug",
			Parameters: map[string]string{},
			Value:      "hello there",
		}
		actual, err = githubactions.ParseWorkflowCommandString(command)
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if expected.String() != actual.String() {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestParseCommandOneParam(t *testing.T) {
	var (
		command  = "::save-state name=isPost::true"
		expected = &githubactions.WorkflowCommand{
			Command: "save-state",
			Parameters: map[string]string{
				"name": "isPost",
			},
			Value: "true",
		}
		actual, err = githubactions.ParseWorkflowCommandString(command)
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if expected.String() != actual.String() {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestParseCommandManyParams(t *testing.T) {
	var (
		command  = "::save-state name=isPost,otherParam=1::true"
		expected = &githubactions.WorkflowCommand{
			Command: "save-state",
			Parameters: map[string]string{
				"name":       "isPost",
				"otherParam": "1",
			},
			Value: "true",
		}
		actual, err = githubactions.ParseWorkflowCommandString(command)
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if expected.String() != actual.String() {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestParseCommandNoValue(t *testing.T) {
	var (
		command  = "::save-state name=isPost::"
		expected = &githubactions.WorkflowCommand{
			Command: "save-state",
			Parameters: map[string]string{
				"name": "isPost",
			},
			Value: "",
		}
		actual, err = githubactions.ParseWorkflowCommandString(command)
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if expected.String() != actual.String() {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}
