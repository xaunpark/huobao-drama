# Debugging Router — Bug Investigation Context

> Load this when investigating errors, unexpected behavior, or performance issues.

## Retrieval Order

1. **This file** (you're reading it)
2. `ai/memory/risks.md` — known fragile areas & tech debt
3. `ai/skills/debugging.md` — debugging strategies
4. Relevant `ai/systems/*.md` for the affected subsystem

## Debugging Strategy

### Step 1: Identify the Layer

| Symptom | Likely Layer | Start Here |
|---------|-------------|------------|
| HTTP 4xx/5xx | API handler | `api/handlers/` |
| Logic error, wrong data | Service | `application/services/` |
| Data missing/corrupt | Domain/DB | `domain/models/` + `migrations/init.sql` |
| External API failure | Infrastructure | `pkg/ai/`, `pkg/image/`, `pkg/video/` |
| FFmpeg error | Infrastructure | `infrastructure/external/ffmpeg/ffmpeg.go` |
| UI not updating | Frontend | `web/src/views/` or `web/src/stores/` |
| Build failure | Config | `go.mod`, `web/package.json`, `Dockerfile` |

### Step 2: Check Known Issues

Search these locations for prior solutions:
1. `docs/solutions/` — documented solutions with YAML frontmatter
2. `plans/` — implementation plans may document edge cases
3. `ai/memory/risks.md` — known fragile areas

### Step 3: Trace the Data Flow

```
HTTP Request
  → api/routes/routes.go (which handler?)
    → api/handlers/xxx.go (input parsing)
      → application/services/xxx_service.go (business logic)
        → domain/models/xxx.go (data model)
        → pkg/ai/ or pkg/image/ or pkg/video/ (external API)
        → infrastructure/database/ (DB operations)
      ← Response transformation
    ← JSON response
  ← HTTP Response
```

## Common Bug Patterns

### 1. SQLite "database is locked"
- **Cause**: Concurrent write operations
- **Fix**: Ensure WAL mode is enabled, reduce write concurrency
- **Files**: `infrastructure/database/database.go`

### 2. Expired Image/Video URLs
- **Cause**: Provider URLs have TTL, resource transfer didn't complete
- **Fix**: Check `infrastructure/scheduler/resource_transfer_scheduler.go`
- **Files**: `application/services/resource_transfer_service.go`

### 3. AI Response Parse Failures
- **Cause**: AI model returns non-JSON or malformed output
- **Fix**: Check prompt templates in `application/prompts/`, improve parsing
- **Files**: `application/services/storyboard_service.go` (JSON parsing logic)

### 4. FFmpeg Command Failures
- **Cause**: Special characters in file paths, missing codecs, wrong args
- **Fix**: Check FFmpeg command construction
- **Files**: `infrastructure/external/ffmpeg/ffmpeg.go`

### 5. Video Generation Stuck
- **Cause**: Async polling lost track of task, provider API down
- **Fix**: Check video generation status polling logic
- **Files**: `application/services/video_generation_service.go`

### 6. Storyboard Mode Conflicts
- **Cause**: 6 different modes share code paths with conditional logic
- **Fix**: Trace the specific mode path through `storyboard_service.go`
- **Files**: `application/services/storyboard_service.go` (77KB — use targeted search)

## Logging

- Backend uses Zap logger (`pkg/logger/`)
- Many `fmt.Printf` debug statements in AI clients (not structured)
- Frontend: Browser console
- Docker: `docker-compose logs -f`

## No Test Suite

This codebase has minimal tests (`storyboard_parser_test.go` is the only test file).
Verification must be done through:
1. Running the application
2. Checking logs
3. Manual API testing
4. Visual inspection of frontend
