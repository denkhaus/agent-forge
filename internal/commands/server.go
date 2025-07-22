// Package commands provides command-line interface handlers for the MCP Planner application.
package commands

import (
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/denkhaus/agentforge/internal/container"
	"github.com/denkhaus/agentforge/internal/shared"
	"github.com/denkhaus/agentforge/internal/signals"
)

// Package-level variables for mocking in tests.
// Package-level variables for mocking in tests.
// var (
// 	signalsWithInterruptContext = signals.WithInterruptContextFunc
// )

// GetServerCommand returns the server command configuration.
func GetServerCommand() *cli.Command {
	return &cli.Command{
		Name:    "server",
		Aliases: []string{"s", "serve"},
		Usage:   "Start the AgentForge server",
		Action:  HandleServer(),
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   8080,
				Usage:   "Server port",
				EnvVars: []string{"PORT"},
			},
		},
	}
}

// HandleServer creates a new server command handler.
func HandleServer() cli.ActionFunc {
	return func(c *cli.Context) error {
		// Create context that cancels on interrupt signals
		ctx, cancel := signals.WithInterruptContextFunc(c.Context)
		defer cancel()

		sc, err := shared.StartupCommon(c)
		if err != nil {
			return err
		}

		defer container.Cleanup(sc.Injector)
		defer log.Sync()

		cfg := sc.Config

		// Start application

		log.Info("Starting MCP Planner server",
			zap.String("version", c.App.Version),
			zap.Int("port", cfg.Port),
			zap.String("log_level", cfg.LogLevel),
			zap.String("environment", cfg.Environment))

		// Wait for shutdown signal
		<-ctx.Done()

		log.Info("Shutting down MCP Planner server")

		return nil
	}
}
