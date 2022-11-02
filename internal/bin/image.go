package bin

import (
	"io"

	"github.com/frantjc/forge/pkg/runtime/container/native"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"

	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

const (
	ShimPath  = "/" + ShimName
	ImageName = "forge.dev/shim"
)

var (
	ShimEntrypoint      = []string{ShimPath}
	ShimSleepEntrypoint = append(ShimEntrypoint, "-s")
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

var (
	reference = name.MustParseReference(ImageName)
	baseImage = empty.Image
)

func newShimImage() (*native.Image, error) {
	shimLayer, err := tarball.LayerFromOpener(func() (io.ReadCloser, error) {
		return io.NopCloser(NewShimTarArchive()), nil
	}, tarball.WithCompressionLevel(TarArchiveCompression), tarball.WithCompressedCaching)
	if err != nil {
		return nil, err
	}

	img, err := mutate.AppendLayers(baseImage, shimLayer)
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
