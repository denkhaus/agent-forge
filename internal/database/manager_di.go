// Package database provides DI-aware database manager implementation.
package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/denkhaus/agentforge/internal/config"
	"github.com/samber/do"
	"go.uber.org/zap"
)

// diAwareManager is a DI-aware implementation of DatabaseManager
type diAwareManager struct {
	client      DatabaseClient
	config      Config
	diContainer *do.Injector
}

// NewManagerWithDI creates a new database manager using dependency injection
func NewManagerWithDI(cfg *config.Config, diContainer *do.Injector) (DatabaseManager, error) {
	// Get user home directory for database location
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create AgentForge data directory
	agentForgeDir := filepath.Join(homeDir, ".agentforge")
	if err := os.MkdirAll(agentForgeDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create AgentForge directory: %w", err)
	}

	// Database configuration
	dbConfig := Config{
		DatabasePath: filepath.Join(agentForgeDir, "agentforge.db"),
	}

	// Create database client directly (since it's the foundation)
	client, err := NewClient(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %w", err)
	}

	return &diAwareManager{
		client:      client,
		config:      dbConfig,
		diContainer: diContainer,
	}, nil
}

// Initialize connects to the database and performs any necessary setup.
func (m *diAwareManager) Initialize(ctx context.Context) error {
	log.Info("Initializing DI-aware database manager")

	// Connect to database and run migrations
	if err := m.client.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	return nil
}

// Shutdown gracefully closes database connections.
func (m *diAwareManager) Shutdown() error {
	log.Info("Shutting down DI-aware database manager")

	if err := m.client.Disconnect(); err != nil {
		return fmt.Errorf("failed to disconnect from database: %w", err)
	}

	log.Info("DI-aware database manager shutdown complete")
	return nil
}

// GetRepositoryService returns the repository service from DI container
func (m *diAwareManager) GetRepositoryService() RepositoryService {
	// Get from DI container instead of storing directly
	service, err := do.Invoke[RepositoryService](m.diContainer)
	if err != nil {
		log.Error("Failed to get repository service from DI container", zap.Error(err))
		// Fallback to direct creation for backward compatibility
		return NewRepositoryService(m.client)
	}
	return service
}

// GetDatabasePath returns the database file path
func (m *diAwareManager) GetDatabasePath() string {
	return m.config.DatabasePath
}

// GetClient returns the database client
func (m *diAwareManager) GetClient() DatabaseClient {
	return m.client
}

// GetConfigService returns the config service from DI container
func (m *diAwareManager) GetConfigService() ConfigService {
	// Get from DI container instead of storing directly
	service, err := do.Invoke[ConfigService](m.diContainer)
	if err != nil {
		log.Error("Failed to get config service from DI container", zap.Error(err))
		// Fallback to direct creation for backward compatibility
		return NewConfigService(m.client)
	}
	return service
}