# System Map — Subsystems → Folders

> Maps conceptual subsystems to their folder locations.

## Core Subsystems

| System | Primary Location | Secondary |
|--------|-----------------|-----------|
| **HTTP API** | `api/handlers/` (23 files) | `api/routes/routes.go`, `api/middlewares/` |
| **Business Logic** | `application/services/` (29 files) | — |
| **Domain Models** | `domain/models/` (12 files) | — |
| **AI Text** | `pkg/ai/` (4 files) | `application/services/ai_service.go` |
| **AI Image** | `pkg/image/` (5 files) | `application/services/image_generation_service.go` |
| **AI Video** | `pkg/video/` (6 files) | `application/services/video_generation_service.go` |
| **Prompt Templates** | `application/prompts/` (34 files) | `application/services/prompt_i18n.go` |
| **Database** | `infrastructure/database/` (2 files) | `domain/models/*.go` |
| **FFmpeg** | `infrastructure/external/ffmpeg/` (1 file) | `application/services/video_merge_service.go` |
| **Storage** | `infrastructure/storage/` (1 file) | — |
| **Scheduler** | `infrastructure/scheduler/` (1 file) | — |
| **Config** | `pkg/config/`, `configs/` | — |
| **Logging** | `pkg/logger/` | — |
| **Frontend SPA** | `web/src/` | — |
| **Channel Templates** | `docs/*_template.md` | `docs/features/` |

## Data Flow: Script → Video

```
1. Script Input
   └─ web/src/views/script/ → api/handlers/drama.go → services/drama_service.go

2. Character Extraction
   └─ handlers/character_library.go → services/character_library_service.go
   └─ prompts/character_extraction.txt → pkg/ai/ (AI call)

3. Scene Extraction
   └─ handlers/scene.go → services/
   └─ prompts/scene_extraction.txt → pkg/ai/

4. Storyboard Generation
   └─ handlers/storyboard.go → services/storyboard_service.go (77KB)
   └─ prompts/storyboard_*.txt (mode-specific) → pkg/ai/

5. Image Generation
   └─ handlers/image_generation.go → services/image_generation_service.go
   └─ prompts/image_*.txt → pkg/image/ (provider-specific)

6. Video Generation
   └─ handlers/video_generation.go → services/video_generation_service.go
   └─ prompts/video_*.txt → pkg/video/ (provider-specific, async)

7. Video Merge
   └─ handlers/video_merge.go → services/video_merge_service.go
   └─ infrastructure/external/ffmpeg/ffmpeg.go

8. Episode Finalization
   └─ handlers/drama.go:FinalizeEpisode → FFmpeg final compose
```

## Frontend Subsystems

| Subsystem | Location | Purpose |
|-----------|----------|---------|
| API Clients | `web/src/api/` (14 files) | Axios HTTP calls |
| State | `web/src/stores/` | Pinia state management |
| Types | `web/src/types/` (10 files) | TypeScript interfaces |
| Routing | `web/src/router/` | Vue Router page routes |
| i18n | `web/src/locales/` | Translation files |
| Composables | `web/src/composables/` | Vue composition hooks |

## Agent Tooling Subsystems

| System | Location | Purpose |
|--------|----------|---------|
| Compound Learning | `docs/solutions/`, `docs/explorations/` | Knowledge base |
| Implementation Plans | `plans/` | Planned & in-progress work |
| Deferred Work | `todos/` | Action items |
| Workflows | `.agent/workflows/` (28 files) | Agent workflow commands |
| Skills | `skills/` (7 directories) | Agent capabilities |
| Scripts | `scripts/` (49 files) | Automation & validation |
| AI Orchestration | `ai/` | Agent context system (this system) |
