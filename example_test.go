package forge_test

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/githubactions"
	"github.com/frantjc/forge/ore"
	"github.com/frantjc/forge/runtime/docker"
)

// Checkout https://github.com/frantjc/forge
// using https://github.com/actions/checkout,
// grepping to only print debug logs.
func ExampleFoundry_Process() {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	globalContext := githubactions.NewGlobalContextFromEnv()
	globalContext.EnableDebug()
	globalContext.GitHubContext.Repository = "frantjc/forge"

	if err = forge.NewFoundry(docker.New(cli, false)).Process(
		ctx,
		&ore.Lava{
			From: &ore.Action{
				Uses:          "actions/checkout@v4",
				GlobalContext: globalContext,
			},
			To: &ore.Pure{
				Image:      "alpine:3.20",
				Entrypoint: []string{"grep", "debug"},
			},
		},
		forge.StdDrains(),
	); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
