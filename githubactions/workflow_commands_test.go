package githubactions_test

import (
	"testing"

	"github.com/frantjc/forge/githubactions"
)

func TestParseCommandNoParams(t *testing.T) {
	var (
		command  = "::" + githubactions.CommandDebug + "::hello there"
		expected = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandDebug,
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
		t.Error("actual", `"`+actual.String()+`"`, "does not equal expected", `"`+expected.String()+`"`)
		t.FailNow()
	}
}

func TestParseCommandOneParam(t *testing.T) {
	var (
		command  = "::" + githubactions.CommandSaveState + " name=isPost::true"
		expected = &githubactions.WorkflowCommand{
			Command: githubactions.CommandSaveState,
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
		t.Error("actual", `"`+actual.String()+`"`, "does not equal expected", `"`+expected.String()+`"`)
		t.FailNow()
	}
}

func TestParseCommandManyParams(t *testing.T) {
	var (
		command  = "::" + githubactions.CommandSaveState + " name=isPost,otherParam=1::true"
		expected = &githubactions.WorkflowCommand{
			Command: githubactions.CommandSaveState,
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
		t.Error("actual", `"`+actual.String()+`"`, "does not equal expected", `"`+expected.String()+`"`)
		t.FailNow()
	}
}

func TestParseCommandNoValue(t *testing.T) {
	var (
		command  = "::" + githubactions.CommandSaveState + " name=isPost::"
		expected = &githubactions.WorkflowCommand{
			Command: githubactions.CommandSaveState,
			Parameters: map[string]string{
				"name": "isPost",
			},
		}
		actual, err = githubactions.ParseWorkflowCommandString(command)
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if expected.String() != actual.String() {
		t.Error("actual", `"`+actual.String()+`"`, "does not equal expected", `"`+expected.String()+`"`)
		t.FailNow()
	}
}
