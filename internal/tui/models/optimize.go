// Package models provides the optimize model for AI-powered prompt optimization.
package models

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
)

// OptimizationStep represents a step in the optimization process
type OptimizationStep struct {
	Iteration   int
	Prompt      string
	Suggestion  string
	Improvement float64
	Timestamp   time.Time
}

// OptimizeModel handles AI-powered prompt optimization
type OptimizeModel struct {
	promptName string
	logger     *zap.Logger
	
	// Optimization state
	desiredOutcome   string
	currentPrompt    string
	optimizationSteps []OptimizationStep
	maxIterations    int
	currentIteration int
	
	// UI dimensions
	width  int
	height int
	
	// State
	optimizing     bool
	showHistory    bool
	editingOutcome bool
}

// NewOptimizeModel creates a new optimize model
func NewOptimizeModel(promptName string, logger *zap.Logger) *OptimizeModel {
	return &OptimizeModel{
		promptName:       promptName,
		logger:           logger,
		desiredOutcome:   "Generate clear, helpful, and accurate responses",
		currentPrompt:    "",
		optimizationSteps: []OptimizationStep{},
		maxIterations:    5,
		currentIteration: 0,
		optimizing:       false,
		showHistory:      false,
		editingOutcome:   false,
	}
}

// Init initializes the optimize model
func (m *OptimizeModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the optimize model
func (m *OptimizeModel) Update(msg tea.Msg) (*OptimizeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "o":
			// Start optimization
			if !m.optimizing && !m.editingOutcome {
				return m, m.startOptimization()
			}
		case "s":
			// Stop optimization
			if m.optimizing {
				m.optimizing = false
				m.logger.Info("Optimization stopped by user")
			}
		case "h":
			// Toggle history view
			m.showHistory = !m.showHistory
		case "e":
			// Edit desired outcome
			m.editingOutcome = !m.editingOutcome
		case "r":
			// Reset optimization
			m.resetOptimization()
		case "enter":
			if m.editingOutcome {
				m.editingOutcome = false
			}
		}
	case OptimizationStepMsg:
		m.optimizationSteps = append(m.optimizationSteps, msg.Step)
		m.currentIteration++
		
		if m.currentIteration >= m.maxIterations || msg.Step.Improvement < 0.05 {
			m.optimizing = false
			m.logger.Info("Optimization completed", 
				zap.Int("iterations", m.currentIteration),
				zap.Float64("final_improvement", msg.Step.Improvement))
		} else {
			// Continue optimization
			return m, m.continueOptimization()
		}
	}
	
	return m, nil
}

// View renders the optimization interface
func (m *OptimizeModel) View() string {
	if m.showHistory {
		return m.renderHistory()
	}
	
	// Create the main layout
	topPanel := m.renderOptimizationConfig()
	bottomPanel := m.renderCurrentStep()
	
	return lipgloss.JoinVertical(
		lipgloss.Left,
		topPanel,
		bottomPanel,
	)
}

// renderOptimizationConfig renders the optimization configuration
func (m *OptimizeModel) renderOptimizationConfig() string {
	configStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1).
		Width(m.width - 4).
		Height(m.height/2 - 2)
	
	var content []string
	content = append(content, "🚀 AI Prompt Optimization")
	content = append(content, "")
	
	if m.editingOutcome {
		content = append(content, "✏️ Editing Desired Outcome (press enter to save):")
		content = append(content, fmt.Sprintf("▶ %s", m.desiredOutcome))
	} else {
		content = append(content, "🎯 Desired Outcome:")
		content = append(content, fmt.Sprintf("  %s", m.desiredOutcome))
	}
	
	content = append(content, "")
	content = append(content, "📊 Configuration:")
	content = append(content, fmt.Sprintf("  Max Iterations: %d", m.maxIterations))
	content = append(content, fmt.Sprintf("  Current Iteration: %d", m.currentIteration))
	content = append(content, fmt.Sprintf("  Steps Completed: %d", len(m.optimizationSteps)))
	
	content = append(content, "")
	if m.optimizing {
		content = append(content, "🔄 Optimization in progress...")
		content = append(content, "  s - Stop optimization")
	} else {
		content = append(content, "Commands:")
		content = append(content, "  o - Start optimization")
		content = append(content, "  e - Edit desired outcome")
		content = append(content, "  h - Show history")
		content = append(content, "  r - Reset optimization")
	}
	
	configContent := lipgloss.JoinVertical(lipgloss.Left, content...)
	return configStyle.Render(configContent)
}

// renderCurrentStep renders the current optimization step
func (m *OptimizeModel) renderCurrentStep() string {
	stepStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1).
		Width(m.width - 4).
		Height(m.height/2 - 2)
	
	var content []string
	
	if len(m.optimizationSteps) == 0 {
		content = append(content, "📝 Current Step")
		content = append(content, "")
		content = append(content, "No optimization steps yet.")
		content = append(content, "Press 'o' to start optimization.")
	} else {
		lastStep := m.optimizationSteps[len(m.optimizationSteps)-1]
		content = append(content, fmt.Sprintf("📝 Step %d", lastStep.Iteration))
		content = append(content, "")
		content = append(content, "💡 AI Suggestion:")
		content = append(content, fmt.Sprintf("  %s", lastStep.Suggestion))
		content = append(content, "")
		content = append(content, "📈 Improvement Score:")
		content = append(content, fmt.Sprintf("  %.2f%% improvement", lastStep.Improvement*100))
		content = append(content, "")
		content = append(content, "⏰ Timestamp:")
		content = append(content, fmt.Sprintf("  %s", lastStep.Timestamp.Format("15:04:05")))
		
		if m.optimizing {
			content = append(content, "")
			content = append(content, "🔄 Analyzing and generating next step...")
		}
	}
	
	stepContent := lipgloss.JoinVertical(lipgloss.Left, content...)
	return stepStyle.Render(stepContent)
}

// renderHistory renders the optimization history
func (m *OptimizeModel) renderHistory() string {
	historyStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1).
		Width(m.width - 4).
		Height(m.height - 4)
	
	var content []string
	content = append(content, "📚 Optimization History")
	content = append(content, "")
	
	if len(m.optimizationSteps) == 0 {
		content = append(content, "No optimization history available.")
	} else {
		for i, step := range m.optimizationSteps {
			content = append(content, fmt.Sprintf("--- Step %d ---", step.Iteration))
			content = append(content, fmt.Sprintf("Time: %s", step.Timestamp.Format("15:04:05")))
			content = append(content, fmt.Sprintf("Improvement: %.2f%%", step.Improvement*100))
			content = append(content, fmt.Sprintf("Suggestion: %s", step.Suggestion))
			
			if i < len(m.optimizationSteps)-1 {
				content = append(content, "")
			}
		}
	}
	
	content = append(content, "")
	content = append(content, "Press 'h' to go back.")
	
	historyContent := lipgloss.JoinVertical(lipgloss.Left, content...)
	return historyStyle.Render(historyContent)
}

// startOptimization begins the optimization process
func (m *OptimizeModel) startOptimization() tea.Cmd {
	m.optimizing = true
	m.currentIteration = 0
	m.logger.Info("Starting prompt optimization", 
		zap.String("prompt", m.promptName),
		zap.String("desired_outcome", m.desiredOutcome))
	
	return m.continueOptimization()
}

// continueOptimization continues the optimization process
func (m *OptimizeModel) continueOptimization() tea.Cmd {
	return func() tea.Msg {
		// Simulate AI analysis and optimization
		time.Sleep(3 * time.Second)
		
		// Generate mock optimization step
		suggestions := []string{
			"Add more specific context about the task requirements",
			"Include examples to clarify the expected output format",
			"Specify constraints and limitations more clearly",
			"Add role-based instructions for better context",
			"Include error handling instructions",
		}
		
		iteration := m.currentIteration + 1
		suggestion := suggestions[iteration%len(suggestions)]
		improvement := 0.15 - float64(iteration)*0.03 // Diminishing returns
		
		if improvement < 0 {
			improvement = 0.02
		}
		
		step := OptimizationStep{
			Iteration:   iteration,
			Prompt:      fmt.Sprintf("Optimized prompt v%d", iteration),
			Suggestion:  suggestion,
			Improvement: improvement,
			Timestamp:   time.Now(),
		}
		
		return OptimizationStepMsg{Step: step}
	}
}

// resetOptimization resets the optimization state
func (m *OptimizeModel) resetOptimization() {
	m.optimizationSteps = []OptimizationStep{}
	m.currentIteration = 0
	m.optimizing = false
	m.logger.Info("Optimization reset")
}

// OptimizationStepMsg is sent when an optimization step completes
type OptimizationStepMsg struct {
	Step OptimizationStep
}

// SetSize updates the optimize model dimensions
func (m *OptimizeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetPrompt updates the current prompt for optimization
func (m *OptimizeModel) SetPrompt(prompt string) {
	m.currentPrompt = prompt
}

// SetDesiredOutcome updates the desired outcome
func (m *OptimizeModel) SetDesiredOutcome(outcome string) {
	m.desiredOutcome = outcome
}