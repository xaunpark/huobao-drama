# Nursery Rhyme Split Mode

> Created: 2026-04-10
> Status: Implemented ✓ (2026-04-10)

## Summary

Add a new `nursery_rhyme` split mode to the storyboard generation system. This mode takes timestamped lyrics input and produces shot lists optimized for children's nursery rhyme music videos, following patterns extracted from real-world analysis of CoComelon/Little Baby Bum and Super Simple Songs productions.

## Problem Statement

The existing 4 split modes (auto, preserve, breakdown, visual_unit) are designed for drama/documentary content. Nursery rhyme music videos have fundamentally different requirements:
- Timing is dictated by lyrics timestamps (not AI-estimated)
- Visuals must literally illustrate lyrics (not cinematic storytelling)
- Structure follows verse/chorus patterns with cumulative repetition
- Target audience is 0-5 years old (bright, safe, phonetic support)

## Prior Solutions

No prior solutions found in `docs/solutions/` for nursery rhymes, lyrics sync, or music video.

## Research Findings

### Codebase Patterns

1. **Split mode routing**: [storyboard_service.go:264-268](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L264-L268) — `visual_unit` mode uses dedicated `processVisualUnitGeneration()`. Nursery rhyme mode will follow the same pattern.

2. **AI output struct → DB mapping**: The `VoiceoverShot` struct ([L74-111](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L74-L111)) is parsed from AI JSON and mapped to `models.Storyboard` via `saveVoiceoverShots()` ([L1214](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L1214)). Nursery rhyme will follow this exact pattern.

3. **Storyboard DB model**: [drama.go:96-155](file:///g:/VS-Project/huobao-drama/domain/models/drama.go#L96-L155) — Already has `ScriptSegment`, `ShotRole`, `AudioMode`, `VisualType` nullable fields from visual_unit mode. New nursery fields will be added as nullable columns.

4. **Prompt template system**: [prompt_template.go:29-62](file:///g:/VS-Project/huobao-drama/domain/models/prompt_template.go#L29-L62) — Custom prompt templates support per-drama overrides via `PromptTemplatePrompts` struct + `PromptTypeToDefaultFile` mapping.

5. **Frontend split mode UI**: [EpisodeWorkflow.vue:804-821](file:///g:/VS-Project/huobao-drama/web/src/views/drama/EpisodeWorkflow.vue#L804-L821) — Radio button group with 4 options. Add 5th option.

6. **Locales**: [en-US.ts:447-454](file:///g:/VS-Project/huobao-drama/web/src/locales/en-US.ts#L447-L454) — Split mode labels and tips.

### Analysis Data (Ground Truth)

Based on analysis of 2 real nursery rhyme videos:
- **Wheels on the Bus** (Little Baby Bum): 43 shots, 2:00, 2.8s avg
- **Old MacDonald** (Super Simple Songs): 44 shots, 3:17, 3.5s avg
- Full analysis files: `docs/Wheels on the bus - Analysis.txt`, `docs/Old-Mc-Donald-Analysis.txt`

### Key Patterns Discovered

1. **Downbeat editing**: All transitions on musical beats
2. **Literal visual mapping**: 80-85% of shots directly illustrate lyrics
3. **Shot recipe**: Establishing → Reveal → Detail → Group Payoff (per verse)
4. **2 structure types**: Narrative (different events per verse) vs Cumulative (additive subjects)
5. **1 line ≠ 1 shot**: AI must freely split/merge within timestamp ranges
6. **2-5s shot duration**: Sweet spot for children's content

## Proposed Solution

### Approach

Follow the `visual_unit` pattern: new AI output struct → dedicated processing function → dedicated save function → reuse existing `Storyboard` DB model with additional nullable columns.

### Architecture Overview

```
INPUT (user pastes):
  [VERSE 1: The Wheels]
  (0:05 – 0:11) [INSTRUMENTAL] Bus driving establishing shot
  (0:12 – 0:15) The wheels on the bus go round and round
  ...

BACKEND PARSE (no AI needed):
  → LyricsBlock[] with timestamp, text, section, is_instrumental

AI CALL (single prompt):
  → NurseryRhymeShot[] (structured JSON)

SAVE TO DB:
  → models.Storyboard records (reuse existing table)

PROFESSIONAL PRODUCTION:
  → Same flow: extract prompts → generate images → generate videos → download
```

---

## Implementation Steps

### Phase 1: Backend Core (Priority: HIGH)

#### Task 1.1: New Go Structs

**File**: `application/services/storyboard_service.go`

```go
// LyricsBlock represents a parsed segment from timestamped lyrics input
type LyricsBlock struct {
    BlockID        int    // Sequential ID
    StartTimeSec   int    // Start time in seconds (parsed from timestamp)
    EndTimeSec     int    // End time in seconds
    DurationSec    int    // Calculated duration
    LyricsText     string // The lyrics text
    SectionType    string // verse / chorus / bridge / intro / outro / instrumental
    SectionNumber  int    // Which verse/chorus (1, 2, 3...)
    VerseSubject   string // Subject from section header ("[VERSE 1: The Wheels]" → "The Wheels")
    IsInstrumental bool   // True if tagged [INSTRUMENTAL]
}

// NurseryRhymeShot is the AI output struct for nursery_rhyme mode
type NurseryRhymeShot struct {
    ShotID            int      `json:"shot_id"`
    LyricsBlockID     int      `json:"lyrics_block_id"`     // Reference to source lyrics block
    LyricsText        string   `json:"lyrics_text"`         // Exact lyrics for this shot
    TimestampStart    string   `json:"timestamp_start"`
    TimestampEnd      string   `json:"timestamp_end"`
    DurationSec       int      `json:"duration_sec"`
    SectionType       string   `json:"section_type"`        // verse/chorus/bridge/intro/outro/instrumental
    SectionNumber     int      `json:"section_number"`
    VerseSubject      string   `json:"verse_subject"`       // "The Wheels", "The Pig"
    ShotRole          string   `json:"shot_role"`           // establishing/reveal/detail/group_payoff
    IsCallback        bool     `json:"is_callback"`         // Is this repeating a previous shot?
    CallbackToShotID  *int     `json:"callback_to_shot_id"` // Reference to original shot
    VisualDescription string   `json:"visual_description"`  // Detailed visual for image generation
    Title             string   `json:"title"`
    ShotType          string   `json:"shot_type"`           // ELS/LS/MS/CU/ECU
    CameraMovement    string   `json:"camera_movement"`
    AnimationHint     string   `json:"animation_hint"`      // "lip_sync, bounce, speech_bubble_pop"
    OverlayText       string   `json:"overlay_text"`        // On-screen text ("oink", "E-I-E-I-O")
    BgmPrompt         string   `json:"bgm_prompt"`
    SoundEffect       string   `json:"sound_effect"`
    TransitionIn      string   `json:"transition_in"`       // hard_cut / crossfade
    Characters        []uint   `json:"characters"`
    Props             []uint   `json:"props"`
    SceneID           *uint    `json:"scene_id"`
    Location          string   `json:"location"`
    Atmosphere        string   `json:"atmosphere"`
}
```

#### Task 1.2: Lyrics Parser

**File**: `application/services/storyboard_service.go`

```go
// parseLyricsInput parses timestamped lyrics format into LyricsBlock array
// Handles:
//   [VERSE 1: The Wheels]                    → section header
//   (0:05 – 0:11) [INSTRUMENTAL] text        → lyrics line with timestamp
//   (0:12 – 0:15) The wheels on the bus...   → regular lyrics line
func parseLyricsInput(script string) ([]LyricsBlock, string, int) {
    // Returns: blocks, detected section summary, total blocks count
    // Section header pattern: [VERSE N: Subject], [CHORUS], [BRIDGE], [INTRO], [OUTRO]
    // Timestamp pattern: (M:SS – M:SS) or (MM:SS – MM:SS) or (M:SS - M:SS)
    // Instrumental tag: [INSTRUMENTAL]
}

// detectNurseryStructure detects Narrative vs Cumulative structure type
func detectNurseryStructure(blocks []LyricsBlock) (structureType string, reason string) {
    // Analyze lyrics patterns:
    // - If later verses contain words from earlier verses → "cumulative"
    // - If each verse is entirely different content → "narrative"
    // Returns: ("narrative"|"cumulative", human-readable reason)
}

// detectNurseryRhymeInput checks if script content matches nursery_rhyme format
// Used by auto-detect mode
func detectNurseryRhymeInput(script string) bool {
    // Check for: section headers [VERSE...] + timestamp patterns (M:SS – M:SS)
    // Require >= 3 timestamp lines AND >= 1 section header
}
```

#### Task 1.3: Processing Function

**File**: `application/services/storyboard_service.go`

Add new function following existing `processVisualUnitGeneration()` pattern:

```go
func (s *StoryboardService) processNurseryRhymeGeneration(
    taskID, episodeID string, dramaID uint, 
    model, scriptContent, characterList, sceneList, propList string,
) {
    // 1. Parse lyrics input → LyricsBlock[]
    // 2. Detect structure type (narrative/cumulative)
    // 3. Build analysis context string from parsed blocks
    // 4. Load system prompt (storyboard_nursery_rhyme.txt)
    // 5. Build full prompt with lyrics + analysis + characters + scenes + format instructions
    // 6. Call AI
    // 7. Parse JSON → NurseryRhymeShot[]
    // 8. Validate: total duration must match lyrics range
    // 9. Save via saveNurseryRhymeShots()
}
```

#### Task 1.4: Routing

**File**: `application/services/storyboard_service.go` — Line ~264

Add routing case before the preserve/breakdown split:

```go
// Route to nursery_rhyme mode if selected
if splitMode == "nursery_rhyme" {
    s.log.Infow("Using NURSERY_RHYME mode — lyrics-synced shot planning", "task_id", taskID)
    s.processNurseryRhymeGeneration(taskID, episodeID, dramaID, model, scriptContent, characterList, sceneList, propList)
    return
}
```

Also update `detectTimestampPattern()` / auto-detect logic to optionally detect nursery format.

#### Task 1.5: AI Prompt Files

**New files**:
- `application/prompts/storyboard_nursery_rhyme.txt` — System prompt
- `application/prompts/storyboard_nursery_rhyme_format.txt` — Output format spec

System prompt key sections:
- Role: Children's animation storyboard director
- 7 universal rules (from analysis synthesis)
- Structure type handling (narrative vs cumulative)
- 4-step shot recipe per verse
- Duration rules (2-5s, strict to timestamps)
- Literal visual mapping philosophy
- On-screen text/phonetic support rules

#### Task 1.6: New DB Fields (Migration)

**File**: `domain/models/drama.go` — Storyboard struct

Add nullable columns (no migration script needed — GORM AutoMigrate handles this):

```go
// Nursery Rhyme fields
LyricsText      *string `gorm:"type:text" json:"lyrics_text"`              // Lyrics text for this shot
SectionType     *string `gorm:"size:20" json:"section_type"`               // verse/chorus/bridge/intro/outro/instrumental
VerseSubject    *string `gorm:"size:100" json:"verse_subject"`             // "The Wheels", "The Pig"
ShotRoleNursery *string `gorm:"size:30" json:"shot_role_nursery"`          // establishing/reveal/detail/group_payoff
OverlayText     *string `gorm:"size:200" json:"overlay_text"`              // On-screen text ("oink", "E-I-E-I-O")
AnimationHint   *string `gorm:"size:200" json:"animation_hint"`            // "lip_sync, bounce, speech_bubble_pop"
IsCallback      *bool   `json:"is_callback"`                               // Is this a callback/repeat shot?
CallbackShotNum *int    `json:"callback_shot_num"`                         // Reference to original shot number
```

#### Task 1.7: Save Function

**File**: `application/services/storyboard_service.go`

```go
// saveNurseryRhymeShots maps NurseryRhymeShot[] to models.Storyboard and saves to DB
// Follows same pattern as saveVoiceoverShots()
func (s *StoryboardService) saveNurseryRhymeShots(episodeID string, dramaID uint, shots []NurseryRhymeShot) error {
    // Delete existing storyboards for episode
    // Map NurseryRhymeShot fields to models.Storyboard
    // Set ScriptSegment = LyricsText (for display in Professional Production)
    // Set AudioMode = "lyrics_sync"
    // Save with character/prop associations
}
```

#### Task 1.8: Prompt Template Registration

**File**: `domain/models/prompt_template.go`

```go
// Add to PromptTemplatePrompts struct:
NurseryRhymeBreakdown string `json:"nursery_rhyme_breakdown,omitempty"`

// Add to PromptTypeToDefaultFile map:
"nursery_rhyme_breakdown": "storyboard_nursery_rhyme.txt",
```

---

### Phase 2: Frontend (Priority: MEDIUM)

#### Task 2.1: Add Radio Button

**File**: `web/src/views/drama/EpisodeWorkflow.vue`

Add `nursery_rhyme` option to the `el-radio-group` at line ~817, and to the `el-dropdown-menu` at line ~965.

```html
<el-radio-button value="nursery_rhyme">
  <el-icon style="margin-right: 4px;"><Mic /></el-icon>
  {{ $t('workflow.splitModeNurseryRhyme') }}
</el-radio-button>
```

#### Task 2.2: Locale Strings

**Files**: `web/src/locales/en-US.ts`, `web/src/locales/vi-VN.ts` (if exists)

```typescript
splitModeNurseryRhyme: 'Nursery Rhyme',
splitModeNurseryRhymeTip: 'Nursery Rhyme Mode: Designed for children\'s music videos. Input timestamped lyrics with [VERSE], [CHORUS] markers. AI creates child-friendly, literal visual illustrations synced to lyrics timing. Supports both Narrative and Cumulative song structures.',
```

#### Task 2.3: Split Mode Tip Display

**File**: `web/src/views/drama/EpisodeWorkflow.vue` — Line ~825-831

Add nursery_rhyme case to the tip text conditional:

```javascript
shotSplitMode === 'nursery_rhyme' ? $t('workflow.splitModeNurseryRhymeTip') : ...
```

#### Task 2.4: Storyboard Table — Lyrics Display

**File**: `web/src/views/drama/EpisodeWorkflow.vue` — Shot list table

Add conditional columns when nursery_rhyme data is detected:

```html
<!-- Nursery: Lyrics Text column -->
<el-table-column
  v-if="currentEpisode?.storyboards?.some(s => s.lyrics_text)"
  label="🎵 Lyrics" min-width="200" show-overflow-tooltip>
  <template #default="{ row }">
    <span style="font-style: italic; color: var(--el-color-primary-light-3);">
      {{ row.lyrics_text || "-" }}
    </span>
    <el-tag v-if="row.section_type" size="small" type="success" style="margin-left: 4px;">
      {{ row.section_type }}{{ row.verse_subject ? `: ${row.verse_subject}` : '' }}
    </el-tag>
  </template>
</el-table-column>
```

---

### Phase 3: Integration (Priority: MEDIUM)

#### Task 3.1: Frontend API Pass-through

**File**: `web/src/api/generation.ts`

Already passes `split_mode` parameter — no changes needed. The `nursery_rhyme` value will be sent as-is.

#### Task 3.2: Backend Handler

**File**: `api/handlers/storyboard.go`

The handler already accepts any string for `split_mode` — no changes needed. Value flows through to `storyboard_service.GenerateStoryboard()`.

#### Task 3.3: Professional Production Compatibility

The nursery_rhyme shots will be saved as regular `Storyboard` records. The Professional Production UI reads from the same table, so:
- ✅ Image prompt extraction — works (uses `action`, `result`, `visual_description` fields)
- ✅ Image generation — works (same flow)
- ✅ Video generation — works (same flow)
- ✅ Download — works (same flow)

**Lyrics text visibility**: Set `ScriptSegment = LyricsText` during save so it appears in the existing "📝 Narrator" column in Professional Production.

#### Task 3.4: Storyboard Composition Service

**File**: `application/services/storyboard_composition_service.go`

The image prompt extraction uses `Action`, `Result`, `Atmosphere`, `Description` from the storyboard. For nursery_rhyme shots:
- `Action` = NurseryRhymeShot.VisualDescription (primary visual)
- `Result` = NurseryRhymeShot.OverlayText context
- `Atmosphere` = NurseryRhymeShot.Atmosphere
- `Description` = Built from shot metadata

No changes needed to the composition service — it reads generic storyboard fields.

---

## Acceptance Criteria

- [ ] User can select "Nursery Rhyme" mode in the split mode radio group
- [ ] Backend correctly parses timestamped lyrics with `[VERSE N: Subject]` headers
- [ ] Backend correctly parses `[INSTRUMENTAL]` tagged lines
- [ ] Backend correctly detects Narrative vs Cumulative structure
- [ ] AI generates shots that respect lyrics timestamps (total duration matches)
- [ ] Each shot has clear `lyrics_text` showing exactly which lyrics it illustrates
- [ ] Shots are saved as regular Storyboard records and appear in Professional Production
- [ ] `ScriptSegment` field shows lyrics text in the storyboard table
- [ ] Section labels (`section_type`, `verse_subject`) are visible in the shot list
- [ ] On-screen text (overlay_text) is generated for onomatopoeia and key refrains
- [ ] Shot role follows the 4-step recipe (establishing/reveal/detail/group_payoff)
- [ ] Professional Production flow works: extract prompts → generate images → generate videos → download

## Technical Considerations

### Dependencies
- No new Go packages needed
- No new npm packages needed
- Uses existing GORM AutoMigrate for new columns

### Risks
1. **AI prompt quality**: The prompt needs careful tuning with real nursery rhyme inputs. Mitigation: test with the 2 analyzed songs first
2. **Duration accuracy**: AI may not perfectly match timestamp durations. Mitigation: post-process to force-correct durations from parsed timestamps
3. **DB migration**: New nullable columns on `storyboards` table. Risk: minimal (nullable, no data loss). GORM handles automatically.

### Alternatives Considered
1. **Reuse `preserve` mode** — Rejected: doesn't understand lyrics → visual mapping, no child-friendly prompting
2. **Reuse `visual_unit` mode** — Rejected: 3-5s duration constraint, would break timestamps, wrong audio model
3. **Create separate table for nursery shots** — Rejected: would require duplicating entire Professional Production pipeline. Reusing `Storyboard` model is more pragmatic.

### Future Work (Out of Scope)
- Music file upload + waveform-based beat detection
- Audio/music merge in final video export
- Template gallery for common nursery rhymes
- Verse template selection (user picks verse structure from presets)

## References

- Analysis data: [Wheels on the Bus](file:///g:/VS-Project/huobao-drama/docs/Wheels%20on%20the%20bus%20-%20Analysis.txt)
- Analysis data: [Old MacDonald](file:///g:/VS-Project/huobao-drama/docs/Old-Mc-Donald-Analysis.txt)
- Synthesis artifact: `artifacts/nursery_rhymes_synthesis.md`
- Existing visual_unit implementation: [storyboard_service.go:876-1212](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L876-L1212)
