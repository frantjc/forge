package ore

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/frantjc/forge/pkg/actions2container"
	"github.com/frantjc/forge/pkg/github/actions"
)

func (o *Action) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, basin forge.Basin, drains *forge.Drains) (*forge.Cast, error) {
	_ = forge.LoggerFrom(ctx)

	uses, err := actions.Parse(o.Uses)
	if err != nil {
		return nil, err
	}

	volumes, err := actions2container.CreateVolumes(ctx, containerRuntime, uses)
	if err != nil {
		return nil, err
	}

	for _, volume := range volumes {
		defer volume.Remove(ctx) //nolint:errcheck
	}

	containerConfig := actions2container.UsesToConfig(uses)
	containerConfig.Mounts = append(containerConfig.Mounts, contaminate.MountsFrom(ctx)...)

	container, err := containerRuntime.CreateContainer(ctx, actions2container.UsesImage, containerConfig)
	if err != nil {
		return nil, err
	}
	defer container.Stop(ctx)   //nolint:errcheck
	defer container.Remove(ctx) //nolint:errcheck

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
		containerConfig.Mounts = contaminate.OverrideWithMountsFrom(ctx, containerConfig.Mounts...)
		container, err := CreateSleepingContainer(ctx, containerRuntime, image, containerConfig)
		if err != nil {
			break
		}
		defer container.Stop(ctx)   //nolint:errcheck
		defer container.Remove(ctx) //nolint:errcheck

		exitCode, err = container.Exec(ctx, containerConfig, workflowCommandStreams)
		if err != nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	return &forge.Cast{
		ExitCode: int64(exitCode),
	}, nil
}
