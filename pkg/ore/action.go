package ore

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/actions2container"
	"github.com/frantjc/forge/pkg/github/actions"
)

func (o *Action) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) (*forge.Lava, error) {
	uses, err := actions.Parse(o.Uses)
	if err != nil {
		return nil, err
	}

	volumes, err := actions2container.CreateVolumes(ctx, containerRuntime, uses)
	if err != nil {
		return nil, err
	}

	for _, volume := range volumes {
		defer volume.Remove(ctx)
	}

	container, err := actions2container.CreateContainerForUses(ctx, containerRuntime, uses)
	if err != nil {
		return nil, err
	}
	defer container.Remove(ctx)

	var (
		stdout = new(bytes.Buffer)
		stderr = new(bytes.Buffer)
	)
	exitCode, err := container.Run(ctx, &forge.Streams{
		Drains: &forge.Drains{
			Out: stdout,
			Err: stderr,
		},
	})
	if err != nil {
		return nil, err
	}

	actionMetadata := &actions.Metadata{}
	if err = json.NewDecoder(stdout).Decode(actionMetadata); err != nil {
		return nil, err
	}

	image, err := actions2container.PullImageForMetadata(ctx, containerRuntime, actionMetadata)
	if err != nil {
		return nil, err
	}

	if o.GetGlobalContext() == nil {
		o.GlobalContext = actions.NewGlobalContextFromEnv()
	}
	defer func() {
		ctx = actions.WithGlobalContext(ctx, o.GlobalContext)
	}()

	conatinerConfigs, err := actions2container.ActionToConfigs(o.GlobalContext, uses, o.With, o.Env, actionMetadata)
	if err != nil {
		return nil, err
	}

	workflowCommandStreams := actions2container.NewWorkflowCommandStreams(o.GlobalContext, o.GetId(), drains)
	for _, containerConfig := range conatinerConfigs {
		container, err = actions2container.CreateContainer(ctx, containerRuntime, image, containerConfig)
		if err != nil {
			break
		}
		defer container.Remove(ctx)

		exitCode, err = container.Run(ctx, workflowCommandStreams)
		if err != nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	return &forge.Lava{
		ExitCode: int64(exitCode),
	}, nil
}
