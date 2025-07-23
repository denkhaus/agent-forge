# Dependency Resolution

## Overview

AgentForge's dependency resolution system ensures reproducible builds through commit-based linking while supporting semantic versioning for development flexibility. Every dependency is resolved to an exact Git commit, creating a deterministic dependency graph.

## Dependency Types

### Direct Dependencies
Components explicitly declared in a component's specification.

```yaml
# Agent component declaring direct dependencies
dependencies:
  tools:
    - repository: "github.com/company/crm-tools"
      component: "salesforce-lookup"
      version: "^2.1.0"
      required: true
      purpose: "Customer data retrieval"
  
  prompts:
    - repository: "github.com/ai-prompts/sales-templates"
      component: "sales-system-v3"
      version: ">=3.0.0"
      role: "system"
      required: true
```

### Transitive Dependencies
Dependencies of dependencies, automatically discovered and resolved.

```yaml
# Tool component with its own dependencies
dependencies:
  external:
    - name: "salesforce-api"
      version: ">=58.0"
  
  forge_components:
    - repository: "github.com/tools/auth"
      component: "oauth-handler"
      version: "~1.2.0"
```

### Peer Dependencies
Dependencies that must be provided by the consuming environment.

```yaml
# Component requiring peer dependencies
dependencies:
  peer:
    - repository: "github.com/tools/logging"
      component: "structured-logger"
      version: "^2.0.0"
      reason: "Required for audit logging"
```

## Version Constraints

### Semantic Version Constraints
```yaml
version_constraints:
  exact: "2.1.0"           # Exactly version 2.1.0
  caret: "^2.1.0"          # Compatible with 2.x.x (>=2.1.0, <3.0.0)
  tilde: "~2.1.0"          # Compatible with 2.1.x (>=2.1.0, <2.2.0)
  range: ">=2.1.0 <3.0.0"  # Explicit range
  latest: "latest"         # Latest stable version
  prerelease: "2.2.0-beta.1" # Specific prerelease
  wildcard: "2.*"          # Any 2.x version
```

### Git-Based Constraints
```yaml
git_constraints:
  branch: "main"           # Latest commit on branch
  tag: "v2.1.0"           # Specific Git tag
  commit: "abc123..."      # Exact commit hash
  ref: "refs/heads/feature" # Specific Git reference
```

## Resolution Algorithm

### Phase 1: Dependency Discovery
```go
type DependencyResolver struct {
    repoManager *RepositoryManager
    cache       *ResolutionCache
    logger      *zap.Logger
}

func (dr *DependencyResolver) DiscoverDependencies(ctx context.Context, rootComponent *Component) (*DependencyGraph, error) {
    graph := NewDependencyGraph()
    visited := make(map[string]bool)
    
    // Start with root component
    if err := dr.discoverRecursive(ctx, rootComponent, graph, visited, 0); err != nil {
        return nil, err
    }
    
    return graph, nil
}

func (dr *DependencyResolver) discoverRecursive(ctx context.Context, component *Component, graph *DependencyGraph, visited map[string]bool, depth int) error {
    componentKey := fmt.Sprintf("%s/%s", component.Repository, component.Name)
    
    // Prevent infinite recursion
    if visited[componentKey] {
        return nil
    }
    visited[componentKey] = true
    
    // Prevent excessive depth
    if depth > MaxDependencyDepth {
        return fmt.Errorf("dependency depth exceeded maximum of %d", MaxDependencyDepth)
    }
    
    // Add component to graph
    graph.AddNode(component)
    
    // Process each dependency
    for _, dep := range component.Dependencies {
        depComponent, err := dr.loadComponent(ctx, dep)
        if err != nil {
            return fmt.Errorf("failed to load dependency %s: %w", dep.Component, err)
        }
        
        // Add edge to graph
        graph.AddEdge(component, depComponent, dep)
        
        // Recurse into dependency
        if err := dr.discoverRecursive(ctx, depComponent, graph, visited, depth+1); err != nil {
            return err
        }
    }
    
    return nil
}
```

### Phase 2: Constraint Satisfaction
```go
type ConstraintSolver struct {
    graph    *DependencyGraph
    resolver *DependencyResolver
}

func (cs *ConstraintSolver) SolveConstraints(ctx context.Context) (*ResolvedDependencySet, error) {
    // Group dependencies by component
    componentGroups := cs.groupDependenciesByComponent()
    
    resolved := &ResolvedDependencySet{
        Components: make(map[string]*ResolvedComponent),
        Conflicts:  []Conflict{},
    }
    
    // Resolve each component group
    for componentKey, deps := range componentGroups {
        resolvedComponent, err := cs.resolveComponentGroup(ctx, componentKey, deps)
        if err != nil {
            if conflict, ok := err.(*ConstraintConflict); ok {
                resolved.Conflicts = append(resolved.Conflicts, *conflict)
                continue
            }
            return nil, err
        }
        
        resolved.Components[componentKey] = resolvedComponent
    }
    
    return resolved, nil
}

func (cs *ConstraintSolver) resolveComponentGroup(ctx context.Context, componentKey string, deps []*Dependency) (*ResolvedComponent, error) {
    // Find all available versions
    versions, err := cs.resolver.GetAvailableVersions(ctx, componentKey)
    if err != nil {
        return nil, err
    }
    
    // Find version that satisfies all constraints
    for _, version := range versions {
        satisfiesAll := true
        for _, dep := range deps {
            if !cs.satisfiesConstraint(version, dep.VersionConstraint) {
                satisfiesAll = false
                break
            }
        }
        
        if satisfiesAll {
            // Resolve to exact commit
            commit, err := cs.resolver.GetCommitForVersion(ctx, componentKey, version)
            if err != nil {
                continue // Try next version
            }
            
            return &ResolvedComponent{
                ComponentKey: componentKey,
                Version:      version,
                Commit:       commit,
                ResolvedAt:   time.Now(),
            }, nil
        }
    }
    
    // No version satisfies all constraints
    return nil, &ConstraintConflict{
        ComponentKey: componentKey,
        Constraints:  deps,
        Message:      "No version satisfies all constraints",
    }
}
```

### Phase 3: Commit Resolution
```go
func (dr *DependencyResolver) ResolveToCommits(ctx context.Context, resolved *ResolvedDependencySet) (*CommitResolvedSet, error) {
    commitResolved := &CommitResolvedSet{
        Components: make(map[string]*CommitResolvedComponent),
        LockFile:   &LockFile{},
    }
    
    for key, component := range resolved.Components {
        // Get exact commit for resolved version
        commit, err := dr.getExactCommit(ctx, component)
        if err != nil {
            return nil, fmt.Errorf("failed to resolve commit for %s: %w", key, err)
        }
        
        // Load component definition at exact commit
        definition, err := dr.loadComponentAtCommit(ctx, component.ComponentKey, commit)
        if err != nil {
            return nil, fmt.Errorf("failed to load component definition: %w", err)
        }
        
        commitResolvedComponent := &CommitResolvedComponent{
            ComponentKey: component.ComponentKey,
            Version:      component.Version,
            Commit:       commit,
            Definition:   definition,
            ResolvedAt:   time.Now(),
        }
        
        commitResolved.Components[key] = commitResolvedComponent
        
        // Add to lock file
        commitResolved.LockFile.Dependencies[key] = &LockEntry{
            Repository: component.Repository,
            Component:  component.Component,
            Version:    component.Version,
            Commit:     commit,
            ResolvedAt: time.Now(),
            Checksum:   dr.calculateChecksum(definition),
        }
    }
    
    return commitResolved, nil
}
```

## Lock Files

### Lock File Format
```yaml
# forge-lock.yaml
lockfile_version: "1.0"
generated_at: "2024-01-15T10:30:00Z"
forge_version: "0.2.0"
root_component: "sales-specialist"

# Resolved dependencies with exact commits
dependencies:
  "github.com/company/crm-tools/salesforce-lookup":
    repository: "github.com/company/crm-tools"
    component: "salesforce-lookup"
    version: "2.1.3"
    version_constraint: "^2.1.0"
    commit: "a1b2c3d4e5f6789012345678901234567890abcd"
    resolved_at: "2024-01-15T10:30:00Z"
    checksum: "sha256:abc123def456..."
    
  "github.com/tools/communication/email-composer":
    repository: "github.com/tools/communication"
    component: "email-composer"
    version: "1.3.2"
    version_constraint: "~1.3.0"
    commit: "b2c3d4e5f6789012345678901234567890abcdef"
    resolved_at: "2024-01-15T10:30:00Z"
    checksum: "sha256:def456ghi789..."
    
  "github.com/ai-prompts/sales-templates/sales-system-v3":
    repository: "github.com/ai-prompts/sales-templates"
    component: "sales-system-v3"
    version: "3.2.0"
    version_constraint: ">=3.0.0"
    commit: "c3d4e5f6789012345678901234567890abcdef12"
    resolved_at: "2024-01-15T10:30:00Z"
    checksum: "sha256:ghi789jkl012..."

# Transitive dependencies
transitive_dependencies:
  "github.com/tools/auth/oauth-handler":
    repository: "github.com/tools/auth"
    component: "oauth-handler"
    version: "1.2.1"
    commit: "d4e5f6789012345678901234567890abcdef1234"
    required_by: ["github.com/company/crm-tools/salesforce-lookup"]
    resolved_at: "2024-01-15T10:30:00Z"
    checksum: "sha256:jkl012mno345..."

# Resolution metadata
resolution_metadata:
  total_dependencies: 4
  resolution_time: "1.2s"
  conflicts_resolved: 0
  cache_hits: 2
  network_requests: 6
```

### Lock File Operations
```bash
# Generate lock file
forge lock generate
forge lock generate --output custom-lock.yaml

# Validate lock file
forge lock validate
forge lock validate --lock-file custom-lock.yaml

# Update specific dependency in lock file
forge lock update salesforce-lookup
forge lock update --all --check-breaking

# Install from lock file
forge install --from-lock
forge install --from-lock custom-lock.yaml --verify-checksums
```

## Conflict Resolution

### Conflict Types
```go
type ConflictType string

const (
    ConflictTypeVersion     ConflictType = "version"     // Version constraint conflicts
    ConflictTypeDependency  ConflictType = "dependency"  // Dependency requirement conflicts
    ConflictTypeCircular    ConflictType = "circular"    // Circular dependency
    ConflictTypeMissing     ConflictType = "missing"     // Missing dependency
    ConflictTypeIncompatible ConflictType = "incompatible" // Incompatible requirements
)

type Conflict struct {
    Type         ConflictType
    ComponentKey string
    Description  string
    Constraints  []*Dependency
    Suggestions  []ConflictSuggestion
}

type ConflictSuggestion struct {
    Action      string // "upgrade", "downgrade", "exclude", "override"
    Component   string
    FromVersion string
    ToVersion   string
    Reason      string
}
```

### Automatic Conflict Resolution
```go
func (cs *ConstraintSolver) ResolveConflicts(ctx context.Context, conflicts []Conflict) (*ConflictResolution, error) {
    resolution := &ConflictResolution{
        Resolutions: make(map[string]ConflictAction),
        Manual:      []Conflict{},
    }
    
    for _, conflict := range conflicts {
        action, err := cs.resolveConflictAutomatically(ctx, conflict)
        if err != nil {
            // Cannot resolve automatically, requires manual intervention
            resolution.Manual = append(resolution.Manual, conflict)
            continue
        }
        
        resolution.Resolutions[conflict.ComponentKey] = action
    }
    
    return resolution, nil
}

func (cs *ConstraintSolver) resolveConflictAutomatically(ctx context.Context, conflict Conflict) (ConflictAction, error) {
    switch conflict.Type {
    case ConflictTypeVersion:
        return cs.resolveVersionConflict(ctx, conflict)
    case ConflictTypeDependency:
        return cs.resolveDependencyConflict(ctx, conflict)
    case ConflictTypeCircular:
        return ConflictAction{}, fmt.Errorf("circular dependencies require manual resolution")
    default:
        return ConflictAction{}, fmt.Errorf("unknown conflict type: %s", conflict.Type)
    }
}
```

### Manual Conflict Resolution
```bash
# List conflicts
forge deps conflicts
# Conflicts found:
# 1. Version conflict: salesforce-lookup
#    - sales-specialist requires ^2.1.0
#    - crm-integration requires ~2.0.0
#    - No version satisfies both constraints
#    
#    Suggestions:
#    - Upgrade crm-integration to use ^2.1.0
#    - Downgrade sales-specialist to use ~2.0.0
#    - Override with specific version

# Resolve specific conflict
forge deps resolve salesforce-lookup --strategy upgrade-dependents
forge deps resolve salesforce-lookup --override-version 2.1.0
forge deps resolve salesforce-lookup --exclude-from crm-integration

# Interactive resolution
forge deps resolve --interactive
# Conflict 1/3: salesforce-lookup version conflict
# Choose resolution:
#   1. Upgrade crm-integration constraint to ^2.1.0
#   2. Downgrade sales-specialist constraint to ~2.0.0  
#   3. Override with specific version 2.1.0
#   4. Skip this conflict
# Choice: 1
```

## Dependency Graph Analysis

### Graph Visualization
```bash
# Show dependency tree
forge deps tree sales-specialist
# sales-specialist@1.5.0
# ├── salesforce-lookup@2.1.3
# │   └── oauth-handler@1.2.1
# ├── email-composer@1.3.2
# │   ├── smtp-client@2.0.1
# │   └── template-engine@1.1.0
# └── sales-system-v3@3.2.0

# Show dependency graph
forge deps graph sales-specialist --format dot
forge deps graph sales-specialist --format json
forge deps graph sales-specialist --output deps.svg

# Analyze dependency health
forge deps analyze sales-specialist
# Dependency Analysis for sales-specialist:
# ✓ No circular dependencies
# ✓ All dependencies resolved
# ⚠ 2 dependencies have newer versions available
# ⚠ 1 dependency has security advisory
# 
# Recommendations:
# - Update email-composer: 1.3.2 → 1.3.4 (security fix)
# - Update oauth-handler: 1.2.1 → 1.2.3 (bug fixes)
```

### Dependency Metrics
```go
type DependencyMetrics struct {
    TotalDependencies    int
    DirectDependencies   int
    TransitiveDependencies int
    MaxDepth            int
    CircularDependencies []string
    OutdatedDependencies []OutdatedDependency
    SecurityAdvisories   []SecurityAdvisory
    LicenseIssues       []LicenseIssue
}

type OutdatedDependency struct {
    Component      string
    CurrentVersion string
    LatestVersion  string
    VersionsBehind int
    SecurityRisk   bool
}
```

## Performance Optimization

### Parallel Resolution
```go
func (dr *DependencyResolver) ResolveParallel(ctx context.Context, components []string) (*ResolvedDependencySet, error) {
    semaphore := make(chan struct{}, MaxConcurrentResolutions)
    results := make(chan ResolutionResult, len(components))
    
    var wg sync.WaitGroup
    for _, component := range components {
        wg.Add(1)
        go func(comp string) {
            defer wg.Done()
            semaphore <- struct{}{} // Acquire
            defer func() { <-semaphore }() // Release
            
            resolved, err := dr.resolveComponent(ctx, comp)
            results <- ResolutionResult{
                Component: comp,
                Resolved:  resolved,
                Error:     err,
            }
        }(component)
    }
    
    wg.Wait()
    close(results)
    
    // Collect results
    resolvedSet := &ResolvedDependencySet{
        Components: make(map[string]*ResolvedComponent),
    }
    
    for result := range results {
        if result.Error != nil {
            return nil, result.Error
        }
        resolvedSet.Components[result.Component] = result.Resolved
    }
    
    return resolvedSet, nil
}
```

### Caching Strategy
```go
type ResolutionCache struct {
    versionCache    *cache.Cache // Cache available versions
    commitCache     *cache.Cache // Cache commit resolutions
    componentCache  *cache.Cache // Cache component definitions
    constraintCache *cache.Cache // Cache constraint satisfactions
}

func (rc *ResolutionCache) GetCachedResolution(key string) (*ResolvedComponent, bool) {
    if cached := rc.constraintCache.Get(key); cached != nil {
        return cached.(*ResolvedComponent), true
    }
    return nil, false
}

func (rc *ResolutionCache) CacheResolution(key string, resolved *ResolvedComponent) {
    // Cache for 1 hour by default, longer for stable versions
    ttl := 1 * time.Hour
    if resolved.IsStableVersion() {
        ttl = 24 * time.Hour
    }
    
    rc.constraintCache.Set(key, resolved, ttl)
}
```

### Incremental Resolution
```go
func (dr *DependencyResolver) IncrementalResolve(ctx context.Context, changes []ComponentChange) (*ResolvedDependencySet, error) {
    // Only re-resolve affected components
    affectedComponents := dr.findAffectedComponents(changes)
    
    // Load existing resolution
    existing, err := dr.loadExistingResolution(ctx)
    if err != nil {
        return nil, err
    }
    
    // Re-resolve only affected components
    for _, component := range affectedComponents {
        resolved, err := dr.resolveComponent(ctx, component)
        if err != nil {
            return nil, err
        }
        existing.Components[component] = resolved
    }
    
    return existing, nil
}
```