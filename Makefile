GO = go
GOLANGCI-LINT = golangci-lint
BUF = buf
INSTALL = sudo install
GORELEASER = goreleaser

BIN = /usr/local/bin

VERSION ?= 0.0.0
PRERELEASE ?=

GOOS = $(shell $(GO) env GOOS)
GOARCH = $(shell $(GO) env GOARCH)

.DEFAULT: install

install: build
	@$(INSTALL) ./dist/forge_$(GOOS)_$(GOARCH)*/forge $(BIN)

build:
	@$(GORELEASER) release --snapshot --rm-dist

protos:
	@$(BUF) format -w
	@$(BUF) generate .

fmt generate test:
	@$(GO) $@ ./...

download vendor verify:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run --fix

proto: protos
buf: proto
gen: generate
dl: download
ven: vendor
ver: verify
format: fmt
	@$(BUF) format -w

.PHONY: install build protos fmt generate test download vendor verify lint proto buf gen dl ven ver format
