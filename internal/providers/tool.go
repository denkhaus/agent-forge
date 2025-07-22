// Package providers contains implementations of various provider interfaces.
package providers

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/samber/do"
	"github.com/tmc/langchaingo/tools"
	"go.uber.org/zap"

	"github.com/denkhaus/agentforge/internal/types"
)

// toolProvider is a private implementation of types.ToolProvider interface.
type toolProvider struct {
	log   *zap.Logger
	tools map[string]tools.Tool
	mutex sync.RWMutex

	// Performance optimizations
	toolsSlice []tools.Tool            // Pre-built slice for GetTools()
	toolsCache map[string][]tools.Tool // Cache tools by agent required tools hash
}

// NewToolProvider creates a new tool provider instance.
func NewToolProvider(injector *do.Injector) (types.ToolProvider, error) {
	log, err := do.Invoke[*zap.Logger](injector)
	if err != nil {
		return nil, fmt.Errorf("failed to get logger: %w", err)
	}
	
	toolsSlice, err := do.Invoke[[]tools.Tool](injector)
	if err != nil {
		return nil, fmt.Errorf("failed to get tools slice: %w", err)
	}
	
	// Convert slice to map for efficient lookup
	toolsMap := make(map[string]tools.Tool)
	for _, tool := range toolsSlice {
		toolsMap[tool.Name()] = tool
	}

	provider := &toolProvider{
		log:        log,
		tools:      toolsMap,
		toolsCache: make(map[string][]tools.Tool),
	}

	// Pre-build tools slice for better performance
	provider.rebuildToolsSlice()

	provider.log.Info("Tool provider initialized", zap.Int("tool_count", len(provider.tools)))

	return provider, nil
}

// Removed Startup method - Provider interface eliminated

// GetTools returns all available tools as langchain-go tools.
func (p *toolProvider) GetTools() []tools.Tool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// Return pre-built slice for better performance
	result := make([]tools.Tool, len(p.toolsSlice))
	copy(result, p.toolsSlice)
	return result
}

// GetToolsForAgent returns tools required by the agent, returns error if any tools are missing.
func (p *toolProvider) GetToolsForAgent(agent types.Agent) ([]tools.Tool, error) {
	requiredTools := agent.GetRequiredTools()

	// Create cache key from required tools
	cacheKey := p.createRequiredToolsHash(requiredTools)

	p.mutex.RLock()
	// Check cache first
	if cachedTools, exists := p.toolsCache[cacheKey]; exists {
		p.mutex.RUnlock()
		// Return copy of cached tools
		result := make([]tools.Tool, len(cachedTools))
		copy(result, cachedTools)
		return result, nil
	}
	p.mutex.RUnlock()

	// Cache miss - validate and build tools list
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

	// Cache the result
	p.mutex.Lock()
	p.toolsCache[cacheKey] = agentTools
	p.mutex.Unlock()

	p.log.Info("Tools loaded for agent",
		zap.String("agent", agent.GetName()),
		zap.Strings("required_tools", requiredTools),
		zap.Int("loaded_tools", len(agentTools)))

	// Return copy
	result := make([]tools.Tool, len(agentTools))
	copy(result, agentTools)
	return result, nil
}

// HasTool checks if a tool with the given name is available.
func (p *toolProvider) HasTool(name string) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	
	_, exists := p.tools[name]
	return exists
}

// ValidateAgentRequirements checks if all required tools for an agent are available.
func (p *toolProvider) ValidateAgentRequirements(agent types.Agent) error {
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
		return fmt.Errorf("agent %s requires missing tools: %v", agent.GetName(), missingTools)
	}
	
	return nil
}

// GetToolNames returns the names of all available tools.
func (p *toolProvider) GetToolNames() []string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	
	names := make([]string, 0, len(p.tools))
	for name := range p.tools {
		names = append(names, name)
	}
	
	return names
}

// ExecuteTool executes a specific tool with string input (JSON format).
func (p *toolProvider) ExecuteTool(ctx context.Context, name string, input string) (string, error) {
	p.mutex.RLock()
	tool, exists := p.tools[name]
	p.mutex.RUnlock()

	if !exists {
		return "", fmt.Errorf("tool %s not found", name)
	}

	p.log.Info("Executing tool",
		zap.String("name", name),
		zap.String("input", input))

	result, err := tool.Call(ctx, input)
	if err != nil {
		p.log.Error("Tool execution failed",
			zap.String("name", name),
			zap.Error(err))
		return "", fmt.Errorf("tool execution failed: %w", err)
	}

	p.log.Info("Tool execution completed",
		zap.String("name", name),
		zap.String("result", result))
	return result, nil
}

// Removed ExecuteToolForAgent - capability validation now handled automatically in ExecuteTool

// RegisterTool registers a standard langchain-go tool.
func (p *toolProvider) RegisterTool(tool tools.Tool) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	name := tool.Name()
	if _, exists := p.tools[name]; exists {
		return fmt.Errorf("tool %s already exists", name)
	}

	p.tools[name] = tool

	// Invalidate caches when new tool is registered
	p.toolsCache = make(map[string][]tools.Tool)
	p.rebuildToolsSlice()

	p.log.Info("Tool registered",
		zap.String("name", name),
		zap.String("description", tool.Description()))

	return nil
}

// rebuildToolsSlice rebuilds the pre-built tools slice for GetTools().
func (p *toolProvider) rebuildToolsSlice() {
	p.toolsSlice = make([]tools.Tool, 0, len(p.tools))
	for _, tool := range p.tools {
		p.toolsSlice = append(p.toolsSlice, tool)
	}
}

// createRequiredToolsHash creates a hash key from agent required tools for caching.
func (p *toolProvider) createRequiredToolsHash(requiredTools []string) string {
	if len(requiredTools) == 0 {
		return "no_tools"
	}

	// Sort tools for consistent hashing
	sorted := make([]string, len(requiredTools))
	copy(sorted, requiredTools)
	slices.Sort(sorted)

	// Simple string concatenation for hash (could use crypto/hash for better distribution)
	hash := ""
	for i, tool := range sorted {
		if i > 0 {
			hash += "|"
		}
		hash += tool
	}
	return hash
}
