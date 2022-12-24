package forgeactions

import (
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/githubactions"
)

type WorkflowCommandWriter struct {
	*githubactions.GlobalContext
	ID                 string
	StopCommandsTokens map[string]bool
	State              map[string]string
	Debug              bool
}

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
		if w.GlobalContext.StepsContext[w.ID] == nil {
			w.GlobalContext.StepsContext[w.ID] = &githubactions.StepContext{
				Outputs: map[string]string{},
			}
		}

		w.GlobalContext.StepsContext[w.ID].Outputs[wc.GetName()] = wc.Value
	case githubactions.CommandStopCommands:
		w.StopCommandsTokens[wc.Value] = true
	case githubactions.CommandSaveState:
		w.State[wc.GetName()] = wc.Value
	case githubactions.CommandEcho:
		w.Debug = !w.Debug
	case githubactions.CommandEndGroup:
		return []byte("[endgroup]")
	case githubactions.CommandDebug:
		if w.Debug {
			return []byte("[debug] " + wc.Value)
		}
	default:
		return []byte("[" + wc.Command + "] " + wc.Value)
	}

	return make([]byte, 0)
}

func NewWorkflowCommandStreams(globalContext *githubactions.GlobalContext, id string, drains *forge.Drains) *forge.Streams {
	if globalContext == nil {
		globalContext = ConfigureGlobalContext(githubactions.NewGlobalContextFromEnv())
	}

	w := &WorkflowCommandWriter{
		GlobalContext:      globalContext,
		ID:                 id,
		StopCommandsTokens: map[string]bool{},
		State:              map[string]string{},
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
