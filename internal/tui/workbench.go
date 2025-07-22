// Package tui provides the main workbench constructor.
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/denkhaus/agentforge/internal/tui/models"
	"go.uber.org/zap"
)

// NewPromptWorkbench creates a new prompt workbench model
func NewPromptWorkbench(promptName string, logger *zap.Logger) tea.Model {
	return models.NewWorkbenchModel(promptName, logger)
}