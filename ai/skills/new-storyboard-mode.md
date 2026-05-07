# Skill: New Storyboard Mode — Adding Production Modes

> Use when implementing a new storyboard generation mode.

## Prerequisites
- Load `ai/systems/storyboard-system.md` — system architecture
- Load `ai/memory/conventions.md` — naming conventions
- Check `plans/` for existing mode plans as reference

## Implementation Steps

### 1. Create Prompt Templates
```
application/prompts/storyboard_{mode}.txt         → Main system prompt
application/prompts/storyboard_{mode}_format.txt   → Output format spec
```

Pattern: Copy nearest existing mode's prompts and adapt.

### 2. Create Mode-Specific Service (Optional)
```
application/services/storyboard_{mode}_service.go
```

Only create if the mode has significantly different logic from standard.
If the mode is just a different prompt, handle it in the dispatcher.

### 3. Add Dispatch Case
In `application/services/storyboard_service.go`:
- Add case in the mode dispatch switch
- Wire up prompt template loading
- Configure mode-specific fields on Storyboard model

### 4. Update Domain Model (If Needed)
In `domain/models/drama.go`:
- Add mode-specific fields to the `Storyboard` struct
- Use nullable pointers for optional fields
- GORM auto-migrate will add columns on restart

### 5. Update Frontend
- Add mode option to storyboard generation UI
- Handle mode-specific fields in shot display
- Update TypeScript types if new fields added

### 6. Create Implementation Plan
```
plans/{mode}-mode.md
```

### 7. Test
- Test with at least 2 different AI providers
- Test with short and long scripts
- Verify JSON parsing handles all expected fields
- Check storyboard display in frontend

## Existing Modes for Reference

| Mode | Service File | Prompt Files | Added Fields |
|------|-------------|-------------|-------------|
| Standard | `storyboard_service.go` | `storyboard_story_breakdown.txt` | Core fields |
| Visual Unit | `storyboard_service.go` | `storyboard_visual_unit*.txt` | — |
| Nursery Rhyme | `storyboard_nursery_service.go` | `storyboard_nursery_rhyme*.txt` | lyrics_text, section_type, verse_subject, etc. |
| MV Maker | `storyboard_mv_service.go` | `storyboard_mv_*.txt` | — |
| Narrative MV | `storyboard_narrative_service.go` | `storyboard_narrative_*.txt` | narrative_part, has_music, music_segment, etc. |
| Voiceover | `storyboard_composition_service.go` | Multiple | script_segment, narrator_script, audio fields |
| Rapid Cut | `rapid_cut_service.go` | `rapid_cut_merge.txt` | is_production, pacing_mode, source_shot_ids |

## Common Pitfalls

- ❌ Don't add mode logic directly in `storyboard_service.go` — it's already 77KB
- ❌ Don't create new model structs — add fields to existing `Storyboard`
- ❌ Don't forget to update `prompt_i18n.go` if adding new prompt keys
- ❌ Don't assume AI returns valid JSON — always implement cleanup before parsing

## Instrumentation

```bash
./scripts/log-skill.sh "new-storyboard-mode" "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
```
