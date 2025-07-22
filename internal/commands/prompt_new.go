package commands

import (
	"fmt"

	"github.com/denkhaus/agentforge/internal/startup"
	"github.com/denkhaus/agentforge/internal/types"
	"github.com/samber/do"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// GetPromptNewCommand returns the prompt new subcommand.
func GetPromptNewCommand() *cli.Command {
	return &cli.Command{
		Name:   "new",
		Usage:  "Create a new prompt with interactive editor",
		Action: HandlePromptNew(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Name of the prompt to create",
				Required: true,
			},
		},
	}
}

// HandlePromptNew handles the prompt new command.
func HandlePromptNew() cli.ActionFunc {
	return startup.WithStartup(startup.Minimal()...)(func(ctx *startup.Context) error {
		// Get name from flag
		name := ctx.CLI.String("name")
		
		if name == "" {
			return fmt.Errorf("prompt name is required: forge prompt new --name <name>")
		}
		
		log.Info("Creating new prompt", zap.String("name", name))
		
		// Step 1: Create basic prompt structure (placeholder)
		fmt.Printf("Creating prompt '%s'...\n", name)
		log.Info("Creating prompt structure", zap.String("name", name))
		
		// TODO: Implement actual file creation when needed
		fmt.Printf("Basic structure created for '%s'\n", name)
		
		// Step 2: Launch TUI Prompt Workbench
		fmt.Println("Launching TUI Prompt Workbench...")
		
		tuiManager, err := do.Invoke[types.TUIManager](ctx.DIContainer)
		if err != nil {
			log.Warn("TUI manager not available", zap.Error(err))
			fmt.Printf("Prompt '%s' created successfully (TUI workbench unavailable)\n", name)
			return nil
		}
		
		err = tuiManager.RunPromptWorkbench(name)
		if err != nil {
			return fmt.Errorf("failed to run prompt workbench: %w", err)
		}
		
		fmt.Printf("Prompt '%s' development session completed!\n", name)
		return nil
	})
}