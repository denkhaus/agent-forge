package decorators

import (
	"context"
	"errors"
	"testing"

	"github.com/tmc/langchaingo/tools"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// Mocks from test_mocks.go
// mockToolProvider
// mockAgent
// mockToolImpl

func TestNewLoggingToolProviderDecorator(t *testing.T) {
	mockInner := &mockToolProvider{}
	log := zap.NewNop()

	decorator := NewLoggingToolProviderDecorator(mockInner, log)
	loggingDecorator, ok := decorator.(*LoggingToolProviderDecorator)
	if !ok {
		t.Fatal("Expected decorator to be LoggingToolProviderDecorator")
	}

	if loggingDecorator.inner != mockInner {
		t.Error("Inner tool provider not set correctly")
	}
	if loggingDecorator.log != log {
		t.Error("Logger not set correctly")
	}
}

func TestLoggingToolProviderDecorator_ExecuteTool(t *testing.T) {
	tests := []struct {
		name             string
		executeToolErr   error
		expectedLogMsg   string
		expectedLogLevel zapcore.Level
	}{
		{
			name:             "Successful execution",
			executeToolErr:   nil,
			expectedLogMsg:   "Tool execution completed",
			expectedLogLevel: zapcore.InfoLevel,
		},
		{
			name:             "Failed execution",
			executeToolErr:   errors.New("simulated error"),
			expectedLogMsg:   "Tool execution failed",
			expectedLogLevel: zapcore.ErrorLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zapcore.DebugLevel) // Capture all levels
			log := zap.New(core)

			callCount := 0
			mockInner := &mockToolProvider{
				executeToolFunc: func(_ context.Context, _ string, _ string) (string, error) {
					callCount++
					return "mock_result", tt.executeToolErr
				},
			}
			decorator := NewLoggingToolProviderDecorator(mockInner, log)

			ctx := context.Background()
			toolName := "test_tool"
			input := `{"param": "value"}`

			_, err := decorator.ExecuteTool(ctx, toolName, input)
			if !errors.Is(err, tt.executeToolErr) {
				t.Errorf("Expected error %v, got %v", tt.executeToolErr, err)
			}
			if callCount != 1 {
				t.Error("Expected inner tool to be called once")
			}

			// Verify logs
			if recorded.Len() != 2 { // Start and end log
				t.Fatalf("Expected 2 log entries, got %d", recorded.Len())
			}

			startLog := recorded.All()[0]
			if startLog.Message != "Tool execution started" {
				t.Error("Expected 'Tool execution started' log")
			}
			if startLog.Context[0].Key != "tool_name" || startLog.Context[0].String != toolName {
				t.Errorf("Expected tool_name %q, got %q", toolName, startLog.Context[0].String)
			}

			endLog := recorded.All()[1]
			if endLog.Message != tt.expectedLogMsg {
				t.Errorf("Expected log message %q, got %q", tt.expectedLogMsg, endLog.Message)
			}
			if endLog.Level != tt.expectedLogLevel {
				t.Errorf("Expected log level %s, got %s", tt.expectedLogLevel, endLog.Level)
			}
		})
	}
}

func TestLoggingToolProviderDecorator_GetTools(t *testing.T) {
	core, recorded := observer.New(zapcore.DebugLevel)
	log := zap.New(core)
	mockInner := &mockToolProvider{
		tools: []tools.Tool{&mockToolImpl{name: "test_tool", description: "Test tool"}},
	}
	decorator := NewLoggingToolProviderDecorator(mockInner, log)

	tools := decorator.GetTools()
	if len(tools) != 1 { // Check length, not content
		t.Errorf("Expected 1 tool, got %d", len(tools))
	}
	if recorded.Len() != 1 || recorded.All()[0].Message != "Tools retrieved" {
		t.Error("Expected 'Tools retrieved' log message")
	}
	if recorded.All()[0].Context[0].Key != "count" || recorded.All()[0].Context[0].Integer != 1 {
		t.Error("Expected 'count' field with value 1")
	}
}

func TestLoggingToolProviderDecorator_GetToolsForAgent(t *testing.T) {
	core, recorded := observer.New(zapcore.DebugLevel)
	log := zap.New(core)
	mockInner := &mockToolProvider{
		toolsForAgent: []tools.Tool{&mockToolImpl{name: "test_tool", description: "Test tool"}},
	}
	decorator := NewLoggingToolProviderDecorator(mockInner, log)
	mockAgent := &mockAgent{}

	tools, _ := decorator.GetToolsForAgent(mockAgent)
	if len(tools) != 1 { // Check length, not content
		t.Errorf("Expected 1 tool, got %d", len(tools))
	}
	if recorded.Len() != 1 || recorded.All()[0].Message != "Retrieved tools for agent" {
		t.Error("Expected 'Retrieved tools for agent' log message")
	}
	if recorded.All()[0].Context[0].Key != "agent" || recorded.All()[0].Context[0].String != mockAgent.GetName() {
		t.Errorf(
			"Expected 'agent' field with correct name %q, got %q",
			mockAgent.GetName(),
			recorded.All()[0].Context[0].String,
		)
	}
	if recorded.All()[0].Context[1].Key != "count" || recorded.All()[0].Context[1].Integer != 1 {
		t.Error("Expected 'count' field with value 1")
	}
}

func TestLoggingToolProviderDecorator_RegisterTool(t *testing.T) {
	tests := []struct {
		name             string
		registerToolErr  error
		expectedLogMsg   string
		expectedLogLevel zapcore.Level
	}{
		{
			name:             "Successful registration",
			registerToolErr:  nil,
			expectedLogMsg:   "Tool registered successfully",
			expectedLogLevel: zapcore.InfoLevel,
		},
		{
			name:             "Failed registration",
			registerToolErr:  errors.New("registration failed"),
			expectedLogMsg:   "Tool registration failed",
			expectedLogLevel: zapcore.ErrorLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zapcore.DebugLevel) // Capture all levels
			log := zap.New(core)

			mockInner := &mockToolProvider{
				registerToolErr: tt.registerToolErr,
			}
			decorator := NewLoggingToolProviderDecorator(mockInner, log)

			mockTool := &mockToolImpl{
				name:        "test_tool_register",
				description: "A tool to be registered",
			}

			err := decorator.RegisterTool(mockTool)
			if !errors.Is(err, tt.registerToolErr) {
				t.Errorf("Expected error %v, got %v", tt.registerToolErr, err)
			}

			// Verify logs
			if recorded.Len() != 2 { // Start and end log
				t.Fatalf("Expected 2 log entries, got %d", recorded.Len())
			}

			startLog := recorded.All()[0]
			if startLog.Message != "Registering tool" {
				t.Error("Expected 'Registering tool' log")
			}
			if startLog.Context[0].Key != "name" || startLog.Context[0].String != mockTool.Name() {
				t.Errorf("Expected tool name %q, got %q", mockTool.Name(), startLog.Context[0].String)
			}

			endLog := recorded.All()[1]
			if endLog.Message != tt.expectedLogMsg {
				t.Errorf("Expected log message %q, got %q", tt.expectedLogMsg, endLog.Message)
			}
			if endLog.Level != tt.expectedLogLevel {
				t.Errorf("Expected log level %s, got %s", tt.expectedLogLevel, endLog.Level)
			}
		})
	}
}
