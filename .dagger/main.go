package main

import (
	"context"
	"strings"

	"github.com/frantjc/forge/.dagger/internal/dagger"
	xslices "github.com/frantjc/x/slices"
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

func (m *ForgeDev) Fmt(ctx context.Context) *dagger.Changeset {
	goModules := []string{
		".dagger/",
	}

	root := dag.Go(dagger.GoOpts{
		Module: m.Source.Filter(dagger.DirectoryFilterOpts{
			Exclude: goModules,
		}),
	}).
		Container().
		WithExec([]string{"go", "fmt", "./..."}).
		Directory(".")

	for _, module := range goModules {
		root = root.WithDirectory(
			module,
			dag.Go(dagger.GoOpts{
				Module: m.Source.Directory(module).Filter(dagger.DirectoryFilterOpts{
					Exclude: xslices.Filter(goModules, func(m string, _ int) bool {
						return strings.HasPrefix(m, module)
					}),
				}),
			}).
				Container().
				WithExec([]string{"go", "fmt", "./..."}).
				Directory("."),
		)
	}

	return root.Changes(m.Source)
}

func (m *ForgeDev) Test(ctx context.Context) (*dagger.Container, error) {
	return dag.Go(dagger.GoOpts{
		Module: m.Source,
	}).
		Container().
		WithExec([]string{"go", "test", "-race", "-cover", "-test.v", "./..."}), nil
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
		Module: module,
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

	return g0.WithSource(module).
		Build(dagger.GoBuildOpts{
			Pkg:     "./cmd/forge",
			Ldflags: "-s -w -X main.version=" + version,
			Goos:    goos,
			Goarch:  goarch,
		})
}

func (m *ForgeDev) Vulncheck(ctx context.Context) (string, error) {
	return dag.Go(dagger.GoOpts{
		Module: m.Source,
	}).
		Container().
		WithExec([]string{"go", "install", "golang.org/x/vuln/cmd/govulncheck@v1.1.4"}).
		WithExec([]string{"govulncheck", "./..."}).
		CombinedOutput(ctx)
}

func (m *ForgeDev) Vet(ctx context.Context) (string, error) {
	return dag.Go(dagger.GoOpts{
		Module: m.Source,
	}).
		Container().
		WithExec([]string{"go", "vet", "./..."}).
		CombinedOutput(ctx)
}

func (m *ForgeDev) Staticcheck(ctx context.Context) (string, error) {
	return dag.Go(dagger.GoOpts{
		Module: m.Source,
	}).
		Container().
		WithExec([]string{"go", "install", "honnef.co/go/tools/cmd/staticcheck@v0.6.1"}).
		WithExec([]string{"staticcheck", "./..."}).
		CombinedOutput(ctx)
}
