package integration

import (
	"context"
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/denkhaus/agentforge/internal/config"
	"github.com/denkhaus/agentforge/internal/container"
	"github.com/denkhaus/agentforge/internal/types"
)

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

// TestMCPIntegration_DisabledByDefault tests that MCP is disabled by default
func TestMCPIntegration_DisabledByDefault(t *testing.T) {
	// Create config with MCP disabled (default)
	cfg := &config.Config{
		LogLevel:    "info",
		Port:        8080,
		Environment: "test",
		// MCP fields use default values (disabled)
	}

	// Setup container
	injector := container.Setup(cfg)
	defer func() {
		if injector != nil {
			container.Cleanup(injector)
		}
	}()

	// Get tool provider (should be aggregated provider)
	toolProvider, err := do.Invoke[types.ToolProvider](injector)
	if err != nil {
		t.Logf("Error getting tool provider: %v", err)
		require.NoError(t, err)
	}
	assert.NotNil(t, toolProvider)

	// Should have internal tools but no MCP tools
	tools := toolProvider.GetTools()
	assert.NotEmpty(t, tools) // Should have internal tools like getCurrentWeather, createTask

	toolNames := toolProvider.GetToolNames()
	assert.Contains(t, toolNames, "getCurrentWeather")
	assert.Contains(t, toolNames, "createTask")

	// Should not have any MCP-prefixed tools
	for _, name := range toolNames {
		assert.NotContains(t, name, "mcp_", "Should not have MCP tools when disabled")
	}
}

// TestMCPIntegration_EnabledWithInvalidConfig tests MCP enabled but with invalid config
func TestMCPIntegration_EnabledWithInvalidConfig(t *testing.T) {
	// Create config with MCP enabled but invalid config path
	cfg := &config.Config{
		LogLevel:      "info",
		Port:          8080,
		Environment:   "test",
		MCPEnabled:    true,
		MCPConfigPath: "nonexistent-config.json",
	}

	// Setup container - should handle MCP initialization gracefully
	injector := container.Setup(cfg)
	defer container.Cleanup(injector)

	// Get tool provider - should fail gracefully when MCP config is invalid
	toolProvider, err := do.Invoke[types.ToolProvider](injector)
	// Expect error due to invalid MCP config
	if err != nil {
		t.Logf("Expected error due to invalid MCP config: %v", err)
		return // Test passes - invalid config properly handled
	}
	assert.NotNil(t, toolProvider)

	// Should still have internal tools
	tools := toolProvider.GetTools()
	assert.NotEmpty(t, tools)

	toolNames := toolProvider.GetToolNames()
	assert.Contains(t, toolNames, "getCurrentWeather")
	assert.Contains(t, toolNames, "createTask")
}

// TestMCPIntegration_ToolProviderInterface tests that the aggregated provider implements the interface correctly
func TestMCPIntegration_ToolProviderInterface(t *testing.T) {
	cfg := &config.Config{
		LogLevel:    "info",
		Port:        8080,
		Environment: "test",
		MCPEnabled:  false, // Disabled for this test
	}

	injector := container.Setup(cfg)
	defer container.Cleanup(injector)

	toolProvider, err := do.Invoke[types.ToolProvider](injector)
	require.NoError(t, err)

	// Test all interface methods
	t.Run("GetTools", func(t *testing.T) {
		tools := toolProvider.GetTools()
		assert.NotNil(t, tools)
	})

	t.Run("GetToolNames", func(t *testing.T) {
		names := toolProvider.GetToolNames()
		assert.NotNil(t, names)
	})

	t.Run("HasTool", func(t *testing.T) {
		// Should have internal tools
		assert.True(t, toolProvider.HasTool("getCurrentWeather"))
		assert.False(t, toolProvider.HasTool("nonexistent_tool"))
	})

	t.Run("ExecuteTool", func(t *testing.T) {
		// Test with existing tool
		result, err := toolProvider.ExecuteTool(context.Background(), "getCurrentWeather", `{"location": "London"}`)
		assert.NoError(t, err)
		assert.NotEmpty(t, result)

		// Test with nonexistent tool
		_, err = toolProvider.ExecuteTool(context.Background(), "nonexistent_tool", "{}")
		assert.Error(t, err)
	})

	t.Run("ValidateAgentRequirements", func(t *testing.T) {
		// Create mock agent
		agent := &MockAgent{}
		agent.On("GetRequiredTools").Return([]string{"getCurrentWeather"})
		agent.On("GetName").Return("test_agent")

		err := toolProvider.ValidateAgentRequirements(agent)
		assert.NoError(t, err)

		// Test with missing tool
		agent2 := &MockAgent{}
		agent2.On("GetRequiredTools").Return([]string{"nonexistent_tool"})
		agent2.On("GetName").Return("test_agent2")

		err = toolProvider.ValidateAgentRequirements(agent2)
		assert.Error(t, err)
	})
}

// TestMCPIntegration_ConfigurationFlow tests the complete configuration flow
func TestMCPIntegration_ConfigurationFlow(t *testing.T) {
	testCases := []struct {
		name           string
		mcpEnabled     bool
		mcpConfigPath  string
		mcpToolPrefix  string
		expectedError  bool
		expectMCPTools bool
	}{
		{
			name:           "MCP Disabled",
			mcpEnabled:     false,
			mcpConfigPath:  "mcp-servers.json",
			mcpToolPrefix:  "",
			expectedError:  false,
			expectMCPTools: false,
		},
		{
			name:           "MCP Enabled Invalid Config",
			mcpEnabled:     true,
			mcpConfigPath:  "invalid.json",
			mcpToolPrefix:  "",
			expectedError:  true, // Should error due to invalid config
			expectMCPTools: false,
		},
		{
			name:           "MCP Enabled With Prefix",
			mcpEnabled:     true,
			mcpConfigPath:  "invalid.json", // Will fail but test prefix config
			mcpToolPrefix:  "mcp_",
			expectedError:  true,
			expectMCPTools: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{
				LogLevel:      "info",
				Port:          8080,
				Environment:   "test",
				MCPEnabled:    tc.mcpEnabled,
				MCPConfigPath: tc.mcpConfigPath,
				MCPToolPrefix: tc.mcpToolPrefix,
			}

			// Test MCP config extraction
			mcpConfig := cfg.GetMCPConfig()
			assert.Equal(t, tc.mcpEnabled, mcpConfig.Enabled)
			assert.Equal(t, tc.mcpConfigPath, mcpConfig.ConfigPath)
			assert.Equal(t, tc.mcpToolPrefix, mcpConfig.ToolPrefix)

			// Test container setup
			injector := container.Setup(cfg)
			defer container.Cleanup(injector)

			toolProvider, err := do.Invoke[types.ToolProvider](injector)
			if tc.expectedError {
				assert.Error(t, err)
				return // Skip further checks if error expected
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, toolProvider)

				// Should always have internal tools
				assert.True(t, toolProvider.HasTool("getCurrentWeather"))
				assert.True(t, toolProvider.HasTool("createTask"))
			}
		})
	}
}

// TestMCPIntegration_ProviderIsolation tests that providers are properly isolated
func TestMCPIntegration_ProviderIsolation(t *testing.T) {
	cfg := &config.Config{
		LogLevel:    "info",
		Port:        8080,
		Environment: "test",
		MCPEnabled:  false,
	}

	injector := container.Setup(cfg)
	defer container.Cleanup(injector)

	// Get individual providers by name
	internalProvider, err := do.InvokeNamed[types.ToolProvider](injector, "internalProvider")
	require.NoError(t, err)

	mcpProvider, err := do.InvokeNamed[types.ToolProvider](injector, "mcpProvider")
	require.NoError(t, err)

	// Get the main tool provider (which is the aggregated provider)
	aggregatedProvider, err := do.Invoke[types.ToolProvider](injector)
	require.NoError(t, err)

	// Test that they are different instances
	assert.NotEqual(t, internalProvider, mcpProvider)
	assert.NotEqual(t, internalProvider, aggregatedProvider)
	assert.NotEqual(t, mcpProvider, aggregatedProvider)

	// Test that aggregated provider has tools from internal provider
	internalTools := internalProvider.GetTools()
	aggregatedTools := aggregatedProvider.GetTools()

	assert.GreaterOrEqual(t, len(aggregatedTools), len(internalTools))

	// All internal tools should be available in aggregated provider
	for _, tool := range internalTools {
		assert.True(t, aggregatedProvider.HasTool(tool.Name()),
			"Aggregated provider should have internal tool: %s", tool.Name())
	}

	// MCP provider should be empty when disabled
	mcpTools := mcpProvider.GetTools()
	assert.Empty(t, mcpTools, "MCP provider should be empty when disabled")
}

// TestMCPIntegration_ErrorHandling tests error handling scenarios
func TestMCPIntegration_ErrorHandling(t *testing.T) {
	cfg := &config.Config{
		LogLevel:    "info",
		Port:        8080,
		Environment: "test",
		MCPEnabled:  false,
	}

	injector := container.Setup(cfg)
	defer container.Cleanup(injector)

	toolProvider, err := do.Invoke[types.ToolProvider](injector)
	require.NoError(t, err)

	t.Run("ExecuteNonexistentTool", func(t *testing.T) {
		_, err := toolProvider.ExecuteTool(context.Background(), "nonexistent", "{}")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("ValidateAgentWithMissingTools", func(t *testing.T) {
		agent := &MockAgent{}
		agent.On("GetRequiredTools").Return([]string{"missing_tool1", "missing_tool2"})
		agent.On("GetName").Return("test_agent")

		err := toolProvider.ValidateAgentRequirements(agent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing tools")
		assert.Contains(t, err.Error(), "missing_tool1")
		assert.Contains(t, err.Error(), "missing_tool2")
	})

	t.Run("GetToolsForAgentWithMissingTools", func(t *testing.T) {
		agent := &MockAgent{}
		agent.On("GetRequiredTools").Return([]string{"getCurrentWeather", "missing_tool"})
		agent.On("GetName").Return("test_agent")

		tools, err := toolProvider.GetToolsForAgent(agent)
		assert.Error(t, err)
		assert.Nil(t, tools)
	})
}