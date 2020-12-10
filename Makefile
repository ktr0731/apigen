SHELL := /bin/bash

export GOBIN := $(PWD)/_tools
export PATH := $(GOBIN):$(PATH)

.PHONY: tools
tools:
	@cat tools/tools.go | grep -E '^\s*_\s.*' | awk '{ print $$2 }' | xargs go install

.PHONY: build
build:
	go build ./...

.PHONY: test
test: format gotest

.PHONY: format
format:
	go mod tidy

.PHONY: gotest
gotest: lint
	go test -race ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY:
release:
	goreleaser --rm-dist
