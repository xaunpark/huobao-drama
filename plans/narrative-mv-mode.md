# Narrative MV Mode — Cinematic Short Film + Music

> Created: 2026-04-21
> Status: Implemented
> Reviewed: 2026-04-21 (plan_review concerns addressed)

## Summary

Implement `narrative_mv` — a new split mode that generates storyboards for **cinematic short films combined with music**. Unlike existing modes (nursery_rhyme, mv_maker) where visuals illustrate lyrics line-by-line, this mode treats **music as emotional atmosphere** and **visuals as independent storytelling**. Output is a 3-part unified film:

- **Part 1 — Prologue** (60–90s): Pure film, no music. Establishes world, characters, inciting event.
- **Part 2 — Music Film** (full song duration): Story continues while music plays. Music and visuals run in parallel, converging at key dramatic moments. Shot timing anchored to lyric timestamps for post-production convenience.
- **Part 3 — Epilogue** (≈60s): Pure film, no music. Resolution / aftermath.

---

## Problem Statement

Current split modes force a 1:1 relationship between lyrics/text blocks and shots ("Karaoke illustration"). This produces technically correct storyboards but cinematically shallow results. The user needs a mode that generates **film-quality shot lists** where:
- Characters **act** (express, react, move through space) rather than sing at camera
- **Visual Irony** is possible (happy music + dark visuals)
- **Narrative continuity** spans all 3 parts via shared motifs and character arcs
- **Lyric timestamps** serve as time anchors for post-production, not content drivers

---

## Research Findings

### Existing Routing Pattern
`application/services/storyboard_service.go:267–291` — `processStoryboardGeneration()` uses a simple `if splitMode == "..."` chain. Adding `narrative_mv` follows the exact same pattern.

### Existing DB Model
`domain/models/drama.go:96–162` — `Storyboard` already has extensible fields from previous modes (nursery rhyme added 7 fields, voice-over added 8 fields). New fields appended to the same struct per established pattern.

### Existing 2-File Prompt Architecture
`storyboard_mv_service.go:69` — MV Maker uses `systemPrompt` (genre-specific) + `formatInstructions` (shared format file). Narrative MV will use a 3-file architecture: `planner` + `director` + `format`.

### Established Parser Pattern
`storyboard_nursery_service.go:102–191` — `parseLyricsInput()` uses regex marker matching on a single textarea. The Story Bible parser follows the same design using `[MARKER]` sections.

### Prompt Template Registration
`domain/models/prompt_template.go:45,65` — New genre fields appended to `PromptTemplatePrompts` struct and `PromptTypeToDefaultFile` map.

---

## DB Schema Decision

**Decision: Add new columns via GORM AutoMigrate** (same as every previous feature addition in this project).

New fields appended to `models.Storyboard`:

```go
// Narrative MV fields
NarrativePart    *string `gorm:"size:20" json:"narrative_part"`     // "prologue" | "music_film" | "epilogue"
HasMusic         *bool   `json:"has_music"`                         // false for parts 1 & 3
MusicSegment     *string `gorm:"size:50" json:"music_segment"`      // "VERSE 1", "CHORUS", null for no-music parts
MusicSyncType    *string `gorm:"size:20" json:"music_sync_type"`    // "parallel" | "convergent" | "irony" | null
ActingNote       *string `gorm:"type:text" json:"acting_note"`      // Director's acting instruction for AI video gen
LyricsAnchor     *string `gorm:"type:text" json:"lyrics_anchor"`    // Lyric line this shot is time-anchored to
```

**Why these 6 fields are columns (not JSON):**
- `NarrativePart` → queried for GROUP BY grouping in the Professional Editor list
- `HasMusic` → filtered when rendering color bands in the shot list
- The others are display fields; added as columns for consistency with existing field pattern

GORM AutoMigrate is safe — it only ADDs columns, never drops or modifies existing ones. Zero risk to existing data.

---

## Input Format: Story Bible (Single Textarea)

Follows the same `[MARKER]` convention as `parseLyricsInput()`. Users type this into the existing script textarea when `narrative_mv` is selected:

```
[STORY_BIBLE]
Setting: Hanoi, 1985. A ballet rehearsal room at night.
Emotional core: Not about death. About loving someone without letting them know you're saying goodbye.

[CHARACTERS]
An | Female, 25, ballet dancer | Role: Protagonist, hiding terminal illness
Minh | Male, 28, orchestra conductor | Role: An's lover, witnesses from the shadows

[PROLOGUE] duration: 75s
An enters the empty rehearsal room alone. She moves slowly, touching the barre.
She puts on her ballet shoes with deliberate care — as if memorizing the feeling.
An stands center stage under a single spotlight. The room is silent.

[MUSIC_SEGMENTS]
(0:00 - 0:45) INTRO — piano only — emotion: solitude, false calm
(0:45 - 1:30) VERSE 1 — emotion: nostalgia, talking to herself
(1:30 - 2:00) PRE-CHORUS — emotion: suppressed grief, about to break
(2:00 - 2:30) CHORUS — emotion: release, pain erupts
  [SYNC_POINT] An begins dancing for the first time — convergent
(2:30 - 3:00) VERSE 2 — emotion: resignation
  [SYNC_POINT] Minh appears at the doorway, watching unseen — irony
(3:00 - 3:30) CHORUS 2 — emotion: emotional peak
(3:30 - 4:00) OUTRO — piano fade — emotion: acceptance

[LYRICS]
(0:45 - 0:52) I'm still here though there's no one left
(0:52 - 1:05) This dusty room with its final steps

[EPILOGUE] duration: 60s
Next morning. Empty room. Only An's ballet shoes remain on the floor.
Minh enters, picks up a shoe, holds it briefly.
He looks out the window at the rain. Slow fade to white.
```

---

## Proposed Solution

### Story Bible Struct (Parser Output)

```go
type StoryBible struct {
    WorldDescription string           // Content of [STORY_BIBLE] section
    EmotionalCore    string           // Extracted from "Emotional core:" label if present
    Characters       []NarrativeCharacter
    PrologueDesc     string           // Full text of [PROLOGUE] block (free text, no sub-parsing)
    PrologueDuration int              // Parsed from "duration: Ns" on [PROLOGUE] header line
    MusicSegments    []MusicSegment
    LyricsBlocks     []LyricsBlock    // Reused from storyboard_nursery_service.go
    EpilogueDesc     string           // Full text of [EPILOGUE] block
    EpilogueDuration int              // Parsed from "duration: Ns" on [EPILOGUE] header
    TotalMusicSec    int              // Computed: sum of all MusicSegment durations
}

type NarrativeCharacter struct {
    Name        string // Pipe-delimited field 1
    Description string // Pipe-delimited field 2
    Role        string // Pipe-delimited field 3
}

type MusicSegment struct {
    StartSec   int
    EndSec     int
    Name       string      // "INTRO", "VERSE 1", "CHORUS"...
    Emotion    string      // text after "emotion:"
    SyncPoints []SyncPoint // indented [SYNC_POINT] lines belonging to this segment
}

type SyncPoint struct {
    Description string // text before " — "
    SyncType    string // "convergent" | "irony" | "parallel"
}
```

**Parser Notes:**
- `[PROLOGUE]` and `[EPILOGUE]` blocks are treated as **free text** — no sub-parsing. The entire block content becomes `PrologueDesc`/`EpilogueDesc`. Users may write naturally; no special tags like `END:` are required or parsed.
- `[SYNC_POINT]` detection regex:
  ```go
  syncPointPattern = regexp.MustCompile(`^\s{2,}\[SYNC_POINT\]\s+(.+?)\s+—\s+(convergent|irony|parallel)$`)
  ```
  Must be indented by ≥2 spaces/1 tab to be recognized as belonging to the preceding `[MUSIC_SEGMENTS]` line.

---

### Phase 1: AI Story Planning (AI Call 1)

**Prompt file:** `application/prompts/storyboard_narrative_planner.txt`

Input: Full Story Bible text
Output: `NarrativePlan` JSON

> **Note (from plan review):** Shot counts (`prologue_shots`, etc.) are removed. They were false precision — Phase 2 AI owns shot count decisions. Only `total_music_duration_sec` is authoritative, used for backend duration validation.

```json
{
  "narrative_thread": [
    {"motif": "ballet shoes", "appears_in": ["prologue", "music_film", "epilogue"], "meaning": "An's identity and mortality"}
  ],
  "lighting_arc": "Warm artificial light (prologue) → dims progressively (music_film) → cold morning light (epilogue)",
  "character_arcs": {
    "An": "Performing strength → Genuine release through dance → Absent (represented by shoes)",
    "Minh": "Absent → Silent witness in shadow → Grief made physical"
  },
  "sync_points": [
    {"music_timestamp": "2:00", "segment": "CHORUS", "sync_type": "convergent", "description": "An begins dancing as chorus erupts"},
    {"music_timestamp": "2:30", "segment": "VERSE 2", "sync_type": "irony", "description": "Music is wistful but Minh is sobbing unseen"}
  ],
  "total_music_duration_sec": 240
}
```

### Phase 2: AI Shot Generation (AI Call 2)

**Prompt file:** `application/prompts/storyboard_narrative_director.txt`

Input: Full Story Bible + NarrativePlan JSON
Output: `[]NarrativeShot`

**Duration rules (enforced in format file):**
- Prologue shots: 4–8s (acting space)
- Music Film shots: 3–6s (music-paced)
- Epilogue shots: 4–10s (deliberate final images)
- Validation: sum(music_film shots) MUST approximately equal `total_music_duration_sec` ± 10%

**AI timestamp behaviour:** Phase 2 AI generates only `duration_sec`. It does NOT generate `timestamp_start`/`timestamp_end` — these are computed by backend post-processing (Step A3.5).

**Format file:** `application/prompts/storyboard_narrative_format.txt`

---

## New Output Struct: `NarrativeShot`

```go
type NarrativeShot struct {
    ShotID            int    `json:"shot_id"`
    NarrativePart     string `json:"narrative_part"`       // "prologue" | "music_film" | "epilogue"
    HasMusic          bool   `json:"has_music"`
    TimestampStart    string `json:"timestamp_start"`      // absolute video time (e.g. "1:15")
    TimestampEnd      string `json:"timestamp_end"`
    DurationSec       int    `json:"duration_sec"`
    MusicSegment      string `json:"music_segment"`        // "VERSE 1", "CHORUS", "" for no-music parts
    MusicSyncType     string `json:"music_sync_type"`      // "parallel" | "convergent" | "irony" | ""
    LyricsAnchor      string `json:"lyrics_anchor"`        // lyric line this shot is anchored to (post-prod reference)
    VisualDescription string `json:"visual_description"`   // template-agnostic scene description
    ActingNote        string `json:"acting_note"`          // micro-expression / body language direction
    ShotType          string `json:"shot_type"`
    CameraAngle       string `json:"camera_angle"`
    CameraMovement    string `json:"camera_movement"`
    NarrativeFunction string `json:"narrative_function"`   // setup_character | plot_reveal | emotional_peak | resolution
    Subtext           string `json:"subtext"`              // what character hides (enables visual irony)
    Location          string `json:"location"`
    Atmosphere        string `json:"atmosphere"`
    Characters        []uint `json:"characters"`
    SceneID           *uint  `json:"scene_id"`
    Title             string `json:"title"`
}
```

---

## Implementation Phases

### Phase A — Backend Core (Day 1)

**A1. DB Schema** — `domain/models/drama.go`
- Append 6 new `NarrativePart`, `HasMusic`, `MusicSegment`, `MusicSyncType`, `ActingNote`, `LyricsAnchor` fields to `Storyboard` struct

**A2. Story Bible Parser + Pipeline** — `application/services/storyboard_narrative_service.go` (NEW FILE)
- `parseStoryBible(script string) (*StoryBible, error)` — regex-based section splitter
  - Returns `*StoryBible` struct (see struct definitions in Proposed Solution section)
  - `[STORY_BIBLE]`: captured as `WorldDescription` free text
  - `[CHARACTERS]`: `Name | Description | Role` pipe-delimited, each line = one character
  - `[PROLOGUE]` / `[EPILOGUE]`: header line parsed for `duration: Ns`; body = free text (no sub-parsing)
  - `[MUSIC_SEGMENTS]`: `(ts - ts) SEGMENT_NAME — emotion: ...` with indented `[SYNC_POINT]` lines
    - SyncPoint regex: `^\s{2,}\[SYNC_POINT\]\s+(.+?)\s+—\s+(convergent|irony|parallel)$`
  - `[LYRICS]`: delegates to existing `parseLyricsInput()` — full reuse
- `processNarrativeMVGeneration(...)` — 2-phase AI pipeline
  - Phase 1: planner prompt → parse `NarrativePlan` JSON
  - Phase 2: director prompt + NarrativePlan context → parse `[]NarrativeShot`
  - Post-validate: Part 2 duration sum within ±10% of `StoryBible.TotalMusicSec`
- `saveNarrativeShots()` — maps `NarrativeShot` → `models.Storyboard` (same pattern as `saveNurseryRhymeShots()`)

**A3. Routing** — `application/services/storyboard_service.go` (line ~268)
```go
if splitMode == "narrative_mv" {
    s.log.Infow("Using NARRATIVE_MV mode", "task_id", taskID)
    s.processNarrativeMVGeneration(...)
    return
}
```

**A3.5. Backend Timestamp Post-Processing** *(critical step — not AI responsibility)*

After Phase 2 AI returns `[]NarrativeShot`, backend assigns absolute video timestamps by walking shots sequentially:
```go
cumulativeSec := 0
for i := range shots {
    shots[i].TimestampStart = formatSeconds(cumulativeSec)
    cumulativeSec += shots[i].DurationSec
    shots[i].TimestampEnd = formatSeconds(cumulativeSec)
}
```
This is the same pattern used in `storyboard_nursery_service.go` for duration correction. AI does NOT output timestamps — only `duration_sec`.

**A4. Prompt Template Registration** — `domain/models/prompt_template.go`
```go
// PromptTemplatePrompts struct:
NarrativeMVPlanner  string `json:"narrative_mv_planner,omitempty"`
NarrativeMVDirector string `json:"narrative_mv_director,omitempty"`

// PromptTypeToDefaultFile map:
"narrative_mv_planner":  "storyboard_narrative_planner.txt",
"narrative_mv_director": "storyboard_narrative_director.txt",
```

**A5. Prompt I18n Wiring** — `application/services/prompt_i18n.go`
- Add `WithDramaNarrativeMVPlannerPrompt(dramaID uint) string` method
- Add `WithDramaNarrativeMVDirectorPrompt(dramaID uint) string` method
- Same pattern as existing `WithDramaMVMakerSystemPrompt()`

---

### Phase B — Prompt Files (Day 1–2)

| File | Purpose |
|---|---|
| `storyboard_narrative_planner.txt` | Role: Story analyst. Extract narrative threads, lighting arc, sync points. Output: NarrativePlan JSON. MUST be template-agnostic. |
| `storyboard_narrative_director.txt` | Role: Film director. Receives Story Bible + NarrativePlan. Generates NarrativeShot[]. Characters act naturally, NO 4th wall. Visual content driven by story arc, NOT lyric meaning. |
| `storyboard_narrative_format.txt` | JSON schema + field definitions + duration validation rules |

---

### Phase C — Frontend (Day 2–3)

**C1. Dropdown** — `EpisodeWorkflow.vue`
- Add `narrative_mv` to split mode dropdown (same pattern as `mv_maker`)

**C2. Locale** — `en-US.ts` / `zh-CN.ts`
- `splitModeNarrativeMV`: `'🎬 Narrative MV (Cinematic short film + music)'`
- `splitModeNarrativeMVTip`: full input format description

**C3. Professional Editor — 3-Part Grouping**

*Condition: activated only when `narrative_part != null` on ANY shot in episode*

Shot list left panel:
- Group Header separators: `🎬 PROLOGUE · Shot 1–8 · 75s`
- Left border per shot: Amber `#E67E22` / Violet `#8E44AD` / Steel Blue `#2980B9`

Shot detail panel:
- Part badge: `🎵 MUSIC FILM · CHORUS · 🎯 convergent`
- `acting_note` field: italic, below atmosphere
- `lyrics_anchor` field: dimmed, bottom of card (post-prod reference)
- All new fields conditional on `!= null` — zero impact on other modes

---

## Acceptance Criteria

- [ ] Story Bible parser handles all 6 sections correctly
- [ ] `[CHARACTERS]` pipe-delimited format (`Name | Description | Role`) parsed into `NarrativeCharacter[]`
- [ ] `[PROLOGUE]` / `[EPILOGUE]` `duration: Ns` parsed; body treated as free text (no `END:` or special tags required)
- [ ] `[MUSIC_SEGMENTS]` with indented `[SYNC_POINT]` lines parsed using defined regex
- [ ] `[SYNC_POINT]` with missing `sync_type` suffix handled gracefully (default to `"parallel"`)
- [ ] `[LYRICS]` delegated to existing `parseLyricsInput()` without modification
- [ ] Phase 1 AI call produces valid `NarrativePlan` JSON (no shot count fields)
- [ ] Phase 2 AI call produces `NarrativeShot[]` with `duration_sec` only (no timestamp fields from AI)
- [ ] Backend post-processing (A3.5) correctly assigns cumulative `timestamp_start`/`timestamp_end`
- [ ] Sum of Part 2 shot durations within ±10% of `StoryBible.TotalMusicSec`; warning logged if outside
- [ ] All 6 new DB columns present after AutoMigrate
- [ ] `saveNarrativeShots()` correctly populates all 6 new fields
- [ ] `WithDramaNarrativeMVPlannerPrompt()` and `WithDramaNarrativeMVDirectorPrompt()` wired in `prompt_i18n.go`
- [ ] `narrative_mv` option visible in EpisodeWorkflow dropdown
- [ ] Professional Editor shows 3-part group headers when `narrative_part` detected
- [ ] Left border color: Amber (prologue) / Violet (music_film) / Steel Blue (epilogue)
- [ ] `acting_note` and `lyrics_anchor` visible in shot detail panel (when not null)
- [ ] `docs/features/narrative_mv_input_guide.md` created with worked example
- [ ] **All existing split modes (auto, preserve, breakdown, visual_unit, nursery_rhyme, mv_maker, cinematic_movie) completely unaffected**

---

## Technical Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| 2-Phase AI latency | Medium | Phase 1 is compact; total latency comparable to single large prompt |
| Music Film duration mismatch | Medium | ±10% tolerance + backend warning log (not hard error) |
| Template coupling in prompts | Low | Strict review: zero visual style references in prompt files |
| Story Bible format too complex for users | Medium | Input guide doc + tooltip with worked example |
| `[SYNC_POINT]` indentation varies by editor | Low | Regex uses `\s{2,}` (≥2 spaces OR tab) — tolerant of common editors |
| AI omits `duration_sec` for some shots | Low | Backend fallback: if `duration_sec <= 0`, assign based on `NarrativePart` (4s prologue, 4s music_film, 5s epilogue) |

---

## Files Changed

| File | Action |
|---|---|
| `domain/models/drama.go` | ADD 6 fields to `Storyboard` |
| `application/services/storyboard_service.go` | ADD 1 routing case |
| `application/services/storyboard_narrative_service.go` | CREATE — parser, pipeline, save |
| `application/services/prompt_i18n.go` | ADD `WithDramaNarrativeMVPlannerPrompt()` + `WithDramaNarrativeMVDirectorPrompt()` |
| `application/prompts/storyboard_narrative_planner.txt` | CREATE |
| `application/prompts/storyboard_narrative_director.txt` | CREATE |
| `application/prompts/storyboard_narrative_format.txt` | CREATE |
| `domain/models/prompt_template.go` | ADD 2 struct fields + 2 map entries |
| `web/src/views/drama/EpisodeWorkflow.vue` | ADD dropdown item + conditional 3-part UI in editor |
| `web/src/locales/en-US.ts` | ADD 2 locale strings |
| `web/src/locales/zh-CN.ts` | ADD 2 locale strings |
| `docs/features/narrative_mv_input_guide.md` | CREATE — user guide with worked example |

---

## References

- Routing pattern: `storyboard_service.go:267–291`
- DB model: `domain/models/drama.go:96–162`
- Parser to reuse: `storyboard_nursery_service.go:102`
- MV Maker prompt pattern: `storyboard_mv_service.go`
- Template registration: `domain/models/prompt_template.go`
