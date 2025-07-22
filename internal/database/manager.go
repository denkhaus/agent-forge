package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/denkhaus/agentforge/internal/config"
)

// manager handles database lifecycle and service initialization.
type manager struct {
	client            DatabaseClient
	repositoryService RepositoryService
	configService     ConfigService
	config            Config
}

// NewManager creates a new database manager.
// TODO: rename this to NewDatabaseManager
func NewManager(cfg *config.Config) (DatabaseManager, error) {
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

	// Create database client
	client, err := NewClient(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %w", err)
	}

	// Create services
	repositoryService := NewRepositoryService(client)
	configService := NewConfigService(client)

	return &manager{
		client:            client,
		repositoryService: repositoryService,
		configService:     configService,
		config:            dbConfig,
	}, nil
}

// Initialize connects to the database and performs any necessary setup.
func (m *manager) Initialize(ctx context.Context) error {
	log.Info("Initializing database manager")

	// Connect to database and run migrations
	if err := m.client.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	return nil
}

// Shutdown gracefully closes database connections.
func (m *manager) Shutdown() error {
	log.Info("Shutting down database manager")

	if err := m.client.Disconnect(); err != nil {
		return fmt.Errorf("failed to disconnect from database: %w", err)
	}

	log.Info("Database manager shutdown complete")
	return nil
}

// GetRepositoryService returns the repository service.
func (m *manager) GetRepositoryService() RepositoryService {
	return m.repositoryService
}

// GetDatabasePath returns the database file path
func (m *manager) GetDatabasePath() string {
	return "agentforge.db"
}

// GetClient returns the database client (for DI support)
func (m *manager) GetClient() DatabaseClient {
	return m.client
}
