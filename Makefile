GO = go
GIT = git
GOLANGCI-LINT = golangci-lint
INSTALL = sudo install
GORELEASER = goreleaser
UPX = upx
YARN = yarn

BIN = /usr/local/bin

GOOS = $(shell $(GO) env GOOS)
GOARCH = $(shell $(GO) env GOARCH)

SEMVER ?= 0.6.0

.DEFAULT: install

install: build
	@$(INSTALL) ./dist/forge_$(GOOS)_$(GOARCH)*/forge $(BIN)

build:
	@$(GORELEASER) release --snapshot --clean

.github/action:
	@cd .github/action && \
		$(YARN) build

generate:
	@$(GO) $@ ./...

fmt test:
	@$(GO) $@ ./...
	@cd .github/action && \
		$(YARN) $@

download:
	@$(GO) mod $@
	@cd .github/action && \
		$(YARN)

vendor verify:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run --fix

shim:
	@GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-s -w" -o ./internal/bin/shim ./internal/cmd/shim
	@$(UPX) --ultra-brute ./internal/bin/shim

clean:
	@rm -rf dist/ privileged version internal/bin/shim.*

release:
	@cd .github/action && \
		$(YARN) version --new-version $(SEMVER)
	@$(GIT) tag -a v$(SEMVER) -m v$(SEMVER)
	@$(GIT) push --follow-tags

action: .github/action
gen: generate
dl: download
ven: vendor
ver: verify
format: fmt
i: install

.PHONY: .github/action action i install build fmt generate test download vendor verify lint shim clean gen dl ven ver format release

-include docs/docs.mk
