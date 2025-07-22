package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	// Primary colors
	PrimaryColor   = lipgloss.Color("#7C3AED") // Purple
	SecondaryColor = lipgloss.Color("#10B981") // Green
	AccentColor    = lipgloss.Color("#F59E0B") // Amber
	ErrorColor     = lipgloss.Color("#EF4444") // Red
	WarningColor   = lipgloss.Color("#F97316") // Orange
	InfoColor      = lipgloss.Color("#3B82F6") // Blue
	
	// Neutral colors
	TextColor      = lipgloss.Color("#1F2937") // Dark gray
	MutedColor     = lipgloss.Color("#6B7280") // Medium gray
	BorderColor    = lipgloss.Color("#D1D5DB") // Light gray
	BackgroundColor = lipgloss.Color("#F9FAFB") // Very light gray
)

// Common styles
var (
	// Title styles
	TitleStyle = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Padding(1, 0)
	
	SubtitleStyle = lipgloss.NewStyle().
		Foreground(TextColor).
		Bold(true)
	
	// Text styles
	BodyStyle = lipgloss.NewStyle().
		Foreground(TextColor)
	
	MutedStyle = lipgloss.NewStyle().
		Foreground(MutedColor)
	
	ErrorStyle = lipgloss.NewStyle().
		Foreground(ErrorColor).
		Bold(true)
	
	SuccessStyle = lipgloss.NewStyle().
		Foreground(SecondaryColor).
		Bold(true)
	
	WarningStyle = lipgloss.NewStyle().
		Foreground(WarningColor).
		Bold(true)
	
	InfoStyle = lipgloss.NewStyle().
		Foreground(InfoColor).
		Bold(true)
	
	// Container styles
	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(1, 2)
	
	HighlightBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(1, 2)
	
	// Input styles
	FocusedInputStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(0, 1)
	
	BlurredInputStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(0, 1)
	
	// Button styles
	PrimaryButtonStyle = lipgloss.NewStyle().
		Background(PrimaryColor).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 2).
		Bold(true)
	
	SecondaryButtonStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Foreground(PrimaryColor).
		Padding(0, 2)
	
	// Progress styles
	ProgressBarStyle = lipgloss.NewStyle().
		Background(BackgroundColor).
		Foreground(PrimaryColor)
	
	// Header and footer
	HeaderStyle = lipgloss.NewStyle().
		Background(PrimaryColor).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Bold(true)
	
	FooterStyle = lipgloss.NewStyle().
		Foreground(MutedColor).
		Padding(1, 0)
)

// Helper functions for consistent spacing
func Spacer(height int) string {
	return lipgloss.NewStyle().Height(height).Render("")
}

func Divider(width int) string {
	return lipgloss.NewStyle().
		Width(width).
		Foreground(BorderColor).
		Render(lipgloss.NewStyle().Width(width).Render("─"))
}