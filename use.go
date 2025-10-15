// Run reusable steps from various proprietary CI systems.

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
	// +private
	Action
}

// Action has a container that's prepared to execute an action and the subpath to that
// action, but has not yet executed the main step.
type Action struct {
	// +private
	PostAction
}

// PostAction has a container that's prepared to execute an action and the subpath to that
// action, but has not yet executed the post-step.
type PostAction struct {
	// +private
	FinalizedAction
	// +private
	Ref *dagger.GitRef
	// +private
	Metadata string
	// +private
	Subpath string
	// +private
	Inputs []string
}

// FinalizedAction has a container that has fully executed its action.
type FinalizedAction struct {
	// +private
	Ctr *dagger.Container
}

// Use creates a container to execute a GitHub Action in.
func (m *Forge) Use(
	ctx context.Context,
	// The action to use
	action string,
	// The workspace to act on
	// +defaultPath="."
	workspace *dagger.Directory,
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

	rawMetadata, err := actionMetadata(ctx, actn, subpath)
	if err != nil {
		return nil, err
	}

	metadata, err := parseActionMetadata(rawMetadata)
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

			container = actn.Directory(subpath).DockerBuild(dagger.DirectoryDockerBuildOpts{Dockerfile: dockerfile})
		} else {
			container = container.From(strings.TrimPrefix(metadata.Runs.Image, githubactions.RunsUsingDockerImagePrefix))
		}
	case githubactions.RunsUsingNode12:
		container = container.From(Node12ImageReference)
	case githubactions.RunsUsingNode16:
		container = container.From(Node16ImageReference)
	case githubactions.RunsUsingNode20:
		container = container.From(Node20ImageReference)
	case githubactions.RunsUsingNode24:
		container = container.From(Node24ImageReference)
	default:
		return nil, fmt.Errorf("actions that run using %s are not supported", metadata.Runs.Using)
	}

	container = withAction(container, actn)
	container = withActionsCache(container)
	container = withDefaultGitHubEnvVars(container)
	container = withGitHubEnv(container)
	container = withGitHubPath(container)
	container = withGitHubOutput(container)
	container = withGitHubState(container)
	container = withRunnerTmp(container)
	container = withRunnerToolcache(container)
	container = withHome(container)
	container = withWorkspace(container, workspace)

	return &PreAction{
		Action: Action{
			PostAction: PostAction{
				FinalizedAction: FinalizedAction{
					Ctr: container,
				},
				Metadata: rawMetadata,
				Subpath:  subpath,
			},
		},
	}, nil
}

// Pre executes the pre-step of the GitHub Action in the underlying container.
func (a *PreAction) Pre(ctx context.Context) (*Action, error) {
	metadata, err := parseActionMetadata(a.Metadata)
	if err != nil {
		return nil, err
	}

	if a.Ref != nil {
		a.Ctr = withGitHubEnvVarsFromRef(ctx, a.Container(), a.Ref)
	}

	wkv, err := parseKeyValuePairs(a.Inputs)
	if err != nil {
		return nil, err
	}

	a.Ctr, err = withInputs(a.Container(), metadata, wkv)
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
			a.Ctr = a.Container().WithExec(append([]string{shimPath, "--", metadata.Runs.PreEntrypoint}, metadata.Runs.Args...))

			a.Ctr, err = withExportedEnv(ctx, a.Container())
			if err != nil {
				return nil, err
			}
		}
	case metadata.IsNode():
		if metadata.Runs.Pre != "" {
			a.Ctr = a.Container().WithExec(append([]string{shimPath, "--", "node", path.Join(actionPath, a.Subpath, metadata.Runs.Pre)}, metadata.Runs.Args...))

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
	metadata, err := parseActionMetadata(a.Metadata)
	if err != nil {
		return nil, err
	}

	switch {
	case metadata.IsDocker():
		if metadata.Runs.Entrypoint != "" {
			a.Ctr = a.Container().WithExec(append([]string{shimPath, "--", metadata.Runs.Entrypoint}, metadata.Runs.Args...))
		} else {
			entrypoint, err := a.Container().Entrypoint(ctx)
			if err != nil {
				return nil, err
			}

			a.Ctr = a.Container().WithExec(append(append([]string{shimPath, "--"}, entrypoint...), metadata.Runs.Args...))
		}
	case metadata.IsNode():
		a.Ctr = a.Container().WithExec(append([]string{shimPath, "--", "node", path.Join(actionPath, a.Subpath, metadata.Runs.Main)}, metadata.Runs.Args...))
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
func (a *PostAction) Post(ctx context.Context) (*FinalizedAction, error) {
	metadata, err := parseActionMetadata(a.Metadata)
	if err != nil {
		return nil, err
	}

	switch {
	case metadata.IsDocker():
		if metadata.Runs.PostEntrypoint != "" {
			a.Ctr = a.Container().WithExec(append([]string{shimPath, "--", metadata.Runs.PostEntrypoint}, metadata.Runs.Args...))

			a.Ctr, err = withExportedEnv(ctx, a.Container())
			if err != nil {
				return nil, err
			}
		}
	case metadata.IsNode():
		if metadata.Runs.Post != "" {
			a.Ctr = a.Container().WithExec(append([]string{shimPath, "--", "node", path.Join(actionPath, a.Subpath, metadata.Runs.Post)}, metadata.Runs.Args...))

			a.Ctr, err = withExportedEnv(ctx, a.Container())
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("actions that run using %s are not supported", metadata.Runs.Using)
	}

	return &a.FinalizedAction, nil
}

// Post executes the pre-, main and post-steps of the GitHub Action in the underlying container.
func (a *PreAction) Post(ctx context.Context) (*FinalizedAction, error) {
	main, err := a.Pre(ctx)
	if err != nil {
		return nil, err
	}

	return main.Post(ctx)
}

// Post executes the main and post-steps of the GitHub Action in the underlying container.
func (a *Action) Post(ctx context.Context) (*FinalizedAction, error) {
	post, err := a.Main(ctx)
	if err != nil {
		return nil, err
	}

	return post.Post(ctx)
}

// WithInput set an input (i.e. a key-value pair from the `with` object of a GitHub Action).
func (a *PreAction) WithInput(name, value string) *PreAction {
	a.Action = *a.Action.WithInput(name, value)
	return a
}

// WithEnv sets a new environment variable for the action.
func (a *PreAction) WithEnv(name, value string) *PreAction {
	a.Action = *a.Action.WithEnv(name, value)
	return a
}

// WithToken sets the GitHub token for the action.
func (a *PreAction) WithToken(token *dagger.Secret) *PreAction {
	a.Action = *a.Action.WithToken(token)
	return a
}

// WithDebug enables debug for the action.
func (a *PreAction) WithDebug() *PreAction {
	a.Action = *a.Action.WithDebug()
	return a
}

// WithRef sets the GitHub environment variables for the action from the given ref.
func (a *PreAction) WithRef(ref *dagger.GitRef) *PreAction {
	a.Action = *a.Action.WithRef(ref)
	return a
}

// WithInput set an input (i.e. a key-value pair from the `with` object of a GitHub Action).
func (a *Action) WithInput(name, value string) *Action {
	a.PostAction = *a.PostAction.WithInput(name, value)
	return a
}

// WithEnv sets a new environment variable for the action.
func (a *Action) WithEnv(name, value string) *Action {
	a.PostAction = *a.PostAction.WithEnv(name, value)
	return a
}

// WithToken sets the GitHub token for the action.
func (a *Action) WithToken(token *dagger.Secret) *Action {
	a.PostAction = *a.PostAction.WithToken(token)
	return a
}

// WithDebug enables debug for the action.
func (a *Action) WithDebug() *Action {
	a.PostAction = *a.PostAction.WithDebug()
	return a
}

// WithRef sets the GitHub environment variables for the action from the given ref.
func (a *Action) WithRef(ref *dagger.GitRef) *Action {
	a.PostAction = *a.PostAction.WithRef(ref)
	return a
}

// WithInput set an input (i.e. a key-value pair from the `with` object of a GitHub Action).
func (a *PostAction) WithInput(name, value string) *PostAction {
	a.Inputs = append(a.Inputs, fmt.Sprintf("%s=%s", name, value))
	return a
}

// WithEnv sets a new environment variable for the action.
func (a *PostAction) WithEnv(name, value string) *PostAction {
	a.Ctr = a.Ctr.WithEnvVariable(name, value)
	return a
}

// WithToken sets the GitHub token for the action.
func (a *PostAction) WithToken(token *dagger.Secret) *PostAction {
	a.Ctr = withToken(a.Container(), token)
	return a
}

// WithDebug enables debug for the action.
func (a *PostAction) WithDebug() *PostAction {
	a.Ctr = withDebug(a.Container())
	return a
}

// WithRef sets the GitHub environment variables for the action from the given ref.
func (a *PostAction) WithRef(ref *dagger.GitRef) *PostAction {
	a.Ref = ref
	return a
}

// Container gives access to the underlying container.
func (a *FinalizedAction) Container() *dagger.Container {
	return a.Ctr
}

// Terminal is a convenient alias for Container().Terminal().
func (a *FinalizedAction) Terminal() *dagger.Container {
	return a.Container().Terminal()
}

// Sync is a convenient alias for Container().Sync().
func (a *FinalizedAction) Sync(ctx context.Context) (*dagger.Container, error) {
	return a.Container().Sync(ctx)
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
	state, err := gitHubState(ctx, a.Container())
	if err != nil {
		return "", err
	}

	return strings.Join(envconv.MapToArr(state), "\n"), nil
}

// Output returns the parsed key-value pairs that were saved to GITHUB_OUTPUT.
func (a *FinalizedAction) Output(ctx context.Context) (string, error) {
	output, err := gitHubOutput(ctx, a.Container())
	if err != nil {
		return "", err
	}

	return strings.Join(envconv.MapToArr(output), "\n"), nil
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

	workdir := path.Join("$GOPATH", "src", gomod.Module.Mod.Path)

	shim := golang.
		WithEnvVariable("CGO_ENABLED", "0").
		WithDirectory(workdir, src, dagger.ContainerWithDirectoryOpts{
			Expand: true,
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
		WithWorkdir(workdir, dagger.ContainerWithWorkdirOpts{Expand: true}).
		WithExec([]string{"go", "build", "-o", shimPath, "./cmd/shim"}).
		File(shimPath)

	return container.WithFile(shimPath, shim), nil
}

func actionMetadata(ctx context.Context, dir *dagger.Directory, subpath string) (string, error) {
	var errs error

	for _, actionYAMLFileName := range githubactions.ActionYAMLFilenames {
		contents, err := dir.File(path.Join(subpath, actionYAMLFileName)).Contents(ctx)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		return contents, nil
	}

	return "", fmt.Errorf("find action metadata: %w", errs)
}

func parseActionMetadata(rawMetadata string) (*githubactions.Metadata, error) {
	metadata := &githubactions.Metadata{}

	if err := yaml.Unmarshal([]byte(rawMetadata), metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func withDebug(container *dagger.Container) *dagger.Container {
	return container.
		WithEnvVariable(githubactions.SecretActionsRunnerDebug, fmt.Sprint(true)).
		WithEnvVariable(githubactions.SecretActionsStepDebug, fmt.Sprint(true)).
		WithEnvVariable(githubactions.SecretRunnerDebug, fmt.Sprint(1))
}

func withToken(container *dagger.Container, token *dagger.Secret) *dagger.Container {
	return container.WithSecretVariable(githubactions.EnvVarToken, token)
}

func withHome(container *dagger.Container) *dagger.Container {
	return container.
		WithEnvVariable("HOME", homePath).
		WithMountedDirectory(homePath, dag.Directory())
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
		// Ensure trailing slash. Ref: https://gha-cache-server.falcondev.io/getting-started#_2-self-hosted-runner-setup.
		WithEnvVariable(githubactions.EnvVarActionsResultsURL, fmt.Sprintf("%s/", actionsResultsURL)).
		WithEnvVariable(githubactions.EnvVarActionsCacheURL, fmt.Sprintf("%s/", actionsResultsURL)).
		WithEnvVariable(githubactions.EnvVarActionsCacheServiceV2, fmt.Sprint(true)).
		// Ref: https://github.com/actions/toolkit/blob/f58042f9cc16bcaa87afaa86c2974a8c771ce1ea/packages/cache/src/internal/cacheUtils.ts#L162.
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
		WithMountedDirectory(tmpPath, dag.Directory())
}

func withRunnerToolcache(container *dagger.Container) *dagger.Container {
	return container.
		WithEnvVariable(githubactions.EnvVarRunnerToolCache, toolcachePath).
		WithMountedDirectory(toolcachePath, dag.Directory())
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

	return container.WithEnvVariable(
		"PATH", fmt.Sprintf("%s:$PATH", newPath),
		dagger.ContainerWithEnvVariableOpts{Expand: true},
	), nil
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

func withInputs(container *dagger.Container, metadata *githubactions.Metadata, with map[string]string) (*dagger.Container, error) {
	with, err := metadata.InputsFromWith(with)
	if err != nil {
		return nil, err
	}

	for k, v := range with {
		container = container.WithEnvVariable("INPUT_"+strings.ToUpper(k), v)
	}

	return container, nil
}
