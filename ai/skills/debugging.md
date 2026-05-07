# Skill: Debugging — Strategies for This Codebase

> Reference: `ai/routing/debugging-router.md` for routing context first.

## Prerequisites
- Load `ai/memory/risks.md` — known fragile areas
- Know the layer (api/services/domain/infra) before deep-diving

## Strategy 1: Layer Identification

Most bugs fall into one of these patterns:

| Symptom | Layer | Entry Point |
|---------|-------|-------------|
| 4xx HTTP errors | Handler | `api/handlers/` — check input validation |
| Wrong data returned | Service | `application/services/` — check business logic |
| Missing/null data | Model/DB | `domain/models/` — check GORM associations, nullable fields |
| External API failure | Infra | `pkg/ai/`, `pkg/image/`, `pkg/video/` — check provider response |
| FFmpeg crash | Infra | `infrastructure/external/ffmpeg/ffmpeg.go` — check command args |
| UI not rendering | Frontend | `web/src/views/` — check API response handling |

## Strategy 2: Log Tracing

### Backend Logs
- Zap structured logging: `logr.Info()`, `logr.Error()`, `logr.Warnw()`
- AI client debug prints: `fmt.Printf("OpenAI: ...")` in `pkg/ai/openai_client.go`
- Run `go run main.go` with `debug: true` in config for verbose output

### Frontend Logs
- Browser DevTools → Console
- Network tab: check API request/response payloads
- Vue DevTools: inspect component state and Pinia stores

## Strategy 3: AI Response Debugging

Common AI response issues:
1. **Markdown-wrapped JSON**: AI returns ````json\n{...}\n```` — strip markdown fences
2. **Truncated response**: `max_tokens` too low — check response `finish_reason`
3. **Content filtered**: `finish_reason: "content_filter"` — modify prompt to avoid triggers
4. **Wrong format**: AI ignores JSON format instructions — improve prompt with examples
5. **Provider difference**: Same prompt works on GPT-4 but fails on Gemini — test both

## Strategy 4: Database Debugging

```sql
-- Check via SQLite CLI
sqlite3 data/drama_generator.db

-- Common queries:
.tables                             -- List all tables
.schema storyboards                 -- Show table structure
SELECT COUNT(*) FROM storyboards;   -- Count records
SELECT * FROM storyboards WHERE episode_id = ? LIMIT 5;
```

Watch for:
- Soft-deleted records (WHERE deleted_at IS NULL)
- JSON fields that contain serialized data
- Foreign key relationships (GORM handles in app layer, not DB constraints)

## Strategy 5: FFmpeg Debugging

```bash
# Test FFmpeg directly
ffmpeg -version

# Check the constructed command in logs
# FFmpeg commands are logged in video_merge_service.go

# Common failures:
# - File path with spaces/special chars
# - Missing input files
# - Unsupported codec
# - Insufficient disk space
```

## Anti-Patterns (Don't Do These)

- ❌ Don't read entire 77KB `storyboard_service.go` — search for specific function names
- ❌ Don't modify AI prompt templates without updating the parsing logic in services
- ❌ Don't assume a bug is in the frontend when the API response itself is wrong
- ❌ Don't debug production data directly — use `check_db.go` scripts for safe inspection

## Instrumentation

```bash
# Log this skill usage
./scripts/log-skill.sh "debugging" "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
```
