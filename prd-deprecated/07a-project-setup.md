# Project Setup & Architecture

## Overview

This document covers the initial project setup, Go module configuration, and architectural foundations for the MCP-Planner system.

## Project Initialization

### Go Module Setup

```bash
# Create project directory
mkdir mcp-planner
cd mcp-planner

# Initialize Go module
go mod init github.com/your-org/mcp-planner

# Create basic directory structure
mkdir -p {cmd/server,internal/{config,database,models,mcp,ai,complexity,disputes,utils},pkg/{types,errors},web,docs,scripts,migrations,tests,docker,.github/workflows}

# Create initial files
touch {cmd/server/main.go,internal/config/config.go,pkg/types/types.go,README.md,Makefile,.gitignore,.env.example}
```

### Directory Structure

```
mcp-planner/
├── cmd/
│   └── server/
│       ├── main.go              # Application entry point
│       └── wire.go              # Dependency injection (optional)
├── internal/
│   ├── config/
│   │   ├── config.go            # Configuration management
│   │   └── validation.go        # Config validation
│   ├── database/
│   │   ├── connection.go        # Database connection
│   │   ├── migrations.go        # Migration runner
│   │   └── health.go            # Database health checks
│   ├── models/
│   │   ├── project.go           # Project business logic
│   │   ├── task.go              # Task business logic
│   │   ├── step.go              # Step business logic
│   │   ├── dispute.go           # Dispute business logic
│   │   └── audit.go             # Audit logging
│   ├── mcp/
│   │   ├── server.go            # MCP server implementation
│   │   ├── handlers.go          # MCP function handlers
│   │   ├── transport.go         # Transport layer
│   │   └── middleware.go        # MCP middleware
│   ├── ai/
│   │   ├── providers/           # AI provider implementations
│   │   │   ├── openai.go
│   │   │   ├── anthropic.go
│   │   │   └── local.go
│   │   ├── collaboration.go     # Dual-AI workflow
│   │   ├── prompts.go           # Prompt templates
│   │   └── client.go            # AI client interface
│   ├── complexity/
│   │   ├── analyzer.go          # Complexity analysis
│   │   ├── promotion.go         # Step promotion logic
│   │   └── optimization.go      # Project optimization
│   ├── disputes/
│   │   ├── resolver.go          # Dispute resolution
│   │   ├── presentation.go      # UI data preparation
│   │   └── workflow.go          # Resolution workflow
│   └── utils/
│       ├── uuid.go              # UUID utilities
│       ├── validation.go        # Input validation
│       └── navigation.go        # Navigation helpers
├── pkg/
│   ├── types/
│   │   ├── project.go           # Project types
│   │   ├── task.go              # Task types
│   │   ├── step.go              # Step types
│   │   ├── ai.go                # AI-related types
│   │   └── mcp.go               # MCP types
│   └── errors/
│       ├── types.go             # Error type definitions
│       └── codes.go             # Error codes
├── web/                         # Frontend (optional)
├── docs/                        # Documentation
├── scripts/
│   ├── build.sh                 # Build script
│   ├── test.sh                  # Test script
│   └── deploy.sh                # Deployment script
├── migrations/                  # Database migrations
├── tests/
│   ├── integration/             # Integration tests
│   ├── unit/                    # Unit tests
│   └── fixtures/                # Test data
├── docker/
│   ├── Dockerfile               # Production image
│   ├── Dockerfile.dev           # Development image
│   └── docker-compose.yml       # Local development
├── .github/
│   └── workflows/
│       ├── ci.yml               # Continuous integration
│       └── release.yml          # Release automation
├── Makefile                     # Build automation
├── go.mod                       # Go dependencies
├── go.sum                       # Dependency checksums
├── .gitignore                   # Git ignore rules
├── .env.example                 # Environment template
└── README.md                    # Project documentation
```

## Dependencies

### Core Dependencies

```go
// go.mod
module github.com/your-org/mcp-planner

go 1.21

require (
    // Web framework
    github.com/gin-gonic/gin v1.9.1

    // Database
    github.com/lib/pq v1.10.9
    github.com/steebchen/prisma-client-go v0.25.0
    github.com/golang-migrate/migrate/v4 v4.16.2

    // Configuration
    github.com/spf13/viper v1.16.0
    github.com/spf13/cobra v1.7.0

    // Utilities
    github.com/google/uuid v1.3.0
    github.com/gorilla/websocket v1.5.0

    // Logging
    go.uber.org/zap v1.25.0

    // AI Providers
    github.com/sashabaranov/go-openai v1.15.3

    // Validation
    github.com/go-playground/validator/v10 v10.15.1

    // Testing
    github.com/stretchr/testify v1.8.4
    github.com/testcontainers/testcontainers-go v0.23.0
)

require (
    // Development tools
    github.com/air-verse/air v1.45.0 // Hot reload
    github.com/golangci/golangci-lint v1.54.2 // Linting
    github.com/securecodewarrior/gosec/v2 v2.18.2 // Security
)
```

### Installation Script

```bash
#!/bin/bash
# scripts/install-deps.sh

echo "Installing Go dependencies..."
go mod download
go mod verify

echo "Installing development tools..."
go install github.com/air-verse/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

echo "Installing Prisma CLI..."
npm install -g prisma

echo "Setting up pre-commit hooks..."
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
make lint
make test-unit
EOF
chmod +x .git/hooks/pre-commit

echo "Dependencies installed successfully!"
```

## Configuration Management

### Configuration Structure

```go
// internal/config/config.go
package config

import (
    "fmt"
    "time"

    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    AI       AIConfig       `mapstructure:"ai"`
    MCP      MCPConfig      `mapstructure:"mcp"`
    Logging  LoggingConfig  `mapstructure:"logging"`
    Security SecurityConfig `mapstructure:"security"`
}

type ServerConfig struct {
    Host         string        `mapstructure:"host" validate:"required"`
    Port         int           `mapstructure:"port" validate:"required,min=1,max=65535"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
    IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DatabaseConfig struct {
    URL             string        `mapstructure:"url" validate:"required"`
    MaxConnections  int           `mapstructure:"max_connections"`
    MaxIdleTime     time.Duration `mapstructure:"max_idle_time"`
    MaxLifetime     time.Duration `mapstructure:"max_lifetime"`
    MigrationsPath  string        `mapstructure:"migrations_path"`
}

type AIConfig struct {
    OpenAI    OpenAIConfig    `mapstructure:"openai"`
    Anthropic AnthropicConfig `mapstructure:"anthropic"`
    Local     LocalAIConfig   `mapstructure:"local"`
    Timeout   time.Duration   `mapstructure:"timeout"`
    Retries   int             `mapstructure:"retries"`
}

type OpenAIConfig struct {
    APIKey      string  `mapstructure:"api_key" validate:"required"`
    Model       string  `mapstructure:"model"`
    Temperature float32 `mapstructure:"temperature"`
    MaxTokens   int     `mapstructure:"max_tokens"`
}

type AnthropicConfig struct {
    APIKey      string  `mapstructure:"api_key"`
    Model       string  `mapstructure:"model"`
    Temperature float32 `mapstructure:"temperature"`
    MaxTokens   int     `mapstructure:"max_tokens"`
}

type LocalAIConfig struct {
    Endpoint string `mapstructure:"endpoint"`
    Model    string `mapstructure:"model"`
}

type MCPConfig struct {
    Transport      string        `mapstructure:"transport"`
    MaxConnections int           `mapstructure:"max_connections"`
    Timeout        time.Duration `mapstructure:"timeout"`
    BufferSize     int           `mapstructure:"buffer_size"`
}

type LoggingConfig struct {
    Level      string `mapstructure:"level"`
    Format     string `mapstructure:"format"`
    Output     string `mapstructure:"output"`
    MaxSize    int    `mapstructure:"max_size"`
    MaxBackups int    `mapstructure:"max_backups"`
    MaxAge     int    `mapstructure:"max_age"`
}

type SecurityConfig struct {
    APIKeys        []string      `mapstructure:"api_keys"`
    RateLimit      int           `mapstructure:"rate_limit"`
    RateLimitWindow time.Duration `mapstructure:"rate_limit_window"`
    CORS           CORSConfig    `mapstructure:"cors"`
}

type CORSConfig struct {
    AllowedOrigins []string `mapstructure:"allowed_origins"`
    AllowedMethods []string `mapstructure:"allowed_methods"`
    AllowedHeaders []string `mapstructure:"allowed_headers"`
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")
    viper.AddConfigPath("/etc/mcp-planner")

    // Environment variable support
    viper.AutomaticEnv()
    viper.SetEnvPrefix("MCP_TASK")

    // Set defaults
    setDefaults()

    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, fmt.Errorf("failed to read config file: %w", err)
        }
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func setDefaults() {
    // Server defaults
    viper.SetDefault("server.host", "0.0.0.0")
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("server.read_timeout", "30s")
    viper.SetDefault("server.write_timeout", "30s")
    viper.SetDefault("server.idle_timeout", "120s")

    // Database defaults
    viper.SetDefault("database.max_connections", 25)
    viper.SetDefault("database.max_idle_time", "15m")
    viper.SetDefault("database.max_lifetime", "1h")
    viper.SetDefault("database.migrations_path", "./migrations")

    // AI defaults
    viper.SetDefault("ai.timeout", "30s")
    viper.SetDefault("ai.retries", 3)
    viper.SetDefault("ai.openai.model", "gpt-4")
    viper.SetDefault("ai.openai.temperature", 0.3)
    viper.SetDefault("ai.openai.max_tokens", 2000)
    viper.SetDefault("ai.anthropic.model", "claude-3-sonnet-20240229")
    viper.SetDefault("ai.anthropic.temperature", 0.3)
    viper.SetDefault("ai.anthropic.max_tokens", 2000)

    // MCP defaults
    viper.SetDefault("mcp.transport", "stdio")
    viper.SetDefault("mcp.max_connections", 100)
    viper.SetDefault("mcp.timeout", "30s")
    viper.SetDefault("mcp.buffer_size", 1024)

    // Logging defaults
    viper.SetDefault("logging.level", "info")
    viper.SetDefault("logging.format", "json")
    viper.SetDefault("logging.output", "stdout")
    viper.SetDefault("logging.max_size", 100)
    viper.SetDefault("logging.max_backups", 3)
    viper.SetDefault("logging.max_age", 28)

    // Security defaults
    viper.SetDefault("security.rate_limit", 100)
    viper.SetDefault("security.rate_limit_window", "1m")
    viper.SetDefault("security.cors.allowed_origins", []string{"*"})
    viper.SetDefault("security.cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
    viper.SetDefault("security.cors.allowed_headers", []string{"*"})
}
```

### Configuration Validation

```go
// internal/config/validation.go
package config

import (
    "fmt"
    "net/url"
    "strings"

    "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
    validate = validator.New()

    // Register custom validators
    validate.RegisterValidation("database_url", validateDatabaseURL)
    validate.RegisterValidation("log_level", validateLogLevel)
}

func validateConfig(config *Config) error {
    if err := validate.Struct(config); err != nil {
        return formatValidationError(err)
    }

    // Custom validation logic
    if err := validateAIConfig(&config.AI); err != nil {
        return err
    }

    return nil
}

func validateDatabaseURL(fl validator.FieldLevel) bool {
    dbURL := fl.Field().String()
    if dbURL == "" {
        return false
    }

    u, err := url.Parse(dbURL)
    if err != nil {
        return false
    }

    return u.Scheme == "postgresql" || u.Scheme == "postgres"
}

func validateLogLevel(fl validator.FieldLevel) bool {
    level := strings.ToLower(fl.Field().String())
    validLevels := []string{"debug", "info", "warn", "error", "fatal"}

    for _, validLevel := range validLevels {
        if level == validLevel {
            return true
        }
    }

    return false
}

func validateAIConfig(ai *AIConfig) error {
    hasProvider := false

    if ai.OpenAI.APIKey != "" {
        hasProvider = true
    }

    if ai.Anthropic.APIKey != "" {
        hasProvider = true
    }

    if ai.Local.Endpoint != "" {
        hasProvider = true
    }

    if !hasProvider {
        return fmt.Errorf("at least one AI provider must be configured")
    }

    return nil
}

func formatValidationError(err error) error {
    var messages []string

    for _, err := range err.(validator.ValidationErrors) {
        switch err.Tag() {
        case "required":
            messages = append(messages, fmt.Sprintf("%s is required", err.Field()))
        case "min":
            messages = append(messages, fmt.Sprintf("%s must be at least %s", err.Field(), err.Param()))
        case "max":
            messages = append(messages, fmt.Sprintf("%s must be at most %s", err.Field(), err.Param()))
        case "database_url":
            messages = append(messages, fmt.Sprintf("%s must be a valid PostgreSQL URL", err.Field()))
        case "log_level":
            messages = append(messages, fmt.Sprintf("%s must be one of: debug, info, warn, error, fatal", err.Field()))
        default:
            messages = append(messages, fmt.Sprintf("%s is invalid", err.Field()))
        }
    }

    return fmt.Errorf("validation errors: %s", strings.Join(messages, ", "))
}
```

## Application Entry Point

### Main Application

```go
// cmd/server/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/your-org/mcp-planner/internal/config"
    "github.com/your-org/mcp-planner/internal/database"
    "github.com/your-org/mcp-planner/internal/mcp"
    "github.com/your-org/mcp-planner/pkg/errors"

    "go.uber.org/zap"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Initialize logger
    logger, err := initLogger(cfg.Logging)
    if err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
    defer logger.Sync()

    // Initialize database
    db, err := database.Connect(cfg.Database)
    if err != nil {
        logger.Fatal("Failed to connect to database", zap.Error(err))
    }
    defer db.Close()

    // Run migrations
    if err := database.Migrate(cfg.Database); err != nil {
        logger.Fatal("Failed to run migrations", zap.Error(err))
    }

    // Initialize MCP server
    mcpServer, err := mcp.NewServer(cfg, db, logger)
    if err != nil {
        logger.Fatal("Failed to create MCP server", zap.Error(err))
    }

    // Start server
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        if err := mcpServer.Start(ctx); err != nil {
            logger.Error("MCP server error", zap.Error(err))
            cancel()
        }
    }()

    // Wait for shutdown signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    select {
    case sig := <-sigChan:
        logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
    case <-ctx.Done():
        logger.Info("Context cancelled")
    }

    // Graceful shutdown
    logger.Info("Shutting down server...")
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer shutdownCancel()

    if err := mcpServer.Shutdown(shutdownCtx); err != nil {
        logger.Error("Error during shutdown", zap.Error(err))
    }

    logger.Info("Server shutdown complete")
}

func initLogger(cfg config.LoggingConfig) (*zap.Logger, error) {
    var zapConfig zap.Config

    switch cfg.Format {
    case "json":
        zapConfig = zap.NewProductionConfig()
    case "console":
        zapConfig = zap.NewDevelopmentConfig()
    default:
        return nil, fmt.Errorf("invalid log format: %s", cfg.Format)
    }

    // Set log level
    level, err := zap.ParseAtomicLevel(cfg.Level)
    if err != nil {
        return nil, fmt.Errorf("invalid log level: %s", cfg.Level)
    }
    zapConfig.Level = level

    // Configure output
    if cfg.Output != "stdout" {
        zapConfig.OutputPaths = []string{cfg.Output}
    }

    return zapConfig.Build()
}
```

## Build Automation

### Makefile

```makefile
# Makefile
.PHONY: help setup build test lint clean docker dev

# Variables
APP_NAME := mcp-planner
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Setup development environment
	@echo "Setting up development environment..."
	@./scripts/install-deps.sh
	@cp .env.example .env
	@echo "Setup complete! Edit .env file with your configuration."

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@go build $(LDFLAGS) -o bin/$(APP_NAME) cmd/server/main.go

build-linux: ## Build for Linux
	@echo "Building $(APP_NAME) for Linux..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME)-linux cmd/server/main.go

test: ## Run all tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	@go test -v -race -short ./...

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	@go test -v -race -run Integration ./...

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

security: ## Run security scan
	@echo "Running security scan..."
	@gosec ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

docker: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):$(VERSION) -f docker/Dockerfile .

docker-dev: ## Build development Docker image
	@echo "Building development Docker image..."
	@docker build -t $(APP_NAME):dev -f docker/Dockerfile.dev .

dev: ## Start development environment
	@echo "Starting development environment..."
	@docker-compose -f docker/docker-compose.yml up -d
	@air

dev-stop: ## Stop development environment
	@echo "Stopping development environment..."
	@docker-compose -f docker/docker-compose.yml down

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	@migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	@migrate -path migrations -database "$(DATABASE_URL)" down

migrate-create: ## Create new migration (usage: make migrate-create NAME=migration_name)
	@echo "Creating migration: $(NAME)"
	@migrate create -ext sql -dir migrations $(NAME)

generate: ## Generate code (Prisma, etc.)
	@echo "Generating code..."
	@prisma generate

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

update-deps: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

install: build ## Install the application
	@echo "Installing $(APP_NAME)..."
	@cp bin/$(APP_NAME) /usr/local/bin/

uninstall: ## Uninstall the application
	@echo "Uninstalling $(APP_NAME)..."
	@rm -f /usr/local/bin/$(APP_NAME)
```

## Environment Configuration

### Environment Template

```bash
# .env.example

# Database Configuration
DATABASE_URL=postgresql://mcp_task:password@localhost:5432/mcp_task?sslmode=disable
DATABASE_MAX_CONNECTIONS=25
DATABASE_MAX_IDLE_TIME=15m
DATABASE_MAX_LIFETIME=1h

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_IDLE_TIMEOUT=120s

# AI Provider Configuration
OPENAI_API_KEY=sk-your-openai-api-key
ANTHROPIC_API_KEY=sk-ant-your-anthropic-api-key
AI_TIMEOUT=30s
AI_RETRIES=3

# OpenAI Specific
OPENAI_MODEL=gpt-4
OPENAI_TEMPERATURE=0.3
OPENAI_MAX_TOKENS=2000

# Anthropic Specific
ANTHROPIC_MODEL=claude-3-sonnet-20240229
ANTHROPIC_TEMPERATURE=0.3
ANTHROPIC_MAX_TOKENS=2000

# Local AI (optional)
LOCAL_AI_ENDPOINT=http://localhost:11434
LOCAL_AI_MODEL=llama2

# MCP Configuration
MCP_TRANSPORT=stdio
MCP_MAX_CONNECTIONS=100
MCP_TIMEOUT=30s
MCP_BUFFER_SIZE=1024

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout
LOG_MAX_SIZE=100
LOG_MAX_BACKUPS=3
LOG_MAX_AGE=28

# Security Configuration
API_KEYS=key1,key2,key3
RATE_LIMIT=100
RATE_LIMIT_WINDOW=1m

# CORS Configuration
CORS_ALLOWED_ORIGINS=*
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=*
```

### Git Configuration

```gitignore
# .gitignore

# Binaries
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out
coverage.html

# Go workspace file
go.work

# Environment files
.env
.env.local
.env.production

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS files
.DS_Store
Thumbs.db

# Logs
*.log
logs/

# Database files
*.db
*.sqlite

# Temporary files
tmp/
temp/

# Node modules (if using frontend)
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Docker
.dockerignore

# Build artifacts
dist/
build/
```

---

*Next: [Database Implementation](./07b-database-implementation.md)*
