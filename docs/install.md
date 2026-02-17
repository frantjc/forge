# Install

Using Homebrew:

```sh
brew install --cask frantjc/tap/forge
```

From a [release](https://github.com/frantjc/forge/releases).

From source, after [installing Dagger](https://docs.dagger.io/getting-started/installation):

```sh
git clone https://github.com/frantjc/forge
cd forge
make install
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
