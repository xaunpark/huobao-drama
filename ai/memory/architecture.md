# Architecture — Architectural Philosophy & Patterns

> Inferred from real code analysis. Not generic descriptions.

## Architectural Philosophy

This project follows **pragmatic DDD** — Domain-Driven Design structure without strict DDD ceremony:

1. **Clear layer separation** — API → Application → Domain → Infrastructure
2. **No repository pattern** — Services use GORM `*gorm.DB` directly (pragmatic choice)
3. **Manual DI** — All dependency injection happens in `api/routes/routes.go:SetupRouter()`, no DI framework
4. **Single route file** — All routes registered in one function, not distributed
5. **Monolith design** — Single binary serves API + static frontend, no microservices

## Key Design Decisions

### 1. SQLite Over PostgreSQL/MySQL
- **Why**: Zero-config, single-file DB, easy backup, sufficient for single-instance
- **Tradeoff**: No concurrent write scaling, file-level locking
- **Mitigation**: WAL mode enabled, pure Go driver (no CGO)

### 2. Pure Go SQLite (`modernc.org/sqlite`)
- **Why**: Enables `CGO_ENABLED=0` cross-compilation
- **Tradeoff**: Slightly slower than C SQLite
- **Benefit**: Simpler Docker builds, Alpine compatibility

### 3. FFmpeg as External Process
- **Why**: Full FFmpeg capabilities without Go bindings
- **Tradeoff**: Must be installed on system, subprocess management complexity
- **Pattern**: Shell out via `exec.Command`, parse stdout/stderr

### 4. Multi-Provider AI Abstraction
- **Pattern**: Interface per concern (`AIClient`, `ImageClient`, `VideoClient`)
- **Providers configurable at runtime** via DB-stored configs
- **User manages API keys** in web UI, not config files
- **Factory pattern** in `pkg/image/` and `pkg/video/` selects provider by config

### 5. Prompt Template Architecture
- **Static templates** in `application/prompts/*.txt` (compiled into binary)
- **User overrides** via DB-stored `prompt_template` records
- **I18n layer** (`prompt_i18n.go`) wraps both
- **Channel templates** in `docs/*.md` — reference docs, not runtime

### 6. Async Video Generation
- **All video providers** return task IDs, not immediate results
- **Polling model**: Service polls provider API until completion
- **Resource transfer**: Background cron downloads remote URLs to local storage

### 7. Frontend-Backend Coupling
- **SPA embedded in Go binary** for production
- **Dev mode**: Separate servers with Vite proxy
- **TypeScript types manually mirror Go models** — no code generation
- **No API schema** (OpenAPI/Swagger) — types kept in sync manually

## Code Organization Patterns

### Handler File Pattern
```
handlers/{resource}.go
  - type XxxHandler struct { db, cfg, log, ...services }
  - func NewXxxHandler(...) *XxxHandler
  - func (h *XxxHandler) CreateXxx(c *gin.Context)
  - func (h *XxxHandler) GetXxx(c *gin.Context)
  - func (h *XxxHandler) UpdateXxx(c *gin.Context)
  - func (h *XxxHandler) DeleteXxx(c *gin.Context)
  - func (h *XxxHandler) SpecialAction(c *gin.Context)
```

### Service File Pattern
```
services/{resource}_service.go
  - type XxxService struct { db, log, cfg, ...deps }
  - func NewXxxService(...) *XxxService
  - func (s *XxxService) Create/Get/Update/Delete(...)
  - func (s *XxxService) BusinessLogic(...)
  - Private helper functions at bottom
```

### Model File Pattern
```
models/{resource}.go
  - GORM struct with tags: gorm, json
  - TableName() method
  - Relationships via GORM associations
  - Runtime fields with gorm:"-" tag
```

## Scaling Assumptions

- **Single instance** — no horizontal scaling designed
- **Moderate data** — thousands of dramas, not millions
- **File-based storage** — local filesystem, no S3/cloud storage
- **Sequential batch processing** — no job queue or worker pool
- **Single-tenant** — no multi-user authentication system
