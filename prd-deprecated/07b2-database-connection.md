# Database Connection & Setup

## Overview

This document covers database connection management, connection pooling, health checks, and initialization procedures for the MCP-Planner system using PostgreSQL and Prisma.

## Connection Management

### Database Client Setup

```go
// internal/database/connection.go
package database

import (
    "context"
    "fmt"
    "time"

    "github.com/your-org/mcp-planner/internal/config"
    "github.com/your-org/mcp-planner/internal/database/generated/db"

    "go.uber.org/zap"
)

type Client struct {
    prisma *db.PrismaClient
    config config.DatabaseConfig
    logger *zap.Logger
}

type ConnectionOptions struct {
    MaxConnections  int
    MaxIdleTime     time.Duration
    MaxLifetime     time.Duration
    ConnectTimeout  time.Duration
    QueryTimeout    time.Duration
}

func NewClient(cfg config.DatabaseConfig, logger *zap.Logger) (*Client, error) {
    client := db.NewClient()

    // Configure connection options
    if err := configureConnection(client, cfg); err != nil {
        return nil, fmt.Errorf("failed to configure database connection: %w", err)
    }

    return &Client{
        prisma: client,
        config: cfg,
        logger: logger,
    }, nil
}

func (c *Client) Connect(ctx context.Context) error {
    c.logger.Info("Connecting to database...")

    if err := c.prisma.Connect(); err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    // Test connection
    if err := c.Ping(ctx); err != nil {
        return fmt.Errorf("database ping failed: %w", err)
    }

    c.logger.Info("Database connection established")
    return nil
}

func (c *Client) Disconnect() error {
    c.logger.Info("Disconnecting from database...")

    if err := c.prisma.Disconnect(); err != nil {
        c.logger.Error("Error disconnecting from database", zap.Error(err))
        return err
    }

    c.logger.Info("Database disconnected")
    return nil
}

func (c *Client) Ping(ctx context.Context) error {
    // Use a simple query to test connectivity
    _, err := c.prisma.Prisma.QueryRaw("SELECT 1").Exec(ctx)
    return err
}

func (c *Client) Client() *db.PrismaClient {
    return c.prisma
}

func configureConnection(client *db.PrismaClient, cfg config.DatabaseConfig) error {
    // Note: Prisma Go client connection configuration is typically done
    // through environment variables or connection string parameters

    // Set connection pool parameters via environment if needed
    // This is a placeholder for actual Prisma configuration
    return nil
}
```

### Connection Pool Configuration

```go
// internal/database/pool.go
package database

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    "github.com/lib/pq"
    "go.uber.org/zap"
)

type PoolConfig struct {
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration
}

type ConnectionPool struct {
    db     *sql.DB
    config PoolConfig
    logger *zap.Logger
}

func NewConnectionPool(databaseURL string, config PoolConfig, logger *zap.Logger) (*ConnectionPool, error) {
    db, err := sql.Open("postgres", databaseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to open database connection: %w", err)
    }

    // Configure connection pool
    db.SetMaxOpenConns(config.MaxOpenConns)
    db.SetMaxIdleConns(config.MaxIdleConns)
    db.SetConnMaxLifetime(config.ConnMaxLifetime)
    db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

    pool := &ConnectionPool{
        db:     db,
        config: config,
        logger: logger,
    }

    return pool, nil
}

func (p *ConnectionPool) Ping(ctx context.Context) error {
    return p.db.PingContext(ctx)
}

func (p *ConnectionPool) Stats() sql.DBStats {
    return p.db.Stats()
}

func (p *ConnectionPool) Close() error {
    return p.db.Close()
}

func (p *ConnectionPool) LogStats() {
    stats := p.Stats()
    p.logger.Info("Database connection pool stats",
        zap.Int("open_connections", stats.OpenConnections),
        zap.Int("in_use", stats.InUse),
        zap.Int("idle", stats.Idle),
        zap.Int64("wait_count", stats.WaitCount),
        zap.Duration("wait_duration", stats.WaitDuration),
        zap.Int64("max_idle_closed", stats.MaxIdleClosed),
        zap.Int64("max_idle_time_closed", stats.MaxIdleTimeClosed),
        zap.Int64("max_lifetime_closed", stats.MaxLifetimeClosed),
    )
}

// Monitor connection pool health
func (p *ConnectionPool) StartMonitoring(ctx context.Context, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            p.LogStats()

            // Check for potential issues
            stats := p.Stats()
            if stats.OpenConnections >= p.config.MaxOpenConns {
                p.logger.Warn("Connection pool at maximum capacity",
                    zap.Int("open_connections", stats.OpenConnections),
                    zap.Int("max_open_conns", p.config.MaxOpenConns),
                )
            }

            if stats.WaitCount > 0 {
                p.logger.Warn("Connections waiting for available connection",
                    zap.Int64("wait_count", stats.WaitCount),
                    zap.Duration("wait_duration", stats.WaitDuration),
                )
            }
        }
    }
}
```

### Health Checks

```go
// internal/database/health.go
package database

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    "go.uber.org/zap"
)

type HealthChecker struct {
    client *Client
    pool   *ConnectionPool
    logger *zap.Logger
}

type HealthStatus struct {
    Healthy     bool          `json:"healthy"`
    Latency     time.Duration `json:"latency"`
    Error       string        `json:"error,omitempty"`
    PoolStats   *sql.DBStats  `json:"pool_stats,omitempty"`
    Timestamp   time.Time     `json:"timestamp"`
}

func NewHealthChecker(client *Client, pool *ConnectionPool, logger *zap.Logger) *HealthChecker {
    return &HealthChecker{
        client: client,
        pool:   pool,
        logger: logger,
    }
}

func (h *HealthChecker) Check(ctx context.Context) *HealthStatus {
    start := time.Now()

    status := &HealthStatus{
        Timestamp: start,
    }

    // Test basic connectivity
    if err := h.client.Ping(ctx); err != nil {
        status.Healthy = false
        status.Error = fmt.Sprintf("ping failed: %v", err)
        status.Latency = time.Since(start)
        return status
    }

    // Test query execution
    if err := h.testQuery(ctx); err != nil {
        status.Healthy = false
        status.Error = fmt.Sprintf("query test failed: %v", err)
        status.Latency = time.Since(start)
        return status
    }

    status.Healthy = true
    status.Latency = time.Since(start)

    // Add pool stats if available
    if h.pool != nil {
        stats := h.pool.Stats()
        status.PoolStats = &stats
    }

    return status
}

func (h *HealthChecker) testQuery(ctx context.Context) error {
    // Test with a simple query that exercises the database
    _, err := h.client.prisma.Prisma.QueryRaw(`
        SELECT
            COUNT(*) as project_count,
            NOW() as current_time
    `).Exec(ctx)

    return err
}

func (h *HealthChecker) DeepCheck(ctx context.Context) *HealthStatus {
    start := time.Now()

    status := &HealthStatus{
        Timestamp: start,
    }

    // Basic health check first
    basicStatus := h.Check(ctx)
    if !basicStatus.Healthy {
        return basicStatus
    }

    // Test write operations
    if err := h.testWriteOperation(ctx); err != nil {
        status.Healthy = false
        status.Error = fmt.Sprintf("write test failed: %v", err)
        status.Latency = time.Since(start)
        return status
    }

    // Test complex queries
    if err := h.testComplexQuery(ctx); err != nil {
        status.Healthy = false
        status.Error = fmt.Sprintf("complex query test failed: %v", err)
        status.Latency = time.Since(start)
        return status
    }

    status.Healthy = true
    status.Latency = time.Since(start)

    if h.pool != nil {
        stats := h.pool.Stats()
        status.PoolStats = &stats
    }

    return status
}

func (h *HealthChecker) testWriteOperation(ctx context.Context) error {
    // Create and immediately delete a test record
    testID := fmt.Sprintf("health-check-%d", time.Now().UnixNano())

    // Create test project
    project, err := h.client.prisma.Project.CreateOne(
        db.Project.Name.Set(testID),
        db.Project.Description.Set("Health check test project"),
    ).Exec(ctx)

    if err != nil {
        return fmt.Errorf("failed to create test project: %w", err)
    }

    // Delete test project
    _, err = h.client.prisma.Project.FindUnique(
        db.Project.ID.Equals(project.ID),
    ).Delete().Exec(ctx)

    if err != nil {
        return fmt.Errorf("failed to delete test project: %w", err)
    }

    return nil
}

func (h *HealthChecker) testComplexQuery(ctx context.Context) error {
    // Test a complex query that exercises joins and aggregations
    _, err := h.client.prisma.Prisma.QueryRaw(`
        SELECT
            p.id,
            p.name,
            COUNT(t.id) as task_count,
            AVG(t.progress) as avg_progress
        FROM projects p
        LEFT JOIN tasks t ON p.id = t.project_id
        WHERE p.created_at > NOW() - INTERVAL '1 day'
        GROUP BY p.id, p.name
        LIMIT 10
    `).Exec(ctx)

    return err
}

// Continuous health monitoring
func (h *HealthChecker) StartMonitoring(ctx context.Context, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            status := h.Check(ctx)

            if status.Healthy {
                h.logger.Debug("Database health check passed",
                    zap.Duration("latency", status.Latency),
                )
            } else {
                h.logger.Error("Database health check failed",
                    zap.String("error", status.Error),
                    zap.Duration("latency", status.Latency),
                )
            }

            // Log pool stats if available
            if status.PoolStats != nil {
                h.logger.Debug("Database pool stats",
                    zap.Int("open_connections", status.PoolStats.OpenConnections),
                    zap.Int("in_use", status.PoolStats.InUse),
                    zap.Int("idle", status.PoolStats.Idle),
                )
            }
        }
    }
}
```

### Database Initialization

```go
// internal/database/init.go
package database

import (
    "context"
    "fmt"
    "time"

    "github.com/your-org/mcp-planner/internal/config"

    "go.uber.org/zap"
)

type DatabaseManager struct {
    client       *Client
    pool         *ConnectionPool
    healthChecker *HealthChecker
    config       config.DatabaseConfig
    logger       *zap.Logger
}

func NewDatabaseManager(cfg config.DatabaseConfig, logger *zap.Logger) (*DatabaseManager, error) {
    // Create Prisma client
    client, err := NewClient(cfg, logger)
    if err != nil {
        return nil, fmt.Errorf("failed to create database client: %w", err)
    }

    // Create connection pool for monitoring
    poolConfig := PoolConfig{
        MaxOpenConns:    cfg.MaxConnections,
        MaxIdleConns:    cfg.MaxConnections / 2,
        ConnMaxLifetime: cfg.MaxLifetime,
        ConnMaxIdleTime: cfg.MaxIdleTime,
    }

    pool, err := NewConnectionPool(cfg.URL, poolConfig, logger)
    if err != nil {
        return nil, fmt.Errorf("failed to create connection pool: %w", err)
    }

    // Create health checker
    healthChecker := NewHealthChecker(client, pool, logger)

    return &DatabaseManager{
        client:        client,
        pool:          pool,
        healthChecker: healthChecker,
        config:        cfg,
        logger:        logger,
    }, nil
}

func (dm *DatabaseManager) Initialize(ctx context.Context) error {
    dm.logger.Info("Initializing database...")

    // Connect to database
    if err := dm.client.Connect(ctx); err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    // Test connection pool
    if err := dm.pool.Ping(ctx); err != nil {
        return fmt.Errorf("connection pool ping failed: %w", err)
    }

    // Run initial health check
    status := dm.healthChecker.Check(ctx)
    if !status.Healthy {
        return fmt.Errorf("initial health check failed: %s", status.Error)
    }

    dm.logger.Info("Database initialized successfully",
        zap.Duration("initial_latency", status.Latency),
    )

    return nil
}

func (dm *DatabaseManager) StartMonitoring(ctx context.Context) {
    // Start health monitoring
    go dm.healthChecker.StartMonitoring(ctx, 30*time.Second)

    // Start pool monitoring
    go dm.pool.StartMonitoring(ctx, 60*time.Second)

    dm.logger.Info("Database monitoring started")
}

func (dm *DatabaseManager) Shutdown(ctx context.Context) error {
    dm.logger.Info("Shutting down database connections...")

    // Close Prisma client
    if err := dm.client.Disconnect(); err != nil {
        dm.logger.Error("Error disconnecting Prisma client", zap.Error(err))
    }

    // Close connection pool
    if err := dm.pool.Close(); err != nil {
        dm.logger.Error("Error closing connection pool", zap.Error(err))
        return err
    }

    dm.logger.Info("Database connections closed")
    return nil
}

func (dm *DatabaseManager) Client() *Client {
    return dm.client
}

func (dm *DatabaseManager) HealthChecker() *HealthChecker {
    return dm.healthChecker
}

func (dm *DatabaseManager) GetStats() map[string]interface{} {
    stats := dm.pool.Stats()

    return map[string]interface{}{
        "open_connections":      stats.OpenConnections,
        "in_use":               stats.InUse,
        "idle":                 stats.Idle,
        "wait_count":           stats.WaitCount,
        "wait_duration":        stats.WaitDuration.String(),
        "max_idle_closed":      stats.MaxIdleClosed,
        "max_idle_time_closed": stats.MaxIdleTimeClosed,
        "max_lifetime_closed":  stats.MaxLifetimeClosed,
    }
}
```

### Connection Retry Logic

```go
// internal/database/retry.go
package database

import (
    "context"
    "fmt"
    "math"
    "time"

    "go.uber.org/zap"
)

type RetryConfig struct {
    MaxAttempts     int
    InitialDelay    time.Duration
    MaxDelay        time.Duration
    BackoffFactor   float64
    RetryableErrors []string
}

type RetryableOperation func() error

func WithRetry(ctx context.Context, config RetryConfig, logger *zap.Logger, operation RetryableOperation) error {
    var lastErr error

    for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
        err := operation()
        if err == nil {
            if attempt > 1 {
                logger.Info("Operation succeeded after retry",
                    zap.Int("attempt", attempt),
                )
            }
            return nil
        }

        lastErr = err

        // Check if error is retryable
        if !isRetryableError(err, config.RetryableErrors) {
            logger.Error("Non-retryable error encountered",
                zap.Error(err),
                zap.Int("attempt", attempt),
            )
            return err
        }

        if attempt == config.MaxAttempts {
            logger.Error("Max retry attempts reached",
                zap.Error(err),
                zap.Int("max_attempts", config.MaxAttempts),
            )
            break
        }

        // Calculate delay with exponential backoff
        delay := calculateDelay(attempt, config)

        logger.Warn("Operation failed, retrying",
            zap.Error(err),
            zap.Int("attempt", attempt),
            zap.Duration("delay", delay),
        )

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(delay):
            // Continue to next attempt
        }
    }

    return fmt.Errorf("operation failed after %d attempts: %w", config.MaxAttempts, lastErr)
}

func calculateDelay(attempt int, config RetryConfig) time.Duration {
    delay := float64(config.InitialDelay) * math.Pow(config.BackoffFactor, float64(attempt-1))

    if delay > float64(config.MaxDelay) {
        delay = float64(config.MaxDelay)
    }

    return time.Duration(delay)
}

func isRetryableError(err error, retryableErrors []string) bool {
    if len(retryableErrors) == 0 {
        // Default retryable errors for PostgreSQL
        retryableErrors = []string{
            "connection refused",
            "connection reset",
            "timeout",
            "temporary failure",
            "server closed the connection",
        }
    }

    errStr := err.Error()
    for _, retryableErr := range retryableErrors {
        if contains(errStr, retryableErr) {
            return true
        }
    }

    return false
}

func contains(s, substr string) bool {
    return len(s) >= len(substr) && (s == substr ||
        (len(s) > len(substr) &&
         (s[:len(substr)] == substr ||
          s[len(s)-len(substr):] == substr ||
          containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
    for i := 0; i <= len(s)-len(substr); i++ {
        if s[i:i+len(substr)] == substr {
            return true
        }
    }
    return false
}

// Connection wrapper with retry logic
func (c *Client) ConnectWithRetry(ctx context.Context) error {
    config := RetryConfig{
        MaxAttempts:   5,
        InitialDelay:  1 * time.Second,
        MaxDelay:      30 * time.Second,
        BackoffFactor: 2.0,
    }

    return WithRetry(ctx, config, c.logger, func() error {
        return c.Connect(ctx)
    })
}

func (c *Client) PingWithRetry(ctx context.Context) error {
    config := RetryConfig{
        MaxAttempts:   3,
        InitialDelay:  500 * time.Millisecond,
        MaxDelay:      5 * time.Second,
        BackoffFactor: 2.0,
    }

    return WithRetry(ctx, config, c.logger, func() error {
        return c.Ping(ctx)
    })
}
```

### Environment-Specific Configuration

```go
// internal/database/config.go
package database

import (
    "fmt"
    "net/url"
    "strconv"
    "strings"
    "time"

    "github.com/your-org/mcp-planner/internal/config"
)

func BuildDatabaseConfig(cfg config.DatabaseConfig) (*DatabaseConfig, error) {
    // Parse database URL
    dbURL, err := url.Parse(cfg.URL)
    if err != nil {
        return nil, fmt.Errorf("invalid database URL: %w", err)
    }

    // Extract connection parameters
    params := dbURL.Query()

    config := &DatabaseConfig{
        Host:     dbURL.Hostname(),
        Port:     getPort(dbURL.Port()),
        Database: strings.TrimPrefix(dbURL.Path, "/"),
        Username: dbURL.User.Username(),
        SSLMode:  getSSLMode(params.Get("sslmode")),

        // Connection pool settings
        MaxConnections:  cfg.MaxConnections,
        MaxIdleTime:     cfg.MaxIdleTime,
        MaxLifetime:     cfg.MaxLifetime,
        ConnectTimeout:  30 * time.Second,
        QueryTimeout:    30 * time.Second,
    }

    if password, ok := dbURL.User.Password(); ok {
        config.Password = password
    }

    return config, nil
}

type DatabaseConfig struct {
    Host           string
    Port           int
    Database       string
    Username       string
    Password       string
    SSLMode        string
    MaxConnections int
    MaxIdleTime    time.Duration
    MaxLifetime    time.Duration
    ConnectTimeout time.Duration
    QueryTimeout   time.Duration
}

func (dc *DatabaseConfig) ConnectionString() string {
    return fmt.Sprintf(
        "host=%s port=%d dbname=%s user=%s password=%s sslmode=%s connect_timeout=%d",
        dc.Host,
        dc.Port,
        dc.Database,
        dc.Username,
        dc.Password,
        dc.SSLMode,
        int(dc.ConnectTimeout.Seconds()),
    )
}

func getPort(portStr string) int {
    if portStr == "" {
        return 5432 // Default PostgreSQL port
    }

    port, err := strconv.Atoi(portStr)
    if err != nil {
        return 5432
    }

    return port
}

func getSSLMode(sslMode string) string {
    if sslMode == "" {
        return "prefer"
    }
    return sslMode
}

// Environment-specific configurations
func DevelopmentConfig() DatabaseConfig {
    return DatabaseConfig{
        MaxConnections:  10,
        MaxIdleTime:     15 * time.Minute,
        MaxLifetime:     1 * time.Hour,
        ConnectTimeout:  10 * time.Second,
        QueryTimeout:    30 * time.Second,
        SSLMode:        "disable",
    }
}

func ProductionConfig() DatabaseConfig {
    return DatabaseConfig{
        MaxConnections:  25,
        MaxIdleTime:     15 * time.Minute,
        MaxLifetime:     1 * time.Hour,
        ConnectTimeout:  30 * time.Second,
        QueryTimeout:    60 * time.Second,
        SSLMode:        "require",
    }
}

func TestConfig() DatabaseConfig {
    return DatabaseConfig{
        MaxConnections:  5,
        MaxIdleTime:     5 * time.Minute,
        MaxLifetime:     30 * time.Minute,
        ConnectTimeout:  5 * time.Second,
        QueryTimeout:    10 * time.Second,
        SSLMode:        "disable",
    }
}
```

---

*Next: [Migrations & Schema Evolution](./07b3-migrations.md)*
