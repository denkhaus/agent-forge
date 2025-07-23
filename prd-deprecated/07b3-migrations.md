# Migrations & Schema Evolution

## Overview

Database migrations for the MCP-Planner system manage schema evolution, data transformations, and version control. This guide is broken down into focused migration areas for systematic database management.

## Migration Structure

The migration implementation is organized into the following documents:

- **[Migration Framework & Setup](./07b3a-migration-framework.md)** - Migration tooling, framework setup, and execution
- **[Initial Schema Migrations](./07b3b-initial-schema.md)** - Core table creation and initial schema setup
- **[Index & Constraint Migrations](./07b3c-indexes-constraints.md)** - Performance indexes and data integrity constraints
- **[Function & Trigger Migrations](./07b3d-functions-triggers.md)** - Database functions, triggers, and automation
- **[Data Migration Strategies](./07b3e-data-migrations.md)** - Data transformation and migration patterns
- **[Schema Versioning](./07b3f-schema-versioning.md)** - Version control and rollback strategies
- **[Migration Testing](./07b3g-migration-testing.md)** - Testing migration scripts and validation

## Migration Philosophy

### Version-Controlled Schema
- All schema changes tracked in version control
- Sequential migration numbering
- Rollback capability for all changes
- Environment-specific migration handling

### Zero-Downtime Migrations
- Backward-compatible changes when possible
- Multi-phase migrations for breaking changes
- Online schema modifications
- Data migration strategies

### Data Integrity
- Constraint validation during migrations
- Data transformation verification
- Rollback data preservation
- Migration transaction safety

## Migration Tools

### Primary Tools
- **golang-migrate**: Primary migration framework
- **Prisma Migrate**: Schema generation and management
- **Custom Scripts**: Complex data transformations
- **SQL Scripts**: Direct database modifications

### Migration Commands
```bash
# Create new migration
make migrate-create NAME=add_complexity_analysis

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Check migration status
make migrate-status

# Force migration version
make migrate-force VERSION=001

# Validate migrations
make migrate-validate
```

## Migration Workflow

### Development Workflow
```mermaid
flowchart TD
    A[Schema Change Required] --> B[Create Migration]
    B --> C[Write Up Migration]
    C --> D[Write Down Migration]
    D --> E[Test Locally]
    E --> F{Tests Pass?}
    F -->|No| G[Fix Migration]
    G --> E
    F -->|Yes| H[Commit Migration]
    H --> I[Code Review]
    I --> J[Deploy to Staging]
    J --> K[Validate in Staging]
    K --> L[Deploy to Production]
```

### Migration Types

#### Schema Migrations
- Table creation/modification
- Column additions/changes
- Index management
- Constraint updates

#### Data Migrations
- Data transformations
- Reference updates
- Cleanup operations
- Seed data insertion

#### Function Migrations
- Stored procedures
- Triggers
- Custom functions
- Database automation

## Migration Standards

### Naming Convention
```
{version}_{type}_{description}.{direction}.sql

Examples:
001_create_projects_table.up.sql
001_create_projects_table.down.sql
002_add_complexity_fields.up.sql
002_add_complexity_fields.down.sql
003_create_navigation_indexes.up.sql
003_create_navigation_indexes.down.sql
```

### Migration Structure
```sql
-- Migration Header
-- Version: 001
-- Description: Create projects table
-- Author: Developer Name
-- Date: 2024-01-01
-- Dependencies: None

-- Migration Body
BEGIN;

-- Schema changes
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    -- ... other fields
);

-- Validation
DO $$
BEGIN
    -- Verify table exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'projects') THEN
        RAISE EXCEPTION 'Migration failed: projects table not created';
    END IF;
END $$;

COMMIT;
```

### Rollback Strategy
```sql
-- Rollback Header
-- Version: 001
-- Description: Drop projects table
-- Rollback for: Create projects table

BEGIN;

-- Backup critical data if needed
CREATE TEMP TABLE projects_backup AS SELECT * FROM projects;

-- Rollback changes
DROP TABLE IF EXISTS projects CASCADE;

-- Validation
DO $$
BEGIN
    -- Verify table is dropped
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'projects') THEN
        RAISE EXCEPTION 'Rollback failed: projects table still exists';
    END IF;
END $$;

COMMIT;
```

## Migration Safety

### Pre-Migration Checks
- Database backup verification
- Disk space availability
- Connection pool capacity
- Lock timeout configuration

### Transaction Management
- Atomic migrations when possible
- Long-running operation handling
- Lock management
- Timeout configuration

### Validation
- Schema validation after migration
- Data integrity checks
- Performance impact assessment
- Rollback capability verification

## Environment Management

### Development Environment
```bash
# Local development database
DATABASE_URL=postgresql://dev:password@localhost:5432/mcp_task_dev

# Run migrations
migrate -path migrations -database $DATABASE_URL up
```

### Staging Environment
```bash
# Staging database with production-like data
DATABASE_URL=postgresql://staging:password@staging-db:5432/mcp_task_staging

# Validate migrations before production
migrate -path migrations -database $DATABASE_URL up
```

### Production Environment
```bash
# Production database with backup and monitoring
DATABASE_URL=postgresql://prod:password@prod-db:5432/mcp_task_prod

# Production migration with safety checks
./scripts/production-migrate.sh
```

## Migration Monitoring

### Migration Metrics
- Migration execution time
- Data volume processed
- Lock duration
- Error rates

### Alerting
- Migration failures
- Long-running migrations
- Lock timeouts
- Data integrity violations

### Logging
- Migration start/completion
- Error details
- Performance metrics
- Rollback events

## Best Practices

### Schema Design
- Use appropriate data types
- Plan for future extensibility
- Consider performance implications
- Maintain referential integrity

### Migration Writing
- Keep migrations atomic
- Write comprehensive rollbacks
- Include validation checks
- Document complex changes

### Testing
- Test on production-like data
- Validate rollback procedures
- Performance test large migrations
- Test concurrent access scenarios

### Deployment
- Backup before migration
- Monitor during execution
- Validate after completion
- Have rollback plan ready

---

*Continue with detailed migration implementation files:*
- *[Migration Framework & Setup](./07b3a-migration-framework.md)*
- *[Initial Schema Migrations](./07b3b-initial-schema.md)*
- *[Index & Constraint Migrations](./07b3c-indexes-constraints.md)*
- *[Function & Trigger Migrations](./07b3d-functions-triggers.md)*
- *[Data Migration Strategies](./07b3e-data-migrations.md)*
- *[Schema Versioning](./07b3f-schema-versioning.md)*
- *[Migration Testing](./07b3g-migration-testing.md)*
