package forge

import (
	"context"
)

// Ore represents one or more sequential containerized commands.
// Ores are meant to represent the entire input to said commands,
// so that if two Ore's match, their results should be the same.
// Because of this, Ores can be cached.
type Ore interface {
	Liquify(context.Context, ContainerRuntime, *Drains) error
}
