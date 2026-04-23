# Episode Character Variant System

> Created: 2026-04-23
> Status: Reviewed ✓ (minor revisions applied)
> Review: See artifacts/plan_review_character_variant.md

## Summary

Nâng cấp `character_extraction` output để mỗi nhân vật có 3 fields mới: `character_prompt` (full text-to-image), `variant_prompt` (delta từ base ref), và `episode_descriptor` (mô tả ngắn cho shot prompts). Sau đó inject dynamic descriptor table vào storyboard/shot prompt generation thay cho bảng hardcoded hiện tại.

## Problem Statement

Khi Bubi mặc đồ chef/pirate/teacher trong 1 episode:
1. **Descriptor cứng**: Template hardcode `Bubi → "a bald toddler in yellow romper"` cho MỌI episode
2. **Không có costume ref**: Pipeline chỉ có ảnh turnaround BASE, không có step gen ảnh Bubi mặc costume mới
3. **Reference conflict**: Upload ref Bubi (yellow romper) + prompt "wearing chef hat" → AI output không nhất quán
4. **Chỉ hỗ trợ 1 workflow**: Hiện tại chỉ có field `appearance` (dùng cho text-to-image), không phân biệt "có base ref" vs "không có base ref"

## Research Findings

### Codebase Patterns

#### Character Model (`domain/models/drama.go:40-64`)
```go
type Character struct {
    Name        string   // Character name
    Role        *string  // main/supporting/animal
    Appearance  *string  // Current: full text-to-image prompt OR "See reference image"
    Personality *string  // Movement style
    Description *string  // Narrative role
    VoiceStyle  *string  // Voice description
    ImageURL    *string  // Generated/uploaded character ref image
    LocalPath   *string  // Local file path
}
```
**Gap**: No `variant_prompt` or `episode_descriptor` fields.

#### Character Extraction (`character_library_service.go:607-614`)
```go
type extractedChar struct {
    Name        string `json:"name"`
    Role        string `json:"role"`
    Appearance  string `json:"appearance"`
    Personality string `json:"personality"`
    Description string `json:"description"`
    VoiceStyle  string `json:"voice_style"`
}
```
**Gap**: No `character_prompt`, `variant_prompt`, or `episode_descriptor` in extraction output.

#### Character Image Generation (`character_library_service.go:334-361`)
- `buildCharacterPrompt()` uses `Appearance` field + appends "t-pose, character sheet..." suffix
- If appearance contains "character sheet" → use as-is (user previously edited the full prompt)
- Supports reference image (I2I mode) via `referenceImageURL` param
**Gap**: No logic to differentiate "generate from text" vs "generate variant from base ref"

#### Storyboard Prompt Injection
- `storyboard_breakdown` template has HARDCODED descriptor table in the template markdown
- Backend resolves via `WithDramaStoryboardSystemPrompt()` → no dynamic character injection
- The descriptor table lives in `cocomelon_template.md` lines 648-661 — static text, NOT built from DB

#### Frontend Character Panel
- `ProfessionalEditor.vue`, `EpisodeWorkflow.vue`, `DramaWorkflow.vue` handle character editing
- Currently: one "Generate Image" button per character
- Uses `GenerateCharacterImage` API with optional reference image

### Key Insight: Static vs Dynamic Descriptors

**Current flow (broken)**:
```
Template hardcodes: Bubi → "bald toddler in yellow romper"
                    ↓ (baked into template text)
storyboard_breakdown uses this regardless of episode
                    ↓
shot prompts always say "yellow romper" even if Bubi is a pirate
```

**Target flow**:
```
character_extraction outputs: episode_descriptor per character
                    ↓
Backend builds dynamic descriptor table from DB
                    ↓
Injects into storyboard_breakdown prompt
                    ↓
shot prompts correctly say "bald toddler in pirate coat"
```

## Proposed Solution

### Phase 1: Data Model (DB Migration)

Add 3 new nullable columns to `characters` table:

```go
type Character struct {
    // ... existing fields ...
    
    // NEW: Full standalone text-to-image prompt (200-400 words)
    // Used when NO base reference image exists
    CharacterPrompt  *string `gorm:"type:text" json:"character_prompt"`
    
    // NEW: Variant delta prompt describing ONLY changes from base ref
    // Used when base reference image EXISTS (costume, injury, transformation...)
    VariantPrompt    *string `gorm:"type:text" json:"variant_prompt"`
    
    // NEW: Short descriptor for this episode (used in shot prompts)
    // e.g. "a bald toddler in chef hat and white jacket"
    EpisodeDescriptor *string `gorm:"type:varchar(500)" json:"episode_descriptor"`
}
```

**Migration**: GORM AutoMigrate adds nullable columns — no data loss, backward compatible.

### Phase 2: Template Update (character_extraction)

Update `character_extraction` template output schema to require 3 new fields:

```json
{
  "name": "Bubi",
  "role": "main",
  "appearance": "See reference image",
  "personality": "Energetic, curious chef",
  "description": "Main toddler learning to cook",
  "voice_style": "Excited toddler",
  "character_prompt": "A bald toddler boy, 2-3 years old, warm honey skin with SSS glow, perfectly smooth round head, very large expressive eyes with star-shaped catchlights, rosy pink gradient cheeks, wearing a white chef toque hat, white double-breasted chef jacket with blue lightning bolt emblem, blue apron, white pants, brown shoes. Standing in t-pose. 3D CGI render, Pixar quality, toy aesthetic, character turnaround sheet, front/side/back views, white background, no text",
  "variant_prompt": "Same character now wearing: white chef toque hat, white double-breasted chef jacket with blue lightning bolt emblem, blue apron, white pants. Keep IDENTICAL bald head, face, skin tone, eye style, and body proportions from reference image.",
  "episode_descriptor": "a bald toddler in chef hat and white chef jacket with blue apron"
}
```

**Rules for AI**:
- `character_prompt`: ALWAYS generated. Extremely detailed (200-400 words). Self-contained.
- `variant_prompt`: ONLY generated when character is from roster AND has costume/state changes. Empty string if no changes.
- `episode_descriptor`: ALWAYS generated. Short (10-30 words). No character names. Used for shot prompts.

### Phase 3: Backend — Fix charDescMap at ALL 3 Save Locations

> **Review Finding**: The actual appearance leak happens at SAVE TIME, not prompt time.
> Three `charDescMap` builders inject `c.Appearance` into stored ImagePrompt/VideoPrompt.

**Step 3a: Extract shared utility** (currently duplicated 3 times):

```go
// buildCharDescMap loads character descriptions with EpisodeDescriptor priority
func buildCharDescMap(tx *gorm.DB, charIDs []uint) map[uint]string {
    descMap := make(map[uint]string)
    if len(charIDs) == 0 { return descMap }
    var chars []models.Character
    if err := tx.Where("id IN ?", charIDs).Find(&chars).Error; err != nil {
        return descMap
    }
    for _, c := range chars {
        desc := c.Name
        if c.EpisodeDescriptor != nil && *c.EpisodeDescriptor != "" {
            desc += " (" + *c.EpisodeDescriptor + ")"
        } else if c.Appearance != nil && *c.Appearance != "" {
            desc += " (" + *c.Appearance + ")"
        } else if c.Description != nil && *c.Description != "" {
            desc += " (" + *c.Description + ")"
        }
        descMap[c.ID] = desc
    }
    return descMap
}
```

**Step 3b: Update all 3 callers:**

| # | File | Line | Method |
|---|---|---|---|
| 1 | `storyboard_service.go` | 1297-1326 | `saveVoiceoverShots` |
| 2 | `storyboard_nursery_service.go` | 555-580 | `saveNurseryRhymeShots` |
| 3 | `storyboard_service.go` | 1867-1896 | `saveStoryboards` |

### Phase 4: Character Image Generation — Dual Mode

Update `GenerateCharacterImage` API to support a `mode` parameter:

| Mode | Prompt Source | Reference Image | Use Case |
|---|---|---|---|
| `text_to_image` | `character_prompt` | None | New characters, MV Maker |
| `variant_from_ref` | `variant_prompt` | Base ref image | Roster + costume/state change |

### Phase 5: Frontend UI

Character panel shows contextual buttons:

```
┌─────────────────────────────────────┐
│  🏴‍☠️ Bubi (Pirate Captain)           │
│                                     │
│  Episode Descriptor:                │
│  [a bald toddler in pirate coat  ] │  ← editable textarea
│                                     │
│  [📝 Generate from Text]           │  ← uses character_prompt
│  [🖼️ Generate with Base Ref]       │  ← uses variant_prompt + ref
│       ↑ disabled if no base ref     │
└─────────────────────────────────────┘
```

## Acceptance Criteria

- [ ] `characters` table has 3 new columns: `character_prompt`, `variant_prompt`, `episode_descriptor`
- [ ] `character_extraction` template outputs all 3 new fields for every character
- [ ] `character_prompt` is detailed enough (200-400 words) for standalone text-to-image
- [ ] `variant_prompt` correctly describes only changes from base ref
- [ ] `episode_descriptor` is short (10-30 words), uses visual descriptors not names
- [ ] Storyboard breakdown receives dynamic descriptor table from DB (not hardcoded)
- [ ] Shot image/video prompts use `episode_descriptor` instead of hardcoded descriptor
- [ ] Character image generation supports dual mode (text-to-image / variant-from-ref)
- [ ] Frontend shows 2 generation buttons with correct availability logic
- [ ] Backward compatible: old characters without new fields still work (fallback to `appearance`)
- [ ] Works for ALL modes: Nursery Rhyme, MV Maker, Narrative MV

## Review Decisions (Resolved)

### Decision 1: `episode_descriptor` scope — Character vs Join Table
**Decision**: Store on `Character` model (not `episode_characters` join table).
**Rationale**: Users work on one episode at a time. Fields represent "current working state" — same as how `appearance` already works. Moving to a proper per-episode join table is a v2 optimization.
**Limitation accepted**: Re-extracting Episode 3 after Episode 4 will overwrite fields. This matches existing behavior.

### Decision 2: Keep 3 fields (not 2)
**Decision**: Keep `character_prompt` separate from `appearance`.
**Rationale**: `appearance` for roster characters is "See reference image" — not a valid T2I prompt. `character_prompt` is ALWAYS a complete standalone prompt. Different semantic purpose.

### Decision 3: StyleDistillService impact
**Decision**: Check `BatchDistillStyles` — if it reads character appearance, update to prefer `EpisodeDescriptor`.

## Technical Considerations

### Dependencies
- GORM AutoMigrate for DB schema change
- Template files: `cocomelon_template.md` (and all other templates)
- Frontend: `ProfessionalEditor.vue`, `EpisodeWorkflow.vue`

### Risks
- **charDescMap duplication**: 3 locations must be updated atomically — extract to shared utility
- **Template update scope**: ALL templates need updated character_extraction section
- **Backward compatibility**: Old episodes without `episode_descriptor` must fallback gracefully
- **Prompt length**: `character_prompt` (200-400 words) could be long for some models

### Alternatives Considered
1. **Separate `costume_prompt` + `state_prompt`**: Rejected — adds complexity without benefit
2. **New `costume_reference_generation` step**: Rejected — only serves 1 use case, not universal
3. **Store descriptors per-shot in storyboard**: Rejected — descriptor is per-character-per-episode
4. **Store on join table**: Deferred to v2 — current approach matches existing patterns

## Implementation Steps

| Task | Description | Estimated |
|---|---|---|
| 1 | DB Migration: Add 3 columns to Character model | 15 min |
| 2 | Template Updates: character_extraction output schema | 45 min |
| 3 | Backend: Parse & store new fields from AI response | 30 min |
| 4 | Backend: Dynamic descriptor injection into storyboard | 1 hour |
| 5 | Backend: Dual mode image generation | 30 min |
| 6 | Frontend: Character panel UI with 2 buttons | 1 hour |
| 7 | API: Update endpoints for new fields | 30 min |
| **Total** | | **~4-5 hours** |

## References
- Character model: `domain/models/drama.go:40-64`
- Character extraction: `application/services/character_library_service.go:557-687`
- Character image gen: `application/services/character_library_service.go:334-417`
- Prompt injection: `application/services/prompt_i18n.go:482-493`
- Template: `docs/cocomelon_template.md:490-536`
- Storyboard descriptor rule: `docs/cocomelon_template.md:648-661`
