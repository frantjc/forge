package ore

import (
	"context"

	"github.com/frantjc/forge"
	cfs "github.com/frantjc/forge/internal/containerfs"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/google/uuid"
)

// Alloy is an Ore made up of other Ores that share
// a volume that is mounted to their WorkingDir.
type Alloy struct {
	ID   string      `json:"id,omitempty"`
	Ores []forge.Ore `json:"ores,omitempty"`
}

func (o *Alloy) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) (metal *forge.Metal, err error) {
	volumeName := o.ID
	if volumeName == "" {
		volumeName = uuid.NewString()
	}

	volume, err := containerRuntime.CreateVolume(ctx, volumeName)
	if err != nil {
		return nil, err
	}
	defer volume.Remove(ctx) //nolint:errcheck

	for _, ore := range o.Ores {
		if metal, err = ore.Liquify(contaminate.WithMounts(ctx, &forge.Mount{
			Source:      volumeName,
			Destination: cfs.WorkingDir,
		}), containerRuntime, drains); err != nil {
			break
		}
	}

	return metal, err
}
