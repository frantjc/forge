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
	if resourceType != nil && resourceType.GetSource() != nil {
		tag := resourceType.GetSource().GetTag()
		if tag == "" {
			tag = DefaultTag
		}

		return resourceType.GetSource().GetRepository() + ":" + tag
	}

	return ""
}
