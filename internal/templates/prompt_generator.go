package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/denkhaus/agentforge/internal/logger"
	"github.com/denkhaus/agentforge/internal/templates/prompts"
	"go.uber.org/zap"
)

var promptLog *zap.Logger

func init() {
	promptLog = logger.WithPackage("templates.prompts")
}

// PromptTemplateGenerator provides prompt template generation functionality.
type PromptTemplateGenerator interface {
	GeneratePromptFiles(data PromptTemplateData, outputDir string) error
}

// PromptTemplateData contains all data needed for prompt template generation
type PromptTemplateData struct {
	// Basic metadata
	Name         string
	DisplayName  string
	Namespace    string
	Version      string
	Description  string
	Author       string
	License      string
	Homepage     string
	Documentation string
	Tags         []string
	Categories   []string
	Keywords     []string
	Created      string
	
	// Prompt-specific
	PromptType        string
	Language          string
	ContextWindow     int
	SupportsStreaming bool
	Instructions      string
	OutputFormat      string
	PrimaryVariable   string
	
	// Variables
	Variables []PromptVariable
	
	// Configuration
	Temperature       float64
	MaxTokens         int
	StopSequences     []string
	ModelPreferences  []string
	ExpectedOutputs   []string
	
	// System
	Stability     string
	Maturity      string
	ForgeVersion  string
	Platforms     []string
	SchemaVersion string
	GeneratedAt   string
	Generator     string
}

// PromptVariable represents a template variable
type PromptVariable struct {
	Name        string
	Type        string
	Description string
	Required    bool
	Default     interface{}
}

// promptTemplateGenerator implements PromptTemplateGenerator interface.
type promptTemplateGenerator struct{}

// NewPromptTemplateGenerator creates a new prompt template generator.
func NewPromptTemplateGenerator() PromptTemplateGenerator {
	return &promptTemplateGenerator{}
}

// GeneratePromptFiles generates all prompt files using virtual filesystem templates
func (ptg *promptTemplateGenerator) GeneratePromptFiles(data PromptTemplateData, outputDir string) error {
	promptLog.Info("Generating prompt files", 
		zap.String("name", data.Name),
		zap.String("output_dir", outputDir))
	
	// Set default values if not provided
	if data.Created == "" {
		data.Created = time.Now().Format(time.RFC3339)
	}
	if data.GeneratedAt == "" {
		data.GeneratedAt = time.Now().Format(time.RFC3339)
	}
	if data.Generator == "" {
		data.Generator = "agentforge-template-generator"
	}
	if data.SchemaVersion == "" {
		data.SchemaVersion = "1.0.0"
	}
	if data.Namespace == "" {
		data.Namespace = "default"
	}
	if data.Version == "" {
		data.Version = "1.0.0"
	}
	if data.DisplayName == "" {
		data.DisplayName = data.Name
	}
	if data.ContextWindow == 0 {
		data.ContextWindow = 4096
	}
	if data.Temperature == 0 {
		data.Temperature = 0.7
	}
	if data.MaxTokens == 0 {
		data.MaxTokens = 2048
	}
	if data.Stability == "" {
		data.Stability = "experimental"
	}
	if data.Maturity == "" {
		data.Maturity = "alpha"
	}
	if data.ForgeVersion == "" {
		data.ForgeVersion = "0.1.0"
	}
	if len(data.Platforms) == 0 {
		data.Platforms = []string{"linux", "darwin", "windows"}
	}
	if len(data.ModelPreferences) == 0 {
		data.ModelPreferences = []string{"gpt-4", "claude-3-sonnet", "gemini-pro"}
	}
	if len(data.StopSequences) == 0 {
		data.StopSequences = []string{}
	}
	
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Template files to generate
	templateFiles := map[string]string{
		"component.yaml.tmpl": "component.yaml",
		"template.txt.tmpl":   "template.txt",
		"variables.json.tmpl": "variables.json",
		"README.md.tmpl":      "README.md",
	}
	
	// Generate each file
	for templateFile, outputFile := range templateFiles {
		promptLog.Debug("Generating file", 
			zap.String("template", templateFile),
			zap.String("output", outputFile))
		
		// Read template from embedded filesystem
		templateContent, err := prompts.Templates.ReadFile(templateFile)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", templateFile, err)
		}
		
		// Parse template
		tmpl, err := template.New(templateFile).Parse(string(templateContent))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", templateFile, err)
		}
		
		// Execute template
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return fmt.Errorf("failed to execute template %s: %w", templateFile, err)
		}
		
		// Write to output file
		outputPath := filepath.Join(outputDir, outputFile)
		if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", outputPath, err)
		}
		
		promptLog.Debug("Generated file successfully", 
			zap.String("path", outputPath),
			zap.Int("size", buf.Len()))
	}
	
	promptLog.Info("Prompt files generated successfully", 
		zap.String("name", data.Name),
		zap.Int("files", len(templateFiles)))
	
	return nil
}