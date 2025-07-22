// Package config provides MCP-specific configuration management.
package config

// MCPConfig holds configuration for MCP server integration.
type MCPConfig struct {
	// Enabled controls whether MCP integration is active
	Enabled bool `envconfig:"MCP_ENABLED" default:"false"`
	
	// ConfigPath is the path to the MCP servers configuration file
	ConfigPath string `envconfig:"MCP_CONFIG_PATH" default:"mcp-servers.json"`
	
	// ServerTimeout is the timeout for MCP server operations in seconds
	ServerTimeout int `envconfig:"MCP_SERVER_TIMEOUT" default:"30"`
	
	// ToolPrefix is an optional prefix for MCP tool names to avoid conflicts
	ToolPrefix string `envconfig:"MCP_TOOL_PREFIX" default:""`
	
	// EnableHotReload enables hot-reloading of MCP server configuration
	EnableHotReload bool `envconfig:"MCP_HOT_RELOAD" default:"false"`
}

// GetMCPConfig returns the MCP configuration from the main config.
func (c *Config) GetMCPConfig() *MCPConfig {
	configPath := c.MCPConfigPath
	if configPath == "" {
		configPath = "mcp-servers.json"
	}
	
	serverTimeout := c.MCPServerTimeout
	if serverTimeout == 0 {
		serverTimeout = 30
	}
	
	return &MCPConfig{
		Enabled:         c.MCPEnabled,
		ConfigPath:      configPath,
		ServerTimeout:   serverTimeout,
		ToolPrefix:      c.MCPToolPrefix,
		EnableHotReload: c.MCPHotReload,
	}
}

// Helper methods to extract MCP config from environment
func (c *Config) getMCPEnabled() bool {
	return c.MCPEnabled
}

func (c *Config) getMCPConfigPath() string {
	return c.MCPConfigPath
}

func (c *Config) getMCPServerTimeout() int {
	return c.MCPServerTimeout
}

func (c *Config) getMCPToolPrefix() string {
	return c.MCPToolPrefix
}

func (c *Config) getMCPHotReload() bool {
	return c.MCPHotReload
}