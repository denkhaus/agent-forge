#!/bin/bash
# Install pre-commit hooks for AgentForge development

set -e

echo "Installing pre-commit hooks for AgentForge..."

# Check if pre-commit is installed
if ! command -v pre-commit &> /dev/null; then
    echo "Installing pre-commit..."
    pip3 install pre-commit || pip install pre-commit
fi

# Install pre-commit hooks
echo "Installing pre-commit hooks..."
pre-commit install

# Install commit-msg hook
echo "Installing commit-msg hook..."
pre-commit install --hook-type commit-msg

# Copy custom git hooks
echo "Installing custom git hooks..."
mkdir -p .git/hooks

# Install pre-commit hook
cp scripts/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

# Install commit-msg hook
cp scripts/commit-msg .git/hooks/commit-msg
chmod +x .git/hooks/commit-msg

# Create secrets baseline if it doesn't exist
if [ ! -f .secrets.baseline ]; then
    echo "Creating secrets baseline..."
    detect-secrets scan --baseline .secrets.baseline || echo "{}" > .secrets.baseline
fi

echo "✅ Pre-commit hooks installed successfully!"
echo ""
echo "Available commands:"
echo "  make install-hooks  - Install/update hooks"
echo "  pre-commit run --all-files  - Run all hooks on all files"
echo "  git commit --no-verify  - Bypass hooks (emergency only)"
echo ""
echo "Hooks will now run automatically on every commit."