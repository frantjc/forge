package actions2container

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/fn"
	"github.com/frantjc/forge/pkg/github/actions"
)

var (
	Node12ImageReference = "docker.io/library/node:12" // -alpine"
	Node16ImageReference = "docker.io/library/node:16" // -alpine"
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

	return RunsUsingImage(actionMetadata.Runs.Using, actionMetadata.Runs.Image)
}

func RunsUsingImage(runsUsing string, fallbacks ...string) string {
	switch runsUsing {
	case actions.RunsUsingNode12:
		return Node12ImageReference
	case actions.RunsUsingNode16:
		return Node16ImageReference
	}

	return fn.Coalesce(fallbacks...)
}
