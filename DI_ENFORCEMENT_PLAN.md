# Dependency Injection Enforcement Plan

## 🎯 **Vision**
Establish and enforce consistent dependency injection patterns throughout the AgentForge codebase, ensuring maintainable, testable, and loosely coupled architecture following SOLID principles.

## 📊 **Current State Assessment**

### ✅ **Well-Implemented DI Areas**
- **Container System**: `internal/container/container.go` with samber/do
- **Database Layer**: Complete interface-first design with private implementations
- **Provider System**: Proper DI with named providers to avoid circular dependencies
- **Session Factory**: Complex DI with LLMService injection
- **Startup System**: Fine-grained options with composable service initialization

### ⚠️ **Areas Needing DI Enforcement**
- **Commands Package**: Mixed DI patterns, some direct instantiation
- **TUI Components**: Inconsistent interface usage
- **Tools Package**: Some direct dependencies
- **Services Layer**: Incomplete interface abstractions
- **Testing**: Mock implementations not always following interfaces

## 🏗️ **DI Architecture Principles**

### Core Pattern
```
DI Container -> Interfaces -> Private Implementations
```

### Interface-First Design Rules
1. **All public constructors return interfaces, never implementations**
2. **All implementations are private (lowercase)**
3. **Dependencies are injected as interfaces, not concrete types**
4. **Use samber/do for service registration and resolution**

### Naming Conventions
```go
// ✅ Correct Pattern
type ServiceInterface interface { ... }
type privateImplementation struct { ... }
func NewService() ServiceInterface { return &privateImplementation{} }

// ❌ Anti-Pattern
type ServiceStruct struct { ... }
func NewServiceStruct() *ServiceStruct { return &ServiceStruct{} }
```

## 📋 **Enforcement Checklist**

### Phase 1: Interface Standardization
- [ ] **Audit all public constructors** - Ensure they return interfaces
- [ ] **Create missing interfaces** - For services without proper abstractions
- [ ] **Make implementations private** - Convert public structs to private
- [ ] **Update DI registrations** - Ensure all services are properly registered

### Phase 2: Dependency Cleanup
- [ ] **Remove direct instantiation** - Replace `NewService()` calls with DI resolution
- [ ] **Fix circular dependencies** - Use named providers where needed
- [ ] **Standardize injection patterns** - Consistent `do.Invoke` usage
- [ ] **Update error handling** - Proper error propagation from DI resolution

### Phase 3: Testing Alignment
- [ ] **Create interface mocks** - For all service interfaces
- [ ] **Update test dependencies** - Use mocks instead of real implementations
- [ ] **Verify test isolation** - Ensure tests don't depend on DI container state
- [ ] **Add DI container tests** - Verify service registration and resolution

## 🔧 **Implementation Standards**

### Service Registration Pattern
```go
// Standard service registration
do.Provide(injector, func(i *do.Injector) (ServiceInterface, error) {
    dep := do.MustInvoke[DependencyInterface](i)
    logger := do.MustInvoke[*zap.Logger](i)
    return NewService(dep, logger), nil
})

// Named providers for circular dependency prevention
do.ProvideNamed(injector, "serviceName", func(i *do.Injector) (ServiceInterface, error) {
    return NewService(), nil
})
```

### Service Resolution Pattern
```go
// Safe resolution with error handling
service, err := do.Invoke[ServiceInterface](injector)
if err != nil {
    return fmt.Errorf("failed to resolve service: %w", err)
}

// Must invoke for critical services (panics on failure)
service := do.MustInvoke[ServiceInterface](injector)

// Named resolution
service := do.MustInvokeNamed[ServiceInterface](injector, "serviceName")
```

### Constructor Pattern
```go
// ✅ Correct: Returns interface, accepts interface dependencies
func NewService(dep DependencyInterface, logger *zap.Logger) ServiceInterface {
    return &privateService{
        dep:    dep,
        logger: logger.WithPackage("service"),
    }
}

// ❌ Incorrect: Returns concrete type, accepts concrete dependencies
func NewService(dep *ConcreteDependency) *ConcreteService {
    return &ConcreteService{dep: dep}
}
```

## 📁 **Package-by-Package Enforcement**

### High Priority Packages
1. **`internal/commands/`** - Standardize startup context usage
2. **`internal/tui/`** - Implement proper interface abstractions
3. **`internal/services/`** - Create missing service interfaces
4. **`internal/tools/`** - Convert to interface-based DI

### Medium Priority Packages
1. **`internal/agents/`** - Implement agent service interface
2. **`internal/session/`** - Verify session factory DI patterns
3. **`internal/decorators/`** - Ensure decorator pattern compatibility

### Low Priority Packages
1. **`internal/config/`** - Already well-structured
2. **`internal/database/`** - Excellent DI implementation (reference example)
3. **`internal/providers/`** - Good DI patterns with named providers

## 🧪 **Testing Strategy**

### Mock Generation
```go
// Generate mocks for all interfaces
//go:generate mockgen -source=interfaces.go -destination=mocks/mock_interfaces.go

// Use in tests
func TestService(t *testing.T) {
    mockDep := mocks.NewMockDependencyInterface(ctrl)
    service := NewService(mockDep, zap.NewNop())
    // Test implementation
}
```

### DI Container Testing
```go
func TestDIContainer(t *testing.T) {
    injector := container.Setup(&config.Config{})
    defer injector.Shutdown()
    
    // Verify all services can be resolved
    _, err := do.Invoke[ServiceInterface](injector)
    assert.NoError(t, err)
}
```

## 🚨 **Anti-Patterns to Eliminate**

### Direct Instantiation
```go
// ❌ Anti-pattern
service := &ConcreteService{}

// ✅ Correct
service := do.MustInvoke[ServiceInterface](injector)
```

### Concrete Dependencies
```go
// ❌ Anti-pattern
func NewService(db *sql.DB) ServiceInterface

// ✅ Correct
func NewService(client DatabaseClient) ServiceInterface
```

### Public Implementations
```go
// ❌ Anti-pattern
type PublicService struct { ... }

// ✅ Correct
type ServiceInterface interface { ... }
type privateService struct { ... }
```

## 📈 **Success Metrics**

### Code Quality Indicators
- [ ] **100% interface coverage** - All services have proper interfaces
- [ ] **Zero direct instantiation** - All services resolved through DI
- [ ] **Clean dependency graph** - No circular dependencies
- [ ] **Consistent patterns** - Same DI patterns across all packages

### Testing Improvements
- [ ] **Mockable interfaces** - All dependencies can be mocked
- [ ] **Isolated tests** - Tests don't depend on real implementations
- [ ] **Fast test execution** - No database/network dependencies in unit tests
- [ ] **High test coverage** - Improved testability leads to better coverage

### Maintainability Benefits
- [ ] **Easy refactoring** - Interface changes don't break consumers
- [ ] **Pluggable implementations** - Easy to swap service implementations
- [ ] **Clear dependencies** - Explicit dependency declarations
- [ ] **Reduced coupling** - Services depend on abstractions, not concretions

## 🔄 **Enforcement Process**

### Daily Development
1. **Pre-commit checks** - Verify DI patterns in changed files
2. **Code review focus** - Ensure new code follows DI principles
3. **Refactoring opportunities** - Identify and fix anti-patterns

### Weekly Reviews
1. **Package audits** - Systematic review of package DI compliance
2. **Interface coverage** - Track progress on interface implementation
3. **Dependency analysis** - Identify and resolve circular dependencies

### Monthly Assessments
1. **Architecture review** - Evaluate overall DI architecture health
2. **Performance impact** - Measure DI overhead and optimization opportunities
3. **Developer feedback** - Gather team input on DI patterns and tooling

## 🛠️ **Tools and Automation**

### Static Analysis
```bash
# Check for direct instantiation patterns
grep -r "New.*(" internal/ | grep -v "interface"

# Find public struct types (should be interfaces)
grep -r "type.*struct" internal/ | grep -v "private"

# Verify DI registration completeness
grep -r "do.Provide" internal/container/
```

### Linting Rules
```yaml
# .golangci.yml additions for DI enforcement
linters:
  enable:
    - interfacer      # Suggest interface usage
    - prealloc        # Find slice preallocation opportunities
    - unconvert       # Remove unnecessary type conversions
```

### Build Integration
```makefile
# Makefile targets for DI enforcement
.PHONY: check-di
check-di:
	@echo "Checking DI patterns..."
	@./scripts/check-di-patterns.sh

.PHONY: enforce-di
enforce-di: check-di
	@echo "DI patterns verified ✅"
```

## 📚 **Reference Implementation**

The `internal/database/` package serves as the gold standard for DI implementation:

- **Interfaces**: Clear service abstractions
- **Private implementations**: All structs are private
- **DI registration**: Proper samber/do usage
- **Error handling**: Comprehensive error propagation
- **Testing**: Mockable interfaces with proper test coverage

Use this package as a reference when implementing DI patterns in other areas of the codebase.

## 🎯 **Next Steps**

1. **Start with Commands Package** - High impact, visible improvements
2. **Create Missing Interfaces** - Systematic interface creation
3. **Update DI Registrations** - Ensure all services are properly registered
4. **Implement Testing Strategy** - Mock generation and test updates
5. **Monitor and Iterate** - Continuous improvement of DI patterns

---

**Remember**: Dependency Injection is not just about tools and patterns—it's about creating a maintainable, testable, and flexible architecture that enables rapid development and easy refactoring.