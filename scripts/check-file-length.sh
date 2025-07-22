#!/bin/bash

# Script to check Go file length constraint (max 500 lines)
# Usage: ./scripts/check-file-length.sh [max_lines]

set -e

MAX_LINES=${1:-500}
EXIT_CODE=0

echo "Checking Go files for maximum length of ${MAX_LINES} lines..."

# Find all .go files excluding vendor, generated files, and test files for now
while IFS= read -r -d '' file; do
    line_count=$(wc -l < "$file")
    
    if [ "$line_count" -gt "$MAX_LINES" ]; then
        echo "❌ $file: $line_count lines (exceeds $MAX_LINES)"
        EXIT_CODE=1
    else
        echo "✅ $file: $line_count lines"
    fi
done < <(find . -name "*.go" \
    -not -path "./vendor/*" \
    -not -path "./.git/*" \
    -not -name "*_test.go" \
    -not -name "*.pb.go" \
    -not -name "*_gen.go" \
    -print0)

if [ $EXIT_CODE -eq 0 ]; then
    echo "✅ All Go files are within the $MAX_LINES line limit"
else
    echo "❌ Some files exceed the $MAX_LINES line limit"
    echo ""
    echo "Consider refactoring large files by:"
    echo "  1. Extracting functions into separate files"
    echo "  2. Moving related types to dedicated files"
    echo "  3. Splitting large structs into smaller components"
    echo "  4. Creating sub-packages for related functionality"
fi

exit $EXIT_CODE