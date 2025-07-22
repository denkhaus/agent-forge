// Rich Data Display with Table Component
// Source: https://github.com/charmbracelet/bubbletea/blob/main/examples/table/main.go
// Adapted for AgentForge variable management

package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Key Pattern: Base styling for table container
var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type TableModel struct {
	table table.Model
}

func (m TableModel) Init() tea.Cmd { return nil }

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Key Pattern: Toggle focus with escape
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			// Key Pattern: Action on selected row
			selectedRow := m.table.SelectedRow()
			if len(selectedRow) > 1 {
				return m, tea.Batch(
					tea.Printf("Selected: %s", selectedRow[1]),
				)
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m TableModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

// AgentForge Variable Table Example
func NewVariableTable() TableModel {
	// Key Pattern: Define columns with appropriate widths
	columns := []table.Column{
		{Title: "Name", Width: 15},
		{Title: "Type", Width: 10},
		{Title: "Default", Width: 20},
		{Title: "Required", Width: 8},
		{Title: "Description", Width: 30},
	}

	// Key Pattern: Sample data for variables
	rows := []table.Row{
		{"task", "string", "analyze", "✓", "The task to perform"},
		{"input", "string", "code", "✓", "Type of input to process"},
		{"content", "string", "", "✓", "The content to analyze"},
		{"language", "string", "auto", "", "Programming language"},
		{"format", "string", "markdown", "", "Output format"},
		{"max_length", "number", "1000", "", "Maximum response length"},
		{"temperature", "number", "0.7", "", "AI model temperature"},
		{"model", "string", "gpt-4", "", "AI model to use"},
	}

	// Key Pattern: Create table with configuration
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Key Pattern: Custom styling for professional appearance
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).                              // Make headers bold
		Foreground(lipgloss.Color("15"))         // White text

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).       // Yellow text
		Background(lipgloss.Color("57")).        // Purple background
		Bold(false)

	// Key Pattern: Apply custom styles
	t.SetStyles(s)

	return TableModel{table: t}
}

// AgentForge Test Results Table Example
func NewTestResultsTable() TableModel {
	columns := []table.Column{
		{Title: "Provider", Width: 12},
		{Title: "Model", Width: 15},
		{Title: "Status", Width: 8},
		{Title: "Latency", Width: 10},
		{Title: "Tokens", Width: 8},
		{Title: "Cost", Width: 10},
		{Title: "Score", Width: 8},
	}

	rows := []table.Row{
		{"OpenAI", "gpt-4", "✓", "1.2s", "45", "$0.002", "9.2"},
		{"OpenAI", "gpt-3.5-turbo", "✓", "0.8s", "42", "$0.001", "8.7"},
		{"Anthropic", "claude-3-sonnet", "✓", "1.5s", "48", "$0.003", "9.0"},
		{"Google", "gemini-pro", "✗", "timeout", "-", "-", "-"},
		{"Azure", "gpt-4", "✓", "1.1s", "44", "$0.002", "9.1"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(8),
	)

	// Key Pattern: Status-aware styling
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("39")).  // Cyan for test results
		BorderBottom(true).
		Bold(true)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("0")).         // Black text
		Background(lipgloss.Color("11")).        // Bright yellow background
		Bold(true)

	t.SetStyles(s)
	return TableModel{table: t}
}

// Key Pattern: Dynamic table updates
func (m *TableModel) UpdateRows(newRows []table.Row) {
	m.table.SetRows(newRows)
}

func (m *TableModel) AddRow(row table.Row) {
	currentRows := m.table.Rows()
	currentRows = append(currentRows, row)
	m.table.SetRows(currentRows)
}

func (m *TableModel) GetSelectedRow() table.Row {
	return m.table.SelectedRow()
}

func (m *TableModel) SetFocus(focused bool) {
	if focused {
		m.table.Focus()
	} else {
		m.table.Blur()
	}
}

// AgentForge Adaptation Notes:
// 1. Use for variable management in Variables tab
// 2. Use for test results display in Test tab
// 3. Use for optimization history in Optimize tab
// 4. Implement row editing with Enter key
// 5. Add context menu with right-click or menu key
// 6. Color-code rows based on status (success/error/pending)
// 7. Support sorting by clicking column headers
// 8. Add filtering and search capabilities

// Example integration with AgentForge workbench:
/*
type VariablesComponent struct {
	table TableModel
	editing bool
	editForm *VariableEditForm
}

func (v *VariablesComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if !v.editing {
				// Start editing selected variable
				selectedRow := v.table.GetSelectedRow()
				v.editForm = NewVariableEditForm(selectedRow)
				v.editing = true
			}
		case "n":
			// Add new variable
			v.table.AddRow([]string{"new_var", "string", "", "", "New variable"})
		case "d":
			// Delete selected variable
			// Implementation here
		}
	}
	
	if v.editing {
		// Update edit form
		return v.editForm.Update(msg)
	} else {
		// Update table
		return v.table.Update(msg)
	}
}
*/