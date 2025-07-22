package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/denkhaus/agentforge/internal/config"
	"github.com/denkhaus/agentforge/internal/container"
)

// TestSimpleContainerSetup tests basic container setup without complex dependencies
func TestSimpleContainerSetup(t *testing.T) {
	cfg := &config.Config{
		LogLevel:    "info",
		Port:        8080,
		Environment: "test",
		MCPEnabled:  false,
	}

	// Test container setup
	injector := container.Setup(cfg)
	require.NotNil(t, injector)

	// Test cleanup
	container.Cleanup(injector)
}

// TestMCPConfigOnly tests just the MCP configuration without container
func TestMCPConfigOnly(t *testing.T) {
	cfg := &config.Config{
		LogLevel:    "info",
		Port:        8080,
		Environment: "test",
		MCPEnabled:  false,
	}

	mcpConfig := cfg.GetMCPConfig()
	assert.NotNil(t, mcpConfig)
	assert.False(t, mcpConfig.Enabled)
	assert.Equal(t, "mcp-servers.json", mcpConfig.ConfigPath)
}