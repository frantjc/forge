package ore

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/concourse"
	"github.com/frantjc/forge/forgeconcourse"
	"github.com/frantjc/forge/internal/containerutil"
	"github.com/frantjc/forge/internal/contaminate"
	xos "github.com/frantjc/x/os"
)

// Resource is an Ore representing a Concourse Resource--
// any of get, put or check.
type Resource struct {
	Method       string                  `json:"method,omitempty"`
	Version      map[string]any          `json:"version,omitempty"`
	Params       map[string]any          `json:"params,omitempty"`
	Resource     *concourse.Resource     `json:"resource,omitempty"`
	ResourceType *concourse.ResourceType `json:"resource_type,omitempty"`
}

func (o *Resource) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) error {
	image, err := forgeconcourse.PullImageForResourceType(ctx, containerRuntime, o.ResourceType)
	if err != nil {
		return err
	}

	containerConfig := forgeconcourse.ResourceToConfig(o.Resource, o.ResourceType, o.Method)
	containerConfig.Mounts = contaminate.OverrideWithMountsFrom(ctx, containerConfig.Mounts...)

	container, err := containerutil.CreateSleepingContainer(ctx, containerRuntime, image, containerConfig)
	if err != nil {
		return err
	}
	defer container.Stop(ctx)   //nolint:errcheck
	defer container.Remove(ctx) //nolint:errcheck

	if exitCode, err := container.Exec(ctx, containerConfig, forgeconcourse.NewStreams(drains, &concourse.Input{
		Params: func() map[string]any {
			if o.Method == forgeconcourse.MethodCheck {
				return nil
			}

			return o.Params
		}(),
		Source: o.Resource.Source,
		Version: func() map[string]any {
			if o.Method == forgeconcourse.MethodPut {
				return nil
			}

			return o.Version
		}(),
	})); err != nil {
		return err
	} else if exitCode > 0 {
		return xos.NewExitCodeError(ErrContainerExitedWithNonzeroExitCode, exitCode)
	}

	return nil
}
