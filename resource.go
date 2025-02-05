package forge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/frantjc/forge/concourse"
	xos "github.com/frantjc/x/os"
)

// Resource is an Ore representing a Concourse resource--
// any of get, put or check.
type Resource struct {
	Method       string
	Version      map[string]any
	Params       map[string]any
	Resource     *concourse.Resource
	ResourceType *concourse.ResourceType
}

func (o *Resource) Liquify(ctx context.Context, containerRuntime ContainerRuntime, opts ...OreOpt) error {
	opt := oreOptsWithDefaults(opts...)

	image, err := containerRuntime.PullImage(ctx, resourceTypeToImageReference(o.ResourceType))
	if err != nil {
		return err
	}

	containerConfig := resourceToConfig(o.Resource, o.ResourceType, o.Method, opt)
	containerConfig.Mounts = overrideMounts(containerConfig.Mounts, opt.Mounts...)

	container, err := createSleepingContainer(ctx, containerRuntime, image, containerConfig, opt)
	if err != nil {
		return err
	}
	defer container.Stop(ctx)   //nolint:errcheck
	defer container.Remove(ctx) //nolint:errcheck

	streams, err := resourceStreams(opt.Streams, &concourse.Input{
		Params: func() map[string]any {
			if o.Method == concourse.MethodCheck {
				return nil
			}

			return o.Params
		}(),
		Source: o.Resource.Source,
		Version: func() map[string]any {
			if o.Method == concourse.MethodPut {
				return nil
			}

			return o.Version
		}(),
	})
	if err != nil {
		return err
	}

	if exitCode, err := container.Exec(ctx, containerConfig, streams); err != nil {
		return err
	} else if exitCode > 0 {
		return xos.NewExitCodeError(ErrContainerExitedWithNonzeroExitCode, exitCode)
	}

	return nil
}

func resourceStreams(streams *Streams, input *concourse.Input) (*Streams, error) {
	if streams.In != nil {
		return nil, fmt.Errorf("stdin not supported for Concourse Resources")
	}

	in := new(bytes.Buffer)

	if err := json.NewEncoder(in).Encode(input); err != nil {
		return nil, err
	}

	return &Streams{
		In:         in,
		Out:        streams.Out,
		Err:        streams.Err,
		Tty:        streams.Tty,
		DetachKeys: streams.DetachKeys,
	}, nil
}

func resourceToConfig(resource *concourse.Resource, resourceType *concourse.ResourceType, method string, opt *OreOpts) *ContainerConfig {
	return &ContainerConfig{
		Entrypoint: concourse.GetEntrypoint(method),
		Cmd:        []string{filepath.Join(ConcourseResourceWorkingDir(opt.WorkingDir), resource.Name)},
		Privileged: resourceType.Privileged,
		Mounts: []Mount{
			{
				Destination: filepath.Join(ConcourseResourceWorkingDir(opt.WorkingDir), resource.Name),
			},
		},
	}
}

func resourceTypeToImageReference(resourceType *concourse.ResourceType) string {
	if resourceType != nil && resourceType.Source != nil {
		tag := resourceType.Source.Tag
		if tag == "" {
			tag = "latest"
		}

		return resourceType.Source.Repository + ":" + tag
	}

	return ""
}
