package session

import (
	"context"
	"fmt"
	"sync"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
	"go.uber.org/zap"

	"github.com/denkhaus/agentforge/internal/types"
)

// agentSession implements the AgentSession interface for interacting with a specific agent.
// This is a private implementation of types.AgentSession interface.
type agentSession struct {
	config         types.Config
	agent          types.Agent
	llm            llms.Model
	messageHistory []llms.MessageContent
	toolProvider   types.ToolProvider
	agentProvider  types.AgentProvider
	sessionConfig  types.AgentSessionConfig
	llmService     types.LLMService // Added for dynamic LLM reinitialization
	mutex          sync.RWMutex

	// Performance optimizations
	messagePool    sync.Pool   // Reuse message slices
	toolsCache     []llms.Tool // Cache tools for agent
	toolsCacheTime int64       // Cache timestamp
}

// NewAgentSession creates a new agent session for the specified agent with default configuration.
func NewAgentSession(
	config types.Config,
	agent types.Agent,
	llm llms.Model,
	toolProvider types.ToolProvider,
	agentProvider types.AgentProvider,
	llmService types.LLMService,
) (types.AgentSession, error) {
	return NewAgentSessionWithConfig(config, agent, llm, toolProvider, agentProvider, llmService, types.NewAgentSessionConfig())
}

// NewAgentSessionWithConfig creates a new agent session with custom session configuration.
func NewAgentSessionWithConfig(
	config types.Config,
	agent types.Agent,
	llm llms.Model,
	toolProvider types.ToolProvider,
	agentProvider types.AgentProvider,
	llmService types.LLMService,
	sessionConfig types.AgentSessionConfig,
) (types.AgentSession, error) {
	// Pre-allocate message history with reasonable capacity
	messageHistory := make([]llms.MessageContent, 1, 16) // Start with 1, capacity for 16
	messageHistory[0] = llms.TextParts(llms.ChatMessageTypeSystem, agent.GetSystemPrompt())

	session := &agentSession{
		config:         config,
		agent:          agent,
		llm:            llm,
		toolProvider:   toolProvider,
		agentProvider:  agentProvider,
		sessionConfig:  sessionConfig,
		llmService:     llmService,
		messageHistory: messageHistory,
		messagePool: sync.Pool{
			New: func() interface{} {
				// Pre-allocate slice with reasonable capacity
				return make([]llms.MessageContent, 0, 16)
			},
		},
	}

	log.Info("Agent session created",
		zap.String("agent", agent.GetName()))

	return session, nil
}

// GetAgent returns the agent associated with this session.
func (s *agentSession) GetAgent() types.Agent {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.agent
}

// GenerateResponse generates a response using the AI agent.
func (s *agentSession) GenerateResponse(
	ctx context.Context,
	messages []llms.MessageContent,
	agentTools []tools.Tool,
) (*llms.ContentResponse, error) {
	log.Info("Generating response",
		zap.Int("message_count", len(messages)),
		zap.Int("tool_count", len(agentTools)),
		zap.String("agent", s.agent.GetName()),
		zap.String("execution_mode", s.sessionConfig.ExecutionMode.String()))

	switch s.sessionConfig.ExecutionMode {
	case types.ExecutionModeDirect:
		return s.generateDirectResponse(ctx, messages, agentTools)
	case types.ExecutionModeAgent:
		return s.generateAgentResponse(ctx, messages, agentTools)
	default:
		return nil, fmt.Errorf("unsupported execution mode: %s", s.sessionConfig.ExecutionMode.String())
	}
}

// generateDirectResponse uses the current direct LLM approach.
func (s *agentSession) generateDirectResponse(
	ctx context.Context,
	messages []llms.MessageContent,
	agentTools []tools.Tool,
) (*llms.ContentResponse, error) {
	// Convert tools.Tool to llms.Tool for direct LLM calls
	llmTools := s.convertToLLMTools(agentTools)

	resp, err := s.llm.GenerateContent(ctx, messages, llms.WithTools(llmTools))
	if err != nil {
		log.Error("Failed to generate direct response", zap.Error(err))
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	log.Info("Direct response generated", zap.Int("choices", len(resp.Choices)))
	return resp, nil
}

// generateAgentResponse uses OneShotAgent with automatic reasoning loops.
func (s *agentSession) generateAgentResponse(
	ctx context.Context,
	messages []llms.MessageContent,
	agentTools []tools.Tool,
) (*llms.ContentResponse, error) {
	// Create OneShotAgent with tools directly (no conversion needed)
	agent := agents.NewOneShotAgent(s.llm, agentTools, agents.WithMaxIterations(s.sessionConfig.MaxIterations))
	executor := agents.NewExecutor(agent)

	// Extract the user message from the messages (last non-system message)
	userMessage := s.extractUserMessage(messages)
	if userMessage == "" {
		return nil, fmt.Errorf("no user message found in message history")
	}

	log.Info("Running OneShotAgent",
		zap.String("user_message", userMessage),
		zap.Int("max_iterations", s.sessionConfig.MaxIterations))

	// Run the agent
	answer, err := chains.Run(ctx, executor, userMessage)
	if err != nil {
		log.Error("Failed to run OneShotAgent", zap.Error(err))
		return nil, fmt.Errorf("failed to run agent: %w", err)
	}

	// Convert the string response back to ContentResponse format
	response := &llms.ContentResponse{
		Choices: []*llms.ContentChoice{
			{
				Content: answer,
			},
		},
	}

	log.Info("Agent response generated", zap.String("answer", answer))
	return response, nil
}

// Chat starts an interactive chat session.
func (s *agentSession) Chat(ctx context.Context, initialMessage string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Delegate complex chat logic to chat manager
	chatMgr := newChatManager(s, s.toolProvider)
	return chatMgr.processChat(ctx, initialMessage)
}

// addMessageToHistory adds a message to the session's message history.
// This is a helper method for the chat manager to maintain encapsulation.
func (s *agentSession) addMessageToHistory(message llms.MessageContent) {
	// Check capacity and grow efficiently if needed
	if cap(s.messageHistory) == len(s.messageHistory) {
		// Grow by 50% to reduce allocations
		newCap := cap(s.messageHistory) + cap(s.messageHistory)/2
		if newCap < 8 {
			newCap = 8
		}
		newHistory := make([]llms.MessageContent, len(s.messageHistory), newCap)
		copy(newHistory, s.messageHistory)
		s.messageHistory = newHistory
	}
	s.messageHistory = append(s.messageHistory, message)
}

// GetMessageHistory returns the current message history.
func (s *agentSession) GetMessageHistory() []llms.MessageContent {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Use pool for better memory management
	history := s.messagePool.Get().([]llms.MessageContent)
	history = history[:0] // Reset length but keep capacity

	// Append instead of copy for better performance
	history = append(history, s.messageHistory...)

	// Return copy and put slice back in pool
	result := make([]llms.MessageContent, len(history))
	copy(result, history)
	s.messagePool.Put(&history) // Adjusted line length

	return result
}

// ClearMessageHistory clears the current message history.
func (s *agentSession) ClearMessageHistory() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Keep only the current agent's system prompt
	s.messageHistory = []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, s.agent.GetSystemPrompt()),
	}

	log.Info("Message history cleared")
}

// AddMessage adds a message to the history.
func (s *agentSession) AddMessage(message llms.MessageContent) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.addMessageToHistory(message)
	log.Info("Message added to history", zap.String("role", string(message.Role)))
}

// SwitchAgent switches to a different agent.
func (s *agentSession) SwitchAgent(ctx context.Context, agentName string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	newAgent, err := s.agentProvider.GetAgent(agentName)
	if err != nil {
		return fmt.Errorf("failed to get agent %s: %w", agentName, err)
	}

	// Initialize new LLM if the agent's config is different
	if s.agent.GetLLMConfig() != newAgent.GetLLMConfig() {
		log.Info("Agent LLM config change detected, reinitializing LLM",
			zap.String("agent", agentName),
			zap.String("old_provider", s.agent.GetLLMConfig().GetProvider()),
			zap.String("new_provider", newAgent.GetLLMConfig().GetProvider()))

		newLLM, err := s.llmService.InitializeLLM(ctx, s.config, newAgent.GetLLMConfig())
		if err != nil {
			return fmt.Errorf("failed to initialize LLM for agent %s: %w", agentName, err)
		}
		s.llm = newLLM
		log.Info("LLM successfully reinitialized for new agent", zap.String("agent", agentName))
	}

	oldAgentName := s.agent.GetName()
	s.agent = newAgent

	// Clear tools cache when switching agents
	s.toolsCache = nil
	s.toolsCacheTime = 0

	// Update system prompt in message history efficiently
	systemPrompt := llms.TextParts(llms.ChatMessageTypeSystem, newAgent.GetSystemPrompt())
	if len(s.messageHistory) > 0 && s.messageHistory[0].Role == llms.ChatMessageTypeSystem {
		s.messageHistory[0] = systemPrompt
	} else {
		// Efficient prepend: grow slice and shift elements
		s.messageHistory = append(s.messageHistory, llms.MessageContent{})
		copy(s.messageHistory[1:], s.messageHistory[0:])
		s.messageHistory[0] = systemPrompt
	}

	log.Info("Agent switched",
		zap.String("from", oldAgentName),
		zap.String("to", newAgent.GetName()))

	return nil
}

// GetExecutionMode returns the current execution mode.
func (s *agentSession) GetExecutionMode() types.ExecutionMode {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.sessionConfig.ExecutionMode
}

// GetSessionConfig returns the session configuration.
func (s *agentSession) GetSessionConfig() types.AgentSessionConfig {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.sessionConfig
}

// convertToLLMTools converts tools.Tool to llms.Tool for direct LLM calls.
func (s *agentSession) convertToLLMTools(agentTools []tools.Tool) []llms.Tool {
	llmTools := make([]llms.Tool, 0, len(agentTools))

	for _, tool := range agentTools {
		llmTool := llms.Tool{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        tool.Name(),
				Description: tool.Description(),
				// Note: Parameters are not available from tools.Tool interface
				// LLM will rely on the description for input format
				Parameters: nil,
			},
		}
		llmTools = append(llmTools, llmTool)
	}

	return llmTools
}

// extractUserMessage extracts the last user message from the message history.
func (s *agentSession) extractUserMessage(messages []llms.MessageContent) string {
	// Look for the last user message
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == llms.ChatMessageTypeHuman {
			// Extract text content from the message
			for _, part := range messages[i].Parts {
				if textPart, ok := part.(llms.TextContent); ok {
					return textPart.Text
				}
			}
		}
	}
	return ""
}
