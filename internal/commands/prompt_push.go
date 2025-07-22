package commands

import (
	"fmt"
	"strings"

	"github.com/denkhaus/agentforge/internal/startup"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// GetPromptPushCommand returns the prompt push subcommand.
func GetPromptPushCommand() *cli.Command {
	return &cli.Command{
		Name:      "push",
		Usage:     "Push a prompt to GitHub",
		Action:    HandlePromptPush(),
		ArgsUsage: "<repo-url>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Prompt name",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "message",
				Aliases: []string{"m"},
				Usage: "Commit message",
				Value: "Update prompt",
			},
			&cli.StringFlag{
				Name:  "tag",
				Usage: "Create a git tag for this version",
			},
		},
	}
}

// HandlePromptPush handles the prompt push command.
func HandlePromptPush() cli.ActionFunc {
	return startup.WithStartup(startup.Minimal()...)(func(ctx *startup.Context) error {
		name := ctx.CLI.String("name")
		message := ctx.CLI.String("message")
		tag := ctx.CLI.String("tag")
		
		args := ctx.CLI.Args()
		if args.Len() == 0 {
			return fmt.Errorf("repository required: forge prompt push --name %s <repo-url>", name)
		}
		
		repo := args.First()
		
		// Parse repo and version
		var version string
		if strings.Contains(repo, "@") {
			parts := strings.Split(repo, "@")
			repo, version = parts[0], parts[1]
		}
		
		// Auto-expand GitHub URLs if needed
		if !strings.HasPrefix(repo, "github.com/") && !strings.HasPrefix(repo, "http") {
			repo = "github.com/" + repo
		}
		
		log.Info("Pushing prompt", 
			zap.String("name", name),
			zap.String("repo", repo),
			zap.String("version", version),
			zap.String("message", message),
			zap.String("tag", tag))

		promptService, err := getPromptServiceFromDI(ctx.DIContainer)
		if err != nil {
			log.Warn("Failed to get prompt service from DI, using direct instantiation", zap.Error(err))
			promptService = getPromptService()
		}
		
		// Push the prompt using the prompt service
		fmt.Printf("Pushing prompt '%s' to %s", name, repo)
		if version != "" {
			fmt.Printf("@%s", version)
		}
		fmt.Println("...")

		err = promptService.PushPrompt(name, repo, message, tag)
		if err != nil {
			return fmt.Errorf("failed to push prompt: %w", err)
		}

		fmt.Printf("✓ Prompt '%s' successfully pushed to %s\n", name, repo)
		if tag != "" {
			fmt.Printf("✓ Tagged as %s\n", tag)
		}
		
		return nil
	})
}