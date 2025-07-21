.PHONY: default
default: build lint test

GOBIN=$(shell go env GOBIN)

include ../../Makefiles/build.mk
include ../../Makefiles/lint.mk
include ../../Makefiles/test.mk
include ../../Makefiles/cov-unit.mk
include ../../Makefiles/wait4x.mk

.PHONY: test
test: test-unit

MYSQL_VERSION=9.3.0
MYSQL_DATABASE=querybm-test
MYSQL_PORT=3306

ENVS=\
	MYSQL_VERSION=$(MYSQL_VERSION) \
	MYSQL_DATABASE=$(MYSQL_DATABASE) \
	MYSQL_PORT=$(MYSQL_PORT)

.PHONY: up
up: container-up mysql-wait schema-load
.PHONY: down
down: container-down

.PHONY: container-up
container-up:
	$(ENVS) docker compose -f docker-compose.yml up -d

.PHONY: container-down
container-down:
	$(ENVS) docker compose -f docker-compose.yml down

.PHONY: mysql-wait
mysql-wait: $(WAIT4X_CLI)
	$(WAIT4X_CLI) mysql 'root@tcp(localhost:$(MYSQL_PORT))/$(MYSQL_DATABASE)'

.PHONY: schema-load
schema-load:
	docker run --rm -v ./schema.sql:/schema.sql \
		--network querybm-mysql-test_network1 \
		mysql:$(MYSQL_VERSION) \
		/bin/bash -c 'mysql -h mysql -P $(MYSQL_PORT) $(MYSQL_DATABASE) < /schema.sql'

.PHONY: dbconsole
dbconsole:
	docker run --rm -it \
		--network querybm-mysql-test_network1 \
		mysql:$(MYSQL_VERSION) \
		mysql -h mysql -P $(MYSQL_PORT) $(MYSQL_DATABASE)

include ../../Makefiles/sqlc.mk

.PHONY: generate
generate: $(SQLC_CLI)
	$(SQLC_CLI) generate

# COVERAGE_GO_PACKAGES_CSV is used in Makefiles/cov-unit.mk
COVERAGE_GO_PACKAGES_CSV=$(shell make -C ../.. coverage-go-packages --no-print-directory)
