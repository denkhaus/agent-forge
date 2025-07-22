// Package database provides database connectivity and operations for AgentForge.
package database

import (
	"context"

	"github.com/denkhaus/agentforge/internal/database/ent"
)

// DatabaseManager defines the interface for database lifecycle management.
type DatabaseManagerNew interface {
	Initialize(ctx context.Context) error
	Shutdown() error
	GetRepositoryService() RepositoryService
	GetToolService() ToolService
	GetPromptService() PromptService
	GetAgentService() AgentService
	GetConfigService() ConfigService
	GetDatabasePath() string
}

// ToolService defines the interface for tool management operations.
type ToolService interface {
	CreateTool(ctx context.Context, req CreateToolRequest) (*ent.Tool, error)
	GetTool(ctx context.Context, id string) (*ent.Tool, error)
	GetToolByName(ctx context.Context, name, version string, repositoryID string) (*ent.Tool, error)
	ListTools(ctx context.Context, opts ListToolsOptions) ([]*ent.Tool, error)
	UpdateTool(ctx context.Context, id string, updates map[string]interface{}) (*ent.Tool, error)
	DeleteTool(ctx context.Context, id string) error
	InstallTool(ctx context.Context, id string, installPath string) (*ent.Tool, error)
	UninstallTool(ctx context.Context, id string) (*ent.Tool, error)
	SearchTools(ctx context.Context, query string, opts SearchToolsOptions) ([]*ent.Tool, error)
	
	// CLI-specific methods
	ListToolsForCLI(ctx context.Context, opts interface{}) ([]interface{}, error)
	CreateToolFiles(ctx context.Context, name string) error
	PullTool(ctx context.Context, repo string) error
	PushTool(ctx context.Context, name string) error
}

// PromptService defines the interface for prompt management operations.
type PromptService interface {
	CreatePrompt(ctx context.Context, req CreatePromptRequest) (*ent.Prompt, error)
	GetPrompt(ctx context.Context, id string) (*ent.Prompt, error)
	GetPromptByName(ctx context.Context, name, version string, repositoryID string) (*ent.Prompt, error)
	ListPrompts(ctx context.Context, opts ListPromptsOptions) ([]*ent.Prompt, error)
	UpdatePrompt(ctx context.Context, id string, updates map[string]interface{}) (*ent.Prompt, error)
	DeletePrompt(ctx context.Context, id string) error
	InstallPrompt(ctx context.Context, id string, installPath string) (*ent.Prompt, error)
	UninstallPrompt(ctx context.Context, id string) (*ent.Prompt, error)
	SearchPrompts(ctx context.Context, query string, opts SearchPromptsOptions) ([]*ent.Prompt, error)
	
	// CLI-specific methods
	ListPromptsForCLI(ctx context.Context, opts interface{}) ([]interface{}, error)
	CreatePromptFiles(ctx context.Context, name string) error
	PullPrompt(ctx context.Context, repo string) error
	PushPrompt(ctx context.Context, name string) error
}

// AgentService defines the interface for agent management operations.
type AgentService interface {
	CreateAgent(ctx context.Context, req CreateAgentRequest) (*ent.Agent, error)
	GetAgent(ctx context.Context, id string) (*ent.Agent, error)
	GetAgentByName(ctx context.Context, name, version string, repositoryID string) (*ent.Agent, error)
	ListAgents(ctx context.Context, opts ListAgentsOptions) ([]*ent.Agent, error)
	UpdateAgent(ctx context.Context, id string, updates map[string]interface{}) (*ent.Agent, error)
	DeleteAgent(ctx context.Context, id string) error
	InstallAgent(ctx context.Context, id string, installPath string) (*ent.Agent, error)
	UninstallAgent(ctx context.Context, id string) (*ent.Agent, error)
	SearchAgents(ctx context.Context, query string, opts SearchAgentsOptions) ([]*ent.Agent, error)
	
	// CLI-specific methods
	ListAgentsForCLI(ctx context.Context, opts interface{}) ([]interface{}, error)
	CreateAgentFiles(ctx context.Context, name string) error
	PullAgent(ctx context.Context, repo string) error
	PushAgent(ctx context.Context, name string) error
}

// Request types for tool operations
type CreateToolRequest struct {
	Name              string
	Namespace         string
	Version           string
	Description       string
	Author            string
	License           string
	Homepage          *string
	Documentation     *string
	Tags              []string
	Categories        []string
	Keywords          []string
	Stability         string
	Maturity          string
	ForgeVersion      string
	Platforms         []string
	Spec              string
	SpecHash          string
	RepositoryID      string
	CommitHash        string
	Branch            string
	ExecutionType     string
	SchemaPath        *string
	ServerConfig      map[string]interface{}
	Capabilities      []string
	EntryPoint        *string
	EnvironmentVars   map[string]string
	RequiredPerms     []string
	TimeoutSeconds    int
	SupportsStreaming bool
	InputSchema       map[string]interface{}
	OutputSchema      map[string]interface{}
}

type ListToolsOptions struct {
	ExecutionType   *string
	Stability       *string
	IsInstalled     *bool
	RepositoryID    *string
	SupportsStreaming *bool
	Limit           int
	Offset          int
}

type SearchToolsOptions struct {
	ExecutionType *string
	Stability     *string
	Limit         int
}

// Request types for prompt operations
type CreatePromptRequest struct {
	Name              string
	Namespace         string
	Version           string
	Description       string
	Author            string
	License           string
	Homepage          *string
	Documentation     *string
	Tags              []string
	Categories        []string
	Keywords          []string
	Stability         string
	Maturity          string
	ForgeVersion      string
	Platforms         []string
	Spec              string
	SpecHash          string
	RepositoryID      string
	CommitHash        string
	Branch            string
	TemplatePath      *string
	TemplateContent   *string
	VariablesSchema   map[string]interface{}
	PromptType        string
	ContextWindow     *int
	DefaultVariables  map[string]interface{}
	RequiredVariables []string
	Language          string
	SupportsStreaming bool
	ModelPreferences  []string
	Temperature       *float64
	MaxTokens         *int
	StopSequences     []string
}

type ListPromptsOptions struct {
	PromptType      *string
	Stability       *string
	IsInstalled     *bool
	RepositoryID    *string
	Language        *string
	Limit           int
	Offset          int
}

type SearchPromptsOptions struct {
	PromptType *string
	Stability  *string
	Language   *string
	Limit      int
}

// Request types for agent operations
type CreateAgentRequest struct {
	Name                    string
	Namespace               string
	Version                 string
	Description             string
	Author                  string
	License                 string
	Homepage                *string
	Documentation           *string
	Tags                    []string
	Categories              []string
	Keywords                []string
	Stability               string
	Maturity                string
	ForgeVersion            string
	Platforms               []string
	Spec                    string
	SpecHash                string
	RepositoryID            string
	CommitHash              string
	Branch                  string
	ConfigPath              *string
	AgentConfig             map[string]interface{}
	LLMProvider             *string
	SystemPromptID          *string
	ToolDependencies        []string
	PromptDependencies      []string
	AgentDependencies       []string
	AgentType               string
	Capabilities            []string
	SupportedLanguages      []string
	SupportsMemory          bool
	SupportsTools           bool
	SupportsMultimodal      bool
	ModelPreferences        []string
	DefaultTemperature      *float64
	DefaultMaxTokens        *int
	SessionTimeoutMinutes   int
}

type ListAgentsOptions struct {
	AgentType       *string
	LLMProvider     *string
	Stability       *string
	IsInstalled     *bool
	RepositoryID    *string
	SupportsTools   *bool
	SupportsMemory  *bool
	Limit           int
	Offset          int
}

type SearchAgentsOptions struct {
	AgentType     *string
	LLMProvider   *string
	Stability     *string
	SupportsTools *bool
	Limit         int
}