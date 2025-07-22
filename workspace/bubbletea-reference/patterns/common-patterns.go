// Common TUI Patterns and Best Practices
// Compiled from Bubble Tea examples and AgentForge experience

package patterns

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Pattern 1: Message-Based Communication
// Use custom message types for component communication

type ComponentFocusMsg struct{ ComponentID string }
type DataUpdateMsg struct{ Data interface{} }
type StatusChangeMsg struct{ Status string }
type ProgressUpdateMsg struct{ Percent float64 }

// Pattern 2: Comprehensive Key Mapping
type UniversalKeyMap struct {
	// Navigation
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Tab   key.Binding
	
	// Actions
	Enter  key.Binding
	Escape key.Binding
	Save   key.Binding
	Delete key.Binding
	Edit   key.Binding
	
	// Application
	Quit key.Binding
	Help key.Binding
	
	// Custom actions
	Test     key.Binding
	Optimize key.Binding
	Export   key.Binding
}

func NewUniversalKeyMap() UniversalKeyMap {
	return UniversalKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d", "delete"),
			key.WithHelp("d", "delete"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Test: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "test"),
		),
		Optimize: key.NewBinding(
			key.WithKeys("ctrl+o"),
			key.WithHelp("ctrl+o", "optimize"),
		),
		Export: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "export"),
		),
	}
}

// Pattern 3: Component Interface
type TUIComponent interface {
	Init() tea.Cmd
	Update(tea.Msg) (TUIComponent, tea.Cmd)
	View() string
	Focus() tea.Cmd
	Blur()
	SetSize(width, height int)
	GetID() string
}

// Pattern 4: State Management
type ComponentState int

const (
	StateIdle ComponentState = iota
	StateLoading
	StateEditing
	StateError
	StateSuccess
)

func (s ComponentState) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StateLoading:
		return "loading"
	case StateEditing:
		return "editing"
	case StateError:
		return "error"
	case StateSuccess:
		return "success"
	default:
		return "unknown"
	}
}

// Pattern 5: Responsive Layout Manager
type LayoutManager struct {
	width  int
	height int
	areas  map[string]Area
}

type Area struct {
	X, Y          int
	Width, Height int
}

func NewLayoutManager(width, height int) *LayoutManager {
	return &LayoutManager{
		width:  width,
		height: height,
		areas:  make(map[string]Area),
	}
}

func (lm *LayoutManager) SetSize(width, height int) {
	lm.width = width
	lm.height = height
	lm.calculateAreas()
}

func (lm *LayoutManager) calculateAreas() {
	// Header area
	lm.areas["header"] = Area{
		X: 0, Y: 0,
		Width: lm.width, Height: 3,
	}
	
	// Tab area
	lm.areas["tabs"] = Area{
		X: 0, Y: 3,
		Width: lm.width, Height: 3,
	}
	
	// Content area
	lm.areas["content"] = Area{
		X: 0, Y: 6,
		Width: lm.width, Height: lm.height - 9,
	}
	
	// Footer area
	lm.areas["footer"] = Area{
		X: 0, Y: lm.height - 3,
		Width: lm.width, Height: 3,
	}
}

func (lm *LayoutManager) GetArea(name string) Area {
	if area, exists := lm.areas[name]; exists {
		return area
	}
	return Area{} // Empty area
}

// Pattern 6: Error Handling
type ErrorHandler struct {
	errors []error
	maxErrors int
}

func NewErrorHandler(maxErrors int) *ErrorHandler {
	return &ErrorHandler{
		errors: make([]error, 0),
		maxErrors: maxErrors,
	}
}

func (eh *ErrorHandler) AddError(err error) {
	eh.errors = append(eh.errors, err)
	if len(eh.errors) > eh.maxErrors {
		eh.errors = eh.errors[1:] // Remove oldest error
	}
}

func (eh *ErrorHandler) GetLatestError() error {
	if len(eh.errors) == 0 {
		return nil
	}
	return eh.errors[len(eh.errors)-1]
}

func (eh *ErrorHandler) HasErrors() bool {
	return len(eh.errors) > 0
}

func (eh *ErrorHandler) ClearErrors() {
	eh.errors = eh.errors[:0]
}

// Pattern 7: Animation and Transitions
type AnimationState struct {
	frame    int
	maxFrame int
	active   bool
}

func NewAnimationState(maxFrame int) *AnimationState {
	return &AnimationState{
		frame:    0,
		maxFrame: maxFrame,
		active:   false,
	}
}

func (as *AnimationState) Start() {
	as.active = true
	as.frame = 0
}

func (as *AnimationState) Stop() {
	as.active = false
	as.frame = 0
}

func (as *AnimationState) NextFrame() bool {
	if !as.active {
		return false
	}
	
	as.frame++
	if as.frame >= as.maxFrame {
		as.frame = 0
		return true // Animation cycle complete
	}
	return false
}

func (as *AnimationState) GetProgress() float64 {
	if as.maxFrame == 0 {
		return 0
	}
	return float64(as.frame) / float64(as.maxFrame)
}

// Pattern 8: Safe Styling Utilities
func SafeRender(style lipgloss.Style, content string, maxWidth int) string {
	// Ensure content doesn't exceed terminal width
	if len(content) > maxWidth {
		content = content[:maxWidth-3] + "..."
	}
	
	// Apply style safely
	return style.Width(maxWidth).Render(content)
}

func SafeJoinHorizontal(parts ...string) string {
	// Join with safety checks for terminal width
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

func SafeJoinVertical(parts ...string) string {
	// Join with safety checks for terminal height
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// Pattern 9: Configuration Management
type ComponentConfig struct {
	ID           string
	Title        string
	Focusable    bool
	Visible      bool
	MinWidth     int
	MinHeight    int
	MaxWidth     int
	MaxHeight    int
	BorderStyle  lipgloss.Border
	BorderColor  lipgloss.Color
	Padding      int
	Margin       int
}

func DefaultComponentConfig(id, title string) ComponentConfig {
	return ComponentConfig{
		ID:          id,
		Title:       title,
		Focusable:   true,
		Visible:     true,
		MinWidth:    20,
		MinHeight:   5,
		MaxWidth:    120,
		MaxHeight:   40,
		BorderStyle: lipgloss.RoundedBorder(),
		BorderColor: lipgloss.Color("240"),
		Padding:     1,
		Margin:      0,
	}
}

// Pattern 10: Performance Optimization
type RenderCache struct {
	cache    map[string]string
	maxSize  int
	enabled  bool
}

func NewRenderCache(maxSize int) *RenderCache {
	return &RenderCache{
		cache:   make(map[string]string),
		maxSize: maxSize,
		enabled: true,
	}
}

func (rc *RenderCache) Get(key string) (string, bool) {
	if !rc.enabled {
		return "", false
	}
	
	value, exists := rc.cache[key]
	return value, exists
}

func (rc *RenderCache) Set(key, value string) {
	if !rc.enabled {
		return
	}
	
	if len(rc.cache) >= rc.maxSize {
		// Simple eviction: clear cache when full
		rc.cache = make(map[string]string)
	}
	
	rc.cache[key] = value
}

func (rc *RenderCache) Clear() {
	rc.cache = make(map[string]string)
}

// AgentForge Integration Examples:
/*
// Example 1: Using Component Interface
type EditorComponent struct {
	id     string
	state  ComponentState
	config ComponentConfig
	// ... other fields
}

func (e *EditorComponent) GetID() string { return e.id }
func (e *EditorComponent) Focus() tea.Cmd { return nil }
func (e *EditorComponent) Blur() {}
// ... implement other interface methods

// Example 2: Using Layout Manager
func (w *Workbench) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w.layout.SetSize(msg.Width, msg.Height)
		
		contentArea := w.layout.GetArea("content")
		w.editor.SetSize(contentArea.Width, contentArea.Height)
	}
	return w, nil
}

// Example 3: Using Error Handler
func (w *Workbench) handleError(err error) {
	w.errorHandler.AddError(err)
	// Show error in status bar or notification
}
*/