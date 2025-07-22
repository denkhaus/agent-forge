// Package git provides Git client interfaces for dependency injection.
package git

import (
	"context"

	"go.uber.org/zap"
)

// GitClient defines the interface for Git operations (matching actual implementation)
type GitClient interface {
	Clone(ctx context.Context, opts CloneOptions) error
	Push(ctx context.Context, repoPath string, opts PushOptions) error
	InitRepository(ctx context.Context, path string) error
	AddAndCommit(ctx context.Context, repoPath, message string) error
	AddRemote(ctx context.Context, repoPath, name, url string) error
}

// NewGitClientFromDI creates a GitClient using dependency injection
func NewGitClientFromDI(logger *zap.Logger) GitClient {
	return NewClient(logger)
}

// Ensure Client implements GitClient interface
var _ GitClient = (*Client)(nil)