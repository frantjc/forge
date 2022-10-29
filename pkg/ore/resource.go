package ore

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/frantjc/forge/pkg/concourse"
	fc "github.com/frantjc/forge/pkg/forgeconcourse"
)

func (o *Resource) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) (*forge.Metal, error) {
	image, err := fc.PullImageForResourceType(ctx, containerRuntime, o.GetResourceType())
	if err != nil {
		return nil, err
	}

	containerConfig := fc.ResourceToConfig(o.GetResource(), o.GetResourceType(), o.GetMethod())
	containerConfig.Mounts = contaminate.OverrideWithMountsFrom(ctx, containerConfig.GetMounts()...)

	container, err := CreateSleepingContainer(ctx, containerRuntime, image, containerConfig)
	if err != nil {
		return nil, err
	}
	defer container.Stop(ctx)   //nolint:errcheck
	defer container.Remove(ctx) //nolint:errcheck

	exitCode, err := container.Exec(ctx, containerConfig, fc.NewStreams(drains, &concourse.Input{
		Params:  o.GetParams(),
		Source:  o.GetResource().GetSource(),
		Version: o.GetVersion(),
	}))
	if err != nil {
		return nil, err
	}

	return &forge.Metal{
		ExitCode: int64(exitCode),
	}, nil
}
