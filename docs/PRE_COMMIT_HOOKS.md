# Pre-commit Hooks Guide

This document explains the pre-commit hook system for AgentForge development.

## Overview

Pre-commit hooks automatically validate code quality before commits are made, ensuring consistent standards across the codebase.

## Installation

### Quick Setup

```bash
make install-hooks
```

This will:
- Install the pre-commit framework
- Set up all configured hooks
- Install custom Git hooks
- Create necessary configuration files

### Manual Installation

```bash
# Install pre-commit framework
pip3 install pre-commit

# Install hooks
pre-commit install
pre-commit install --hook-type commit-msg

# Install custom hooks
./scripts/install-hooks.sh
```

## Hook Types

### Automated Hooks (via pre-commit framework)

1. **Go Formatting** - Ensures consistent Go code formatting
2. **Go Vet** - Runs Go vet for potential issues
3. **Go Mod Tidy** - Keeps go.mod and go.sum clean
4. **Unit Tests** - Runs Go unit tests
5. **Linting** - Runs golangci-lint with project configuration
6. **Trailing Whitespace** - Removes trailing whitespace
7. **File Endings** - Ensures proper file endings
8. **YAML/JSON Validation** - Validates configuration files
9. **Large File Check** - Prevents committing large files (>1MB)
10. **Secret Detection** - Scans for potential secrets

### Custom Hooks

1. **File Size Check** - Enforces 500-line file size limit
2. **Build Verification** - Ensures code compiles
3. **Security Scan** - Runs gosec security analysis
4. **Commit Message Format** - Validates conventional commit format

## Validation Process

When you commit, hooks run in this order:

1. **File Size Validation** - Check 500-line limit
2. **Build Check** - Verify compilation
3. **Test Execution** - Run unit tests
4. **Linting** - Code style validation
5. **Security Scan** - Security vulnerability check
6. **Commit Message** - Format validation

## Commit Message Format

Commits must follow conventional commit format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `perf`: Performance improvements
- `ci`: CI/CD changes

### Examples
```
feat(tui): add prompt optimization feature
fix(database): resolve connection timeout issue
docs(api): update authentication documentation
```

## Bypass Procedures

### Emergency Bypass

For urgent commits that can't pass validation:

```bash
git commit --no-verify -m "emergency: critical hotfix"
```

**⚠️ Use sparingly and fix issues immediately after!**

### WIP Commits

Work-in-progress commits automatically skip validation:

```bash
git commit -m "wip: implementing new feature"
git commit -m "work in progress: refactoring auth"
```

### Merge Commits

Merge commits automatically skip pre-commit validation.

## Troubleshooting

### Hook Installation Issues

```bash
# Check hook status
make check-hooks

# Reinstall hooks
make install-hooks

# Manual hook installation
pre-commit install --install-hooks
```

### Pre-commit Framework Issues

```bash
# Update pre-commit
pip3 install --upgrade pre-commit

# Clean and reinstall
pre-commit clean
pre-commit install
```

### Hook Execution Issues

```bash
# Run specific hook
pre-commit run <hook-name>

# Run all hooks on all files
make pre-commit-all

# Skip specific hook
SKIP=<hook-name> git commit -m "message"
```

### Common Issues

1. **Build Failures**: Fix compilation errors before committing
2. **Test Failures**: Ensure all tests pass
3. **Linting Issues**: Run `make lint` and fix style issues
4. **File Size**: Split large files into smaller modules
5. **Secrets Detected**: Remove sensitive data and update baseline

## Configuration

### Pre-commit Configuration

Edit `.pre-commit-config.yaml` to:
- Add new hooks
- Modify hook arguments
- Exclude files from specific hooks

### Custom Hook Configuration

Modify scripts in `scripts/` directory:
- `scripts/pre-commit` - Main validation logic
- `scripts/commit-msg` - Commit message validation
- `scripts/check-file-length.sh` - File size validation

## Best Practices

1. **Install hooks early** in development process
2. **Run hooks locally** before pushing: `make pre-commit-all`
3. **Fix issues immediately** rather than bypassing
4. **Keep commits small** to make validation faster
5. **Use WIP commits** for work-in-progress
6. **Update hooks regularly** with `make install-hooks`

## Integration with CI/CD

Pre-commit hooks complement CI/CD pipelines:
- **Local validation** catches issues early
- **CI validation** provides comprehensive testing
- **Consistent standards** across all environments

The same validation runs locally and in CI, ensuring no surprises.

## Support

For issues with pre-commit hooks:
1. Check this documentation
2. Run `make check-hooks` for diagnostics
3. Review hook logs in `.git/hooks/`
4. Consult the team for complex issues

Remember: Hooks are there to help maintain code quality and catch issues early!