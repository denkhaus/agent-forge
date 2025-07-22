# MCP Implementation Patterns & Best Practices

This document provides advanced patterns and best practices for implementing robust MCP servers using the `github.com/mark3labs/mcp-go` library.

## Advanced Parameter Handling Patterns

### Complex Parameter Validation

#### Array Parameter with Type Validation
```go
func validateStringArray(request mcp.CallToolRequest, paramName string, required bool) ([]string, error) {
    if required {
        itemsInterface, err := request.RequireArray(paramName)
        if err != nil {
            return nil, err
        }

        items := make([]string, len(itemsInterface))
        for i, item := range itemsInterface {
            if itemStr, ok := item.(string); ok {
                items[i] = itemStr
            } else {
                return nil, fmt.Errorf("all %s must be strings, got %T at index %d", paramName, item, i)
            }
        }

        if len(items) == 0 {
            return nil, fmt.Errorf("%s cannot be empty", paramName)
        }

        return items, nil
    }

    // Optional array handling
    itemsInterface, _ := request.GetArray(paramName)
    if itemsInterface == nil {
        return nil, nil
    }

    topics := make([]string, len(topicsInterface))
    for i, topic := range topicsInterface {
        if topicStr, ok := topic.(string); ok {
            topics[i] = topicStr
        } else {
            return nil, fmt.Errorf("all %s must be strings, got %T at index %d", paramName, topic, i)
        }
    }

    return topics, nil
}
```

#### Number Parameter with Range Validation
```go
func validateNumberRange(request mcp.CallToolRequest, paramName string, defaultVal, min, max float64) (float64, error) {
    value := request.GetNumberOrDefault(paramName, defaultVal)

    if value < min || value > max {
        return 0, fmt.Errorf("%s must be between %.1f and %.1f, got %.1f", paramName, min, max, value)
    }

    return value, nil
}

// Usage in handler
threshold, err := validateNumberRange(request, "threshold", 0.7, 0.0, 1.0)
if err != nil {
    return mcp.NewToolResultError(err.Error()), nil
}
```

### Advanced Object Parameter Patterns

#### Nested Object Validation
```go
func validateNestedConfig(request mcp.CallToolRequest) (*Config, error) {
    configObj := helpers.GetObject(request, "config", nil)
    if configObj == nil {
        return &Config{}, nil // Return default config
    }

    config := &Config{}

    // Validate search settings
    if searchObj, ok := configObj["search"].(map[string]interface{}); ok {
        if threshold, ok := searchObj["threshold"].(float64); ok {
            if threshold < 0.0 || threshold > 1.0 {
                return nil, fmt.Errorf("search.threshold must be between 0.0 and 1.0")
            }
            config.SearchThreshold = threshold
        }

        if limit, ok := searchObj["limit"].(float64); ok {
            if limit < 1 || limit > 100 {
                return nil, fmt.Errorf("search.limit must be between 1 and 100")
            }
            config.SearchLimit = int(limit)
        }
    }

    return config, nil
}
```

#### Object Merging Pattern
```go
func mergeMetadata(existing, updates map[string]interface{}) map[string]interface{} {
    if existing == nil {
        existing = make(map[string]interface{})
    }

    result := make(map[string]interface{})

    // Copy existing
    for k, v := range existing {
        result[k] = v
    }

    // Apply updates
    for k, v := range updates {
        if v == nil {
            delete(result, k) // nil values remove keys
        } else {
            result[k] = v
        }
    }

    return result
}
```

## Response Formatting Patterns

### Standardized Success Response
```go
type StandardResponse struct {
    Success   bool                   `json:"success"`
    Data      interface{}            `json:"data,omitempty"`
    Message   string                 `json:"message,omitempty"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
    Timestamp string                 `json:"timestamp,omitempty"`
}

func createSuccessResponse(data interface{}, message string) *StandardResponse {
    return &StandardResponse{
        Success:   true,
        Data:      data,
        Message:   message,
        Timestamp: time.Now().UTC().Format(time.RFC3339),
    }
}

func marshalResponse(response *StandardResponse) (*mcp.CallToolResult, error) {
    resultBytes, err := json.Marshal(response)
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response: %v", err)), err
    }
    return mcp.NewToolResultText(string(resultBytes)), nil
}
```

### Paginated Response Pattern
```go
type PaginatedResponse struct {
    Success    bool                   `json:"success"`
    Data       interface{}            `json:"data"`
    Pagination PaginationInfo         `json:"pagination"`
    Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type PaginationInfo struct {
    Total      int    `json:"total"`
    Count      int    `json:"count"`
    Limit      int    `json:"limit"`
    HasMore    bool   `json:"has_more"`
    NextCursor string `json:"next_cursor,omitempty"`
}

func createPaginatedResponse(data interface{}, total, count, limit int, nextCursor string) *PaginatedResponse {
    return &PaginatedResponse{
        Success: true,
        Data:    data,
        Pagination: PaginationInfo{
            Total:      total,
            Count:      count,
            Limit:      limit,
            HasMore:    count == limit, // Assume more if we hit the limit
            NextCursor: nextCursor,
        },
    }
}
```

## Error Handling Patterns

### Structured Error Response
```go
type ErrorDetails struct {
    Code      string                 `json:"code"`
    Message   string                 `json:"message"`
    Details   map[string]interface{} `json:"details,omitempty"`
    Timestamp string                 `json:"timestamp"`
}

func createErrorResponse(code, message string, details map[string]interface{}) *mcp.CallToolResult {
    errorResp := ErrorDetails{
        Code:      code,
        Message:   message,
        Details:   details,
        Timestamp: time.Now().UTC().Format(time.RFC3339),
    }

    errorBytes, _ := json.Marshal(errorResp)
    return mcp.NewToolResultError(string(errorBytes))
}

// Usage examples
return createErrorResponse("VALIDATION_ERROR", "Invalid parameter", map[string]interface{}{
    "parameter": "threshold",
    "value":     threshold,
    "expected":  "0.0 <= value <= 1.0",
}), nil

return createErrorResponse("SERVICE_ERROR", "Vector service unavailable", map[string]interface{}{
    "service": "qdrant",
    "retry":   true,
}), err
```

### Error Classification
```go
const (
    ErrorCodeValidation   = "VALIDATION_ERROR"
    ErrorCodeNotFound     = "NOT_FOUND"
    ErrorCodeService      = "SERVICE_ERROR"
    ErrorCodePermission   = "PERMISSION_ERROR"
    ErrorCodeRateLimit    = "RATE_LIMIT"
    ErrorCodeInternal     = "INTERNAL_ERROR"
)

func classifyError(err error) string {
    switch {
    case strings.Contains(err.Error(), "not found"):
        return ErrorCodeNotFound
    case strings.Contains(err.Error(), "validation"):
        return ErrorCodeValidation
    case strings.Contains(err.Error(), "connection"):
        return ErrorCodeService
    default:
        return ErrorCodeInternal
    }
}
```

## Logging Patterns

### Contextual Logging
```go
func (s *VectorServer) logOperation(ctx context.Context, operation string, params map[string]interface{}) {
    fields := []zap.Field{
        zap.String("operation", operation),
        zap.String("mcp_tool", operation),
    }

    for key, value := range params {
        switch v := value.(type) {
        case string:
            fields = append(fields, zap.String(key, v))
        case int:
            fields = append(fields, zap.Int(key, v))
        case float64:
            fields = append(fields, zap.Float64(key, v))
        case bool:
            fields = append(fields, zap.Bool(key, v))
        default:
            fields = append(fields, zap.Any(key, v))
        }
    }

    s.logger.Debug("MCP operation started", fields...)
}

// Usage
s.logOperation(ctx, "store_memory", map[string]interface{}{
    "project_id":    projectID,
    "title_length":  len(title),
    "topics_count":  len(topics),
    "content_length": len(content),
    "has_metadata":  metadata != nil,
})
```

### Performance Logging
```go
func (s *VectorServer) withPerformanceLogging(operation string, fn func() (*mcp.CallToolResult, error)) (*mcp.CallToolResult, error) {
    start := time.Now()

    result, err := fn()

    duration := time.Since(start)

    logFields := []zap.Field{
        zap.String("operation", operation),
        zap.Duration("duration", duration),
        zap.Bool("success", err == nil),
    }

    if err != nil {
        logFields = append(logFields, zap.Error(err))
        s.logger.Warn("MCP operation completed with error", logFields...)
    } else {
        s.logger.Debug("MCP operation completed successfully", logFields...)
    }

    return result, err
}

// Usage in handler
func (s *VectorServer) StoreMemoryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    return s.withPerformanceLogging("store_memory", func() (*mcp.CallToolResult, error) {
        // Handler implementation
        return s.storeMemoryImpl(ctx, request)
    })
}
```

## Testing Patterns

### Mock Service Setup
```go
func setupMockService() *MockBusinessService {
    mockService := &MockBusinessService{}

    // Default successful responses
    mockService.On("CreateItem", mock.Anything, mock.AnythingOfType("string"),
        mock.AnythingOfType("[]string"), mock.AnythingOfType("string"),
        mock.Anything).Return(
        &Item{
            ID:      "test-item-id",
            Title:   "Test Item",
            Tags:    []string{"test"},
            Content: "Test content",
        }, nil)

    return mockService
}

func setupErrorMockService() *MockBusinessService {
    mockService := &MockBusinessService{}

    // Error responses
    mockService.On("CreateItem", mock.Anything, mock.Anything,
        mock.Anything, mock.Anything, mock.Anything).Return(
        (*Item)(nil), errors.New("service error"))

    return mockService
}
```

### Test Data Builders
```go
type TestRequestBuilder struct {
    args map[string]interface{}
}

func NewTestRequest() *TestRequestBuilder {
    return &TestRequestBuilder{
        args: make(map[string]interface{}),
    }
}

func (b *TestRequestBuilder) WithString(key, value string) *TestRequestBuilder {
    b.args[key] = value
    return b
}

func (b *TestRequestBuilder) WithArray(key string, value []interface{}) *TestRequestBuilder {
    b.args[key] = value
    return b
}

func (b *TestRequestBuilder) WithObject(key string, value map[string]interface{}) *TestRequestBuilder {
    b.args[key] = value
    return b
}

func (b *TestRequestBuilder) Build() mcp.CallToolRequest {
    return &mockCallToolRequest{args: b.args}
}

// Usage
request := NewTestRequest().
    WithString("project_id", "test-project").
    WithString("title", "Test Memory").
    WithArray("topics", []interface{}{"test", "memory"}).
    WithString("content", "Test content").
    Build()
```

### Assertion Helpers
```go
func assertSuccessResponse(t *testing.T, result *mcp.CallToolResult) map[string]interface{} {
    assert.False(t, result.IsError)

    var response map[string]interface{}
    err := json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &response)
    assert.NoError(t, err)

    assert.True(t, response["success"].(bool))
    return response
}

func assertErrorResponse(t *testing.T, result *mcp.CallToolResult, expectedMessage string) {
    assert.True(t, result.IsError)

    errorText := result.Content[0].(mcp.TextContent).Text
    assert.Contains(t, errorText, expectedMessage)
}
```

These patterns provide a comprehensive foundation for building robust, maintainable MCP implementations with proper error handling, logging, and testing practices.
