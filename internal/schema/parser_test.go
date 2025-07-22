package schema

import (
	"os"
	"testing"
)

func TestParseToolComponent(t *testing.T) {
	parser := NewComponentParser()
	
	// Read the example tool YAML
	content, err := os.ReadFile("../../examples/tool.yaml")
	if err != nil {
		t.Fatalf("Failed to read tool example: %v", err)
	}
	
	// Parse the component
	component, err := parser.ParseComponent(content)
	if err != nil {
		t.Fatalf("Failed to parse tool component: %v", err)
	}
	
	// Verify it's a tool
	tool, ok := component.(*Tool)
	if !ok {
		t.Fatalf("Expected Tool, got %T", component)
	}
	
	// Verify basic properties
	if tool.Metadata.Name != "weather-lookup" {
		t.Errorf("Expected name 'weather-lookup', got '%s'", tool.Metadata.Name)
	}
	
	if tool.Metadata.Version != "1.2.0" {
		t.Errorf("Expected version '1.2.0', got '%s'", tool.Metadata.Version)
	}
	
	if len(tool.Spec.Functions) != 2 {
		t.Errorf("Expected 2 functions, got %d", len(tool.Spec.Functions))
	}
}

func TestParsePromptComponent(t *testing.T) {
	parser := NewComponentParser()
	
	// Read the example prompt YAML
	content, err := os.ReadFile("../../examples/prompt.yaml")
	if err != nil {
		t.Fatalf("Failed to read prompt example: %v", err)
	}
	
	// Parse the component
	component, err := parser.ParseComponent(content)
	if err != nil {
		t.Fatalf("Failed to parse prompt component: %v", err)
	}
	
	// Verify it's a prompt
	prompt, ok := component.(*Prompt)
	if !ok {
		t.Fatalf("Expected Prompt, got %T", component)
	}
	
	// Verify basic properties
	if prompt.Metadata.Name != "code-reviewer" {
		t.Errorf("Expected name 'code-reviewer', got '%s'", prompt.Metadata.Name)
	}
	
	if len(prompt.Spec.Variables) != 3 {
		t.Errorf("Expected 3 variables, got %d", len(prompt.Spec.Variables))
	}
}

func TestParseAgentComponent(t *testing.T) {
	parser := NewComponentParser()
	
	// Read the example agent YAML
	content, err := os.ReadFile("../../examples/agent.yaml")
	if err != nil {
		t.Fatalf("Failed to read agent example: %v", err)
	}
	
	// Parse the component
	component, err := parser.ParseComponent(content)
	if err != nil {
		t.Fatalf("Failed to parse agent component: %v", err)
	}
	
	// Verify it's an agent
	agent, ok := component.(*Agent)
	if !ok {
		t.Fatalf("Expected Agent, got %T", component)
	}
	
	// Verify basic properties
	if agent.Metadata.Name != "customer-support-agent" {
		t.Errorf("Expected name 'customer-support-agent', got '%s'", agent.Metadata.Name)
	}
	
	if len(agent.Spec.Tools) != 4 {
		t.Errorf("Expected 4 tools, got %d", len(agent.Spec.Tools))
	}
}

func TestComponentSerialization(t *testing.T) {
	parser := NewComponentParser()
	
	// Create a simple tool
	tool := NewTool("test-tool", "1.0.0")
	tool.Metadata.Description = "A test tool"
	tool.Metadata.Author = "Test Author"
	tool.Metadata.License = "MIT"
	tool.Spec.Type = ToolTypeMCPServer
	tool.Spec.Runtime = RuntimeGo
	tool.Spec.EntryPoint = "./test-tool"
	tool.Spec.Functions = []ToolFunction{
		{
			Name:        "test_function",
			Description: "A test function",
			Parameters:  []ToolParameter{},
		},
	}
	
	// Serialize to YAML
	yamlData, err := parser.SerializeComponent(tool)
	if err != nil {
		t.Fatalf("Failed to serialize tool: %v", err)
	}
	
	// Parse it back
	component, err := parser.ParseComponent(yamlData)
	if err != nil {
		t.Fatalf("Failed to parse serialized tool: %v", err)
	}
	
	// Verify it's still a tool with correct properties
	parsedTool, ok := component.(*Tool)
	if !ok {
		t.Fatalf("Expected Tool, got %T", component)
	}
	
	if parsedTool.Metadata.Name != "test-tool" {
		t.Errorf("Expected name 'test-tool', got '%s'", parsedTool.Metadata.Name)
	}
}