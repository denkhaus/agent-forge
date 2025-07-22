// Package models provides the test model for prompt testing.
package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

// testModel handles prompt testing functionality (private implementation)
type testModel struct {
	promptName string
	logger     *zap.Logger
	width      int
	height     int
}

// NewTestModel creates a new test model
func NewTestModel(promptName string, logger *zap.Logger) *testModel {
	return &testModel{
		promptName: promptName,
		logger:     logger,
	}
}

// Init initializes the test model
func (m *testModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the test model
func (m *testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View renders the test model
func (m *testModel) View() string {
	return "Test Tab - Multi-Model Testing\n\nThis will allow you to:\n- Test prompts with multiple AI models\n- Compare responses\n- Measure performance\n- Analyze results\n\n(Implementation coming soon)"
}

// SetSize updates the test model dimensions
func (m *testModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}