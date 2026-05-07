# Onboarding — Quick Start for New AI Agents

> Read this if you're a new agent encountering this codebase for the first time.

## 30-Second Summary

**Huobao Drama** is an AI-powered short drama production platform.
- **Backend**: Go 1.23, Gin web framework, GORM ORM, SQLite database
- **Frontend**: Vue 3, TypeScript, Vite, Element Plus, TailwindCSS, Pinia
- **Purpose**: Automate Script → Characters → Storyboards → Images → Videos → Final Episode
- **AI Providers**: OpenAI, Gemini, Doubao, MiniMax, Sora, and more

## First Steps

1. **Read `AGENTS.md`** — the routing brain (you're already past this if you're here)
2. **Identify your task type** — use Task-Based Routing in AGENTS.md
3. **Load the relevant router** from `ai/routing/`
4. **Check prior work** — search `docs/solutions/` and `plans/`

## Key Facts

- The largest file is `storyboard_service.go` (77KB) — never read it entirely
- Chinese comments throughout backend — preserve them
- No automated tests — manual testing only
- AI provider configs stored in database, not config files
- Video generation is async — all providers use task polling
- SQLite is the only database — no PostgreSQL/MySQL in production

## File Discovery

- `ai/indexes/feature-map.md` — features → files
- `ai/indexes/system-map.md` — subsystems → folders
- `ai/indexes/dependency-map.md` — external dependencies
- `ai/indexes/pipeline-map.md` — data flow pipelines

## Constraints

- `ai/rules/forbidden-patterns.md` — things you must NOT do
- `ai/rules/architecture-rules.md` — layer dependency rules
- `ai/rules/naming.md` — naming conventions
- `ai/rules/api-rules.md` — API endpoint design rules

## Getting Help

- `ai/memory/risks.md` — known fragile areas
- `ai/memory/decisions.md` — why things were done this way
- `ai/skills/debugging.md` — debugging strategies
- `docs/solutions/` — previously solved problems
