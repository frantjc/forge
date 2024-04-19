<p align="center">
  <img src="https://raw.githubusercontent.com/frantjc/forge/main/docs/demo.gif">
</p>

Forge is a library and CLI for running reusable steps from various proprietary CI systems using a pluggable container runtime. This, for example, makes the functionality provided to GitHub Actions easily consumable (or testable) by users of other CI systems.

Forge currently exposes running [GitHub Actions](https://docs.github.com/en/actions/learn-github-actions/finding-and-customizing-actions) (e.g. [`actions/setup-go`](https://github.com/actions/setup-go)), [Concourse Resources](https://concourse-ci.org/resources.html) (e.g. [`concourse/git-resource`](https://github.com/concourse/git-resource)) and [Google Cloudbuild Steps](https://cloud.google.com/build/docs/configuring-builds/create-basic-configuration) (e.g. [`gcr.io/cloud-builders/docker`](https://cloud.google.com/build/docs/building/build-containers)).
