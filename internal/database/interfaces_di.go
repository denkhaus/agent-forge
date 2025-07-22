// Package database provides DI-aware interfaces and factory functions.
package database

import (
	"github.com/samber/do"
)

// DIAwareDatabaseManager extends DatabaseManager with DI capabilities
type DIAwareDatabaseManager interface {
	DatabaseManager
	GetClient() DatabaseClient
}

// NewDatabaseManagerFromDI creates a DatabaseManager using dependency injection
func NewDatabaseManagerFromDI(injector *do.Injector) (DatabaseManager, error) {
	return do.Invoke[DatabaseManager](injector)
}

// NewDatabaseClientFromDI creates a DatabaseClient using dependency injection
func NewDatabaseClientFromDI(injector *do.Injector) (DatabaseClient, error) {
	// Get manager and extract client
	manager, err := do.Invoke[DatabaseManager](injector)
	if err != nil {
		return nil, err
	}
	
	// Extract client from manager
	return manager.GetClient(), nil
}

// NewRepositoryServiceFromDI creates a RepositoryService using dependency injection
func NewRepositoryServiceFromDI(injector *do.Injector) (RepositoryService, error) {
	client, err := NewDatabaseClientFromDI(injector)
	if err != nil {
		return nil, err
	}
	return NewRepositoryService(client), nil
}

// NewConfigServiceFromDI creates a ConfigService using dependency injection
func NewConfigServiceFromDI(injector *do.Injector) (ConfigService, error) {
	client, err := NewDatabaseClientFromDI(injector)
	if err != nil {
		return nil, err
	}
	return NewConfigService(client), nil
}