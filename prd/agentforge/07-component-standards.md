# Component Standards

## Overview

AgentForge component standards define the structure, metadata, and interfaces for tools, prompts, and agents to ensure interoperability, discoverability, and quality across the ecosystem.

## Universal Component Structure

### Base Component Schema
```yaml
# Universal fields for all component types
apiVersion: "forge.dev/v1"
kind: "Tool" | "Prompt" | "Agent"
metadata:
  name: string              # Unique identifier within repository
  version: string           # Semantic version (e.g., "2.1.0")
  description: string       # Human-readable description
  author: string            # Author or organization
  license: string           # SPDX license identifier
  homepage: string?         # Project homepage URL
  documentation: string?    # Documentation URL
  repository: string?       # Source code repository
  
  # Categorization
  tags: string[]           # Searchable tags
  categories: string[]     # Hierarchical categories
  keywords: string[]       # Search keywords
  
  # Quality indicators
  stability: "experimental" | "beta" | "stable" | "deprecated"
  maturity: "alpha" | "beta" | "stable" | "mature"
  
  # Compatibility
  forge_version: string    # Minimum AgentForge version
  platforms: string[]     # Supported platforms
  
spec:
  # Component-specific specification
  # Defined per component type below
```

## Tool Components

### Tool Specification
```yaml
apiVersion: "forge.dev/v1"
kind: "Tool"
metadata:
  name: "salesforce-lookup"
  version: "2.1.0"
  description: "Query Salesforce CRM for customer data with advanced filtering"
  author: "enterprise-tools-team"
  license: "MIT"
  homepage: "https://tools.company.com/salesforce-lookup"
  documentation: "https://docs.company.com/tools/salesforce-lookup"
  
  tags: ["crm", "salesforce", "customer-data", "enterprise"]
  categories: ["crm", "data-retrieval"]
  keywords: ["salesforce", "customer", "lookup", "query", "crm"]
  
  stability: "stable"
  maturity: "mature"
  forge_version: ">=0.2.0"
  platforms: ["linux/amd64", "linux/arm64", "darwin/amd64", "darwin/arm64"]

spec:
  # Tool type and runtime
  type: "mcp-server"        # mcp-server, function, webhook, grpc
  runtime: "go"             # go, python, node, rust, etc.
  
  # Execution configuration
  execution:
    binary: "./bin/salesforce-server"
    args: ["--config", "${CONFIG_FILE}"]
    env_file: ".env"
    working_directory: "."
    timeout: 30             # seconds
    memory_limit: "256MB"
    cpu_limit: "0.5"
  
  # Input/Output schema
  schema:
    input:
      type: "object"
      properties:
        query_type:
          type: "string"
          enum: ["account", "contact", "opportunity", "lead"]
          description: "Type of Salesforce object to query"
        search_term:
          type: "string"
          minLength: 1
          maxLength: 255
          description: "Search term for the query"
        fields:
          type: "array"
          items:
            type: "string"
          description: "Fields to retrieve"
          default: ["Id", "Name"]
        limit:
          type: "integer"
          minimum: 1
          maximum: 1000
          default: 10
          description: "Maximum number of results"
        offset:
          type: "integer"
          minimum: 0
          default: 0
          description: "Number of results to skip"
      required: ["query_type", "search_term"]
      additionalProperties: false
    
    output:
      type: "object"
      properties:
        results:
          type: "array"
          items:
            type: "object"
            additionalProperties: true
          description: "Query results"
        total_count:
          type: "integer"
          description: "Total number of matching records"
        has_more:
          type: "boolean"
          description: "Whether more results are available"
        next_offset:
          type: "integer"
          description: "Offset for next page of results"
      required: ["results", "total_count", "has_more"]
  
  # Configuration
  configuration:
    environment:
      - name: "SALESFORCE_URL"
        required: true
        description: "Salesforce instance URL"
        example: "https://company.my.salesforce.com"
      - name: "SALESFORCE_TOKEN"
        required: true
        secret: true
        description: "Salesforce API token or session ID"
      - name: "SALESFORCE_VERSION"
        required: false
        default: "58.0"
        description: "Salesforce API version"
    
    parameters:
      timeout:
        type: "integer"
        default: 30
        minimum: 5
        maximum: 300
        description: "Request timeout in seconds"
      retry_count:
        type: "integer"
        default: 3
        minimum: 0
        maximum: 10
        description: "Number of retry attempts"
      cache_ttl:
        type: "integer"
        default: 300
        minimum: 0
        description: "Cache TTL in seconds (0 = no cache)"
  
  # Capabilities and features
  capabilities:
    - "customer_lookup"
    - "crm_integration"
    - "data_retrieval"
    - "pagination"
    - "caching"
  
  # Dependencies
  dependencies:
    external:
      - name: "salesforce-api"
        version: ">=58.0"
        description: "Salesforce REST API"
    
    forge_components: []
    
    system:
      - "network.outbound"
      - "env.read"
  
  # Security
  security:
    permissions:
      - "network.outbound"
      - "env.read"
    
    sandboxed: true
    trusted_domains:
      - "*.salesforce.com"
      - "*.force.com"
    
    secrets:
      - name: "SALESFORCE_TOKEN"
        description: "Salesforce authentication token"
        rotation_period: "90d"
  
  # Quality and testing
  quality:
    test_coverage: 95
    has_integration_tests: true
    has_performance_tests: true
    benchmark_results:
      avg_response_time: "150ms"
      max_response_time: "500ms"
      throughput: "100 req/s"
  
  # Monitoring and observability
  monitoring:
    health_check:
      endpoint: "/health"
      interval: "30s"
      timeout: "5s"
    
    metrics:
      - name: "requests_total"
        type: "counter"
        description: "Total number of requests"
      - name: "request_duration"
        type: "histogram"
        description: "Request duration in seconds"
      - name: "active_connections"
        type: "gauge"
        description: "Number of active connections"
    
    logs:
      level: "info"
      format: "json"
      fields: ["timestamp", "level", "message", "component", "request_id"]
```

### Tool Implementation Requirements

#### MCP Server Tools
```go
// Required interface for MCP server tools
type MCPServerTool interface {
    // Server lifecycle
    Start(ctx context.Context, config Config) error
    Stop(ctx context.Context) error
    Health(ctx context.Context) error
    
    // Tool execution
    Execute(ctx context.Context, input string) (string, error)
    
    // Metadata
    Name() string
    Description() string
    Schema() *Schema
}
```

#### Function Tools
```go
// Required interface for function-based tools
type FunctionTool interface {
    // Direct execution
    Call(ctx context.Context, input string) (string, error)
    
    // Metadata
    Name() string
    Description() string
    Schema() *Schema
    
    // Configuration
    Configure(config map[string]interface{}) error
    Validate() error
}
```

## Prompt Components

### Prompt Specification
```yaml
apiVersion: "forge.dev/v1"
kind: "Prompt"
metadata:
  name: "sales-system-v3"
  version: "3.2.0"
  description: "Advanced sales assistant system prompt with objection handling"
  author: "ai-prompts-team"
  license: "CC-BY-4.0"
  
  tags: ["sales", "system-prompt", "objection-handling", "crm"]
  categories: ["sales", "customer-service"]
  keywords: ["sales", "assistant", "crm", "objection", "closing"]
  
  stability: "stable"
  maturity: "mature"
  forge_version: ">=0.2.0"

spec:
  # Prompt type and category
  type: "system"            # system, user, assistant, function
  category: "sales"         # Domain-specific category
  
  # Template configuration
  template:
    engine: "go-template"   # go-template, jinja2, handlebars, mustache
    file: "./template.txt"  # Template file path
    syntax_version: "1.0"   # Template syntax version
  
  # Variables and parameters
  variables:
    - name: "company_name"
      type: "string"
      required: true
      description: "Name of the company"
      example: "Acme Corporation"
      validation:
        min_length: 1
        max_length: 100
    
    - name: "product_catalog"
      type: "array"
      required: false
      description: "Available products and services"
      items:
        type: "object"
        properties:
          name: { type: "string" }
          description: { type: "string" }
          price: { type: "number" }
      default: []
    
    - name: "sales_methodology"
      type: "string"
      required: false
      default: "consultative"
      enum: ["consultative", "solution", "challenger", "spin"]
      description: "Sales methodology to follow"
    
    - name: "target_audience"
      type: "string"
      required: false
      description: "Primary target audience"
      example: "enterprise customers"
    
    - name: "objection_strategies"
      type: "object"
      required: false
      description: "Strategies for handling common objections"
      properties:
        price_objections: { type: "string" }
        timing_objections: { type: "string" }
        authority_objections: { type: "string" }
      default:
        price_objections: "Focus on value and ROI"
        timing_objections: "Highlight urgency and opportunity cost"
        authority_objections: "Identify decision makers"
  
  # Content validation
  validation:
    max_length: 4000
    min_length: 100
    required_sections:
      - "role_definition"
      - "capabilities"
      - "guidelines"
      - "objection_handling"
    
    forbidden_content:
      - "inappropriate_language"
      - "personal_information"
      - "competitor_disparagement"
    
    tone_requirements:
      - "professional"
      - "helpful"
      - "confident"
  
  # Localization support
  localization:
    supported_languages: ["en", "es", "fr", "de", "ja"]
    default_language: "en"
    translation_files:
      es: "./translations/es.yaml"
      fr: "./translations/fr.yaml"
      de: "./translations/de.yaml"
      ja: "./translations/ja.yaml"
  
  # Usage context
  context:
    use_cases:
      - "lead_qualification"
      - "objection_handling"
      - "product_demonstration"
      - "closing_techniques"
    
    target_roles:
      - "sales_representative"
      - "account_manager"
      - "sales_engineer"
    
    interaction_types:
      - "phone_calls"
      - "video_meetings"
      - "email_communication"
      - "chat_support"
  
  # Performance characteristics
  performance:
    token_efficiency: 0.85    # Tokens used / total tokens
    response_quality: 4.7     # Average quality rating (1-5)
    consistency_score: 0.92   # Response consistency
    
    benchmarks:
      - metric: "objection_resolution_rate"
        value: 0.78
        description: "Rate of successful objection handling"
      - metric: "conversation_flow_score"
        value: 4.5
        description: "Natural conversation flow rating"
  
  # Integration requirements
  integration:
    compatible_llms:
      - provider: "openai"
        models: ["gpt-4", "gpt-3.5-turbo"]
        optimal_temperature: 0.7
      - provider: "anthropic"
        models: ["claude-3-sonnet", "claude-3-haiku"]
        optimal_temperature: 0.6
      - provider: "google"
        models: ["gemini-pro"]
        optimal_temperature: 0.8
    
    required_context_window: 8000
    recommended_max_tokens: 2048
```

### Prompt Template Format
```go
// template.txt - Go template example
You are {{.company_name}}'s AI sales assistant, specialized in {{.sales_methodology}} selling.

## Your Role
You are an expert sales professional representing {{.company_name}}. Your primary goal is to help potential customers understand how our solutions can address their specific needs and challenges.

## Available Products
{{range .product_catalog}}
- **{{.name}}**: {{.description}} ({{.price}})
{{end}}

## Sales Methodology: {{.sales_methodology | title}}
{{if eq .sales_methodology "consultative"}}
Focus on understanding the customer's needs through thoughtful questions before presenting solutions.
{{else if eq .sales_methodology "solution"}}
Position yourself as a trusted advisor who provides comprehensive solutions to business problems.
{{else if eq .sales_methodology "challenger"}}
Challenge the customer's thinking and provide unique insights about their industry.
{{end}}

## Objection Handling Strategies
{{with .objection_strategies}}
- **Price Objections**: {{.price_objections}}
- **Timing Objections**: {{.timing_objections}}
- **Authority Objections**: {{.authority_objections}}
{{end}}

## Guidelines
1. Always maintain a professional and helpful tone
2. Ask qualifying questions to understand needs
3. Present solutions that directly address stated problems
4. Handle objections with empathy and evidence
5. Guide conversations toward next steps

Remember: Your goal is to help customers make informed decisions that benefit their business.
```

## Agent Components

### Agent Specification
```yaml
apiVersion: "forge.dev/v1"
kind: "Agent"
metadata:
  name: "sales-specialist"
  version: "1.5.0"
  description: "Specialized AI agent for enterprise sales operations"
  author: "agents-team"
  license: "Apache-2.0"
  
  tags: ["sales", "enterprise", "crm", "customer-engagement"]
  categories: ["sales", "customer-service"]
  keywords: ["sales", "agent", "crm", "enterprise", "automation"]
  
  stability: "stable"
  maturity: "mature"
  forge_version: ">=0.2.0"

spec:
  # Agent dependencies
  dependencies:
    tools:
      - repository: "github.com/company/crm-tools"
        component: "salesforce-lookup"
        version: "^2.1.0"
        required: true
        purpose: "Customer data retrieval"
      
      - repository: "github.com/tools/communication"
        component: "email-composer"
        version: "~1.3.0"
        required: true
        purpose: "Email communication"
      
      - repository: "github.com/tools/calendar"
        component: "meeting-scheduler"
        version: ">=2.0.0"
        required: false
        purpose: "Meeting scheduling"
    
    prompts:
      - repository: "github.com/ai-prompts/sales-templates"
        component: "sales-system-v3"
        version: ">=3.0.0"
        role: "system"
        required: true
        purpose: "Primary system prompt"
      
      - repository: "github.com/ai-prompts/sales-templates"
        component: "objection-handling"
        version: "latest"
        role: "task"
        required: true
        purpose: "Objection handling guidance"
      
      - repository: "github.com/ai-prompts/sales-templates"
        component: "closing-techniques"
        version: "^2.0.0"
        role: "task"
        required: false
        purpose: "Sales closing assistance"
    
    agents: []  # Can inherit from other agents
  
  # LLM configuration
  llm:
    primary:
      provider: "openai"
      model: "gpt-4"
      temperature: 0.7
      max_tokens: 2048
      top_p: 0.9
      frequency_penalty: 0.0
      presence_penalty: 0.0
    
    fallback:
      provider: "anthropic"
      model: "claude-3-sonnet"
      temperature: 0.7
      max_tokens: 2048
    
    # Model-specific optimizations
    optimizations:
      openai:
        function_calling: true
        parallel_function_calls: true
        response_format: "auto"
      anthropic:
        prefill_assistant_message: true
        stop_sequences: ["Human:", "Assistant:"]
  
  # Behavioral configuration
  behavior:
    execution_mode: "agent"      # agent, direct, hybrid
    max_iterations: 5
    iteration_timeout: 30        # seconds
    total_timeout: 300          # seconds
    
    # Memory and context
    memory_limit: "100MB"
    context_window: 8000
    context_strategy: "sliding"  # sliding, summarize, truncate
    
    # Tool usage
    tool_timeout: 30
    max_parallel_tools: 3
    tool_retry_count: 2
    
    # Response characteristics
    response_style: "professional"
    verbosity: "balanced"       # concise, balanced, detailed
    creativity: "moderate"      # low, moderate, high
  
  # Capabilities and skills
  capabilities:
    primary:
      - "lead_qualification"
      - "objection_handling"
      - "product_recommendation"
      - "follow_up_scheduling"
      - "crm_data_retrieval"
    
    secondary:
      - "market_research"
      - "competitive_analysis"
      - "proposal_generation"
      - "contract_negotiation"
    
    languages:
      - "en"  # Primary
      - "es"  # Secondary
      - "fr"  # Secondary
  
  # Configuration options
  configuration:
    customizable_fields:
      - field: "llm.primary.temperature"
        min: 0.0
        max: 2.0
        step: 0.1
        description: "Controls response creativity"
      
      - field: "llm.primary.max_tokens"
        min: 512
        max: 4096
        step: 128
        description: "Maximum response length"
      
      - field: "behavior.max_iterations"
        min: 1
        max: 10
        step: 1
        description: "Maximum reasoning iterations"
      
      - field: "behavior.verbosity"
        enum: ["concise", "balanced", "detailed"]
        description: "Response detail level"
    
    presets:
      - name: "conservative"
        description: "Lower temperature, fewer iterations"
        overrides:
          llm.primary.temperature: 0.5
          behavior.max_iterations: 3
      
      - name: "creative"
        description: "Higher temperature, more iterations"
        overrides:
          llm.primary.temperature: 0.9
          behavior.max_iterations: 7
      
      - name: "enterprise"
        description: "Optimized for enterprise customers"
        overrides:
          llm.primary.temperature: 0.6
          behavior.verbosity: "detailed"
  
  # Performance and monitoring
  performance:
    expected_metrics:
      avg_response_time: "2.5s"
      success_rate: 0.95
      customer_satisfaction: 4.6
      conversion_rate: 0.23
    
    sla:
      max_response_time: "10s"
      min_availability: 0.99
      max_error_rate: 0.05
  
  # Integration requirements
  integration:
    required_apis:
      - "salesforce"
      - "email_service"
    
    optional_apis:
      - "calendar_service"
      - "document_generation"
    
    webhooks:
      - event: "lead_qualified"
        endpoint: "/webhooks/lead-qualified"
      - event: "meeting_scheduled"
        endpoint: "/webhooks/meeting-scheduled"
  
  # Security and compliance
  security:
    data_handling:
      - "customer_data"
      - "sales_data"
      - "communication_logs"
    
    compliance:
      - "GDPR"
      - "CCPA"
      - "SOX"
    
    audit_requirements:
      - "conversation_logging"
      - "decision_tracking"
      - "data_access_logging"
```

## Quality Standards

### Component Quality Checklist
```yaml
quality_requirements:
  documentation:
    - complete_readme: true
    - api_documentation: true
    - usage_examples: true
    - changelog: true
    - license_file: true
  
  testing:
    - unit_tests: true
    - integration_tests: true
    - performance_tests: true
    - security_tests: true
    - min_coverage: 80
  
  code_quality:
    - linting_passed: true
    - security_scan_passed: true
    - dependency_audit_passed: true
    - no_hardcoded_secrets: true
  
  metadata:
    - complete_manifest: true
    - valid_semver: true
    - proper_categorization: true
    - accurate_dependencies: true
```

### Validation Rules
```go
type ComponentValidator struct {
    schemaValidator *jsonschema.Validator
    contentValidator *ContentValidator
}

func (cv *ComponentValidator) ValidateComponent(component *Component) []ValidationError {
    var errors []ValidationError
    
    // Schema validation
    if err := cv.schemaValidator.Validate(component.Definition); err != nil {
        errors = append(errors, ValidationError{
            Type: "schema",
            Message: err.Error(),
        })
    }
    
    // Content validation
    if contentErrors := cv.contentValidator.Validate(component); len(contentErrors) > 0 {
        errors = append(errors, contentErrors...)
    }
    
    // Dependency validation
    if depErrors := cv.validateDependencies(component); len(depErrors) > 0 {
        errors = append(errors, depErrors...)
    }
    
    return errors
}
```

## Versioning and Compatibility

### Semantic Versioning Rules
- **Major (X.0.0)**: Breaking changes to component interface
- **Minor (0.X.0)**: New features, backward compatible
- **Patch (0.0.X)**: Bug fixes, no interface changes

### Compatibility Matrix
```yaml
compatibility:
  forge_versions:
    "0.1.x": ["component_v1"]
    "0.2.x": ["component_v1", "component_v2"]
    "0.3.x": ["component_v2", "component_v3"]
  
  breaking_changes:
    "2.0.0":
      - "Changed input schema format"
      - "Removed deprecated fields"
      - "Updated authentication method"
    "3.0.0":
      - "New output format"
      - "Required environment variables changed"
```