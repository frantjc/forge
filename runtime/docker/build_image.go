package docker

import (
	"context"
	"io"
	"strings"

	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/frantjc/forge"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
)

func (d *ContainerRuntime) BuildDockerfile(ctx context.Context, dir, reference string) (forge.Image, error) {
	ref, err := name.ParseReference(strings.NewReplacer(".", "").Replace(reference))
	if err != nil {
		return nil, err
	}

	excludes, err := build.ReadDockerignore(dir)
	if err != nil {
		return nil, err
	}

	buildCtx, err := archive.TarWithOptions(dir, &archive.TarOptions{
		ExcludePatterns: excludes,
		ChownOpts:       &idtools.Identity{UID: 0, GID: 0},
	})
	if err != nil {
		return nil, err
	}

	if bc, err := build.Compress(buildCtx); err == nil {
		buildCtx = bc
	}

	ibr, err := d.Client.ImageBuild(ctx, buildCtx, types.ImageBuildOptions{
		Tags:       []string{ref.Name()},
		PullParent: true,
		Remove:     true,
	})
	if err != nil {
		return nil, err
	}

	if err := jsonmessage.DisplayJSONMessagesStream(ibr.Body, io.Discard, 0, false, nil); err != nil {
		if jerr, ok := err.(*jsonmessage.JSONError); ok {
			return nil, jerr
		}

		return nil, err
	}

	if _, err = io.Copy(io.Discard, ibr.Body); err != nil {
		return nil, err
	}

	if err = ibr.Body.Close(); err != nil {
		return nil, err
	}

	img, err := daemon.Image(ref, daemon.WithClient(d), daemon.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return &Image{
		Image:     img,
		Reference: ref,
	}, nil
}
