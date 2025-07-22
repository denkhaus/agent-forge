// Package config provides configuration management for the MCP Planner application.
package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	cli "github.com/urfave/cli/v2"
)

// EnvProcessFunc is a package-level variable that can be overridden for testing.
var EnvProcessFunc = envconfig.Process

// Config holds all configuration parameters for the application.
type Config struct {
	// LogLevel sets the logging level (debug, info, warn, error)
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	// Port is the HTTP server port
	Port int `envconfig:"PORT" default:"8080"`

	// DatabaseURL is the connection string for the database (optional, for future use)
	DatabaseURL string `envconfig:"DATABASE_URL"`

	// Environment specifies the deployment environment (development, staging, production)
	Environment string `envconfig:"ENVIRONMENT" default:"development"`

	// API Keys for various LLM providers
	GoogleAPIKey    string `envconfig:"GOOGLE_API_KEY"`
	OpenAIAPIKey    string `envconfig:"OPENAI_API_KEY"`
	AnthropicAPIKey string `envconfig:"ANTHROPIC_API_KEY"`
	AzureAPIKey     string `envconfig:"AZURE_API_KEY"`
	
	// GitHub API configuration
	GitHubAPIKey    string `envconfig:"GITHUB_API_KEY"`
	GitHubToken     string `envconfig:"GITHUB_TOKEN"`
	

	// Generic API key mapping for extensibility
	apiKeys map[string]string

	// MCP Configuration
	MCPEnabled       bool   `envconfig:"MCP_ENABLED" default:"false"`
	MCPConfigPath    string `envconfig:"MCP_CONFIG_PATH" default:"mcp-servers.json"`
	MCPServerTimeout int    `envconfig:"MCP_SERVER_TIMEOUT" default:"30"`
	MCPToolPrefix    string `envconfig:"MCP_TOOL_PREFIX" default:""`
	MCPHotReload     bool   `envconfig:"MCP_HOT_RELOAD" default:"false"`
}

// Load reads configuration from environment variables and returns a Config struct.
//
// Returns:
//   - *Config: The loaded configuration
//   - error: Returns an error if environment variable processing fails or required variables are missing
func Load(c *cli.Context) (*Config, error) {
	var cfg Config
	if err := EnvProcessFunc("", &cfg); err != nil {
		return nil, err
	}

	// Override config with CLI flags if provided
	if c.IsSet("log-level") {
		cfg.LogLevel = c.String("log-level")
	}
	if c.IsSet("port") {
		cfg.Port = c.Int("port")
	}

	// Initialize API keys map
	cfg.initializeAPIKeys()

	// Validate that at least one API key is configured for commands that need it
	// Skip validation for component commands that only need GitHub access
	if err := cfg.ValidateAPIKeysIfNeeded(); err != nil {
		return nil, err
	}

	return &cfg, nil
}


// initializeAPIKeys populates the API keys map from environment variables.
func (c *Config) initializeAPIKeys() {
	c.apiKeys = make(map[string]string)

	// Direct mapping of API keys to provider names
	if c.GoogleAPIKey != "" {
		c.apiKeys["googleai"] = c.GoogleAPIKey
		c.apiKeys["google"] = c.GoogleAPIKey
	}
	if c.OpenAIAPIKey != "" {
		c.apiKeys["openai"] = c.OpenAIAPIKey
	}
	if c.AnthropicAPIKey != "" {
		c.apiKeys["anthropic"] = c.AnthropicAPIKey
		c.apiKeys["claude"] = c.AnthropicAPIKey
	}
	if c.AzureAPIKey != "" {
		c.apiKeys["azure"] = c.AzureAPIKey
		c.apiKeys["azureopenai"] = c.AzureAPIKey
	}
	if c.GitHubAPIKey != "" {
		c.apiKeys["github"] = c.GitHubAPIKey
		c.apiKeys["github_api"] = c.GitHubAPIKey
	}
	if c.GitHubToken != "" {
		c.apiKeys["github_token"] = c.GitHubToken
		c.apiKeys["gh_token"] = c.GitHubToken
	}
}


// GetLogLevel returns the configured log level.
func (c *Config) GetLogLevel() string {
	return c.LogLevel
}

// GetPort returns the configured port.
func (c *Config) GetPort() int {
	return c.Port
}

// GetDatabaseURL returns the configured database URL.
func (c *Config) GetDatabaseURL() string {
	return c.DatabaseURL
}

// GetEnvironment returns the configured environment.
func (c *Config) GetEnvironment() string {
	return c.Environment
}

// GetAPIKey returns the API key for a specific provider.
func (c *Config) GetAPIKey(provider string) string {
	if c.apiKeys == nil {
		c.initializeAPIKeys()
	}
	return c.apiKeys[provider]
}

// GetAPIKeys returns all configured API keys.
func (c *Config) GetAPIKeys() map[string]string {
	if c.apiKeys == nil {
		c.initializeAPIKeys()
	}

	// Return a copy to prevent external modification
	result := make(map[string]string)
	for k, v := range c.apiKeys {
		result[k] = v
	}
	return result
}

// GetTimeout returns a default timeout for operations.
func (c *Config) GetTimeout() time.Duration {
	return 30 * time.Second
}

