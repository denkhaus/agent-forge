package types

import "context"

// Factory interfaces for essential functionality only

// SessionFactory defines the interface for creating sessions.
type SessionFactory interface {
	// CreateSession creates a new agent session with flexible options
	CreateSession(ctx context.Context, opts SessionOptions) (AgentSession, error)
}

// SessionOptions contains all options for session creation.
type SessionOptions struct {
	Config        Config
	Agent         Agent
	ToolProvider  ToolProvider
	AgentProvider AgentProvider
	AgentType     string         // optional agent type hint for optimization
	SessionConfig *SessionConfig // optional session configuration override
}

// Removed FactoryManager - over-engineering, use DI container directly

// SessionConfig represents simplified configuration for session creation.
type SessionConfig struct {
	// EnableLogging enables logging for session operations
	EnableLogging bool
}
