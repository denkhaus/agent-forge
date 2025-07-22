// Enhanced Workbench V3 - Professional TUI using Bubble Tea best practices
// Implements patterns from workspace/bubbletea-reference/

package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
)

// Enhanced color scheme with adaptive support
var (
	primaryColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	accentColor  = lipgloss.AdaptiveColor{Light: "#FF6B6B", Dark: "#FF8E8E"}
	successColor = lipgloss.AdaptiveColor{Light: "#51CF66", Dark: "#69DB7C"}
	warningColor = lipgloss.AdaptiveColor{Light: "#FFD43B", Dark: "#FFE066"}

	// Seamless tab borders (from reference/tabs/seamless-tabs.go)
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")

	// Professional tab styles
	inactiveTabStyle = lipgloss.NewStyle().
				Border(inactiveTabBorder, true).
				BorderForeground(primaryColor).
				Padding(0, 1)

		// TODO: Copy is deprecated
	activeTabStyle = inactiveTabStyle.Copy().
			Border(activeTabBorder, true).
			Foreground(lipgloss.Color("229")).
			Background(primaryColor)

	// Content window with seamless connection
	windowStyle = lipgloss.NewStyle().
			BorderForeground(primaryColor).
			Padding(1, 2).
			Border(lipgloss.NormalBorder()).
			UnsetBorderTop() // Key: seamless connection to tabs

	// Focus-aware styling (from reference/split-editors/)
	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor)

	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240"))
)

// Tab types for the workbench
type TabType int

const (
	EditorTab TabType = iota
	VariablesTab
	TestTab
	OptimizeTab
)

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

// Enhanced key mapping (from reference/patterns/common-patterns.go)
type workbenchKeyMap struct {
	Tab      key.Binding
	ShiftTab key.Binding
	Enter    key.Binding
	Escape   key.Binding
	Save     key.Binding
	Test     key.Binding
	Optimize key.Binding
	Quit     key.Binding
	Help     key.Binding
}

// Implement help.KeyMap interface
func (k workbenchKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Enter, k.Save, k.Test, k.Optimize, k.Quit, k.Help}
}

func (k workbenchKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.ShiftTab, k.Enter, k.Escape},
		{k.Save, k.Test, k.Optimize, k.Quit, k.Help},
	}
}

func newWorkbenchKeyMap() workbenchKeyMap {
	return workbenchKeyMap{
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
			key.WithHelp("enter", "select/edit"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save"),
		),
		Test: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "run tests"),
		),
		Optimize: key.NewBinding(
			key.WithKeys("ctrl+o"),
			key.WithHelp("ctrl+o", "optimize"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}

// Enhanced Workbench V3 model
type WorkbenchV3 struct {
	// Core state
	promptName string
	activeTab  TabType
	// TODO: remove that, we hav a package logger
	logger *zap.Logger

	// UI dimensions
	width  int
	height int

	// Components (from reference patterns)
	editor     textarea.Model
	variables  table.Model
	testRunner *TestRunner
	optimizer  *OptimizerV3

	// Progress tracking (from reference/progress/)
	progress progress.Model

	// UI components
	help help.Model
	keys workbenchKeyMap

	// State management
	focused    bool
	quitting   bool
	testing    bool
	optimizing bool
}

// Test runner with progress tracking
type TestRunner struct {
	providers []ModelProvider
	results   []TestResultV3
	progress  float64
	active    bool
}

type ModelProvider struct {
	Name    string
	Model   string
	Enabled bool
	Status  string
}

type TestResultV3 struct {
	Provider string
	Model    string
	Response string
	Latency  string
	Success  bool
	Tokens   int
	Cost     float64
}

// Enhanced optimizer with iteration tracking
type OptimizerV3 struct {
	iterations    int
	maxIterations int
	progress      float64
	active        bool
	suggestions   []OptimizationStep
	bestScore     float64
}

type OptimizationStep struct {
	Iteration   int
	Suggestion  string
	Score       float64
	Improvement float64
}

// Custom border function (from reference/tabs/seamless-tabs.go)
func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

// NewWorkbenchV3 creates an enhanced workbench with professional patterns
func NewWorkbenchV3(promptName string, logger *zap.Logger) *WorkbenchV3 {
	// Enhanced textarea with professional styling
	editor := textarea.New()
	editor.Placeholder = "Enter your prompt here..."
	editor.Focus()
	editor.CharLimit = 10000
	editor.ShowLineNumbers = true

	// Apply focus-aware styling
	editor.FocusedStyle.Base = focusedBorderStyle
	editor.BlurredStyle.Base = blurredBorderStyle

	// Set initial content
	editor.SetValue(fmt.Sprintf(`# %s Prompt

You are an AI assistant that helps with {{task}}.

Please analyze the following {{input}}:

{{content}}

Provide a detailed response that includes:
1. Analysis of the input
2. Key findings or insights
3. Recommendations or next steps

Format your response in {{format}} format.`, promptName))

	// Enhanced variables table (from reference/table/rich-data-display.go)
	columns := []table.Column{
		{Title: "Name", Width: 15},
		{Title: "Type", Width: 10},
		{Title: "Default", Width: 20},
		{Title: "Required", Width: 8},
		{Title: "Description", Width: 25},
	}

	rows := []table.Row{
		{"task", "string", "analyze", "✓", "The task to perform"},
		{"input", "string", "code", "✓", "Type of input to process"},
		{"content", "string", "", "✓", "The content to analyze"},
		{"format", "string", "markdown", "", "Output format"},
	}

	variables := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(8),
	)

	// Professional table styling
	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(accentColor).
		BorderBottom(true).
		Bold(true)
	tableStyles.Selected = tableStyles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(accentColor).
		Bold(false)
	variables.SetStyles(tableStyles)

	// Initialize test runner
	testRunner := &TestRunner{
		providers: []ModelProvider{
			{"OpenAI", "gpt-4", true, "Ready"},
			{"OpenAI", "gpt-3.5-turbo", true, "Ready"},
			{"Anthropic", "claude-3-sonnet", false, "Disabled"},
			{"Google", "gemini-pro", false, "Disabled"},
		},
		results: []TestResultV3{},
		active:  false,
	}

	// Initialize optimizer
	optimizer := &OptimizerV3{
		maxIterations: 5,
		suggestions:   []OptimizationStep{},
		active:        false,
	}

	// Progress bar with gradient
	progressBar := progress.New(progress.WithDefaultGradient())

	return &WorkbenchV3{
		promptName: promptName,
		activeTab:  EditorTab,
		logger:     logger,
		editor:     editor,
		variables:  variables,
		testRunner: testRunner,
		optimizer:  optimizer,
		progress:   progressBar,
		help:       help.New(),
		keys:       newWorkbenchKeyMap(),
		focused:    true,
	}
}

// Init initializes the enhanced workbench
func (m *WorkbenchV3) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages with enhanced patterns
func (m *WorkbenchV3) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		m.updateComponentSizes()

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
			m.logger.Info("Save requested", zap.String("prompt", m.promptName))
			// TODO: Implement save functionality

		case key.Matches(msg, m.keys.Test):
			if m.activeTab == TestTab && !m.testing {
				return m, m.startTesting()
			}

		case key.Matches(msg, m.keys.Optimize):
			if m.activeTab == OptimizeTab && !m.optimizing {
				return m, m.startOptimization()
			}

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}

	// Progress animation (from reference/progress/animated-progress.go)
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
	}

	// Update active component (from reference/split-editors/multi-component-focus.go)
	switch m.activeTab {
	case EditorTab:
		m.editor, cmd = m.editor.Update(msg)
		cmds = append(cmds, cmd)
	case VariablesTab:
		m.variables, cmd = m.variables.Update(msg)
		cmds = append(cmds, cmd)
	case TestTab:
		// Update test runner
	case OptimizeTab:
		// Update optimizer
	}

	return m, tea.Batch(cmds...)
}

// View renders the enhanced workbench with seamless styling
func (m *WorkbenchV3) View() string {
	if m.quitting {
		return lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			Render("Thanks for using AgentForge Enhanced Workbench V3!")
	}

	if m.width < 100 || m.height < 30 {
		return lipgloss.NewStyle().
			Foreground(warningColor).
			Render("Terminal too small. Minimum size: 100x30")
	}

	// Render seamless tabs (from reference/tabs/seamless-tabs.go)
	tabs := m.renderSeamlessTabs()

	// Render content with seamless connection
	content := m.renderTabContent()

	// Render help
	helpView := m.help.View(m.keys)

	return tabs + "\n" + content + "\n" + helpView
}

// Missing methods implementation

// updateComponentSizes updates component dimensions based on terminal size
func (m *WorkbenchV3) updateComponentSizes() {
	contentWidth := m.width - 6
	contentHeight := m.height - 10

	m.editor.SetWidth(contentWidth)
	m.editor.SetHeight(contentHeight)

	m.progress.Width = contentWidth - 4
}

// nextTab switches to the next tab
func (m *WorkbenchV3) nextTab() {
	m.activeTab = (m.activeTab + 1) % 4
	m.updateFocus()
	m.logger.Debug("Switched to next tab", zap.String("tab", m.activeTab.String()))
}

// prevTab switches to the previous tab
func (m *WorkbenchV3) prevTab() {
	m.activeTab = (m.activeTab + 3) % 4 // +3 is equivalent to -1 in modulo 4
	m.updateFocus()
	m.logger.Debug("Switched to previous tab", zap.String("tab", m.activeTab.String()))
}

// updateFocus manages focus between components
func (m *WorkbenchV3) updateFocus() {
	switch m.activeTab {
	case EditorTab:
		m.editor.Focus()
		m.variables.Blur()
	case VariablesTab:
		m.editor.Blur()
		m.variables.Focus()
	default:
		m.editor.Blur()
		m.variables.Blur()
	}
}

// startTesting initiates the testing process
func (m *WorkbenchV3) startTesting() tea.Cmd {
	m.testing = true
	m.testRunner.active = true
	m.logger.Info("Starting prompt testing")

	// TODO: Implement actual testing logic
	return func() tea.Msg {
		return tea.Msg("test_started")
	}
}

// startOptimization initiates the optimization process
func (m *WorkbenchV3) startOptimization() tea.Cmd {
	m.optimizing = true
	m.optimizer.active = true
	m.logger.Info("Starting prompt optimization")

	// TODO: Implement actual optimization logic
	return func() tea.Msg {
		return tea.Msg("optimization_started")
	}
}

// renderSeamlessTabs renders tabs with seamless borders
func (m *WorkbenchV3) renderSeamlessTabs() string {
	var renderedTabs []string

	for i := EditorTab; i <= OptimizeTab; i++ {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == OptimizeTab, i == m.activeTab

		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}

		// Dynamic border adjustment for seamless appearance
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}

		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(i.String()))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

// renderTabContent renders the content area with seamless connection
func (m *WorkbenchV3) renderTabContent() string {
	var content string

	switch m.activeTab {
	case EditorTab:
		content = m.renderEditorContent()
	case VariablesTab:
		content = m.renderVariablesContent()
	case TestTab:
		content = m.renderTestContent()
	case OptimizeTab:
		content = m.renderOptimizeContent()
	}

	// Calculate content width based on tabs
	tabsWidth := lipgloss.Width(m.renderSeamlessTabs())
	contentWidth := max(tabsWidth, m.width-4)

	return windowStyle.Width(contentWidth).Render(content)
}

// Content rendering methods
func (m *WorkbenchV3) renderEditorContent() string {
	return m.editor.View()
}

func (m *WorkbenchV3) renderVariablesContent() string {
	return m.variables.View()
}

func (m *WorkbenchV3) renderTestContent() string {
	content := "Multi-Model Testing\n\n"

	for _, provider := range m.testRunner.providers {
		status := "❌"
		if provider.Enabled {
			status = "✅"
		}
		content += fmt.Sprintf("%s %s - %s (%s)\n", status, provider.Name, provider.Model, provider.Status)
	}

	if m.testing {
		content += "\n" + m.progress.View()
	}

	return content
}

func (m *WorkbenchV3) renderOptimizeContent() string {
	content := "AI-Powered Optimization\n\n"
	content += fmt.Sprintf("Iterations: %d/%d\n", m.optimizer.iterations, m.optimizer.maxIterations)
	content += fmt.Sprintf("Best Score: %.2f\n", m.optimizer.bestScore)

	if m.optimizing {
		content += "\n" + m.progress.View()
	}

	return content
}

// Helper function
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
