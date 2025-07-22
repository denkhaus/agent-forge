package session

import (
	"context"

	"go.uber.org/zap"

	"github.com/denkhaus/agentforge/internal/types"
)

// sessionFactory creates agent sessions with proper dependency injection.
// This is a private implementation of types.SessionFactory interface.
type sessionFactory struct {
	config     types.SessionConfig
	llmService types.LLMService
}

// NewFactory creates a new session factory.
func NewFactory(llmService types.LLMService) types.SessionFactory {
	return &sessionFactory{
		config: types.SessionConfig{
			EnableLogging: true,
		},
		llmService: llmService,
	}
}

// NewFactoryWithConfig creates a new session factory with custom configuration.
func NewFactoryWithConfig(config types.SessionConfig, llmService types.LLMService) types.SessionFactory {
	return &sessionFactory{
		config:     config,
		llmService: llmService,
	}
}

// CreateSession creates a new agent session with flexible options.
func (f *sessionFactory) CreateSession(ctx context.Context, opts types.SessionOptions) (types.AgentSession, error) {
	if f.config.EnableLogging {
		log.Info("Creating session",
			zap.String("agent_name", opts.Agent.GetName()),
			zap.String("agent_type", opts.AgentType),
		)
	}

	llm, err := f.llmService.InitializeLLM(ctx, opts.Config, opts.Agent.GetLLMConfig())
	if err != nil {
		return nil, err
	}

	// Use session config from options if provided, otherwise use default
	agentSessionConfig := types.NewAgentSessionConfig()
	if opts.SessionConfig != nil {
		// Update factory's session config for logging preferences
		f.config = *opts.SessionConfig
	}

	return NewAgentSessionWithConfig(opts.Config, opts.Agent, llm, opts.ToolProvider, opts.AgentProvider, f.llmService, agentSessionConfig)
}
