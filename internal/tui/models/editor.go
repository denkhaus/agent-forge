// Package models provides the editor model for prompt editing.
package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
)

// EditorModel handles prompt editing functionality
type EditorModel struct {
	promptName string
	logger     *zap.Logger
	
	// UI components
	textarea textarea.Model
	
	// UI dimensions
	width  int
	height int
	
	// State
	focused bool
}

// NewEditorModel creates a new editor model
func NewEditorModel(promptName string, logger *zap.Logger) *EditorModel {
	// Initialize textarea
	ta := textarea.New()
	ta.Placeholder = "Enter your prompt here..."
	ta.Focus()
	ta.CharLimit = 10000
	ta.SetWidth(80)
	ta.SetHeight(20)
	
	// Load existing prompt content if available
	// TODO: Load from file system
	ta.SetValue(fmt.Sprintf("# %s Prompt\n\nYou are an AI assistant that helps with...\n\nPlease {{task}} the following {{input}}:\n\n{{content}}", promptName))
	
	return &EditorModel{
		promptName: promptName,
		logger:     logger,
		textarea:   ta,
		focused:    true,
	}
}

// Init initializes the editor model
func (m *EditorModel) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages for the editor
func (m *EditorModel) Update(msg tea.Msg) (*EditorModel, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			// Save the prompt
			m.logger.Info("Saving prompt", zap.String("name", m.promptName))
			// TODO: Implement save functionality
			return m, nil
		}
	}
	
	// Update textarea
	m.textarea, cmd = m.textarea.Update(msg)
	
	return m, cmd
}

// View renders the editor
func (m *EditorModel) View() string {
	// Create the editor section
	editorStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1).
		Width(m.width - 4).
		Height(m.height - 4)
	
	// Create info panel
	info := m.renderInfo()
	
	// Create preview panel
	preview := m.renderPreview()
	
	// Layout: editor on left, info and preview on right
	leftPanel := editorStyle.Render(m.textarea.View())
	
	rightPanelStyle := lipgloss.NewStyle().
		Width(40).
		Height(m.height - 4)
	
	rightPanel := rightPanelStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			info,
			preview,
		),
	)
	
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		rightPanel,
	)
}

// renderInfo renders the prompt information panel
func (m *EditorModel) renderInfo() string {
	content := m.textarea.Value()
	
	// Count variables (simple {{variable}} detection)
	variables := m.extractVariables(content)
	
	infoStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1).
		Margin(0, 0, 1, 0)
	
	info := fmt.Sprintf("📊 Prompt Info\n\n"+
		"Name: %s\n"+
		"Length: %d chars\n"+
		"Lines: %d\n"+
		"Variables: %d\n\n"+
		"Variables found:\n%s",
		m.promptName,
		len(content),
		strings.Count(content, "\n")+1,
		len(variables),
		strings.Join(variables, ", "),
	)
	
	return infoStyle.Render(info)
}

// renderPreview renders the prompt preview panel
func (m *EditorModel) renderPreview() string {
	previewStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1)
	
	// Simple preview with sample variable substitution
	content := m.textarea.Value()
	preview := m.generatePreview(content)
	
	previewContent := fmt.Sprintf("🔍 Live Preview\n\n%s", preview)
	
	return previewStyle.Render(previewContent)
}

// extractVariables extracts variable names from the prompt content
func (m *EditorModel) extractVariables(content string) []string {
	var variables []string
	
	// Simple regex-like extraction for {{variable}} patterns
	parts := strings.Split(content, "{{")
	for i := 1; i < len(parts); i++ {
		if closingIndex := strings.Index(parts[i], "}}"); closingIndex != -1 {
			variable := strings.TrimSpace(parts[i][:closingIndex])
			if variable != "" {
				// Check if already exists
				exists := false
				for _, v := range variables {
					if v == variable {
						exists = true
						break
					}
				}
				if !exists {
					variables = append(variables, variable)
				}
			}
		}
	}
	
	return variables
}

// generatePreview generates a preview with sample variable values
func (m *EditorModel) generatePreview(content string) string {
	// Sample variable substitutions
	sampleVars := map[string]string{
		"task":    "analyze",
		"input":   "code",
		"content": "function hello() { console.log('Hello, World!'); }",
		"name":    "example",
		"type":    "JavaScript",
	}
	
	preview := content
	for variable, value := range sampleVars {
		placeholder := fmt.Sprintf("{{%s}}", variable)
		preview = strings.ReplaceAll(preview, placeholder, value)
	}
	
	// Truncate if too long
	if len(preview) > 300 {
		preview = preview[:300] + "..."
	}
	
	return preview
}

// SetSize updates the editor dimensions
func (m *EditorModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	
	// Update textarea size
	m.textarea.SetWidth(width - 50) // Leave space for right panel
	m.textarea.SetHeight(height - 8) // Leave space for borders and padding
}

// GetContent returns the current prompt content
func (m *EditorModel) GetContent() string {
	return m.textarea.Value()
}

// SetContent sets the prompt content
func (m *EditorModel) SetContent(content string) {
	m.textarea.SetValue(content)
}