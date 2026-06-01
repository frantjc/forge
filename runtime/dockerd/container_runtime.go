package dockerd

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	clibuild "github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/frantjc/forge"
	xslices "github.com/frantjc/x/slices"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/moby/go-archive"
)

func New(c *client.Client, dindPath string) *ContainerRuntime {
	return &ContainerRuntime{c, dindPath}
}

// ContainerRuntime implements github.com/frantjc/forge.ContainerRuntime.
type ContainerRuntime struct {
	// Client interacts with a Docker daemon.
	*client.Client
	// DockerInDockerPath signals whether or not to mount the docker.sock of the
	// *github.com/docker/docker/client.Client and configuration to direct
	// `docker` to it into each container that it runs.
	DockerInDockerPath string
}

func (d *ContainerRuntime) GoString() string {
	return "&ContainerRuntime{" + d.DaemonHost() + "}"
}

func (d *ContainerRuntime) PullImage(ctx context.Context, reference string) (forge.Image, error) {
	ref, err := name.ParseReference(reference)
	if err != nil {
		return nil, err
	}

	r, err := d.ImagePull(ctx, ref.Name(), image.PullOptions{})
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(io.Discard, r); err != nil {
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

func (d *ContainerRuntime) BuildDockerfile(ctx context.Context, dockerfile, reference string) (forge.Image, error) {
	ref, err := name.ParseReference(reference)
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(dockerfile)

	excludes, err := clibuild.ReadDockerignore(dir)
	if err != nil {
		return nil, err
	}

	buildCtx, err := archive.TarWithOptions(dir, &archive.TarOptions{
		ExcludePatterns: excludes,
		ChownOpts:       &archive.ChownOpts{UID: 0, GID: 0},
	})
	if err != nil {
		return nil, err
	}

	if bc, err := clibuild.Compress(buildCtx); err == nil {
		buildCtx = bc
	}

	ibr, err := d.ImageBuild(ctx, buildCtx, build.ImageBuildOptions{
		Tags:       []string{ref.Name()},
		Dockerfile: filepath.Base(dockerfile),
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

func (d *ContainerRuntime) CreateContainer(ctx context.Context, image forge.Image, config *forge.ContainerConfig) (forge.Container, error) {
	// If the Docker daemon already has the image,
	// don't bother loading it in again.
	if _, err := d.ImageInspect(ctx, image.Name()); err != nil {
		if ilr, err := d.ImageLoad(ctx, image.Blob(), client.ImageLoadWithQuiet(true)); err != nil {
			return nil, err
		} else if err = ilr.Body.Close(); err != nil {
			return nil, err
		}
	}

	var (
		addr            = d.DaemonHost()
		containerConfig = &container.Config{
			User:         config.User,
			Env:          config.Env,
			Cmd:          config.Cmd,
			WorkingDir:   config.WorkingDir,
			Entrypoint:   config.Entrypoint,
			Image:        image.Name(),
			AttachStdin:  true,
			OpenStdin:    true,
			StdinOnce:    true,
			AttachStdout: true,
			AttachStderr: true,
		}
		hostConfig = &container.HostConfig{
			Privileged: config.Privileged,
		}
	)

	if d.DockerInDockerPath != "" {
		// Because this is the Docker runtime...
		// Mount the Docker daemon into the container for use by the process inside the container.
		if strings.HasPrefix(addr, "unix://") {
			sock := filepath.Join(d.DockerInDockerPath, "/docker.sock")
			hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
				Source: strings.TrimPrefix(addr, "unix://"),
				Target: sock,
				Type:   mount.TypeBind,
			})
			containerConfig.Env = append(containerConfig.Env, fmt.Sprintf("%s=unix://%s", client.EnvOverrideHost, sock))
		}
	}

	hostConfig.Mounts = append(hostConfig.Mounts, xslices.Map(
		config.Mounts,
		func(m forge.Mount, _ int) mount.Mount {
			mountType := mount.TypeVolume
			switch {
			case m.Source == "":
				mountType = mount.TypeTmpfs
			case filepath.IsAbs(m.Source):
				mountType = mount.TypeBind
			}

			return mount.Mount{
				Type:   mountType,
				Source: m.Source,
				Target: m.Destination,
			}
		},
	)...)

	cccb, err := d.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return nil, err
	}

	return &Container{cccb.ID, d.Client}, nil
}
