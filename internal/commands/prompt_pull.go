package commands

import (
	"fmt"
	"strings"

	"github.com/denkhaus/agentforge/internal/startup"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// GetPromptPullCommand returns the prompt pull subcommand.
func GetPromptPullCommand() *cli.Command {
	return &cli.Command{
		Name:      "pull",
		Usage:     "Pull a prompt from GitHub",
		Action:    HandlePromptPull(),
		ArgsUsage: "<user/repo[@version]>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "version",
				Usage: "Specific version to pull",
			},
			&cli.BoolFlag{
				Name:  "force",
				Usage: "Force overwrite if prompt already exists",
			},
		},
	}
}

// HandlePromptPull handles the prompt pull command.
func HandlePromptPull() cli.ActionFunc {
	return startup.WithStartup(startup.Minimal()...)(func(ctx *startup.Context) error {
		args := ctx.CLI.Args()
		if args.Len() == 0 {
			return fmt.Errorf("repository required: forge prompt pull <user/repo[@version]>")
		}
		
		repo := args.First()
		version := ctx.CLI.String("version")
		force := ctx.CLI.Bool("force")
		
		// Parse repo and version
		if strings.Contains(repo, "@") && version == "" {
			parts := strings.Split(repo, "@")
			repo, version = parts[0], parts[1]
		}
		
		// Auto-expand GitHub URLs if needed
		if !strings.HasPrefix(repo, "github.com/") && !strings.HasPrefix(repo, "http") {
			repo = "github.com/" + repo
		}
		
		log.Info("Pulling prompt", 
			zap.String("repo", repo),
			zap.String("version", version),
			zap.Bool("force", force))

		promptService, err := getPromptServiceFromDI(ctx.DIContainer)
		if err != nil {
			log.Warn("Failed to get prompt service from DI, using direct instantiation", zap.Error(err))
			promptService = getPromptService()
		}
		
		// Pull the prompt using the prompt service
		fmt.Printf("Pulling prompt from %s", repo)
		if version != "" {
			fmt.Printf("@%s", version)
		}
		fmt.Println("...")

		err = promptService.PullPrompt(repo, version, force)
		if err != nil {
			return fmt.Errorf("failed to pull prompt: %w", err)
		}

		fmt.Printf("✓ Prompt successfully pulled from %s\n", repo)
		
		return nil
	})
}