# Technical Stack

> Last Updated: 2025-01-23
> Version: 1.0.0

## Core Technologies

### Application Framework
- **Go 1.24.0** - Primary language with modern generics support
- **urfave/cli/v2** - Command-line interface framework
- **Bubble Tea + Lipgloss** - Terminal user interface framework

### Database System
- **SQLite** - Local-first database for fast iteration
- **Ent ORM** - Type-safe Go ORM with code generation
- **84,559 lines** - Comprehensive schema implementation

### Architecture Patterns
- **Dependency Injection** - samber/do container for clean architecture
- **Interface-First Design** - All services implement interfaces
- **Provider Pattern** - Pluggable tool and prompt providers

### Integration Layer
- **Git Integration** - go-git/go-git/v5 for repository operations
- **GitHub API** - google/go-github/v57 for ecosystem integration
- **MCP Protocol** - denkhaus/mcp-server-adapter for tool providers

### Development Tools
- **Logging** - Uber Zap structured logging
- **Testing** - Go standard testing with comprehensive coverage
- **Build System** - Make-based build pipeline
- **Linting** - golangci-lint with comprehensive rules

## Data Architecture

### Entity Schema
- **Agents** - Conversational/Task-oriented/Specialized/Composite types
- **Prompts** - System/User/Assistant/Function/Template types
- **Tools** - MCP/HTTP/Binary/Function execution types
- **Repositories** - Git-based component distribution

### Storage Strategy
- **Local SQLite** - Fast queries and offline capability
- **Git Repositories** - Distributed component storage
- **JSON Metadata** - Flexible configuration and schemas

## Deployment Architecture

### Application Hosting
- **Local Binary** - Single executable with embedded assets
- **Cross-Platform** - Windows, macOS, Linux support

### Database Hosting
- **Local SQLite File** - No external dependencies
- **Backup Strategy** - Git-based versioning

### Asset Hosting
- **GitHub Repositories** - Component distribution
- **Local Cache** - Fast component access

### Deployment Solution
- **Go Install** - Direct installation from source
- **GitHub Releases** - Binary distribution
- **Package Managers** - Future Homebrew/Chocolatey support

## Code Repository
- **URL:** https://github.com/denkhaus/agentforge
- **License:** Open source (license TBD)
- **CI/CD:** GitHub Actions with comprehensive testing