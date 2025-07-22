# Audit Logging Triggers

## Overview

This document covers the implementation of comprehensive audit logging triggers for the MCP-Planner system. These triggers automatically track all changes to entities, providing a complete audit trail for compliance, debugging, and analytics.

## Migration 011: Audit Logging System

### Up Migration

```sql
-- migrations/011_create_audit_triggers.up.sql
-- Migration: Create audit logging triggers
-- Version: 011
-- Description: Implement comprehensive audit logging for all entity changes

BEGIN;

-- Core audit logging function
CREATE OR REPLACE FUNCTION audit_trigger_function()
RETURNS TRIGGER AS $$
DECLARE
    entity_type_val entity_type;
    action_val audit_action;
    old_data JSONB;
    new_data JSONB;
    changed_fields JSONB;
    audit_metadata JSONB;
BEGIN
    -- Determine entity type based on table name
    CASE TG_TABLE_NAME
        WHEN 'projects' THEN entity_type_val := 'PROJECT';
        WHEN 'tasks' THEN entity_type_val := 'TASK';
        WHEN 'steps' THEN entity_type_val := 'STEP';
        WHEN 'disputes' THEN entity_type_val := 'DISPUTE';
        ELSE
            RAISE EXCEPTION 'Unsupported table for audit logging: %', TG_TABLE_NAME;
    END CASE;

    -- Determine action type
    CASE TG_OP
        WHEN 'INSERT' THEN action_val := 'CREATED';
        WHEN 'UPDATE' THEN action_val := 'UPDATED';
        WHEN 'DELETE' THEN action_val := 'DELETED';
        ELSE
            RAISE EXCEPTION 'Unsupported operation for audit logging: %', TG_OP;
    END CASE;

    -- Prepare data for logging
    CASE TG_OP
        WHEN 'INSERT' THEN
            old_data := NULL;
            new_data := to_jsonb(NEW);
            changed_fields := new_data;
        WHEN 'UPDATE' THEN
            old_data := to_jsonb(OLD);
            new_data := to_jsonb(NEW);
            changed_fields := get_changed_fields(old_data, new_data);
        WHEN 'DELETE' THEN
            old_data := to_jsonb(OLD);
            new_data := NULL;
            changed_fields := old_data;
    END CASE;

    -- Build metadata
    audit_metadata := jsonb_build_object(
        'table_name', TG_TABLE_NAME,
        'operation', TG_OP,
        'trigger_name', TG_NAME,
        'session_user', session_user,
        'current_user', current_user,
        'client_addr', inet_client_addr(),
        'client_port', inet_client_port(),
        'application_name', current_setting('application_name', true),
        'transaction_id', txid_current(),
        'changed_fields_count', jsonb_array_length(jsonb_object_keys(changed_fields))
    );

    -- Insert audit record
    INSERT INTO audit_logs (
        entity_type,
        entity_id,
        action,
        old_values,
        new_values,
        metadata,
        timestamp
    ) VALUES (
        entity_type_val,
        COALESCE(NEW.id, OLD.id),
        action_val,
        old_data,
        new_data,
        audit_metadata,
        NOW()
    );

    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

-- Function to identify changed fields between old and new JSON
CREATE OR REPLACE FUNCTION get_changed_fields(old_data JSONB, new_data JSONB)
RETURNS JSONB AS $$
DECLARE
    changed_fields JSONB := '{}';
    key TEXT;
    old_value JSONB;
    new_value JSONB;
BEGIN
    -- Compare each field in the new data
    FOR key IN SELECT jsonb_object_keys(new_data)
    LOOP
        old_value := old_data -> key;
        new_value := new_data -> key;

        -- Check if value changed (handle NULL comparisons)
        IF (old_value IS NULL AND new_value IS NOT NULL) OR
           (old_value IS NOT NULL AND new_value IS NULL) OR
           (old_value != new_value) THEN

            changed_fields := changed_fields || jsonb_build_object(
                key, jsonb_build_object(
                    'old', old_value,
                    'new', new_value
                )
            );
        END IF;
    END LOOP;

    RETURN changed_fields;
END;
$$ LANGUAGE plpgsql;

-- Specialized audit function for step status changes
CREATE OR REPLACE FUNCTION audit_step_status_change()
RETURNS TRIGGER AS $$
DECLARE
    status_metadata JSONB;
BEGIN
    -- Only log if status actually changed
    IF TG_OP = 'UPDATE' AND OLD.status = NEW.status THEN
        RETURN NEW;
    END IF;

    -- Build status-specific metadata
    status_metadata := jsonb_build_object(
        'status_transition', jsonb_build_object(
            'from', OLD.status,
            'to', NEW.status
        ),
        'iteration_count', NEW.iteration_count,
        'server_ready', NEW.server_ready,
        'client_approved', NEW.client_approved,
        'complexity_score', NEW.complexity_score
    );

    -- Insert specialized audit record
    INSERT INTO audit_logs (
        entity_type,
        entity_id,
        action,
        old_values,
        new_values,
        metadata,
        timestamp
    ) VALUES (
        'STEP',
        NEW.id,
        'UPDATED',
        jsonb_build_object('status', OLD.status),
        jsonb_build_object('status', NEW.status),
        status_metadata,
        NOW()
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Audit function for step promotion to task
CREATE OR REPLACE FUNCTION audit_step_promotion()
RETURNS TRIGGER AS $$
BEGIN
    -- This function is called manually during step promotion
    -- Log the promotion event
    INSERT INTO audit_logs (
        entity_type,
        entity_id,
        action,
        old_values,
        new_values,
        metadata,
        timestamp
    ) VALUES (
        'STEP',
        OLD.id,
        'PROMOTED',
        to_jsonb(OLD),
        jsonb_build_object('promoted_to_task_id', NEW.id),
        jsonb_build_object(
            'promotion_reason', 'complexity_threshold_exceeded',
            'original_complexity_score', OLD.complexity_score,
            'new_task_id', NEW.id
        ),
        NOW()
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Audit function for dispute creation and resolution
CREATE OR REPLACE FUNCTION audit_dispute_lifecycle()
RETURNS TRIGGER AS $$
DECLARE
    dispute_metadata JSONB;
    action_val audit_action;
BEGIN
    -- Determine specific action based on status changes
    IF TG_OP = 'INSERT' THEN
        action_val := 'DISPUTED';
        dispute_metadata := jsonb_build_object(
            'step_id', NEW.step_id,
            'project_id', NEW.project_id,
            'dispute_created', true
        );
    ELSIF TG_OP = 'UPDATE' AND OLD.status != NEW.status THEN
        IF NEW.status = 'RESOLVED' THEN
            action_val := 'RESOLVED';
            dispute_metadata := jsonb_build_object(
                'resolution_type', NEW.user_resolution,
                'resolution_time', NEW.resolved_at,
                'dispute_duration', EXTRACT(EPOCH FROM (NEW.resolved_at - NEW.created_at))
            );
        ELSE
            action_val := 'UPDATED';
            dispute_metadata := jsonb_build_object(
                'status_change', jsonb_build_object(
                    'from', OLD.status,
                    'to', NEW.status
                )
            );
        END IF;
    ELSE
        -- Regular update, use standard audit
        RETURN NEW;
    END IF;

    -- Insert specialized dispute audit record
    INSERT INTO audit_logs (
        entity_type,
        entity_id,
        action,
        old_values,
        new_values,
        metadata,
        timestamp
    ) VALUES (
        'DISPUTE',
        COALESCE(NEW.id, OLD.id),
        action_val,
        CASE WHEN TG_OP = 'INSERT' THEN NULL ELSE to_jsonb(OLD) END,
        to_jsonb(NEW),
        dispute_metadata,
        NOW()
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to log AI agent activities
CREATE OR REPLACE FUNCTION log_ai_agent_activity(
    agent_type ai_agent,
    entity_type_param entity_type,
    entity_id_param UUID,
    action_param audit_action,
    activity_data JSONB DEFAULT NULL
)
RETURNS VOID AS $$
BEGIN
    INSERT INTO audit_logs (
        entity_type,
        entity_id,
        action,
        ai_agent,
        metadata,
        timestamp
    ) VALUES (
        entity_type_param,
        entity_id_param,
        action_param,
        agent_type,
        COALESCE(activity_data, '{}'),
        NOW()
    );
END;
$$ LANGUAGE plpgsql;

-- Function to log user activities
CREATE OR REPLACE FUNCTION log_user_activity(
    user_id_param VARCHAR(255),
    entity_type_param entity_type,
    entity_id_param UUID,
    action_param audit_action,
    activity_data JSONB DEFAULT NULL
)
RETURNS VOID AS $$
BEGIN
    INSERT INTO audit_logs (
        entity_type,
        entity_id,
        action,
        user_id,
        metadata,
        timestamp
    ) VALUES (
        entity_type_param,
        entity_id_param,
        action_param,
        user_id_param,
        COALESCE(activity_data, '{}'),
        NOW()
    );
END;
$$ LANGUAGE plpgsql;

-- Create audit triggers for all tables

-- Projects audit trigger
CREATE TRIGGER audit_projects_trigger
    AFTER INSERT OR UPDATE OR DELETE ON projects
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

-- Tasks audit trigger
CREATE TRIGGER audit_tasks_trigger
    AFTER INSERT OR UPDATE OR DELETE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

-- Steps audit trigger (general)
CREATE TRIGGER audit_steps_trigger
    AFTER INSERT OR UPDATE OR DELETE ON steps
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

-- Steps status change trigger (specialized)
CREATE TRIGGER audit_steps_status_trigger
    AFTER UPDATE OF status ON steps
    FOR EACH ROW
    EXECUTE FUNCTION audit_step_status_change();

-- Disputes audit trigger (specialized)
CREATE TRIGGER audit_disputes_lifecycle_trigger
    AFTER INSERT OR UPDATE ON disputes
    FOR EACH ROW
    EXECUTE FUNCTION audit_dispute_lifecycle();

-- Function to get audit trail for an entity
CREATE OR REPLACE FUNCTION get_entity_audit_trail(
    entity_type_param entity_type,
    entity_id_param UUID,
    limit_param INTEGER DEFAULT 100
)
RETURNS TABLE(
    id UUID,
    action audit_action,
    old_values JSONB,
    new_values JSONB,
    user_id VARCHAR(255),
    ai_agent ai_agent,
    metadata JSONB,
    timestamp TIMESTAMPTZ
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        al.id,
        al.action,
        al.old_values,
        al.new_values,
        al.user_id,
        al.ai_agent,
        al.metadata,
        al.timestamp
    FROM audit_logs al
    WHERE al.entity_type = entity_type_param
    AND al.entity_id = entity_id_param
    ORDER BY al.timestamp DESC
    LIMIT limit_param;
END;
$$ LANGUAGE plpgsql;

-- Function to get recent activity across all entities
CREATE OR REPLACE FUNCTION get_recent_activity(
    project_id_param UUID DEFAULT NULL,
    limit_param INTEGER DEFAULT 50,
    since_timestamp TIMESTAMPTZ DEFAULT NOW() - INTERVAL '24 hours'
)
RETURNS TABLE(
    id UUID,
    entity_type entity_type,
    entity_id UUID,
    action audit_action,
    user_id VARCHAR(255),
    ai_agent ai_agent,
    timestamp TIMESTAMPTZ,
    summary TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        al.id,
        al.entity_type,
        al.entity_id,
        al.action,
        al.user_id,
        al.ai_agent,
        al.timestamp,
        CASE
            WHEN al.entity_type = 'PROJECT' THEN
                'Project ' || (al.new_values->>'name') || ' was ' || LOWER(al.action::TEXT)
            WHEN al.entity_type = 'TASK' THEN
                'Task ' || (al.new_values->>'title') || ' was ' || LOWER(al.action::TEXT)
            WHEN al.entity_type = 'STEP' THEN
                'Step ' || (al.new_values->>'title') || ' was ' || LOWER(al.action::TEXT)
            WHEN al.entity_type = 'DISPUTE' THEN
                'Dispute was ' || LOWER(al.action::TEXT)
            ELSE
                'Entity was ' || LOWER(al.action::TEXT)
        END as summary
    FROM audit_logs al
    LEFT JOIN tasks t ON al.entity_type = 'TASK' AND al.entity_id = t.id
    LEFT JOIN steps s ON al.entity_type = 'STEP' AND al.entity_id = s.id
    LEFT JOIN tasks st ON s.task_id = st.id
    WHERE al.timestamp >= since_timestamp
    AND (project_id_param IS NULL OR
         al.entity_type = 'PROJECT' AND al.entity_id = project_id_param OR
         al.entity_type = 'TASK' AND t.project_id = project_id_param OR
         al.entity_type = 'STEP' AND st.project_id = project_id_param)
    ORDER BY al.timestamp DESC
    LIMIT limit_param;
END;
$$ LANGUAGE plpgsql;

-- Function to get audit statistics
CREATE OR REPLACE FUNCTION get_audit_statistics(
    since_timestamp TIMESTAMPTZ DEFAULT NOW() - INTERVAL '30 days'
)
RETURNS TABLE(
    total_events BIGINT,
    events_by_type JSONB,
    events_by_action JSONB,
    events_by_agent JSONB,
    top_active_entities JSONB
) AS $$
DECLARE
    stats_record RECORD;
BEGIN
    SELECT
        COUNT(*) as total_events,
        jsonb_object_agg(entity_type, type_count) as events_by_type,
        jsonb_object_agg(action, action_count) as events_by_action,
        jsonb_object_agg(
            COALESCE(ai_agent::TEXT, user_id, 'system'),
            agent_count
        ) as events_by_agent,
        jsonb_agg(
            jsonb_build_object(
                'entity_type', entity_type,
                'entity_id', entity_id,
                'event_count', entity_event_count
            ) ORDER BY entity_event_count DESC
        ) FILTER (WHERE entity_rank <= 10) as top_active_entities
    INTO stats_record
    FROM (
        SELECT
            entity_type,
            action,
            ai_agent,
            user_id,
            entity_id,
            COUNT(*) OVER (PARTITION BY entity_type) as type_count,
            COUNT(*) OVER (PARTITION BY action) as action_count,
            COUNT(*) OVER (PARTITION BY COALESCE(ai_agent::TEXT, user_id, 'system')) as agent_count,
            COUNT(*) OVER (PARTITION BY entity_type, entity_id) as entity_event_count,
            ROW_NUMBER() OVER (PARTITION BY entity_type, entity_id ORDER BY timestamp DESC) as entity_rank
        FROM audit_logs
        WHERE timestamp >= since_timestamp
    ) subq
    WHERE entity_rank = 1;

    RETURN QUERY SELECT
        stats_record.total_events,
        stats_record.events_by_type,
        stats_record.events_by_action,
        stats_record.events_by_agent,
        stats_record.top_active_entities;
END;
$$ LANGUAGE plpgsql;

-- Function to cleanup old audit logs
CREATE OR REPLACE FUNCTION cleanup_old_audit_logs(
    retention_days INTEGER DEFAULT 365
)
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM audit_logs
    WHERE timestamp < NOW() - (retention_days || ' days')::INTERVAL;

    GET DIAGNOSTICS deleted_count = ROW_COUNT;

    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Add function comments
COMMENT ON FUNCTION audit_trigger_function() IS 'Main audit trigger function for all entity changes';
COMMENT ON FUNCTION get_changed_fields(JSONB, JSONB) IS 'Identifies changed fields between old and new JSON data';
COMMENT ON FUNCTION audit_step_status_change() IS 'Specialized audit logging for step status changes';
COMMENT ON FUNCTION log_ai_agent_activity(ai_agent, entity_type, UUID, audit_action, JSONB) IS 'Logs AI agent activities';
COMMENT ON FUNCTION log_user_activity(VARCHAR, entity_type, UUID, audit_action, JSONB) IS 'Logs user activities';
COMMENT ON FUNCTION get_entity_audit_trail(entity_type, UUID, INTEGER) IS 'Retrieves audit trail for a specific entity';
COMMENT ON FUNCTION get_recent_activity(UUID, INTEGER, TIMESTAMPTZ) IS 'Retrieves recent activity across entities';
COMMENT ON FUNCTION cleanup_old_audit_logs(INTEGER) IS 'Removes audit logs older than specified retention period';

-- Validation
DO $$
DECLARE
    function_count INTEGER;
    trigger_count INTEGER;
BEGIN
    -- Count audit functions
    SELECT COUNT(*) INTO function_count
    FROM pg_proc
    WHERE (proname LIKE 'audit_%' OR proname LIKE 'log_%' OR proname LIKE 'get_%audit%' OR proname LIKE 'cleanup_%')
    AND pronamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'public');

    IF function_count < 8 THEN
        RAISE EXCEPTION 'Expected at least 8 audit functions, found %', function_count;
    END IF;

    -- Count audit triggers
    SELECT COUNT(*) INTO trigger_count
    FROM pg_trigger
    WHERE tgname LIKE 'audit_%';

    IF trigger_count < 5 THEN
        RAISE EXCEPTION 'Expected at least 5 audit triggers, found %', trigger_count;
    END IF;

    RAISE NOTICE 'Audit logging system created successfully:';
    RAISE NOTICE '  Functions: %', function_count;
    RAISE NOTICE '  Triggers: %', trigger_count;
END $$;

COMMIT;
```

### Down Migration

```sql
-- migrations/011_create_audit_triggers.down.sql
-- Rollback: Create audit logging triggers
-- Version: 011
-- Description: Drop audit logging system

BEGIN;

-- Drop triggers
DROP TRIGGER IF EXISTS audit_disputes_lifecycle_trigger ON disputes;
DROP TRIGGER IF EXISTS audit_steps_status_trigger ON steps;
DROP TRIGGER IF EXISTS audit_steps_trigger ON steps;
DROP TRIGGER IF EXISTS audit_tasks_trigger ON tasks;
DROP TRIGGER IF EXISTS audit_projects_trigger ON projects;

-- Drop functions
DROP FUNCTION IF EXISTS cleanup_old_audit_logs(INTEGER);
DROP FUNCTION IF EXISTS get_audit_statistics(TIMESTAMPTZ);
DROP FUNCTION IF EXISTS get_recent_activity(UUID, INTEGER, TIMESTAMPTZ);
DROP FUNCTION IF EXISTS get_entity_audit_trail(entity_type, UUID, INTEGER);
DROP FUNCTION IF EXISTS log_user_activity(VARCHAR, entity_type, UUID, audit_action, JSONB);
DROP FUNCTION IF EXISTS log_ai_agent_activity(ai_agent, entity_type, UUID, audit_action, JSONB);
DROP FUNCTION IF EXISTS audit_dispute_lifecycle();
DROP FUNCTION IF EXISTS audit_step_promotion();
DROP FUNCTION IF EXISTS audit_step_status_change();
DROP FUNCTION IF EXISTS get_changed_fields(JSONB, JSONB);
DROP FUNCTION IF EXISTS audit_trigger_function();

-- Validation
DO $$
BEGIN
    RAISE NOTICE 'Audit logging system dropped successfully';
END $$;

COMMIT;
```

## Usage Examples

### Manual Audit Logging

```sql
-- Log AI agent activity
SELECT log_ai_agent_activity(
    'SERVER_AI',
    'STEP',
    '123e4567-e89b-12d3-a456-426614174000',
    'UPDATED',
    '{"content_generated": true, "complexity_assessed": "medium"}'
);

-- Log user activity
SELECT log_user_activity(
    'user123',
    'DISPUTE',
    '987fcdeb-51a2-43d1-9f4e-123456789abc',
    'RESOLVED',
    '{"resolution_type": "custom", "resolution_time_minutes": 15}'
);
```

### Audit Trail Queries

```sql
-- Get audit trail for a specific step
SELECT * FROM get_entity_audit_trail('STEP', '123e4567-e89b-12d3-a456-426614174000');

-- Get recent activity for a project
SELECT * FROM get_recent_activity('987fcdeb-51a2-43d1-9f4e-123456789abc', 20);

-- Get audit statistics
SELECT * FROM get_audit_statistics();
```

### Maintenance Operations

```sql
-- Cleanup old audit logs (keep last 180 days)
SELECT cleanup_old_audit_logs(180);

-- Check audit log size
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size,
    pg_stat_get_tuples_inserted(c.oid) as inserts,
    pg_stat_get_tuples_updated(c.oid) as updates,
    pg_stat_get_tuples_deleted(c.oid) as deletes
FROM pg_tables pt
JOIN pg_class c ON c.relname = pt.tablename
WHERE tablename = 'audit_logs';
```

## Performance Monitoring

### Audit Performance Queries

```sql
-- Monitor audit trigger performance
SELECT
    schemaname,
    tablename,
    n_tup_ins + n_tup_upd + n_tup_del as total_changes,
    n_tup_ins as inserts,
    n_tup_upd as updates,
    n_tup_del as deletes
FROM pg_stat_user_tables
WHERE tablename IN ('projects', 'tasks', 'steps', 'disputes')
ORDER BY total_changes DESC;

-- Check audit log growth rate
SELECT
    DATE_TRUNC('day', timestamp) as day,
    COUNT(*) as audit_events,
    COUNT(DISTINCT entity_id) as unique_entities,
    COUNT(DISTINCT COALESCE(user_id, ai_agent::TEXT)) as unique_actors
FROM audit_logs
WHERE timestamp > NOW() - INTERVAL '7 days'
GROUP BY DATE_TRUNC('day', timestamp)
ORDER BY day;
```

---

*Next: [Navigation Helper Functions](./07b3d3-navigation-functions.md)*
