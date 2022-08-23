package forge

import (
	"context"
	"fmt"
)

type ContainerRuntime interface {
	fmt.GoStringer

	CreateContainer(context.Context, Image, *ContainerConfig) (Container, error)
	PullImage(context.Context, string) (Image, error)
	CreateVolume(context.Context, string) (Volume, error)
	Close() error
}
