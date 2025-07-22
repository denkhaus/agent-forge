package startup

import (
	"context"
	"fmt"

	"github.com/denkhaus/agentforge/internal/config"
	"github.com/denkhaus/agentforge/internal/database"
	"github.com/denkhaus/agentforge/internal/git"
	"github.com/denkhaus/agentforge/internal/github"
	"github.com/denkhaus/agentforge/internal/prompts"
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
			opts.Context = context.Background()
		}
		
		// Create startup context
		startupCtx := &Context{
			Context: opts.Context,
			CLI:     cliCtx,
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
	log.Info("Initializing database", zap.String("path", dbPath))
	
	// TODO: Use DI container to get DatabaseManager instead of direct creation
	// For now, keep existing logic but mark for refactoring
	if dbPath == "" {
		dbPath = "agentforge.db"
	}
	
	cfg := &config.Config{
		DatabaseURL: dbPath,
	}
	
	dbManager, err := database.NewManager(cfg)
	if err != nil {
		return fmt.Errorf("failed to create database manager: %w", err)
	}
	
	if err := dbManager.Initialize(ctx.Context); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	
	ctx.DatabaseManager = dbManager
	log.Info("Database initialized successfully")
	return nil
}

// initializeGit initializes the Git client using DI
func (ctx *Context) initializeGit() error {
	log.Info("Initializing Git client")
	
	// TODO: Use DI container to get GitClient instead of direct creation
	// For now, keep existing logic but mark for refactoring
	gitClient := git.NewClient(log)
	ctx.GitClient = gitClient
	
	log.Info("Git client initialized successfully")
	return nil
}

// initializeGitHub initializes the GitHub client using DI
func (ctx *Context) initializeGitHub() error {
	log.Info("Initializing GitHub client")
	
	// TODO: Use DI container to get GitHubClient instead of direct creation
	// For now, keep existing logic but mark for refactoring
	githubClient := github.NewClient("")
	ctx.GitHubClient = githubClient
	
	log.Info("GitHub client initialized successfully")
	return nil
}

// initializePromptService initializes the prompt service
func (ctx *Context) initializePromptService() error {
	log.Info("Initializing prompt service")
	
	promptService := prompts.NewPromptService()
	ctx.PromptService = promptService
	
	log.Info("Prompt service initialized successfully")
	return nil
}