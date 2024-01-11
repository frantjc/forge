package githubactions_test

import (
	"bytes"
	"testing"

	"github.com/frantjc/forge/githubactions"
	"github.com/google/uuid"
)

func TestCommandDebugOff(t *testing.T) {
	var (
		out           = new(bytes.Buffer)
		commandWriter = &githubactions.WorkflowCommandWriter{
			Out: out,
		}
		command = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandDebug,
			Parameters: map[string]string{},
			Value:      "hello there",
		}
		_, err = commandWriter.IssueCommand(command)
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if out.Len() > 0 {
		t.Error("actual", `"`+out.String()+`"`, `does not equal expected ""`)
		t.FailNow()
	}
}

func TestCommandEcho(t *testing.T) {
	var (
		out           = new(bytes.Buffer)
		commandWriter = &githubactions.WorkflowCommandWriter{
			Out: out,
		}
		command = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandEcho,
			Parameters: map[string]string{},
			Value:      "on",
		}
		_, err = commandWriter.IssueCommand(command)
	)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if out.Len() > 0 {
		t.Error("actual", `"`+out.String()+`"`, `does not equal expected ""`)
		t.FailNow()
	}

	if !commandWriter.Debug {
		t.Error("debug is not true")
		t.FailNow()
	}
}

func TestCommandDebugOn(t *testing.T) {
	var (
		out           = new(bytes.Buffer)
		commandWriter = &githubactions.WorkflowCommandWriter{
			Debug: true,
			Out:   out,
		}
		value   = "hello there"
		command = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandDebug,
			Parameters: map[string]string{},
			Value:      value,
		}
		_, err   = commandWriter.IssueCommand(command)
		expected = "[" + githubactions.CommandDebug + "] " + value + "\n"
		actual   = out.String()
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if actual != expected {
		t.Error("actual", `"`+actual+`"`, "does not equal expected", `"`+expected+`"`)
		t.FailNow()
	}
}

func TestCommandNotice(t *testing.T) {
	var (
		out           = new(bytes.Buffer)
		commandWriter = &githubactions.WorkflowCommandWriter{
			Out: out,
		}
		value   = "hello there"
		command = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandNotice,
			Parameters: map[string]string{},
			Value:      value,
		}
		_, err   = commandWriter.IssueCommand(command)
		expected = "[" + githubactions.CommandNotice + "] " + value + "\n"
		actual   = out.String()
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if actual != expected {
		t.Error("actual", `"`+actual+`"`, "does not equal expected", `"`+expected+`"`)
		t.FailNow()
	}
}

func TestCommandWarning(t *testing.T) {
	var (
		out           = new(bytes.Buffer)
		commandWriter = &githubactions.WorkflowCommandWriter{
			Debug: true,
			Out:   out,
		}
		value   = "hello there"
		command = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandWarning,
			Parameters: map[string]string{},
			Value:      value,
		}
		_, err   = commandWriter.IssueCommand(command)
		expected = "[" + githubactions.CommandWarning + "] " + value + "\n"
		actual   = out.String()
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if actual != expected {
		t.Error("actual", `"`+actual+`"`, "does not equal expected", `"`+expected+`"`)
		t.FailNow()
	}
}

func TestCommandError(t *testing.T) {
	var (
		out           = new(bytes.Buffer)
		commandWriter = &githubactions.WorkflowCommandWriter{
			Debug: true,
			Out:   out,
		}
		value   = "hello there"
		command = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandError,
			Parameters: map[string]string{},
			Value:      value,
		}
		_, err   = commandWriter.IssueCommand(command)
		expected = "[" + githubactions.CommandError + "] " + value + "\n"
		actual   = out.String()
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if actual != expected {
		t.Error("actual", `"`+actual+`"`, "does not equal expected", `"`+expected+`"`)
		t.FailNow()
	}
}

func TestMask(t *testing.T) {
	var (
		out           = new(bytes.Buffer)
		value         = uuid.NewString()
		commandWriter = &githubactions.WorkflowCommandWriter{
			Out:   out,
			Masks: []string{value},
		}
		command = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandWarning,
			Parameters: map[string]string{},
			Value:      value,
		}
		_, err   = commandWriter.IssueCommand(command)
		expected = "[" + githubactions.CommandWarning + "] " + "***\n"
		actual   = out.String()
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if actual != expected {
		t.Error("actual", `"`+actual+`"`, "does not equal expected", `"`+expected+`"`)
		t.FailNow()
	}
}

func TestAddMask(t *testing.T) {
	var (
		out           = new(bytes.Buffer)
		value         = uuid.NewString()
		commandWriter = &githubactions.WorkflowCommandWriter{
			Out: out,
		}
		command = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandAddMask,
			Parameters: map[string]string{},
			Value:      value,
		}
		_, err = commandWriter.IssueCommand(command)
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if out.Len() > 0 {
		t.Error("actual", `"`+out.String()+`"`, `does not equal expected ""`)
		t.FailNow()
	}

	if len(commandWriter.Masks) != 1 || commandWriter.Masks[0] != value {
		t.Error("mask was not added")
		t.FailNow()
	}
}
