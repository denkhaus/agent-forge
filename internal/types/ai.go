package types

import (
	"context"
	"io"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
)

// Prompt represents a template prompt that can be used by the AI agent.
type Prompt interface {
	// GetName returns the prompt's name
	GetName() string

	// GetDescription returns the prompt's description
	GetDescription() string

	// GetTemplate returns the prompt's template
	GetTemplate() string

	// Render renders the prompt template with the given variables
	Render(variables map[string]any) (string, error)
}

// Agent represents an AI agent with specific tool requirements and behavior.
type Agent interface {
	// GetName returns the agent's name
	GetName() string

	// GetDescription returns the agent's description
	GetDescription() string

	// GetSystemPrompt returns the agent's system prompt
	GetSystemPrompt() string

	// GetRequiredTools returns the names of tools this agent requires
	GetRequiredTools() []string

	// GetLLMConfig returns the agent's LLM configuration
	GetLLMConfig() LLMConfig

	// HasRequiredTool checks if the agent requires a specific tool
	HasRequiredTool(toolName string) bool

	// Clone creates a copy of the agent with optional overrides
	Clone(overrides map[string]any) Agent

	// Note: Session creation is handled by the session package factory
	// This keeps agents focused on their core responsibility: representing agent data
}

// LLMConfig represents the configuration for the underlying LLM.
type LLMConfig interface {
	// GetProvider returns the LLM provider name (e.g., "googleai", "openai", "anthropic")
	GetProvider() string

	// GetModel returns the model name (e.g., "gemini-pro", "gpt-4", "claude-3")
	GetModel() string

	// GetTemperature returns the creativity/randomness setting (0.0 - 1.0)
	GetTemperature() float64

	// GetMaxTokens returns the maximum response length
	GetMaxTokens() int

	// GetParameters returns provider-specific parameters
	GetParameters() map[string]any

	// GetParameter returns a specific parameter value
	GetParameter(key string) (any, bool)

	// Validate checks if the configuration is valid
	Validate() error
}

// ToolProvider defines the interface for managing AI tools.
type ToolProvider interface {
	// GetTools returns all available tools as langchain-go tools
	GetTools() []tools.Tool

	// GetToolsForAgent returns tools required by the agent, returns error if any tools are missing
	GetToolsForAgent(agent Agent) ([]tools.Tool, error)

	// ExecuteTool executes a specific tool with string input (JSON format)
	ExecuteTool(ctx context.Context, name string, input string) (string, error)

	// RegisterTool registers a standard langchain-go tool
	RegisterTool(tool tools.Tool) error

	// HasTool checks if a tool with the given name is available
	HasTool(name string) bool

	// ValidateAgentRequirements checks if all required tools for an agent are available
	ValidateAgentRequirements(agent Agent) error

	// GetToolNames returns the names of all available tools
	GetToolNames() []string
}

// PromptProvider defines the interface for managing AI prompts and templates.
type PromptProvider interface {
	// GetSystemPrompt returns the system prompt for the AI agent
	GetSystemPrompt() string

	// FormatPrompt formats a prompt template with the given variables
	FormatPrompt(template string, vars map[string]any) (string, error)

	// GetPrompts returns all available prompts
	GetPrompts() map[string]Prompt

	// GetPrompt retrieves a prompt by name, returns nil if not found
	GetPrompt(name string) Prompt

	// RegisterPrompt registers a new prompt with the provider
	RegisterPrompt(prompt Prompt) error
}

// AgentProvider defines the interface for providing and managing agents.
type AgentProvider interface {
	// GetAgents returns all available agents
	GetAgents() map[string]Agent

	// GetAgent returns a specific agent by name
	GetAgent(name string) (Agent, error)

	// GetDefaultAgent returns the default agent (typically "planner")
	GetDefaultAgent() (Agent, error)
}

// LLMService defines the interface for LLM initialization and management.
type LLMService interface {
	// InitializeLLM creates an LLM instance based on the configuration
	InitializeLLM(ctx context.Context, config Config, llmConfig LLMConfig) (llms.Model, error)

	// ValidateConfig validates the LLM configuration
	ValidateConfig(config Config, llmConfig LLMConfig) error
}

// ComponentParser defines the interface for parsing AgentForge component manifests.
type ComponentParser interface {
	// ParseComponent parses a component manifest from YAML content
	ParseComponent(content []byte) (interface{}, error)

	// ParseComponentFromReader parses a component manifest from an io.Reader
	ParseComponentFromReader(reader io.Reader) (interface{}, error)

	// ParseComponentFromString parses a component manifest from a string
	ParseComponentFromString(content string) (interface{}, error)

	// SerializeComponent serializes a component to YAML
	SerializeComponent(component interface{}) ([]byte, error)

	// SerializeComponentToString serializes a component to a YAML string
	SerializeComponentToString(component interface{}) (string, error)

	// ValidateComponentKind validates that a component kind is supported
	ValidateComponentKind(kind string) error

	// GetComponentInfo extracts basic information from a component
	GetComponentInfo(component interface{}) (*ComponentInfo, error)
}

// ComponentInfo represents basic information about a component.
type ComponentInfo struct {
	APIVersion  string   `json:"apiVersion"`
	Kind        string   `json:"kind"`
	Name        string   `json:"name"`
	Namespace   string   `json:"namespace,omitempty"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	License     string   `json:"license"`
	Tags        []string `json:"tags,omitempty"`
	Categories  []string `json:"categories,omitempty"`
	Stability   string   `json:"stability"`
	Maturity    string   `json:"maturity"`
}

// GitHubRepository represents a GitHub repository with AgentForge component metadata.
type GitHubRepository struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	HTMLURL     string    `json:"html_url"`
	CloneURL    string    `json:"clone_url"`
	StarCount   int       `json:"stargazers_count"`
	Language    string    `json:"language"`
	UpdatedAt   time.Time `json:"updated_at"`
	Topics      []string  `json:"topics"`
}

// GitHubComponent represents an AgentForge component found in a repository.
type GitHubComponent struct {
	Type        string    `json:"type"` // tool, prompt, agent
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Repository  string    `json:"repository"`
	Path        string    `json:"path"`
	Version     string    `json:"version"`
	Stars       int       `json:"stars"`
	UpdatedAt   time.Time `json:"updated_at"`
	Topics      []string  `json:"topics"`
}

// GitHubSearchOptions represents options for searching GitHub repositories.
type GitHubSearchOptions struct {
	Sort      string // stars, updated, created
	Order     string // asc, desc
	PerPage   int    // results per page (max 100)
	Page      int    // page number
	Language  string // filter by language
	Topic     string // filter by topic
	MinStars  int    // minimum star count
	MaxStars  int    // maximum star count
	CreatedAt string // created date filter (YYYY-MM-DD)
	UpdatedAt string // updated date filter (YYYY-MM-DD)
}

// GitHubClient defines the interface for GitHub API integration.
type GitHubClient interface {
	// SearchRepositories searches for repositories containing AgentForge components
	SearchRepositories(ctx context.Context, query string, opts *GitHubSearchOptions) ([]*GitHubRepository, error)

	// GetRepository gets detailed information about a specific repository
	GetRepository(ctx context.Context, owner, repo string) (*GitHubRepository, error)

	// ListComponents lists all AgentForge components in a repository
	ListComponents(ctx context.Context, owner, repo string) ([]*GitHubComponent, error)

	// GetComponentContent gets the raw content of a component file
	GetComponentContent(ctx context.Context, owner, repo, path string) ([]byte, error)
}

// AgentSession defines the interface for interacting with a specific agent.
type AgentSession interface {
	// GetAgent returns the agent associated with this session
	GetAgent() Agent

	// GenerateResponse generates a response using the AI agent with tools.Tool
	GenerateResponse(ctx context.Context, messages []llms.MessageContent, agentTools []tools.Tool) (*llms.ContentResponse, error)

	// Chat starts an interactive chat session
	Chat(ctx context.Context, initialMessage string) error

	// GetMessageHistory returns the current message history
	GetMessageHistory() []llms.MessageContent

	// ClearMessageHistory clears the current message history
	ClearMessageHistory()

	// AddMessage adds a message to the history
	AddMessage(message llms.MessageContent)

	// SwitchAgent switches to a different agent
	SwitchAgent(ctx context.Context, agentName string) error

	// GetExecutionMode returns the current execution mode
	GetExecutionMode() ExecutionMode

	// GetSessionConfig returns the session configuration
	GetSessionConfig() AgentSessionConfig
}
