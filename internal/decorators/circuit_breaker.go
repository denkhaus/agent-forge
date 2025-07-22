// Package decorators provides various decorators for enhancing functionality.
package decorators

import (
	"context"
	"sync"
	"time"

	"github.com/tmc/langchaingo/tools"
	"go.uber.org/zap"

	"github.com/denkhaus/agentforge/internal/errors"
	"github.com/denkhaus/agentforge/internal/types"
)

// CircuitState represents the state of the circuit breaker.
type CircuitState int

const (
	// CircuitClosed - normal operation.
	CircuitClosed CircuitState = iota
	// CircuitOpen - circuit is open, requests are rejected.
	CircuitOpen
	// CircuitHalfOpen - testing if the service has recovered.
	CircuitHalfOpen
)

func (s CircuitState) String() string {
	switch s {
	case CircuitClosed:
		return "CLOSED"
	case CircuitOpen:
		return "OPEN"
	case CircuitHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreakerConfig holds configuration for the circuit breaker.
type CircuitBreakerConfig struct {
	FailureThreshold int           // Number of failures before opening circuit
	RecoveryTimeout  time.Duration // Time to wait before attempting recovery
	SuccessThreshold int           // Number of successes needed to close circuit in half-open state
}

// DefaultCircuitBreakerConfig returns a default configuration.
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		FailureThreshold: 5,
		RecoveryTimeout:  30 * time.Second,
		SuccessThreshold: 3,
	}
}

// CircuitBreakerToolProviderDecorator adds circuit breaker pattern to tool operations.
type CircuitBreakerToolProviderDecorator struct {
	inner           types.ToolProvider
	log             *zap.Logger
	config          CircuitBreakerConfig
	state           CircuitState
	failureCount    int
	successCount    int
	lastFailureTime time.Time
	mutex           sync.RWMutex
}

// NewCircuitBreakerToolProviderDecorator creates a new circuit breaker decorator.
func NewCircuitBreakerToolProviderDecorator(
	inner types.ToolProvider,
	log *zap.Logger,
	config CircuitBreakerConfig,
) types.ToolProvider {
	return &CircuitBreakerToolProviderDecorator{
		inner:  inner,
		log:    log,
		config: config,
		state:  CircuitClosed,
	}
}

// NewCircuitBreakerToolProviderDecoratorWithDefaults creates a circuit breaker with default config.
func NewCircuitBreakerToolProviderDecoratorWithDefaults(inner types.ToolProvider, log *zap.Logger) types.ToolProvider {
	return NewCircuitBreakerToolProviderDecorator(inner, log, DefaultCircuitBreakerConfig())
}

// ExecuteTool executes a tool, applying circuit breaker logic.
func (d *CircuitBreakerToolProviderDecorator) ExecuteTool(
	ctx context.Context,
	name string,
	input string,
) (string, error) {
	// Check if circuit allows execution
	if !d.canExecute() {
		d.log.Warn("Circuit breaker is open, rejecting tool execution",
			zap.String("tool_name", name),
			zap.String("circuit_state", d.getState().String()))
		return "", errors.ErrCircuitBreakerOpen
	}

	// Execute the tool
	result, err := d.inner.ExecuteTool(ctx, name, input)

	// Record the result
	d.recordResult(name, err == nil)

	return result, err
}

func (d *CircuitBreakerToolProviderDecorator) canExecute() bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	switch d.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		// Check if recovery timeout has passed
		if time.Since(d.lastFailureTime) > d.config.RecoveryTimeout {
			// Transition to half-open
			d.mutex.RUnlock()
			d.mutex.Lock()
			if d.state == CircuitOpen { // Double-check after acquiring write lock
				d.state = CircuitHalfOpen
				d.successCount = 0
				d.log.Info("Circuit breaker transitioning to half-open state")
			}
			d.mutex.Unlock()
			d.mutex.RLock()
			return true
		}
		return false
	case CircuitHalfOpen:
		return true
	default:
		return false
	}
}

func (d *CircuitBreakerToolProviderDecorator) recordResult(toolName string, success bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if success {
		d.handleSuccess(toolName)
	} else {
		d.handleFailure(toolName)
	}
}

func (d *CircuitBreakerToolProviderDecorator) handleSuccess(toolName string) {
	d.failureCount = 0
	if d.state == CircuitHalfOpen {
		d.successCount++
		if d.successCount >= d.config.SuccessThreshold {
			d.state = CircuitClosed
			d.log.Info("Circuit breaker closed after successful recovery",
				zap.String("tool_name", toolName),
				zap.Int("success_count", d.successCount))
		}
	}
}

func (d *CircuitBreakerToolProviderDecorator) handleFailure(toolName string) {
	d.failureCount++
	d.lastFailureTime = time.Now()
	d.successCount = 0 // Reset success count on failure

	// Determine new state based on current state and failure count
	if d.state == CircuitClosed && d.failureCount >= d.config.FailureThreshold {
		d.state = CircuitOpen
		d.log.Error("Circuit breaker opened due to failures",
			zap.String("tool_name", toolName),
			zap.Int("failure_count", d.failureCount),
			zap.Int("failure_threshold", d.config.FailureThreshold))
	} else if d.state == CircuitHalfOpen {
		d.state = CircuitOpen // Reopen immediately on failure in half-open state
		d.log.Error("Circuit breaker reopened during recovery attempt",
			zap.String("tool_name", toolName))
	}
}

func (d *CircuitBreakerToolProviderDecorator) getState() CircuitState {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.state
}

// GetCircuitBreakerStats returns current circuit breaker statistics.
func (d *CircuitBreakerToolProviderDecorator) GetCircuitBreakerStats() map[string]any {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return map[string]any{
		"state":             d.state.String(),
		"failure_count":     d.failureCount,
		"success_count":     d.successCount,
		"failure_threshold": d.config.FailureThreshold,
		"success_threshold": d.config.SuccessThreshold,
		"recovery_timeout":  d.config.RecoveryTimeout.String(),
		"last_failure_time": d.lastFailureTime,
	}
}

// ResetCircuitBreaker manually resets the circuit breaker to closed state.
func (d *CircuitBreakerToolProviderDecorator) ResetCircuitBreaker() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	oldState := d.state
	d.state = CircuitClosed
	d.failureCount = 0
	d.successCount = 0

	d.log.Info("Circuit breaker manually reset",
		zap.String("previous_state", oldState.String()),
		zap.String("new_state", d.state.String()))
}

// GetTools returns a list of tools.
func (d *CircuitBreakerToolProviderDecorator) GetTools() []tools.Tool {
	return d.inner.GetTools()
}

// GetToolsForAgent returns tools required by the agent, returns error if any tools are missing.
func (d *CircuitBreakerToolProviderDecorator) GetToolsForAgent(agent types.Agent) ([]tools.Tool, error) {
	return d.inner.GetToolsForAgent(agent)
}

// Removed ExecuteToolForAgent - use ExecuteTool instead

// RegisterTool registers a standard langchain-go tool.
func (d *CircuitBreakerToolProviderDecorator) RegisterTool(tool tools.Tool) error {
	return d.inner.RegisterTool(tool)
}

// HasTool checks if a tool with the given name is available.
func (d *CircuitBreakerToolProviderDecorator) HasTool(name string) bool {
	return d.inner.HasTool(name)
}

// ValidateAgentRequirements checks if all required tools for an agent are available.
func (d *CircuitBreakerToolProviderDecorator) ValidateAgentRequirements(agent types.Agent) error {
	return d.inner.ValidateAgentRequirements(agent)
}

// GetToolNames returns the names of all available tools.
func (d *CircuitBreakerToolProviderDecorator) GetToolNames() []string {
	return d.inner.GetToolNames()
}

// Removed Startup method - Provider interface eliminated
