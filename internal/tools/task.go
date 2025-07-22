// Package tools provides implementations of various tools for the AI agent.
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/tmc/langchaingo/tools"
)

// Task represents a simple task structure.
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// taskTool implements the tools.Tool interface for task management.
type taskTool struct {
	name        string
	description string
	tasks       map[string]*Task
	mutex       sync.RWMutex
	nextID      int
}

// Name returns the tool's name.
func (t *taskTool) Name() string {
	return "createTask"
}

// Description returns the tool's description with input format details.
func (t *taskTool) Description() string {
	return `Create a new task with title and description.
Input: JSON object with required "title" field (string) and optional "description" field (string).
Example: {"title": "Review code", "description": "Review the new feature implementation"}
Returns: JSON object with the created task details including ID and status.`
}

// Call executes the task creation tool with JSON string input.
func (t *taskTool) Call(ctx context.Context, input string) (string, error) {
	var args struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.Unmarshal([]byte(input), &args); err != nil {
		return "", fmt.Errorf("invalid input format, expected JSON with 'title' field: %w", err)
	}

	if args.Title == "" {
		return "", fmt.Errorf("title parameter is required")
	}

	return t.createTaskImpl(ctx, args.Title, args.Description)
}

// createTaskImpl is the implementation function for task creation.
func (t *taskTool) createTaskImpl(ctx context.Context, title, description string) (string, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Generate next ID
	t.nextID++
	taskID := fmt.Sprintf("task_%d", t.nextID)

	// Create new task
	task := &Task{
		ID:          taskID,
		Title:       title,
		Description: description,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	// Store task
	t.tasks[taskID] = task

	log.Info("Task created",
		zap.String("id", taskID),
		zap.String("title", title),
		zap.String("description", description))

	// Return task as JSON
	result, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("failed to marshal task: %w", err)
	}

	return string(result), nil
}

// NewTaskTool creates a new task management tool.
func NewTaskTool() tools.Tool {
	return &taskTool{
		tasks:  make(map[string]*Task),
		nextID: 0,
	}
}