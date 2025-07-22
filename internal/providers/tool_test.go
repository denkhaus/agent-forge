package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/samber/do"
	"github.com/tmc/langchaingo/tools"
	"go.uber.org/zap"

	"github.com/denkhaus/agentforge/internal/types"
)

func TestNewToolProvider(t *testing.T) {
	// Setup DI container for testing
	injector := do.New()
	defer func() {
		if err := injector.Shutdown(); err != nil {
			t.Logf("Failed to shutdown injector: %v", err)
		}
	}()

	// Register test dependencies
	do.ProvideValue(injector, zap.NewNop())

	do.ProvideValue(injector, []tools.Tool{
		&testTool{name: "test_tool", description: "Test tool for unit tests"},
	})

	provider, err := NewToolProvider(injector)
	if err != nil {
		t.Fatalf("Failed to create tool provider: %v", err)
	}

	if provider == nil {
		t.Fatal("Expected non-nil tool provider")
	}

	// Test that provider implements the simplified interface
	var _ = provider
}

func TestToolProviderGetTools(t *testing.T) {
	provider := createTestToolProvider(t)

	tools := provider.GetTools()
	if len(tools) == 0 {
		t.Error("Expected at least one tool")
	}

	// Verify tool structure
	found := false
	for _, tool := range tools {
		if tool.Name() == "test_tool" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find 'test_tool' in tools list")
	}
}

func TestToolProviderGetToolsForAgent(t *testing.T) {
	provider := createTestToolProvider(t)

	// Create test agent with specific tool requirements
	agent := &testAgent{
		name:          "test_agent",
		requiredTools: []string{"test_tool"},
	}

	tools, err := provider.GetToolsForAgent(agent)
	if err != nil {
		t.Fatalf("Failed to get tools for agent: %v", err)
	}

	// Should return tools that match agent requirements
	// The agent requires "test_tool" which exists in the provider
	if len(tools) == 0 {
		t.Error("Expected at least one tool for agent with matching requirements")
	}

	// Test agent without required tools
	agentNoTools := &testAgent{
		name:          "limited_agent",
		requiredTools: []string{}, // No required tools
	}

	limitedTools, err := provider.GetToolsForAgent(agentNoTools)
	if err != nil {
		t.Fatalf("Failed to get tools for agent: %v", err)
	}

	// Should return empty since agent has no required tools
	if len(limitedTools) != 0 {
		t.Errorf("Expected no tools for agent with no requirements, got %d", len(limitedTools))
	}
}

func TestToolProviderExecuteTool(t *testing.T) {
	provider := createTestToolProvider(t)

	ctx := context.Background()
	inputArgs := map[string]any{
		"input": "test_value",
	}
	inputJSON, err := json.Marshal(inputArgs)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}
	result, err := provider.ExecuteTool(ctx, "test_tool", string(inputJSON))
	if err != nil {
		t.Fatalf("Failed to execute tool: %v", err)
	}

	expected := "test result: test_value"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestToolProviderExecuteNonexistentTool(t *testing.T) {
	provider := createTestToolProvider(t)

	ctx := context.Background()
	inputJSON, err := json.Marshal(map[string]any{})
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}
	_, err = provider.ExecuteTool(ctx, "nonexistent_tool", string(inputJSON))
	if err == nil {
		t.Error("Expected error when executing nonexistent tool")
	}
}

func TestToolProviderRegisterTool(t *testing.T) {
	provider := createTestToolProvider(t)

	newTool := &testTool{
		name:        "new_tool",
		description: "A new test tool",
	}

	err := provider.RegisterTool(newTool)
	if err != nil {
		t.Fatalf("Failed to register tool: %v", err)
	}

	// Verify tool was registered by trying to execute it
	ctx := context.Background()
	inputArgs := map[string]any{"input": "new_test"}
	inputJSON, err := json.Marshal(inputArgs)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}
	result, err := provider.ExecuteTool(ctx, "new_tool", string(inputJSON))
	if err != nil {
		t.Fatalf("Failed to execute newly registered tool: %v", err)
	}

	expected := "test result: new_test"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestToolProviderRegisterDuplicateTool(t *testing.T) {
	provider := createTestToolProvider(t)

	duplicateTool := &testTool{
		name:        "test_tool", // Same as existing tool
		description: "Duplicate tool",
	}

	err := provider.RegisterTool(duplicateTool)
	if err == nil {
		t.Error("Expected error when registering duplicate tool")
	}
}

func TestToolProviderCaching(t *testing.T) {
	provider := createTestToolProvider(t)

	// Create test agent
	agent := &testAgent{
		name:          "cache_test_agent",
		requiredTools: []string{"test_tool"},
	}

	// First call - should populate cache
	tools1, err := provider.GetToolsForAgent(agent)
	if err != nil {
		t.Fatalf("Failed to get tools for agent: %v", err)
	}

	// Second call - should use cache
	tools2, err := provider.GetToolsForAgent(agent)
	if err != nil {
		t.Fatalf("Failed to get tools for agent: %v", err)
	}

	// Results should be identical
	if len(tools1) != len(tools2) {
		t.Errorf("Expected same number of tools from cache, got %d vs %d", len(tools1), len(tools2))
	}
}

// Helper functions and test types

type testTool struct {
	name        string
	description string
}

func (t *testTool) Name() string        { return t.name }
func (t *testTool) Description() string { return t.description }
func (t *testTool) Schema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"input": map[string]any{
				"type":        "string",
				"description": "Test input parameter",
			},
		},
		"required": []string{"input"},
	}
}

func (t *testTool) Call(ctx context.Context, input string) (string, error) {
	// For testing, we'll parse the input string back to a map.
	// In real tools, you'd directly use the input string or define a specific input struct.
	var args map[string]any
	if err := json.Unmarshal([]byte(input), &args); err != nil {
		return "", fmt.Errorf("invalid input format for testTool: %w", err)
	}

	inputValue, ok := args["input"].(string)
	if !ok {
		inputValue = "default"
	}
	return "test result: " + inputValue, nil
}

type testAgent struct {
	name          string
	requiredTools []string
}

func (a *testAgent) GetName() string                            { return a.name }
func (a *testAgent) GetDescription() string                     { return "Test agent" }
func (a *testAgent) GetSystemPrompt() string                    { return "Test system prompt" }
func (a *testAgent) GetRequiredTools() []string                 { return a.requiredTools }
func (a *testAgent) GetLLMConfig() types.LLMConfig              { return nil }
func (a *testAgent) Clone(overrides map[string]any) types.Agent { return a }
func (a *testAgent) HasRequiredTool(toolName string) bool {
	for _, tool := range a.requiredTools {
		if tool == toolName {
			return true
		}
	}
	return false
}

func createTestToolProvider(t *testing.T) types.ToolProvider {
	injector := do.New()
	t.Cleanup(func() {
		if err := injector.Shutdown(); err != nil {
			t.Logf("Failed to shutdown injector: %v", err)
		}
	})

	// Register test dependencies
	do.ProvideValue(injector, zap.NewNop())

	do.ProvideValue(injector, []tools.Tool{
		&testTool{
			name:        "test_tool",
			description: "Test tool for unit tests",
		},
	})

	provider, err := NewToolProvider(injector)
	if err != nil {
		t.Fatalf("Failed to create test tool provider: %v", err)
	}

	return provider
}
