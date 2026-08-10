# Kanna Architecture

**Status:** Draft  
**Target:** August 31, 2026 MVP  
**Architecture style:** Modular monolith  
**Primary language:** Go  
**HTTP framework:** Gin  

## 1. Vision

Kanna is a personal AI assistant designed to securely connect to the user's digital life, home environment, and external services. It should understand natural-language requests, use tools to perform actions, remember useful context, and eventually automate recurring tasks.

The initial goal is not to build a fully autonomous assistant. The goal is to establish a clean, extensible foundation that can evolve into one.

## 2. Goals

For the August 2026 MVP, Kanna should:

- Accept natural-language requests through an HTTP API and simple UI.
- Use an LLM to understand requests and generate responses.
- Support a provider abstraction so the LLM implementation can change.
- Initially support a local LLM through Ollama.
- Provide a generic tool abstraction and tool registry.
- Execute simple tools and at least one real external integration.
- Persist conversations and basic memories.
- Use PostgreSQL for persistent application state.
- Provide structured logging, configuration, tests, and graceful shutdown.
- Keep security and user approval in mind before sensitive actions.

## 3. Non-goals for the August MVP

The following are explicitly deferred:

- Voice/wake-word support.
- Mobile application.
- Smart-home device control.
- Banking/financial integrations.
- Camera/vision workflows.
- Multi-agent orchestration.
- Kafka or distributed event-driven infrastructure.
- Kubernetes deployment.
- Autonomous background agents.
- Complex RAG pipelines.
- Knowledge graphs.
- Production-grade multi-user SaaS functionality.

These can be evaluated after the MVP proves the core architecture.

## 4. Architecture Principles

### 4.1 Modular monolith first

Kanna will initially be a modular monolith. Components will have clear interfaces and package boundaries, but they will run in a single Go process.

This allows the project to prioritize capabilities over distributed-systems complexity. Components can be extracted into services later if there is a concrete scaling, isolation, or operational reason.

### 4.2 Provider independence

LLM providers and external integrations should be accessed through interfaces rather than being coupled directly to application logic.

### 4.3 Tool-driven capabilities

The agent should interact with external systems through explicit tools. A tool should expose its name, description, input contract, execution behavior, and result.

### 4.4 Local-first where practical

The initial LLM runtime will use Ollama and a local model. Sensitive information should remain local whenever practical.

### 4.5 Human control

Actions that can have meaningful consequences should eventually support policy checks and explicit user confirmation.

### 4.6 Observability from the beginning

Tool calls, errors, request IDs, and important agent decisions should be observable without logging secrets or sensitive user data.

## 5. High-level Architecture

```text
                         ┌──────────────────────┐
                         │      User / UI       │
                         └──────────┬───────────┘
                                    │
                                    ▼
                         ┌──────────────────────┐
                         │      Gin API         │
                         │  HTTP / REST Layer   │
                         └──────────┬───────────┘
                                    │
                                    ▼
                         ┌──────────────────────┐
                         │ Conversation Service │
                         └──────────┬───────────┘
                                    │
                                    ▼
                         ┌──────────────────────┐
                         │    Agent Runtime     │
                         └──────┬───────┬───────┘
                                │       │
                     ┌──────────┘       └──────────┐
                     ▼                             ▼
              ┌──────────────┐              ┌──────────────┐
              │ LLM Provider │              │ Tool Registry │
              └──────┬───────┘              └──────┬───────┘
                     │                             │
                     ▼                       ┌─────┴──────────┐
                ┌─────────┐                  │                │
                │ Ollama  │              Calculator      GitHub
                └─────────┘                  │                │
                                             └─────┬──────────┘
                                                   │
                                                   ▼
                                            External APIs

                         ┌──────────────────────┐
                         │   Memory / Storage   │
                         │     PostgreSQL       │
                         └──────────────────────┘
```

## 6. Core Components

### 6.1 API

Gin handles HTTP transport, request validation, middleware, routing, and response serialization.

Initial endpoints:

```text
GET  /health
POST /api/v1/chat
GET  /api/v1/conversations
GET  /api/v1/memories
```

### 6.2 Conversation Service

Responsible for:

- Creating and retrieving conversations.
- Persisting messages.
- Passing user input to the agent runtime.
- Returning the final assistant response.

The conversation service should not contain provider-specific LLM logic.

### 6.3 Agent Runtime

The agent runtime coordinates reasoning and tool execution.

Initial flow:

```text
User request
    │
    ▼
Agent Runtime
    │
    ▼
LLM
    │
    ├── Direct answer ───────────────► Response
    │
    └── Tool required
             │
             ▼
        Tool Registry
             │
             ▼
        Execute Tool
             │
             ▼
        Tool Result
             │
             ▼
            LLM
             │
             ▼
          Response
```

The first version should remain intentionally simple. More advanced planning and multi-step execution can be added later.

### 6.4 LLM Provider

Define an abstraction similar to:

```go
type LLM interface {
    Generate(ctx context.Context, request LLMRequest) (LLMResponse, error)
}
```

Initial implementation:

```text
LLM
└── OllamaProvider
```

Future implementations can include other cloud or local providers without changing the agent runtime.

### 6.5 Tool System

Tools represent capabilities available to the agent.

Conceptually:

```go
type Tool interface {
    Name() string
    Description() string
    Execute(ctx context.Context, input ToolInput) (ToolOutput, error)
}
```

Initial tools:

- Calculator
- Date/time
- Weather
- GitHub

The tool registry is responsible for discovering available tools and routing execution to the correct implementation.

### 6.6 Memory

Memory should initially distinguish between:

- Conversation history.
- Explicit user memories/facts.
- Tool execution history.

Initial interface:

```go
type Memory interface {
    Store(ctx context.Context, memory MemoryEntry) error
    Search(ctx context.Context, query string) ([]MemoryEntry, error)
}
```

The first implementation can use PostgreSQL without introducing a vector database immediately.

### 6.7 Persistence

PostgreSQL is the initial persistent datastore.

Initial conceptual entities:

```text
users
conversations
messages
memories
tool_executions
```

The schema should remain minimal until actual requirements emerge.

## 7. Suggested Go Package Structure

```text
kanna/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── agent/
│   ├── llm/
│   ├── memory/
│   ├── tools/
│   ├── conversation/
│   ├── api/
│   ├── handlers/
│   ├── middleware/
│   ├── config/
│   └── models/
├── pkg/
├── configs/
├── docs/
├── scripts/
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

The exact package boundaries can change as implementation experience exposes better boundaries.

## 8. Security

Security is a first-class concern because Kanna will eventually have access to personal accounts and devices.

Initial requirements:

- Never commit API keys or credentials.
- Use environment variables or a secrets mechanism.
- Do not log tokens, passwords, or full sensitive payloads.
- Validate tool inputs.
- Give tools explicit permissions.
- Require confirmation for destructive or high-impact operations.
- Keep external connector credentials isolated from the LLM prompt where possible.
- Treat LLM-generated tool arguments as untrusted input.

Financial integrations, smart locks, shell execution, and other high-impact tools should not be added until a proper authorization and approval model exists.

## 9. Observability

Initial observability:

- Structured logs.
- Request IDs.
- Tool execution logs.
- LLM latency.
- External API latency.
- Error tracking.
- Basic health endpoint.

Do not log sensitive user data by default.

## 10. Deployment

For the August MVP:

```text
Developer machine
    │
    ├── Kanna Go API
    ├── Ollama
    └── PostgreSQL
```

Docker Compose can later package the components into a reproducible local environment.

Kubernetes is intentionally deferred.

## 11. Future Architecture

Once the modular monolith becomes constrained by actual requirements, components may be extracted.

Potential future services:

```text
API Gateway
    │
    ├── Conversation Service
    ├── Agent Runtime
    ├── Memory Service
    ├── Tool/Connector Service
    └── Scheduler
```

Extraction should be driven by real operational requirements rather than anticipated scale.

## 12. August MVP Definition of Done

The MVP is considered successful when a user can:

1. Send a natural-language request to Kanna.
2. Have Kanna process the request using a local LLM.
3. Have the agent decide whether a tool is needed.
4. Execute a registered tool.
5. Incorporate the tool result into the response.
6. Persist the conversation.
7. Store and retrieve at least basic user memory.
8. Interact with Kanna through a simple UI or API.
9. Run the complete system locally with documented setup instructions.

## 13. Key Architectural Decisions

### ADR-001: Modular monolith

Kanna will begin as a modular monolith to maximize development velocity and minimize operational overhead.

### ADR-002: Go + Gin

Go is the primary backend language and Gin provides the initial HTTP layer.

### ADR-003: Ollama as initial LLM provider

Ollama provides a local LLM runtime suitable for development and privacy-sensitive personal use.

### ADR-004: PostgreSQL as initial datastore

PostgreSQL provides reliable persistence for conversations, memories, and application metadata without adding unnecessary infrastructure.

### ADR-005: Explicit tool abstraction

External capabilities will be exposed through a tool interface so that the agent runtime remains independent from individual integrations.

## 14. Open Questions

These should be resolved as the project develops:

- Which LLM/tool-calling protocol should Kanna standardize on?
- Should tools eventually use MCP?
- How should tool permissions be modeled?
- What actions require confirmation?
- How should long-term memory be ranked and retrieved?
- When is a vector database justified?
- What scheduling model should be used?
- How should secrets be stored?
- When should individual components become separate services?
