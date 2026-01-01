ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

GO ?= go
GIT ?= git
UPX ?= upx

GOARCH ?= $(shell $(GO) env GOARCH)

.PHONY: fmt test generate
fmt test generate:
	@$(GO) $@ ./...

.PHONY: internal/bin/shim_amd64
internal/bin/shim_amd64:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -ldflags "-s -w" -o $@ ./internal/cmd/shim
	@$(UPX) --ultra-brute $@

.PHONY: internal/bin/fs_amd64.go
internal/bin/fs_amd64.go:
	@cat internal/bin/fs.go.tpl | sed -e "s|GOARCH|amd64|g" > $@

.PHONY: internal/bin/shim_arm64
internal/bin/shim_arm64:
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GO) build -ldflags "-s -w" -o $@ ./internal/cmd/shim
	@$(UPX) --ultra-brute $@

.PHONY: internal/bin/fs_arm64.go
internal/bin/fs_arm64.go:
	@cat internal/bin/fs.go.tpl | sed -e "s|GOARCH|arm64|g" > $@

.PHONY: clean
clean:
	@rm -rf dist/ rootfs/ vendor/ privileged version internal/bin/shim*.*

.PHONY: release
ifeq (,$(findstring -,$(SEMVER)))
MAJOR = $(word 1,$(subst ., ,$(SEMVER)))
MINOR = $(word 2,$(subst ., ,$(SEMVER)))

release:
	@$(GIT) tag $(SEMVER)
	@$(GIT) push
	@$(GIT) tag -f v$(MAJOR)
	@$(GIT) tag -f v$(MAJOR).$(MINOR)
	@$(GIT) push --tags -f
else
release:
	@$(GIT) tag $(SEMVER)
	@$(GIT) push
	@$(GIT) push --tags
endif

.PHONY: action gen dl ven ver format i shims shim shim shim_amd64 shim_arm64
action: .github/actions/setup-forge
dist: build
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

.PHONY: i install build fmt generate test download vendor verify lint shims shim shim_$(GOARCH) internal/bin/fs_$(GOARCH).go internal/bin/shim_$(GOARCH) clean gen dl ven ver format release

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

FORGE ?= $(LOCALBIN)/forge

.PHONY: forge
forge: $(FORGE)
$(FORGE): $(LOCALBIN)
	@$(GO) build -o $@ ./cmd/forge

BIN ?= /usr/local/bin
INSTALL ?= sudo install

.PHONY: install
install: forge
	@$(INSTALL) $(FORGE) $(BIN)

-include docs/gifs.mk
