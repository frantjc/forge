package forgeactions

import (
	"context"
	"errors"
	"strings"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/githubactions"
)

const (
	// DefaultNode12ImageReference is the default image to use
	// when an action specifies that it runs using "node12".
	DefaultNode12ImageReference = "docker.io/library/node:12"
	// DefaultNode16ImageReference is the default image to use
	// when an action specifies that it runs using "node16".
	DefaultNode16ImageReference = "docker.io/library/node:16"
	// DefaultNode16ImageReference is the default image to use
	// when an action specifies that it runs using "node20".
	DefaultNode20ImageReference = "docker.io/library/node:20"
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
	// Node16ImageReference is the image to use when an action
	// when an action specifies that it runs using "node20".
	// var so as to be overridable.
	Node20ImageReference = DefaultNode16ImageReference
)

// GetImageForMetadata is a re-export of DefaultMapping.GetImageForMetadata
// for convenience purposes.
func GetImageForMetadata(ctx context.Context, containerRuntime forge.ContainerRuntime, actionMetadata *githubactions.Metadata, uses *githubactions.Uses) (forge.Image, error) {
	return DefaultMapping.GetImageForMetadata(ctx, containerRuntime, actionMetadata, uses)
}

// ImageBuilder is for a ContainerRuntime to implement building a Dockerfile.
// Because building an OCI image is not ubiquitous, forge.ContainerRuntimes are
// not required to implement this, but they may. The default runtime (Docker)
// happens to so as to support GitHub Actions that run using "docker".
type ImageBuilder interface {
	BuildDockerfile(context.Context, string, string) (forge.Image, error)
}

// ErrCantBuildDockerfile will be returned when a forge.ContainerRuntime
// does not implement ImageBuilder.
var ErrCantBuildDockerfile = errors.New("runtime cannot build Dockerfile")

// GetImageForMetadata takes an action.yml and returns the OCI image that forge
// should run it inside of. If the action.yml runs using "dockerfile" and the
// forge.ContainerRuntime does not implement ImageBuilder, returns ErrCantBuildDockerfile.
func (m *Mapping) GetImageForMetadata(ctx context.Context, containerRuntime forge.ContainerRuntime, actionMetadata *githubactions.Metadata, uses *githubactions.Uses) (forge.Image, error) {
	if actionMetadata.IsDockerfile() {
		dir, err := m.UsesToActionDirectory(uses)
		if err != nil {
			return nil, err
		}

		reference := "ghcr.io/" + uses.GetRepository() + ":" + uses.Version
		if uses.IsLocal() {
			// dir will always be an absolute path here
			reference = "forge.dev" + strings.ToLower(dir)
		}

		if imageBuilder, ok := containerRuntime.(ImageBuilder); ok {
			return imageBuilder.BuildDockerfile(ctx, dir, reference)
		}

		return nil, ErrCantBuildDockerfile
	}

	return containerRuntime.PullImage(ctx, MetadataToImageReference(actionMetadata))
}

// MetadataToImageReference takes an action.yaml and finds the reference
// to the OCI image that forge should run it inside of.
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
	case githubactions.RunsUsingNode20:
		return Node20ImageReference
	case githubactions.RunsUsingDocker:
		if strings.HasPrefix(actionMetadata.Runs.Image, githubactions.RunsUsingDockerImagePrefix) {
			return strings.TrimPrefix(actionMetadata.Runs.Image, githubactions.RunsUsingDockerImagePrefix)
		}
	}

	return ""
}
