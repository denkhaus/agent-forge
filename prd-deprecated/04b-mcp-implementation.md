# MCP Server Implementation Guide

## Table of Contents
1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Server Setup](#server-setup)
4. [Tool Registration](#tool-registration)
5. [Handler Implementation](#handler-implementation)
6. [Helper Functions](#helper-functions)
7. [Best Practices](#best-practices)
8. [Code Examples](#code-examples)

## Overview

This guide teaches you how to implement a Model Context Protocol (MCP) server using the `github.com/mark3labs/mcp-go` library. You'll learn to create robust MCP servers with proper tool registration, parameter validation, and error handling.

### Key Features You'll Learn
- **Tool Registration**: How to register MCP tools with proper parameter definitions
- **Parameter Validation**: Type-safe parameter extraction and validation
- **Handler Implementation**: Best practices for tool handler functions
- **Error Handling**: Proper MCP error response patterns
- **Structured Responses**: JSON-formatted tool results

### Required Dependencies
```go
import (
    "context"
    "encoding/json"
    "fmt"
    
    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)
```

## Architecture

### Recommended Project Structure
```
internal/mcp/
├── server.go                    # Main server setup and tool registration
├── handlers.go                  # Tool handler implementations
├── helpers.go                   # Parameter extraction utilities
├── instructions.go              # Server instructions and documentation
└── server_test.go              # Test suite
```

### Core Components

#### Server Structure
```go
type MCPServer struct {
    config     *Config
    logger     Logger
    service    BusinessService  // Your business logic interface
    mcpServer  *server.MCPServer
}
```

#### Key Dependencies
- **BusinessService**: Interface for your application's business logic
- **Logger**: Logging interface (can be any logger)
- **Config**: Your application configuration
- **MCPServer**: Core MCP protocol server from the library

## Server Setup

### Server Initialization
```go
func NewMCPServer(config *Config, logger Logger, service BusinessService) (*MCPServer, error) {
    // Create MCP server with options
    mcpServer := server.NewMCPServer(config.Name, config.Version,
        server.WithToolCapabilities(true),
        server.WithLogging(),
        server.WithInstructions(getServerInstructions()),
    )

    s := &MCPServer{
        config:    config,
        logger:    logger,
        service:   service,
        mcpServer: mcpServer,
    }

    // Register all tools
    s.registerTools()
    return s, nil
}
```

### Server Startup
```go
func (s *MCPServer) Start(ctx context.Context) error {
    s.logger.Info("Starting MCP server",
        "transport", "http",
        "name", s.config.Name,
        "version", s.config.Version)

    srv := server.NewStreamableHTTPServer(s.mcpServer,
        server.WithLogger(s.logger),
    )
    return srv.Start(s.config.ServerAddr)
}
```

## Tool Registration

### Tool Registration Pattern
The MCP server uses a consistent pattern for tool registration using the builder pattern with `mcp.NewTool` and parameter definitions.

#### Basic Tool Registration
```go
s.mcpServer.AddTool(
    mcp.NewTool("tool_name",
        mcp.WithDescription("Tool description"),
        // Parameter definitions...
    ),
    s.HandlerFunction,
)
```

### Parameter Types and Validation

#### String Parameters
```go
mcp.WithString("parameter_name",
    mcp.Description("Parameter description"),
    mcp.Required(), // For required parameters
)
```

#### Number Parameters with Constraints
```go
mcp.WithNumber("threshold",
    mcp.Description("Similarity threshold (0.0-1.0) for filtering results"),
    mcp.DefaultNumber(0.7),
    mcp.Min(0.0),
    mcp.Max(1.0),
)
```

#### Array Parameters
```go
mcp.WithArray("tags",
    mcp.Description("List of tags for categorization"),
    mcp.Required(),
)
```

#### Object Parameters
```go
mcp.WithObject("metadata",
    mcp.Description("Additional metadata (optional)"),
)
```

### Complete Tool Examples

#### Example Tool: Create Item
```go
s.mcpServer.AddTool(
    mcp.NewTool("create_item",
        mcp.WithDescription("Create a new item with title, tags, and content"),
        mcp.WithString("title",
            mcp.Description("Title of the item"),
            mcp.Required(),
        ),
        mcp.WithArray("tags",
            mcp.Description("List of tags for categorization"),
            mcp.Required(),
        ),
        mcp.WithString("content",
            mcp.Description("Main content of the item"),
            mcp.Required(),
        ),
        mcp.WithObject("metadata",
            mcp.Description("Additional metadata (optional)"),
        ),
    ),
    s.CreateItemHandler,
)
```

#### Example Tool: Search Items
```go
s.mcpServer.AddTool(
    mcp.NewTool("search_items",
        mcp.WithDescription("Search for items using a query string"),
        mcp.WithString("query",
            mcp.Description("Search query string"),
            mcp.Required(),
        ),
        mcp.WithNumber("limit",
            mcp.Description("Maximum number of results to return"),
            mcp.DefaultNumber(10),
        ),
        mcp.WithNumber("threshold",
            mcp.Description("Relevance threshold (0.0-1.0) for filtering results"),
            mcp.DefaultNumber(0.7),
            mcp.Min(0.0),
            mcp.Max(1.0),
        ),
        mcp.WithArray("tags",
            mcp.Description("Filter by specific tags (optional)"),
        ),
    ),
    s.SearchItemsHandler,
)
```

#### Example Tool: Update Item
```go
s.mcpServer.AddTool(
    mcp.NewTool("update_item",
        mcp.WithDescription("Update an existing item"),
        mcp.WithString("item_id",
            mcp.Description("ID of the item to update"),
            mcp.Required(),
        ),
        mcp.WithString("title",
            mcp.Description("New title (optional)"),
        ),
        mcp.WithArray("tags",
            mcp.Description("New tags (optional)"),
        ),
        mcp.WithString("content",
            mcp.Description("New content (optional)"),
        ),
        mcp.WithNumber("priority",
            mcp.Description("Priority level (1-10)"),
            mcp.DefaultNumber(5),
            mcp.Min(1),
            mcp.Max(10),
        ),
        mcp.WithObject("metadata",
            mcp.Description("Additional metadata (optional)"),
        ),
    ),
    s.UpdateItemHandler,
)
```

## Handler Implementation

### Handler Function Signature
All MCP tool handlers follow this signature:
```go
func (s *MCPServer) HandlerName(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
```

### Parameter Extraction Patterns

#### Required String Parameters
```go
func (s *MCPServer) ExampleHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Extract required string parameter
    itemID, err := request.RequireString("item_id")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), err
    }

    // Handler logic...
}
```

#### Optional Parameters with Defaults
```go
// Extract optional number with default
limit := request.GetNumberOrDefault("limit", 10)
priority := request.GetNumberOrDefault("priority", 5)

// Extract optional string with default
status := request.GetStringOrDefault("status", "active")
```

#### Array Parameters
```go
// Extract required array parameter
tagsInterface, err := request.RequireArray("tags")
if err != nil {
    return mcp.NewToolResultError(err.Error()), err
}

// Convert to string slice
tags := make([]string, len(tagsInterface))
for i, tag := range tagsInterface {
    if tagStr, ok := tag.(string); ok {
        tags[i] = tagStr
    } else {
        return mcp.NewToolResultError("All tags must be strings"), nil
    }
}

// Validate array is not empty
if len(tags) == 0 {
    return mcp.NewToolResultError("Tags cannot be empty"), nil
}
```

### Response Patterns

#### Success Response
```go
result := map[string]interface{}{
    "success": true,
    "data":    responseData,
    "message": "Operation completed successfully",
}

resultBytes, err := json.Marshal(result)
if err != nil {
    return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal result: %v", err)), err
}

return mcp.NewToolResultText(string(resultBytes)), nil
```

#### Error Response
```go
if err != nil {
    return mcp.NewToolResultError(fmt.Sprintf("Operation failed: %v", err)), err
}
```

### Complete Handler Example

#### Complete Handler Example
```go
func (s *MCPServer) CreateItemHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    s.logger.Debug("Creating item via MCP")

    // Extract required parameters
    title, err := request.RequireString("title")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), err
    }

    content, err := request.RequireString("content")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), err
    }

    // Extract and validate tags array
    tagsInterface, err := request.RequireArray("tags")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), err
    }

    tags := make([]string, len(tagsInterface))
    for i, tag := range tagsInterface {
        if tagStr, ok := tag.(string); ok {
            tags[i] = tagStr
        } else {
            return mcp.NewToolResultError("All tags must be strings"), nil
        }
    }

    if len(tags) == 0 {
        return mcp.NewToolResultError("Tags cannot be empty"), nil
    }

    // Extract optional metadata using helper
    metadata := GetObject(request, "metadata", nil)

    // Call business service
    item, err := s.service.CreateItem(ctx, title, tags, content, metadata)
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("Failed to create item: %v", err)), err
    }

    s.logger.Debug("Item created via MCP", 
        "item_id", item.ID,
        "title", title)

    // Prepare response
    result := map[string]interface{}{
        "success": true,
        "item":    item,
    }

    resultBytes, err := json.Marshal(result)
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal result: %v", err)), err
    }

    return mcp.NewToolResultText(string(resultBytes)), nil
}
```

## Helper Functions

### MCP Helper Utilities
Create utility functions for parameter extraction and validation to reduce code duplication.

#### GetObject Helper
```go
// GetObject extracts an object (map[string]interface{}) parameter from the request
// Returns the default value if the parameter doesn't exist or is not a valid object
func GetObject(request mcp.CallToolRequest, key string, defaultValue ...map[string]interface{}) map[string]interface{} {
    args := request.GetArguments()
    if value, exists := args[key]; exists {
        if objectMap, ok := value.(map[string]interface{}); ok {
            return objectMap
        }
    }

    // Return default value if provided, otherwise nil
    if len(defaultValue) > 0 {
        return defaultValue[0]
    }
    return nil
}
```

#### RequireObject Helper
```go
// RequireObject extracts a required object (map[string]interface{}) parameter from the request
// Returns an error if the parameter doesn't exist or is not a valid object
func RequireObject(request mcp.CallToolRequest, key string) (map[string]interface{}, error) {
    args := request.GetArguments()
    value, exists := args[key]
    if !exists {
        return nil, fmt.Errorf("missing required parameter: %s", key)
    }

    objectMap, ok := value.(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("parameter %s must be an object, got %T", key, value)
    }

    return objectMap, nil
}
```

#### Usage Examples
```go
// Optional metadata with default
metadata := GetObject(request, "metadata", nil)

// Required configuration object
config, err := RequireObject(request, "config")
if err != nil {
    return mcp.NewToolResultError(err.Error()), err
}

// Optional settings with default values
defaultSettings := map[string]interface{}{
    "enabled": true,
    "timeout": 30,
}
settings := GetObject(request, "settings", defaultSettings)
```

## Best Practices

### 1. Use MCP Library Functions
**Always prefer the official MCP library functions over custom implementations:**

#### Tool Creation
```go
// ✅ PREFERRED: Use mcp.NewTool with builder pattern
tool := mcp.NewTool("tool_name",
    mcp.WithDescription("Tool description"),
    mcp.WithString("param", mcp.Description("Parameter description"), mcp.Required()),
)

// ❌ AVOID: Custom tool creation
```

#### Parameter Definition
```go
// ✅ PREFERRED: Use mcp parameter builders
mcp.WithString("name", mcp.Description("Name parameter"), mcp.Required())
mcp.WithNumber("count", mcp.Description("Count parameter"), mcp.DefaultNumber(10), mcp.Min(1), mcp.Max(100))
mcp.WithArray("items", mcp.Description("List of items"))
mcp.WithObject("metadata", mcp.Description("Metadata object"))

// ❌ AVOID: Manual parameter validation
```

### 2. Parameter Validation Patterns

#### Required Parameters
```go
// ✅ PREFERRED: Use built-in validation
projectID, err := request.RequireString("project_id")
if err != nil {
    return mcp.NewToolResultError(err.Error()), err
}
```

#### Optional Parameters
```go
// ✅ PREFERRED: Use default values
limit := request.GetNumberOrDefault("limit", 10)
threshold := request.GetNumberOrDefault("threshold", 0.7)
```

#### Object Parameters
```go
// ✅ PREFERRED: Use helper functions
metadata := helpers.GetObject(request, "metadata", nil)

// For required objects
config, err := helpers.RequireObject(request, "config")
if err != nil {
    return mcp.NewToolResultError(err.Error()), err
}
```

### 3. Error Handling

#### Consistent Error Responses
```go
// ✅ PREFERRED: Descriptive error messages
if err != nil {
    return mcp.NewToolResultError(fmt.Sprintf("Failed to create item: %v", err)), err
}

// ✅ PREFERRED: Validation errors
if len(tags) == 0 {
    return mcp.NewToolResultError("Tags cannot be empty"), nil
}
```

### 4. Response Structure

#### Consistent Success Response
```go
result := map[string]interface{}{
    "success": true,
    "data":    responseData,
    "message": "Operation completed successfully",
}

resultBytes, err := json.Marshal(result)
if err != nil {
    return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal result: %v", err)), err
}

return mcp.NewToolResultText(string(resultBytes)), nil
```

### 5. Logging

#### Structured Logging
```go
s.logger.Debug("Operation started via MCP",
    "operation", "create_item",
    "title", title)

s.logger.Error("Operation failed",
    "error", err,
    "operation", "create_item")
```

## Code Examples

### Complete Tool Registration
```go
func (s *MCPServer) registerTools() {
    // Core operations
    s.registerCoreTools()
    
    // Additional operations  
    s.registerAdvancedTools()
}

func (s *MCPServer) registerCoreTools() {
    // Create item
    s.mcpServer.AddTool(
        mcp.NewTool("create_item",
            mcp.WithDescription("Create a new item"),
        ),
        s.CreateItemHandler,
    )

    // Search items
    s.mcpServer.AddTool(
        mcp.NewTool("search_items",
            mcp.WithDescription("Search for items"),
            mcp.WithString("query", mcp.Description("Search query"), mcp.Required()),
            mcp.WithNumber("limit", mcp.Description("Result limit"), mcp.DefaultNumber(10)),
            mcp.WithArray("tags", mcp.Description("Filter tags")),
            mcp.WithObject("options", mcp.Description("Search options")),
        ),
        s.SearchItemsHandler,
    )

    // Additional tools...
}
```

### Advanced Parameter Handling
```go
func (s *MCPServer) AdvancedHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Required parameters
    itemID, err := request.RequireString("item_id")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), err
    }

    // Optional parameters with validation
    maxResults := request.GetNumberOrDefault("max_results", 10)
    if maxResults < 1 || maxResults > 100 {
        return mcp.NewToolResultError("max_results must be between 1 and 100"), nil
    }

    // Array parameter with type checking
    categoriesInterface, _ := request.GetArray("categories")
    var categories []string
    if categoriesInterface != nil {
        categories = make([]string, len(categoriesInterface))
        for i, c := range categoriesInterface {
            if categoryStr, ok := c.(string); ok {
                categories[i] = categoryStr
            } else {
                return mcp.NewToolResultError("All categories must be strings"), nil
            }
        }
    }

    // Object parameter with helper
    options := GetObject(request, "options", map[string]interface{}{
        "include_metadata": true,
        "sort_by": "relevance",
    })

    // Business logic...
    result := map[string]interface{}{
        "success": true,
        "item_id": itemID,
        "max_results": maxResults,
        "categories": categories,
        "options": options,
    }

    resultBytes, err := json.Marshal(result)
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal result: %v", err)), err
    }

    return mcp.NewToolResultText(string(resultBytes)), nil
}
```

This comprehensive guide covers all aspects of MCP server implementation, providing patterns, best practices, and complete examples for building robust MCP tools.
