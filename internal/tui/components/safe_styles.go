// Package components provides safe TUI styling components.
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// SafeColors provides terminal-safe color constants
var SafeColors = struct {
	Reset     string
	Bold      string
	Dim       string
	Black     string
	Red       string
	Green     string
	Yellow    string
	Blue      string
	Magenta   string
	Cyan      string
	White     string
	BrightRed string
}{
	Reset:     "\033[0m",
	Bold:      "\033[1m",
	Dim:       "\033[2m",
	Black:     "\033[30m",
	Red:       "\033[31m",
	Green:     "\033[32m",
	Yellow:    "\033[33m",
	Blue:      "\033[34m",
	Magenta:   "\033[35m",
	Cyan:      "\033[36m",
	White:     "\033[37m",
	BrightRed: "\033[91m",
}

// SafeHeader creates a safe header with basic styling
func SafeHeader(title string, width int) string {
	// Center the title
	padding := (width - len(title)) / 2
	if padding < 0 {
		padding = 0
	}
	
	// Create header with safe colors
	header := fmt.Sprintf("%s%s%s%s%s",
		SafeColors.Bold,
		SafeColors.Blue,
		strings.Repeat(" ", padding)+title+strings.Repeat(" ", width-len(title)-padding),
		SafeColors.Reset,
		"")
	
	return header
}

// SafeTabs creates safe tab navigation
func SafeTabs(tabs []string, activeIndex int) string {
	var result []string
	
	for i, tab := range tabs {
		if i == activeIndex {
			// Active tab with highlighting
			result = append(result, fmt.Sprintf("%s%s[%s]%s",
				SafeColors.Bold, SafeColors.Yellow, tab, SafeColors.Reset))
		} else {
			// Inactive tab
			result = append(result, fmt.Sprintf("%s %s %s",
				SafeColors.Dim, tab, SafeColors.Reset))
		}
	}
	
	return strings.Join(result, " ")
}

// SafeBox creates a simple box around content
func SafeBox(content string, width, height int, title string) string {
	lines := strings.Split(content, "\n")
	
	// Ensure we have enough lines
	for len(lines) < height-2 {
		lines = append(lines, "")
	}
	
	// Truncate if too many lines
	if len(lines) > height-2 {
		lines = lines[:height-2]
	}
	
	var result []string
	
	// Top border with title
	if title != "" {
		topBorder := fmt.Sprintf("┌─ %s %s┐", title, strings.Repeat("─", width-len(title)-5))
		result = append(result, topBorder)
	} else {
		result = append(result, "┌"+strings.Repeat("─", width-2)+"┐")
	}
	
	// Content lines
	for _, line := range lines {
		// Truncate or pad line to fit
		if len(line) > width-4 {
			line = line[:width-4]
		}
		paddedLine := fmt.Sprintf("│ %-*s │", width-4, line)
		result = append(result, paddedLine)
	}
	
	// Bottom border
	result = append(result, "└"+strings.Repeat("─", width-2)+"┘")
	
	return strings.Join(result, "\n")
}

// SafeStatus creates a safe status line
func SafeStatus(text string, width int) string {
	if len(text) > width-4 {
		text = text[:width-7] + "..."
	}
	
	return fmt.Sprintf("%s%s%s", SafeColors.Dim, text, SafeColors.Reset)
}

// LipglossHeader creates a safe lipgloss header
func LipglossHeader(title string, width int) string {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")). // White
		Background(lipgloss.Color("4")).  // Blue
		Width(width).
		Align(lipgloss.Center).
		Padding(0, 1)
	
	return style.Render(title)
}

// LipglossTab creates safe tab styling
func LipglossTab(text string, active bool) string {
	if active {
		style := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).  // Black
			Background(lipgloss.Color("11")). // Bright Yellow
			Padding(0, 1)
		return style.Render(text)
	} else {
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")). // Gray
			Padding(0, 1)
		return style.Render(text)
	}
}