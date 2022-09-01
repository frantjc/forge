package forge

import (
	"context"
	"fmt"
)

// Volume represents a volume created by a ContainerRuntime
// which can be attached to a Container via its ContainerConfig.Mounts
type Volume interface {
	fmt.GoStringer

	Remove(context.Context) error
}
