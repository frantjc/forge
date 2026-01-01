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
		commandWriter = githubactions.NewWorkflowCommandWriter(out, nil)
		command = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandDebug,
			Parameters: map[string]string{},
			Value:      "hello there",
		}
		_, err = commandWriter.Write([]byte(command.String()))
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
		globalContext = githubactions.NewGlobalContextFromEnv()
		commandWriter = githubactions.NewWorkflowCommandWriter(out, globalContext)
		command = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandEcho,
			Parameters: map[string]string{},
			Value:      "on",
		}
		_, err = command.WriteTo(commandWriter)
	)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if out.Len() > 0 {
		t.Error("actual", `"`+out.String()+`"`, `does not equal expected ""`)
		t.FailNow()
	}

	if !globalContext.DebugEnabled() {
		t.Error("debug is not true")
		t.FailNow()
	}
}

func TestCommandDebugOn(t *testing.T) {
	var (
		out           = new(bytes.Buffer)
		globalContext = githubactions.NewGlobalContextFromEnv().EnableDebug()
		commandWriter = githubactions.NewWorkflowCommandWriter(out, globalContext)
		value   = "hello there"
		command = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandDebug,
			Parameters: map[string]string{},
			Value:      value,
		}
		_, err = command.WriteTo(commandWriter)
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
		globalContext = githubactions.NewGlobalContextFromEnv()
		commandWriter = githubactions.NewWorkflowCommandWriter(out, globalContext)
		value   = "hello there"
		command = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandNotice,
			Parameters: map[string]string{},
			Value:      value,
		}
		_, err = command.WriteTo(commandWriter)
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
		globalContext = githubactions.NewGlobalContextFromEnv()
		commandWriter = githubactions.NewWorkflowCommandWriter(out, globalContext)
		value   = "hello there"
		command = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandWarning,
			Parameters: map[string]string{},
			Value:      value,
		}
		_, err = command.WriteTo(commandWriter)
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
		globalContext = githubactions.NewGlobalContextFromEnv().EnableDebug()
		commandWriter = githubactions.NewWorkflowCommandWriter(out, globalContext)
		value   = "hello there"
		command = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandError,
			Parameters: map[string]string{},
			Value:      value,
		}
		_, err = command.WriteTo(commandWriter)
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
		globalContext = githubactions.NewGlobalContextFromEnv()
		commandWriter = githubactions.NewWorkflowCommandWriter(out, globalContext)
		maskCommand = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandAddMask,
			Parameters: map[string]string{},
			Value:      value,
		}
	)

	if _, err := maskCommand.WriteTo(commandWriter); err != nil {
		t.Error(err)
		t.FailNow()
	}

	var (
		command = &githubactions.WorkflowCommand{
			Command:    githubactions.CommandWarning,
			Parameters: map[string]string{},
			Value:      value,
		}
		expected = "[" + githubactions.CommandWarning + "] " + "***\n"
	)
	if _, err := command.WriteTo(commandWriter); err != nil {
		t.Error(err)
		t.FailNow()
	}
	
	actual := out.String()

	if actual != expected {
		t.Error("actual", `"`+actual+`"`, "does not equal expected", `"`+expected+`"`)
		t.FailNow()
	}
}
