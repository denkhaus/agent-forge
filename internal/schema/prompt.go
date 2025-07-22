package schema

import "fmt"

// PromptType represents the type of prompt.
type PromptType string

const (
	PromptTypeSystem      PromptType = "system"
	PromptTypeUser        PromptType = "user"
	PromptTypeAssistant   PromptType = "assistant"
	PromptTypeTemplate    PromptType = "template"
	PromptTypeConversation PromptType = "conversation"
)

// PromptFormat represents the format of the prompt content.
type PromptFormat string

const (
	FormatText     PromptFormat = "text"
	FormatMarkdown PromptFormat = "markdown"
	FormatJSON     PromptFormat = "json"
	FormatYAML     PromptFormat = "yaml"
	FormatXML      PromptFormat = "xml"
)

// PromptVariable represents a variable that can be used in the prompt template.
type PromptVariable struct {
	Name        string      `yaml:"name" json:"name" validate:"required"`
	Type        string      `yaml:"type" json:"type" validate:"required,oneof=string number boolean object array"`
	Description string      `yaml:"description" json:"description" validate:"required"`
	Required    bool        `yaml:"required" json:"required"`
	Default     interface{} `yaml:"default,omitempty" json:"default,omitempty"`
	Enum        []string    `yaml:"enum,omitempty" json:"enum,omitempty"`
	Pattern     string      `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	MinLength   *int        `yaml:"minLength,omitempty" json:"minLength,omitempty"`
	MaxLength   *int        `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
	Minimum     *float64    `yaml:"minimum,omitempty" json:"minimum,omitempty"`
	Maximum     *float64    `yaml:"maximum,omitempty" json:"maximum,omitempty"`
}

// PromptMessage represents a single message in a conversation prompt.
type PromptMessage struct {
	Role    string `yaml:"role" json:"role" validate:"required,oneof=system user assistant"`
	Content string `yaml:"content" json:"content" validate:"required"`
}

// PromptExample represents an example usage of the prompt.
type PromptExample struct {
	Name        string                 `yaml:"name" json:"name" validate:"required"`
	Description string                 `yaml:"description" json:"description"`
	Variables   map[string]interface{} `yaml:"variables,omitempty" json:"variables,omitempty"`
	Expected    string                 `yaml:"expected,omitempty" json:"expected,omitempty"`
	Context     string                 `yaml:"context,omitempty" json:"context,omitempty"`
}

// PromptValidation represents validation rules for the prompt.
type PromptValidation struct {
	MaxTokens     *int     `yaml:"maxTokens,omitempty" json:"maxTokens,omitempty"`
	MinTokens     *int     `yaml:"minTokens,omitempty" json:"minTokens,omitempty"`
	RequiredVars  []string `yaml:"requiredVars,omitempty" json:"requiredVars,omitempty"`
	ForbiddenVars []string `yaml:"forbiddenVars,omitempty" json:"forbiddenVars,omitempty"`
}

// PromptMetrics represents metrics and performance data for the prompt.
type PromptMetrics struct {
	AverageTokens    *int     `yaml:"averageTokens,omitempty" json:"averageTokens,omitempty"`
	SuccessRate      *float64 `yaml:"successRate,omitempty" json:"successRate,omitempty"`
	AverageLatency   *string  `yaml:"averageLatency,omitempty" json:"averageLatency,omitempty"`
	RecommendedModel []string `yaml:"recommendedModel,omitempty" json:"recommendedModel,omitempty"`
}

// PromptSpec defines the specification for Prompt components.
type PromptSpec struct {
	// Type of prompt
	Type PromptType `yaml:"type" json:"type" validate:"required"`
	
	// Format of the prompt content
	Format PromptFormat `yaml:"format" json:"format" validate:"required"`
	
	// Template content (for template type prompts)
	Template string `yaml:"template,omitempty" json:"template,omitempty"`
	
	// Content for simple prompts
	Content string `yaml:"content,omitempty" json:"content,omitempty"`
	
	// Messages for conversation type prompts
	Messages []PromptMessage `yaml:"messages,omitempty" json:"messages,omitempty"`
	
	// Variables that can be used in the prompt
	Variables []PromptVariable `yaml:"variables,omitempty" json:"variables,omitempty"`
	
	// Examples of prompt usage
	Examples []PromptExample `yaml:"examples,omitempty" json:"examples,omitempty"`
	
	// Validation rules
	Validation *PromptValidation `yaml:"validation,omitempty" json:"validation,omitempty"`
	
	// Performance metrics
	Metrics *PromptMetrics `yaml:"metrics,omitempty" json:"metrics,omitempty"`
	
	// Model compatibility
	Models []string `yaml:"models,omitempty" json:"models,omitempty"`
	
	// Temperature setting recommendation
	Temperature *float64 `yaml:"temperature,omitempty" json:"temperature,omitempty"`
	
	// Max tokens recommendation
	MaxTokens *int `yaml:"maxTokens,omitempty" json:"maxTokens,omitempty"`
	
	// Stop sequences
	StopSequences []string `yaml:"stopSequences,omitempty" json:"stopSequences,omitempty"`
}

// Prompt represents a complete Prompt component manifest.
type Prompt struct {
	BaseComponent `yaml:",inline"`
	Spec          PromptSpec `yaml:"spec" json:"spec" validate:"required"`
}

// NewPrompt creates a new Prompt with default values.
func NewPrompt(name, version string) *Prompt {
	return &Prompt{
		BaseComponent: BaseComponent{
			APIVersion: APIVersion,
			Kind:       KindPrompt,
			Metadata: BaseMetadata{
				Name:      name,
				Version:   version,
				Stability: StabilityExperimental,
				Maturity:  MaturityAlpha,
			},
		},
	}
}

// Validate performs Prompt-specific validation.
func (p *Prompt) Validate() error {
	if err := p.BaseComponent.Validate(); err != nil {
		return err
	}
	
	if p.Kind != KindPrompt {
		return fmt.Errorf("invalid kind for Prompt: %s", p.Kind)
	}
	
	// Validate that at least one content field is provided
	hasContent := p.Spec.Content != "" || p.Spec.Template != "" || len(p.Spec.Messages) > 0
	if !hasContent {
		return fmt.Errorf("prompt must have content, template, or messages")
	}
	
	// Validate variable names are unique
	variableNames := make(map[string]bool)
	for _, variable := range p.Spec.Variables {
		if variableNames[variable.Name] {
			return fmt.Errorf("duplicate variable name: %s", variable.Name)
		}
		variableNames[variable.Name] = true
	}
	
	return nil
}

// GetVariableByName returns a variable by name.
func (p *Prompt) GetVariableByName(name string) *PromptVariable {
	for _, variable := range p.Spec.Variables {
		if variable.Name == name {
			return &variable
		}
	}
	return nil
}

// HasVariable checks if a variable exists.
func (p *Prompt) HasVariable(name string) bool {
	return p.GetVariableByName(name) != nil
}