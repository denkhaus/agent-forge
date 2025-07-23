#!/bin/bash
# Check commit message format for pre-commit hook

# This script is used by pre-commit to validate commit messages
# It reads the commit message from git and validates the format

commit_msg=$(git log --format=%B -n 1 HEAD 2>/dev/null || echo "")

if [ -z "$commit_msg" ]; then
    echo "No commit message found"
    exit 1
fi

commit_regex='^(feat|fix|docs|style|refactor|test|chore|perf|ci)(\(.+\))?: .{1,50}'

if ! echo "$commit_msg" | head -n 1 | grep -qE "$commit_regex"; then
    echo "Invalid commit message format"
    echo ""
    echo "Format: <type>(<scope>): <subject>"
    echo ""
    echo "Types: feat, fix, docs, style, refactor, test, chore, perf, ci"
    echo "Example: feat(tui): add prompt optimization feature"
    echo ""
    echo "Your commit message:"
    echo "$commit_msg"
    echo ""
    exit 1
fi

echo "Commit message format is valid"