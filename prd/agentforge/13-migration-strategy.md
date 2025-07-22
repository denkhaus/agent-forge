# Migration Strategy

## Overview

This document outlines the comprehensive migration strategy from MCP Planner to AgentForge, ensuring a smooth transition while preserving existing functionality and enabling new capabilities. The migration follows a phased approach with backward compatibility and minimal disruption.

## Migration Phases

### Phase 1: Foundation and Compatibility (Weeks 1-4)
**Goal**: Establish AgentForge foundation while maintaining MCP Planner functionality

#### 1.1 Project Rename and Rebranding
```bash
# Repository migration
git clone https://github.com/denkhaus/mcp-planner agentforge
cd agentforge

# Update project metadata
find . -name "*.go" -exec sed -i 's/mcp-planner/agentforge/g' {} \;
find . -name "*.md" -exec sed -i 's/MCP Planner/AgentForge/g' {} \;
find . -name "*.yaml" -exec sed -i 's/mcp-planner/agentforge/g' {} \;

# Update module name
go mod edit -module github.com/denkhaus/agentforge
```

#### 1.2 CLI Rename and Alias
```go
// cmd/main.go - Support both names during transition
func main() {
    execName := filepath.Base(os.Args[0])
    
    var appName, usage string
    switch execName {
    case "mcp-planner", "mcp":
        appName = "mcp-planner"
        usage = "A Model Context Protocol (MCP) based planning system (legacy)"
        fmt.Fprintf(os.Stderr, "Warning: 'mcp-planner' is deprecated. Please use 'forge' instead.\n")
    case "forge":
        appName = "forge"
        usage = "AgentForge - Git-native AI agent development platform"
    default:
        appName = "forge"
        usage = "AgentForge - Git-native AI agent development platform"
    }
    
    app := createApp(appName, usage)
    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
```

#### 1.3 Database Schema Evolution
```sql
-- Add migration tracking table
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);

-- Migration 001: Add AgentForge tables alongside existing ones
-- Keep existing tables for backward compatibility
-- Add new tables with 'af_' prefix

-- New repository management
CREATE TABLE af_repositories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    url VARCHAR(500) NOT NULL,
    type VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- New component management
CREATE TABLE af_components (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    version VARCHAR(50) NOT NULL,
    repository_id UUID REFERENCES af_repositories(id),
    commit_hash VARCHAR(40) NOT NULL,
    definition JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Migration tracking
INSERT INTO schema_migrations (version, description) 
VALUES ('001', 'Add AgentForge foundation tables');
```

#### 1.4 Configuration Migration
```go
// internal/migration/config.go
type ConfigMigrator struct {
    oldConfig *mcpplanner.Config
    newConfig *agentforge.Config
}

func (cm *ConfigMigrator) MigrateConfig() error {
    // Map old configuration to new structure
    cm.newConfig = &agentforge.Config{
        // Direct mappings
        LogLevel:    cm.oldConfig.LogLevel,
        Port:        cm.oldConfig.Port,
        DatabaseURL: cm.oldConfig.DatabaseURL,
        Environment: cm.oldConfig.Environment,
        
        // API key mappings
        GoogleAPIKey:    cm.oldConfig.GoogleAPIKey,
        OpenAIAPIKey:    cm.oldConfig.OpenAIAPIKey,
        AnthropicAPIKey: cm.oldConfig.AnthropicAPIKey,
        
        // New AgentForge specific settings
        Git: GitConfig{
            Username: "", // Will be configured during setup
            Email:    "",
        },
        
        Repository: RepositoryConfig{
            CacheDir:     "~/.forge/cache",
            MaxCacheSize: "1GB",
        },
        
        Security: SecurityConfig{
            EnableSandboxing: false, // Disabled initially for compatibility
            TrustedDomains:   []string{"*"}, // Permissive initially
        },
    }
    
    return nil
}
```

### Phase 2: Component System Implementation (Weeks 5-8)
**Goal**: Implement Git-native component system with existing component migration

#### 2.1 Legacy Component Wrapper
```go
// internal/migration/components.go
type LegacyComponentWrapper struct {
    legacyTool   types.Tool
    legacyAgent  types.Agent
    legacyPrompt types.Prompt
}

func (lcw *LegacyComponentWrapper) WrapAsForgeComponent() *forge.Component {
    return &forge.Component{
        Metadata: forge.ComponentMetadata{
            Name:        lcw.extractName(),
            Version:     "1.0.0-legacy",
            Description: lcw.extractDescription(),
            Author:      "migrated-from-mcp-planner",
            Type:        lcw.determineType(),
        },
        
        Spec: forge.ComponentSpec{
            Definition: lcw.convertDefinition(),
            Dependencies: lcw.extractDependencies(),
        },
        
        Source: forge.ComponentSource{
            Type:       "legacy",
            Repository: "local://migrated",
            Commit:     "legacy",
        },
    }
}

func (lcw *LegacyComponentWrapper) convertDefinition() map[string]interface{} {
    switch {
    case lcw.legacyTool != nil:
        return map[string]interface{}{
            "type": "tool",
            "runtime": "embedded",
            "schema": lcw.legacyTool.Schema(),
            "implementation": "legacy_wrapper",
        }
    case lcw.legacyAgent != nil:
        return map[string]interface{}{
            "type": "agent",
            "llm": lcw.legacyAgent.GetLLMConfig(),
            "required_tools": lcw.legacyAgent.GetRequiredTools(),
        }
    case lcw.legacyPrompt != nil:
        return map[string]interface{}{
            "type": "prompt",
            "template": lcw.legacyPrompt.GetTemplate(),
            "variables": lcw.legacyPrompt.GetVariables(),
        }
    }
    return nil
}
```

#### 2.2 Hybrid Component Provider
```go
// internal/providers/hybrid.go
type HybridComponentProvider struct {
    legacyProvider types.ToolProvider
    forgeProvider  *forge.ComponentProvider
    migrator       *ComponentMigrator
}

func (hcp *HybridComponentProvider) GetComponent(ctx context.Context, name string) (*forge.Component, error) {
    // Try forge provider first
    if component, err := hcp.forgeProvider.GetComponent(ctx, name); err == nil {
        return component, nil
    }
    
    // Fall back to legacy provider
    if legacyTool, err := hcp.legacyProvider.GetTool(name); err == nil {
        // Migrate on-demand
        return hcp.migrator.MigrateTool(legacyTool), nil
    }
    
    return nil, fmt.Errorf("component %s not found", name)
}

func (hcp *HybridComponentProvider) ListComponents(ctx context.Context) ([]*forge.Component, error) {
    var components []*forge.Component
    
    // Get forge components
    forgeComponents, _ := hcp.forgeProvider.ListComponents(ctx)
    components = append(components, forgeComponents...)
    
    // Get legacy components and migrate
    legacyTools, _ := hcp.legacyProvider.GetTools()
    for _, tool := range legacyTools {
        migrated := hcp.migrator.MigrateTool(tool)
        components = append(components, migrated)
    }
    
    return components, nil
}
```

#### 2.3 Migration Commands
```bash
# Migration CLI commands
forge migrate init
# Initializes migration from MCP Planner
# - Scans existing configuration
# - Creates migration plan
# - Backs up existing data

forge migrate components
# Migrates existing components to AgentForge format
# - Converts tools, agents, prompts
# - Creates local repository structure
# - Maintains backward compatibility

forge migrate config
# Migrates configuration files
# - Updates environment variables
# - Converts configuration format
# - Preserves API keys and settings

forge migrate validate
# Validates migration completeness
# - Checks all components migrated
# - Verifies functionality
# - Reports any issues
```

### Phase 3: Git Integration and Repository System (Weeks 9-12)
**Goal**: Implement Git-native repository system with gradual migration

#### 3.1 Local Repository Creation
```go
// internal/migration/repository.go
type RepositoryMigrator struct {
    sourceDir   string
    targetRepo  string
    gitClient   *git.Client
}

func (rm *RepositoryMigrator) CreateLocalRepository() error {
    // Create local repository structure
    repoPath := filepath.Join(os.Getenv("HOME"), ".forge", "repositories", "migrated-components")
    
    // Initialize Git repository
    repo, err := rm.gitClient.Init(repoPath)
    if err != nil {
        return fmt.Errorf("failed to initialize repository: %w", err)
    }
    
    // Create standard structure
    dirs := []string{
        "components/tools",
        "components/prompts", 
        "components/agents",
        "examples",
        "docs",
    }
    
    for _, dir := range dirs {
        if err := os.MkdirAll(filepath.Join(repoPath, dir), 0755); err != nil {
            return err
        }
    }
    
    // Create manifest
    manifest := &forge.RepositoryManifest{
        APIVersion: "forge.dev/v1",
        Kind:       "ComponentRepository",
        Metadata: forge.ManifestMetadata{
            Name:        "migrated-components",
            Description: "Components migrated from MCP Planner",
            Version:     "1.0.0",
            Author:      "migration-tool",
        },
        Repository: forge.RepositoryInfo{
            Type:       "mixed",
            Categories: []string{"migrated", "legacy"},
        },
    }
    
    // Write manifest
    manifestPath := filepath.Join(repoPath, "forge-manifest.yaml")
    if err := rm.writeManifest(manifestPath, manifest); err != nil {
        return err
    }
    
    // Initial commit
    if err := repo.AddAll(); err != nil {
        return err
    }
    
    if err := repo.Commit("Initial migration from MCP Planner"); err != nil {
        return err
    }
    
    return nil
}
```

#### 3.2 Component Export
```go
func (rm *RepositoryMigrator) ExportComponents(components []*forge.Component) error {
    for _, component := range components {
        componentDir := filepath.Join(
            rm.targetRepo, 
            "components", 
            string(component.Metadata.Type)+"s", 
            component.Metadata.Name,
        )
        
        if err := os.MkdirAll(componentDir, 0755); err != nil {
            return err
        }
        
        // Write component definition
        componentFile := filepath.Join(componentDir, "component.yaml")
        if err := rm.writeComponentDefinition(componentFile, component); err != nil {
            return err
        }
        
        // Write additional files based on type
        switch component.Metadata.Type {
        case forge.ComponentTypeTool:
            if err := rm.exportToolFiles(componentDir, component); err != nil {
                return err
            }
        case forge.ComponentTypePrompt:
            if err := rm.exportPromptFiles(componentDir, component); err != nil {
                return err
            }
        case forge.ComponentTypeAgent:
            if err := rm.exportAgentFiles(componentDir, component); err != nil {
                return err
            }
        }
        
        // Write README
        readmeFile := filepath.Join(componentDir, "README.md")
        if err := rm.generateReadme(readmeFile, component); err != nil {
            return err
        }
    }
    
    return nil
}
```

### Phase 4: Advanced Features and Optimization (Weeks 13-16)
**Goal**: Enable advanced AgentForge features and optimize performance

#### 4.1 Sync System Implementation
```go
// internal/migration/sync.go
type SyncMigrator struct {
    localRepo   string
    remoteRepo  string
    syncEngine  *forge.SyncEngine
}

func (sm *SyncMigrator) EnableSyncCapabilities() error {
    // Create remote repository if needed
    if sm.remoteRepo != "" {
        if err := sm.createRemoteRepository(); err != nil {
            return err
        }
        
        // Push local repository to remote
        if err := sm.pushToRemote(); err != nil {
            return err
        }
    }
    
    // Configure sync settings
    syncConfig := &forge.SyncConfig{
        AutoSync:           false, // Manual initially
        ConflictResolution: forge.ConflictResolutionManual,
        SyncInterval:       time.Hour,
    }
    
    return sm.syncEngine.Configure(syncConfig)
}
```

#### 4.2 Security Migration
```go
// internal/migration/security.go
type SecurityMigrator struct {
    policyManager *forge.PolicyManager
    scanner       *forge.SecurityScanner
}

func (sm *SecurityMigrator) MigrateSecurity() error {
    // Create permissive policy for migrated components
    policy := &forge.SecurityPolicy{
        Name:        "migrated-components",
        Description: "Permissive policy for components migrated from MCP Planner",
        
        ComponentRequirements: forge.ComponentRequirements{
            Signing: forge.SigningRequirements{
                Required: false, // Disabled for migrated components
            },
            Scanning: forge.ScanningRequirements{
                Required:         true,
                MaxRiskScore:     10.0, // Permissive
                BlockOnHighSeverity: false,
            },
        },
        
        RuntimeRestrictions: forge.RuntimeRestrictions{
            Sandboxing: forge.SandboxingConfig{
                Required: false, // Disabled initially
            },
            ResourceLimits: forge.ResourceLimits{
                MaxCPU:    "unlimited",
                MaxMemory: "unlimited",
                MaxDisk:   "unlimited",
            },
        },
    }
    
    return sm.policyManager.CreatePolicy(policy)
}
```

## Migration Tools and Automation

### Migration CLI
```bash
#!/bin/bash
# scripts/migrate-to-agentforge.sh

set -e

echo "🚀 Starting MCP Planner to AgentForge migration..."

# Step 1: Backup existing installation
echo "📦 Creating backup..."
forge migrate backup --output ./mcp-planner-backup-$(date +%Y%m%d).tar.gz

# Step 2: Initialize AgentForge
echo "🏗️  Initializing AgentForge..."
forge migrate init --from mcp-planner

# Step 3: Migrate configuration
echo "⚙️  Migrating configuration..."
forge migrate config --validate

# Step 4: Migrate components
echo "🧩 Migrating components..."
forge migrate components --create-local-repo

# Step 5: Validate migration
echo "✅ Validating migration..."
forge migrate validate --comprehensive

# Step 6: Test functionality
echo "🧪 Testing migrated functionality..."
forge migrate test --all-components

echo "✨ Migration completed successfully!"
echo "📖 See migration report: ./migration-report.html"
echo "🔧 Next steps:"
echo "   1. Review migrated components: forge list"
echo "   2. Test your workflows: forge compose test"
echo "   3. Configure Git integration: forge config git"
echo "   4. Explore new features: forge help"
```

### Automated Migration Script
```go
// cmd/migrate/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "github.com/denkhaus/agentforge/internal/migration"
    "github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name:  "forge-migrate",
        Usage: "Migration tool for MCP Planner to AgentForge",
        Commands: []*cli.Command{
            {
                Name:  "init",
                Usage: "Initialize migration process",
                Action: handleInit,
                Flags: []cli.Flag{
                    &cli.StringFlag{
                        Name:  "from",
                        Value: "mcp-planner",
                        Usage: "Source system to migrate from",
                    },
                    &cli.StringFlag{
                        Name:  "config",
                        Usage: "Path to existing configuration",
                    },
                },
            },
            {
                Name:  "backup",
                Usage: "Create backup of existing installation",
                Action: handleBackup,
                Flags: []cli.Flag{
                    &cli.StringFlag{
                        Name:     "output",
                        Required: true,
                        Usage:    "Backup file path",
                    },
                },
            },
            {
                Name:  "components",
                Usage: "Migrate components",
                Action: handleComponents,
                Flags: []cli.Flag{
                    &cli.BoolFlag{
                        Name:  "create-local-repo",
                        Usage: "Create local Git repository for components",
                    },
                    &cli.StringFlag{
                        Name:  "repo-path",
                        Usage: "Path for local repository",
                    },
                },
            },
            {
                Name:  "validate",
                Usage: "Validate migration",
                Action: handleValidate,
                Flags: []cli.Flag{
                    &cli.BoolFlag{
                        Name:  "comprehensive",
                        Usage: "Run comprehensive validation",
                    },
                },
            },
        },
    }
    
    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

func handleInit(c *cli.Context) error {
    migrator := migration.NewMigrator()
    
    ctx := context.Background()
    config := &migration.Config{
        SourceSystem: c.String("from"),
        ConfigPath:   c.String("config"),
    }
    
    return migrator.Initialize(ctx, config)
}

func handleComponents(c *cli.Context) error {
    migrator := migration.NewMigrator()
    
    ctx := context.Background()
    config := &migration.ComponentMigrationConfig{
        CreateLocalRepo: c.Bool("create-local-repo"),
        RepoPath:        c.String("repo-path"),
    }
    
    return migrator.MigrateComponents(ctx, config)
}
```

## Backward Compatibility

### Legacy Command Support
```go
// internal/commands/legacy.go
type LegacyCommandHandler struct {
    forgeHandler *forge.CommandHandler
    logger       *zap.Logger
}

func (lch *LegacyCommandHandler) HandleLegacyCommand(ctx context.Context, cmd string, args []string) error {
    lch.logger.Warn("Using legacy command", 
        zap.String("command", cmd),
        zap.Strings("args", args))
    
    // Map legacy commands to new forge commands
    switch cmd {
    case "mcp-planner chat":
        return lch.forgeHandler.HandleChat(ctx, args)
    case "mcp-planner server":
        return lch.forgeHandler.HandleServer(ctx, args)
    case "mcp-planner version":
        return lch.forgeHandler.HandleVersion(ctx, args)
    default:
        return fmt.Errorf("unknown legacy command: %s", cmd)
    }
}
```

### Configuration Compatibility
```go
// internal/config/legacy.go
type LegacyConfigLoader struct {
    newLoader *forge.ConfigLoader
}

func (lcl *LegacyConfigLoader) LoadConfig(path string) (*forge.Config, error) {
    // Try to load as new format first
    if config, err := lcl.newLoader.Load(path); err == nil {
        return config, nil
    }
    
    // Fall back to legacy format
    legacyConfig, err := lcl.loadLegacyConfig(path)
    if err != nil {
        return nil, err
    }
    
    // Convert to new format
    return lcl.convertLegacyConfig(legacyConfig), nil
}
```

## Migration Validation

### Validation Framework
```go
// internal/migration/validation.go
type MigrationValidator struct {
    testRunner    *TestRunner
    configChecker *ConfigChecker
    dataVerifier  *DataVerifier
}

type ValidationResult struct {
    Success     bool
    Errors      []error
    Warnings    []string
    TestResults map[string]*TestResult
    Report      *ValidationReport
}

func (mv *MigrationValidator) ValidateMigration(ctx context.Context) (*ValidationResult, error) {
    result := &ValidationResult{
        Success:     true,
        TestResults: make(map[string]*TestResult),
    }
    
    // Validate configuration migration
    if err := mv.validateConfig(); err != nil {
        result.Errors = append(result.Errors, err)
        result.Success = false
    }
    
    // Validate component migration
    if err := mv.validateComponents(ctx); err != nil {
        result.Errors = append(result.Errors, err)
        result.Success = false
    }
    
    // Validate data integrity
    if err := mv.validateData(ctx); err != nil {
        result.Errors = append(result.Errors, err)
        result.Success = false
    }
    
    // Run functional tests
    testResults, err := mv.runFunctionalTests(ctx)
    if err != nil {
        result.Errors = append(result.Errors, err)
        result.Success = false
    }
    result.TestResults = testResults
    
    // Generate report
    result.Report = mv.generateReport(result)
    
    return result, nil
}
```

### Test Suite
```bash
# Migration test suite
forge migrate test --suite basic
# Tests:
# ✅ Configuration loading
# ✅ Component discovery
# ✅ Basic functionality
# ✅ API compatibility

forge migrate test --suite comprehensive
# Tests:
# ✅ All basic tests
# ✅ Performance benchmarks
# ✅ Security validation
# ✅ Integration tests
# ✅ Stress tests

forge migrate test --suite compatibility
# Tests:
# ✅ Legacy command support
# ✅ Configuration compatibility
# ✅ Data format compatibility
# ✅ API backward compatibility
```

## Rollback Strategy

### Rollback Procedures
```bash
# Emergency rollback
forge migrate rollback --to-backup ./mcp-planner-backup-20240115.tar.gz

# Partial rollback
forge migrate rollback --components-only
forge migrate rollback --config-only

# Validation before rollback
forge migrate rollback --dry-run --validate
```

### Rollback Implementation
```go
// internal/migration/rollback.go
type RollbackManager struct {
    backupPath   string
    currentState *SystemState
    logger       *zap.Logger
}

func (rm *RollbackManager) Rollback(ctx context.Context, options *RollbackOptions) error {
    rm.logger.Info("Starting rollback process", 
        zap.String("backup", rm.backupPath),
        zap.Any("options", options))
    
    // Create current state snapshot
    if err := rm.snapshotCurrentState(); err != nil {
        return fmt.Errorf("failed to snapshot current state: %w", err)
    }
    
    // Restore from backup
    if err := rm.restoreFromBackup(options); err != nil {
        return fmt.Errorf("failed to restore from backup: %w", err)
    }
    
    // Validate rollback
    if err := rm.validateRollback(); err != nil {
        rm.logger.Error("Rollback validation failed", zap.Error(err))
        return fmt.Errorf("rollback validation failed: %w", err)
    }
    
    rm.logger.Info("Rollback completed successfully")
    return nil
}
```

## Communication and Documentation

### Migration Guide
```markdown
# MCP Planner to AgentForge Migration Guide

## Overview
This guide walks you through migrating from MCP Planner to AgentForge.

## Prerequisites
- MCP Planner v0.1.x installed
- Go 1.21+ installed
- Git installed and configured

## Step-by-Step Migration

### 1. Backup Your Current Setup
```bash
forge migrate backup --output ./backup-$(date +%Y%m%d).tar.gz
```

### 2. Install AgentForge
```bash
go install github.com/denkhaus/agentforge/cmd/forge@latest
```

### 3. Initialize Migration
```bash
forge migrate init --from mcp-planner
```

### 4. Migrate Components
```bash
forge migrate components --create-local-repo
```

### 5. Validate Migration
```bash
forge migrate validate --comprehensive
```

## What Changes
- CLI command: `mcp-planner` → `forge`
- Configuration format: Enhanced with Git and security settings
- Component storage: File-based → Git-native repositories
- New features: Bidirectional sync, marketplace, security sandboxing

## What Stays the Same
- Core functionality preserved
- Existing API keys and configurations
- Component behavior and interfaces
- Database schemas (with additions)
```

### User Communication
```bash
# Migration announcement
forge announce migration
# Displays:
# 🎉 Welcome to AgentForge!
# 
# MCP Planner has evolved into AgentForge, bringing you:
# ✨ Git-native component management
# 🔄 Bidirectional sync with repositories  
# 🛡️ Enhanced security and sandboxing
# 🌍 Community marketplace
# 
# Your existing setup will continue to work.
# Run 'forge migrate init' to start using new features.
```

## Success Metrics

### Migration KPIs
- **Migration Success Rate**: >95% of users successfully migrate
- **Functionality Preservation**: 100% of existing functionality works
- **Migration Time**: <30 minutes for typical installation
- **User Satisfaction**: >4.5/5 rating for migration experience
- **Support Tickets**: <5% of users require migration support

### Monitoring and Feedback
```bash
# Migration telemetry (opt-in)
forge migrate telemetry --enable
forge migrate feedback --submit

# Migration analytics
forge migrate stats
# Shows:
# - Migration completion rate
# - Common issues
# - User feedback
# - Performance metrics
```