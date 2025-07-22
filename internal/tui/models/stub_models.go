// Package models provides stub implementations for TUI models.
package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

// stubVariablesModel is a placeholder implementation for variables model
type stubVariablesModel struct {
	promptName string
	logger     *zap.Logger
	width      int
	height     int
}

func (m *stubVariablesModel) Init() tea.Cmd {
	return nil
}

func (m *stubVariablesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *stubVariablesModel) View() string {
	return "Variables Tab - Coming Soon!\n\nThis will allow you to:\n- Define prompt variables\n- Set default values\n- Configure variable types"
}

func (m *stubVariablesModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// stubTestModel is a placeholder implementation for test model
type stubTestModel struct {
	promptName string
	logger     *zap.Logger
	width      int
	height     int
}

func (m *stubTestModel) Init() tea.Cmd {
	return nil
}

func (m *stubTestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *stubTestModel) View() string {
	return "Test Tab - Coming Soon!\n\nThis will allow you to:\n- Test prompts with multiple AI models\n- Compare responses\n- Measure performance"
}

func (m *stubTestModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// stubOptimizeModel is a placeholder implementation for optimize model
type stubOptimizeModel struct {
	promptName string
	logger     *zap.Logger
	width      int
	height     int
}

func (m *stubOptimizeModel) Init() tea.Cmd {
	return nil
}

func (m *stubOptimizeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *stubOptimizeModel) View() string {
	return "Optimize Tab - Coming Soon!\n\nThis will allow you to:\n- AI-powered prompt optimization\n- Performance analysis\n- Iterative improvements"
}

func (m *stubOptimizeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}