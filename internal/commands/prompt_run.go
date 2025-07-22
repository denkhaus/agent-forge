package commands

import (
	"fmt"

	"github.com/denkhaus/agentforge/internal/startup"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// GetPromptRunCommand returns the prompt run subcommand.
func GetPromptRunCommand() *cli.Command {
	return &cli.Command{
		Name:      "run",
		Usage:     "Run a prompt with variables",
		Action:    HandlePromptRun(),
		ArgsUsage: "<prompt-name>",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "var",
				Aliases: []string{"v"},
				Usage:   "Set template variables (e.g., --var task='review this code')",
			},
			&cli.StringFlag{
				Name:  "output",
				Aliases: []string{"o"},
				Usage: "Output format (text, json, yaml)",
				Value: "text",
			},
			&cli.BoolFlag{
				Name:  "interactive",
				Aliases: []string{"i"},
				Usage: "Interactive mode to set variables",
			},
		},
	}
}

// HandlePromptRun handles the prompt run command.
func HandlePromptRun() cli.ActionFunc {
	return startup.WithStartup(startup.Minimal()...)(func(ctx *startup.Context) error {
		args := ctx.CLI.Args()
		if args.Len() == 0 {
			return fmt.Errorf("prompt name required: forge prompt run <prompt-name>")
		}
		
		name := args.First()
		variables := ctx.CLI.StringSlice("var")
		output := ctx.CLI.String("output")
		interactive := ctx.CLI.Bool("interactive")
		
		log.Info("Running prompt", 
			zap.String("name", name),
			zap.Strings("variables", variables),
			zap.String("output", output),
			zap.Bool("interactive", interactive))

		promptService := getPromptService()
		
		// Load prompt data
		promptData, err := promptService.LoadPromptData(name)
		if err != nil {
			return fmt.Errorf("failed to load prompt '%s': %w", name, err)
		}
		
		fmt.Printf("Running prompt: %s\n", promptData.Name)
		fmt.Printf("Description: %s\n\n", promptData.Description)
		
		if interactive {
			// Launch interactive variable setter
			fmt.Println("Interactive mode - launching variable editor...")
			
			tui := getTUIModule()
			err = tui.RunPromptVariableEditor(name, promptData)
			if err != nil {
				return fmt.Errorf("failed to run variable editor: %w", err)
			}
		} else {
			// Parse variables from command line
			varMap := make(map[string]string)
			for _, variable := range variables {
				// Parse key=value format
				// TODO: Implement variable parsing
				_ = variable // Suppress unused variable warning
			}
			
			// Execute prompt with variables
			result, err := promptService.ExecutePrompt(name, varMap)
			if err != nil {
				return fmt.Errorf("failed to execute prompt: %w", err)
			}
			
			// Output result
			switch output {
			case "json":
				fmt.Printf("{\n  \"prompt\": \"%s\",\n  \"result\": \"%s\"\n}\n", name, result)
			case "yaml":
				fmt.Printf("prompt: %s\nresult: |\n  %s\n", name, result)
			default:
				fmt.Printf("Result:\n%s\n", result)
			}
		}
		
		return nil
	})
}