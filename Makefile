GO = go

downlaod tidy vendor:
	@$(GO) mod $@

fmt generate test:
	@$(GO) $@ ./...
