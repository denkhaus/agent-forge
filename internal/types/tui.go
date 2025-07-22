package types

import (
	tea "github.com/charmbracelet/bubbletea"
)

// TUIModel represents a Bubble Tea model interface
type TUIModel interface {
	tea.Model
}

// PromptCreator defines the interface for interactive prompt creation
type PromptCreator interface {
	Run() error
}

// TUIManager manages TUI operations
type TUIManager interface {
	RunPromptCreator() error
	RunPromptCreatorWithName(name string) error
	RunPromptEditor(name string) error
	RunPromptVariableEditor(name string, data interface{}) error
	RunPromptWorkbench(name string) error
}