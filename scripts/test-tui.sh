#!/bin/bash

# TUI Testing Script for AgentForge Prompt Workbench
# This script provides various ways to test the TUI without blocking development

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[TUI TEST]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to run a test mode
run_test() {
    local mode=$1
    local prompt=$2
    local duration=$3
    
    print_status "Running TUI test in $mode mode..."
    
    if go run ./cmd/test-tui -mode="$mode" -prompt="$prompt" -duration="$duration" -verbose; then
        print_success "$mode test completed successfully"
        return 0
    else
        print_error "$mode test failed"
        return 1
    fi
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  headless    - Run headless tests (no UI, fast)"
    echo "  automated   - Run automated tests (UI with scripted actions)"
    echo "  demo        - Run demo mode (showcase all features)"
    echo "  interactive - Run interactive mode (manual testing)"
    echo "  all         - Run all non-interactive tests"
    echo "  quick       - Run quick headless test"
    echo "  help        - Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 headless"
    echo "  $0 demo"
    echo "  $0 all"
}

# Main script logic
case "${1:-help}" in
    "headless")
        print_status "Starting headless TUI tests..."
        run_test "headless" "headless-test" "5s"
        ;;
    
    "automated")
        print_status "Starting automated TUI tests..."
        run_test "automated" "automated-test" "10s"
        ;;
    
    "demo")
        print_status "Starting TUI demo..."
        print_warning "This will show a live demo of the workbench features"
        run_test "demo" "demo-prompt" "15s"
        ;;
    
    "interactive")
        print_status "Starting interactive TUI test..."
        print_warning "This will open the full TUI - press 'q' or Ctrl+C to exit"
        run_test "interactive" "interactive-test" "60s"
        ;;
    
    "quick")
        print_status "Running quick headless test..."
        run_test "headless" "quick-test" "3s"
        ;;
    
    "all")
        print_status "Running all non-interactive TUI tests..."
        
        failed=0
        
        print_status "1/3 Running headless tests..."
        run_test "headless" "test-1" "5s" || failed=$((failed + 1))
        
        print_status "2/3 Running automated tests..."
        run_test "automated" "test-2" "8s" || failed=$((failed + 1))
        
        print_status "3/3 Running demo tests..."
        run_test "demo" "test-3" "12s" || failed=$((failed + 1))
        
        echo ""
        if [ $failed -eq 0 ]; then
            print_success "All TUI tests passed! 🎉"
        else
            print_error "$failed test(s) failed"
            exit 1
        fi
        ;;
    
    "help"|*)
        show_usage
        ;;
esac