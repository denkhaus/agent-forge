

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

# Full pre-commit validation
pre-commit-full:
	@echo "Running full pre-commit validation..."
	@$(MAKE) check-file-length
	@$(MAKE) build
	@$(MAKE) test
	@$(MAKE) lint
	@echo "All validation checks passed"

# Install git hooks
install-hooks:
	@echo "Installing git hooks..."
	@mkdir -p .git/hooks
	@cp scripts/pre-commit .git/hooks/pre-commit
	@cp scripts/commit-msg .git/hooks/commit-msg
	@chmod +x .git/hooks/pre-commit .git/hooks/commit-msg
	@echo "Git hooks installed"

