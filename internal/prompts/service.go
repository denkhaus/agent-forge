package prompts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/denkhaus/agentforge/internal/templates"
	"go.uber.org/zap"
)

// PromptService defines the interface for prompt operations
type PromptService interface {
	CreatePromptStructure(name string) error
	LoadPromptData(name string) (*PromptData, error)
	SavePromptData(name string, data *PromptData) error
	ValidatePromptName(name string) error
	ListLocalPrompts() ([]*PromptData, error)
	PullPrompt(repo, version string, force bool) error
	PushPrompt(name, repo, message, tag string) error
	ExecutePrompt(name string, variables map[string]string) (string, error)
}

// promptService implements PromptService interface
type promptService struct {
	manifestGenerator ManifestGenerator
	templateGenerator templates.PromptTemplateGenerator
}

// NewPromptService creates a new prompt service
func NewPromptService() PromptService {
	return &promptService{
		manifestGenerator: NewManifestGenerator(),
		templateGenerator: templates.NewPromptTemplateGenerator(),
	}
}

// CreatePromptStructure creates the basic filesystem structure for a prompt
func (ps *promptService) CreatePromptStructure(name string) error {
	log.Info("Creating prompt structure", zap.String("name", name))

	// Validate name
	if err := ps.ValidatePromptName(name); err != nil {
		return err
	}

	// Create directory
	promptDir := filepath.Join("prompts", name)
	if err := os.MkdirAll(promptDir, 0755); err != nil {
		return fmt.Errorf("failed to create prompt directory: %w", err)
	}

	// Create default prompt data
	defaultData := &PromptData{
		Name:        name,
		Description: fmt.Sprintf("AI prompt for %s", name),
		Author:      "Your Name <your.email@example.com>",
		License:     "MIT",
		PromptType:  "TEMPLATE",
		Language:    "en",
		Template:    fmt.Sprintf("# %s\n\nYou are an AI assistant. Please help with:\n\n{{task}}\n\nProvide a detailed and helpful response.", name),
		Variables:   []string{"task"},
	}

	// Convert PromptData to PromptTemplateData and use template generator
	templateData := templates.PromptTemplateData{
		Name:         defaultData.Name,
		DisplayName:  strings.Title(strings.ReplaceAll(defaultData.Name, "-", " ")),
		Description:  defaultData.Description,
		Author:       defaultData.Author,
		License:      defaultData.License,
		PromptType:   defaultData.PromptType,
		Language:     defaultData.Language,
		Instructions: "Please analyze the following input and provide a comprehensive response.",
		OutputFormat: "Provide a structured and detailed response.",
		PrimaryVariable: "task",
	}
	
	// Convert variables
	for _, varName := range defaultData.Variables {
		templateData.Variables = append(templateData.Variables, templates.PromptVariable{
			Name:        varName,
			Type:        "string",
			Description: fmt.Sprintf("Description for %s", varName),
			Required:    true,
		})
	}
	
	// Set expected outputs
	templateData.ExpectedOutputs = []string{
		"Comprehensive analysis",
		"Actionable recommendations", 
		"Clear explanations",
		"Structured results",
	}
	
	// Use template generator to create all files
	if err := ps.templateGenerator.GeneratePromptFiles(templateData, promptDir); err != nil {
		return fmt.Errorf("failed to generate prompt files: %w", err)
	}

	log.Info("Prompt structure created successfully",
		zap.String("name", name),
		zap.String("path", promptDir))

	return nil
}

// LoadPromptData loads existing prompt data from filesystem
func (ps *promptService) LoadPromptData(name string) (*PromptData, error) {
	log.Info("Loading prompt data", zap.String("name", name))

	promptDir := filepath.Join("prompts", name)

	// Check if prompt exists
	if _, err := os.Stat(promptDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("prompt '%s' does not exist", name)
	}

	// Load template content
	templatePath := filepath.Join(promptDir, "template.txt")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template.txt: %w", err)
	}

	// TODO: Parse component.yaml and variables.json for complete data
	// For now, return basic data structure
	data := &PromptData{
		Name:        name,
		Description: fmt.Sprintf("AI prompt for %s", name),
		Author:      "Your Name <your.email@example.com>",
		License:     "MIT",
		PromptType:  "TEMPLATE",
		Language:    "en",
		Template:    string(templateContent),
		Variables:   []string{}, // TODO: Parse from variables.json
	}

	return data, nil
}

// SavePromptData saves prompt data to filesystem
func (ps *promptService) SavePromptData(name string, data *PromptData) error {
	log.Info("Saving prompt data", zap.String("name", name))

	promptDir := filepath.Join("prompts", name)

	// Update all files
	if err := ps.manifestGenerator.CreateComponentYAML(promptDir, *data); err != nil {
		return fmt.Errorf("failed to update component.yaml: %w", err)
	}

	if err := ps.manifestGenerator.CreateTemplateFile(promptDir, data.Template); err != nil {
		return fmt.Errorf("failed to update template.txt: %w", err)
	}

	if err := ps.manifestGenerator.CreateVariablesFile(promptDir, data.Variables); err != nil {
		return fmt.Errorf("failed to update variables.json: %w", err)
	}

	if err := ps.manifestGenerator.CreateReadmeFile(promptDir, *data); err != nil {
		return fmt.Errorf("failed to update README.md: %w", err)
	}

	log.Info("Prompt data saved successfully", zap.String("name", name))
	return nil
}

// ValidatePromptName validates the prompt name
func (ps *promptService) ValidatePromptName(name string) error {
	if name == "" {
		return fmt.Errorf("prompt name cannot be empty")
	}

	// Check for valid characters (lowercase, hyphens, numbers)
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-') {
			return fmt.Errorf("prompt name must contain only lowercase letters, numbers, and hyphens")
		}
	}

	// Check if already exists
	promptDir := filepath.Join("prompts", name)
	if _, err := os.Stat(promptDir); err == nil {
		return fmt.Errorf("prompt '%s' already exists", name)
	}

	return nil
}

// ListLocalPrompts lists all locally installed prompts
func (ps *promptService) ListLocalPrompts() ([]*PromptData, error) {
	log.Info("Listing local prompts")

	promptsDir := "prompts"

	// Check if prompts directory exists
	if _, err := os.Stat(promptsDir); os.IsNotExist(err) {
		return []*PromptData{}, nil
	}

	// Read directory entries
	entries, err := os.ReadDir(promptsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read prompts directory: %w", err)
	}

	var prompts []*PromptData
	for _, entry := range entries {
		if entry.IsDir() {
			// Try to load prompt data
			if data, err := ps.LoadPromptData(entry.Name()); err == nil {
				prompts = append(prompts, data)
			}
		}
	}

	log.Info("Found local prompts", zap.Int("count", len(prompts)))
	return prompts, nil
}

// PullPrompt pulls a prompt from a remote repository
func (ps *promptService) PullPrompt(repo, version string, force bool) error {
	log.Info("Pulling prompt",
		zap.String("repo", repo),
		zap.String("version", version),
		zap.Bool("force", force))

	// TODO: Implement actual Git pulling logic
	// For now, this is a placeholder
	return fmt.Errorf("pull functionality not yet implemented")
}

// PushPrompt pushes a prompt to a remote repository
func (ps *promptService) PushPrompt(name, repo, message, tag string) error {
	log.Info("Pushing prompt",
		zap.String("name", name),
		zap.String("repo", repo),
		zap.String("message", message),
		zap.String("tag", tag))

	// TODO: Implement actual Git pushing logic
	// For now, this is a placeholder
	return fmt.Errorf("push functionality not yet implemented")
}

// ExecutePrompt executes a prompt with given variables
func (ps *promptService) ExecutePrompt(name string, variables map[string]string) (string, error) {
	log.Info("Executing prompt",
		zap.String("name", name),
		zap.Any("variables", variables))

	// Load prompt data
	data, err := ps.LoadPromptData(name)
	if err != nil {
		return "", err
	}

	// Replace variables in template
	result := data.Template
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}

	log.Info("Prompt executed successfully", zap.String("name", name))
	return result, nil
}
