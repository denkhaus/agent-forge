// Package models provides Bubble Tea models for the TUI workbench.
package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/denkhaus/agentforge/internal/tui/components"
	"go.uber.org/zap"
)

// TabType represents the different tabs in the workbench
type TabType int

const (
	EditorTab TabType = iota
	VariablesTab
	TestTab
	OptimizeTab
)

// String returns the string representation of a tab type
func (t TabType) String() string {
	switch t {
	case EditorTab:
		return "Editor"
	case VariablesTab:
		return "Variables"
	case TestTab:
		return "Test"
	case OptimizeTab:
		return "Optimize"
	default:
		return "Unknown"
	}
}

// WorkbenchModel represents the main prompt workbench
type WorkbenchModel struct {
	// Core state
	promptName string
	activeTab  TabType
	logger     *zap.Logger

	// UI dimensions
	width  int
	height int

	// Sub-models
	editorModel    *editorModel
	variablesModel *VariablesModel
	testModel      *testModel
	optimizeModel  *optimizeModel

	// UI components
	help help.Model
	keys keyMap

	// State
	quitting bool
}

// keyMap defines the key bindings for the workbench
type keyMap struct {
	Tab      key.Binding
	ShiftTab key.Binding
	Enter    key.Binding
	Escape   key.Binding
	Quit     key.Binding
	Help     key.Binding
	Save     key.Binding
}

// ShortHelp returns the short help for the key map
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Enter, k.Save, k.Quit, k.Help}
}

// FullHelp returns the full help for the key map
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.ShiftTab, k.Enter, k.Escape},
		{k.Save, k.Quit, k.Help},
	}
}

// NewWorkbenchModel creates a new workbench model
func NewWorkbenchModel(promptName string, logger *zap.Logger) *WorkbenchModel {
	// Initialize key bindings
	keys := keyMap{
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab"),
		),
		ShiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev tab"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save"),
		),
	}

	// Initialize sub-models
	editorModel := NewEditorModel(promptName, logger)
	variablesModel := NewVariablesModel(promptName, logger)
	testModel := NewTestModel(promptName, logger)
	optimizeModel := NewOptimizeModel(promptName, logger)

	return &WorkbenchModel{
		promptName:     promptName,
		activeTab:      EditorTab,
		logger:         logger,
		editorModel:    editorModel,
		variablesModel: variablesModel,
		testModel:      testModel,
		optimizeModel:  optimizeModel,
		help:           help.New(),
		keys:           keys,
		quitting:       false,
	}
}

// Init initializes the workbench model
func (m *WorkbenchModel) Init() tea.Cmd {
	m.logger.Info("Initializing prompt workbench", zap.String("prompt", m.promptName))

	// Initialize all sub-models
	var cmds []tea.Cmd
	cmds = append(cmds, m.editorModel.Init())
	cmds = append(cmds, m.variablesModel.Init())
	cmds = append(cmds, m.testModel.Init())
	cmds = append(cmds, m.optimizeModel.Init())

	return tea.Batch(cmds...)
}

// Update handles messages and updates the model
func (m *WorkbenchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width

		// Update sub-models with new dimensions
		m.editorModel.SetSize(msg.Width, msg.Height-6) // Reserve space for tabs and help
		m.variablesModel.SetSize(msg.Width, msg.Height-6)
		m.testModel.SetSize(msg.Width, msg.Height-6)
		m.optimizeModel.SetSize(msg.Width, msg.Height-6)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.keys.Tab):
			m.nextTab()

		case key.Matches(msg, m.keys.ShiftTab):
			m.prevTab()

		case key.Matches(msg, m.keys.Save):
			// TODO: Implement save functionality
			m.logger.Info("Save requested")

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	// Update the active tab's model
	switch m.activeTab {
	case EditorTab:
		updatedModel, cmd := m.editorModel.Update(msg)
		if editorModel, ok := updatedModel.(*editorModel); ok {
			m.editorModel = editorModel
		}
		cmds = append(cmds, cmd)
	case VariablesTab:
		updatedModel, cmd := m.variablesModel.Update(msg)
		if variablesModel, ok := updatedModel.(*VariablesModel); ok {
			m.variablesModel = variablesModel
		}
		cmds = append(cmds, cmd)
	case TestTab:
		updatedModel, cmd := m.testModel.Update(msg)
		if testModel, ok := updatedModel.(*testModel); ok {
			m.testModel = testModel
		}
		cmds = append(cmds, cmd)
	case OptimizeTab:
		updatedModel, cmd := m.optimizeModel.Update(msg)
		if optimizeModel, ok := updatedModel.(*optimizeModel); ok {
			m.optimizeModel = optimizeModel
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the workbench
func (m *WorkbenchModel) View() string {
	if m.quitting {
		return "Thanks for using AgentForge Prompt Workbench!\n"
	}

	// Render header
	header := m.renderHeader()

	// Render tabs
	tabs := m.renderTabs()

	// Render active tab content
	content := m.renderActiveTab()

	// Render help
	helpView := m.help.View(m.keys)

	// Combine all parts
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tabs,
		content,
		helpView,
	)
}

// renderHeader renders the workbench header
func (m *WorkbenchModel) renderHeader() string {
	title := fmt.Sprintf("AgentForge Prompt Workbench - %s", m.promptName)

	// Use safe lipgloss header
	return components.LipglossHeader(title, m.width)
}

// renderTabs renders the tab navigation
func (m *WorkbenchModel) renderTabs() string {
	var tabs []string

	for i := EditorTab; i <= OptimizeTab; i++ {
		tab := i.String()
		// Use safe tab styling
		styledTab := components.LipglossTab(tab, i == m.activeTab)
		tabs = append(tabs, styledTab)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

// renderActiveTab renders the content of the active tab
func (m *WorkbenchModel) renderActiveTab() string {
	switch m.activeTab {
	case EditorTab:
		return m.editorModel.View()
	case VariablesTab:
		return m.variablesModel.View()
	case TestTab:
		return m.testModel.View()
	case OptimizeTab:
		return m.optimizeModel.View()
	default:
		return "Unknown tab"
	}
}

// nextTab switches to the next tab
func (m *WorkbenchModel) nextTab() {
	m.activeTab = (m.activeTab + 1) % 4
	m.logger.Debug("Switched to tab", zap.String("tab", m.activeTab.String()))
}

// prevTab switches to the previous tab
func (m *WorkbenchModel) prevTab() {
	m.activeTab = (m.activeTab + 3) % 4 // +3 is equivalent to -1 in modulo 4
	m.logger.Debug("Switched to tab", zap.String("tab", m.activeTab.String()))
}
