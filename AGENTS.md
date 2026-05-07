# AGENTS.md — AI Context Orchestration Hub

> **This file is the entrypoint for ALL AI agents.**
> It is NOT documentation. It is a **bootloader**, **router**, and **retrieval map**.
> Load this file first. Then follow the routing rules to load only what you need.

---

## Project Identity

| Key | Value |
|-----|-------|
| **Name** | Huobao Drama — AI Short Drama Production Platform |
| **Module** | `github.com/drama-generator/backend` |
| **Stack** | Go 1.23 (Gin + GORM + SQLite) · Vue 3 + TypeScript (Vite + Element Plus + TailwindCSS + Pinia) |
| **Architecture** | DDD (Domain-Driven Design) — 4-layer |
| **AI Integration** | OpenAI, Gemini, Doubao, MiniMax, Sora, FlowTool, Volcengine |
| **Port** | Backend `:5678` · Frontend dev `:3012` |
| **DB** | SQLite (WAL mode, `./data/drama_generator.db`) |

## What This System Does

Automates the full AI short drama production pipeline:

```
Script → Character Extraction → Scene Design → Storyboard Breakdown
→ Image Generation (text-to-image) → Video Generation (image-to-video)
→ Video Merge (FFmpeg) → Final Episode
```

Multiple production modes: Standard, Rapid Cut, Nursery Rhyme, MV Maker, Narrative MV, Voiceover Director.

---

## Architecture Summary (4-Layer DDD)

```
api/                    → HTTP handlers + routes + middleware (Gin)
  handlers/             → 23 handler files (one per resource)
  routes/routes.go      → Single route registry
  middlewares/           → CORS, rate limiting, logging

application/            → Business logic layer
  services/             → 29 service files (core business logic)
  prompts/              → 34 prompt template files (.txt)

domain/                 → Pure domain models (GORM structs)
  models/               → 12 model files (Drama, Episode, Storyboard, Character, Scene, Prop, etc.)

infrastructure/         → External concerns
  database/             → SQLite connection + custom logger
  external/ffmpeg/      → FFmpeg video processing (29KB, complex)
  scheduler/            → Resource transfer cron
  storage/              → Local file storage

pkg/                    → Shared packages
  ai/                   → AI client interface + OpenAI/Gemini implementations
  config/               → Viper-based config loading
  image/                → Image client factory (OpenAI, Gemini, Volcengine, FlowTool)
  video/                → Video client factory (Doubao, MiniMax, Sora, FlowTool, Chatfire)
  logger/               → Zap logger wrapper
  response/             → HTTP response helpers
  utils/                → Utility functions

web/                    → Vue 3 SPA frontend
  src/api/              → 14 API client modules
  src/views/            → 8 view groups (dashboard, drama, editor, generation, script, settings, storyboard, workflow)
  src/stores/           → Pinia store (episode.ts)
  src/components/       → Shared UI components
  src/types/            → 10 TypeScript type definition files
```

---

## Context Loading Strategy

### Level 1 — Always Loaded (Lightweight, <500 tokens each)

| File | Purpose |
|------|---------|
| `AGENTS.md` | This file. Routing brain. |
| `CLAUDE.md` | Claude-specific vendor adapter |
| `GEMINI.md` | Gemini-specific vendor adapter |

### Level 2 — Task Routers (Load based on task type)

| Router | When to Load |
|--------|-------------|
| `ai/routing/backend-router.md` | Modifying Go backend code |
| `ai/routing/frontend-router.md` | Modifying Vue frontend code |
| `ai/routing/ai-pipeline-router.md` | Working with AI generation pipeline |
| `ai/routing/infrastructure-router.md` | Database, Docker, FFmpeg, deployment |
| `ai/routing/debugging-router.md` | Investigating bugs or errors |
| `ai/routing/prompt-router.md` | Modifying AI prompt templates |

### Level 3 — Deep Context (Load only when specifically needed)

| Category | Path | When |
|----------|------|------|
| Skills | `ai/skills/*.md` | Need domain-specific capabilities |
| Memory | `ai/memory/*.md` | Need architectural/historical context |
| Systems | `ai/systems/*.md` | Working on specific subsystem |
| Indexes | `ai/indexes/*.md` | Need to discover files/features |
| Rules | `ai/rules/*.md` | Need constraint/convention guidance |
| Workflows | `ai/workflows/*.md` | Executing multi-step procedures |

---

## Task-Based Routing

### IF debugging a bug:
```
1. Load ai/routing/debugging-router.md
2. Load ai/memory/risks.md (known fragile areas)
3. Load relevant ai/systems/*.md for the affected subsystem
```

### IF implementing backend feature:
```
1. Load ai/routing/backend-router.md
2. Load ai/rules/architecture-rules.md
3. Load ai/memory/conventions.md
4. Load ai/indexes/feature-map.md (check for existing patterns)
```

### IF implementing frontend feature:
```
1. Load ai/routing/frontend-router.md
2. Load ai/rules/naming.md
3. Load ai/memory/coding-style.md
```

### IF modifying AI prompts/templates:
```
1. Load ai/routing/prompt-router.md
2. Load ai/memory/conventions.md (prompt conventions)
3. Load docs/ template examples for context
```

### IF working on video/image generation pipeline:
```
1. Load ai/routing/ai-pipeline-router.md
2. Load ai/systems/video-pipeline.md
3. Load ai/systems/image-pipeline.md
4. Load ai/memory/risks.md (known API gotchas)
```

### IF modifying Docker/deployment/infrastructure:
```
1. Load ai/routing/infrastructure-router.md
2. Load ai/memory/infrastructure.md
3. Load ai/memory/decisions.md
```

### IF adding a new storyboard mode:
```
1. Load ai/routing/prompt-router.md
2. Load ai/systems/storyboard-system.md
3. Load plans/ for existing mode implementation plans
4. Load application/prompts/ for prompt template patterns
```

---

## Retrieval Priorities

1. **Check `docs/solutions/`** before solving any problem (prior knowledge)
2. **Check `plans/`** before starting significant work (existing plans)
3. **Check `todos/`** for deferred work items
4. **Check `ai/memory/decisions.md`** before making architectural choices
5. **Follow `ai/rules/`** constraints always

---

## Critical Files (High-Impact, Touch Carefully)

| File | Risk | Why |
|------|------|-----|
| `api/routes/routes.go` | HIGH | Single route registry, all endpoints |
| `application/services/storyboard_service.go` | CRITICAL | 77KB, core business logic, 6 production modes |
| `application/services/image_generation_service.go` | HIGH | 43KB, all image generation logic |
| `application/services/video_generation_service.go` | HIGH | 39KB, all video generation logic |
| `infrastructure/external/ffmpeg/ffmpeg.go` | HIGH | 29KB, video processing, FFmpeg commands |
| `domain/models/drama.go` | HIGH | Core domain model, 70+ fields on Storyboard alone |
| `application/services/prompt_i18n.go` | HIGH | 27KB, all prompt internationalization |
| `migrations/init.sql` | CRITICAL | 20KB, database schema source of truth |

---

## System Map

```
ai/                         → AI agent orchestration system
  indexes/                  → Discoverability maps
    feature-map.md          → Features → files mapping
    system-map.md           → Subsystems → folders mapping
    dependency-map.md       → External dependencies
    pipeline-map.md         → Data flow pipelines
  memory/                   → Long-term project intelligence
    architecture.md         → Architectural philosophy & patterns
    decisions.md            → Key technical decisions made
    coding-style.md         → Code style conventions
    infrastructure.md       → Infra patterns & config
    conventions.md          → Naming & structural conventions
    risks.md                → Known fragile areas & tech debt
  skills/                   → Reusable agent capabilities
    debugging.md            → Debugging strategies for this codebase
    refactor.md             → Safe refactoring patterns
    testing.md              → Testing approach
    prompt-engineering.md   → Writing AI prompts for this system
    new-storyboard-mode.md  → Adding new production modes
  workflows/                → Multi-step execution procedures
    feature-workflow.md     → Implementing new features
    bugfix-workflow.md      → Fixing bugs
    release-workflow.md     → Release process
    prompt-workflow.md      → Creating/modifying prompt templates
  rules/                    → Constraints & conventions
    forbidden-patterns.md   → Anti-patterns to avoid
    naming.md               → Naming conventions
    architecture-rules.md   → Architecture constraints
    api-rules.md            → API design rules
  routing/                  → Task-specific context routers
    backend-router.md       → Go backend context
    frontend-router.md      → Vue frontend context
    ai-pipeline-router.md   → AI generation pipeline context
    infrastructure-router.md → Infra/deploy context
    debugging-router.md     → Bug investigation context
    prompt-router.md        → Prompt template context
  systems/                  → Subsystem deep documentation
    storyboard-system.md    → Storyboard generation engine
    video-pipeline.md       → Video generation & merge
    image-pipeline.md       → Image generation pipeline
    frontend-system.md      → Vue SPA architecture
    backend-system.md       → Go API server architecture
  glossary/                 → Domain terminology
    terms.md                → Key terms & abbreviations
  prompts/                  → Agent interaction prompts
    onboarding.md           → New agent quick start
```

---

## Existing Knowledge System

This repo already has a compound learning system:

| Path | Purpose |
|------|---------|
| `docs/solutions/` | Solved problems with YAML frontmatter |
| `plans/` | Implementation plans (19 active, many archived) |
| `todos/` | Deferred work items |
| `skills/` | Agent skills (code-review, compound-docs, debug, testing, session-resume) |
| `.agent/workflows/` | 28 workflow commands (`/plan`, `/work`, `/review`, `/compound`, etc.) |
| `scripts/` | 49 automation scripts |

**Always check these BEFORE creating new solutions.**

---

## AI Workflow Expectations

1. **Read AGENTS.md first** (you're doing this now)
2. **Route to the right context** using Task-Based Routing above
3. **Search before solving** — check `docs/solutions/`, `plans/`, `ai/memory/`
4. **Follow rules** — check `ai/rules/` for constraints
5. **Document after solving** — update relevant `ai/memory/` or `docs/solutions/`
6. **Minimize token usage** — load only what you need, never load all context at once
