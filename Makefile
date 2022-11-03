GO = go
GOLANGCI-LINT = golangci-lint
BUF = buf
INSTALL = sudo install
GORELEASER = goreleaser
UPX = upx

BIN = /usr/local/bin

VERSION ?= 0.0.0
PRERELEASE ?=

GOOS = $(shell $(GO) env GOOS)
GOARCH = $(shell $(GO) env GOARCH)

-include docs/docs.mk

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

shim:
	@GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-s -w" -o ./internal/bin/shim ./internal/cmd/shim
	@$(UPX) --ultra-brute ./internal/bin/shim

proto: protos
buf: proto
gen: generate
dl: download
ven: vendor
ver: verify
format: fmt
	@$(BUF) format -w

.PHONY: install build protos fmt generate test download vendor verify lint shim proto buf gen dl ven ver format
