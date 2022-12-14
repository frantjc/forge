# forge [![CI](https://github.com/frantjc/forge/actions/workflows/push.yml/badge.svg?branch=main&event=push)](https://github.com/frantjc/forge/actions) [![godoc](https://pkg.go.dev/badge/github.com/frantjc/forge.svg)](https://pkg.go.dev/github.com/frantjc/forge) [![goreportcard](https://goreportcard.com/badge/github.com/frantjc/forge)](https://goreportcard.com/report/github.com/frantjc/forge) ![license](https://shields.io/github/license/frantjc/forge)

<p align="center">
  <img src="https://raw.githubusercontent.com/frantjc/forge/main/docs/demo.gif">
</p>

Forge is a library and CLI for running reusable steps from various proprietary CI systems using a pluggable container runtime. This, for example, makes the functionality provided to GitHub Actions easily consumable (or testable) by users of other CI systems.

Forge currently exposes running [GitHub Actions](https://docs.github.com/en/actions/learn-github-actions/finding-and-customizing-actions) (e.g. [`actions/setup-go`](https://github.com/actions/setup-go)) and [Concourse Resources](https://concourse-ci.org/resources.html) (e.g. [`concourse/git-resource`](https://github.com/concourse/git-resource)).

## install

### macOS

```sh
brew install frantjc/tap/forge
```

## usage

### GitHub Actions

For GitHub Actions, Forge will try to source the GitHub Actions variables from the working directory's Git configuration as well as GitHub's [default environment variables](https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables).

```sh
forge use actions/setup-go@v3 -w go-version=1.19
```

Forge mounts the current working directory to the Action's `GITHUB_WORKSPACE` as well as cache directories respecting the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) to the Action's `RUNNER_TOOLCACHE` and `RUNNER_TEMP`.

That is to say, after running the above command, `go` should be installed to `XDG_CACHE_HOME/.forge/runner/toolcache`.

You can also use local GitHub Actions by starting the reference with `"/"` or `"./"` to signify that it is an absolute or relative local filepath, respectively.

```sh
forge use ./testdata/actions/mock/docker
```

### Concourse Resources

For Concourse Resources, Forge will source `resource_types` and `resources` from the working directory's [`forge.yml`](forge.yml) (overridable with `-c`). This schema is conveniently compatible with Concourse's [pipeline](https://concourse-ci.org/pipelines.html) schema.

```sh
forge get mock -V version=v0.0.0
```

## why?

Automation begins with a shell script that executes a bunch of CLI commands often to test, build and publish some code. The next step is to set up some CI system that executes that script, for example, on every commit to a repository's `main` branch. Such CI systems tend to identify that all of the scripts that they are executing do a lot of the same things--checkout a Git repository, setup a tool and so on.

In an effort to make their platform easier to use and to refactor the shared functionality out of all of the aforementioned scripts, CI systems in the past have introduced reusable "plugins"/"Actions"/"Resources" which take minimal configuration to do a complex task. GitHub Actions' [`actions/checkout`](https://github.com/actions/checkout), for example, takes one short line of code to invoke and accepts much optional configuration to modify its functionality to fulfill many related use cases.

Unfortunately, using such powerful plugins outside of the the system they were built for can be wildly difficult. This makes debugging the use of these plugins painful due to the long feedback loops. It also makes migrating from one CI system to another treacherous, having to replace uses of one system's plugins with another's.

Forge aims to remedy this.

## developing

- `make` is recommended - version 3.81 is tested
- `golang` is _required_ - version 1.18.x or above is required for [generics](https://go.dev/doc/tutorial/generics)
- `docker` is _required_ - version 20.10.x is tested
- [`upx`](https://github.com/upx/upx) is _required for_ compressing [`shim`](internal/cmd/shim)
