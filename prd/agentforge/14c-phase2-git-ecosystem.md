# Phase 2: Git-Native Ecosystem (v0.2.0)

## Overview
**Timeline**: Weeks 9-16  
**Goal**: Implement full Git-native component lifecycle with external repository support and robust dependency management

## Objectives

### Primary Goals
1. **External Repositories**: Seamless integration with GitHub and other Git providers
2. **Dependency Management**: Robust resolution with semantic versioning and commit-based linking
3. **Security Foundation**: Basic security scanning and permission management
4. **Developer Workflow**: Intuitive local development with change tracking

### Success Criteria
- [ ] External repositories work seamlessly with 100+ components
- [ ] Dependency resolution handles complex scenarios without conflicts
- [ ] Security prevents obvious threats and vulnerabilities
- [ ] Local development workflow feels natural to developers
- [ ] Performance remains acceptable with 1000+ components

## Detailed Milestones

### Week 9-10: Repository System

#### Milestone: External Repositories
**Deliverables:**
```bash
✅ Remote repository support
  - GitHub, GitLab, and generic Git repository support
  - Authentication and access token management
  - Repository discovery and metadata caching

✅ Repository cloning and caching
  - Intelligent caching with TTL and invalidation
  - Shallow clones for performance
  - Incremental updates and sync

✅ Component discovery across repos
  - Cross-repository component search
  - Metadata aggregation and indexing
  - Performance optimization for large catalogs

✅ Repository metadata and manifests
  - forge-manifest.yaml standard implementation
  - Repository health and quality metrics
  - Compatibility and version tracking

✅ Repository health monitoring
  - Availability and response time monitoring
  - Quality metrics calculation
  - Automated health reporting
```

**Technical Implementation:**
```go
type RepositoryManager struct {
    cache       *RepositoryCache
    auth        *AuthManager
    indexer     *ComponentIndexer
    monitor     *HealthMonitor
}

type Repository struct {
    URL         string
    Name        string
    Type        RepositoryType
    Manifest    *Manifest
    LastSync    time.Time
    HealthScore float64
}

func (rm *RepositoryManager) AddRepository(url string) error
func (rm *RepositoryManager) SyncRepository(name string) error
func (rm *RepositoryManager) DiscoverComponents(repo *Repository) ([]Component, error)
```

**CLI Commands:**
```bash
# Repository management
forge repo add <name> <url>
forge repo list [--type tools|prompts|agents]
forge repo sync <name>
forge repo health <name>
forge repo remove <name>

# Component discovery
forge search <query> [--repo <name>]
forge search --type tools --category crm
forge list --repo github.com/company/tools
```

### Week 11-12: Dependency Resolution

#### Milestone: Dependency Management
**Deliverables:**
```bash
✅ Semantic version constraints
  - Support for ^, ~, >=, exact version specifications
  - Version range validation and parsing
  - Compatibility checking across constraints

✅ Dependency graph resolution
  - Topological sorting for dependency order
  - Circular dependency detection
  - Conflict identification and reporting

✅ Commit-based exact resolution
  - Every dependency resolved to exact Git commit
  - Reproducible builds across environments
  - Lock file generation and validation

✅ Lock file generation
  - forge-lock.yaml with exact commit references
  - Checksum validation for integrity
  - Human-readable format with metadata

✅ Conflict detection and reporting
  - Version constraint conflicts
  - Dependency requirement conflicts
  - Suggested resolution strategies
```

**Dependency Resolution Algorithm:**
```go
type DependencyResolver struct {
    repoManager *RepositoryManager
    cache       *ResolutionCache
    validator   *ConstraintValidator
}

type Dependency struct {
    Repository  string
    Component   string
    Constraint  string  // "^2.1.0", "~1.3.0", "latest"
    Required    bool
}

type ResolvedDependency struct {
    Repository string
    Component  string
    Version    string
    Commit     string
    ResolvedAt time.Time
}

func (dr *DependencyResolver) Resolve(deps []Dependency) (*DependencySet, error)
func (dr *DependencyResolver) GenerateLockFile(resolved *DependencySet) (*LockFile, error)
```

**Lock File Format:**
```yaml
# forge-lock.yaml
lockfile_version: "1.0"
generated_at: "2024-01-15T10:30:00Z"
forge_version: "0.2.0"

dependencies:
  "github.com/company/crm-tools/salesforce-lookup":
    repository: "github.com/company/crm-tools"
    component: "salesforce-lookup"
    version: "2.1.3"
    constraint: "^2.1.0"
    commit: "a1b2c3d4e5f6789012345678901234567890abcd"
    checksum: "sha256:abc123..."
    resolved_at: "2024-01-15T10:30:00Z"
```

**CLI Commands:**
```bash
# Dependency management
forge deps resolve <component>
forge deps tree <component>
forge deps outdated
forge deps update <component>
forge lock generate
forge lock validate
```

### Week 13-14: Security Foundation

#### Milestone: Basic Security
**Deliverables:**
```bash
✅ Component scanning framework
  - Static code analysis for common vulnerabilities
  - Dependency vulnerability scanning
  - Secret detection in component code
  - Malware and suspicious pattern detection

✅ Permission system foundation
  - Component permission declarations
  - Runtime permission enforcement
  - Permission grant/revoke mechanisms
  - Audit logging for permission usage

✅ Basic sandboxing (containers)
  - Container-based component isolation
  - Resource limits (CPU, memory, disk)
  - Network access controls
  - Filesystem access restrictions

✅ Security policy framework
  - Policy definition and management
  - Component compliance checking
  - Violation reporting and remediation
  - Enterprise policy templates

✅ Audit logging
  - Security event logging
  - Component usage tracking
  - Permission access logs
  - Compliance reporting
```

**Security Architecture:**
```go
type SecurityManager struct {
    scanner     *SecurityScanner
    permissions *PermissionManager
    sandbox     *SandboxManager
    auditor     *AuditLogger
}

type SecurityPolicy struct {
    Name                string
    ComponentRequirements ComponentRequirements
    RuntimeRestrictions   RuntimeRestrictions
    AuditRequirements    AuditRequirements
}

type Permission struct {
    Type     PermissionType  // network, filesystem, environment
    Resource string          // specific resource identifier
    Actions  []string        // read, write, execute
}
```

**Security Commands:**
```bash
# Security operations
forge security scan <component>
forge security policy list
forge security policy apply <policy> --to <component>
forge security permissions <component>
forge security audit --component <component> --since 7d
```

### Week 15-16: Local Development

#### Milestone: Development Workflow
**Deliverables:**
```bash
✅ Local component editing
  - Interactive component editor
  - File-based editing with validation
  - Real-time syntax checking
  - Template and example generation

✅ Change tracking and diffing
  - Git-style diff for component changes
  - Change history and attribution
  - Rollback and revert capabilities
  - Branch-based development workflow

✅ Local testing framework
  - Component unit testing
  - Integration testing with dependencies
  - Performance benchmarking
  - Validation and linting

✅ Component validation
  - Schema validation against standards
  - Dependency requirement checking
  - Security policy compliance
  - Best practice recommendations

✅ Development workspace management
  - Isolated development environments
  - Workspace sharing and collaboration
  - Component versioning and tagging
  - Export and publication workflows
```

**Development Workflow:**
```go
type DevelopmentManager struct {
    editor      *ComponentEditor
    tracker     *ChangeTracker
    tester      *TestRunner
    validator   *ComponentValidator
    workspace   *WorkspaceManager
}

type Workspace struct {
    Name        string
    Path        string
    Components  []LocalComponent
    Config      WorkspaceConfig
}

type LocalComponent struct {
    Name            string
    OriginalVersion string
    LocalChanges    []Change
    TestResults     *TestResults
}
```

**Development Commands:**
```bash
# Local development
forge dev edit <component>
forge dev test <component>
forge dev validate <component>
forge dev diff <component>
forge dev commit <component> --message "..."

# Workspace management
forge workspace create <name>
forge workspace switch <name>
forge workspace list
forge workspace export <name>
```

## Technical Architecture

### Enhanced System Architecture
```
AgentForge v0.2.0 Architecture:
├── CLI Layer
│   ├── Command Router
│   ├── Repository Commands
│   ├── Dependency Commands
│   └── Development Commands
├── Core Services
│   ├── Repository Manager
│   ├── Dependency Resolver
│   ├── Security Manager
│   ├── Development Manager
│   └── Component Manager
├── Git Integration
│   ├── Repository Cache
│   ├── Clone Manager
│   ├── Sync Engine
│   └── Commit Tracker
├── Security Layer
│   ├── Scanner Framework
│   ├── Permission System
│   ├── Sandbox Manager
│   └── Audit Logger
└── Storage Layer
    ├── Local Cache
    ├── Component Index
    ├── Lock Files
    └── Workspace Data
```

### Database Schema Extensions
```sql
-- Repository management
CREATE TABLE af_repositories (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE,
    url VARCHAR(500),
    type VARCHAR(50),
    manifest JSONB,
    last_sync TIMESTAMPTZ,
    health_score FLOAT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Component dependencies
CREATE TABLE af_component_dependencies (
    id UUID PRIMARY KEY,
    component_id UUID REFERENCES af_components(id),
    dependency_repository VARCHAR(500),
    dependency_component VARCHAR(255),
    version_constraint VARCHAR(100),
    resolved_version VARCHAR(50),
    resolved_commit VARCHAR(40),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Security policies
CREATE TABLE af_security_policies (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE,
    definition JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Audit logs
CREATE TABLE af_audit_logs (
    id UUID PRIMARY KEY,
    event_type VARCHAR(100),
    component_id UUID,
    user_id VARCHAR(255),
    details JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Configuration Enhancements
```yaml
# ~/.forge/config.yaml (v0.2.0)
repositories:
  cache_dir: "~/.forge/cache"
  max_cache_size: "5GB"
  sync_interval: "1h"
  default_timeout: "30s"

dependencies:
  resolution_strategy: "latest_compatible"
  allow_prerelease: false
  max_resolution_time: "60s"

security:
  enable_scanning: true
  default_policy: "standard"
  sandbox_runtime: "docker"
  audit_level: "info"

development:
  default_editor: "$EDITOR"
  auto_validate: true
  auto_test: false
  workspace_dir: "~/.forge/workspaces"
```

## Performance Targets

### Repository Operations
- **Repository Clone**: <30 seconds for typical repository
- **Component Discovery**: <5 seconds for 1000+ components
- **Repository Sync**: <10 seconds for incremental updates
- **Search Operations**: <2 seconds across all repositories

### Dependency Resolution
- **Simple Resolution**: <5 seconds for <10 dependencies
- **Complex Resolution**: <30 seconds for 50+ dependencies
- **Lock File Generation**: <10 seconds for any complexity
- **Conflict Detection**: <15 seconds for complex scenarios

### Security Operations
- **Component Scan**: <60 seconds for typical component
- **Permission Check**: <100ms for any permission
- **Policy Validation**: <5 seconds for any component
- **Audit Log Query**: <2 seconds for typical queries

### Development Workflow
- **Component Edit**: <1 second to open editor
- **Validation**: <5 seconds for any component
- **Testing**: <30 seconds for unit tests
- **Diff Generation**: <2 seconds for any changes

## Quality Assurance

### Testing Strategy
- **Unit Tests**: >95% coverage for new components
- **Integration Tests**: All repository and dependency workflows
- **Performance Tests**: All performance targets validated
- **Security Tests**: Vulnerability scanning and penetration testing

### Security Requirements
- **Vulnerability Scanning**: Zero critical vulnerabilities
- **Permission Enforcement**: 100% permission checks enforced
- **Sandbox Isolation**: Complete process and network isolation
- **Audit Completeness**: All security events logged

### Compatibility Requirements
- **Git Providers**: GitHub, GitLab, Bitbucket, generic Git
- **Operating Systems**: Linux, macOS, Windows
- **Container Runtimes**: Docker, Podman, containerd
- **Database Systems**: PostgreSQL 12+, SQLite (development)

## Risk Mitigation

### Technical Risks
| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Git performance with large repos | High | Medium | Shallow clones, LFS support, caching |
| Dependency resolution complexity | High | Medium | Comprehensive testing, fallback strategies |
| Security sandbox bypass | High | Low | Multiple isolation layers, regular audits |
| Repository availability issues | Medium | Medium | Caching, fallback repositories, mirrors |

### Integration Risks
| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| GitHub API rate limits | Medium | Medium | Caching, authentication, request optimization |
| Network connectivity issues | Medium | High | Offline mode, cached operations, retries |
| Authentication complexity | Medium | Low | Multiple auth methods, clear documentation |
| Repository format variations | Low | Medium | Flexible parsing, validation, standards |

## Success Metrics

### Functional Metrics
- **Repository Integration**: 100% success rate with major Git providers
- **Dependency Resolution**: >99% success rate for valid constraints
- **Security Scanning**: <1% false positive rate
- **Development Workflow**: <5 minutes from idea to tested component

### Performance Metrics
- **Repository Operations**: All targets met consistently
- **Memory Usage**: <500MB for typical operations
- **Disk Usage**: <2GB for full installation with cache
- **Network Efficiency**: <10MB for typical repository sync

### User Experience Metrics
- **Command Success Rate**: >98% of commands succeed
- **Error Message Quality**: <10% of errors require support
- **Documentation Coverage**: 100% of features documented
- **User Satisfaction**: >4.2/5 rating for Phase 2 features

## Deliverable Checklist

### Week 9-10 Deliverables
- [ ] Remote repository support functional
- [ ] Repository caching system operational
- [ ] Component discovery across repositories
- [ ] Repository health monitoring active
- [ ] Performance targets met

### Week 11-12 Deliverables
- [ ] Semantic versioning fully supported
- [ ] Dependency resolution handles complex scenarios
- [ ] Lock file generation and validation working
- [ ] Conflict detection and reporting functional
- [ ] Performance targets achieved

### Week 13-14 Deliverables
- [ ] Security scanning framework operational
- [ ] Permission system enforcing access controls
- [ ] Container-based sandboxing working
- [ ] Security policies configurable and enforced
- [ ] Audit logging capturing all events

### Week 15-16 Deliverables
- [ ] Local component editing intuitive and functional
- [ ] Change tracking and diffing operational
- [ ] Testing framework supporting all component types
- [ ] Validation catching common issues
- [ ] Workspace management supporting team workflows

### Final Phase 2 Validation
- [ ] All success criteria met
- [ ] Performance benchmarks achieved
- [ ] Security requirements satisfied
- [ ] Integration tests passing
- [ ] User acceptance testing completed

## Next Phase Preparation

### Phase 3 Prerequisites
- [ ] Repository system stable and performant
- [ ] Dependency resolution robust and reliable
- [ ] Security foundation solid and extensible
- [ ] Development workflow proven and adopted
- [ ] Performance targets consistently met

### Technical Debt to Address
- [ ] Basic security model (enhance in Phase 3)
- [ ] Limited collaboration features (expand in Phase 3)
- [ ] Manual conflict resolution (automate in Phase 3)
- [ ] Basic marketplace features (enhance in Phase 3)
- [ ] Limited analytics (expand in Phase 3)