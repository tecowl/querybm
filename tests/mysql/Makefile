.PHONY: default
default: build lint test

GOBIN=$(shell go env GOBIN)

include ../../Makefiles/build.mk
include ../../Makefiles/lint.mk
include ../../Makefiles/test.mk
include ../../Makefiles/cov-unit.mk

ENVS=\
	MYSQL_DATABASE=querybm-test \
	MYSQL_PORT=3306

upd:
	$(ENVS) docker compose -f docker-compose.yml up -d

down:
	$(ENVS) docker compose -f docker-compose.yml down
