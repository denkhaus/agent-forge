# Progress Calculation Functions

## Overview

This document covers the implementation of automatic progress calculation functions and triggers for the MCP-Planner system. These functions maintain hierarchical progress consistency across projects, tasks, and steps.

## Migration 010: Progress Calculation Functions

### Up Migration

```sql
-- migrations/010_create_progress_functions.up.sql
-- Migration: Create progress calculation functions
-- Version: 010
-- Description: Implement hierarchical progress calculation and automatic updates

BEGIN;

-- Core progress calculation function for tasks
CREATE OR REPLACE FUNCTION calculate_task_progress(task_id_param UUID)
RETURNS FLOAT AS $$
DECLARE
    step_progress FLOAT := 0;
    subtask_progress FLOAT := 0;
    step_count INTEGER := 0;
    subtask_count INTEGER := 0;
    total_progress FLOAT := 0;
    total_items INTEGER := 0;
BEGIN
    -- Calculate progress from direct steps
    SELECT
        COALESCE(AVG(progress), 0),
        COUNT(*)
    INTO step_progress, step_count
    FROM steps
    WHERE task_id = task_id_param;

    -- Calculate progress from sub-tasks (recursive)
    SELECT
        COALESCE(AVG(calculate_task_progress(id)), 0),
        COUNT(*)
    INTO subtask_progress, subtask_count
    FROM tasks
    WHERE parent_task_id = task_id_param;

    -- Combine step and subtask progress
    total_items := step_count + subtask_count;

    IF total_items = 0 THEN
        RETURN 0;
    END IF;

    total_progress := (step_progress * step_count + subtask_progress * subtask_count) / total_items;

    RETURN LEAST(1.0, GREATEST(0.0, total_progress));
END;
$$ LANGUAGE plpgsql;

-- Core progress calculation function for projects
CREATE OR REPLACE FUNCTION calculate_project_progress(project_id_param UUID)
RETURNS FLOAT AS $$
DECLARE
    total_progress FLOAT := 0;
    task_count INTEGER := 0;
BEGIN
    -- Calculate average progress of root tasks
    SELECT
        COALESCE(AVG(calculate_task_progress(id)), 0),
        COUNT(*)
    INTO total_progress, task_count
    FROM tasks
    WHERE project_id = project_id_param
    AND parent_task_id IS NULL;

    IF task_count = 0 THEN
        RETURN 0;
    END IF;

    RETURN LEAST(1.0, GREATEST(0.0, total_progress));
END;
$$ LANGUAGE plpgsql;

-- Efficient batch progress calculation for multiple tasks
CREATE OR REPLACE FUNCTION calculate_multiple_task_progress(task_ids UUID[])
RETURNS TABLE(task_id UUID, progress FLOAT) AS $$
DECLARE
    task_id_param UUID;
BEGIN
    FOREACH task_id_param IN ARRAY task_ids
    LOOP
        RETURN QUERY SELECT task_id_param, calculate_task_progress(task_id_param);
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Update task progress and propagate to parent
CREATE OR REPLACE FUNCTION update_task_progress_cascade(task_id_param UUID)
RETURNS VOID AS $$
DECLARE
    new_progress FLOAT;
    parent_task_id_var UUID;
    project_id_var UUID;
BEGIN
    -- Calculate new progress for the task
    new_progress := calculate_task_progress(task_id_param);

    -- Update the task's progress
    UPDATE tasks
    SET progress = new_progress, updated_at = NOW()
    WHERE id = task_id_param
    RETURNING parent_task_id, project_id INTO parent_task_id_var, project_id_var;

    -- If this task has a parent, update parent's progress recursively
    IF parent_task_id_var IS NOT NULL THEN
        PERFORM update_task_progress_cascade(parent_task_id_var);
    ELSE
        -- This is a root task, update project progress
        PERFORM update_project_progress(project_id_var);
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Update project progress
CREATE OR REPLACE FUNCTION update_project_progress(project_id_param UUID)
RETURNS VOID AS $$
DECLARE
    new_progress FLOAT;
BEGIN
    -- Calculate new progress for the project
    new_progress := calculate_project_progress(project_id_param);

    -- Update the project's progress
    UPDATE projects
    SET progress = new_progress, updated_at = NOW()
    WHERE id = project_id_param;
END;
$$ LANGUAGE plpgsql;

-- Recalculate progress for entire project hierarchy
CREATE OR REPLACE FUNCTION recalculate_project_hierarchy_progress(project_id_param UUID)
RETURNS VOID AS $$
DECLARE
    task_record RECORD;
BEGIN
    -- Update all tasks in the project (bottom-up approach)
    -- First, update leaf tasks (tasks with no sub-tasks)
    FOR task_record IN
        WITH RECURSIVE task_hierarchy AS (
            -- Start with leaf tasks (no children)
            SELECT id, parent_task_id, 0 as level
            FROM tasks t1
            WHERE project_id = project_id_param
            AND NOT EXISTS (
                SELECT 1 FROM tasks t2 WHERE t2.parent_task_id = t1.id
            )

            UNION ALL

            -- Move up the hierarchy
            SELECT t.id, t.parent_task_id, th.level + 1
            FROM tasks t
            JOIN task_hierarchy th ON t.id = th.parent_task_id
        )
        SELECT DISTINCT id FROM task_hierarchy ORDER BY level
    LOOP
        PERFORM update_task_progress_cascade(task_record.id);
    END LOOP;

    -- Finally, update project progress
    PERFORM update_project_progress(project_id_param);
END;
$$ LANGUAGE plpgsql;

-- Get progress statistics for a project
CREATE OR REPLACE FUNCTION get_project_progress_stats(project_id_param UUID)
RETURNS TABLE(
    total_tasks INTEGER,
    completed_tasks INTEGER,
    total_steps INTEGER,
    completed_steps INTEGER,
    overall_progress FLOAT,
    task_progress_distribution JSONB
) AS $$
DECLARE
    stats_record RECORD;
BEGIN
    -- Calculate comprehensive progress statistics
    SELECT
        COUNT(DISTINCT t.id) as total_tasks,
        COUNT(DISTINCT CASE WHEN t.progress >= 1.0 THEN t.id END) as completed_tasks,
        COUNT(DISTINCT s.id) as total_steps,
        COUNT(DISTINCT CASE WHEN s.progress >= 1.0 THEN s.id END) as completed_steps,
        calculate_project_progress(project_id_param) as overall_progress,
        jsonb_build_object(
            'not_started', COUNT(DISTINCT CASE WHEN t.progress = 0 THEN t.id END),
            'in_progress', COUNT(DISTINCT CASE WHEN t.progress > 0 AND t.progress < 1.0 THEN t.id END),
            'completed', COUNT(DISTINCT CASE WHEN t.progress >= 1.0 THEN t.id END)
        ) as task_progress_distribution
    INTO stats_record
    FROM tasks t
    LEFT JOIN steps s ON s.task_id = t.id
    WHERE t.project_id = project_id_param;

    RETURN QUERY SELECT
        stats_record.total_tasks,
        stats_record.completed_tasks,
        stats_record.total_steps,
        stats_record.completed_steps,
        stats_record.overall_progress,
        stats_record.task_progress_distribution;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for automatic step progress updates
CREATE OR REPLACE FUNCTION trigger_update_step_progress()
RETURNS TRIGGER AS $$
BEGIN
    -- Only update if progress actually changed
    IF TG_OP = 'UPDATE' AND OLD.progress = NEW.progress THEN
        RETURN NEW;
    END IF;

    -- Update parent task progress
    PERFORM update_task_progress_cascade(NEW.task_id);

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for automatic task progress updates
CREATE OR REPLACE FUNCTION trigger_update_task_progress()
RETURNS TRIGGER AS $$
DECLARE
    affected_project_id UUID;
    affected_parent_task_id UUID;
BEGIN
    -- Handle different trigger operations
    CASE TG_OP
        WHEN 'INSERT' THEN
            affected_project_id := NEW.project_id;
            affected_parent_task_id := NEW.parent_task_id;
        WHEN 'UPDATE' THEN
            -- Only update if progress changed or hierarchy changed
            IF OLD.progress = NEW.progress AND
               OLD.parent_task_id = NEW.parent_task_id AND
               OLD.project_id = NEW.project_id THEN
                RETURN NEW;
            END IF;
            affected_project_id := NEW.project_id;
            affected_parent_task_id := NEW.parent_task_id;
        WHEN 'DELETE' THEN
            affected_project_id := OLD.project_id;
            affected_parent_task_id := OLD.parent_task_id;
    END CASE;

    -- Update parent task if exists
    IF affected_parent_task_id IS NOT NULL THEN
        PERFORM update_task_progress_cascade(affected_parent_task_id);
    ELSE
        -- Update project progress for root tasks
        PERFORM update_project_progress(affected_project_id);
    END IF;

    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

-- Trigger function for step status changes affecting progress
CREATE OR REPLACE FUNCTION trigger_step_status_progress_sync()
RETURNS TRIGGER AS $$
BEGIN
    -- Automatically set progress based on status
    CASE NEW.status
        WHEN 'AGREED' THEN
            NEW.progress := 1.0;
        WHEN 'PENDING' THEN
            IF NEW.progress > 0 THEN
                NEW.progress := 0.0;
            END IF;
        ELSE
            -- For other statuses, keep current progress
            NULL;
    END CASE;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for automatic progress updates

-- Step progress triggers
CREATE TRIGGER trigger_step_progress_update
    AFTER INSERT OR UPDATE OF progress ON steps
    FOR EACH ROW
    EXECUTE FUNCTION trigger_update_step_progress();

CREATE TRIGGER trigger_step_status_progress_update
    BEFORE INSERT OR UPDATE OF status ON steps
    FOR EACH ROW
    EXECUTE FUNCTION trigger_step_status_progress_sync();

-- Task progress triggers
CREATE TRIGGER trigger_task_hierarchy_progress_update
    AFTER INSERT OR UPDATE OF parent_task_id, project_id OR DELETE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION trigger_update_task_progress();

-- Performance optimization: batch progress updates
CREATE OR REPLACE FUNCTION batch_update_project_progress(project_ids UUID[])
RETURNS TABLE(project_id UUID, old_progress FLOAT, new_progress FLOAT) AS $$
DECLARE
    project_id_param UUID;
    old_prog FLOAT;
    new_prog FLOAT;
BEGIN
    FOREACH project_id_param IN ARRAY project_ids
    LOOP
        -- Get current progress
        SELECT progress INTO old_prog FROM projects WHERE id = project_id_param;

        -- Calculate new progress
        new_prog := calculate_project_progress(project_id_param);

        -- Update if changed
        IF old_prog != new_prog THEN
            UPDATE projects
            SET progress = new_prog, updated_at = NOW()
            WHERE id = project_id_param;
        END IF;

        RETURN QUERY SELECT project_id_param, old_prog, new_prog;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Progress validation function
CREATE OR REPLACE FUNCTION validate_progress_consistency(project_id_param UUID)
RETURNS TABLE(
    entity_type TEXT,
    entity_id UUID,
    calculated_progress FLOAT,
    stored_progress FLOAT,
    difference FLOAT
) AS $$
BEGIN
    -- Check project progress consistency
    RETURN QUERY
    SELECT
        'project'::TEXT,
        p.id,
        calculate_project_progress(p.id),
        p.progress,
        ABS(calculate_project_progress(p.id) - p.progress)
    FROM projects p
    WHERE p.id = project_id_param
    AND ABS(calculate_project_progress(p.id) - p.progress) > 0.001;

    -- Check task progress consistency
    RETURN QUERY
    SELECT
        'task'::TEXT,
        t.id,
        calculate_task_progress(t.id),
        t.progress,
        ABS(calculate_task_progress(t.id) - t.progress)
    FROM tasks t
    WHERE t.project_id = project_id_param
    AND ABS(calculate_task_progress(t.id) - t.progress) > 0.001;
END;
$$ LANGUAGE plpgsql;

-- Add function comments
COMMENT ON FUNCTION calculate_task_progress(UUID) IS 'Calculates task progress from steps and sub-tasks';
COMMENT ON FUNCTION calculate_project_progress(UUID) IS 'Calculates project progress from root tasks';
COMMENT ON FUNCTION update_task_progress_cascade(UUID) IS 'Updates task progress and cascades to parents';
COMMENT ON FUNCTION recalculate_project_hierarchy_progress(UUID) IS 'Recalculates progress for entire project hierarchy';
COMMENT ON FUNCTION get_project_progress_stats(UUID) IS 'Returns comprehensive progress statistics for a project';
COMMENT ON FUNCTION validate_progress_consistency(UUID) IS 'Validates progress consistency across hierarchy';

-- Validation
DO $$
DECLARE
    function_count INTEGER;
    trigger_count INTEGER;
BEGIN
    -- Count progress functions
    SELECT COUNT(*) INTO function_count
    FROM pg_proc
    WHERE proname LIKE '%progress%'
    AND pronamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'public');

    IF function_count < 10 THEN
        RAISE EXCEPTION 'Expected at least 10 progress functions, found %', function_count;
    END IF;

    -- Count progress triggers
    SELECT COUNT(*) INTO trigger_count
    FROM pg_trigger
    WHERE tgname LIKE '%progress%';

    IF trigger_count < 3 THEN
        RAISE EXCEPTION 'Expected at least 3 progress triggers, found %', trigger_count;
    END IF;

    RAISE NOTICE 'Progress functions created successfully:';
    RAISE NOTICE '  Functions: %', function_count;
    RAISE NOTICE '  Triggers: %', trigger_count;
END $$;

COMMIT;
```

### Down Migration

```sql
-- migrations/010_create_progress_functions.down.sql
-- Rollback: Create progress calculation functions
-- Version: 010
-- Description: Drop progress calculation functions and triggers

BEGIN;

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_task_hierarchy_progress_update ON tasks;
DROP TRIGGER IF EXISTS trigger_step_status_progress_update ON steps;
DROP TRIGGER IF EXISTS trigger_step_progress_update ON steps;

-- Drop functions
DROP FUNCTION IF EXISTS validate_progress_consistency(UUID);
DROP FUNCTION IF EXISTS batch_update_project_progress(UUID[]);
DROP FUNCTION IF EXISTS trigger_step_status_progress_sync();
DROP FUNCTION IF EXISTS trigger_update_task_progress();
DROP FUNCTION IF EXISTS trigger_update_step_progress();
DROP FUNCTION IF EXISTS get_project_progress_stats(UUID);
DROP FUNCTION IF EXISTS recalculate_project_hierarchy_progress(UUID);
DROP FUNCTION IF EXISTS update_project_progress(UUID);
DROP FUNCTION IF EXISTS update_task_progress_cascade(UUID);
DROP FUNCTION IF EXISTS calculate_multiple_task_progress(UUID[]);
DROP FUNCTION IF EXISTS calculate_project_progress(UUID);
DROP FUNCTION IF EXISTS calculate_task_progress(UUID);

-- Validation
DO $$
BEGIN
    RAISE NOTICE 'Progress functions dropped successfully';
END $$;

COMMIT;
```

## Usage Examples

### Manual Progress Calculation

```sql
-- Calculate progress for a specific task
SELECT calculate_task_progress('123e4567-e89b-12d3-a456-426614174000');

-- Calculate progress for a project
SELECT calculate_project_progress('987fcdeb-51a2-43d1-9f4e-123456789abc');

-- Get comprehensive project statistics
SELECT * FROM get_project_progress_stats('987fcdeb-51a2-43d1-9f4e-123456789abc');
```

### Progress Validation

```sql
-- Check for progress inconsistencies
SELECT * FROM validate_progress_consistency('987fcdeb-51a2-43d1-9f4e-123456789abc');

-- Recalculate entire project hierarchy
SELECT recalculate_project_hierarchy_progress('987fcdeb-51a2-43d1-9f4e-123456789abc');
```

### Batch Operations

```sql
-- Update multiple projects efficiently
SELECT * FROM batch_update_project_progress(ARRAY[
    '987fcdeb-51a2-43d1-9f4e-123456789abc',
    '123e4567-e89b-12d3-a456-426614174000'
]);

-- Calculate progress for multiple tasks
SELECT * FROM calculate_multiple_task_progress(ARRAY[
    'task-id-1',
    'task-id-2',
    'task-id-3'
]);
```

## Performance Considerations

### Optimization Strategies

1. **Recursive Function Optimization**
   - Use iterative approaches for deep hierarchies
   - Implement caching for frequently accessed calculations
   - Batch operations to reduce function call overhead

2. **Trigger Optimization**
   - Conditional trigger execution
   - Batch updates during bulk operations
   - Async processing for heavy calculations

3. **Index Usage**
   - Ensure proper indexes for hierarchy queries
   - Optimize for parent-child lookups
   - Index progress columns for filtering

### Monitoring Queries

```sql
-- Monitor progress calculation performance
SELECT
    schemaname,
    funcname,
    calls,
    total_time,
    mean_time,
    stddev_time
FROM pg_stat_user_functions
WHERE funcname LIKE '%progress%'
ORDER BY total_time DESC;

-- Check for long-running progress calculations
SELECT
    pid,
    now() - pg_stat_activity.query_start AS duration,
    query
FROM pg_stat_activity
WHERE query LIKE '%progress%'
AND state = 'active'
ORDER BY duration DESC;
```

## Testing Functions

### Progress Calculation Tests

```sql
-- Test basic progress calculation
CREATE OR REPLACE FUNCTION test_basic_progress_calculation()
RETURNS BOOLEAN AS $$
DECLARE
    test_project_id UUID;
    test_task_id UUID;
    test_step_id UUID;
    calculated_progress FLOAT;
BEGIN
    -- Create test data
    INSERT INTO projects (id, name, description)
    VALUES (gen_random_uuid(), 'Test Project', 'Test Description')
    RETURNING id INTO test_project_id;

    INSERT INTO tasks (id, project_id, title, objective)
    VALUES (gen_random_uuid(), test_project_id, 'Test Task', 'Test Objective')
    RETURNING id INTO test_task_id;

    INSERT INTO steps (id, task_id, title, progress)
    VALUES (gen_random_uuid(), test_task_id, 'Test Step', 0.5)
    RETURNING id INTO test_step_id;

    -- Test calculation
    calculated_progress := calculate_task_progress(test_task_id);

    -- Cleanup
    DELETE FROM projects WHERE id = test_project_id;

    -- Verify result
    RETURN calculated_progress = 0.5;
END;
$$ LANGUAGE plpgsql;

-- Test hierarchical progress calculation
CREATE OR REPLACE FUNCTION test_hierarchical_progress()
RETURNS BOOLEAN AS $$
DECLARE
    test_project_id UUID;
    parent_task_id UUID;
    child_task_id UUID;
    step1_id UUID;
    step2_id UUID;
    calculated_progress FLOAT;
BEGIN
    -- Create test hierarchy
    INSERT INTO projects (id, name, description)
    VALUES (gen_random_uuid(), 'Test Project', 'Test Description')
    RETURNING id INTO test_project_id;

    INSERT INTO tasks (id, project_id, title, objective)
    VALUES (gen_random_uuid(), test_project_id, 'Parent Task', 'Parent Objective')
    RETURNING id INTO parent_task_id;

    INSERT INTO tasks (id, project_id, parent_task_id, title, objective)
    VALUES (gen_random_uuid(), test_project_id, parent_task_id, 'Child Task', 'Child Objective')
    RETURNING id INTO child_task_id;

    -- Add steps with different progress
    INSERT INTO steps (id, task_id, title, progress)
    VALUES (gen_random_uuid(), child_task_id, 'Step 1', 1.0)
    RETURNING id INTO step1_id;

    INSERT INTO steps (id, task_id, title, progress)
    VALUES (gen_random_uuid(), child_task_id, 'Step 2', 0.0)
    RETURNING id INTO step2_id;

    -- Test calculation (should be 0.5)
    calculated_progress := calculate_task_progress(parent_task_id);

    -- Cleanup
    DELETE FROM projects WHERE id = test_project_id;

    -- Verify result
    RETURN calculated_progress = 0.5;
END;
$$ LANGUAGE plpgsql;
```

---

*Next: [Audit Logging Triggers](./07b3d2-audit-triggers.md)*
