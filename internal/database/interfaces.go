// Package database provides database connectivity and operations for AgentForge.
package database

import (
	"context"
	"time"

	"github.com/denkhaus/agentforge/internal/database/ent"
)

// DatabaseManager defines the interface for database lifecycle management.
type DatabaseManager interface {
	Initialize(ctx context.Context) error
	Shutdown() error
	GetRepositoryService() RepositoryService
	GetDatabasePath() string
	GetClient() DatabaseClient
}

// RepositoryService defines the interface for repository management operations.
type RepositoryService interface {
	CreateRepository(ctx context.Context, req CreateRepositoryRequest) (*ent.Repository, error)
	GetRepository(ctx context.Context, id string) (*ent.Repository, error)
	GetRepositoryByName(ctx context.Context, name string) (*ent.Repository, error)
	ListRepositories(ctx context.Context, opts ListRepositoriesOptions) ([]*ent.Repository, error)
	UpdateRepository(ctx context.Context, id string, req UpdateRepositoryRequest) (*ent.Repository, error)
	DeleteRepository(ctx context.Context, id string) error
}

// ComponentService is an alias to the components package interface to avoid duplication

// ConfigService defines the interface for configuration management operations.
type ConfigService interface {
	GetConfig(ctx context.Context, key string) (string, error)
	SetConfig(ctx context.Context, key, value, configType string) error
	GetBoolConfig(ctx context.Context, key string) (bool, error)
	GetIntConfig(ctx context.Context, key string) (int, error)
	SetBoolConfig(ctx context.Context, key string, value bool) error
	SetIntConfig(ctx context.Context, key string, value int) error
}

// DatabaseClient defines the interface for low-level database operations.
type DatabaseClient interface {
	Connect(ctx context.Context) error
	Disconnect() error
	Close() error
	GetEnt() *ent.Client
	GetDatabasePath() string
}

// Request and response types for repository operations
type CreateRepositoryRequest struct {
	Name           string
	URL            string
	Type           string
	IsActive       bool
	DefaultBranch  string
	HasWriteAccess bool
	AccessToken    *string
}

type UpdateRepositoryRequest struct {
	LastSync     *time.Time
	SyncStatus   *string
	Manifest     interface{}
	ManifestHash *string
	IsActive     *bool
}

type ListRepositoriesOptions struct {
	IsActive *bool
	Type     *string
	Limit    int
	Offset   int
}

// Request and response types for component operations
type CreateComponentRequest struct {
	Name          string
	Namespace     string
	Version       string
	Kind          string
	Description   string
	Author        string
	License       string
	Homepage      *string
	Documentation *string
	Tags          []string
	Categories    []string
	Keywords      []string
	Stability     string
	Maturity      string
	ForgeVersion  string
	Platforms     []string
	Spec          string
	SpecHash      string
	RepositoryID  string
	CommitHash    string
	Branch        string
}

type ListComponentsOptions struct {
	Kind         *string
	Stability    *string
	IsInstalled  *bool
	RepositoryID *string
	Limit        int
	Offset       int
}

type SearchComponentsOptions struct {
	Kind      *string
	Stability *string
	Limit     int
}