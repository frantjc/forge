package ore

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/google/uuid"
)

type Alloy struct {
	Id   string      `json:"id,omitempty"`
	Ores []forge.Ore `json:"ores,omitempty"`
}

func (o *Alloy) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) (lava *forge.Lava, err error) {
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
	defer volume.Remove(ctx)

	for _, ore := range o.Ores {
		if lava, err = ore.Liquify(contaminate.WithMounts(ctx, &forge.Mount{
			Source:      volumeName,
			Destination: forge.WorkingDir,
		}), containerRuntime, drains); err != nil {
			break
		}
	}

	return lava, err
}

func (o *Alloy) GetId() string {
	return o.Id
}

func (o *Alloy) GetOres() []forge.Ore {
	return o.Ores
}
