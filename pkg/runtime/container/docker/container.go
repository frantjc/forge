package docker

import "github.com/docker/docker/client"

type Container struct {
	ID string
	*client.Client
}

func (c *Container) GoString() string {
	return "&Container{ID: " + c.ID + "}"
}
