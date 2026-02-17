package main

import (
	"context"

	"github.com/frantjc/forge/.dagger/internal/dagger"
)

type ForgeDev struct {
	Source *dagger.Directory
}

func New(
	ctx context.Context,
	// +optional
	// +defaultPath="."
	src *dagger.Directory,
) (*ForgeDev, error) {
	return &ForgeDev{
		Source: src,
	}, nil
}

func (m *ForgeDev) Release(
	ctx context.Context,
	githubRepo string,
	githubToken *dagger.Secret,
) error {
	return dag.Release(m.Source.AsGit().LatestVersion()).Create(ctx, githubToken, githubRepo, "forge", dagger.ReleaseCreateOpts{Brew: true})
}

func (m *ForgeDev) Binary(
	ctx context.Context,
	// +default=v0.0.0-unknown
	version string,
	// +optional
	goarch string,
	// +optional
	goos string,
) *dagger.File {
	module := m.Source

	g0 := dag.Go(dagger.GoOpts{
		Source: m.Source,
	})
	upx := dag.Upx()

	module = module.WithFile(
		"internal/bin/shim",
		upx.
			Pack(
				g0.
					Build(dagger.GoBuildOpts{
						Pkg:    "./internal/cmd/shim",
						Goarch: goarch,
					}),
				dagger.UpxPackOpts{Brute: true},
			),
	)

	return g0.
		Build(dagger.GoBuildOpts{
			Pkg:     "./cmd/forge",
			Ldflags: "-s -w -X main.version=" + version,
			Goos:    goos,
			Goarch:  goarch,
		})
}
