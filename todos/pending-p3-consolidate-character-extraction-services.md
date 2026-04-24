---
id: pending-p3-consolidate-character-extraction-services
status: pending
priority: p3
title: Consolidate duplicate character extraction services
created: 2026-04-24
dependencies: []
---

## Problem

There are **two separate services** that both handle character extraction from episode scripts. They diverged over time and became out of sync, causing bugs when new fields were added to one but not the other.

### Service 1 — `ScriptGenerationService` (ACTIVE)
- **File:** `application/services/script_generation_service.go`
- **Endpoint:** `POST /api/v1/generation/characters`
- **Called from:** `generationAPI.generateCharacters()` in `EpisodeWorkflow.vue:2108`
- **Status:** This is the code path currently used by the "Extract Characters & Scenes" button.

### Service 2 — `CharacterLibraryService.ExtractCharactersFromScript` (UNUSED from workflow)
- **File:** `application/services/character_library_service.go`
- **Endpoint:** `POST /api/v1/episodes/:id/characters/extract`
- **Called from:** `characterLibraryAPI.extractFromEpisode()` — defined in `character-library.ts` but **not called from EpisodeWorkflow.vue**
- **Status:** Dead code path from the main extraction workflow. The endpoint exists and works, but the frontend doesn't route to it.

## Root Cause of the Bug We Fixed

When `character_prompt`, `variant_prompt`, `episode_descriptor` were added to the `Character` model, they were added to Service 2's `extractedChar` struct and save logic — but **not to Service 1**. Since Service 1 is what actually runs, the new fields were silently ignored on every extraction.

**Fix applied (2026-04-24):** Updated Service 1's `extractedChar` struct and save logic to include the 3 new fields. Both services are now in sync.

## Recommended Action

Migrate `EpisodeWorkflow.vue` to call Service 2's endpoint instead of Service 1:

1. **Update `EpisodeWorkflow.vue`** — Replace `generationAPI.generateCharacters(...)` with `characterLibraryAPI.extractFromEpisode(episodeId)`. Service 2 is more feature-rich (base-name matching, narrator profile, proper variant update logic).

2. **Verify parity** — Ensure Service 2 supports the `model` selection parameter that Service 1 currently accepts.

3. **Deprecate Service 1's character generation** — Mark `ScriptGenerationService.GenerateCharacters` and `POST /api/v1/generation/characters` as deprecated (or remove them) once migration is confirmed working.

4. **Clean up dead code** — Remove `extractFromEpisode` from `character-library.ts` API client if it remains unused, or remove `generateCharacters` from `generation.ts` after migration.

## Risk

- Low. Both services produce the same end result. The migration is a frontend-only change in which endpoint URL is called.
- Service 2 already handles the same data model and task polling pattern.

## Related

- Bug discovered while implementing character variant UI (session 2026-04-24)
- `review-character-variants.md` — P3 note about `charDescMap` duplication (separate issue)
