# Claude Agent Configuration

> **Start here → then load [AGENTS.md](AGENTS.md) for full routing.**

## Quick Context

- **Project**: Huobao Drama — AI Short Drama Production Platform
- **Stack**: Go 1.23 (Gin/GORM/SQLite) + Vue 3 (TypeScript/Vite/Element Plus/TailwindCSS/Pinia)
- **Architecture**: DDD 4-layer (api → application → domain → infrastructure)
- **Full routing & context loading**: See [AGENTS.md](AGENTS.md)

## Claude-Specific Notes

1. **Memory**: This repo has a compound learning system. Check `docs/solutions/` before solving problems.
2. **Plans**: Check `plans/` before starting significant work — there are 19+ active implementation plans.
3. **Chinese comments**: Backend code has Chinese comments (中文注释). Preserve them. Don't translate or remove.
4. **Large files**: `storyboard_service.go` is 77KB. Use targeted reads, not full-file loads.
5. **Prompt templates**: Located in `application/prompts/*.txt`. These are AI system prompts, not code.
6. **Existing workflows**: 28 workflow commands available via `.agent/workflows/`. Use `/plan` → `/work` flow.

## Retrieval Rules

- For routing logic: Load [AGENTS.md](AGENTS.md) → follow Task-Based Routing section
- For deep context: Load from `ai/` directory based on task type
- For prior solutions: Search `docs/solutions/` first
- For architecture: Load `ai/memory/architecture.md`
