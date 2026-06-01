//go:build shim && dockerd

package forge_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/githubactions"
	"github.com/stretchr/testify/require"
)

func TestActionRunDockerfile(t *testing.T) {
	cr := Runtime(t)

	uses := Uses(t, &githubactions.Metadata{
		Name: "test-dockerfile",
		Runs: &githubactions.MetadataRuns{
			Using: githubactions.RunsUsingDocker,
			Image: "Dockerfile",
		},
	})
	require.NoError(t, os.WriteFile(filepath.Join(uses, "Dockerfile"), []byte(`FROM public.ecr.aws/docker/library/alpine
CMD ["/bin/sh", "-c", "exit 0"]
`), 0o644))

	action := &forge.Action{Uses: uses}

	require.NoError(t, action.Run(t.Context(), cr, forge.WithStreams(Streams(t)), MountShim(t)))
}

func TestActionRunDockerfileWithArgs(t *testing.T) {
	cr := Runtime(t)

	uses := Uses(t, &githubactions.Metadata{
		Name: "test-dockerfile-args",
		Runs: &githubactions.MetadataRuns{
			Using:      githubactions.RunsUsingDocker,
			Image:      "Dockerfile",
			Entrypoint: "/bin/sh",
			Args:       []string{"-c", "exit 0"},
		},
	})
	require.NoError(t, os.WriteFile(filepath.Join(uses, "Dockerfile"), []byte(`FROM public.ecr.aws/docker/library/alpine
`), 0o644))

	action := &forge.Action{Uses: uses}

	require.NoError(t, action.Run(t.Context(), cr, forge.WithStreams(Streams(t)), MountShim(t)))
}

func TestActionRunDockerfileWithUnusualName(t *testing.T) {
	cr := Runtime(t)

	uses := Uses(t, &githubactions.Metadata{
		Name: "test-dockerfile-args",
		Runs: &githubactions.MetadataRuns{
			Using:      githubactions.RunsUsingDocker,
			Image:      "test.Dockerfile",
			Entrypoint: "/bin/sh",
			Args:       []string{"-c", "exit 0"},
		},
	})
	require.NoError(t, os.WriteFile(filepath.Join(uses, "test.Dockerfile"), []byte(`FROM public.ecr.aws/docker/library/alpine
`), 0o644))

	action := &forge.Action{Uses: uses}

	require.NoError(t, action.Run(t.Context(), cr, forge.WithStreams(Streams(t)), MountShim(t)))
}
