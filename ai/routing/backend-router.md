# Backend Router — Go API Server Context

> Load this when modifying any Go backend code.

## Retrieval Order

1. **This file** (you're reading it)
2. `ai/rules/architecture-rules.md` — DDD constraints
3. `ai/memory/conventions.md` — naming & patterns
4. `ai/memory/coding-style.md` — code style

## Architecture Quick Ref

```
api/handlers/       → HTTP request handling, input validation, response formatting
application/services/ → Business logic, AI orchestration, data transformation
domain/models/      → GORM structs, pure domain entities
infrastructure/     → Database, FFmpeg, storage, schedulers
pkg/                → Shared utilities (AI clients, config, logging, image/video factories)
```

## Key Patterns

### Handler Pattern
- One file per resource (`drama.go`, `storyboard.go`, `video_generation.go`, etc.)
- Constructor: `NewXxxHandler(db, cfg, log, ...deps)` — returns struct
- Methods: `(h *XxxHandler) VerbNoun(c *gin.Context)`
- Always call services layer — handlers do NOT contain business logic

### Service Pattern
- One file per domain (`drama_service.go`, `storyboard_service.go`, etc.)
- Constructor: `NewXxxService(db, log, ...deps)` — returns struct
- Services contain ALL business logic — AI calls, data transformation, orchestration
- Services call `pkg/ai/`, `pkg/image/`, `pkg/video/` for external APIs

### Route Registration
- ALL routes defined in `api/routes/routes.go` — single file
- Group pattern: `api.Group("/resource")` with CRUD verbs
- All services instantiated in `SetupRouter()` — manual DI, no framework

### AI Client Pattern
- `pkg/ai/client.go` — `AIClient` interface: `GenerateText()`, `GenerateImage()`, `TestConnection()`
- `pkg/ai/openai_client.go` — OpenAI-compatible implementation (also used for local models)
- `pkg/ai/gemini_client.go` — Google Gemini
- `pkg/image/image_client.go` — Image generation factory
- `pkg/video/video_client.go` — Video generation factory (5 providers)

## When You Need Deeper Context

| Task | Load |
|------|------|
| Adding new API endpoint | `api/routes/routes.go` + existing handler as template |
| Adding new service | Existing service file as template + `ai/rules/architecture-rules.md` |
| Modifying domain models | `domain/models/drama.go` + `migrations/init.sql` |
| Adding new AI provider | `pkg/ai/client.go` interface + existing client as template |
| Working with prompts | `ai/routing/prompt-router.md` |
| Storyboard modes | `ai/systems/storyboard-system.md` |

## Critical Warnings

- `storyboard_service.go` is 77KB — the largest file. Read specific functions, not the whole file.
- `prompt_i18n.go` is 27KB — handles all prompt localization. Very complex.
- Chinese comments (中文注释) throughout — preserve them, don't translate.
- SQLite uses WAL mode — concurrent reads OK, writes serialize. Mind `database is locked` errors.
- No test suite — rely on manual testing and log inspection.
