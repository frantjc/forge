package native

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func PullImage(ctx context.Context, reference string) (forge.Image, error) {
	ref, err := name.ParseReference(reference)
	if err != nil {
		return nil, err
	}

	img, err := remote.Image(ref, remote.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return &Image{img, ref}, nil
}
