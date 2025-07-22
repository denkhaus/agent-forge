package decorators

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/tools"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewCachingToolProviderDecorator(t *testing.T) {
	mockInner := &mockToolProvider{}
	log := zap.NewNop()

	decorator := NewCachingToolProviderDecorator(mockInner, log, 10*time.Minute)
	cachingDecorator, ok := decorator.(*CachingToolProviderDecorator)
	if !ok {
		t.Fatal("Expected decorator to be CachingToolProviderDecorator")
	}

	if cachingDecorator.inner != mockInner {
		t.Error("Inner tool provider not set correctly")
	}
	if cachingDecorator.log != log {
		t.Error("Logger not set correctly")
	}
	if cachingDecorator.defaultTTL != 10*time.Minute {
		t.Errorf("Default TTL not set correctly, got %v", cachingDecorator.defaultTTL)
	}
	if cachingDecorator.cache == nil {
		t.Error("Cache map not initialized")
	}
}

func TestCachingToolProviderDecorator_ExecuteTool_Caching(t *testing.T) {
	core, _ := observer.New(zapcore.InfoLevel)
	log := zap.New(core)

	callCount := 0
	mockInner := &mockToolProvider{
		executeToolFunc: func(_ context.Context, name string, input string) (string, error) {
			callCount++
			return fmt.Sprintf("result_for_%s_%s", name, input), nil
		},
	}

	decorator := NewCachingToolProviderDecorator(mockInner, log, 1*time.Second)
	ctx := context.Background()
	toolName := "test_tool"
	input := `{"param": "value1"}`

	result, err := decorator.ExecuteTool(ctx, toolName, input)
	if err != nil {
		t.Fatalf("ExecuteTool failed: %v", err)
	}
	expectedResult := `result_for_test_tool_{"param": "value1"}`
	if result != expectedResult {
		t.Errorf("Expected result %s, got %s", expectedResult, result)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call to inner provider, got %d", callCount)
	}

	// Second call should use cache
	result2, err2 := decorator.ExecuteTool(ctx, toolName, input)
	if err2 != nil {
		t.Fatalf("Second ExecuteTool failed: %v", err2)
	}
	if result2 != expectedResult {
		t.Errorf("Expected cached result %s, got %s", expectedResult, result2)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call to inner provider (cached), got %d", callCount)
	}
}

func TestCachingToolProviderDecorator_ExecuteTool_CacheExpiry(t *testing.T) {
	core, _ := observer.New(zapcore.InfoLevel)
	log := zap.New(core)

	callCount := 0
	mockInner := &mockToolProvider{
		executeToolFunc: func(_ context.Context, name string, input string) (string, error) {
			callCount++
			return fmt.Sprintf("result_for_%s_%s", name, input), nil
		},
	}

	decorator := NewCachingToolProviderDecorator(mockInner, log, 50*time.Millisecond)
	ctx := context.Background()
	toolName := "test_tool"
	input := `{"param": "value1"}`

	// First call
	_, err := decorator.ExecuteTool(ctx, toolName, input)
	if err != nil {
		t.Fatalf("ExecuteTool failed: %v", err)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call to inner provider, got %d", callCount)
	}

	// Wait for cache to expire
	time.Sleep(100 * time.Millisecond)

	// Second call should not use cache
	_, err = decorator.ExecuteTool(ctx, toolName, input)
	if err != nil {
		t.Fatalf("Second ExecuteTool failed: %v", err)
	}
	if callCount != 2 {
		t.Errorf("Expected 2 calls to inner provider (cache expired), got %d", callCount)
	}
}

func TestCachingToolProviderDecorator_ExecuteTool_DifferentInputs(t *testing.T) {
	core, _ := observer.New(zapcore.InfoLevel)
	log := zap.New(core)

	callCount := 0
	mockInner := &mockToolProvider{
		executeToolFunc: func(_ context.Context, name string, input string) (string, error) {
			callCount++
			return fmt.Sprintf("result_for_%s_%s", name, input), nil
		},
	}

	decorator := NewCachingToolProviderDecorator(mockInner, log, 1*time.Minute)
	ctx := context.Background()
	toolName := "test_tool"

	// First call
	input1 := `{"param": "value1"}`
	result1, err := decorator.ExecuteTool(ctx, toolName, input1)
	if err != nil {
		t.Fatalf("ExecuteTool failed: %v", err)
	}
	expectedResult1 := `result_for_test_tool_{"param": "value1"}`
	if result1 != expectedResult1 {
		t.Errorf("Expected result %s, got %s", expectedResult1, result1)
	}

	// Second call with different input
	input2 := `{"param": "value2"}`
	result2, err := decorator.ExecuteTool(ctx, toolName, input2)
	if err != nil {
		t.Fatalf("Second ExecuteTool failed: %v", err)
	}
	expectedResult2 := `result_for_test_tool_{"param": "value2"}`
	if result2 != expectedResult2 {
		t.Errorf("Expected result %s, got %s", expectedResult2, result2)
	}

	if callCount != 2 {
		t.Errorf("Expected 2 calls to inner provider (different inputs), got %d", callCount)
	}
}

func TestCachingToolProviderDecorator_ClearCache(t *testing.T) {
	mockInner := &mockToolProvider{
		executeToolFunc: func(_ context.Context, name string, input string) (string, error) {
			return "test_result", nil
		},
	}

	core, observed := observer.New(zapcore.InfoLevel)
	log := zap.New(core)

	decorator := NewCachingToolProviderDecorator(mockInner, log, 1*time.Minute)
	cachingDecorator := decorator.(*CachingToolProviderDecorator)

	ctx := context.Background()
	// Add some entries to cache
	_, _ = decorator.ExecuteTool(ctx, "tool1", `{"param": "value1"}`)
	_, _ = decorator.ExecuteTool(ctx, "tool2", `{"param": "value2"}`)

	// Verify cache has entries
	stats := cachingDecorator.GetCacheStats()
	assert.Equal(t, 2, stats["total_entries"])

	// Clear cache
	cachingDecorator.ClearCache()

	// Verify cache is empty
	stats = cachingDecorator.GetCacheStats()
	assert.Equal(t, 0, stats["total_entries"])

	// Verify log message
	logs := observed.All()
	found := false
	for _, log := range logs {
		if log.Message == "Cache cleared" {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected 'Cache cleared' log message")
}

func TestCachingToolProviderDecorator_GetCacheStats(t *testing.T) {
	mockInner := &mockToolProvider{
		executeToolFunc: func(_ context.Context, name string, input string) (string, error) {
			return "test_result", nil
		},
	}

	log := zap.NewNop()
	decorator := NewCachingToolProviderDecorator(mockInner, log, 1*time.Minute)
	cachingDecorator := decorator.(*CachingToolProviderDecorator)

	ctx := context.Background()
	// Add some entries to cache
	_, _ = decorator.ExecuteTool(ctx, "tool1", `{"param": "value1"}`)
	_, _ = decorator.ExecuteTool(ctx, "tool2", `{"param": "value2"}`)

	stats := cachingDecorator.GetCacheStats()
	assert.Equal(t, 2, stats["total_entries"])
	assert.Equal(t, 0, stats["expired_entries"])
	assert.Equal(t, 2, stats["active_entries"])
	assert.Equal(t, "1m0s", stats["default_ttl"])
}

func TestCachingToolProviderDecorator_PassthroughMethods(t *testing.T) {
	mockInner := &mockToolProvider{
		tools: []tools.Tool{},
	}

	log := zap.NewNop()
	decorator := NewCachingToolProviderDecorator(mockInner, log, 1*time.Minute)

	// Test passthrough methods
	tools := decorator.GetTools()
	assert.Equal(t, mockInner.tools, tools)

	agent := &mockAgent{}
	agentTools, err := decorator.GetToolsForAgent(agent)
	assert.NoError(t, err)
	assert.Equal(t, mockInner.toolsForAgent, agentTools)

	assert.True(t, decorator.HasTool("any_tool"))
	assert.NoError(t, decorator.ValidateAgentRequirements(agent))

	toolNames := decorator.GetToolNames()
	assert.Equal(t, []string{}, toolNames)
}