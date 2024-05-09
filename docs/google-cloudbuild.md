# Run Google CloudBuild steps

Forge can be used to run, test and debug [`Google CloudBuild steps`](https://cloud.google.com/build/docs/configuring-builds/create-basic-configuration) (e.g. [`gcr.io/cloud-builders/docker`](https://cloud.google.com/build/docs/building/build-containers)).

Forge is specifically focused on running individual CloudBuild _steps_ (e.g. `gcr.io/cloud-builders/docker`), not entire [build config files](https://cloud.google.com/build/docs/configuring-builds/create-basic-configuration). For example, Forge is **not** intended to run a build config file such as:

```yml
steps:
  - name: gcr.io/cloud-builders/docker
    args:
      - build
      - -t
      - 'us-central1-docker.pkg.dev/${PROJECT_ID}/my-docker-repo/my-image'
      - .
```

Rather, it is intended to run an individual step from within a build config file, like so:

```sh
forge cloudbuild gcr.io/cloud-builders/docker -- build -t 'us-central1-docker.pkg.dev/${PROJECT_ID}/my-docker-repo/my-image' .
```

> In the previous example, the `--` is important to signify to Forge that the rest of the arguments are meant to be passed to the underlying step, not parsed by Forge itself. Additionally, the `''` are important to keep your shell from doing the substitution before Forge can get ahold of it.

Forge mounts the current working directory to the step's as well as a directory respecting the XDG Base Directory Specification to the step's `/workspace`.

Forge will try to source the [default substitutions](https://cloud.google.com/build/docs/configuring-builds/substitute-variable-values#using_default_substitutions) (e.g. for `PROJECT_ID` above) from the working directory's Git configuration as well as `~/.config/gcloud`.

Forge can modify the entrypoint to a Google CloudBuild step as well, like so:

```sh
forge cloudbuild --entrypoint bash gcr.io/cloud-builders/docker -- -c "docker build -t 'us-central1-docker.pkg.dev/${PROJECT_ID}/my-docker-repo/my-image' ."
```

> In the previous example, `""` are important to pass the entire `docker` command as the value to `bash`'s `-c` flag.

For additional debugging, you can attach to the container running the step to snoop around:

```sh
forge cloudbuild --attach gcr.io/cloud-builders/docker -- build -t 'us-central1-docker.pkg.dev/${PROJECT_ID}/my-docker-repo/my-image' .
```

> The step's image must have `bash` or `sh` on its `PATH` for the attach to work.
