# Git Branch Strategy - Technical Details

> **Spec:** Sophisticated Git Development Workflow  
> **Created:** 2025-01-23  
> **Focus:** Branch management and protection rules  

## Branch Protection Configuration

### Main Branch Protection

```yaml
# GitHub API configuration for main branch protection
{
  "required_status_checks": {
    "strict": true,
    "contexts": [
      "ci/lint",
      "ci/test", 
      "ci/security",
      "ci/build"
    ]
  },
  "enforce_admins": true,
  "required_pull_request_reviews": {
    "required_approving_review_count": 1,
    "dismiss_stale_reviews": true,
    "require_code_owner_reviews": false,
    "require_last_push_approval": true
  },
  "restrictions": null,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "block_creations": false
}
```

### Develop Branch Protection

```yaml
{
  "required_status_checks": {
    "strict": true,
    "contexts": [
      "ci/lint",
      "ci/test",
      "ci/integration"
    ]
  },
  "enforce_admins": false,
  "required_pull_request_reviews": {
    "required_approving_review_count": 1,
    "dismiss_stale_reviews": true
  },
  "allow_force_pushes": false,
  "allow_deletions": false
}
```

## Branch Naming Conventions

### Feature Branches
```bash
# Format: feature/brief-description-with-hyphens
feature/tui-prompt-workbench
feature/component-discovery
feature/github-integration
feature/ai-optimization-engine
feature/session-persistence

# Invalid examples:
feature/TUI_Prompt_Workbench  # No underscores or capitals
feature/fix-bug               # Use fix/ prefix for bug fixes
feature/123-add-feature       # No leading numbers
```

### Release Branches
```bash
# Format: release/vMAJOR.MINOR.PATCH
release/v0.1.0
release/v0.2.0
release/v1.0.0

# Pre-release versions
release/v0.2.0-alpha.1
release/v0.2.0-beta.1
release/v0.2.0-rc.1
```

### Hotfix Branches
```bash
# Format: hotfix/brief-description
hotfix/critical-memory-leak
hotfix/security-vulnerability
hotfix/data-corruption-fix

# Emergency hotfixes
hotfix/emergency-auth-bypass
```

### Support Branches
```bash
# Format: support/version-number
support/v1.0.x
support/v0.9.x

# For maintaining older versions
```

## Git Configuration

### Repository Settings

```bash
# Configure Git for the repository
git config branch.main.description "Production-ready code only"
git config branch.develop.description "Integration branch for features"

# Set up default branch
git config init.defaultBranch main

# Configure merge strategy
git config merge.ours.driver true
git config pull.rebase true

# Set up commit template
git config commit.template .gitmessage
```

### Global Git Configuration

```bash
# Developer setup commands
git config --global user.name "Developer Name"
git config --global user.email "developer@agentforge.dev"
git config --global init.defaultBranch main
git config --global pull.rebase true
git config --global rebase.autoStash true
git config --global merge.conflictstyle diff3
```

## Branch Lifecycle Management

### Feature Branch Lifecycle

```bash
# 1. Create feature branch
git checkout develop
git pull origin develop
git checkout -b feature/new-feature

# 2. Development work
# ... make changes ...
git add .
git commit -m "feat: implement new feature"

# 3. Keep up to date
git fetch origin
git rebase origin/develop

# 4. Push and create PR
git push origin feature/new-feature
# Create PR via GitHub UI or CLI

# 5. After merge, cleanup
git checkout develop
git pull origin develop
git branch -d feature/new-feature
git push origin --delete feature/new-feature
```

### Release Branch Lifecycle

```bash
# 1. Create release branch
git checkout develop
git pull origin develop
git checkout -b release/v0.2.0

# 2. Release preparation
# Update version numbers
# Update CHANGELOG.md
# Final bug fixes only

# 3. Merge to main
git checkout main
git merge --no-ff release/v0.2.0
git tag -a v0.2.0 -m "Release version 0.2.0"
git push origin main --tags

# 4. Merge back to develop
git checkout develop
git merge --no-ff release/v0.2.0
git push origin develop

# 5. Cleanup
git branch -d release/v0.2.0
git push origin --delete release/v0.2.0
```

## Automated Branch Management

### GitHub Actions for Branch Management

```yaml
# .github/workflows/branch-management.yml
name: Branch Management

on:
  push:
    branches: [ 'feature/*', 'hotfix/*', 'release/*' ]
  pull_request:
    branches: [ 'develop', 'main' ]

jobs:
  validate-branch-name:
    runs-on: ubuntu-latest
    steps:
      - name: Validate branch name
        run: |
          branch=${GITHUB_REF#refs/heads/}
          if [[ ! $branch =~ ^(feature|hotfix|release)/.+ ]]; then
            echo "Invalid branch name: $branch"
            exit 1
          fi
  
  auto-cleanup:
    runs-on: ubuntu-latest
    if: github.event.pull_request.merged == true
    steps:
      - name: Delete merged branch
        run: |
          gh api repos/${{ github.repository }}/git/refs/heads/${{ github.head_ref }} \
            -X DELETE
```

### Branch Protection Scripts

```bash
#!/bin/bash
# scripts/setup-branch-protection.sh

# Setup main branch protection
gh api repos/:owner/:repo/branches/main/protection \
  --method PUT \
  --field required_status_checks='{"strict":true,"contexts":["ci/lint","ci/test","ci/security","ci/build"]}' \
  --field enforce_admins=true \
  --field required_pull_request_reviews='{"required_approving_review_count":1,"dismiss_stale_reviews":true}' \
  --field restrictions=null \
  --field allow_force_pushes=false \
  --field allow_deletions=false

# Setup develop branch protection  
gh api repos/:owner/:repo/branches/develop/protection \
  --method PUT \
  --field required_status_checks='{"strict":true,"contexts":["ci/lint","ci/test","ci/integration"]}' \
  --field enforce_admins=false \
  --field required_pull_request_reviews='{"required_approving_review_count":1,"dismiss_stale_reviews":true}' \
  --field allow_force_pushes=false \
  --field allow_deletions=false
```

This branch strategy provides the foundation for controlled, professional Git workflow management.