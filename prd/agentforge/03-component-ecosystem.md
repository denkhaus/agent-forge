# Component Ecosystem

## Git-Native Component System

AgentForge's component ecosystem is built on Git repositories, where each repository can contain multiple related components with strong versioning and dependency management.

## Repository Structure

### Standard Repository Layout
```
github.com/org/component-repo/
├── forge-manifest.yaml        # Repository metadata and component catalog
├── components/
│   ├── tools/
│   │   ├── tool-name/
│   │   │   ├── component.yaml # Component definition
│   │   │   ├── server.go      # Implementation (for tools)
│   │   │   ├── schema.json    # Input/output schema
│   │   │   └── README.md      # Documentation
│   │   └── another-tool/
│   ├── prompts/
│   │   ├── prompt-name/
│   │   │   ├── component.yaml
│   │   │   ├── template.txt   # Prompt template
│   │   │   └── variables.json # Template variables
│   │   └── another-prompt/
│   └── agents/
│       ├── agent-name/
│       │   ├── component.yaml
│       │   ├── config.yaml    # Default configuration
│       │   └── dependencies.yaml
│       └── another-agent/
├── examples/                  # Usage examples
├── tests/                     # Component tests
├── docs/                      # Documentation
└── CHANGELOG.md              # Version history
```

## Component Types

### 1. Tools
**Purpose**: Executable functions that agents can call
**Examples**: API integrations, data processors, calculators

```yaml
# components/tools/salesforce-lookup/component.yaml
apiVersion: "forge.dev/v1"
kind: "Tool"
metadata:
  name: "salesforce-lookup"
  version: "2.1.0"
  description: "Query Salesforce CRM for customer data"
  author: "company-team"
  license: "MIT"

spec:
  type: "mcp-server"
  runtime: "go"
  binary: "./bin/salesforce-server"
  
  schema:
    input:
      type: "object"
      properties:
        query_type:
          type: "string"
          enum: ["account", "contact", "opportunity"]
        search_term:
          type: "string"
          minLength: 1
        fields:
          type: "array"
          items:
            type: "string"
      required: ["query_type", "search_term"]
    
    output:
      type: "object"
      properties:
        results:
          type: "array"
        total_count:
          type: "integer"
        has_more:
          type: "boolean"

  config:
    environment:
      - name: "SALESFORCE_URL"
        required: true
        description: "Salesforce instance URL"
      - name: "SALESFORCE_TOKEN"
        required: true
        secret: true
        description: "Salesforce API token"
    
    parameters:
      timeout:
        type: "integer"
        default: 30
        description: "Request timeout in seconds"
      retry_count:
        type: "integer"
        default: 3
        description: "Number of retry attempts"

  capabilities:
    - "customer_lookup"
    - "crm_integration"
    - "data_retrieval"

  dependencies:
    external:
      - "salesforce-api"
    forge_components: []

  security:
    permissions:
      - "network.outbound"
      - "env.read"
    sandboxed: true
```

### 2. Prompts
**Purpose**: Template-based text generation for AI interactions
**Examples**: System prompts, task templates, conversation starters

```yaml
# components/prompts/sales-system/component.yaml
apiVersion: "forge.dev/v1"
kind: "Prompt"
metadata:
  name: "sales-system-v3"
  version: "3.2.0"
  description: "Advanced sales assistant system prompt"
  author: "ai-prompts-team"
  license: "CC-BY-4.0"

spec:
  type: "system"
  category: "sales"
  
  template_engine: "go-template"
  template_file: "./template.txt"
  
  variables:
    - name: "company_name"
      type: "string"
      required: true
      description: "Name of the company"
    - name: "product_catalog"
      type: "array"
      required: false
      description: "Available products"
    - name: "sales_methodology"
      type: "string"
      default: "consultative"
      enum: ["consultative", "solution", "challenger"]

  validation:
    max_length: 4000
    required_sections:
      - "role_definition"
      - "capabilities"
      - "guidelines"

  localization:
    supported_languages: ["en", "es", "fr", "de"]
    default_language: "en"

  metadata:
    tags: ["sales", "crm", "customer-service"]
    use_cases: ["lead_qualification", "objection_handling", "closing"]
```

### 3. Agents
**Purpose**: Complete AI agent configurations with dependencies
**Examples**: Specialized assistants, domain experts, workflow agents

```yaml
# components/agents/sales-specialist/component.yaml
apiVersion: "forge.dev/v1"
kind: "Agent"
metadata:
  name: "sales-specialist"
  version: "1.5.0"
  description: "Specialized AI agent for sales operations"
  author: "agents-team"
  license: "Apache-2.0"

spec:
  dependencies:
    tools:
      - repository: "github.com/company/crm-tools"
        component: "salesforce-lookup"
        version: "^2.1.0"
        commit: "a1b2c3d4e5f6789012345678901234567890abcd"
      - repository: "github.com/tools/communication"
        component: "email-composer"
        version: "~1.3.0"
        commit: "b2c3d4e5f6789012345678901234567890abcdef"
    
    prompts:
      - repository: "github.com/ai-prompts/sales-templates"
        component: "sales-system-v3"
        version: ">=3.0.0"
        commit: "c3d4e5f6789012345678901234567890abcdef12"
        role: "system"
      - repository: "github.com/ai-prompts/sales-templates"
        component: "objection-handling"
        version: "latest"
        commit: "d4e5f6789012345678901234567890abcdef1234"
        role: "task"

  llm:
    provider: "openai"
    model: "gpt-4"
    temperature: 0.7
    max_tokens: 2048
    
    fallback:
      provider: "anthropic"
      model: "claude-3-sonnet"
      temperature: 0.7

  behavior:
    execution_mode: "agent"  # or "direct"
    max_iterations: 5
    tool_timeout: 30
    memory_limit: "100MB"
    
    capabilities:
      - "lead_qualification"
      - "objection_handling"
      - "product_recommendation"
      - "follow_up_scheduling"

  configuration:
    customizable_fields:
      - "llm.temperature"
      - "llm.max_tokens"
      - "behavior.max_iterations"
    
    validation_rules:
      - field: "llm.temperature"
        min: 0.0
        max: 2.0
      - field: "behavior.max_iterations"
        min: 1
        max: 10

  metadata:
    tags: ["sales", "crm", "customer-engagement"]
    use_cases: ["lead_qualification", "sales_calls", "follow_up"]
    performance_metrics:
      - "conversion_rate"
      - "response_time"
      - "customer_satisfaction"
```

## Repository Manifest

### Forge Manifest Format
```yaml
# forge-manifest.yaml
apiVersion: "forge.dev/v1"
kind: "ComponentRepository"
metadata:
  name: "crm-tools"
  description: "Enterprise CRM integration tools and utilities"
  version: "2.1.0"
  author: "company-team"
  license: "MIT"
  homepage: "https://github.com/company/crm-tools"
  documentation: "https://docs.company.com/crm-tools"

repository:
  type: "mixed"  # tools, prompts, agents, or mixed
  categories: ["crm", "sales", "customer-service"]
  
components:
  tools:
    - name: "salesforce-lookup"
      version: "2.1.0"
      path: "components/tools/salesforce-lookup"
      description: "Query Salesforce CRM data"
      stability: "stable"
      
    - name: "hubspot-sync"
      version: "2.0.5"
      path: "components/tools/hubspot-sync"
      description: "Sync data with HubSpot"
      stability: "beta"

  prompts:
    - name: "crm-system-prompt"
      version: "1.2.0"
      path: "components/prompts/crm-system"
      description: "System prompt for CRM operations"
      stability: "stable"

compatibility:
  forge_version: ">=0.2.0"
  mcp_protocol: ">=1.0.0"
  
  platforms:
    - "linux/amd64"
    - "linux/arm64"
    - "darwin/amd64"
    - "darwin/arm64"
    - "windows/amd64"

security:
  permissions:
    - "network.outbound"
    - "env.read"
  
  secrets:
    - name: "SALESFORCE_TOKEN"
      description: "Salesforce API authentication token"
    - name: "HUBSPOT_API_KEY"
      description: "HubSpot API key"

quality:
  test_coverage: 85
  documentation_score: 90
  community_rating: 4.7
  download_count: 15420
  
  badges:
    - "verified"
    - "enterprise-ready"
    - "well-documented"

maintenance:
  status: "active"
  last_updated: "2024-01-15T10:30:00Z"
  maintainers:
    - "john.doe@company.com"
    - "jane.smith@company.com"
  
  support:
    issues: "https://github.com/company/crm-tools/issues"
    discussions: "https://github.com/company/crm-tools/discussions"
    documentation: "https://docs.company.com/crm-tools"
```

## Component Discovery

### Search and Discovery
```bash
# Search across all repositories
forge search tools --keyword "crm" --category "sales"
forge search prompts --tag "customer-service" --language "en"
forge search agents --capability "lead_qualification"

# Repository-specific search
forge search --repo github.com/company/crm-tools --type tools
forge search --repo github.com/ai-prompts/sales-templates --category "objection"

# Advanced filtering
forge search tools \
  --stability stable \
  --min-rating 4.0 \
  --compatible-with forge@0.2.0 \
  --has-tests
```

### Component Metadata
```bash
# Get component information
forge info github.com/company/crm-tools/salesforce-lookup
forge info github.com/ai-prompts/sales-templates/sales-system-v3@v3.2.0

# Show dependencies
forge deps github.com/agents/specialized/sales-specialist
forge deps --tree --depth 3 sales-specialist

# Compatibility check
forge compat check sales-specialist --with-forge 0.2.0
forge compat matrix --components salesforce-lookup,email-composer
```

## Version Management

### Semantic Versioning
- **Major**: Breaking changes to component interface
- **Minor**: New features, backward compatible
- **Patch**: Bug fixes, no interface changes

### Commit-Based References
```yaml
dependencies:
  tools:
    - repository: "github.com/company/crm-tools"
      component: "salesforce-lookup"
      version: "^2.1.0"                    # Semantic version constraint
      commit: "a1b2c3d4e5f6789..."         # Exact commit for reproducibility
      resolved_version: "2.1.3"            # Actually resolved version
```

### Version Constraints
- `^2.1.0`: Compatible with 2.x.x (>=2.1.0, <3.0.0)
- `~2.1.0`: Compatible with 2.1.x (>=2.1.0, <2.2.0)
- `>=2.1.0`: Any version >= 2.1.0
- `latest`: Latest stable version
- `main`: Latest commit on main branch
- `@commit`: Specific commit hash