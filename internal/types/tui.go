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

// EditorModel defines the interface for prompt editor functionality
type EditorModel interface {
	TUIModel
	SetSize(width, height int)
	GetContent() string
	SetContent(content string)
}

// VariablesModel defines the interface for variables management
type VariablesModel interface {
	TUIModel
	SetSize(width, height int)
}

// TestModel defines the interface for prompt testing
type TestModel interface {
	TUIModel
	SetSize(width, height int)
}

// OptimizeModel defines the interface for prompt optimization
type OptimizeModel interface {
	TUIModel
	SetSize(width, height int)
}

// WorkbenchModel defines the interface for the main workbench
type WorkbenchModel interface {
	TUIModel
}