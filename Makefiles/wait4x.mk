WAIT4X_CLI_VERSION?=latest
WAIT4X_CLI_MODULE=wait4x.dev/v3/cmd/wait4x
WAIT4X_CLI=$(GOBIN)/wait4x
$(WAIT4X_CLI):
	$(MAKE) wait4x-cli-install
wait4x-cli-install:
	go install $(WAIT4X_CLI_MODULE)@$(WAIT4X_CLI_VERSION)
