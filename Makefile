GO = go
GOLANGCI-LINT = golangci-lint
BUF = buf

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

.PHONY: protos fmt generate test download tidy vendor lint proto buf gen dl td ven format
