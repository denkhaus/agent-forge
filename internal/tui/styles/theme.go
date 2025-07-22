// Package styles provides consistent styling for the TUI workbench.
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	// Primary colors
	Primary   = lipgloss.Color("#7C3AED") // Purple
	Secondary = lipgloss.Color("#3B82F6") // Blue
	Success   = lipgloss.Color("#10B981") // Green
	Warning   = lipgloss.Color("#F59E0B") // Yellow
	Error     = lipgloss.Color("#EF4444") // Red
	
	// Neutral colors
	Background = lipgloss.Color("#1F2937") // Dark gray
	Surface    = lipgloss.Color("#374151") // Medium gray
	Border     = lipgloss.Color("#6B7280") // Light gray
	Text       = lipgloss.Color("#F9FAFB") // White
	TextMuted  = lipgloss.Color("#9CA3AF") // Muted gray
)

// Base styles
var (
	// Container styles
	BaseContainer = lipgloss.NewStyle().
			Padding(1).
			Margin(0).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Border)
	
	// Header styles
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Text).
			Background(Primary).
			Padding(0, 1).
			Margin(0)
	
	// Tab styles
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Text).
			Background(Primary).
			Padding(0, 2).
			Margin(0)
	
	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(TextMuted).
				Background(Surface).
				Padding(0, 2).
				Margin(0)
	
	// Content styles
	ContentStyle = lipgloss.NewStyle().
			Padding(1).
			Margin(0).
			Height(20).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Border)
	
	// Status styles
	StatusStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			Padding(0, 1).
			Margin(0)
	
	// Help styles
	HelpStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			Padding(0, 1).
			Margin(0)
)

// GetContainerStyle returns a styled container with specified dimensions
func GetContainerStyle(width, height int) lipgloss.Style {
	return BaseContainer.Copy().
		Width(width - 4).
		Height(height - 4)
}

// GetHeaderStyle returns a styled header with specified width
func GetHeaderStyle(width int) lipgloss.Style {
	return HeaderStyle.Copy().Width(width)
}

// GetTabStyle returns the appropriate tab style based on active state
func GetTabStyle(active bool) lipgloss.Style {
	if active {
		return ActiveTabStyle
	}
	return InactiveTabStyle
}

// GetContentStyle returns a styled content area with specified dimensions
func GetContentStyle(width, height int) lipgloss.Style {
	return ContentStyle.Copy().
		Width(width - 4).
		Height(height - 6) // Reserve space for header and tabs
}

// GetStatusStyle returns a styled status bar with specified width
func GetStatusStyle(width int) lipgloss.Style {
	return StatusStyle.Copy().Width(width)
}

// GetHelpStyle returns a styled help area with specified width
func GetHelpStyle(width int) lipgloss.Style {
	return HelpStyle.Copy().Width(width)
}