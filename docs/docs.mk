GIFSICLE = gifsicle
TERMINALIZER = terminalizer

CWD = $(abspath $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST))))))
INTERMEDIATE = $(CWD)/tmp.gif

gif:
	@$(TERMINALIZER) render -o=$(INTERMEDIATE) $(CWD)/demo.yml
	@$(GIFSICLE) --optimize=3 --colors=256 -o=$(CWD)/demo.gif $(INTERMEDIATE)
	@rm $(INTERMEDIATE)
