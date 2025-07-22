// Package providers contains implementations of various provider interfaces.
package providers

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"text/template"

	"github.com/samber/do"
	"go.uber.org/zap"

	"github.com/denkhaus/agentforge/internal/types"
)

// promptProvider is a private implementation of types.PromptProvider interface.
type promptProvider struct {
	log     *zap.Logger
	prompts map[string]types.Prompt
	mutex   sync.RWMutex
}

// NewPromptProvider creates a new prompt provider instance.
func NewPromptProvider(injector *do.Injector) (types.PromptProvider, error) {
	log, err := do.Invoke[*zap.Logger](injector)
	if err != nil {
		return nil, fmt.Errorf("failed to get logger: %w", err)
	}
	
	prompts, err := do.Invoke[map[string]types.Prompt](injector)
	if err != nil {
		return nil, fmt.Errorf("failed to get prompts map: %w", err)
	}

	provider := &promptProvider{
		log:     log,
		prompts: prompts,
	}

	provider.log.Info("Prompt provider initialized", zap.Int("prompt_count", len(provider.prompts)))

	return provider, nil
}

// Startup performs any initialization required before the service can run.
// Removed Startup method - Provider interface eliminated

// GetSystemPrompt returns the system prompt for the AI agent.
func (p *promptProvider) GetSystemPrompt() string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if systemPrompt, exists := p.prompts["system"]; exists {
		return systemPrompt.GetTemplate()
	}

	// Fallback system prompt if not found
	p.log.Warn("System prompt not found, using fallback")
	return "You are a helpful AI assistant."
}

// FormatPrompt formats a prompt template with the given variables.
func (p *promptProvider) FormatPrompt(templateStr string, vars map[string]any) (string, error) {
	tmpl, err := template.New("prompt").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, vars); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}

// Removed GetPromptTemplate - use GetPrompt().GetTemplate() instead

// Note: basePrompt implementation moved to internal/prompts/base.go to avoid duplication

// Removed SetPromptTemplate - use RegisterPrompt() instead

// NewCustomPrompt creates a new custom prompt - delegated to prompts package.
func NewCustomPrompt(name, description, template string) types.Prompt {
	return &customPrompt{
		name:        name,
		description: description,
		template:    template,
	}
}

// customPrompt is a simple implementation for dynamic prompts created by the provider.
type customPrompt struct {
	name        string
	description string
	template    string
}

func (p *customPrompt) GetName() string        { return p.name }
func (p *customPrompt) GetDescription() string { return p.description }
func (p *customPrompt) GetTemplate() string    { return p.template }

func (p *customPrompt) Clone(overrides map[string]any) types.Prompt {
	clone := &customPrompt{
		name:        p.name,
		description: p.description,
		template:    p.template,
	}

	if name, ok := overrides["name"].(string); ok {
		clone.name = name
	}
	if description, ok := overrides["description"].(string); ok {
		clone.description = description
	}
	if template, ok := overrides["template"].(string); ok {
		clone.template = template
	}

	return clone
}

func (p *customPrompt) Render(variables map[string]any) (string, error) {
	tmpl, err := template.New(p.name).Parse(p.template)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, variables); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GetPrompts returns all available prompts.
func (p *promptProvider) GetPrompts() map[string]types.Prompt {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	result := make(map[string]types.Prompt)
	for name, prompt := range p.prompts {
		result[name] = prompt
	}
	return result
}

// RegisterPrompt registers a new prompt with the provider.
func (p *promptProvider) RegisterPrompt(prompt types.Prompt) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, exists := p.prompts[prompt.GetName()]; exists {
		return fmt.Errorf("prompt %s already exists", prompt.GetName())
	}

	p.prompts[prompt.GetName()] = prompt
	p.log.Info("Prompt registered",
		zap.String("name", prompt.GetName()),
		zap.String("description", prompt.GetDescription()))

	return nil
}

// Removed RenderPrompt - use GetPrompt().Render() instead

// GetPrompt retrieves a prompt by name, with fallback to a default prompt if not found.
func (p *promptProvider) GetPrompt(promptName string) types.Prompt {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if prompt, exists := p.prompts[promptName]; exists {
		return prompt
	}

	// Fallback: create a basic prompt if not found
	p.log.Warn("Prompt not found, using fallback", zap.String("prompt_name", promptName))
	return &fallbackPrompt{
		name:        promptName,
		description: "Fallback prompt for " + promptName,
		template:    "You are a helpful AI assistant.",
	}
}

// fallbackPrompt is a simple fallback implementation.
type fallbackPrompt struct {
	name        string
	description string
	template    string
}

func (p *fallbackPrompt) GetName() string                         { return p.name }
func (p *fallbackPrompt) GetDescription() string                  { return p.description }
func (p *fallbackPrompt) GetTemplate() string                     { return p.template }
func (p *fallbackPrompt) Clone(_ map[string]any) types.Prompt     { return p }
func (p *fallbackPrompt) Render(_ map[string]any) (string, error) { return p.template, nil }
