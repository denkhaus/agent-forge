package startup

import (
	cli "github.com/urfave/cli/v2"
)

// CommandHandler represents a command handler function that uses startup context.
type CommandHandler func(*Context) error

// WithStartup wraps a command handler with startup context initialization using fine-grained options.
func WithStartup(options ...Option) func(CommandHandler) cli.ActionFunc {
	return func(handler CommandHandler) cli.ActionFunc {
		return func(cliCtx *cli.Context) error {
			// Initialize startup context with options
			startupCtx, err := NewStartupContext(options...)(cliCtx)
			if err != nil {
				return err
			}
			
			// Ensure cleanup happens
			defer startupCtx.Cleanup()
			
			// Execute the handler
			return handler(startupCtx)
		}
	}
}