package main

import (
	"context"
	"encoding/json"
	"fmt"
	"path"

	"github.com/frantjc/forge/concourse"
	"github.com/frantjc/forge/internal/dagger"
	"github.com/frantjc/forge/internal/envconv"
	"sigs.k8s.io/yaml"
)

// Resource has a container that's prepared to execute a resource method.
type Resource struct {
	// +private
	FinalizedResource
	// +private
	Source []string
}

const (
	resourcePath = "/forge/resource"
)

// Resource creates a container to execute a Concourse resource in.
func (m *Forge) Resource(
	ctx context.Context,
	// The resource to execute
	resource string,
	// The pipeline file to find the resource in
	// +defaultPath=".forge.yml"
	pipeline *dagger.File,
	// The workdir for the resource to execute in
	// +defaultPath="."
	workdir *dagger.Directory,
) (*Resource, error) {
	contents, err := pipeline.Contents(ctx)
	if err != nil {
		return nil, err
	}

	parsed := &concourse.Pipeline{}

	if err := yaml.Unmarshal([]byte(contents), parsed); err != nil {
		return nil, err
	}

	parsed.ResourceTypes = append(parsed.ResourceTypes, concourse.BuiltinResourceTypes...)

	r, rt, err := resourceAndType(parsed, resource)
	if err != nil {
		return nil, err
	}

	container := dag.Container().
		From(fmt.Sprintf("%s:%s", rt.Source.Repository, rt.Source.Tag)).
		WithWorkdir(path.Join(resourcePath, r.Name)).
		WithMountedDirectory(path.Join(resourcePath, r.Name), workdir)

	return &Resource{
		FinalizedResource: FinalizedResource{
			Ctr:  container,
			Name: resource,
		},
		Source: envconv.MapToArr(r.Source),
	}, nil
}

func resourceAndType(pipeline *concourse.Pipeline, resource string) (r *concourse.Resource, rt *concourse.ResourceType, err error) {
	for _, _r := range pipeline.Resources {
		if _r.Name == resource {
			r = &_r
			break
		}
	}

	if r == nil {
		err = fmt.Errorf("resource %s not found", resource)
		return
	}

	for _, _rt := range pipeline.ResourceTypes {
		if _rt.Name == r.Type {
			rt = &_rt
			break
		}
	}

	if rt == nil {
		r = nil
		err = fmt.Errorf("resource type %s not found", resource)
		return
	}

	return
}

// Get runs the get step.
func (r *Resource) Get(
	// +optional
	version []string,
	// +optional
	param []string,
) (*FinalizedResource, error) {
	stdin, err := getStdin(r.Source, version, param)
	if err != nil {
		return nil, err
	}

	r.Ctr = r.Container().WithExec([]string{concourse.EntrypointGet, path.Join(resourcePath, r.Name)}, dagger.ContainerWithExecOpts{
		Stdin: stdin,
	})

	return &r.FinalizedResource, nil
}

// Resource has a container that's executed a resource method.
type FinalizedResource struct {
	// +private
	Ctr *dagger.Container
	// +private
	Name string
}

// Check runs the check step.
func (r *Resource) Check(
	// +optional
	version []string,
) (*FinalizedResource, error) {
	stdin, err := getStdin(r.Source, version, nil)
	if err != nil {
		return nil, err
	}

	r.Ctr = r.Container().WithExec([]string{concourse.EntrypointCheck, path.Join(resourcePath, r.Name)}, dagger.ContainerWithExecOpts{
		Stdin: stdin,
	})

	return &r.FinalizedResource, nil
}

// Put runs the put step.
func (r *Resource) Put(
	// +optional
	param []string,
) (*FinalizedResource, error) {
	stdin, err := getStdin(r.Source, nil, param)
	if err != nil {
		return nil, err
	}

	r.Ctr = r.Container().WithExec([]string{concourse.EntrypointPut, path.Join(resourcePath, r.Name)}, dagger.ContainerWithExecOpts{
		Stdin: stdin,
	})

	return &r.FinalizedResource, nil
}

func getStdin(source, version, params []string) (string, error) {
	input := &concourse.Input{}

	skv, err := parseKeyValuePairs(source)
	if err != nil {
		return "", err
	}

	if len(skv) > 0 {
		input.Source = map[string]any{}

		for k, v := range skv {
			input.Source[k] = v
		}
	}

	vkv, err := parseKeyValuePairs(version)
	if err != nil {
		return "", err
	}

	if len(vkv) > 0 {
		input.Version = map[string]any{}

		for k, v := range vkv {
			input.Version[k] = v
		}
	}

	pkv, err := parseKeyValuePairs(params)
	if err != nil {
		return "", err
	}

	if len(pkv) > 0 {
		input.Params = map[string]any{}

		for k, v := range pkv {
			input.Params[k] = v
		}
	}

	b, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// Container gives access to the underlying container.
func (r *FinalizedResource) Container() *dagger.Container {
	return r.Ctr
}

// Sync is a convenient alias for Container().Sync().
func (a *FinalizedResource) Sync(ctx context.Context) (*dagger.Container, error) {
	return a.Container().Sync(ctx)
}

// Stdout is a convenient alias for Container().Stdout().
func (r *FinalizedResource) Stdout(ctx context.Context) (string, error) {
	return r.Container().Stdout(ctx)
}

// Stderr is a convenient alias for Container().Stderr().
func (r *FinalizedResource) Stderr(ctx context.Context) (string, error) {
	return r.Container().Stderr(ctx)
}

// CombinedOutput is a convenient alias for Container().CombinedOutput().
func (r *FinalizedResource) CombinedOutput(ctx context.Context) (string, error) {
	return r.Container().CombinedOutput(ctx)
}

// Terminal is a convenient alias for Container().Terminal().
func (r *FinalizedResource) Terminal() *dagger.Container {
	return r.Container().Terminal()
}

// Workdir returns the current state of the working directory.
func (r *FinalizedResource) Workdir() *dagger.Directory {
	return r.Container().Directory(path.Join(resourcePath, r.Name))
}
