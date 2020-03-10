ifeq ($(origin GITROOT), undefined)
GITROOT := $(shell git rev-parse --show-toplevel)
endif

# Commands
GO         ?= go
GOLINT     ?= golint
GOIMPORTS  ?= goimports

GO_IMPORT_PATH = github.com/popodidi/log

.PHONY: dep precommit format lint vet test tidy

dep:
	$(GO) get -u golang.org/x/lint/golint

precommit: tidy format lint vet test

lint:
	$(GOLINT) $(GITROOT)/...

format:
	@find . -name "*.go" | xargs $(GOIMPORTS) -w -local $(GO_IMPORT_PATH)

vet:
	$(GO) vet $(GITROOT)/...

test:
	$(GO) test $(GITROOT)/...

tidy:
	$(GO) mod tidy
