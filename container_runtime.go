package forge

import (
	"context"
	"fmt"
)

type ContainerRuntime interface {
	fmt.GoStringer

	GetContainer(context.Context, string) (Container, error)
	CreateContainer(context.Context, Image, *ContainerConfig) (Container, error)
	PullImage(context.Context, string) (Image, error)
	CreateVolume(context.Context, string) (Volume, error)
	Close() error
}
