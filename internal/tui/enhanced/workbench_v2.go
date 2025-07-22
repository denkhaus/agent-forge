// Package enhanced provides advanced TUI components inspired by Bubble Tea examples.
package enhanced

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
)

// Enhanced color scheme inspired by Bubble Tea examples
var (
	highlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	accentColor    = lipgloss.AdaptiveColor{Light: "#FF6B6B", Dark: "#FF8E8E"}
	successColor   = lipgloss.AdaptiveColor{Light: "#51CF66", Dark: "#69DB7C"}
	warningColor   = lipgloss.AdaptiveColor{Light: "#FFD43B", Dark: "#FFE066"}
	
	// Enhanced tab styles with seamless borders
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	
	inactiveTabStyle = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(highlightColor).
		Padding(0, 1)
	
	activeTabStyle = inactiveTabStyle.Copy().
		Border(activeTabBorder, true).
		Foreground(lipgloss.Color("229")).
		Background(highlightColor)
	
	// Content window style that connects to tabs
	windowStyle = lipgloss.NewStyle().
		BorderForeground(highlightColor).
		Padding(1, 2).
		Border(lipgloss.NormalBorder()).
		UnsetBorderTop()
	
	// Enhanced component styles
	focusedBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(highlightColor)
	
	blurredBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))
)

// TabType represents workbench tabs
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
		return "📝 Editor"
	case VariablesTab:
		return "🔧 Variables"
	case TestTab:
		return "🧪 Test"
	case OptimizeTab:
		return "🚀 Optimize"
	default:
		return "Unknown"
	}
}

// EnhancedWorkbench represents the enhanced prompt workbench
type EnhancedWorkbench struct {
	// Core state
	promptName string
	activeTab  TabType
	logger     *zap.Logger
	
	// UI dimensions
	width  int
	height int
	
	// Enhanced components
	editor     textarea.Model
	variables  table.Model
	testRunner *TestRunner
	optimizer  *OptimizerModel
	
	// Progress tracking
	progress progress.Model
	
	// UI components
	help help.Model
	keys keyMap
	
	// State
	focused  bool
	quitting bool
}

// TestRunner manages multi-model testing
type TestRunner struct {
	providers []ModelProvider
	results   []TestResult
	progress  float64
	testing   bool
}

// OptimizerModel manages AI-powered optimization
type OptimizerModel struct {
	iterations    int
	maxIterations int
	progress      float64
	optimizing    bool
	suggestions   []string
}

// ModelProvider represents an AI model provider
type ModelProvider struct {
	Name    string
	Model   string
	Enabled bool
	Status  string
}

// TestResult represents test results
type TestResult struct {
	Provider string
	Model    string
	Response string
	Latency  string
	Success  bool
}

// keyMap defines enhanced key bindings
type keyMap struct {
	Tab      key.Binding
	ShiftTab key.Binding
	Enter    key.Binding
	Escape   key.Binding
	Quit     key.Binding
	Help     key.Binding
	Save     key.Binding
	Test     key.Binding
	Optimize key.Binding
}

// NewEnhancedWorkbench creates a new enhanced workbench
func NewEnhancedWorkbench(promptName string, logger *zap.Logger) *EnhancedWorkbench {
	// Initialize enhanced textarea
	editor := textarea.New()
	editor.Placeholder = "Enter your prompt here..."
	editor.Focus()
	editor.CharLimit = 10000
	editor.ShowLineNumbers = true
	
	// Enhanced styling for editor
	editor.FocusedStyle.Base = focusedBorderStyle
	editor.BlurredStyle.Base = blurredBorderStyle
	
	// Initialize variables table
	columns := []table.Column{
		{Title: "Name", Width: 15},
		{Title: "Type", Width: 10},
		{Title: "Default", Width: 20},
		{Title: "Required", Width: 8},
	}
	
	rows := []table.Row{
		{"task", "string", "analyze", "✓"},
		{"input", "string", "code", "✓"},
		{"content", "string", "", "✓"},
	}
	
	variables := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(10),
	)
	
	// Enhanced table styling
	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(highlightColor).
		BorderBottom(true).
		Bold(true)
	tableStyles.Selected = tableStyles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(highlightColor).
		Bold(false)
	variables.SetStyles(tableStyles)
	
	// Initialize progress bar
	progressBar := progress.New(progress.WithDefaultGradient())
	
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
			key.WithHelp("enter", "select/edit"),
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
		Test: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "run tests"),
		),
		Optimize: key.NewBinding(
			key.WithKeys("ctrl+o"),
			key.WithHelp("ctrl+o", "optimize"),
		),
	}
	
	return &EnhancedWorkbench{
		promptName: promptName,
		activeTab:  EditorTab,
		logger:     logger,
		editor:     editor,
		variables:  variables,
		testRunner: &TestRunner{
			providers: []ModelProvider{
				{"OpenAI", "gpt-4", true, "Ready"},
				{"OpenAI", "gpt-3.5-turbo", true, "Ready"},
				{"Anthropic", "claude-3-sonnet", false, "Disabled"},
			},
		},
		optimizer: &OptimizerModel{
			maxIterations: 5,
		},
		progress: progressBar,
		help:     help.New(),
		keys:     keys,
		focused:  true,
	}
}

// Helper function for custom tab borders
func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

// Init initializes the enhanced workbench
func (m *EnhancedWorkbench) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages for the enhanced workbench
func (m *EnhancedWorkbench) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		
		// Update component sizes
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
			// TODO: Implement save functionality
			m.logger.Info("Save requested")
			
		case key.Matches(msg, m.keys.Test):
			if m.activeTab == TestTab {
				return m, m.runTests()
			}
			
		case key.Matches(msg, m.keys.Optimize):
			if m.activeTab == OptimizeTab {
				return m, m.startOptimization()
			}
			
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	
	// Update active component
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

// View renders the enhanced workbench
func (m *EnhancedWorkbench) View() string {
	if m.quitting {
		return "Thanks for using AgentForge Enhanced Workbench!\n"
	}
	
	if m.width < 100 || m.height < 30 {
		return "Terminal too small. Minimum size: 100x30\n"
	}
	
	// Render enhanced tabs
	tabs := m.renderEnhancedTabs()
	
	// Render content with seamless connection to tabs
	content := m.renderTabContent()
	
	// Render help
	helpView := m.help.View(m.keys)
	
	return tabs + "\n" + content + "\n" + helpView
}

// renderEnhancedTabs renders tabs with seamless borders
func (m *EnhancedWorkbench) renderEnhancedTabs() string {
	var renderedTabs []string
	
	for i := EditorTab; i <= OptimizeTab; i++ {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == OptimizeTab, i == m.activeTab
		
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		
		// Adjust borders for seamless connection
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
func (m *EnhancedWorkbench) renderTabContent() string {
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
	tabsWidth := lipgloss.Width(m.renderEnhancedTabs())
	contentWidth := max(tabsWidth, m.width-4)
	
	return windowStyle.Width(contentWidth).Render(content)
}

// Content rendering methods
func (m *EnhancedWorkbench) renderEditorContent() string {
	return m.editor.View()
}

func (m *EnhancedWorkbench) renderVariablesContent() string {
	return m.variables.View()
}

func (m *EnhancedWorkbench) renderTestContent() string {
	var content strings.Builder
	content.WriteString("🧪 Multi-Model Testing\n\n")
	
	for _, provider := range m.testRunner.providers {
		status := "❌"
		if provider.Enabled {
			status = "✅"
		}
		content.WriteString(fmt.Sprintf("%s %s - %s (%s)\n", status, provider.Name, provider.Model, provider.Status))
	}
	
	if m.testRunner.testing {
		content.WriteString("\n" + m.progress.View())
	}
	
	return content.String()
}

func (m *EnhancedWorkbench) renderOptimizeContent() string {
	var content strings.Builder
	content.WriteString("🚀 AI-Powered Optimization\n\n")
	content.WriteString(fmt.Sprintf("Iterations: %d/%d\n", m.optimizer.iterations, m.optimizer.maxIterations))
	
	if m.optimizer.optimizing {
		content.WriteString("\n" + m.progress.View())
	}
	
	return content.String()
}

// Helper methods
func (m *EnhancedWorkbench) nextTab() {
	m.activeTab = (m.activeTab + 1) % 4
	m.updateFocus()
}

func (m *EnhancedWorkbench) prevTab() {
	m.activeTab = (m.activeTab + 3) % 4
	m.updateFocus()
}

func (m *EnhancedWorkbench) updateFocus() {
	// Update component focus based on active tab
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

func (m *EnhancedWorkbench) updateComponentSizes() {
	// Update component sizes based on terminal dimensions
	contentWidth := m.width - 6
	contentHeight := m.height - 10
	
	m.editor.SetWidth(contentWidth)
	m.editor.SetHeight(contentHeight)
	
	m.progress.Width = contentWidth - 4
}

func (m *EnhancedWorkbench) runTests() tea.Cmd {
	m.testRunner.testing = true
	// TODO: Implement actual testing
	return nil
}

func (m *EnhancedWorkbench) startOptimization() tea.Cmd {
	m.optimizer.optimizing = true
	// TODO: Implement actual optimization
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}