// Multi-Component Focus Management Pattern
// Source: https://github.com/charmbracelet/bubbletea/blob/main/examples/split-editors/main.go
// Adapted for AgentForge prompt workbench

package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Key Pattern: Focus-aware styling
var (
	focusedBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238"))

	blurredBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.HiddenBorder())

	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	focusedPlaceholderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("99"))

	placeholderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("238"))
)

// Key Pattern: Comprehensive keymap for multi-component navigation
type keymap struct {
	next, prev, add, remove, quit key.Binding
}

// Key Pattern: Factory function for consistent component creation
func newTextarea() textarea.Model {
	t := textarea.New()
	t.Prompt = ""
	t.Placeholder = "Type something"
	t.ShowLineNumbers = true
	t.Cursor.Style = cursorStyle
	
	// Key Pattern: Different styles for focused/blurred states
	t.FocusedStyle.Placeholder = focusedPlaceholderStyle
	t.BlurredStyle.Placeholder = placeholderStyle
	t.FocusedStyle.Base = focusedBorderStyle
	t.BlurredStyle.Base = blurredBorderStyle
	
	// Key Pattern: Custom key bindings
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	t.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
	
	t.Blur() // Start blurred
	return t
}

type MultiComponentModel struct {
	width  int
	height int
	keymap keymap
	help   help.Model
	inputs []textarea.Model
	focus  int // Key Pattern: Track which component has focus
}

func NewMultiComponentModel() MultiComponentModel {
	m := MultiComponentModel{
		inputs: make([]textarea.Model, 2), // Start with 2 components
		help:   help.New(),
		keymap: keymap{
			next: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "next"),
			),
			prev: key.NewBinding(
				key.WithKeys("shift+tab"),
				key.WithHelp("shift+tab", "prev"),
			),
			add: key.NewBinding(
				key.WithKeys("ctrl+n"),
				key.WithHelp("ctrl+n", "add component"),
			),
			remove: key.NewBinding(
				key.WithKeys("ctrl+w"),
				key.WithHelp("ctrl+w", "remove component"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
		},
	}
	
	// Initialize components
	for i := range m.inputs {
		m.inputs[i] = newTextarea()
	}
	
	// Key Pattern: Focus the first component
	m.inputs[m.focus].Focus()
	m.updateKeybindings()
	return m
}

func (m MultiComponentModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m MultiComponentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			// Key Pattern: Blur all components before quitting
			for i := range m.inputs {
				m.inputs[i].Blur()
			}
			return m, tea.Quit
			
		case key.Matches(msg, m.keymap.next):
			// Key Pattern: Focus management - blur current, focus next
			m.inputs[m.focus].Blur()
			m.focus++
			if m.focus > len(m.inputs)-1 {
				m.focus = 0
			}
			cmd := m.inputs[m.focus].Focus()
			cmds = append(cmds, cmd)
			
		case key.Matches(msg, m.keymap.prev):
			// Key Pattern: Focus management - blur current, focus previous
			m.inputs[m.focus].Blur()
			m.focus--
			if m.focus < 0 {
				m.focus = len(m.inputs) - 1
			}
			cmd := m.inputs[m.focus].Focus()
			cmds = append(cmds, cmd)
			
		case key.Matches(msg, m.keymap.add):
			// Key Pattern: Dynamic component addition
			m.inputs = append(m.inputs, newTextarea())
			
		case key.Matches(msg, m.keymap.remove):
			// Key Pattern: Dynamic component removal with focus adjustment
			if len(m.inputs) > 1 {
				m.inputs = m.inputs[:len(m.inputs)-1]
				if m.focus > len(m.inputs)-1 {
					m.focus = len(m.inputs) - 1
				}
			}
		}
		
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	m.updateKeybindings()
	m.sizeInputs()

	// Key Pattern: Update all components, but only focused one handles input
	for i := range m.inputs {
		newModel, cmd := m.inputs[i].Update(msg)
		m.inputs[i] = newModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// Key Pattern: Dynamic sizing based on number of components
func (m *MultiComponentModel) sizeInputs() {
	for i := range m.inputs {
		m.inputs[i].SetWidth(m.width / len(m.inputs))
		m.inputs[i].SetHeight(m.height - 5) // Reserve space for help
	}
}

// Key Pattern: Dynamic keybinding updates based on state
func (m *MultiComponentModel) updateKeybindings() {
	m.keymap.add.SetEnabled(len(m.inputs) < 6)    // Max 6 components
	m.keymap.remove.SetEnabled(len(m.inputs) > 1) // Min 1 component
}

func (m MultiComponentModel) View() string {
	// Key Pattern: Help view with current keybindings
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.next,
		m.keymap.prev,
		m.keymap.add,
		m.keymap.remove,
		m.keymap.quit,
	})

	// Key Pattern: Horizontal layout of all components
	var views []string
	for i := range m.inputs {
		views = append(views, m.inputs[i].View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, views...) + "\n\n" + help
}

// AgentForge Adaptation Notes:
// 1. Use this pattern for Editor + Variables + Test + Optimize components
// 2. Replace textarea.Model with our custom component interfaces
// 3. Add component-specific keyboard shortcuts
// 4. Implement focus indicators in our safe styling system
// 5. Use for managing focus between prompt editor, variable table, test results, etc.

// Example adaptation for AgentForge:
/*
type WorkbenchComponent interface {
	Focus() tea.Cmd
	Blur()
	Update(tea.Msg) (WorkbenchComponent, tea.Cmd)
	View() string
	SetSize(width, height int)
}

type AgentForgeWorkbench struct {
	components []WorkbenchComponent
	focus      int
	// editor, variables, test, optimize components
}
*/