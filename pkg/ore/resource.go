package ore

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/concourse"
	"github.com/frantjc/forge/pkg/concourse2container"
)

func (o *Resource) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) (*forge.Lava, error) {
	container, err := concourse2container.CreateContainerForResource(ctx, containerRuntime, o.GetResource(), o.GetResourceType(), o.GetMethod())
	if err != nil {
		return nil, err
	}

	exitCode, err := container.Run(ctx, concourse2container.NewStreams(drains, &concourse.Input{
		Params: o.GetParams(),
		Source: o.GetResource().GetSource(),
	}))
	if err != nil {
		return nil, err
	}

	return &forge.Lava{
		ExitCode: int64(exitCode),
	}, nil
}
