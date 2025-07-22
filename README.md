# AgentForge

A Git-native AI agent development platform that enables developers to build, share, and manage AI agent components through a collaborative ecosystem.

## Features

- 🌍 **Git-Native**: Components distributed via GitHub repositories
- ⚡ **Local-First**: Fast iteration with local storage
- 🎨 **Beautiful CLI**: Interactive interfaces for component management
- 🔄 **Simple Sync**: Easy push/pull with GitHub
- 🧩 **Three Component Types**: Tools, Prompts, and Agents
- 🤝 **Community-Driven**: 30-second component sharing
- 🏗️ **Modular Architecture**: Clean, testable, and extensible design
- 📊 **Structured Logging**: Comprehensive logging with Zap

## Quick Start

### Prerequisites

- Go 1.21 or later
- PostgreSQL database
- Google AI API key

### Installation

1. Clone the repository:
```bash
git clone https://github.com/denkhaus/mcp-planner.git
cd mcp-planner
```

2. Install dependencies:
```bash
make deps
```

3. Set up environment variables:
```bash
export GOOGLE_API_KEY="your-google-ai-api-key"
export DATABASE_URL="postgres://user:password@localhost/mcp_planner?sslmode=disable"
export LOG_LEVEL="info"
export PORT="8080"
export ENVIRONMENT="development"
```

4. Build the application:
```bash
make build
```

5. Start the interactive chat:
```bash
./bin/mcp-planner chat
```

## Usage

### Commands

- **Server Mode**: `./bin/mcp-planner server --port 8080`
- **Interactive Chat**: `./bin/mcp-planner chat --clear`
- **Version Info**: `./bin/mcp-planner version`
- **Help**: `./bin/mcp-planner --help`

### Chat Examples

```
You: What's the weather like in Chicago?
Assistant: The weather in Chicago is currently 64°F and sunny with light winds.

You: Create a task to review project documentation with high priority
Assistant: Task created successfully: {"id":"task-42","title":"review project documentation","priority":"high","status":"created"}

You: Help me plan a team meeting for next week
Assistant: I'll help you plan a team meeting. Here are some suggestions...
```

### Available Tools

- **Weather Information**: Get current weather for any city
- **Task Management**: Create and manage tasks with priorities
- **Planning Assistance**: AI-powered planning and organization help

## Development

### Setup Development Environment

```bash
make dev-setup
```

### Code Quality

```bash
# Run linter
make lint

# Fix linting issues
make lint-fix

# Format code
make fmt

# Run tests
make test

# Generate coverage report
make coverage

# Run security checks
make security

# Pre-commit checks
make pre-commit
```

### Project Structure

```
├── main.go                 # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── container/         # Dependency injection
│   ├── logger/            # Structured logging
│   ├── signals/           # Signal handling utilities
│   ├── types/             # Type definitions and interfaces
│   └── providers/         # Service providers
│       ├── agent.go       # AI agent provider
│       ├── prompt.go      # Prompt management provider
│       └── tool.go        # Tool execution provider
├── .golangci.yml          # Linting configuration
├── Makefile              # Build automation
└── .github/workflows/    # CI/CD pipelines
```

## Configuration

The application uses environment variables for configuration:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `GOOGLE_API_KEY` | Google AI API key | - | ✅ |
| `DATABASE_URL` | PostgreSQL connection string | - | ✅ |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | info | ❌ |
| `PORT` | HTTP server port | 8080 | ❌ |
| `ENVIRONMENT` | Deployment environment | development | ❌ |

### CLI Flag Overrides

CLI flags can override environment variables:

```bash
./bin/mcp-planner server --port 9090 --log-level debug
./bin/mcp-planner chat --log-level debug --clear
```

## Architecture

### Core Components

1. **Agent Provider**: Manages AI interactions using Google AI
2. **Tool Provider**: Handles tool registration and execution
3. **Prompt Provider**: Manages prompt templates and system prompts
4. **Configuration**: Environment-based config with CLI overrides
5. **Dependency Injection**: Clean service management with samber/do

### Design Principles

- **Interface-First**: All providers implement well-defined interfaces
- **Dependency Injection**: Clean separation of concerns
- **Error Handling**: Comprehensive error wrapping and logging
- **Testability**: Mockable interfaces and dependency injection
- **Observability**: Structured logging throughout the application

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes following the coding standards in `go-coding.md`
4. Run pre-commit checks: `make pre-commit`
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Code Standards

This project follows strict Go coding standards defined in `rules/go-coding.md`:

- **Error Handling**: Custom error types with proper wrapping
- **Naming Conventions**: Clear, descriptive names following Go conventions
- **Package Organization**: Single responsibility principle
- **Interface Design**: Small, focused interfaces
- **Dependency Injection**: Clean service management
- **Testing**: Comprehensive test coverage with testify
- **Documentation**: GoDoc comments for all exported functions

## CI/CD

The project includes comprehensive CI/CD pipelines:

- **Linting**: golangci-lint with extensive rule set
- **Testing**: Unit tests with coverage reporting
- **Security**: Vulnerability scanning with gosec
- **Build**: Multi-platform binary generation

## License

[Add your license here]

## Support

For questions, issues, or contributions, please:

1. Check the [Issues](https://github.com/denkhaus/mcp-planner/issues) page
2. Create a new issue with detailed information
3. Follow the contributing guidelines

---

Built with ❤️ using Go and Git-native workflows