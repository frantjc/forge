GO = go
GOLANGCI-LINT = golangci-lint
BUF = buf
INSTALL = sudo install

BIN = /usr/local/bin

VERSION ?= 0.0.0
PRERELEASE ?=

install: build
	@$(INSTALL) ./bin/4ge $(BIN)

build:
	@$(GO) $@ \
		-ldflags " \
			-s -w \
			-X github.com/frantjc/forge.Version=$(VERSION) \
			-X github.com/frantjc/forge.Prerelease=$(PRERELEASE) \
		" \
		-o ./bin \
		./cmd/4ge

protos:
	@$(BUF) format -w
	@$(BUF) generate .

fmt generate test:
	@$(GO) $@ ./...

download vendor:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run --fix

proto: protos
buf: proto
gen: generate
dl: download
ven: vendor
format: fmt
	@$(BUF) format -w

.PHONY: install build protos fmt generate test download vendor lint proto buf gen dl ven format
