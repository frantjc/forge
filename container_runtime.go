package forge

import "context"

// ContainerRuntime represents the functionality needed by Ores
// to pull OCI images and run containers when being processed.
type ContainerRuntime interface {
	GetContainer(context.Context, string) (Container, error)
	CreateContainer(context.Context, Image, *ContainerConfig) (Container, error)
	PullImage(context.Context, string) (Image, error)
	CreateVolume(context.Context, string) (Volume, error)
	Close() error
}
