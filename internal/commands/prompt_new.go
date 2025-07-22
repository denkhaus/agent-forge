package commands

import (
	"fmt"

	"github.com/denkhaus/agentforge/internal/startup"
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
	return startup.WithStartup(startup.WithPromptService())(func(ctx *startup.Context) error {
		// Get name from flag
		name := ctx.CLI.String("name")
		
		if name == "" {
			return fmt.Errorf("prompt name is required: forge prompt new --name <name>")
		}
		
		log.Info("Creating new prompt", zap.String("name", name))
		
		// Step 1: Create filesystem structure first
		fmt.Printf("Creating prompt '%s' filesystem structure...\n", name)
		
		err := ctx.PromptService.CreatePromptStructure(name)
		if err != nil {
			return fmt.Errorf("failed to create prompt structure: %w", err)
		}
		
		fmt.Printf("✓ Basic structure created in workspace/prompts/%s/\n", name)
		
		// Step 2: Launch TUI for interactive modifications
		fmt.Println("Starting interactive editor...")
		
		tui := getTUIModule()
		err = tui.RunPromptEditor(name)
		if err != nil {
			return fmt.Errorf("failed to run prompt editor: %w", err)
		}
		
		fmt.Printf("🎉 Prompt '%s' created successfully!\n", name)
		return nil
	})
}