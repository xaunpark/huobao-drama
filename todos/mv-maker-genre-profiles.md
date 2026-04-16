---
priority: p3
status: ready
created: 2026-04-16
source: plans/mv-maker-mode.md
---

# MV Maker — Additional Genre Profiles

Deferred from MV Maker Phase 1 implementation. These genre profiles only need a prompt file + PromptTemplate registration + dropdown option each.

## Genre Profiles to Add

1. **gaming_parody** — LHUGUENY-style musical parodies (comedy + horror blend, narrative-heavy, character parody)
2. **general** — Mainstream pop/rock/hip-hop MVs (standard MTV-style editing)
3. **anime_opening** — Anime OP/ED style (sakuga cuts, speed lines, transformation sequences)
4. **lyric_video** — Typography-focused lyric videos (minimal scene visuals, text as primary visual)

## Steps per Genre

- [ ] Analyze 2-3 reference videos for each genre (extract ASL, energy curve, shot recipes)
- [ ] Create `application/prompts/storyboard_mv_{genre}.txt` with genre-specific rules
- [ ] Add to `PromptTemplatePrompts` struct and `PromptTypeToDefaultFile` map
- [ ] Add `<el-option>` to genre dropdown in EpisodeWorkflow.vue
- [ ] Add locale strings for en-US and zh-CN
- [ ] Test with real lyrics input

## Context

See [plans/mv-maker-mode.md](../plans/mv-maker-mode.md) for architecture details.
Zero backend logic changes needed — only prompt files + registration.
