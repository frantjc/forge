# forge [![CI](https://github.com/frantjc/forge/actions/workflows/ci.yml/badge.svg?branch=main&event=push)](https://github.com/frantjc/forge/actions) [![godoc](https://pkg.go.dev/badge/github.com/frantjc/forge.svg)](https://pkg.go.dev/github.com/frantjc/forge) [![goreportcard](https://goreportcard.com/badge/github.com/frantjc/forge)](https://goreportcard.com/report/github.com/frantjc/forge) ![license](https://shields.io/github/license/frantjc/forge)

<p align="center">
  <img src="https://raw.githubusercontent.com/frantjc/forge/main/docs/demo.gif">
</p>

Forge is a library and CLI for running reusable steps from various proprietary CI systems using a pluggable container runtime. This, for example, makes the functionality provided to GitHub Actions easily consumable (or testable) by users of other CI systems.

Forge currently exposes running [GitHub Actions](https://docs.github.com/en/actions/learn-github-actions/finding-and-customizing-actions) (e.g. [`actions/setup-go`](https://github.com/actions/setup-go)), [Concourse Resources](https://concourse-ci.org/resources.html) (e.g. [`concourse/git-resource`](https://github.com/concourse/git-resource)) and [Google Cloudbuild Steps](https://cloud.google.com/build/docs/configuring-builds/create-basic-configuration) (e.g. [gcr.io/cloud-builders/docker](https://cloud.google.com/build/docs/building/build-containers)).

## install

From a [release](https://github.com/frantjc/forge/releases).

Using `brew`:

```sh
brew install frantjc/tap/forge
```

From source:

```sh
git clone https://github.com/frantjc/forge
cd forge
make
```

Using `go`:

```sh
go install github.com/frantjc/forge/cmd/forge
```

In GitHub Actions:

```yml
  - uses: frantjc/forge@v0
```

As a library:

```sh
go get -u github.com/frantjc/forge
```

## usage

### GitHub Actions

For GitHub Actions, Forge will try to source the GitHub Actions variables from the working directory's Git configuration as well as GitHub's [default environment variables](https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables).

```sh
forge use actions/setup-go@v3 -w go-version=1.22
```

Forge mounts the current working directory to the Action's `GITHUB_WORKSPACE` as well as cache directories respecting the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) to the Action's `RUNNER_TOOLCACHE` and `RUNNER_TEMP`.

That is to say, after running the above command, `go` should be installed to `XDG_CACHE_HOME/forge/runner/toolcache`.

You can also use local GitHub Actions by starting the reference with `"/"` or `"."` to signify that it is an absolute or relative local filepath, respectively.

```sh
forge use ./testdata/actions/docker
```

For additional debugging, you can attach to the container running the Action:

```sh
forge use -a ./testdata/actions/dockerfile
```

> If the Action runs using a custom image, that image must have `bash` or `sh` on its `PATH` for the attach to work.

> Local Actions cannot refer to files outside of the action metadata file's directory.

### Concourse Resources

For Concourse Resources, Forge will source `resource_types` and `resources` from the working directory's [`.forge.yml`](.forge.yml) (overridable with `-c`). This schema is conveniently compatible with [Concourse's pipeline schema](https://concourse-ci.org/pipelines.html).

> Just like Concourse itself, Forge ships with [some Resource Types builtin](concourse/builtin.go) that can be overridden.

```sh
forge get mock -v version=v0.0.0
```

You can also attach to the container executing the Resource to snoop around:

```sh
forge get -a mock -v version=v0.0.0
```

> The Resource's image must have `bash` or `sh` on its `PATH` for the attach to work.

### Google Cloudbuild Tasks

For Google Cloudbuild, Forge will try to source the default substitutions from the working directory's Git configuration as well as `~/.config/gcloud`.

```sh
forge cloudbuild gcr.io/cloud-builders/docker -- build -t 'gcr.io/${PROJECT_ID}/my-image:${SHORT_SHA}' .
```

For additional debugging, you can attach to the container running the Cloudbuild:

```sh
forge cloudbuild -a gcr.io/cloud-builders/docker -- build -t 'gcr.io/${PROJECT_ID}/my-image:${SHORT_SHA}' .
```

> The Cloudbuild's image must have `bash` or `sh` on its `PATH` for the attach to work.

## why?

Automation begins with a shell script that executes a bunch of CLI commands often to test, build and publish some code. The next step is to set up some continuous integration (CI) system that executes that script in response to some event such as a commit to a Git repository's `main` branch. Such CI systems tend to identify that all of the scripts that they are executing do a lot of the same things--checkout a Git repository, setup a tool and so on.

In an effort to make their platform easier to use and to refactor the shared functionality out of all of the aforementioned scripts, CI systems in the past have introduced reusable "plugins"/"Actions"/"Resources"/"Tasks"/"Orbs" which take minimal configuration to do a complex task. GitHub Actions' [`actions/checkout`](https://github.com/actions/checkout), for example, takes one short line of code to invoke and accepts a bunch of optional configuration to fulfill many related use cases.

Unfortunately, using such powerful plugins outside of the the system they were built for can be wildly difficult. This makes debugging the use of these plugins require long feedback loops. It also makes migrating from one CI system to another treacherous, having to replace uses of one system's plugins with another's.

Forge aims to remedy this.

## developing

- `git` is _required_
- `make` is _required_
- `go` 1.20 is _required_ for multi-error handling
- `docker` is _required_ to test as it is its only runtime
- [`upx`](https://github.com/upx/upx) is _required_ for compressing [`shim`](internal/cmd/shim/main.go)
- `node` 20 is _required_ for developing the [`action`](.github/action)
