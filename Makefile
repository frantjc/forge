GO = go
GOLANGCI-LINT = golangci-lint
BUF = buf
INSTALL = sudo install

BIN = /usr/local/bin

install: build
	@$(INSTALL) ./bin/4ge $(BIN)

build:
	@$(GO) $@ -o ./bin ./cmd/4ge

protos:
	@$(BUF) format -w
	@$(BUF) generate .

fmt generate test:
	@$(GO) $@ ./...

download tidy vendor:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run --fix

proto: protos
buf: proto
gen: generate
dl: download
td: tidy
ven: vendor
format: fmt
	@$(BUF) format -w

.PHONY: install build protos fmt generate test download tidy vendor lint proto buf gen dl td ven format
