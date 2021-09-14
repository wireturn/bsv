# Common makefile commands & variables between projects
include .make/common.mk

# Common Golang makefile commands & variables between projects
include .make/go.mk

## Not defined? Use default repo name which is the application
ifeq ($(REPO_NAME),)
	REPO_NAME="go-bsvrates"
endif

## Not defined? Use default repo owner
ifeq ($(REPO_OWNER),)
	REPO_OWNER="tonicpow"
endif

.PHONY: clean

all: ## Runs lint, test and vet
	@$(MAKE) test

clean: ## Remove previous builds and any test cache data
	@go clean -cache -testcache -i -r
	@test $(DISTRIBUTIONS_DIR)
	@if [ -d $(DISTRIBUTIONS_DIR) ]; then rm -r $(DISTRIBUTIONS_DIR); fi

release:: ## Runs common.release then runs godocs
	@$(MAKE) godocs

run-examples: ## Runs the basic example
	@go run examples/get_rates/get_rates.go
