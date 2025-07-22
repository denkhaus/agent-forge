package commands

import (
	"github.com/denkhaus/agentforge/internal/prompts"
	"github.com/samber/do"
)

// getPromptServiceFromDI returns a prompt service from the DI container
func getPromptServiceFromDI(injector *do.Injector) (prompts.PromptService, error) {
	return do.Invoke[prompts.PromptService](injector)
}

// Deprecated: Use getPromptServiceFromDI instead
// getPromptService returns a prompt service instance (legacy direct instantiation)
func getPromptService() prompts.PromptService {
	return prompts.NewPromptService()
}