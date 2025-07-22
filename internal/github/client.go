package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v57/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/denkhaus/agentforge/internal/types"
)


// client is a private implementation of types.GitHubClient interface.
type client struct {
	gh *github.Client
}

// NewClient creates a new GitHub API client.
func NewClient(token string) types.GitHubClient {
	var tc *github.Client
	
	if token != "" {
		// Create OAuth2 token source
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = github.NewClient(oauth2.NewClient(context.Background(), ts))
	} else {
		// Create unauthenticated client (rate limited)
		tc = github.NewClient(nil)
	}
	
	return &client{
		gh: tc,
	}
}

// SearchRepositories searches for repositories containing AgentForge components.
func (c *client) SearchRepositories(ctx context.Context, query string, opts *types.GitHubSearchOptions) ([]*types.GitHubRepository, error) {
	log.Info("Searching GitHub repositories", zap.String("query", query))

	// Build search query with AgentForge-specific terms
	searchQuery := fmt.Sprintf("%s agentforge OR \"agent forge\" OR mcp OR \"model context protocol\"", query)
	
	if opts != nil {
		if opts.Language != "" {
			searchQuery += fmt.Sprintf(" language:%s", opts.Language)
		}
		if opts.Topic != "" {
			searchQuery += fmt.Sprintf(" topic:%s", opts.Topic)
		}
		if opts.MinStars > 0 {
			searchQuery += fmt.Sprintf(" stars:>=%d", opts.MinStars)
		}
		if opts.MaxStars > 0 {
			searchQuery += fmt.Sprintf(" stars:<=%d", opts.MaxStars)
		}
		if opts.CreatedAt != "" {
			searchQuery += fmt.Sprintf(" created:%s", opts.CreatedAt)
		}
		if opts.UpdatedAt != "" {
			searchQuery += fmt.Sprintf(" pushed:%s", opts.UpdatedAt)
		}
	}

	// Set up search options
	searchOpts := &github.SearchOptions{}
	if opts != nil {
		if opts.Sort != "" {
			searchOpts.Sort = opts.Sort
		}
		if opts.Order != "" {
			searchOpts.Order = opts.Order
		}
		searchOpts.ListOptions = github.ListOptions{
			Page:    opts.Page,
			PerPage: opts.PerPage,
		}
	}

	// Search repositories using official client
	result, _, err := c.gh.Search.Repositories(ctx, searchQuery, searchOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to search repositories: %w", err)
	}

	// Convert to our types
	repositories := make([]*types.GitHubRepository, len(result.Repositories))
	for i, repo := range result.Repositories {
		repositories[i] = &types.GitHubRepository{
			ID:          repo.GetID(),
			Name:        repo.GetName(),
			FullName:    repo.GetFullName(),
			Description: repo.GetDescription(),
			HTMLURL:     repo.GetHTMLURL(),
			CloneURL:    repo.GetCloneURL(),
			StarCount:   repo.GetStargazersCount(),
			Language:    repo.GetLanguage(),
			UpdatedAt:   repo.GetUpdatedAt().Time,
			Topics:      repo.Topics,
		}
	}

	log.Info("Found repositories", zap.Int("count", len(repositories)))
	return repositories, nil
}

// GetRepository gets detailed information about a specific repository.
func (c *client) GetRepository(ctx context.Context, owner, repo string) (*types.GitHubRepository, error) {
	log.Info("Getting repository details", zap.String("owner", owner), zap.String("repo", repo))

	// Get repository using official client
	repository, _, err := c.gh.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	// Convert to our type
	return &types.GitHubRepository{
		ID:          repository.GetID(),
		Name:        repository.GetName(),
		FullName:    repository.GetFullName(),
		Description: repository.GetDescription(),
		HTMLURL:     repository.GetHTMLURL(),
		CloneURL:    repository.GetCloneURL(),
		StarCount:   repository.GetStargazersCount(),
		Language:    repository.GetLanguage(),
		UpdatedAt:   repository.GetUpdatedAt().Time,
		Topics:      repository.Topics,
	}, nil
}

// ListComponents lists all AgentForge components in a repository.
func (c *client) ListComponents(ctx context.Context, owner, repo string) ([]*types.GitHubComponent, error) {
	log.Info("Listing components in repository", zap.String("owner", owner), zap.String("repo", repo))

	// Search for component files (YAML files that might be components)
	query := fmt.Sprintf("repo:%s/%s filename:*.yaml OR filename:*.yml", owner, repo)
	
	searchOpts := &github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	result, _, err := c.gh.Search.Code(ctx, query, searchOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to search for component files: %w", err)
	}

	var components []*types.GitHubComponent
	
	// Get repository info for metadata
	repoInfo, _, err := c.gh.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository info: %w", err)
	}

	for _, file := range result.CodeResults {
		// Basic component detection - could be enhanced with actual file parsing
		if strings.Contains(strings.ToLower(file.GetName()), "component") ||
		   strings.Contains(strings.ToLower(file.GetPath()), "agent") ||
		   strings.Contains(strings.ToLower(file.GetPath()), "tool") ||
		   strings.Contains(strings.ToLower(file.GetPath()), "prompt") {
			
			component := &types.GitHubComponent{
				Name:        file.GetName(),
				Repository:  repoInfo.GetFullName(),
				Path:        file.GetPath(),
				Stars:       repoInfo.GetStargazersCount(),
				UpdatedAt:   repoInfo.GetUpdatedAt().Time,
				Topics:      repoInfo.Topics,
			}
			
			// Determine component type from path/name
			path := strings.ToLower(file.GetPath())
			if strings.Contains(path, "tool") {
				component.Type = "tool"
			} else if strings.Contains(path, "prompt") {
				component.Type = "prompt"
			} else if strings.Contains(path, "agent") {
				component.Type = "agent"
			} else {
				component.Type = "unknown"
			}
			
			components = append(components, component)
		}
	}

	log.Info("Found components", zap.Int("count", len(components)))
	return components, nil
}

// GetComponentContent gets the raw content of a component file.
func (c *client) GetComponentContent(ctx context.Context, owner, repo, path string) ([]byte, error) {
	log.Info("Getting component content", 
		zap.String("owner", owner), 
		zap.String("repo", repo), 
		zap.String("path", path))

	// Get file content using official client
	fileContent, _, _, err := c.gh.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get file content: %w", err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return nil, fmt.Errorf("failed to decode file content: %w", err)
	}

	return []byte(content), nil
}