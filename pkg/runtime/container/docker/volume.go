package docker

import "github.com/docker/docker/client"

type Volume struct {
	ID string
	*client.Client
}
