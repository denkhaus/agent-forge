package git

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"go.uber.org/zap"
)

// Client provides Git operations for component management.
type Client struct {
	logger *zap.Logger
}

// NewClient creates a new Git client.
func NewClient(logger *zap.Logger) *Client {
	return &Client{
		// TODO: get rid of this logger, we have a Package logger!
		logger: logger,
	}
}

// CloneOptions contains options for cloning a repository.
type CloneOptions struct {
	URL         string
	Destination string
	Branch      string
	Depth       int
	Shallow     bool
}

// PushOptions contains options for pushing to a repository.
type PushOptions struct {
	Repository string
	Branch     string
	Remote     string
}

// Clone clones a Git repository to the specified destination.
func (c *Client) Clone(ctx context.Context, opts CloneOptions) error {
	c.logger.Info("Cloning repository",
		zap.String("url", opts.URL),
		zap.String("destination", opts.Destination),
		zap.String("branch", opts.Branch))

	// Prepare clone options
	cloneOpts := &git.CloneOptions{
		URL: opts.URL,
	}

	// Add branch if specified
	if opts.Branch != "" {
		cloneOpts.ReferenceName = plumbing.ReferenceName("refs/heads/" + opts.Branch)
		cloneOpts.SingleBranch = true
	}

	// Add shallow clone if requested
	if opts.Shallow || opts.Depth > 0 {
		depth := opts.Depth
		if depth == 0 {
			depth = 1
		}
		cloneOpts.Depth = depth
	}

	// Clone the repository
	_, err := git.PlainCloneContext(ctx, opts.Destination, false, cloneOpts)
	if err != nil {
		return fmt.Errorf("failed to clone repository %s: %w", opts.URL, err)
	}

	c.logger.Info("Repository cloned successfully",
		zap.String("url", opts.URL),
		zap.String("destination", opts.Destination))

	return nil
}

// Push pushes changes to a remote repository.
func (c *Client) Push(ctx context.Context, repoPath string, opts PushOptions) error {
	c.logger.Info("Pushing to repository",
		zap.String("path", repoPath),
		zap.String("remote", opts.Remote),
		zap.String("branch", opts.Branch))

	// Set default values
	remote := opts.Remote
	if remote == "" {
		remote = "origin"
	}

	// TODO: Implement push with authentication
	// For now, we'll just log that push is prepared
	c.logger.Info("Repository prepared for push",
		zap.String("remote", remote),
		zap.String("path", repoPath))

	return nil
}

// InitRepository initializes a new Git repository.
func (c *Client) InitRepository(ctx context.Context, path string) error {
	c.logger.Info("Initializing Git repository", zap.String("path", path))

	// Initialize the repository
	_, err := git.PlainInit(path, false)
	if err != nil {
		return fmt.Errorf("failed to initialize Git repository: %w", err)
	}

	c.logger.Info("Git repository initialized", zap.String("path", path))
	return nil
}

// AddAndCommit adds all changes and commits them.
func (c *Client) AddAndCommit(ctx context.Context, repoPath, message string) error {
	c.logger.Info("Adding and committing changes",
		zap.String("path", repoPath),
		zap.String("message", message))

	// Open the repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the working tree
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Add all changes
	_, err = worktree.Add(".")
	if err != nil {
		return fmt.Errorf("failed to add changes: %w", err)
	}

	// Create commit
	commit, err := worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "AgentForge",
			Email: "agentforge@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	c.logger.Info("Changes committed successfully", 
		zap.String("message", message),
		zap.String("commit", commit.String()))
	return nil
}

// AddRemote adds a remote repository.
func (c *Client) AddRemote(ctx context.Context, repoPath, name, url string) error {
	c.logger.Info("Adding remote",
		zap.String("path", repoPath),
		zap.String("name", name),
		zap.String("url", url))

	// Open the repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	// Create remote config
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{url},
	})
	if err != nil {
		return fmt.Errorf("failed to add remote %s: %w", name, err)
	}

	c.logger.Info("Remote added successfully",
		zap.String("name", name),
		zap.String("url", url))

	return nil
}

// ParseRepositoryURL parses a repository URL and extracts useful information.
func ParseRepositoryURL(url string) (owner, repo string, err error) {
	// Handle GitHub URLs
	if strings.Contains(url, "github.com") {
		// Remove .git suffix if present
		url = strings.TrimSuffix(url, ".git")

		// Handle both HTTPS and SSH formats
		if strings.HasPrefix(url, "https://github.com/") {
			parts := strings.Split(strings.TrimPrefix(url, "https://github.com/"), "/")
			if len(parts) >= 2 {
				return parts[0], parts[1], nil
			}
		} else if strings.HasPrefix(url, "git@github.com:") {
			parts := strings.Split(strings.TrimPrefix(url, "git@github.com:"), "/")
			if len(parts) >= 2 {
				return parts[0], parts[1], nil
			}
		}
	}

	return "", "", fmt.Errorf("unsupported repository URL format: %s", url)
}

// GetComponentPath returns the local path for a component based on its repository info.
func GetComponentPath(cacheDir, owner, repo, componentName string) string {
	return filepath.Join(cacheDir, "components", owner, repo, componentName)
}

// GetRepositoryPath returns the local path for a repository.
func GetRepositoryPath(cacheDir, owner, repo string) string {
	return filepath.Join(cacheDir, "repositories", owner, repo)
}
