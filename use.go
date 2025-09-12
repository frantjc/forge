// A generated module for Forge functions

package main

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/frantjc/forge/githubactions"
	"github.com/frantjc/forge/internal/dagger"
	"github.com/frantjc/forge/internal/envconv"
	xos "github.com/frantjc/x/os"
	"golang.org/x/mod/modfile"
	"sigs.k8s.io/yaml"
)

const (
	// DefaultNode12ImageReference is the default container image reference that
	// is used by GitHub Actions that run using node12.
	DefaultNode12ImageReference = "docker.io/library/node:12"
	// DefaultNode16ImageReference is the default container image reference that
	// is used by GitHub Actions that run using node16.
	DefaultNode16ImageReference = "docker.io/library/node:16"
	// DefaultNode20ImageReference is the default container image reference that
	// is used by GitHub Actions that run using node20.
	DefaultNode20ImageReference = "docker.io/library/node:20"
	// DefaultNode24ImageReference is the default container image reference that
	// is used by GitHub Actions that run using node24.
	DefaultNode24ImageReference = "docker.io/library/node:24"
	// DefaultActionsCacheImageReference is the default container image reference
	// that is used as the actions cache service for GitHub Actions.
	DefaultActionsCacheImageReference = "ghcr.io/falcondev-oss/github-actions-cache-server:8"
)

var (
	// Node12ImageReference is the container image reference that
	// is used by GitHub Actions that run using node12.
	Node12ImageReference = DefaultNode12ImageReference
	// Node16ImageReference is the container image reference that
	// is used by GitHub Actions that run using node16.
	Node16ImageReference = DefaultNode16ImageReference
	// Node20ImageReference is the container image reference that
	// is used by GitHub Actions that run using node20.
	Node20ImageReference = DefaultNode20ImageReference
	// Node24ImageReference is the container image reference that
	// is used by GitHub Actions that run using node24.
	Node24ImageReference = DefaultNode24ImageReference
	// ActionsCacheImageReference is the container image reference that
	// is used as the actions cache service for GitHub Actions.
	ActionsCacheImageReference = DefaultActionsCacheImageReference
)

// Forge is the struct that methods are defined on for forge's Dagger module.
type Forge struct{}

const (
	actionPath    = "/forge/github/action"
	workspacePath = "/forge/github/workspace"
	tmpPath       = "/forge/github/runner/tmp"
	toolcachePath = "/forge/github/runner/toolcache"
	envPath       = "/forge/github/env"
	pathPath      = "/forge/github/path"
	outputPath    = "/forge/github/output"
	statePath     = "/forge/github/state"
	shimPath      = "/forge/shim"
	homePath      = "/forge/home"
)

// PreAction has a container that's prepared to execute an action and the subpath to that
// action, but has not yet executed the pre-step.
type PreAction struct {
	Action
}

// Action has a container that's prepared to execute an action and the subpath to that
// action, but has not yet executed the main step.
type Action struct {
	PostAction
}

// PostAction has a container that's prepared to execute an action and the subpath to that
// action, but has not yet executed the post-step.
type PostAction struct {
	FinalizedAction
	Subpath string
}

// FinalizedAction has a container that's prepared to execute an action and has executed that action.
type FinalizedAction struct {
	Ctr *dagger.Container
}

// Use creates a container to execute a GitHub Action in.
func (a *Forge) Use(
	ctx context.Context,
	action string,
	// +defaultPath="."
	repo *dagger.Directory,
	// +defaultPath="."
	workspace *dagger.Directory,
	// +optional
	with []string,
	// +optional
	env []string,
	// +optional
	debug bool,
	// +optional
	token *dagger.Secret,
) (*PreAction, error) {
	uses, err := githubactions.Parse(action)
	if err != nil {
		return nil, err
	}

	actn := workspace
	subpath := uses.Path

	if uses.IsRemote() {
		actn = dag.Git(githubactions.GetServerURL().JoinPath(uses.GetOwner(), uses.GetRepository()).String()).Ref(uses.Version).Tree()
		subpath = uses.GetActionPath()
	}

	metadata, err := actionMetadata(ctx, withAction(dag.Container(), actn), subpath)
	if err != nil {
		return nil, err
	}

	container := dag.Container()

	switch metadata.Runs.Using {
	case githubactions.RunsUsingDocker:
		if metadata.IsDockerfile() {
			dockerfile := metadata.Runs.Image

			if dockerfile == "" {
				dockerfile = "Dockerfile"
			}

			container = actn.DockerBuild(dagger.DirectoryDockerBuildOpts{Dockerfile: dockerfile})
		} else {
			container = container.From(strings.TrimPrefix(metadata.Runs.Image, githubactions.RunsUsingDockerImagePrefix))
		}
	case githubactions.RunsUsingNode12:
		container = withAction(container.From(Node12ImageReference), actn)
	case githubactions.RunsUsingNode16:
		container = withAction(container.From(Node16ImageReference), actn)
	case githubactions.RunsUsingNode20:
		container = withAction(container.From(Node20ImageReference), actn)
	case githubactions.RunsUsingNode24:
		container = withAction(container.From(Node24ImageReference), actn)
	default:
		return nil, fmt.Errorf("actions that run using %s are not supported", metadata.Runs.Using)
	}

	container = withActionsCache(container)
	container = withGitHubEnvVarsFromRef(ctx, container, repo.AsGit().Head())
	container = withDefaultGitHubEnvVars(container)

	ekv, err := parseKeyValuePairs(env)
	if err != nil {
		return nil, err
	}

	for k, v := range ekv {
		container = container.WithEnvVariable(k, v)
	}

	wkv, err := parseKeyValuePairs(with)
	if err != nil {
		return nil, err
	}

	container, err = withWith(container, metadata, wkv)
	if err != nil {
		return nil, err
	}

	container = withGitHubEnv(container)
	container = withGitHubPath(container)
	container = withGitHubOutput(container)
	container = withGitHubState(container)
	container = withRunnerTmp(container)
	container = withRunnerToolcache(container)
	container = withHome(container)
	container = withToken(ctx, container, token)
	container = withWorkspace(container, workspace)

	if debug {
		container = withDebug(container)
	}

	return &PreAction{
		Action: Action{
			PostAction: PostAction{
				FinalizedAction: FinalizedAction{
					Ctr: container,
				},
				Subpath: subpath,
			},
		},
	}, nil
}

// Pre executes the pre-step of the GitHub Action in the underlying container.
func (a *PreAction) Pre(ctx context.Context) (*Action, error) {
	metadata, err := actionMetadata(ctx, a.Container(), a.Subpath)
	if err != nil {
		return nil, err
	}

	a.Ctr, err = withShim(ctx, a.Container())
	if err != nil {
		return nil, err
	}

	switch {
	case metadata.IsDocker():
		if metadata.Runs.PreEntrypoint != "" {
			a.Ctr = a.Container().WithExec(append([]string{shimPath, metadata.Runs.PreEntrypoint}, metadata.Runs.Args...))

			a.Ctr, err = withExportedEnv(ctx, a.Container())
			if err != nil {
				return nil, err
			}
		}
	case metadata.IsNode():
		if metadata.Runs.Pre != "" {
			a.Ctr = a.Container().WithExec(append([]string{shimPath, "node", path.Join(actionPath, a.Subpath, metadata.Runs.Pre)}, metadata.Runs.Args...))

			a.Ctr, err = withExportedEnv(ctx, a.Container())
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("actions that run using %s are not supported", metadata.Runs.Using)
	}

	return &a.Action, nil
}

// Main executes the main step of the GitHub Action in the underlying container.
func (a *Action) Main(ctx context.Context) (*PostAction, error) {
	metadata, err := actionMetadata(ctx, a.Container(), a.Subpath)
	if err != nil {
		return nil, err
	}

	switch {
	case metadata.IsDocker():
		if metadata.Runs.Entrypoint != "" {
			a.Ctr = a.Container().WithExec(append([]string{shimPath, metadata.Runs.Entrypoint}, metadata.Runs.Args...))
		} else {
			a.Ctr = a.Container().WithExec(metadata.Runs.Args)
		}
	case metadata.IsNode():
		a.Ctr = a.Container().WithExec(append([]string{shimPath, "node", path.Join(actionPath, a.Subpath, metadata.Runs.Main)}, metadata.Runs.Args...))
	default:
		return nil, fmt.Errorf("actions that run using %s are not supported", metadata.Runs.Using)
	}

	a.Ctr, err = withExportedEnv(ctx, a.Container())
	if err != nil {
		return nil, err
	}

	return &a.PostAction, nil
}

// Main executes the pre- and main steps of the GitHub Action in the underlying container.
func (a *PreAction) Main(ctx context.Context) (*PostAction, error) {
	main, err := a.Pre(ctx)
	if err != nil {
		return nil, err
	}

	return main.Main(ctx)
}

// Post executes the post-step of the GitHub Action in the underlying container.
func (a *PostAction) Post(ctx context.Context) (*dagger.Container, error) {
	metadata, err := actionMetadata(ctx, a.Container(), a.Subpath)
	if err != nil {
		return nil, err
	}

	switch {
	case metadata.IsDocker():
		if metadata.Runs.PostEntrypoint != "" {
			a.Ctr = a.Container().WithExec(append([]string{shimPath, metadata.Runs.PostEntrypoint}, metadata.Runs.Args...))

			a.Ctr, err = withExportedEnv(ctx, a.Container())
			if err != nil {
				return nil, err
			}
		}
	case metadata.IsNode():
		if metadata.Runs.Post != "" {
			a.Ctr = a.Container().WithExec(append([]string{shimPath, "node", path.Join(actionPath, a.Subpath, metadata.Runs.Post)}, metadata.Runs.Args...))

			a.Ctr, err = withExportedEnv(ctx, a.Container())
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("actions that run using %s are not supported", metadata.Runs.Using)
	}

	return a.Container(), nil
}

// Post executes the pre-, main and post-steps of the GitHub Action in the underlying container.
func (a *PreAction) Post(ctx context.Context) (*dagger.Container, error) {
	main, err := a.Pre(ctx)
	if err != nil {
		return nil, err
	}

	return main.Post(ctx)
}

// Post executes the main and post-steps of the GitHub Action in the underlying container.
func (a *Action) Post(ctx context.Context) (*dagger.Container, error) {
	post, err := a.Main(ctx)
	if err != nil {
		return nil, err
	}

	return post.Post(ctx)
}

// Container gives access to the underlying container.
func (a *FinalizedAction) Container() *dagger.Container {
	return a.Ctr
}

// Terminal is a convenient alias for Container().Terminal().
func (a *FinalizedAction) Terminal() *dagger.Container {
	return a.Container().Terminal()
}

// Stdout is a convenient alias for Container().Stdout().
func (a *FinalizedAction) Stdout(ctx context.Context) (string, error) {
	return a.Container().Stdout(ctx)
}

// Stderr is a convenient alias for Container().Stderr().
func (a *FinalizedAction) Stderr(ctx context.Context) (string, error) {
	return a.Container().Stderr(ctx)
}

// CombinedOutput is a convenient alias for Container().CombinedOutput().
func (a *FinalizedAction) CombinedOutput(ctx context.Context) (string, error) {
	return a.Container().CombinedOutput(ctx)
}

// Workspace returns the current state of the GITHUB_WORKSPACE directory.
func (a *FinalizedAction) Workspace() *dagger.Directory {
	return a.Container().Directory(workspacePath)
}

// Toolcache returns the current state of the RUNNER_TOOLCACHE directory.
func (a *FinalizedAction) Toolcache() *dagger.Directory {
	return a.Container().Directory(toolcachePath)
}

// Action returns the current state of the GITHUB_ACTION_PATH directory.
func (a *FinalizedAction) Action() *dagger.Directory {
	return a.Container().Directory(actionPath)
}

// Home returns the current state of the HOME directory.
func (a *FinalizedAction) Home() *dagger.Directory {
	return a.Container().Directory(homePath)
}

// Env returns the parsed key-value pairs that were saved to GITHUB_ENV.
func (a *FinalizedAction) Env(ctx context.Context) (string, error) {
	env, err := gitHubEnv(ctx, a.Container())
	if err != nil {
		return "", err
	}

	return strings.Join(envconv.MapToArr(env), "\n"), nil
}

// State returns the parsed key-value pairs that were saved to GITHUB_STATE.
func (a *FinalizedAction) State(ctx context.Context) (string, error) {
	env, err := gitHubState(ctx, a.Container())
	if err != nil {
		return "", err
	}

	return strings.Join(envconv.MapToArr(env), "\n"), nil
}

// Output returns the parsed key-value pairs that were saved to GITHUB_OUTPUT.
func (a *FinalizedAction) Output(ctx context.Context) (string, error) {
	env, err := gitHubOutput(ctx, a.Container())
	if err != nil {
		return "", err
	}

	return strings.Join(envconv.MapToArr(env), "\n"), nil
}

func withShim(ctx context.Context, container *dagger.Container) (*dagger.Container, error) {
	src := dag.CurrentModule().Source()

	contents, err := src.File("go.mod").Contents(ctx)
	if err != nil {
		return nil, err
	}

	gomod, err := modfile.Parse("go.mod", []byte(contents), nil)
	if err != nil {
		return nil, err
	}

	golang := dag.Container().
		From(fmt.Sprintf("docker.io/library/golang:%s", gomod.Go.Version))

	gopath, err := golang.EnvVariable(ctx, "GOPATH")
	if err != nil {
		return nil, err
	}

	workdir := path.Join(gopath, "src", gomod.Module.Mod.Path)

	shim := golang.
		WithDirectory(workdir, src, dagger.ContainerWithDirectoryOpts{
			Include: []string{
				"go.mod",
				"go.sum",
				"cmd/shim/**",
				"command/shim.go",
				"githubactions/**",
				"internal/envconv/**",
				"internal/yaml/**",
				"internal/rangemap/**",
			},
		}).
		WithWorkdir(workdir).
		WithExec([]string{"go", "build", "-o", shimPath, "./cmd/shim"}).
		File(shimPath)

	return container.WithFile(shimPath, shim), nil
}

func actionMetadata(ctx context.Context, container *dagger.Container, subpath string) (*githubactions.Metadata, error) {
	var (
		errs     error
		metadata = &githubactions.Metadata{}
	)

	dir := container.Directory(path.Join(actionPath, subpath))

	for _, actionYAMLFileName := range githubactions.ActionYAMLFilenames {
		contents, err := dir.File(actionYAMLFileName).Contents(ctx)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		if err := yaml.Unmarshal([]byte(contents), metadata); err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		return metadata, nil
	}

	return nil, fmt.Errorf("find action metadata: %w", errs)
}

func withDebug(container *dagger.Container) *dagger.Container {
	return container.
		WithEnvVariable(githubactions.SecretActionsRunnerDebug, fmt.Sprint(true)).
		WithEnvVariable(githubactions.SecretActionsStepDebug, fmt.Sprint(true)).
		WithEnvVariable(githubactions.SecretRunnerDebug, fmt.Sprint(1))
}

func withToken(ctx context.Context, container *dagger.Container, secret *dagger.Secret) *dagger.Container {
	if secret != nil {
		token, err := secret.Plaintext(ctx)
		if err == nil {
			return container.
				WithEnvVariable(githubactions.EnvVarToken, token)
		}
	}

	return container
}

func withHome(container *dagger.Container) *dagger.Container {
	return container.
		WithEnvVariable("HOME", homePath).
		WithMountedCache(homePath, dag.CacheVolume("home"))
}

func withAction(container *dagger.Container, action *dagger.Directory) *dagger.Container {
	return container.
		WithEnvVariable(githubactions.EnvVarActionPath, actionPath).
		WithDirectory(actionPath, action)
}

func withEmptyFile(container *dagger.Container, fullpath string) *dagger.Container {
	return withFile(container, fullpath, "")
}

func withFile(container *dagger.Container, fullpath, contents string, opts ...dagger.ContainerWithFileOpts) *dagger.Container {
	return container.WithFile(path.Dir(fullpath), dag.File(path.Base(fullpath), contents), opts...)
}

func withGitHubEnvVarsFromRef(ctx context.Context, container *dagger.Container, gitRef *dagger.GitRef) *dagger.Container {
	ref, err := gitRef.Ref(ctx)
	if err != nil {
		return container
	}

	container = container.WithEnvVariable(githubactions.EnvVarRef, ref)

	commit, err := gitRef.Commit(ctx)
	if err != nil {
		return container
	}

	return container.WithEnvVariable(githubactions.EnvVarSha, commit)
}

func withDefaultGitHubEnvVars(container *dagger.Container) *dagger.Container {
	return container.
		WithEnvVariable(githubactions.EnvVarActions, fmt.Sprint(true)).
		WithEnvVariable(githubactions.EnvVarCI, fmt.Sprint(true)).
		WithEnvVariable(githubactions.EnvVarServerURL, githubactions.GetServerURL().String()).
		WithEnvVariable(githubactions.EnvVarAPIURL, githubactions.GetAPIURL().String()).
		WithEnvVariable(githubactions.EnvVarGraphQLURL, githubactions.GetGraphQLURL().String()).
		WithEnvVariable(githubactions.EnvVarRunnerOS, githubactions.OSLinux).
		WithEnvVariable(githubactions.EnvVarRunnerArch, githubactions.RunnerArch()).
		WithEnvVariable(githubactions.EnvVarRunNumber, fmt.Sprint(1)).
		WithEnvVariable(githubactions.EnvVarRunnerName, "forge")
}

func withActionsCache(container *dagger.Container) *dagger.Container {
	var (
		actionsResultsURL     = "http://actions-cache:3000"
		storageFilesystemPath = "/data"
	)

	return container.
		WithServiceBinding(
			"actions-cache",
			dag.Container().
				From(ActionsCacheImageReference).
				WithEnvVariable("API_BASE_URL", actionsResultsURL).
				WithEnvVariable("STORAGE_FILESYSTEM_PATH", storageFilesystemPath).
				WithEnvVariable("DEBUG", fmt.Sprint(true)).
				WithMountedCache(storageFilesystemPath, dag.CacheVolume("actions-cache")).
				WithExposedPort(3000).
				AsService(),
		).
		WithEnvVariable(githubactions.EnvVarActionsResultsURL, fmt.Sprintf("%s/", actionsResultsURL)).
		WithEnvVariable(githubactions.EnvVarActionsCacheURL, fmt.Sprintf("%s/", actionsResultsURL)).
		WithEnvVariable(githubactions.EnvVarActionsCacheServiceV2, fmt.Sprint(true)).
		WithEnvVariable(githubactions.EnvVarActionsRuntimeToken, "fake")
}

func withGitHubState(container *dagger.Container) *dagger.Container {
	return withEmptyFile(
		container.
			WithEnvVariable(githubactions.EnvVarState, statePath),
		statePath,
	)
}

func withGitHubOutput(container *dagger.Container) *dagger.Container {
	return withEmptyFile(
		container.
			WithEnvVariable(githubactions.EnvVarOutput, outputPath),
		outputPath,
	)
}

func withGitHubPath(container *dagger.Container) *dagger.Container {
	return withEmptyFile(
		container.
			WithEnvVariable(githubactions.EnvVarPath, pathPath),
		pathPath,
	)
}

func withGitHubEnv(container *dagger.Container) *dagger.Container {
	return withEmptyFile(
		container.
			WithEnvVariable(githubactions.EnvVarEnv, envPath),
		envPath,
	)
}

func withRunnerTmp(container *dagger.Container) *dagger.Container {
	return container.
		WithEnvVariable(githubactions.EnvVarRunnerTemp, tmpPath).
		WithMountedTemp(tmpPath)
}

func withRunnerToolcache(container *dagger.Container) *dagger.Container {
	return container.
		WithEnvVariable(githubactions.EnvVarRunnerToolCache, toolcachePath).
		WithMountedCache(toolcachePath, dag.CacheVolume("runner-toolcache"))
}

func withWorkspace(container *dagger.Container, workspace *dagger.Directory) *dagger.Container {
	return container.
		WithWorkdir(workspacePath).
		WithEnvVariable(githubactions.EnvVarWorkspace, workspacePath).
		WithMountedDirectory(workspacePath, workspace)
}

func gitHubEnv(ctx context.Context, container *dagger.Container) (map[string]string, error) {
	contents, err := container.File(envPath).Contents(ctx)
	if err != nil {
		return nil, err
	}

	env, err := githubactions.ParseEnvFile(strings.NewReader(contents))
	if err != nil {
		return nil, err
	}

	return env, nil
}

func withExportedGitHubEnv(ctx context.Context, container *dagger.Container) (*dagger.Container, error) {
	env, err := gitHubEnv(ctx, container)
	if err != nil {
		return nil, err
	}

	for k, v := range env {
		container = container.WithEnvVariable(k, v)
	}

	return container, nil
}

func withExportedGitHubPath(ctx context.Context, container *dagger.Container) (*dagger.Container, error) {
	contents, err := container.File(pathPath).Contents(ctx)
	if err != nil {
		return nil, err
	}

	newPath, err := githubactions.ParsePathFile(strings.NewReader(contents))
	if err != nil {
		return nil, err
	}

	oldPath, err := container.EnvVariable(ctx, "PATH")
	if err != nil {
		return nil, err
	}

	return container.WithEnvVariable("PATH", xos.JoinPath(newPath, oldPath)), nil
}

func gitHubState(ctx context.Context, container *dagger.Container) (map[string]string, error) {
	contents, err := container.File(statePath).Contents(ctx)
	if err != nil {
		return nil, err
	}

	state, err := githubactions.ParseEnvFile(strings.NewReader(contents))
	if err != nil {
		return nil, err
	}

	return state, nil
}

func withExportedGitHubState(ctx context.Context, container *dagger.Container) (*dagger.Container, error) {
	state, err := gitHubState(ctx, container)
	if err != nil {
		return nil, err
	}

	for k, v := range state {
		container = container.WithEnvVariable("STATE_"+k, v)
	}

	return container, nil
}

func gitHubOutput(ctx context.Context, container *dagger.Container) (map[string]string, error) {
	contents, err := container.File(outputPath).Contents(ctx)
	if err != nil {
		return nil, err
	}

	state, err := githubactions.ParseEnvFile(strings.NewReader(contents))
	if err != nil {
		return nil, err
	}

	return state, nil
}

func withExportedEnv(ctx context.Context, container *dagger.Container) (*dagger.Container, error) {
	container, err := withExportedGitHubEnv(ctx, container)
	if err != nil {
		return nil, err
	}

	container, err = withExportedGitHubPath(ctx, container)
	if err != nil {
		return nil, err
	}

	container, err = withExportedGitHubState(ctx, container)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func parseKeyValuePairs(args []string) (map[string]string, error) {
	with := map[string]string{}

	for _, a := range args {
		k, v, ok := strings.Cut(a, "=")
		if !ok {
			return nil, fmt.Errorf("malformed with: %s", a)
		}

		with[k] = v
	}

	return with, nil
}

func withWith(container *dagger.Container, metadata *githubactions.Metadata, with map[string]string) (*dagger.Container, error) {
	with, err := metadata.InputsFromWith(with)
	if err != nil {
		return nil, err
	}

	for k, v := range with {
		container = container.WithEnvVariable("INPUT_"+strings.ToUpper(k), v)
	}

	return container, nil
}
