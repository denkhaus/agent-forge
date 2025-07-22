// Package tui provides the main workbench constructor.
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

// NewPromptWorkbench creates a new prompt workbench model
func NewPromptWorkbench(promptName string, logger *zap.Logger) tea.Model {
	// Use enhanced WorkbenchV3 with professional patterns
	return NewWorkbenchV3(promptName, logger)
}