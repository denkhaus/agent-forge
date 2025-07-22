package providers

import (
	"strings"
	"testing"

	"github.com/samber/do"
	"go.uber.org/zap"

	"github.com/denkhaus/agentforge/internal/types"
)

func TestNewPromptProvider(t *testing.T) {
	// Setup DI container for testing
	injector := do.New()
	defer func() {
		if err := injector.Shutdown(); err != nil {
			t.Logf("Failed to shutdown injector: %v", err)
		}
	}()

	// Register test dependencies
	do.ProvideValue(injector, zap.NewNop())

	do.ProvideValue(injector, map[string]types.Prompt{
		"test": &testPrompt{name: "test", template: "Hello {{.name}}"},
	})

	provider, err := NewPromptProvider(injector)
	if err != nil {
		t.Fatalf("Failed to create prompt provider: %v", err)
	}

	if provider == nil {
		t.Fatal("Expected non-nil prompt provider")
	}

	// Test that provider implements the simplified interface
	var _ = provider
}

func TestPromptProviderGetSystemPrompt(t *testing.T) {
	provider := createTestPromptProvider(t)

	systemPrompt := provider.GetSystemPrompt()
	if systemPrompt == "" {
		t.Error("Expected non-empty system prompt")
	}

	// Should return fallback if no system prompt exists
	if systemPrompt != "You are a helpful AI assistant." {
		t.Errorf("Expected fallback system prompt, got: %s", systemPrompt)
	}
}

func TestPromptProviderFormatPrompt(t *testing.T) {
	provider := createTestPromptProvider(t)

	template := "Hello {{.name}}, welcome to {{.app}}"
	vars := map[string]any{
		"name": "Alice",
		"app":  "MCP Planner",
	}

	result, err := provider.FormatPrompt(template, vars)
	if err != nil {
		t.Fatalf("Failed to format prompt: %v", err)
	}

	expected := "Hello Alice, welcome to MCP Planner"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestPromptProviderFormatPromptError(t *testing.T) {
	provider := createTestPromptProvider(t)

	// Invalid template syntax
	template := "Hello {{.name"
	vars := map[string]any{"name": "Alice"}

	_, err := provider.FormatPrompt(template, vars)
	if err == nil {
		t.Error("Expected error for invalid template syntax")
	}
}

func TestPromptProviderGetPrompts(t *testing.T) {
	provider := createTestPromptProvider(t)

	prompts := provider.GetPrompts()
	if len(prompts) == 0 {
		t.Error("Expected at least one prompt")
	}

	if _, exists := prompts["test"]; !exists {
		t.Error("Expected 'test' prompt to exist")
	}
}

func TestPromptProviderGetPrompt(t *testing.T) {
	provider := createTestPromptProvider(t)

	// Test existing prompt
	prompt := provider.GetPrompt("test")
	if prompt == nil {
		t.Fatal("Expected non-nil prompt")
	}

	if prompt.GetName() != "test" {
		t.Errorf("Expected prompt name 'test', got %q", prompt.GetName())
	}

	// Test non-existing prompt (should return fallback)
	fallbackPrompt := provider.GetPrompt("nonexistent")
	if fallbackPrompt == nil {
		t.Fatal("Expected non-nil fallback prompt")
	}

	if fallbackPrompt.GetName() != "nonexistent" {
		t.Errorf("Expected fallback prompt name 'nonexistent', got %q", fallbackPrompt.GetName())
	}
}

func TestPromptProviderRegisterPrompt(t *testing.T) {
	provider := createTestPromptProvider(t)

	newPrompt := &testPrompt{
		name:        "new_prompt",
		description: "A new test prompt",
		template:    "New template: {{.value}}",
	}

	err := provider.RegisterPrompt(newPrompt)
	if err != nil {
		t.Fatalf("Failed to register prompt: %v", err)
	}

	// Verify prompt was registered
	retrievedPrompt := provider.GetPrompt("new_prompt")
	if retrievedPrompt == nil {
		t.Fatal("Expected registered prompt to be retrievable")
	}

	if retrievedPrompt.GetName() != "new_prompt" {
		t.Errorf("Expected prompt name 'new_prompt', got %q", retrievedPrompt.GetName())
	}
}

func TestPromptProviderRegisterDuplicatePrompt(t *testing.T) {
	provider := createTestPromptProvider(t)

	duplicatePrompt := &testPrompt{
		name:        "test", // Same as existing prompt
		description: "Duplicate prompt",
		template:    "Duplicate template",
	}

	err := provider.RegisterPrompt(duplicatePrompt)
	if err == nil {
		t.Error("Expected error when registering duplicate prompt")
	}
}

// Helper functions and test types

type testPrompt struct {
	name        string
	description string
	template    string
}

func (p *testPrompt) GetName() string        { return p.name }
func (p *testPrompt) GetDescription() string { return p.description }
func (p *testPrompt) GetTemplate() string    { return p.template }

func (p *testPrompt) Render(variables map[string]any) (string, error) {
	// Simple template rendering for testing
	result := p.template
	for key, value := range variables {
		placeholder := "{{." + key + "}}"
		if str, ok := value.(string); ok {
			result = strings.ReplaceAll(result, placeholder, str)
		}
	}
	return result, nil
}

func createTestPromptProvider(t *testing.T) types.PromptProvider {
	injector := do.New()
	t.Cleanup(func() {
		if err := injector.Shutdown(); err != nil {
			t.Logf("Failed to shutdown injector: %v", err)
		}
	})

	// Register test dependencies
	do.ProvideValue(injector, zap.NewNop())

	do.ProvideValue(injector, map[string]types.Prompt{
		"test": &testPrompt{
			name:        "test",
			description: "Test prompt",
			template:    "Hello {{.name}}",
		},
	})

	provider, err := NewPromptProvider(injector)
	if err != nil {
		t.Fatalf("Failed to create test prompt provider: %v", err)
	}

	return provider
}
