GO = go
GOLANGCI-LINT = golangci-lint

fmt generate test:
	@$(GO) $@ ./...

download tidy vendor:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run

.PHONY: fmt generate test download tidy vendor lint
