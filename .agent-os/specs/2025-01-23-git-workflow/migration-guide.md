# Git Workflow Migration Guide

> **Spec:** Sophisticated Git Development Workflow  
> **Created:** 2025-01-23  
> **Focus:** Step-by-step migration from current workflow  

## Current State Analysis

### Existing Workflow Assessment

**Current Git State:**
- Single `main` branch with all development
- Direct commits to main branch
- Good commit message patterns already established
- Build validation before commits (via DEV_WORKFLOW.md)
- No feature branch isolation

**Recent Commit Analysis:**
```bash
f5be934 feat: DI Enforcement Phase 3 Complete - Services Layer Implementation
fe8055d feat: DI Enforcement Phase 1 & 2 Complete + Code Cleanup
e3ac035 🚀 Enhanced Workbench V3: Professional TUI Implementation Complete
```

**Strengths to Preserve:**
- Commit message quality and structure
- Build validation discipline
- Recovery procedures and Git knowledge
- Structured development approach

## Migration Strategy

### Phase 1: Infrastructure Setup (Day 1-2)

#### Step 1: Create Develop Branch
```bash
# Backup current state
git tag backup-before-workflow-migration

# Create develop branch from main
git checkout main
git pull origin main
git checkout -b develop
git push origin develop

# Set develop as default branch for new PRs
gh repo edit --default-branch develop
```

#### Step 2: Set Up Branch Protection
```bash
# Install GitHub CLI if not available
# Run branch protection setup script
./scripts/setup-branch-protection.sh

# Verify protection rules
gh api repos/:owner/:repo/branches/main/protection
gh api repos/:owner/:repo/branches/develop/protection
```

#### Step 3: Update CI/CD
```bash
# Backup existing workflow
cp .github/workflows/ci.yml .github/workflows/ci.yml.backup

# Copy new workflow files
cp .agent-os/specs/2025-01-23-git-workflow/workflows/* .github/workflows/

# Test workflows
git add .github/workflows/
git commit -m "ci: update workflows for feature branch support"
git push origin develop
```

### Phase 2: Current Work Migration (Day 3-4)

#### Step 1: Identify Current Changes
```bash
# Check current status
git status
git diff HEAD

# Identify logical feature groups from recent commits
git log --oneline -10
```

#### Step 2: Create Feature Branches for Active Work
```bash
# For Agent OS work
git checkout develop
git checkout -b feature/agent-os-installation
# Cherry-pick or move relevant commits
git cherry-pick <commit-hash>

# For TUI workbench work  
git checkout develop
git checkout -b feature/tui-prompt-workbench
# Move TUI-related changes

# For any other ongoing features
git checkout develop
git checkout -b feature/di-enforcement-completion
```

#### Step 3: Clean Up Main Branch
```bash
# Ensure main is clean and matches last stable state
git checkout main
git reset --hard <last-stable-commit>
git push origin main --force-with-lease

# Update develop to match main
git checkout develop
git reset --hard main
git push origin develop --force-with-lease
```

### Phase 3: Workflow Implementation (Day 5-7)

#### Step 1: Install Development Tools
```bash
# Install pre-commit hooks
make install-hooks

# Test pre-commit functionality
echo "test" > test-file.txt
git add test-file.txt
git commit -m "test: verify pre-commit hooks"
# Should run validation

# Clean up test
git reset HEAD~1
rm test-file.txt
```

#### Step 2: Update Makefile
```bash
# Add new workflow commands to Makefile
# Test each command
make feature-start
# Enter test feature name
make feature-sync
make feature-finish
```

#### Step 3: Test Complete Workflow
```bash
# Test full feature development cycle
make feature-start
# Enter: test-workflow

# Make a small change
echo "# Test Feature" > test-feature.md
git add test-feature.md
git commit -m "feat: add test feature documentation"

# Test sync and finish
make feature-sync
make feature-finish

# Create PR via GitHub UI
# Merge PR and verify automation
```

## Migration Checklist

### Pre-Migration Preparation
- [ ] Backup current repository state with tag
- [ ] Document current uncommitted changes
- [ ] Identify active feature work in progress
- [ ] Ensure all team members are informed
- [ ] Verify GitHub CLI access and permissions

### Infrastructure Setup
- [ ] Create and push develop branch
- [ ] Configure branch protection rules for main
- [ ] Configure branch protection rules for develop
- [ ] Update default branch settings
- [ ] Deploy new CI/CD workflows
- [ ] Test workflow execution

### Work Migration
- [ ] Create feature branches for active work
- [ ] Migrate uncommitted changes to appropriate branches
- [ ] Clean up main branch to stable state
- [ ] Verify all work is preserved in feature branches
- [ ] Test feature branch CI/CD execution

### Tool Installation
- [ ] Install pre-commit hooks
- [ ] Update Makefile with workflow commands
- [ ] Test all new make commands
- [ ] Verify pre-commit validation works
- [ ] Test complete feature development cycle

### Documentation Updates
- [ ] Update DEV_WORKFLOW.md with new process
- [ ] Create quick reference guide
- [ ] Update contributing guidelines
- [ ] Document troubleshooting procedures
- [ ] Create workflow examples

### Team Onboarding
- [ ] Conduct workflow training session
- [ ] Provide hands-on practice with new commands
- [ ] Create workflow cheat sheet
- [ ] Set up mentoring for adoption period
- [ ] Gather feedback and address concerns

## Rollback Plan

### Emergency Rollback Procedure
```bash
# If migration causes critical issues
git checkout main
git reset --hard backup-before-workflow-migration
git push origin main --force-with-lease

# Restore original CI
cp .github/workflows/ci.yml.backup .github/workflows/ci.yml
git add .github/workflows/ci.yml
git commit -m "rollback: restore original CI workflow"
git push origin main

# Remove branch protection temporarily
gh api repos/:owner/:repo/branches/main/protection -X DELETE
```

### Partial Rollback Options
- **Keep infrastructure, revert work migration:** Maintain branches but move work back to main
- **Keep branches, revert automation:** Maintain feature branches but disable hooks/CI
- **Gradual adoption:** Use new workflow for new features only

## Success Validation

### Technical Validation
```bash
# Verify branch structure
git branch -a
# Should show: main, develop, feature branches

# Test CI/CD
git push origin feature/test-branch
# Should trigger feature branch CI

# Test protection
git push origin main
# Should be rejected

# Test workflow commands
make feature-start
make feature-sync  
make feature-finish
# All should work correctly
```

### Process Validation
- [ ] Feature branches are created correctly
- [ ] CI runs on all feature branches
- [ ] PRs are required for main and develop
- [ ] Pre-commit hooks prevent bad commits
- [ ] Release process works end-to-end
- [ ] Team can use workflow effectively

### Quality Validation
- [ ] Build success rate >95%
- [ ] No broken builds on main branch
- [ ] Test coverage maintained or improved
- [ ] Security scans pass
- [ ] Code quality metrics stable

## Troubleshooting

### Common Issues

#### "Branch protection prevents push"
```bash
# Solution: Use PR process
git push origin feature/my-feature
# Create PR via GitHub UI
```

#### "Pre-commit hooks failing"
```bash
# Bypass for emergency (use sparingly)
git commit --no-verify -m "emergency: bypass hooks"

# Fix and re-commit properly
make pre-commit-full
git commit --amend
```

#### "Merge conflicts during rebase"
```bash
# Resolve conflicts
git status
# Edit conflicted files
git add .
git rebase --continue

# Or abort and merge instead
git rebase --abort
git merge origin/develop
```

#### "CI failing on feature branch"
```bash
# Check CI logs in GitHub Actions
# Fix issues locally
make build test lint
git add .
git commit -m "fix: resolve CI issues"
git push origin feature/branch-name
```

This migration guide ensures a smooth transition from the current direct-to-main workflow to the sophisticated feature branch workflow while preserving all existing work and maintaining development velocity.