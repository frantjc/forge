package main

import (
	"context"

	"github.com/frantjc/forge/.dagger/internal/dagger"
)

type ForgeDev struct {}

// +check
func (m *ForgeDev) Test(
	ctx context.Context,
	workspace *dagger.Workspace,
	// +optional
	goarch string,
	// +optional
	dockerSock *dagger.Socket,
	// +optional
	docker *dagger.File,
) error {
	tags := []string{"shim"}
	g0, err := m.goWithShim(ctx, workspace, goarch, func(r *dagger.Container) *dagger.Container {
		if dockerSock != nil {
			tags = append(tags, "dockerd")
			return r.
				WithUnixSocket("/var/run/docker.sock", dockerSock).
				With(func(s *dagger.Container) *dagger.Container {
					if docker != nil {
						// FIXME(frantjc): This fails in Dagger.
						tags = append(tags, "docker")
						return s.WithFile("/usr/local/bin/docker", docker)
					}
					return s
				})
		}
		return r
	})
	if err != nil {
		return err
	}
	return g0.
		Test(ctx, dagger.GoTestOpts{
			Race: true,
			Tags: tags,
		})
}

func (m *ForgeDev) Release(
	ctx context.Context,
	workspace *dagger.Workspace,
	githubRepo string,
	githubToken *dagger.Secret,
) error {
	return dag.Release(
		workspace.Directory(".").AsGit().LatestVersion(),
	).
		Create(ctx, githubToken, githubRepo, "forge", dagger.ReleaseCreateOpts{Brew: true})
}

// +check
func (m *ForgeDev) Binary(
	ctx context.Context,
	workspace *dagger.Workspace,
	// +default=v0.0.0-unknown
	version,
	// +optional
	goarch,
	// +optional
	goos string,
) (*dagger.File, error) {
	g0, err := m.goWithShim(ctx, workspace, goarch, func(r *dagger.Container) *dagger.Container { return r })
	if err != nil {
		return nil, err
	}
	return g0.Build(dagger.GoBuildOpts{
		Pkg:     "./cmd/forge",
		Ldflags: "-s -w -X main.version=" + version,
		Goos:    goos,
		Goarch:  goarch,
	}), nil
}

func (m *ForgeDev) goWithShim(
	ctx context.Context,
	workspace *dagger.Workspace,
	// +optional
	goarch string,
	f dagger.WithContainerFunc,
) (*dagger.Go, error) {
	g0 := dag.Go(dagger.GoOpts{
		Workspace: workspace,
	})
	shim, err := dag.Upx().
		Pack(
			g0.
				Build(dagger.GoBuildOpts{
					Pkg: "./internal/cmd/shim",
					Goarch: goarch,
				}),
		).Sync(ctx)
	if err != nil {
		return nil, err
	}
	return dag.Go(dagger.GoOpts{
		Container: g0.Container().
			WithFile("internal/bin/shim", shim).
			With(f),
	}), nil
}
