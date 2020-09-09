## Project Parameters
PKG_NAME = github.com/ewangplay/serval

# SUBDIRS are components that have their own Makefiles that we can invoke
SUBDIRS = gotools
SUBDIRS:=$(strip $(SUBDIRS))

all: srvc checks

checks: linter unit-test

srvc: serval

.PHONY: $(SUBDIRS)
$(SUBDIRS):
	cd $@ && $(MAKE)

.PHONY: serval
serval: build/bin/serval

unit-test:
	@sh ./scripts/gotest.sh

linter: gotools
	@echo "Running golint"
	@sh ./scripts/golint.sh
	@echo "Running goimports"
	@sh ./scripts/goimports.sh

build/bin:
	mkdir -p $@

build/bin/%:
	@mkdir -p $(@D)
	@echo "$@"
	GOBIN=$(abspath $(@D)) go install $(PKG_NAME)
	@echo "Binary available as $@"
	@touch $@

.PHONY: $(SUBDIRS:=-clean)
$(SUBDIRS:=-clean):
	cd $(patsubst %-clean,%,$@) && $(MAKE) clean

.PHONY: clean
clean:
	-@rm -rf build ||:

.PHONY: dist-clean
dist-clean: clean gotools-clean
	-@rm -rf /opt/serval/* ||:

.PHONY: install
install: serval
	@mkdir -p /opt/serval/bin
	@cp build/bin/serval /opt/serval/bin/
	@mkdir -p /opt/serval/etc 
	@cp sampleconfig/serval.yaml /opt/serval/etc/
	@cp -r sampleconfig/blockchain /opt/serval/etc/
	@mkdir -p /opt/serval/etc/blockchain/wallet
	@mkdir -p /opt/serval/log
	@echo "Serval installed successfully"
