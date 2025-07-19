# Directory path from COVERAGES_DIR is used in test in sub directory. So COVERAGE_DIR must be an absolute path.
COVERAGES_DIR=$(CURDIR)/coverages
$(COVERAGES_DIR):
	mkdir -p $(COVERAGES_DIR)

COVERAGE_INTEGRATED_DIR=$(COVERAGES_DIR)/integrated
$(COVERAGE_INTEGRATED_DIR):
	mkdir -p $(COVERAGE_INTEGRATED_DIR)

.PHONY: test-with-coverage
test-with-coverage: test-unit-with-coverage

COVERAGE_DIRS_CSV=$(UNIT_COVERAGE_DIR)

COVERAGE_PROFILE?=$(COVERAGES_DIR)/coverage.txt
$(COVERAGE_PROFILE): $(COVERAGE_INTEGRATED_DIR)
	$(MAKE) test-coverage-profile

.PHONY: test-coverage-profile
test-coverage-profile: $(COVERAGE_INTEGRATED_DIR)
	go tool covdata merge \
		-i $(COVERAGE_DIRS_CSV) \
		-o $(COVERAGE_INTEGRATED_DIR)
	go tool covdata percent -i=$(COVERAGE_INTEGRATED_DIR) -o $(COVERAGE_PROFILE)

COVERAGE_HTML?=$(COVERAGES_DIR)/coverage.html
$(COVERAGE_HTML): $(COVERAGE_PROFILE)
	go tool covdata html -i=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)

.PHONY: test-coverage
test-coverage: test-coverage-profile
	go tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	@command -v open && open $(COVERAGE_HTML) || echo "open $(COVERAGE_HTML)"

UNIT_COVERAGE_DIR=$(COVERAGES_DIR)/unit
$(UNIT_COVERAGE_DIR):
	mkdir -p $(UNIT_COVERAGE_DIR)

# See https://app.codecov.io/github/akm/go-requestid/new
.PHONY: test-unit-with-coverage
test-unit-with-coverage: $(UNIT_COVERAGE_DIR)
	go test -cover ./... -args -test.gocoverdir="$(UNIT_COVERAGE_DIR)"
