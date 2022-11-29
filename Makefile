GO = go
GOLANGCI-LINT = golangci-lint
INSTALL = sudo install
GORELEASER = goreleaser
UPX = upx

BIN = /usr/local/bin

GOOS = $(shell $(GO) env GOOS)
GOARCH = $(shell $(GO) env GOARCH)

-include docs/docs.mk

.DEFAULT: install

install: build
	@$(INSTALL) ./dist/forge_$(GOOS)_$(GOARCH)*/forge $(BIN)

build:
	@$(GORELEASER) release --snapshot --rm-dist

fmt generate test:
	@$(GO) $@ ./...

download vendor verify:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run --fix

shim:
	@GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-s -w" -o ./internal/bin/shim ./internal/cmd/shim
	@$(UPX) --ultra-brute ./internal/bin/shim

clean:
	@rm -rf dist/ privileged version

gen: generate
dl: download
ven: vendor
ver: verify
format: fmt

.PHONY: install build fmt generate test download vendor verify lint shim clean gen dl ven ver format
