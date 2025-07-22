# MVP System Architecture

## High-Level Architecture

AgentForge MVP follows a **simple, local-first architecture** that builds on the existing MCP Planner foundation while adding GitHub integration for component sharing.

```
┌─────────────────────────────────────────────────────────────┐
│                    AgentForge Ecosystem                     │
├─────────────────────────────────────────────────────────────┤
│  GitHub Repositories (Distributed Component Sources)       │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
│  │ Tools Repo  │ │Prompts Repo │ │ Agents Repo │   ...    │
│  │ @commit123  │ │ @commit456  │ │ @commit789  │          │
│  └─────────────┘ └─────────────┘ └─────────────┘          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Local AgentForge                        │
├─────────────────────────────────────────────────────────────┤
│  CLI Interface (forge)                                     │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
│  │ repo mgmt   │ │ dev tools   │ │ sync engine │          │
│  └─────────────┘ └─────────────┘ └─────────────┘          │
├─────────────────────────────────────────────────────────────┤
│  Core Services                                             │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
│  │ Component   │ │ Dependency  │ │ Composition │          │
│  │ Resolver    │ │ Manager     │ │ Engine      │          │
│  └─────────────┘ └─────────────┘ └─────────────┘          │
├─────────────────────────────────────────────────────────────┤
│  Local Database (PostgreSQL)                              │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
│  │ Components  │ │ Repositories│ │ Compositions│          │
│  │ & Changes   │ │ & Sync      │ │ & Configs   │          │
│  └─────────────┘ └─────────────┘ └─────────────┘          │
├─────────────────────────────────────────────────────────────┤
│  Runtime Engine                                           │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
│  │ Agent       │ │ Tool        │ │ MCP Server  │          │
│  │ Factory     │ │ Provider    │ │ Manager     │          │
│  └─────────────┘ └─────────────┘ └─────────────┘          │
└─────────────────────────────────────────────────────────────┘