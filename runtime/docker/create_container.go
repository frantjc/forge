package docker

import (
	"context"
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/containerfs"
	xslice "github.com/frantjc/x/slice"
)

func (d *ContainerRuntime) CreateContainer(ctx context.Context, image forge.Image, config *forge.ContainerConfig) (forge.Container, error) {
	// If the Docker daemon already has the image,
	// don't bother loading it in again.
	ii, _, err := d.ImageInspectWithRaw(ctx, image.Name())
	if err != nil {
		if ilr, err := d.ImageLoad(ctx, image.Blob(), true); err != nil {
			return nil, err
		} else if err = ilr.Body.Close(); err != nil {
			return nil, err
		}
	}

	var (
		addr            = d.Client.DaemonHost()
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

	// Because this is the Docker runtime...
	// Mount the Docker daemon into the container for use by the process inside the container.
	if addr != "" {
		sock := filepath.Join(containerfs.WorkingDir, "/docker.sock")
		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Source: strings.TrimPrefix(addr, "unix://"),
			Target: sock,
			Type:   mount.TypeBind,
		})
		containerConfig.Env = append(containerConfig.Env, "DOCKER_HOST=unix://"+sock)
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
				bin       = filepath.Join(containerfs.WorkingDir, "/bin")
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
					containerConfig.Env[i] = e + ":" + bin
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
						containerConfig.Env = append(containerConfig.Env, e+":"+bin)
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
				if ii.Config == nil || len(ii.Config.Env) == 0 {
					ii, _, _ = d.ImageInspectWithRaw(ctx, image.Name())
				}

				if ii.Config != nil {
					findPATHAppendBinFn(ii.Config.Env)
				}

				if !addedPath && ii.ContainerConfig != nil {
					findPATHAppendBinFn(ii.ContainerConfig.Env)
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
				containerConfig.Env = append(containerConfig.Env, "PATH="+bin)
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
