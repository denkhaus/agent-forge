// Adaptive Themes and Professional Styling Patterns
// Compiled from multiple Bubble Tea examples
// Adapted for AgentForge prompt workbench

package main

import (
	"github.com/charmbracelet/lipgloss"
)

// Key Pattern: Adaptive colors for light/dark theme support
var (
	// Primary color palette with adaptive support
	primaryColor = lipgloss.AdaptiveColor{
		Light: "#874BFD", // Purple in light mode
		Dark:  "#7D56F4", // Slightly different purple in dark mode
	}
	
	accentColor = lipgloss.AdaptiveColor{
		Light: "#FF6B6B", // Red accent in light mode
		Dark:  "#FF8E8E", // Softer red in dark mode
	}
	
	successColor = lipgloss.AdaptiveColor{
		Light: "#51CF66", // Green in light mode
		Dark:  "#69DB7C", // Brighter green in dark mode
	}
	
	warningColor = lipgloss.AdaptiveColor{
		Light: "#FFD43B", // Yellow in light mode
		Dark:  "#FFE066", // Brighter yellow in dark mode
	}
	
	errorColor = lipgloss.AdaptiveColor{
		Light: "#FF6B6B", // Red in light mode
		Dark:  "#FF8E8E", // Softer red in dark mode
	}
	
	// Neutral colors for backgrounds and borders
	backgroundColor = lipgloss.AdaptiveColor{
		Light: "#FFFFFF", // White in light mode
		Dark:  "#1A1A1A", // Dark gray in dark mode
	}
	
	surfaceColor = lipgloss.AdaptiveColor{
		Light: "#F8F9FA", // Light gray in light mode
		Dark:  "#2D2D2D", // Medium gray in dark mode
	}
	
	borderColor = lipgloss.AdaptiveColor{
		Light: "#E9ECEF", // Light border in light mode
		Dark:  "#404040", // Dark border in dark mode
	}
	
	textColor = lipgloss.AdaptiveColor{
		Light: "#212529", // Dark text in light mode
		Dark:  "#F8F9FA", // Light text in dark mode
	}
	
	mutedTextColor = lipgloss.AdaptiveColor{
		Light: "#6C757D", // Gray text in light mode
		Dark:  "#ADB5BD", // Light gray text in dark mode
	}
)

// Key Pattern: Semantic styling system
var (
	// Base styles
	BaseStyle = lipgloss.NewStyle().
		Foreground(textColor).
		Background(backgroundColor)
	
	// Header styles
	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(primaryColor).
		Padding(0, 1).
		Align(lipgloss.Center)
	
	// Tab styles with adaptive colors
	ActiveTabStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(backgroundColor).
		Background(primaryColor).
		Padding(0, 1)
	
	InactiveTabStyle = lipgloss.NewStyle().
		Foreground(mutedTextColor).
		Background(surfaceColor).
		Padding(0, 1)
	
	// Content area styles
	ContentStyle = lipgloss.NewStyle().
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Background(backgroundColor)
	
	// Focus indicators
	FocusedStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Background(backgroundColor)
	
	BlurredStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Background(surfaceColor)
	
	// Status styles
	SuccessStyle = lipgloss.NewStyle().
		Foreground(successColor).
		Bold(true)
	
	WarningStyle = lipgloss.NewStyle().
		Foreground(warningColor).
		Bold(true)
	
	ErrorStyle = lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true)
	
	// Help and info styles
	HelpStyle = lipgloss.NewStyle().
		Foreground(mutedTextColor).
		Italic(true)
	
	InfoStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true)
)

// Key Pattern: Responsive styling functions
func GetResponsiveHeaderStyle(width int) lipgloss.Style {
	return HeaderStyle.Copy().Width(width)
}

func GetResponsiveContentStyle(width, height int) lipgloss.Style {
	return ContentStyle.Copy().
		Width(width - 4).  // Account for borders
		Height(height - 4) // Account for borders
}

func GetResponsiveTabStyle(active bool, width int) lipgloss.Style {
	if active {
		return ActiveTabStyle.Copy().Width(width)
	}
	return InactiveTabStyle.Copy().Width(width)
}

// Key Pattern: Status-aware styling
func GetStatusStyle(status string) lipgloss.Style {
	switch status {
	case "success", "complete", "passed":
		return SuccessStyle
	case "warning", "pending", "in-progress":
		return WarningStyle
	case "error", "failed", "timeout":
		return ErrorStyle
	default:
		return BaseStyle
	}
}

// Key Pattern: Component-specific styling
type ComponentTheme struct {
	Primary    lipgloss.Style
	Secondary  lipgloss.Style
	Focused    lipgloss.Style
	Blurred    lipgloss.Style
	Header     lipgloss.Style
	Content    lipgloss.Style
	Footer     lipgloss.Style
}

func NewEditorTheme() ComponentTheme {
	return ComponentTheme{
		Primary: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1),
		
		Secondary: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1),
		
		Focused: FocusedStyle,
		Blurred: BlurredStyle,
		
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(borderColor),
		
		Content: ContentStyle,
		
		Footer: lipgloss.NewStyle().
			Foreground(mutedTextColor).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(borderColor),
	}
}

func NewVariablesTheme() ComponentTheme {
	return ComponentTheme{
		Primary: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(1),
		
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor),
		
		// Table-specific styles
		Content: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(borderColor),
	}
}

func NewTestTheme() ComponentTheme {
	return ComponentTheme{
		Primary: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(warningColor).
			Padding(1),
		
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(warningColor),
	}
}

func NewOptimizeTheme() ComponentTheme {
	return ComponentTheme{
		Primary: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(successColor).
			Padding(1),
		
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(successColor),
	}
}

// Key Pattern: Dynamic theme switching
type ThemeManager struct {
	currentTheme string
	themes       map[string]ComponentTheme
}

func NewThemeManager() *ThemeManager {
	return &ThemeManager{
		currentTheme: "default",
		themes: map[string]ComponentTheme{
			"editor":    NewEditorTheme(),
			"variables": NewVariablesTheme(),
			"test":      NewTestTheme(),
			"optimize":  NewOptimizeTheme(),
		},
	}
}

func (tm *ThemeManager) GetTheme(component string) ComponentTheme {
	if theme, exists := tm.themes[component]; exists {
		return theme
	}
	return tm.themes["editor"] // Default fallback
}

// AgentForge Adaptation Notes:
// 1. Use adaptive colors for automatic light/dark theme support
// 2. Apply component-specific themes for visual hierarchy
// 3. Use status-aware styling for test results and optimization
// 4. Implement responsive styling for different terminal sizes
// 5. Combine with our safe styling components for compatibility
// 6. Add theme switching capability for user preferences
// 7. Use semantic styling for consistent appearance

// Example usage in AgentForge workbench:
/*
type AgentForgeWorkbench struct {
	themeManager *ThemeManager
	// ... other fields
}

func (w *AgentForgeWorkbench) renderEditor() string {
	theme := w.themeManager.GetTheme("editor")
	content := w.editor.View()
	return theme.Primary.Render(content)
}

func (w *AgentForgeWorkbench) renderWithStatus(content, status string) string {
	statusStyle := GetStatusStyle(status)
	return statusStyle.Render(content)
}
*/