package forgeactions

import (
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/github/actions"
)

type WorkflowCommandWriter struct {
	*actions.GlobalContext
	ID                 string
	StopCommandsTokens map[string]bool
	State              map[string]string
	Debug              bool
}

func (w *WorkflowCommandWriter) Callback(wc *actions.WorkflowCommand) []byte {
	if _, ok := w.StopCommandsTokens[wc.GetCommand()]; ok {
		w.StopCommandsTokens[wc.GetCommand()] = false
		return make([]byte, 0)
	}

	for _, stop := range w.StopCommandsTokens {
		if stop {
			return []byte(wc.CommandString())
		}
	}

	switch wc.Command {
	case actions.CommandSetOutput:
		if w.GlobalContext.StepsContext[w.ID] == nil {
			w.GlobalContext.StepsContext[w.ID] = &actions.StepContext{
				Outputs: map[string]string{},
			}
		}

		w.GlobalContext.StepsContext[w.ID].Outputs[wc.GetName()] = wc.GetValue()
	case actions.CommandStopCommands:
		w.StopCommandsTokens[wc.GetValue()] = true
	case actions.CommandSaveState:
		w.State[wc.GetName()] = wc.GetValue()
	case actions.CommandEcho:
		w.Debug = !w.Debug
	case actions.CommandEndGroup:
		return []byte("[endgroup]")
	case actions.CommandDebug:
		if w.Debug {
			return []byte("[debug] " + wc.GetValue())
		}
	default:
		return []byte("[" + wc.GetCommand() + "] " + wc.GetValue())
	}

	return make([]byte, 0)
}

func NewWorkflowCommandStreams(globalContext *actions.GlobalContext, id string, drains *forge.Drains) *forge.Streams {
	if globalContext == nil {
		globalContext = ConfigureGlobalContext(actions.NewGlobalContextFromEnv())
	}

	w := &WorkflowCommandWriter{
		GlobalContext:      globalContext,
		ID:                 id,
		StopCommandsTokens: map[string]bool{},
		State:              map[string]string{},
		Debug:              globalContext.SecretsContext[actions.SecretActionsStepDebug] == actions.SecretDebugValue,
	}

	return &forge.Streams{
		Drains: &forge.Drains{
			Out: actions.NewWorkflowCommandWriter(w.Callback, drains.Out),
			Err: actions.NewWorkflowCommandWriter(w.Callback, drains.Err),
			Tty: drains.Tty,
		},
	}
}
