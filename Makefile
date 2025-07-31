.PHONY: build run clean test lint fmt help

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=techircd
CMD_PATH=.
VERSION?=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

## Build Commands

# Build the IRC server
build: ## Build the binary
	go build $(LDFLAGS) -o $(BINARY_NAME) $(CMD_PATH)

# Build for different platforms
build-all: ## Build for multiple platforms
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 $(CMD_PATH)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe $(CMD_PATH)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 $(CMD_PATH)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 $(CMD_PATH)

## Development Commands

# Run the server
run: build ## Build and run the server
	./$(BINARY_NAME)

# Run with custom config
run-config: build ## Run with custom config file
	./$(BINARY_NAME) -config configs/config.dev.json

# Install dependencies
deps: ## Download and install dependencies
	go mod download
	go mod tidy

# Format code
fmt: ## Format Go code
	go fmt ./...
	goimports -w -local github.com/ComputerTech312/TechIRCd .

# Lint code
lint: ## Run linters
	golangci-lint run

# Fix linting issues
lint-fix: ## Fix auto-fixable linting issues
	golangci-lint run --fix

## Testing Commands

# Run all tests
test: ## Run all tests
	go test -v -race ./...

# Run tests with coverage
test-coverage: ## Run tests with coverage report
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
benchmark: ## Run benchmark tests
	go test -bench=. -benchmem ./...

# Test with the test client
test-client: build ## Test with the simple client
	go run tests/test_client.go

## Release Commands

# Create a release build
release: clean ## Create optimized release build
	CGO_ENABLED=0 go build $(LDFLAGS) -a -installsuffix cgo -o $(BINARY_NAME) $(CMD_PATH)

# Tag a new version
tag: ## Tag a new version (usage: make tag VERSION=v1.0.1)
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)

## Utility Commands

# Clean build artifacts
clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)*
	rm -f coverage.out coverage.html

# Show git status and recent commits
status: ## Show git status and recent commits
	@echo "Git Status:"
	@git status --short
	@echo "\nRecent Commits:"
	@git log --oneline -10

# Install development tools
install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Generate documentation
docs: ## Generate documentation
	go doc -all > docs/API.md

# Show help
help: ## Show this help message
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	go fmt ./...

# Run with race detection
run-race: build
	go run -race *.go
