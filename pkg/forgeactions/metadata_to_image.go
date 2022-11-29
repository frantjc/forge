package forgeactions

import (
	"context"
	"path"
	"strings"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/githubactions"
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

func GetImageForMetadata(ctx context.Context, containerRuntime forge.ContainerRuntime, actionMetadata *githubactions.Metadata, uses *githubactions.Uses) (forge.Image, error) {
	return DefaultMapping.GetImageForMetadata(ctx, containerRuntime, actionMetadata, uses)
}

func (m *Mapping) GetImageForMetadata(ctx context.Context, containerRuntime forge.ContainerRuntime, actionMetadata *githubactions.Metadata, uses *githubactions.Uses) (forge.Image, error) {
	if actionMetadata.IsDockerfile() {
		dir, err := m.UsesToActionDirectory(uses)
		if err != nil {
			return nil, err
		}

		reference := "ghcr.io/" + uses.GetRepository() + ":" + uses.Version
		if uses.IsLocal() {
			reference = path.Join("forge.dev", strings.ToLower(actionMetadata.Name))
		}

		return containerRuntime.BuildDockerfile(ctx, dir, reference)
	}

	return containerRuntime.PullImage(ctx, MetadataToImageReference(actionMetadata))
}

func MetadataToImageReference(actionMetadata *githubactions.Metadata) string {
	if actionMetadata == nil {
		return ""
	}

	if actionMetadata.Runs == nil {
		return ""
	}

	switch actionMetadata.Runs.Using {
	case githubactions.RunsUsingNode12:
		return Node12ImageReference
	case githubactions.RunsUsingNode16:
		return Node16ImageReference
	case githubactions.RunsUsingDocker:
		if strings.HasPrefix(actionMetadata.Runs.Image, githubactions.RunsUsingDockerImagePrefix) {
			return strings.TrimPrefix(actionMetadata.Runs.Image, githubactions.RunsUsingDockerImagePrefix)
		}
	}

	return ""
}
