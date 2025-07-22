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

.PHONY: linters-enabled
linters-enabled: $(GOLANGCI_LINT_CLI)
	@golangci-lint linters | awk '/^Enabled by your configuration linters:$$/{flag=1;next}/^Disabled by your configuration linters:$$/{flag=0}flag{print}'
