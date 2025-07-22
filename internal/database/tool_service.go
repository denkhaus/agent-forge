package database

import (
	"context"
	"fmt"

	"github.com/denkhaus/agentforge/internal/database/ent"
	"github.com/denkhaus/agentforge/internal/database/ent/repository"
	"github.com/denkhaus/agentforge/internal/database/ent/tool"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// toolService provides tool management operations (private implementation)
type toolService struct {
	client DatabaseClient
}

// NewToolService creates a new tool service
func NewToolService(client DatabaseClient) ToolService {
	return &toolService{
		client: client,
	}
}

// CreateTool creates a new tool record
func (ts *toolService) CreateTool(ctx context.Context, req CreateToolRequest) (*ent.Tool, error) {
	log.Info("Creating tool", 
		zap.String("name", req.Name),
		zap.String("version", req.Version))
	
	tool, err := ts.client.GetEnt().Tool.Create().
		SetID(uuid.New().String()).
		SetName(req.Name).
		SetNamespace(req.Namespace).
		SetVersion(req.Version).
		SetDescription(req.Description).
		SetAuthor(req.Author).
		SetLicense(req.License).
		SetNillableHomepage(req.Homepage).
		SetNillableDocumentation(req.Documentation).
		SetTags(req.Tags).
		SetCategories(req.Categories).
		SetKeywords(req.Keywords).
		SetStability(tool.Stability(req.Stability)).
		SetMaturity(tool.Maturity(req.Maturity)).
		SetForgeVersion(req.ForgeVersion).
		SetPlatforms(req.Platforms).
		SetSpec(req.Spec).
		SetSpecHash(req.SpecHash).
		// Repository relationship handled via edge
		SetCommitHash(req.CommitHash).
		SetBranch(req.Branch).
		SetExecutionType(tool.ExecutionType(req.ExecutionType)).
		SetNillableSchemaPath(req.SchemaPath).
		SetCapabilities(req.Capabilities).
		SetNillableEntryPoint(req.EntryPoint).
		SetTimeoutSeconds(req.TimeoutSeconds).
		SetSupportsStreaming(req.SupportsStreaming).
		Save(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create tool: %w", err)
	}
	
	log.Info("Tool created", zap.String("id", tool.ID))
	return tool, nil
}

// GetTool retrieves a tool by ID
func (ts *toolService) GetTool(ctx context.Context, id string) (*ent.Tool, error) {
	tool, err := ts.client.GetEnt().Tool.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tool: %w", err)
	}
	return tool, nil
}

// GetToolByName retrieves a tool by name, version, and repository
func (ts *toolService) GetToolByName(ctx context.Context, name, version string, repositoryID string) (*ent.Tool, error) {
	query := ts.client.GetEnt().Tool.Query().
		Where(tool.Name(name))
	
	if version != "" {
		query = query.Where(tool.Version(version))
	}
	if repositoryID != "" {
		query = query.Where(tool.HasRepositoryWith(repository.ID(repositoryID)))
	}
	
	tool, err := query.Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tool by name: %w", err)
	}
	return tool, nil
}

// ListTools lists all tools with optional filtering
func (ts *toolService) ListTools(ctx context.Context, opts ListToolsOptions) ([]*ent.Tool, error) {
	query := ts.client.GetEnt().Tool.Query()
	
	// Apply filters
	if opts.ExecutionType != nil {
		query = query.Where(tool.ExecutionTypeEQ(tool.ExecutionType(*opts.ExecutionType)))
	}
	if opts.Stability != nil {
		query = query.Where(tool.StabilityEQ(tool.Stability(*opts.Stability)))
	}
	if opts.IsInstalled != nil {
		query = query.Where(tool.IsInstalled(*opts.IsInstalled))
	}
	if opts.RepositoryID != nil {
		query = query.Where(tool.HasRepositoryWith(repository.ID(*opts.RepositoryID)))
	}
	if opts.SupportsStreaming != nil {
		query = query.Where(tool.SupportsStreaming(*opts.SupportsStreaming))
	}
	
	// Apply ordering
	query = query.Order(ent.Desc(tool.FieldCreatedAt))
	
	// Apply pagination
	if opts.Limit > 0 {
		query = query.Limit(opts.Limit)
	}
	if opts.Offset > 0 {
		query = query.Offset(opts.Offset)
	}
	
	tools, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}
	
	return tools, nil
}

// UpdateTool updates tool information
func (ts *toolService) UpdateTool(ctx context.Context, id string, updates map[string]interface{}) (*ent.Tool, error) {
	log.Info("Updating tool", zap.String("id", id))
	
	update := ts.client.GetEnt().Tool.UpdateOneID(id)
	
	// Apply updates based on the map
	for key, value := range updates {
		switch key {
		case "description":
			if desc, ok := value.(string); ok {
				update = update.SetDescription(desc)
			}
		case "is_installed":
			if installed, ok := value.(bool); ok {
				update = update.SetIsInstalled(installed)
			}
		case "install_path":
			if path, ok := value.(*string); ok {
				update = update.SetNillableInstallPath(path)
			}
		}
	}
	
	tool, err := update.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update tool: %w", err)
	}
	
	log.Info("Tool updated", zap.String("id", id))
	return tool, nil
}

// DeleteTool deletes a tool
func (ts *toolService) DeleteTool(ctx context.Context, id string) error {
	log.Info("Deleting tool", zap.String("id", id))
	
	err := ts.client.GetEnt().Tool.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete tool: %w", err)
	}
	
	log.Info("Tool deleted", zap.String("id", id))
	return nil
}

// InstallTool marks a tool as installed
func (ts *toolService) InstallTool(ctx context.Context, id string, installPath string) (*ent.Tool, error) {
	return ts.UpdateTool(ctx, id, map[string]interface{}{
		"is_installed": true,
		"install_path": &installPath,
	})
}

// UninstallTool marks a tool as uninstalled
func (ts *toolService) UninstallTool(ctx context.Context, id string) (*ent.Tool, error) {
	return ts.UpdateTool(ctx, id, map[string]interface{}{
		"is_installed": false,
		"install_path": (*string)(nil),
	})
}

// SearchTools searches for tools based on query
func (ts *toolService) SearchTools(ctx context.Context, query string, opts SearchToolsOptions) ([]*ent.Tool, error) {
	dbQuery := ts.client.GetEnt().Tool.Query().
		Where(tool.Or(
			tool.NameContains(query),
			tool.DescriptionContains(query),
		))
	
	// Apply filters
	if opts.ExecutionType != nil {
		dbQuery = dbQuery.Where(tool.ExecutionTypeEQ(tool.ExecutionType(*opts.ExecutionType)))
	}
	if opts.Stability != nil {
		dbQuery = dbQuery.Where(tool.StabilityEQ(tool.Stability(*opts.Stability)))
	}
	
	// Apply limit
	if opts.Limit > 0 {
		dbQuery = dbQuery.Limit(opts.Limit)
	}
	
	tools, err := dbQuery.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to search tools: %w", err)
	}
	
	return tools, nil
}

// CLI-specific methods (placeholder implementations)
func (ts *toolService) ListToolsForCLI(ctx context.Context, opts interface{}) ([]interface{}, error) {
	// TODO: Implement CLI-specific tool listing
	return []interface{}{}, fmt.Errorf("CLI tool listing not yet implemented")
}

func (ts *toolService) CreateToolFiles(ctx context.Context, name string) error {
	// TODO: Implement tool file creation
	return fmt.Errorf("tool file creation not yet implemented")
}

func (ts *toolService) PullTool(ctx context.Context, repo string) error {
	// TODO: Implement tool pulling from repository
	return fmt.Errorf("tool pulling not yet implemented")
}

func (ts *toolService) PushTool(ctx context.Context, name string) error {
	// TODO: Implement tool pushing to repository
	return fmt.Errorf("tool pushing not yet implemented")
}