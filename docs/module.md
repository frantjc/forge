# Import as a Go module

Forge can be used as a library as well. Each type of step from a proprietary CI system is represented as a [`Runnable`](https://github.com/frantjc/forge/blob/main/runnable.go). For example: a GitHub [`Action`](https://github.com/frantjc/forge/blob/main/action.go).

There are a few additional helper `Runnables`s: [`Pipe`](https://github.com/frantjc/forge/blob/main/pipe.go), which pipes the stdout of one ore to the stdin of another and [`Pure`](https://github.com/frantjc/forge/blob/main/ore/pure.go), which simply runs one containerized command.

In this example, a `Pipe` is used to pipe the stdout of a GitHub Action, which is using `actions/checkout` to check out [github.com/frantjc/forge](https://github.com/frantjc/forge), to a `Pure` running `grep` to only printing lines that contain the string `"debug"`:

```go
import (
	// Some stdlib imports omitted for brevity.

	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/githubactions"
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
	if err = (&forge.Pipe{
		From: &forge.Action{
			Uses:          "actions/checkout@v4",
			GlobalContext: globalContext,
		},
		To: &forge.Pure{
			Image:      "alpine:3.20",
			Entrypoint: []string{"grep", "debug"},
		},
	}).Run(ctx, docker.New(cli, ""), forge.WithStdStreams()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
```
