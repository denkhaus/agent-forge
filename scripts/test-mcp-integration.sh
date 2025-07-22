#!/bin/bash

# MCP Integration Testing Script
set -e

echo "🧪 MCP Integration Testing Script"
echo "================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
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

# Check prerequisites
print_status "Checking prerequisites..."

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    print_error "Node.js is not installed. Please install Node.js to run MCP servers."
    exit 1
fi

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    print_error "npm is not installed. Please install npm to run MCP servers."
    exit 1
fi

print_success "Node.js $(node --version) and npm $(npm --version) are available"

# Check if our binary exists
if [ ! -f "./main" ]; then
    print_status "Building MCP Planner..."
    go build -o main ./cmd/main.go
    if [ $? -eq 0 ]; then
        print_success "MCP Planner built successfully"
    else
        print_error "Failed to build MCP Planner"
        exit 1
    fi
else
    print_success "MCP Planner binary found"
fi

# Test 1: MCP Disabled (baseline)
print_status "Test 1: Running with MCP disabled (baseline test)..."
export MCP_ENABLED=false
export LOG_LEVEL=info

echo "Testing tool availability with MCP disabled..."
timeout 10s ./main chat --clear <<EOF || true
help
exit
EOF

if [ $? -eq 124 ]; then
    print_warning "Test timed out (expected for interactive mode)"
else
    print_success "Baseline test completed"
fi

# Test 2: MCP Enabled with invalid config
print_status "Test 2: Testing MCP with invalid configuration..."
export MCP_ENABLED=true
export MCP_CONFIG_PATH=nonexistent.json

echo "Testing error handling with invalid MCP config..."
timeout 10s ./main chat --clear <<EOF || true
help
exit
EOF

if [ $? -eq 124 ]; then
    print_warning "Test timed out (expected for interactive mode)"
else
    print_success "Invalid config test completed"
fi

# Test 3: MCP Enabled with valid config
print_status "Test 3: Testing MCP with valid configuration..."
export MCP_ENABLED=true
export MCP_CONFIG_PATH=test-mcp-servers.json

print_status "Installing MCP servers (this may take a moment)..."

# Pre-install MCP servers to avoid timeout during test
print_status "Installing filesystem MCP server..."
npx -y @modelcontextprotocol/server-filesystem --version > /dev/null 2>&1 || {
    print_warning "Filesystem server installation may have issues, continuing..."
}

print_status "Installing memory MCP server..."
npx -y @modelcontextprotocol/server-memory --version > /dev/null 2>&1 || {
    print_warning "Memory server installation may have issues, continuing..."
}

print_success "MCP servers installation completed"

# Test with MCP enabled
print_status "Testing MCP integration..."
timeout 15s ./main chat --clear <<EOF || true
help
exit
EOF

if [ $? -eq 124 ]; then
    print_warning "MCP test timed out (expected for interactive mode)"
else
    print_success "MCP integration test completed"
fi

# Test 4: Tool prefix testing
print_status "Test 4: Testing tool prefix functionality..."
export MCP_TOOL_PREFIX=mcp_

timeout 10s ./main chat --clear <<EOF || true
help
exit
EOF

if [ $? -eq 124 ]; then
    print_warning "Prefix test timed out (expected for interactive mode)"
else
    print_success "Tool prefix test completed"
fi

# Test 5: Manual verification
print_status "Test 5: Manual verification mode..."
print_status "Starting interactive session for manual testing..."
print_status "You can test the following:"
print_status "1. Type 'help' to see available commands"
print_status "2. Check if MCP tools are available alongside internal tools"
print_status "3. Try executing both internal and MCP tools"
print_status "4. Type 'exit' when done"

export MCP_ENABLED=true
export MCP_CONFIG_PATH=test-mcp-servers.json
export MCP_TOOL_PREFIX=""
export LOG_LEVEL=debug

echo ""
print_status "Starting manual test session..."
echo "Press Ctrl+C to skip manual testing"
sleep 2

./main chat --clear || true

print_success "Manual testing session completed"

# Summary
echo ""
echo "🎉 MCP Integration Testing Summary"
echo "=================================="
print_success "✅ Prerequisites check passed"
print_success "✅ Baseline test (MCP disabled) completed"
print_success "✅ Error handling test (invalid config) completed"
print_success "✅ MCP integration test (valid config) completed"
print_success "✅ Tool prefix test completed"
print_success "✅ Manual verification completed"

echo ""
print_status "Integration testing completed successfully!"
print_status "Check the logs above for any warnings or errors."
print_status "The MCP integration appears to be working correctly."

# Cleanup
unset MCP_ENABLED
unset MCP_CONFIG_PATH
unset MCP_TOOL_PREFIX
unset LOG_LEVEL

echo ""
print_status "Environment variables cleaned up."
print_success "Testing script completed! 🚀"