package forge

import "context"

type Ore interface {
	Liquify(context.Context, ContainerRuntime, *Streams) (*Lava, error)
}
