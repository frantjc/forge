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
YARN ?= yarn

GOARCH ?= $(shell $(GO) env GOARCH)

.PHONY: .github/actions/setup-forge/node_modules .github/actions/setup-forge/node_modules/
.github/actions/setup-forge/node_modules .github/actions/setup-forge/node_modules/:
	@cd .github/actions/setup-forge && $(YARN)

.PHONY: .github/actions/setup-forge .github/actions/setup-forge/
.github/actions/setup-forge .github/actions/setup-forge/: .github/actions/setup-forge/node_modules
	@cd $@ && $(YARN) all

.PHONY: fmt test generate
fmt test generate:
	@$(GO) $@ ./...

.PHONY: download tidy
download tidy:
	@$(GO) mod $@

.PHONY: lint
lint: golangci-lint
	@$(GOLANGCI_LINT) run --fix

.PHONY: release
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

.PHONY: .github/actions/setup-forge .github/actions/setup-forge/ action i install build fmt generate test download vendor verify lint shims shim shim_$(GOARCH) internal/bin/fs_$(GOARCH).go internal/bin/shim_$(GOARCH) clean gen dl ven ver format release

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

FORGE ?= $(LOCALBIN)/forge
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint

GOLANGCI_LINT_VERSION ?= v2.1.5

.PHONY: forge
forge: $(FORGE)
$(FORGE): $(LOCALBIN)
	@$(GO) build -o $@ ./cmd/forge

BIN ?= /usr/local/bin
INSTALL ?= sudo install

.PHONY: install
install: forge
	@$(INSTALL) $(FORGE) $(BIN)

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT)
$(GOLANGCI_LINT): $(LOCALBIN)
	@$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

define go-install-tool
@[ -f "$(1)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
} ;
endef

-include docs/gifs.mk
