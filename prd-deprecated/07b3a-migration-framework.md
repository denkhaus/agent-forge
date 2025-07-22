# Migration Framework & Setup

## Overview

This document covers the migration framework setup, tooling configuration, and execution infrastructure for the MCP-Planner database migrations using golang-migrate and custom migration utilities.

## Migration Framework Setup

### golang-migrate Installation

```bash
# Install golang-migrate CLI
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Verify installation
migrate -version

# Alternative: Install via Go
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Project Structure

```
migrations/
├── 001_create_projects_table.up.sql
├── 001_create_projects_table.down.sql
├── 002_create_tasks_table.up.sql
├── 002_create_tasks_table.down.sql
├── 003_create_steps_table.up.sql
├── 003_create_steps_table.down.sql
├── 004_create_disputes_table.up.sql
├── 004_create_disputes_table.down.sql
├── 005_create_audit_logs_table.up.sql
├── 005_create_audit_logs_table.down.sql
├── 006_add_indexes.up.sql
├── 006_add_indexes.down.sql
├── 007_add_constraints.up.sql
├── 007_add_constraints.down.sql
├── 008_create_functions.up.sql
├── 008_create_functions.down.sql
├── 009_create_triggers.up.sql
├── 009_create_triggers.down.sql
└── 010_seed_data.up.sql
    └── 010_seed_data.down.sql
```

### Migration Manager

```go
// internal/database/migrations.go
package database

import (
    "context"
    "database/sql"
    "fmt"
    "io/fs"
    "path/filepath"
    "time"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"

    "go.uber.org/zap"
)

type MigrationManager struct {
    db            *sql.DB
    migrate       *migrate.Migrate
    migrationsDir string
    logger        *zap.Logger
}

type MigrationConfig struct {
    DatabaseURL     string
    MigrationsDir   string
    LockTimeout     time.Duration
    StatementTimeout time.Duration
}

func NewMigrationManager(config MigrationConfig, logger *zap.Logger) (*MigrationManager, error) {
    // Open database connection
    db, err := sql.Open("postgres", config.DatabaseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Configure connection for migrations
    db.SetMaxOpenConns(1) // Single connection for migrations
    db.SetMaxIdleConns(1)
    db.SetConnMaxLifetime(time.Hour)

    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    // Create postgres driver instance
    driver, err := postgres.WithInstance(db, &postgres.Config{
        MigrationsTable:  "schema_migrations",
        DatabaseName:     "mcp_task",
        StatementTimeout: config.StatementTimeout,
        LockTimeout:      config.LockTimeout,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create postgres driver: %w", err)
    }

    // Create file source
    sourceURL := fmt.Sprintf("file://%s", config.MigrationsDir)
    m, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
    if err != nil {
        return nil, fmt.Errorf("failed to create migrate instance: %w", err)
    }

    return &MigrationManager{
        db:            db,
        migrate:       m,
        migrationsDir: config.MigrationsDir,
        logger:        logger,
    }, nil
}

func (mm *MigrationManager) Up(ctx context.Context) error {
    mm.logger.Info("Running database migrations...")

    // Get current version
    currentVersion, dirty, err := mm.migrate.Version()
    if err != nil && err != migrate.ErrNilVersion {
        return fmt.Errorf("failed to get current version: %w", err)
    }

    if dirty {
        return fmt.Errorf("database is in dirty state at version %d", currentVersion)
    }

    mm.logger.Info("Current migration version", zap.Uint("version", currentVersion))

    // Run migrations
    if err := mm.migrate.Up(); err != nil {
        if err == migrate.ErrNoChange {
            mm.logger.Info("No new migrations to apply")
            return nil
        }
        return fmt.Errorf("failed to run migrations: %w", err)
    }

    // Get new version
    newVersion, _, err := mm.migrate.Version()
    if err != nil {
        return fmt.Errorf("failed to get new version: %w", err)
    }

    mm.logger.Info("Migrations completed successfully", zap.Uint("new_version", newVersion))
    return nil
}

func (mm *MigrationManager) Down(ctx context.Context, steps int) error {
    mm.logger.Info("Rolling back database migrations", zap.Int("steps", steps))

    currentVersion, dirty, err := mm.migrate.Version()
    if err != nil && err != migrate.ErrNilVersion {
        return fmt.Errorf("failed to get current version: %w", err)
    }

    if dirty {
        return fmt.Errorf("database is in dirty state at version %d", currentVersion)
    }

    // Rollback specified number of steps
    if err := mm.migrate.Steps(-steps); err != nil {
        return fmt.Errorf("failed to rollback migrations: %w", err)
    }

    newVersion, _, err := mm.migrate.Version()
    if err != nil && err != migrate.ErrNilVersion {
        return fmt.Errorf("failed to get new version: %w", err)
    }

    mm.logger.Info("Rollback completed successfully", zap.Uint("new_version", newVersion))
    return nil
}

func (mm *MigrationManager) Goto(ctx context.Context, version uint) error {
    mm.logger.Info("Migrating to specific version", zap.Uint("target_version", version))

    if err := mm.migrate.Migrate(version); err != nil {
        return fmt.Errorf("failed to migrate to version %d: %w", version, err)
    }

    mm.logger.Info("Migration to version completed", zap.Uint("version", version))
    return nil
}

func (mm *MigrationManager) Force(ctx context.Context, version int) error {
    mm.logger.Warn("Forcing migration version", zap.Int("version", version))

    if err := mm.migrate.Force(version); err != nil {
        return fmt.Errorf("failed to force version %d: %w", version, err)
    }

    mm.logger.Info("Version forced successfully", zap.Int("version", version))
    return nil
}

func (mm *MigrationManager) Status(ctx context.Context) (*MigrationStatus, error) {
    version, dirty, err := mm.migrate.Version()
    if err != nil && err != migrate.ErrNilVersion {
        return nil, fmt.Errorf("failed to get version: %w", err)
    }

    // Get available migrations
    availableMigrations, err := mm.getAvailableMigrations()
    if err != nil {
        return nil, fmt.Errorf("failed to get available migrations: %w", err)
    }

    // Get applied migrations
    appliedMigrations, err := mm.getAppliedMigrations(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get applied migrations: %w", err)
    }

    return &MigrationStatus{
        CurrentVersion:        version,
        Dirty:                dirty,
        AvailableMigrations:   availableMigrations,
        AppliedMigrations:     appliedMigrations,
        PendingMigrations:     mm.getPendingMigrations(availableMigrations, appliedMigrations),
    }, nil
}

func (mm *MigrationManager) Validate(ctx context.Context) error {
    mm.logger.Info("Validating migrations...")

    // Check for missing down migrations
    if err := mm.validateDownMigrations(); err != nil {
        return fmt.Errorf("down migration validation failed: %w", err)
    }

    // Check migration sequence
    if err := mm.validateSequence(); err != nil {
        return fmt.Errorf("sequence validation failed: %w", err)
    }

    // Check for duplicate versions
    if err := mm.validateDuplicates(); err != nil {
        return fmt.Errorf("duplicate validation failed: %w", err)
    }

    mm.logger.Info("Migration validation passed")
    return nil
}

func (mm *MigrationManager) Close() error {
    if mm.migrate != nil {
        if sourceErr, dbErr := mm.migrate.Close(); sourceErr != nil || dbErr != nil {
            return fmt.Errorf("failed to close migrate instance: source=%v, db=%v", sourceErr, dbErr)
        }
    }

    if mm.db != nil {
        return mm.db.Close()
    }

    return nil
}

type MigrationStatus struct {
    CurrentVersion      uint                `json:"current_version"`
    Dirty              bool                `json:"dirty"`
    AvailableMigrations []MigrationInfo     `json:"available_migrations"`
    AppliedMigrations   []AppliedMigration  `json:"applied_migrations"`
    PendingMigrations   []MigrationInfo     `json:"pending_migrations"`
}

type MigrationInfo struct {
    Version     uint   `json:"version"`
    Description string `json:"description"`
    Filename    string `json:"filename"`
}

type AppliedMigration struct {
    Version uint      `json:"version"`
    Dirty   bool      `json:"dirty"`
    AppliedAt time.Time `json:"applied_at"`
}

func (mm *MigrationManager) getAvailableMigrations() ([]MigrationInfo, error) {
    var migrations []MigrationInfo

    err := filepath.WalkDir(mm.migrationsDir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() {
            return nil
        }

        filename := d.Name()
        if filepath.Ext(filename) != ".sql" {
            return nil
        }

        // Parse migration filename
        version, description, direction, err := parseMigrationFilename(filename)
        if err != nil {
            mm.logger.Warn("Failed to parse migration filename",
                zap.String("filename", filename),
                zap.Error(err),
            )
            return nil
        }

        // Only include up migrations in available list
        if direction == "up" {
            migrations = append(migrations, MigrationInfo{
                Version:     version,
                Description: description,
                Filename:    filename,
            })
        }

        return nil
    })

    return migrations, err
}

func (mm *MigrationManager) getAppliedMigrations(ctx context.Context) ([]AppliedMigration, error) {
    query := `
        SELECT version, dirty
        FROM schema_migrations
        ORDER BY version
    `

    rows, err := mm.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var migrations []AppliedMigration
    for rows.Next() {
        var migration AppliedMigration
        if err := rows.Scan(&migration.Version, &migration.Dirty); err != nil {
            return nil, err
        }
        migrations = append(migrations, migration)
    }

    return migrations, rows.Err()
}

func (mm *MigrationManager) getPendingMigrations(available []MigrationInfo, applied []AppliedMigration) []MigrationInfo {
    appliedMap := make(map[uint]bool)
    for _, migration := range applied {
        appliedMap[migration.Version] = true
    }

    var pending []MigrationInfo
    for _, migration := range available {
        if !appliedMap[migration.Version] {
            pending = append(pending, migration)
        }
    }

    return pending
}

func (mm *MigrationManager) validateDownMigrations() error {
    upMigrations := make(map[uint]bool)
    downMigrations := make(map[uint]bool)

    err := filepath.WalkDir(mm.migrationsDir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() {
            return nil
        }

        filename := d.Name()
        if filepath.Ext(filename) != ".sql" {
            return nil
        }

        version, _, direction, err := parseMigrationFilename(filename)
        if err != nil {
            return nil // Skip invalid filenames
        }

        if direction == "up" {
            upMigrations[version] = true
        } else if direction == "down" {
            downMigrations[version] = true
        }

        return nil
    })

    if err != nil {
        return err
    }

    // Check that every up migration has a corresponding down migration
    for version := range upMigrations {
        if !downMigrations[version] {
            return fmt.Errorf("missing down migration for version %d", version)
        }
    }

    return nil
}

func (mm *MigrationManager) validateSequence() error {
    var versions []uint

    err := filepath.WalkDir(mm.migrationsDir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() {
            return nil
        }

        filename := d.Name()
        if filepath.Ext(filename) != ".sql" {
            return nil
        }

        version, _, direction, err := parseMigrationFilename(filename)
        if err != nil {
            return nil
        }

        if direction == "up" {
            versions = append(versions, version)
        }

        return nil
    })

    if err != nil {
        return err
    }

    // Check for gaps in sequence
    if len(versions) == 0 {
        return nil
    }

    // Sort versions
    for i := 0; i < len(versions)-1; i++ {
        for j := i + 1; j < len(versions); j++ {
            if versions[i] > versions[j] {
                versions[i], versions[j] = versions[j], versions[i]
            }
        }
    }

    // Check for gaps
    for i := 1; i < len(versions); i++ {
        if versions[i] != versions[i-1]+1 {
            return fmt.Errorf("gap in migration sequence: version %d follows %d", versions[i], versions[i-1])
        }
    }

    return nil
}

func (mm *MigrationManager) validateDuplicates() error {
    versions := make(map[uint][]string)

    err := filepath.WalkDir(mm.migrationsDir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() {
            return nil
        }

        filename := d.Name()
        if filepath.Ext(filename) != ".sql" {
            return nil
        }

        version, _, _, err := parseMigrationFilename(filename)
        if err != nil {
            return nil
        }

        versions[version] = append(versions[version], filename)

        return nil
    })

    if err != nil {
        return err
    }

    // Check for duplicates (should have exactly 2 files per version: up and down)
    for version, files := range versions {
        if len(files) != 2 {
            return fmt.Errorf("version %d has %d files, expected 2 (up and down): %v", version, len(files), files)
        }
    }

    return nil
}
```

### Migration Utilities

```go
// internal/database/migration_utils.go
package database

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
)

var migrationFilenameRegex = regexp.MustCompile(`^(\d+)_(.+)\.(up|down)\.sql$`)

func parseMigrationFilename(filename string) (version uint, description, direction string, err error) {
    matches := migrationFilenameRegex.FindStringSubmatch(filename)
    if len(matches) != 4 {
        return 0, "", "", fmt.Errorf("invalid migration filename format: %s", filename)
    }

    versionInt, err := strconv.Atoi(matches[1])
    if err != nil {
        return 0, "", "", fmt.Errorf("invalid version number in filename: %s", filename)
    }

    return uint(versionInt), matches[2], matches[3], nil
}

func generateMigrationFilename(version uint, description, direction string) string {
    // Clean description for filename
    cleanDescription := strings.ReplaceAll(description, " ", "_")
    cleanDescription = strings.ToLower(cleanDescription)

    // Remove special characters
    reg := regexp.MustCompile(`[^a-z0-9_]`)
    cleanDescription = reg.ReplaceAllString(cleanDescription, "")

    return fmt.Sprintf("%03d_%s.%s.sql", version, cleanDescription, direction)
}

func CreateMigrationFiles(migrationsDir string, description string) error {
    // Get next version number
    nextVersion, err := getNextMigrationVersion(migrationsDir)
    if err != nil {
        return fmt.Errorf("failed to get next version: %w", err)
    }

    // Generate filenames
    upFilename := generateMigrationFilename(nextVersion, description, "up")
    downFilename := generateMigrationFilename(nextVersion, description, "down")

    // Create up migration file
    upTemplate := fmt.Sprintf(`-- Migration: %s
-- Version: %d
-- Description: %s
-- Author: %s
-- Date: %s

BEGIN;

-- Add your migration SQL here

COMMIT;
`, description, nextVersion, description, "TODO", "TODO")

    // Create down migration file
    downTemplate := fmt.Sprintf(`-- Rollback: %s
-- Version: %d
-- Description: Rollback %s

BEGIN;

-- Add your rollback SQL here

COMMIT;
`, description, nextVersion, description)

    // Write files
    upPath := filepath.Join(migrationsDir, upFilename)
    downPath := filepath.Join(migrationsDir, downFilename)

    if err := os.WriteFile(upPath, []byte(upTemplate), 0644); err != nil {
        return fmt.Errorf("failed to write up migration: %w", err)
    }

    if err := os.WriteFile(downPath, []byte(downTemplate), 0644); err != nil {
        return fmt.Errorf("failed to write down migration: %w", err)
    }

    fmt.Printf("Created migration files:\n")
    fmt.Printf("  %s\n", upPath)
    fmt.Printf("  %s\n", downPath)

    return nil
}

func getNextMigrationVersion(migrationsDir string) (uint, error) {
    var maxVersion uint = 0

    err := filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() {
            return nil
        }

        filename := d.Name()
        if filepath.Ext(filename) != ".sql" {
            return nil
        }

        version, _, _, err := parseMigrationFilename(filename)
        if err != nil {
            return nil // Skip invalid filenames
        }

        if version > maxVersion {
            maxVersion = version
        }

        return nil
    })

    if err != nil {
        return 0, err
    }

    return maxVersion + 1, nil
}
```

### CLI Commands

```go
// cmd/migrate/main.go
package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    "strconv"
    "time"

    "github.com/your-org/mcp-planner/internal/config"
    "github.com/your-org/mcp-planner/internal/database"

    "go.uber.org/zap"
)

func main() {
    var (
        configPath    = flag.String("config", "config.yaml", "Path to configuration file")
        migrationsDir = flag.String("migrations", "./migrations", "Path to migrations directory")
        command       = flag.String("command", "up", "Migration command: up, down, goto, force, status, validate, create")
        steps         = flag.Int("steps", 1, "Number of steps for down command")
        version       = flag.String("version", "", "Version for goto/force commands")
        description   = flag.String("description", "", "Description for create command")
    )
    flag.Parse()

    // Load configuration
    cfg, err := config.Load(*configPath)
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Initialize logger
    logger, err := zap.NewProduction()
    if err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
    defer logger.Sync()

    // Create migration manager
    migrationConfig := database.MigrationConfig{
        DatabaseURL:      cfg.Database.URL,
        MigrationsDir:    *migrationsDir,
        LockTimeout:      30 * time.Second,
        StatementTimeout: 60 * time.Second,
    }

    manager, err := database.NewMigrationManager(migrationConfig, logger)
    if err != nil {
        log.Fatalf("Failed to create migration manager: %v", err)
    }
    defer manager.Close()

    ctx := context.Background()

    // Execute command
    switch *command {
    case "up":
        if err := manager.Up(ctx); err != nil {
            log.Fatalf("Migration up failed: %v", err)
        }

    case "down":
        if err := manager.Down(ctx, *steps); err != nil {
            log.Fatalf("Migration down failed: %v", err)
        }

    case "goto":
        if *version == "" {
            log.Fatal("Version required for goto command")
        }
        v, err := strconv.ParseUint(*version, 10, 32)
        if err != nil {
            log.Fatalf("Invalid version: %v", err)
        }
        if err := manager.Goto(ctx, uint(v)); err != nil {
            log.Fatalf("Migration goto failed: %v", err)
        }

    case "force":
        if *version == "" {
            log.Fatal("Version required for force command")
        }
        v, err := strconv.Atoi(*version)
        if err != nil {
            log.Fatalf("Invalid version: %v", err)
        }
        if err := manager.Force(ctx, v); err != nil {
            log.Fatalf("Migration force failed: %v", err)
        }

    case "status":
        status, err := manager.Status(ctx)
        if err != nil {
            log.Fatalf("Failed to get migration status: %v", err)
        }
        printStatus(status)

    case "validate":
        if err := manager.Validate(ctx); err != nil {
            log.Fatalf("Migration validation failed: %v", err)
        }
        fmt.Println("Migration validation passed")

    case "create":
        if *description == "" {
            log.Fatal("Description required for create command")
        }
        if err := database.CreateMigrationFiles(*migrationsDir, *description); err != nil {
            log.Fatalf("Failed to create migration files: %v", err)
        }

    default:
        log.Fatalf("Unknown command: %s", *command)
    }
}

func printStatus(status *database.MigrationStatus) {
    fmt.Printf("Current Version: %d\n", status.CurrentVersion)
    fmt.Printf("Dirty: %t\n", status.Dirty)
    fmt.Printf("\nAvailable Migrations: %d\n", len(status.AvailableMigrations))
    for _, migration := range status.AvailableMigrations {
        fmt.Printf("  %d: %s\n", migration.Version, migration.Description)
    }

    fmt.Printf("\nApplied Migrations: %d\n", len(status.AppliedMigrations))
    for _, migration := range status.AppliedMigrations {
        fmt.Printf("  %d (dirty: %t)\n", migration.Version, migration.Dirty)
    }

    fmt.Printf("\nPending Migrations: %d\n", len(status.PendingMigrations))
    for _, migration := range status.PendingMigrations {
        fmt.Printf("  %d: %s\n", migration.Version, migration.Description)
    }
}
```

### Makefile Integration

```makefile
# Migration commands
MIGRATIONS_DIR := ./migrations
DATABASE_URL := $(shell grep DATABASE_URL .env | cut -d '=' -f2)

migrate-create: ## Create new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then echo "NAME is required"; exit 1; fi
	@go run cmd/migrate/main.go -command=create -description="$(NAME)" -migrations=$(MIGRATIONS_DIR)

migrate-up: ## Run all pending migrations
	@echo "Running migrations..."
	@go run cmd/migrate/main.go -command=up -migrations=$(MIGRATIONS_DIR)

migrate-down: ## Rollback last migration (usage: make migrate-down STEPS=1)
	@echo "Rolling back $(or $(STEPS),1) migration(s)..."
	@go run cmd/migrate/main.go -command=down -steps=$(or $(STEPS),1) -migrations=$(MIGRATIONS_DIR)

migrate-goto: ## Migrate to specific version (usage: make migrate-goto VERSION=5)
	@if [ -z "$(VERSION)" ]; then echo "VERSION is required"; exit 1; fi
	@echo "Migrating to version $(VERSION)..."
	@go run cmd/migrate/main.go -command=goto -version=$(VERSION) -migrations=$(MIGRATIONS_DIR)

migrate-force: ## Force migration version (usage: make migrate-force VERSION=5)
	@if [ -z "$(VERSION)" ]; then echo "VERSION is required"; exit 1; fi
	@echo "Forcing version $(VERSION)..."
	@go run cmd/migrate/main.go -command=force -version=$(VERSION) -migrations=$(MIGRATIONS_DIR)

migrate-status: ## Show migration status
	@go run cmd/migrate/main.go -command=status -migrations=$(MIGRATIONS_DIR)

migrate-validate: ## Validate migrations
	@go run cmd/migrate/main.go -command=validate -migrations=$(MIGRATIONS_DIR)

migrate-reset: ## Reset database (drop all tables and re-run migrations)
	@echo "WARNING: This will drop all tables and data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo ""; \
		migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" drop -f; \
		make migrate-up; \
	else \
		echo ""; \
		echo "Cancelled."; \
	fi
```

---

*Next: [Initial Schema Migrations](./07b3b-initial-schema.md)*
