# AgentForge MVP - Product Requirements Documentation

## Overview

AgentForge MVP is a Git-native component sharing system for AI agents, built on the existing MCP Planner foundation. We focus on making component sharing delightful and simple, leveraging GitHub as our distribution mechanism.

## MVP Documentation Structure

### Core MVP Documents
- [01-vision-and-goals.md](./01-vision-and-goals.md) - MVP vision and focused goals
- [02-system-architecture.md](./02-system-architecture.md) - Simplified local-first architecture
- [03-component-ecosystem.md](./03-component-ecosystem.md) - Basic component system (Tools, Prompts, Agents)

### MVP Implementation
- [04-database-schema.md](./04-database-schema.md) - Essential database tables only
- [05-cli-interface.md](./05-cli-interface.md) - 11 core CLI commands
- [06-sync-engine.md](./06-sync-engine.md) - Simple GitHub API integration
- [07-component-standards.md](./07-component-standards.md) - Component YAML format
- [09-local-development.md](./09-local-development.md) - Basic development workflow
- [13-migration-strategy.md](./13-migration-strategy.md) - MCP Planner to Forge migration
- [14b-phase1-foundation.md](./14b-phase1-foundation.md) - 6-week MVP implementation plan

### Post-MVP Features (Postponed)
- [08-dependency-resolution.md](./08-dependency-resolution.md) - **POSTPONED**: Complex dependency management
- [10-composition-system.md](./10-composition-system.md) - **POSTPONED**: Advanced agent composition
- [11-security-sandboxing.md](./11-security-sandboxing.md) - **POSTPONED**: Enterprise security features
- [12-ecosystem-marketplace.md](./12-ecosystem-marketplace.md) - **POSTPONED**: Custom marketplace (use GitHub)
- [14a-roadmap-overview.md](./14a-roadmap-overview.md) - **POSTPONED**: Long-term roadmap
- [14c-phase2-git-ecosystem.md](./14c-phase2-git-ecosystem.md) - **POSTPONED**: Advanced Git features
- [14d-phase3-collaboration.md](./14d-phase3-collaboration.md) - **POSTPONED**: Team collaboration features
- [14e-success-metrics.md](./14e-success-metrics.md) - **POSTPONED**: Enterprise metrics
- [15-api-reference.md](./15-api-reference.md) - **POSTPONED**: Complete API documentation

## MVP Quick Start

```bash
# Install Forge
go install github.com/denkhaus/agentforge/cmd/forge@latest

# Initialize local forge system
forge init

# Configure GitHub token
forge config

# Discover components
forge component ls --type tool --min-stars 10

# Pull a component
forge component pull --type tool github.com/user/awesome-weather-tool

# Create new component
forge component new --type prompt

# Push to share with community
forge component push --type prompt
```

## MVP Key Features

- **🌍 Git-Native**: Components distributed via GitHub repositories
- **⚡ Local-First**: Fast iteration with PostgreSQL storage
- **🎨 Beautiful TUI**: Interactive interfaces using bubbletea
- **🔄 Simple Sync**: Easy push/pull with GitHub
- **🧩 Three Component Types**: Tools, Prompts, and Agents
- **🤝 Community-Driven**: 30-second component sharing

## MVP Target Audience

- **Individual Developers**: Building and sharing AI components
- **Open Source Community**: Contributing to AI agent ecosystem
- **Early Adopters**: Developers wanting Git-native AI workflows

## MVP Success Metrics

- **Time to First Success**: <5 minutes from install to working component
- **Sharing Friction**: <30 seconds to share a component
- **Community Growth**: 50+ components by month 3, 100+ developers by month 6
- **User Delight**: "I can't believe how easy this is" feedback