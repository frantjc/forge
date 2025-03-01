GIFSICLE ?= gifsicle
TERMINALIZER ?= terminalizer

CWD = $(abspath $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST))))))
INTERMEDIATE = $(CWD)/tmp.gif

.PHONY: github-actions.gif
github-actions.gif:
	@$(TERMINALIZER) render -o=$(INTERMEDIATE) $(CWD)/$@.yml
	@$(GIFSICLE) --optimize=3 --colors=256 -o=$(CWD)/$@ $(INTERMEDIATE)
	@rm $(INTERMEDIATE)
