package githubactions

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	xos "github.com/frantjc/x/os"
)

// WorkflowCommandWriter holds the state of GitHub Actions
// workflow commands throughout the execution of a step.
type WorkflowCommandWriter struct {
	*GlobalContext
	ID                 string
	StopCommandsTokens map[string]bool
	Debug              bool
	Masks              []string
	Out                io.Writer

	saveStateDeprecationWarned bool
	setOutputDeprecationWarned bool
}

// handleCommand takes a *WorkflowCommand and processes it by storing
// a value in the appropriate location in its *GlobalContext if
// necessary. It returns the bytes that should be written for the workflow command.
func (w *WorkflowCommandWriter) handleCommand(wc *WorkflowCommand) []byte {
	if w.StopCommandsTokens == nil {
		w.StopCommandsTokens = make(map[string]bool)
	}

	if _, ok := w.StopCommandsTokens[wc.Command]; ok {
		w.StopCommandsTokens[wc.Command] = false
		return make([]byte, 0)
	}

	for _, stop := range w.StopCommandsTokens {
		if stop {
			return []byte(wc.String())
		}
	}

	if w.GlobalContext == nil {
		w.GlobalContext = NewGlobalContextFromEnv()
	}

	switch wc.Command {
	case CommandSetOutput:
		if _, ok := w.GlobalContext.StepsContext[w.ID]; !ok {
			w.GlobalContext.StepsContext[w.ID] = StepContext{
				Outputs: make(map[string]string),
			}
		}

		w.GlobalContext.StepsContext[w.ID].Outputs[wc.GetName()] = wc.Value

		if !w.setOutputDeprecationWarned {
			return []byte("[" + CommandWarning + "] The `" + wc.Command + "` command is deprecated and will be disabled soon. Please upgrade to using Environment Files. For more information see: https://github.blog/changelog/2022-10-11-github-actions-deprecating-save-state-and-set-output-commands/")
		}
	case CommandStopCommands:
		w.StopCommandsTokens[wc.Value] = true
	case CommandSaveState:
		w.GlobalContext.EnvContext["STATE_"+wc.GetName()] = wc.Value

		if !w.saveStateDeprecationWarned {
			return []byte("[" + CommandWarning + "] The `" + wc.Command + "` command is deprecated and will be disabled soon. Please upgrade to using Environment Files. For more information see: https://github.blog/changelog/2022-10-11-github-actions-deprecating-save-state-and-set-output-commands/")
		}
	case CommandEcho:
		switch wc.Value {
		case "on":
			w.Debug = true
		case "off":
			w.Debug = false
		default:
			// Not sure if this was ever correct.
			// Keeping it for backwards compatibility.
			w.Debug = !w.Debug
		}
	case CommandAddMask:
		w.Masks = append(w.Masks, wc.Value)
	case CommandAddPath:
		w.GlobalContext.EnvContext["PATH"] = xos.JoinPath(wc.Value, w.GlobalContext.EnvContext["PATH"])
	case CommandEndGroup:
		return []byte("[" + CommandEndGroup + "]")
	case CommandDebug:
		if w.Debug {
			return []byte("[" + CommandDebug + "] " + wc.Value)
		}
	default:
		return []byte("[" + wc.Command + "] " + wc.Value)
	}

	return make([]byte, 0)
}

func (w *WorkflowCommandWriter) IssueCommand(wc *WorkflowCommand) (int, error) {
	return fmt.Fprintln(w, wc.String())
}

func (w *WorkflowCommandWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	scanner := bufio.NewScanner(bytes.NewReader(p))

	for scanner.Scan() {
		line := scanner.Text()
		for _, mask := range w.Masks {
			line = strings.ReplaceAll(line, mask, "***")
		}

		b := []byte(line)

		//nolint:revive
		if len(line) == 0 || strings.HasPrefix(line, "##[add-matcher]") {
		} else if c, err := ParseWorkflowCommandString(line); err == nil {
			b = w.handleCommand(c)
		}

		if len(b) > 0 {
			b = append(b, '\n')

			if n, err := w.Out.Write(b); err != nil {
				return n, err
			}
		}
	}

	return len(p), nil
}
