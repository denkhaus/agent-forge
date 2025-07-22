// Seamless Tab Navigation Pattern
// Source: https://github.com/charmbracelet/bubbletea/blob/main/examples/tabs/main.go
// Adapted for AgentForge prompt workbench

package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Key Pattern: Custom border function for seamless tab connection
func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

// Key Pattern: Adaptive colors for light/dark theme support
var (
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	
	// Tab styles with seamless borders
	inactiveTabStyle = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(highlightColor).
		Padding(0, 1)
	
	activeTabStyle = inactiveTabStyle.Copy().
		Border(activeTabBorder, true)
	
	// Window style that connects seamlessly to tabs
	windowStyle = lipgloss.NewStyle().
		BorderForeground(highlightColor).
		Padding(2, 0).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		UnsetBorderTop() // Key: Remove top border for seamless connection
)

type TabModel struct {
	Tabs       []string
	TabContent []string
	activeTab  int
}

func (m TabModel) Init() tea.Cmd {
	return nil
}

func (m TabModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}
	return m, nil
}

// Key Pattern: Dynamic border adjustment for first/last/active tabs
func (m TabModel) View() string {
	doc := strings.Builder{}
	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		
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
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	// Join tabs horizontally and connect to content window
	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	
	// Content window with dynamic width matching tabs
	contentWidth := lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize()
	doc.WriteString(windowStyle.Width(contentWidth).Render(m.TabContent[m.activeTab]))
	
	return doc.String()
}

// AgentForge Adaptation Notes:
// 1. Use this pattern for Editor, Variables, Test, Optimize tabs
// 2. Replace TabContent with actual component views
// 3. Add keyboard shortcuts: Ctrl+S (save), Ctrl+T (test), Ctrl+O (optimize)
// 4. Integrate with our safe styling components for terminal compatibility

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}