# Go parameters
PROJECT_NAME := $(shell echo $${PWD\#\#*/})
PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

all: lint vet install tags

install: ## Run install
	@go install && echo Installed `date` && echo

lint: ## Run lint
	@golint -set_exit_status ${PKG_LIST}

vet: ## Run go vet
	@go vet ./...

check: ## Run gosimple and staticcheck
	@gosimple && staticcheck

test: ## Run unittests
	@go test -short ${PKG_LIST}

race: ## Run data race detector
	@go test -race -short ${PKG_LIST}

build: ## Build the binary file
	@go build -i -v

clean: ## Remove previous build
	@go clean ./...

watch:
	@echo Watching for changes...
	@fswatch -or . -e ".*" -i "\\.go$$" | xargs -n1 -I{} make

tags:
	@gotags -R *.go . > tags

linux:
	@env GOOS=linux GOARCH=amd64 go build -v -o ./build/$(PROJECT_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: all install lint vet check test race build clean watch tags linux help
