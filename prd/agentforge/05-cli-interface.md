# MVP CLI Interface

## Overview

The `forge` CLI provides 11 core commands for simple, delightful component management. Built with `urfave/cli` and enhanced with `bubbletea` for interactive experiences.

## Command Structure

```
forge [global-flags] <command> [command-flags] [arguments]
```

### Global Flags
```bash
--config, -c     Configuration file path (default: ~/.forge/config.yaml)
--verbose, -v    Verbose output
--quiet, -q      Quiet mode (errors only)
--help, -h       Show help
--version        Show version information
```

## MVP Core Commands (11 Total)

### 1. System Commands

#### `forge init` - Initialize local forge system
```bash
# Initialize forge in current directory
forge init

# Creates:
# - ~/.forge/ directory structure
# - PostgreSQL database
# - Default configuration
# - Component directories (tools/, prompts/, agents/)
```

#### `forge config` - Interactive configuration
```bash
# Opens bubbletea TUI for configuration
forge config

# Configure specific values directly
forge config --github-token <token>
forge config --github-username <username>
```

#### `forge lint` - Validate local components
```bash
# Check all components for issues
forge lint

# Check specific component type
forge lint --type tool
forge lint --type prompt
forge lint --type agent

# Example output:
# ✓ weather-tool: Valid
# ✗ broken-prompt: Missing required field 'template'
# ⚠ old-agent: Uses deprecated format
```

### 2. Component Management (8 Commands)

#### `forge component new` - Create new component
```bash
# Create new component with interactive prompts
forge component new --type tool
forge component new --type prompt
forge component new --type agent

# Create with name
forge component new --type tool weather-checker

# Creates template file in ~/.forge/components/{type}/{name}.yaml
```

#### `forge component rm` - Remove component
```bash
# Remove component from local database
forge component rm --type tool weather-checker
forge component rm --type prompt system-prompt

# Confirmation prompt before deletion
# Removes from database but keeps local file
```

#### `forge component pull` - Pull component from GitHub
```bash
# Pull specific component
forge component pull --type tool github.com/user/repo/weather-tool

# Pull with specific version/tag
forge component pull --type tool github.com/user/repo/weather-tool@v1.2.0

# Pull and rename locally
forge component pull --type tool github.com/user/repo/weather-tool --as my-weather-tool
```

#### `forge component push` - Push component to GitHub
```bash
# Push component to GitHub (creates repo if needed)
forge component push --type tool weather-checker

# Push to specific repository
forge component push --type tool weather-checker --repo my-components

# Push with custom message
forge component push --type tool weather-checker --message "Add weather checking functionality"
```

#### `forge component status` - Show component status
```bash
# Show status of specific component
forge component status --type tool weather-checker

# Show status of all components
forge component status

# Example output:
# weather-checker (tool)
# ├─ Local: ✓ Present, modified 2 hours ago
# ├─ Remote: ✓ github.com/user/weather-tools
# ├─ Sync: ⚠ Local changes not pushed
# └─ Validation: ✓ Valid
```

#### `forge component sync` - Bidirectional sync
```bash
# Sync specific component
forge component sync --type tool weather-checker

# Sync all components
forge component sync

# Sync with conflict resolution
forge component sync --type tool weather-checker --resolve local
forge component sync --type tool weather-checker --resolve remote
```

#### `forge component ls` - List available components
```bash
# List all local components
forge component ls

# List specific type
forge component ls --type tool
forge component ls --type prompt
forge component ls --type agent

# List remote components (GitHub search)
forge component ls --type tool --remote
forge component ls --type tool --remote --min-stars 10

# Interactive TUI for browsing
forge component ls --interactive
```

## Interactive TUI Features

### Configuration TUI (`forge config`)
```
┌─ Forge Configuration ─────────────────────────────────┐
│                                                       │
│  GitHub Token:     [●●●●●●●●●●●●●●●●●●●●] ✓           │
│  GitHub Username:  [john-doe              ] ✓           │
│  Database URL:     [postgres://localhost/forge] ✓       │
│  Components Dir:   [~/.forge/components   ] ✓           │
│  Log Level:        [info ▼                ] ✓           │
│                                                       │
│  [Save] [Cancel] [Test Connection]                    │
└───────────────────────────────────────────────────────┘
```

### Component Browser (`forge component ls --interactive`)
```
┌─ Component Browser ───────────────────────────────────┐
│ Filter: [weather] Type: [tool ▼] Source: [remote ▼]  │
├───────────────────────────────────────────────────────┤
│ ✓ weather-api-tool        ⭐ 45  📅 2 days ago       │
│ ✓ openweather-connector   ⭐ 23  📅 1 week ago       │
│ ✓ weather-forecast        ⭐ 12  📅 3 weeks ago      │
│   simple-weather          ⭐ 8   📅 1 month ago      │
│   weather-alerts          ⭐ 3   📅 2 months ago     │
├───────────────────────────────────────────────────────┤
│ [Pull] [View Details] [Star] [q: Quit]               │
└───────────────────────────────────────────────────────┘
```

## Component File Structure

### Local Directory Layout
```
~/.forge/
├── config.yaml
├── forge.db (PostgreSQL)
└── components/
    ├── tools/
    │   ├── weather-checker.yaml
    │   └── calendar-tool.yaml
    ├── prompts/
    │   ├── system-prompt.yaml
    │   └── user-assistant.yaml
    └── agents/
        ├── customer-service.yaml
        └── code-reviewer.yaml
```

### Component YAML Format
```yaml
apiVersion: forge.dev/v1
kind: Tool
metadata:
  name: weather-checker
  description: "Get current weather information"
  version: "1.0.0"
  tags: ["weather", "api"]
  author: "github.com/username"
spec:
  mcp:
    server: "weather-server"
    tools: ["get_weather", "get_forecast"]
  config:
    api_key_required: true
    rate_limit: 100
  dependencies: []
```

## Error Handling

### User-Friendly Error Messages
```bash
# Missing GitHub token
$ forge component push --type tool weather-checker
forge: GITHUB_TOKEN_MISSING
GitHub token not configured

Hint: Run 'forge config' to set your GitHub token

# Component not found
$ forge component status --type tool nonexistent
forge: COMPONENT_NOT_FOUND
Component 'nonexistent' of type 'tool' not found

Hint: Use 'forge component ls --type tool' to see available components

# Network error
$ forge component pull --type tool github.com/user/repo/tool
forge: NETWORK_ERROR
Failed to connect to GitHub API

Hint: Check your internet connection and GitHub token permissions
```

## Performance Targets

- **CLI Response Time**: <1 second for all commands
- **Component Discovery**: <3 seconds for GitHub search
- **Local Operations**: <500ms for database queries
- **TUI Responsiveness**: <100ms for interface updates

## Testing Commands

```bash
# Test all CLI commands
make test-cli

# Test specific command
go test ./internal/commands -run TestComponentNew

# Integration test with GitHub
./scripts/test-github-integration.sh

# Performance test
./scripts/benchmark-cli.sh
```

This MVP CLI interface focuses on the essential 11 commands needed for delightful component sharing, with beautiful TUI interfaces and clear error handling.