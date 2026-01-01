FORGE ?= $(LOCALBIN)/forge

.PHONY: forge
forge: $(FORGE)
$(FORGE): $(LOCALBIN)
	@dagger call binary export --path $(FORGE)

.PHONY: .git/hooks .git/hooks/ .git/hooks/pre-commit
.git/hooks .git/hooks/ .git/hooks/pre-commit:
	@cp .githooks/* .git/hooks

.PHONY: internal/bin/fs_amd64.go
internal/bin/fs_amd64.go:
	@cat internal/bin/fs.go.tpl | sed -e "s|GOARCH|amd64|g" > $@

.PHONY: internal/bin/fs_arm64.go
internal/bin/fs_arm64.go:
	@cat internal/bin/fs.go.tpl | sed -e "s|GOARCH|arm64|g" > $@

.PHONY: clean
clean:
	@rm -rf rootfs/ vendor/ privileged version internal/bin/shim*.*

.PHONY: release
release:
	@git tag $(SEMVER)
	@git push
	@git push --tags

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

BIN ?= /usr/local/bin
INSTALL ?= sudo install

.PHONY: install
install: forge
	@$(INSTALL) $(FORGE) $(BIN)

-include docs/gifs.mk
