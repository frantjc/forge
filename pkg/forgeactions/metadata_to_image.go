package forgeactions

import (
	"context"
	"strings"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/githubactions"
	"github.com/frantjc/go-fn"
)

const (
	// DefaultNode12ImageReference is the default image to use
	// when an action specifies that it runs using "node12".
	DefaultNode12ImageReference = "docker.io/library/node:12"
	// DefaultNode16ImageReference is the default image to use
	// when an action specifies that it runs using "node16".
	DefaultNode16ImageReference = "docker.io/library/node:16"
)

var (
	// Node12ImageReference is the image to use when an action
	// when an action specifies that it runs using "node12".
	// var so as to be overridable.
	Node12ImageReference = DefaultNode12ImageReference
	// Node16ImageReference is the image to use when an action
	// when an action specifies that it runs using "node16".
	// var so as to be overridable.
	Node16ImageReference = DefaultNode16ImageReference
)

func PullImageForMetadata(ctx context.Context, containerRuntime forge.ContainerRuntime, actionMetadata *githubactions.Metadata) (forge.Image, error) {
	return containerRuntime.PullImage(ctx, MetadataToImageReference(actionMetadata))
}

func MetadataToImageReference(actionMetadata *githubactions.Metadata) string {
	if actionMetadata == nil {
		return ""
	}

	if actionMetadata.Runs == nil {
		return ""
	}

	return RunsUsingImage(actionMetadata.Runs.Using, strings.TrimPrefix(actionMetadata.Runs.Image, githubactions.RunsUsingDockerImagePrefix))
}

func RunsUsingImage(runsUsing string, fallbacks ...string) string {
	switch runsUsing {
	case githubactions.RunsUsingNode12:
		return Node12ImageReference
	case githubactions.RunsUsingNode16:
		return Node16ImageReference
	}

	return fn.Coalesce(fallbacks...)
}
