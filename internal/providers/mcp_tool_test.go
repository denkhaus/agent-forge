package providers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/denkhaus/agentforge/internal/config"
	"github.com/denkhaus/agentforge/internal/types"
)

func TestNewMCPToolProvider_Disabled(t *testing.T) {
	log := zaptest.NewLogger(t)
	cfg := &config.MCPConfig{
		Enabled: false,
	}

	provider, err := NewMCPToolProvider(log, cfg)
	require.NoError(t, err)
	assert.NotNil(t, provider)

	// Should have no tools when disabled
	tools := provider.GetTools()
	assert.Empty(t, tools)

	toolNames := provider.GetToolNames()
	assert.Empty(t, toolNames)
}

func TestNewMCPToolProvider_InvalidConfig(t *testing.T) {
	log := zaptest.NewLogger(t)
	cfg := &config.MCPConfig{
		Enabled:    true,
		ConfigPath: "nonexistent-config.json",
	}

	provider, err := NewMCPToolProvider(log, cfg)
	// Should return error for invalid config path
	assert.Error(t, err)
	assert.Nil(t, provider)
}

func TestMCPToolProvider_HasTool(t *testing.T) {
	log := zaptest.NewLogger(t)
	cfg := &config.MCPConfig{
		Enabled: false, // Disabled for this test
	}

	provider, err := NewMCPToolProvider(log, cfg)
	require.NoError(t, err)

	// Should return false for any tool when disabled
	assert.False(t, provider.HasTool("any_tool"))
}

func TestMCPToolProvider_RegisterTool(t *testing.T) {
	log := zaptest.NewLogger(t)
	cfg := &config.MCPConfig{
		Enabled: false,
	}

	provider, err := NewMCPToolProvider(log, cfg)
	require.NoError(t, err)

	// Should return error when trying to register tools
	err = provider.RegisterTool(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not supported for MCP provider")
}

// MockAgent implements types.Agent for testing
type MockAgent struct {
	mock.Mock
}

func (m *MockAgent) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockAgent) GetDescription() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockAgent) GetSystemPrompt() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockAgent) GetRequiredTools() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockAgent) GetLLMConfig() types.LLMConfig {
	args := m.Called()
	return args.Get(0).(types.LLMConfig)
}

func (m *MockAgent) HasRequiredTool(toolName string) bool {
	args := m.Called(toolName)
	return args.Bool(0)
}

func (m *MockAgent) Clone(overrides map[string]any) types.Agent {
	args := m.Called(overrides)
	return args.Get(0).(types.Agent)
}

func TestMCPToolProvider_GetToolsForAgent_Disabled(t *testing.T) {
	log := zaptest.NewLogger(t)
	cfg := &config.MCPConfig{
		Enabled: false,
	}

	provider, err := NewMCPToolProvider(log, cfg)
	require.NoError(t, err)

	agent := &MockAgent{}
	agent.On("GetRequiredTools").Return([]string{"some_tool"})
	agent.On("GetName").Return("test_agent")

	// Should return error for missing tools when disabled
	tools, err := provider.GetToolsForAgent(agent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing MCP tools")
	assert.Empty(t, tools)
}

func TestMCPToolProvider_ValidateAgentRequirements_Disabled(t *testing.T) {
	log := zaptest.NewLogger(t)
	cfg := &config.MCPConfig{
		Enabled: false,
	}

	provider, err := NewMCPToolProvider(log, cfg)
	require.NoError(t, err)

	agent := &MockAgent{}
	agent.On("GetRequiredTools").Return([]string{"some_tool"})
	agent.On("GetName").Return("test_agent")

	// Should return error for missing tools when disabled
	err = provider.ValidateAgentRequirements(agent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing MCP tools")
}

func TestMCPToolProvider_ExecuteTool_Disabled(t *testing.T) {
	log := zaptest.NewLogger(t)
	cfg := &config.MCPConfig{
		Enabled: false,
	}

	provider, err := NewMCPToolProvider(log, cfg)
	require.NoError(t, err)

	// Should return error for nonexistent tool when disabled
	result, err := provider.ExecuteTool(context.Background(), "nonexistent", "input")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Empty(t, result)
}

func TestMCPToolProvider_GetToolNames_Disabled(t *testing.T) {
	log := zaptest.NewLogger(t)
	cfg := &config.MCPConfig{
		Enabled: false,
	}

	provider, err := NewMCPToolProvider(log, cfg)
	require.NoError(t, err)

	// Should return empty list when disabled
	names := provider.GetToolNames()
	assert.Empty(t, names)
}

func TestMCPToolWrapper(t *testing.T) {
	log := zaptest.NewLogger(t)
	
	// Create a mock tool
	mockTool := &MockTool{
		name:        "original_name",
		description: "Test tool description",
	}

	// Create wrapper with prefix
	wrapper := &mcpToolWrapper{
		mcpTool: mockTool,
		name:    "mcp_original_name",
		log:     log,
	}

	// Test Name method
	assert.Equal(t, "mcp_original_name", wrapper.Name())

	// Test Description method
	assert.Equal(t, "Test tool description", wrapper.Description())

	// Test Call method
	result, err := wrapper.Call(context.Background(), "test input")
	require.NoError(t, err)
	assert.Equal(t, "mock result", result)
}

func TestMCPConfig_GetMethods(t *testing.T) {
	cfg := &config.Config{
		MCPEnabled:       true,
		MCPConfigPath:    "test-config.json",
		MCPServerTimeout: 60,
		MCPToolPrefix:    "test_",
		MCPHotReload:     true,
	}

	mcpConfig := cfg.GetMCPConfig()

	assert.True(t, mcpConfig.Enabled)
	assert.Equal(t, "test-config.json", mcpConfig.ConfigPath)
	assert.Equal(t, 60, mcpConfig.ServerTimeout)
	assert.Equal(t, "test_", mcpConfig.ToolPrefix)
	assert.True(t, mcpConfig.EnableHotReload)
}

func TestMCPConfig_Defaults(t *testing.T) {
	cfg := &config.Config{
		// Use default values - these will be empty, but GetMCPConfig should return defaults
	}

	mcpConfig := cfg.GetMCPConfig()

	// Test that defaults are applied correctly
	assert.False(t, mcpConfig.Enabled)
	assert.Equal(t, "mcp-servers.json", mcpConfig.ConfigPath)
	assert.Equal(t, 30, mcpConfig.ServerTimeout)
	assert.Equal(t, "", mcpConfig.ToolPrefix)
	assert.False(t, mcpConfig.EnableHotReload)
}