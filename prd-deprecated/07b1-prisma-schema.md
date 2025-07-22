# Prisma Schema & Models

## Overview

This document covers the complete Prisma schema definition for the MCP-Planner system, including all models, relationships, and database configuration.

## Complete Prisma Schema

### Schema Configuration

```prisma
// schema.prisma
generator client {
  provider = "prisma-client-go"
  output   = "../internal/database/generated"
}

datasource db {
  provider = "postgresql"
  url      = env("DATABASE_URL")
}

// Enable PostgreSQL extensions
generator dbml {
  provider = "prisma-dbml-generator"
  output   = "../docs/database"
}
```

### Core Models

#### Project Model

```prisma
model Project {
  id                    String    @id @default(uuid()) @db.Uuid
  name                  String    @db.VarChar(255)
  description           String    @db.Text
  progress              Float     @default(0) @db.DoublePrecision
  complexityThreshold   Float     @default(0.7) @db.DoublePrecision
  maxIterations         Int       @default(3) @db.Integer
  createdAt            DateTime  @default(now()) @db.Timestamptz
  updatedAt            DateTime  @updatedAt @db.Timestamptz

  // Relationships
  tasks                Task[]
  disputes             Dispute[]
  auditLogs           AuditLog[] @relation("ProjectAuditLogs")

  // Constraints
  @@check([progress >= 0 && progress <= 1], name: "valid_progress")
  @@check([complexityThreshold >= 0 && complexityThreshold <= 1], name: "valid_complexity_threshold")
  @@check([maxIterations > 0], name: "positive_max_iterations")

  // Indexes
  @@index([createdAt])
  @@index([progress])
  @@index([complexityThreshold])

  @@map("projects")
}
```

#### Task Model

```prisma
model Task {
  id           String    @id @default(uuid()) @db.Uuid
  title        String    @db.VarChar(255)
  objective    String    @db.Text
  progress     Float     @default(0) @db.DoublePrecision
  projectId    String    @db.Uuid
  parentTaskId String?   @db.Uuid
  prevId       String?   @db.VarChar(50)  // Polymorphic: "task://uuid" or "step://uuid"
  nextId       String?   @db.VarChar(50)  // Polymorphic: "task://uuid" or "step://uuid"
  createdAt    DateTime  @default(now()) @db.Timestamptz
  updatedAt    DateTime  @updatedAt @db.Timestamptz

  // Relationships
  project      Project   @relation(fields: [projectId], references: [id], onDelete: Cascade)
  parentTask   Task?     @relation("TaskHierarchy", fields: [parentTaskId], references: [id], onDelete: Cascade)
  subTasks     Task[]    @relation("TaskHierarchy")
  steps        Step[]
  auditLogs    AuditLog[] @relation("TaskAuditLogs")

  // Constraints
  @@check([progress >= 0 && progress <= 1], name: "valid_progress")
  @@check([prevId IS NULL OR prevId ~ '^(task|step)://[0-9a-f-]{36}$'], name: "valid_prev_id_format")
  @@check([nextId IS NULL OR nextId ~ '^(task|step)://[0-9a-f-]{36}$'], name: "valid_next_id_format")

  // Indexes
  @@index([projectId])
  @@index([parentTaskId])
  @@index([prevId])
  @@index([nextId])
  @@index([progress])
  @@index([createdAt])

  @@map("tasks")
}
```

#### Step Model

```prisma
model Step {
  id               String    @id @default(uuid()) @db.Uuid
  title            String    @db.VarChar(255)
  taskId           String    @db.Uuid
  prevId           String?   @db.VarChar(50)  // Polymorphic reference
  nextId           String?   @db.VarChar(50)  // Polymorphic reference
  progress         Float     @default(0) @db.DoublePrecision

  // Dual-AI collaboration content
  serverContent    String?   @db.Text
  clientContent    String?   @db.Text
  finalContent     String?   @db.Text

  // Iteration tracking
  iterationCount   Int       @default(0) @db.Integer
  maxIterationsReached Boolean @default(false) @db.Boolean

  // Collaboration state
  status           StepStatus @default(PENDING)
  serverReady      Boolean   @default(false) @db.Boolean
  clientApproved   Boolean   @default(false) @db.Boolean

  // Complexity assessment
  serverComplexity ComplexityLevel?
  clientComplexity ComplexityLevel?
  agreedComplexity ComplexityLevel?
  complexityScore  Float?    @db.DoublePrecision
  shouldPromote    Boolean?  @db.Boolean

  // Metadata
  createdAt        DateTime  @default(now()) @db.Timestamptz
  updatedAt        DateTime  @updatedAt @db.Timestamptz

  // Relationships
  task             Task      @relation(fields: [taskId], references: [id], onDelete: Cascade)
  disputes         Dispute[]
  auditLogs        AuditLog[] @relation("StepAuditLogs")

  // Constraints
  @@check([progress >= 0 && progress <= 1], name: "valid_progress")
  @@check([iterationCount >= 0], name: "non_negative_iteration_count")
  @@check([complexityScore IS NULL OR (complexityScore >= 0 AND complexityScore <= 1)], name: "valid_complexity_score")
  @@check([prevId IS NULL OR prevId ~ '^(task|step)://[0-9a-f-]{36}$'], name: "valid_prev_id_format")
  @@check([nextId IS NULL OR nextId ~ '^(task|step)://[0-9a-f-]{36}$'], name: "valid_next_id_format")

  // Indexes
  @@index([taskId])
  @@index([status])
  @@index([serverReady])
  @@index([clientApproved])
  @@index([prevId])
  @@index([nextId])
  @@index([complexityScore])
  @@index([shouldPromote])
  @@index([progress])
  @@index([createdAt])

  @@map("steps")
}
```

#### Dispute Model

```prisma
model Dispute {
  id              String    @id @default(uuid()) @db.Uuid
  projectId       String    @db.Uuid
  stepId          String    @db.Uuid
  serverContent   String    @db.Text
  clientContent   String    @db.Text
  serverReasoning String    @db.Text
  clientReasoning String    @db.Text
  status          DisputeStatus @default(PENDING)
  userResolution  ResolutionType?
  resolvedContent String?   @db.Text
  resolvedAt      DateTime? @db.Timestamptz
  createdAt       DateTime  @default(now()) @db.Timestamptz

  // Relationships
  project         Project   @relation(fields: [projectId], references: [id], onDelete: Cascade)
  step            Step      @relation(fields: [stepId], references: [id], onDelete: Cascade)
  auditLogs       AuditLog[] @relation("DisputeAuditLogs")

  // Indexes
  @@index([projectId])
  @@index([stepId])
  @@index([status])
  @@index([createdAt])
  @@index([resolvedAt])

  @@map("disputes")
}
```

#### Audit Log Model

```prisma
model AuditLog {
  id          String    @id @default(uuid()) @db.Uuid
  entityType  EntityType
  entityId    String    @db.Uuid
  action      AuditAction
  oldValues   Json?     @db.JsonB
  newValues   Json?     @db.JsonB
  userId      String?   @db.VarChar(255)
  aiAgent     AIAgent?
  metadata    Json?     @db.JsonB
  timestamp   DateTime  @default(now()) @db.Timestamptz

  // Relationships (optional foreign keys for easier querying)
  project     Project?  @relation("ProjectAuditLogs", fields: [entityId], references: [id], onDelete: Cascade)
  task        Task?     @relation("TaskAuditLogs", fields: [entityId], references: [id], onDelete: Cascade)
  step        Step?     @relation("StepAuditLogs", fields: [entityId], references: [id], onDelete: Cascade)
  dispute     Dispute?  @relation("DisputeAuditLogs", fields: [entityId], references: [id], onDelete: Cascade)

  // Indexes
  @@index([entityType, entityId])
  @@index([timestamp])
  @@index([action])
  @@index([userId])
  @@index([aiAgent])

  @@map("audit_logs")
}
```

### Enums

```prisma
enum StepStatus {
  PENDING
  SERVER_DRAFT
  CLIENT_REVIEW
  AGREED
  DISPUTED
  USER_RESOLUTION

  @@map("step_status")
}

enum ComplexityLevel {
  LOW
  MEDIUM
  HIGH

  @@map("complexity_level")
}

enum DisputeStatus {
  PENDING
  RESOLVED

  @@map("dispute_status")
}

enum ResolutionType {
  SERVER
  CLIENT
  CUSTOM
  HYBRID

  @@map("resolution_type")
}

enum EntityType {
  PROJECT
  TASK
  STEP
  DISPUTE

  @@map("entity_type")
}

enum AuditAction {
  CREATED
  UPDATED
  DELETED
  PROMOTED
  DISPUTED
  RESOLVED
  COMPLETED

  @@map("audit_action")
}

enum AIAgent {
  SERVER_AI
  CLIENT_AI

  @@map("ai_agent")
}
```

## Model Extensions

### Custom Types

```go
// internal/database/types.go
package database

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"
)

// ItemReference represents a polymorphic reference
type ItemReference struct {
    Type string `json:"type"` // "task" or "step"
    ID   string `json:"id"`   // UUID
}

func (ir ItemReference) String() string {
    if ir.Type == "" || ir.ID == "" {
        return ""
    }
    return fmt.Sprintf("%s://%s", ir.Type, ir.ID)
}

func ParseItemReference(s string) (*ItemReference, error) {
    if s == "" {
        return nil, nil
    }

    parts := strings.Split(s, "://")
    if len(parts) != 2 {
        return nil, fmt.Errorf("invalid reference format: %s", s)
    }

    if parts[0] != "task" && parts[0] != "step" {
        return nil, fmt.Errorf("invalid reference type: %s", parts[0])
    }

    return &ItemReference{
        Type: parts[0],
        ID:   parts[1],
    }, nil
}

// ComplexityMetrics for JSON storage
type ComplexityMetrics struct {
    ContentLength     int     `json:"content_length"`
    CodeBlocks        int     `json:"code_blocks"`
    ListItems         int     `json:"list_items"`
    ExternalLinks     int     `json:"external_links"`
    SubtaskKeywords   int     `json:"subtask_keywords"`
    TechnicalTerms    int     `json:"technical_terms"`
    ConfigurationSteps int    `json:"configuration_steps"`
}

// Implement sql.Scanner and driver.Valuer for JSON fields
func (cm *ComplexityMetrics) Scan(value interface{}) error {
    if value == nil {
        return nil
    }

    bytes, ok := value.([]byte)
    if !ok {
        return fmt.Errorf("cannot scan %T into ComplexityMetrics", value)
    }

    return json.Unmarshal(bytes, cm)
}

func (cm ComplexityMetrics) Value() (driver.Value, error) {
    return json.Marshal(cm)
}
```

### Model Methods

```go
// internal/models/project.go
package models

import (
    "context"
    "time"

    "github.com/your-org/mcp-planner/internal/database/generated/db"
)

type ProjectModel struct {
    client *db.PrismaClient
}

func NewProjectModel(client *db.PrismaClient) *ProjectModel {
    return &ProjectModel{client: client}
}

func (pm *ProjectModel) Create(ctx context.Context, params CreateProjectParams) (*db.Project, error) {
    return pm.client.Project.CreateOne(
        db.Project.Name.Set(params.Name),
        db.Project.Description.Set(params.Description),
        db.Project.ComplexityThreshold.SetIfPresent(params.ComplexityThreshold),
        db.Project.MaxIterations.SetIfPresent(params.MaxIterations),
    ).Exec(ctx)
}

func (pm *ProjectModel) GetByID(ctx context.Context, id string) (*db.Project, error) {
    return pm.client.Project.FindUnique(
        db.Project.ID.Equals(id),
    ).With(
        db.Project.Tasks.Fetch().With(
            db.Task.Steps.Fetch(),
        ),
    ).Exec(ctx)
}

func (pm *ProjectModel) UpdateProgress(ctx context.Context, id string, progress float64) error {
    _, err := pm.client.Project.FindUnique(
        db.Project.ID.Equals(id),
    ).Update(
        db.Project.Progress.Set(progress),
        db.Project.UpdatedAt.Set(time.Now()),
    ).Exec(ctx)

    return err
}

func (pm *ProjectModel) CalculateProgress(ctx context.Context, id string) (float64, error) {
    // Get all root tasks for the project
    tasks, err := pm.client.Task.FindMany(
        db.Task.ProjectID.Equals(id),
        db.Task.ParentTaskID.IsNull(),
    ).With(
        db.Task.Steps.Fetch(),
        db.Task.SubTasks.Fetch().With(
            db.Task.Steps.Fetch(),
        ),
    ).Exec(ctx)

    if err != nil {
        return 0, err
    }

    if len(tasks) == 0 {
        return 0, nil
    }

    totalProgress := 0.0
    for _, task := range tasks {
        taskProgress := pm.calculateTaskProgress(task)
        totalProgress += taskProgress
    }

    return totalProgress / float64(len(tasks)), nil
}

func (pm *ProjectModel) calculateTaskProgress(task db.Task) float64 {
    totalItems := len(task.Steps()) + len(task.SubTasks())
    if totalItems == 0 {
        return 0.0
    }

    completedProgress := 0.0

    // Sum step progress
    for _, step := range task.Steps() {
        completedProgress += step.Progress
    }

    // Sum sub-task progress (recursive)
    for _, subTask := range task.SubTasks() {
        completedProgress += pm.calculateTaskProgress(subTask)
    }

    return completedProgress / float64(totalItems)
}

type CreateProjectParams struct {
    Name                string
    Description         string
    ComplexityThreshold *float64
    MaxIterations       *int
}
```

### Validation Extensions

```go
// internal/models/validation.go
package models

import (
    "fmt"
    "regexp"
    "strings"

    "github.com/go-playground/validator/v10"
)

var (
    validate = validator.New()

    // Regex for polymorphic references
    polymorphicRefRegex = regexp.MustCompile(`^(task|step)://[0-9a-f-]{36}$`)
)

func init() {
    // Register custom validators
    validate.RegisterValidation("polymorphic_ref", validatePolymorphicRef)
    validate.RegisterValidation("complexity_level", validateComplexityLevel)
    validate.RegisterValidation("step_status", validateStepStatus)
}

func validatePolymorphicRef(fl validator.FieldLevel) bool {
    ref := fl.Field().String()
    if ref == "" {
        return true // Allow empty references
    }
    return polymorphicRefRegex.MatchString(ref)
}

func validateComplexityLevel(fl validator.FieldLevel) bool {
    level := strings.ToUpper(fl.Field().String())
    validLevels := []string{"LOW", "MEDIUM", "HIGH"}

    for _, validLevel := range validLevels {
        if level == validLevel {
            return true
        }
    }
    return false
}

func validateStepStatus(fl validator.FieldLevel) bool {
    status := strings.ToUpper(fl.Field().String())
    validStatuses := []string{
        "PENDING", "SERVER_DRAFT", "CLIENT_REVIEW",
        "AGREED", "DISPUTED", "USER_RESOLUTION",
    }

    for _, validStatus := range validStatuses {
        if status == validStatus {
            return true
        }
    }
    return false
}

// Validation structs for API inputs
type CreateProjectRequest struct {
    Name                string   `json:"name" validate:"required,min=1,max=255"`
    Description         string   `json:"description" validate:"required,min=10"`
    ComplexityThreshold *float64 `json:"complexity_threshold,omitempty" validate:"omitempty,min=0,max=1"`
    MaxIterations       *int     `json:"max_iterations,omitempty" validate:"omitempty,min=1,max=10"`
}

type CreateTaskRequest struct {
    Title        string  `json:"title" validate:"required,min=1,max=255"`
    Objective    string  `json:"objective" validate:"required,min=10"`
    ParentTaskID *string `json:"parent_task_id,omitempty" validate:"omitempty,uuid"`
    PrevID       *string `json:"prev_id,omitempty" validate:"omitempty,polymorphic_ref"`
    NextID       *string `json:"next_id,omitempty" validate:"omitempty,polymorphic_ref"`
}

type CreateStepRequest struct {
    Title  string  `json:"title" validate:"required,min=1,max=255"`
    PrevID *string `json:"prev_id,omitempty" validate:"omitempty,polymorphic_ref"`
    NextID *string `json:"next_id,omitempty" validate:"omitempty,polymorphic_ref"`
}

type UpdateStepContentRequest struct {
    ServerContent    *string `json:"server_content,omitempty"`
    ClientContent    *string `json:"client_content,omitempty"`
    ServerComplexity *string `json:"server_complexity,omitempty" validate:"omitempty,complexity_level"`
    ClientComplexity *string `json:"client_complexity,omitempty" validate:"omitempty,complexity_level"`
    Status           *string `json:"status,omitempty" validate:"omitempty,step_status"`
}

func ValidateStruct(s interface{}) error {
    if err := validate.Struct(s); err != nil {
        return formatValidationError(err)
    }
    return nil
}

func formatValidationError(err error) error {
    var messages []string

    for _, err := range err.(validator.ValidationErrors) {
        switch err.Tag() {
        case "required":
            messages = append(messages, fmt.Sprintf("%s is required", err.Field()))
        case "min":
            messages = append(messages, fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param()))
        case "max":
            messages = append(messages, fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param()))
        case "uuid":
            messages = append(messages, fmt.Sprintf("%s must be a valid UUID", err.Field()))
        case "polymorphic_ref":
            messages = append(messages, fmt.Sprintf("%s must be in format 'task://uuid' or 'step://uuid'", err.Field()))
        case "complexity_level":
            messages = append(messages, fmt.Sprintf("%s must be 'low', 'medium', or 'high'", err.Field()))
        case "step_status":
            messages = append(messages, fmt.Sprintf("%s must be a valid step status", err.Field()))
        default:
            messages = append(messages, fmt.Sprintf("%s is invalid", err.Field()))
        }
    }

    return fmt.Errorf("validation errors: %s", strings.Join(messages, ", "))
}
```

## Database Triggers

### Automatic Progress Updates

```sql
-- Trigger to automatically update parent progress when child changes
CREATE OR REPLACE FUNCTION update_parent_progress()
RETURNS TRIGGER AS $$
BEGIN
    -- Update task progress when step changes
    IF TG_TABLE_NAME = 'steps' THEN
        UPDATE tasks
        SET progress = (
            SELECT COALESCE(AVG(progress), 0)
            FROM steps
            WHERE task_id = NEW.task_id
        ),
        updated_at = NOW()
        WHERE id = NEW.task_id;
    END IF;

    -- Update parent task progress when sub-task changes
    IF TG_TABLE_NAME = 'tasks' AND NEW.parent_task_id IS NOT NULL THEN
        UPDATE tasks
        SET progress = (
            SELECT COALESCE(
                (
                    SELECT AVG(progress) FROM steps WHERE task_id = NEW.parent_task_id
                ) + (
                    SELECT AVG(progress) FROM tasks WHERE parent_task_id = NEW.parent_task_id
                )
            ) / 2, 0)
        ),
        updated_at = NOW()
        WHERE id = NEW.parent_task_id;
    END IF;

    -- Update project progress when root task changes
    IF TG_TABLE_NAME = 'tasks' AND NEW.parent_task_id IS NULL THEN
        UPDATE projects
        SET progress = (
            SELECT COALESCE(AVG(progress), 0)
            FROM tasks
            WHERE project_id = NEW.project_id
            AND parent_task_id IS NULL
        ),
        updated_at = NOW()
        WHERE id = NEW.project_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers
CREATE TRIGGER step_progress_update
    AFTER UPDATE OF progress ON steps
    FOR EACH ROW
    EXECUTE FUNCTION update_parent_progress();

CREATE TRIGGER task_progress_update
    AFTER UPDATE OF progress ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_parent_progress();
```

### Audit Logging Trigger

```sql
-- Automatic audit logging
CREATE OR REPLACE FUNCTION audit_changes()
RETURNS TRIGGER AS $$
DECLARE
    entity_type_val entity_type;
    action_val audit_action;
BEGIN
    -- Determine entity type
    CASE TG_TABLE_NAME
        WHEN 'projects' THEN entity_type_val := 'PROJECT';
        WHEN 'tasks' THEN entity_type_val := 'TASK';
        WHEN 'steps' THEN entity_type_val := 'STEP';
        WHEN 'disputes' THEN entity_type_val := 'DISPUTE';
    END CASE;

    -- Determine action
    CASE TG_OP
        WHEN 'INSERT' THEN action_val := 'CREATED';
        WHEN 'UPDATE' THEN action_val := 'UPDATED';
        WHEN 'DELETE' THEN action_val := 'DELETED';
    END CASE;

    -- Insert audit record
    INSERT INTO audit_logs (
        entity_type,
        entity_id,
        action,
        old_values,
        new_values,
        timestamp
    ) VALUES (
        entity_type_val,
        COALESCE(NEW.id, OLD.id),
        action_val,
        CASE WHEN TG_OP = 'DELETE' THEN to_jsonb(OLD) ELSE NULL END,
        CASE WHEN TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN to_jsonb(NEW) ELSE NULL END,
        NOW()
    );

    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

-- Create audit triggers for all tables
CREATE TRIGGER projects_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON projects
    FOR EACH ROW EXECUTE FUNCTION audit_changes();

CREATE TRIGGER tasks_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON tasks
    FOR EACH ROW EXECUTE FUNCTION audit_changes();

CREATE TRIGGER steps_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON steps
    FOR EACH ROW EXECUTE FUNCTION audit_changes();

CREATE TRIGGER disputes_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON disputes
    FOR EACH ROW EXECUTE FUNCTION audit_changes();
```

---

*Next: [Database Connection & Setup](./07b2-database-connection.md)*
