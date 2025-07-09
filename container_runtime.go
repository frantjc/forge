package forge

import (
	"context"
	"io"

	xslices "github.com/frantjc/x/slices"
	imagespecsv1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type Mount struct {
	Source      string `json:"source,omitempty"`
	Destination string `json:"destination,omitempty"`
}

func overrideMounts(oldMounts []Mount, newMounts ...Mount) []Mount {
	return append(xslices.Filter(oldMounts, func(m Mount, _ int) bool {
		return !xslices.Some(newMounts, func(n Mount, _ int) bool {
			return m.Destination == n.Destination
		})
	}), newMounts...)
}

// ContainerConfig is the configuration that is used to
// create a container or an exec in a running container.
type ContainerConfig struct {
	Entrypoint []string
	Cmd        []string
	WorkingDir string
	Env        []string
	User       string
	Privileged bool
	Mounts     []Mount
}

// Container represents a container created by a ContainerRuntime.
type Container interface {
	GetID() string
	CopyTo(context.Context, string, io.Reader) error
	CopyFrom(context.Context, string) (io.ReadCloser, error)
	Start(context.Context) error
	Exec(context.Context, *ContainerConfig, *Streams) (int, error)
	Stop(context.Context) error
	Remove(context.Context) error
}

// Image represents a image pulled by a ContainerRuntime.
// Used to create Containers from.
type Image interface {
	Config() (*imagespecsv1.ImageConfig, error)
	Blob() io.Reader
	Name() string
}

// ContainerRuntime represents the functionality needed by Runnables
// to pull OCI images and run containers when being processed.
type ContainerRuntime interface {
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
