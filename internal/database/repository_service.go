package database

import (
	"context"
	"fmt"

	"github.com/denkhaus/agentforge/internal/database/ent"
	"github.com/denkhaus/agentforge/internal/database/ent/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// repositoryService provides repository management operations.
type repositoryService struct {
	client DatabaseClient
}

// NewRepositoryService creates a new repository service.
func NewRepositoryService(client DatabaseClient) RepositoryService {
	return &repositoryService{
		client: client,
	}
}

// CreateRepository creates a new repository record.
func (rs *repositoryService) CreateRepository(ctx context.Context, req CreateRepositoryRequest) (*ent.Repository, error) {
	log.Info("Creating repository", 
		zap.String("name", req.Name),
		zap.String("url", req.URL))
	
	repo, err := rs.client.GetEnt().Repository.Create().
		SetID(uuid.New().String()).
		SetName(req.Name).
		SetURL(req.URL).
		SetType(repository.Type(req.Type)).
		SetIsActive(req.IsActive).
		SetDefaultBranch(req.DefaultBranch).
		SetHasWriteAccess(req.HasWriteAccess).
		SetNillableAccessToken(req.AccessToken).
		Save(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}
	
	log.Info("Repository created", zap.String("id", repo.ID))
	return repo, nil
}

// GetRepository retrieves a repository by ID.
func (rs *repositoryService) GetRepository(ctx context.Context, id string) (*ent.Repository, error) {
	repo, err := rs.client.GetEnt().Repository.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}
	return repo, nil
}

// GetRepositoryByName retrieves a repository by name.
func (rs *repositoryService) GetRepositoryByName(ctx context.Context, name string) (*ent.Repository, error) {
	repo, err := rs.client.GetEnt().Repository.Query().
		Where(repository.Name(name)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository by name: %w", err)
	}
	return repo, nil
}

// ListRepositories lists all repositories with optional filtering.
func (rs *repositoryService) ListRepositories(ctx context.Context, opts ListRepositoriesOptions) ([]*ent.Repository, error) {
	query := rs.client.GetEnt().Repository.Query()
	
	// Apply filters
	if opts.IsActive != nil {
		query = query.Where(repository.IsActive(*opts.IsActive))
	}
	
	if opts.Type != nil {
		query = query.Where(repository.TypeEQ(repository.Type(*opts.Type)))
	}
	
	// Apply ordering
	query = query.Order(ent.Desc(repository.FieldCreatedAt))
	
	// Apply pagination
	if opts.Limit > 0 {
		query = query.Limit(opts.Limit)
	}
	
	if opts.Offset > 0 {
		query = query.Offset(opts.Offset)
	}
	
	repos, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}
	
	return repos, nil
}

// UpdateRepository updates repository information.
func (rs *repositoryService) UpdateRepository(ctx context.Context, id string, req UpdateRepositoryRequest) (*ent.Repository, error) {
	log.Info("Updating repository", zap.String("id", id))
	
	update := rs.client.GetEnt().Repository.UpdateOneID(id)
	
	if req.LastSync != nil {
		update = update.SetNillableLastSync(req.LastSync)
	}
	
	if req.SyncStatus != nil {
		update = update.SetSyncStatus(repository.SyncStatus(*req.SyncStatus))
	}
	
	if req.Manifest != nil {
		manifestStr := fmt.Sprintf("%v", req.Manifest)
		update = update.SetNillableManifest(&manifestStr)
	}
	
	if req.ManifestHash != nil {
		update = update.SetNillableManifestHash(req.ManifestHash)
	}
	
	if req.IsActive != nil {
		update = update.SetIsActive(*req.IsActive)
	}
	
	repo, err := update.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update repository: %w", err)
	}
	
	log.Info("Repository updated", zap.String("id", id))
	return repo, nil
}

// DeleteRepository deletes a repository and all its components.
func (rs *repositoryService) DeleteRepository(ctx context.Context, id string) error {
	log.Info("Deleting repository", zap.String("id", id))
	
	err := rs.client.GetEnt().Repository.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}
	
	log.Info("Repository deleted", zap.String("id", id))
	return nil
}

