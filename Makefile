GO = go

fmt generate test:
	@$(GO) $@ ./...

download tidy vendor:
	@$(GO) mod $@

.PHONY: fmt generate test download tidy vendor
