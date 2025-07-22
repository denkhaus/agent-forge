package decorators

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// Mocks from test_mocks.go
// mockToolProvider
// mockAgent
// mockToolImpl

const (
	mockResult   = "mock_result"
	testToolName = "test_tool"
)

func TestNewToolMetrics(t *testing.T) {
	metrics := NewToolMetrics()
	if metrics == nil {
		t.Fatal("Expected non-nil ToolMetrics")
	}
	if metrics.ExecutionCounts == nil {
		t.Error("ExecutionCounts not initialized")
	}
	if metrics.ExecutionTimes == nil {
		t.Error("ExecutionTimes not initialized")
	}
	if metrics.ErrorCounts == nil {
		t.Error("ErrorCounts not initialized")
	}
	if metrics.TotalExecutions != 0 {
		t.Errorf("Expected TotalExecutions 0, got %d", metrics.TotalExecutions)
	}
	if metrics.TotalErrors != 0 {
		t.Errorf("Expected TotalErrors 0, got %d", metrics.TotalErrors)
	}
}

func TestNewMetricsToolProviderDecorator(t *testing.T) {
	mockInner := &mockToolProvider{}
	log := zap.NewNop()

	decorator := NewMetricsToolProviderDecorator(mockInner, log)
	metricsDecorator, ok := decorator.(*MetricsToolProviderDecorator)
	if !ok {
		t.Fatal("Expected decorator to be MetricsToolProviderDecorator")
	}

	if metricsDecorator.inner != mockInner {
		t.Error("Inner tool provider not set correctly")
	}
	if metricsDecorator.log != log {
		t.Error("Logger not set correctly")
	}
	if metricsDecorator.metrics == nil {
		t.Error("Metrics not initialized")
	}
}

func TestMetricsToolProviderDecorator_ExecuteTool(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	log := zap.New(core)

	callCount := 0
	mockInner := &mockToolProvider{
		executeToolFunc: func(_ context.Context, _ string, _ string) (string, error) {
			callCount++
			return mockResult, nil
		},
	}
	decorator := NewMetricsToolProviderDecorator(mockInner, log)
	metricsDecorator := decorator.(*MetricsToolProviderDecorator)

	toolName := testToolName
	input := "{}"

	// First execution (success)
	_, err := decorator.ExecuteTool(context.Background(), toolName, input)
	if err != nil {
		t.Fatalf("ExecuteTool failed: %v", err)
	}
	if callCount != 1 {
		t.Error("Expected inner tool to be called once")
	}

	// Verify metrics after first call
	metrics := metricsDecorator.GetMetrics()
	if metrics.TotalExecutions != 1 {
		t.Errorf("Expected TotalExecutions 1, got %d", metrics.TotalExecutions)
	}
	if metrics.ExecutionCounts[toolName] != 1 {
		t.Errorf("Expected ExecutionCounts %s 1, got %d", toolName, metrics.ExecutionCounts[toolName])
	}
	if len(metrics.ExecutionTimes[toolName]) != 1 {
		t.Errorf("Expected 1 ExecutionTime for %s, got %d", toolName, len(metrics.ExecutionTimes[toolName]))
	}
	if metrics.TotalErrors != 0 {
		t.Errorf("Expected TotalErrors 0, got %d", metrics.TotalErrors)
	}

	// Second execution (failure)
	mockInner.executeToolFunc = func(_ context.Context, _ string, _ string) (string, error) {
		callCount++
		return "", errors.New("simulated error")
	}
	_, err = decorator.ExecuteTool(context.Background(), toolName, input)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if callCount != 2 {
		t.Error("Expected inner tool to be called twice")
	}

	// Verify metrics after second call (including error)
	metrics = metricsDecorator.GetMetrics()
	if metrics.TotalExecutions != 2 {
		t.Errorf("Expected TotalExecutions 2, got %d", metrics.TotalExecutions)
	}
	if metrics.ExecutionCounts[toolName] != 2 {
		t.Errorf("Expected ExecutionCounts %s 2, got %d", toolName, metrics.ExecutionCounts[toolName])
	}
	if len(metrics.ExecutionTimes[toolName]) != 2 {
		t.Errorf("Expected 2 ExecutionTimes for %s, got %d", toolName, len(metrics.ExecutionTimes[toolName]))
	}
	if metrics.TotalErrors != 1 {
		t.Errorf("Expected TotalErrors 1, got %d", metrics.TotalErrors)
	}
	if metrics.ErrorCounts[toolName] != 1 {
		t.Errorf("Expected ErrorCounts %s 1, got %d", toolName, metrics.ErrorCounts[toolName])
	}

	// Test logMetrics is called periodically (e.g., every 10 executions)
	// Execute 8 more times (total 10)
	for i := 0; i < 8; i++ {
		_, _ = decorator.ExecuteTool(context.Background(), toolName, input)
	}

	foundLogMetrics := false
	for _, entry := range recorded.All() {
		if strings.Contains(entry.Message, "Tool execution metrics") {
			foundLogMetrics = true
			break
		}
	}
	if !foundLogMetrics {
		t.Error("Expected 'Tool execution metrics' log message not found")
	}
}

func TestMetricsToolProviderDecorator_GetMetrics(t *testing.T) {
	log := zap.NewNop()
	mockInner := &mockToolProvider{}
	decorator := NewMetricsToolProviderDecorator(mockInner, log)
	metricsDecorator := decorator.(*MetricsToolProviderDecorator)

	// Manually add some data to metrics
	metricsDecorator.metrics.mutex.Lock()
	metricsDecorator.metrics.ExecutionCounts["toolA"] = 10
	metricsDecorator.metrics.ExecutionCounts["toolB"] = 5
	metricsDecorator.metrics.TotalExecutions = 15
	metricsDecorator.metrics.ErrorCounts["toolA"] = 2
	metricsDecorator.metrics.TotalErrors = 2
	metricsDecorator.metrics.ExecutionTimes["toolA"] = []time.Duration{10 * time.Millisecond, 20 * time.Millisecond}
	metricsDecorator.metrics.mutex.Unlock()

	snapshot := metricsDecorator.GetMetrics()

	if snapshot.TotalExecutions != 15 {
		t.Errorf("Expected TotalExecutions 15, got %d", snapshot.TotalExecutions)
	}
	if snapshot.TotalErrors != 2 {
		t.Errorf("Expected TotalErrors 2, got %d", snapshot.TotalErrors)
	}
	if snapshot.ExecutionCounts["toolA"] != 10 {
		t.Errorf("Expected ExecutionCounts toolA 10, got %d", snapshot.ExecutionCounts["toolA"])
	}
	if snapshot.ExecutionCounts["toolB"] != 5 {
		t.Errorf("Expected ExecutionCounts toolB 5, got %d", snapshot.ExecutionCounts["toolB"])
	}
	if snapshot.ErrorCounts["toolA"] != 2 {
		t.Errorf("Expected ErrorCounts toolA 2, got %d", snapshot.ErrorCounts["toolA"])
	}
	if len(snapshot.ExecutionTimes["toolA"]) != 2 {
		t.Errorf("Expected 2 ExecutionTimes for toolA, got %d", len(snapshot.ExecutionTimes["toolA"]))
	}
}

// Pass-through methods are intentionally not tested in a way that requires
// direct instantiation of llms.ToolFunction, as it is not exported.
// We assume that the underlying inner.GetTools(), inner.GetToolsForAgent(),
// and inner.RegisterTool() calls are correct.
