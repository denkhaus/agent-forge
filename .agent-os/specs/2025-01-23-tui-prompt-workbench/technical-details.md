# TUI Prompt Workbench - Technical Implementation Details

> **Spec:** TUI Prompt Workbench Enhancement  
> **Created:** 2025-01-23  
> **Architecture:** Local-first with AI integration  

## Architecture Overview

### System Components

```
┌─────────────────────────────────────────────────────────────┐
│                    TUI Prompt Workbench                     │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────┐ │
│  │   Editor    │ │ Variables   │ │    Test     │ │Optimize │ │
│  │    Tab      │ │     Tab     │ │     Tab     │ │   Tab   │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────┘ │
├─────────────────────────────────────────────────────────────┤
│                    Workbench Manager                        │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────┐ │
│  │   Session   │ │   Model     │ │Optimization │ │Variable │ │
│  │  Manager    │ │  Provider   │ │   Engine    │ │Manager  │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────┘ │
├─────────────────────────────────────────────────────────────┤
│                      Core Services                          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────┐ │
│  │   Prompt    │ │    File     │ │   Config    │ │   DI    │ │
│  │  Service    │ │   System    │ │  Manager    │ │Container│ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Core Interfaces & Types

### Enhanced TUI Manager Interface

```go
// TUIManager defines the interface for TUI management operations
type TUIManager interface {
    RunPromptWorkbench(promptName string) error
    CreateWorkbenchSession(config WorkbenchConfig) (*WorkbenchSession, error)
    SaveSession(session *WorkbenchSession) error
    LoadSession(sessionID string) (*WorkbenchSession, error)
    GetAvailableModels() ([]ModelInfo, error)
}

// WorkbenchConfig defines configuration for workbench sessions
type WorkbenchConfig struct {
    PromptName      string
    InitialContent  string
    Variables       []VariableDefinition
    ModelProviders  []string
    AutoSave        bool
    SessionPath     string
}

// WorkbenchSession represents a complete workbench session
type WorkbenchSession struct {
    ID              string
    PromptName      string
    Content         string
    Variables       []VariableDefinition
    TestResults     []TestResult
    OptimizationRuns []OptimizationRun
    CreatedAt       time.Time
    UpdatedAt       time.Time
    Version         string
}
```

### Model Provider System

```go
// ModelProvider defines the interface for AI model integration
type ModelProvider interface {
    Name() string
    DisplayName() string
    TestPrompt(ctx context.Context, request TestRequest) (*TestResult, error)
    EstimateCost(prompt string, variables map[string]interface{}) (*CostEstimate, error)
    GetCapabilities() ModelCapabilities
    ValidateConfig() error
}

// TestRequest encapsulates a prompt testing request
type TestRequest struct {
    Prompt      string
    Variables   map[string]interface{}
    Temperature float64
    MaxTokens   int
    Model       string
    Metadata    map[string]interface{}
}

// TestResult contains the response from a model test
type TestResult struct {
    ID           string
    ModelName    string
    ModelVersion string
    Request      TestRequest
    Response     string
    ResponseTime time.Duration
    TokensUsed   TokenUsage
    Cost         float64
    Timestamp    time.Time
    Success      bool
    Error        error
    Metadata     map[string]interface{}
}

// TokenUsage tracks token consumption
type TokenUsage struct {
    PromptTokens     int
    CompletionTokens int
    TotalTokens      int
}

// CostEstimate provides cost information
type CostEstimate struct {
    EstimatedCost   float64
    Currency        string
    TokensEstimate  int
    Confidence      float64
}

// ModelCapabilities describes what a model can do
type ModelCapabilities struct {
    MaxTokens        int
    SupportsStreaming bool
    SupportsFunctions bool
    SupportedFormats []string
    CostPerToken     float64
}
```

### Optimization Engine

```go
// PromptOptimizer defines the interface for AI-powered prompt optimization
type PromptOptimizer interface {
    AnalyzeGap(ctx context.Context, analysis GapAnalysisRequest) (*GapAnalysis, error)
    EnhancePrompt(ctx context.Context, request EnhancementRequest) (*EnhancementResult, error)
    EvaluateImprovement(before, after TestResult) (*ImprovementMetrics, error)
}

// GapAnalysisRequest defines parameters for gap analysis
type GapAnalysisRequest struct {
    OriginalPrompt  string
    DesiredOutcome  string
    ActualResults   []TestResult
    Context         map[string]interface{}
}

// GapAnalysis contains the results of analyzing prompt performance gaps
type GapAnalysis struct {
    IdentifiedGaps    []PerformanceGap
    Recommendations   []Recommendation
    ConfidenceScore   float64
    AnalysisMetadata  map[string]interface{}
}

// PerformanceGap describes a specific area where the prompt underperforms
type PerformanceGap struct {
    Category     string // "clarity", "specificity", "format", "context"
    Description  string
    Severity     float64 // 0.0 to 1.0
    Examples     []string
    SuggestedFix string
}

// Recommendation provides specific improvement suggestions
type Recommendation struct {
    Type        string // "add_context", "improve_clarity", "add_examples", etc.
    Description string
    Priority    int
    Impact      float64
    Effort      float64
}

// EnhancementRequest defines parameters for prompt enhancement
type EnhancementRequest struct {
    OriginalPrompt string
    GapAnalysis    *GapAnalysis
    Constraints    EnhancementConstraints
    Strategy       EnhancementStrategy
}

// EnhancementConstraints limit how the prompt can be modified
type EnhancementConstraints struct {
    MaxLength        int
    PreserveFormat   bool
    RequiredElements []string
    ForbiddenChanges []string
}

// EnhancementStrategy defines the approach to improvement
type EnhancementStrategy struct {
    FocusAreas      []string // Which gaps to prioritize
    Aggressiveness  float64  // How much change to allow (0.0 to 1.0)
    IterationLimit  int
    ConvergenceThreshold float64
}

// EnhancementResult contains the improved prompt and metadata
type EnhancementResult struct {
    EnhancedPrompt   string
    ChangesApplied   []ChangeDescription
    ExpectedImprovement float64
    Rationale        string
    Metadata         map[string]interface{}
}

// ChangeDescription explains what was changed and why
type ChangeDescription struct {
    Type        string
    Description string
    Location    string // Where in the prompt
    Rationale   string
}

// ImprovementMetrics quantifies the improvement between versions
type ImprovementMetrics struct {
    OverallScore      float64
    CategoryScores    map[string]float64
    ResponseQuality   float64
    TaskCompletion    float64
    Consistency       float64
    Efficiency        float64
}
```

### Variable Management System

```go
// VariableManager handles prompt variable operations
type VariableManager interface {
    DefineVariable(definition VariableDefinition) error
    ValidateVariables(variables map[string]interface{}) error
    SubstituteVariables(prompt string, variables map[string]interface{}) (string, error)
    CreateTestSet(name string, variables map[string]interface{}) error
    GetTestSets() ([]TestSet, error)
    ImportVariables(source io.Reader, format string) error
    ExportVariables(dest io.Writer, format string) error
}

// VariableDefinition describes a prompt variable
type VariableDefinition struct {
    Name         string
    Type         VariableType
    Description  string
    DefaultValue interface{}
    Required     bool
    Constraints  VariableConstraints
    Examples     []interface{}
    Metadata     map[string]interface{}
}

// VariableType defines the data type of a variable
type VariableType string

const (
    StringType  VariableType = "string"
    NumberType  VariableType = "number"
    BooleanType VariableType = "boolean"
    ArrayType   VariableType = "array"
    ObjectType  VariableType = "object"
)

// VariableConstraints define validation rules
type VariableConstraints struct {
    MinLength    *int
    MaxLength    *int
    Pattern      *string
    MinValue     *float64
    MaxValue     *float64
    AllowedValues []interface{}
    CustomValidator string
}

// TestSet represents a collection of variable values for testing
type TestSet struct {
    Name        string
    Description string
    Variables   map[string]interface{}
    CreatedAt   time.Time
    Tags        []string
}
```

## Data Persistence

### Session File Format

```yaml
# .agentforge/workbench/sessions/code-reviewer-20250123.yaml
session:
  id: "sess_20250123_143022"
  version: "1.0"
  created_at: "2025-01-23T14:30:22Z"
  updated_at: "2025-01-23T15:45:18Z"
  
  prompt:
    name: "code-reviewer"
    content: |
      You are an expert code reviewer. Please review the following {{language}} code:
      
      ```{{language}}
      {{code}}
      ```
      
      Provide feedback on:
      1. Code quality and best practices
      2. Potential bugs or issues
      3. Performance improvements
      4. Security considerations
      
      Format your response as structured feedback with specific suggestions.
    
    variables:
      - name: "code"
        type: "string"
        description: "The code to review"
        required: true
        constraints:
          min_length: 10
          max_length: 10000
      
      - name: "language"
        type: "string"
        description: "Programming language"
        default: "go"
        constraints:
          allowed_values: ["go", "python", "javascript", "rust", "java"]
  
  test_sets:
    - name: "go-example"
      description: "Simple Go function"
      variables:
        code: |
          func fibonacci(n int) int {
              if n <= 1 {
                  return n
              }
              return fibonacci(n-1) + fibonacci(n-2)
          }
        language: "go"
    
    - name: "python-example"
      description: "Python class with issues"
      variables:
        code: |
          class Calculator:
              def add(self, a, b):
                  return a + b
              def divide(self, a, b):
                  return a / b  # Potential division by zero
        language: "python"
  
  test_results:
    - id: "test_20250123_143045"
      model_name: "gpt-4"
      model_version: "gpt-4-0125-preview"
      test_set: "go-example"
      response: |
        ## Code Review: Fibonacci Function
        
        ### Issues Identified:
        1. **Performance**: Exponential time complexity O(2^n)
        2. **Stack overflow**: Deep recursion for large inputs
        
        ### Suggestions:
        1. Use memoization or iterative approach
        2. Add input validation for negative numbers
        
        ### Improved Version:
        ```go
        func fibonacci(n int) int {
            if n < 0 {
                return -1 // or handle error appropriately
            }
            if n <= 1 {
                return n
            }
            
            a, b := 0, 1
            for i := 2; i <= n; i++ {
                a, b = b, a+b
            }
            return b
        }
        ```
      
      response_time: "2.3s"
      tokens_used:
        prompt_tokens: 156
        completion_tokens: 234
        total_tokens: 390
      cost: 0.0195
      timestamp: "2025-01-23T14:30:45Z"
      success: true
      
      metrics:
        relevance: 0.95
        completeness: 0.88
        accuracy: 0.92
        helpfulness: 0.90
  
  optimization:
    desired_outcome: |
      Provide comprehensive code reviews that identify specific issues,
      explain the problems clearly, and offer concrete solutions with
      improved code examples. Reviews should be structured and actionable.
    
    runs:
      - iteration: 1
        original_prompt: "Review this code: {{code}}"
        enhanced_prompt: |
          You are an expert code reviewer. Please review the following {{language}} code:
          
          ```{{language}}
          {{code}}
          ```
          
          Provide feedback on:
          1. Code quality and best practices
          2. Potential bugs or issues
          3. Performance improvements
          4. Security considerations
        
        gap_analysis:
          identified_gaps:
            - category: "specificity"
              description: "Prompt lacks specific review criteria"
              severity: 0.7
              suggested_fix: "Add structured review categories"
            
            - category: "format"
              description: "No output format specified"
              severity: 0.6
              suggested_fix: "Specify structured response format"
        
        test_results:
          - model: "gpt-4"
            improvement_score: 0.65
            metrics:
              relevance: 0.90
              completeness: 0.75
              structure: 0.60
        
        improvement_score: 0.65
        timestamp: "2025-01-23T14:35:12Z"
      
      - iteration: 2
        enhanced_prompt: |
          You are an expert code reviewer. Please review the following {{language}} code:
          
          ```{{language}}
          {{code}}
          ```
          
          Provide feedback on:
          1. Code quality and best practices
          2. Potential bugs or issues
          3. Performance improvements
          4. Security considerations
          
          Format your response as structured feedback with specific suggestions.
        
        improvement_score: 0.85
        timestamp: "2025-01-23T14:40:33Z"
        converged: true

  metadata:
    total_tests: 12
    total_optimizations: 2
    models_used: ["gpt-4", "claude-3-sonnet", "gpt-3.5-turbo"]
    session_duration: "1h 15m"
    auto_saved: true
```

### Configuration Management

```yaml
# .agentforge/config/workbench.yaml
workbench:
  default_models:
    - "gpt-4"
    - "claude-3-sonnet"
  
  auto_save:
    enabled: true
    interval: "30s"
  
  optimization:
    max_iterations: 10
    convergence_threshold: 0.05
    default_strategy: "balanced"
  
  ui:
    theme: "dark"
    tab_width: 4
    show_line_numbers: true
    syntax_highlighting: true
  
  providers:
    openai:
      api_key_env: "OPENAI_API_KEY"
      default_model: "gpt-4"
      max_tokens: 4096
      temperature: 0.7
    
    anthropic:
      api_key_env: "ANTHROPIC_API_KEY"
      default_model: "claude-3-sonnet-20240229"
      max_tokens: 4096
      temperature: 0.7
  
  export:
    default_format: "yaml"
    include_test_results: true
    include_optimization_history: true
```

## Performance Considerations

### TUI Optimization

```go
// Efficient state management for large sessions
type WorkbenchState struct {
    // Core state
    currentTab    TabType
    promptContent string
    variables     map[string]interface{}
    
    // Lazy-loaded data
    testResults   *LazyTestResults
    optHistory    *LazyOptimizationHistory
    
    // UI state
    viewport      ViewportState
    isDirty       bool
    lastSaved     time.Time
}

// LazyTestResults loads test results on demand
type LazyTestResults struct {
    loaded   bool
    results  []TestResult
    loader   func() ([]TestResult, error)
}

// Efficient rendering with viewport management
type ViewportState struct {
    width       int
    height      int
    scrollY     int
    visibleRows int
    needsRedraw bool
}
```

### Caching Strategy

```go
// Multi-level caching for performance
type CacheManager struct {
    // In-memory cache for active session
    sessionCache *sync.Map
    
    // Disk cache for test results
    resultCache *DiskCache
    
    // Model response cache
    modelCache *ModelResponseCache
}

// Cache test results to avoid redundant API calls
type ModelResponseCache struct {
    cache map[string]CachedResponse
    mutex sync.RWMutex
    ttl   time.Duration
}

type CachedResponse struct {
    Response  TestResult
    Timestamp time.Time
    Hash      string
}
```

## Error Handling & Resilience

### Graceful Degradation

```go
// Fallback strategies for different failure modes
type FallbackManager struct {
    offlineMode     bool
    mockProviders   []ModelProvider
    cachedResponses map[string]TestResult
}

// Handle API failures gracefully
func (f *FallbackManager) HandleProviderFailure(provider string, err error) error {
    switch {
    case isRateLimitError(err):
        return f.activateRateLimit(provider)
    case isNetworkError(err):
        return f.activateOfflineMode()
    case isAuthError(err):
        return f.promptForCredentials(provider)
    default:
        return f.logAndContinue(err)
    }
}
```

### Validation Framework

```go
// Comprehensive validation for all inputs
type ValidationEngine struct {
    rules map[string][]ValidationRule
}

type ValidationRule interface {
    Validate(value interface{}) error
    GetErrorMessage() string
}

// Built-in validation rules
type LengthRule struct {
    Min, Max int
}

type PatternRule struct {
    Pattern *regexp.Regexp
}

type CustomRule struct {
    Validator func(interface{}) error
    Message   string
}
```

This technical specification provides the detailed implementation guidance needed to build a robust, performant, and user-friendly TUI prompt workbench that meets all the requirements outlined in the main specification.