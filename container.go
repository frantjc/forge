package forge

import (
	"context"
	"io"
)

// Container represents a container created by a ContainerRuntime.
type Container interface {
	GetID() string

	CopyTo(context.Context, string, io.Reader) error
	CopyFrom(context.Context, string) (io.ReadCloser, error)

	Run(context.Context, *Streams) (int, error)
	Start(context.Context) error
	Restart(context.Context) error
	Exec(context.Context, *ContainerConfig, *Streams) (int, error)

	Stop(context.Context) error
	Remove(context.Context) error
	Kill(context.Context) error
}
