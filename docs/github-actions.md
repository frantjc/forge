# Use GitHub Actions

Forge can be used to run, test and debug [GitHub Actions](https://docs.github.com/en/actions/learn-github-actions/finding-and-customizing-actions) (e.g. [`actions/setup-go`](https://github.com/actions/setup-go)) or to utilize them within other CI systems.

<p align="center">
  <img src="https://raw.githubusercontent.com/frantjc/forge/main/docs/github-actions.gif">
</p>

Forge is specifically focused on running individual GitHub _Actions_ (e.g. [`actions/checkout`](https://github.com/actions/checkout)), not entire GitHub Actions' _workflows_. For example, Forge is **not** intended to run a workflow such as:

```yml
on: push
jobs:
  example:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/setup-go@v5
        with:
          go-version: 1.22
```

Rather, it is intended to run an Action from within a workflow, like so:

```sh
forge use actions/setup-go@v5 --with go-version=1.22
```

When running an Action, Forge mounts the current working directory to the Action's `GITHUB_WORKSPACE` as well as directories respecting the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) to the Action's `RUNNER_TOOLCACHE` and `RUNNER_TEMP`. So, after the above command is ran, if it succeeds, `go` should be installed somewhere in `XDG_CACHE_HOME/forge/runner/toolcache`. This can be found more easily by running:

```sh
forge cache toolcache
```

Some Actions rely heavily on some default variables provided by GitHub. For example, `actions/checkout` requires on the environment variable `GITHUB_REPOSITORY` to be set to know which repository it should checkout.

Forge does its best to source such variables from the working directory's Git configuration as well as GitHub's [default environment variables](https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables) in its own environment.

However, in the event that an action errors and reports that a variable that it relies on is not set, it's likely that Forge did not find a value for that variable. For example, `actions/checkout` errors and reports by saying: `"[error] context.repo requires a GITHUB_REPOSITORY environment variable like 'owner/repo'"`.

In such cases, you can provide the value as an environment variable:

```sh
GITHUB_REPOSITORY=frantjc/forge forge use actions/checkout@v4
```

Authentication to GitHub is provided to Forge in a similar way--through the `GITHUB_TOKEN` environment variable. The value for the environment variable should be a [personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens) and can be injected into the environment more safely via something like `~/.bash_profile` like so:

```sh
export GITHUB_TOKEN=yourtokenhere
```

Forge can also execute local GitHub Actions. This helps custom Action developers more quickly and easily test out their actions while developing them locally.

To signify to Forge that a GitHub Action can be found on the filesystem as opposed to in a remote GitHub repository, start the reference with `"/"` or `"."` for absolute or relative filepaths, respectively. For [example](https://github.com/frantjc/forge/blob/main/testdata/actions/docker/action.yml):

```sh
forge use ./testdata/actions/docker
```

> Local Actions cannot refer to files outside of the `action.yml`'s directory.

For additional assistance with debugging, you can attach to the container running the Action to snoop around as in this [example](https://github.com/frantjc/forge/blob/main/testdata/actions/dockerfile/action.yml):

```sh
forge use --attach ./testdata/actions/dockerfile
```

> If the Action runs using a custom image, that image must have `bash` or `sh` on its `PATH` for the attach to work.
