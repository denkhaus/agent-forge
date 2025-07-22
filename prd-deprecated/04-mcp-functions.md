# MCP Functions

## Overview

The MCP-Planner server exposes a comprehensive set of functions through the Model Context Protocol (MCP) interface. These functions enable AI agents and clients to manage projects, tasks, steps, and collaborate on content creation with dispute resolution capabilities.

## Function Categories

### 1. Project Management
### 2. Task Management
### 3. Step Management
### 4. Dual-AI Collaboration
### 5. Navigation & Workflow
### 6. Complexity Management
### 7. Dispute Resolution
### 8. Analytics & Reporting

---

## 1. Project Management

### CreateProject
Creates a new project with mandatory description for AI task planning.

```json
{
  "name": "createProject",
  "description": "Create a new project with description and configuration",
  "inputSchema": {
    "type": "object",
    "properties": {
      "name": {
        "type": "string",
        "description": "Project name"
      },
      "description": {
        "type": "string",
        "description": "Mandatory project description - the initial idea for task tree building"
      },
      "complexityThreshold": {
        "type": "number",
        "minimum": 0,
        "maximum": 1,
        "default": 0.7,
        "description": "Threshold for step complexity promotion (0.0-1.0)"
      },
      "maxIterations": {
        "type": "integer",
        "minimum": 1,
        "default": 3,
        "description": "Maximum AI collaboration iterations per step"
      }
    },
    "required": ["name", "description"]
  }
}
```

**Response:**
```json
{
  "id": "uuid",
  "name": "E-commerce Platform",
  "description": "Build a modern e-commerce platform...",
  "progress": 0.0,
  "complexityThreshold": 0.7,
  "maxIterations": 3,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### GetProject
Retrieves project details with current progress and statistics.

```json
{
  "name": "getProject",
  "description": "Get project details and current status",
  "inputSchema": {
    "type": "object",
    "properties": {
      "projectId": {
        "type": "string",
        "description": "Project UUID"
      },
      "includeStats": {
        "type": "boolean",
        "default": false,
        "description": "Include detailed statistics"
      }
    },
    "required": ["projectId"]
  }
}
```

### UpdateProjectSettings
Updates project configuration including complexity threshold and iteration limits.

```json
{
  "name": "updateProjectSettings",
  "description": "Update project configuration settings",
  "inputSchema": {
    "type": "object",
    "properties": {
      "projectId": {
        "type": "string",
        "description": "Project UUID"
      },
      "complexityThreshold": {
        "type": "number",
        "minimum": 0,
        "maximum": 1,
        "description": "New complexity threshold - triggers re-evaluation"
      },
      "maxIterations": {
        "type": "integer",
        "minimum": 1,
        "description": "New maximum iteration limit"
      }
    },
    "required": ["projectId"]
  }
}
```

### ListProjects
Lists all projects with optional filtering and pagination.

```json
{
  "name": "listProjects",
  "description": "List projects with filtering options",
  "inputSchema": {
    "type": "object",
    "properties": {
      "limit": {
        "type": "integer",
        "default": 50,
        "description": "Maximum number of projects to return"
      },
      "offset": {
        "type": "integer",
        "default": 0,
        "description": "Number of projects to skip"
      },
      "status": {
        "type": "string",
        "enum": ["active", "completed", "all"],
        "default": "all",
        "description": "Filter by project status"
      }
    }
  }
}
```

---

## 2. Task Management

### CreateTask
Creates a new task within a project or as a sub-task.

```json
{
  "name": "createTask",
  "description": "Create a new task or sub-task",
  "inputSchema": {
    "type": "object",
    "properties": {
      "projectId": {
        "type": "string",
        "description": "Project UUID"
      },
      "parentTaskId": {
        "type": "string",
        "description": "Parent task UUID (null for root tasks)"
      },
      "title": {
        "type": "string",
        "description": "Task title"
      },
      "objective": {
        "type": "string",
        "description": "What this task should accomplish"
      },
      "prevId": {
        "type": "string",
        "description": "Previous item reference (task://uuid or step://uuid)"
      },
      "nextId": {
        "type": "string",
        "description": "Next item reference (task://uuid or step://uuid)"
      }
    },
    "required": ["projectId", "title", "objective"]
  }
}
```

### DefineRootTaskObjectives
Client AI function to define initial root tasks from project description.

```json
{
  "name": "defineRootTaskObjectives",
  "description": "Define root task objectives from project description (Client AI)",
  "inputSchema": {
    "type": "object",
    "properties": {
      "projectId": {
        "type": "string",
        "description": "Project UUID"
      },
      "taskObjectives": {
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "title": {
              "type": "string",
              "description": "Task title"
            },
            "objective": {
              "type": "string",
              "description": "Task objective description"
            }
          },
          "required": ["title", "objective"]
        },
        "description": "Array of root task objectives"
      }
    },
    "required": ["projectId", "taskObjectives"]
  }
}
```

### GetTask
Retrieves task details with progress and relationships.

```json
{
  "name": "getTask",
  "description": "Get task details and current status",
  "inputSchema": {
    "type": "object",
    "properties": {
      "taskId": {
        "type": "string",
        "description": "Task UUID"
      },
      "includeSteps": {
        "type": "boolean",
        "default": false,
        "description": "Include task steps"
      },
      "includeSubTasks": {
        "type": "boolean",
        "default": false,
        "description": "Include sub-tasks"
      }
    },
    "required": ["taskId"]
  }
}
```

### UpdateTask
Updates task properties and relationships.

```json
{
  "name": "updateTask",
  "description": "Update task properties",
  "inputSchema": {
    "type": "object",
    "properties": {
      "taskId": {
        "type": "string",
        "description": "Task UUID"
      },
      "title": {
        "type": "string",
        "description": "New task title"
      },
      "objective": {
        "type": "string",
        "description": "New task objective"
      },
      "prevId": {
        "type": "string",
        "description": "New previous item reference"
      },
      "nextId": {
        "type": "string",
        "description": "New next item reference"
      }
    },
    "required": ["taskId"]
  }
}
```

---

## 3. Step Management

### CreateStep
Creates a new step within a task.

```json
{
  "name": "createStep",
  "description": "Create a new step within a task",
  "inputSchema": {
    "type": "object",
    "properties": {
      "taskId": {
        "type": "string",
        "description": "Parent task UUID"
      },
      "title": {
        "type": "string",
        "description": "Step title"
      },
      "prevId": {
        "type": "string",
        "description": "Previous item reference (task://uuid or step://uuid)"
      },
      "nextId": {
        "type": "string",
        "description": "Next item reference (task://uuid or step://uuid)"
      }
    },
    "required": ["taskId", "title"]
  }
}
```

### GetStep
Retrieves step details with collaboration status.

```json
{
  "name": "getStep",
  "description": "Get step details and collaboration status",
  "inputSchema": {
    "type": "object",
    "properties": {
      "stepId": {
        "type": "string",
        "description": "Step UUID"
      },
      "includeContent": {
        "type": "boolean",
        "default": true,
        "description": "Include all content versions"
      }
    },
    "required": ["stepId"]
  }
}
```

**Response:**
```json
{
  "id": "uuid",
  "title": "Setup JWT library",
  "taskId": "uuid",
  "prevId": "task://uuid",
  "nextId": "step://uuid",
  "progress": 0.0,
  "serverContent": "Install jsonwebtoken package...",
  "clientContent": "Install and configure jsonwebtoken...",
  "finalContent": null,
  "iterationCount": 2,
  "maxIterationsReached": false,
  "status": "client_review",
  "serverReady": true,
  "clientApproved": false,
  "serverComplexity": "medium",
  "clientComplexity": "low",
  "agreedComplexity": null,
  "complexityScore": 0.6,
  "shouldPromote": false,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### MarkStepComplete
Marks a step as completed and triggers progress recalculation.

```json
{
  "name": "markStepComplete",
  "description": "Mark step as completed (progress = 1.0)",
  "inputSchema": {
    "type": "object",
    "properties": {
      "stepId": {
        "type": "string",
        "description": "Step UUID"
      }
    },
    "required": ["stepId"]
  }
}
```

---

## 4. Dual-AI Collaboration

### CraftStepContentWithContext
Server AI function to create step content with project context.

```json
{
  "name": "craftStepContentWithContext",
  "description": "Create step content with project context (Server AI)",
  "inputSchema": {
    "type": "object",
    "properties": {
      "stepId": {
        "type": "string",
        "description": "Step UUID"
      },
      "focusPrompt": {
        "type": "string",
        "description": "Specific guidance for step focus"
      },
      "aiProvider": {
        "type": "string",
        "enum": ["openai", "anthropic", "local"],
        "default": "openai",
        "description": "AI provider to use"
      }
    },
    "required": ["stepId"]
  }
}
```

**Response:**
```json
{
  "stepId": "uuid",
  "content": "## Setup JWT Library\n\n1. Install package...",
  "reasoning": "This step focuses on library setup because...",
  "complexity": "medium",
  "iterationCount": 1,
  "canContinue": true,
  "serverReady": true
}
```

### ReviewStepContent
Client AI function to review and refine server-generated content.

```json
{
  "name": "reviewStepContent",
  "description": "Review and refine step content (Client AI)",
  "inputSchema": {
    "type": "object",
    "properties": {
      "stepId": {
        "type": "string",
        "description": "Step UUID"
      },
      "clientContent": {
        "type": "string",
        "description": "Client AI refined content (markdown)"
      },
      "clientReasoning": {
        "type": "string",
        "description": "Client AI reasoning for changes"
      },
      "approved": {
        "type": "boolean",
        "description": "Whether client AI approves the content"
      },
      "aiProvider": {
        "type": "string",
        "enum": ["openai", "anthropic", "local"],
        "default": "openai",
        "description": "AI provider to use"
      }
    },
    "required": ["stepId", "approved"]
  }
}
```

### SignalContentReady
Server AI signals that content is ready for client review.

```json
{
  "name": "signalContentReady",
  "description": "Signal that server content is ready for review",
  "inputSchema": {
    "type": "object",
    "properties": {
      "stepId": {
        "type": "string",
        "description": "Step UUID"
      }
    },
    "required": ["stepId"]
  }
}
```

### ApproveStepContent
Client AI approves step content for finalization.

```json
{
  "name": "approveStepContent",
  "description": "Approve step content for finalization (Client AI)",
  "inputSchema": {
    "type": "object",
    "properties": {
      "stepId": {
        "type": "string",
        "description": "Step UUID"
      },
      "finalContent": {
        "type": "string",
        "description": "Final approved content (markdown)"
      }
    },
    "required": ["stepId"]
  }
}
```

### RequestContentRevision
Client AI requests revision with specific feedback.

```json
{
  "name": "requestContentRevision",
  "description": "Request content revision with feedback",
  "inputSchema": {
    "type": "object",
    "properties": {
      "stepId": {
        "type": "string",
        "description": "Step UUID"
      },
      "feedback": {
        "type": "string",
        "description": "Specific feedback for revision"
      }
    },
    "required": ["stepId", "feedback"]
  }
}
```

---

## 5. Navigation & Workflow

### GetNextActionableItem
Gets the next step or task that needs to be worked on.

```json
{
  "name": "getNextActionableItem",
  "description": "Get the next actionable item in project workflow",
  "inputSchema": {
    "type": "object",
    "properties": {
      "projectId": {
        "type": "string",
        "description": "Project UUID"
      }
    },
    "required": ["projectId"]
  }
}
```

**Response:**
```json
{
  "type": "step",
  "id": "uuid",
  "title": "Setup JWT library",
  "taskId": "uuid",
  "taskTitle": "Implement Authentication",
  "status": "pending",
  "progress": 0.0,
  "isBlocked": false,
  "blockingReason": null
}
```

### GetProjectWorkflow
Gets the complete workflow structure for a project.

```json
{
  "name": "getProjectWorkflow",
  "description": "Get complete project workflow structure",
  "inputSchema": {
    "type": "object",
    "properties": {
      "projectId": {
        "type": "string",
        "description": "Project UUID"
      },
      "includeCompleted": {
        "type": "boolean",
        "default": false,
        "description": "Include completed items"
      }
    },
    "required": ["projectId"]
  }
}
```

### GetNavigationChain
Gets the navigation chain from a specific item.

```json
{
  "name": "getNavigationChain",
  "description": "Get navigation chain from specific item",
  "inputSchema": {
    "type": "object",
    "properties": {
      "itemId": {
        "type": "string",
        "description": "Starting item UUID"
      },
      "itemType": {
        "type": "string",
        "enum": ["task", "step"],
        "description": "Type of starting item"
      },
      "direction": {
        "type": "string",
        "enum": ["forward", "backward", "both"],
        "default": "forward",
        "description": "Navigation direction"
      },
      "maxDepth": {
        "type": "integer",
        "default": 10,
        "description": "Maximum chain depth"
      }
    },
    "required": ["itemId", "itemType"]
  }
}
```

---

## 6. Complexity Management

### AnalyzeStepComplexity
Analyzes step complexity using AI.

```json
{
  "name": "analyzeStepComplexity",
  "description": "Analyze step complexity using AI",
  "inputSchema": {
    "type": "object",
    "properties": {
      "stepId": {
        "type": "string",
        "description": "Step UUID"
      },
      "aiProvider": {
        "type": "string",
        "enum": ["openai", "anthropic", "local"],
        "default": "openai",
        "description": "AI provider for analysis"
      },
      "forceReanalysis": {
        "type": "boolean",
        "default": false,
        "description": "Force re-analysis even if already analyzed"
      }
    },
    "required": ["stepId"]
  }
}
```

**Response:**
```json
{
  "stepId": "uuid",
  "complexity": "high",
  "complexityScore": 0.85,
  "reasoning": "This step involves multiple sub-tasks including...",
  "shouldPromote": true,
  "suggestedSubtasks": [
    "Install JWT library",
    "Configure JWT settings",
    "Create JWT middleware"
  ],
  "analyzedAt": "2024-01-01T00:00:00Z"
}
```

### PromoteStepToTask
Promotes a complex step to a sub-task.

```json
{
  "name": "promoteStepToTask",
  "description": "Promote a step to a sub-task",
  "inputSchema": {
    "type": "object",
    "properties": {
      "stepId": {
        "type": "string",
        "description": "Step UUID to promote"
      },
      "preserveContent": {
        "type": "boolean",
        "default": true,
        "description": "Preserve step content as task objective"
      }
    },
    "required": ["stepId"]
  }
}
```

### OptimizeProjectComplexity
Optimizes project complexity by promoting complex steps.

```json
{
  "name": "optimizeProjectComplexity",
  "description": "Optimize project by promoting complex steps to tasks",
  "inputSchema": {
    "type": "object",
    "properties": {
      "projectId": {
        "type": "string",
        "description": "Project UUID"
      },
      "maxIterations": {
        "type": "integer",
        "default": 10,
        "description": "Maximum optimization iterations"
      },
      "forceReanalysis": {
        "type": "boolean",
        "default": false,
        "description": "Force re-analysis of all steps"
      }
    },
    "required": ["projectId"]
  }
}
```

**Response:**
```json
{
  "projectId": "uuid",
  "optimizationComplete": true,
  "iterationsUsed": 3,
  "promotedSteps": [
    {
      "stepId": "uuid",
      "newTaskId": "uuid",
      "reason": "High complexity score (0.85)"
    }
  ],
  "remainingComplexSteps": [],
  "convergenceReason": "no_changes"
}
```

---

## 7. Dispute Resolution

### GetPendingDisputes
Gets disputes awaiting user resolution.

```json
{
  "name": "getPendingDisputes",
  "description": "Get disputes awaiting user resolution",
  "inputSchema": {
    "type": "object",
    "properties": {
      "projectId": {
        "type": "string",
        "description": "Project UUID (optional filter)"
      },
      "limit": {
        "type": "integer",
        "default": 50,
        "description": "Maximum disputes to return"
      }
    }
  }
}
```

**Response:**
```json
{
  "disputes": [
    {
      "id": "uuid",
      "stepId": "uuid",
      "stepTitle": "Setup JWT library",
      "projectId": "uuid",
      "projectName": "E-commerce Platform",
      "serverContent": "Install jsonwebtoken package...",
      "clientContent": "Install and configure jsonwebtoken...",
      "serverReasoning": "Focused on basic installation...",
      "clientReasoning": "Need configuration details...",
      "iterationHistory": [...],
      "createdAt": "2024-01-01T00:00:00Z"
    }
  ],
  "totalCount": 1
}
```

### ResolveDispute
Resolves a dispute with user decision.

```json
{
  "name": "resolveDispute",
  "description": "Resolve dispute with user decision",
  "inputSchema": {
    "type": "object",
    "properties": {
      "disputeId": {
        "type": "string",
        "description": "Dispute UUID"
      },
      "resolution": {
        "type": "string",
        "enum": ["server", "client", "custom", "hybrid"],
        "description": "Resolution type"
      },
      "customContent": {
        "type": "string",
        "description": "Custom content (required if resolution=custom)"
      },
      "hybridContent": {
        "type": "string",
        "description": "Hybrid content (required if resolution=hybrid)"
      }
    },
    "required": ["disputeId", "resolution"]
  }
}
```

---

## 8. Analytics & Reporting

### GetProjectStats
Gets comprehensive project statistics.

```json
{
  "name": "getProjectStats",
  "description": "Get comprehensive project statistics",
  "inputSchema": {
    "type": "object",
    "properties": {
      "projectId": {
        "type": "string",
        "description": "Project UUID"
      }
    },
    "required": ["projectId"]
  }
}
```

**Response:**
```json
{
  "projectId": "uuid",
  "progress": 0.45,
  "totalTasks": 12,
  "completedTasks": 5,
  "totalSteps": 48,
  "completedSteps": 22,
  "pendingSteps": 15,
  "disputedSteps": 2,
  "averageIterationsPerStep": 2.1,
  "complexityDistribution": {
    "low": 30,
    "medium": 15,
    "high": 3
  },
  "collaborationMetrics": {
    "agreementRate": 0.85,
    "averageResolutionTime": "2.5 hours",
    "disputeRate": 0.04
  },
  "estimatedCompletion": "2024-02-15T00:00:00Z"
}
```

### GetWorkQueue
Gets work queue for AI agents.

```json
{
  "name": "getWorkQueue",
  "description": "Get work queue for AI agents",
  "inputSchema": {
    "type": "object",
    "properties": {
      "agentType": {
        "type": "string",
        "enum": ["server", "client"],
        "description": "AI agent type"
      },
      "projectId": {
        "type": "string",
        "description": "Project UUID (optional filter)"
      },
      "limit": {
        "type": "integer",
        "default": 10,
        "description": "Maximum items to return"
      }
    },
    "required": ["agentType"]
  }
}
```

---

## Error Handling

All MCP functions return standardized error responses:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid step ID format",
    "details": {
      "field": "stepId",
      "value": "invalid-uuid",
      "expected": "Valid UUID"
    }
  }
}
```

### Error Codes
- `VALIDATION_ERROR`: Input validation failed
- `NOT_FOUND`: Resource not found
- `PERMISSION_DENIED`: Access denied
- `CONFLICT`: Resource conflict (e.g., circular references)
- `AI_PROVIDER_ERROR`: AI service unavailable
- `DATABASE_ERROR`: Database operation failed
- `ITERATION_LIMIT_EXCEEDED`: Max iterations reached

---

*Next: [Complexity Management](./05-complexity-management.md)*
