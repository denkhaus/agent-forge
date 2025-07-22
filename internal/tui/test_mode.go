// Package tui provides test mode functionality for non-blocking TUI testing.
package tui

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/denkhaus/agentforge/internal/tui/models"
	"go.uber.org/zap"
)

// TestMode represents different testing modes for the TUI
type TestMode int

const (
	InteractiveMode TestMode = iota
	AutomatedMode
	HeadlessMode
	DemoMode
)

// TestConfig configures TUI testing behavior
type TestConfig struct {
	Mode           TestMode
	Duration       time.Duration
	AutoActions    []string
	OutputCapture  bool
	ScreenshotPath string
}

// TestManager handles TUI testing operations
type TestManager struct {
	logger *zap.Logger
	config TestConfig
}

// NewTestManager creates a new test manager
func NewTestManager(logger *zap.Logger, config TestConfig) *TestManager {
	return &TestManager{
		logger: logger,
		config: config,
	}
}

// RunWorkbenchTest runs the workbench in test mode
func (tm *TestManager) RunWorkbenchTest(promptName string) (*TestResult, error) {
	tm.logger.Info("Starting workbench test", 
		zap.String("prompt", promptName),
		zap.Int("mode", int(tm.config.Mode)))

	switch tm.config.Mode {
	case InteractiveMode:
		return tm.runInteractiveTest(promptName)
	case AutomatedMode:
		return tm.runAutomatedTest(promptName)
	case HeadlessMode:
		return tm.runHeadlessTest(promptName)
	case DemoMode:
		return tm.runDemoTest(promptName)
	default:
		return nil, fmt.Errorf("unknown test mode: %d", tm.config.Mode)
	}
}

// runInteractiveTest runs the workbench with user interaction but captures output
func (tm *TestManager) runInteractiveTest(promptName string) (*TestResult, error) {
	result := &TestResult{
		PromptName: promptName,
		StartTime:  time.Now(),
		Mode:       InteractiveMode,
	}

	// Create workbench model
	workbench := models.NewWorkbenchModel(promptName, tm.logger)

	// Create program with test options
	program := tea.NewProgram(workbench, 
		tea.WithAltScreen(),
		tea.WithMouseCellMotion())

	// Run in a goroutine to allow for timeout
	done := make(chan error, 1)
	go func() {
		_, err := program.Run()
		done <- err
	}()

	// Wait for completion or timeout
	select {
	case err := <-done:
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		result.Success = err == nil
		if err != nil {
			result.Error = err.Error()
		}
		return result, err
	case <-time.After(tm.config.Duration):
		program.Quit()
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		result.Success = false
		result.Error = "test timeout"
		return result, fmt.Errorf("test timeout after %v", tm.config.Duration)
	}
}

// runAutomatedTest runs the workbench with automated actions
func (tm *TestManager) runAutomatedTest(promptName string) (*TestResult, error) {
	result := &TestResult{
		PromptName: promptName,
		StartTime:  time.Now(),
		Mode:       AutomatedMode,
		Actions:    tm.config.AutoActions,
	}

	// Create workbench model
	workbench := models.NewWorkbenchModel(promptName, tm.logger)

	// Create program
	program := tea.NewProgram(workbench)

	// Send automated actions
	go func() {
		time.Sleep(100 * time.Millisecond) // Let UI initialize
		
		for _, action := range tm.config.AutoActions {
			tm.sendAction(program, action)
			time.Sleep(200 * time.Millisecond) // Delay between actions
		}
		
		// Quit after all actions
		time.Sleep(500 * time.Millisecond)
		program.Quit()
	}()

	// Run program
	_, err := program.Run()
	
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Success = err == nil
	if err != nil {
		result.Error = err.Error()
	}

	return result, err
}

// runHeadlessTest runs the workbench without UI for testing logic
func (tm *TestManager) runHeadlessTest(promptName string) (*TestResult, error) {
	result := &TestResult{
		PromptName: promptName,
		StartTime:  time.Now(),
		Mode:       HeadlessMode,
	}

	tm.logger.Info("Running headless test", zap.String("prompt", promptName))

	// Create workbench model
	workbench := models.NewWorkbenchModel(promptName, tm.logger)

	// Initialize model
	workbench.Init()

	// Test basic functionality without UI
	ctx, cancel := context.WithTimeout(context.Background(), tm.config.Duration)
	defer cancel()

	// Simulate user interactions
	testActions := []tea.Msg{
		tea.WindowSizeMsg{Width: 120, Height: 40},
		tea.KeyMsg{Type: tea.KeyTab},     // Switch to Variables tab
		tea.KeyMsg{Type: tea.KeyTab},     // Switch to Test tab
		tea.KeyMsg{Type: tea.KeyTab},     // Switch to Optimize tab
		tea.KeyMsg{Type: tea.KeyTab},     // Back to Editor tab
	}

	for _, action := range testActions {
		select {
		case <-ctx.Done():
			result.Error = "test timeout"
			result.Success = false
			break
		default:
			_, _ = workbench.Update(action)
		}
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	if result.Error == "" {
		result.Success = true
	}

	return result, nil
}

// runDemoTest runs the workbench in demo mode with automated showcase
func (tm *TestManager) runDemoTest(promptName string) (*TestResult, error) {
	// Demo actions that showcase all features
	demoActions := []string{
		"tab",           // Go to Variables tab
		"wait:2s",       // Wait 2 seconds
		"n",             // Add new variable
		"wait:1s",
		"tab",           // Go to Test tab
		"wait:2s",
		"t",             // Run test
		"wait:3s",
		"tab",           // Go to Optimize tab
		"wait:2s",
		"o",             // Start optimization
		"wait:3s",
		"tab",           // Back to Editor
		"wait:1s",
		"quit",          // Exit
	}

	// Override config for demo
	tm.config.AutoActions = demoActions
	return tm.runAutomatedTest(promptName)
}

// sendAction sends an action to the tea program
func (tm *TestManager) sendAction(program *tea.Program, action string) {
	switch action {
	case "tab":
		program.Send(tea.KeyMsg{Type: tea.KeyTab})
	case "enter":
		program.Send(tea.KeyMsg{Type: tea.KeyEnter})
	case "quit":
		program.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	case "n":
		program.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	case "t":
		program.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	case "o":
		program.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	default:
		if action[:5] == "wait:" {
			duration, _ := time.ParseDuration(action[5:])
			time.Sleep(duration)
		}
	}
}

// TestResult contains the results of a TUI test
type TestResult struct {
	PromptName string
	Mode       TestMode
	StartTime  time.Time
	EndTime    time.Time
	Duration   time.Duration
	Success    bool
	Error      string
	Actions    []string
	Output     string
}

// String returns a string representation of the test result
func (tr *TestResult) String() string {
	status := "✅ PASS"
	if !tr.Success {
		status = "❌ FAIL"
	}

	return fmt.Sprintf("%s %s (%s) - %v", 
		status, tr.PromptName, tr.Mode.String(), tr.Duration)
}

// String returns the string representation of TestMode
func (tm TestMode) String() string {
	switch tm {
	case InteractiveMode:
		return "Interactive"
	case AutomatedMode:
		return "Automated"
	case HeadlessMode:
		return "Headless"
	case DemoMode:
		return "Demo"
	default:
		return "Unknown"
	}
}