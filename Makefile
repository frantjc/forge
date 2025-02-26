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

SEMVER ?= 1.0.0

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

internal/bin/shim_amd64:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -ldflags "-s -w" -o $@ ./internal/cmd/shim
	@$(UPX) --ultra-brute $@

internal/bin/fs_amd64.go:
	@cat internal/bin/fs.go.tpl | sed -e "s|GOARCH|amd64|g" > $@

internal/bin/shim_arm64:
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GO) build -ldflags "-s -w" -o $@ ./internal/cmd/shim
	@$(UPX) --ultra-brute $@

internal/bin/fs_arm64.go:
	@cat internal/bin/fs.go.tpl | sed -e "s|GOARCH|arm64|g" > $@

clean:
	@rm -rf dist/ rootfs/ vendor/ privileged version internal/bin/shim*.*

ifeq (,$(findstring -,$(SEMVER)))
MAJOR = $(word 1,$(subst ., ,$(SEMVER)))
MINOR = $(word 2,$(subst ., ,$(SEMVER)))

release:
	@cd .github/actions/setup-forge && \
		$(YARN) version --new-version $(SEMVER)
	@$(GIT) push
	@$(GIT) tag -f v$(MAJOR)
	@$(GIT) tag -f v$(MAJOR).$(MINOR)
	@$(GIT) push --tags -f
else
release:
	@cd .github/actions/setup-forge && \
		$(YARN) version --new-version $(SEMVER)
	@$(GIT) push
	@$(GIT) push --tags
endif

action: .github/actions/setup-forge
gen: generate
dl: download
ven: vendor
ver: verify
format: fmt
i: install
shims: shim_amd64 shim_arm64
shim: shim_$(GOARCH)
shim_amd64: internal/bin/shim_amd64 internal/bin/fs_amd64.go
shim_arm64: internal/bin/shim_arm64 internal/bin/fs_arm64.go

.PHONY: .github/actions/setup-forge .github/actions/setup-forge/ action i install build fmt generate test download vendor verify lint shims shim shim_$(GOARCH) internal/bin/fs_$(GOARCH).go internal/bin/shim_$(GOARCH) clean gen dl ven ver format release

-include docs/gifs.mk
