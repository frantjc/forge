package forge

import "context"

type ContainerRuntime interface {
	GetContainer(context.Context, string) (Container, error)
	CreateContainer(context.Context, Image, *ContainerConfig) (Container, error)
	PullImage(context.Context, string) (Image, error)
	BuildImage(context.Context, string, string) (Image, error)
	CreateVolume(context.Context, string) (Volume, error)
	Close() error
}
