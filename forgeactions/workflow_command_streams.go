package forgeactions

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/githubactions"
	"github.com/frantjc/forge/internal/contaminate"
)

// NewWorkflowCommandStreams takes io.Writers and returns wrapped writers to pass to a process executing
// a GitHub Action as stdout and stderr. These streams process workflow commands that are written to them
// and write any corresponding bytes to the underlying writers. They write any non-workflow command bytes
// directly to the underlying writers.
func NewWorkflowCommandStreams(ctx context.Context, globalContext *githubactions.GlobalContext, id string, drains *forge.Drains) *forge.Streams {
	globalContext = ConfigureGlobalContext(globalContext)
	debug := globalContext.DebugEnabled()

	return &forge.Streams{
		In: contaminate.StdinFrom(ctx),
		Drains: &forge.Drains{
			Out: &githubactions.WorkflowCommandWriter{
				GlobalContext:      globalContext,
				ID:                 id,
				StopCommandsTokens: map[string]bool{},
				Debug:              debug,
				Out:                drains.Out,
			},
			Err: &githubactions.WorkflowCommandWriter{
				GlobalContext:      globalContext,
				ID:                 id,
				StopCommandsTokens: map[string]bool{},
				Debug:              debug,
				Out:                drains.Err,
			},
			Tty: drains.Tty,
		},
	}
}
