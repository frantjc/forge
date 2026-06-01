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
	source *dagger.Directory,
) (*ForgeDev, error) {
	return &ForgeDev{
		Source: source,
	}, nil
}

func (m *ForgeDev) SourceWithShim(
	// +optional
	goarch string,
) *dagger.Directory {
	return m.Source.WithFile(
		"internal/bin/shim",
		dag.Upx().
			Pack(
				dag.Go(dagger.GoOpts{
					Source: m.Source,
				}).
					Build(dagger.GoBuildOpts{
						Cgo: 	false,
						Goarch: goarch,
					}),
			),
	)
}

func (m *ForgeDev) Test(
	ctx context.Context,
	// +optional
	dockerSock *dagger.Socket,
	// +optional
	docker *dagger.File,
) error {
	tags := []string{"shim"}
	return dag.Go(dagger.GoOpts{
		Container: dag.Go(dagger.GoOpts{
			Source: m.SourceWithShim(""),
		}).
			Container().
			With(func(r *dagger.Container) *dagger.Container {
				if dockerSock != nil {
					tags = append(tags, "dockerd")
					return r.
						WithUnixSocket("/var/run/docker.sock", dockerSock).
						With(func(s *dagger.Container) *dagger.Container {
							if docker != nil {
								tags = append(tags, "docker")
								return s.WithFile("/usr/local/bin/docker", docker)
							}
							return s
						})
				}
				return r
			}),
	}).
		Test(ctx, dagger.GoTestOpts{
			Race: true,
			Tags: tags,
		})
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
	version,
	// +optional
	goarch,
	// +optional
	goos string,
) *dagger.File {
	return dag.Go(dagger.GoOpts{
		Source: m.SourceWithShim(goarch),
	}).
		Build(dagger.GoBuildOpts{
			Pkg:     "./cmd/forge",
			Ldflags: "-s -w -X main.version=" + version,
			Goos:    goos,
			Goarch:  goarch,
		})
}
