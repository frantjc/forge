package fa

import (
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/github/actions"
)

type DiscardWorkflowCommandWriter struct {
	*actions.GlobalContext
	ID                 string
	StopCommandsTokens map[string]bool
	State              map[string]string
}

func (w *DiscardWorkflowCommandWriter) Callback(wc *actions.WorkflowCommand) []byte {
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
	case actions.CommandSetOutput:
		if w.GlobalContext.StepsContext[w.ID] == nil {
			w.GlobalContext.StepsContext[w.ID] = &actions.StepContext{
				Outputs: map[string]string{},
			}
		}

		w.GlobalContext.StepsContext[w.ID].Outputs[wc.GetName()] = wc.Value
	case actions.CommandStopCommands:
		w.StopCommandsTokens[wc.Value] = true
	case actions.CommandSaveState:
		w.State[wc.GetName()] = wc.Value
	}

	return make([]byte, 0)
}

func NewWorkflowCommandStreams(globalContext *actions.GlobalContext, id string, drains *forge.Drains) *forge.Streams {
	if globalContext == nil {
		globalContext = ConfigureGlobalContext(actions.NewGlobalContextFromEnv())
	}

	w := &DiscardWorkflowCommandWriter{
		GlobalContext:      globalContext,
		ID:                 id,
		StopCommandsTokens: map[string]bool{},
		State:              map[string]string{},
	}

	return &forge.Streams{
		Drains: &forge.Drains{
			Out: actions.NewWorkflowCommandWriter(w.Callback, drains.Out),
			Err: actions.NewWorkflowCommandWriter(w.Callback, drains.Err),
			Tty: drains.Tty,
		},
	}
}
