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
	@$(GORELEASER) release --snapshot --rm-dist

action:
	@cd action/ && \
		$(YARN) build

generate:
	@$(GO) $@ ./...

fmt test:
	@$(GO) $@ ./...
	@cd action/ && \
		$(YARN) $@

download vendor verify:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run --fix

shim:
	@GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-s -w" -o ./internal/bin/shim ./internal/cmd/shim
	@$(UPX) --ultra-brute ./internal/bin/shim

clean:
	@rm -rf dist/ privileged version internal/bin/shim.*

release:
	@cd action/ && \
		$(YARN) version --new-version $(SEMVER)
	@$(GIT) tag -a v$(SEMVER) -m v$(SEMVER)
	@$(GIT) push --follow-tags

gen: generate
dl: download
ven: vendor
ver: verify
format: fmt
i: install

.PHONY: action i install build fmt generate test download vendor verify lint shim clean gen dl ven ver format release

-include docs/docs.mk
