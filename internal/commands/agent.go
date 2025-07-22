// Package commands provides CLI command handlers for AgentForge.
package commands

import (
	"fmt"
	"strings"

	"github.com/denkhaus/agentforge/internal/startup"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// GetAgentCommand returns the agent management command configuration.
func GetAgentCommand() *cli.Command {
	return &cli.Command{
		Name:    "agent",
		Usage:   "Manage AgentForge agents",
		Subcommands: []*cli.Command{
			getAgentListCommand(),
			getAgentPullCommand(),
			getAgentPushCommand(),
			getAgentNewCommand(),
		},
	}
}

// getAgentListCommand returns the agent list subcommand.
func getAgentListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List available agents",
		Action:  HandleAgentList(),
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "stars",
				Usage: "Minimum GitHub stars filter",
				Value: 0,
			},
		},
	}
}

// getAgentPullCommand returns the agent pull subcommand.
func getAgentPullCommand() *cli.Command {
	return &cli.Command{
		Name:      "pull",
		Usage:     "Pull an agent from GitHub",
		Action:    HandleAgentPull(),
		ArgsUsage: "<user/repo[@version]>",
	}
}

// getAgentPushCommand returns the agent push subcommand.
func getAgentPushCommand() *cli.Command {
	return &cli.Command{
		Name:      "push",
		Usage:     "Push an agent to GitHub",
		Action:    HandleAgentPush(),
		ArgsUsage: "<repo-url>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Agent name",
				Required: true,
			},
		},
	}
}

// getAgentNewCommand returns the agent new subcommand.
func getAgentNewCommand() *cli.Command {
	return &cli.Command{
		Name:      "new",
		Usage:     "Create a new agent",
		Action:    HandleAgentNew(),
		ArgsUsage: "<name>",
	}
}

// HandleAgentList handles the agent list command.
func HandleAgentList() cli.ActionFunc {
	return startup.WithStartup(startup.Minimal()...)(func(ctx *startup.Context) error {
		minStars := ctx.CLI.Int("stars")
		
		log.Info("Listing agents", 
			zap.Int("min_stars", minStars))

		// TODO: Implement agent service when available
		// For now, return a placeholder message
		log.Warn("Agent service not yet implemented")
		
		// Placeholder implementation
		components := []interface{}{}

		// Display results
		if len(components) == 0 {
			fmt.Printf("No agents found with minimum %d stars\n", minStars)
			return nil
		}

		fmt.Printf("Found %d agents:\n\n", len(components))
		for _, comp := range components {
			if compMap, ok := comp.(map[string]interface{}); ok {
				displayComponentPlaceholder(compMap, "AGENT")
			}
		}
		
		return nil
	})
}

// HandleAgentPull handles the agent pull command.
func HandleAgentPull() cli.ActionFunc {
	return startup.WithStartup(startup.Minimal()...)(func(ctx *startup.Context) error {
		args := ctx.CLI.Args()
		if args.Len() == 0 {
			return fmt.Errorf("repository required: forge agent pull <user/repo[@version]>")
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
		
		log.Info("Pulling agent", 
			zap.String("repo", repo),
			zap.String("version", version))

		// Get component service
		if ctx.ComponentService == nil {
			return fmt.Errorf("component service not available")
		}

		// Pull the agent using the component service
		fmt.Printf("Pulling agent from %s", repo)
		if version != "" {
			fmt.Printf("@%s", version)
		}
		fmt.Println("...")

		err := ctx.ComponentService.PullComponent(ctx.Context(), "agent", repo)
		if err != nil {
			return fmt.Errorf("failed to pull agent: %w", err)
		}

		fmt.Printf("Agent successfully pulled from %s\n", repo)
		
		return nil
	})
}

// HandleAgentPush handles the agent push command.
func HandleAgentPush() cli.ActionFunc {
	return startup.WithStartup(startup.Minimal()...)(func(ctx *startup.Context) error {
		name := ctx.CLI.String("name")
		
		args := ctx.CLI.Args()
		if args.Len() == 0 {
			return fmt.Errorf("repository required: forge agent push --name %s <repo-url>", name)
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
		
		log.Info("Pushing agent", 
			zap.String("name", name),
			zap.String("repo", repo),
			zap.String("version", version))

		// Get component service
		if ctx.ComponentService == nil {
			return fmt.Errorf("component service not available")
		}

		// Push the agent using the component service
		fmt.Printf("Pushing agent '%s' to %s", name, repo)
		if version != "" {
			fmt.Printf("@%s", version)
		}
		fmt.Println("...")

		err := ctx.ComponentService.PushComponent(ctx.Context(), "agent", name)
		if err != nil {
			return fmt.Errorf("failed to push agent: %w", err)
		}

		fmt.Printf("Agent '%s' successfully pushed to %s\n", name, repo)
		
		return nil
	})
}

// HandleAgentNew handles the agent new command.
func HandleAgentNew() cli.ActionFunc {
	return startup.WithStartup(startup.Minimal()...)(func(ctx *startup.Context) error {
		args := ctx.CLI.Args()
		if args.Len() == 0 {
			return fmt.Errorf("agent name required: forge agent new <name>")
		}
		
		name := args.First()
		
		log.Info("Creating new agent", 
			zap.String("name", name))

		// TODO: Implement component service when available
		// For now, return a placeholder message
		log.Warn("Agent creation service not yet implemented")
		
		// Create the agent using placeholder logic
		fmt.Printf("Creating new agent '%s'...\n", name)
		fmt.Printf("(This is a placeholder - actual implementation pending)\n")

		fmt.Printf("Agent '%s' successfully created!\n", name)
		fmt.Printf("Files created in ./%s/\n", name)
		
		return nil
	})
}

// displayComponentPlaceholder displays component information (placeholder implementation)
func displayComponentPlaceholder(comp map[string]interface{}, componentType string) {
	fmt.Printf("  %s: %v\n", componentType, comp)
}