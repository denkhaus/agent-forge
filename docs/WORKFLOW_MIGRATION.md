# Git Workflow Migration Guide

This document provides a comprehensive guide for migrating existing development work to the new sophisticated Git workflow.

## Overview

The AgentForge project has implemented a sophisticated Git workflow with:
- Protected main and develop branches
- Feature branch development
- Automated CI/CD pipelines
- Pre-commit hooks and quality gates
- Comprehensive workflow commands

## Migration Process

### 1. Current State Analysis

**Completed Migration Items:**
- ✅ **Git Workflow Foundation** - Branch protection, CI/CD, pre-commit hooks
- ✅ **Enhanced Makefile Commands** - Complete workflow command set
- ✅ **TUI Prompt Workbench Spec** - Migrated to dedicated feature branch
- ✅ **Agent OS Installation** - Already completed and integrated

**Branch Organization:**
- `main` - Production-ready code (protected)
- `develop` - Integration branch (protected)
- `feature/tui-prompt-workbench` - TUI workbench development
- `feature/enhanced-makefile-commands` - Workflow commands (ready for merge)

### 2. Workflow Commands Available

#### Feature Development
```bash
make feature-start    # Start new feature branch from develop
make feature-sync     # Sync feature branch with develop
make feature-finish   # Complete feature development and push
```

#### Release Management
```bash
make release-start    # Start release branch from develop
make release-finish   # Complete release preparation
```

#### Emergency Hotfixes
```bash
make hotfix-start     # Start emergency hotfix from main
make hotfix-finish    # Complete hotfix development
```

#### Quality Assurance
```bash
make pre-commit-full  # Run comprehensive validation
make install-hooks    # Install pre-commit hooks
make check-hooks      # Verify hook installation
```

### 3. Developer Workflow

#### Starting New Work
1. **Ensure clean state**: `git status`
2. **Switch to develop**: `git checkout develop && git pull origin develop`
3. **Start feature**: `make feature-start`
4. **Enter feature name** when prompted

#### During Development
1. **Make changes** and test locally
2. **Commit frequently**: `git add . && git commit -m "descriptive message"`
3. **Sync with develop**: `make feature-sync` (periodically)
4. **Push regularly**: `git push origin feature/your-feature`

#### Completing Work
1. **Final validation**: `make pre-commit-full`
2. **Finish feature**: `make feature-finish`
3. **Create PR**: Use GitHub interface or `gh pr create`
4. **Review and merge**: Follow PR review process

### 4. Branch Protection Rules

#### Main Branch
- ✅ Requires PR reviews (1 approver minimum)
- ✅ Requires status checks: `lint`, `test`, `security`, `build`
- ✅ Dismisses stale reviews
- ✅ Enforces for administrators

#### Develop Branch
- ✅ Requires PR reviews (1 approver minimum)
- ✅ Requires status checks: `integration`
- ✅ Dismisses stale reviews
- ✅ Enforces for administrators

### 5. Quality Gates

#### Pre-commit Hooks
- File size validation (500 line limit)
- Build verification
- Test execution
- Code formatting and linting
- Security scanning

#### CI/CD Pipeline
- **Feature Branches**: Lint, test, security scan, build
- **Develop Branch**: Integration tests, performance benchmarks
- **Main Branch**: Full validation, release preparation

### 6. Migration Best Practices

#### For Existing Work
1. **Create feature branches** for any work in progress
2. **Use descriptive names**: `feature/user-authentication`, `feature/api-integration`
3. **Keep branches focused** on single features or fixes
4. **Sync regularly** with develop to avoid conflicts

#### For New Development
1. **Always start from develop**: `make feature-start`
2. **Follow conventional commits**: `feat:`, `fix:`, `docs:`, etc.
3. **Test before committing**: Pre-commit hooks will enforce this
4. **Create PRs early** for feedback and collaboration

#### For Emergency Fixes
1. **Use hotfix workflow**: `make hotfix-start`
2. **Work from main branch** for critical production issues
3. **Merge to both main and develop** after completion
4. **Keep hotfixes minimal** and focused

### 7. Troubleshooting

#### Common Issues

**Branch Protection Errors**
```bash
# Error: Protected branch update failed
# Solution: Create PR instead of direct push
gh pr create --base develop --title "Your Feature"
```

**Pre-commit Hook Failures**
```bash
# Error: Build failed, tests failed, etc.
# Solution: Fix issues locally first
make build          # Check build
make test           # Check tests
make lint           # Check linting
```

**Workflow Command Issues**
```bash
# Error: Not on correct branch type
# Solution: Check current branch
git branch --show-current
git checkout develop  # Switch to correct branch
```

#### Recovery Procedures

**Corrupted Working Directory**
```bash
git status          # Check what changed
git restore .       # Restore all files
git reset --hard HEAD  # Nuclear option
```

**Wrong Branch Development**
```bash
# Move commits to correct branch
git checkout correct-branch
git cherry-pick commit-hash
git checkout wrong-branch
git reset --hard HEAD~1
```

### 8. Team Adoption

#### Onboarding New Developers
1. **Clone repository**: `git clone <repo-url>`
2. **Install hooks**: `make install-hooks`
3. **Read documentation**: This guide + BRANCH_STRATEGY.md
4. **Practice workflow**: Create test feature branch

#### Training Checklist
- [ ] Understand branch structure (main/develop/feature)
- [ ] Know workflow commands (`make feature-start`, etc.)
- [ ] Familiar with pre-commit hooks
- [ ] Can create and manage PRs
- [ ] Understands quality gates and CI/CD

### 9. Success Metrics

The migration is successful when:
- ✅ All development happens on feature branches
- ✅ No direct pushes to main or develop
- ✅ Pre-commit hooks prevent quality issues
- ✅ CI/CD pipeline catches integration problems
- ✅ Team follows consistent workflow patterns

### 10. Support and Resources

#### Documentation
- `BRANCH_STRATEGY.md` - Detailed branch strategy
- `PRE_COMMIT_HOOKS.md` - Hook installation and usage
- `.github/pull_request_template.md` - PR template

#### Commands Reference
```bash
make help           # Show all available commands
make check-hooks    # Verify hook installation
make pre-commit-all # Run all hooks manually
```

#### Getting Help
1. Check this documentation first
2. Review error messages carefully
3. Use `git status` to understand current state
4. Ask team for assistance with complex scenarios

---

**Remember**: The workflow is designed to help maintain code quality and enable safe collaboration. When in doubt, create a feature branch and ask for review!