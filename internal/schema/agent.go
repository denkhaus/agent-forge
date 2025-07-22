package schema

import "fmt"

// AgentType represents the type of agent implementation.
type AgentType string

const (
	AgentTypeConversational AgentType = "conversational"
	AgentTypeTask          AgentType = "task"
	AgentTypeWorkflow      AgentType = "workflow"
	AgentTypeReactive      AgentType = "reactive"
	AgentTypeAutonomous    AgentType = "autonomous"
)

// AgentCapability represents a capability that the agent possesses.
type AgentCapability struct {
	Name        string            `yaml:"name" json:"name" validate:"required"`
	Description string            `yaml:"description" json:"description" validate:"required"`
	Type        string            `yaml:"type" json:"type" validate:"required,oneof=tool prompt reasoning memory planning"`
	Config      map[string]string `yaml:"config,omitempty" json:"config,omitempty"`
}

// AgentTool represents a tool that the agent can use.
type AgentTool struct {
	Name      string            `yaml:"name" json:"name" validate:"required"`
	Type      string            `yaml:"type" json:"type" validate:"required"`
	Source    string            `yaml:"source" json:"source" validate:"required"`
	Version   string            `yaml:"version,omitempty" json:"version,omitempty"`
	Config    map[string]string `yaml:"config,omitempty" json:"config,omitempty"`
	Required  bool              `yaml:"required" json:"required"`
	Condition string            `yaml:"condition,omitempty" json:"condition,omitempty"`
}

// AgentPrompt represents a prompt that the agent uses.
type AgentPrompt struct {
	Name     string            `yaml:"name" json:"name" validate:"required"`
	Type     string            `yaml:"type" json:"type" validate:"required,oneof=system instruction template"`
	Source   string            `yaml:"source" json:"source" validate:"required"`
	Version  string            `yaml:"version,omitempty" json:"version,omitempty"`
	Priority int               `yaml:"priority,omitempty" json:"priority,omitempty"`
	Config   map[string]string `yaml:"config,omitempty" json:"config,omitempty"`
}

// AgentMemory represents memory configuration for the agent.
type AgentMemory struct {
	Type       string `yaml:"type" json:"type" validate:"required,oneof=short-term long-term episodic semantic working"`
	Provider   string `yaml:"provider,omitempty" json:"provider,omitempty"`
	Capacity   *int   `yaml:"capacity,omitempty" json:"capacity,omitempty"`
	Persistence bool  `yaml:"persistence" json:"persistence"`
	TTL        string `yaml:"ttl,omitempty" json:"ttl,omitempty"`
}

// AgentModel represents model configuration for the agent.
type AgentModel struct {
	Provider    string             `yaml:"provider" json:"provider" validate:"required"`
	Model       string             `yaml:"model" json:"model" validate:"required"`
	Temperature *float64           `yaml:"temperature,omitempty" json:"temperature,omitempty"`
	MaxTokens   *int               `yaml:"maxTokens,omitempty" json:"maxTokens,omitempty"`
	TopP        *float64           `yaml:"topP,omitempty" json:"topP,omitempty"`
	TopK        *int               `yaml:"topK,omitempty" json:"topK,omitempty"`
	Fallback    []AgentModelFallback `yaml:"fallback,omitempty" json:"fallback,omitempty"`
}

// AgentModelFallback represents fallback model configuration.
type AgentModelFallback struct {
	Provider string   `yaml:"provider" json:"provider" validate:"required"`
	Model    string   `yaml:"model" json:"model" validate:"required"`
	Condition string  `yaml:"condition,omitempty" json:"condition,omitempty"`
}

// AgentBehavior represents behavioral configuration for the agent.
type AgentBehavior struct {
	Personality   string            `yaml:"personality,omitempty" json:"personality,omitempty"`
	Goals         []string          `yaml:"goals,omitempty" json:"goals,omitempty"`
	Constraints   []string          `yaml:"constraints,omitempty" json:"constraints,omitempty"`
	Ethics        []string          `yaml:"ethics,omitempty" json:"ethics,omitempty"`
	Communication map[string]string `yaml:"communication,omitempty" json:"communication,omitempty"`
}

// AgentWorkflow represents workflow configuration for task-based agents.
type AgentWorkflow struct {
	Steps       []AgentWorkflowStep `yaml:"steps" json:"steps" validate:"required"`
	ErrorPolicy string              `yaml:"errorPolicy,omitempty" json:"errorPolicy,omitempty"`
	Timeout     string              `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Retries     *int                `yaml:"retries,omitempty" json:"retries,omitempty"`
}

// AgentWorkflowStep represents a single step in an agent workflow.
type AgentWorkflowStep struct {
	Name        string            `yaml:"name" json:"name" validate:"required"`
	Type        string            `yaml:"type" json:"type" validate:"required,oneof=tool prompt decision loop condition"`
	Action      string            `yaml:"action" json:"action" validate:"required"`
	Input       map[string]string `yaml:"input,omitempty" json:"input,omitempty"`
	Output      string            `yaml:"output,omitempty" json:"output,omitempty"`
	Condition   string            `yaml:"condition,omitempty" json:"condition,omitempty"`
	OnSuccess   string            `yaml:"onSuccess,omitempty" json:"onSuccess,omitempty"`
	OnFailure   string            `yaml:"onFailure,omitempty" json:"onFailure,omitempty"`
	Timeout     string            `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Retries     *int              `yaml:"retries,omitempty" json:"retries,omitempty"`
}

// AgentInterface represents interface configuration for the agent.
type AgentInterface struct {
	Type     string            `yaml:"type" json:"type" validate:"required,oneof=cli web api webhook mcp"`
	Protocol string            `yaml:"protocol,omitempty" json:"protocol,omitempty"`
	Port     *int              `yaml:"port,omitempty" json:"port,omitempty"`
	Path     string            `yaml:"path,omitempty" json:"path,omitempty"`
	Config   map[string]string `yaml:"config,omitempty" json:"config,omitempty"`
}

// AgentSecurity represents security configuration for the agent.
type AgentSecurity struct {
	Sandbox         bool     `yaml:"sandbox" json:"sandbox"`
	AllowedDomains  []string `yaml:"allowedDomains,omitempty" json:"allowedDomains,omitempty"`
	BlockedDomains  []string `yaml:"blockedDomains,omitempty" json:"blockedDomains,omitempty"`
	MaxExecutionTime string  `yaml:"maxExecutionTime,omitempty" json:"maxExecutionTime,omitempty"`
	ResourceLimits  *AgentResourceLimits `yaml:"resourceLimits,omitempty" json:"resourceLimits,omitempty"`
}

// AgentResourceLimits represents resource limits for the agent.
type AgentResourceLimits struct {
	CPU       string `yaml:"cpu,omitempty" json:"cpu,omitempty"`
	Memory    string `yaml:"memory,omitempty" json:"memory,omitempty"`
	Disk      string `yaml:"disk,omitempty" json:"disk,omitempty"`
	Network   string `yaml:"network,omitempty" json:"network,omitempty"`
	Processes *int   `yaml:"processes,omitempty" json:"processes,omitempty"`
}

// AgentSpec defines the specification for Agent components.
type AgentSpec struct {
	// Type of agent
	Type AgentType `yaml:"type" json:"type" validate:"required"`
	
	// Capabilities that the agent possesses
	Capabilities []AgentCapability `yaml:"capabilities,omitempty" json:"capabilities,omitempty"`
	
	// Tools that the agent can use
	Tools []AgentTool `yaml:"tools,omitempty" json:"tools,omitempty"`
	
	// Prompts that the agent uses
	Prompts []AgentPrompt `yaml:"prompts,omitempty" json:"prompts,omitempty"`
	
	// Memory configuration
	Memory []AgentMemory `yaml:"memory,omitempty" json:"memory,omitempty"`
	
	// Model configuration
	Model AgentModel `yaml:"model" json:"model" validate:"required"`
	
	// Behavioral configuration
	Behavior *AgentBehavior `yaml:"behavior,omitempty" json:"behavior,omitempty"`
	
	// Workflow configuration (for task/workflow agents)
	Workflow *AgentWorkflow `yaml:"workflow,omitempty" json:"workflow,omitempty"`
	
	// Interface configuration
	Interface AgentInterface `yaml:"interface" json:"interface" validate:"required"`
	
	// Security configuration
	Security *AgentSecurity `yaml:"security,omitempty" json:"security,omitempty"`
	
	// Environment variables
	Environment map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
	
	// Configuration files
	ConfigFiles []string `yaml:"configFiles,omitempty" json:"configFiles,omitempty"`
}

// Agent represents a complete Agent component manifest.
type Agent struct {
	BaseComponent `yaml:",inline"`
	Spec          AgentSpec `yaml:"spec" json:"spec" validate:"required"`
}

// NewAgent creates a new Agent with default values.
func NewAgent(name, version string) *Agent {
	return &Agent{
		BaseComponent: BaseComponent{
			APIVersion: APIVersion,
			Kind:       KindAgent,
			Metadata: BaseMetadata{
				Name:      name,
				Version:   version,
				Stability: StabilityExperimental,
				Maturity:  MaturityAlpha,
			},
		},
	}
}

// Validate performs Agent-specific validation.
func (a *Agent) Validate() error {
	if err := a.BaseComponent.Validate(); err != nil {
		return err
	}
	
	if a.Kind != KindAgent {
		return fmt.Errorf("invalid kind for Agent: %s", a.Kind)
	}
	
	// Validate that model is configured
	if a.Spec.Model.Provider == "" || a.Spec.Model.Model == "" {
		return fmt.Errorf("agent must have model provider and model configured")
	}
	
	// Validate tool names are unique
	toolNames := make(map[string]bool)
	for _, tool := range a.Spec.Tools {
		if toolNames[tool.Name] {
			return fmt.Errorf("duplicate tool name: %s", tool.Name)
		}
		toolNames[tool.Name] = true
	}
	
	// Validate prompt names are unique
	promptNames := make(map[string]bool)
	for _, prompt := range a.Spec.Prompts {
		if promptNames[prompt.Name] {
			return fmt.Errorf("duplicate prompt name: %s", prompt.Name)
		}
		promptNames[prompt.Name] = true
	}
	
	// Validate workflow steps if workflow is defined
	if a.Spec.Workflow != nil {
		stepNames := make(map[string]bool)
		for _, step := range a.Spec.Workflow.Steps {
			if stepNames[step.Name] {
				return fmt.Errorf("duplicate workflow step name: %s", step.Name)
			}
			stepNames[step.Name] = true
		}
	}
	
	return nil
}

// GetToolByName returns a tool by name.
func (a *Agent) GetToolByName(name string) *AgentTool {
	for _, tool := range a.Spec.Tools {
		if tool.Name == name {
			return &tool
		}
	}
	return nil
}

// GetPromptByName returns a prompt by name.
func (a *Agent) GetPromptByName(name string) *AgentPrompt {
	for _, prompt := range a.Spec.Prompts {
		if prompt.Name == name {
			return &prompt
		}
	}
	return nil
}

// HasTool checks if a tool exists.
func (a *Agent) HasTool(name string) bool {
	return a.GetToolByName(name) != nil
}

// HasPrompt checks if a prompt exists.
func (a *Agent) HasPrompt(name string) bool {
	return a.GetPromptByName(name) != nil
}