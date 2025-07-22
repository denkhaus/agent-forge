// Package tui provides Terminal User Interface components using Bubble Tea.
package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/denkhaus/agentforge/internal/types"
	"go.uber.org/zap"
)

// manager implements the TUIManager interface
type manager struct {
	logger *zap.Logger
}

// NewManager creates a new TUI manager
func NewManager(logger *zap.Logger) types.TUIManager {
	return &manager{
		logger: logger,
	}
}

// RunPromptCreator runs the interactive prompt creator
func (m *manager) RunPromptCreator() error {
	m.logger.Info("Starting prompt creator TUI")
	
	// TODO: Implement prompt creator
	return fmt.Errorf("prompt creator not yet implemented")
}

// RunPromptCreatorWithName runs the prompt creator with a specific name
func (m *manager) RunPromptCreatorWithName(name string) error {
	m.logger.Info("Starting prompt creator TUI", zap.String("name", name))
	
	// TODO: Implement prompt creator with name
	return fmt.Errorf("prompt creator with name not yet implemented")
}

// RunPromptEditor runs the prompt editor TUI
func (m *manager) RunPromptEditor(name string) error {
	m.logger.Info("Starting prompt editor TUI", zap.String("name", name))
	
	// Create and run the prompt workbench
	workbench := NewPromptWorkbench(name, m.logger)
	
	program := tea.NewProgram(workbench, tea.WithAltScreen())
	_, err := program.Run()
	
	if err != nil {
		m.logger.Error("TUI program failed", zap.Error(err))
		return fmt.Errorf("failed to run prompt editor: %w", err)
	}
	
	m.logger.Info("Prompt editor TUI completed successfully")
	return nil
}

// RunPromptVariableEditor runs the variable editor TUI
func (m *manager) RunPromptVariableEditor(name string, data interface{}) error {
	m.logger.Info("Starting variable editor TUI", 
		zap.String("name", name),
		zap.Any("data", data))
	
	// TODO: Implement variable editor
	return fmt.Errorf("variable editor not yet implemented")
}

// RunPromptWorkbench runs the full prompt workbench (new method for enhanced functionality)
func (m *manager) RunPromptWorkbench(name string) error {
	m.logger.Info("Starting prompt workbench TUI", zap.String("name", name))
	
	// Create and run the enhanced prompt workbench
	workbench := NewPromptWorkbench(name, m.logger)
	
	program := tea.NewProgram(workbench, tea.WithAltScreen())
	_, err := program.Run()
	
	if err != nil {
		m.logger.Error("TUI workbench failed", zap.Error(err))
		return fmt.Errorf("failed to run prompt workbench: %w", err)
	}
	
	m.logger.Info("Prompt workbench completed successfully")
	return nil
}