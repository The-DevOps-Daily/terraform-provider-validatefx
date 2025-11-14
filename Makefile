SHELL := /bin/bash
.DEFAULT_GOAL := help

COMPOSE ?= docker compose
GOLANGCI_LINT ?= golangci-lint
TFPLUGINDOCS ?= $(shell go env GOPATH)/bin/tfplugindocs

.PHONY: help deps tidy fmt build test lint docs integration docker-build clean validate

help: ## Show available targets and short descriptions
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "%-18s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download Go module dependencies
	go mod download

tidy: ## Ensure go.mod and go.sum are in sync
	go mod tidy

fmt: ## Format the Go sources
	go fmt ./...

build: ## Compile the provider binary
	go build ./...

test: ## Run unit tests
	go test ./...

lint: ## Run golangci-lint (requires golangci-lint in PATH)
	@if ! command -v $(GOLANGCI_LINT) >/dev/null 2>&1; then \
		echo "golangci-lint not found. Install with:" >&2; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0" >&2; \
		exit 1; \
	fi
	$(GOLANGCI_LINT) run ./...

docs: ## Generate provider documentation using tfplugindocs
	@if [ ! -x "$(TFPLUGINDOCS)" ]; then \
		echo "tfplugindocs not found. Install with:" >&2; \
		echo "  go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.19.2" >&2; \
		exit 1; \
	fi
	"$(TFPLUGINDOCS)" generate
	go run ./scripts/update-readme-functions-table.go

integration: ## Execute Terraform integration scenario via Docker Compose
	@status=0; \
	$(COMPOSE) build terraform || status=$$?; \
	$(COMPOSE) run --rm terraform || status=$$?; \
	$(COMPOSE) down -v || true; \
	exit $$status


docker-build: ## Build the provider Docker image
	docker build -t local/terraform-provider-validatefx -f Dockerfile .

clean: ## Remove build artifacts and local Terraform state
	rm -rf bin
	$(COMPOSE) down -v 2>/dev/null || true
	rm -rf integration/.terraform
	rm -rf .terraform .terraform.lock.hcl

validate: ## Run local pre-flight checks before pushing
	go fmt ./...
	terraform fmt -recursive
	go mod tidy
	go vet ./...
	$(MAKE) lint
	go test ./...
	$(MAKE) docs
	go run ./scripts/check-function-coverage.go examples integration
