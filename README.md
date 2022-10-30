# forge [![CI](https://github.com/frantjc/forge/actions/workflows/push.yml/badge.svg?branch=main&event=push)](https://github.com/frantjc/forge/actions)

Library and CLI for running reusable steps from various proprietary CI systems using a pluggable container runtime. This, for example, makes the functionality provided to GitHub Actions easily consumable (or testable) by users of other CI systems.

Currently exposes running [GitHub Actions](https://docs.github.com/en/actions/learn-github-actions/finding-and-customizing-actions) (e.g. [`actions/checkout`](https://github.com/actions/checkout)) and [Concourse Resources](https://concourse-ci.org/resources.html) (e.g. [`concourse/git-resource`](https://github.com/concourse/git-resource))

## usage

### GitHub Actions

For GitHub Actions, Forge will try to source the GitHub Actions variables from the working directory's Git configuration as well as [environment variables](https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables)

```sh
4ge use actions/setup-go@v3 -w go-version=1.19
```

Forge mounts the current working directory to the Action's `GITHUB_WORKSPACE` as well as cache directories respecting the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) to the Action's `RUNNER_TOOLCACHE` and `RUNNER_TEMP`.

That is to say, after running the above command, `go` should be installed to `${XDG_CACHE_HOME}/.forge/runner/toolcache`.

### Concourse Resources

For Concourse Resources, Forge will source `resource_types` and `resources` from the working directory's [`forge.json`](forge.json) file.

```sh
4ge get mock -w version=v0.0.0
```
