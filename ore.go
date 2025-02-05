package forge

import (
	"context"
	"path/filepath"

	"github.com/google/uuid"
)

func CloudBuildWorkingDir(workingDir string) string {
	return filepath.Join(workingDir, "cloudbuild")
}

func InterceptingDockerSock(workingDir string) string {
	return filepath.Join(workingDir, "forge.sock")
}

func GitHubWorkspace(workingDir string) string {
	return filepath.Join(workingDir, "github/workspace")
}

func GitHubActionPath(workingDir string) string {
	return filepath.Join(workingDir, "github/action")
}

func GitHubRunnerTmp(workingDir string) string {
	return filepath.Join(workingDir, "github/runner/tmp")
}

func GitHubRunnerToolCache(workingDir string) string {
	return filepath.Join(workingDir, "github/runner/toolcache")
}

func GitHubPath(workingDir string) string {
	return filepath.Join(workingDir, "github/add_path")
}

func GitHubEnv(workingDir string) string {
	return filepath.Join(workingDir, "github/set_env")
}

func GitHubOutput(workingDir string) string {
	return filepath.Join(workingDir, "github/set_output")
}

func GitHubState(workingDir string) string {
	return filepath.Join(workingDir, "github/save_state")
}

func ConcourseResourceWorkingDir(workingDir string) string {
	return filepath.Join(workingDir, "concourse/resource")
}

func AzureDevOpsTaskWorkingDir(workingDir string) string {
	return filepath.Join(workingDir, "task")
}

func oreOptsWithDefaults(opts ...OreOpt) *OreOpts {
	o := &OreOpts{
		WorkingDir: "/" + uuid.NewString(),
	}

	for _, opt := range opts {
		opt.Apply(o)
	}

	return o
}

type OreOpts struct {
	Streams             *Streams
	Mounts              []Mount
	InterceptDockerSock bool
	WorkingDir          string
}

func (o *OreOpts) Apply(opts *OreOpts) {
	if opts == nil {
		opts = &OreOpts{}
	}
	if o.Streams != nil {
		opts.Streams = o.Streams
	}
	opts.Mounts = overrideMounts(opts.Mounts, o.Mounts...)
	if o.InterceptDockerSock {
		opts.InterceptDockerSock = true
	}
	if o.WorkingDir != "" {
		opts.WorkingDir = o.WorkingDir
	}
}

type OreOpt interface {
	Apply(*OreOpts)
}

// Ore represents one or more sequential containerized commands.
type Ore interface {
	Liquify(context.Context, ContainerRuntime, ...OreOpt) error
}
