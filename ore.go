package forge

import "context"

// Ore represents one or more sequential containerized commands.
// Ores are meant to represent the entire input to said commands,
// so that if two Ore's digests match, their resulting Lava
// should be the same. Because of this, Ores can be cached,
// using said Digest as the key
type Ore interface {
	Liquify(context.Context, ContainerRuntime, *Drains) (*Lava, error)
}
