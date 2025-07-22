package commands

import (
	"github.com/denkhaus/agentforge/internal/prompts"
)

// getPromptService returns a prompt service instance
func getPromptService() prompts.PromptService {
	return prompts.NewPromptService()
}