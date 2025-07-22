package decorators

import (
	"context"
	"errors"
	"strings" // Added import
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	internalErrors "github.com/denkhaus/agentforge/internal/errors"
)

// Mocks from test_mocks.go
// mockToolProvider
// mockAgent
// mockToolImpl

func TestCircuitState_String(t *testing.T) {
	tests := []struct {
		state    CircuitState
		expected string
	}{
		{CircuitClosed, "CLOSED"},
		{CircuitOpen, "OPEN"},
		{CircuitHalfOpen, "HALF_OPEN"},
		{CircuitState(99), "UNKNOWN"}, // Test unknown state
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.state.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, tt.state.String())
			}
		})
	}
}

func TestDefaultCircuitBreakerConfig(t *testing.T) {
	cfg := DefaultCircuitBreakerConfig()
	if cfg.FailureThreshold != 5 {
		t.Errorf("Expected FailureThreshold 5, got %d", cfg.FailureThreshold)
	}
	if cfg.RecoveryTimeout != 30*time.Second {
		t.Errorf("Expected RecoveryTimeout 30s, got %v", cfg.RecoveryTimeout)
	}
	if cfg.SuccessThreshold != 3 {
		t.Errorf("Expected SuccessThreshold 3, got %d", cfg.SuccessThreshold)
	}
}

func TestNewCircuitBreakerToolProviderDecorator(t *testing.T) {
	mockInner := &mockToolProvider{}
	log := zap.NewNop()
	cfg := CircuitBreakerConfig{FailureThreshold: 2, RecoveryTimeout: 1 * time.Minute, SuccessThreshold: 1}

	decorator := NewCircuitBreakerToolProviderDecorator(mockInner, log, cfg)
	cbDecorator, ok := decorator.(*CircuitBreakerToolProviderDecorator)
	if !ok {
		t.Fatal("Expected decorator to be CircuitBreakerToolProviderDecorator")
	}

	if cbDecorator.inner != mockInner {
		t.Error("Inner tool provider not set correctly")
	}
	if cbDecorator.log != log {
		t.Error("Logger not set correctly")
	}
	if cbDecorator.config != cfg {
		t.Error("Config not set correctly")
	}
	if cbDecorator.state != CircuitClosed {
		t.Errorf("Expected initial state CLOSED, got %s", cbDecorator.state)
	}
}

func TestNewCircuitBreakerToolProviderDecoratorWithDefaults(t *testing.T) {
	mockInner := &mockToolProvider{}
	log := zap.NewNop()

	decorator := NewCircuitBreakerToolProviderDecoratorWithDefaults(mockInner, log)
	cbDecorator, ok := decorator.(*CircuitBreakerToolProviderDecorator)
	if !ok {
		t.Fatal("Expected decorator to be CircuitBreakerToolProviderDecorator")
	}

	defaultCfg := DefaultCircuitBreakerConfig()
	if cbDecorator.config != defaultCfg {
		t.Error("Default config not applied correctly")
	}
}

func TestCircuitBreaker_ClosedToOpen(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	log := zap.New(core)
	cfg := CircuitBreakerConfig{FailureThreshold: 2, RecoveryTimeout: 1 * time.Minute, SuccessThreshold: 1}
	mockInner := &mockToolProvider{
		executeToolFunc: func(_ context.Context, _ string, _ string) (string, error) {
			return "", errors.New("simulated error")
		},
	}
	decorator := NewCircuitBreakerToolProviderDecorator(mockInner, log, cfg)
	cbDecorator := decorator.(*CircuitBreakerToolProviderDecorator)

	// First failure - still closed
	_, err := decorator.ExecuteTool(context.Background(), "test_tool", "")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if cbDecorator.getState() != CircuitClosed {
		t.Errorf("Expected state CLOSED, got %s", cbDecorator.getState())
	}
	if cbDecorator.failureCount != 1 {
		t.Errorf("Expected failureCount 1, got %d", cbDecorator.failureCount)
	}

	// Second failure - transitions to Open
	_, err = decorator.ExecuteTool(context.Background(), "test_tool", "")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if cbDecorator.getState() != CircuitOpen {
		t.Errorf("Expected state OPEN, got %s", cbDecorator.getState())
	}
	if cbDecorator.failureCount != 2 {
		t.Errorf("Expected failureCount 2, got %d", cbDecorator.failureCount)
	}

	// Verify log message for state change
	foundLog := false
	for _, entry := range recorded.All() {
		if strings.Contains(entry.Message, "Circuit breaker opened due to failures") {
			foundLog = true
			break
		}
	}
	if !foundLog {
		t.Error("Expected log message 'Circuit breaker opened due to failures' not found")
	}
}

func TestCircuitBreaker_OpenToHalfOpenToClosed(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	log := zap.New(core)
	cfg := CircuitBreakerConfig{FailureThreshold: 1, RecoveryTimeout: 10 * time.Millisecond, SuccessThreshold: 1}
	mockInner := &mockToolProvider{} // Will return success by default
	decorator := NewCircuitBreakerToolProviderDecorator(mockInner, log, cfg)
	cbDecorator := decorator.(*CircuitBreakerToolProviderDecorator)

	// Force open state
	cbDecorator.state = CircuitOpen
	cbDecorator.failureCount = 1
	cbDecorator.lastFailureTime = time.Now().Add(-20 * time.Millisecond) // Ensure timeout passed

	// Attempt execution - should transition to Half-Open
	_, err := decorator.ExecuteTool(context.Background(), "test_tool", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if cbDecorator.getState() != CircuitClosed { // Should transition directly to Closed if success threshold met
		t.Errorf("Expected state CLOSED, got %s", cbDecorator.getState())
	}
	if cbDecorator.successCount != 1 {
		t.Errorf("Expected successCount 1, got %d", cbDecorator.successCount)
	}

	// Verify log messages
	foundHalfOpenLog := false
	foundClosedLog := false
	for _, entry := range recorded.All() {
		if strings.Contains(entry.Message, "Circuit breaker transitioning to half-open state") {
			foundHalfOpenLog = true
		}
		if strings.Contains(entry.Message, "Circuit breaker closed after successful recovery") {
			foundClosedLog = true
		}
	}
	if !foundHalfOpenLog {
		t.Error("Expected log message 'Circuit breaker transitioning to half-open state' not found")
	}
	if !foundClosedLog {
		t.Error("Expected log message 'Circuit breaker closed after successful recovery' not found")
	}
}

func TestCircuitBreaker_HalfOpenToClosed(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	log := zap.New(core)
	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		RecoveryTimeout:  10 * time.Millisecond,
		SuccessThreshold: 1, // Only one success needed to close
	}

	decorator := NewCircuitBreakerToolProviderDecorator(nil, log, cfg)
	cbDecorator := decorator.(*CircuitBreakerToolProviderDecorator)

	mockInner := &mockToolProvider{
		executeToolFunc: func(_ context.Context, _ string, _ string) (string, error) {
			return "success", nil
		},
	}
	cbDecorator.inner = mockInner

	// Force open state, ensure timeout has passed
	cbDecorator.state = CircuitOpen
	cbDecorator.failureCount = 1
	cbDecorator.lastFailureTime = time.Now().Add(-20 * time.Millisecond)

	// Attempt execution - should transition to Half-Open and then immediately to Closed if successful
	_, err := decorator.ExecuteTool(context.Background(), "test_tool", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if cbDecorator.getState() != CircuitClosed {
		t.Errorf("Expected state CLOSED, got %s", cbDecorator.getState())
	}
	if cbDecorator.successCount != 1 {
		t.Errorf("Expected successCount 1, got %d", cbDecorator.successCount)
	}

	// Verify log messages
	if !strings.Contains(recorded.All()[0].Message, "Circuit breaker transitioning to half-open state") {
		t.Error("Expected log message 'Circuit breaker transitioning to half-open state' not found")
	}
	if !strings.Contains(recorded.All()[1].Message, "Circuit breaker closed after successful recovery") {
		t.Error("Expected log message 'Circuit breaker closed after successful recovery' not found")
	}
}

func TestCircuitBreaker_HalfOpenToOpen(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	log := zap.New(core)
	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		RecoveryTimeout:  10 * time.Millisecond,
		SuccessThreshold: 2, // Need 2 successes to close, but we'll fail after 1
	}

	decorator := NewCircuitBreakerToolProviderDecorator(nil, log, cfg)
	cbDecorator := decorator.(*CircuitBreakerToolProviderDecorator)

	mockInner := &mockToolProvider{
		executeToolFunc: func(_ context.Context, _ string, _ string) (string, error) {
			if cbDecorator.successCount == 0 {
				return "success", nil // First call in half-open succeeds
			}
			return "", errors.New("simulated error in half-open") // Second call fails
		},
	}
	cbDecorator.inner = mockInner

	// Force open state, ensure timeout has passed
	cbDecorator.state = CircuitOpen
	cbDecorator.failureCount = 1
	cbDecorator.lastFailureTime = time.Now().Add(-20 * time.Millisecond)

	// First attempt in half-open (success)
	_, err := decorator.ExecuteTool(context.Background(), "test_tool", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if cbDecorator.getState() != CircuitHalfOpen {
		t.Errorf("Expected state HALF_OPEN, got %s", cbDecorator.getState())
	}
	if cbDecorator.successCount != 1 {
		t.Errorf("Expected successCount 1, got %d", cbDecorator.successCount)
	}

	// Second attempt in half-open (failure) - should transition back to Open
	_, err = decorator.ExecuteTool(context.Background(), "test_tool", "")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if cbDecorator.getState() != CircuitOpen {
		t.Errorf("Expected state OPEN, got %s", cbDecorator.getState())
	}
	if cbDecorator.failureCount != 1 { // Failure count resets on state change, then increments
		t.Errorf("Expected failureCount 1, got %d", cbDecorator.failureCount)
	}

	// Verify log messages
	logs := recorded.All()
	if len(logs) < 2 {
		t.Fatalf("Expected at least 2 log messages, got %d", len(logs))
	}
	
	foundHalfOpen := false
	foundReopened := false
	for _, log := range logs {
		if strings.Contains(log.Message, "Circuit breaker transitioning to half-open state") {
			foundHalfOpen = true
		}
		if strings.Contains(log.Message, "Circuit breaker reopened during recovery attempt") {
			foundReopened = true
		}
	}
	
	if !foundHalfOpen {
		t.Error("Expected log message 'Circuit breaker transitioning to half-open state' not found")
	}
	if !foundReopened {
		t.Error("Expected log message 'Circuit breaker reopened during recovery attempt' not found")
	}
}

func TestCircuitBreaker_OpenStateRejection(t *testing.T) {
	core, recorded := observer.New(zapcore.WarnLevel) // Capture warnings
	log := zap.New(core)
	cfg := CircuitBreakerConfig{FailureThreshold: 1, RecoveryTimeout: 1 * time.Hour, SuccessThreshold: 1} // Long timeout
	mockInner := &mockToolProvider{
		executeToolFunc: func(_ context.Context, _ string, _ string) (string, error) {
			t.Fatal("Inner tool should not be called when circuit is open")
			return "", nil
		},
	}
	decorator := NewCircuitBreakerToolProviderDecorator(mockInner, log, cfg)
	cbDecorator := decorator.(*CircuitBreakerToolProviderDecorator)

	// Force open state, ensure recovery timeout has NOT passed
	cbDecorator.state = CircuitOpen
	cbDecorator.failureCount = 1
	cbDecorator.lastFailureTime = time.Now().Add(-10 * time.Minute) // Still within recovery timeout

	_, err := decorator.ExecuteTool(context.Background(), "test_tool", "")
	if err == nil || !errors.Is(err, internalErrors.ErrCircuitBreakerOpen) {
		t.Errorf("Expected ErrCircuitBreakerOpen, got %v", err)
	}
	if cbDecorator.getState() != CircuitOpen {
		t.Errorf("Expected state OPEN, got %s", cbDecorator.getState())
	}

	// Verify log message
	foundLog := false
	for _, entry := range recorded.All() {
		if strings.Contains(entry.Message, "Circuit breaker is open, rejecting tool execution") {
			foundLog = true
			break
		}
	}
	if !foundLog {
		t.Error("Expected log message 'Circuit breaker is open, rejecting tool execution' not found")
	}
}

func TestCircuitBreaker_ResetCircuitBreaker(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	log := zap.New(core)
	cfg := CircuitBreakerConfig{FailureThreshold: 1, RecoveryTimeout: 1 * time.Minute, SuccessThreshold: 1}
	mockInner := &mockToolProvider{}
	decorator := NewCircuitBreakerToolProviderDecorator(mockInner, log, cfg)
	cbDecorator := decorator.(*CircuitBreakerToolProviderDecorator)

	// Force open state
	cbDecorator.state = CircuitOpen
	cbDecorator.failureCount = 5
	cbDecorator.successCount = 5
	cbDecorator.lastFailureTime = time.Now().Add(-1 * time.Hour) // irrelevant for manual reset

	cbDecorator.ResetCircuitBreaker() // Call on concrete type

	if cbDecorator.getState() != CircuitClosed {
		t.Errorf("Expected state CLOSED after reset, got %s", cbDecorator.getState())
	}
	if cbDecorator.failureCount != 0 {
		t.Errorf("Expected failureCount 0 after reset, got %d", cbDecorator.failureCount)
	}
	if cbDecorator.successCount != 0 {
		t.Errorf("Expected successCount 0 after reset, got %d", cbDecorator.successCount)
	}

	foundLog := false
	for _, entry := range recorded.All() {
		if strings.Contains(entry.Message, "Circuit breaker manually reset") {
			foundLog = true
			break
		}
	}
	if !foundLog {
		t.Error("Expected log message 'Circuit breaker manually reset' not found")
	}
}

func TestCircuitBreaker_GetCircuitBreakerStats(t *testing.T) {
	log := zap.NewNop()
	cfg := CircuitBreakerConfig{FailureThreshold: 5, RecoveryTimeout: 30 * time.Second, SuccessThreshold: 3}
	mockInner := &mockToolProvider{}
	decorator := NewCircuitBreakerToolProviderDecorator(mockInner, log, cfg)
	cbDecorator := decorator.(*CircuitBreakerToolProviderDecorator)

	// Manually set some states for testing stats
	cbDecorator.state = CircuitHalfOpen
	cbDecorator.failureCount = 2
	cbDecorator.successCount = 1
	cbDecorator.lastFailureTime = time.Now().Add(-15 * time.Second)

	stats := cbDecorator.GetCircuitBreakerStats()

	expectedStats := map[string]any{
		"state":             "HALF_OPEN",
		"failure_count":     2,
		"success_count":     1,
		"failure_threshold": 5,
		"success_threshold": 3,
	}

	for key, expected := range expectedStats {
		actual, ok := stats[key]
		if !ok {
			t.Errorf("Expected stat %q not found", key)
			continue
		}
		if actual != expected {
			t.Errorf("Expected %q for %q, got %v", expected, key, actual)
		}
	}

	// Check time-related fields separately due to type assertion
	if _, ok := stats["recovery_timeout"].(string); !ok {
		t.Error("Expected recovery_timeout to be present and of type string")
	}
	if _, ok := stats["last_failure_time"].(time.Time); !ok {
		t.Error("Expected last_failure_time to be present and of type time.Time")
	}
}

// Pass-through methods are intentionally not tested in a way that requires
// direct instantiation of llms.ToolFunction, as it is not exported.
// We assume that the underlying inner.GetTools(), inner.GetToolsForAgent(),
// and inner.RegisterTool() calls are correct.
