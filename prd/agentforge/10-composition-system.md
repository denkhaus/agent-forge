# Composition System

## Overview

AgentForge's composition system enables users to create, manage, and deploy complete AI agent solutions by combining tools, prompts, and agents from distributed repositories. Compositions provide a declarative way to define complex agent behaviors with precise dependency management.

## Composition Architecture

### Composition Hierarchy
```
Composition (Top Level)
├── Metadata (name, version, description)
├── Dependencies
│   ├── Tools (executable functions)
│   ├── Prompts (templates and instructions)
│   └── Agents (complete agent configurations)
├── Configuration
│   ├── Global Settings
│   ├── Component Overrides
│   └── Environment Variables
├── Deployment
│   ├── Runtime Configuration
│   ├── Scaling Parameters
│   └── Health Checks
└── Instances (running deployments)
```

## Composition Definition

### Complete Composition Specification
```yaml
apiVersion: "forge.dev/v1"
kind: "Composition"
metadata:
  name: "enterprise-sales-stack"
  version: "2.1.0"
  description: "Complete enterprise sales automation system with CRM integration"
  author: "sales-automation-team"
  license: "MIT"
  
  tags: ["sales", "enterprise", "crm", "automation"]
  categories: ["sales", "customer-service"]
  
  # Composition-specific metadata
  target_environment: "production"
  complexity_level: "advanced"
  estimated_cost: "medium"
  
spec:
  # Component dependencies with exact versioning
  dependencies:
    tools:
      - name: "salesforce-lookup"
        repository: "github.com/company/crm-tools"
        version: "^2.1.0"
        commit: "a1b2c3d4e5f6789012345678901234567890abcd"
        required: true
        purpose: "Customer data retrieval and CRM integration"
        
      - name: "email-composer"
        repository: "github.com/tools/communication"
        version: "~1.3.0"
        commit: "b2c3d4e5f6789012345678901234567890abcdef"
        required: true
        purpose: "Email generation and sending"
        
      - name: "calendar-scheduler"
        repository: "github.com/tools/productivity"
        version: ">=2.0.0"
        commit: "c3d4e5f6789012345678901234567890abcdef12"
        required: false
        purpose: "Meeting scheduling and calendar management"
        
      - name: "document-generator"
        repository: "github.com/tools/documents"
        version: "latest"
        commit: "d4e5f6789012345678901234567890abcdef1234"
        required: false
        purpose: "Proposal and contract generation"
    
    prompts:
      - name: "sales-system-v3"
        repository: "github.com/ai-prompts/sales-templates"
        version: ">=3.0.0"
        commit: "e5f6789012345678901234567890abcdef123456"
        role: "system"
        required: true
        purpose: "Primary sales assistant system prompt"
        
      - name: "objection-handling"
        repository: "github.com/ai-prompts/sales-templates"
        version: "^2.1.0"
        commit: "f6789012345678901234567890abcdef1234567"
        role: "task"
        required: true
        purpose: "Objection handling strategies"
        
      - name: "closing-techniques"
        repository: "github.com/ai-prompts/sales-templates"
        version: "~1.5.0"
        commit: "g789012345678901234567890abcdef12345678"
        role: "task"
        required: false
        purpose: "Sales closing assistance"
    
    agents:
      - name: "lead-qualifier"
        repository: "github.com/agents/specialized"
        version: "^1.2.0"
        commit: "h89012345678901234567890abcdef123456789"
        required: true
        purpose: "Initial lead qualification and scoring"
        
      - name: "sales-specialist"
        repository: "github.com/agents/specialized"
        version: "^1.5.0"
        commit: "i9012345678901234567890abcdef1234567890"
        required: true
        purpose: "Primary sales agent for customer interactions"
  
  # Global configuration
  configuration:
    # LLM settings
    llm:
      primary_provider: "openai"
      primary_model: "gpt-4"
      fallback_provider: "anthropic"
      fallback_model: "claude-3-sonnet"
      
      # Global LLM parameters
      temperature: 0.7
      max_tokens: 2048
      top_p: 0.9
      
    # Runtime behavior
    runtime:
      execution_mode: "agent"
      max_iterations: 5
      timeout: 300  # seconds
      memory_limit: "512MB"
      
    # Monitoring and observability
    monitoring:
      enable_metrics: true
      enable_tracing: true
      log_level: "info"
      health_check_interval: 30  # seconds
      
    # Security settings
    security:
      enable_sandboxing: true
      allowed_domains: ["*.salesforce.com", "*.hubspot.com"]
      secret_management: "env_vars"  # env_vars, vault, k8s_secrets
  
  # Component-specific overrides
  overrides:
    tools:
      salesforce-lookup:
        config:
          timeout: 45
          retry_count: 3
          cache_ttl: 600
        environment:
          SALESFORCE_API_VERSION: "58.0"
          
      email-composer:
        config:
          template_engine: "handlebars"
          max_recipients: 50
        
    prompts:
      sales-system-v3:
        variables:
          company_name: "Enterprise Corp"
          sales_methodology: "consultative"
          target_audience: "enterprise customers"
          
    agents:
      sales-specialist:
        llm:
          temperature: 0.8  # Override global setting
          max_tokens: 3000
        behavior:
          max_iterations: 7
          verbosity: "detailed"
  
  # Environment-specific configurations
  environments:
    development:
      llm:
        primary_provider: "openai"
        primary_model: "gpt-3.5-turbo"  # Cheaper for dev
      runtime:
        timeout: 60
        memory_limit: "256MB"
      monitoring:
        log_level: "debug"
        
    staging:
      llm:
        primary_provider: "openai"
        primary_model: "gpt-4"
      runtime:
        timeout: 180
        memory_limit: "512MB"
      monitoring:
        log_level: "info"
        
    production:
      llm:
        primary_provider: "openai"
        primary_model: "gpt-4"
        fallback_provider: "anthropic"
      runtime:
        timeout: 300
        memory_limit: "1GB"
      monitoring:
        log_level: "warn"
        enable_metrics: true
        enable_tracing: true
  
  # Deployment configuration
  deployment:
    # Scaling configuration
    scaling:
      min_instances: 1
      max_instances: 10
      target_cpu_utilization: 70
      target_memory_utilization: 80
      scale_up_cooldown: 300    # seconds
      scale_down_cooldown: 600  # seconds
      
    # Health checks
    health_checks:
      liveness_probe:
        endpoint: "/health/live"
        initial_delay: 30
        period: 10
        timeout: 5
        failure_threshold: 3
        
      readiness_probe:
        endpoint: "/health/ready"
        initial_delay: 10
        period: 5
        timeout: 3
        failure_threshold: 2
    
    # Resource requirements
    resources:
      requests:
        cpu: "500m"
        memory: "512Mi"
      limits:
        cpu: "2000m"
        memory: "2Gi"
    
    # Networking
    networking:
      ports:
        - name: "http"
          port: 8080
          protocol: "TCP"
        - name: "metrics"
          port: 9090
          protocol: "TCP"
      
      ingress:
        enabled: true
        host: "sales-agent.company.com"
        tls: true
        annotations:
          nginx.ingress.kubernetes.io/rate-limit: "100"
  
  # Workflow definitions
  workflows:
    lead_processing:
      description: "Complete lead processing workflow"
      steps:
        - name: "qualify_lead"
          agent: "lead-qualifier"
          input: "${workflow.input}"
          output: "qualification_result"
          
        - name: "initial_contact"
          agent: "sales-specialist"
          input: 
            lead_data: "${steps.qualify_lead.output}"
            action: "initial_contact"
          condition: "${steps.qualify_lead.output.qualified == true}"
          output: "contact_result"
          
        - name: "schedule_follow_up"
          agent: "sales-specialist"
          input:
            contact_result: "${steps.initial_contact.output}"
            action: "schedule_follow_up"
          condition: "${steps.initial_contact.output.interested == true}"
          output: "follow_up_scheduled"
    
    objection_handling:
      description: "Handle customer objections"
      steps:
        - name: "analyze_objection"
          agent: "sales-specialist"
          input: "${workflow.input}"
          output: "objection_analysis"
          
        - name: "generate_response"
          agent: "sales-specialist"
          input:
            objection: "${workflow.input.objection}"
            analysis: "${steps.analyze_objection.output}"
            action: "handle_objection"
          output: "objection_response"
```

## Composition Operations

### Creation and Management
```bash
# Create new composition
forge compose create enterprise-sales-stack
forge compose create enterprise-sales-stack --from-template sales-automation
forge compose create enterprise-sales-stack --interactive

# Create from existing components
forge compose create my-stack \
  --agent github.com/agents/sales-specialist@v1.5.0 \
  --tool github.com/company/crm-tools/salesforce-lookup@v2.1.0 \
  --prompt github.com/prompts/sales-system@v3.0.0

# Import composition from file
forge compose import ./enterprise-sales-stack.yaml
forge compose import --from-url https://github.com/company/compositions/raw/main/sales-stack.yaml
```

### Configuration Management
```bash
# Edit composition
forge compose edit enterprise-sales-stack
forge compose edit enterprise-sales-stack --add-tool email-composer
forge compose edit enterprise-sales-stack --remove-tool old-tool
forge compose edit enterprise-sales-stack --update-agent sales-specialist@v1.6.0

# Configure environments
forge compose config enterprise-sales-stack --environment production
forge compose config enterprise-sales-stack --set llm.temperature=0.8
forge compose config enterprise-sales-stack --set-env SALESFORCE_URL=https://prod.salesforce.com

# Override component settings
forge compose override enterprise-sales-stack \
  --component salesforce-lookup \
  --config timeout=60

# Validate composition
forge compose validate enterprise-sales-stack
forge compose validate enterprise-sales-stack --environment production
forge compose validate --all
```

### Dependency Management
```bash
# Show composition dependencies
forge compose deps enterprise-sales-stack
forge compose deps enterprise-sales-stack --tree
forge compose deps enterprise-sales-stack --outdated

# Update dependencies
forge compose update enterprise-sales-stack
forge compose update enterprise-sales-stack --component salesforce-lookup
forge compose update enterprise-sales-stack --check-breaking

# Lock dependencies
forge compose lock enterprise-sales-stack
forge compose lock enterprise-sales-stack --output enterprise-sales-lock.yaml

# Install dependencies
forge compose install enterprise-sales-stack
forge compose install enterprise-sales-stack --from-lock
```

## Deployment and Runtime

### Deployment Operations
```bash
# Deploy composition
forge compose deploy enterprise-sales-stack
forge compose deploy enterprise-sales-stack --environment production
forge compose deploy enterprise-sales-stack --instance sales-prod-1

# Deploy with overrides
forge compose deploy enterprise-sales-stack \
  --environment production \
  --set scaling.min_instances=3 \
  --set resources.requests.memory=1Gi

# Rolling deployment
forge compose deploy enterprise-sales-stack --strategy rolling
forge compose deploy enterprise-sales-stack --strategy blue-green
forge compose deploy enterprise-sales-stack --strategy canary --canary-weight 10
```

### Instance Management
```bash
# List instances
forge compose instances
forge compose instances enterprise-sales-stack
forge compose instances --environment production

# Instance details
forge compose instance info sales-prod-1
forge compose instance logs sales-prod-1
forge compose instance logs sales-prod-1 --follow --since 1h

# Instance operations
forge compose instance start sales-prod-1
forge compose instance stop sales-prod-1
forge compose instance restart sales-prod-1
forge compose instance scale sales-prod-1 --replicas 5

# Health checks
forge compose instance health sales-prod-1
forge compose instance health --all
```

### Monitoring and Observability
```bash
# Composition metrics
forge compose metrics enterprise-sales-stack
forge compose metrics enterprise-sales-stack --instance sales-prod-1
forge compose metrics enterprise-sales-stack --since 1h

# Performance monitoring
forge compose performance enterprise-sales-stack
# Metrics:
# - Average response time: 2.3s
# - Success rate: 98.5%
# - Throughput: 45 requests/minute
# - Error rate: 1.5%
# - Resource utilization: CPU 45%, Memory 67%

# Alerts and notifications
forge compose alerts enterprise-sales-stack
forge compose alerts enterprise-sales-stack --configure
```

## Composition Templates

### Template System
```yaml
# templates/sales-automation.yaml
apiVersion: "forge.dev/v1"
kind: "CompositionTemplate"
metadata:
  name: "sales-automation"
  version: "1.0.0"
  description: "Template for sales automation compositions"
  
template:
  metadata:
    name: "{{ .name }}"
    description: "{{ .description | default "Sales automation system" }}"
    
  spec:
    dependencies:
      tools:
        - name: "crm-lookup"
          repository: "{{ .crm_repository | default "github.com/company/crm-tools" }}"
          version: "{{ .crm_version | default "^2.0.0" }}"
          
      prompts:
        - name: "sales-system"
          repository: "{{ .prompt_repository | default "github.com/prompts/sales" }}"
          version: "{{ .prompt_version | default "^3.0.0" }}"
          
    configuration:
      llm:
        primary_provider: "{{ .llm_provider | default "openai" }}"
        primary_model: "{{ .llm_model | default "gpt-4" }}"
        temperature: {{ .temperature | default 0.7 }}

parameters:
  - name: "name"
    type: "string"
    required: true
    description: "Composition name"
    
  - name: "description"
    type: "string"
    required: false
    description: "Composition description"
    
  - name: "crm_repository"
    type: "string"
    required: false
    description: "CRM tools repository"
    
  - name: "llm_provider"
    type: "string"
    enum: ["openai", "anthropic", "google"]
    default: "openai"
    description: "Primary LLM provider"
```

### Template Usage
```bash
# List available templates
forge template list
forge template list --category sales

# Create composition from template
forge compose create my-sales-stack --template sales-automation
forge compose create my-sales-stack --template sales-automation \
  --param name="My Sales Stack" \
  --param llm_provider=anthropic \
  --param temperature=0.8

# Interactive template usage
forge compose create --template sales-automation --interactive
# Template: sales-automation
# Name: My Sales Stack
# Description: [Sales automation system] 
# CRM Repository: [github.com/company/crm-tools]
# LLM Provider: [openai] anthropic, google
# Temperature: [0.7]
```

## Composition Lifecycle

### Development Lifecycle
```go
type CompositionLifecycle struct {
    db          *database.Client
    deployer    *Deployer
    monitor     *Monitor
    logger      *zap.Logger
}

func (cl *CompositionLifecycle) CreateComposition(ctx context.Context, req CreateCompositionRequest) (*Composition, error) {
    // Validate composition definition
    if err := cl.validateComposition(req.Definition); err != nil {
        return nil, fmt.Errorf("invalid composition: %w", err)
    }
    
    // Resolve dependencies
    resolved, err := cl.resolveDependencies(ctx, req.Definition.Dependencies)
    if err != nil {
        return nil, fmt.Errorf("dependency resolution failed: %w", err)
    }
    
    // Create composition record
    composition := &Composition{
        Name:        req.Definition.Metadata.Name,
        Version:     req.Definition.Metadata.Version,
        Definition:  req.Definition,
        Resolved:    resolved,
        Status:      CompositionStatusCreated,
        CreatedAt:   time.Now(),
    }
    
    if err := cl.db.CreateComposition(ctx, composition); err != nil {
        return nil, fmt.Errorf("failed to create composition: %w", err)
    }
    
    cl.logger.Info("Composition created",
        zap.String("name", composition.Name),
        zap.String("version", composition.Version))
    
    return composition, nil
}

func (cl *CompositionLifecycle) DeployComposition(ctx context.Context, req DeployCompositionRequest) (*Deployment, error) {
    composition, err := cl.db.GetComposition(ctx, req.CompositionName)
    if err != nil {
        return nil, err
    }
    
    // Prepare deployment
    deployment := &Deployment{
        CompositionID: composition.ID,
        Environment:   req.Environment,
        Configuration: cl.mergeConfiguration(composition.Definition.Configuration, req.Overrides),
        Status:        DeploymentStatusPending,
        CreatedAt:     time.Now(),
    }
    
    // Deploy using strategy
    switch req.Strategy {
    case DeploymentStrategyRolling:
        return cl.deployRolling(ctx, deployment)
    case DeploymentStrategyBlueGreen:
        return cl.deployBlueGreen(ctx, deployment)
    case DeploymentStrategyCanary:
        return cl.deployCanary(ctx, deployment, req.CanaryWeight)
    default:
        return cl.deployDirect(ctx, deployment)
    }
}
```

### Monitoring and Health
```go
type CompositionMonitor struct {
    metrics     *metrics.Client
    alerts      *alerts.Manager
    healthCheck *health.Checker
}

func (cm *CompositionMonitor) MonitorComposition(ctx context.Context, compositionID string) error {
    composition, err := cm.getComposition(ctx, compositionID)
    if err != nil {
        return err
    }
    
    // Start health monitoring
    go cm.healthMonitorLoop(ctx, composition)
    
    // Start metrics collection
    go cm.metricsCollectionLoop(ctx, composition)
    
    // Start alert monitoring
    go cm.alertMonitorLoop(ctx, composition)
    
    return nil
}

func (cm *CompositionMonitor) healthMonitorLoop(ctx context.Context, composition *Composition) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            health, err := cm.healthCheck.CheckComposition(ctx, composition.ID)
            if err != nil {
                cm.logger.Error("Health check failed", zap.Error(err))
                continue
            }
            
            if health.Status != HealthStatusHealthy {
                cm.alerts.TriggerAlert(ctx, &Alert{
                    Type:          AlertTypeHealthCheck,
                    CompositionID: composition.ID,
                    Severity:      AlertSeverityWarning,
                    Message:       fmt.Sprintf("Composition health: %s", health.Status),
                })
            }
        }
    }
}
```

## Advanced Features

### Composition Inheritance
```yaml
# Base composition
apiVersion: "forge.dev/v1"
kind: "Composition"
metadata:
  name: "base-sales-stack"
  version: "1.0.0"
  
spec:
  dependencies:
    tools:
      - name: "basic-crm"
        repository: "github.com/tools/crm"
        version: "^1.0.0"
        
  configuration:
    llm:
      primary_provider: "openai"
      primary_model: "gpt-3.5-turbo"

---
# Extended composition
apiVersion: "forge.dev/v1"
kind: "Composition"
metadata:
  name: "enterprise-sales-stack"
  version: "2.0.0"
  
spec:
  # Inherit from base composition
  extends: "base-sales-stack@1.0.0"
  
  dependencies:
    tools:
      # Override base tool
      - name: "basic-crm"
        repository: "github.com/enterprise/crm-tools"
        version: "^2.0.0"
        
      # Add new tools
      - name: "advanced-analytics"
        repository: "github.com/tools/analytics"
        version: "^1.5.0"
        
  configuration:
    llm:
      # Override base LLM
      primary_model: "gpt-4"
```

### Multi-Environment Compositions
```bash
# Deploy to multiple environments
forge compose deploy enterprise-sales-stack --environments dev,staging,prod

# Environment-specific configurations
forge compose config enterprise-sales-stack \
  --environment dev \
  --set llm.primary_model=gpt-3.5-turbo

forge compose config enterprise-sales-stack \
  --environment prod \
  --set scaling.min_instances=5

# Promote between environments
forge compose promote enterprise-sales-stack --from staging --to production
```

### Composition Versioning
```bash
# Version management
forge compose version enterprise-sales-stack
forge compose version enterprise-sales-stack --bump minor
forge compose version enterprise-sales-stack --set 2.1.0

# Version comparison
forge compose diff enterprise-sales-stack@2.0.0 enterprise-sales-stack@2.1.0
forge compose changelog enterprise-sales-stack

# Rollback
forge compose rollback enterprise-sales-stack --to-version 2.0.0
forge compose rollback enterprise-sales-stack --to-deployment deployment-123
```