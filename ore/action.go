package ore

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/forgeactions"
	"github.com/frantjc/forge/githubactions"
	"github.com/frantjc/forge/internal/containerutil"
	"github.com/frantjc/forge/internal/contaminate"
	errorcode "github.com/frantjc/go-error-code"
)

// Action is an Ore representing a GitHub Action.
// That is--a step in a GitHub Actions Workflow that
// uses the `uses` key.
type Action struct {
	ID            string                       `json:"id,omitempty"`
	Uses          string                       `json:"uses,omitempty"`
	With          map[string]string            `json:"with,omitempty"`
	Env           map[string]string            `json:"env,omitempty"`
	GlobalContext *githubactions.GlobalContext `json:"global_context,omitempty"`
}

func (o *Action) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) error {
	_ = forge.LoggerFrom(ctx)

	uses, err := githubactions.Parse(o.Uses)
	if err != nil {
		return err
	}

	actionMetadata, err := forgeactions.GetUsesMetadata(ctx, uses)
	if err != nil {
		return err
	}

	image, err := forgeactions.GetImageForMetadata(ctx, containerRuntime, actionMetadata, uses)
	if err != nil {
		return err
	}

	if o.GlobalContext == nil {
		o.GlobalContext = githubactions.NewGlobalContextFromEnv()
	}
	defer func() {
		ctx = githubactions.WithGlobalContext(ctx, o.GlobalContext)
	}()

	containerConfigs, err := forgeactions.ActionToConfigs(o.GlobalContext, uses, o.With, o.Env, actionMetadata, image)
	if err != nil {
		return err
	}

	workflowCommandStreams := forgeactions.NewWorkflowCommandStreams(o.GlobalContext, o.ID, drains)
	for _, containerConfig := range containerConfigs {
		containerConfig.Mounts = contaminate.OverrideWithMountsFrom(ctx, containerConfig.Mounts...)
		container, err := containerutil.CreateSleepingContainer(ctx, containerRuntime, image, containerConfig)
		if err != nil {
			return err
		}
		defer container.Stop(ctx)   //nolint:errcheck
		defer container.Remove(ctx) //nolint:errcheck

		if exitCode, err := container.Exec(ctx, containerConfig, workflowCommandStreams); err != nil {
			return err
		} else if exitCode > 0 {
			return errorcode.New(ErrContainerExitedWithNonzeroExitCode, errorcode.WithExitCode(exitCode))
		}
	}

	return nil
}
