package forgeconcourse

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/concourse"
)

var DefaultTag = "latest"

func PullImageForResourceType(ctx context.Context, containerRuntime forge.ContainerRuntime, resourceType *concourse.ResourceType) (forge.Image, error) {
	return containerRuntime.PullImage(ctx, ResourceTypeToImageReference(resourceType))
}

func ResourceTypeToImageReference(resourceType *concourse.ResourceType) string {
	if resourceType != nil && resourceType.Source != nil {
		tag := resourceType.Source.Tag
		if tag == "" {
			tag = DefaultTag
		}

		return resourceType.Source.Repository + ":" + tag
	}

	return ""
}
