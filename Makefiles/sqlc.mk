SQLC_CLI_VERSION?=latest
SQLC_CLI_MODULE=github.com/sqlc-dev/sqlc/cmd/sqlc
SQLC_CLI=$(GOBIN)/sqlc
$(SQLC_CLI):
	$(MAKE) sqlc-cli-install
sqlc-cli-install:
	go install $(SQLC_CLI_MODULE)@$(SQLC_CLI_VERSION)
