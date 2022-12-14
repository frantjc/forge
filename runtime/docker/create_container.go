package docker

import (
	"context"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/frantjc/forge"
	"github.com/frantjc/go-fn"
)

func (d *ContainerRuntime) CreateContainer(ctx context.Context, image forge.Image, config *forge.ContainerConfig) (forge.Container, error) {
	// if the Docker daemon already has the image, don't bother loading it in
	if _, _, err := d.ImageInspectWithRaw(ctx, image.Name()); err != nil {
		ilr, err := d.ImageLoad(ctx, image.Blob(), true)
		if err != nil {
			return nil, err
		}

		if err = ilr.Body.Close(); err != nil {
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

	if addr != "" {
		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Source: strings.TrimPrefix(addr, "unix://"),
			Target: "/var/run/docker.sock",
			Type:   "bind",
		})
		containerConfig.Env = append(containerConfig.Env, "DOCKER_HOST=/var/run/docker.sock")
	}

	if goruntime.GOOS == "linux" {
		if docker, err := exec.LookPath("docker"); err == nil {
			hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
				Source: docker,
				Target: "/usr/bin/docker",
				Type:   "bind",
			})
		}
	}

	hostConfig.Mounts = append(hostConfig.Mounts, fn.Map(
		config.Mounts,
		func(m *forge.Mount, _ int) mount.Mount {
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
