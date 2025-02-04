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

const (
	DefaultNode10ImageReference = "docker.io/library/node:10"
	DefaultNode12ImageReference = "docker.io/library/node:12"
	DefaultNode16ImageReference = "docker.io/library/node:16"
	DefaultNode20ImageReference = "docker.io/library/node:20"
	DefaultNodeImageReference   = DefaultNode16ImageReference
)

var (
	Node10ImageReference = DefaultNode10ImageReference
	Node12ImageReference = DefaultNode12ImageReference
	Node16ImageReference = DefaultNode16ImageReference
	Node20ImageReference = DefaultNode20ImageReference
	NodeImageReference   = DefaultNodeImageReference
)
