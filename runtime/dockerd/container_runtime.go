package dockerd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/frantjc/forge"
	xslice "github.com/frantjc/x/slice"
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

	excludes, err := build.ReadDockerignore(dir)
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

	if bc, err := build.Compress(buildCtx); err == nil {
		buildCtx = bc
	}

	ibr, err := d.ImageBuild(ctx, buildCtx, types.ImageBuildOptions{
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
	ir, err := d.ImageInspect(ctx, image.Name())
	if err != nil {
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

		// Also because this is the Docker runtime...
		// If we're on linux, mount the Docker CLI into the container since then executables
		// on the host can also be used by the container because they have a common OS.
		if runtime.GOOS == "linux" {
			docker, err := exec.LookPath("docker")
			if errors.Is(err, exec.ErrDot) {
				docker, err = filepath.Abs(docker)
			}

			if err == nil {
				var (
					bin       = filepath.Join(d.DockerInDockerPath, "bin")
					addedPath = false
				)

				hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
					Source: docker,
					Target: filepath.Join(bin, "docker"),
					Type:   mount.TypeBind,
				})

				// If there already is a PATH, add to it.
				for i, e := range containerConfig.Env {
					if strings.HasPrefix(e, "PATH=") {
						containerConfig.Env[i] = fmt.Sprintf("%s:%s", e, bin)
						addedPath = true
						break
					}
				}

				// findPATHAppendBinFn iterates an env array searching for PATH.
				// If it finds it, it appends it to containerConfig.Env with bin
				// appended to the end and marks addedPath as true so we know to
				// stop this absurd PATH hunt.
				//
				// Note that we append to the very end so that we don't override
				// a Docker CLI that already exists on the PATH and cause unexpected
				// behavior with our arguably over the top helpfulness here.
				findPATHAppendBinFn := func(env []string) {
					for _, e := range env {
						if strings.HasPrefix(e, "PATH=") {
							containerConfig.Env = append(containerConfig.Env, fmt.Sprintf("%s:%s", e, bin))
							addedPath = true
							break
						}
					}
				}

				// If we didn't find the PATH on the containerConfig to modify,
				// then we want to modify it as it appears on the image config.
				//
				// Unfortunately, image.Config() is unbearably slow, so we use this
				// workaround to try and get the image config from elsewhere first.
				if !addedPath {
					// We may already have successfully called ImageInspectWithRaw
					// previously to check if we needed to call ImageLoad. If we didn't,
					// retry and it should work now that we've definitely loaded the image.
					if ir.Config == nil || len(ir.Config.Env) == 0 {
						ir, _ = d.ImageInspect(ctx, image.Name())
					}

					if ir.Config != nil {
						findPATHAppendBinFn(ir.Config.Env)
					}

					if !addedPath {
						if imageConfig, err := image.Config(); err == nil {
							findPATHAppendBinFn(imageConfig.Env)
						}
					}
				}

				// If we still didn't find the PATH on the imageConfig to modify,
				// we just add PATH ourselves.
				if !addedPath {
					containerConfig.Env = append(containerConfig.Env, fmt.Sprintf("PATH=%s", bin))
				}
			}
		}
	}

	hostConfig.Mounts = append(hostConfig.Mounts, xslice.Map(
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

func (d *ContainerRuntime) GetContainer(ctx context.Context, id string) (forge.Container, error) {
	if _, err := d.ContainerInspect(ctx, id); err != nil {
		return nil, err
	}

	return &Container{id, d.Client}, nil
}
