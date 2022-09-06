package forge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/sync/errgroup"
)

// Foundry is a wrapper around a ContainerRuntime
// and a Basin meant to process and cache Ores.
type Foundry struct {
	ContainerRuntime
	Basin
}

// Process checks if its Basin already has the result of an Ore.
// If so, it returns the Metal from the Depoist. Otherwise,
// it Liquifies the Ore, caches the Metal and returns it.
func (f *Foundry) Process(ctx context.Context, ore Ore, drains *Drains) (*Metal, error) {
	if f.ContainerRuntime == nil {
		return nil, fmt.Errorf("nil ContainerRuntime")
	}

	var (
		logr        = LoggerFrom(ctx)
		stdout      = drains.Out
		stderr      = drains.Err
		digest, err = Digest(ore)
	)
	if err != nil {
		return nil, err
	}

	if f.Basin != nil {
		if stdoutCache, err := f.Basin.NewReader(ctx, digest.Encoded()+"/stdout.txt"); err == nil {
			defer stdoutCache.Close()

			if stderrCache, err := f.Basin.NewReader(ctx, digest.Encoded()+"/stderr.txt"); err == nil {
				defer stderrCache.Close()

				if metalCache, err := f.Basin.NewReader(ctx, digest.Encoded()+"/metal.json"); err == nil {
					defer metalCache.Close()

					logr.Info("[cached] " + digest.String())

					var (
						eg, _ = errgroup.WithContext(ctx)
						metal = &Metal{}
					)

					eg.Go(func() error {
						_, gerr := io.Copy(stdout, stdoutCache)
						return gerr
					})

					eg.Go(func() error {
						_, gerr := io.Copy(stderr, stderrCache)
						return gerr
					})

					eg.Go(func() error {
						return json.NewDecoder(metalCache).Decode(metal)
					})

					return metal, eg.Wait()
				}
			}
		}

		stdoutCache, err := f.Basin.NewWriter(ctx, digest.Encoded()+"/stdout.txt")
		if err != nil {
			return nil, err
		}
		defer stdoutCache.Close()

		stderrCache, err := f.Basin.NewWriter(ctx, digest.Encoded()+"/stderr.txt")
		if err != nil {
			return nil, err
		}
		defer stderrCache.Close()

		stdout = io.MultiWriter(stdout, stdoutCache)
		stderr = io.MultiWriter(stderr, stderrCache)
	}

	lava, err := ore.Liquify(ctx, f, f, &Drains{
		Out: stdout,
		Err: stderr,
		Tty: drains.Tty,
	})
	if err != nil {
		return nil, err
	}

	metal := &Metal{
		ExitCode: lava.GetExitCode(),
	}

	if f.Basin != nil {
		metalCache, err := f.Basin.NewWriter(ctx, digest.Encoded()+"/metal.json")
		if err != nil {
			return nil, err
		}
		defer metalCache.Close()

		return metal, json.NewEncoder(metalCache).Encode(metal)
	}

	return metal, nil
}

// GoString implements fmt.GoStringer.
func (f *Foundry) GoString() string {
	return "&Foundry{ContainerRuntime: " + f.ContainerRuntime.GoString() + ", Basin: " + f.Basin.GoString() + "}"
}
