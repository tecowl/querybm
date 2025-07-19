GO_TEST_OPTIONS?=

.PHONY: test
test: test-unit

.PHONY: test-unit
test-unit:
	go test $(GO_TEST_OPTIONS) ./...
