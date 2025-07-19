.PHONY: default
default: build lint test

.PHONY: build
build:
	go build ./...

GOBIN=$(shell go env GOBIN)
GOLANGCI_LINT_CLI_VERSION?=latest
GOLANGCI_LINT_CLI_MODULE=github.com/golangci/golangci-lint/cmd/golangci-lint
GOLANGCI_LINT_CLI=$(GOBIN)/golangci-lint
$(GOLANGCI_LINT_CLI):
	$(MAKE) golangci-lint-cli-install
golangci-lint-cli-install:
	go install $(GOLANGCI_LINT_CLI_MODULE)@$(GOLANGCI_LINT_CLI_VERSION)

.PHONY: lint
lint: $(GOLANGCI_LINT_CLI)
	golangci-lint run


GODOC_CLI_VERSION=latest
GODOC_CLI_MODULE=golang.org/x/tools/cmd/godoc
GODOC_CLI=$(GOBIN)/godoc
$(GODOC_CLI):
	$(MAKE) godoc-cli-install
godoc-cli-install:
	go install $(GODOC_CLI_MODULE)@$(GODOC_CLI_VERSION)

.PHONY: godoc
godoc: $(GODOC_CLI)
	@echo "Open http://localhost:6060/pkg/github.com/akm/sql-slog"
	godoc -http=:6060

GO_TEST_OPTIONS?=

.PHONY: test
test: test-unit

.PHONY: test-unit
test-unit:
	go test $(GO_TEST_OPTIONS) ./...
