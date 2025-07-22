# Database Schema

## Overview

AgentForge uses PostgreSQL for local storage of component metadata, sync state, and composition configurations. The schema supports the full Git-native workflow with bidirectional sync capabilities.

## Core Schema

### Repository Management

```prisma
model Repository {
  id          String   @id @default(uuid()) @db.Uuid
  name        String   @unique @db.VarChar(255)
  url         String   @db.VarChar(500)
  type        RepoType
  isActive    Boolean  @default(true)
  
  // Git information
  defaultBranch String  @default("main") @db.VarChar(100)
  lastSync     DateTime?
  syncStatus   RepoSyncStatus @default(NEVER_SYNCED)
  
  // Cached manifest data
  manifest     Json?
  manifestHash String? @db.VarChar(64)
  
  // Access control
  hasWriteAccess Boolean @default(false)
  accessToken    String? @db.Text
  
  // Relationships
  components   Component[]
  forks        Fork[]
  
  // Metadata
  createdAt    DateTime @default(now()) @db.Timestamptz
  updatedAt    DateTime @updatedAt @db.Timestamptz
  
  @@index([type])
  @@index([isActive])
  @@index([lastSync])
  @@map("repositories")
}

model Fork {
  id           String @id @default(uuid()) @db.Uuid
  originalRepo String @db.VarChar(500)
  forkRepo     String @db.VarChar(500)
  createdAt    DateTime @default(now()) @db.Timestamptz
  
  // Relationships
  repository   Repository @relation(fields: [originalRepo], references: [url])
  
  @@unique([originalRepo, forkRepo])
  @@map("forks")
}

enum RepoType {
  TOOLS
  PROMPTS
  AGENTS
  COMPOSITIONS
  MIXED
}

enum RepoSyncStatus {
  NEVER_SYNCED
  SYNCED
  SYNC_FAILED
  SYNC_IN_PROGRESS
  OUTDATED
}
```

### Component Management

```prisma
model Component {
  id           String        @id @default(uuid()) @db.Uuid
  name         String        @db.VarChar(255)
  type         ComponentType
  version      String        @db.VarChar(50)
  
  // Source information
  repositoryId String        @db.Uuid
  repositoryUrl String       @db.VarChar(500)
  path         String        @db.VarChar(500)
  commitHash   String        @db.VarChar(40)
  
  // Component definition
  definition   Json
  schema       Json?
  metadata     Json?
  
  // Installation state
  isInstalled  Boolean       @default(false)
  installedAt  DateTime?     @db.Timestamptz
  installPath  String?       @db.VarChar(500)
  
  // Local modifications
  hasLocalChanges Boolean    @default(false)
  localDefinition Json?
  originalDefinition Json?
  
  // Sync state
  syncStatus   ComponentSyncStatus @default(REMOTE_ONLY)
  lastSyncAt   DateTime?     @db.Timestamptz
  lastSyncCommit String?     @db.VarChar(40)
  upstreamCommit String?     @db.VarChar(40)
  
  // Fork information
  forkRepo     String?       @db.VarChar(500)
  forkBranch   String?       @db.VarChar(100)
  
  // Relationships
  repository   Repository    @relation(fields: [repositoryId], references: [id])
  dependencies ComponentDependency[] @relation("ComponentDeps")
  dependents   ComponentDependency[] @relation("DependentComponents")
  changes      ComponentChange[]
  syncOps      SyncOperation[]
  
  // Metadata
  createdAt    DateTime      @default(now()) @db.Timestamptz
  updatedAt    DateTime      @updatedAt @db.Timestamptz
  
  @@unique([repositoryUrl, path, commitHash])
  @@index([type])
  @@index([isInstalled])
  @@index([hasLocalChanges])
  @@index([syncStatus])
  @@map("components")
}

enum ComponentType {
  TOOL
  PROMPT
  AGENT
}

enum ComponentSyncStatus {
  REMOTE_ONLY     // Only exists in remote repository
  LOCAL_ONLY      // Only exists locally (new component)
  SYNCED          // In sync with remote
  MODIFIED        // Has local changes not synced
  FORKED          // User created fork
  CONFLICT        // Merge conflict with upstream
  OUTDATED        // Remote has newer version
  SYNC_FAILED     // Last sync attempt failed
}
```

### Dependency Management

```prisma
model ComponentDependency {
  id            String @id @default(uuid()) @db.Uuid
  
  // Dependency relationship
  componentId   String @db.Uuid
  dependencyId  String @db.Uuid
  
  // Version constraints
  versionConstraint String @db.VarChar(100)  // "^2.1.0", "~1.3.0", "latest"
  resolvedVersion   String @db.VarChar(50)   // Actually resolved version
  resolvedCommit    String @db.VarChar(40)   // Exact commit used
  
  // Dependency type
  dependencyType DependencyType
  isOptional     Boolean @default(false)
  
  // Resolution metadata
  resolvedAt     DateTime @default(now()) @db.Timestamptz
  resolutionPath String[] // How this dependency was resolved
  
  // Relationships
  component      Component @relation("ComponentDeps", fields: [componentId], references: [id])
  dependency     Component @relation("DependentComponents", fields: [dependencyId], references: [id])
  
  @@unique([componentId, dependencyId])
  @@index([dependencyType])
  @@map("component_dependencies")
}

enum DependencyType {
  DIRECT      // Directly specified dependency
  TRANSITIVE  // Dependency of a dependency
  PEER        // Peer dependency (must be provided by user)
  DEV         // Development-only dependency
}
```

### Change Tracking

```prisma
model ComponentChange {
  id          String     @id @default(uuid()) @db.Uuid
  componentId String     @db.Uuid
  
  // Change information
  changeType  ChangeType
  fieldPath   String     @db.VarChar(500)  // JSON path of changed field
  oldValue    Json?
  newValue    Json?
  
  // Change metadata
  reason      String?    @db.Text
  author      String?    @db.VarChar(255)
  timestamp   DateTime   @default(now()) @db.Timestamptz
  
  // Git information
  commitHash  String?    @db.VarChar(40)
  branchName  String?    @db.VarChar(100)
  
  // Relationships
  component   Component  @relation(fields: [componentId], references: [id])
  
  @@index([componentId])
  @@index([changeType])
  @@index([timestamp])
  @@map("component_changes")
}

enum ChangeType {
  CREATE
  UPDATE
  DELETE
  MOVE
  RENAME
}
```

### Sync Operations

```prisma
model SyncOperation {
  id           String              @id @default(uuid()) @db.Uuid
  componentId  String              @db.Uuid
  
  // Operation details
  operation    SyncOperationType
  status       SyncOperationStatus
  direction    SyncDirection
  
  // Git information
  sourceCommit String              @db.VarChar(40)
  targetCommit String?             @db.VarChar(40)
  branchName   String?             @db.VarChar(100)
  
  // Pull request information
  pullRequestUrl String?           @db.VarChar(500)
  pullRequestId  String?           @db.VarChar(50)
  
  // Operation metadata
  message      String?             @db.Text
  errorMessage String?             @db.Text
  conflictFiles String[]          // Files with conflicts
  
  // Timing
  createdAt    DateTime            @default(now()) @db.Timestamptz
  startedAt    DateTime?           @db.Timestamptz
  completedAt  DateTime?           @db.Timestamptz
  
  // Relationships
  component    Component           @relation(fields: [componentId], references: [id])
  
  @@index([componentId])
  @@index([status])
  @@index([operation])
  @@map("sync_operations")
}

enum SyncOperationType {
  PULL_UPSTREAM
  PUSH_TO_ORIGIN
  PUSH_TO_FORK
  CREATE_FORK
  CREATE_PR
  MERGE_UPSTREAM
  RESOLVE_CONFLICT
}

enum SyncOperationStatus {
  PENDING
  IN_PROGRESS
  COMPLETED
  FAILED
  CANCELLED
  REQUIRES_MANUAL_RESOLUTION
}

enum SyncDirection {
  UPSTREAM_TO_LOCAL
  LOCAL_TO_UPSTREAM
  LOCAL_TO_FORK
  FORK_TO_UPSTREAM
}
```

### Composition Management

```prisma
model Composition {
  id          String   @id @default(uuid()) @db.Uuid
  name        String   @unique @db.VarChar(255)
  description String?  @db.Text
  version     String   @db.VarChar(50)
  
  // Composition definition
  definition  Json     // Complete composition specification
  lockFile    Json?    // Resolved dependencies with exact commits
  
  // State
  isActive    Boolean  @default(true)
  isDeployed  Boolean  @default(false)
  
  // Local overrides
  localOverrides Json?
  hasLocalChanges Boolean @default(false)
  
  // Git information
  sourceRepo  String?  @db.VarChar(500)
  sourceCommit String? @db.VarChar(40)
  
  // Deployment information
  deployedAt  DateTime? @db.Timestamptz
  environment String?   @db.VarChar(100)
  
  // Relationships
  instances   CompositionInstance[]
  
  // Metadata
  createdAt   DateTime @default(now()) @db.Timestamptz
  updatedAt   DateTime @updatedAt @db.Timestamptz
  
  @@index([isActive])
  @@index([isDeployed])
  @@map("compositions")
}

model CompositionInstance {
  id            String      @id @default(uuid()) @db.Uuid
  compositionId String      @db.Uuid
  name          String      @db.VarChar(255)
  
  // Instance configuration
  config        Json
  overrides     Json?
  
  // Runtime state
  status        InstanceStatus @default(STOPPED)
  lastStarted   DateTime?      @db.Timestamptz
  lastStopped   DateTime?      @db.Timestamptz
  
  // Performance metrics
  requestCount  BigInt         @default(0)
  errorCount    BigInt         @default(0)
  avgResponseTime Float?       @db.DoublePrecision
  
  // Relationships
  composition   Composition    @relation(fields: [compositionId], references: [id])
  
  @@unique([compositionId, name])
  @@index([status])
  @@map("composition_instances")
}

enum InstanceStatus {
  STOPPED
  STARTING
  RUNNING
  STOPPING
  ERROR
  UPDATING
}
```

### Configuration and Settings

```prisma
model UserSettings {
  id       String @id @default(uuid()) @db.Uuid
  userId   String @unique @db.VarChar(255)
  
  // Git configuration
  gitUsername String? @db.VarChar(255)
  gitEmail    String? @db.VarChar(255)
  
  // Default preferences
  defaultLLMProvider String? @db.VarChar(100)
  defaultLLMModel    String? @db.VarChar(100)
  
  // CLI preferences
  outputFormat       String  @default("table") @db.VarChar(50)
  verboseLogging     Boolean @default(false)
  autoSync           Boolean @default(true)
  
  // Security settings
  allowUnsignedComponents Boolean @default(false)
  trustedRepositories     String[] // List of trusted repo URLs
  
  // Custom settings
  customSettings Json?
  
  @@map("user_settings")
}

model SystemConfig {
  key         String   @id @db.VarChar(255)
  value       Json
  description String?  @db.Text
  updatedAt   DateTime @updatedAt @db.Timestamptz
  
  @@map("system_config")
}
```

## Indexes and Performance

### Primary Indexes
```sql
-- Repository lookups
CREATE INDEX idx_repositories_url ON repositories(url);
CREATE INDEX idx_repositories_type_active ON repositories(type, is_active);

-- Component lookups
CREATE INDEX idx_components_repo_path ON components(repository_url, path);
CREATE INDEX idx_components_name_type ON components(name, type);
CREATE INDEX idx_components_sync_status ON components(sync_status);

-- Dependency resolution
CREATE INDEX idx_dependencies_component ON component_dependencies(component_id);
CREATE INDEX idx_dependencies_type ON component_dependencies(dependency_type);

-- Change tracking
CREATE INDEX idx_changes_component_time ON component_changes(component_id, timestamp);
CREATE INDEX idx_changes_type_time ON component_changes(change_type, timestamp);

-- Sync operations
CREATE INDEX idx_sync_ops_component_status ON sync_operations(component_id, status);
CREATE INDEX idx_sync_ops_created ON sync_operations(created_at);
```

### Query Optimization
```sql
-- Composite indexes for common queries
CREATE INDEX idx_components_installed_type ON components(is_installed, type) WHERE is_installed = true;
CREATE INDEX idx_components_local_changes ON components(has_local_changes, sync_status) WHERE has_local_changes = true;
CREATE INDEX idx_sync_ops_pending ON sync_operations(status, created_at) WHERE status IN ('PENDING', 'IN_PROGRESS');
```

## Data Integrity

### Constraints
```sql
-- Ensure valid Git commit hashes
ALTER TABLE components ADD CONSTRAINT check_commit_hash 
  CHECK (commit_hash ~ '^[a-f0-9]{40}$');

-- Ensure valid version strings
ALTER TABLE components ADD CONSTRAINT check_version_format
  CHECK (version ~ '^[0-9]+\.[0-9]+\.[0-9]+');

-- Ensure valid repository URLs
ALTER TABLE repositories ADD CONSTRAINT check_repo_url
  CHECK (url ~ '^https://github\.com/[^/]+/[^/]+$');
```

### Triggers
```sql
-- Update component updated_at when local changes are made
CREATE OR REPLACE FUNCTION update_component_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_component_update
  BEFORE UPDATE ON components
  FOR EACH ROW
  EXECUTE FUNCTION update_component_timestamp();

-- Log component changes
CREATE OR REPLACE FUNCTION log_component_change()
RETURNS TRIGGER AS $$
BEGIN
  IF TG_OP = 'UPDATE' AND OLD.local_definition IS DISTINCT FROM NEW.local_definition THEN
    INSERT INTO component_changes (component_id, change_type, field_path, old_value, new_value)
    VALUES (NEW.id, 'UPDATE', 'local_definition', OLD.local_definition, NEW.local_definition);
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_log_component_change
  AFTER UPDATE ON components
  FOR EACH ROW
  EXECUTE FUNCTION log_component_change();
```