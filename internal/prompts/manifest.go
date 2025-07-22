package prompts

// PromptData contains all the data needed to generate prompt files
type PromptData struct {
	Name        string
	Description string
	Author      string
	License     string
	PromptType  string
	Language    string
	Template    string
	Variables   []string
}

// ManifestGenerator handles prompt manifest file generation (legacy interface)
type ManifestGenerator interface {
	CreateComponentYAML(promptDir string, data PromptData) error
	CreateTemplateFile(promptDir string, template string) error
	CreateVariablesFile(promptDir string, variables []string) error
	CreateReadmeFile(promptDir string, data PromptData) error
}

// manifestGenerator implements ManifestGenerator (empty implementation)
type manifestGenerator struct{}

// NewManifestGenerator creates a new manifest generator
func NewManifestGenerator() ManifestGenerator {
	return &manifestGenerator{}
}

// Legacy methods - empty implementations
func (mg *manifestGenerator) CreateComponentYAML(promptDir string, data PromptData) error {
	return nil
}

func (mg *manifestGenerator) CreateTemplateFile(promptDir string, template string) error {
	return nil
}

func (mg *manifestGenerator) CreateVariablesFile(promptDir string, variables []string) error {
	return nil
}

func (mg *manifestGenerator) CreateReadmeFile(promptDir string, data PromptData) error {
	return nil
}