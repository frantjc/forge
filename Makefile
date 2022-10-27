GO = go
GOLANGCI-LINT = golangci-lint
BUF = buf

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
	@$(GOLANGCI-LINT) run

proto: protos
buf: proto
gen: generate
dl: download
td: tidy
ven: vendor
format: fmt
	@$(BUF) format -w

.PHONY: build protos fmt generate test download tidy vendor lint proto buf gen dl td ven format
