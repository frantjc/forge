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

SEMVER ?= 0.15.0-alpha0

.DEFAULT: install

install: build
	@$(INSTALL) ./dist/forge_$(GOOS)_$(GOARCH)*/forge $(BIN)

build:
	@$(GORELEASER) release --snapshot --clean

.github/actions/setup-forge .github/actions/setup-forge/:
	@cd $@ && $(YARN) all

generate:
	@$(GO) $@ ./...

fmt test:
	@$(GO) $@ ./...

download:
	@$(GO) mod $@
	@cd .github/actions/setup-forge && $(YARN)

vendor verify:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run --fix

internal/bin/shim_$(GOARCH):
	@GOOS=linux GOARCH=$(GOARCH) CGO_ENABLED=0 $(GO) build -ldflags "-s -w -X github.com/frantjc/forge.VersionCore=$(SEMVER)" -o $@ ./internal/cmd/shim
	@$(UPX) --ultra-brute $@

internal/bin/fs_$(GOARCH).go:
	@cat internal/bin/fs.go.tpl | sed -e "s|GOARCH|$(GOARCH)|g" > $@

clean:
	@rm -rf dist/ rootfs/ vendor/ privileged version internal/bin/shim*.*

MAJOR = $(word 1,$(subst ., ,$(SEMVER)))
MINOR = $(word 2,$(subst ., ,$(SEMVER)))

release:
	@cd .github/actions/setup-forge && \
		$(YARN) version --new-version $(SEMVER)
	@$(GIT) push
	@$(GIT) tag -f v$(MAJOR)
	@$(GIT) tag -f v$(MAJOR).$(MINOR)
	@$(GIT) push --tags -f

action: .github/actions/setup-forge
gen: generate
dl: download
ven: vendor
ver: verify
format: fmt
i: install
shim: shim_$(GOARCH)
shim_$(GOARCH): internal/bin/shim_$(GOARCH) internal/bin/fs_$(GOARCH).go

.PHONY: .github/actions/setup-forge .github/actions/setup-forge/ action i install build fmt generate test download vendor verify lint shim shim_$(GOARCH) internal/bin/fs_$(GOARCH).go internal/bin/shim_$(GOARCH) clean gen dl ven ver format release

-include docs/docs.mk
