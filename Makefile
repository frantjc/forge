GO = go
GIT = git
GOLANGCI-LINT = golangci-lint
INSTALL = sudo install
GORELEASER = goreleaser
UPX = upx
YARN = yarn

BIN = /usr/local/bin

GOOS = $(shell $(GO) env GOOS)
GOARCH = $(shell $(GO) env GOARCH)

SEMVER ?= 0.10.0

.DEFAULT: install

install: build
	@$(INSTALL) ./dist/forge_$(GOOS)_$(GOARCH)*/forge $(BIN)

build:
	@$(GORELEASER) release --snapshot --clean

.github/action .github/action/:
	@cd .github/action && $(YARN) all

generate:
	@$(GO) $@ ./...

fmt test:
	@$(GO) $@ ./...

download:
	@$(GO) mod $@
	@cd .github/action && $(YARN)

vendor verify:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run --fix

internal/bin/shim_$(GOARCH):
	@GOOS=linux $(GO) build -ldflags "-s -w" -o $@ ./internal/cmd/shim
	@$(UPX) --ultra-brute $@
	@cat internal/bin/fs.go.tpl | sed -e "s|GOARCH|$(GOARCH)|g" > internal/bin/fs_$(GOARCH).go

clean:
	@rm -rf dist/ vendor/ privileged version internal/bin/shim*.*

MAJOR = $(word 1,$(subst ., ,$(SEMVER)))
MINOR = $(word 2,$(subst ., ,$(SEMVER)))

release:
	@cd .github/action && \
		$(YARN) version --new-version $(SEMVER)
	@$(GIT) push
	@$(GIT) tag -f v$(MAJOR)
	@$(GIT) tag -f v$(MAJOR).$(MINOR)
	@$(GIT) push --tags -f

action: .github/action
gen: generate
dl: download
ven: vendor
ver: verify
format: fmt
i: install
shim: shim_$(GOARCH)
shim_$(GOARCH): internal/bin/shim_$(GOARCH)

.PHONY: .github/action .github/action/ action i install build fmt generate test download vendor verify lint shim shim_$(GOARCH) internal/bin/shim_$(GOARCH) clean gen dl ven ver format release

-include docs/docs.mk
