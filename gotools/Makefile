GOTOOLS = golint goimports
GOTOOLS_BIN = $(patsubst %,$(GOPATH)/bin/%, $(GOTOOLS))

# go tool->path mapping
go.fqp.golint    := golang.org/x/lint/golint
go.fqp.goimports := golang.org/x/tools/cmd/goimports

all: $(GOTOOLS_BIN)

# Default rule for gotools uses the name->path map for a generic 'go get' style build
gotool.%:
	$(eval TOOL = ${subst gotool.,,${@}})
	@echo "Building $(TOOL)"
	go get ${go.fqp.${TOOL}}

$(GOPATH)/bin/%:
	$(eval TOOL = ${subst $(GOPATH)/bin/,,${@}})
	@$(MAKE) gotool.$(TOOL)

.PHONY: clean
clean:
	-@rm -f $(GOTOOLS_BIN) ||:
