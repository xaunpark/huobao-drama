# System: Backend — Go API Server Architecture

> Deep documentation of the Go backend subsystem.

## Entry Point

`main.go` — Initializes config, logger, database, storage, router, then starts HTTP server.

## Layer Details

### API Layer (`api/`)

#### Handlers (23 files)
Each handler: struct with dependencies, constructor, HTTP methods.

Key handlers by size/complexity:
| Handler | Size | Domain |
|---------|------|--------|
| `storyboard.go` | 5KB | Storyboard CRUD + generation |
| `image_generation.go` | 7KB | Image generation + batch |
| `video_generation.go` | 5KB | Video generation + batch |
| `drama.go` | 9KB | Drama CRUD + episode management |
| `character_library.go` | 8KB | Character library + image gen |
| `scene.go` | 5KB | Scene CRUD + image gen |

#### Routes (`routes.go`, single file)
- All routes registered in `SetupRouter()`
- All services instantiated here (manual DI)
- CORS, rate limiting, logging middleware applied

#### Middlewares
- `LoggerMiddleware` — request logging with Zap
- `CORSMiddleware` — configurable CORS origins
- `RateLimitMiddleware` — rate limiting

### Application Layer (`application/`)

#### Services (29 files, core business logic)
Key services by size:
| Service | Size | Complexity |
|---------|------|-----------|
| `storyboard_service.go` | 77KB | 6 production modes, central dispatcher |
| `image_generation_service.go` | 43KB | Multi-provider, batch, reference images |
| `video_generation_service.go` | 39KB | Multi-provider, async, polling |
| `storyboard_narrative_service.go` | 29KB | Narrative MV mode |
| `prompt_i18n.go` | 28KB | Bilingual prompt resolution |
| `frame_prompt_service.go` | 26KB | Frame prompt generation |
| `storyboard_composition_service.go` | 25KB | Voiceover director mode |
| `storyboard_nursery_service.go` | 24KB | Nursery rhyme mode |
| `character_library_service.go` | 23KB | Character management |
| `drama_service.go` | 22KB | Drama CRUD + episode finalization |
| `style_distill_service.go` | 22KB | Per-shot style generation |
| `video_merge_service.go` | 22KB | FFmpeg video merging |

#### Prompts (34 template files)
Plain text `.txt` files loaded at runtime. See `ai/routing/prompt-router.md`.

### Domain Layer (`domain/models/`)

12 model files. Primary entity: `Drama` with nested associations.

```
Drama
  ├── Episode (1:N)
  │     ├── Storyboard (1:N, 50+ fields)
  │     │     ├── Character (M:N via junction)
  │     │     └── Prop (M:N via junction)
  │     └── Scene (1:N)
  ├── Character (1:N, also M:N with Episodes)
  ├── Scene (1:N)
  └── Prop (1:N)
```

Other models: AIConfig, FramePrompt, ImageGeneration, VideoGeneration, VideoMerge, VideoReview, PromptTemplate, Asset, Task, Timeline.

### Infrastructure Layer

- `database/` — SQLite setup + GORM custom logger
- `external/ffmpeg/` — FFmpeg subprocess wrapper (29KB)
- `scheduler/` — Resource transfer cron job
- `storage/` — Local file storage

### Package Layer (`pkg/`)

- `ai/` — AI text client interface + OpenAI/Gemini implementations
- `image/` — Image generation client factory (4 providers)
- `video/` — Video generation client factory (5 providers)
- `config/` — Viper config loading
- `logger/` — Zap logger wrapper
- `response/` — HTTP response helpers
- `utils/` — Utility functions

## Startup Sequence

```
1. config.LoadConfig()           → Read configs/config.yaml
2. logger.NewLogger()            → Initialize Zap
3. database.NewDatabase()        → Connect SQLite (WAL mode)
4. database.AutoMigrate()        → GORM auto-migration
5. storage.NewLocalStorage()     → Initialize file storage
6. routes.SetupRouter()          → Wire all dependencies, register routes
7. http.Server.ListenAndServe()  → Start on configured port
8. Signal handler                → Graceful shutdown (close DB, HTTP)
```
