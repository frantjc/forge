package ore

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/forgeactions"
	"github.com/frantjc/forge/githubactions"
	"github.com/frantjc/forge/internal/containerutil"
	"github.com/frantjc/forge/internal/contaminate"
	xos "github.com/frantjc/x/os"
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

	o.GlobalContext = forgeactions.ConfigureGlobalContext(o.GlobalContext)
	o.GlobalContext.StepsContext[o.ID] = githubactions.StepContext{Outputs: make(map[string]string)}
	defer func() {
		ctx = githubactions.WithGlobalContext(ctx, o.GlobalContext)
	}()

	containerConfigs, err := forgeactions.ActionToConfigs(o.GlobalContext, uses, o.With, o.Env, actionMetadata, image)
	if err != nil {
		return err
	}

	workflowCommandStreams := forgeactions.NewWorkflowCommandStreams(o.GlobalContext, o.ID, drains)
	for _, containerConfig := range containerConfigs {
		cc := containerConfig
		cc.Mounts = contaminate.OverrideWithMountsFrom(ctx, containerConfig.Mounts...)
		cc.Env = append(cc.Env, o.GlobalContext.Env()...)

		container, err := containerutil.CreateSleepingContainer(ctx, containerRuntime, image, &cc)
		if err != nil {
			return err
		}

		if exitCode, err := container.Exec(ctx, &cc, workflowCommandStreams); err != nil {
			return err
		} else if exitCode > 0 {
			return xos.NewExitCodeError(ErrContainerExitedWithNonzeroExitCode, exitCode)
		}

		if err = container.Stop(ctx); err != nil {
			return err
		}

		if err = forgeactions.SetGlobalContextFromEnvFiles(ctx, o.GlobalContext, o.ID, container); err != nil {
			return err
		}

		if err = container.Remove(ctx); err != nil {
			return err
		}
	}

	return nil
}
