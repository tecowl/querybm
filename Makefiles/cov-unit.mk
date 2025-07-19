# Directory path from COVERAGES_DIR is used in test in sub directory. So COVERAGE_DIR must be an absolute path.
COVERAGES_DIR=$(CURDIR)/coverages
$(COVERAGES_DIR):
	mkdir -p $(COVERAGES_DIR)

UNIT_COVERAGE_DIR=$(COVERAGES_DIR)/unit
$(UNIT_COVERAGE_DIR):
	mkdir -p $(UNIT_COVERAGE_DIR)

# See https://app.codecov.io/github/akm/go-requestid/new
.PHONY: test-unit-with-coverage
test-unit-with-coverage: $(UNIT_COVERAGE_DIR)
	go test -cover ./... -args -test.gocoverdir="$(UNIT_COVERAGE_DIR)"
