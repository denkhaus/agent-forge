# Index & Constraint Migrations

## Overview

This document covers performance optimization indexes and additional data integrity constraints for the MCP-Planner system. These migrations enhance query performance and ensure data consistency beyond the basic schema.

## Migration 007: Performance Indexes

### Up Migration

```sql
-- migrations/007_add_performance_indexes.up.sql
-- Migration: Add performance indexes
-- Version: 007
-- Description: Add specialized indexes for query optimization and performance

BEGIN;

-- Navigation chain indexes for polymorphic references
-- These are critical for traversing task/step chains efficiently

-- Index for finding items that reference a specific task
CREATE INDEX idx_navigation_task_references ON tasks(prev_id)
WHERE prev_id LIKE 'task://%';

CREATE INDEX idx_navigation_task_references_steps ON steps(prev_id)
WHERE prev_id LIKE 'task://%';

-- Index for finding items that reference a specific step
CREATE INDEX idx_navigation_step_references ON tasks(prev_id)
WHERE prev_id LIKE 'step://%';

CREATE INDEX idx_navigation_step_references_steps ON steps(prev_id)
WHERE prev_id LIKE 'step://%';

-- Composite indexes for common query patterns

-- Project dashboard queries
CREATE INDEX idx_projects_dashboard ON projects(id, progress, created_at);

-- Task hierarchy queries
CREATE INDEX idx_tasks_hierarchy ON tasks(project_id, parent_task_id, progress);

-- Step workflow queries
CREATE INDEX idx_steps_workflow ON steps(task_id, status, progress, created_at);

-- AI collaboration queries
CREATE INDEX idx_steps_collaboration ON steps(status, server_ready, client_approved, iteration_count);

-- Dispute resolution queries
CREATE INDEX idx_disputes_resolution ON disputes(project_id, status, created_at);

-- Progress calculation indexes
-- These support efficient hierarchical progress calculation

-- Root tasks for project progress
CREATE INDEX idx_tasks_project_roots ON tasks(project_id, progress)
WHERE parent_task_id IS NULL;

-- Sub-tasks for task progress
CREATE INDEX idx_tasks_subtasks ON tasks(parent_task_id, progress)
WHERE parent_task_id IS NOT NULL;

-- Steps for task progress
CREATE INDEX idx_steps_task_progress ON steps(task_id, progress);

-- Complexity analysis indexes

-- Steps needing complexity analysis
CREATE INDEX idx_steps_complexity_pending ON steps(task_id, complexity_score, created_at)
WHERE complexity_score IS NULL;

-- Steps above complexity threshold (for promotion)
CREATE INDEX idx_steps_complex ON steps(task_id, complexity_score, should_promote)
WHERE complexity_score > 0.7;

-- Steps ready for promotion
CREATE INDEX idx_steps_promotion_ready ON steps(task_id, should_promote, status)
WHERE should_promote = TRUE;

-- Audit and monitoring indexes

-- Recent audit logs (for monitoring)
CREATE INDEX idx_audit_logs_recent_activity ON audit_logs(timestamp DESC, entity_type, action)
WHERE timestamp > NOW() - INTERVAL '7 days';

-- User activity tracking
CREATE INDEX idx_audit_logs_user_activity ON audit_logs(user_id, timestamp DESC)
WHERE user_id IS NOT NULL;

-- AI agent activity tracking
CREATE INDEX idx_audit_logs_ai_activity ON audit_logs(ai_agent, timestamp DESC)
WHERE ai_agent IS NOT NULL;

-- Search and filtering indexes

-- Text search preparation (for future full-text search)
CREATE INDEX idx_projects_text_search ON projects USING GIN(to_tsvector('english', name || ' ' || description));
CREATE INDEX idx_tasks_text_search ON tasks USING GIN(to_tsvector('english', title || ' ' || objective));
CREATE INDEX idx_steps_text_search ON steps USING GIN(to_tsvector('english', title || ' ' || COALESCE(final_content, '')));

-- Status-based filtering
CREATE INDEX idx_steps_by_status ON steps(status, task_id, created_at);
CREATE INDEX idx_disputes_by_status ON disputes(status, project_id, created_at);

-- Time-based queries
CREATE INDEX idx_projects_by_date ON projects(created_at, updated_at);
CREATE INDEX idx_tasks_by_date ON tasks(created_at, updated_at);
CREATE INDEX idx_steps_by_date ON steps(created_at, updated_at);

-- Partial indexes for specific use cases

-- Active projects (non-completed)
CREATE INDEX idx_projects_active ON projects(id, name, progress)
WHERE progress < 1.0;

-- Pending steps (not completed)
CREATE INDEX idx_steps_pending ON steps(task_id, title, created_at)
WHERE progress < 1.0;

-- Disputed steps requiring attention
CREATE INDEX idx_steps_disputed ON steps(task_id, title, created_at)
WHERE status = 'DISPUTED';

-- Steps in AI collaboration
CREATE INDEX idx_steps_ai_collaboration ON steps(task_id, status, iteration_count)
WHERE status IN ('SERVER_DRAFT', 'CLIENT_REVIEW');

-- Covering indexes for read-heavy queries

-- Project summary (covers most dashboard queries)
CREATE INDEX idx_projects_summary ON projects(id)
INCLUDE (name, description, progress, complexity_threshold, created_at);

-- Task summary (covers most task list queries)
CREATE INDEX idx_tasks_summary ON tasks(project_id)
INCLUDE (id, title, objective, progress, parent_task_id, created_at);

-- Step summary (covers most step list queries)
CREATE INDEX idx_steps_summary ON steps(task_id)
INCLUDE (id, title, progress, status, created_at);

-- Add comments for documentation
COMMENT ON INDEX idx_navigation_task_references IS 'Optimizes polymorphic navigation queries for task references';
COMMENT ON INDEX idx_projects_dashboard IS 'Optimizes project dashboard queries';
COMMENT ON INDEX idx_tasks_hierarchy IS 'Optimizes task hierarchy traversal';
COMMENT ON INDEX idx_steps_workflow IS 'Optimizes step workflow queries';
COMMENT ON INDEX idx_steps_collaboration IS 'Optimizes AI collaboration status queries';
COMMENT ON INDEX idx_projects_text_search IS 'Enables full-text search on projects';
COMMENT ON INDEX idx_steps_complex IS 'Identifies steps above complexity threshold';

-- Validation
DO $$
DECLARE
    index_count INTEGER;
BEGIN
    -- Count new indexes created in this migration
    SELECT COUNT(*) INTO index_count
    FROM pg_indexes
    WHERE schemaname = 'public'
    AND indexname LIKE 'idx_%'
    AND indexname NOT IN (
        -- Exclude indexes from previous migrations
        'idx_projects_created_at', 'idx_projects_progress', 'idx_projects_complexity_threshold',
        'idx_tasks_project_id', 'idx_tasks_parent_task_id', 'idx_tasks_prev_id', 'idx_tasks_next_id',
        'idx_tasks_progress', 'idx_tasks_created_at', 'idx_tasks_root_tasks',
        'idx_steps_task_id', 'idx_steps_status', 'idx_steps_server_ready', 'idx_steps_client_approved',
        'idx_steps_prev_id', 'idx_steps_next_id', 'idx_steps_complexity_score', 'idx_steps_should_promote',
        'idx_steps_progress', 'idx_steps_created_at', 'idx_steps_task_status', 'idx_steps_task_progress',
        'idx_disputes_project_id', 'idx_disputes_step_id', 'idx_disputes_status', 'idx_disputes_created_at',
        'idx_disputes_resolved_at', 'idx_disputes_project_status', 'idx_disputes_pending',
        'idx_audit_logs_entity', 'idx_audit_logs_timestamp', 'idx_audit_logs_action',
        'idx_audit_logs_user_id', 'idx_audit_logs_ai_agent', 'idx_audit_logs_entity_timestamp',
        'idx_audit_logs_recent', 'idx_audit_logs_old_values_gin', 'idx_audit_logs_new_values_gin',
        'idx_audit_logs_metadata_gin'
    );

    IF index_count < 25 THEN
        RAISE EXCEPTION 'Expected at least 25 new indexes, found %', index_count;
    END IF;

    RAISE NOTICE 'Performance indexes created successfully: % new indexes', index_count;
END $$;

COMMIT;
```

### Down Migration

```sql
-- migrations/007_add_performance_indexes.down.sql
-- Rollback: Add performance indexes
-- Version: 007
-- Description: Drop performance optimization indexes

BEGIN;

-- Drop covering indexes
DROP INDEX IF EXISTS idx_steps_summary;
DROP INDEX IF EXISTS idx_tasks_summary;
DROP INDEX IF EXISTS idx_projects_summary;

-- Drop partial indexes
DROP INDEX IF EXISTS idx_steps_ai_collaboration;
DROP INDEX IF EXISTS idx_steps_disputed;
DROP INDEX IF EXISTS idx_steps_pending;
DROP INDEX IF EXISTS idx_projects_active;

-- Drop time-based indexes
DROP INDEX IF EXISTS idx_steps_by_date;
DROP INDEX IF EXISTS idx_tasks_by_date;
DROP INDEX IF EXISTS idx_projects_by_date;

-- Drop status-based indexes
DROP INDEX IF EXISTS idx_disputes_by_status;
DROP INDEX IF EXISTS idx_steps_by_status;

-- Drop text search indexes
DROP INDEX IF EXISTS idx_steps_text_search;
DROP INDEX IF EXISTS idx_tasks_text_search;
DROP INDEX IF EXISTS idx_projects_text_search;

-- Drop audit and monitoring indexes
DROP INDEX IF EXISTS idx_audit_logs_ai_activity;
DROP INDEX IF EXISTS idx_audit_logs_user_activity;
DROP INDEX IF EXISTS idx_audit_logs_recent_activity;

-- Drop complexity analysis indexes
DROP INDEX IF EXISTS idx_steps_promotion_ready;
DROP INDEX IF EXISTS idx_steps_complex;
DROP INDEX IF EXISTS idx_steps_complexity_pending;

-- Drop progress calculation indexes
DROP INDEX IF EXISTS idx_steps_task_progress;
DROP INDEX IF EXISTS idx_tasks_subtasks;
DROP INDEX IF EXISTS idx_tasks_project_roots;

-- Drop composite indexes
DROP INDEX IF EXISTS idx_disputes_resolution;
DROP INDEX IF EXISTS idx_steps_collaboration;
DROP INDEX IF EXISTS idx_steps_workflow;
DROP INDEX IF EXISTS idx_tasks_hierarchy;
DROP INDEX IF EXISTS idx_projects_dashboard;

-- Drop navigation indexes
DROP INDEX IF EXISTS idx_navigation_step_references_steps;
DROP INDEX IF EXISTS idx_navigation_step_references;
DROP INDEX IF EXISTS idx_navigation_task_references_steps;
DROP INDEX IF EXISTS idx_navigation_task_references;

-- Validation
DO $$
BEGIN
    RAISE NOTICE 'Performance indexes dropped successfully';
END $$;

COMMIT;
```

## Migration 008: Advanced Constraints

### Up Migration

```sql
-- migrations/008_add_advanced_constraints.up.sql
-- Migration: Add advanced constraints
-- Version: 008
-- Description: Add complex business logic constraints and data integrity rules

BEGIN;

-- Circular reference prevention constraints

-- Prevent circular task hierarchies
CREATE OR REPLACE FUNCTION check_task_hierarchy_circular()
RETURNS TRIGGER AS $$
DECLARE
    current_id UUID;
    depth INTEGER := 0;
    max_depth INTEGER := 100; -- Prevent infinite loops
BEGIN
    -- Only check if parent_task_id is being set
    IF NEW.parent_task_id IS NULL THEN
        RETURN NEW;
    END IF;

    -- Start from the new parent and traverse up the hierarchy
    current_id := NEW.parent_task_id;

    WHILE current_id IS NOT NULL AND depth < max_depth LOOP
        -- Check if we've reached the task being updated (circular reference)
        IF current_id = NEW.id THEN
            RAISE EXCEPTION 'Circular reference detected in task hierarchy: task % cannot be a descendant of itself', NEW.id;
        END IF;

        -- Move up one level
        SELECT parent_task_id INTO current_id
        FROM tasks
        WHERE id = current_id;

        depth := depth + 1;
    END LOOP;

    -- Check for maximum depth exceeded
    IF depth >= max_depth THEN
        RAISE EXCEPTION 'Task hierarchy too deep (max % levels)', max_depth;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for circular reference check
CREATE TRIGGER trigger_check_task_hierarchy_circular
    BEFORE INSERT OR UPDATE OF parent_task_id ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION check_task_hierarchy_circular();

-- Navigation chain integrity constraints

-- Function to validate polymorphic references
CREATE OR REPLACE FUNCTION validate_polymorphic_reference(ref_value TEXT)
RETURNS BOOLEAN AS $$
BEGIN
    IF ref_value IS NULL THEN
        RETURN TRUE;
    END IF;

    -- Check format
    IF ref_value !~ '^(task|step)://[0-9a-f-]{36}$' THEN
        RETURN FALSE;
    END IF;

    -- Extract type and ID
    DECLARE
        ref_type TEXT;
        ref_id UUID;
    BEGIN
        ref_type := split_part(ref_value, '://', 1);
        ref_id := split_part(ref_value, '://', 2)::UUID;

        -- Check if referenced entity exists
        IF ref_type = 'task' THEN
            IF NOT EXISTS (SELECT 1 FROM tasks WHERE id = ref_id) THEN
                RETURN FALSE;
            END IF;
        ELSIF ref_type = 'step' THEN
            IF NOT EXISTS (SELECT 1 FROM steps WHERE id = ref_id) THEN
                RETURN FALSE;
            END IF;
        END IF;

        RETURN TRUE;
    EXCEPTION
        WHEN OTHERS THEN
            RETURN FALSE;
    END;
END;
$$ LANGUAGE plpgsql;

-- Function to check navigation chain integrity
CREATE OR REPLACE FUNCTION check_navigation_integrity()
RETURNS TRIGGER AS $$
BEGIN
    -- Validate prev_id reference
    IF NOT validate_polymorphic_reference(NEW.prev_id) THEN
        RAISE EXCEPTION 'Invalid prev_id reference: %', NEW.prev_id;
    END IF;

    -- Validate next_id reference
    IF NOT validate_polymorphic_reference(NEW.next_id) THEN
        RAISE EXCEPTION 'Invalid next_id reference: %', NEW.next_id;
    END IF;

    -- Prevent self-reference
    IF NEW.prev_id = (TG_TABLE_NAME || '://' || NEW.id::TEXT) THEN
        RAISE EXCEPTION 'Self-reference not allowed in prev_id';
    END IF;

    IF NEW.next_id = (TG_TABLE_NAME || '://' || NEW.id::TEXT) THEN
        RAISE EXCEPTION 'Self-reference not allowed in next_id';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for navigation integrity
CREATE TRIGGER trigger_check_tasks_navigation_integrity
    BEFORE INSERT OR UPDATE OF prev_id, next_id ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION check_navigation_integrity();

CREATE TRIGGER trigger_check_steps_navigation_integrity
    BEFORE INSERT OR UPDATE OF prev_id, next_id ON steps
    FOR EACH ROW
    EXECUTE FUNCTION check_navigation_integrity();

-- Business logic constraints

-- Ensure steps belong to tasks in the same project
CREATE OR REPLACE FUNCTION check_step_task_project_consistency()
RETURNS TRIGGER AS $$
DECLARE
    task_project_id UUID;
BEGIN
    -- Get the project_id of the task
    SELECT project_id INTO task_project_id
    FROM tasks
    WHERE id = NEW.task_id;

    IF task_project_id IS NULL THEN
        RAISE EXCEPTION 'Task % does not exist', NEW.task_id;
    END IF;

    -- Store project_id in step for future reference (if we add this column)
    -- For now, just validate the relationship exists

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_step_task_project_consistency
    BEFORE INSERT OR UPDATE OF task_id ON steps
    FOR EACH ROW
    EXECUTE FUNCTION check_step_task_project_consistency();

-- Ensure sub-tasks belong to the same project as parent
CREATE OR REPLACE FUNCTION check_subtask_project_consistency()
RETURNS TRIGGER AS $$
DECLARE
    parent_project_id UUID;
BEGIN
    -- Only check if parent_task_id is set
    IF NEW.parent_task_id IS NULL THEN
        RETURN NEW;
    END IF;

    -- Get parent task's project_id
    SELECT project_id INTO parent_project_id
    FROM tasks
    WHERE id = NEW.parent_task_id;

    IF parent_project_id IS NULL THEN
        RAISE EXCEPTION 'Parent task % does not exist', NEW.parent_task_id;
    END IF;

    -- Ensure same project
    IF NEW.project_id != parent_project_id THEN
        RAISE EXCEPTION 'Sub-task must belong to same project as parent task';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_subtask_project_consistency
    BEFORE INSERT OR UPDATE OF project_id, parent_task_id ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION check_subtask_project_consistency();

-- AI collaboration workflow constraints

-- Ensure proper step status transitions
CREATE OR REPLACE FUNCTION check_step_status_transition()
RETURNS TRIGGER AS $$
BEGIN
    -- Allow any transition on INSERT
    IF TG_OP = 'INSERT' THEN
        RETURN NEW;
    END IF;

    -- Define valid transitions
    CASE OLD.status
        WHEN 'PENDING' THEN
            IF NEW.status NOT IN ('PENDING', 'SERVER_DRAFT') THEN
                RAISE EXCEPTION 'Invalid status transition from PENDING to %', NEW.status;
            END IF;
        WHEN 'SERVER_DRAFT' THEN
            IF NEW.status NOT IN ('SERVER_DRAFT', 'CLIENT_REVIEW', 'DISPUTED') THEN
                RAISE EXCEPTION 'Invalid status transition from SERVER_DRAFT to %', NEW.status;
            END IF;
        WHEN 'CLIENT_REVIEW' THEN
            IF NEW.status NOT IN ('CLIENT_REVIEW', 'AGREED', 'SERVER_DRAFT', 'DISPUTED') THEN
                RAISE EXCEPTION 'Invalid status transition from CLIENT_REVIEW to %', NEW.status;
            END IF;
        WHEN 'AGREED' THEN
            IF NEW.status NOT IN ('AGREED') THEN
                RAISE EXCEPTION 'Cannot change status from AGREED to %', NEW.status;
            END IF;
        WHEN 'DISPUTED' THEN
            IF NEW.status NOT IN ('DISPUTED', 'USER_RESOLUTION') THEN
                RAISE EXCEPTION 'Invalid status transition from DISPUTED to %', NEW.status;
            END IF;
        WHEN 'USER_RESOLUTION' THEN
            IF NEW.status NOT IN ('USER_RESOLUTION', 'AGREED') THEN
                RAISE EXCEPTION 'Invalid status transition from USER_RESOLUTION to %', NEW.status;
            END IF;
    END CASE;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_step_status_transition
    BEFORE UPDATE OF status ON steps
    FOR EACH ROW
    EXECUTE FUNCTION check_step_status_transition();

-- Ensure final_content is set when status is AGREED
CREATE OR REPLACE FUNCTION check_step_final_content()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'AGREED' AND (NEW.final_content IS NULL OR LENGTH(TRIM(NEW.final_content)) = 0) THEN
        RAISE EXCEPTION 'final_content must be set when status is AGREED';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_step_final_content
    BEFORE INSERT OR UPDATE ON steps
    FOR EACH ROW
    EXECUTE FUNCTION check_step_final_content();

-- Dispute resolution constraints

-- Ensure dispute resolution data consistency
CREATE OR REPLACE FUNCTION check_dispute_resolution_consistency()
RETURNS TRIGGER AS $$
BEGIN
    -- When status changes to RESOLVED, ensure required fields are set
    IF NEW.status = 'RESOLVED' THEN
        IF NEW.user_resolution IS NULL THEN
            RAISE EXCEPTION 'user_resolution must be set when dispute is resolved';
        END IF;

        IF NEW.resolved_content IS NULL OR LENGTH(TRIM(NEW.resolved_content)) = 0 THEN
            RAISE EXCEPTION 'resolved_content must be set when dispute is resolved';
        END IF;

        IF NEW.resolved_at IS NULL THEN
            NEW.resolved_at := NOW();
        END IF;
    END IF;

    -- When status is not RESOLVED, clear resolution fields
    IF NEW.status != 'RESOLVED' THEN
        NEW.user_resolution := NULL;
        NEW.resolved_content := NULL;
        NEW.resolved_at := NULL;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_dispute_resolution_consistency
    BEFORE INSERT OR UPDATE ON disputes
    FOR EACH ROW
    EXECUTE FUNCTION check_dispute_resolution_consistency();

-- Data consistency constraints

-- Ensure only one dispute per step at a time
ALTER TABLE disputes ADD CONSTRAINT unique_active_dispute_per_step
    EXCLUDE (step_id WITH =) WHERE (status = 'PENDING');

-- Ensure progress values are consistent with status
CREATE OR REPLACE FUNCTION check_progress_status_consistency()
RETURNS TRIGGER AS $$
BEGIN
    -- For steps: progress should be 1.0 when status is AGREED
    IF TG_TABLE_NAME = 'steps' THEN
        IF NEW.status = 'AGREED' AND NEW.progress != 1.0 THEN
            NEW.progress := 1.0;
        ELSIF NEW.status != 'AGREED' AND NEW.progress = 1.0 THEN
            -- Allow manual progress setting, but warn
            RAISE NOTICE 'Step % marked as complete but status is %', NEW.id, NEW.status;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_progress_status_consistency
    BEFORE INSERT OR UPDATE ON steps
    FOR EACH ROW
    EXECUTE FUNCTION check_progress_status_consistency();

-- Add comments for documentation
COMMENT ON FUNCTION check_task_hierarchy_circular() IS 'Prevents circular references in task hierarchy';
COMMENT ON FUNCTION validate_polymorphic_reference(TEXT) IS 'Validates polymorphic reference format and existence';
COMMENT ON FUNCTION check_navigation_integrity() IS 'Ensures navigation chain integrity';
COMMENT ON FUNCTION check_step_status_transition() IS 'Enforces valid step status transitions';
COMMENT ON CONSTRAINT unique_active_dispute_per_step ON disputes IS 'Ensures only one active dispute per step';

-- Validation
DO $$
DECLARE
    function_count INTEGER;
    trigger_count INTEGER;
    constraint_count INTEGER;
BEGIN
    -- Count new functions
    SELECT COUNT(*) INTO function_count
    FROM pg_proc
    WHERE proname LIKE 'check_%'
    AND pronamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'public');

    IF function_count < 8 THEN
        RAISE EXCEPTION 'Expected at least 8 constraint functions, found %', function_count;
    END IF;

    -- Count new triggers
    SELECT COUNT(*) INTO trigger_count
    FROM pg_trigger
    WHERE tgname LIKE 'trigger_check_%';

    IF trigger_count < 8 THEN
        RAISE EXCEPTION 'Expected at least 8 constraint triggers, found %', trigger_count;
    END IF;

    RAISE NOTICE 'Advanced constraints created successfully:';
    RAISE NOTICE '  Functions: %', function_count;
    RAISE NOTICE '  Triggers: %', trigger_count;
END $$;

COMMIT;
```

### Down Migration

```sql
-- migrations/008_add_advanced_constraints.down.sql
-- Rollback: Add advanced constraints
-- Version: 008
-- Description: Drop advanced business logic constraints

BEGIN;

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_check_progress_status_consistency ON steps;
DROP TRIGGER IF EXISTS trigger_check_dispute_resolution_consistency ON disputes;
DROP TRIGGER IF EXISTS trigger_check_step_final_content ON steps;
DROP TRIGGER IF EXISTS trigger_check_step_status_transition ON steps;
DROP TRIGGER IF EXISTS trigger_check_subtask_project_consistency ON tasks;
DROP TRIGGER IF EXISTS trigger_check_step_task_project_consistency ON steps;
DROP TRIGGER IF EXISTS trigger_check_steps_navigation_integrity ON steps;
DROP TRIGGER IF EXISTS trigger_check_tasks_navigation_integrity ON tasks;
DROP TRIGGER IF EXISTS trigger_check_task_hierarchy_circular ON tasks;

-- Drop constraints
ALTER TABLE disputes DROP CONSTRAINT IF EXISTS unique_active_dispute_per_step;

-- Drop functions
DROP FUNCTION IF EXISTS check_progress_status_consistency();
DROP FUNCTION IF EXISTS check_dispute_resolution_consistency();
DROP FUNCTION IF EXISTS check_step_final_content();
DROP FUNCTION IF EXISTS check_step_status_transition();
DROP FUNCTION IF EXISTS check_subtask_project_consistency();
DROP FUNCTION IF EXISTS check_step_task_project_consistency();
DROP FUNCTION IF EXISTS check_navigation_integrity();
DROP FUNCTION IF EXISTS validate_polymorphic_reference(TEXT);
DROP FUNCTION IF EXISTS check_task_hierarchy_circular();

-- Validation
DO $$
BEGIN
    RAISE NOTICE 'Advanced constraints dropped successfully';
END $$;

COMMIT;
```

## Migration 009: Unique Constraints

### Up Migration

```sql
-- migrations/009_add_unique_constraints.up.sql
-- Migration: Add unique constraints
-- Version: 009
-- Description: Add unique constraints for data integrity

BEGIN;

-- Unique constraints for business logic

-- Project names should be unique (optional - depends on requirements)
-- ALTER TABLE projects ADD CONSTRAINT unique_project_name UNIQUE (name);

-- Task titles should be unique within a project
ALTER TABLE tasks ADD CONSTRAINT unique_task_title_per_project
    UNIQUE (project_id, title);

-- Step titles should be unique within a task
ALTER TABLE steps ADD CONSTRAINT unique_step_title_per_task
    UNIQUE (task_id, title);

-- Navigation chain constraints

-- Ensure no duplicate prev_id references within the same context
-- (A task/step can only be the "next" item for one other item)
CREATE UNIQUE INDEX unique_prev_id_tasks ON tasks(prev_id)
WHERE prev_id IS NOT NULL;

CREATE UNIQUE INDEX unique_prev_id_steps ON steps(prev_id)
WHERE prev_id IS NOT NULL;

-- Ensure no duplicate next_id references within the same context
-- (A task/step can only be the "previous" item for one other item)
CREATE UNIQUE INDEX unique_next_id_tasks ON tasks(next_id)
WHERE next_id IS NOT NULL;

CREATE UNIQUE INDEX unique_next_id_steps ON steps(next_id)
WHERE next_id IS NOT NULL;

-- Audit log constraints

-- Ensure audit log entries are unique for simultaneous operations
-- (Prevent duplicate audit entries for the same entity/action/timestamp)
CREATE UNIQUE INDEX unique_audit_log_entry ON audit_logs(
    entity_type, entity_id, action, timestamp, COALESCE(user_id, ''), COALESCE(ai_agent::TEXT, '')
);

-- Add comments
COMMENT ON CONSTRAINT unique_task_title_per_project ON tasks IS 'Ensures task titles are unique within each project';
COMMENT ON CONSTRAINT unique_step_title_per_task ON steps IS 'Ensures step titles are unique within each task';
COMMENT ON INDEX unique_prev_id_tasks IS 'Ensures navigation chain integrity - no duplicate prev references';
COMMENT ON INDEX unique_next_id_tasks IS 'Ensures navigation chain integrity - no duplicate next references';
COMMENT ON INDEX unique_audit_log_entry IS 'Prevents duplicate audit log entries';

-- Validation
DO $$
DECLARE
    constraint_count INTEGER;
    unique_index_count INTEGER;
BEGIN
    -- Count unique constraints
    SELECT COUNT(*) INTO constraint_count
    FROM information_schema.table_constraints
    WHERE constraint_type = 'UNIQUE'
    AND table_schema = 'public'
    AND constraint_name LIKE 'unique_%';

    -- Count unique indexes
    SELECT COUNT(*) INTO unique_index_count
    FROM pg_indexes
    WHERE schemaname = 'public'
    AND indexname LIKE 'unique_%';

    RAISE NOTICE 'Unique constraints created successfully:';
    RAISE NOTICE '  Constraints: %', constraint_count;
    RAISE NOTICE '  Unique indexes: %', unique_index_count;
END $$;

COMMIT;
```

### Down Migration

```sql
-- migrations/009_add_unique_constraints.down.sql
-- Rollback: Add unique constraints
-- Version: 009
-- Description: Drop unique constraints

BEGIN;

-- Drop unique indexes
DROP INDEX IF EXISTS unique_audit_log_entry;
DROP INDEX IF EXISTS unique_next_id_steps;
DROP INDEX IF EXISTS unique_next_id_tasks;
DROP INDEX IF EXISTS unique_prev_id_steps;
DROP INDEX IF EXISTS unique_prev_id_tasks;

-- Drop unique constraints
ALTER TABLE steps DROP CONSTRAINT IF EXISTS unique_step_title_per_task;
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS unique_task_title_per_project;
-- ALTER TABLE projects DROP CONSTRAINT IF EXISTS unique_project_name;

-- Validation
DO $$
BEGIN
    RAISE NOTICE 'Unique constraints dropped successfully';
END $$;

COMMIT;
```

---

*Next: [Function & Trigger Migrations](./07b3d-functions-triggers.md)*
