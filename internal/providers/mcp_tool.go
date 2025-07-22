// Package providers contains MCP tool provider implementation.
package providers

import (
	"context"
	"fmt"
	"sync"

	mcpadapter "github.com/denkhaus/mcp-server-adapter"
	"github.com/tmc/langchaingo/tools"
	"go.uber.org/zap"

	"github.com/denkhaus/agentforge/internal/config"
	"github.com/denkhaus/agentforge/internal/types"
)

// mcpToolProvider is a private implementation of types.ToolProvider interface for MCP servers.
type mcpToolProvider struct {
	log     *zap.Logger
	adapter mcpadapter.MCPAdapter
	config  *config.MCPConfig
	tools   map[string]tools.Tool
	mutex   sync.RWMutex

	// Performance optimizations
	toolsSlice []tools.Tool // Pre-built slice for GetTools()
}

// NewMCPToolProvider creates a new MCP tool provider instance.
func NewMCPToolProvider(log *zap.Logger, cfg *config.MCPConfig) (types.ToolProvider, error) {
	if !cfg.Enabled {
		log.Info("MCP integration disabled")
		return &mcpToolProvider{
			log:    log,
			config: cfg,
			tools:  make(map[string]tools.Tool),
		}, nil
	}

	log.Info("Initializing MCP tool provider",
		zap.String("config_path", cfg.ConfigPath),
		zap.String("tool_prefix", cfg.ToolPrefix))

	// Create MCP adapter
	mcpAdapter, err := mcpadapter.New(mcpadapter.WithConfigPath(cfg.ConfigPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP adapter: %w", err)
	}

	provider := &mcpToolProvider{
		log:     log,
		adapter: mcpAdapter,
		config:  cfg,
		tools:   make(map[string]tools.Tool),
	}

	// Load tools from MCP servers
	if err := provider.loadMCPTools(); err != nil {
		return nil, fmt.Errorf("failed to load MCP tools: %w", err)
	}

	log.Info("MCP tool provider initialized",
		zap.Int("tool_count", len(provider.tools)))

	return provider, nil
}

// GetTools returns all available MCP tools as langchain-go tools.
func (p *mcpToolProvider) GetTools() []tools.Tool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// Return copy of pre-built slice for better performance
	result := make([]tools.Tool, len(p.toolsSlice))
	copy(result, p.toolsSlice)
	return result
}

// GetToolsForAgent returns tools required by the agent from MCP servers.
func (p *mcpToolProvider) GetToolsForAgent(agent types.Agent) ([]tools.Tool, error) {
	requiredTools := agent.GetRequiredTools()

	// Validate that all required tools are available
	if err := p.ValidateAgentRequirements(agent); err != nil {
		return nil, err
	}

	agentTools := make([]tools.Tool, 0, len(requiredTools))

	p.mutex.RLock()
	for _, toolName := range requiredTools {
		if tool, exists := p.tools[toolName]; exists {
			agentTools = append(agentTools, tool)
		}
	}
	p.mutex.RUnlock()

	p.log.Info("MCP tools loaded for agent",
		zap.String("agent", agent.GetName()),
		zap.Strings("required_tools", requiredTools),
		zap.Int("loaded_tools", len(agentTools)))

	return agentTools, nil
}

// ExecuteTool executes a specific MCP tool with string input (JSON format).
func (p *mcpToolProvider) ExecuteTool(ctx context.Context, name string, input string) (string, error) {
	p.mutex.RLock()
	tool, exists := p.tools[name]
	p.mutex.RUnlock()

	if !exists {
		return "", fmt.Errorf("MCP tool %s not found", name)
	}

	p.log.Info("Executing MCP tool",
		zap.String("name", name),
		zap.String("input", input))

	result, err := tool.Call(ctx, input)
	if err != nil {
		p.log.Error("MCP tool execution failed",
			zap.String("name", name),
			zap.Error(err))
		return "", fmt.Errorf("MCP tool execution failed: %w", err)
	}

	p.log.Info("MCP tool execution completed",
		zap.String("name", name),
		zap.String("result", result))
	return result, nil
}

// RegisterTool registers a new MCP tool (not supported for MCP provider).
func (p *mcpToolProvider) RegisterTool(tool tools.Tool) error {
	return fmt.Errorf("registering tools not supported for MCP provider - tools are loaded from MCP servers")
}

// HasTool checks if an MCP tool with the given name is available.
func (p *mcpToolProvider) HasTool(name string) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	_, exists := p.tools[name]
	return exists
}

// ValidateAgentRequirements checks if all required tools for an agent are available from MCP servers.
func (p *mcpToolProvider) ValidateAgentRequirements(agent types.Agent) error {
	requiredTools := agent.GetRequiredTools()
	missingTools := make([]string, 0)

	p.mutex.RLock()
	for _, toolName := range requiredTools {
		if _, exists := p.tools[toolName]; !exists {
			missingTools = append(missingTools, toolName)
		}
	}
	p.mutex.RUnlock()

	if len(missingTools) > 0 {
		return fmt.Errorf("agent %s requires missing MCP tools: %v", agent.GetName(), missingTools)
	}

	return nil
}

// GetToolNames returns the names of all available MCP tools.
func (p *mcpToolProvider) GetToolNames() []string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	names := make([]string, 0, len(p.tools))
	for name := range p.tools {
		names = append(names, name)
	}

	return names
}

// loadMCPTools loads tools from MCP servers via the adapter.
func (p *mcpToolProvider) loadMCPTools() error {
	if p.adapter == nil {
		// MCP disabled, no tools to load
		p.rebuildToolsSlice()
		return nil
	}

	// Get tools from MCP adapter
	mcpTools, err := p.adapter.GetAllLangChainTools(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get tools from MCP adapter: %w", err)
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Convert MCP tools to our internal format
	for _, mcpTool := range mcpTools {
		toolName := mcpTool.Name()

		// Apply prefix if configured
		if p.config.ToolPrefix != "" {
			toolName = p.config.ToolPrefix + toolName
		}

		// Create wrapper tool that implements tools.Tool
		wrappedTool := &mcpToolWrapper{
			mcpTool: mcpTool,
			name:    toolName,
			log:     p.log,
		}

		p.tools[toolName] = wrappedTool
	}

	// Rebuild tools slice for GetTools()
	p.rebuildToolsSlice()

	return nil
}

// rebuildToolsSlice rebuilds the pre-built tools slice for GetTools().
func (p *mcpToolProvider) rebuildToolsSlice() {
	p.toolsSlice = make([]tools.Tool, 0, len(p.tools))
	for _, tool := range p.tools {
		p.toolsSlice = append(p.toolsSlice, tool)
	}
}

// mcpToolWrapper wraps an MCP tool to implement the langchain-go tools.Tool interface.
type mcpToolWrapper struct {
	mcpTool tools.Tool // The underlying MCP tool from the adapter
	name    string     // Potentially prefixed name
	log     *zap.Logger
}

// Name returns the tool name (potentially with prefix).
func (w *mcpToolWrapper) Name() string {
	return w.name
}

// Description returns the tool description from the underlying MCP tool.
func (w *mcpToolWrapper) Description() string {
	return w.mcpTool.Description()
}

// Call executes the underlying MCP tool.
func (w *mcpToolWrapper) Call(ctx context.Context, input string) (string, error) {
	w.log.Debug("Calling MCP tool wrapper",
		zap.String("tool", w.name),
		zap.String("input", input))

	return w.mcpTool.Call(ctx, input)
}
