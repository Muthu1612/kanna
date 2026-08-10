# Kanna — August 2026 Sprint

**Sprint deadline:** August 31, 2026  
**Available development time:** Weekends  
**Sprint type:** MVP / Foundation  
**Primary objective:** Build a working agentic assistant foundation rather than a collection of disconnected integrations.

## 1. Sprint Goal

By August 31, Kanna should be able to receive a natural-language request, use a local LLM to reason about the request, invoke a registered tool when necessary, persist conversation/memory state, and return a useful response.

The system should run locally and be understandable enough that future integrations can be added without redesigning the entire application.

## 2. Success Criteria

The sprint is successful if the following end-to-end flow works:

```text
User
  │
  ▼
Kanna API / UI
  │
  ▼
Conversation Service
  │
  ▼
Agent Runtime
  │
  ▼
Local LLM (Ollama)
  │
  ├── Answer directly
  │
  └── Call Tool
        │
        ▼
     Tool Result
        │
        ▼
       LLM
        │
        ▼
    Final Response
        │
        ▼
 PostgreSQL persistence
```

## 3. Scope

### P0 — Must Have

- Go project foundation.
- Gin HTTP server.
- `/health` endpoint.
- Configuration management.
- Structured logging.
- Graceful shutdown.
- Basic unit tests.
- LLM abstraction.
- Ollama provider.
- Agent runtime.
- Tool interface.
- Tool registry.
- At least 2 simple tools.
- At least 1 external integration.
- PostgreSQL persistence.
- Conversation persistence.
- Basic memory persistence/retrieval.
- `/api/v1/chat`.
- Basic local UI or API-based demo.
- Docker/Compose local setup.
- Architecture documentation.
- Setup documentation.

### P1 — Should Have

- Weather tool.
- GitHub integration.
- Conversation history endpoint.
- Memory management endpoint.
- Tool execution tracing.
- Request IDs.
- Better error handling.
- Basic authentication boundary for future use.

### P2 — Nice to Have

- Streaming responses.
- WebSocket/SSE support.
- Scheduled tasks.
- More sophisticated memory retrieval.
- Tool approval flow.
- Basic dashboard.
- LLM provider switching.

## 4. Explicitly Out of Scope

Do not allow these to expand the sprint:

- Voice assistant.
- Mobile application.
- Smart-home integrations.
- Chase/Amex financial integration.
- Banking aggregation.
- Camera/vision.
- Multi-agent architecture.
- Kafka.
- Kubernetes.
- Distributed microservices.
- Knowledge graph.
- Complex RAG.
- Autonomous background agents.
- Production SaaS deployment.

These are future roadmap items.

## 5. Weekend Plan

### Weekend 1 — Foundation

**Goal:** Make Kanna a clean, runnable Go service.

Tasks:

- [ ] Initialize Go module.
- [ ] Configure Gin.
- [ ] Implement `/health`.
- [ ] Create package structure.
- [ ] Add configuration loading.
- [ ] Add structured logging.
- [ ] Add graceful HTTP shutdown.
- [ ] Add basic tests.
- [ ] Add `.gitignore`.
- [ ] Add `.gitattributes`.
- [ ] Create Dockerfile.
- [ ] Create initial Docker Compose setup.
- [ ] Update README with local setup.

Expected result:

```bash
go run ./cmd/server
```

and:

```bash
curl localhost:8080/health
```

returns a successful health response.

---

### Weekend 2 — AI Core

**Goal:** Kanna can understand a request and generate a response.

Tasks:

- [ ] Define `LLM` interface.
- [ ] Implement Ollama provider.
- [ ] Define `LLMRequest`.
- [ ] Define `LLMResponse`.
- [ ] Implement agent runtime.
- [ ] Implement conversation service.
- [ ] Add `POST /api/v1/chat`.
- [ ] Persist messages.
- [ ] Handle LLM errors/timeouts.
- [ ] Add basic prompt construction.

Expected result:

```text
User:
Hello Kanna

Kanna:
Hello! How can I help?
```

and:

```text
POST /api/v1/chat
```

successfully invokes the local model.

---

### Weekend 3 — Tools

**Goal:** Kanna can perform actions instead of only generating text.

Tasks:

- [ ] Define `Tool` interface.
- [ ] Build tool registry.
- [ ] Implement calculator tool.
- [ ] Implement date/time tool.
- [ ] Implement weather tool if time permits.
- [ ] Implement GitHub integration.
- [ ] Add tool invocation to agent runtime.
- [ ] Validate tool inputs.
- [ ] Persist tool execution metadata.
- [ ] Add tool execution logs.

Expected demo:

```text
User:
What are my open GitHub PRs?

Kanna:
I found 3 open PRs...
```

The response must actually come from the GitHub tool rather than fabricated LLM output.

---

### Weekend 4 — Memory + Demo

**Goal:** Make Kanna feel like an assistant rather than an API wrapper.

Tasks:

- [ ] Add PostgreSQL.
- [ ] Define memory model.
- [ ] Implement memory storage.
- [ ] Implement memory retrieval.
- [ ] Allow explicit memory creation.
- [ ] Add memory to agent context.
- [ ] Add conversation history.
- [ ] Build a minimal UI if time permits.
- [ ] Add end-to-end tests.
- [ ] Improve README.
- [ ] Record/demo the final workflow.

Expected demo:

```text
User:
Remember that my favorite coffee is cold brew.

Kanna:
Got it. I'll remember that.
```

Later:

```text
User:
What's my favorite coffee?

Kanna:
Your favorite coffee is cold brew.
```

## 6. Final Demo

The ideal August 31 demo should include:

### Demo 1 — Basic reasoning

```text
You:
What can you help me with?
```

### Demo 2 — Tool use

```text
You:
What's the weather today?
```

Kanna invokes the weather tool.

### Demo 3 — External integration

```text
You:
Show me my open GitHub PRs.
```

Kanna invokes GitHub.

### Demo 4 — Memory

```text
You:
Remember that I prefer cold brew.
```

Then later:

```text
You:
What coffee do I prefer?
```

### Demo 5 — Agentic workflow

```text
You:
Check my GitHub PRs and tell me which ones need my attention.
```

Expected flow:

```text
Natural language request
        ↓
Agent reasoning
        ↓
GitHub tool
        ↓
PR data
        ↓
LLM analysis
        ↓
Prioritized response
```

This should be the centerpiece of the sprint demo.

## 7. Definition of Done

### Backend

- [ ] Application starts reliably.
- [ ] `/health` works.
- [ ] `/api/v1/chat` works.
- [ ] Configuration is externalized.
- [ ] Logs are structured.
- [ ] Shutdown is graceful.
- [ ] Core packages have tests.

### AI

- [ ] Ollama works locally.
- [ ] LLM provider is abstracted.
- [ ] Agent runtime works.
- [ ] Tool calling works.
- [ ] Tool results are fed back into the model.

### Persistence

- [ ] PostgreSQL runs locally.
- [ ] Conversations persist.
- [ ] Messages persist.
- [ ] Memories persist.
- [ ] Memories can be retrieved.

### Tooling

- [ ] Tool registry works.
- [ ] Calculator works.
- [ ] Date/time works.
- [ ] At least one external API integration works.

### Documentation

- [ ] `architecture.md` is current.
- [ ] README contains setup instructions.
- [ ] Environment variables are documented.
- [ ] Example API requests are documented.

### Demo

- [ ] Full system runs locally.
- [ ] Natural-language request works.
- [ ] Tool invocation works.
- [ ] Memory works.
- [ ] At least one external integration works.

## 8. Risk Management

### Risk: Too much infrastructure

**Mitigation:** Keep the system as a modular monolith. No Kubernetes, Kafka, or microservices during this sprint.

### Risk: LLM/tool-calling complexity

**Mitigation:** Start with a single LLM provider and a small number of deterministic tools.

### Risk: Integration authentication takes too long

**Mitigation:** Use a low-risk API such as GitHub for the first real integration. Defer banking integrations.

### Risk: Memory becomes a research project

**Mitigation:** Start with PostgreSQL and explicit memories. Vector search can come later.

### Risk: Weekend time is limited

**Mitigation:** P0 scope takes priority. P1/P2 work is optional.

## 9. Post-Sprint Roadmap

Potential September+ work:

```text
Phase 2
├── Tool permission system
├── Scheduling
├── Gmail
├── Calendar
├── Home Assistant
└── Better memory

Phase 3
├── Voice
├── Multimodal input
├── Financial tracking
├── Proactive notifications
└── Background workflows

Phase 4
├── Advanced planning
├── Multi-agent workflows
├── Knowledge graph
├── Distributed execution
└── Kubernetes deployment
```

## 10. Sprint Principle

> Build the smallest system that proves the architecture.

Do not optimize for the number of integrations completed.

Optimize for having one clean end-to-end path:

```text
Request → Reason → Tool → Result → Memory → Response
```

If that works reliably by August 31, the sprint is complete.
