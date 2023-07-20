package forgeactions

import (
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/githubactions"
)

// WorkflowCommandWriter holds the state of GitHub Actions
// workflow commands throughout the execution of a step.
type WorkflowCommandWriter struct {
	*githubactions.GlobalContext
	ID                         string
	StopCommandsTokens         map[string]bool
	Debug                      bool
	SaveStateDeprecationWarned bool
	SetOutputDeprecationWarned bool
}

// Callback takes a *githubactions.WorkflowCommand and processes it by storing
// a value in the appropriate location in its *githubactions.GlobalContext if
// necessary. It returns the bytes that should be written for the workflow command.
func (w *WorkflowCommandWriter) Callback(wc *githubactions.WorkflowCommand) []byte {
	if _, ok := w.StopCommandsTokens[wc.Command]; ok {
		w.StopCommandsTokens[wc.Command] = false
		return make([]byte, 0)
	}

	for _, stop := range w.StopCommandsTokens {
		if stop {
			return []byte(wc.String())
		}
	}

	switch wc.Command {
	case githubactions.CommandSetOutput:
		if _, ok := w.GlobalContext.StepsContext[w.ID]; !ok {
			w.GlobalContext.StepsContext[w.ID] = githubactions.StepContext{
				Outputs: make(map[string]string),
			}
		}

		w.GlobalContext.StepsContext[w.ID].Outputs[wc.GetName()] = wc.Value

		if !w.SetOutputDeprecationWarned {
			return []byte("[" + githubactions.CommandWarning + "] The `" + wc.Command + "` command is deprecated and will be disabled soon. Please upgrade to using Environment Files. For more information see: https://github.blog/changelog/2022-10-11-github-actions-deprecating-save-state-and-set-output-commands/")
		}
	case githubactions.CommandStopCommands:
		w.StopCommandsTokens[wc.Value] = true
	case githubactions.CommandSaveState:
		w.GlobalContext.EnvContext["STATE_"+wc.GetName()] = wc.Value

		if !w.SaveStateDeprecationWarned {
			return []byte("[" + githubactions.CommandWarning + "] The `" + wc.Command + "` command is deprecated and will be disabled soon. Please upgrade to using Environment Files. For more information see: https://github.blog/changelog/2022-10-11-github-actions-deprecating-save-state-and-set-output-commands/")
		}
	case githubactions.CommandEcho:
		w.Debug = !w.Debug
	case githubactions.CommandEndGroup:
		return []byte("[" + githubactions.CommandEndGroup + "]")
	case githubactions.CommandDebug:
		if w.Debug {
			return []byte("[" + githubactions.CommandDebug + "] " + wc.Value)
		}
	default:
		return []byte("[" + wc.Command + "] " + wc.Value)
	}

	return make([]byte, 0)
}

// NewWorkflowCommandStreams takes io.Writers and returns wrapped writers to pass to a process executing
// a GitHub Action as stdout and stderr. These streams process workflow commands that are written to them
// and write any corresponding bytes to the underlying writers. They write any non-workflow command bytes
// directly to the underlying writers.
func NewWorkflowCommandStreams(globalContext *githubactions.GlobalContext, id string, drains *forge.Drains) *forge.Streams {
	w := &WorkflowCommandWriter{
		GlobalContext:      ConfigureGlobalContext(globalContext),
		ID:                 id,
		StopCommandsTokens: map[string]bool{},
		Debug:              globalContext.SecretsContext[githubactions.SecretActionsStepDebug] == githubactions.SecretDebugValue,
	}

	return &forge.Streams{
		Drains: &forge.Drains{
			Out: githubactions.NewWorkflowCommandWriter(w.Callback, drains.Out),
			Err: githubactions.NewWorkflowCommandWriter(w.Callback, drains.Err),
			Tty: drains.Tty,
		},
	}
}
