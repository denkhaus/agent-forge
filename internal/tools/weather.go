// Package tools provides implementations of various tools for the AI agent.
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/tmc/langchaingo/tools"
)

// weatherTool implements the EnhancedTool interface for weather information.
type weatherTool struct {
	name        string
	description string
}

// Name returns the tool's name.
func (w *weatherTool) Name() string {
	return "getCurrentWeather"
}

// Description returns the tool's description with input format details.
func (w *weatherTool) Description() string {
	return `Get the current weather in a given location.
Input: JSON object with required "location" field (string).
The location should be a city and state, e.g. "San Francisco, CA".
Example: {"location": "San Francisco, CA"}
Returns: Current weather conditions as a descriptive string.`
}

// Call executes the weather tool with JSON string input.
func (w *weatherTool) Call(ctx context.Context, input string) (string, error) {
	var args struct {
		Location string `json:"location"`
	}

	if err := json.Unmarshal([]byte(input), &args); err != nil {
		return "", fmt.Errorf("invalid input format, expected JSON with 'location' field: %w", err)
	}

	if args.Location == "" {
		return "", fmt.Errorf("location parameter is required")
	}

	return getCurrentWeatherImpl(ctx, args.Location)
}

// NewWeatherTool creates a new weather tool with the given configuration.
func NewWeatherTool() tools.Tool {
	return &weatherTool{}
}

// getCurrentWeatherImpl is the implementation function for weather retrieval.
func getCurrentWeatherImpl(ctx context.Context, location string) (string, error) {

	log.Info("Getting weather information", zap.String("location", location))

	var weather string
	// Mock weather data based on location
	switch strings.ToLower(location) {
	case "chicago":
		weather = "64°F and sunny with light winds"
	case "new york":
		weather = "58°F and cloudy with a chance of rain"
	case "san francisco":
		weather = "62°F and foggy with mild winds"
	default:
		weather = fmt.Sprintf("Weather data not available for %s", location)
		log.Warn("Weather data not available for location", zap.String("location", location))

		return weather, nil
	}

	log.Info("Weather information retrieved",
		zap.String("location", location),
		zap.String("weather", weather))

	return weather, nil
}
