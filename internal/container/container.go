// Package container provides dependency injection setup and management.
package container

import (
	"github.com/samber/do"
	"go.uber.org/zap"

	"github.com/denkhaus/agentforge/internal/config"
	"github.com/denkhaus/agentforge/internal/database"
	"github.com/denkhaus/agentforge/internal/git"
	"github.com/denkhaus/agentforge/internal/github"
	"github.com/denkhaus/agentforge/internal/logger"
	"github.com/denkhaus/agentforge/internal/providers"
	"github.com/denkhaus/agentforge/internal/session"
	"github.com/denkhaus/agentforge/internal/types"
)

var (
	injector = do.New()
)

func GetInjector() *do.Injector {
	return injector
}

// Setup initializes the dependency injection container with all required services.
func Setup(cfg *config.Config) *do.Injector {
	// Create a new injector for each setup to avoid conflicts in tests
	newInjector := do.New()

	// Register configuration
	do.ProvideValue(newInjector, cfg)

	// Provide the *zap.Logger instance - single logger registration
	do.Provide(newInjector, func(i *do.Injector) (*zap.Logger, error) {
		return logger.Create(cfg.LogLevel)
	})

	// Register session factory (only factory we keep - adds real value)
	do.Provide(newInjector, func(i *do.Injector) (types.SessionFactory, error) {
		// Correctly invoke LLMService
		llmService := do.MustInvoke[types.LLMService](i)
		return session.NewFactory(llmService), nil
	})

	// Register individual tool providers with specific names to avoid circular dependency
	do.ProvideNamed(newInjector, "internalProvider", func(i *do.Injector) (types.ToolProvider, error) {
		return providers.NewToolProvider(i)
	})

	// Register MCP tool provider
	do.ProvideNamed(newInjector, "mcpProvider", func(i *do.Injector) (types.ToolProvider, error) {
		log := do.MustInvoke[*zap.Logger](i)
		config := do.MustInvoke[*config.Config](i)
		mcpConfig := config.GetMCPConfig()
		return providers.NewMCPToolProvider(log, mcpConfig)
	})

	// Register aggregated tool provider (what consumers actually use)
	do.Provide(newInjector, func(i *do.Injector) (types.ToolProvider, error) {
		log := do.MustInvoke[*zap.Logger](i)

		// Get individual providers by name to avoid circular dependency
		internalProvider := do.MustInvokeNamed[types.ToolProvider](i, "internalProvider")
		mcpProvider := do.MustInvokeNamed[types.ToolProvider](i, "mcpProvider")

		// Create aggregated provider
		return providers.NewAggregatedToolProvider(log, internalProvider, mcpProvider), nil
	})

	do.Provide(newInjector, func(i *do.Injector) (types.PromptProvider, error) {
		return providers.NewPromptProvider(i)
	})

	// Register database services
	do.Provide(newInjector, func(i *do.Injector) (database.DatabaseManager, error) {
		cfg := do.MustInvoke[*config.Config](i)
		return database.NewManagerWithDI(cfg, newInjector)
	})

	// Register DatabaseClient through DatabaseManager
	do.Provide(newInjector, func(i *do.Injector) (database.DatabaseClient, error) {
		manager := do.MustInvoke[database.DatabaseManager](i)
		return manager.GetClient(), nil
	})

	// Register Git services
	do.Provide(newInjector, func(i *do.Injector) (*git.Client, error) {
		log := do.MustInvoke[*zap.Logger](i)
		return git.NewClient(log), nil
	})

	// Register GitHub services  
	do.Provide(newInjector, func(i *do.Injector) (types.GitHubClient, error) {
		cfg := do.MustInvoke[*config.Config](i)
		return github.NewClient(cfg.GitHubToken), nil
	})

	// Register database services (RepositoryService and ConfigService)
	do.Provide(newInjector, func(i *do.Injector) (database.RepositoryService, error) {
		client := do.MustInvoke[database.DatabaseClient](i)
		return database.NewRepositoryService(client), nil
	})

	do.Provide(newInjector, func(i *do.Injector) (database.ConfigService, error) {
		client := do.MustInvoke[database.DatabaseClient](i)
		return database.NewConfigService(client), nil
	})

	return newInjector
}

// Cleanup performs cleanup operations on the dependency injection container.
func Cleanup(injector *do.Injector) {
	// Get logger from container if available, otherwise use the internal one
	var log *zap.Logger
	if injector != nil {
		if containerLogger, err := do.Invoke[*zap.Logger](injector); err == nil {
			log = containerLogger
		}
	}
	if log == nil {
		log = do.MustInvoke[*zap.Logger](injector)
	}

	if log != nil {
		log.Info("Cleaning up dependency injection container")
	}
	if err := injector.Shutdown(); err != nil {
		if log != nil {
			log.Error("Failed to shutdown injector", zap.Error(err))
		}
	}
}
