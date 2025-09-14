package main

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/frantjc/forge/cloudbuild"
	"github.com/frantjc/forge/internal/dagger"
)

type CloudBuild struct {
	FinalizedCloudBuild
}

type FinalizedCloudBuild struct {
	Ctr *dagger.Container
}

const (
	scriptPath  = "/forge/script"
	workdirPath = "/forge/workdir"
)

func (f *Forge) CloudBuild(
	ctx context.Context,
	name string,
	// +defaultPath="."
	workdir *dagger.Directory,
	// +optional
	entrypoint []string,
	// +optional
	args []string,
	// +optional
	env []string,
	// +optional
	gcloudConfig *dagger.Directory,
	// +optional
	script *dagger.File,
	// +optional
	substitutions []string,
	// +optional
	dynamicSubstitutions bool,
	// +optional
	automapSubstitutions bool,
) (*CloudBuild, error) {
	container := dag.Container().
		From(name).
		WithMountedCache(
			cloudbuild.WorkspacePath,
			dag.CacheVolume("workspace"),
		).
		WithMountedDirectory(workdirPath, workdir)

	container = withHome(container)

	if gcloudConfig != nil {
		container = container.WithMountedDirectory(path.Join(homePath, ".config", "gcloud"), gcloudConfig)
	}

	if script != nil {
		contents, err := script.Contents(ctx)
		if err != nil {
			return nil, err
		}

		if len(entrypoint) > 0 || len(args) > 0 {
			return nil, fmt.Errorf("cannot specify entrypoint or args with script")
		}

		container = withFile(
			container,
			scriptPath,
			contents,
			dagger.ContainerWithFileOpts{
				Permissions: 0700,
			},
		).WithExec([]string{scriptPath})
	} else {
		if len(entrypoint) == 0 {
			var err error
			entrypoint, err = container.Entrypoint(ctx)
			if err != nil {
				return nil, err
			}
		}

		container.WithExec(append(entrypoint, args...))
	}

	skv, err := parseKeyValuePairs(substitutions)
	if err != nil {
		return nil, err
	}

	if dynamicSubstitutions {
		for range []byte{0, 0} {
			for k, v := range skv {
				skv[k] = os.Expand(v, func(s string) string {
					if substitution, ok := skv[s]; ok {
						return substitution
					}

					return fmt.Sprintf("$%s", s)
				})
			}
		}
	}

	if automapSubstitutions {
		for k, v := range skv {
			container = container.WithEnvVariable(k, v)
		}
	}

	return &CloudBuild{
		FinalizedCloudBuild: FinalizedCloudBuild{
			Ctr: container,
		},
	}, nil
}

// Run executes the CloudBuild.
func (c *CloudBuild) Run(ctx context.Context) (*FinalizedCloudBuild, error) {
	if _, err := c.Stdout(ctx); err != nil {
		return nil, err
	}

	return &c.FinalizedCloudBuild, nil
}

// Container gives access to the underlying container.
func (c *FinalizedCloudBuild) Container() *dagger.Container {
	return c.Ctr
}

// Terminal is a convenient alias for Container().Terminal().
func (c *FinalizedCloudBuild) Terminal() *dagger.Container {
	return c.Container().Terminal()
}

// Stdout is a convenient alias for Container().Stdout().
func (c *FinalizedCloudBuild) Stdout(ctx context.Context) (string, error) {
	return c.Container().Stdout(ctx)
}

// Stderr is a convenient alias for Container().Stderr().
func (c *FinalizedCloudBuild) Stderr(ctx context.Context) (string, error) {
	return c.Container().Stderr(ctx)
}

// CombinedOutput is a convenient alias for Container().CombinedOutput().
func (c *FinalizedCloudBuild) CombinedOutput(ctx context.Context) (string, error) {
	return c.Container().CombinedOutput(ctx)
}

// Workspace returns the current state of the /workspace directory.
func (c *FinalizedCloudBuild) Workspace() *dagger.Directory {
	return c.Container().Directory(cloudbuild.WorkspacePath)
}

// Workdir returns the current state of the working directory.
func (c *FinalizedCloudBuild) Workdir() *dagger.Directory {
	return c.Container().Directory(workdirPath)
}
