package ore

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/containerutil"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/frantjc/forge/pkg/forgeactions"
	"github.com/frantjc/forge/pkg/githubactions"
)

type Action struct {
	ID            string                       `json:"id,omitempty"`
	Uses          string                       `json:"uses,omitempty"`
	With          map[string]string            `json:"with,omitempty"`
	Env           map[string]string            `json:"env,omitempty"`
	GlobalContext *githubactions.GlobalContext `json:"global_context,omitempty"`
}

func (o *Action) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) (*forge.Metal, error) {
	var (
		_        = forge.LoggerFrom(ctx)
		exitCode = -1
	)

	uses, err := githubactions.Parse(o.Uses)
	if err != nil {
		return nil, err
	}

	actionMetadata, err := forgeactions.GetUsesMetadata(ctx, uses)
	if err != nil {
		return nil, err
	}

	image, err := forgeactions.GetImageForMetadata(ctx, containerRuntime, actionMetadata, uses)
	if err != nil {
		return nil, err
	}

	if o.GlobalContext == nil {
		o.GlobalContext = githubactions.NewGlobalContextFromEnv()
	}
	defer func() {
		ctx = githubactions.WithGlobalContext(ctx, o.GlobalContext)
	}()

	containerConfigs, err := forgeactions.ActionToConfigs(o.GlobalContext, uses, o.With, o.Env, actionMetadata, image)
	if err != nil {
		return nil, err
	}

	workflowCommandStreams := forgeactions.NewWorkflowCommandStreams(o.GlobalContext, o.ID, drains)
	for _, containerConfig := range containerConfigs {
		containerConfig.Mounts = contaminate.OverrideWithMountsFrom(ctx, containerConfig.Mounts...)
		container, err := containerutil.CreateSleepingContainer(ctx, containerRuntime, image, containerConfig)
		if err != nil {
			return nil, err
		}
		defer container.Stop(ctx)   //nolint:errcheck
		defer container.Remove(ctx) //nolint:errcheck

		exitCode, err = container.Exec(ctx, containerConfig, workflowCommandStreams)
		if err != nil {
			return nil, err
		}
	}

	return &forge.Metal{
		ExitCode: int64(exitCode),
	}, nil
}
