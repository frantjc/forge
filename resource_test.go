//go:build shim && (docker || dockerd)

package forge_test

import (
	"testing"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/concourse"
	"github.com/stretchr/testify/require"
)

func TestResourceCheck(t *testing.T) {
	ctx := t.Context()
	cr := Runtime(t)

	r := &forge.Resource{
		Method: concourse.MethodCheck,
		Resource: &concourse.Resource{
			Name: "mock",
			Source: map[string]any{
				"create_files": map[string]any{},
			},
		},
		ResourceType: &concourse.ResourceType{
			Name: "mock",
			Source: &concourse.ResourceTypeSource{
				Repository: "concourse/mock-resource",
				Tag:        "latest",
			},
		},
	}

	require.NoError(t, r.Run(ctx, cr, forge.WithStreams(Streams(t)), MountShim(t)))
}

func TestResourceGet(t *testing.T) {
	ctx := t.Context()
	cr := Runtime(t)

	r := &forge.Resource{
		Method: concourse.MethodGet,
		Version: map[string]any{
			"version": "mock",
		},
		Resource: &concourse.Resource{
			Name: "mock",
			Source: map[string]any{
				"create_files": map[string]any{},
			},
		},
		ResourceType: &concourse.ResourceType{
			Name: "mock",
			Source: &concourse.ResourceTypeSource{
				Repository: "concourse/mock-resource",
				Tag:        "latest",
			},
		},
	}

	require.NoError(t, r.Run(ctx, cr, forge.WithStreams(Streams(t)), MountShim(t)))
}

func TestResourcePut(t *testing.T) {
	ctx := t.Context()
	cr := Runtime(t)

	r := &forge.Resource{
		Method: concourse.MethodPut,
		Params: map[string]any{
			"version": "mock",
		},
		Resource: &concourse.Resource{
			Name: "mock",
			Source: map[string]any{
				"create_files": map[string]any{},
			},
		},
		ResourceType: &concourse.ResourceType{
			Name: "mock",
			Source: &concourse.ResourceTypeSource{
				Repository: "concourse/mock-resource",
				Tag:        "latest",
			},
		},
	}

	require.NoError(t, r.Run(ctx, cr, forge.WithStreams(Streams(t)), MountShim(t)))
}
