package main

import (
	"context"
	"fmt"
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
		Module:                  m.Source,
	}).
		Container().
		WithExec([]string{"go", "test", "-race", "-cover", "-test.v", "./..."}), nil
}

func (m *ForgeDev) Release(ctx context.Context, githubToken *dagger.Secret) error {
	gitRef := m.Source.AsGit().LatestVersion()

	ref, err := gitRef.Ref(ctx)
	if err != nil {
		return err
	}

	tag := strings.TrimPrefix(ref, "refs/tags/")

	release := dag.Wolfi().
		Container(dagger.WolfiContainerOpts{
			Packages: []string{"gh"},
		}).
		WithSecretVariable("GITHUB_TOKEN", githubToken).
		WithExec([]string{"gh", "release", "-R=frantjc/forge", "create", tag, "--generate-notes", "--draft"})

	g0 := dag.Go(dagger.GoOpts{
		Module: gitRef.Tree(),
	})

	for _, goos := range []string{"darwin", "linux"} {
		for _, goarch := range []string{"amd64", "arm64"} {
			file := fmt.Sprintf("forge_%s_%s_%s", tag, goos, goarch)

			release = release.
				WithFile(
					file,
					g0.Build(dagger.GoBuildOpts{
						Pkg:     "./cmd/forge",
						Ldflags: "-s -w -X main.version=" + tag,
						Goos:    goos,
						Goarch:  goarch,
					}),
				).
				WithExec([]string{
					"gh", "release", "-R=frantjc/forge", "upload", tag, file,
				})
		}
	}

	_, err = release.
		WithExec([]string{"gh", "release", "-R=frantjc/forge", "edit", tag, "--latest", "--draft=false"}).
		Sync(ctx)
	if err != nil {
		return err
	}

	return nil
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
	return dag.Go(dagger.GoOpts{
		Module: m.Source,
	}).
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
