package startup

import (
	"fmt"

	"github.com/denkhaus/agentforge/internal/config"
	"github.com/denkhaus/agentforge/internal/container"
	"github.com/denkhaus/agentforge/internal/database"
	"github.com/denkhaus/agentforge/internal/git"
	"github.com/denkhaus/agentforge/internal/prompts"
	"github.com/denkhaus/agentforge/internal/types"
	"github.com/samber/do"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// NewStartupContext creates a startup context with fine-grained options
func NewStartupContext(options ...Option) func(*cli.Context) (*Context, error) {
	return func(cliCtx *cli.Context) (*Context, error) {
		// Apply options
		opts := &StartupOptions{}
		for _, option := range options {
			option(opts)
		}

		// Set default context if not provided
		if opts.Context == nil {
			opts.Context = cliCtx.Context
		}

		// Setup DI container if not provided
		var diContainer *do.Injector
		if opts.DIContainer != nil {
			diContainer = opts.DIContainer
		} else {
			// Load configuration for DI container setup
			cfg, err := config.Load(cliCtx)
			if err != nil {
				return nil, fmt.Errorf("failed to load configuration: %w", err)
			}
			diContainer = container.Setup(cfg)
		}

		// Create startup context
		startupCtx := &Context{
			Context:     opts.Context,
			CLI:         cliCtx,
			DIContainer: diContainer,
		}

		log.Info("Initializing startup context",
			zap.Bool("database", opts.RequireDatabase),
			zap.Bool("git", opts.RequireGit),
			zap.Bool("github", opts.RequireGitHub),
			zap.Bool("prompt_service", opts.RequirePromptService))

		// Initialize services based on options
		if err := startupCtx.initializeServices(opts); err != nil {
			startupCtx.Cleanup()
			return nil, fmt.Errorf("failed to initialize services: %w", err)
		}

		log.Info("Startup context initialized successfully")
		return startupCtx, nil
	}
}

// initializeServices initializes services based on startup options
func (ctx *Context) initializeServices(opts *StartupOptions) error {
	// Initialize database if required
	if opts.RequireDatabase {
		if err := ctx.initializeDatabase(opts.DatabasePath); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
	}

	// Initialize Git client if required
	if opts.RequireGit {
		if err := ctx.initializeGit(); err != nil {
			return fmt.Errorf("failed to initialize git: %w", err)
		}
	}

	// Initialize GitHub client if required
	if opts.RequireGitHub {
		if err := ctx.initializeGitHub(); err != nil {
			return fmt.Errorf("failed to initialize github: %w", err)
		}
	}

	// Initialize entity services
	if opts.RequirePromptService {
		if err := ctx.initializePromptService(); err != nil {
			return fmt.Errorf("failed to initialize prompt service: %w", err)
		}
	}

	// TODO: Add tool and agent service initialization when ready
	if opts.RequireToolService {
		log.Warn("Tool service not yet implemented")
	}

	if opts.RequireAgentService {
		log.Warn("Agent service not yet implemented")
	}

	return nil
}

// initializeDatabase initializes the database manager using DI
func (ctx *Context) initializeDatabase(dbPath string) error {
	log.Info("Initializing database using DI container")

	// Get DatabaseManager from DI container
	dbManager, err := do.Invoke[database.DatabaseManager](ctx.DIContainer)
	if err != nil {
		return fmt.Errorf("failed to get database manager from DI container: %w", err)
	}

	// Initialize the database
	if err := dbManager.Initialize(ctx.Context); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Store for backward compatibility (deprecated)
	ctx.DatabaseManager = dbManager
	log.Info("Database initialized successfully via DI")
	return nil
}

// initializeGit initializes the Git client using DI
func (ctx *Context) initializeGit() error {
	log.Info("Initializing Git client using DI container")

	// Get GitClient from DI container
	gitClient, err := do.Invoke[*git.Client](ctx.DIContainer)
	if err != nil {
		return fmt.Errorf("failed to get git client from DI container: %w", err)
	}

	// Store for backward compatibility (deprecated)
	ctx.GitClient = gitClient
	log.Info("Git client initialized successfully via DI")
	return nil
}

// initializeGitHub initializes the GitHub client using DI
func (ctx *Context) initializeGitHub() error {
	log.Info("Initializing GitHub client using DI container")

	// Get GitHubClient from DI container
	githubClient, err := do.Invoke[types.GitHubClient](ctx.DIContainer)
	if err != nil {
		return fmt.Errorf("failed to get github client from DI container: %w", err)
	}

	// Store for backward compatibility (deprecated)
	ctx.GitHubClient = githubClient
	log.Info("GitHub client initialized successfully via DI")
	return nil
}

// initializePromptService initializes the prompt service using DI
func (ctx *Context) initializePromptService() error {
	log.Info("Initializing prompt service using DI container")

	// Get PromptProvider from DI container
	promptProvider, err := do.Invoke[types.PromptProvider](ctx.DIContainer)
	if err != nil {
		return fmt.Errorf("failed to get prompt provider from DI container: %w", err)
	}

	// Store for backward compatibility (deprecated)
	// Note: This assumes PromptProvider can be used as PromptService
	// TODO: Align interfaces or create adapter
	ctx.PromptService = promptProvider.(prompts.PromptService)
	log.Info("Prompt service initialized successfully via DI")
	return nil
}
