package docker

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/frantjc/forge"
	"github.com/frantjc/go-js"
)

var (
	errNonJSONResponse = errors.New("non-JSON response from image load")
)

func (d *ContainerRuntime) CreateContainer(ctx context.Context, image forge.Image, config *forge.ContainerConfig) (forge.Container, error) {
	ilr, err := d.ImageLoad(ctx, image.Blob(), true)
	if err != nil {
		return nil, err
	}

	if !ilr.JSON {
		return nil, errNonJSONResponse
	}

	// Body typically contains the response
	//
	//  '{"stream":"Loaded image[ ID]: <image reference>\n"}'
	//
	// Here, we extract it so we can tell the Docker daemon to create
	// a new container from that image reference
	m := map[string]string{}
	if err = json.NewDecoder(ilr.Body).Decode(&m); err != nil {
		return nil, errNonJSONResponse
	}

	reference := strings.TrimSpace(
		strings.TrimPrefix(
			strings.TrimPrefix(m["stream"], "Loaded image:"), "Loaded image ID:",
		),
	)

	defer d.Client.ImageRemove(ctx, reference, types.ImageRemoveOptions{Force: true}) //nolint:errcheck

	if err = ilr.Body.Close(); err != nil {
		return nil, err
	}

	var (
		addr            = d.Client.DaemonHost()
		containerConfig = &container.Config{
			User:        config.User,
			Env:         config.Env,
			Cmd:         config.Cmd,
			WorkingDir:  config.WorkingDir,
			Entrypoint:  config.Entrypoint,
			Image:       reference,
			AttachStdin: true,
			OpenStdin:   true,
			StdinOnce:   true,
			// AttachStdout: true,
			// AttachStderr: true,
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
		containerConfig.Env = append(containerConfig.Env, "DOCKER_HOST="+addr)
	}

	if runtime.GOOS == "linux" {
		if docker, err := exec.LookPath("docker"); err == nil {
			hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
				Source: docker,
				Target: "/usr/bin/docker",
				Type:   "bind",
			})
		}
	}

	hostConfig.Mounts = append(hostConfig.Mounts, js.Map(
		config.Mounts,
		func(m *forge.ContainerConfig_Mount, _ int, _ []*forge.Mount) mount.Mount {
			var (
				mountType = mount.TypeVolume
			)
			switch {
			case m.GetSource() == "":
				mountType = mount.TypeTmpfs
			case filepath.IsAbs(m.GetSource()):
				mountType = mount.TypeBind
			}

			return mount.Mount{
				Type:   mountType,
				Source: m.GetSource(),
				Target: m.GetDestination(),
			}
		},
	)...)

	cccb, err := d.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return nil, err
	}

	return &Container{cccb.ID, d.Client}, nil
}
