package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/image"
	"github.com/frantjc/forge"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
)

func (d *ContainerRuntime) PullImage(ctx context.Context, reference string) (forge.Image, error) {
	ref, err := name.ParseReference(reference)
	if err != nil {
		return nil, err
	}

	r, err := d.Client.ImagePull(ctx, ref.Name(), image.PullOptions{})
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(io.Discard, r); err != nil {
		return nil, err
	}

	img, err := daemon.Image(ref, daemon.WithClient(d), daemon.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return &Image{
		Image:     img,
		Reference: ref,
	}, nil
}
