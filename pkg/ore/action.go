package ore

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/containerutil"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/frantjc/forge/pkg/forgeactions"
	"github.com/frantjc/forge/pkg/githubactions"
)

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

	if o.GetGlobalContext() == nil {
		o.GlobalContext = githubactions.NewGlobalContextFromEnv()
	}
	defer func() {
		ctx = githubactions.WithGlobalContext(ctx, o.GlobalContext)
	}()

	containerConfigs, err := forgeactions.ActionToConfigs(o.GetGlobalContext(), uses, o.GetWith(), o.GetEnv(), actionMetadata, image)
	if err != nil {
		return nil, err
	}

	workflowCommandStreams := forgeactions.NewWorkflowCommandStreams(o.GetGlobalContext(), o.GetId(), drains)
	for _, containerConfig := range containerConfigs {
		containerConfig.Mounts = contaminate.OverrideWithMountsFrom(ctx, containerConfig.GetMounts()...)
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
