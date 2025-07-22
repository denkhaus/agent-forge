package tools

import (
	"github.com/tmc/langchaingo/tools"
)

// GetTools returns all available tools.
func GetTools() []tools.Tool {
	weatherTool := NewWeatherTool()
	taskTool := NewTaskTool()

	return []tools.Tool{
		weatherTool,
		taskTool,
	}
}
