# MCP Architecture & File Organization

## File Structure & Logical Separation

Organize your MCP implementation into logical modules for maintainability and clarity:

```
internal/mcp/
├── server.go                    # Server setup, initialization, and tool registration
├── handlers.go                  # Tool handler implementations
├── helpers.go                   # MCP utility functions (GetObject, RequireObject)
├── instructions.go              # Server instructions and documentation
└── server_test.go              # Test suite
```

## Architectural Principles

### 1. Separation of Concerns

#### Server Layer (`server.go`)
- **Responsibility**: Server initialization, configuration, tool registration
- **Key Functions**:
  - `NewMCPServer()` - Server creation with dependency injection
  - `Start()` - HTTP server startup with logging
  - `registerTools()` - Central tool registration
  - `Health()` - Health check coordination

#### Handler Layer (`handlers.go`)
- **Responsibility**: Tool handler implementations
- **Pattern**: One handler function per tool
- **Structure**: Consistent parameter extraction and response formatting

#### Helper Layer (`helpers.go`)
- **Responsibility**: Reusable parameter extraction utilities
- **Functions**: `GetObject()`, `RequireObject()`

### 2. Dependency Injection Pattern

```go
// Clean dependency injection
func NewMCPServer(config *Config, logger Logger, service BusinessService) (*MCPServer, error) {
    // Server creation with injected dependencies
    mcpServer := server.NewMCPServer(config.Name, config.Version,
        server.WithToolCapabilities(true),
        server.WithLogging(),
    )

    return &MCPServer{
        config:    config,
        logger:    logger,
        service:   service,
        mcpServer: mcpServer,
    }, nil
}
```

**Benefits**:
- Testable components with mock injection
- Clear dependency boundaries
- Configuration-driven setup

### 3. Interface-First Design

```go
type BusinessService interface {
    CreateItem(ctx context.Context, title string, tags []string, content string, metadata map[string]interface{}) (*Item, error)
    SearchItems(ctx context.Context, query string, limit int, threshold float32) ([]*Item, error)
    GetItem(ctx context.Context, itemID string) (*Item, error)
    UpdateItem(ctx context.Context, itemID string, updates map[string]interface{}) (*Item, error)
    DeleteItem(ctx context.Context, itemID string) error
    // ... additional methods
}
```

**Benefits**:
- Easy mocking for tests
- Implementation flexibility
- Clear contracts between layers

## Tool Organization Strategy

### Core CRUD Tools
Located in `handlers.go`:

1. **create_item** - Item creation with validation
2. **search_items** - Search with filters and thresholds
3. **get_item** - Single item retrieval
4. **list_items** - Item listing with pagination
5. **update_item** - Flexible item updates
6. **delete_item** - Item removal

### Advanced Tools (Examples)
Additional functionality you might implement:

1. **batch_create_items** - Bulk item creation
2. **export_items** - Data export functionality
3. **import_items** - Data import with validation
4. **get_item_stats** - Analytics and statistics
5. **archive_items** - Soft deletion/archiving
6. **restore_items** - Restore archived items

## Configuration Architecture

### Server Configuration
```go
mcpServer := server.NewMCPServer(config.Name, config.Version,
    server.WithToolCapabilities(true),
    server.WithLogging(),
    server.WithInstructions(getServerInstructions()),
)
```

## Error Handling Architecture

### Consistent Error Response Pattern
```go
// Service layer error
if err != nil {
    return mcp.NewToolResultError(fmt.Sprintf("Failed to create item: %v", err)), err
}

// Validation error
if len(topics) == 0 {
    return mcp.NewToolResultError("Topics cannot be empty"), nil
}
```

### Error Categories
1. **Validation Errors**: Parameter validation failures
2. **Service Errors**: Vector service operation failures
3. **Marshaling Errors**: JSON serialization failures
4. **System Errors**: Infrastructure or dependency failures

## Logging Architecture

### Structured Logging
```go
s.logger.Debug("Item created via MCP",
    "item_id", item.ID,
    "title", title,
    "tags_count", len(tags))

s.logger.Error("Failed to create item",
    "error", err,
    "title", title)
```

### Log Levels
- **Debug**: Operation details, parameter values
- **Info**: Server startup, major operations
- **Error**: Operation failures, system errors

## Testing Architecture

### Test Organization
```go
// Test structure mirrors handler organization
func TestStoreMemoryHandler(t *testing.T) {
    // Setup with mocks
    mockService := &MockVectorService{}
    server := createTestServer(mockService)

    // Test scenarios
    t.Run("successful_store", func(t *testing.T) { /* ... */ })
    t.Run("validation_error", func(t *testing.T) { /* ... */ })
    t.Run("service_error", func(t *testing.T) { /* ... */ })
}
```

### Mock Patterns
```go
type MockBusinessService struct {
    mock.Mock
}

func (m *MockBusinessService) CreateItem(ctx context.Context, title string, tags []string, content string, metadata map[string]interface{}) (*Item, error) {
    args := m.Called(ctx, title, tags, content, metadata)
    return args.Get(0).(*Item), args.Error(1)
}
```

## Performance Considerations

### Concurrent Safety
- All handlers are stateless and thread-safe
- Business service handles concurrent operations
- No shared mutable state in handlers

### Resource Management
- Context-based cancellation support
- Proper error handling prevents resource leaks
- Structured logging for monitoring

### Scalability Patterns
- Interface-based design allows service swapping
- Configuration-driven behavior
- Stateless handler design supports horizontal scaling

This architectural approach ensures maintainable, testable, and scalable MCP server implementation.
