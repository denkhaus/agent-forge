package config

import (
	"fmt"
	"os"
)

// ValidateAPIKeys ensures at least one API key is configured for LLM providers.
func (c *Config) ValidateAPIKeys() error {
	if c.apiKeys == nil {
		c.initializeAPIKeys()
	}

	if len(c.apiKeys) == 0 {
		return fmt.Errorf(`no API keys configured for any LLM provider.
Please set at least one of:
  - GOOGLE_API_KEY
  - OPENAI_API_KEY  
  - ANTHROPIC_API_KEY
  - AZURE_API_KEY`)
	}

	return nil
}


// HasProviderConfigured checks if a specific provider has an API key configured.
func (c *Config) HasProviderConfigured(provider string) bool {
	return c.GetAPIKey(provider) != ""
}

// ValidateAPIKeysIfNeeded validates API keys only for commands that require them.
func (c *Config) ValidateAPIKeysIfNeeded() error {
	// Check command line arguments to determine which command is being run
	if len(os.Args) > 1 {
		commandName := os.Args[1]
		
		// Component commands only need GitHub access, not LLM providers
		if commandName == "component" || commandName == "comp" {
			return nil
		}
		
		// Entity-specific commands don't need LLM API keys
		if commandName == "tool" || commandName == "prompt" || commandName == "agent" {
			return nil
		}
		
		// Version command doesn't need any API keys
		if commandName == "version" || commandName == "v" {
			return nil
		}
	}
	
	// For other commands (server, chat), require LLM API keys
	return c.ValidateAPIKeys()
}
