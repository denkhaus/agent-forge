// Package commands provides CLI command handlers for AgentForge.
package commands

import (
	cli "github.com/urfave/cli/v2"
)

// GetPromptCommand returns the prompt management command configuration.
func GetPromptCommand() *cli.Command {
	return &cli.Command{
		Name:    "prompt",
		Usage:   "Manage AgentForge prompts",
		Subcommands: []*cli.Command{
			GetPromptListCommand(),
			GetPromptPullCommand(),
			GetPromptPushCommand(),
			GetPromptNewCommand(),
			GetPromptRunCommand(),
		},
	}
}
