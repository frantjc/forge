package main

import (
	"context"
	"os"

	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/github/actions"
	"github.com/frantjc/forge/pkg/ore"
	"github.com/frantjc/forge/pkg/runtime/container/docker"
)

var (
	ctx = forge.WithLogger(
		context.Background(),
		forge.NewLogger(),
	)
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

	foundry := &forge.Foundry{ContainerRuntime: docker.New(c)}
	if _, err = foundry.Process(
		actions.WithGlobalContext(ctx, globalContext),
		&ore.Pure{
			Image:      "alpine",
			Entrypoint: []string{"base64"},
		},
		// &ore.Action{
		// 	Uses: "actions/setup-go@v3",
		// 	With: map[string]string{
		// 		"go-version": "1.19",
		// 	},
		// }
		forge.StdStreams(),
	); err != nil {
		panic(err)
	}
}
