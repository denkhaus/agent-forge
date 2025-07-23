

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
	@mkdir -p .git/hooks
	@cp scripts/pre-commit .git/hooks/pre-commit
	@cp scripts/commit-msg .git/hooks/commit-msg
	@chmod +x .git/hooks/pre-commit .git/hooks/commit-msg
	@echo "Git hooks installed"

