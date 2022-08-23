package actions2container

import (
	"io"

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
			w.GlobalContext.StepsContext[w.ID] = &actions.StepsContext{
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

func NewWorkflowCommandStreams(globalContext *actions.GlobalContext, id string, stdout, stderr io.Writer) *forge.Streams {
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
		Out: actions.NewWorkflowCommandWriter(w.Callback, stdout),
		Err: actions.NewWorkflowCommandWriter(w.Callback, stderr),
	}
}
