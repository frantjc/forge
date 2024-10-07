# Import as a Go module

Forge can be used as a library as well. Each type of step from a proprietary CI system is represented as an [`Ore`](https://github.com/frantjc/forge/blob/main/ore.go) and live in the package [`ore`](https://github.com/frantjc/forge/blob/main/ore). For example: a GitHub [`Action`](https://github.com/frantjc/forge/blob/main/ore/action.go).

There are a few additional helper `Ore`s: [`Lava`](https://github.com/frantjc/forge/blob/main/ore/lava.go), which pipes the stdout of one ore to the stdin of another; [`Alloy`](https://github.com/frantjc/forge/blob/main/ore/alloy.go), which runs other Ores sequentially that will share a working directory; and [`Pure`](https://github.com/frantjc/forge/blob/main/ore/pure.go), which simply runs one containerized command.

In this example, a `Lava` is used to pipe the stdout of a GitHub Action, which is using `actions/checkout` to check out [github.com/frantjc/forge](https://github.com/frantjc/forge), to a `Pure` running `grep` to only printing lines that contain the string `"debug"`:

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
```
