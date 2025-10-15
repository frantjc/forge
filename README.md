# forge [![CI](https://github.com/frantjc/forge/actions/workflows/ci.yml/badge.svg?branch=main&event=push)](https://github.com/frantjc/forge/actions) [![godoc](https://pkg.go.dev/badge/github.com/frantjc/forge.svg)](https://pkg.go.dev/github.com/frantjc/forge) [![goreportcard](https://goreportcard.com/badge/github.com/frantjc/forge)](https://goreportcard.com/report/github.com/frantjc/forge)

Forge is a Dagger module for running reusable steps from various proprietary CI systems. This, for example, makes the functionality provided to GitHub Actions easily consumable by users of Dagger. For example, using [actions/setup-go](https://github.com/actions/setup-go):

```sh
dagger -m github.com/frantjc/forge -c 'use actions/setup-go@v5 | withInput go-version 1.25 | post | combinedOutput'
```
