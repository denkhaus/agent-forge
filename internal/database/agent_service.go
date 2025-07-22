package database

import (
	"context"
	"fmt"

	"github.com/denkhaus/agentforge/internal/database/ent"
	"github.com/denkhaus/agentforge/internal/database/ent/agent"
	"github.com/denkhaus/agentforge/internal/database/ent/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// agentService provides agent management operations (private implementation)
type agentService struct {
	client DatabaseClient
}

// NewAgentService creates a new agent service
func NewAgentService(client DatabaseClient) AgentService {
	return &agentService{
		client: client,
	}
}

// CreateAgent creates a new agent record
func (as *agentService) CreateAgent(ctx context.Context, req CreateAgentRequest) (*ent.Agent, error) {
	log.Info("Creating agent", 
		zap.String("name", req.Name),
		zap.String("version", req.Version))
	
	create := as.client.GetEnt().Agent.Create().
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
		SetStability(agent.Stability(req.Stability)).
		SetMaturity(agent.Maturity(req.Maturity)).
		SetForgeVersion(req.ForgeVersion).
		SetPlatforms(req.Platforms).
		SetSpec(req.Spec).
		SetSpecHash(req.SpecHash).
		SetCommitHash(req.CommitHash).
		SetBranch(req.Branch).
		SetNillableConfigPath(req.ConfigPath).
		SetNillableLlmProvider(req.LLMProvider).
		SetNillableSystemPromptID(req.SystemPromptID).
		SetToolDependencies(req.ToolDependencies).
		SetPromptDependencies(req.PromptDependencies).
		SetAgentDependencies(req.AgentDependencies).
		SetAgentType(agent.AgentType(req.AgentType)).
		SetCapabilities(req.Capabilities).
		SetSupportedLanguages(req.SupportedLanguages).
		SetSupportsMemory(req.SupportsMemory).
		SetSupportsTools(req.SupportsTools).
		SetSupportsMultimodal(req.SupportsMultimodal).
		SetModelPreferences(req.ModelPreferences).
		SetNillableDefaultTemperature(req.DefaultTemperature).
		SetNillableDefaultMaxTokens(req.DefaultMaxTokens).
		SetSessionTimeoutMinutes(req.SessionTimeoutMinutes)
	
	agent, err := create.Save(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}
	
	log.Info("Agent created", zap.String("id", agent.ID))
	return agent, nil
}

// GetAgent retrieves an agent by ID
func (as *agentService) GetAgent(ctx context.Context, id string) (*ent.Agent, error) {
	agent, err := as.client.GetEnt().Agent.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}
	return agent, nil
}

// GetAgentByName retrieves an agent by name, version, and repository
func (as *agentService) GetAgentByName(ctx context.Context, name, version string, repositoryID string) (*ent.Agent, error) {
	query := as.client.GetEnt().Agent.Query().
		Where(agent.Name(name))
	
	if version != "" {
		query = query.Where(agent.Version(version))
	}
	if repositoryID != "" {
		query = query.Where(agent.HasRepositoryWith(repository.ID(repositoryID)))
	}
	
	agent, err := query.Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent by name: %w", err)
	}
	return agent, nil
}

// ListAgents lists all agents with optional filtering
func (as *agentService) ListAgents(ctx context.Context, opts ListAgentsOptions) ([]*ent.Agent, error) {
	query := as.client.GetEnt().Agent.Query()
	
	// Apply filters
	if opts.AgentType != nil {
		query = query.Where(agent.AgentTypeEQ(agent.AgentType(*opts.AgentType)))
	}
	if opts.LLMProvider != nil {
		query = query.Where(agent.LlmProviderEQ(*opts.LLMProvider))
	}
	if opts.Stability != nil {
		query = query.Where(agent.StabilityEQ(agent.Stability(*opts.Stability)))
	}
	if opts.IsInstalled != nil {
		query = query.Where(agent.IsInstalled(*opts.IsInstalled))
	}
	if opts.RepositoryID != nil {
		query = query.Where(agent.HasRepositoryWith(repository.ID(*opts.RepositoryID)))
	}
	if opts.SupportsTools != nil {
		query = query.Where(agent.SupportsTools(*opts.SupportsTools))
	}
	if opts.SupportsMemory != nil {
		query = query.Where(agent.SupportsMemory(*opts.SupportsMemory))
	}
	
	// Apply ordering
	query = query.Order(ent.Desc(agent.FieldCreatedAt))
	
	// Apply pagination
	if opts.Limit > 0 {
		query = query.Limit(opts.Limit)
	}
	if opts.Offset > 0 {
		query = query.Offset(opts.Offset)
	}
	
	agents, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list agents: %w", err)
	}
	
	return agents, nil
}

// UpdateAgent updates agent information
func (as *agentService) UpdateAgent(ctx context.Context, id string, updates map[string]interface{}) (*ent.Agent, error) {
	log.Info("Updating agent", zap.String("id", id))
	
	update := as.client.GetEnt().Agent.UpdateOneID(id)
	
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
	
	agent, err := update.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update agent: %w", err)
	}
	
	log.Info("Agent updated", zap.String("id", id))
	return agent, nil
}

// DeleteAgent deletes an agent
func (as *agentService) DeleteAgent(ctx context.Context, id string) error {
	log.Info("Deleting agent", zap.String("id", id))
	
	err := as.client.GetEnt().Agent.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete agent: %w", err)
	}
	
	log.Info("Agent deleted", zap.String("id", id))
	return nil
}

// InstallAgent marks an agent as installed
func (as *agentService) InstallAgent(ctx context.Context, id string, installPath string) (*ent.Agent, error) {
	return as.UpdateAgent(ctx, id, map[string]interface{}{
		"is_installed": true,
		"install_path": &installPath,
	})
}

// UninstallAgent marks an agent as uninstalled
func (as *agentService) UninstallAgent(ctx context.Context, id string) (*ent.Agent, error) {
	return as.UpdateAgent(ctx, id, map[string]interface{}{
		"is_installed": false,
		"install_path": (*string)(nil),
	})
}

// SearchAgents searches for agents based on query
func (as *agentService) SearchAgents(ctx context.Context, query string, opts SearchAgentsOptions) ([]*ent.Agent, error) {
	dbQuery := as.client.GetEnt().Agent.Query().
		Where(agent.Or(
			agent.NameContains(query),
			agent.DescriptionContains(query),
		))
	
	// Apply filters
	if opts.AgentType != nil {
		dbQuery = dbQuery.Where(agent.AgentTypeEQ(agent.AgentType(*opts.AgentType)))
	}
	if opts.LLMProvider != nil {
		dbQuery = dbQuery.Where(agent.LlmProviderEQ(*opts.LLMProvider))
	}
	if opts.Stability != nil {
		dbQuery = dbQuery.Where(agent.StabilityEQ(agent.Stability(*opts.Stability)))
	}
	if opts.SupportsTools != nil {
		dbQuery = dbQuery.Where(agent.SupportsTools(*opts.SupportsTools))
	}
	
	// Apply limit
	if opts.Limit > 0 {
		dbQuery = dbQuery.Limit(opts.Limit)
	}
	
	agents, err := dbQuery.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to search agents: %w", err)
	}
	
	return agents, nil
}

// CLI-specific methods (placeholder implementations)
func (as *agentService) ListAgentsForCLI(ctx context.Context, opts interface{}) ([]interface{}, error) {
	// TODO: Implement CLI-specific agent listing
	return []interface{}{}, fmt.Errorf("CLI agent listing not yet implemented")
}

func (as *agentService) CreateAgentFiles(ctx context.Context, name string) error {
	// TODO: Implement agent file creation
	return fmt.Errorf("agent file creation not yet implemented")
}

func (as *agentService) PullAgent(ctx context.Context, repo string) error {
	// TODO: Implement agent pulling from repository
	return fmt.Errorf("agent pulling not yet implemented")
}

func (as *agentService) PushAgent(ctx context.Context, name string) error {
	// TODO: Implement agent pushing to repository
	return fmt.Errorf("agent pushing not yet implemented")
}