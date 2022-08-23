GO = go

tidy downlaod:
	@$(GO) mod $@

test generate fmt:
	@$(GO) $@ ./...
