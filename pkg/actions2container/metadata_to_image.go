package actions2container

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/github/actions"
)

type RunsUsing string

var (
	Node12ImageReference = "docker.io/library/node:12"
	Node16ImageReference = "docker.io/library/node:16"
)

func PullImageForMetadata(ctx context.Context, containerRuntime forge.ContainerRuntime, actionMetadata *actions.Metadata) (forge.Image, error) {
	return containerRuntime.PullImage(ctx, MetadataToImageReference(actionMetadata))
}

func MetadataToImageReference(actionMetadata *actions.Metadata) string {
	if actionMetadata == nil {
		return ""
	}

	if actionMetadata.Runs == nil {
		return ""
	}

	return RunsUsing(actionMetadata.Runs.Using).ImageReference(actionMetadata.Runs.Image)
}

func (r RunsUsing) ImageReference(image string) string {
	switch r {
	case actions.RunsUsingNode12:
		return Node12ImageReference
	case actions.RunsUsingNode16:
		return Node16ImageReference
	}

	return image
}
