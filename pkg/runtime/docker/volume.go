package docker

import "github.com/docker/docker/client"

type Volume struct {
	ID string
	*client.Client
}

func (v *Volume) GetID() string {
	return v.ID
}

func (v *Volume) GoString() string {
	return "&Volume{" + v.ID + "}"
}
