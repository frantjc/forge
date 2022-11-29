package docker

import "github.com/docker/docker/client"

type Container struct {
	ID string
	*client.Client
}

func (c *Container) GetID() string {
	return c.ID
}

func (c *Container) GoString() string {
	return "&Container{" + c.GetID() + "}"
}
