# Local Development

## Overview

AgentForge's local development workflow enables developers to modify any component locally, test changes in isolation, and seamlessly contribute improvements back to the ecosystem. This creates a Git-like development experience for AI agent configurations.

## Development Workflow

### 1. Component Installation for Development

```bash
# Install component for development (creates local copy)
forge dev install github.com/company/crm-tools/salesforce-lookup@v2.1.0

# Install with development options
forge dev install salesforce-lookup \
  --editable \
  --with-dependencies \
  --create-workspace ./my-workspace

# Install from local composition
forge dev install --from-composition sales-stack --dev-mode
```

### 2. Local Editing and Modification

#### Interactive Editing
```bash
# Interactive component editing
forge dev edit salesforce-lookup
# Opens interactive editor with current configuration
# Allows field-by-field modification with validation

# Edit specific fields
forge dev edit salesforce-lookup --field config.timeout --value 60
forge dev edit salesforce-lookup --field schema.input.properties.limit.maximum --value 1000

# Edit with wizard mode
forge dev edit salesforce-lookup --wizard
# Guided editing with explanations and suggestions
```

#### File-Based Editing
```bash
# Edit component definition file
forge dev edit salesforce-lookup --file
# Opens component.yaml in $EDITOR

# Edit specific files
forge dev edit salesforce-lookup --edit-schema
forge dev edit salesforce-lookup --edit-config
forge dev edit salesforce-lookup --edit-template

# Import from external file
forge dev edit salesforce-lookup --import-config ./new-config.yaml
forge dev edit sales-prompt --import-template ./new-template.txt
```

#### Programmatic Editing
```bash
# Batch modifications using JSON/YAML
forge dev edit salesforce-lookup --patch '{"config": {"timeout": 60, "retry_count": 5}}'

# Apply modifications from file
forge dev edit salesforce-lookup --patch-file ./modifications.yaml

# Use jq-style queries for complex modifications
forge dev edit salesforce-lookup --jq '.config.timeout = 60 | .config.retry_count = 5'
```

### 3. Local Testing and Validation

#### Component Testing
```bash
# Test component with sample input
forge dev test salesforce-lookup
forge dev test salesforce-lookup --input '{"query_type": "account", "search_term": "test"}'
forge dev test salesforce-lookup --input-file ./test-cases.json

# Test with different configurations
forge dev test salesforce-lookup --config ./test-config.yaml
forge dev test salesforce-lookup --env-file ./test.env

# Batch testing
forge dev test --all-modified
forge dev test --composition sales-stack
forge dev test --pattern "*-lookup"
```

#### Validation and Linting
```bash
# Validate component definition
forge dev validate salesforce-lookup
forge dev validate salesforce-lookup --strict
forge dev validate --all --fix-issues

# Schema validation
forge dev validate salesforce-lookup --check-schema
forge dev validate salesforce-lookup --check-dependencies
forge dev validate salesforce-lookup --check-security

# Performance testing
forge dev benchmark salesforce-lookup
forge dev benchmark salesforce-lookup --iterations 100 --concurrency 10
```

#### Integration Testing
```bash
# Test component integration
forge dev test-integration salesforce-lookup
forge dev test-integration --composition sales-stack
forge dev test-integration --with-dependencies

# Test against different LLM providers
forge dev test salesforce-lookup --llm-provider openai
forge dev test salesforce-lookup --llm-provider anthropic --llm-model claude-3-sonnet
```

## Change Tracking and Management

### Change Detection
```go
type ChangeTracker struct {
    db     *database.Client
    fs     *filesystem.Watcher
    logger *zap.Logger
}

func (ct *ChangeTracker) TrackChanges(componentID string) error {
    component, err := ct.db.GetLocalComponent(ctx, componentID)
    if err != nil {
        return err
    }
    
    // Watch for file system changes
    watcher, err := ct.fs.Watch(component.LocalPath)
    if err != nil {
        return err
    }
    
    go func() {
        for event := range watcher.Events {
            if event.Op&fsnotify.Write == fsnotify.Write {
                ct.handleFileChange(componentID, event.Name)
            }
        }
    }()
    
    return nil
}

func (ct *ChangeTracker) handleFileChange(componentID, filePath string) {
    // Calculate diff
    diff, err := ct.calculateDiff(componentID, filePath)
    if err != nil {
        ct.logger.Error("Failed to calculate diff", zap.Error(err))
        return
    }
    
    // Record change
    change := &ComponentChange{
        ComponentID: componentID,
        ChangeType:  ChangeTypeUpdate,
        FilePath:    filePath,
        Diff:        diff,
        Timestamp:   time.Now(),
    }
    
    if err := ct.db.RecordChange(ctx, change); err != nil {
        ct.logger.Error("Failed to record change", zap.Error(err))
    }
}
```

### Change Visualization
```bash
# Show local changes
forge dev diff salesforce-lookup
# Output:
# Modified: components/tools/salesforce-lookup/component.yaml
# @@ -15,7 +15,7 @@
#    configuration:
#      parameters:
# -      timeout: 30
# +      timeout: 60
#        retry_count: 3
# 
# Modified: components/tools/salesforce-lookup/schema.json
# @@ -25,6 +25,9 @@
#        "limit": {
#          "type": "integer",
# -        "maximum": 100
# +        "maximum": 1000
#        }
# +      "offset": {
# +        "type": "integer",
# +        "default": 0
# +      }

# Show change summary
forge dev status
# Component Status:
# ┌─────────────────┬─────────┬─────────────┬──────────────┐
# │ Component       │ Status  │ Changes     │ Last Modified│
# ├─────────────────┼─────────┼─────────────┼──────────────┤
# │ salesforce-lookup│ modified│ 2 files     │ 5 min ago    │
# │ sales-system    │ modified│ 1 file      │ 1 hour ago   │
# │ email-composer  │ clean   │ -           │ -            │
# └─────────────────┴─────────┴─────────────┴──────────────┘

# Detailed change log
forge dev log salesforce-lookup
# Change History for salesforce-lookup:
# 2024-01-15 14:30:00 - Modified config.timeout (30 → 60)
# 2024-01-15 14:32:15 - Modified schema.input.properties.limit.maximum (100 → 1000)
# 2024-01-15 14:33:45 - Added schema.input.properties.offset
# 
# Total changes: 3
# Files modified: 2 (component.yaml, schema.json)
```

### Change Management
```bash
# Revert specific changes
forge dev revert salesforce-lookup --change-id abc123
forge dev revert salesforce-lookup --file schema.json
forge dev revert salesforce-lookup --field config.timeout

# Revert to specific state
forge dev revert salesforce-lookup --to-original
forge dev revert salesforce-lookup --to-commit def456
forge dev revert salesforce-lookup --to-timestamp "2024-01-15 14:00:00"

# Stash changes temporarily
forge dev stash salesforce-lookup --message "WIP: adding pagination"
forge dev stash pop salesforce-lookup
forge dev stash list
forge dev stash drop stash-id-123

# Create checkpoints
forge dev checkpoint salesforce-lookup --message "Working timeout configuration"
forge dev checkpoint list salesforce-lookup
forge dev restore salesforce-lookup --checkpoint checkpoint-456
```

## Local Workspace Management

### Workspace Structure
```
~/.forge/workspaces/
├── default/
│   ├── components/
│   │   ├── tools/
│   │   │   ├── salesforce-lookup/
│   │   │   │   ├── component.yaml
│   │   │   │   ├── schema.json
│   │   │   │   ├── server.go
│   │   │   │   └── .forge/
│   │   │   │       ├── metadata.yaml
│   │   │   │       ├── changes.log
│   │   │   │       └── checkpoints/
│   │   │   └── email-composer/
│   │   ├── prompts/
│   │   └── agents/
│   ├── compositions/
│   │   ├── sales-stack.yaml
│   │   └── support-stack.yaml
│   └── .forge/
│       ├── workspace.yaml
│       ├── dependencies.lock
│       └── cache/
└── project-alpha/
    └── ...
```

### Workspace Operations
```bash
# Create new workspace
forge workspace create project-alpha
forge workspace create project-alpha --from-template enterprise

# Switch workspaces
forge workspace switch project-alpha
forge workspace switch default

# List workspaces
forge workspace list
# Workspaces:
# * default (active)
#   project-alpha
#   experimental

# Workspace information
forge workspace info
# Workspace: default
# Path: ~/.forge/workspaces/default
# Components: 15 (3 modified)
# Compositions: 2
# Last activity: 5 minutes ago

# Clean workspace
forge workspace clean --unused-components
forge workspace clean --old-checkpoints --keep 5
forge workspace clean --cache

# Export/import workspace
forge workspace export --output ./workspace-backup.tar.gz
forge workspace import ./workspace-backup.tar.gz --name restored-workspace
```

## Development Tools and Utilities

### Component Generator
```bash
# Generate new component from template
forge dev generate tool my-new-tool --template basic
forge dev generate tool my-new-tool --template mcp-server --runtime go
forge dev generate prompt my-prompt --template system --category sales

# Generate from existing component
forge dev generate tool my-tool --from salesforce-lookup --customize

# Interactive generation
forge dev generate --interactive
# Component type: tool, prompt, agent
# Template: basic, mcp-server, function, webhook
# Runtime: go, python, node, rust
# Features: authentication, caching, monitoring
```

### Development Server
```bash
# Start development server for testing
forge dev server --port 8080
forge dev server --components salesforce-lookup,email-composer
forge dev server --composition sales-stack

# Server provides:
# - Component testing endpoints
# - Real-time change monitoring
# - Integration testing interface
# - Performance metrics
# - Log streaming
```

### Hot Reloading
```bash
# Enable hot reloading for development
forge dev watch salesforce-lookup
# Watching for changes to salesforce-lookup...
# Change detected: component.yaml
# Reloading component... ✓
# Running tests... ✓
# Component ready

# Watch multiple components
forge dev watch --all-modified
forge dev watch --composition sales-stack

# Watch with custom actions
forge dev watch salesforce-lookup --on-change "forge dev test {component}"
forge dev watch --pattern "*-lookup" --on-change "./custom-test.sh {component}"
```

## Collaboration Features

### Sharing Local Changes
```bash
# Create shareable patch
forge dev patch create salesforce-lookup --output ./my-changes.patch
forge dev patch create --all-modified --output ./workspace-changes.patch

# Apply patch from colleague
forge dev patch apply ./colleague-changes.patch
forge dev patch apply ./colleague-changes.patch --dry-run

# Share via temporary branch
forge dev share salesforce-lookup --create-branch feature/timeout-increase
forge dev share --all-modified --branch feature/workspace-improvements

# Create pull request draft
forge dev share salesforce-lookup --create-pr --draft \
  --title "Increase timeout for enterprise environments" \
  --body-file ./pr-description.md
```

### Code Review Integration
```bash
# Request review for local changes
forge dev review request salesforce-lookup \
  --reviewers john.doe,jane.smith \
  --message "Please review timeout configuration changes"

# Review colleague's changes
forge dev review show review-id-123
forge dev review approve review-id-123 --comment "LGTM"
forge dev review request-changes review-id-123 --comment "Please add tests"

# Apply reviewed changes
forge dev review apply review-id-123
```

## Performance and Optimization

### Local Caching
```go
type LocalCache struct {
    componentCache *cache.Cache
    testCache      *cache.Cache
    buildCache     *cache.Cache
}

func (lc *LocalCache) CacheTestResults(componentID string, results *TestResults) {
    key := fmt.Sprintf("test:%s:%s", componentID, results.ConfigHash)
    lc.testCache.Set(key, results, 1*time.Hour)
}

func (lc *LocalCache) GetCachedTestResults(componentID, configHash string) (*TestResults, bool) {
    key := fmt.Sprintf("test:%s:%s", componentID, configHash)
    if cached := lc.testCache.Get(key); cached != nil {
        return cached.(*TestResults), true
    }
    return nil, false
}
```

### Incremental Building
```bash
# Build only changed components
forge dev build --incremental
forge dev build salesforce-lookup --if-changed

# Parallel building
forge dev build --all --parallel --max-workers 4

# Build with caching
forge dev build --cache --cache-dir ~/.forge/build-cache
```

### Development Metrics
```bash
# Show development metrics
forge dev metrics
# Development Metrics:
# Components modified: 5
# Total changes: 23
# Tests run: 47 (45 passed, 2 failed)
# Build time: 2.3s (avg)
# Test time: 8.7s (avg)
# Cache hit rate: 78%

# Component-specific metrics
forge dev metrics salesforce-lookup
# salesforce-lookup Metrics:
# Changes: 8
# Last test: 5 minutes ago (passed)
# Last build: 3 minutes ago (success)
# Performance: 150ms avg response time
# Test coverage: 92%
```

## Integration with IDEs

### VS Code Extension
```json
{
  "name": "agentforge",
  "displayName": "AgentForge",
  "description": "AgentForge component development",
  "version": "0.2.0",
  "engines": {
    "vscode": "^1.74.0"
  },
  "categories": ["Other"],
  "activationEvents": [
    "workspaceContains:**/component.yaml",
    "workspaceContains:**/forge-manifest.yaml"
  ],
  "main": "./out/extension.js",
  "contributes": {
    "commands": [
      {
        "command": "agentforge.editComponent",
        "title": "Edit Component",
        "category": "AgentForge"
      },
      {
        "command": "agentforge.testComponent",
        "title": "Test Component",
        "category": "AgentForge"
      },
      {
        "command": "agentforge.syncComponent",
        "title": "Sync Component",
        "category": "AgentForge"
      }
    ],
    "languages": [
      {
        "id": "forge-component",
        "aliases": ["AgentForge Component", "component"],
        "extensions": [".component.yaml"],
        "configuration": "./language-configuration.json"
      }
    ],
    "grammars": [
      {
        "language": "forge-component",
        "scopeName": "source.forge-component",
        "path": "./syntaxes/forge-component.tmGrammar.json"
      }
    ]
  }
}
```

### IntelliJ Plugin
```kotlin
class AgentForgePlugin : Plugin {
    override fun initComponent() {
        // Register file types
        FileTypeManager.getInstance().registerFileType(
            ComponentFileType(), "component.yaml"
        )
        
        // Register actions
        ActionManager.getInstance().registerAction(
            "AgentForge.EditComponent",
            EditComponentAction()
        )
    }
}

class EditComponentAction : AnAction() {
    override fun actionPerformed(e: AnActionEvent) {
        val project = e.project ?: return
        val file = e.getData(CommonDataKeys.VIRTUAL_FILE) ?: return
        
        if (file.name.endsWith(".component.yaml")) {
            ComponentEditor.open(project, file)
        }
    }
}
```

## Best Practices

### Development Guidelines
1. **Small, Focused Changes**: Make incremental changes that are easy to review and test
2. **Test-Driven Development**: Write tests before implementing changes
3. **Documentation**: Update documentation alongside code changes
4. **Validation**: Always validate changes before committing
5. **Backup**: Create checkpoints before major changes

### Performance Tips
1. **Use Incremental Builds**: Only rebuild changed components
2. **Cache Test Results**: Avoid re-running unchanged tests
3. **Parallel Operations**: Use parallel testing and building when possible
4. **Local Dependencies**: Use local copies of dependencies for faster iteration
5. **Hot Reloading**: Use watch mode for rapid feedback

### Collaboration Best Practices
1. **Clear Commit Messages**: Describe what and why, not just what
2. **Small Pull Requests**: Easier to review and merge
3. **Review Early**: Share work-in-progress for early feedback
4. **Document Decisions**: Explain complex changes and trade-offs
5. **Test Integration**: Ensure changes work with dependent components