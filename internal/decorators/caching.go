// Package decorators provides various decorators for enhancing functionality.
package decorators

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/tmc/langchaingo/tools"
	"go.uber.org/zap"

	"github.com/denkhaus/agentforge/internal/types"
)

// CacheEntry represents a cached tool execution result.
type CacheEntry struct {
	Result    string
	Error     error
	Timestamp time.Time
	TTL       time.Duration
}

// IsExpired checks if the cache entry has expired.
func (e *CacheEntry) IsExpired() bool {
	return time.Since(e.Timestamp) > e.TTL
}

// CachingToolProviderDecorator adds caching to tool operations.
type CachingToolProviderDecorator struct {
	inner      types.ToolProvider
	log        *zap.Logger
	cache      map[string]*CacheEntry
	mutex      sync.RWMutex
	defaultTTL time.Duration
}

// NewCachingToolProviderDecorator creates a new caching decorator for tool providers.
func NewCachingToolProviderDecorator(
	inner types.ToolProvider,
	log *zap.Logger,
	defaultTTL time.Duration,
) types.ToolProvider {
	if defaultTTL == 0 {
		defaultTTL = 5 * time.Minute // Default 5 minutes
	}

	return &CachingToolProviderDecorator{
		inner:      inner,
		log:        log,
		cache:      make(map[string]*CacheEntry),
		defaultTTL: defaultTTL,
	}
}

// ExecuteTool executes a tool, returning a cached result if available and not expired.
func (d *CachingToolProviderDecorator) ExecuteTool(
	ctx context.Context,
	name string,
	input string,
) (string, error) {
	// Generate cache key
	cacheKey := d.generateCacheKey(name, input)

	// Check cache first
	d.mutex.RLock()
	if entry, exists := d.cache[cacheKey]; exists && !entry.IsExpired() {
		d.mutex.RUnlock()
		d.log.Debug("Tool result served from cache",
			zap.String("tool_name", name),
			zap.String("cache_key", cacheKey))
		return entry.Result, entry.Error
	}
	d.mutex.RUnlock()

	// Execute tool
	result, err := d.inner.ExecuteTool(ctx, name, input)

	// Cache the result
	d.mutex.Lock()
	d.cache[cacheKey] = &CacheEntry{
		Result:    result,
		Error:     err,
		Timestamp: time.Now(),
		TTL:       d.defaultTTL,
	}
	d.mutex.Unlock()

	d.log.Debug("Tool result cached",
		zap.String("tool_name", name),
		zap.String("cache_key", cacheKey),
		zap.Duration("ttl", d.defaultTTL))

	return result, err
}

func (d *CachingToolProviderDecorator) generateCacheKey(name string, input string) string {
	// Create a deterministic key from tool name and input
	data := map[string]string{
		"tool":  name,
		"input": input,
	}

	jsonData, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonData)
	return fmt.Sprintf("%x", hash)
}

// ClearCache removes all cached entries.
func (d *CachingToolProviderDecorator) ClearCache() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	count := len(d.cache)
	d.cache = make(map[string]*CacheEntry)

	d.log.Info("Cache cleared", zap.Int("entries_removed", count))
}

// ClearExpiredEntries removes expired cache entries.
func (d *CachingToolProviderDecorator) ClearExpiredEntries() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	var removed int
	for key, entry := range d.cache {
		if entry.IsExpired() {
			delete(d.cache, key)
			removed++
		}
	}

	if removed > 0 {
		d.log.Debug("Expired cache entries removed", zap.Int("count", removed))
	}
}

// GetCacheStats returns cache statistics.
func (d *CachingToolProviderDecorator) GetCacheStats() map[string]any {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	var expired int
	for _, entry := range d.cache {
		if entry.IsExpired() {
			expired++
		}
	}

	return map[string]any{
		"total_entries":   len(d.cache),
		"expired_entries": expired,
		"active_entries":  len(d.cache) - expired,
		"default_ttl":     d.defaultTTL.String(),
	}
}

// GetTools returns a list of tools.
func (d *CachingToolProviderDecorator) GetTools() []tools.Tool {
	return d.inner.GetTools()
}

// GetToolsForAgent returns tools required by the agent, returns error if any tools are missing.
func (d *CachingToolProviderDecorator) GetToolsForAgent(agent types.Agent) ([]tools.Tool, error) {
	return d.inner.GetToolsForAgent(agent)
}

// Removed ExecuteToolForAgent - use ExecuteTool instead

// RegisterTool registers a standard langchain-go tool.
func (d *CachingToolProviderDecorator) RegisterTool(tool tools.Tool) error {
	return d.inner.RegisterTool(tool)
}

// HasTool checks if a tool with the given name is available.
func (d *CachingToolProviderDecorator) HasTool(name string) bool {
	return d.inner.HasTool(name)
}

// ValidateAgentRequirements checks if all required tools for an agent are available.
func (d *CachingToolProviderDecorator) ValidateAgentRequirements(agent types.Agent) error {
	return d.inner.ValidateAgentRequirements(agent)
}

// GetToolNames returns the names of all available tools.
func (d *CachingToolProviderDecorator) GetToolNames() []string {
	return d.inner.GetToolNames()
}

// Removed Startup method - Provider interface eliminated
// Cache cleanup can be handled separately if needed
