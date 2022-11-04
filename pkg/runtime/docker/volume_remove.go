package docker

import "context"

func (v *Volume) Remove(ctx context.Context) error {
	return v.VolumeRemove(ctx, v.ID, true)
}
