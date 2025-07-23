# Sync Engine

## Overview

The AgentForge Sync Engine provides bidirectional synchronization between local components and Git repositories, enabling seamless collaboration while maintaining strong versioning guarantees through commit-based linking.

## Architecture

### Core Components

```
┌─────────────────────────────────────────────────────────────┐
│                    Sync Engine                             │
├─────────────────────────────────────────────────────────────┤
│  Sync Coordinator                                          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
│  │ Operation   │ │ Conflict    │ │ Status      │          │
│  │ Manager     │ │ Resolver    │ │ Tracker     │          │
│  └─────────────┘ └─────────────┘ └─────────────┘          │
├─────────────────────────────────────────────────────────────┤
│  Git Integration Layer                                     │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
│  │ Repository  │ │ Branch      │ │ Commit      │          │
│  │ Manager     │ │ Manager     │ │ Tracker     │          │
│  └─────────────┘ └─────────────┘ └─────────────┘          │
├─────────────────────────────────────────────────────────────┤
│  GitHub Integration                                        │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
│  │ API Client  │ │ Fork        │ │ PR          │          │
│  │             │ │ Manager     │ │ Manager     │          │
│  └─────────────┘ └─────────────┘ └─────────────┘          │
├─────────────────────────────────────────────────────────────┤
│  Local Storage                                             │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
│  │ Component   │ │ Change      │ │ Sync        │          │
│  │ Database    │ │ Tracker     │ │ Operations  │          │
│  └─────────────┘ └─────────────┘ └─────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

## Sync Operations

### 1. Pull Operations (Upstream → Local)

#### Simple Pull (No Local Changes)
```go
type PullOperation struct {
    ComponentID    string
    SourceRepo     string
    SourceCommit   string
    TargetCommit   string
    Strategy       PullStrategy
}

type PullStrategy string

const (
    PullStrategyFastForward PullStrategy = "fast_forward"
    PullStrategyMerge      PullStrategy = "merge"
    PullStrategyRebase     PullStrategy = "rebase"
    PullStrategyOverwrite  PullStrategy = "overwrite"
)
```

#### Pull Workflow
```bash
# 1. Check current state
forge sync status salesforce-lookup
# Output: Component 'salesforce-lookup' is outdated (local: abc123, upstream: def456)

# 2. Pull latest changes
forge sync pull salesforce-lookup
# Fetching latest from github.com/company/crm-tools...
# Updating salesforce-lookup: abc123 → def456
# ✓ Component updated successfully

# 3. Verify update
forge sync status salesforce-lookup
# Output: Component 'salesforce-lookup' is synced (commit: def456)
```

#### Pull with Local Changes
```bash
# Component has local modifications
forge sync pull salesforce-lookup
# Warning: Component has local changes
# Options:
#   1. Stash changes and pull: forge sync pull salesforce-lookup --stash
#   2. Merge with upstream: forge sync pull salesforce-lookup --merge
#   3. Rebase on upstream: forge sync pull salesforce-lookup --rebase
#   4. Overwrite local changes: forge sync pull salesforce-lookup --overwrite

# Choose merge strategy
forge sync pull salesforce-lookup --merge
# Merging upstream changes with local modifications...
# Auto-merge successful for: component.yaml
# Conflict in: schema.json
# Run 'forge sync resolve salesforce-lookup' to resolve conflicts
```

### 2. Push Operations (Local → Upstream/Fork)

#### Direct Push (Write Access)
```go
type PushOperation struct {
    ComponentID     string
    TargetRepo      string
    BranchName      string
    CommitMessage   string
    CreatePR        bool
    PRTitle         string
    PRBody          string
}
```

#### Push Workflow
```bash
# 1. Make local changes
forge dev edit salesforce-lookup --field config.timeout --value 60

# 2. Review changes
forge dev diff salesforce-lookup
# Modified: config.timeout (30 → 60)

# 3. Push changes
forge sync push salesforce-lookup --message "Increase timeout for enterprise use"
# Checking repository access...
# ✓ Write access confirmed for github.com/company/crm-tools
# Creating branch: forge/salesforce-lookup-1642089600
# Committing changes...
# Pushing to origin...
# ✓ Changes pushed successfully (commit: ghi789)
```

#### Fork-Based Push (No Write Access)
```bash
# Push to fork when no write access
forge sync push salesforce-lookup --create-fork
# No write access to github.com/company/crm-tools
# Creating fork: github.com/myuser/crm-tools
# ✓ Fork created successfully
# Creating branch: forge/salesforce-lookup-enhancement
# Committing changes...
# Pushing to fork...
# ✓ Changes pushed to fork

# Optionally create pull request
forge sync push salesforce-lookup --create-pr \
  --title "Increase timeout for enterprise environments" \
  --body-file ./pr-description.md
# ✓ Pull request created: https://github.com/company/crm-tools/pull/123
```

### 3. Conflict Resolution

#### Conflict Types
```go
type ConflictType string

const (
    ConflictTypeContent    ConflictType = "content"     // File content conflicts
    ConflictTypeSchema     ConflictType = "schema"      // Schema validation conflicts
    ConflictTypeVersion    ConflictType = "version"     // Version constraint conflicts
    ConflictTypeDependency ConflictType = "dependency"  // Dependency conflicts
)

type Conflict struct {
    ID           string
    ComponentID  string
    Type         ConflictType
    FilePath     string
    LocalContent string
    UpstreamContent string
    BaseContent  string
    Resolution   ConflictResolution
}
```

#### Conflict Resolution Strategies
```bash
# List conflicts
forge sync conflicts list
# Component: salesforce-lookup
# Conflicts:
#   1. content: components/tools/salesforce-lookup/component.yaml
#   2. schema: components/tools/salesforce-lookup/schema.json

# Interactive resolution
forge sync resolve salesforce-lookup --interactive
# Conflict 1/2: component.yaml
# 
# <<<<<<< LOCAL
# timeout: 60
# retry_count: 5
# =======
# timeout: 45
# retry_count: 3
# max_retries: 3
# >>>>>>> UPSTREAM
# 
# Choose resolution:
#   1. Use local version
#   2. Use upstream version
#   3. Edit manually
#   4. Skip this conflict
# Choice: 3

# Manual editing opens in $EDITOR
# After editing, continue with next conflict

# Automatic resolution strategies
forge sync resolve salesforce-lookup --prefer-local
forge sync resolve salesforce-lookup --prefer-upstream
forge sync resolve salesforce-lookup --strategy merge
```

#### Three-Way Merge
```go
type MergeStrategy string

const (
    MergeStrategyOurs     MergeStrategy = "ours"      // Prefer local changes
    MergeStrategyTheirs   MergeStrategy = "theirs"    // Prefer upstream changes
    MergeStrategyUnion    MergeStrategy = "union"     // Combine both (when possible)
    MergeStrategyManual   MergeStrategy = "manual"    // Require manual resolution
)

func (se *SyncEngine) ResolveConflict(ctx context.Context, conflict *Conflict, strategy MergeStrategy) error {
    switch strategy {
    case MergeStrategyOurs:
        return se.applyLocalContent(conflict)
    case MergeStrategyTheirs:
        return se.applyUpstreamContent(conflict)
    case MergeStrategyUnion:
        return se.attemptUnionMerge(conflict)
    case MergeStrategyManual:
        return se.requestManualResolution(conflict)
    }
}
```

## Commit-Based Dependency Resolution

### Exact Commit Tracking
```yaml
# Component definition with exact commit references
dependencies:
  tools:
    - repository: "github.com/company/crm-tools"
      component: "salesforce-lookup"
      version_constraint: "^2.1.0"
      resolved_version: "2.1.3"
      commit: "a1b2c3d4e5f6789012345678901234567890abcd"
      resolved_at: "2024-01-15T10:30:00Z"
      
    - repository: "github.com/tools/communication"
      component: "email-composer"
      version_constraint: "~1.3.0"
      resolved_version: "1.3.2"
      commit: "b2c3d4e5f6789012345678901234567890abcdef"
      resolved_at: "2024-01-15T10:30:00Z"
```

### Dependency Resolution Algorithm
```go
type DependencyResolver struct {
    repoManager *RepositoryManager
    cache       *cache.Cache
    logger      *zap.Logger
}

func (dr *DependencyResolver) ResolveExactDependencies(ctx context.Context, deps []Dependency) (*ResolvedDependencySet, error) {
    resolved := &ResolvedDependencySet{
        Components: make(map[string]*ResolvedComponent),
        LockFile:   &LockFile{},
    }
    
    // Build dependency graph
    graph, err := dr.buildDependencyGraph(ctx, deps)
    if err != nil {
        return nil, fmt.Errorf("failed to build dependency graph: %w", err)
    }
    
    // Topological sort for resolution order
    sortedDeps, err := graph.TopologicalSort()
    if err != nil {
        return nil, fmt.Errorf("circular dependency detected: %w", err)
    }
    
    // Resolve each dependency to exact commit
    for _, dep := range sortedDeps {
        component, err := dr.resolveToExactCommit(ctx, dep)
        if err != nil {
            return nil, fmt.Errorf("failed to resolve %s: %w", dep.Name, err)
        }
        
        resolved.Components[dep.Name] = component
        resolved.LockFile.Dependencies[dep.Name] = &LockEntry{
            Repository: dep.Repository,
            Commit:     component.CommitHash,
            Version:    component.Version,
            ResolvedAt: time.Now(),
        }
    }
    
    return resolved, nil
}

func (dr *DependencyResolver) resolveToExactCommit(ctx context.Context, dep Dependency) (*ResolvedComponent, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("resolve:%s:%s:%s", dep.Repository, dep.Component, dep.VersionConstraint)
    if cached := dr.cache.Get(cacheKey); cached != nil {
        return cached.(*ResolvedComponent), nil
    }
    
    // Fetch available versions from repository
    versions, err := dr.repoManager.GetComponentVersions(ctx, dep.Repository, dep.Component)
    if err != nil {
        return nil, err
    }
    
    // Find best matching version
    bestVersion, err := dr.findBestVersion(dep.VersionConstraint, versions)
    if err != nil {
        return nil, err
    }
    
    // Get exact commit for this version
    commit, err := dr.repoManager.GetCommitForVersion(ctx, dep.Repository, bestVersion.Tag)
    if err != nil {
        return nil, err
    }
    
    // Load component definition at exact commit
    component, err := dr.loadComponentAtCommit(ctx, dep.Repository, dep.Component, commit)
    if err != nil {
        return nil, err
    }
    
    resolved := &ResolvedComponent{
        Name:       dep.Component,
        Repository: dep.Repository,
        Version:    bestVersion.Version,
        CommitHash: commit,
        Definition: component,
    }
    
    // Cache result
    dr.cache.Set(cacheKey, resolved, 24*time.Hour)
    
    return resolved, nil
}
```

## Fork Management

### Automatic Fork Creation
```go
type ForkManager struct {
    github *github.Client
    db     *database.Client
    logger *zap.Logger
}

func (fm *ForkManager) EnsureFork(ctx context.Context, originalRepo string) (*Fork, error) {
    // Check if fork already exists
    existingFork, err := fm.db.GetFork(ctx, originalRepo)
    if err == nil {
        return existingFork, nil
    }
    
    // Create new fork
    fork, err := fm.github.CreateFork(ctx, originalRepo)
    if err != nil {
        return nil, fmt.Errorf("failed to create fork: %w", err)
    }
    
    // Store fork information
    forkRecord := &Fork{
        OriginalRepo: originalRepo,
        ForkRepo:     fork.FullName,
        CreatedAt:    time.Now(),
    }
    
    if err := fm.db.StoreFork(ctx, forkRecord); err != nil {
        fm.logger.Warn("Failed to store fork information", zap.Error(err))
    }
    
    return forkRecord, nil
}
```

### Fork Synchronization
```bash
# Keep fork in sync with upstream
forge fork sync my-crm-tools --with-upstream
# Fetching upstream changes...
# Merging upstream/main into fork/main...
# ✓ Fork synchronized with upstream

# Sync all forks
forge fork sync --all
# Synchronizing 3 forks...
# ✓ my-crm-tools: up to date
# ✓ my-sales-tools: 2 commits behind, merged
# ✗ my-prompt-tools: merge conflict, manual resolution required
```

## Pull Request Management

### Automated PR Creation
```go
type PRManager struct {
    github *github.Client
    db     *database.Client
    logger *zap.Logger
}

func (pm *PRManager) CreatePullRequest(ctx context.Context, req CreatePRRequest) (*PullRequest, error) {
    // Validate request
    if err := pm.validatePRRequest(req); err != nil {
        return nil, err
    }
    
    // Create pull request
    pr, err := pm.github.CreatePullRequest(ctx, github.CreatePullRequestRequest{
        BaseRepo: req.BaseRepo,
        HeadRepo: req.HeadRepo,
        Branch:   req.Branch,
        Title:    req.Title,
        Body:     pm.generatePRBody(req),
        Draft:    req.Draft,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create pull request: %w", err)
    }
    
    // Store PR information
    prRecord := &PullRequest{
        ID:         pr.ID,
        URL:        pr.URL,
        Title:      pr.Title,
        BaseRepo:   req.BaseRepo,
        HeadRepo:   req.HeadRepo,
        Branch:     req.Branch,
        Status:     PRStatusOpen,
        CreatedAt:  time.Now(),
    }
    
    if err := pm.db.StorePullRequest(ctx, prRecord); err != nil {
        pm.logger.Warn("Failed to store PR information", zap.Error(err))
    }
    
    return prRecord, nil
}

func (pm *PRManager) generatePRBody(req CreatePRRequest) string {
    var body strings.Builder
    
    body.WriteString("## Changes\n\n")
    body.WriteString(req.Description)
    body.WriteString("\n\n")
    
    if len(req.Changes) > 0 {
        body.WriteString("### Modified Components\n\n")
        for _, change := range req.Changes {
            body.WriteString(fmt.Sprintf("- **%s**: %s\n", change.Component, change.Description))
        }
        body.WriteString("\n")
    }
    
    body.WriteString("---\n")
    body.WriteString("*This pull request was created automatically by AgentForge*\n")
    
    return body.String()
}
```

## Sync Status Tracking

### Component Sync States
```go
type SyncStatus string

const (
    SyncStatusSynced    SyncStatus = "synced"     // Local matches upstream
    SyncStatusModified  SyncStatus = "modified"   // Local changes not pushed
    SyncStatusOutdated  SyncStatus = "outdated"   // Upstream has newer version
    SyncStatusConflict  SyncStatus = "conflict"   // Merge conflict exists
    SyncStatusForked    SyncStatus = "forked"     // Changes in user's fork
    SyncStatusDetached  SyncStatus = "detached"   // No upstream tracking
)

type ComponentSyncInfo struct {
    ComponentID    string
    Status         SyncStatus
    LocalCommit    string
    UpstreamCommit string
    LastSyncAt     time.Time
    HasLocalChanges bool
    ConflictCount   int
    ForkRepo        *string
}
```

### Status Reporting
```bash
# Comprehensive sync status
forge sync status --all
# Component Status Report
# ┌─────────────────┬─────────┬──────────┬─────────────┬──────────────┐
# │ Component       │ Status  │ Local    │ Upstream    │ Last Sync    │
# ├─────────────────┼─────────┼──────────┼─────────────┼──────────────┤
# │ salesforce-lookup│ synced  │ abc123   │ abc123      │ 2h ago       │
# │ sales-system    │ modified│ def456   │ def456      │ 1d ago       │
# │ email-composer  │ outdated│ ghi789   │ jkl012      │ 3d ago       │
# │ sales-specialist│ conflict│ mno345   │ pqr678      │ 5d ago       │
# └─────────────────┴─────────┴──────────┴─────────────┴──────────────┘
# 
# Summary:
# ✓ 1 synced, ⚠ 1 modified, ⬆ 1 outdated, ✗ 1 conflict

# Detailed status for specific component
forge sync status salesforce-lookup --detailed
# Component: salesforce-lookup
# Status: modified
# Repository: github.com/company/crm-tools
# Local commit: abc123 (2024-01-15 10:30:00)
# Upstream commit: abc123 (2024-01-15 10:30:00)
# Local changes: 2 files modified
#   - component.yaml: config.timeout (30 → 60)
#   - schema.json: added pagination fields
# Last sync: 2 hours ago
# Fork: github.com/myuser/crm-tools (if applicable)
```

## Performance Optimization

### Parallel Operations
```go
type SyncCoordinator struct {
    maxConcurrency int
    semaphore      chan struct{}
    wg             sync.WaitGroup
}

func (sc *SyncCoordinator) SyncMultipleComponents(ctx context.Context, components []string) error {
    sc.semaphore = make(chan struct{}, sc.maxConcurrency)
    
    for _, component := range components {
        sc.wg.Add(1)
        go func(comp string) {
            defer sc.wg.Done()
            sc.semaphore <- struct{}{} // Acquire
            defer func() { <-sc.semaphore }() // Release
            
            if err := sc.syncComponent(ctx, comp); err != nil {
                sc.logger.Error("Failed to sync component", 
                    zap.String("component", comp), 
                    zap.Error(err))
            }
        }(component)
    }
    
    sc.wg.Wait()
    return nil
}
```

### Incremental Sync
```go
func (se *SyncEngine) IncrementalSync(ctx context.Context, componentID string) error {
    component, err := se.db.GetComponent(ctx, componentID)
    if err != nil {
        return err
    }
    
    // Only sync if changes detected
    upstreamCommit, err := se.getLatestUpstreamCommit(ctx, component.RepositoryURL)
    if err != nil {
        return err
    }
    
    if upstreamCommit == component.UpstreamCommit {
        se.logger.Debug("Component already up to date", zap.String("component", componentID))
        return nil
    }
    
    // Perform incremental sync
    return se.performIncrementalSync(ctx, component, upstreamCommit)
}
```

### Caching Strategy
```go
type SyncCache struct {
    commitCache    *cache.Cache // Cache commit information
    contentCache   *cache.Cache // Cache file contents
    resolveCache   *cache.Cache // Cache dependency resolutions
}

func (sc *SyncCache) GetCachedCommit(repo, ref string) (*Commit, bool) {
    key := fmt.Sprintf("commit:%s:%s", repo, ref)
    if cached := sc.commitCache.Get(key); cached != nil {
        return cached.(*Commit), true
    }
    return nil, false
}

func (sc *SyncCache) CacheCommit(repo, ref string, commit *Commit) {
    key := fmt.Sprintf("commit:%s:%s", repo, ref)
    sc.commitCache.Set(key, commit, 1*time.Hour)
}
```