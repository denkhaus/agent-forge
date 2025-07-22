package startup

import (
	"context"

	"github.com/denkhaus/agentforge/internal/database"
	"github.com/denkhaus/agentforge/internal/prompts"
	cli "github.com/urfave/cli/v2"
)

// StartupOptions defines fine-grained startup configuration
type StartupOptions struct {
	// Core services
	RequireDatabase bool
	RequireGit      bool
	RequireGitHub   bool
	
	// Entity services
	RequirePromptService bool
	RequireToolService   bool
	RequireAgentService  bool
	
	// Additional options
	DatabasePath string
	Context      context.Context
}

// Option represents a startup configuration option
type Option func(*StartupOptions)

// WithDatabase enables database initialization
func WithDatabase() Option {
	return func(opts *StartupOptions) {
		opts.RequireDatabase = true
	}
}

// WithDatabasePath enables database with custom path
func WithDatabasePath(path string) Option {
	return func(opts *StartupOptions) {
		opts.RequireDatabase = true
		opts.DatabasePath = path
	}
}

// WithGit enables Git client initialization
func WithGit() Option {
	return func(opts *StartupOptions) {
		opts.RequireGit = true
	}
}

// WithGitHub enables GitHub client initialization
func WithGitHub() Option {
	return func(opts *StartupOptions) {
		opts.RequireGitHub = true
	}
}

// WithPromptService enables prompt service initialization
func WithPromptService() Option {
	return func(opts *StartupOptions) {
		opts.RequirePromptService = true
	}
}

// WithToolService enables tool service initialization
func WithToolService() Option {
	return func(opts *StartupOptions) {
		opts.RequireToolService = true
	}
}

// WithAgentService enables agent service initialization
func WithAgentService() Option {
	return func(opts *StartupOptions) {
		opts.RequireAgentService = true
	}
}

// WithContext sets a custom context
func WithContext(ctx context.Context) Option {
	return func(opts *StartupOptions) {
		opts.Context = ctx
	}
}

// Context provides access to initialized services and dependencies
type Context struct {
	Context context.Context
	CLI     *cli.Context
	
	// Core services
	DatabaseManager database.DatabaseManager
	GitClient       interface{} // TODO: Define proper interface when git package is ready
	GitHubClient    interface{} // TODO: Define proper interface when github package is ready
	
	// Entity services
	PromptService prompts.PromptService
	// ToolService   tools.ToolService     // TODO: Implement when tools package is ready
	// AgentService  agents.AgentService   // TODO: Implement when agents package is ready
}

// Cleanup releases resources
func (ctx *Context) Cleanup() {
	if ctx.DatabaseManager != nil {
		ctx.DatabaseManager.Shutdown()
	}
}

// Minimal returns minimal startup options (no services)
func Minimal() []Option {
	return []Option{}
}

// Database returns database-only startup options
func Database() []Option {
	return []Option{WithDatabase()}
}

// Full returns all available startup options
func Full() []Option {
	return []Option{
		WithDatabase(),
		WithGit(),
		WithGitHub(),
		WithPromptService(),
		WithToolService(),
		WithAgentService(),
	}
}