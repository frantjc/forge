package docker

import (
	"context"

	"github.com/docker/docker/api/types/volume"
	"github.com/frantjc/forge"
)

const (
	VolumeDriver = "local"
)

func (r *ContainerRuntime) CreateVolume(ctx context.Context, name string) (forge.Volume, error) {
	v, err := r.VolumeCreate(ctx, volume.CreateOptions{
		Driver: VolumeDriver,
		Name:   name,
	})
	if err != nil {
		return nil, err
	}

	return &Volume{v.Name, r.Client}, nil
}
