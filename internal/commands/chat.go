// Package commands provides command-line interface handlers for the MCP Planner application.
package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/samber/do"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/denkhaus/agentforge/internal/errors"
	"github.com/denkhaus/agentforge/internal/shared"
	"github.com/denkhaus/agentforge/internal/types"
)

// chatHandler handles the chat command functionality.
type chatHandler struct {
}

// GetChatCommand returns the chat command configuration.
func GetChatCommand() *cli.Command {
	return &cli.Command{
		Name:    "chat",
		Aliases: []string{"c"},
		Usage:   "Start an interactive chat session with the AI agent",
		Action:  HandleChat(),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "clear",
				Aliases: []string{"cl"},
				Usage:   "Clear message history before starting",
			},
		},
	}
}

// HandleChat creates a new chat command handler.
func HandleChat() cli.ActionFunc {
	return func(c *cli.Context) error {
		ctx := c.Context

		sc, err := shared.StartupCommon(c)
		if err != nil {
			return err
		}

		defer sc.Cleanup()

		h := &chatHandler{}
		// Create session
		session, err := h.createSession(ctx, sc)
		if err != nil {
			return err
		}

		// Clear history if requested
		if c.Bool("clear") {
			session.ClearMessageHistory()
			log.Info("Message history cleared by user request")
			log.Info("Message history cleared.")
		}

		// Start interactive chat
		return h.runInteractiveChat(ctx, session)
	}
}

// createSession creates and configures an agent session.
func (h *chatHandler) createSession(ctx context.Context, startupContext *shared.StartupContext) (types.AgentSession, error) {
	// Get agent provider
	agentProvider := do.MustInvoke[types.AgentProvider](startupContext.Injector)
	cfg := startupContext.Config

	// Create a session with the default agent
	agent, err := agentProvider.GetAgent("planner")
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("default agent 'planner' not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}

	// Get dependencies for session creation
	toolProvider := do.MustInvoke[types.ToolProvider](startupContext.Injector)

	// Create session using the session factory
	sessionFactory := do.MustInvoke[types.SessionFactory](startupContext.Injector)
	agentSession, err := sessionFactory.CreateSession(ctx, types.SessionOptions{
		Config:        cfg,
		Agent:         agent,
		ToolProvider:  toolProvider,
		AgentProvider: agentProvider,
		AgentType:     "planner", // hint for potential optimizations
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create agent session: %w", err)
	}

	return agentSession, nil
}

// printWelcomeMessage displays the initial chat interface information.
func (h *chatHandler) printWelcomeMessage(session types.AgentSession) {
	log.Info("MCP Planner Chat Interface")
	log.Info("Using agent", zap.String("agent_name", session.GetAgent().GetName()))
	log.Info("Type 'exit', 'quit', or 'bye' to end the session.")
	log.Info("Type 'clear' to clear the message history.")
	log.Info("Type 'history' to view the current message history.")
	log.Info("Type 'switch to <agent>' to switch agents.")
	log.Info("Type 'help' for more commands.")
	log.Info("---")
}

// handleSpecialCommands processes special chat commands.
// Returns true if the command was handled, false otherwise.
func (h *chatHandler) handleSpecialCommands(ctx context.Context, input string, session types.AgentSession) (bool, error) {
	lowerInput := strings.ToLower(input)

	switch lowerInput {
	case "exit", "quit", "bye":
		log.Info("Goodbye!")
		os.Exit(0)
		return true, nil
	case "clear":
		session.ClearMessageHistory()
		log.Info("Message history cleared by user command")
		log.Info("Message history cleared.")
		return true, nil
	case "history":
		h.printMessageHistory(session)
		return true, nil
	case "help":
		h.printChatHelp()
		return true, nil
	}

	// Handle agent switching commands
	if strings.HasPrefix(lowerInput, "switch to ") {
		agentName := strings.TrimSpace(input[10:])
		return true, h.handleAgentSwitch(ctx, session, agentName)
	}

	return false, nil
}

// handleAgentSwitch processes agent switching commands.
func (h *chatHandler) handleAgentSwitch(ctx context.Context, session types.AgentSession, agentName string) error {
	if err := session.SwitchAgent(ctx, agentName); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Agent not found", zap.String("agent_name", agentName),
				zap.String("hint", "Available agents can be listed with 'agents' command."))
		} else {
			log.Error("Failed to switch agent", zap.Error(err), zap.String("agent_name", agentName))
			log.Error("Error switching agent", zap.String("agent_name", agentName), zap.Error(err))
		}
		return err
	}

	log.Info("Switched to agent", zap.String("agent_name", session.GetAgent().GetName()))
	return nil
}

// runInteractiveChat handles the interactive chat loop.
func (h *chatHandler) runInteractiveChat(ctx context.Context, session types.AgentSession) error {
	h.printWelcomeMessage(session)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		log.Info("You: ", zap.String("prompt", "true"))
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Handle special commands
		if handled, err := h.handleSpecialCommands(ctx, input, session); handled {
			if err != nil {
				log.Error("Error during command handling", zap.Error(err))
			}
			continue
		}

		// Process user message
		log.Info("Assistant: ", zap.String("prompt", "true"))
		if err := session.Chat(ctx, input); err != nil {
			h.handleChatError(err, input)
			continue
		}
		log.Info("")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil

}

// handleChatError handles different types of chat errors appropriately.
func (h *chatHandler) handleChatError(err error, userInput string) {
	if errors.IsConfiguration(err) {
		log.Error("Configuration error during chat", zap.Error(err))
		log.Error("Configuration error", zap.Error(err))
	} else if errors.IsProvider(err) {
		log.Error("Provider error during chat", zap.Error(err))
		log.Error("Service error", zap.Error(err))
	} else {
		log.Error("Chat error occurred", zap.Error(err), zap.String("user_input", userInput))
		log.Error("Chat error", zap.Error(err))
	}
}

// printMessageHistory displays the current message history.
func (h *chatHandler) printMessageHistory(session types.AgentSession) {
	history := session.GetMessageHistory()
	log.Info("Message History", zap.Int("message_count", len(history)))

	for i, msg := range history {
		role := string(msg.Role)
		if len(msg.Parts) > 0 {
			if textPart, ok := msg.Parts[0].(interface{ GetText() string }); ok {
				content := textPart.GetText()
				if len(content) > 100 {
					content = content[:100] + "..."
				}
				log.Info("Message", zap.Int("index", i+1), zap.String("role", role), zap.String("content", content))
			}
		}
	}
	log.Info("---")
}

// printChatHelp displays available chat commands.
func (h *chatHandler) printChatHelp() {
	log.Info("Available Commands:")
	log.Info("  exit, quit, bye - End the chat session")
	log.Info("  clear           - Clear message history")
	log.Info("  history         - View current message history")
	log.Info("  help            - Show this help message")
	log.Info("")
	log.Info("Example queries:")
	log.Info("  - What's the weather like in Chicago?")
	log.Info("  - Create a task to review project documentation")
	log.Info("  - Help me plan a team meeting for next week")
	log.Info("---")
}
