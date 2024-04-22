# Execute Concourse resources

Forge can be used to run, test and debug [Concourse resources](https://concourse-ci.org/resources.html) (e.g. [`concourse/git-resource`](https://github.com/concourse/git-resource)) or to utilize them within other CI systems.

Forge is specifically focused on running individual Concourse _resources_ (e.g. `concourse/git-resource`), not entire [Concourse pipelines](https://concourse-ci.org/pipelines.html). For example, Forge is **not** intended to run a pipeline such as:

```yml
# pipeline.yml
resources:
  - name: forge-git
    type: git
    source:
      uri: https://github.com/frantjc/forge
jobs:
  - name: example
    plan:
      - get: forge-git
        trigger: true
```

Rather, it is intended to run a single resource from within a pipeline, like so:

```sh
forge get forge-git --config pipeline.yml
```

Forge mounts the current working directory to the resource's. So, after the previous command is ran, if it succeeds, the working directory should have [`frantjc/forge`](https://github.com/frantjc/forge) cloned into it.

Forge can also execute a Concourse resource's `check` and `put` commands by simply replacing the `get` in the above command, as seen below:

```sh
forge check forge-git --config pipeline.yml
```

When executing Concourse resources, Forge will source `resource_types` and `resources` from the working directory's [`.forge.yml`](../.forge.yml) (overridable with `--config` as seen previously). This schema is conveniently compatible with Concourse's pipeline schema.

Just like Concourse itself, Forge ships with [some resource types builtin](../concourse/builtin.go) which can also be overridden. This is why the `git` resource type did not need explicitly defined in the above examples.

You can also attach to the container executing the resource to snoop around:

```sh
forge get --attach forge-git
```

> The resource's image must have `bash` or `sh` on its `PATH` for the attach to work.
