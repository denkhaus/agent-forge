// Package session provides session management for the AI agent.
package session

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
	"go.uber.org/zap"

	"github.com/denkhaus/agentforge/internal/errors"
	"github.com/denkhaus/agentforge/internal/types"
)

// chatManager handles the complex chat orchestration logic.
type chatManager struct {
	session      *agentSession
	toolProvider types.ToolProvider

	// Performance optimizations
	argsPool sync.Pool // Reuse argument maps
}

// newChatManager creates a new chat manager for the given session.
func newChatManager(session *agentSession, toolProvider types.ToolProvider) *chatManager {
	return &chatManager{
		session:      session,
		toolProvider: toolProvider,
		argsPool: sync.Pool{
			New: func() interface{} {
				return make(map[string]any, 8) // Pre-allocate with reasonable capacity
			},
		},
	}
}

// processChat handles the complete chat flow including tool execution.
func (cm *chatManager) processChat(ctx context.Context, initialMessage string) error {
	// Add user message to history
	userMessage := llms.TextParts(llms.ChatMessageTypeHuman, initialMessage)
	cm.session.addMessageToHistory(userMessage)

	// Get tools filtered by agent capabilities (with caching)
	tools, err := cm.getToolsForAgentCached()
	if err != nil {
		return fmt.Errorf("failed to get tools for agent: %w", err)
	}

	// Convert llms.Tool to tools.Tool for session
	agentTools := cm.convertToAgentTools(tools)
	
	// Generate initial response
	resp, err := cm.session.GenerateResponse(ctx, cm.session.messageHistory, agentTools)
	if err != nil {
		return fmt.Errorf("failed to generate initial response: %w", err)
	}

	if len(resp.Choices) == 0 {
		return errors.ErrServiceUnavailable
	}

	choice := resp.Choices[0]

	// Add assistant response to history
	assistantResponse := cm.buildAssistantResponse(choice)
	cm.session.addMessageToHistory(assistantResponse)

	// Execute tool calls if any
	if len(choice.ToolCalls) > 0 {
		if err := cm.executeToolCalls(ctx, choice.ToolCalls); err != nil {
			return fmt.Errorf("tool execution failed: %w", err)
		}

		// Generate final response after tool execution
		return cm.generateFinalResponse(ctx, tools)
	}

	log.Info("Direct response", zap.String("content", choice.Content))
	return nil
}

// buildAssistantResponse creates an assistant message with tool calls.
func (cm *chatManager) buildAssistantResponse(choice *llms.ContentChoice) llms.MessageContent {
	assistantResponse := llms.TextParts(llms.ChatMessageTypeAI, choice.Content)
	for _, tc := range choice.ToolCalls {
		assistantResponse.Parts = append(assistantResponse.Parts, tc)
	}
	return assistantResponse
}

// executeToolCalls executes all tool calls and adds results to history.
func (cm *chatManager) executeToolCalls(ctx context.Context, toolCalls []llms.ToolCall) error { // Adjusted line length
	var err error // Declare error once
	for _, tc := range toolCalls {
		if currentErr := cm.executeToolCall(ctx, tc); currentErr != nil {
			log.Error("Tool call execution failed",
				zap.String("tool", tc.FunctionCall.Name),
				zap.Error(currentErr))
			if err == nil { // Store the first error encountered
				err = currentErr
			}
		}
	}
	return err // Return the first error or nil if all succeeded
}

// executeToolCall executes a single tool call and adds the result to message history.
func (cm *chatManager) executeToolCall(ctx context.Context, tc llms.ToolCall) error {
	// Use pool for argument maps to reduce allocations
	args := cm.argsPool.Get().(map[string]any)
	defer func() {
		// Clear map and return to pool
		for k := range args {
			delete(args, k)
		}
		cm.argsPool.Put(args)
	}()

	if err := json.Unmarshal([]byte(tc.FunctionCall.Arguments), &args); err != nil {
		return fmt.Errorf("failed to unmarshal tool arguments: %w", err)
	}

	// Convert args to JSON string for new ExecuteTool signature
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("failed to marshal tool arguments: %w", err)
	}
	
	result, err := cm.toolProvider.ExecuteTool(ctx, tc.FunctionCall.Name, string(argsJSON))
	if err != nil {
		return fmt.Errorf("tool execution failed: %w", err)
	}

	// Add tool response to message history
	toolResponse := llms.MessageContent{
		Role: llms.ChatMessageTypeTool,
		Parts: []llms.ContentPart{
			llms.ToolCallResponse{
				Name:    tc.FunctionCall.Name,
				Content: result,
			},
		},
	}
	cm.session.addMessageToHistory(toolResponse)

	return nil
}

// generateFinalResponse generates a final response after tool execution.
func (cm *chatManager) generateFinalResponse(ctx context.Context, tools []llms.Tool) error {
	// Convert llms.Tool to tools.Tool for session
	agentTools := cm.convertToAgentTools(tools)
	
	finalResp, err := cm.session.GenerateResponse(ctx, cm.session.messageHistory, agentTools)
	if err != nil {
		return fmt.Errorf("failed to generate final response: %w", err)
	}

	if len(finalResp.Choices) > 0 {
		finalChoice := finalResp.Choices[0]
		finalAssistantResponse := llms.TextParts(llms.ChatMessageTypeAI, finalChoice.Content)
		cm.session.addMessageToHistory(finalAssistantResponse)

		log.Info("Final response", zap.String("content", finalChoice.Content))
	}

	return nil
}

// getToolsForAgentCached returns tools for the agent with caching for performance.
func (cm *chatManager) getToolsForAgentCached() ([]llms.Tool, error) {
	const cacheValiditySeconds = 300 // 5 minutes cache

	now := time.Now().Unix()

	// Check if cache is valid
	if cm.session.toolsCache != nil && (now-cm.session.toolsCacheTime) < cacheValiditySeconds {
		return cm.session.toolsCache, nil
	}

	// Cache miss or expired - fetch fresh tools
	agentTools, err := cm.toolProvider.GetToolsForAgent(cm.session.agent)
	if err != nil {
		return nil, err
	}
	
	// Convert to llms.Tool for caching
	tools := cm.convertToLLMTools(agentTools)

	// Update cache
	cm.session.toolsCache = tools
	cm.session.toolsCacheTime = now

	return tools, nil
}

// convertToAgentTools converts llms.Tool to tools.Tool for session usage.
func (cm *chatManager) convertToAgentTools(llmTools []llms.Tool) []tools.Tool {
	agentTools := make([]tools.Tool, 0, len(llmTools))
	
	for _, llmTool := range llmTools {
		// Create a simple wrapper that implements tools.Tool
		agentTool := &llmToolWrapper{
			name:        llmTool.Function.Name,
			description: llmTool.Function.Description,
		}
		agentTools = append(agentTools, agentTool)
	}
	
	return agentTools
}

// convertToLLMTools converts tools.Tool to llms.Tool for caching.
func (cm *chatManager) convertToLLMTools(agentTools []tools.Tool) []llms.Tool {
	llmTools := make([]llms.Tool, 0, len(agentTools))
	
	for _, agentTool := range agentTools {
		llmTool := llms.Tool{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        agentTool.Name(),
				Description: agentTool.Description(),
				Parameters:  nil, // Schema not available from tools.Tool
			},
		}
		llmTools = append(llmTools, llmTool)
	}
	
	return llmTools
}

// llmToolWrapper wraps llms.Tool data for tools.Tool interface.
type llmToolWrapper struct {
	name        string
	description string
}

func (w *llmToolWrapper) Name() string {
	return w.name
}

func (w *llmToolWrapper) Description() string {
	return w.description
}

func (w *llmToolWrapper) Call(ctx context.Context, input string) (string, error) {
	// This is a placeholder - actual tool execution happens through toolProvider
	return fmt.Sprintf("Tool %s called with input: %s", w.name, input), nil
}
