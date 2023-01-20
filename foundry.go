package forge

import (
	"context"
	"fmt"
)

// NewFoundry returns a Foundry.
func NewFoundry(containerRuntime ContainerRuntime) *Foundry {
	return &Foundry{containerRuntime}
}

// Foundry is a wrapper around a ContainerRuntime for processing Ores.
type Foundry struct {
	ContainerRuntime
}

// Process Liquifies the Ore and returns the resulting Metal.
func (f *Foundry) Process(ctx context.Context, ore Ore, drains *Drains) error {
	if f.ContainerRuntime == nil {
		return fmt.Errorf("nil container runtime")
	}

	_ = LoggerFrom(ctx)

	return ore.Liquify(ctx, f.ContainerRuntime, drains)
}

// GoString implements fmt.GoStringer.
func (f *Foundry) GoString() string {
	return fmt.Sprint("&Foundry{ContainerRuntime: ", f.ContainerRuntime, "}")
}
