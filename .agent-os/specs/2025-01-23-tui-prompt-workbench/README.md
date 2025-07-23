# TUI Prompt Workbench Enhancement Specification

> **Created:** 2025-01-23  
> **Status:** Planning  
> **Priority:** High  
> **Roadmap Reference:** Phase 1 - MVP Foundation  

## Overview

Transform the existing basic TUI prompt workbench into a comprehensive prompt development environment that enables interactive editing, variable management, multi-model testing, and AI-powered optimization with feedback loops.

## Context

### Current State Analysis

**✅ Existing Infrastructure:**
- Basic TUI framework with Bubble Tea implemented (`internal/tui/workbench_v3.go`)
- Four-tab structure: Editor, Variables, Test, Optimize (`internal/tui/models/workbench.go`)
- Prompt new command with TUI integration placeholder (`internal/commands/prompt_new.go`)
- DI container ready for TUI manager injection
- Comprehensive planning document (`PROMPT_WORKBENCH_PLAN.md`)

**❌ Current Limitations:**
- Test and Optimize tabs show placeholder content only
- No actual model integration for testing
- No AI optimization engine implementation
- No variable management functionality
- No session persistence
- No feedback loop system

### Problem Statement

After running `forge prompt new --name example`, users should see a fully functional TUI workbench where they can:
1. **Edit prompts** with live preview and syntax highlighting
2. **Manage variables** with type validation and test data sets
3. **Test prompts** across multiple AI models with comparison
4. **Optimize prompts** using AI feedback loops comparing desired vs actual outcomes
5. **Iterate enhancement** in configurable steps until convergence

## Specification

### User Experience Flow

```bash
# 1. Create new prompt
forge prompt new --name code-reviewer
# → Creates basic structure
# → Launches enhanced TUI workbench automatically

# 2. Workbench opens with 4 tabs:
# Tab 1: Editor - Rich text editing with live preview
# Tab 2: Variables - Manage prompt variables and test sets  
# Tab 3: Test - Multi-model testing with comparison
# Tab 4: Optimize - AI-powered optimization feedback loop
```

### Core Features

#### 1. Enhanced Editor Tab
- **Rich Text Editing:** Multi-line prompt editor with syntax highlighting
- **Live Preview:** Real-time preview with variable substitution
- **Template Support:** Common prompt patterns and structures
- **Validation:** Real-time prompt validation and suggestions

#### 2. Variables Management Tab
- **Variable Definition:** Name, type, description, default values
- **Test Data Sets:** Multiple test scenarios with different variable values
- **Type Validation:** String, number, boolean, array validation
- **Import/Export:** Load variables from JSON/YAML files

#### 3. Multi-Model Testing Tab
- **Model Selection:** Choose from multiple AI providers (OpenAI, Anthropic, etc.)
- **Parallel Testing:** Run same prompt across multiple models simultaneously
- **Response Comparison:** Side-by-side comparison of model outputs
- **Performance Metrics:** Response time, token usage, cost analysis
- **Test History:** Track all test runs with timestamps and results

#### 4. AI Optimization Tab
- **Desired Outcome Definition:** User specifies expected output characteristics
- **Gap Analysis:** AI compares actual vs desired outcomes
- **Enhancement Suggestions:** AI proposes specific prompt improvements
- **Iterative Testing:** Automated testing of enhanced prompts
- **Convergence Detection:** Stop optimization when improvement plateaus
- **Optimization History:** Track all optimization iterations

### Technical Architecture

#### Enhanced TUI Models

```go
// Enhanced test model with actual functionality
type TestModel struct {
    promptContent   string
    variables       map[string]interface{}
    selectedModels  []string
    testResults     []TestResult
    isRunning       bool
    modelProviders  []types.ModelProvider
}

// Enhanced optimize model with AI integration
type OptimizeModel struct {
    promptContent     string
    desiredOutcome    string
    optimizationRuns  []OptimizationRun
    currentIteration  int
    maxIterations     int
    optimizer         PromptOptimizer
    isOptimizing      bool
}

// New optimization types
type OptimizationRun struct {
    Iteration       int
    OriginalPrompt  string
    EnhancedPrompt  string
    TestResults     []TestResult
    ImprovementScore float64
    Suggestions     []string
}

type PromptOptimizer interface {
    AnalyzeGap(desired, actual string) OptimizationSuggestion
    EnhancePrompt(prompt string, suggestion OptimizationSuggestion) string
    EvaluateImprovement(before, after TestResult) float64
}
```

#### Model Integration

```go
// Enhanced model provider interface
type ModelProvider interface {
    Name() string
    TestPrompt(prompt string, variables map[string]interface{}) (TestResult, error)
    GetCapabilities() ModelCapabilities
    EstimateCost(prompt string) CostEstimate
}

type TestResult struct {
    ModelName     string
    Response      string
    ResponseTime  time.Duration
    TokensUsed    int
    Cost          float64
    Timestamp     time.Time
    Success       bool
    Error         error
}
```

#### Session Persistence

```yaml
# .agentforge/workbench/sessions/prompt-name.yaml
session:
  prompt:
    name: "code-reviewer"
    content: "Review this code for..."
    variables:
      - name: "code"
        type: "string"
        default: ""
      - name: "language"
        type: "string"
        default: "go"
  
  test_results:
    - model: "gpt-4"
      response: "..."
      metrics: {...}
  
  optimization:
    desired_outcome: "Detailed code review with suggestions"
    iterations: 3
    history: [...]
```

### Implementation Plan

#### Phase 1: Enhanced Testing (Week 1)
- Implement actual model integration in TestModel
- Add parallel testing across multiple models
- Create response comparison UI
- Add performance metrics tracking

#### Phase 2: AI Optimization Engine (Week 2)
- Implement PromptOptimizer interface
- Create optimization feedback loop
- Add gap analysis functionality
- Implement iterative enhancement

#### Phase 3: Advanced Features (Week 3)
- Add session persistence
- Implement variable management
- Create template system
- Add export capabilities

#### Phase 4: Polish & Integration (Week 4)
- Performance optimization
- Error handling and validation
- Integration testing
- Documentation and examples

### Success Criteria

#### Functional Requirements
- [ ] Users can edit prompts with live preview
- [ ] Variables can be defined and managed with validation
- [ ] Prompts can be tested across 3+ AI models simultaneously
- [ ] AI optimization improves prompts in 3-5 iterations
- [ ] Sessions persist and can be resumed
- [ ] Results can be exported in multiple formats

#### Performance Requirements
- [ ] TUI responds in <100ms for all interactions
- [ ] Model testing completes in <30 seconds
- [ ] Optimization iterations complete in <60 seconds
- [ ] Session save/load completes in <5 seconds

#### User Experience Requirements
- [ ] Time to first successful test: <2 minutes
- [ ] Optimization shows measurable improvement: >50%
- [ ] Interface is intuitive without documentation
- [ ] Error messages are clear and actionable

### Dependencies

#### External Dependencies
- AI model provider APIs (OpenAI, Anthropic, etc.)
- API keys and authentication setup
- Network connectivity for model testing

#### Internal Dependencies
- Enhanced DI container with model providers
- Prompt service with persistence capabilities
- Configuration management for API keys
- Error handling and logging framework

### Risks and Mitigations

#### Technical Risks
- **Model API Rate Limits:** Implement request queuing and retry logic
- **TUI Performance:** Use efficient rendering and state management
- **Session Corruption:** Add backup and recovery mechanisms

#### User Experience Risks
- **Complexity Overload:** Progressive disclosure of advanced features
- **Learning Curve:** Provide guided tutorials and examples
- **Model Costs:** Clear cost estimation and warnings

### Testing Strategy

#### Unit Testing
- Individual TUI model components
- Optimization algorithm logic
- Model provider integrations
- Session persistence functionality

#### Integration Testing
- End-to-end workbench workflows
- Multi-model testing scenarios
- Optimization feedback loops
- Error handling and recovery

#### User Testing
- Prompt development workflows
- Optimization effectiveness
- Interface usability
- Performance under load

This specification transforms the basic TUI workbench into a comprehensive prompt development environment that significantly enhances developer productivity and prompt quality through AI-powered optimization and multi-model testing capabilities.