# Git Branch Strategy & Workflow

> **Implementation Date:** 2025-01-23  
> **Status:** Active  
> **Related Spec:** Sophisticated Git Development Workflow  

## Overview

This document defines the Git branch strategy and workflow for AgentForge, implementing a sophisticated development process with automated quality gates and protected branches.

## Branch Structure

### Primary Branches

#### `main` Branch
- **Purpose:** Production-ready code
- **Protection:** ✅ Fully protected
- **Merge Method:** Squash merge only
- **Required Checks:** `lint`, `test`, `security`, `build`
- **Review Requirements:** 1 approving review minimum

#### `develop` Branch  
- **Purpose:** Integration branch for ongoing development
- **Protection:** ✅ Protected
- **Required Checks:** `integration`
- **Review Requirements:** 1 approving review minimum

### Supporting Branches

#### Feature Branches (`feature/*`)
- **Purpose:** New feature development
- **Naming:** `feature/descriptive-name`
- **Base:** `develop`
- **Merge Target:** `develop`
- **Lifecycle:** Created → Developed → PR → Merged → Deleted

#### Hotfix Branches (`hotfix/*`)
- **Purpose:** Critical production fixes
- **Naming:** `hotfix/issue-description`
- **Base:** `main`
- **Merge Target:** Both `main` and `develop`
- **Lifecycle:** Created → Fixed → PR → Merged → Deleted

## Branch Naming Conventions

### Feature Branches
```
feature/user-authentication
feature/mcp-integration
feature/tui-improvements
feature/api-endpoints
```

### Bug Fix Branches
```
bugfix/login-validation
bugfix/memory-leak-fix
bugfix/dependency-resolution
```

### Hotfix Branches
```
hotfix/security-vulnerability
hotfix/critical-crash
hotfix/data-corruption
```

### Documentation Branches
```
docs/api-documentation
docs/user-guide-update
docs/architecture-diagrams
```

## Workflow Process

### 1. Starting New Feature

```bash
# Switch to develop and update
git checkout develop
git pull origin develop

# Create feature branch
git checkout -b feature/your-feature-name

# Push branch to establish tracking
git push -u origin feature/your-feature-name
```

### 2. Development Cycle

```bash
# Make changes and test
make build
make test

# Commit frequently with descriptive messages
git add .
git commit -m "feat: implement user authentication logic"

# Push regularly to backup work
git push origin feature/your-feature-name
```

### 3. Creating Pull Request

```bash
# Ensure branch is up to date
git checkout develop
git pull origin develop
git checkout feature/your-feature-name
git rebase develop

# Push final changes
git push origin feature/your-feature-name

# Create PR using GitHub CLI
gh pr create \
  --title "Feature: User Authentication System" \
  --body "Implements secure user authentication with JWT tokens" \
  --base develop \
  --assignee @me
```

### 4. Hotfix Process

```bash
# Create hotfix from main
git checkout main
git pull origin main
git checkout -b hotfix/critical-security-fix

# Make minimal fix
# ... implement fix ...
make build
make test

# Commit and push
git add .
git commit -m "hotfix: resolve security vulnerability in auth"
git push -u origin hotfix/critical-security-fix

# Create PR to main
gh pr create \
  --title "Hotfix: Critical Security Vulnerability" \
  --body "Resolves CVE-XXXX in authentication module" \
  --base main

# After merge to main, also merge to develop
git checkout develop
git pull origin develop
git merge main
git push origin develop
```

## Quality Gates

### Automated Checks

All branches must pass these CI/CD checks:

#### Main Branch Requirements
- ✅ **Lint Check** - Code style and formatting validation
- ✅ **Test Suite** - All unit tests must pass
- ✅ **Security Scan** - No high/critical vulnerabilities
- ✅ **Build Verification** - Successful compilation

#### Develop Branch Requirements  
- ✅ **Integration Tests** - Full integration test suite
- ✅ **Performance Benchmarks** - No significant regressions

#### Feature Branch Requirements
- ✅ **Basic Validation** - Lint, test, build
- ✅ **Coverage Check** - Maintain test coverage levels

### Manual Review Process

1. **Code Review** - At least one approving review required
2. **Functional Testing** - Reviewer tests functionality
3. **Documentation Review** - Ensure docs are updated
4. **Breaking Change Assessment** - Evaluate impact

## Repository Settings

### Merge Configuration
- ✅ **Squash Merge Enabled** - Clean commit history
- ❌ **Merge Commits Disabled** - Avoid merge bubbles  
- ❌ **Rebase Merge Disabled** - Consistent squash approach
- ✅ **Auto-delete Head Branches** - Clean up after merge

### Branch Protection
- ✅ **Dismiss Stale Reviews** - New commits require re-review
- ✅ **Require Up-to-date Branches** - Must rebase before merge
- ✅ **Enforce for Administrators** - No exceptions

## Best Practices

### Commit Messages
Use conventional commit format:
```
feat: add user authentication system
fix: resolve memory leak in session handler  
docs: update API documentation
test: add integration tests for auth flow
refactor: simplify dependency injection setup
```

### Branch Management
- Keep feature branches small and focused
- Rebase feature branches regularly on develop
- Delete merged branches promptly
- Use descriptive branch names

### Code Quality
- Run `make build` before every commit
- Execute `make test` to verify functionality
- Use `make lint` to check code style
- Run `make pre-commit` for full validation

### Collaboration
- Create draft PRs early for feedback
- Request specific reviewers for expertise areas
- Respond to review comments promptly
- Test reviewer suggestions before implementing

## Recovery Procedures

### Corrupted Working Directory
```bash
# Check status
git status
git diff

# Restore specific files
git restore path/to/file

# Nuclear option - restore everything
git reset --hard HEAD
```

### Wrong Branch Development
```bash
# Move commits to correct branch
git checkout correct-branch
git cherry-pick commit-hash

# Remove from wrong branch
git checkout wrong-branch
git reset --hard HEAD~1
```

### Failed Merge
```bash
# Abort merge
git merge --abort

# Or resolve conflicts and continue
git add resolved-files
git commit
```

## Monitoring & Metrics

Track these workflow metrics:
- Average PR review time
- Build success rate
- Test coverage trends
- Security scan results
- Branch lifecycle duration

---

This branch strategy ensures code quality, enables safe collaboration, and maintains a clean project history while supporting rapid development cycles.