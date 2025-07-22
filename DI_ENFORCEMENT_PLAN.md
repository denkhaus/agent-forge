# Dependency Injection Enforcement Plan

## 🎯 **Current DI Violations**

### 1. **Direct Constructor Calls in startup/context.go**
- `database.NewManager(cfg)` - Line 107
- `git.NewClient(log)` - Line 125  
- `github.NewClient("")` - Line 136

### 2. **Direct Constructor Calls in database/manager.go**
- `NewClient(dbConfig)` - Line 41
- `NewRepositoryService(client)` - Line 46
- `NewConfigService(client)` - Line 47

### 3. **Missing DI Registration**
- DatabaseManager not registered in container
- GitClient not registered in container
- GitHubClient not registered in container
- DatabaseClient not registered in container
- RepositoryService not registered in container
- ConfigService not registered in container

## 🔧 **Enforcement Strategy**

### Phase 1: Container Enhancement
1. Register all missing services in container.go
2. Create proper interfaces for all services
3. Add factory functions that use DI

### Phase 2: Startup Context Refactoring
1. Replace direct constructor calls with DI lookups
2. Remove service initialization logic from startup context
3. Make startup context purely a DI consumer

### Phase 3: Database Layer Refactoring
1. Register database services in container
2. Remove direct constructor calls from manager
3. Use DI for service creation

### Phase 4: Interface Enforcement
1. Ensure all implementations are private
2. Expose only interfaces publicly
3. Add interface compliance checks

## 🚀 **Implementation Order**

1. **Update container.go** - Add all missing service registrations
2. **Create service interfaces** - Define proper contracts
3. **Refactor startup/context.go** - Use DI instead of direct calls
4. **Refactor database/manager.go** - Use DI for service creation
5. **Update all NewXXX functions** - Make them DI-aware
6. **Add compliance tests** - Ensure DI patterns are followed

## ✅ **Success Criteria**

- [ ] No direct constructor calls outside of DI container
- [ ] All services registered in container
- [ ] All implementations are private
- [ ] All public APIs use interfaces
- [ ] Startup context uses only DI lookups
- [ ] Database manager uses DI for service creation
- [ ] Build and tests pass
- [ ] Git commit with working state