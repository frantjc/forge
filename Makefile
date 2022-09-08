GO = go
GOLANGCI-LINT = golangci-lint
BUF = buf

proto:
	@$(BUF) format -w
	@$(BUF) generate .

fmt generate test:
	@$(GO) $@ ./...

download tidy vendor:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run

.PHONY: fmt generate test download tidy vendor lint proto

