package ore

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/contaminate"
	fa "github.com/frantjc/forge/pkg/forgeactions"
	"github.com/frantjc/forge/pkg/github/actions"
)

func (o *Action) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) (*forge.Cast, error) {
	var (
		_        = forge.LoggerFrom(ctx)
		exitCode = -1
	)

	uses, err := actions.Parse(o.Uses)
	if err != nil {
		return nil, err
	}

	actionMetadata, err := fa.GetUsesMetadata(ctx, uses)
	if err != nil {
		return nil, err
	}

	image, err := fa.PullImageForMetadata(ctx, containerRuntime, actionMetadata)
	if err != nil {
		return nil, err
	}

	if o.GetGlobalContext() == nil {
		o.GlobalContext = actions.NewGlobalContextFromEnv()
	}
	defer func() {
		ctx = actions.WithGlobalContext(ctx, o.GlobalContext)
	}()

	conatinerConfigs, err := fa.ActionToConfigs(o.GetGlobalContext(), uses, o.GetWith(), o.GetEnv(), actionMetadata)
	if err != nil {
		return nil, err
	}

	workflowCommandStreams := fa.NewWorkflowCommandStreams(o.GetGlobalContext(), o.GetId(), drains)
	for _, containerConfig := range conatinerConfigs {
		containerConfig.Mounts = contaminate.OverrideWithMountsFrom(ctx, containerConfig.GetMounts()...)
		container, err := CreateSleepingContainer(ctx, containerRuntime, image, containerConfig)
		if err != nil {
			break
		}
		defer container.Stop(ctx)   //nolint:errcheck
		defer container.Remove(ctx) //nolint:errcheck

		exitCode, err = container.Exec(ctx, containerConfig, workflowCommandStreams)
		if err != nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	return &forge.Cast{
		ExitCode: int64(exitCode),
	}, nil
}
