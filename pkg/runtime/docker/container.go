package docker

import "github.com/docker/docker/client"

type Container struct {
	ID string
	*client.Client
}

//nolint:revive // matching protobuf style
func (c *Container) GetId() string {
	return c.ID
}

func (c *Container) GoString() string {
	return "&Container{" + c.GetId() + "}"
}
