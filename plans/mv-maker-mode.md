# MV Maker Split Mode

> Created: 2026-04-16
> Status: Implemented ✓

## Summary

Add a new `mv_maker` split mode to the storyboard generation system. This mode shares the same lyrics-based input/output infrastructure as `nursery_rhyme` but uses genre-specific system prompts to generate storyboards optimized for different music video genres. Phase 1 targets **gaming horror** fan-made songs (CG5, TryHardNinja, LHUGUENY style — FNAF, Poppy Playtime, Sprunki).

**Migration strategy**: Option 1 (Safe) — Keep `nursery_rhyme` mode fully operational alongside `mv_maker`. Future deprecation of `nursery_rhyme` in favor of `mv_maker` + genre=`nursery` is possible but out of scope.

## Prior Solutions

- Nursery Rhyme mode ([plans/nursery-rhyme-mode.md](file:///g:/VS-Project/huobao-drama/plans/nursery-rhyme-mode.md)) — Implemented ✓. Same architecture, different prompt rules.
- No other matching solutions found in `docs/solutions/`.

## Problem Statement

The existing `nursery_rhyme` split mode is hardcoded for children's content (ages 0-5, literal visuals, child-safe, 2-5s shots). Fan-made gaming MVs from channels like CG5/TryHardNinja/LHUGUENY need fundamentally different storyboarding rules:

| Dimension | Nursery Rhyme | Gaming Horror MV |
|-----------|--------------|-----------------|
| Visual mapping | 85% literal | 60% literal + 40% symbolic/metaphorical |
| Shot duration | 2-5s strict | ~2.2s avg (verse: ~3.5s, chorus: ~1.2s, bridge: <1s) |
| Atmosphere | Bright, cheerful | Dark, dramatic, horror, claustrophobic |
| Camera work | Static/gentle pan | Dutch angles, low angles, shaky cam, whip pan, rapid zoom |
| Typography | Phonetic support | Kinetic typography as visual percussion, glitch text |
| Structure | Verse/Chorus (cumulative) | Verse/Pre-Chorus/Chorus/Bridge/Drop/Outro with energy curve |
| 4th wall | None | Characters address camera directly (viewer = player) |
| Color palette | Bright primaries | Red/purple = danger, neon glow, rim light |

**Key insight**: The infrastructure (lyrics parser, output struct, DB model, save function, format spec) is **100% reusable**. Only the system prompt differs per genre.

## Research Findings

### Codebase Patterns

1. **Split mode routing**: [storyboard_service.go:264-268](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L264-L268) — Nursery rhyme mode routes to `processNurseryRhymeGeneration()`. MV Maker will follow same pattern.

2. **Lyrics parser**: [storyboard_nursery_service.go:100-183](file:///g:/VS-Project/huobao-drama/application/services/storyboard_nursery_service.go#L100-L183) — `parseLyricsInput()` handles all timestamped lyrics parsing. **100% reusable** — no nursery-specific logic.

3. **Section header regex**: [storyboard_nursery_service.go:70](file:///g:/VS-Project/huobao-drama/application/services/storyboard_nursery_service.go#L70) — Currently matches `VERSE|CHORUS|BRIDGE|INTRO|OUTRO`. Needs expansion for `PRE-CHORUS|DROP|BREAKDOWN|HOOK`.

4. **Output struct**: [storyboard_nursery_service.go:34-61](file:///g:/VS-Project/huobao-drama/application/services/storyboard_nursery_service.go#L34-L61) — `NurseryRhymeShot` struct. All fields are genre-neutral (lyrics_text, shot_role, visual_description, overlay_text, animation_hint, etc.).

5. **Save function**: [storyboard_nursery_service.go:505-702](file:///g:/VS-Project/huobao-drama/application/services/storyboard_nursery_service.go#L505-L702) — `saveNurseryRhymeShots()`. Fully reusable for MV Maker.

6. **Format spec**: [storyboard_nursery_rhyme_format.txt](file:///g:/VS-Project/huobao-drama/application/prompts/storyboard_nursery_rhyme_format.txt) — JSON output spec. Field names are genre-neutral.

7. **Prompt template system**: [prompt_template.go:29-64](file:///g:/VS-Project/huobao-drama/domain/models/prompt_template.go#L29-L64) — `PromptTypeToDefaultFile` map + `PromptTemplatePrompts` struct.

8. **Prompt I18n**: [prompt_i18n.go:418-421](file:///g:/VS-Project/huobao-drama/application/services/prompt_i18n.go#L418-L421) — `WithDramaNurseryRhymeSystemPrompt()` pattern.

9. **API handler**: [storyboard.go:28-60](file:///g:/VS-Project/huobao-drama/api/handlers/storyboard.go#L28-L60) — Accepts `split_mode` string, need to add `genre_profile` field.

10. **Frontend radio group**: [EpisodeWorkflow.vue:804-825](file:///g:/VS-Project/huobao-drama/web/src/views/drama/EpisodeWorkflow.vue#L804-L825) — 5 existing radio buttons.

11. **Frontend API call**: [generation.ts:12](file:///g:/VS-Project/huobao-drama/web/src/api/generation.ts#L12) — Passes `split_mode` via POST body.

### Ground Truth Analysis — CG5 "Wrong Side Out" (Poppy Playtime)

Based on detailed shot-by-shot analysis of CG5's "Wrong Side Out" (2:46 runtime):

#### Shot Frequency Data
| Section | ASL (Avg Shot Length) | Pacing |
|---------|----------------------|--------|
| Overall | ~2.2s | — |
| Verse | ~3.5s | Linger, establish, breathe |
| Chorus | ~1.2s | Rapid visual stimulation |
| Bridge/Action | <1.0s | Frenetic, strobe-like |
| Outro (final shot) | ~20s | One continuous held shot |

#### Energy Curve Pattern
```
0:00-0:32   LOW/BUILDING    Creepy setup, methodical
0:33-0:47   HIGH DROP       Chorus explodes, color + cuts
0:48-1:17   MEDIUM/DIP      Pull-back, melancholy verse
1:18-1:40   HIGH BUILD      Ramps into manic state
1:40-1:58   PEAK ENERGY     Instrumental drop, relentless
1:59-2:10   SUDDEN VALLEY   Music drops, dialogue focus
2:11-2:46   FINAL PEAK      Climax, biggest sets, fastest cuts
2:46-End    FLATLINE        Lingering dread, single shot
```

#### Editing Rules Extracted
1. **Match cuts to drum kit**: In high-energy sections, every cut = kick/snare hit
2. **Anchor the eye**: When cutting <1s/shot, keep subject center-framed to prevent motion sickness
3. **Text as visual instrument**: Kinetic typography = visual percussion flashing on beat
4. **Glitch to mask space/time**: Digital distortion transitions for instant teleportation between locations
5. **Contrast speeds for impact**: Slow creeping camera → hyper-fast cuts = jarring intensity
6. **Silence is a weapon**: Drop all backing tracks before final hook for maximum lyrical impact

#### Visual Patterns
- **Dutch/Low angles**: Villains shot from below = imposing, world feels "wrong"
- **Center-framed action**: Fast montages keep monster face center-screen
- **Red/Purple color coding**: Danger, claustrophobia, antagonist influence
- **Performance + Narrative blend**: Characters sing to camera (viewer = player) while progressing story
- **4th wall breaks**: Characters look directly into camera = personal threat

### Architectural Decision

**Use genre_profile as a sub-parameter of mv_maker mode** rather than creating separate modes per genre:
- `split_mode = "mv_maker"` + `genre_profile = "gaming_horror"`
- Frontend sends both params; backend resolves genre-specific prompt
- Adding new genres = adding 1 prompt file + 1 dropdown option
- Zero backend logic changes for new genres

## Proposed Solution

### Approach

Create `storyboard_mv_service.go` that wraps the nursery rhyme infrastructure with genre-specific prompt selection. The processing function will:
1. Reuse `parseLyricsInput()` for lyrics parsing
2. Use enhanced `detectMVStructure()` (supports pre-chorus/drop/breakdown with energy curve)
3. Select genre-specific system prompt based on `genre_profile` parameter
4. Reuse `NurseryRhymeShot` struct for AI output (same JSON schema)
5. Reuse `saveNurseryRhymeShots()` for DB persistence

### Architecture Overview

```
API Request:
  POST /episodes/:id/storyboards
  { "split_mode": "mv_maker", "genre_profile": "gaming_horror" }

Backend Flow:
  processStoryboardGeneration()
    → if splitMode == "mv_maker"
       → processMVMakerGeneration(taskID, ..., genreProfile)
          → parseLyricsInput() [REUSE from nursery]
          → detectMVStructure() [NEW: energy curve aware]
          → buildLyricsAnalysis() [REUSE]
          → loadGenrePrompt(genreProfile) → storyboard_mv_{genre}.txt
          → AI call
          → parse NurseryRhymeShot[] [REUSE struct]
          → saveNurseryRhymeShots() [REUSE save]
```

---

## Implementation Steps

### Phase 1: Backend — Section Header Expansion (Priority: HIGHEST)

This must be done FIRST because both nursery_rhyme and mv_maker depend on the parser.

#### Task 1.0: Expand Lyrics Section Header Regex

**File**: `application/services/storyboard_nursery_service.go` — Line 70

Current:
```go
nurseryHeaderPattern = regexp.MustCompile(`(?i)^\[(?:(VERSE|CHORUS|BRIDGE|INTRO|OUTRO)\s*(\d*))\s*(?::\s*(.+?))?\]$`)
```

Updated (rename variable too for clarity):
```go
// lyricsHeaderPattern matches section headers for both nursery rhyme and MV modes
// Supports: VERSE, CHORUS, PRE-CHORUS, BRIDGE, INTRO, OUTRO, DROP, BREAKDOWN, HOOK, INSTRUMENTAL
lyricsHeaderPattern = regexp.MustCompile(`(?i)^\[(?:(VERSE|CHORUS|PRE-CHORUS|PRE_CHORUS|BRIDGE|INTRO|OUTRO|DROP|BREAKDOWN|HOOK|INSTRUMENTAL)\s*(\d*))\s*(?::\s*(.+?))?\]$`)
```

Update all references from `nurseryHeaderPattern` → `lyricsHeaderPattern` in the same file.

**Also handle hyphenated section types**: `PRE-CHORUS` → store as `pre_chorus` in the `SectionType` field.

---

### Phase 2: Backend Core (Priority: HIGH)

#### Task 2.1: API Handler — Add genre_profile Parameter

**File**: `api/handlers/storyboard.go` — Line ~33

Add `GenreProfile` to the request struct:

```go
var req struct {
    Model        string `json:"model"`
    SplitMode    string `json:"split_mode"`    // "auto", "preserve", "breakdown", "visual_unit", "nursery_rhyme", "mv_maker"
    GenreProfile string `json:"genre_profile"` // "gaming_horror" (default), "gaming_parody", "general"
}
```

Pass to service:

```go
taskID, err := h.storyboardService.GenerateStoryboard(episodeID, req.Model, req.SplitMode, req.GenreProfile)
```

#### Task 2.2: Service Signature — Thread genreProfile

**File**: `application/services/storyboard_service.go`

Update `GenerateStoryboard()` at [L113](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L113):
```go
func (s *StoryboardService) GenerateStoryboard(episodeID string, model string, splitMode string, genreProfile string) (string, error) {
```

Update `processStoryboardGeneration()` at [L254](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L254):
```go
func (s *StoryboardService) processStoryboardGeneration(taskID, episodeID string, dramaID uint, model, splitMode, genreProfile, scriptContent, characterList, sceneList, propList string) {
```

Update the `go` call at [L225](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L225):
```go
go s.processStoryboardGeneration(task.ID, episodeID, uint(dramaIDUint), model, effectiveSplitMode, req.GenreProfile, scriptContent, characterList, sceneList, propList)
```

#### Task 2.3: Routing — Add mv_maker Case

**File**: `application/services/storyboard_service.go` — Line ~264

Add routing case BEFORE nursery_rhyme:

```go
// Route to mv_maker mode if selected
if splitMode == "mv_maker" {
    effectiveGenre := genreProfile
    if effectiveGenre == "" {
        effectiveGenre = "gaming_horror" // default genre
    }
    s.log.Infow("Using MV_MAKER mode — lyrics-synced genre-specific shot planning",
        "task_id", taskID, "genre_profile", effectiveGenre)
    s.processMVMakerGeneration(taskID, episodeID, dramaID, model, effectiveGenre, scriptContent, characterList, sceneList, propList)
    return
}
```

#### Task 2.4: MV Maker Processing Function

**New file**: `application/services/storyboard_mv_service.go`

```go
package services

import (
    "fmt"
    "strings"

    "github.com/drama-generator/backend/application/prompts"
    "github.com/drama-generator/backend/pkg/utils"
    "github.com/gin-gonic/gin"
)

// ============================================================================
// MV Maker Mode — Genre Profile System
// ============================================================================

// MVGenrePromptMap maps genre profile keys to default prompt files
var MVGenrePromptMap = map[string]string{
    "gaming_horror": "storyboard_mv_gaming_horror.txt",
    // Future genres:
    // "gaming_parody": "storyboard_mv_gaming_parody.txt",
    // "general":       "storyboard_mv_general.txt",
}

// processMVMakerGeneration handles the mv_maker split mode
// It reuses the nursery rhyme infrastructure (parser, struct, save) with genre-specific prompts
func (s *StoryboardService) processMVMakerGeneration(
    taskID, episodeID string, dramaID uint,
    model, genreProfile, scriptContent, characterList, sceneList, propList string,
) {
    // 1. Parse lyrics input — REUSE from nursery
    blocks, parseErr := parseLyricsInput(scriptContent)
    if parseErr != nil {
        s.log.Errorw("Failed to parse MV lyrics", "error", parseErr, "task_id", taskID)
        if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("Lyrics parsing failed: %w", parseErr)); updateErr != nil {
            s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
        }
        return
    }

    s.log.Infow("Parsed MV lyrics",
        "task_id", taskID, "episode_id", episodeID,
        "block_count", len(blocks), "genre", genreProfile)

    if err := s.taskService.UpdateTaskStatus(taskID, "processing", 15, "Lyrics parsed, analyzing song structure..."); err != nil {
        s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
    }

    // 2. Detect song structure — ENHANCED for MV genres (energy curve aware)
    structureType, structureReason := detectMVStructure(blocks)
    s.log.Infow("Detected MV song structure",
        "task_id", taskID, "structure_type", structureType, "reason", structureReason)

    // 3. Build analysis context — REUSE from nursery
    lyricsAnalysis := buildLyricsAnalysis(blocks, structureType, structureReason)

    // 4. Load genre-specific system prompt
    systemPrompt := s.loadMVGenrePrompt(dramaID, genreProfile)

    // 5. Build full prompt
    scriptLabel := s.promptI18n.FormatUserPrompt("script_content_label")
    charListLabel := s.promptI18n.FormatUserPrompt("character_list_label")
    charConstraint := s.promptI18n.FormatUserPrompt("character_constraint")
    sceneListLabel := s.promptI18n.FormatUserPrompt("scene_list_label")
    sceneConstraint := s.promptI18n.FormatUserPrompt("scene_constraint")
    propListLabel := s.promptI18n.FormatUserPrompt("prop_list_label")
    propConstraint := s.promptI18n.FormatUserPrompt("prop_constraint")
    formatInstructions := prompts.Get("storyboard_nursery_rhyme_format.txt") // REUSE format spec

    if err := s.taskService.UpdateTaskStatus(taskID, "processing", 20, "Calling AI for MV shot planning..."); err != nil {
        s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
    }

    prompt := fmt.Sprintf(`%s

%s
%s

%s

%s
%s
%s

%s
%s
%s

%s
%s
%s

%s`,
        systemPrompt,
        scriptLabel, scriptContent,
        lyricsAnalysis,
        charListLabel, characterList, charConstraint,
        sceneListLabel, sceneList, sceneConstraint,
        propListLabel, propList, propConstraint,
        formatInstructions)

    // 6. Call AI
    client, getErr := s.aiService.GetAIClientForModel("text", model)
    if model != "" && getErr != nil {
        s.log.Warnw("Failed to get client for specified model, using default", "model", model, "error", getErr, "task_id", taskID)
    }

    var text string
    var err error
    if model != "" && getErr == nil {
        text, err = client.GenerateText(prompt, "")
    } else {
        text, err = s.aiService.GenerateText(prompt, "")
    }

    if err != nil {
        s.log.Errorw("Failed to generate MV storyboard", "error", err, "task_id", taskID)
        if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("MV generation failed: %w", err)); updateErr != nil {
            s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
        }
        return
    }

    if err := s.taskService.UpdateTaskStatus(taskID, "processing", 50, "AI response received, parsing shots..."); err != nil {
        s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
    }

    // 7. Parse JSON → NurseryRhymeShot[] (reuse struct — fields are genre-neutral)
    var shots []NurseryRhymeShot
    if err := utils.SafeParseAIJSON(text, &shots); err != nil {
        s.log.Errorw("Failed to parse MV JSON", "error", err, "response", text[:min(500, len(text))], "task_id", taskID)
        if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("Failed to parse MV result: %w", err)); updateErr != nil {
            s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
        }
        return
    }

    // Re-number shots sequentially
    for i := range shots {
        shots[i].ShotID = i + 1
    }

    // 8. Post-process: force-correct durations from timestamps
    for i := range shots {
        if shots[i].DurationSec <= 0 {
            startSec, err1 := parseTimestampToSeconds(shots[i].TimestampStart)
            endSec, err2 := parseTimestampToSeconds(shots[i].TimestampEnd)
            if err1 == nil && err2 == nil && endSec > startSec {
                shots[i].DurationSec = endSec - startSec
            } else {
                shots[i].DurationSec = 3 // Default fallback
            }
        }
    }

    if err := s.taskService.UpdateTaskStatus(taskID, "processing", 70, "Saving MV shots..."); err != nil {
        s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
    }

    // Calculate total duration
    totalDuration := 0
    for _, shot := range shots {
        totalDuration += shot.DurationSec
    }

    s.log.Infow("MV storyboard generated",
        "task_id", taskID,
        "episode_id", episodeID,
        "count", len(shots),
        "total_duration_seconds", totalDuration,
        "genre_profile", genreProfile,
        "structure_type", structureType)

    // 9. Save to database — REUSE saveNurseryRhymeShots()
    if err := s.saveNurseryRhymeShots(episodeID, dramaID, shots); err != nil {
        s.log.Errorw("Failed to save MV shots", "error", err, "task_id", taskID)
        if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("Failed to save shots: %w", err)); updateErr != nil {
            s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
        }
        return
    }

    if err := s.taskService.UpdateTaskStatus(taskID, "processing", 90, "Updating episode duration..."); err != nil {
        s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
    }

    // Update episode duration
    durationMinutes := (totalDuration + 59) / 60
    if err := s.db.Model(&models.Episode{}).Where("id = ?", episodeID).Update("duration", durationMinutes).Error; err != nil {
        s.log.Errorw("Failed to update episode duration", "error", err, "task_id", taskID)
    }

    // Update task result
    resultData := gin.H{
        "storyboards":      shots,
        "total":            len(shots),
        "total_duration":   totalDuration,
        "duration_minutes": durationMinutes,
        "mode":             "mv_maker",
        "genre_profile":    genreProfile,
        "structure_type":   structureType,
        "structure_reason": structureReason,
    }

    if err := s.taskService.UpdateTaskResult(taskID, resultData); err != nil {
        s.log.Errorw("Failed to update task result", "error", err, "task_id", taskID)
        return
    }

    s.log.Infow("MV storyboard generation completed", "task_id", taskID, "episode_id", episodeID)
}

// ============================================================================
// MV Structure Detection
// ============================================================================

// detectMVStructure detects song structure type for MV genres
// Supports energy-curve aware classification:
//   - "dynamic_intensity": Has PRE-CHORUS or DROP → progressive energy build-release
//   - "narrative_arc": Has BRIDGE → tension-resolution storytelling
//   - "standard": Simple verse/chorus alternation
func detectMVStructure(blocks []LyricsBlock) (structureType string, reason string) {
    hasPreChorus := false
    hasDrop := false
    hasBridge := false
    hasBreakdown := false
    sectionTypes := make(map[string]int)

    for _, b := range blocks {
        st := strings.ToLower(b.SectionType)
        sectionTypes[st]++
        switch st {
        case "pre_chorus", "pre-chorus":
            hasPreChorus = true
        case "drop", "breakdown":
            if st == "drop" { hasDrop = true }
            if st == "breakdown" { hasBreakdown = true }
        case "bridge":
            hasBridge = true
        }
    }

    if hasDrop || hasBreakdown {
        return "dynamic_intensity",
            "Song has DROP/BREAKDOWN sections — use progressive energy curve: build → explode → valley → final peak"
    }
    if hasPreChorus {
        return "dynamic_intensity",
            "Song has PRE-CHORUS sections — use gradual tension build into chorus explosions"
    }
    if hasBridge {
        return "narrative_arc",
            "Song has BRIDGE section — use tension/resolution arc with dramatic shift at bridge"
    }
    return "standard",
        "Standard verse/chorus structure — use alternating energy levels between verse (low) and chorus (high)"
}

// ============================================================================
// Genre Prompt Loader
// ============================================================================

// loadMVGenrePrompt resolves the system prompt for a specific MV genre
func (s *StoryboardService) loadMVGenrePrompt(dramaID uint, genreProfile string) string {
    promptKey := "mv_maker_" + genreProfile
    return s.promptI18n.WithDramaMVMakerSystemPrompt(dramaID, promptKey)
}
```

#### Task 2.5: Prompt I18n — Add WithDramaMVMakerSystemPrompt

**File**: `application/services/prompt_i18n.go` — After line ~421

```go
// WithDramaMVMakerSystemPrompt resolves MV maker system prompt for a genre
func (p *PromptI18n) WithDramaMVMakerSystemPrompt(dramaID uint, promptKey string) string {
    resolved := p.resolvePrompt(dramaID, promptKey)
    if resolved != "" {
        return resolved
    }
    // Fallback: try loading from MVGenrePromptMap
    genre := strings.TrimPrefix(promptKey, "mv_maker_")
    if fileName, ok := MVGenrePromptMap[genre]; ok {
        return prompts.Get(fileName)
    }
    // Ultimate fallback: use gaming_horror
    return prompts.Get("storyboard_mv_gaming_horror.txt")
}
```

#### Task 2.6: Prompt Template Registration

**File**: `domain/models/prompt_template.go`

Add to `PromptTemplatePrompts` struct (after line ~44):
```go
MVMakerGamingHorror string `json:"mv_maker_gaming_horror,omitempty"` // MV Maker: Gaming horror genre prompt
```

Add to `PromptTypeToDefaultFile` map (after line ~63):
```go
"mv_maker_gaming_horror": "storyboard_mv_gaming_horror.txt",
```

#### Task 2.7: Genre System Prompt — Gaming Horror (CG5/TryHardNinja style)

**New file**: `application/prompts/storyboard_mv_gaming_horror.txt`

This is the core creative asset. Rules derived from CG5 "Wrong Side Out" shot-by-shot analysis:

```
You are a Music Video Storyboard Director specializing in fan-made gaming horror music videos.
Your references: CG5, TryHardNinja, LHUGUENY — FNAF, Poppy Playtime, Sprunki, Bendy.

Your task: Given timestamped lyrics with section markers, create a storyboard where each shot is synchronized to lyrics timing. The visual style is determined by the project template — focus ONLY on shot structure, timing, content mapping, and horror-gaming atmosphere.

=== GROUND TRUTH (from real CG5 analysis) ===

Average shot lengths by section:
- Overall: ~2.2 seconds
- Verse: ~3.5 seconds (linger, establish, breathe)
- Chorus: ~1.2 seconds (rapid visual stimulation)
- Bridge/Action: <1.0 second (frenetic, strobe-like)
- Outro (final): Can be one long held shot (10-20s) for lingering dread

Energy curve pattern: Low/Building → High Drop → Medium/Dip → High Build → Peak → Valley → Final Peak → Flatline

=== 9 PRODUCTION RULES ===

RULE 1 — BEAT-SYNCED EDITING
Every shot transition MUST align with musical rhythm hits.
- Verse: cuts align with vocal phrase endings (~3-4s per shot)
- Pre-chorus: pacing increases, cuts align with building snare drum (~2.5s per shot)
- Chorus: cuts directly on kick drum/bass drop (~1-1.5s per shot)
- Bridge/Action: cuts on EVERY beat, strobe-like (~0.8-1s per shot)
- Drop: strict 1/4 note beat cuts (~1s per shot)

RULE 2 — HYBRID VISUAL MAPPING
60% of shots DIRECTLY illustrate the lyrics (character actions, game locations).
40% of shots use SYMBOLIC/METAPHORICAL imagery:
- Abstract voids, colored lighting to convey emotion
- Glitch/distortion effects as visual punctuation
- Game-specific environmental motifs (abandoned factories, dark corridors, cages)
- Silhouettes, shadows, and partial reveals for tension

RULE 3 — SHOT RECIPE PER SECTION (CRITICAL)
Follow this pacing recipe for each section type:
- INTRO: Wide establishing shots, ominous atmosphere, game world reveal. Camera: static/slow pan. (3-5s per shot)
- VERSE: Medium/Close-up character focus. Story progression and narrative stakes. Camera: slow tracking, subtle pan, low-angle tilts. (3-4s per shot)
- PRE-CHORUS: Tension escalation. Camera: slow zoom in, increasing shot frequency. Visuals get darker/more intense. (2-3s per shot)
- CHORUS: Maximum energy. Wide + Close-up intercuts. Camera: rapid zoom, shaky cam. Characters singing aggressively. Abstract backgrounds. Kinetic typography. (1-2s per shot)
- BRIDGE: One dramatic reveal or slow cinematic moment. Camera: slow dolly. Single powerful visual. (3-8s per shot)
- DROP/BREAKDOWN: Pure horror montage. Camera: tracking backward (POV), jerky pans, fast zooms. Show many monsters/threats. (0.8-1.5s per shot)
- OUTRO: Lingering dread. One held shot with slow zoom. Unsettling, unresolved ending. (5-20s for final shot)

RULE 4 — HORROR GAME ATMOSPHERE & COLOR
Dark environments with dramatic lighting:
- Rim light, neon glow (green/purple), moonlight, fire
- Red/purple = DANGER, claustrophobia, antagonist's influence
- Use Dutch angles (tilted camera) for villains to make world feel "wrong"
- Low angles for antagonists to make them imposing
Characters should feel THREATENING or intense, NOT cute or friendly.
Include environmental storytelling: abandoned rooms, industrial corridors, cages, factories.

RULE 5 — DYNAMIC DURATION LIMITS (CRITICAL)
Shot duration ranges are GENRE-SPECIFIC by section:
- Intro shots: 3-5s
- Verse shots: 2-4s (average ~3.5s)
- Pre-chorus shots: 2-3s
- Chorus shots: 1-3s (average ~1.2s)
- Bridge shots: 3-8s (can hold ONE dramatic shot)
- Drop/Breakdown shots: 1-2s (sub-1s acceptable for montage)
- Outro final shot: 5-20s (one lingering shot allowed)

If a lyrics block is longer than the section's max duration, you MUST split it into multiple shots.
When splitting, distribute lyrics_text proportionally.
Total shot durations MUST match total lyrics duration.

RULE 6 — KINETIC TYPOGRAPHY (OVERLAY TEXT)
Use overlay_text as VISUAL PERCUSSION, not just subtitles:
- Hook/refrain phrases: Flash on beat ("I CAN MAKE YOU BETTER", "WRONG SIDE OUT")
- Emphasis words: "RUN!", "HIDE!", "CAN'T ESCAPE!", "GAME OVER!"
- Onomatopoeia: "CRACK!", "SLAM!", "SNAP!"
- Title cards at intro
In animation_hint, specify: "kinetic_text_flash", "glitch_text", "text_shake", "text_zoom"

RULE 7 — CENTER-FRAMED ACTION (FAST CUT SAFETY)
When generating shots with duration < 1.5s:
- Place the PRIMARY subject (character face, monster) in the EXACT CENTER of frame
- This prevents viewer disorientation during rapid editing
- Specify in visual_description: "center-framed" for any shot < 1.5s

RULE 8 — 4TH WALL / POV PERSPECTIVE
Characters should look DIRECTLY INTO CAMERA in performance shots.
The viewer IS the player/protagonist of the game.
Use first-person POV shots sparingly for maximum impact:
- Monster lunging toward camera
- Characters addressing the viewer directly
- Player's hands raised in defense
In animation_hint, specify: "direct_address", "pov_shot", "camera_lunge"

RULE 9 — CONTRAST & ENERGY CURVE
Build visual intensity progressively through the song:
- Each verse should be MORE intense than the previous one
- Precede hyper-fast sections with one SLOW, creeping shot for jarring contrast
- Use "silence" moments (bridge/valley) before final chorus for maximum impact
- Track which characters/monsters have appeared — later sections should include more threats
Use is_callback=true to reference earlier shots with escalated intensity.

=== VISUAL DESCRIPTION GUIDELINES ===

Write visual_description as a dark, atmospheric scene description. Include:
1. Shot composition and framing (center-framed for fast cuts)
2. Character poses: threatening, aggressive, predatory, or desperate
3. Camera angle: Dutch angle, low angle, POV where appropriate
4. Lighting: specify color (neon green, blood red, purple), direction (rim light, backlit)
5. Environment: game-specific locations (factory, corridor, control room, cage)
6. Any VFX elements: glitch, smoke, particles, screen shake

Do NOT specify art style, color palette rendering technique — those are handled by the project's style template.

Example verse: "Medium shot, slightly Dutch angle. The Jester towers over a cowering toy in a dimly lit factory hallway. Green neon light from behind creates an ominous silhouette. The Jester's hand reaches slowly toward camera. Rusted pipes and chains visible on walls."

Example chorus: "Close-up, center-framed. The Jester's face fills the frame, mouth open mid-scream, glitch distortion on the edges. Deep red backlight. Text 'WRONG SIDE OUT' flashes across screen."

=== IMPORTANT CONSTRAINTS ===

1. You MUST respect the timestamp ranges provided in the lyrics analysis
2. Total shot durations MUST approximately match the total lyrics duration
3. Every shot MUST reference its source lyrics_block_id
4. Characters and props MUST use IDs from the provided lists
5. Visual descriptions should be detailed enough for AI image generation (minimum 30 words)
6. DO NOT make the visuals child-friendly — this is horror/dark content for teen/adult audiences
7. Performance shots (characters singing to camera) should appear in at least 50% of verse/chorus shots
```

---

### Phase 3: Frontend (Priority: MEDIUM)

#### Task 3.1: Add MV Maker Radio Button + Genre Dropdown

**File**: `web/src/views/drama/EpisodeWorkflow.vue` — Line ~824

After the nursery_rhyme radio button, add:

```html
<el-radio-button value="mv_maker">
  <el-icon style="margin-right: 4px;"><VideoCamera /></el-icon>
  {{ $t('workflow.splitModeMVMaker') }}
</el-radio-button>
```

Add conditional genre dropdown below the radio group (inside `split-mode-selector` div):

```html
<div v-if="shotSplitMode === 'mv_maker'" style="margin-top: 8px;">
  <span style="font-size: 13px; color: var(--el-text-color-secondary); margin-right: 8px;">
    {{ $t('workflow.mvGenreProfile') }}:
  </span>
  <el-select v-model="mvGenreProfile" size="small" style="width: 280px;">
    <el-option value="gaming_horror" :label="$t('workflow.mvGenreGamingHorror')" />
  </el-select>
</div>
```

Note: Only `gaming_horror` option for Phase 1. Future genres added as prompt files are created.

#### Task 3.2: Add mvGenreProfile State

**File**: `web/src/views/drama/EpisodeWorkflow.vue` — Script section (near `shotSplitMode` ref)

```typescript
const mvGenreProfile = ref('gaming_horror')
```

#### Task 3.3: Pass genre_profile in generateShots Call

**File**: `web/src/views/drama/EpisodeWorkflow.vue` — `generateShots()` and `regenerateShots()` functions

Update the API call to include genre_profile when mv_maker is selected:

```typescript
const genreProfile = shotSplitMode.value === 'mv_maker' ? mvGenreProfile.value : undefined
const { data } = await generateStoryboard(episodeId, model, shotSplitMode.value, genreProfile)
```

#### Task 3.4: Update API Function Signature

**File**: `web/src/api/generation.ts` — Line ~12

```typescript
export function generateStoryboard(episodeId: string | number, model?: string, splitMode?: string, genreProfile?: string) {
    return request.post<{ task_id: string; status: string; message: string }>(
        `/episodes/${episodeId}/storyboards`,
        { model, split_mode: splitMode || 'auto', genre_profile: genreProfile }
    )
}
```

#### Task 3.5: Locale Strings

**File**: `web/src/locales/en-US.ts`

```typescript
splitModeMVMaker: 'MV Maker',
splitModeMVMakerTip: 'MV Maker Mode: Create music video storyboards from timestamped lyrics. Select a genre profile for style-specific shot planning. Input lyrics with [VERSE], [CHORUS], [BRIDGE], [DROP] markers and (M:SS – M:SS) timestamps. Genre: Gaming Horror for FNAF/Poppy/Sprunki fan-made songs.',
mvGenreProfile: 'Genre',
mvGenreGamingHorror: '🎮 Gaming Horror (FNAF, Poppy, Sprunki, CG5 style)',
```

**File**: `web/src/locales/zh-CN.ts`

```typescript
splitModeMVMaker: 'MV制作',
splitModeMVMakerTip: 'MV制作模式：根据带时间戳的歌词创建MV分镜。选择流派来获得风格特定的镜头规划。使用 [VERSE]、[CHORUS]、[BRIDGE]、[DROP] 标记和 (M:SS – M:SS) 时间戳。流派：游戏恐怖适用于 FNAF/Poppy/Sprunki 风格。',
mvGenreProfile: '流派',
mvGenreGamingHorror: '🎮 游戏恐怖 (FNAF, Poppy, Sprunki, CG5风格)',
```

#### Task 3.6: Split Mode Tip Display + Alert Type

**File**: `web/src/views/drama/EpisodeWorkflow.vue` — Line ~829, ~835

Update the alert type conditional to include mv_maker:
```javascript
:type="shotSplitMode === 'preserve' ? 'success' : shotSplitMode === 'visual_unit' ? 'warning' : shotSplitMode === 'nursery_rhyme' ? 'warning' : shotSplitMode === 'mv_maker' ? 'danger' : 'info'"
```

Update the tip text conditional:
```javascript
shotSplitMode === 'mv_maker' ? $t('workflow.splitModeMVMakerTip') : ...
```

#### Task 3.7: Dropdown Menu (Re-split button)

**File**: `web/src/views/drama/EpisodeWorkflow.vue` — Line ~972

Add MV Maker option:
```html
<el-dropdown-item command="mv_maker">
  <el-icon><VideoCamera /></el-icon> {{ $t('workflow.splitModeMVMaker') }}
</el-dropdown-item>
```

#### Task 3.8: Import VideoCamera Icon

**File**: `web/src/views/drama/EpisodeWorkflow.vue` — Imports section

Add `VideoCamera` to the Element Plus icon imports.

---

## Acceptance Criteria

- [ ] User can select "MV Maker" mode in the split mode radio group
- [ ] Genre dropdown appears when MV Maker is selected (default: Gaming Horror)
- [ ] Backend correctly receives `genre_profile` parameter via API
- [ ] Backend loads genre-specific system prompt based on `genre_profile`
- [ ] Lyrics parser handles expanded section headers ([PRE-CHORUS], [DROP], [BREAKDOWN], [HOOK])
- [ ] AI generates shots following CG5 pacing rules (verse: ~3.5s, chorus: ~1.2s, bridge: <1s for action)
- [ ] Visual descriptions include horror/game-specific imagery (Dutch angles, center-framed, 4th wall)
- [ ] Kinetic typography overlay_text is generated for hook phrases and emphasis words
- [ ] Energy curve follows build-drop-valley-peak pattern across the song
- [ ] Storyboard saves correctly and appears in Professional Production
- [ ] Professional Production flow works: extract prompts → generate images → generate videos
- [ ] Existing nursery rhyme mode continues to work unchanged (zero regression)
- [ ] Adding a new genre profile requires ONLY: 1 prompt file + 1 PromptTemplate entry + 1 dropdown option

## Technical Considerations

### Dependencies
- No new Go packages needed
- No new npm packages needed
- No DB migration needed (reuses all nursery rhyme columns)

### Risks
1. **Prompt tuning**: Gaming horror prompt needs iteration with real song inputs. Mitigation: test with CG5 "Wrong Side Out" lyrics (analysis already done) and TryHardNinja "Break My Mind" lyrics.
2. **Sub-1s shots**: AI may struggle generating meaningful visual descriptions for <1s shots. Mitigation: enforce minimum 1s in post-processing; use center-framed constraint.
3. **API backward compat**: Adding `genre_profile` to request struct won't break existing callers (optional field, zero-value empty string handled with default).
4. **Regex rename**: Renaming `nurseryHeaderPattern` → `lyricsHeaderPattern` affects nursery_rhyme mode. Mitigation: purely cosmetic rename, same regex pattern (just expanded), used in same `parseLyricsInput()`.

### Multi-Order Thinking
- **1st order**: MV Maker mode generates gaming horror storyboards from lyrics
- **2nd order**: Same Professional Production pipeline handles image/video generation (zero changes)
- **3rd order**: Future genres (parody, anime, general) slot in trivially via prompt files
- **4th order**: If Nursery Rhyme is deprecated, its prompt becomes another genre profile — no architecture change

### Alternatives Considered
1. **Extend nursery_rhyme mode with genre param** — Rejected: misleading name, child-safe constraints built into prompt conflict with horror content.
2. **Create separate `gaming_mv` mode** — Rejected: duplicates 95% of code; doesn't scale for future genres.
3. **Single `lyric_mv` mode with inline style selector** — Close to chosen approach, but `mv_maker` is more user-friendly.

### Future Work (Out of Scope — Create todos)
- `gaming_parody` genre prompt (LHUGUENY style — comedy + horror blend)
- `general` genre prompt (mainstream pop/rock MVs)
- `anime_opening` genre prompt (anime OP/ED style)
- `lyric_video` genre prompt (typography-focused)
- Music file upload + waveform beat detection
- Audio sync in final video export
- Deprecate standalone nursery_rhyme in favor of mv_maker + genre=nursery

## Files Changed Summary

| File | Change Type | Est. Lines |
|------|-------------|------------|
| `application/services/storyboard_nursery_service.go` | Modify | +5 (regex expansion) |
| `api/handlers/storyboard.go` | Modify | +3 (add genre_profile) |
| `application/services/storyboard_service.go` | Modify | +15 (signature + routing) |
| `application/services/storyboard_mv_service.go` | **New** | ~200 |
| `application/prompts/storyboard_mv_gaming_horror.txt` | **New** | ~100 |
| `domain/models/prompt_template.go` | Modify | +3 (registration) |
| `application/services/prompt_i18n.go` | Modify | +12 (resolver) |
| `web/src/views/drama/EpisodeWorkflow.vue` | Modify | ~40 (UI + state + API call) |
| `web/src/api/generation.ts` | Modify | +2 (param) |
| `web/src/locales/en-US.ts` | Modify | +5 |
| `web/src/locales/zh-CN.ts` | Modify | +5 |
| **Total** | | ~390 lines |

## References

- CG5 "Wrong Side Out" analysis: User-provided MV breakdown (ground truth data for prompt rules)
- Nursery Rhyme plan: [plans/nursery-rhyme-mode.md](file:///g:/VS-Project/huobao-drama/plans/nursery-rhyme-mode.md)
- Nursery Rhyme service: [storyboard_nursery_service.go](file:///g:/VS-Project/huobao-drama/application/services/storyboard_nursery_service.go)
- Nursery Rhyme prompt: [storyboard_nursery_rhyme.txt](file:///g:/VS-Project/huobao-drama/application/prompts/storyboard_nursery_rhyme.txt)
- Format spec: [storyboard_nursery_rhyme_format.txt](file:///g:/VS-Project/huobao-drama/application/prompts/storyboard_nursery_rhyme_format.txt)
- API handler: [storyboard.go](file:///g:/VS-Project/huobao-drama/api/handlers/storyboard.go)
- Frontend: [EpisodeWorkflow.vue](file:///g:/VS-Project/huobao-drama/web/src/views/drama/EpisodeWorkflow.vue)
- MV Maker analysis artifact: [mv_maker_analysis.md](file:///C:/Users/dinht/.gemini/antigravity/brain/0b9f8a1d-c2b4-4847-a89d-11fe82a63bde/artifacts/mv_maker_analysis.md)
