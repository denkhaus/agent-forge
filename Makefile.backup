# Makefile for AgentForge

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Build parameters
BINARY_NAME=forge
BINARY_PATH=bin/$(BINARY_NAME)
BUILD_FLAGS=-buildvcs=false

# Version information
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Linker flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

.PHONY: all build clean test coverage lint fmt vet deps help install run-server run-chat

# Default target
all: clean lint test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_PATH) ./cmd
	@echo "Build complete: $(BINARY_PATH)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf bin/
	@rm -rf dist/
	@rm -f tmp_rovodev_*
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
coverage: test
	@echo "Generating coverage report..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# TUI testing targets
test-tui-headless:
	@echo "Running headless TUI tests..."
	go run ./cmd/test-tui -mode=headless -prompt=test-prompt -duration=5s -verbose

test-tui-automated:
	@echo "Running automated TUI tests..."
	go run ./cmd/test-tui -mode=automated -prompt=test-prompt -duration=10s -verbose

test-tui-demo:
	@echo "Running TUI demo..."
	go run ./cmd/test-tui -mode=demo -prompt=demo-prompt -duration=15s -verbose

test-tui-interactive:
	@echo "Running interactive TUI test..."
	go run ./cmd/test-tui -mode=interactive -prompt=interactive-test -duration=30s -verbose

test-tui-quick:
	@echo "Running quick TUI test..."
	./scripts/test-tui.sh quick

test-tui-all:
	@echo "Running all TUI tests..."
	./scripts/test-tui.sh all

# Run linter
lint:
	@echo "Running linter..."
	GOFLAGS="-buildvcs=false" $(GOLINT) run --config .golangci.yml --build-tags=""

# Fix linting issues
lint-fix:
	@echo "Fixing linting issues..."
	$(GOLINT) run --config .golangci.yml --fix

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	goimports -w -local github.com/denkhaus/agentforge .

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BINARY_PATH) $(GOPATH)/bin/$(BINARY_NAME)
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

# Run the server
run-server: build
	@echo "Starting AgentForge server..."
	./$(BINARY_PATH) server

# Run the chat interface
run-chat: build
	@echo "Starting AgentForge chat..."
	./$(BINARY_PATH) chat

# Run with development environment
run-dev: build
	@echo "Starting AgentForge in development mode..."
	@export LOG_LEVEL=debug && \
	export ENVIRONMENT=development && \
	./$(BINARY_PATH) server --log-level debug

# Component management commands
run-component-list: build
	@echo "Listing AgentForge components..."
	./$(BINARY_PATH) component list

run-component-new: build
	@echo "Creating new AgentForge component..."
	./$(BINARY_PATH) component new --type tool

# Check for security vulnerabilities
security:
	@echo "Checking for security vulnerabilities..."
	@which govulncheck >/dev/null 2>&1 || $(GOGET) golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

# Generate documentation
docs:
	@echo "Generating documentation..."
	@which godoc >/dev/null 2>&1 || $(GOGET) golang.org/x/tools/cmd/godoc@latest
	@echo "Run 'godoc -http=:6060' and visit http://localhost:6060/pkg/github.com/denkhaus/agentforge/"

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	@which golangci-lint >/dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2
	@which goimports >/dev/null 2>&1 || $(GOGET) golang.org/x/tools/cmd/goimports@latest
	@which govulncheck >/dev/null 2>&1 || $(GOGET) golang.org/x/vuln/cmd/govulncheck@latest
	$(GOMOD) download
	@echo "Development environment setup complete"

# Check file length constraint
check-file-length:
	@echo "Checking file length constraints..."
	@./scripts/check-file-length.sh 500

# Pre-commit checks
pre-commit: fmt lint vet check-file-length test
	@echo "Pre-commit checks passed"

# CI pipeline
ci: deps lint vet check-file-length test security build
	@echo "CI pipeline completed successfully"

# Release build (optimized)
release:
	@echo "Building release version..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -a -installsuffix cgo -o bin/$(BINARY_NAME)-linux-amd64 ./cmd
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -a -installsuffix cgo -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -a -installsuffix cgo -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd
	@echo "Release builds complete"

# Help
help:
	@echo "Available targets:"
	@echo "  all               - Clean, lint, test, and build"
	@echo "  build             - Build the binary"
	@echo "  clean             - Clean build artifacts"
	@echo "  test              - Run tests"
	@echo "  coverage          - Run tests with coverage report"
	@echo "  lint              - Run linter"
	@echo "  lint-fix          - Fix linting issues automatically"
	@echo "  fmt               - Format code"
	@echo "  vet               - Run go vet"
	@echo "  deps              - Download dependencies"
	@echo "  install           - Install binary to GOPATH/bin"
	@echo "  run-server        - Build and run server"
	@echo "  run-chat          - Build and run chat interface"
	@echo "  run-dev           - Run in development mode"
	@echo "  security          - Check for security vulnerabilities"
	@echo "  docs              - Generate documentation"
	@echo "  dev-setup         - Setup development environment"
	@echo "  check-file-length - Check file length constraints (max 500 lines)"
	@echo "  pre-commit        - Run pre-commit checks"
	@echo "  ci                - Run CI pipeline"
	@echo "  release           - Build release versions for multiple platforms"
	@echo "  help              - Show this help message"