package forge

import "context"

// Volume represents a volume created by a ContainerRuntime
// which can be attached to a Container via its ContainerConfig.Mounts.
type Volume interface {
	GetID() string
	Remove(context.Context) error
}
