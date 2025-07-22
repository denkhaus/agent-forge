// Animated Progress Bar Pattern
// Source: https://github.com/charmbracelet/bubbletea/blob/main/examples/progress-animated/main.go
// Adapted for AgentForge testing and optimization progress

package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

// Key Pattern: Custom message types for progress updates
type tickMsg time.Time
type progressCompleteMsg struct{}

type ProgressModel struct {
	progress progress.Model
	label    string
	complete bool
}

func NewProgressModel(label string) ProgressModel {
	return ProgressModel{
		progress: progress.New(progress.WithDefaultGradient()),
		label:    label,
		complete: false,
	}
}

func (m ProgressModel) Init() tea.Cmd {
	return tickCmd()
}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		// Key Pattern: Responsive width with constraints
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		// Key Pattern: Progress completion check
		if m.progress.Percent() >= 1.0 {
			m.complete = true
			return m, func() tea.Msg { return progressCompleteMsg{} }
		}

		// Key Pattern: Incremental progress update with animation
		cmd := m.progress.IncrPercent(0.25)
		return m, tea.Batch(tickCmd(), cmd)

	// Key Pattern: Handle progress animation frames
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case progressCompleteMsg:
		// Handle completion
		return m, tea.Quit

	default:
		return m, nil
	}
}

func (m ProgressModel) View() string {
	pad := strings.Repeat(" ", padding)
	
	status := "In Progress..."
	if m.complete {
		status = "Complete!"
	}
	
	return "\n" +
		pad + m.label + "\n" +
		pad + m.progress.View() + "\n" +
		pad + status + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

// Key Pattern: Timed command for progress updates
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// AgentForge Testing Progress Example
type TestProgressModel struct {
	progress     progress.Model
	currentTest  string
	testsTotal   int
	testsComplete int
	results      []TestResult
}

type TestResult struct {
	Provider string
	Model    string
	Success  bool
	Duration time.Duration
}

func NewTestProgressModel(totalTests int) TestProgressModel {
	return TestProgressModel{
		progress:      progress.New(progress.WithDefaultGradient()),
		currentTest:   "Initializing...",
		testsTotal:    totalTests,
		testsComplete: 0,
		results:       make([]TestResult, 0),
	}
}

// Key Pattern: Custom progress update messages
type testStartMsg struct{ testName string }
type testCompleteMsg struct{ result TestResult }

func (m TestProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case testStartMsg:
		m.currentTest = msg.testName
		return m, nil
		
	case testCompleteMsg:
		m.results = append(m.results, msg.result)
		m.testsComplete++
		
		// Key Pattern: Calculate progress percentage
		percent := float64(m.testsComplete) / float64(m.testsTotal)
		cmd := m.progress.SetPercent(percent)
		
		if m.testsComplete >= m.testsTotal {
			m.currentTest = "All tests complete!"
		}
		
		return m, cmd
		
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}
	
	return m, nil
}

func (m TestProgressModel) View() string {
	pad := strings.Repeat(" ", padding)
	
	// Progress bar with current test info
	progressView := pad + fmt.Sprintf("Testing Progress (%d/%d)", m.testsComplete, m.testsTotal) + "\n" +
		pad + m.progress.View() + "\n" +
		pad + "Current: " + m.currentTest + "\n\n"
	
	// Results summary
	if len(m.results) > 0 {
		successCount := 0
		for _, result := range m.results {
			if result.Success {
				successCount++
			}
		}
		
		progressView += pad + fmt.Sprintf("Results: %d passed, %d failed", 
			successCount, len(m.results)-successCount) + "\n"
	}
	
	return progressView
}

// AgentForge Optimization Progress Example
type OptimizationProgressModel struct {
	progress       progress.Model
	currentStep    string
	iteration      int
	maxIterations  int
	improvements   []float64
	bestScore      float64
}

func NewOptimizationProgressModel(maxIter int) OptimizationProgressModel {
	return OptimizationProgressModel{
		progress:      progress.New(progress.WithDefaultGradient()),
		currentStep:   "Starting optimization...",
		iteration:     0,
		maxIterations: maxIter,
		improvements:  make([]float64, 0),
		bestScore:     0.0,
	}
}

type optimizationStepMsg struct {
	step  string
	score float64
}

func (m OptimizationProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case optimizationStepMsg:
		m.currentStep = msg.step
		m.iteration++
		m.improvements = append(m.improvements, msg.score)
		
		if msg.score > m.bestScore {
			m.bestScore = msg.score
		}
		
		// Key Pattern: Progress based on iterations
		percent := float64(m.iteration) / float64(m.maxIterations)
		cmd := m.progress.SetPercent(percent)
		
		if m.iteration >= m.maxIterations {
			m.currentStep = "Optimization complete!"
		}
		
		return m, cmd
		
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}
	
	return m, nil
}

func (m OptimizationProgressModel) View() string {
	pad := strings.Repeat(" ", padding)
	
	return pad + fmt.Sprintf("AI Optimization (Iteration %d/%d)", m.iteration, m.maxIterations) + "\n" +
		pad + m.progress.View() + "\n" +
		pad + "Current: " + m.currentStep + "\n" +
		pad + fmt.Sprintf("Best Score: %.2f", m.bestScore) + "\n\n"
}

// AgentForge Adaptation Notes:
// 1. Use TestProgressModel for multi-model testing in Test tab
// 2. Use OptimizationProgressModel for AI enhancement in Optimize tab
// 3. Add progress bars for file operations (save/load)
// 4. Use for long-running operations with user feedback
// 5. Combine with our safe styling for terminal compatibility
// 6. Add cancel functionality with Escape key
// 7. Show detailed results after completion

// Example integration commands:
/*
func (m *WorkbenchModel) startTesting() tea.Cmd {
	return func() tea.Msg {
		return testStartMsg{testName: "OpenAI GPT-4"}
	}
}

func (m *WorkbenchModel) startOptimization() tea.Cmd {
	return func() tea.Msg {
		return optimizationStepMsg{
			step: "Analyzing prompt structure...",
			score: 0.75,
		}
	}
}
*/