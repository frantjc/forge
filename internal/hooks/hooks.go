package hooks

import "github.com/frantjc/forge"

var (
	ContainerCreated = new(Hook[forge.Container])
	ContainerStarted = new(Hook[forge.Container])
)
