# Initial Schema Migrations

## Overview

This document contains the initial database schema migrations for the MCP-Planner system, creating all core tables, enums, and basic relationships. These migrations establish the foundation for the hierarchical task management system.

## Migration 001: Create Projects Table

### Up Migration

```sql
-- migrations/001_create_projects_table.up.sql
-- Migration: Create projects table
-- Version: 001
-- Description: Create the projects table with basic fields and constraints
-- Author: MCP-Planner Team
-- Date: 2024-01-01

BEGIN;

-- Create projects table
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    progress DOUBLE PRECISION NOT NULL DEFAULT 0,
    complexity_threshold DOUBLE PRECISION NOT NULL DEFAULT 0.7,
    max_iterations INTEGER NOT NULL DEFAULT 3,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add constraints
ALTER TABLE projects ADD CONSTRAINT check_projects_progress
    CHECK (progress >= 0 AND progress <= 1);

ALTER TABLE projects ADD CONSTRAINT check_projects_complexity_threshold
    CHECK (complexity_threshold >= 0 AND complexity_threshold <= 1);

ALTER TABLE projects ADD CONSTRAINT check_projects_max_iterations
    CHECK (max_iterations > 0);

ALTER TABLE projects ADD CONSTRAINT check_projects_name_not_empty
    CHECK (LENGTH(TRIM(name)) > 0);

ALTER TABLE projects ADD CONSTRAINT check_projects_description_not_empty
    CHECK (LENGTH(TRIM(description)) > 0);

-- Add comments
COMMENT ON TABLE projects IS 'Root entity containing tasks and managing project-level settings';
COMMENT ON COLUMN projects.id IS 'Unique project identifier';
COMMENT ON COLUMN projects.name IS 'Human-readable project name';
COMMENT ON COLUMN projects.description IS 'Detailed project description and requirements';
COMMENT ON COLUMN projects.progress IS 'Overall project completion (0.0 to 1.0)';
COMMENT ON COLUMN projects.complexity_threshold IS 'Threshold for step complexity promotion (0.0 to 1.0)';
COMMENT ON COLUMN projects.max_iterations IS 'Maximum AI collaboration iterations per step';

-- Create basic indexes
CREATE INDEX idx_projects_created_at ON projects(created_at);
CREATE INDEX idx_projects_progress ON projects(progress);
CREATE INDEX idx_projects_complexity_threshold ON projects(complexity_threshold);

-- Validation
DO $$
BEGIN
    -- Verify table exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'projects') THEN
        RAISE EXCEPTION 'Migration failed: projects table not created';
    END IF;

    -- Verify constraints exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.check_constraints WHERE constraint_name = 'check_projects_progress') THEN
        RAISE EXCEPTION 'Migration failed: progress constraint not created';
    END IF;

    RAISE NOTICE 'Projects table created successfully';
END $$;

COMMIT;
```

### Down Migration

```sql
-- migrations/001_create_projects_table.down.sql
-- Rollback: Create projects table
-- Version: 001
-- Description: Drop projects table and related objects

BEGIN;

-- Drop indexes
DROP INDEX IF EXISTS idx_projects_complexity_threshold;
DROP INDEX IF EXISTS idx_projects_progress;
DROP INDEX IF EXISTS idx_projects_created_at;

-- Drop table (CASCADE to handle any future dependencies)
DROP TABLE IF EXISTS projects CASCADE;

-- Validation
DO $$
BEGIN
    -- Verify table is dropped
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'projects') THEN
        RAISE EXCEPTION 'Rollback failed: projects table still exists';
    END IF;

    RAISE NOTICE 'Projects table dropped successfully';
END $$;

COMMIT;
```

## Migration 002: Create Enums

### Up Migration

```sql
-- migrations/002_create_enums.up.sql
-- Migration: Create enum types
-- Version: 002
-- Description: Create all enum types used throughout the system

BEGIN;

-- Step status enum
CREATE TYPE step_status AS ENUM (
    'PENDING',
    'SERVER_DRAFT',
    'CLIENT_REVIEW',
    'AGREED',
    'DISPUTED',
    'USER_RESOLUTION'
);

-- Complexity level enum
CREATE TYPE complexity_level AS ENUM (
    'LOW',
    'MEDIUM',
    'HIGH'
);

-- Dispute status enum
CREATE TYPE dispute_status AS ENUM (
    'PENDING',
    'RESOLVED'
);

-- Resolution type enum
CREATE TYPE resolution_type AS ENUM (
    'SERVER',
    'CLIENT',
    'CUSTOM',
    'HYBRID'
);

-- Entity type enum (for audit logs)
CREATE TYPE entity_type AS ENUM (
    'PROJECT',
    'TASK',
    'STEP',
    'DISPUTE'
);

-- Audit action enum
CREATE TYPE audit_action AS ENUM (
    'CREATED',
    'UPDATED',
    'DELETED',
    'PROMOTED',
    'DISPUTED',
    'RESOLVED',
    'COMPLETED'
);

-- AI agent enum
CREATE TYPE ai_agent AS ENUM (
    'SERVER_AI',
    'CLIENT_AI'
);

-- Add comments
COMMENT ON TYPE step_status IS 'Status of step in dual-AI collaboration workflow';
COMMENT ON TYPE complexity_level IS 'Complexity assessment levels for steps';
COMMENT ON TYPE dispute_status IS 'Status of content disputes';
COMMENT ON TYPE resolution_type IS 'Type of dispute resolution chosen';
COMMENT ON TYPE entity_type IS 'Types of entities for audit logging';
COMMENT ON TYPE audit_action IS 'Actions that can be audited';
COMMENT ON TYPE ai_agent IS 'AI agents in the system';

-- Validation
DO $$
BEGIN
    -- Verify enums exist
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'step_status') THEN
        RAISE EXCEPTION 'Migration failed: step_status enum not created';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'complexity_level') THEN
        RAISE EXCEPTION 'Migration failed: complexity_level enum not created';
    END IF;

    RAISE NOTICE 'Enums created successfully';
END $$;

COMMIT;
```

### Down Migration

```sql
-- migrations/002_create_enums.down.sql
-- Rollback: Create enum types
-- Version: 002
-- Description: Drop all enum types

BEGIN;

-- Drop enums (in reverse order to handle dependencies)
DROP TYPE IF EXISTS ai_agent CASCADE;
DROP TYPE IF EXISTS audit_action CASCADE;
DROP TYPE IF EXISTS entity_type CASCADE;
DROP TYPE IF EXISTS resolution_type CASCADE;
DROP TYPE IF EXISTS dispute_status CASCADE;
DROP TYPE IF EXISTS complexity_level CASCADE;
DROP TYPE IF EXISTS step_status CASCADE;

-- Validation
DO $$
BEGIN
    -- Verify enums are dropped
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'step_status') THEN
        RAISE EXCEPTION 'Rollback failed: step_status enum still exists';
    END IF;

    RAISE NOTICE 'Enums dropped successfully';
END $$;

COMMIT;
```

## Migration 003: Create Tasks Table

### Up Migration

```sql
-- migrations/003_create_tasks_table.up.sql
-- Migration: Create tasks table
-- Version: 003
-- Description: Create tasks table with hierarchical relationships and polymorphic navigation

BEGIN;

-- Create tasks table
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    objective TEXT NOT NULL,
    progress DOUBLE PRECISION NOT NULL DEFAULT 0,
    project_id UUID NOT NULL,
    parent_task_id UUID,
    prev_id VARCHAR(50),  -- Polymorphic: "task://uuid" or "step://uuid"
    next_id VARCHAR(50),  -- Polymorphic: "task://uuid" or "step://uuid"
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add foreign key constraints
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_project_id
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;

ALTER TABLE tasks ADD CONSTRAINT fk_tasks_parent_task_id
    FOREIGN KEY (parent_task_id) REFERENCES tasks(id) ON DELETE CASCADE;

-- Add check constraints
ALTER TABLE tasks ADD CONSTRAINT check_tasks_progress
    CHECK (progress >= 0 AND progress <= 1);

ALTER TABLE tasks ADD CONSTRAINT check_tasks_title_not_empty
    CHECK (LENGTH(TRIM(title)) > 0);

ALTER TABLE tasks ADD CONSTRAINT check_tasks_objective_not_empty
    CHECK (LENGTH(TRIM(objective)) > 0);

-- Polymorphic reference format validation
ALTER TABLE tasks ADD CONSTRAINT check_tasks_prev_id_format
    CHECK (prev_id IS NULL OR prev_id ~ '^(task|step)://[0-9a-f-]{36}$');

ALTER TABLE tasks ADD CONSTRAINT check_tasks_next_id_format
    CHECK (next_id IS NULL OR next_id ~ '^(task|step)://[0-9a-f-]{36}$');

-- Prevent self-reference in parent_task_id
ALTER TABLE tasks ADD CONSTRAINT check_tasks_no_self_parent
    CHECK (parent_task_id != id);

-- Add comments
COMMENT ON TABLE tasks IS 'Hierarchical tasks that can contain steps and sub-tasks';
COMMENT ON COLUMN tasks.id IS 'Unique task identifier';
COMMENT ON COLUMN tasks.title IS 'Task title';
COMMENT ON COLUMN tasks.objective IS 'What this task should accomplish';
COMMENT ON COLUMN tasks.progress IS 'Task completion progress (0.0 to 1.0)';
COMMENT ON COLUMN tasks.project_id IS 'Reference to parent project';
COMMENT ON COLUMN tasks.parent_task_id IS 'Reference to parent task (NULL for root tasks)';
COMMENT ON COLUMN tasks.prev_id IS 'Previous item in workflow (polymorphic reference)';
COMMENT ON COLUMN tasks.next_id IS 'Next item in workflow (polymorphic reference)';

-- Create indexes
CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_parent_task_id ON tasks(parent_task_id);
CREATE INDEX idx_tasks_prev_id ON tasks(prev_id);
CREATE INDEX idx_tasks_next_id ON tasks(next_id);
CREATE INDEX idx_tasks_progress ON tasks(progress);
CREATE INDEX idx_tasks_created_at ON tasks(created_at);

-- Create index for root tasks (commonly queried)
CREATE INDEX idx_tasks_root_tasks ON tasks(project_id) WHERE parent_task_id IS NULL;

-- Validation
DO $$
BEGIN
    -- Verify table exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'tasks') THEN
        RAISE EXCEPTION 'Migration failed: tasks table not created';
    END IF;

    -- Verify foreign key constraints
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                   WHERE constraint_name = 'fk_tasks_project_id') THEN
        RAISE EXCEPTION 'Migration failed: project_id foreign key not created';
    END IF;

    RAISE NOTICE 'Tasks table created successfully';
END $$;

COMMIT;
```

### Down Migration

```sql
-- migrations/003_create_tasks_table.down.sql
-- Rollback: Create tasks table
-- Version: 003
-- Description: Drop tasks table and related objects

BEGIN;

-- Drop indexes
DROP INDEX IF EXISTS idx_tasks_root_tasks;
DROP INDEX IF EXISTS idx_tasks_created_at;
DROP INDEX IF EXISTS idx_tasks_progress;
DROP INDEX IF EXISTS idx_tasks_next_id;
DROP INDEX IF EXISTS idx_tasks_prev_id;
DROP INDEX IF EXISTS idx_tasks_parent_task_id;
DROP INDEX IF EXISTS idx_tasks_project_id;

-- Drop table
DROP TABLE IF EXISTS tasks CASCADE;

-- Validation
DO $$
BEGIN
    -- Verify table is dropped
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'tasks') THEN
        RAISE EXCEPTION 'Rollback failed: tasks table still exists';
    END IF;

    RAISE NOTICE 'Tasks table dropped successfully';
END $$;

COMMIT;
```

## Migration 004: Create Steps Table

### Up Migration

```sql
-- migrations/004_create_steps_table.up.sql
-- Migration: Create steps table
-- Version: 004
-- Description: Create steps table with dual-AI collaboration fields

BEGIN;

-- Create steps table
CREATE TABLE steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    task_id UUID NOT NULL,
    prev_id VARCHAR(50),  -- Polymorphic: "task://uuid" or "step://uuid"
    next_id VARCHAR(50),  -- Polymorphic: "task://uuid" or "step://uuid"
    progress DOUBLE PRECISION NOT NULL DEFAULT 0,

    -- Dual-AI collaboration content
    server_content TEXT,
    client_content TEXT,
    final_content TEXT,

    -- Iteration tracking
    iteration_count INTEGER NOT NULL DEFAULT 0,
    max_iterations_reached BOOLEAN NOT NULL DEFAULT FALSE,

    -- Collaboration state
    status step_status NOT NULL DEFAULT 'PENDING',
    server_ready BOOLEAN NOT NULL DEFAULT FALSE,
    client_approved BOOLEAN NOT NULL DEFAULT FALSE,

    -- Complexity assessment
    server_complexity complexity_level,
    client_complexity complexity_level,
    agreed_complexity complexity_level,
    complexity_score DOUBLE PRECISION,
    should_promote BOOLEAN,

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add foreign key constraints
ALTER TABLE steps ADD CONSTRAINT fk_steps_task_id
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE;

-- Add check constraints
ALTER TABLE steps ADD CONSTRAINT check_steps_progress
    CHECK (progress >= 0 AND progress <= 1);

ALTER TABLE steps ADD CONSTRAINT check_steps_title_not_empty
    CHECK (LENGTH(TRIM(title)) > 0);

ALTER TABLE steps ADD CONSTRAINT check_steps_iteration_count
    CHECK (iteration_count >= 0);

ALTER TABLE steps ADD CONSTRAINT check_steps_complexity_score
    CHECK (complexity_score IS NULL OR (complexity_score >= 0 AND complexity_score <= 1));

-- Polymorphic reference format validation
ALTER TABLE steps ADD CONSTRAINT check_steps_prev_id_format
    CHECK (prev_id IS NULL OR prev_id ~ '^(task|step)://[0-9a-f-]{36}$');

ALTER TABLE steps ADD CONSTRAINT check_steps_next_id_format
    CHECK (next_id IS NULL OR next_id ~ '^(task|step)://[0-9a-f-]{36}$');

-- Business logic constraints
ALTER TABLE steps ADD CONSTRAINT check_steps_final_content_when_agreed
    CHECK (status != 'AGREED' OR final_content IS NOT NULL);

-- Add comments
COMMENT ON TABLE steps IS 'Individual steps within tasks containing detailed implementation instructions';
COMMENT ON COLUMN steps.id IS 'Unique step identifier';
COMMENT ON COLUMN steps.title IS 'Step title';
COMMENT ON COLUMN steps.task_id IS 'Reference to parent task';
COMMENT ON COLUMN steps.prev_id IS 'Previous item in workflow (polymorphic reference)';
COMMENT ON COLUMN steps.next_id IS 'Next item in workflow (polymorphic reference)';
COMMENT ON COLUMN steps.progress IS 'Step completion progress (0.0 to 1.0)';
COMMENT ON COLUMN steps.server_content IS 'Content generated by server AI';
COMMENT ON COLUMN steps.client_content IS 'Content refined by client AI';
COMMENT ON COLUMN steps.final_content IS 'Final agreed upon content';
COMMENT ON COLUMN steps.iteration_count IS 'Number of AI collaboration iterations';
COMMENT ON COLUMN steps.status IS 'Current step status in collaboration workflow';
COMMENT ON COLUMN steps.complexity_score IS 'AI-assessed complexity score (0.0 to 1.0)';

-- Create indexes
CREATE INDEX idx_steps_task_id ON steps(task_id);
CREATE INDEX idx_steps_status ON steps(status);
CREATE INDEX idx_steps_server_ready ON steps(server_ready);
CREATE INDEX idx_steps_client_approved ON steps(client_approved);
CREATE INDEX idx_steps_prev_id ON steps(prev_id);
CREATE INDEX idx_steps_next_id ON steps(next_id);
CREATE INDEX idx_steps_complexity_score ON steps(complexity_score);
CREATE INDEX idx_steps_should_promote ON steps(should_promote);
CREATE INDEX idx_steps_progress ON steps(progress);
CREATE INDEX idx_steps_created_at ON steps(created_at);

-- Create composite indexes for common queries
CREATE INDEX idx_steps_task_status ON steps(task_id, status);
CREATE INDEX idx_steps_task_progress ON steps(task_id, progress);

-- Validation
DO $$
BEGIN
    -- Verify table exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'steps') THEN
        RAISE EXCEPTION 'Migration failed: steps table not created';
    END IF;

    -- Verify foreign key constraints
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                   WHERE constraint_name = 'fk_steps_task_id') THEN
        RAISE EXCEPTION 'Migration failed: task_id foreign key not created';
    END IF;

    -- Verify enum usage
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'steps' AND column_name = 'status'
                   AND udt_name = 'step_status') THEN
        RAISE EXCEPTION 'Migration failed: status column not using step_status enum';
    END IF;

    RAISE NOTICE 'Steps table created successfully';
END $$;

COMMIT;
```

### Down Migration

```sql
-- migrations/004_create_steps_table.down.sql
-- Rollback: Create steps table
-- Version: 004
-- Description: Drop steps table and related objects

BEGIN;

-- Drop indexes
DROP INDEX IF EXISTS idx_steps_task_progress;
DROP INDEX IF EXISTS idx_steps_task_status;
DROP INDEX IF EXISTS idx_steps_created_at;
DROP INDEX IF EXISTS idx_steps_progress;
DROP INDEX IF EXISTS idx_steps_should_promote;
DROP INDEX IF EXISTS idx_steps_complexity_score;
DROP INDEX IF EXISTS idx_steps_next_id;
DROP INDEX IF EXISTS idx_steps_prev_id;
DROP INDEX IF EXISTS idx_steps_client_approved;
DROP INDEX IF EXISTS idx_steps_server_ready;
DROP INDEX IF EXISTS idx_steps_status;
DROP INDEX IF EXISTS idx_steps_task_id;

-- Drop table
DROP TABLE IF EXISTS steps CASCADE;

-- Validation
DO $$
BEGIN
    -- Verify table is dropped
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'steps') THEN
        RAISE EXCEPTION 'Rollback failed: steps table still exists';
    END IF;

    RAISE NOTICE 'Steps table dropped successfully';
END $$;

COMMIT;
```

## Migration 005: Create Disputes Table

### Up Migration

```sql
-- migrations/005_create_disputes_table.up.sql
-- Migration: Create disputes table
-- Version: 005
-- Description: Create disputes table for AI collaboration conflict resolution

BEGIN;

-- Create disputes table
CREATE TABLE disputes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL,
    step_id UUID NOT NULL,
    server_content TEXT NOT NULL,
    client_content TEXT NOT NULL,
    server_reasoning TEXT NOT NULL,
    client_reasoning TEXT NOT NULL,
    status dispute_status NOT NULL DEFAULT 'PENDING',
    user_resolution resolution_type,
    resolved_content TEXT,
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add foreign key constraints
ALTER TABLE disputes ADD CONSTRAINT fk_disputes_project_id
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;

ALTER TABLE disputes ADD CONSTRAINT fk_disputes_step_id
    FOREIGN KEY (step_id) REFERENCES steps(id) ON DELETE CASCADE;

-- Add check constraints
ALTER TABLE disputes ADD CONSTRAINT check_disputes_content_not_empty
    CHECK (LENGTH(TRIM(server_content)) > 0 AND LENGTH(TRIM(client_content)) > 0);

ALTER TABLE disputes ADD CONSTRAINT check_disputes_reasoning_not_empty
    CHECK (LENGTH(TRIM(server_reasoning)) > 0 AND LENGTH(TRIM(client_reasoning)) > 0);

-- Business logic constraints
ALTER TABLE disputes ADD CONSTRAINT check_disputes_resolved_content_when_resolved
    CHECK (status != 'RESOLVED' OR (resolved_content IS NOT NULL AND resolved_at IS NOT NULL));

ALTER TABLE disputes ADD CONSTRAINT check_disputes_resolution_type_when_resolved
    CHECK (status != 'RESOLVED' OR user_resolution IS NOT NULL);

-- Add comments
COMMENT ON TABLE disputes IS 'Content disputes between AI agents requiring human resolution';
COMMENT ON COLUMN disputes.id IS 'Unique dispute identifier';
COMMENT ON COLUMN disputes.project_id IS 'Reference to project for context';
COMMENT ON COLUMN disputes.step_id IS 'Reference to disputed step';
COMMENT ON COLUMN disputes.server_content IS 'Server AI final content version';
COMMENT ON COLUMN disputes.client_content IS 'Client AI final content version';
COMMENT ON COLUMN disputes.server_reasoning IS 'Server AI reasoning for their version';
COMMENT ON COLUMN disputes.client_reasoning IS 'Client AI reasoning for their version';
COMMENT ON COLUMN disputes.status IS 'Current dispute status';
COMMENT ON COLUMN disputes.user_resolution IS 'Type of resolution chosen by user';
COMMENT ON COLUMN disputes.resolved_content IS 'Final content after user resolution';

-- Create indexes
CREATE INDEX idx_disputes_project_id ON disputes(project_id);
CREATE INDEX idx_disputes_step_id ON disputes(step_id);
CREATE INDEX idx_disputes_status ON disputes(status);
CREATE INDEX idx_disputes_created_at ON disputes(created_at);
CREATE INDEX idx_disputes_resolved_at ON disputes(resolved_at);

-- Create composite indexes for common queries
CREATE INDEX idx_disputes_project_status ON disputes(project_id, status);
CREATE INDEX idx_disputes_pending ON disputes(status) WHERE status = 'PENDING';

-- Validation
DO $$
BEGIN
    -- Verify table exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'disputes') THEN
        RAISE EXCEPTION 'Migration failed: disputes table not created';
    END IF;

    -- Verify foreign key constraints
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                   WHERE constraint_name = 'fk_disputes_project_id') THEN
        RAISE EXCEPTION 'Migration failed: project_id foreign key not created';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                   WHERE constraint_name = 'fk_disputes_step_id') THEN
        RAISE EXCEPTION 'Migration failed: step_id foreign key not created';
    END IF;

    RAISE NOTICE 'Disputes table created successfully';
END $$;

COMMIT;
```

### Down Migration

```sql
-- migrations/005_create_disputes_table.down.sql
-- Rollback: Create disputes table
-- Version: 005
-- Description: Drop disputes table and related objects

BEGIN;

-- Drop indexes
DROP INDEX IF EXISTS idx_disputes_pending;
DROP INDEX IF EXISTS idx_disputes_project_status;
DROP INDEX IF EXISTS idx_disputes_resolved_at;
DROP INDEX IF EXISTS idx_disputes_created_at;
DROP INDEX IF EXISTS idx_disputes_status;
DROP INDEX IF EXISTS idx_disputes_step_id;
DROP INDEX IF EXISTS idx_disputes_project_id;

-- Drop table
DROP TABLE IF EXISTS disputes CASCADE;

-- Validation
DO $$
BEGIN
    -- Verify table is dropped
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'disputes') THEN
        RAISE EXCEPTION 'Rollback failed: disputes table still exists';
    END IF;

    RAISE NOTICE 'Disputes table dropped successfully';
END $$;

COMMIT;
```

## Migration 006: Create Audit Logs Table

### Up Migration

```sql
-- migrations/006_create_audit_logs_table.up.sql
-- Migration: Create audit logs table
-- Version: 006
-- Description: Create audit logs table for tracking all system changes

BEGIN;

-- Create audit logs table
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type entity_type NOT NULL,
    entity_id UUID NOT NULL,
    action audit_action NOT NULL,
    old_values JSONB,
    new_values JSONB,
    user_id VARCHAR(255),
    ai_agent ai_agent,
    metadata JSONB,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add check constraints
ALTER TABLE audit_logs ADD CONSTRAINT check_audit_logs_actor
    CHECK (user_id IS NOT NULL OR ai_agent IS NOT NULL);

-- Add comments
COMMENT ON TABLE audit_logs IS 'Comprehensive audit trail for all system changes';
COMMENT ON COLUMN audit_logs.id IS 'Unique audit log identifier';
COMMENT ON COLUMN audit_logs.entity_type IS 'Type of entity that was modified';
COMMENT ON COLUMN audit_logs.entity_id IS 'ID of the entity that was modified';
COMMENT ON COLUMN audit_logs.action IS 'Action that was performed';
COMMENT ON COLUMN audit_logs.old_values IS 'Previous state of the entity (JSON)';
COMMENT ON COLUMN audit_logs.new_values IS 'New state of the entity (JSON)';
COMMENT ON COLUMN audit_logs.user_id IS 'User who performed the action (if applicable)';
COMMENT ON COLUMN audit_logs.ai_agent IS 'AI agent that performed the action (if applicable)';
COMMENT ON COLUMN audit_logs.metadata IS 'Additional context and metadata (JSON)';

-- Create indexes
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_ai_agent ON audit_logs(ai_agent);

-- Create composite indexes for common queries
CREATE INDEX idx_audit_logs_entity_timestamp ON audit_logs(entity_type, entity_id, timestamp);
CREATE INDEX idx_audit_logs_recent ON audit_logs(timestamp DESC) WHERE timestamp > NOW() - INTERVAL '30 days';

-- Create GIN indexes for JSON columns
CREATE INDEX idx_audit_logs_old_values_gin ON audit_logs USING GIN(old_values);
CREATE INDEX idx_audit_logs_new_values_gin ON audit_logs USING GIN(new_values);
CREATE INDEX idx_audit_logs_metadata_gin ON audit_logs USING GIN(metadata);

-- Validation
DO $$
BEGIN
    -- Verify table exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'audit_logs') THEN
        RAISE EXCEPTION 'Migration failed: audit_logs table not created';
    END IF;

    -- Verify enum usage
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'audit_logs' AND column_name = 'entity_type'
                   AND udt_name = 'entity_type') THEN
        RAISE EXCEPTION 'Migration failed: entity_type column not using entity_type enum';
    END IF;

    RAISE NOTICE 'Audit logs table created successfully';
END $$;

COMMIT;
```

### Down Migration

```sql
-- migrations/006_create_audit_logs_table.down.sql
-- Rollback: Create audit logs table
-- Version: 006
-- Description: Drop audit logs table and related objects

BEGIN;

-- Drop indexes
DROP INDEX IF EXISTS idx_audit_logs_metadata_gin;
DROP INDEX IF EXISTS idx_audit_logs_new_values_gin;
DROP INDEX IF EXISTS idx_audit_logs_old_values_gin;
DROP INDEX IF EXISTS idx_audit_logs_recent;
DROP INDEX IF EXISTS idx_audit_logs_entity_timestamp;
DROP INDEX IF EXISTS idx_audit_logs_ai_agent;
DROP INDEX IF EXISTS idx_audit_logs_user_id;
DROP INDEX IF EXISTS idx_audit_logs_action;
DROP INDEX IF EXISTS idx_audit_logs_timestamp;
DROP INDEX IF EXISTS idx_audit_logs_entity;

-- Drop table
DROP TABLE IF EXISTS audit_logs CASCADE;

-- Validation
DO $$
BEGIN
    -- Verify table is dropped
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'audit_logs') THEN
        RAISE EXCEPTION 'Rollback failed: audit_logs table still exists';
    END IF;

    RAISE NOTICE 'Audit logs table dropped successfully';
END $$;

COMMIT;
```

## Migration Validation Script

```sql
-- migrations/validate_initial_schema.sql
-- Validation script for initial schema migrations
-- Run after all initial migrations to verify schema integrity

DO $$
DECLARE
    table_count INTEGER;
    enum_count INTEGER;
    constraint_count INTEGER;
    index_count INTEGER;
BEGIN
    -- Count tables
    SELECT COUNT(*) INTO table_count
    FROM information_schema.tables
    WHERE table_schema = 'public'
    AND table_name IN ('projects', 'tasks', 'steps', 'disputes', 'audit_logs');

    IF table_count != 5 THEN
        RAISE EXCEPTION 'Expected 5 tables, found %', table_count;
    END IF;

    -- Count enums
    SELECT COUNT(*) INTO enum_count
    FROM pg_type
    WHERE typname IN ('step_status', 'complexity_level', 'dispute_status',
                      'resolution_type', 'entity_type', 'audit_action', 'ai_agent');

    IF enum_count != 7 THEN
        RAISE EXCEPTION 'Expected 7 enums, found %', enum_count;
    END IF;

    -- Count constraints (approximate check)
    SELECT COUNT(*) INTO constraint_count
    FROM information_schema.table_constraints
    WHERE table_schema = 'public'
    AND constraint_type IN ('CHECK', 'FOREIGN KEY');

    IF constraint_count < 20 THEN
        RAISE EXCEPTION 'Expected at least 20 constraints, found %', constraint_count;
    END IF;

    -- Count indexes (approximate check)
    SELECT COUNT(*) INTO index_count
    FROM pg_indexes
    WHERE schemaname = 'public';

    IF index_count < 30 THEN
        RAISE EXCEPTION 'Expected at least 30 indexes, found %', index_count;
    END IF;

    RAISE NOTICE 'Initial schema validation passed:';
    RAISE NOTICE '  Tables: %', table_count;
    RAISE NOTICE '  Enums: %', enum_count;
    RAISE NOTICE '  Constraints: %', constraint_count;
    RAISE NOTICE '  Indexes: %', index_count;
END $$;
```

---

*Next: [Index & Constraint Migrations](./07b3c-indexes-constraints.md)*
