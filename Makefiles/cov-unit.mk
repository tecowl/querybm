# Directory path from COVERAGES_DIR is used in test in sub directory. So COVERAGE_DIR must be an absolute path.
COVERAGES_DIR=$(CURDIR)/coverages
$(COVERAGES_DIR):
	mkdir -p $(COVERAGES_DIR)

UNIT_COVERAGE_DIR=$(COVERAGES_DIR)/unit
$(UNIT_COVERAGE_DIR):
	mkdir -p $(UNIT_COVERAGE_DIR)

.PHONY: clean-unit-coverage
clean-unit-coverage:
	rm -rf $(UNIT_COVERAGE_DIR)

# COVERAGE_GO_PACKAGES_CSV is defined in Makefile

# See https://app.codecov.io/github/akm/go-requestid/new
.PHONY: test-unit-with-coverage
test-unit-with-coverage: clean-unit-coverage $(UNIT_COVERAGE_DIR)
	go test -cover -coverpkg=$(COVERAGE_GO_PACKAGES_CSV) ./... -args -test.gocoverdir="$(UNIT_COVERAGE_DIR)"
