package ore

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/frantjc/forge/pkg/concourse"
	"github.com/frantjc/forge/pkg/concourse2container"
)

func (o *Resource) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, basin forge.Basin, drains *forge.Drains) (*forge.Cast, error) {
	image, err := concourse2container.PullImageForResourceType(ctx, containerRuntime, o.GetResourceType())
	if err != nil {
		return nil, err
	}

	containerConfig := concourse2container.ResourceToConfig(o.GetResource(), o.GetResourceType(), o.GetMethod())
	containerConfig.Mounts = contaminate.OverrideWithMountsFrom(ctx, containerConfig.Mounts...)

	container, err := CreateSleepingContainer(ctx, containerRuntime, image, containerConfig)
	if err != nil {
		return nil, err
	}
	defer container.Stop(ctx)   //nolint:errcheck
	defer container.Remove(ctx) //nolint:errcheck

	exitCode, err := container.Exec(ctx, containerConfig, concourse2container.NewStreams(drains, &concourse.Input{
		Params: o.GetParams(),
		Source: o.GetResource().GetSource(),
	}))
	if err != nil {
		return nil, err
	}

	return &forge.Cast{
		ExitCode: int64(exitCode),
	}, nil
}
