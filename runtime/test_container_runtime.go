package runtime

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/frantjc/forge"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestContainerRuntimeConformance(t *testing.T, cr forge.ContainerRuntime) {
	t.Run("PullImage", func(t *testing.T) {
		TestPullImage(t, cr)
	})
	t.Run("CreateContainer", func(t *testing.T) {
		TestCreateContainer(t, cr)
	})
	t.Run("StartContainer", func(t *testing.T) {
		TestStartContainer(t, cr)
	})
	t.Run("ContainerExec", func(t *testing.T) {
		TestContainerExec(t, cr)
	})
	t.Run("ContainerCopy", func(t *testing.T) {
		TestContainerCopy(t, cr)
	})
}

func TestPullImage(t *testing.T, cr forge.ContainerRuntime) {
	ctx := t.Context()

	img, err := cr.PullImage(ctx, "public.ecr.aws/docker/library/alpine")
	require.NoError(t, err)
	require.NotNil(t, img)
	require.NotEmpty(t, img.Name())

	cfg, err := img.Config()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	require.NotNil(t, img.Blob())
}

func TestCreateContainer(t *testing.T, cr forge.ContainerRuntime) {
	ctx := t.Context()
	img, err := cr.PullImage(ctx, "public.ecr.aws/docker/library/alpine")
	require.NoError(t, err)

	c, err := cr.CreateContainer(ctx, img, &forge.ContainerConfig{
		Entrypoint: []string{"/bin/sh", "-c"},
		Cmd:        []string{"exit 0"},
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		ctx = context.WithoutCancel(ctx)
		require.NoError(t, c.Remove(ctx))
	})

	require.NotEmpty(t, c.GetID())
}

func TestStartContainer(t *testing.T, cr forge.ContainerRuntime) {
	ctx := t.Context()
	img, err := cr.PullImage(ctx, "public.ecr.aws/docker/library/alpine")
	require.NoError(t, err)

	c, err := cr.CreateContainer(ctx, img, &forge.ContainerConfig{
		Entrypoint: []string{"/bin/sh", "-c"},
		Cmd:        []string{"exit 0"},
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		ctx = context.WithoutCancel(ctx)
		require.NoError(t, c.Remove(ctx))
	})

	require.NoError(t, c.Start(ctx))
}

func TestContainerExec(t *testing.T, cr forge.ContainerRuntime) {
	ctx := t.Context()
	img, err := cr.PullImage(ctx, "public.ecr.aws/docker/library/alpine")
	require.NoError(t, err)

	c, err := cr.CreateContainer(ctx, img, &forge.ContainerConfig{
		Entrypoint: []string{"/bin/sh", "-c"},
		Cmd:        []string{"sleep infinity"},
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		ctx = context.WithoutCancel(ctx)
		require.NoError(t, c.Stop(ctx))
		require.NoError(t, c.Remove(ctx))
	})

	require.NoError(t, c.Start(ctx))

	out := new(bytes.Buffer)
	expected := uuid.NewString()
	exitCode, err := c.Exec(ctx, &forge.ContainerConfig{
		Entrypoint: []string{"echo", expected},
	}, &forge.Streams{
		Out: out,
		Err: io.Discard,
	})
	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
	actual := strings.TrimSpace(out.String())
	require.Equal(t, expected, actual)
}

func TestContainerCopy(t *testing.T, cr forge.ContainerRuntime) {
	ctx := t.Context()
	img, err := cr.PullImage(ctx, "public.ecr.aws/docker/library/alpine")
	require.NoError(t, err)
	c, err := cr.CreateContainer(ctx, img, &forge.ContainerConfig{
		Entrypoint: []string{"/bin/sh", "-c"},
		Cmd:        []string{"sleep"},
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		ctx = context.WithoutCancel(ctx)
		require.NoError(t, c.Stop(ctx))
		require.NoError(t, c.Remove(ctx))
	})

	require.NoError(t, c.Start(ctx))

	expected := []byte(uuid.NewString())
	path := filepath.Join("/tmp", uuid.NewString())

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	require.NoError(t, tw.WriteHeader(&tar.Header{
		Name: filepath.Base(path),
		Mode: 0644,
		Size: int64(len(expected)),
	}))
	_, err = tw.Write(expected)
	require.NoError(t, err)
	require.NoError(t, tw.Close())

	require.NoError(t, c.CopyTo(ctx, filepath.Dir(path), buf))

	rc, err := c.CopyFrom(ctx, path)
	require.NoError(t, err)
	defer rc.Close()

	data, err := io.ReadAll(rc)
	require.NoError(t, err)

	tr := tar.NewReader(bytes.NewReader(data))
	_, err = tr.Next()
	require.NoError(t, err)
	actual, err := io.ReadAll(tr)
	require.NoError(t, err)

	require.Equal(t, expected, actual)
}
