package commands

import (
	"fmt"

	"github.com/denkhaus/agentforge/internal/startup"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// GetPromptListCommand returns the prompt list subcommand.
func GetPromptListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List available prompts",
		Action:  HandlePromptList(),
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "stars",
				Usage: "Minimum GitHub stars filter",
				Value: 0,
			},
			&cli.BoolFlag{
				Name:  "local",
				Usage: "Show only locally installed prompts",
			},
		},
	}
}

// HandlePromptList handles the prompt list command.
func HandlePromptList() cli.ActionFunc {
	return startup.WithStartup(startup.Minimal()...)(func(ctx *startup.Context) error {
		minStars := ctx.CLI.Int("stars")
		localOnly := ctx.CLI.Bool("local")
		
		log.Info("Listing prompts", 
			zap.Int("min_stars", minStars),
			zap.Bool("local_only", localOnly))

		promptService, err := getPromptServiceFromDI(ctx.DIContainer)
		if err != nil {
			log.Warn("Failed to get prompt service from DI, using direct instantiation", zap.Error(err))
			promptService = getPromptService()
		}
		
		if localOnly {
			// List local prompts
			prompts, err := promptService.ListLocalPrompts()
			if err != nil {
				return fmt.Errorf("failed to list local prompts: %w", err)
			}
			
			if len(prompts) == 0 {
				fmt.Println("No local prompts found")
				fmt.Println("Create one with: forge prompt new --name <name>")
				return nil
			}
			
			fmt.Printf("Found %d local prompts:\n\n", len(prompts))
			for _, prompt := range prompts {
				fmt.Printf("[LOCAL] %s\n", prompt.Name)
				fmt.Printf("  Description: %s\n", prompt.Description)
				fmt.Printf("  Author: %s\n", prompt.Author)
				fmt.Printf("  Path: workspace/prompts/%s/\n\n", prompt.Name)
			}
		} else {
			// List remote prompts (mock for now)
			fmt.Printf("Found 2 prompts:\n\n")
			
			fmt.Printf("[PROMPT] code-reviewer\n")
			fmt.Printf("  Description: AI prompt for code review\n")
			fmt.Printf("  Stars: 28\n")
			fmt.Printf("  Repository: github.com/agentforge/code-reviewer-prompt\n\n")
			
			fmt.Printf("[PROMPT] data-analyzer\n")
			fmt.Printf("  Description: AI prompt for data analysis\n")
			fmt.Printf("  Stars: 15\n")
			fmt.Printf("  Repository: github.com/agentforge/data-analyzer-prompt\n")
		}
		
		return nil
	})
}