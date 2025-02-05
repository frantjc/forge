package forge

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/frantjc/forge/envconv"
	"github.com/frantjc/forge/githubactions"
	"github.com/frantjc/forge/internal/hostfs"
	xos "github.com/frantjc/x/os"
	"github.com/opencontainers/go-digest"
)

// Action is an Ore representing a GitHub Action.
// That is--a step in a GitHub Actions workflow that
// uses the `uses` key.
type Action struct {
	ID            string
	Uses          string
	With          map[string]string
	Env           map[string]string
	GlobalContext *githubactions.GlobalContext
}

func (o *Action) Liquify(ctx context.Context, containerRuntime ContainerRuntime, opts ...OreOpt) error {
	opt := oreOptsWithDefaults(opts...)

	uses, err := githubactions.Parse(o.Uses)
	if err != nil {
		return err
	}

	actionMetadata, err := getUsesMetadata(ctx, uses)
	if err != nil {
		return err
	}

	image, err := pullImageForMetadata(ctx, containerRuntime, actionMetadata, uses)
	if err != nil {
		return err
	}

	o.GlobalContext = configureGlobalContext(o.GlobalContext, opt)
	o.GlobalContext.StepsContext[o.ID] = githubactions.StepContext{Outputs: make(map[string]string)}

	containerConfigs, err := actionToConfigs(o.GlobalContext, uses, o.With, o.Env, actionMetadata, image, opt)
	if err != nil {
		return err
	}

	workflowCommandStreams := newWorkflowCommandStreams(o.GlobalContext, o.ID, opt)
	for _, containerConfig := range containerConfigs {
		cc := containerConfig
		cc.Mounts = overrideMounts(cc.Mounts, opt.Mounts...)
		cc.Env = append(cc.Env, o.GlobalContext.Env()...)

		container, err := createSleepingContainer(ctx, containerRuntime, image, &cc, opt)
		if err != nil {
			return err
		}
		defer container.Stop(ctx)   //nolint:errcheck
		defer container.Remove(ctx) //nolint:errcheck

		if exitCode, err := container.Exec(ctx, &cc, workflowCommandStreams); err != nil {
			return err
		} else if exitCode > 0 {
			return xos.NewExitCodeError(ErrContainerExitedWithNonzeroExitCode, exitCode)
		}

		if err = setGlobalContextFromEnvFiles(ctx, o.GlobalContext, o.ID, container, opt); err != nil {
			return err
		}
	}

	return nil
}

func usesToRootDirectory(uses *githubactions.Uses) (string, error) {
	if uses.IsLocal() {
		return filepath.Abs(uses.Path)
	}

	return filepath.Join(hostfs.ActionsCache, uses.GetRepository(), uses.Version), nil
}

func usesToActionDirectory(uses *githubactions.Uses) (string, error) {
	if uses.IsLocal() {
		return usesToRootDirectory(uses)
	}

	return filepath.Join(hostfs.ActionsCache, uses.GetRepository(), uses.Version, uses.GetActionPath()), nil
}

func getUsesMetadata(ctx context.Context, uses *githubactions.Uses) (*githubactions.Metadata, error) {
	dir, err := usesToRootDirectory(uses)
	if err != nil {
		return nil, err
	}

	return githubactions.GetUsesMetadata(ctx, uses, dir)
}

func setGlobalContextFromEnvFiles(ctx context.Context, globalContext *githubactions.GlobalContext, step string, container Container, opt *OreOpts) error {
	var errs []error
	globalContext = configureGlobalContext(globalContext, opt)

	rc, err := container.CopyFrom(ctx, GitHubPath(opt.WorkingDir))
	if err != nil {
		return fmt.Errorf("copying GitHub path from container: %w", err)
	}
	defer rc.Close()

	r := tar.NewReader(rc)
	for {
		header, err := r.Next()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}

		//nolint:gocritic
		switch header.Typeflag {
		case tar.TypeReg:
			switch {
			case strings.HasSuffix(GitHubOutput(opt.WorkingDir), header.Name):
				outputs, err := githubactions.ParseEnvFile(r)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				if stepContext, ok := globalContext.StepsContext[step]; !ok || stepContext.Outputs == nil {
					globalContext.StepsContext[step] = githubactions.StepContext{
						Outputs: outputs,
					}
				} else {
					maps.Copy(globalContext.StepsContext[step].Outputs, outputs)
				}
			case strings.HasSuffix(GitHubState(opt.WorkingDir), header.Name):
				outputs, err := githubactions.ParseEnvFile(r)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				for k, v := range outputs {
					globalContext.EnvContext[fmt.Sprintf("STATE_%s", k)] = v
				}
			case strings.HasSuffix(GitHubEnv(opt.WorkingDir), header.Name):
				env, err := githubactions.ParseEnvFile(r)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				for k, v := range env {
					globalContext.EnvContext[k] = v
				}
			}
		}
	}

	return errors.Join(errs...)
}

func newWorkflowCommandStreams(globalContext *githubactions.GlobalContext, id string, opt *OreOpts) *Streams {
	globalContext = configureGlobalContext(globalContext, opt)
	debug := globalContext.DebugEnabled()

	return &Streams{
		In: opt.Streams.In,
		Out: &githubactions.WorkflowCommandWriter{
			GlobalContext:      globalContext,
			ID:                 id,
			StopCommandsTokens: map[string]bool{},
			Debug:              debug,
			Out:                opt.Streams.Out,
		},
		Err: &githubactions.WorkflowCommandWriter{
			GlobalContext:      globalContext,
			ID:                 id,
			StopCommandsTokens: map[string]bool{},
			Debug:              debug,
			Out:                opt.Streams.Err,
		},
		Tty:        opt.Streams.Tty,
		DetachKeys: opt.Streams.DetachKeys,
	}
}

func configureGlobalContext(globalContext *githubactions.GlobalContext, opt *OreOpts) *githubactions.GlobalContext {
	if globalContext == nil {
		globalContext = githubactions.NewGlobalContextFromEnv()
	}

	if globalContext.GitHubContext == nil {
		globalContext.GitHubContext = &githubactions.GitHubContext{}
	}

	globalContext.GitHubContext.Workspace = GitHubWorkspace(opt.WorkingDir)
	globalContext.GitHubContext.ActionPath = GitHubActionPath(opt.WorkingDir)

	if globalContext.RunnerContext == nil {
		globalContext.RunnerContext = &githubactions.RunnerContext{}
	}

	globalContext.RunnerContext.Temp = GitHubRunnerTmp(opt.WorkingDir)
	globalContext.RunnerContext.ToolCache = GitHubRunnerToolCache(opt.WorkingDir)

	return globalContext
}

func actionToConfigs(globalContext *githubactions.GlobalContext, uses *githubactions.Uses, with, environment map[string]string, actionMetadata *githubactions.Metadata, image Image, opt *OreOpts) ([]ContainerConfig, error) {
	containerConfigs := []ContainerConfig{}
	globalContext = configureGlobalContext(globalContext, opt)

	if actionMetadata != nil {
		if actionMetadata.Runs != nil {
			dir, err := usesToRootDirectory(uses)
			if err != nil {
				return nil, err
			}

			var (
				entrypoint = []string{}
				env        = append(envconv.MapToArr(environment), envconv.MapToArr(actionMetadata.Runs.Env)...)
				cmd        = actionMetadata.Runs.Args
				actionPath = filepath.Join(GitHubActionPath(opt.WorkingDir), uses.GetActionPath())
				mounts     = []Mount{
					{
						Source:      dir,
						Destination: GitHubActionPath(opt.WorkingDir),
					},
					{
						Destination: GitHubWorkspace(opt.WorkingDir),
					},
					{
						Destination: GitHubRunnerToolCache(opt.WorkingDir),
					},
					{
						Destination: GitHubRunnerTmp(opt.WorkingDir),
					},
				}
				entrypoints [][]string
			)

			switch actionMetadata.Runs.Using {
			case githubactions.RunsUsingNode12, githubactions.RunsUsingNode16, githubactions.RunsUsingNode20:
				entrypoint = []string{"node"}

				if pre := actionMetadata.Runs.Pre; pre != "" {
					entrypoints = append(entrypoints, []string{filepath.Join(actionPath, pre)})
				}

				if main := actionMetadata.Runs.Main; main != "" {
					entrypoints = append(entrypoints, []string{filepath.Join(actionPath, main)})
				}
			case githubactions.RunsUsingDocker:
				if pre := actionMetadata.Runs.PreEntrypoint; pre != "" {
					entrypoints = append(entrypoints, []string{pre})
				}

				if main := actionMetadata.Runs.Entrypoint; main != "" {
					entrypoints = append(entrypoints, []string{main})
				} else {
					config, err := image.Config()
					if err != nil {
						return nil, err
					}

					entrypoints = append(entrypoints, config.Entrypoint)
				}
			default:
				return nil, fmt.Errorf("unsupported runs using: %s", actionMetadata.Runs.Using)
			}

			unexpandedInputs, err := actionMetadata.InputsFromWith(with)
			if err != nil {
				return nil, err
			}

			var (
				inputs   = make(map[string]string, len(unexpandedInputs))
				expander = githubactions.ExpandFunc(globalContext.GetString)
			)
			for k, v := range unexpandedInputs {
				inputs[k] = expander.ExpandString(v)
			}

			globalContext.InputsContext = inputs
			env = append(env, globalContext.Env()...)
			env = append(env,
				fmt.Sprintf("%s=%s", githubactions.EnvVarPath, GitHubPath(opt.WorkingDir)),
				fmt.Sprintf("%s=%s", githubactions.EnvVarEnv, GitHubEnv(opt.WorkingDir)),
				fmt.Sprintf("%s=%s", githubactions.EnvVarOutput, GitHubOutput(opt.WorkingDir)),
				fmt.Sprintf("%s=%s", githubactions.EnvVarState, GitHubState(opt.WorkingDir)),
			)

			for _, ep := range entrypoints {
				if len(ep) > 0 {
					if len(entrypoint) > 0 {
						containerConfigs = append(containerConfigs, ContainerConfig{
							Entrypoint: entrypoint,
							Cmd:        append(ep, cmd...),
							Env:        env,
							Mounts:     mounts,
							WorkingDir: GitHubWorkspace(opt.WorkingDir),
						})
					} else {
						containerConfigs = append(containerConfigs, ContainerConfig{
							Entrypoint: ep,
							Cmd:        cmd,
							Env:        env,
							Mounts:     mounts,
							WorkingDir: GitHubWorkspace(opt.WorkingDir),
						})
					}
				}
			}
		}
	}

	return containerConfigs, nil
}

// GetImageForMetadata takes an action.yml and returns the OCI image that forge
// should run it inside of. If the action.yml runs using "dockerfile" and the
// forge.ContainerRuntime does not implement ImageBuilder, returns ErrCannotBuildDockerfile.
func pullImageForMetadata(ctx context.Context, containerRuntime ContainerRuntime, actionMetadata *githubactions.Metadata, uses *githubactions.Uses) (Image, error) {
	if actionMetadata.IsDockerfile() {
		dir, err := usesToActionDirectory(uses)
		if err != nil {
			return nil, err
		}

		var (
			dockerfilePath = filepath.Join(dir, actionMetadata.Runs.Image)
			reference      = fmt.Sprintf("ghcr.io/%s:%s", uses.GetRepository(), uses.Version)
		)
		if uses.IsLocal() {
			dockerfile, err := os.Open(dockerfilePath)
			if err != nil {
				return nil, err
			}
			defer dockerfile.Close()

			digest, err := digest.FromReader(dockerfile)
			if err != nil {
				return nil, err
			}

			reference = fmt.Sprintf(
				"%s:%s",
				filepath.Join(
					"forge.frantj.cc",
					regexp.MustCompile(`[^a-z0-9._/-]`).ReplaceAllString(
						strings.ToLower(dockerfilePath),
						"",
					),
				),
				digest.Hex(),
			)
		}

		if imageBuilder, ok := containerRuntime.(ImageBuilder); ok {
			return imageBuilder.BuildDockerfile(ctx, dockerfilePath, reference)
		}

		return nil, ErrCannotBuildDockerfile
	}

	return containerRuntime.PullImage(ctx, metadataToImageReference(actionMetadata))
}

// MetadataToImageReference takes an action.yaml and finds the reference
// to the OCI image that forge should run it inside of.
func metadataToImageReference(actionMetadata *githubactions.Metadata) string {
	if actionMetadata == nil {
		return ""
	}

	if actionMetadata.Runs == nil {
		return ""
	}

	switch actionMetadata.Runs.Using {
	case githubactions.RunsUsingNode12:
		return Node12ImageReference
	case githubactions.RunsUsingNode16:
		return Node16ImageReference
	case githubactions.RunsUsingNode20:
		return Node20ImageReference
	case githubactions.RunsUsingDocker:
		if !actionMetadata.IsDockerfile() {
			return strings.TrimPrefix(actionMetadata.Runs.Image, githubactions.RunsUsingDockerImagePrefix)
		}
	}

	return ""
}
