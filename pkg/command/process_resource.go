package command

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/frantjc/forge/pkg/concourse2container"
	"github.com/frantjc/forge/pkg/config"
	"github.com/frantjc/forge/pkg/ore"
	"github.com/frantjc/forge/pkg/runtime/container/docker"
)

func processResource(ctx context.Context, method, name string, params map[string]string) error {
	_ = forge.LoggerFrom(ctx)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	resources := &config.Resources{}
	for _, filename := range []string{"forge.json"} {
		filepath := path.Join(wd, filename)
		if file, err := os.Open(filepath); err == nil {
			if err = json.NewDecoder(file).Decode(resources); err == nil {
				break
			}
		}
	}
	if err != nil {
		return err
	}

	o := &ore.Resource{
		Method: method,
		Params: params,
	}
	for _, r := range resources.GetResources() {
		if r.GetName() == name {
			o.Resource = r
		}
	}
	if o.GetResource() == nil {
		return errors.New("resource not found: " + name)
	}

	for _, t := range resources.GetResourceTypes() {
		if t.GetName() == o.GetResource().GetType() {
			o.ResourceType = t
		}
	}
	if o.GetResourceType() == nil {
		return errors.New("resource type not found: " + o.GetResource().GetType())
	}

	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	_, err = forge.NewFoundry(docker.New(c), nil).Process(
		contaminate.WithMounts(ctx, &forge.Mount{
			Source:      wd,
			Destination: concourse2container.DefaultRootPath + "/" + o.GetResource().GetName(),
		}), o, forge.StdDrains(),
	)
	return err
}
