package concourse2container

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/frantjc/forge/pkg/concourse"
)

func CreateContainerForResource(ctx context.Context, containerRuntime forge.ContainerRuntime, resource *concourse.Resource, resourceType *concourse.ResourceType, method string) (forge.Container, error) {
	image, err := PullImageForResourceType(ctx, containerRuntime, resourceType)
	if err != nil {
		return nil, err
	}

	containerConfig := ResourceToConfig(resource, resourceType, method)
	containerConfig.Mounts = append(containerConfig.Mounts, contaminate.MountsFrom(ctx)...)

	return containerRuntime.CreateContainer(ctx, image, containerConfig)
}
