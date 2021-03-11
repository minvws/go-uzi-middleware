# @echo off
.SILENT:

# Make sure that globstar is active, this allows bash to use ./**/*.go
SHELL=/bin/bash -O globstar -c

# Default repository
REPO="github.com/minvws/go-uzimiddleware"

# Set environment variables from GO env if not set explicitly already
ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif
ifndef $(GOOS)
    GOOS=$(shell go env GOOS)
    export GOOS
endif
ifndef $(GOARCH)
    GOARCH=$(shell go env GOARCH)
    export GOARCH
endif

# paths to binaries
GO_STATCHECK_BIN = $(GOPATH)/bin/staticcheck
GO_INEFF_BIN = $(GOPATH)/bin/ineffassign
GO_GOCYCLO_BIN = $(GOPATH)/bin/gocyclo
GO_GOIMPORTS_BIN = $(GOPATH)/bin/goimports
GO_LINT_BIN = $(GOPATH)/bin/golint

# ---------------------------------------------------------------------------

# Downloads external tools as they are not available by default
get_test_tools: ## go get all build tools needed to testing
	GO111MODULE=off go get -u honnef.co/go/tools/cmd/staticcheck
	GO111MODULE=off go get -u github.com/gordonklaus/ineffassign
	GO111MODULE=off go get -u github.com/fzipp/gocyclo/cmd/gocyclo
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
	GO111MODULE=off go get -u golang.org/x/lint/golint

lint: ## Formats your go code to specified standards
	$(GO_GOIMPORTS_BIN) -w  --format-only .

## Runs all tests for the whole repository
test: test_goimports test_vet test_golint test_staticcheck test_ineffassign test_gocyclo test_unit

test_goimports:
	echo "Check goimports"
	$(GO_GOIMPORTS_BIN) -l .

test_vet:
	echo "Check vet"
	go vet ./...

test_staticcheck:
	echo "Check static"
	$(GO_STATCHECK_BIN) ./...

test_golint:
	echo "Check lint"
	$(GO_LINT_BIN) ./...

test_ineffassign:
	echo "Check ineffassign"
	$(GO_INEFF_BIN) ./...

test_gocyclo:
	echo "Check gocyclo"
	$(GO_GOCYCLO_BIN) -over 15 .

test_unit:
	echo "Check unit tests"
	go test ./...

all: test  ## Run tests

help: ## Display available commands
	echo "make commands"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
