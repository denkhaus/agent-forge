// Package models provides the variables model for variable management.
package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/denkhaus/agentforge/internal/tui/components"
	"go.uber.org/zap"
)

// VariableType represents the type of a variable
type VariableType string

const (
	StringType  VariableType = "string"
	NumberType  VariableType = "number"
	BooleanType VariableType = "boolean"
	ArrayType   VariableType = "array"
	ObjectType  VariableType = "object"
)

// Variable represents a prompt variable
type Variable struct {
	Name         string
	Type         VariableType
	DefaultValue string
	Description  string
	Required     bool
}

// VariablesModel handles variable management functionality
type VariablesModel struct {
	promptName string
	logger     *zap.Logger
	
	// Variables
	variables []Variable
	
	// UI dimensions
	width  int
	height int
	
	// State
	selectedIndex int
	editing       bool
}

// NewVariablesModel creates a new variables model
func NewVariablesModel(promptName string, logger *zap.Logger) *VariablesModel {
	// Initialize with some sample variables
	variables := []Variable{
		{
			Name:         "task",
			Type:         StringType,
			DefaultValue: "analyze",
			Description:  "The task to perform",
			Required:     true,
		},
		{
			Name:         "input",
			Type:         StringType,
			DefaultValue: "code",
			Description:  "The type of input",
			Required:     true,
		},
		{
			Name:         "content",
			Type:         StringType,
			DefaultValue: "",
			Description:  "The content to process",
			Required:     true,
		},
	}
	
	return &VariablesModel{
		promptName:    promptName,
		logger:        logger,
		variables:     variables,
		selectedIndex: 0,
		editing:       false,
	}
}

// Init initializes the variables model
func (m *VariablesModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the variables model
func (m *VariablesModel) Update(msg tea.Msg) (*VariablesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down", "j":
			if m.selectedIndex < len(m.variables)-1 {
				m.selectedIndex++
			}
		case "enter":
			// Toggle editing mode
			m.editing = !m.editing
		case "n":
			// Add new variable
			if !m.editing {
				m.addNewVariable()
			}
		case "d":
			// Delete selected variable
			if !m.editing && len(m.variables) > 0 {
				m.deleteVariable(m.selectedIndex)
			}
		}
	}
	
	return m, nil
}

// View renders the variables management interface
func (m *VariablesModel) View() string {
	if m.width < 60 || m.height < 15 {
		return "Variables: Terminal too small"
	}
	
	// Create enhanced layout with safe styling
	leftPanel := m.renderVariableList()
	rightPanel := m.renderVariableEditor()
	
	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
}

// renderVariableList renders the list of variables
func (m *VariablesModel) renderVariableList() string {
	// Safe styling for variable list
	listStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("6")). // Cyan
		Padding(1).
		Width(m.width/2 - 2).
		Height(m.height - 4)
	
	var items []string
	items = append(items, components.SafeColors.Bold + "Variables" + components.SafeColors.Reset)
	items = append(items, "")
	
	for i, variable := range m.variables {
		prefix := "  "
		if i == m.selectedIndex {
			prefix = components.SafeColors.Yellow + "> " + components.SafeColors.Reset
		}
		
		required := ""
		if variable.Required {
			required = components.SafeColors.Red + " *" + components.SafeColors.Reset
		}
		
		item := fmt.Sprintf("%s%s (%s)%s", prefix, variable.Name, variable.Type, required)
		if variable.DefaultValue != "" {
			item += fmt.Sprintf(" = %s", variable.DefaultValue)
		}
		
		items = append(items, item)
	}
	
	items = append(items, "")
	items = append(items, components.SafeColors.Dim + "Commands:" + components.SafeColors.Reset)
	items = append(items, "  n - New variable")
	items = append(items, "  d - Delete variable")
	items = append(items, "  enter - Edit variable")
	
	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	return listStyle.Render(content)
}

// renderVariableEditor renders the variable editor panel
func (m *VariablesModel) renderVariableEditor() string {
	editorStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1).
		Width(m.width/2 - 2).
		Height(m.height - 4)
	
	if len(m.variables) == 0 {
		return editorStyle.Render("No variables defined.\nPress 'n' to add a new variable.")
	}
	
	if m.selectedIndex >= len(m.variables) {
		return editorStyle.Render("Invalid selection.")
	}
	
	variable := m.variables[m.selectedIndex]
	
	var content []string
	content = append(content, "📝 Variable Editor")
	content = append(content, "")
	content = append(content, fmt.Sprintf("Name: %s", variable.Name))
	content = append(content, fmt.Sprintf("Type: %s", variable.Type))
	content = append(content, fmt.Sprintf("Required: %t", variable.Required))
	content = append(content, fmt.Sprintf("Default: %s", variable.DefaultValue))
	content = append(content, "")
	content = append(content, "Description:")
	content = append(content, variable.Description)
	
	if m.editing {
		content = append(content, "")
		content = append(content, "🔄 EDITING MODE")
		content = append(content, "Press enter to save")
	}
	
	editorContent := lipgloss.JoinVertical(lipgloss.Left, content...)
	return editorStyle.Render(editorContent)
}

// addNewVariable adds a new variable to the list
func (m *VariablesModel) addNewVariable() {
	newVar := Variable{
		Name:         fmt.Sprintf("variable_%d", len(m.variables)+1),
		Type:         StringType,
		DefaultValue: "",
		Description:  "New variable",
		Required:     false,
	}
	
	m.variables = append(m.variables, newVar)
	m.selectedIndex = len(m.variables) - 1
	m.logger.Info("Added new variable", zap.String("name", newVar.Name))
}

// deleteVariable removes a variable from the list
func (m *VariablesModel) deleteVariable(index int) {
	if index < 0 || index >= len(m.variables) {
		return
	}
	
	deletedVar := m.variables[index]
	m.variables = append(m.variables[:index], m.variables[index+1:]...)
	
	// Adjust selected index
	if m.selectedIndex >= len(m.variables) && len(m.variables) > 0 {
		m.selectedIndex = len(m.variables) - 1
	}
	if len(m.variables) == 0 {
		m.selectedIndex = 0
	}
	
	m.logger.Info("Deleted variable", zap.String("name", deletedVar.Name))
}

// SetSize updates the variables model dimensions
func (m *VariablesModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// GetVariables returns the current variables
func (m *VariablesModel) GetVariables() []Variable {
	return m.variables
}

// SetVariables sets the variables
func (m *VariablesModel) SetVariables(variables []Variable) {
	m.variables = variables
	if m.selectedIndex >= len(variables) {
		m.selectedIndex = 0
	}
}