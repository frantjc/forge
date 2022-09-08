package main

import (
	"context"
	"os"
	"path"

	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/events"
	"github.com/frantjc/forge/pkg/basin/bucket"
	"github.com/frantjc/forge/pkg/concourse"
	"github.com/frantjc/forge/pkg/github/actions"
	"github.com/frantjc/forge/pkg/ore"
	"github.com/frantjc/forge/pkg/runtime/container/docker"

	_ "gocloud.dev/blob/fileblob"
)

var (
	logr = forge.NewLogger()
	ctx  = forge.WithLogger(context.Background(), logr)
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	globalContext, err := actions.NewGlobalContextFromPath(ctx, wd)
	if err != nil {
		panic(err)
	}

	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	cache, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	cache = path.Join(cache, "forge")
	if err = os.MkdirAll(cache, 0777); err != nil {
		panic(err)
	}

	basin, err := bucket.New(ctx, "file://"+cache)
	if err != nil {
		panic(err)
	}

	var (
		conatinerRuntime                  = docker.New(c)
		_                                 = globalContext
		_                *concourse.Input = nil
	)

	events.Listen(ctx, events.ContainerCreated, func(ctx context.Context, e *events.Event) {
		// container, err := conatinerRuntime.GetContainer(ctx, events.ContainerMetadata(e.GetMetadata()).GetId())
		// if err != nil {
		// 	panic(err)
		// }

		// streams, restore := forge.StdTerminalStreams()
		// defer func() {
		// 	if err = restore(); err != nil {
		// 		panic(err)
		// 	}
		// }()

		// if _, err = container.Exec(ctx, &forge.ContainerConfig{
		// 	Entrypoint: []string{"sh"},
		// }, streams); err != nil {
		// 	panic(err)
		// }
	})

	foundry := &forge.Foundry{ContainerRuntime: conatinerRuntime, Basin: basin}
	if _, err = foundry.Process(
		ctx,
		&ore.Alloy{
			Ores: []forge.Ore{
				// &ore.Action{
				// 	Uses: "actions/setup-go@v3",
				// 	With: map[string]string{
				// 		"go-version": "1.19",
				// 	},
				// 	GlobalContext: globalContext,
				// },
				&ore.Resource{
					Method: "get",
					Resource: &concourse.Resource{
						Name: "github.com/frantjc/forge",
						Type: "git",
						Source: map[string]string{
							"uri":    "https://github.com/frantjc/forge",
							"branch": "main",
						},
					},
					ResourceType: &concourse.ResourceType{
						Name: "git",
						Source: &concourse.Source{
							Repository: "docker.io/concourse/git-resource",
							Tag:        "alpine",
						},
					},
				},
				&ore.Pure{
					Image:      "alpine",
					Entrypoint: []string{"ls", "-al", "github.com/frantjc/forge"},
				},
				// &ore.Lava{
				// 	From: &ore.Pure{
				// 		Image:      "alpine",
				// 		Entrypoint: []string{"echo", "hello"},
				// 	},
				// 	To: &ore.Pure{
				// 		Image:      "alpine",
				// 		Entrypoint: []string{"base64"},
				// 	},
				// },
			},
		},
		forge.StdDrains(),
	); err != nil {
		panic(err)
	}
}
