ifeq ($(origin GITROOT), undefined)
GITROOT := $(shell git rev-parse --show-toplevel)
endif

# Commands
GO      ?= go
GOLINT  ?= golint

.PHONY: dep precommit lint vet test tidy

dep:
	$(GO) get -u golang.org/x/lint/golint

precommit: tidy lint vet test

lint:
	$(GOLINT) $(GITROOT)/...

vet:
	$(GO) vet $(GITROOT)/...

test:
	$(GO) test $(GITROOT)/...

tidy:
	$(GO) mod tidy
