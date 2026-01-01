# Install

From a [release](https://github.com/frantjc/forge/releases).

Using `brew`:

```sh
brew install --cask frantjc/tap/forge
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
- uses: frantjc/actions/setup-tool@v1
  with:
    repo: frantjc/forge
```

As a library:

```sh
go get -u github.com/frantjc/forge
```
