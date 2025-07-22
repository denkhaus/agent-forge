// Package main provides the entry point for the AgentForge application.
package main

import (
	"log"
	"os"

	"github.com/denkhaus/agentforge/internal/commands"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := createApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// createApp creates and configures the CLI application.
func createApp() *cli.App {
	return &cli.App{
		Name:    "forge",
		Usage:   "AgentForge - Git-native AI agent development platform",
		Version: "0.1.0",
		Description: `AgentForge enables developers to build, share, and manage AI agent components
		through a Git-native ecosystem. Create tools, prompts, and agents with ease.`,
		Authors:  getAuthors(),
		Flags:    getGlobalFlags(),
		Commands: getCommands(),
		Action:   cli.ShowAppHelp,
	}
}

// getAuthors returns the application authors.
func getAuthors() []*cli.Author {
	return []*cli.Author{
		{
			Name:  "AgentForge Team",
			Email: "team@agentforge.dev",
		},
	}
}

// getGlobalFlags returns the global CLI flags.
func getGlobalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Load configuration from `FILE`",
			EnvVars: []string{"CONFIG_FILE"},
		},
		&cli.StringFlag{
			Name:    "log-level",
			Aliases: []string{"l"},
			Value:   "info",
			Usage:   "Set log level (debug, info, warn, error)",
			EnvVars: []string{"LOG_LEVEL"},
		},
	}
}

// getCommands returns the CLI commands.
func getCommands() []*cli.Command {
	return []*cli.Command{
		commands.GetServerCommand(),
		commands.GetChatCommand(),
		commands.GetVersionCommand(),
		commands.GetPromptCommand(),
		commands.GetAgentCommand(),
	}
}
