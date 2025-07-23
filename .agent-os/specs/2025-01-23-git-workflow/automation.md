# Git Workflow Automation - Technical Details

> **Spec:** Sophisticated Git Development Workflow  
> **Created:** 2025-01-23  
> **Focus:** CI/CD pipelines and automation scripts  

## Enhanced GitHub Actions Workflows

### Feature Branch Workflow

```yaml
# .github/workflows/feature-branch.yml
name: Feature Branch CI

on:
  push:
    branches: [ 'feature/*' ]
  pull_request:
    branches: [ 'develop' ]

env:
  GO_VERSION: '1.24.0'

jobs:
  validate:
    name: Validate Feature
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - name: Cache dependencies
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run linter
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
          args: --config .golangci.yml
      
      - name: Check file sizes
        run: make check-file-length
      
      - name: Run tests with coverage
        run: |
          go test -v -race -coverprofile=coverage.out ./...
          go tool cover -html=coverage.out -o coverage.html
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
          flags: feature-tests
      
      - name: Build application
        run: make build
      
      - name: Security scan
        uses: securecodewarrior/github-action-gosec@master
        with:
          args: '-no-fail -fmt sarif -out results.sarif ./...'
      
      - name: Upload security results
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: results.sarif
```

### Develop Branch Integration

```yaml
# .github/workflows/develop-integration.yml
name: Develop Integration

on:
  push:
    branches: [ 'develop' ]
  pull_request:
    branches: [ 'develop' ]

jobs:
  integration-tests:
    name: Integration Tests
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24.0'
      
      - name: Run integration tests
        run: make test-integration
      
      - name: Run TUI tests
        run: make test-tui
      
      - name: Performance benchmarks
        run: make benchmark
      
      - name: Generate test report
        run: make test-report
      
      - name: Upload test artifacts
        uses: actions/upload-artifact@v3
        with:
          name: integration-test-results
          path: |
            test-results/
            benchmarks/
```

### Release Workflow

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    branches: [ 'release/*' ]
  pull_request:
    branches: [ 'main' ]
    types: [ closed ]

jobs:
  release:
    name: Create Release
    runs-on: ubuntu-latest
    if: github.event.pull_request.merged == true && startsWith(github.head_ref, 'release/')
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Extract version
        id: version
        run: |
          VERSION=${GITHUB_HEAD_REF#release/}
          echo "version=$VERSION" >> $GITHUB_OUTPUT
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24.0'
      
      - name: Build release binaries
        run: make build-release
      
      - name: Generate changelog
        id: changelog
        run: |
          make generate-changelog > CHANGELOG.tmp
          echo "changelog<<EOF" >> $GITHUB_OUTPUT
          cat CHANGELOG.tmp >> $GITHUB_OUTPUT
          echo "EOF" >> $GITHUB_OUTPUT
      
      - name: Create GitHub release
        uses: actions/create-release@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          tag_name: ${{ steps.version.outputs.version }}
          release_name: Release ${{ steps.version.outputs.version }}
          body: ${{ steps.changelog.outputs.changelog }}
          draft: false
          prerelease: false
      
      - name: Upload release assets
        uses: actions/upload-release-asset@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          upload_url: ${{ steps.create_release.outputs.upload_url }}
          asset_path: ./bin/forge
          asset_name: forge-${{ steps.version.outputs.version }}
          asset_content_type: application/octet-stream
```

## Pre-commit Hooks

### Main Pre-commit Hook

```bash
#!/bin/sh
# .git/hooks/pre-commit

set -e

echo "🔍 Running pre-commit validation..."

# Check if this is a merge commit
if git rev-parse -q --verify MERGE_HEAD >/dev/null; then
    echo "ℹ️  Merge commit detected, skipping pre-commit hooks"
    exit 0
fi

# Check for WIP commits
commit_msg=$(git log --format=%B -n 1 HEAD 2>/dev/null || echo "")
if echo "$commit_msg" | grep -qi "wip\|work in progress\|fixup\|squash"; then
    echo "⚠️  WIP commit detected, skipping validation"
    exit 0
fi

# 1. File size validation
echo "📏 Checking file sizes..."
if ! ./scripts/check-file-length.sh; then
    echo "❌ File size check failed"
    exit 1
fi

# 2. Build validation
echo "🔨 Building application..."
if ! make build >/dev/null 2>&1; then
    echo "❌ Build failed"
    echo "Run 'make build' to see detailed errors"
    exit 1
fi

# 3. Test execution
echo "🧪 Running tests..."
if ! make test >/dev/null 2>&1; then
    echo "❌ Tests failed"
    echo "Run 'make test' to see detailed errors"
    exit 1
fi

# 4. Linting
echo "🔍 Running linter..."
if ! make lint >/dev/null 2>&1; then
    echo "❌ Linting failed"
    echo "Run 'make lint' to see detailed errors"
    exit 1
fi

# 5. Security check
echo "🔒 Running security scan..."
if ! make security-scan >/dev/null 2>&1; then
    echo "⚠️  Security scan found issues"
    echo "Run 'make security-scan' to see details"
    # Don't fail on security issues, just warn
fi

echo "✅ All pre-commit checks passed!"
```

### Commit Message Hook

```bash
#!/bin/sh
# .git/hooks/commit-msg

commit_regex='^(feat|fix|docs|style|refactor|test|chore|perf|ci)(\(.+\))?: .{1,50}'

if ! grep -qE "$commit_regex" "$1"; then
    echo "❌ Invalid commit message format"
    echo ""
    echo "Format: <type>(<scope>): <subject>"
    echo ""
    echo "Types: feat, fix, docs, style, refactor, test, chore, perf, ci"
    echo "Example: feat(tui): add prompt optimization feature"
    echo ""
    exit 1
fi

# Check for conventional commit format
if ! head -n 1 "$1" | grep -qE '^(feat|fix|docs|style|refactor|test|chore|perf|ci)(\(.+\))?: .+'; then
    echo "❌ Commit message must follow conventional commit format"
    exit 1
fi

echo "✅ Commit message format is valid"
```

## Makefile Enhancements

### Git Workflow Commands

```makefile
# Git workflow commands
.PHONY: feature-start feature-sync feature-finish release-start release-finish

# Start a new feature branch
feature-start:
	@echo "🚀 Starting new feature branch..."
	@read -p "Feature name (without 'feature/' prefix): " name; \
	if [ -z "$$name" ]; then \
		echo "❌ Feature name cannot be empty"; \
		exit 1; \
	fi; \
	git checkout develop && \
	git pull origin develop && \
	git checkout -b feature/$$name && \
	echo "✅ Feature branch 'feature/$$name' created and checked out"

# Sync feature branch with develop
feature-sync:
	@echo "🔄 Syncing feature branch with develop..."
	@current_branch=$$(git branch --show-current); \
	if [[ ! $$current_branch =~ ^feature/ ]]; then \
		echo "❌ Not on a feature branch"; \
		exit 1; \
	fi; \
	git fetch origin && \
	git rebase origin/develop && \
	echo "✅ Feature branch synced with develop"

# Finish feature development
feature-finish:
	@echo "🏁 Finishing feature development..."
	@current_branch=$$(git branch --show-current); \
	if [[ ! $$current_branch =~ ^feature/ ]]; then \
		echo "❌ Not on a feature branch"; \
		exit 1; \
	fi; \
	make pre-commit-full && \
	git push origin $$current_branch && \
	echo "✅ Feature branch ready for PR" && \
	echo "Create PR: https://github.com/denkhaus/agentforge/compare/develop...$$current_branch"

# Start release branch
release-start:
	@echo "🎯 Starting release branch..."
	@read -p "Version (e.g., v0.2.0): " version; \
	if [ -z "$$version" ]; then \
		echo "❌ Version cannot be empty"; \
		exit 1; \
	fi; \
	git checkout develop && \
	git pull origin develop && \
	git checkout -b release/$$version && \
	echo "✅ Release branch 'release/$$version' created"

# Complete release
release-finish:
	@echo "🚢 Finishing release..."
	@current_branch=$$(git branch --show-current); \
	if [[ ! $$current_branch =~ ^release/ ]]; then \
		echo "❌ Not on a release branch"; \
		exit 1; \
	fi; \
	make pre-commit-full && \
	git push origin $$current_branch && \
	echo "✅ Release branch ready for PR to main"

# Full pre-commit validation
pre-commit-full:
	@echo "🔍 Running full pre-commit validation..."
	@make check-file-length
	@make build
	@make test
	@make lint
	@make security-scan
	@echo "✅ All validation checks passed"

# Install git hooks
install-hooks:
	@echo "🪝 Installing git hooks..."
	@cp scripts/pre-commit .git/hooks/pre-commit
	@cp scripts/commit-msg .git/hooks/commit-msg
	@chmod +x .git/hooks/pre-commit .git/hooks/commit-msg
	@echo "✅ Git hooks installed"
```

This automation framework provides comprehensive CI/CD support for the sophisticated Git workflow while maintaining fast feedback loops and high code quality.