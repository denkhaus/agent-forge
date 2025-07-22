# User Interface

## Overview

The MCP-Planner system provides human interfaces for dispute resolution, project monitoring, and system administration. While the core system operates through AI-to-AI collaboration, human oversight is essential for resolving conflicts and maintaining quality control.

## Interface Types

### 1. Dispute Resolution Interface
### 2. Project Dashboard
### 3. Administrative Interface
### 4. API Integration Points

---

## 1. Dispute Resolution Interface

### Dispute Presentation

When AI agents cannot reach agreement within iteration limits, disputes are escalated to human users for resolution.

```typescript
interface DisputePresentation {
  dispute: {
    id: string;
    stepId: string;
    stepTitle: string;
    projectName: string;
    createdAt: string;
    iterationCount: number;
  };
  context: {
    projectDescription: string;
    taskObjective: string;
    previousSteps: StepSummary[];
    nextSteps: StepSummary[];
  };
  aiVersions: {
    server: {
      content: string;        // Markdown
      reasoning: string;
      complexity: string;
      confidence: number;
    };
    client: {
      content: string;        // Markdown
      reasoning: string;
      complexity: string;
      confidence: number;
    };
  };
  iterationHistory: IterationRecord[];
  resolutionOptions: ResolutionOption[];
}
```

### Resolution Interface Components

#### Side-by-Side Content Comparison
```html
<div class="dispute-resolution">
  <div class="dispute-header">
    <h2>Resolve Content Dispute</h2>
    <div class="dispute-meta">
      <span class="project">{{projectName}}</span>
      <span class="step">{{stepTitle}}</span>
      <span class="iterations">{{iterationCount}} iterations</span>
    </div>
  </div>

  <div class="context-panel">
    <h3>Context</h3>
    <div class="project-description">{{projectDescription}}</div>
    <div class="task-objective">{{taskObjective}}</div>
    <div class="step-sequence">
      <!-- Previous and next steps -->
    </div>
  </div>

  <div class="content-comparison">
    <div class="server-version">
      <h3>Server AI Version</h3>
      <div class="ai-meta">
        <span class="complexity">{{serverComplexity}}</span>
        <span class="confidence">{{serverConfidence}}% confidence</span>
      </div>
      <div class="content-preview">
        <!-- Rendered markdown -->
      </div>
      <div class="reasoning">
        <h4>Reasoning</h4>
        <p>{{serverReasoning}}</p>
      </div>
    </div>

    <div class="client-version">
      <h3>Client AI Version</h3>
      <div class="ai-meta">
        <span class="complexity">{{clientComplexity}}</span>
        <span class="confidence">{{clientConfidence}}% confidence</span>
      </div>
      <div class="content-preview">
        <!-- Rendered markdown -->
      </div>
      <div class="reasoning">
        <h4>Reasoning</h4>
        <p>{{clientReasoning}}</p>
      </div>
    </div>
  </div>

  <div class="resolution-options">
    <h3>Resolution Options</h3>
    <div class="option-buttons">
      <button class="resolution-btn" data-type="server">
        Use Server Version
      </button>
      <button class="resolution-btn" data-type="client">
        Use Client Version
      </button>
      <button class="resolution-btn" data-type="hybrid">
        Create Hybrid
      </button>
      <button class="resolution-btn" data-type="custom">
        Write Custom
      </button>
    </div>
  </div>
</div>
```

#### Hybrid Content Editor
```html
<div class="hybrid-editor" style="display: none;">
  <h3>Create Hybrid Content</h3>
  <div class="source-panels">
    <div class="source-panel">
      <h4>Server Content</h4>
      <div class="selectable-content" data-source="server">
        <!-- Selectable text blocks -->
      </div>
    </div>
    <div class="source-panel">
      <h4>Client Content</h4>
      <div class="selectable-content" data-source="client">
        <!-- Selectable text blocks -->
      </div>
    </div>
  </div>
  <div class="hybrid-preview">
    <h4>Hybrid Preview</h4>
    <textarea class="hybrid-content" placeholder="Combine content from both versions..."></textarea>
    <div class="preview-render"></div>
  </div>
</div>
```

#### Custom Content Editor
```html
<div class="custom-editor" style="display: none;">
  <h3>Write Custom Content</h3>
  <div class="reference-content">
    <details>
      <summary>Reference: Server AI Content</summary>
      <div class="reference-text">{{serverContent}}</div>
    </details>
    <details>
      <summary>Reference: Client AI Content</summary>
      <div class="reference-text">{{clientContent}}</div>
    </details>
  </div>
  <div class="editor-container">
    <textarea class="custom-content" placeholder="Write your custom step content in Markdown..."></textarea>
    <div class="editor-toolbar">
      <button class="format-btn" data-format="bold">Bold</button>
      <button class="format-btn" data-format="italic">Italic</button>
      <button class="format-btn" data-format="code">Code</button>
      <button class="format-btn" data-format="list">List</button>
    </div>
    <div class="preview-pane">
      <!-- Live markdown preview -->
    </div>
  </div>
</div>
```

### Resolution Workflow

```typescript
interface ResolutionWorkflow {
  // User selects resolution type
  selectResolutionType(type: 'server' | 'client' | 'hybrid' | 'custom'): void;

  // For hybrid resolution
  selectContentBlocks(source: 'server' | 'client', blocks: string[]): void;
  combineSelectedBlocks(): string;

  // For custom resolution
  editCustomContent(content: string): void;
  previewContent(content: string): string;

  // Final submission
  submitResolution(resolution: {
    type: string;
    content: string;
    reasoning?: string;
  }): Promise<void>;
}

// Implementation
class DisputeResolver {
  async resolveDispute(disputeId: string, resolution: Resolution): Promise<void> {
    const response = await fetch('/api/disputes/resolve', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        disputeId,
        resolution: resolution.type,
        customContent: resolution.content,
        userReasoning: resolution.reasoning
      })
    });

    if (response.ok) {
      // Redirect to project dashboard or next dispute
      window.location.href = `/projects/${projectId}`;
    }
  }
}
```

---

## 2. Project Dashboard

### Project Overview

```html
<div class="project-dashboard">
  <div class="project-header">
    <h1>{{projectName}}</h1>
    <div class="project-meta">
      <span class="progress">{{progress}}% Complete</span>
      <span class="status">{{status}}</span>
      <span class="created">Created {{createdAt}}</span>
    </div>
  </div>

  <div class="project-description">
    <h2>Project Vision</h2>
    <div class="description-content">{{description}}</div>
  </div>

  <div class="dashboard-grid">
    <div class="progress-panel">
      <h3>Progress Overview</h3>
      <div class="progress-bar">
        <div class="progress-fill" style="width: {{progress}}%"></div>
      </div>
      <div class="progress-stats">
        <div class="stat">
          <span class="label">Tasks</span>
          <span class="value">{{completedTasks}}/{{totalTasks}}</span>
        </div>
        <div class="stat">
          <span class="label">Steps</span>
          <span class="value">{{completedSteps}}/{{totalSteps}}</span>
        </div>
      </div>
    </div>

    <div class="workflow-panel">
      <h3>Current Workflow</h3>
      <div class="next-action">
        <h4>Next Action</h4>
        <div class="action-item">
          <span class="type">{{nextItem.type}}</span>
          <span class="title">{{nextItem.title}}</span>
          <span class="status">{{nextItem.status}}</span>
        </div>
      </div>
      <div class="recent-activity">
        <h4>Recent Activity</h4>
        <ul class="activity-list">
          <!-- Recent completions, disputes, etc. -->
        </ul>
      </div>
    </div>

    <div class="disputes-panel">
      <h3>Active Disputes</h3>
      <div class="dispute-count">{{disputeCount}} pending</div>
      <div class="dispute-list">
        <!-- List of disputes needing resolution -->
      </div>
      <button class="resolve-disputes-btn">Resolve Disputes</button>
    </div>

    <div class="complexity-panel">
      <h3>Complexity Analysis</h3>
      <div class="complexity-chart">
        <!-- Chart showing complexity distribution -->
      </div>
      <div class="complexity-settings">
        <label>Complexity Threshold</label>
        <input type="range" min="0" max="1" step="0.1"
               value="{{complexityThreshold}}"
               class="threshold-slider">
        <button class="retrigger-btn">Re-evaluate</button>
      </div>
    </div>
  </div>
</div>
```

### Task Tree Visualization

```html
<div class="task-tree">
  <h3>Project Structure</h3>
  <div class="tree-controls">
    <button class="expand-all">Expand All</button>
    <button class="collapse-all">Collapse All</button>
    <select class="view-filter">
      <option value="all">All Items</option>
      <option value="pending">Pending Only</option>
      <option value="completed">Completed Only</option>
      <option value="disputed">Disputed Only</option>
    </select>
  </div>

  <div class="tree-container">
    <ul class="tree-root">
      <!-- Recursive task/step tree -->
      <li class="tree-node task-node" data-id="{{taskId}}">
        <div class="node-header">
          <span class="expand-toggle">▼</span>
          <span class="node-type">Task</span>
          <span class="node-title">{{taskTitle}}</span>
          <span class="node-progress">{{taskProgress}}%</span>
          <span class="node-status">{{taskStatus}}</span>
        </div>
        <ul class="node-children">
          <li class="tree-node step-node" data-id="{{stepId}}">
            <div class="node-header">
              <span class="node-type">Step</span>
              <span class="node-title">{{stepTitle}}</span>
              <span class="node-status">{{stepStatus}}</span>
              <span class="complexity-indicator">{{complexity}}</span>
            </div>
          </li>
        </ul>
      </li>
    </ul>
  </div>
</div>
```

### Collaboration Timeline

```html
<div class="collaboration-timeline">
  <h3>AI Collaboration History</h3>
  <div class="timeline-container">
    <div class="timeline-item" data-type="step-created">
      <div class="timeline-marker"></div>
      <div class="timeline-content">
        <h4>Step Created</h4>
        <p>{{stepTitle}} created in {{taskTitle}}</p>
        <span class="timestamp">{{timestamp}}</span>
      </div>
    </div>

    <div class="timeline-item" data-type="server-content">
      <div class="timeline-marker server"></div>
      <div class="timeline-content">
        <h4>Server AI Content</h4>
        <p>Initial content generated</p>
        <span class="complexity">Complexity: {{complexity}}</span>
        <span class="timestamp">{{timestamp}}</span>
      </div>
    </div>

    <div class="timeline-item" data-type="client-review">
      <div class="timeline-marker client"></div>
      <div class="timeline-content">
        <h4>Client AI Review</h4>
        <p>{{approved ? 'Approved' : 'Requested revision'}}</p>
        <span class="timestamp">{{timestamp}}</span>
      </div>
    </div>

    <div class="timeline-item" data-type="dispute">
      <div class="timeline-marker dispute"></div>
      <div class="timeline-content">
        <h4>Dispute Created</h4>
        <p>Max iterations reached, user resolution required</p>
        <span class="timestamp">{{timestamp}}</span>
      </div>
    </div>

    <div class="timeline-item" data-type="resolution">
      <div class="timeline-marker resolution"></div>
      <div class="timeline-content">
        <h4>User Resolution</h4>
        <p>Dispute resolved: {{resolutionType}}</p>
        <span class="timestamp">{{timestamp}}</span>
      </div>
    </div>
  </div>
</div>
```

---

## 3. Administrative Interface

### System Configuration

```html
<div class="admin-panel">
  <div class="admin-nav">
    <ul>
      <li><a href="#projects">Projects</a></li>
      <li><a href="#ai-config">AI Configuration</a></li>
      <li><a href="#system-stats">System Stats</a></li>
      <li><a href="#audit-logs">Audit Logs</a></li>
    </ul>
  </div>

  <div class="admin-content">
    <div id="ai-config" class="admin-section">
      <h2>AI Provider Configuration</h2>

      <div class="provider-config">
        <h3>Server AI Settings</h3>
        <form class="config-form">
          <div class="form-group">
            <label>Provider</label>
            <select name="serverProvider">
              <option value="openai">OpenAI</option>
              <option value="anthropic">Anthropic</option>
              <option value="local">Local Model</option>
            </select>
          </div>
          <div class="form-group">
            <label>Model</label>
            <input type="text" name="serverModel" value="gpt-4">
          </div>
          <div class="form-group">
            <label>Temperature</label>
            <input type="range" min="0" max="1" step="0.1"
                   name="serverTemperature" value="0.3">
          </div>
          <div class="form-group">
            <label>Max Tokens</label>
            <input type="number" name="serverMaxTokens" value="2000">
          </div>
        </form>
      </div>

      <div class="provider-config">
        <h3>Client AI Settings</h3>
        <form class="config-form">
          <!-- Similar configuration for client AI -->
        </form>
      </div>

      <div class="global-settings">
        <h3>Global AI Settings</h3>
        <form class="config-form">
          <div class="form-group">
            <label>Default Max Iterations</label>
            <input type="number" name="defaultMaxIterations" value="3">
          </div>
          <div class="form-group">
            <label>Default Complexity Threshold</label>
            <input type="range" min="0" max="1" step="0.1"
                   name="defaultComplexityThreshold" value="0.7">
          </div>
          <div class="form-group">
            <label>Timeout (minutes)</label>
            <input type="number" name="aiTimeout" value="5">
          </div>
        </form>
      </div>
    </div>

    <div id="system-stats" class="admin-section">
      <h2>System Statistics</h2>

      <div class="stats-grid">
        <div class="stat-card">
          <h3>Projects</h3>
          <div class="stat-value">{{totalProjects}}</div>
          <div class="stat-detail">{{activeProjects}} active</div>
        </div>

        <div class="stat-card">
          <h3>AI Collaboration</h3>
          <div class="stat-value">{{agreementRate}}%</div>
          <div class="stat-detail">Agreement rate</div>
        </div>

        <div class="stat-card">
          <h3>Disputes</h3>
          <div class="stat-value">{{pendingDisputes}}</div>
          <div class="stat-detail">Pending resolution</div>
        </div>

        <div class="stat-card">
          <h3>Performance</h3>
          <div class="stat-value">{{avgResponseTime}}s</div>
          <div class="stat-detail">Avg AI response</div>
        </div>
      </div>

      <div class="charts-section">
        <div class="chart-container">
          <h3>Daily Activity</h3>
          <canvas id="activityChart"></canvas>
        </div>

        <div class="chart-container">
          <h3>Complexity Distribution</h3>
          <canvas id="complexityChart"></canvas>
        </div>
      </div>
    </div>
  </div>
</div>
```

### Audit Logs

```html
<div id="audit-logs" class="admin-section">
  <h2>Audit Logs</h2>

  <div class="log-filters">
    <div class="filter-group">
      <label>Entity Type</label>
      <select name="entityType">
        <option value="">All</option>
        <option value="project">Project</option>
        <option value="task">Task</option>
        <option value="step">Step</option>
        <option value="dispute">Dispute</option>
      </select>
    </div>

    <div class="filter-group">
      <label>Action</label>
      <select name="action">
        <option value="">All</option>
        <option value="created">Created</option>
        <option value="updated">Updated</option>
        <option value="promoted">Promoted</option>
        <option value="disputed">Disputed</option>
        <option value="resolved">Resolved</option>
      </select>
    </div>

    <div class="filter-group">
      <label>Date Range</label>
      <input type="date" name="startDate">
      <input type="date" name="endDate">
    </div>

    <button class="apply-filters">Apply Filters</button>
  </div>

  <div class="log-table">
    <table>
      <thead>
        <tr>
          <th>Timestamp</th>
          <th>Entity</th>
          <th>Action</th>
          <th>User/Agent</th>
          <th>Details</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>{{timestamp}}</td>
          <td>{{entityType}}/{{entityId}}</td>
          <td>{{action}}</td>
          <td>{{userId || aiAgent}}</td>
          <td>
            <button class="view-details" data-log-id="{{logId}}">
              View Details
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</div>
```

---

## 4. API Integration Points

### REST API Endpoints

```typescript
// Dispute Resolution API
interface DisputeAPI {
  // Get pending disputes
  GET /api/disputes?projectId={id}&limit={n}

  // Get specific dispute details
  GET /api/disputes/{disputeId}

  // Resolve dispute
  POST /api/disputes/{disputeId}/resolve
  Body: {
    resolution: 'server' | 'client' | 'custom' | 'hybrid',
    customContent?: string,
    userReasoning?: string
  }

  // Get dispute statistics
  GET /api/disputes/stats?projectId={id}
}

// Project Dashboard API
interface ProjectAPI {
  // Get project dashboard data
  GET /api/projects/{projectId}/dashboard

  // Get project workflow structure
  GET /api/projects/{projectId}/workflow

  // Get collaboration timeline
  GET /api/projects/{projectId}/timeline?limit={n}

  // Update project settings
  PUT /api/projects/{projectId}/settings
  Body: {
    complexityThreshold?: number,
    maxIterations?: number
  }
}

// Administrative API
interface AdminAPI {
  // Get system statistics
  GET /api/admin/stats

  // Get audit logs
  GET /api/admin/audit-logs?filters={...}

  // Update AI configuration
  PUT /api/admin/ai-config
  Body: {
    serverAI: AIConfig,
    clientAI: AIConfig,
    globalSettings: GlobalConfig
  }
}
```

### WebSocket Events

```typescript
// Real-time updates for dashboard
interface WebSocketEvents {
  // Project progress updates
  'project:progress': {
    projectId: string;
    progress: number;
    completedItem: {
      type: 'task' | 'step';
      id: string;
      title: string;
    };
  };

  // New dispute created
  'dispute:created': {
    disputeId: string;
    projectId: string;
    stepTitle: string;
    urgency: 'low' | 'medium' | 'high';
  };

  // Dispute resolved
  'dispute:resolved': {
    disputeId: string;
    projectId: string;
    resolution: string;
    resolvedBy: string;
  };

  // AI collaboration events
  'collaboration:update': {
    stepId: string;
    projectId: string;
    event: 'server_content' | 'client_review' | 'agreement' | 'iteration';
    details: any;
  };
}
```

### Mobile-Responsive Design

```css
/* Mobile-first responsive design */
@media (max-width: 768px) {
  .dispute-resolution {
    .content-comparison {
      flex-direction: column;
    }

    .server-version,
    .client-version {
      width: 100%;
      margin-bottom: 1rem;
    }

    .resolution-options {
      .option-buttons {
        flex-direction: column;
        gap: 0.5rem;
      }
    }
  }

  .project-dashboard {
    .dashboard-grid {
      grid-template-columns: 1fr;
      gap: 1rem;
    }
  }

  .task-tree {
    .tree-controls {
      flex-wrap: wrap;
      gap: 0.5rem;
    }
  }
}

/* Tablet design */
@media (min-width: 769px) and (max-width: 1024px) {
  .dashboard-grid {
    grid-template-columns: 1fr 1fr;
  }
}

/* Desktop design */
@media (min-width: 1025px) {
  .dashboard-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}
```

## Accessibility Features

### ARIA Labels and Roles

```html
<!-- Dispute resolution accessibility -->
<div class="dispute-resolution" role="main" aria-labelledby="dispute-title">
  <h2 id="dispute-title">Resolve Content Dispute</h2>

  <div class="content-comparison" role="region" aria-label="Content comparison">
    <div class="server-version" role="article" aria-labelledby="server-heading">
      <h3 id="server-heading">Server AI Version</h3>
      <div class="content-preview" role="document" aria-label="Server AI content">
        <!-- Content -->
      </div>
    </div>

    <div class="client-version" role="article" aria-labelledby="client-heading">
      <h3 id="client-heading">Client AI Version</h3>
      <div class="content-preview" role="document" aria-label="Client AI content">
        <!-- Content -->
      </div>
    </div>
  </div>

  <div class="resolution-options" role="group" aria-labelledby="resolution-heading">
    <h3 id="resolution-heading">Resolution Options</h3>
    <div class="option-buttons" role="radiogroup" aria-required="true">
      <button role="radio" aria-checked="false" aria-describedby="server-desc">
        Use Server Version
      </button>
      <div id="server-desc" class="sr-only">
        Accept the Server AI's version of the content
      </div>
    </div>
  </div>
</div>
```

### Keyboard Navigation

```typescript
// Keyboard navigation for dispute resolution
class DisputeKeyboardHandler {
  constructor() {
    this.setupKeyboardHandlers();
  }

  setupKeyboardHandlers() {
    document.addEventListener('keydown', (e) => {
      switch (e.key) {
        case '1':
          if (e.ctrlKey) this.selectResolution('server');
          break;
        case '2':
          if (e.ctrlKey) this.selectResolution('client');
          break;
        case '3':
          if (e.ctrlKey) this.selectResolution('hybrid');
          break;
        case '4':
          if (e.ctrlKey) this.selectResolution('custom');
          break;
        case 'Enter':
          if (e.target.classList.contains('resolution-btn')) {
            this.submitResolution();
          }
          break;
      }
    });
  }
}
```

---

*Next: [Implementation Guide](./07-implementation-guide.md)*
