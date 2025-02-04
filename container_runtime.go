package forge

import "context"

// ContainerRuntime represents the functionality needed by Ores
// to pull OCI images and run containers when being processed.
type ContainerRuntime interface {
	GetContainer(context.Context, string) (Container, error)
	CreateContainer(context.Context, Image, *ContainerConfig) (Container, error)
	PullImage(context.Context, string) (Image, error)
	Close() error
}

// ImageBuilder is for a ContainerRuntime to implement building a Dockerfile.
// Because building an OCI image is not ubiquitous, forge.ContainerRuntimes are
// not required to implement this, but they may. The default runtime (Docker)
// happens to so as to support GitHub Actions that run using "docker".
type ImageBuilder interface {
	BuildDockerfile(context.Context, string, string) (Image, error)
}
