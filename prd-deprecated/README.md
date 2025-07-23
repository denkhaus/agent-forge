# MCP-Planner: AI-Collaborative Task Management System

## Overview

MCP-Planner is a sophisticated Model Context Protocol (MCP) server that provides AI-driven task management with hierarchical project structures. The system enables collaborative task planning between two AI agents - a server-side content creator and a client-side supervisor - with human oversight for dispute resolution.

## Document Structure

This PRD is organized into the following documents:

- **[System Architecture](./01-system-architecture.md)** - Core system design, data models, and technical architecture
- **[Dual-AI Collaboration](./02-dual-ai-collaboration.md)** - Detailed workflow for AI-to-AI collaboration
- **[Data Schema](./03-data-schema.md)** - Complete database schema and relationships
- **[MCP Functions](./04-mcp-functions.md)** - All MCP server functions and their specifications
- **[Complexity Management](./05-complexity-management.md)** - Automated complexity analysis and step promotion
- **[User Interface](./06-user-interface.md)** - Human interaction patterns and dispute resolution
- **[Implementation Guide](./07-implementation-guide.md)** - Technical implementation details and Go/Prisma setup

## Key Features

### 🏗️ Hierarchical Task Structure
- **Projects** contain multiple **Tasks**
- **Tasks** can have **Sub-tasks** (unlimited nesting)
- **Steps** belong to tasks and contain detailed implementation instructions
- Polymorphic navigation: Steps and Tasks can reference each other in prev/next chains

### 🤖 Dual-AI Collaboration
- **Server-side AI**: Content creator with project context but step-focused
- **Client-side AI**: Supervisor with full project context and oversight
- **Iterative refinement**: AIs collaborate to improve step content
- **Human oversight**: Users resolve disputes when AIs disagree

### 📊 Intelligent Progress Tracking
- Hierarchical progress calculation (children → parent)
- Real-time progress updates when steps/tasks complete
- Project-level progress aggregation

### 🧠 Complexity Management
- AI-driven complexity analysis for steps
- Automatic promotion of complex steps to tasks
- Configurable complexity thresholds per project
- Iterative optimization until convergence

### 🔄 Smart Navigation
- Get next actionable item in project workflow
- Polymorphic prev/next relationships using namespaced IDs
- Sequential execution with no step skipping

## Quick Start Example

```go
// 1. Create project with AI collaboration
project := CreateProject("Build Authentication System",
    "Create secure JWT-based auth with registration, login, password reset, and RBAC")

// 2. Client AI defines root tasks
tasks := DefineRootTaskObjectives(project.ID, [
    "Implement user registration with email verification",
    "Create secure login/logout functionality"
])

// 3. Server AI crafts step content, Client AI reviews
step := CreateStep(tasks[0].ID, "Setup JWT library", "")
CraftStepContentWithContext(step.ID, "Focus on library selection and configuration")
ReviewStepContent(step.ID, clientRefinedContent)

// 4. Execute optimized workflow
nextItem := GetNextActionableItem(project.ID)
MarkComplete(nextItem.ID, nextItem.Type)
```

## Technology Stack

- **Language**: Go
- **Database**: PostgreSQL
- **ORM**: Prisma
- **Protocol**: Model Context Protocol (MCP)
- **AI Integration**: Configurable AI providers for complexity analysis

## Success Metrics

- **Collaboration Efficiency**: Percentage of steps agreed upon without disputes
- **Complexity Accuracy**: How often AI complexity assessments match human judgment
- **Project Completion**: Time from project creation to completion
- **User Satisfaction**: Dispute resolution effectiveness and ease of use

---

*For detailed technical specifications, please refer to the individual documents in this PRD.*
