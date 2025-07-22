// Package types contains all type definitions and interfaces used throughout the application.
package types

// Logger interface removed - use *zap.Logger directly throughout the system

// Removed Provider interface - over-engineering, providers implement Startup() directly if needed

// Config represents the application configuration interface.
type Config interface {
	GetLogLevel() string
	GetPort() int
	GetEnvironment() string
	GetAPIKey(provider string) string
	GetAPIKeys() map[string]string
}

// Configuration types for factories

// ExecutionMode defines how the agent session executes requests.
type ExecutionMode int

const (
	// ExecutionModeDirect uses direct LLM calls with fine-grained control
	ExecutionModeDirect ExecutionMode = iota
	// ExecutionModeAgent uses OneShotAgent with automatic reasoning loops
	ExecutionModeAgent
)

// String returns the string representation of the execution mode.
func (e ExecutionMode) String() string {
	switch e {
	case ExecutionModeDirect:
		return "direct"
	case ExecutionModeAgent:
		return "agent"
	default:
		return "unknown"
	}
}

// AgentSessionConfig represents configuration for creating agent sessions.
type AgentSessionConfig struct {
	ExecutionMode ExecutionMode
	MaxIterations int // For agents.WithMaxIterations() in agent mode
}

// NewAgentSessionConfig creates a new agent session configuration with defaults.
func NewAgentSessionConfig() AgentSessionConfig {
	return AgentSessionConfig{
		ExecutionMode: ExecutionModeDirect,
		MaxIterations: 3,
	}
}

// AgentConfig represents configuration for creating agents.
type AgentConfig struct {
	Name         string
	Description  string
	Type         string
	SystemPrompt string
	Capabilities []string
	LLMConfig    LLMConfig
}

// ToolConfig represents configuration for creating tools.
type ToolConfig struct {
	Name   string
	Type   string
	Config interface{} // Specific config based on tool type
}

// WeatherToolConfig represents configuration for weather tools.
type WeatherToolConfig struct {
	Name        string
	Description string
	Schema      map[string]any
	APIKey      string
	BaseURL     string
}

// TaskToolConfig represents configuration for task tools.
type TaskToolConfig struct {
	Name        string
	Description string
	Schema      map[string]any
}

// PromptConfig represents configuration for creating prompts.
type PromptConfig struct {
	Name      string
	Type      string
	Template  string
	Variables map[string]string
	Metadata  map[string]interface{}
}
