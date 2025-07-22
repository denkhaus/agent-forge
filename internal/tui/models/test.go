// Package models provides the test model for multi-model testing.
package models

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
)

// ModelProvider represents an AI model provider
type ModelProvider struct {
	Name     string
	Model    string
	Enabled  bool
	APIKey   string
	Endpoint string
}

// TestResult represents the result of a prompt test
type TestResult struct {
	Provider  string
	Model     string
	Response  string
	Latency   time.Duration
	Tokens    int
	Cost      float64
	Timestamp time.Time
	Error     error
}

// TestModel handles multi-model testing functionality
type TestModel struct {
	promptName string
	logger     *zap.Logger
	
	// Test configuration
	providers []ModelProvider
	variables map[string]string
	
	// Test results
	results []TestResult
	
	// UI dimensions
	width  int
	height int
	
	// State
	selectedProvider int
	testing          bool
	showResults      bool
}

// NewTestModel creates a new test model
func NewTestModel(promptName string, logger *zap.Logger) *TestModel {
	// Initialize with sample providers
	providers := []ModelProvider{
		{Name: "OpenAI", Model: "gpt-4", Enabled: true, APIKey: "", Endpoint: ""},
		{Name: "OpenAI", Model: "gpt-3.5-turbo", Enabled: true, APIKey: "", Endpoint: ""},
		{Name: "Anthropic", Model: "claude-3-sonnet", Enabled: false, APIKey: "", Endpoint: ""},
		{Name: "Google", Model: "gemini-pro", Enabled: false, APIKey: "", Endpoint: ""},
	}
	
	// Sample variables
	variables := map[string]string{
		"task":    "analyze",
		"input":   "code",
		"content": "function hello() { console.log('Hello, World!'); }",
	}
	
	return &TestModel{
		promptName:       promptName,
		logger:           logger,
		providers:        providers,
		variables:        variables,
		results:          []TestResult{},
		selectedProvider: 0,
		testing:          false,
		showResults:      false,
	}
}

// Init initializes the test model
func (m *TestModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the test model
func (m *TestModel) Update(msg tea.Msg) (*TestModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedProvider > 0 {
				m.selectedProvider--
			}
		case "down", "j":
			if m.selectedProvider < len(m.providers)-1 {
				m.selectedProvider++
			}
		case "space":
			// Toggle provider enabled/disabled
			if m.selectedProvider < len(m.providers) {
				m.providers[m.selectedProvider].Enabled = !m.providers[m.selectedProvider].Enabled
			}
		case "enter", "t":
			// Run tests
			if !m.testing {
				return m, m.runTests()
			}
		case "r":
			// Show/hide results
			m.showResults = !m.showResults
		case "c":
			// Clear results
			m.results = []TestResult{}
			m.showResults = false
		}
	case TestCompleteMsg:
		m.testing = false
		m.results = append(m.results, msg.Result)
		m.showResults = true
	}
	
	return m, nil
}

// View renders the test interface
func (m *TestModel) View() string {
	if m.showResults {
		return m.renderResults()
	}
	
	// Create the main layout
	leftPanel := m.renderProviderList()
	rightPanel := m.renderTestConfig()
	
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		rightPanel,
	)
}

// renderProviderList renders the list of model providers
func (m *TestModel) renderProviderList() string {
	listStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1).
		Width(m.width/2 - 2).
		Height(m.height - 4)
	
	var items []string
	items = append(items, "🧪 Model Providers")
	items = append(items, "")
	
	for i, provider := range m.providers {
		prefix := "  "
		if i == m.selectedProvider {
			prefix = "▶ "
		}
		
		status := "❌"
		if provider.Enabled {
			status = "✅"
		}
		
		item := fmt.Sprintf("%s%s %s - %s", prefix, status, provider.Name, provider.Model)
		items = append(items, item)
	}
	
	items = append(items, "")
	if m.testing {
		items = append(items, "🔄 Testing in progress...")
	} else {
		items = append(items, "Commands:")
		items = append(items, "  space - Toggle provider")
		items = append(items, "  t/enter - Run tests")
		items = append(items, "  r - Show results")
		items = append(items, "  c - Clear results")
	}
	
	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	return listStyle.Render(content)
}

// renderTestConfig renders the test configuration panel
func (m *TestModel) renderTestConfig() string {
	configStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1).
		Width(m.width/2 - 2).
		Height(m.height - 4)
	
	var content []string
	content = append(content, "⚙️ Test Configuration")
	content = append(content, "")
	content = append(content, "Variables:")
	for key, value := range m.variables {
		content = append(content, fmt.Sprintf("  %s: %s", key, value))
	}
	
	content = append(content, "")
	content = append(content, "Enabled Providers:")
	enabledCount := 0
	for _, provider := range m.providers {
		if provider.Enabled {
			content = append(content, fmt.Sprintf("  ✅ %s - %s", provider.Name, provider.Model))
			enabledCount++
		}
	}
	
	if enabledCount == 0 {
		content = append(content, "  No providers enabled")
	}
	
	content = append(content, "")
	content = append(content, fmt.Sprintf("Test Results: %d", len(m.results)))
	
	if len(m.results) > 0 {
		content = append(content, "")
		content = append(content, "Latest Results:")
		for i, result := range m.results {
			if i >= 3 { // Show only last 3 results
				break
			}
			status := "✅"
			if result.Error != nil {
				status = "❌"
			}
			content = append(content, fmt.Sprintf("  %s %s: %dms", status, result.Model, result.Latency.Milliseconds()))
		}
	}
	
	configContent := lipgloss.JoinVertical(lipgloss.Left, content...)
	return configStyle.Render(configContent)
}

// renderResults renders the test results view
func (m *TestModel) renderResults() string {
	resultsStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1).
		Width(m.width - 4).
		Height(m.height - 4)
	
	var content []string
	content = append(content, "📊 Test Results")
	content = append(content, "")
	
	if len(m.results) == 0 {
		content = append(content, "No test results available.")
		content = append(content, "Press 't' to run tests, 'r' to go back.")
	} else {
		for i, result := range m.results {
			content = append(content, fmt.Sprintf("--- Result %d ---", i+1))
			content = append(content, fmt.Sprintf("Provider: %s", result.Provider))
			content = append(content, fmt.Sprintf("Model: %s", result.Model))
			content = append(content, fmt.Sprintf("Latency: %dms", result.Latency.Milliseconds()))
			content = append(content, fmt.Sprintf("Tokens: %d", result.Tokens))
			content = append(content, fmt.Sprintf("Cost: $%.4f", result.Cost))
			
			if result.Error != nil {
				content = append(content, fmt.Sprintf("Error: %s", result.Error.Error()))
			} else {
				response := result.Response
				if len(response) > 100 {
					response = response[:100] + "..."
				}
				content = append(content, fmt.Sprintf("Response: %s", response))
			}
			content = append(content, "")
		}
		
		content = append(content, "Press 'r' to go back, 'c' to clear results.")
	}
	
	resultsContent := lipgloss.JoinVertical(lipgloss.Left, content...)
	return resultsStyle.Render(resultsContent)
}

// runTests initiates testing across enabled providers
func (m *TestModel) runTests() tea.Cmd {
	m.testing = true
	m.logger.Info("Starting prompt tests", zap.String("prompt", m.promptName))
	
	// Simulate test execution
	return func() tea.Msg {
		// Simulate API call delay
		time.Sleep(2 * time.Second)
		
		// Create mock result
		result := TestResult{
			Provider:  "OpenAI",
			Model:     "gpt-4",
			Response:  "This is a mock response from the AI model. The code looks good with proper syntax and structure.",
			Latency:   1500 * time.Millisecond,
			Tokens:    45,
			Cost:      0.0023,
			Timestamp: time.Now(),
			Error:     nil,
		}
		
		return TestCompleteMsg{Result: result}
	}
}

// TestCompleteMsg is sent when a test completes
type TestCompleteMsg struct {
	Result TestResult
}

// SetSize updates the test model dimensions
func (m *TestModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetVariables updates the test variables
func (m *TestModel) SetVariables(variables map[string]string) {
	m.variables = variables
}