package decorators

import (
	"context"

	"github.com/tmc/langchaingo/tools"

	"github.com/denkhaus/agentforge/internal/types"
)

// Mock ToolProvider for decorator testing.
type mockToolProvider struct {
	executeToolFunc func(ctx context.Context, name string, input string) (string, error)
	tools           []tools.Tool
	toolsForAgent   []tools.Tool
	registerToolErr error
}

func (m *mockToolProvider) ExecuteTool(ctx context.Context, name string, input string) (string, error) {
	if m.executeToolFunc != nil {
		return m.executeToolFunc(ctx, name, input)
	}
	return "mock_result", nil
}
func (m *mockToolProvider) GetTools() []tools.Tool                      { return m.tools }
func (m *mockToolProvider) GetToolsForAgent(_ types.Agent) ([]tools.Tool, error) { return m.toolsForAgent, nil }
func (m *mockToolProvider) RegisterTool(_ tools.Tool) error             { return m.registerToolErr }
func (m *mockToolProvider) HasTool(name string) bool                    { return true }
func (m *mockToolProvider) ValidateAgentRequirements(_ types.Agent) error { return nil }
func (m *mockToolProvider) GetToolNames() []string                     { 
	names := make([]string, len(m.tools))
	for i, tool := range m.tools {
		names[i] = tool.Name()
	}
	return names
}

type mockAgent struct{} // Minimal mock for types.Agent

func (m *mockAgent) GetName() string                    { return "mock_agent" }
func (m *mockAgent) GetDescription() string             { return "" }
func (m *mockAgent) GetSystemPrompt() string            { return "" }
func (m *mockAgent) GetRequiredTools() []string         { return []string{} }
func (m *mockAgent) GetLLMConfig() types.LLMConfig      { return nil }
func (m *mockAgent) HasRequiredTool(_ string) bool      { return false }
func (m *mockAgent) Clone(_ map[string]any) types.Agent { return &mockAgent{} }

// mockToolImpl implements tools.Tool interface for testing.
type mockToolImpl struct {
	name        string
	description string
	schema      map[string]any
}

func (m *mockToolImpl) Name() string        { return m.name }
func (m *mockToolImpl) Description() string { return m.description }
func (m *mockToolImpl) Call(ctx context.Context, input string) (string, error) {
	return "mock_result", nil
}
