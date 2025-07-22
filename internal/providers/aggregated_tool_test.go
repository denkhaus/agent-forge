package providers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/tools"
	"go.uber.org/zap/zaptest"

	"github.com/denkhaus/agentforge/internal/types"
)

// MockToolProvider implements types.ToolProvider for testing
type MockToolProvider struct {
	mock.Mock
}

func (m *MockToolProvider) GetTools() []tools.Tool {
	args := m.Called()
	return args.Get(0).([]tools.Tool)
}

func (m *MockToolProvider) GetToolsForAgent(agent types.Agent) ([]tools.Tool, error) {
	args := m.Called(agent)
	return args.Get(0).([]tools.Tool), args.Error(1)
}

func (m *MockToolProvider) ExecuteTool(ctx context.Context, name string, input string) (string, error) {
	args := m.Called(ctx, name, input)
	return args.String(0), args.Error(1)
}

func (m *MockToolProvider) RegisterTool(tool tools.Tool) error {
	args := m.Called(tool)
	return args.Error(0)
}

func (m *MockToolProvider) HasTool(name string) bool {
	args := m.Called(name)
	return args.Bool(0)
}

func (m *MockToolProvider) ValidateAgentRequirements(agent types.Agent) error {
	args := m.Called(agent)
	return args.Error(0)
}

func (m *MockToolProvider) GetToolNames() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

// MockTool implements tools.Tool for testing
type MockTool struct {
	name        string
	description string
}

func (m *MockTool) Name() string {
	return m.name
}

func (m *MockTool) Description() string {
	return m.description
}

func (m *MockTool) Call(ctx context.Context, input string) (string, error) {
	return "mock result", nil
}

func TestNewAggregatedToolProvider(t *testing.T) {
	log := zaptest.NewLogger(t)
	
	provider1 := &MockToolProvider{}
	provider2 := &MockToolProvider{}
	
	// Setup mock expectations
	tool1 := &MockTool{name: "tool1", description: "Tool 1"}
	tool2 := &MockTool{name: "tool2", description: "Tool 2"}
	
	provider1.On("GetTools").Return([]tools.Tool{tool1})
	provider2.On("GetTools").Return([]tools.Tool{tool2})
	
	aggregated := NewAggregatedToolProvider(log, provider1, provider2)
	assert.NotNil(t, aggregated)
	
	// Verify all tools are available
	allTools := aggregated.GetTools()
	assert.Len(t, allTools, 2)
	
	toolNames := aggregated.GetToolNames()
	assert.Len(t, toolNames, 2)
	assert.Contains(t, toolNames, "tool1")
	assert.Contains(t, toolNames, "tool2")
}

func TestAggregatedToolProvider_HasTool(t *testing.T) {
	log := zaptest.NewLogger(t)
	
	provider1 := &MockToolProvider{}
	tool1 := &MockTool{name: "tool1", description: "Tool 1"}
	
	provider1.On("GetTools").Return([]tools.Tool{tool1})
	
	aggregated := NewAggregatedToolProvider(log, provider1)
	
	assert.True(t, aggregated.HasTool("tool1"))
	assert.False(t, aggregated.HasTool("nonexistent"))
}

func TestAggregatedToolProvider_ExecuteTool(t *testing.T) {
	log := zaptest.NewLogger(t)
	ctx := context.Background()
	
	provider1 := &MockToolProvider{}
	tool1 := &MockTool{name: "tool1", description: "Tool 1"}
	
	provider1.On("GetTools").Return([]tools.Tool{tool1})
	provider1.On("ExecuteTool", ctx, "tool1", "input").Return("result", nil)
	
	aggregated := NewAggregatedToolProvider(log, provider1)
	
	result, err := aggregated.ExecuteTool(ctx, "tool1", "input")
	require.NoError(t, err)
	assert.Equal(t, "result", result)
	
	// Test nonexistent tool
	_, err = aggregated.ExecuteTool(ctx, "nonexistent", "input")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in any provider")
}

func TestAggregatedToolProvider_RegisterTool(t *testing.T) {
	log := zaptest.NewLogger(t)
	
	provider1 := &MockToolProvider{}
	provider2 := &MockToolProvider{}
	
	provider1.On("GetTools").Return([]tools.Tool{})
	provider2.On("GetTools").Return([]tools.Tool{})
	
	newTool := &MockTool{name: "newTool", description: "New Tool"}
	
	// First provider fails, second succeeds
	provider1.On("RegisterTool", newTool).Return(assert.AnError)
	provider2.On("RegisterTool", newTool).Return(nil)
	
	aggregated := NewAggregatedToolProvider(log, provider1, provider2)
	
	err := aggregated.RegisterTool(newTool)
	assert.NoError(t, err)
	
	provider1.AssertExpectations(t)
	provider2.AssertExpectations(t)
}

func TestAggregatedToolProvider_ToolNameConflict(t *testing.T) {
	log := zaptest.NewLogger(t)
	
	provider1 := &MockToolProvider{}
	provider2 := &MockToolProvider{}
	
	// Both providers have a tool with the same name
	tool1 := &MockTool{name: "conflictTool", description: "Tool from provider 1"}
	tool2 := &MockTool{name: "conflictTool", description: "Tool from provider 2"}
	
	provider1.On("GetTools").Return([]tools.Tool{tool1})
	provider2.On("GetTools").Return([]tools.Tool{tool2})
	
	aggregated := NewAggregatedToolProvider(log, provider1, provider2)
	
	// Should only have one tool (first provider wins)
	allTools := aggregated.GetTools()
	assert.Len(t, allTools, 1)
	assert.Equal(t, "Tool from provider 1", allTools[0].Description())
	
	// Should route to first provider
	provider1.On("ExecuteTool", mock.Anything, "conflictTool", "input").Return("result1", nil)
	
	result, err := aggregated.ExecuteTool(context.Background(), "conflictTool", "input")
	require.NoError(t, err)
	assert.Equal(t, "result1", result)
	
	provider1.AssertExpectations(t)
}

func TestAggregatedToolProvider_AddRemoveProvider(t *testing.T) {
	log := zaptest.NewLogger(t)
	
	provider1 := &MockToolProvider{}
	provider2 := &MockToolProvider{}
	
	tool1 := &MockTool{name: "tool1", description: "Tool 1"}
	tool2 := &MockTool{name: "tool2", description: "Tool 2"}
	
	provider1.On("GetTools").Return([]tools.Tool{tool1})
	provider2.On("GetTools").Return([]tools.Tool{tool2})
	
	// Start with one provider
	aggregated := NewAggregatedToolProvider(log, provider1).(*aggregatedToolProvider)
	
	assert.True(t, aggregated.HasTool("tool1"))
	assert.False(t, aggregated.HasTool("tool2"))
	
	// Add second provider
	aggregated.AddProvider(provider2)
	
	// Force cache rebuild by calling GetTools
	_ = aggregated.GetTools()
	
	assert.True(t, aggregated.HasTool("tool1"))
	assert.True(t, aggregated.HasTool("tool2"))
	
	// Remove second provider
	aggregated.RemoveProvider(provider2)
	
	// Force cache rebuild by calling GetTools
	_ = aggregated.GetTools()
	
	assert.True(t, aggregated.HasTool("tool1"))
	assert.False(t, aggregated.HasTool("tool2"))
}

func TestAggregatedToolProvider_GetToolsForAgent(t *testing.T) {
	log := zaptest.NewLogger(t)
	
	provider1 := &MockToolProvider{}
	provider2 := &MockToolProvider{}
	
	tool1 := &MockTool{name: "tool1", description: "Tool 1"}
	tool2 := &MockTool{name: "tool2", description: "Tool 2"}
	
	provider1.On("GetTools").Return([]tools.Tool{tool1})
	provider2.On("GetTools").Return([]tools.Tool{tool2})
	
	// Mock agent requiring both tools
	agent := &MockAgent{}
	agent.On("GetRequiredTools").Return([]string{"tool1", "tool2"})
	agent.On("GetName").Return("test_agent")
	
	aggregated := NewAggregatedToolProvider(log, provider1, provider2)
	
	// Should validate successfully
	err := aggregated.ValidateAgentRequirements(agent)
	assert.NoError(t, err)
	
	// Should return both tools
	agentTools, err := aggregated.GetToolsForAgent(agent)
	require.NoError(t, err)
	assert.Len(t, agentTools, 2)
	
	// Verify tool names
	toolNames := make([]string, len(agentTools))
	for i, tool := range agentTools {
		toolNames[i] = tool.Name()
	}
	assert.Contains(t, toolNames, "tool1")
	assert.Contains(t, toolNames, "tool2")
}

func TestAggregatedToolProvider_GetToolsForAgent_MissingTool(t *testing.T) {
	log := zaptest.NewLogger(t)
	
	provider1 := &MockToolProvider{}
	tool1 := &MockTool{name: "tool1", description: "Tool 1"}
	
	provider1.On("GetTools").Return([]tools.Tool{tool1})
	
	// Mock agent requiring missing tool
	agent := &MockAgent{}
	agent.On("GetRequiredTools").Return([]string{"tool1", "missing_tool"})
	agent.On("GetName").Return("test_agent")
	
	aggregated := NewAggregatedToolProvider(log, provider1)
	
	// Should fail validation
	err := aggregated.ValidateAgentRequirements(agent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing tools")
	assert.Contains(t, err.Error(), "missing_tool")
	
	// Should return error for GetToolsForAgent
	agentTools, err := aggregated.GetToolsForAgent(agent)
	assert.Error(t, err)
	assert.Nil(t, agentTools)
}

func TestAggregatedToolProvider_CacheInvalidation(t *testing.T) {
	log := zaptest.NewLogger(t)
	
	provider1 := &MockToolProvider{}
	tool1 := &MockTool{name: "tool1", description: "Tool 1"}
	
	provider1.On("GetTools").Return([]tools.Tool{tool1})
	
	aggregated := NewAggregatedToolProvider(log, provider1).(*aggregatedToolProvider)
	
	// Initial state
	assert.True(t, aggregated.HasTool("tool1"))
	assert.Len(t, aggregated.GetTools(), 1)
	
	// Add provider should invalidate cache
	provider2 := &MockToolProvider{}
	tool2 := &MockTool{name: "tool2", description: "Tool 2"}
	provider2.On("GetTools").Return([]tools.Tool{tool2})
	
	aggregated.AddProvider(provider2)
	
	// Force cache rebuild by calling a method that triggers it
	_ = aggregated.GetTools()
	
	// Cache should be rebuilt automatically
	assert.True(t, aggregated.HasTool("tool1"))
	assert.True(t, aggregated.HasTool("tool2"))
	assert.Len(t, aggregated.GetTools(), 2)
}

func TestAggregatedToolProvider_EmptyProviders(t *testing.T) {
	log := zaptest.NewLogger(t)
	
	// Create aggregated provider with no providers
	aggregated := NewAggregatedToolProvider(log)
	
	// Should handle empty state gracefully
	assert.Empty(t, aggregated.GetTools())
	assert.Empty(t, aggregated.GetToolNames())
	assert.False(t, aggregated.HasTool("any_tool"))
	
	// Should return error for tool execution
	_, err := aggregated.ExecuteTool(context.Background(), "any_tool", "input")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in any provider")
}

func TestAggregatedToolProvider_RegisterTool_AllProvidersFail(t *testing.T) {
	log := zaptest.NewLogger(t)
	
	provider1 := &MockToolProvider{}
	provider2 := &MockToolProvider{}
	
	provider1.On("GetTools").Return([]tools.Tool{})
	provider2.On("GetTools").Return([]tools.Tool{})
	
	newTool := &MockTool{name: "newTool", description: "New Tool"}
	
	// Both providers fail to register
	provider1.On("RegisterTool", newTool).Return(assert.AnError)
	provider2.On("RegisterTool", newTool).Return(assert.AnError)
	
	aggregated := NewAggregatedToolProvider(log, provider1, provider2)
	
	err := aggregated.RegisterTool(newTool)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to register tool")
	
	provider1.AssertExpectations(t)
	provider2.AssertExpectations(t)
}

func TestAggregatedToolProvider_ConcurrentAccess(t *testing.T) {
	log := zaptest.NewLogger(t)
	
	provider1 := &MockToolProvider{}
	tool1 := &MockTool{name: "tool1", description: "Tool 1"}
	
	provider1.On("GetTools").Return([]tools.Tool{tool1})
	provider1.On("ExecuteTool", mock.Anything, "tool1", mock.Anything).Return("result", nil)
	
	aggregated := NewAggregatedToolProvider(log, provider1)
	
	// Test concurrent access
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			
			// Concurrent reads
			_ = aggregated.GetTools()
			_ = aggregated.GetToolNames()
			_ = aggregated.HasTool("tool1")
			
			// Concurrent execution
			_, _ = aggregated.ExecuteTool(context.Background(), "tool1", "input")
		}()
	}
	
	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Should still work correctly
	assert.True(t, aggregated.HasTool("tool1"))
}