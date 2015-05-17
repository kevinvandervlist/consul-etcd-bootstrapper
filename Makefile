DEPS = $(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
PACKAGES = $(shell go list ./...)

all: deps format test

deps:
	@echo "--> Installing build dependencies"
	@go get -d -v ./... $(DEPS)

updatedeps: deps
	@echo "--> Updating build dependencies"
	@go get -d -f -u ./... $(DEPS)

test: deps
	@go test $(PACKAGES)

format: deps
	@echo "--> Running go fmt"
	@go fmt $(PACKAGES)

.PHONY: all deps test
