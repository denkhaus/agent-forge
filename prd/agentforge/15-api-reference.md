# API Reference

## Overview

This document provides comprehensive reference documentation for all AgentForge APIs, including CLI commands, programmatic interfaces, and configuration formats. The APIs are designed to be consistent, intuitive, and extensible.

## CLI Command Reference

### Global Flags

All `forge` commands support these global flags:

```bash
--config, -c     Configuration file path (default: ~/.forge/config.yaml)
--verbose, -v    Verbose output
--quiet, -q      Quiet mode (errors only)
--format         Output format: table, json, yaml (default: table)
--no-color       Disable colored output
--help, -h       Show help
--version        Show version information
```

### Core Commands

#### `forge init`
Initialize a new AgentForge workspace or project.

```bash
forge init [flags]

Flags:
  --database-url string    Database connection URL
  --git-username string    Git username for commits
  --git-email string       Git email for commits
  --workspace string       Workspace directory (default: current)
  --template string        Initialize from template

Examples:
  forge init
  forge init --database-url postgres://localhost/agentforge
  forge init --template enterprise --workspace ./my-project
```

#### `forge config`
Manage AgentForge configuration.

```bash
forge config <subcommand> [flags]

Subcommands:
  show                     Show current configuration
  set <key> <value>        Set configuration value
  get <key>                Get configuration value
  validate                 Validate configuration
  init                     Initialize configuration
  export                   Export configuration
  import                   Import configuration

Examples:
  forge config show
  forge config set git.username "John Doe"
  forge config get llm.default_provider
  forge config export --output config-backup.yaml
```

### Repository Management

#### `forge repo`
Manage component repositories.

```bash
forge repo <subcommand> [flags]

Subcommands:
  add <name> <url>         Add repository
  list                     List repositories
  info <name>              Show repository information
  update <name>            Update repository
  remove <name>            Remove repository
  sync <name>              Sync repository
  health <name>            Check repository health

Flags:
  --type string            Repository type filter
  --active-only           Show only active repositories
  --detailed              Show detailed information

Examples:
  forge repo add crm-tools github.com/company/crm-tools
  forge repo list --type tools
  forge repo sync crm-tools
  forge repo health --all
```

### Component Management

#### `forge search`
Search for components across repositories.

```bash
forge search <query> [flags]

Flags:
  --type string            Component type (tool, prompt, agent)
  --category string        Component category
  --tag string             Component tag
  --repo string            Repository filter
  --min-rating float       Minimum rating filter
  --stability string       Stability filter (experimental, beta, stable)
  --license string         License filter
  --updated-since string   Updated since filter (e.g., 30d, 1w)
  --limit int              Maximum results (default: 25)

Examples:
  forge search "crm integration"
  forge search --type tools --category sales
  forge search --repo github.com/company/tools --min-rating 4.0
```

#### `forge install`
Install components from repositories.

```bash
forge install <component-ref> [flags]

Flags:
  --version string         Specific version to install
  --resolve-deps          Resolve and install dependencies
  --force                 Force reinstall
  --no-cache             Skip cache
  --config string        Custom configuration file
  --path string          Installation path

Examples:
  forge install github.com/company/crm-tools/salesforce-lookup
  forge install salesforce-lookup@v2.1.0
  forge install salesforce-lookup --resolve-deps --force
```

#### `forge list`
List installed components.

```bash
forge list [flags]

Flags:
  --type string           Component type filter
  --repo string           Repository filter
  --status string         Status filter (installed, outdated, modified)
  --detailed             Show detailed information

Examples:
  forge list
  forge list --type tools --status outdated
  forge list --repo github.com/company/tools --detailed
```

#### `forge info`
Show detailed component information.

```bash
forge info <component> [flags]

Flags:
  --show-deps            Show dependencies
  --show-config          Show configuration
  --show-schema          Show schema
  --show-all             Show all information
  --versions             Show available versions
  --changelog            Show changelog

Examples:
  forge info salesforce-lookup
  forge info salesforce-lookup --show-deps --show-config
  forge info salesforce-lookup --versions
```

### Development Commands

#### `forge dev`
Local development operations.

```bash
forge dev <subcommand> [flags]

Subcommands:
  edit <component>         Edit component locally
  test <component>         Test component
  validate <component>     Validate component
  diff <component>         Show local changes
  log <component>          Show change history
  revert <component>       Revert changes
  stash <component>        Stash changes
  checkpoint <component>   Create checkpoint

Examples:
  forge dev edit salesforce-lookup --field config.timeout --value 60
  forge dev test salesforce-lookup --input '{"query": "test"}'
  forge dev diff salesforce-lookup
  forge dev checkpoint salesforce-lookup --message "Working configuration"
```

#### `forge workspace`
Manage development workspaces.

```bash
forge workspace <subcommand> [flags]

Subcommands:
  create <name>           Create workspace
  switch <name>           Switch workspace
  list                    List workspaces
  info [name]             Show workspace information
  clean                   Clean workspace
  export <name>           Export workspace
  import <file>           Import workspace

Examples:
  forge workspace create project-alpha
  forge workspace switch project-alpha
  forge workspace export project-alpha --output backup.tar.gz
```

### Sync and Collaboration

#### `forge sync`
Synchronize components with repositories.

```bash
forge sync <subcommand> [flags]

Subcommands:
  status [component]       Show sync status
  push <component>         Push local changes
  pull <component>         Pull remote changes
  resolve <component>      Resolve conflicts
  conflicts               List conflicts

Flags:
  --message string        Commit message for push
  --create-pr            Create pull request
  --create-fork          Create fork if no access
  --strategy string      Conflict resolution strategy
  --interactive          Interactive conflict resolution

Examples:
  forge sync status
  forge sync push salesforce-lookup --message "Increase timeout"
  forge sync pull salesforce-lookup --strategy merge
  forge sync resolve salesforce-lookup --interactive
```

#### `forge fork`
Manage repository forks.

```bash
forge fork <subcommand> [flags]

Subcommands:
  create <repo>           Create fork
  list                    List forks
  sync <fork>             Sync fork with upstream
  delete <fork>           Delete fork

Examples:
  forge fork create github.com/company/crm-tools
  forge fork sync my-crm-tools --with-upstream
  forge fork list --show-status
```

### Composition Management

#### `forge compose`
Manage agent compositions.

```bash
forge compose <subcommand> [flags]

Subcommands:
  create <name>           Create composition
  list                    List compositions
  info <name>             Show composition information
  edit <name>             Edit composition
  validate <name>         Validate composition
  deploy <name>           Deploy composition
  status <name>           Show deployment status
  logs <name>             Show logs
  stop <name>             Stop composition

Examples:
  forge compose create sales-stack --template enterprise
  forge compose deploy sales-stack --environment production
  forge compose logs sales-stack --follow
```

### Security Commands

#### `forge security`
Security operations and management.

```bash
forge security <subcommand> [flags]

Subcommands:
  scan <component>        Scan component for vulnerabilities
  policy list             List security policies
  policy apply <policy>   Apply security policy
  permissions <component> Show component permissions
  audit                   Show audit logs
  trust                   Manage trusted repositories

Examples:
  forge security scan salesforce-lookup
  forge security policy apply enterprise --to salesforce-lookup
  forge security audit --component salesforce-lookup --since 7d
```

### Marketplace Commands

#### `forge marketplace`
Interact with the component marketplace.

```bash
forge marketplace <subcommand> [flags]

Subcommands:
  search <query>          Search marketplace
  info <component>        Show marketplace information
  install <component>     Install from marketplace
  submit <component>      Submit component
  rate <component>        Rate component
  analytics <component>   Show analytics

Examples:
  forge marketplace search "crm tools"
  forge marketplace submit my-component --category tools
  forge marketplace rate salesforce-lookup --rating 5
```

## Configuration Reference

### Main Configuration File

The main configuration file is located at `~/.forge/config.yaml`:

```yaml
# Database configuration
database:
  url: "postgres://localhost/agentforge"
  auto_migrate: true
  connection_pool_size: 10
  connection_timeout: "30s"

# Git configuration
git:
  username: "John Doe"
  email: "john@example.com"
  default_branch: "main"
  signing_key: ""

# Repository configuration
repositories:
  cache_dir: "~/.forge/cache"
  max_cache_size: "5GB"
  sync_interval: "1h"
  default_timeout: "30s"
  auto_update: false

# LLM configuration
llm:
  default_provider: "openai"
  default_model: "gpt-4"
  default_temperature: 0.7
  default_max_tokens: 2048

# CLI configuration
cli:
  output_format: "table"  # table, json, yaml
  verbose: false
  color: true
  pager: "less"

# Security configuration
security:
  enable_scanning: true
  default_policy: "standard"
  sandbox_runtime: "docker"
  audit_level: "info"
  trusted_repositories:
    - "github.com/agentforge-official/*"

# Development configuration
development:
  default_editor: "$EDITOR"
  auto_validate: true
  auto_test: false
  workspace_dir: "~/.forge/workspaces"
  hot_reload: false

# Sync configuration
sync:
  auto_sync: false
  conflict_resolution: "manual"
  create_fork_if_no_access: true
  default_commit_message: "forge: update {component}"

# Marketplace configuration
marketplace:
  auto_submit_public: false
  quality_threshold: 7.0
  include_analytics: true
  default_license: "MIT"

# Performance configuration
performance:
  cache_ttl: "24h"
  max_concurrent_operations: 5
  request_timeout: "30s"
  retry_attempts: 3
```

### Environment Variables

AgentForge supports configuration via environment variables:

```bash
# Database
FORGE_DATABASE_URL="postgres://localhost/agentforge"

# Git
FORGE_GIT_USERNAME="John Doe"
FORGE_GIT_EMAIL="john@example.com"

# LLM API Keys
FORGE_OPENAI_API_KEY="sk-..."
FORGE_ANTHROPIC_API_KEY="..."
FORGE_GOOGLE_API_KEY="..."

# Security
FORGE_SECURITY_POLICY="enterprise"
FORGE_AUDIT_LEVEL="debug"

# Development
FORGE_WORKSPACE_DIR="/custom/workspace"
FORGE_EDITOR="code"

# Performance
FORGE_CACHE_TTL="12h"
FORGE_MAX_CONCURRENT="10"
```

## Component Definition Schemas

### Tool Component Schema

```yaml
apiVersion: "forge.dev/v1"
kind: "Tool"
metadata:
  name: string                 # Required: Component name
  version: string              # Required: Semantic version
  description: string          # Required: Human-readable description
  author: string               # Required: Author or organization
  license: string              # Required: SPDX license identifier
  homepage: string             # Optional: Project homepage
  documentation: string        # Optional: Documentation URL
  tags: [string]              # Optional: Searchable tags
  categories: [string]        # Optional: Hierarchical categories
  stability: string           # Optional: experimental|beta|stable|deprecated

spec:
  type: string                # Required: tool type (mcp-server, function, etc.)
  runtime: string             # Required: runtime environment
  
  execution:
    binary: string            # Required: executable path
    args: [string]           # Optional: command arguments
    env_file: string         # Optional: environment file
    working_directory: string # Optional: working directory
    timeout: integer         # Optional: execution timeout (seconds)
    memory_limit: string     # Optional: memory limit (e.g., "256MB")
    cpu_limit: string        # Optional: CPU limit (e.g., "0.5")
  
  schema:
    input:                   # Required: JSON schema for input
      type: object
      properties: {}
      required: []
    output:                  # Required: JSON schema for output
      type: object
      properties: {}
      required: []
  
  configuration:
    environment:             # Optional: environment variables
      - name: string
        required: boolean
        secret: boolean
        description: string
        default: string
    parameters:              # Optional: configuration parameters
      param_name:
        type: string
        default: any
        description: string
  
  capabilities: [string]     # Optional: capability tags
  
  dependencies:
    external: [string]       # Optional: external dependencies
    forge_components: []     # Optional: other forge components
    system: [string]         # Optional: system requirements
  
  security:
    permissions: [string]    # Optional: required permissions
    sandboxed: boolean       # Optional: requires sandboxing
    trusted_domains: [string] # Optional: allowed domains
```

### Prompt Component Schema

```yaml
apiVersion: "forge.dev/v1"
kind: "Prompt"
metadata:
  name: string
  version: string
  description: string
  author: string
  license: string
  tags: [string]
  categories: [string]

spec:
  type: string              # Required: system|user|assistant|function
  category: string          # Required: domain category
  
  template:
    engine: string          # Required: go-template|jinja2|handlebars
    file: string           # Required: template file path
    syntax_version: string  # Optional: template syntax version
  
  variables:               # Optional: template variables
    - name: string
      type: string
      required: boolean
      description: string
      default: any
      validation: {}
  
  validation:
    max_length: integer
    min_length: integer
    required_sections: [string]
    forbidden_content: [string]
  
  localization:
    supported_languages: [string]
    default_language: string
    translation_files: {}
  
  performance:
    token_efficiency: number
    response_quality: number
    consistency_score: number
```

### Agent Component Schema

```yaml
apiVersion: "forge.dev/v1"
kind: "Agent"
metadata:
  name: string
  version: string
  description: string
  author: string
  license: string
  tags: [string]
  categories: [string]

spec:
  dependencies:
    tools:                  # Required: tool dependencies
      - repository: string
        component: string
        version: string
        required: boolean
        purpose: string
    prompts:               # Required: prompt dependencies
      - repository: string
        component: string
        version: string
        role: string
        required: boolean
    agents: []             # Optional: agent dependencies
  
  llm:
    primary:
      provider: string
      model: string
      temperature: number
      max_tokens: integer
    fallback:
      provider: string
      model: string
  
  behavior:
    execution_mode: string  # agent|direct|hybrid
    max_iterations: integer
    timeout: integer
    memory_limit: string
  
  capabilities:
    primary: [string]
    secondary: [string]
    languages: [string]
  
  configuration:
    customizable_fields: []
    presets: []
```

## Programmatic API

### Go SDK

```go
package agentforge

import (
    "context"
    "github.com/denkhaus/agentforge/pkg/client"
)

// Client provides programmatic access to AgentForge
type Client struct {
    config *Config
    repos  *RepositoryManager
    comps  *ComponentManager
    sync   *SyncManager
}

// NewClient creates a new AgentForge client
func NewClient(config *Config) (*Client, error)

// Repository operations
func (c *Client) AddRepository(ctx context.Context, name, url string) error
func (c *Client) ListRepositories(ctx context.Context) ([]*Repository, error)
func (c *Client) SyncRepository(ctx context.Context, name string) error

// Component operations
func (c *Client) SearchComponents(ctx context.Context, query *SearchQuery) ([]*Component, error)
func (c *Client) InstallComponent(ctx context.Context, ref *ComponentRef) error
func (c *Client) ListComponents(ctx context.Context) ([]*Component, error)

// Development operations
func (c *Client) EditComponent(ctx context.Context, name string, changes map[string]interface{}) error
func (c *Client) TestComponent(ctx context.Context, name string, input interface{}) (*TestResult, error)
func (c *Client) ValidateComponent(ctx context.Context, name string) (*ValidationResult, error)

// Sync operations
func (c *Client) SyncComponent(ctx context.Context, name string, direction SyncDirection) error
func (c *Client) ResolveConflicts(ctx context.Context, name string, strategy ConflictStrategy) error
```

### REST API

AgentForge provides a REST API for integration with external systems:

#### Authentication
```http
Authorization: Bearer <api-token>
Content-Type: application/json
```

#### Endpoints

**Repository Management**
```http
GET    /api/v1/repositories
POST   /api/v1/repositories
GET    /api/v1/repositories/{id}
PUT    /api/v1/repositories/{id}
DELETE /api/v1/repositories/{id}
POST   /api/v1/repositories/{id}/sync
```

**Component Management**
```http
GET    /api/v1/components
POST   /api/v1/components/search
GET    /api/v1/components/{id}
POST   /api/v1/components/{id}/install
PUT    /api/v1/components/{id}
DELETE /api/v1/components/{id}
POST   /api/v1/components/{id}/test
POST   /api/v1/components/{id}/validate
```

**Sync Operations**
```http
GET    /api/v1/sync/status
POST   /api/v1/sync/push
POST   /api/v1/sync/pull
POST   /api/v1/sync/resolve
GET    /api/v1/sync/conflicts
```

### WebSocket API

Real-time updates via WebSocket:

```javascript
const ws = new WebSocket('wss://api.agentforge.dev/ws');

// Subscribe to events
ws.send(JSON.stringify({
  type: 'subscribe',
  events: ['component.updated', 'sync.completed', 'test.finished']
}));

// Handle events
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  switch (data.type) {
    case 'component.updated':
      handleComponentUpdate(data.payload);
      break;
    case 'sync.completed':
      handleSyncCompleted(data.payload);
      break;
  }
};
```

## Error Codes and Handling

### CLI Exit Codes
- `0`: Success
- `1`: General error
- `2`: Command line usage error
- `3`: Configuration error
- `4`: Network error
- `5`: Permission error
- `6`: Validation error
- `7`: Sync conflict
- `8`: Component not found
- `9`: Repository error
- `10`: Security error

### API Error Responses

```json
{
  "error": {
    "code": "COMPONENT_NOT_FOUND",
    "message": "Component 'salesforce-lookup' not found",
    "details": {
      "component": "salesforce-lookup",
      "repository": "github.com/company/crm-tools",
      "suggestions": [
        "salesforce-connector",
        "crm-lookup"
      ]
    },
    "timestamp": "2024-01-15T10:30:00Z",
    "request_id": "req_123456"
  }
}
```

### Common Error Codes
- `COMPONENT_NOT_FOUND`: Component does not exist
- `REPOSITORY_UNREACHABLE`: Cannot access repository
- `DEPENDENCY_CONFLICT`: Dependency version conflict
- `VALIDATION_FAILED`: Component validation failed
- `PERMISSION_DENIED`: Insufficient permissions
- `SYNC_CONFLICT`: Merge conflict during sync
- `RATE_LIMIT_EXCEEDED`: API rate limit exceeded
- `AUTHENTICATION_FAILED`: Invalid credentials
- `CONFIGURATION_INVALID`: Invalid configuration

This API reference provides comprehensive documentation for all AgentForge interfaces, enabling developers to integrate and extend the platform effectively.