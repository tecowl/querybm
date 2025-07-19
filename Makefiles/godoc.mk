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
