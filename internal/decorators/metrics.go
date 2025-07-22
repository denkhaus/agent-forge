package decorators

import (
	"context"
	"sync"
	"time"

	"github.com/tmc/langchaingo/tools"
	"go.uber.org/zap"

	"github.com/denkhaus/agentforge/internal/types"
)

// MetricsToolProviderDecorator adds metrics collection to tool operations.
type MetricsToolProviderDecorator struct {
	inner   types.ToolProvider
	log     *zap.Logger
	metrics *ToolMetrics
}

// ToolMetrics holds metrics data for tool operations.
type ToolMetrics struct {
	mutex           sync.RWMutex
	ExecutionCounts map[string]int64
	ExecutionTimes  map[string][]time.Duration
	ErrorCounts     map[string]int64
	TotalExecutions int64
	TotalErrors     int64
}

// ToolMetricsSnapshot represents a snapshot of metrics without mutex for safe copying.
type ToolMetricsSnapshot struct {
	ExecutionCounts map[string]int64
	ExecutionTimes  map[string][]time.Duration
	ErrorCounts     map[string]int64
	TotalExecutions int64
	TotalErrors     int64
}

// NewToolMetrics creates a new ToolMetrics instance.
func NewToolMetrics() *ToolMetrics {
	return &ToolMetrics{
		ExecutionCounts: make(map[string]int64),
		ExecutionTimes:  make(map[string][]time.Duration),
		ErrorCounts:     make(map[string]int64),
	}
}

// NewMetricsToolProviderDecorator creates a new metrics decorator for tool providers.
func NewMetricsToolProviderDecorator(inner types.ToolProvider, log *zap.Logger) types.ToolProvider {
	return &MetricsToolProviderDecorator{
		inner:   inner,
		log:     log,
		metrics: NewToolMetrics(),
	}
}

// ExecuteTool records metrics for tool execution, including counts, times, and errors.
func (d *MetricsToolProviderDecorator) ExecuteTool(
	ctx context.Context,
	name string,
	input string,
) (string, error) {
	start := time.Now()

	result, err := d.inner.ExecuteTool(ctx, name, input)

	duration := time.Since(start)
	d.recordMetrics(name, duration, err != nil)

	return result, err
}

func (d *MetricsToolProviderDecorator) recordMetrics(toolName string, duration time.Duration, isError bool) {
	d.metrics.mutex.Lock()
	defer d.metrics.mutex.Unlock()

	// Record execution count
	d.metrics.ExecutionCounts[toolName]++
	d.metrics.TotalExecutions++

	// Record execution time
	if d.metrics.ExecutionTimes[toolName] == nil {
		d.metrics.ExecutionTimes[toolName] = make([]time.Duration, 0)
	}
	d.metrics.ExecutionTimes[toolName] = append(d.metrics.ExecutionTimes[toolName], duration)

	// Record errors
	if isError {
		d.metrics.ErrorCounts[toolName]++
		d.metrics.TotalErrors++
	}

	// Log metrics periodically
	if d.metrics.TotalExecutions%10 == 0 {
		d.logMetrics()
	}
}

func (d *MetricsToolProviderDecorator) logMetrics() {
	d.log.Info("Tool execution metrics",
		zap.Int64("total_executions", d.metrics.TotalExecutions),
		zap.Int64("total_errors", d.metrics.TotalErrors),
		zap.Any("execution_counts", d.metrics.ExecutionCounts),
		zap.Any("error_counts", d.metrics.ErrorCounts))
}

// GetMetrics returns a snapshot of the current metrics.
func (d *MetricsToolProviderDecorator) GetMetrics() ToolMetricsSnapshot {
	d.metrics.mutex.RLock()
	defer d.metrics.mutex.RUnlock()

	// Create a deep copy without the mutex
	snapshot := ToolMetricsSnapshot{
		ExecutionCounts: make(map[string]int64),
		ExecutionTimes:  make(map[string][]time.Duration),
		ErrorCounts:     make(map[string]int64),
		TotalExecutions: d.metrics.TotalExecutions,
		TotalErrors:     d.metrics.TotalErrors,
	}

	for k, v := range d.metrics.ExecutionCounts {
		snapshot.ExecutionCounts[k] = v
	}

	for k, v := range d.metrics.ErrorCounts {
		snapshot.ErrorCounts[k] = v
	}

	for k, v := range d.metrics.ExecutionTimes {
		snapshot.ExecutionTimes[k] = make([]time.Duration, len(v))
		copy(snapshot.ExecutionTimes[k], v)
	}

	return snapshot
}

// GetTools returns a list of tools.
func (d *MetricsToolProviderDecorator) GetTools() []tools.Tool {
	return d.inner.GetTools()
}

// GetToolsForAgent returns tools required by the agent, returns error if any tools are missing.
func (d *MetricsToolProviderDecorator) GetToolsForAgent(agent types.Agent) ([]tools.Tool, error) {
	return d.inner.GetToolsForAgent(agent)
}

// Removed ExecuteToolForAgent - use ExecuteTool instead

// RegisterTool registers a standard langchain-go tool.
func (d *MetricsToolProviderDecorator) RegisterTool(tool tools.Tool) error {
	return d.inner.RegisterTool(tool)
}

// HasTool checks if a tool with the given name is available.
func (d *MetricsToolProviderDecorator) HasTool(name string) bool {
	return d.inner.HasTool(name)
}

// ValidateAgentRequirements checks if all required tools for an agent are available.
func (d *MetricsToolProviderDecorator) ValidateAgentRequirements(agent types.Agent) error {
	return d.inner.ValidateAgentRequirements(agent)
}

// GetToolNames returns the names of all available tools.
func (d *MetricsToolProviderDecorator) GetToolNames() []string {
	return d.inner.GetToolNames()
}

// Removed Startup method - Provider interface eliminated
