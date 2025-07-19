.PHONY: default
default: build lint test

GOBIN=$(shell go env GOBIN)

include ./Makefiles/build.mk
include ./Makefiles/lint.mk
include ./Makefiles/godoc.mk
include ./Makefiles/test.mk
include ./Makefiles/cov-unit.mk
include ./Makefiles/cov-integration.mk
