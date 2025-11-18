SHELL := /bin/bash
.DEFAULT_GOAL := help

COMPOSE ?= docker compose
GOLANGCI_LINT ?= golangci-lint
TFPLUGINDOCS ?= $(shell go env GOPATH)/bin/tfplugindocs
SHELLCHECK ?= shellcheck

.PHONY: help deps tidy fmt build test lint lint-shell docs integration docker-build clean validate fuzz-quick

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

lint-shell: ## Lint shell scripts with shellcheck
	@if ! command -v $(SHELLCHECK) >/dev/null 2>&1; then \
		echo "⚠️  shellcheck not found. Install from: https://github.com/koalaman/shellcheck#installing" >&2; \
		echo "Skipping shell linting..." >&2; \
	else \
		echo "Running shellcheck on scripts/*.sh..."; \
		$(SHELLCHECK) scripts/*.sh || exit 1; \
		echo "✓ Shell scripts passed shellcheck"; \
	fi

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
	$(MAKE) lint-shell
	go test ./...
	$(MAKE) docs
	go run ./scripts/check-function-coverage.go examples integration
	go run ./scripts/check-fuzz-coverage.go

fuzz-quick: ## Run short fuzz sessions for validators (1m per package)
	@echo "Running short fuzz sessions for internal/validators..."
	@set -euo pipefail; \
	for pkg in $$(go list ./internal/validators); do \
	  echo "==> $$pkg"; \
	  go test $$pkg -run Fuzz -fuzz Fuzz -fuzztime=1m || true; \
	done

coverage: ## Generate and display test coverage report
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./internal/functions ./internal/validators
	@echo ""
	@echo "Coverage by package:"
	@go tool cover -func=coverage.out | grep -E '^(github|total:)' | tail -3
	@echo ""
	@echo "Detailed report saved to coverage.out"
	@echo "View HTML report with: go tool cover -html=coverage.out"

coverage-html: coverage ## Generate and open HTML coverage report
	@echo "Opening coverage report in browser..."
	@go tool cover -html=coverage.out

pre-push: ## Complete pre-push checklist (format, test, lint, coverage, docs)
	@echo "========================================"
	@echo "Running pre-push checks..."
	@echo "========================================"
	@echo ""
	@echo "[1/6] Formatting code..."
	@$(MAKE) fmt >/dev/null 2>&1 || (echo "❌ Format failed" && exit 1)
	@echo "✓ Format passed"
	@echo ""
	@echo "[2/6] Running tests..."
	@$(MAKE) test || (echo "❌ Tests failed" && exit 1)
	@echo "✓ Tests passed"
	@echo ""
	@echo "[3/6] Linting code..."
	@$(MAKE) lint || (echo "❌ Lint failed" && exit 1)
	@$(MAKE) lint-shell || (echo "❌ Shell lint failed" && exit 1)
	@echo "✓ Lint passed"
	@echo ""
	@echo "[4/6] Checking coverage..."
	@$(MAKE) coverage | tail -4
	@echo ""
	@echo "[5/6] Generating docs..."
	@$(MAKE) docs >/dev/null 2>&1 || (echo "❌ Docs generation failed" && exit 1)
	@echo "✓ Docs generated"
	@echo ""
	@echo "[6/6] Running validation checks..."
	@go run ./scripts/check-function-coverage.go examples integration || (echo "❌ Function coverage check failed" && exit 1)
	@go run ./scripts/check-fuzz-coverage.go || (echo "❌ Fuzz coverage check failed" && exit 1)
	@echo "✓ Validation checks passed"
	@echo ""
	@echo "========================================"
	@echo "✅ All pre-push checks passed!"
	@echo "========================================"
