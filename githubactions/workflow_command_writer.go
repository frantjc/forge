package githubactions

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	xos "github.com/frantjc/x/os"
)

func NewWorkflowCommandWriter(w io.Writer, globalContext *GlobalContext) io.Writer {
	if globalContext == nil {
		globalContext = NewGlobalContextFromEnv()
	}

	return &workflowCommandWriter{
		globalContext: *globalContext,
		debug:         globalContext.DebugEnabled(),
		w:             w,
	}
}

// workflowCommandWriter holds the state of GitHub Actions
// workflow commands throughout the execution of a step.
type workflowCommandWriter struct {
	globalContext              GlobalContext
	id                         string
	stopCommandsTokens         map[string]bool
	debug                      bool
	masks                      []string
	w                          io.Writer
	saveStateDeprecationWarned bool
	setOutputDeprecationWarned bool
}

// handleCommand takes a *WorkflowCommand and processes it by storing
// a value in the appropriate location in its *GlobalContext if
// necessary. It returns the bytes that should be written for the workflow command.
func (w *workflowCommandWriter) handleCommand(wc *WorkflowCommand) []byte {
	if w.stopCommandsTokens == nil {
		w.stopCommandsTokens = make(map[string]bool)
	}

	if _, ok := w.stopCommandsTokens[wc.Command]; ok {
		w.stopCommandsTokens[wc.Command] = false
		return make([]byte, 0)
	}

	for _, stop := range w.stopCommandsTokens {
		if stop {
			return []byte(wc.String())
		}
	}

	switch wc.Command {
	case CommandSetOutput:
		if _, ok := w.globalContext.StepsContext[w.id]; !ok {
			w.globalContext.StepsContext[w.id] = StepContext{
				Outputs: make(map[string]string),
			}
		}

		w.globalContext.StepsContext[w.id].Outputs[wc.GetName()] = wc.Value

		if output, err := os.OpenFile(os.Getenv(EnvVarOutput), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666); err == nil {
			defer output.Close()

			_, _ = fmt.Fprintf(output, "%s=%s\n", wc.GetName(), wc.Value)
		}

		if !w.setOutputDeprecationWarned {
			return []byte("[" + CommandWarning + "] The `" + wc.Command + "` command is deprecated and will be disabled soon. Please upgrade to using Environment Files. For more information see: https://github.blog/changelog/2022-10-11-github-actions-deprecating-save-state-and-set-output-commands/")
		}
	case CommandStopCommands:
		w.stopCommandsTokens[wc.Value] = true
	case CommandSaveState:
		w.globalContext.EnvContext[fmt.Sprintf("STATE_%s", wc.GetName())] = wc.Value

		if state, err := os.OpenFile(os.Getenv(EnvVarState), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666); err == nil {
			defer state.Close()

			_, _ = fmt.Fprintf(state, "%s=%s\n", wc.GetName(), wc.Value)
		}

		if !w.saveStateDeprecationWarned {
			return []byte("[" + CommandWarning + "] The `" + wc.Command + "` command is deprecated and will be disabled soon. Please upgrade to using Environment Files. For more information see: https://github.blog/changelog/2022-10-11-github-actions-deprecating-save-state-and-set-output-commands/")
		}
	case CommandEcho:
		switch wc.Value {
		case "on":
			w.debug = true
		case "off":
			w.debug = false
		}
	case CommandAddMask:
		w.masks = append(w.masks, wc.Value)
	case CommandAddPath:
		w.globalContext.EnvContext["PATH"] = xos.JoinPath(wc.Value, w.globalContext.EnvContext["PATH"])

		if path, err := os.OpenFile(os.Getenv(EnvVarPath), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666); err == nil {
			defer path.Close()

			_, _ = fmt.Fprintln(path, wc.Value)
		}
	case CommandEndGroup:
		return []byte("[" + CommandEndGroup + "]")
	case CommandDebug:
		if w.debug {
			return []byte("[" + CommandDebug + "] " + wc.Value)
		}
	default:
		return []byte("[" + wc.Command + "] " + wc.Value)
	}

	return make([]byte, 0)
}

func (w *workflowCommandWriter) IssueCommand(wc *WorkflowCommand) error {
	if _, err := fmt.Fprintln(w, wc.String()); err != nil {
		return err
	}

	return nil
}

func (w *workflowCommandWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	scanner := bufio.NewScanner(bytes.NewReader(p))

	for scanner.Scan() {
		line := scanner.Text()
		for _, mask := range w.masks {
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

			if n, err := w.w.Write(b); err != nil {
				return n, err
			}
		}
	}

	return len(p), nil
}
