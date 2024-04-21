package forge

import (
	"io"

	"github.com/opencontainers/go-digest"
	imagespecsv1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// Image represents a image pulled by a ContainerRuntime.
// Used to create Containers from.
type Image interface {
	Manifest() (*imagespecsv1.Manifest, error)
	Config() (*imagespecsv1.ImageConfig, error)
	Digest() (digest.Digest, error)
	Blob() io.Reader
	Name() string
}
