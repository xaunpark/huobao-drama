# Part-Aware Template for Narrative MV Mode

> Created: 2026-04-21
> Status: Reviewed — Ready for Implementation

## Summary

Create a dedicated "CG5 — Narrative MV" prompt template that separates **Core Visual DNA** (shared across all 3 parts) from **Music-Specific DNA** (only applied to `music_film` shots). This ensures Prologue/Epilogue shots render as pure cinema without 3D text/beat-sync elements, while Music Film shots retain the full CG5 aesthetic.

## Problem Statement

The existing CG5 template (`style_prompt`) contains music-video-specific elements:
- Kinetic Typography (3D floating text)
- Beat-synced camera shake
- Glitch transitions on the beat
- Text that "EXISTS IN 3D SPACE"

When applied uniformly to all shots in Narrative MV mode, Prologue and Epilogue shots (which are **pure film — no music**) incorrectly receive these music elements, breaking the cinematic feel.

The current `StyleDistillService.BatchDistillStyles()` sends the **same** `style_prompt` to the LLM for ALL shots without awareness of `narrative_part` or `has_music`.

## Research Findings

### Current Architecture (Style Injection Flow)

```
Template (style_prompt) ──┬──> StyleDistillService.BatchDistillStyles()
                          │       └── distillImageStyles(style_prompt, allShots)
                          │              └── LLM ──> per-shot image_style (SAME input for all shots)
                          │
                          └──> FramePromptService.generateFirstFrame()
                                 └── resolveStyleForShot(sb)
                                        ├── has image_style? → use it (distilled)
                                        └── no? → fallback to template style_prompt
```

### Key Files

| File | Role | Line Refs |
|------|------|-----------|
| [style_distill_service.go](file:///g:/VS-Project/huobao-drama/application/services/style_distill_service.go#L62-L150) | Batch distillation pipeline | L93: resolves single style_prompt |
| [prompt_template_service.go](file:///g:/VS-Project/huobao-drama/application/services/prompt_template_service.go#L220-L237) | ResolvePromptIfCustom — template field accessor | L253-287: field switch |
| [prompt_template.go](file:///g:/VS-Project/huobao-drama/domain/models/prompt_template.go#L29-L49) | PromptTemplatePrompts struct — all template fields | L41: StylePrompt field |
| [frame_prompt_service.go](file:///g:/VS-Project/huobao-drama/application/services/frame_prompt_service.go#L36-L45) | resolveStyleForShot — uses distilled image_style | L265: distilled path |
| [image_style_distill.txt](file:///g:/VS-Project/huobao-drama/application/prompts/image_style_distill.txt) | LLM prompt template with `%s` for style + shots | L5: style guide placeholder |

### Key Insight

The `image_style_distill.txt` prompt already instructs the LLM to "Omit any style elements that are irrelevant to the specific shot's content" (L19). If we provide `narrative_part` info in the shot context AND conditionally append Music DNA only to music_film shots, the LLM will naturally produce correct per-shot styles.

## Proposed Solution

### Approach: Conditional Style Input per Narrative Part

**No new template field in the database.** Instead:

1. Add a new field `NarrativeMusicDNA` to `PromptTemplatePrompts` — stores music-specific style elements separately from `StylePrompt`.
2. In `BatchDistillStyles`, when shots have `narrative_part`:
   - Build **two style inputs**: `core_style` (from `StylePrompt`) and `full_style` (core + `NarrativeMusicDNA`)
   - Pass `core_style` for prologue/epilogue shots
   - Pass `full_style` for music_film shots
3. Include `narrative_part` in the `shotContext` struct so the distillation LLM has full context.

### Architecture Diagram

```
Template "CG5 — Narrative MV":
    style_prompt         = "3D CG, PBR, volumetric fog, dark noir, colored rim lights..."
    narrative_music_dna  = "Kinetic Typography 3D text, beat-synced camera shake, glitch transitions..."
    video_constraint     = "3D CG animation, natural character movement..."

BatchDistillStyles (narrative_mv mode):
    ┌─────────────────────────────────────────┐
    │  Prologue shots   → style_prompt ONLY   │  → "3D CG, fog, dark noir"
    │  Music Film shots → style_prompt + DNA  │  → "3D CG, fog, dark noir + 3D text, beat sync"
    │  Epilogue shots   → style_prompt ONLY   │  → "3D CG, fog, cold morning light"
    └─────────────────────────────────────────┘
            ↓
    LLM distills per-shot image_style → stored in DB
            ↓
    FramePromptService uses distilled image_style (no change needed)
```

## Acceptance Criteria

- [ ] New `narrative_music_dna` field exists in `PromptTemplatePrompts`
- [ ] `shotContext` struct includes `NarrativePart` field
- [ ] `BatchDistillStyles` conditionally appends Music DNA for `has_music=true` shots
- [ ] Prologue/Epilogue shots receive style WITHOUT music elements (no "3D text", "beat sync", "kinetic typography")
- [ ] Music Film shots receive style WITH music elements
- [ ] Existing non-narrative modes are unaffected (backward compatible)
- [ ] Template UI shows the new `narrative_music_dna` field

## Technical Considerations

### Dependencies
- Requires the Narrative MV backend (already implemented in previous session)
- Requires shots with `narrative_part` and `has_music` fields populated

### Risks
- **Low risk:** Changes are additive. The `narrative_music_dna` field is optional — if empty, distillation behaves exactly as before.
- **Edge case:** What if a template has no `narrative_music_dna` but is used with narrative_mv mode? → Music Film shots will still use `style_prompt` only, which is acceptable (they just won't get music-specific elements). No error.

### Alternatives Considered
- **Dual Template selection** (user picks 2 templates): Rejected — too complex for UX, risk of style inconsistency between parts.
- **Auto-strip via regex** (detect and remove "3D text" keywords): Rejected — brittle, not maintainable.
- **Separate distill calls per part**: Rejected — unnecessary complexity; single call with conditional input is cleaner.

## Implementation Steps

### Task 1: Add `NarrativeMusicDNA` to template model

**File:** `domain/models/prompt_template.go`

Add new field to `PromptTemplatePrompts`:
```go
NarrativeMusicDNA string `json:"narrative_music_dna,omitempty"` // Music-specific style DNA for narrative_mv mode
```

**Note:** Do NOT add an empty string entry to `PromptTypeToDefaultFile` — this field has no default embed file. It only exists in template overrides. The `getPromptFromStruct` case (Task 2) handles retrieval.

### Task 2: Add `getPromptFromStruct` case

**File:** `application/services/prompt_template_service.go`

Add case in the switch:
```go
case "narrative_music_dna":
    return p.NarrativeMusicDNA
```

### Task 3: Add `NarrativePart` to `shotContext`

**File:** `application/services/style_distill_service.go`

```go
type shotContext struct {
    ShotNumber    int    `json:"shot_number"`
    Action        string `json:"action,omitempty"`
    // ... existing fields ...
    NarrativePart string `json:"narrative_part,omitempty"` // NEW
}
```

In `buildShotContexts`:
```go
if sb.NarrativePart != nil {
    ctx.NarrativePart = *sb.NarrativePart
}
```

### Task 4: Conditional style input in `BatchDistillStyles`

**File:** `application/services/style_distill_service.go`

**4a. Add `ResolveNarrativeMusicDNA` helper to `PromptI18n`:**

`PromptI18n` has `templateService` as private field (no getter). Add a public delegating method:

**File:** `application/services/prompt_i18n.go`
```go
// ResolveNarrativeMusicDNA returns the template's narrative_music_dna field if set.
func (p *PromptI18n) ResolveNarrativeMusicDNA(dramaID uint) string {
    if p.templateService != nil && dramaID > 0 {
        return p.templateService.ResolvePromptIfCustom(dramaID, "narrative_music_dna")
    }
    return ""
}
```

**4b. Modify `BatchDistillStyles` in `style_distill_service.go`:**

After resolving `stylePrompt` (line 93), add:
```go
// Resolve narrative music DNA if available
narrativeMusicDNA := s.promptI18n.ResolveNarrativeMusicDNA(dramaID)

// Check if any shots have narrative_part (indicating narrative_mv mode)
hasNarrativeShots := false
for _, sb := range storyboards {
    if sb.NarrativePart != nil && *sb.NarrativePart != "" {
        hasNarrativeShots = true
        break
    }
}
```

Then modify BOTH image and video distillation calls:
- If `hasNarrativeShots && narrativeMusicDNA != ""`:
  - Split shotContexts into 2 groups: `musicShots` (narrative_part == "music_film") vs `nonMusicShots`
  - For image distillation:
    - Call `distillImageStyles(stylePrompt, nonMusicShots)` for prologue/epilogue
    - Call `distillImageStyles(stylePrompt + "\n\n[MUSIC-SPECIFIC STYLE]\n" + narrativeMusicDNA, musicShots)` for music_film
  - For video distillation: apply the same split pattern with `videoConstraint`
  - Merge results arrays
- Else: existing single-call behavior (backward compatible)

### Task 5: Create CG5 Narrative MV template content

Create the template via the Admin UI (or seed script) with:
- `name`: "CG5 — Narrative MV"
- `style_prompt`: Core Visual DNA only (3D CG, PBR, volumetric fog, dark lighting, colored rim lights — NO 3D text, NO beat sync)
- `narrative_music_dna`: Music-specific elements (Kinetic Typography, beat-synced camera shake, glitch transitions, 3D floating text)
- `video_constraint`: Standard CG5 video constraint (3D animation, natural movement etc.)

### Task 6: Frontend — Show `narrative_music_dna` field in Template Editor

**File:** `web/src/views/prompt-template/PromptTemplateEditor.vue` (or equivalent)

Add a new textarea field for `narrative_music_dna` with label "🎵 Music-Specific Style DNA (Narrative MV)" — only visible when editing, optional.

## Verification

1. Create a CG5 Narrative MV template with separated DNA
2. Assign it to a drama project
3. Generate storyboard with `narrative_mv` split mode
4. Trigger style distillation (automatic after split)
5. Verify:
   - Prologue shots' `image_style` contains NO "kinetic typography", "3D text", "beat sync"
   - Music Film shots' `image_style` DOES contain "kinetic typography", "3D text"
   - Epilogue shots' `image_style` contains NO music elements

## References

- Narrative MV implementation plan: [narrative-mv-mode.md](file:///g:/VS-Project/huobao-drama/plans/narrative-mv-mode.md)
- Style distillation plan: `plans/shot-style-distill.md`
- Template style override solution: [template-style-override-distilled](file:///g:/VS-Project/huobao-drama/docs/solutions/logic-errors/template-style-override-distilled-20260419.md)
