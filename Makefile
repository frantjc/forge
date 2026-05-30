FORGE ?= $(LOCALBIN)/forge

.PHONY: forge
forge: $(FORGE)
$(FORGE): $(LOCALBIN)
	@dagger call binary export --path $(FORGE)

.PHONY: clean
clean:
	@rm -rf rootfs/ vendor/ privileged version

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

BIN ?= /usr/local/bin
INSTALL ?= sudo install

.PHONY: install
install: forge
	@$(INSTALL) $(FORGE) $(BIN)

-include docs/gifs.mk
