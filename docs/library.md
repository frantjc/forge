# Import as a Go module

Forge can be used as a library as well. Each type of step from a proprietary CI system is represented as an [`Ore`](../ore.go) and live in the package [`ore`](../ore). For example: a GitHub [`Action`](../ore/action.go).

There are a few additional helper `Ore`s: [`Lava`](../ore/lava.go), which pipes the stdout of one or to the stdin of another; [`Alloy`](../ore/alloy.go), which runs other sequential Ores that will share a working directory; and [`Pure`](../ore/pure.go), which simply runs one containerized command.

In this example, a `Lava` is used to pipe the stdout of a GitHub Action, which is using `actions/checkout` to check out https://github.com/frantjc/forge, to a `Pure` running `grep`, only printing lines that contain the string `"debug"`:

```go
import (
	// Some std imports omitted for brevity.

	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/githubactions"
	"github.com/frantjc/forge/ore"
	"github.com/frantjc/forge/runtime/docker"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	globalContext := githubactions.NewGlobalContextFromEnv().EnableDebug()
	globalContext.GitHubContext.Repository = "frantjc/forge"

  	// Checkout https://github.com/frantjc/forge
    // using https://github.com/actions/checkout,
    // grepping to only print debug logs.
	if err = forge.NewFoundry(docker.New(cli)).Process(
		ctx,
		&ore.Lava{
			From: &ore.Action{
				Uses:          "actions/checkout@v4",
				GlobalContext: globalContext,
			},
			To: &ore.Pure{
				Image:      "alpine:3.19",
				Entrypoint: []string{"grep", "debug"},
			},
		},
		forge.StdDrains(),
	); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
```
