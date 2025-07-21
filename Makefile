.PHONY: default
default: build lint test

GOBIN=$(shell go env GOBIN)

include ./Makefiles/build.mk
include ./Makefiles/lint.mk
include ./Makefiles/godoc.mk
include ./Makefiles/test.mk
include ./Makefiles/cov-unit.mk
include ./Makefiles/cov-integration.mk

.PHONY: test
test: test-unit tests-test

tests-%:
	$(MAKE) -C tests $*

# COVERAGE_GO_PACKAGES_CSV is used in Makefiles/cov-unit.mk
COVERAGE_GO_PACKAGES_CSV=$(shell find . -type d | grep -v '.git' | grep -v tests | grep -v Makefiles | grep -v coverages | grep -v mysql | sed 's|^\.|github.com/tecowl/querybm|' | tr '\n' ',' | sed 's/,$$//')
.PHONY: coverage-go-packages
coverage-go-packages:
	@echo $(COVERAGE_GO_PACKAGES_CSV)
