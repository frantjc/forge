Have you ever had to swap to using a new CI system? Twice? Three times? Done with searching for a replacement for each Action, CloudBuilder or resource that you were using in your old one? Tired of waiting minutes for feedback on each iteration of figuring out the quirks of your new one?

Forge is here to help.

Forge is a library and CLI for running reusable steps from various proprietary CI systems using Docker. This, for example, makes the functionality provided to GitHub Actions easily consumable (or testable) by users of other CI systems.

Forge currently exposes running [GitHub Actions](github-actions.md) (e.g. [`actions/setup-go`](https://github.com/actions/setup-go)), [Concourse resources](concourse.md) (e.g. [`concourse/git-resource`](https://github.com/concourse/git-resource)) and [Google Cloudbuild steps](https://cloud.google.com/build/docs/configuring-builds/create-basic-configuration) (e.g. [`gcr.io/cloud-builders/docker`](https://cloud.google.com/build/docs/building/build-containers)) _and more_ in its [Go module](library.md).
