// Package decorators provides decorator patterns for enhancing providers with cross-cutting concerns.
package decorators

import (
	"context"
	"time"

	"github.com/tmc/langchaingo/tools"
	"go.uber.org/zap"

	"github.com/denkhaus/agentforge/internal/types"
)

// LoggingToolProviderDecorator adds comprehensive logging to tool operations.
type LoggingToolProviderDecorator struct {
	inner types.ToolProvider
	log   *zap.Logger
}

// NewLoggingToolProviderDecorator creates a new logging decorator for tool providers.
func NewLoggingToolProviderDecorator(inner types.ToolProvider, log *zap.Logger) types.ToolProvider {
	return &LoggingToolProviderDecorator{
		inner: inner,
		log:   log,
	}
}

// ExecuteTool logs the start and end of tool execution, including results or errors.
func (d *LoggingToolProviderDecorator) ExecuteTool(
	ctx context.Context,
	name string,
	input string,
) (string, error) {
	start := time.Now()

	d.log.Info("Tool execution started",
		zap.String("tool_name", name),
		zap.String("input", input),
		zap.Time("start_time", start))

	result, err := d.inner.ExecuteTool(ctx, name, input)

	duration := time.Since(start)

	if err != nil {
		d.log.Error("Tool execution failed",
			zap.String("tool_name", name),
			zap.Error(err),
			zap.Duration("duration", duration))
	} else {
		d.log.Info("Tool execution completed",
			zap.String("tool_name", name),
			zap.String("result", result),
			zap.Duration("duration", duration))
	}

	return result, err
}

// GetTools retrieves and logs the available tools.
func (d *LoggingToolProviderDecorator) GetTools() []tools.Tool {
	toolsList := d.inner.GetTools()
	d.log.Debug("Tools retrieved", zap.Int("count", len(toolsList)))
	return toolsList
}

// GetToolsForAgent retrieves and logs tools specific to an agent.
func (d *LoggingToolProviderDecorator) GetToolsForAgent(agent types.Agent) ([]tools.Tool, error) {
	toolsList, err := d.inner.GetToolsForAgent(agent)
	if err != nil {
		d.log.Error("Failed to get tools for agent",
			zap.String("agent", agent.GetName()),
			zap.Error(err))
		return nil, err
	}
	d.log.Debug("Retrieved tools for agent",
		zap.String("agent", agent.GetName()),
		zap.Int("count", len(toolsList)))
	return toolsList, nil
}

// Removed ExecuteToolForAgent - use ExecuteTool instead

// RegisterTool logs the registration of a standard langchain-go tool.
func (d *LoggingToolProviderDecorator) RegisterTool(tool tools.Tool) error {
	d.log.Info("Registering tool",
		zap.String("name", tool.Name()),
		zap.String("description", tool.Description()))

	err := d.inner.RegisterTool(tool)

	if err != nil {
		d.log.Error("Tool registration failed",
			zap.String("name", tool.Name()),
			zap.Error(err))
	} else {
		d.log.Info("Tool registered successfully",
			zap.String("name", tool.Name()))
	}

	return err
}


// HasTool checks if a tool with the given name is available.
func (d *LoggingToolProviderDecorator) HasTool(name string) bool {
	exists := d.inner.HasTool(name)
	d.log.Debug("Tool availability check",
		zap.String("tool_name", name),
		zap.Bool("exists", exists))
	return exists
}

// ValidateAgentRequirements checks if all required tools for an agent are available.
func (d *LoggingToolProviderDecorator) ValidateAgentRequirements(agent types.Agent) error {
	err := d.inner.ValidateAgentRequirements(agent)
	if err != nil {
		d.log.Error("Agent requirements validation failed",
			zap.String("agent", agent.GetName()),
			zap.Error(err))
	} else {
		d.log.Debug("Agent requirements validated successfully",
			zap.String("agent", agent.GetName()),
			zap.Strings("required_tools", agent.GetRequiredTools()))
	}
	return err
}

// GetToolNames returns the names of all available tools.
func (d *LoggingToolProviderDecorator) GetToolNames() []string {
	names := d.inner.GetToolNames()
	d.log.Debug("Retrieved tool names",
		zap.Strings("tool_names", names),
		zap.Int("count", len(names)))
	return names
}

// Removed Startup method - Provider interface eliminated
