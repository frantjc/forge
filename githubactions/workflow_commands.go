package githubactions

import (
	"fmt"
	"strings"
)

const (
	CommandDebug        = "debug"
	CommandGroup        = "group"
	CommandEndGroup     = "endgroup"
	CommandSaveState    = "save-state"
	CommandSetOutput    = "set-output"
	CommandNotice       = "notice"
	CommandWarning      = "warning"
	CommandError        = "error"
	CommandAddMask      = "add-mask"
	CommandAddPath      = "add-path"
	CommandEcho         = "echo"
	CommandStopCommands = "stop-commands"
)

// ParseWorkflowCommandString parses a workflow command from a string such as:
//
// ::set-env name=HELLO::there
//
// This supports a deprecated GitHub Actions function that used such strings written
// to stdout to send commands from a GitHub Action up to GitHub Actions.
func ParseWorkflowCommandString(workflowCommand string) (*WorkflowCommand, error) {
	if !strings.HasPrefix(workflowCommand, "::") {
		return nil, fmt.Errorf("not a workflow command: %s", workflowCommand)
	}

	a := strings.Split(workflowCommand, "::")
	if len(a) < 2 {
		return nil, fmt.Errorf("not a workflow command: %s", workflowCommand)
	}

	cmdAndParams := a[1]
	b := strings.Split(cmdAndParams, " ")
	if len(b) < 1 {
		return nil, fmt.Errorf("not a workflow command: %s", workflowCommand)
	}

	cmd := b[0]
	params := map[string]string{}

	if len(b) > 1 {
		for _, p := range strings.Split(b[1], ",") {
			if f := strings.Split(p, "="); len(f) > 1 {
				params[f[0]] = f[1]
			}
		}
	}

	value := ""
	if len(a) > 2 {
		value = a[2]
	}

	return &WorkflowCommand{
		Command:    cmd,
		Parameters: params,
		Value:      value,
	}, nil
}

// ParseWorkflowCommand parses a workflow command from bytes.
// See ParseWorkflowCommandString for more details.
func ParseWorkflowCommand(b []byte) (*WorkflowCommand, error) {
	return ParseWorkflowCommandString(string(b))
}
