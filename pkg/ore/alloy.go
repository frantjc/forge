package ore

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/google/uuid"
)

type Alloy struct {
	Id   string      `json:"id,omitempty"` //nolint:revive // matching protobuf style
	Ores []forge.Ore `json:"ores,omitempty"`
}

func (o *Alloy) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, basin forge.Basin, drains *forge.Drains) (lava *forge.Cast, err error) {
	var (
		volumeName = o.Id
	)
	if volumeName == "" {
		volumeName = uuid.NewString()
	}

	volume, err := containerRuntime.CreateVolume(ctx, volumeName)
	if err != nil {
		return nil, err
	}
	defer volume.Remove(ctx) //nolint:errcheck

	for _, ore := range o.Ores {
		if lava, err = ore.Liquify(contaminate.WithMounts(ctx, &forge.Mount{
			Source:      volumeName,
			Destination: forge.WorkingDir,
		}), containerRuntime, basin, drains); err != nil {
			break
		}
	}

	return lava, err
}

//nolint:revive // matching protobuf style
func (o *Alloy) GetId() string {
	return o.Id
}

func (o *Alloy) GetOres() []forge.Ore {
	return o.Ores
}
