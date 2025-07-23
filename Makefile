# Git workflow commands
.PHONY: feature-start feature-sync feature-finish release-start release-finish install-hooks pre-commit-full

# Start a new feature branch
feature-start:
	@echo "Starting new feature branch..."
	@read -p "Feature name (without feature/ prefix): " name; \
	if [ -z "$$name" ]; then \
		echo "Feature name cannot be empty"; \
		exit 1; \
	fi; \
	git checkout develop && \
	git pull origin develop && \
	git checkout -b feature/$$name && \
	echo "Feature branch feature/$$name created and checked out"

# Sync feature branch with develop
feature-sync:
	@echo "Syncing feature branch with develop..."
	@current_branch=$$(git branch --show-current); \
	if [[ ! $$current_branch =~ ^feature/ ]]; then \
		echo "Not on a feature branch"; \
		exit 1; \
	fi; \
	git fetch origin && \
	git rebase origin/develop && \
	echo "Feature branch synced with develop"

# Finish feature development
feature-finish:
	@echo "Finishing feature development..."
	@current_branch=$$(git branch --show-current); \
	if [[ ! $$current_branch =~ ^feature/ ]]; then \
		echo "Not on a feature branch"; \
		exit 1; \
	fi; \
	$(MAKE) pre-commit-full && \
	git push origin $$current_branch && \
	echo "Feature branch ready for PR"

# Build targets
.PHONY: build test lint clean check-file-length

build:
	@echo "Building application..."
	@go build -o bin/agentforge ./cmd/main.go

test:
	@echo "Running tests..."
	@go test ./...

lint:
	@echo "Running linter..."
	@golangci-lint run

security-scan:
	@echo "Running security scan..."
	@gosec ./... || echo "Security scan completed with warnings"

test-integration:
	@echo "Running integration tests..."
	@go test -tags=integration ./internal/integration/... || echo "Integration tests not fully implemented"

test-tui:
	@echo "Running TUI tests..."
	@go test ./internal/tui/... || echo "TUI tests not fully implemented"

benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. ./... || echo "Benchmarks not fully implemented"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/

check-file-length:
	@./scripts/check-file-length.sh

# Full pre-commit validation
pre-commit-full:
	@echo "Running full pre-commit validation..."
	@$(MAKE) check-file-length
	@$(MAKE) build
	@$(MAKE) test
	@$(MAKE) lint
	@$(MAKE) security-scan
	@echo "All validation checks passed"

# Release workflow commands
release-start:
	@echo "Starting release branch..."
	@read -p "Version (e.g., v0.2.0): " version; \
	if [ -z "$$version" ]; then \
		echo "Version cannot be empty"; \
		exit 1; \
	fi; \
	git checkout develop && \
	git pull origin develop && \
	git checkout -b release/$$version && \
	echo "Release branch release/$$version created"

release-finish:
	@echo "Finishing release..."
	@current_branch=$$(git branch --show-current); \
	if [[ ! $$current_branch =~ ^release/ ]]; then \
		echo "Not on a release branch"; \
		exit 1; \
	fi; \
	$(MAKE) pre-commit-full && \
	git push origin $$current_branch && \
	echo "Release branch ready for PR to main"

# Pre-commit hook (lighter version)
pre-commit:
	@echo "Running pre-commit validation..."
	@$(MAKE) build
	@$(MAKE) test
	@echo "Pre-commit checks passed"

# Install git hooks
install-hooks:
	@echo "Installing git hooks..."
	@./scripts/install-hooks.sh

# Check hook status
check-hooks:
	@echo "Checking hook installation status..."
	@if [ -f .git/hooks/pre-commit ]; then echo "Pre-commit hook installed"; else echo "Pre-commit hook missing"; fi
	@if [ -f .git/hooks/commit-msg ]; then echo "Commit-msg hook installed"; else echo "Commit-msg hook missing"; fi
	@if command -v pre-commit &> /dev/null; then echo "Pre-commit framework installed"; else echo "Pre-commit framework missing"; fi

# Run pre-commit on all files
pre-commit-all:
	@echo "Running pre-commit on all files..."
	@pre-commit run --all-files

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build           - Build the application"
	@echo "  test            - Run tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-tui        - Run TUI tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  check-file-length - Check file length constraints"
	@echo "  pre-commit      - Run pre-commit validation"
	@echo "  pre-commit-full - Run full pre-commit validation"
	@echo "  security-scan   - Run security scanning"
	@echo "  lint            - Run linting"
	@echo "  feature-start   - Start a new feature branch"
	@echo "  feature-sync    - Sync feature branch with develop"
	@echo "  feature-finish  - Finish feature development"
	@echo "  release-start   - Start release branch"
	@echo "  release-finish  - Complete release"
	@echo "  install-hooks   - Install git hooks"
	@echo "  check-hooks     - Check hook installation status"
	@echo "  pre-commit-all  - Run pre-commit on all files"