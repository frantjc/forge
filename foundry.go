package forge

import (
	"context"
	"fmt"
)

// NewFoundry returns a Foundry.
func NewFoundry(containerRuntime ContainerRuntime) *Foundry {
	return &Foundry{containerRuntime}
}

// Foundry is a wrapper around a ContainerRuntime.
type Foundry struct {
	ContainerRuntime
}

// Process checks if its Basin already has the result of an Ore.
// If so, it returns the Metal from the Depoist. Otherwise,
// it Liquifies the Ore, caches the Metal and returns it.
func (f *Foundry) Process(ctx context.Context, ore Ore, drains *Drains) (*Metal, error) {
	if f.ContainerRuntime == nil {
		return nil, fmt.Errorf("nil ContainerRuntime")
	}

	var (
		_      = LoggerFrom(ctx)
		stdout = drains.Out
		stderr = drains.Err
	)

	cast, err := ore.Liquify(ctx, f, &Drains{
		Out: stdout,
		Err: stderr,
		Tty: drains.Tty,
	})
	if err != nil {
		return nil, err
	}

	return &Metal{
		ExitCode: cast.GetExitCode(),
	}, nil
}

// GoString implements fmt.GoStringer.
func (f *Foundry) GoString() string {
	return fmt.Sprint("&Foundry{ContainerRuntime: ", f.ContainerRuntime, "}")
}
