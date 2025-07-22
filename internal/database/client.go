// Package database provides database connectivity and operations for AgentForge.
package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/denkhaus/agentforge/internal/database/ent"
	"go.uber.org/zap"

	_ "github.com/mattn/go-sqlite3"
)

// client represents the private database client implementation.
type client struct {
	ent    *ent.Client
	dbPath string
}

// Config represents database configuration.
type Config struct {
	DatabasePath string
}

// NewClient creates a new database client.
func NewClient(config Config) (DatabaseClient, error) {
	log.Info("Creating database client", zap.String("path", config.DatabasePath))
	
	// Ensure database directory exists
	dbDir := filepath.Dir(config.DatabasePath)
	if err := ensureDir(dbDir); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}
	
	// Create Ent client with SQLite
	entClient, err := ent.Open("sqlite3", config.DatabasePath+"?_fk=1")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	return &client{
		ent:    entClient,
		dbPath: config.DatabasePath,
	}, nil
}

// Connect establishes connection and runs migrations.
func (c *client) Connect(ctx context.Context) error {
	log.Info("Connecting to database and running migrations")
	
	// Run auto-migration
	if err := c.ent.Schema.Create(ctx); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	
	log.Info("Database connection established and schema migrated")
	return nil
}

// Disconnect closes the database connection.
func (c *client) Disconnect() error {
	log.Info("Disconnecting from database")
	
	if err := c.ent.Close(); err != nil {
		log.Error("Error disconnecting from database", zap.Error(err))
		return err
	}
	
	log.Info("Database disconnected")
	return nil
}

// Close is an alias for Disconnect to satisfy the components.DatabaseClient interface.
func (c *client) Close() error {
	return c.Disconnect()
}

// GetEnt returns the underlying Ent client.
func (c *client) GetEnt() *ent.Client {
	return c.ent
}

// GetDatabasePath returns the path to the database file.
func (c *client) GetDatabasePath() string {
	return c.dbPath
}

// ensureDir creates directory if it doesn't exist.
func ensureDir(dir string) error {
	if dir == "" || dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}