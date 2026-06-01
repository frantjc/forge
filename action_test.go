//go:build shim && (docker || dockerd)

package forge_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/githubactions"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func Uses(t *testing.T, metadata *githubactions.Metadata) string {
	t.Helper()
	require.GreaterOrEqual(t, len(githubactions.ActionYAMLFilenames), 1)
	b, err := yaml.Marshal(metadata)
	require.NoError(t, err)
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, githubactions.ActionYAMLFilenames[0]), b, 0o644))
	return dir
}

func TestActionRunDocker(t *testing.T) {
	cr := Runtime(t)

	uses := Uses(t, &githubactions.Metadata{
		Name: "test-docker",
		Runs: &githubactions.MetadataRuns{
			Using:      githubactions.RunsUsingDocker,
			Image:      "docker://public.ecr.aws/docker/library/alpine",
			Entrypoint: "/bin/sh",
			Args:       []string{"-c", "exit 0"},
		},
	})

	action := &forge.Action{Uses: uses}

	require.NoError(t, action.Run(t.Context(), cr, forge.WithStreams(Streams(t)), MountShim(t)))
}

func TestActionRunDockerNonzeroExitCode(t *testing.T) {
	cr := Runtime(t)

	uses := Uses(t, &githubactions.Metadata{
		Name: "test-docker-fail",
		Runs: &githubactions.MetadataRuns{
			Using:      githubactions.RunsUsingDocker,
			Image:      "docker://public.ecr.aws/docker/library/alpine",
			Entrypoint: "/bin/sh",
			Args:       []string{"-c", "exit 1"},
		},
	})

	action := &forge.Action{Uses: uses}

	require.Error(t, action.Run(t.Context(), cr, forge.WithStreams(Streams(t)), MountShim(t)))
}

func TestActionRunDockerWithEnv(t *testing.T) {
	cr := Runtime(t)

	uses := Uses(t, &githubactions.Metadata{
		Name: "test-docker-env",
		Runs: &githubactions.MetadataRuns{
			Using:      githubactions.RunsUsingDocker,
			Image:      "docker://public.ecr.aws/docker/library/alpine",
			Entrypoint: "/bin/sh",
			Args:       []string{"-c", `[ "$MY_VAR" = "hello" ]`},
			Env:        map[string]string{"MY_VAR": "hello"},
		},
	})

	action := &forge.Action{Uses: uses}

	require.NoError(t, action.Run(t.Context(), cr, forge.WithStreams(Streams(t)), MountShim(t)))
}

func TestActionRunDockerWithUserEnv(t *testing.T) {
	cr := Runtime(t)

	uses := Uses(t, &githubactions.Metadata{
		Name: "test-docker-env",
		Runs: &githubactions.MetadataRuns{
			Using:      githubactions.RunsUsingDocker,
			Image:      "docker://public.ecr.aws/docker/library/alpine",
			Entrypoint: "/bin/sh",
			Args:       []string{"-c", `[ "$MY_VAR" = "hello" ]`},
			Env:        map[string]string{"MY_VAR": "hello"},
		},
	})

	action := &forge.Action{
		Uses: uses,
		Env:  map[string]string{"MY_VAR": "hello"},
	}

	require.NoError(t, action.Run(t.Context(), cr, forge.WithStreams(Streams(t)), MountShim(t)))
}

func TestActionRunDockerWithInputs(t *testing.T) {
	cr := Runtime(t)

	uses := Uses(t, &githubactions.Metadata{
		Name: "test-docker-inputs",
		Inputs: map[string]githubactions.MetadataInput{
			"greeting": {Description: "greeting word", Required: true},
		},
		Runs: &githubactions.MetadataRuns{
			Using:      githubactions.RunsUsingDocker,
			Image:      "docker://public.ecr.aws/docker/library/alpine",
			Entrypoint: "/bin/sh",
			Args:       []string{"-c", `[ "$INPUT_GREETING" = "hello" ]`},
		},
	})

	action := &forge.Action{
		Uses: uses,
		With: map[string]string{"greeting": "hello"},
	}

	require.NoError(t, action.Run(t.Context(), cr, forge.WithStreams(Streams(t)), MountShim(t)))
}

func TestActionRunTestdataNode(t *testing.T) {
	cr := Runtime(t)

	uses := Uses(t, &githubactions.Metadata{
		Name: "node",
		Runs: &githubactions.MetadataRuns{
			Using: githubactions.RunsUsingNode20,
			Main:  "main.js",
			Env:   map[string]string{"HELLO": "world"},
		},
	})
	require.NoError(t, os.WriteFile(filepath.Join(uses, "main.js"), []byte(
		"if (process.env.HELLO !== 'world') { process.exit(1); }\n",
	), 0o644))

	action := &forge.Action{Uses: uses}
	require.NoError(t, action.Run(t.Context(), cr, forge.WithStreams(Streams(t)), MountShim(t)))
}

func TestActionRunTestdataSaveState(t *testing.T) {
	cr := Runtime(t)

	uses := Uses(t, &githubactions.Metadata{
		Name: "save-state",
		Runs: &githubactions.MetadataRuns{
			Using:      githubactions.RunsUsingDocker,
			Image:      "docker://public.ecr.aws/docker/library/alpine",
			Entrypoint: "/bin/sh",
			Args:       []string{"-c", `printf 'general=kenobi\n' >> "$GITHUB_STATE"`},
		},
	})

	gc := githubactions.NewGlobalContextFromEnv()
	action := &forge.Action{Uses: uses, GlobalContext: gc}
	require.NoError(t, action.Run(t.Context(), cr, forge.WithStreams(Streams(t)), MountShim(t)))
	require.Equal(t, "kenobi", gc.EnvContext["STATE_general"])
}

func TestActionRunTestdataSetEnv(t *testing.T) {
	cr := Runtime(t)

	uses := Uses(t, &githubactions.Metadata{
		Name: "set-env",
		Runs: &githubactions.MetadataRuns{
			Using:      githubactions.RunsUsingDocker,
			Image:      "docker://public.ecr.aws/docker/library/alpine",
			Entrypoint: "/bin/sh",
			Args:       []string{"-c", `printf 'GENERAL=kenobi\n' >> "$GITHUB_ENV"`},
		},
	})

	gc := githubactions.NewGlobalContextFromEnv()
	action := &forge.Action{Uses: uses, GlobalContext: gc}
	require.NoError(t, action.Run(t.Context(), cr, forge.WithStreams(Streams(t)), MountShim(t)))
	require.Equal(t, "kenobi", gc.EnvContext["GENERAL"])
}

func TestActionRunTestdataSetOutput(t *testing.T) {
	cr := Runtime(t)

	uses := Uses(t, &githubactions.Metadata{
		Name: "set-output",
		Runs: &githubactions.MetadataRuns{
			Using:      githubactions.RunsUsingDocker,
			Image:      "docker://public.ecr.aws/docker/library/alpine",
			Entrypoint: "/bin/sh",
			Args:       []string{"-c", `printf 'general=kenobi\n' >> "$GITHUB_OUTPUT"`},
		},
	})

	gc := githubactions.NewGlobalContextFromEnv()
	action := &forge.Action{ID: "test", Uses: uses, GlobalContext: gc}
	require.NoError(t, action.Run(t.Context(), cr, forge.WithStreams(Streams(t)), MountShim(t)))
	require.Equal(t, "kenobi", gc.StepsContext["test"].Outputs["general"])
}

func TestActionRunTestdataYmlVYaml(t *testing.T) {
	cr := Runtime(t)

	require.Len(t, githubactions.ActionYAMLFilenames, 2)

	expected := "yml"
	uses := Uses(t, &githubactions.Metadata{
		Name: "yml-v-yaml",
		Runs: &githubactions.MetadataRuns{
			Using:      githubactions.RunsUsingDocker,
			Image:      "docker://public.ecr.aws/docker/library/alpine",
			Entrypoint: "echo",
			Args:       []string{expected},
		},
	})

	b, err := yaml.Marshal(&githubactions.Metadata{
		Name: "yml-v-yaml",
		Runs: &githubactions.MetadataRuns{
			Using:      githubactions.RunsUsingDocker,
			Image:      "docker://public.ecr.aws/docker/library/alpine",
			Entrypoint: "echo",
			Args:       []string{"yaml"},
		},
	})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(uses, githubactions.ActionYAMLFilenames[1]), b, 0o644))

	action := &forge.Action{Uses: uses}
	streams, stdout := StreamsCaptureStdout(t)
	require.NoError(t, action.Run(t.Context(), cr, forge.WithStreams(streams), MountShim(t)))
	actual := strings.TrimSpace(stdout.String())
	require.Equal(t, expected, actual)
}
