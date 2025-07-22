// Package providers contains the aggregated tool provider implementation.
package providers

import (
	"context"
	"fmt"
	"sync"

	"github.com/tmc/langchaingo/tools"
	"go.uber.org/zap"

	"github.com/denkhaus/agentforge/internal/types"
)

// aggregatedToolProvider implements types.ToolProvider by aggregating multiple tool providers.
// aggregatedToolProvider is a private implementation of types.ToolProvider interface that aggregates multiple providers.
type aggregatedToolProvider struct {
	log       *zap.Logger
	providers []types.ToolProvider
	toolCache map[string]types.ToolProvider // tool name -> provider mapping
	mutex     sync.RWMutex

	// Performance optimizations
	allTools   []tools.Tool // Cached result of all tools
	toolNames  []string     // Cached tool names
	cacheValid bool         // Whether cache is valid
}

// NewAggregatedToolProvider creates a new aggregated tool provider.
func NewAggregatedToolProvider(log *zap.Logger, providers ...types.ToolProvider) types.ToolProvider {
	provider := &aggregatedToolProvider{
		log:       log,
		providers: providers,
		toolCache: make(map[string]types.ToolProvider),
	}

	// Build initial cache
	if err := provider.rebuildCache(); err != nil {
		log.Error("Failed to build initial tool cache", zap.Error(err))
	}

	log.Info("Aggregated tool provider initialized",
		zap.Int("provider_count", len(providers)),
		zap.Int("total_tools", len(provider.allTools)))

	return provider
}

// GetTools returns all available tools from all providers.
func (p *aggregatedToolProvider) GetTools() []tools.Tool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if !p.cacheValid {
		p.mutex.RUnlock()
		if err := p.rebuildCache(); err != nil {
			p.log.Error("Failed to rebuild tool cache", zap.Error(err))
			return []tools.Tool{}
		}
		p.mutex.RLock()
	}

	// Return copy of cached tools
	result := make([]tools.Tool, len(p.allTools))
	copy(result, p.allTools)
	return result
}

// GetToolsForAgent returns tools required by the agent from all providers.
func (p *aggregatedToolProvider) GetToolsForAgent(agent types.Agent) ([]tools.Tool, error) {
	requiredTools := agent.GetRequiredTools()
	
	// Validate that all required tools are available
	if err := p.ValidateAgentRequirements(agent); err != nil {
		return nil, err
	}

	agentTools := make([]tools.Tool, 0, len(requiredTools))

	p.mutex.RLock()
	for _, toolName := range requiredTools {
		if provider, exists := p.toolCache[toolName]; exists {
			// Get the specific tool from the provider
			providerTools := provider.GetTools()
			for _, tool := range providerTools {
				if tool.Name() == toolName {
					agentTools = append(agentTools, tool)
					break
				}
			}
		}
	}
	p.mutex.RUnlock()

	p.log.Info("Aggregated tools loaded for agent",
		zap.String("agent", agent.GetName()),
		zap.Strings("required_tools", requiredTools),
		zap.Int("loaded_tools", len(agentTools)))

	return agentTools, nil
}

// ExecuteTool executes a specific tool by routing to the appropriate provider.
func (p *aggregatedToolProvider) ExecuteTool(ctx context.Context, name string, input string) (string, error) {
	p.mutex.RLock()
	provider, exists := p.toolCache[name]
	p.mutex.RUnlock()

	if !exists {
		return "", fmt.Errorf("tool %s not found in any provider", name)
	}

	p.log.Info("Routing tool execution to provider",
		zap.String("tool", name),
		zap.String("provider_type", fmt.Sprintf("%T", provider)))

	return provider.ExecuteTool(ctx, name, input)
}

// RegisterTool registers a tool with the first provider that supports registration.
func (p *aggregatedToolProvider) RegisterTool(tool tools.Tool) error {
	// Try to register with each provider until one succeeds
	var lastErr error
	for i, provider := range p.providers {
		err := provider.RegisterTool(tool)
		if err == nil {
			p.log.Info("Tool registered with provider",
				zap.String("tool", tool.Name()),
				zap.Int("provider_index", i))
			
			// Invalidate cache since we added a new tool
			p.mutex.Lock()
			p.cacheValid = false
			p.mutex.Unlock()
			
			return nil
		}
		lastErr = err
	}

	return fmt.Errorf("failed to register tool %s with any provider: %w", tool.Name(), lastErr)
}

// HasTool checks if a tool with the given name is available in any provider.
func (p *aggregatedToolProvider) HasTool(name string) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	
	_, exists := p.toolCache[name]
	return exists
}

// ValidateAgentRequirements checks if all required tools for an agent are available.
func (p *aggregatedToolProvider) ValidateAgentRequirements(agent types.Agent) error {
	requiredTools := agent.GetRequiredTools()
	missingTools := make([]string, 0)
	
	p.mutex.RLock()
	for _, toolName := range requiredTools {
		if _, exists := p.toolCache[toolName]; !exists {
			missingTools = append(missingTools, toolName)
		}
	}
	p.mutex.RUnlock()
	
	if len(missingTools) > 0 {
		return fmt.Errorf("agent %s requires missing tools: %v", agent.GetName(), missingTools)
	}
	
	return nil
}

// GetToolNames returns the names of all available tools from all providers.
func (p *aggregatedToolProvider) GetToolNames() []string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if !p.cacheValid {
		p.mutex.RUnlock()
		if err := p.rebuildCache(); err != nil {
			p.log.Error("Failed to rebuild tool cache", zap.Error(err))
			return []string{}
		}
		p.mutex.RLock()
	}

	// Return copy of cached tool names
	result := make([]string, len(p.toolNames))
	copy(result, p.toolNames)
	return result
}

// rebuildCache rebuilds the tool cache from all providers.
func (p *aggregatedToolProvider) rebuildCache() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Clear existing cache
	p.toolCache = make(map[string]types.ToolProvider)
	p.allTools = p.allTools[:0] // Reset slice but keep capacity
	p.toolNames = p.toolNames[:0]

	// Collect tools from all providers
	for _, provider := range p.providers {
		tools := provider.GetTools()
		
		for _, tool := range tools {
			toolName := tool.Name()
			
			// Check for tool name conflicts
			if existingProvider, exists := p.toolCache[toolName]; exists {
				p.log.Warn("Tool name conflict detected",
					zap.String("tool", toolName),
					zap.String("existing_provider", fmt.Sprintf("%T", existingProvider)),
					zap.String("new_provider", fmt.Sprintf("%T", provider)))
				// Keep the first provider that registered the tool
				continue
			}
			
			// Add to cache
			p.toolCache[toolName] = provider
			p.allTools = append(p.allTools, tool)
			p.toolNames = append(p.toolNames, toolName)
		}
	}

	p.cacheValid = true

	p.log.Debug("Tool cache rebuilt",
		zap.Int("total_tools", len(p.allTools)),
		zap.Int("provider_count", len(p.providers)))

	return nil
}

// AddProvider adds a new tool provider to the aggregation.
func (p *aggregatedToolProvider) AddProvider(provider types.ToolProvider) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.providers = append(p.providers, provider)
	p.cacheValid = false // Invalidate cache

	p.log.Info("Tool provider added to aggregation",
		zap.String("provider_type", fmt.Sprintf("%T", provider)),
		zap.Int("total_providers", len(p.providers)))
}

// RemoveProvider removes a tool provider from the aggregation.
func (p *aggregatedToolProvider) RemoveProvider(provider types.ToolProvider) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for i, prov := range p.providers {
		if prov == provider {
			// Remove provider from slice
			p.providers = append(p.providers[:i], p.providers[i+1:]...)
			p.cacheValid = false // Invalidate cache
			
			p.log.Info("Tool provider removed from aggregation",
				zap.String("provider_type", fmt.Sprintf("%T", provider)),
				zap.Int("remaining_providers", len(p.providers)))
			break
		}
	}
}