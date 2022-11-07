# forge [![CI](https://github.com/frantjc/forge/actions/workflows/push.yml/badge.svg?branch=main&event=push)](https://github.com/frantjc/forge/actions)

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

That is to say, after running the above command, `go` should be installed to `$XDG_CACHE_HOME/.forge/runner/toolcache`.

You can also use local GitHub Actions by starting the reference with `"/"` or `"./"` to signify that it is an absolute or relative local filepath, respectively.

```sh
forge use ./testdata/actions/mock
```

### Concourse Resources

For Concourse Resources, Forge will source `resource_types` and `resources` from the working directory's [`forge.yml`](forge.yml).

```sh
forge get mock -i version=v0.0.0
```

## why?

Automation begins with a shell script that executes a bunch of CLI commands often to test, build and publish some code. The next step is to set up some CI system that executes that script, for example, on every commit to a repository's `main` branch. Such CI systems often identify that all of the scripts that they are executing often do a lot of the same things--checkout a repository, setup a tool and so on.

In an effort to make their platform easier to use and to refactor the shared functionality out of all of the aforementioend scripts, they often introduce reusible "plugins"/"Actions"/"Resources" which take minimal configuration to do some complex action. Take GitHub Actions' [`actions/checkout`](https://github.com/actions/checkout) for example. It takes one short line of code to invoke and accepts a lot of optional configuration to modify its functionality to fulfill many related use cases.

The shame is that, unfortunately, using such powerful plugins outside of the the system they were built for can be wildly difficult depending on the complexity of the protocol that the designers created (e.g. GitHub Actions are _much_ more complicated to execute than Concourse Resources). This makes debugging the use of these plugins painful, with long feedback loops. It also makes migrating from one CI system to another painful, having to replace all uses of one system's plugins with another's.

Forge aims to remedy this.
