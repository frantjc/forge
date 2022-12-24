package ore

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/concourse"
	"github.com/frantjc/forge/forgeconcourse"
	"github.com/frantjc/forge/internal/containerutil"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/frantjc/go-fn"
)

// Resource is an Ore representing a Concourse Resource--
// any of get, put or check.
type Resource struct {
	Method       string                  `json:"method,omitempty"`
	Version      map[string]string       `json:"version,omitempty"`
	Params       map[string]string       `json:"params,omitempty"`
	Resource     *concourse.Resource     `json:"resource,omitempty"`
	ResourceType *concourse.ResourceType `json:"resource_type,omitempty"`
}

func (o *Resource) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) (*forge.Metal, error) {
	image, err := forgeconcourse.PullImageForResourceType(ctx, containerRuntime, o.ResourceType)
	if err != nil {
		return nil, err
	}

	containerConfig := forgeconcourse.ResourceToConfig(o.Resource, o.ResourceType, o.Method)
	containerConfig.Mounts = contaminate.OverrideWithMountsFrom(ctx, containerConfig.Mounts...)

	container, err := containerutil.CreateSleepingContainer(ctx, containerRuntime, image, containerConfig)
	if err != nil {
		return nil, err
	}
	defer container.Stop(ctx)   //nolint:errcheck
	defer container.Remove(ctx) //nolint:errcheck

	exitCode, err := container.Exec(ctx, containerConfig, forgeconcourse.NewStreams(drains, &concourse.Input{
		Params: fn.Ternary(
			o.Method == forgeconcourse.MethodCheck,
			nil, o.Params,
		),
		Source: o.Resource.Source,
		Version: fn.Ternary(
			o.Method == forgeconcourse.MethodPut,
			nil, o.Version,
		),
	}))
	if err != nil {
		return nil, err
	}

	return &forge.Metal{
		ExitCode: int64(exitCode),
	}, nil
}
