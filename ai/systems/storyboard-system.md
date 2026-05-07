# System: Storyboard — Storyboard Generation Engine

> Deep documentation of the storyboard subsystem — the most complex part of the codebase.

## Overview

The storyboard system takes script text and produces a structured shot list using AI. It supports 6 production modes.

## Entry Point

```
POST /api/v1/episodes/:episode_id/storyboards
→ api/handlers/storyboard.go:GenerateStoryboard()
  → application/services/storyboard_service.go (dispatcher)
```

## Mode Dispatch Architecture

```
storyboard_service.go (77KB dispatcher)
  ├── "standard"    → Internal: storyboard_story_breakdown.txt
  ├── "visual_unit" → Internal: storyboard_visual_unit*.txt
  ├── "nursery"     → storyboard_nursery_service.go (24KB)
  ├── "mv"          → storyboard_mv_service.go (10KB)
  ├── "narrative_mv" → storyboard_narrative_service.go (29KB)
  └── "voiceover"   → storyboard_composition_service.go (25KB)
```

## Data Model: Storyboard

Core fields (all modes):
- `EpisodeID`, `StoryboardNumber`, `Title`, `Location`, `Time`
- `ShotType`, `Angle`, `Movement`, `Action`, `Result`, `Atmosphere`
- `ImagePrompt`, `VideoPrompt`, `VideoPromptSource`
- `Dialogue`, `Description`, `Duration`, `Status`

Mode-specific fields (nullable, only populated in respective modes):
- **Rapid Cut**: `IsProduction`, `PacingMode`, `SourceShotIDs`
- **Voiceover Director**: `ScriptSegment`, `NarratorScript`, `AudioMode`, etc.
- **Nursery Rhyme**: `LyricsText`, `SectionType`, `VerseSubject`, `OverlayText`, etc.
- **Narrative MV**: `NarrativePart`, `HasMusic`, `MusicSegment`, `ActingNote`, etc.
- **Style Distill**: `ImageStyle`, `VideoStyle`, `VideoPromptDistilled`

## Supporting Services

### Style Distillation (`style_distill_service.go`, 22KB)
- Generates per-shot visual style from channel template
- Populates `ImageStyle`, `VideoStyle`, `VideoPromptDistilled`
- Called via: `POST /api/v1/episodes/:id/distill-styles`

### Frame Prompt Generation (`frame_prompt_service.go`, 26KB)
- Generates detailed image prompts per frame type
- Frame types: first_frame, key_frame, last_frame, action_sequence
- Called via: `POST /api/v1/storyboards/:id/frame-prompt`

### Rapid Cut (`rapid_cut_service.go`, 18KB)
- Post-processing mode that merges standard shots into fast-paced sequences
- Creates new production storyboards from editorial ones
- Called via: `POST /api/v1/episodes/:id/rapid-cut`

### Storyboard Update (`storyboard_update_full.go`, 7KB)
- Handles full storyboard CRUD updates
- Preserves associations (characters, props)

## Prompt Templates Used

| Mode | System Prompt | Format Prompt |
|------|--------------|---------------|
| Standard | `storyboard_story_breakdown.txt` | `storyboard_format_instructions.txt` |
| Visual Unit | `storyboard_visual_unit.txt` | `storyboard_visual_unit_format.txt` |
| Visual Unit Structured | `storyboard_visual_unit_structured.txt` | `storyboard_visual_unit_format.txt` |
| Nursery Rhyme | `storyboard_nursery_rhyme.txt` | `storyboard_nursery_rhyme_format.txt` |
| MV Gaming Horror | `storyboard_mv_gaming_horror.txt` | `storyboard_format_instructions.txt` |
| MV Cinematic | `storyboard_mv_cinematic_movie.txt` | `storyboard_format_instructions.txt` |
| Narrative Director | `storyboard_narrative_director.txt` | `storyboard_narrative_format.txt` |
| Narrative Planner | `storyboard_narrative_planner.txt` | — |
| Preserve Shots | `storyboard_preserve_shots.txt` | — |

## Critical Warnings

1. `storyboard_service.go` at 77KB is the largest file — never read entirely
2. Mode logic is deeply interleaved — changes affect other modes
3. JSON parsing from AI is fragile — always strip markdown fences
4. Adding new Storyboard fields requires checking all 6 mode paths
5. Storyboard deletion cascades to characters/props associations
