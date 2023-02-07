package ore

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/containerfs"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/google/uuid"
)

// Alloy is an Ore made up of other Ores that share
// a volume that is mounted to their WorkingDir.
type Alloy struct {
	ID   string      `json:"id,omitempty"`
	Ores []forge.Ore `json:"ores,omitempty"`
}

func (o *Alloy) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) (err error) {
	volumeName := o.ID
	if volumeName == "" {
		volumeName = uuid.NewString()
	}

	volume, err := containerRuntime.CreateVolume(ctx, volumeName)
	if err != nil {
		return err
	}
	defer volume.Remove(ctx) //nolint:errcheck

	for _, ore := range o.Ores {
		if err = ore.Liquify(contaminate.WithMounts(ctx, &forge.Mount{
			Source:      volumeName,
			Destination: containerfs.WorkingDir,
		}), containerRuntime, drains); err != nil {
			return err
		}
	}

	return nil
}
