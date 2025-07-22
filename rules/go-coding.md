# Go Coding Policies

- This file is project agnostic and reflects the users prefered golang coding practices.
- Keep it concise and focused to provide the best developer experience.
- Provide generic samples that show the principle rather than project related examples.

## Policies

- You as the agent have to enforce the **Single Responsibility Principle** throughout the codebase.
- You as the agent have to enforce the **Dependency Inversion Principle (DIP)** throughout the codebase.

## File Size Management

- **Maximum File Length**: 500 lines per Go file (enforced by linter and CI)
- **Function Length**: Maximum 60 lines per function, 40 statements
- **Line Length**: Maximum 120 characters per line
- **Refactoring Triggers**: Files approaching 400+ lines should be considered for refactoring
- **Monitoring**: Use `make check-file-length` to verify compliance
- **Enforcement**: Automated checks in pre-commit hooks and CI/CD pipeline

## Error Handling

- Establish custom error tooling in the errs package
- Use custom error types for complex error cases that need additional data.
- Document all possible errors in function godoc comments.
- Handle errors at the appropriate level - don't ignore them silently.
- Return errors as the last return value in functions.
- Define sentinel errors at package level: `var ErrNotFound = errors.New("not found")`

## Naming Conventions

- Use camelCase for variable and function names.
- Use PascalCase for exported types, functions, and constants.
- Use ALL_CAPS with underscores for constants that are not exported.
- Use descriptive names that clearly indicate purpose and scope.
- Avoid abbreviations unless they are widely understood (e.g., `ctx` for context).
- Use noun phrases for types and verb phrases for functions.
- Prefix interface names with 'I' only when necessary for disambiguation.
- Use consistent naming patterns across the codebase.

## Package Organization

- Keep packages focused on a single responsibility.
- Use clear, descriptive package names that reflect their purpose.
- Avoid generic package names like `util`, `common`, or `helper`.
- Place internal packages under `internal/` to prevent external imports.
- Group related functionality in the same package.
- Minimize dependencies between packages to reduce coupling.
- Use package-level documentation to explain the package's purpose.

### File Organization Within Packages

- **Single Responsibility**: Each file should have a clear, single purpose
- **Logical Grouping**: Group related types, functions, and methods in the same file
- **File Naming**: Use descriptive names that indicate the file's primary responsibility
- **Size Limits**: Keep files under 500 lines to maintain readability and maintainability
- **Split Strategies**: When files grow large, split by:
  - **Functionality**: Move related functions to separate files (e.g., `user_auth.go`, `user_profile.go`)
  - **Types**: Extract large type definitions to dedicated files (e.g., `types.go`, `models.go`)
  - **Interfaces**: Move interface definitions to separate files when appropriate (e.g., `interfaces.go`)
  - **Implementation**: Separate interface definitions from implementations
  - **Domain Logic**: Group by business domain or feature area
- **File Naming Conventions**:
  - Use snake_case for multi-word file names
  - Include the primary type or functionality in the name
  - Use suffixes like `_test.go`, `_mock.go`, `_impl.go` when appropriate

## Types

- Keep all types in the `internal/types` package.
- In this package maintain the type definitions of providers located in the `internal/providers` package and all other types used within the application.
- Use interfaces instead of concrete types whenever possible.
- Use descriptive, specific type names that reflect their purpose.
- Document all exported types with clear usage examples.
- Prefer composition over inheritance.
- Use embedding for extending functionality when appropriate.
- Define interfaces at the point of use, not at the point of implementation.

## Context Usage Patterns

- Always pass `context.Context` as the first parameter in functions that may be cancelled or have timeouts.
- Use `context.WithTimeout` or `context.WithDeadline` for operations with time limits.
- Use `context.WithCancel` for operations that may need to be cancelled.
- Use `context.WithValue` sparingly and only for request-scoped data.
- Check for context cancellation in long-running operations using `ctx.Done()`.
- Propagate context through the call stack consistently.
- Don't store contexts in structs; pass them as function parameters.

## Concurrency Best Practices

- Use channels for communication between goroutines.
- Prefer `sync.WaitGroup` for waiting on multiple goroutines to complete.
- Use `sync.Mutex` and `sync.RWMutex` for protecting shared data.
- Use `sync.Once` for one-time initialization.
- Avoid sharing memory by communicating; communicate by sharing memory.
- Use buffered channels when appropriate to prevent blocking.
- Always handle goroutine lifecycle and cleanup properly.
- Use `context.Context` for cancellation and timeouts in concurrent operations.

## Interface Design Principles

- Keep interfaces small and focused (prefer many small interfaces over few large ones).
- Define interfaces at the consumer side, not the producer side.
- Use composition to build larger interfaces from smaller ones.
- Prefer implicit interface satisfaction over explicit implementation.
- Document interface contracts clearly, including expected behavior.
- Use interface{} (any) sparingly and only when type safety is not critical.
- Consider using type assertions and type switches when working with interfaces.

## Dependency Injection

- Provide all services and providers located in `internal/providers` in the `internal/container` package.
- Provide them as interfaces defined in the `internal/types` package.
- Use `github.com/samber/do` for dependency injection.
- Use `do.Provide` to register dependencies in the injector.
- Use `do.MustInvoke` to resolve dependencies from the injector.
- Avoid circular dependencies between packages by providing dependencies through the injector.
- Use constructor functions that return interfaces, not concrete types.
- Keep dependency graphs shallow and well-documented.

## Providers

- Keep all providers in the `internal/providers` package.
- Create isolated providers that follow a single purpose.
- When needed inject other providers using `do.MustInvoke` for managing dependencies.
- A provider should have a constructor function that returns an instance of the provider and any potential errors.
- A provider should have a `func Startup(ctx context.Context) func() error` function that is called during startup.
- The `Startup` function should perform any initialization or setup required before the service can run.

Example:

```golang
type serviceProvider struct {
    db  *sql.DB
    log *zap.Logger  // Use structured logger
}

func (p *serviceProvider) Start(ctx context.Context) error {
    p.log.Info("Starting service", zap.String("service", "example_service"))

    // Perform startup tasks here
    return nil
}

func (p *serviceProvider) Stop(ctx context.Context) error {
    p.log.Info("Shutting down service", zap.String("service", "example_service"))
    // Perform cleanup tasks here
    return nil
}

func NewServiceProvider(db *sql.DB, logger *zap.Logger) *serviceProvider {
    provider := &serviceProvider{
        db:  db,
        log: logger,
    }

    logger.Info("Service provider initialized",
        zap.String("service", "example_service"),
        zap.String("status", "ready"))

    return provider
}
```

## Performance Considerations

- Use profiling tools (`go tool pprof`) to identify performance bottlenecks.
- Prefer slices over arrays for better performance and flexibility.
- Use `sync.Pool` for frequently allocated objects to reduce GC pressure.
- Avoid premature optimization; measure before optimizing.
- Use benchmarks (`go test -bench`) to validate performance improvements.
- Consider memory allocation patterns and minimize unnecessary allocations.
- Use appropriate data structures for the use case (maps vs slices vs channels).
- Profile memory usage and CPU usage separately.

## Logging

### Structured Logging with Zap

- **Required**: Use `go.uber.org/zap` for all logging throughout the application
- **Logger Type**: Use `*zap.Logger` (structured) instead of `*zap.SugaredLogger`
- **Logger Access**: Get logger from dependency injection container, not global access
- **Initialization**: Always call `logger.Initialize(level)` and defer `logger.Sync()`

### Logging Pattern

Use structured logging with typed fields consistently:

```go
// ✅ Correct - Structured logging with typed fields
log.Info("User authentication successful",
    zap.String("user_id", userID),
    zap.String("method", "oauth"),
    zap.Duration("duration", time.Since(start)),
    zap.Bool("first_login", isFirstLogin))

logger.Log.Error("Database connection failed",
    zap.Error(err),
    zap.String("database", "users"),
    zap.String("operation", "connect"),
    zap.Int("retry_count", retryCount))

logger.Log.Warn("Rate limit approaching",
    zap.String("client_ip", clientIP),
    zap.Int("current_requests", currentReqs),
    zap.Int("limit", rateLimit),
    zap.String("endpoint", endpoint))
```

### Field Types

Use appropriate zap field types for type safety and performance:

- `zap.String("key", stringValue)` - for string values
- `zap.Int("key", intValue)` - for integer values
- `zap.Int64("key", int64Value)` - for int64 values
- `zap.Bool("key", boolValue)` - for boolean values
- `zap.Error(err)` - for error values (no key needed)
- `zap.Duration("key", duration)` - for time.Duration values
- `zap.Time("key", timeValue)` - for time.Time values
- `zap.Any("key", anyValue)` - for complex objects (use sparingly)

### Logging Levels

- **Debug**: Detailed information for debugging (development only)
- **Info**: General operational information (user actions, system events)
- **Warn**: Warning conditions that should be noted but don't stop operation
- **Error**: Error conditions that need attention but application continues
- **Fatal**: Critical errors that cause application termination (avoid in libraries)

### Context-Rich Logging

Always provide relevant context in log messages:

```go
// ✅ Good - Rich context
logger.Log.Info("HTTP request processed",
    zap.String("method", r.Method),
    zap.String("path", r.URL.Path),
    zap.Int("status_code", statusCode),
    zap.Duration("duration", time.Since(start)),
    zap.String("user_agent", r.UserAgent()),
    zap.String("remote_addr", r.RemoteAddr))

// ❌ Bad - No context
logger.Log.Info("Request processed")
```

### Function Logging Pattern

Functions should include structured logging for operations:

```go
func someOperation(ctx context.Context, logger *zap.Logger, userID string) error {
    start := time.Now()
    logger.Info("Starting operation",
        zap.String("operation", "user_update"),
        zap.String("user_id", userID))

    // ... operation logic ...

    if err != nil {
        logger.Error("Operation failed",
            zap.Error(err),
            zap.String("operation", "user_update"),
            zap.String("user_id", userID))
        return err
    }

    logger.Info("Operation completed successfully",
        zap.String("operation", "user_update"),
        zap.String("user_id", userID),
        zap.Duration("duration", time.Since(start)))

    return nil
}
```


### Anti-Patterns

Avoid these logging patterns:

```go
// ❌ Don't use sugared logger
log.Infow("message", "key", value)

// ❌ Don't use string formatting in messages
logger.Log.Info(fmt.Sprintf("User %s logged in", userID))

// ❌ Don't log without context
logger.Log.Error("Something failed")

// ❌ Don't use inconsistent field names
logger.Log.Info("User action", zap.String("userId", id))  // inconsistent naming
logger.Log.Info("User action", zap.String("user_id", id)) // consistent naming

// ❌ Don't log sensitive information
logger.Log.Info("User login", zap.String("password", password)) // Never log passwords!
```

### Performance Considerations

- Structured logging is more performant than string formatting
- Use `zap.Any()` sparingly as it uses reflection
- Consider log level checks for expensive operations:
  ```go
  if logger.Log.Core().Enabled(zap.DebugLevel) {
      logger.Log.Debug("Expensive debug info", zap.Any("complex_object", obj))
  }
  ```

## Configuration

- Use `github.com/kelseyhightower/envconfig` for loading configuration from environment variables.
- Define configuration parameters in a `Config` struct in the `internal/config` package.
- Use `envconfig.Process` to populate the `Config` struct from environment variables.
- Provide sensible defaults for configuration values.
- Validate configuration values at startup.
- Document all configuration options with examples.
- Use environment-specific configuration files when appropriate.

## Testing

- Write tests for all exported functions and methods.
- Use table-driven tests for functions with multiple test cases.
- Prefer `testify/assert` and `testify/require` for assertions in tests.
- Test files should be named `*_test.go` in the same package.
- Test functions should start with `Test` and describe what they test.
- Use test helpers for common test setup/teardown logic.
- Write integration tests for critical paths and external dependencies.
- Use mocks and stubs to isolate units under test.
- Aim for high test coverage but focus on testing behavior, not implementation.
- Use `testify/suite` for complex test scenarios that require setup/teardown.
- Test error conditions and edge cases thoroughly.

## Documentation

- Keep all docs in the `./docs` directory
- Use English language for all documentation and comments.
- Use GoDoc comments for package-level, function-level, and type-level documentation.
- Use complete sentences and proper punctuation in comments.
- Use the `//` comment style for single-line comments and `/* */` for multi-line comments.
- Document exported functions, types, and variables with clear descriptions.
- Include usage examples in documentation when helpful.
- Use `//nolint:...` to disable specific linter warnings when necessary, but use it sparingly and only when justified.
- Use `//go:generate` comments to document code generation commands.
- Keep documentation up-to-date with code changes.
- Document complex algorithms and business logic thoroughly.

## Code Quality and Maintenance

### File Size Management
- **Enforcement**: Automated checks ensure no Go file exceeds 500 lines
- **Monitoring**: Regular reviews of file sizes during development
- **Refactoring Triggers**:
  - Files approaching 500+ lines should be evaluated for splitting
  - Functions exceeding 60 lines should be broken down
  - Complex logic should be extracted into smaller, focused functions
  - High cyclomatic complexity indicates need for refactoring

### Refactoring Strategies

#### When to Refactor
- **File Size**: Files approaching or exceeding 500 lines
- **Function Complexity**: Functions with high cyclomatic complexity (>10)
- **Code Duplication**: Repeated patterns across multiple files
- **Poor Cohesion**: Unrelated functionality grouped together
- **High Coupling**: Excessive dependencies between components

#### How to Refactor
1. **Extract Functions**: Move large blocks of logic into separate, well-named functions
2. **Extract Types**: Move large type definitions to dedicated files
3. **Extract Interfaces**: Separate interface definitions when they become substantial
4. **Create Sub-packages**: For related functionality that grows beyond a single package
5. **Split by Responsibility**: Separate different concerns into different files
6. **Maintain Cohesion**: Keep related functionality together while respecting size limits

#### Refactoring Process
1. **Backup**: Always create backups before major refactoring (`*.go.backup`)
2. **Incremental**: Make small, focused changes rather than large rewrites
3. **Test Coverage**: Ensure comprehensive tests before refactoring
4. **Preserve Behavior**: Refactoring should not change external behavior
5. **Update Documentation**: Keep documentation in sync with structural changes
6. **Review Dependencies**: Ensure refactoring doesn't break import cycles
7. **Validate**: Run full test suite and linting after refactoring

### Quality Metrics
- **File Length**: Maximum 500 lines per file
- **Function Length**: Maximum 60 lines, 40 statements per function
- **Line Length**: Maximum 120 characters per line
- **Cyclomatic Complexity**: Keep functions simple and focused (<10 complexity)
- **Test Coverage**: Maintain high test coverage during refactoring (>80%)
- **Documentation Coverage**: All exported items must be documented
- **Import Cycles**: Zero tolerance for circular dependencies

### Automation and Tooling
- **Linting**: Comprehensive golangci-lint configuration with 40+ linters
- **File Length Checks**: Automated verification of file size constraints
- **Pre-commit Hooks**: Automated quality checks before commits
- **CI/CD Integration**: Quality gates in continuous integration pipeline
- **Formatting**: Automated code formatting with gofmt and goimports
- **Security Scanning**: Automated vulnerability detection with gosec
- **Dependency Analysis**: Regular dependency updates and vulnerability scanning

### Best Practices for Large Codebases
- **Modular Design**: Design for modularity from the beginning
- **Clear Boundaries**: Establish clear package and file boundaries
- **Regular Refactoring**: Don't let technical debt accumulate
- **Code Reviews**: Focus on file size and complexity during reviews
- **Documentation**: Maintain architectural documentation for large systems
- **Dependency Management**: Keep dependencies clean and minimal
- **Testing Strategy**: Comprehensive testing at all levels (unit, integration, e2e)
- **Monitoring**: Track code quality metrics over time
- **Team Standards**: Ensure all team members follow the same standards

### File Organization Examples

#### Good File Organization
```
internal/
├── user/
│   ├── types.go          # User-related types (50 lines)
│   ├── service.go        # User service logic (200 lines)
│   ├── repository.go     # Data access layer (150 lines)
│   └── handlers.go       # HTTP handlers (180 lines)
├── auth/
│   ├── interfaces.go     # Authentication interfaces (30 lines)
│   ├── jwt.go           # JWT implementation (120 lines)
│   └── middleware.go    # Auth middleware (80 lines)
```

#### Poor File Organization
```
internal/
├── user.go              # Everything user-related (800 lines) ❌
├── auth.go              # All auth logic (600 lines) ❌
└── handlers.go          # All HTTP handlers (1200 lines) ❌
```

### Refactoring Checklist
- [ ] File is under 500 lines
- [ ] Functions are under 60 lines
- [ ] Single responsibility per file
- [ ] Clear, descriptive file names
- [ ] Related functionality grouped together
- [ ] Interfaces separated from implementations when appropriate
- [ ] No circular dependencies
- [ ] All exported items documented
- [ ] Tests updated and passing
- [ ] Linting passes without warnings

### Tooling Commands
```bash
# Check file length constraints
make check-file-length

# Run comprehensive linting
make lint

# Fix auto-fixable linting issues
make lint-fix

# Format code
make fmt

# Run all pre-commit checks
make pre-commit

# Run full CI pipeline
make ci
```

### Quality Gates
1. **Development**: File length monitoring during coding
2. **Pre-commit**: Automated checks before commits
3. **Code Review**: Manual review focusing on complexity and size
4. **CI/CD**: Automated quality gates in pipeline
5. **Deployment**: Final validation before release
