package commands

import (
	"fmt"

	cli "github.com/urfave/cli/v2"
)

// GetVersionCommand returns the version command configuration.
func GetVersionCommand() *cli.Command {
	return &cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Show version information",
		Action:  HandleVersion(),
	}
}

// HandleVersion handles the version command.
func HandleVersion() cli.ActionFunc {
	return func(c *cli.Context) error {
		fmt.Printf("AgentForge version %s\n", c.App.Version)
		return nil
	}
}
