package ore

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/digestutil"
	"github.com/frantjc/forge/internal/hostfs"
)

// Cache is an Ore that caches the output of other Ores.
type Cache struct {
	forge.Ore
}

func (o *Cache) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) error {
	var (
		_      = forge.LoggerFrom(ctx)
		cached = true
	)

	d, err := digestutil.JSON(o.Ore)
	if err != nil {
		return err
	}

	cache := filepath.Join(hostfs.OreCache, d.Encoded())
	if err := os.MkdirAll(cache, 0o755); err != nil {
		return err
	}

	var (
		outname = filepath.Join(cache, "stdout")
		errname = filepath.Join(cache, "stderr")
	)
	stdout, err := os.Open(outname)
	if err != nil {
		cached = false
		stdout, err = os.Create(outname)
		if err != nil {
			return err
		}
	}
	defer stdout.Close()

	stderr, err := os.Open(errname)
	if err != nil {
		cached = false
		stderr, err = os.Create(errname)
		if err != nil {
			return err
		}
	}
	defer stderr.Close()

	if cached {
		errC := make(chan error, 2)

		go func() {
			_, gErr := io.Copy(drains.Out, stdout)
			errC <- gErr
		}()

		go func() {
			_, gErr := io.Copy(drains.Err, stderr)
			errC <- gErr
		}()

		for i := 0; i < 2; i++ {
			if err := <-errC; err != nil {
				return err
			}
		}

		return nil
	}

	err = o.Ore.Liquify(ctx, containerRuntime, &forge.Drains{
		Out: io.MultiWriter(stdout, drains.Out),
		Err: io.MultiWriter(stderr, drains.Err),
		Tty: drains.Tty,
	})
	if err != nil {
		defer os.Remove(outname)
		defer os.Remove(errname)
	}

	return err
}
