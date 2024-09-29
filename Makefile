.DEFAULT_GOAL := help

SHELL := bash
PATH := $(CURDIR)/.dev/gopath/bin:$(PATH)

# Load .env file if it exists.
ifneq (,$(wildcard ./.env))
  include .env
  export
endif

.PHONY: help
help: ## Show help
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[/0-9a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'


# --------------------------------------------------------------------------------------
# Development environment
# --------------------------------------------------------------------------------------
.PHONY: setup
setup: ## Setup development environment
	@echo "==> Setting Go tools up..."
	@mkdir -p .dev/gopath
	@export GOPATH=$(CURDIR)/.dev/gopath && \
		go install honnef.co/go/tools/cmd/staticcheck@latest && \
		go install github.com/golang/mock/mockgen@latest && \
		go install github.com/axw/gocov/gocov@latest && \
		go install github.com/matm/gocov-html/cmd/gocov-html@latest
	@export GOPATH=$(CURDIR)/.dev/gopath && go clean -modcache && rm -rf $(CURDIR)/.dev/gopath/pkg

.PHONY: clean
clean: ## Clean up development environment
	@rm -rf .dev

# --------------------------------------------------------------------------------------
# Testing, Formatting and etc.
# --------------------------------------------------------------------------------------
.PHONY: format
format: ## Format go code
	@go fmt ./...

.PHONY: lint
lint: ## Lint source code
	@go vet ./... ; staticcheck ./...

.PHONY: test
test: ## Run tests
	@go test -race -timeout 30m ./...

.PHONY: test/verbos
test/verbose: ## Run tests with verbose outputting
	@go test -race -timeout 30m -v ./...

.PHONY: test/cover
test/cover: ## Run tests with coverage
	@echo "==> Run tests with coverage report..."
	@mkdir -p $(CURDIR)/.dev/coverage
	@go test -coverpkg=./... -coverprofile=$(CURDIR)/.dev/coverage/coverage.out ./...
	@gocov convert $(CURDIR)/.dev/coverage/coverage.out | gocov-html > $(CURDIR)/.dev/coverage/coverage.html
	@echo "==> Open $(CURDIR)/.dev/coverage/coverage.html to see the coverage report."

.PHONY: open/coverage
open/coverage: ## Open coverage report
	@open $(CURDIR)/.dev/coverage/coverage.html

# --------------------------------------------------------------------------------------
# Run example app
# --------------------------------------------------------------------------------------
.PHONY: start
start: ## Run example app
	@cd examples && go run .

.PHONY: go-mod-tidy
go-mod-tidy: ## Run go mod tidy
	@go mod tidy

