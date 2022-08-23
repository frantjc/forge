package bin

import (
	"github.com/frantjc/forge/pkg/runtime/container/native"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"

	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

const (
	ShimPath = "/" + ShimName
)

var (
	reference = name.MustParseReference("forge.dev/shim")
	base      = empty.Image
)

func init() {
	_ = NewShimImage()
}

func NewShimImage() *native.Image {
	img, err := newShimImage()
	if err != nil {
		panic(err)
	}

	return img
}

func newShimImage() (*native.Image, error) {
	shimLayer, err := tarball.LayerFromOpener(newShimTarArchive, tarball.WithCompressionLevel(Compression))
	if err != nil {
		return nil, err
	}

	img, err := mutate.AppendLayers(base, shimLayer)
	if err != nil {
		return nil, err
	}

	if img, err = mutate.Config(img, v1.Config{
		Image: reference.Name(),
	}); err != nil {
		return nil, err
	}

	return &native.Image{Image: img, Reference: reference}, nil
}
