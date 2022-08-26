package forge

import "context"

type Ore interface {
	Liquify(context.Context, ContainerRuntime, *Drains) (*Lava, error)
}
