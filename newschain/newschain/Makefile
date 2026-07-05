

```makefile id="mk7p3d"
# ==============================================================================
# NewsChain Makefile
# ==============================================================================

# ------------------------------------------------------------------------------
# Variables
# ------------------------------------------------------------------------------

APP_NAME := newschain
CMD := ./cmd/server
BUILD_DIR := build
GO := go

# ------------------------------------------------------------------------------
# Default Target
# ------------------------------------------------------------------------------

.DEFAULT_GOAL := help

# ------------------------------------------------------------------------------
# Help
# ------------------------------------------------------------------------------

.PHONY: help
help: ## Show available commands
	@echo ""
	@echo "NewsChain Development Commands"
	@echo "=============================="
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ------------------------------------------------------------------------------
# Development
# ------------------------------------------------------------------------------

.PHONY: run
run: ## Run the application
	$(GO) run $(CMD)

.PHONY: dev
dev: ## Run the application (alias for run)
	$(GO) run $(CMD)

.PHONY: tidy
tidy: ## Download and tidy Go modules
	$(GO) mod tidy

.PHONY: download
download: ## Download Go dependencies
	$(GO) mod download

.PHONY: verify
verify: ## Verify Go module dependencies
	$(GO) mod verify

# ------------------------------------------------------------------------------
# Build
# ------------------------------------------------------------------------------

.PHONY: build
build: ## Build the application
	mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) $(CMD)

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

.PHONY: rebuild
rebuild: clean build ## Clean and rebuild

# ------------------------------------------------------------------------------
# Formatting
# ------------------------------------------------------------------------------

.PHONY: fmt
fmt: ## Format Go code
	gofmt -w .

.PHONY: gofumpt
gofumpt: ## Format using gofumpt
	gofumpt -w .

# ------------------------------------------------------------------------------
# Linting
# ------------------------------------------------------------------------------

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

# ------------------------------------------------------------------------------
# Testing
# ------------------------------------------------------------------------------

.PHONY: test
test: ## Run tests
	$(GO) test ./...

.PHONY: testv
testv: ## Run tests with verbose output
	$(GO) test -v ./...

.PHONY: race
race: ## Run race detector
	$(GO) test -race ./...

.PHONY: coverage
coverage: ## Generate test coverage
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# ------------------------------------------------------------------------------
# Security
# ------------------------------------------------------------------------------

.PHONY: vuln
vuln: ## Run govulncheck
	govulncheck ./...

# ------------------------------------------------------------------------------
# Smart Contracts (Foundry)
# ------------------------------------------------------------------------------

.PHONY: forge-build
forge-build: ## Build Solidity contracts
	forge build

.PHONY: forge-test
forge-test: ## Run Solidity tests
	forge test

.PHONY: forge-clean
forge-clean: ## Clean Foundry build artifacts
	forge clean

.PHONY: forge-format
forge-format: ## Format Solidity files
	forge fmt

# ------------------------------------------------------------------------------
# Quality Checks
# ------------------------------------------------------------------------------

.PHONY: check
check: fmt vet lint test ## Run formatting, vet, lint, and tests

.PHONY: ci
ci: tidy verify fmt vet lint test build ## Run local CI pipeline

# ------------------------------------------------------------------------------
# Docker (optional)
# ------------------------------------------------------------------------------

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(APP_NAME):latest .

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run --rm -p 8080:8080 $(APP_NAME):latest

# ------------------------------------------------------------------------------
# Cleanup
# ------------------------------------------------------------------------------

.PHONY: distclean
distclean: clean forge-clean ## Remove build artifacts and contract outputs
```
