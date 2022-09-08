package docker

import "github.com/docker/docker/client"

type Volume struct {
	ID string
	*client.Client
}

//nolint:revive // matching protobuf style
func (v *Volume) GetId() string {
	return v.ID
}

func (v *Volume) GoString() string {
	return "&Volume{" + v.ID + "}"
}
