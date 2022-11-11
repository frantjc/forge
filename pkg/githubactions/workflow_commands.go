package githubactions

import "strings"

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
	CommandEcho         = "echo"
	CommandStopCommands = "stop-commands"
)

// TODO regexp.
func ParseWorkflowCommandString(workflowCommand string) (*WorkflowCommand, error) {
	if !strings.HasPrefix(workflowCommand, "::") {
		return nil, ErrNotAWorkflowCommand
	}

	a := strings.Split(workflowCommand, "::")
	if len(a) < 2 {
		return nil, ErrNotAWorkflowCommand
	}

	cmdAndParams := a[1]
	b := strings.Split(cmdAndParams, " ")
	if len(b) < 1 {
		return nil, ErrNotAWorkflowCommand
	}

	cmd := b[0]
	params := map[string]string{}

	if len(b) > 1 {
		for _, p := range strings.Split(b[1], ",") {
			if f := strings.Split(p, "="); len(f) > 0 {
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

func ParseWorkflowCommand(b []byte) (*WorkflowCommand, error) {
	return ParseWorkflowCommandString(string(b))
}
