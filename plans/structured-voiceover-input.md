# Structured Voice-Over Input (Shot-Delimited Mode)

> Created: 2026-04-05
> Status: ✅ Complete
> Estimated Effort: ~4-6 hours

## Summary

Add a **structured input sub-mode** to the existing Voice-Over (`visual_unit`) pipeline. When the user's script contains `// SHOT XX` markers, the system preserves those shot boundaries and only asks AI to **enrich** each shot with visual/audio metadata — never re-splitting. When no markers are found, existing free-form behavior is preserved (backward compatible).

## Problem Statement

The current Voice-Over mode conflates two responsibilities:
1. **Splitting** — deciding WHERE to cut the script into shots
2. **Enriching** — adding visual_description, camera, audio_mode, etc.

The AI frequently splits incorrectly:
- Each line becomes its own shot (too granular)
- Section headers and `[SFX]` become standalone shots
- Narrative groupings that belong together get fragmented
- User's intended shot logic is ignored

Users (who build scripts with external AI tools) already know exactly which lines belong in each shot. They need the system to **respect their structure** and only enrich with metadata.

## Research Findings

### Codebase Patterns

1. **Preserve mode precedent** — [storyboard_service.go:183-243](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L183-L243)
   - `detectTimestampPattern()` already auto-detects `(0:00 – 0:06)` patterns and routes to preserve mode
   - Uses separate prompt file `storyboard_preserve_shots.txt`
   - Key instruction: "Preserve each shot/block... Do NOT merge or skip any shots. Enrich each with cinematography metadata."
   - **This is the exact same concept** — we need to replicate it for the `visual_unit` pipeline

2. **parseMarkedScript()** — [storyboard_service.go:434-556](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L434-L556)
   - Already handles `[Character]`, `[CROWD]`, `[SFX]`, `[NARRATOR]` tags
   - Already builds `ScriptSegmentInfo` with type/character/text per segment
   - **Problem:** Currently treats each segment as a flat list — no concept of shot grouping

3. **isMarkdownMetadata()** — [storyboard_service.go:558-622](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L558-L622)
   - Already skips markdown headers (`#`), tables (`|`), horizontal rules (`---`)
   - **Currently skips `//` lines as metadata** (line 601) — this is the key conflict!
   - `// SHOT XX` is currently treated as metadata and **discarded**
   - We need to **carve out** `// SHOT` markers from the metadata filter

4. **processVisualUnitGeneration()** — [storyboard_service.go:624-794](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L624-L794)
   - Routes `visual_unit` mode to separate prompt + parse flow
   - Uses `VoiceoverShot` struct for AI output (richer than `Storyboard`)
   - Saves via `saveVoiceoverShots()` — maps VoiceoverShot → models.Storyboard

5. **Frontend split mode UI** — [EpisodeWorkflow.vue:803-831](file:///g:/VS-Project/huobao-drama/web/src/views/drama/EpisodeWorkflow.vue#L803-L831)
   - 4 radio buttons: Auto / Preserve / Breakdown / Visual Unit
   - Tooltip descriptions explain each mode
   - `visual_unit` has orange warning-style tooltip

### Key Insight: Minimal Surface Area

The structured mode doesn't need a new split mode button. It works **within** the existing `visual_unit` mode via auto-detection (identical to how `auto` detects timestamps → `preserve`). The flow is:

```
User selects "Visual Unit" mode
     │
     ▼
processVisualUnitGeneration() called
     │
     ▼
NEW: detectStructuredShots(script) — looks for "// SHOT" markers
     │
     ├── FOUND → Structured sub-mode
     │      ├── parseStructuredShots() → ShotBlock[]
     │      ├── Use ENRICH-ONLY prompt (no splitting)
     │      └── AI output: same VoiceoverShot[] schema
     │
     └── NOT FOUND → Existing free-form behavior (unchanged)
```

## Proposed Solution

### Phase 1: Backend Parser — `parseStructuredShots()`

**File:** `application/services/storyboard_service.go`

New function that parses `// SHOT XX` markers:

```go
// ShotBlock represents a user-defined shot with pre-grouped content
type ShotBlock struct {
    ShotNumber   int
    Duration     int      // 0 if not specified
    ShotType     string   // "" if not specified
    AudioMode    string   // "" if not specified
    Lines        []ScriptSegmentInfo // All segments within this shot
    RawContent   string   // Original text between markers (for script_segment)
}

// detectStructuredShots checks if script contains "// SHOT" markers
func detectStructuredShots(script string) bool {
    re := regexp.MustCompile(`(?m)^//\s*SHOT\s+\d+`)
    matches := re.FindAllString(script, -1)
    return len(matches) >= 3
}

// parseStructuredShots splits script into ShotBlocks by "// SHOT" markers
func parseStructuredShots(script string) []ShotBlock {
    headerRe := regexp.MustCompile(
        `^//\s*SHOT\s+(\d+)` +
        `(?:\s*\|\s*(\d+)s)?` +
        `(?:\s*\|\s*([\w-]+))?` +
        `(?:\s*\|\s*(narrator_only|dialogue_dominant))?`,
    )
    
    lines := strings.Split(script, "\n")
    var blocks []ShotBlock
    var current *ShotBlock
    var contentLines []string
    
    for _, line := range lines {
        trimmed := strings.TrimSpace(line)
        
        if matches := headerRe.FindStringSubmatch(trimmed); matches != nil {
            if current != nil {
                current.RawContent = strings.TrimSpace(strings.Join(contentLines, "\n"))
                current.Lines = parseSegmentsInBlock(contentLines)
                blocks = append(blocks, *current)
            }
            
            shotNum, _ := strconv.Atoi(matches[1])
            dur := 0
            if matches[2] != "" { dur, _ = strconv.Atoi(matches[2]) }
            
            current = &ShotBlock{
                ShotNumber: shotNum,
                Duration:   dur,
                ShotType:   matches[3],
                AudioMode:  matches[4],
            }
            contentLines = nil
            continue
        }
        
        if current != nil {
            if trimmed != "" {
                contentLines = append(contentLines, line)
            }
        }
    }
    
    if current != nil {
        current.RawContent = strings.TrimSpace(strings.Join(contentLines, "\n"))
        current.Lines = parseSegmentsInBlock(contentLines)
        blocks = append(blocks, *current)
    }
    
    return blocks
}
```

**Critical fix in `isMarkdownMetadata()`:** Currently line 601 marks ALL `//` lines as metadata. We need to carve out `// SHOT`:

```go
// Before:
if strings.HasPrefix(line, "//") {
    return true
}

// After:
if strings.HasPrefix(line, "//") {
    shotMarker := regexp.MustCompile(`^//\s*SHOT\s+\d+`)
    if !shotMarker.MatchString(line) {
        return true
    }
    return false
}
```

### Phase 2: New Prompt — `storyboard_visual_unit_structured.txt`

**File:** `application/prompts/storyboard_visual_unit_structured.txt`

This is the ENRICH-ONLY prompt — AI receives pre-split shots and only adds visual metadata:

```
[Role] You are a Visual Director for voice-over videos. The user has ALREADY defined 
exact shot boundaries. Your ONLY job is to add visual and cinematic metadata to each 
shot. You MUST NOT split, merge, reorder, or skip any shots.

[CRITICAL RULE — PRESERVE STRUCTURE]
The user has provided exactly N shots using "// SHOT XX" markers. You MUST output 
EXACTLY N shots with matching shot_id numbers. Each block of content between markers
is ONE shot. Do NOT split any shot into multiple shots. Do NOT merge shots together.

[YOUR RESPONSIBILITIES (ONLY these)]
For each pre-defined shot, generate ONLY:
1. visual_description — Detailed scene description for image/video generation
2. shot_type — If not provided by user (ELS/LS/MS/CU/ECU)
3. angle — Camera angle (eye-level/high-angle/low-angle)
4. movement — Camera movement (fixed/push-in/pan/etc.)
5. location — Inferred setting description
6. time — Time of day / lighting condition
7. atmosphere — Mood and ambience
8. audio_mode — If not provided by user (inferred from [Tags])
9. narrator_enabled / narrator_ducking — Based on audio_mode
10. dialogue_type — If dialogue tags present (reaction/quote/crowd/etc.)
11. ambience_type / ambience_level / music_mood / music_level
12. sound_effect / bgm_prompt
13. characters / props / scene_id — Match to available entity lists
14. visual_type / shot_role / reason_for_shot / triggered_rules

[HOW TO INFER audio_mode from content]
- Lines WITHOUT [tags] → narrator_only
- Lines with [Character Name] → dialogue_dominant, narrator_enabled=false
- Lines with [CROWD] → dialogue_dominant, dialogue_type=crowd, narrator_ducking=true
- Lines with [SFX] → does NOT change audio_mode, put in sound_effect field
- If shot has BOTH narrator text AND dialogue → dialogue_dominant with narrator_ducking=true

[DURATION HANDLING]
- If user specified duration (e.g., "// SHOT 01 | 6s") → use that exact value
- If not specified → estimate based on word count at ~140 WPM
- Duration WARNING: If a shot exceeds 10s, add "LONG_SHOT" to triggered_rules

[IMPORTANT — WHAT YOU MUST NOT DO]
- Do NOT change the script_segment content
- Do NOT split any shot into multiple shots
- Do NOT merge adjacent shots
- Do NOT reorder shots
- Do NOT skip any shot
- Do NOT invent dialogue that isn't in the script
- Your output array length MUST equal exactly the number of // SHOT markers
```

### Phase 3: Backend Routing — `processVisualUnitGeneration()`

**File:** `application/services/storyboard_service.go`

Modify `processVisualUnitGeneration()` to detect and route:

```go
func (s *StoryboardService) processVisualUnitGeneration(...) {
    // NEW: Check for structured shot markers
    isStructured := detectStructuredShots(scriptContent)
    
    var systemPrompt string
    var scriptAnalysis string
    
    if isStructured {
        s.log.Infow("Detected structured shot markers, using ENRICH-ONLY mode",
            "task_id", taskID, "episode_id", episodeID)
        
        shotBlocks := parseStructuredShots(scriptContent)
        systemPrompt = prompts.Get("storyboard_visual_unit_structured.txt")
        scriptAnalysis = buildStructuredAnalysis(shotBlocks)
    } else {
        systemPrompt = s.promptI18n.WithDramaVisualUnitSystemPrompt(dramaID)
        _, scriptAnalysis = parseMarkedScript(scriptContent)
    }
    
    // ... rest unchanged — prompt assembly, AI call, parse, save
}
```

Add `buildStructuredAnalysis()` that summarizes shot blocks for AI context.

Add post-validation:
```go
if isStructured {
    expectedCount := len(shotBlocks)
    if len(shots) != expectedCount {
        s.log.Warnw("AI output count mismatch with structured input",
            "expected", expectedCount, "actual", len(shots))
        // Re-number to match — fallback safety
    }
}
```

### Phase 4: Frontend & Documentation

**4a. Tooltip update** — `web/src/locales/zh-CN.ts` and `en-US.ts`

Add mention of `// SHOT` syntax to visual unit tooltip.

**4b. PromptI18n** — Add resolver for structured prompt (fallback to file if no custom template).

**4c. Documentation** — Update voiceover script format guide with `// SHOT` examples.

## Acceptance Criteria

### Core Functionality
- [ ] `detectStructuredShots()` returns `true` when script has >= 3 `// SHOT` markers
- [ ] `parseStructuredShots()` correctly parses shot blocks with optional metadata
- [ ] `[Tag]` parsing within each shot block works correctly
- [ ] `isMarkdownMetadata()` no longer strips `// SHOT XX` lines
- [ ] AI output shot count equals **exactly** the input shot count
- [ ] User-specified duration/type/mode are preserved in output
- [ ] Free-form mode (no markers) behavior is **100% unchanged**

### Edge Cases
- [ ] Shot with no content inside → handled gracefully
- [ ] Mixed: some shots have metadata hints, others don't → auto-detect for missing
- [ ] Shot with both narrator text AND dialogue → correct audio_mode inference
- [ ] Very long shot (>10s) → warning but NOT split
- [ ] Markdown headers between shots → ignored properly

### Backward Compatibility
- [ ] Scripts without `// SHOT` markers → identical behavior to current system
- [ ] Custom prompt templates → still override structured prompt
- [ ] Downstream pipeline (image/video/merge) → works normally
- [ ] Preserve mode (`splitMode="preserve"`) → unaffected

## Technical Considerations

### Dependencies
- No new Go packages required
- No database schema changes
- No new API endpoints
- No new frontend components

### Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| AI ignores "don't split" instruction | Medium | Strong prompt + explicit shot count in analysis. Post-validation to reject mismatches |
| User puts `// SHOT` in non-structural context | Low | Require >= 3 markers for detection |
| `isMarkdownMetadata()` regex change breaks `//` filtering | Medium | Only carve out `// SHOT \d+` — all other `//` still filtered |
| Structured prompt gets too long | Low | Analysis is ~1 line per shot. 80 shots = ~80 lines |

### Alternatives Considered

| Alternative | Reason Rejected |
|---|---|
| New split mode button ("Structured Visual Unit") | Adds UI complexity. Auto-detect follows `auto→preserve` precedent |
| Frontend pre-parser that sends structured JSON | Requires frontend + API changes. Backend auto-detect is simpler |
| Modify existing `parseMarkedScript()` for shot blocks | Complicates working parser. Separate function is cleaner |

## Implementation Steps

### Phase 1: Parser (Backend) — ✅ DONE
1. ✅ Add `ShotBlock` struct
2. ✅ Implement `detectStructuredShots()`
3. ✅ Implement `parseStructuredShots()`
4. ✅ Fix `isMarkdownMetadata()` to not strip `// SHOT XX`
5. ✅ Add `parseSegmentsInBlock()` helper (inline in flushBlock)

### Phase 2: Prompt (Backend) — ✅ DONE
1. ✅ Create `storyboard_visual_unit_structured.txt`
2. ✅ Prompt resolver uses custom template fallback logic
3. ✅ Implement `buildStructuredAnalysis()`

### Phase 3: Routing & Validation (Backend) — ✅ DONE
1. ✅ Modify `processVisualUnitGeneration()` routing
2. ✅ Add post-validation for shot count
3. ✅ Wire user metadata into saved `VoiceoverShot`
4. ✅ Add logging

### Phase 4: Frontend & Docs — ✅ DONE
1. ✅ Update tooltip text in locales
2. ✅ Update script format guide
3. (Placeholder text update deferred — low priority)

## Post-Implementation Test

Convert user's example to structured format and verify exactly 5 shots output:

```
// SHOT 01 | 6s | CU | narrator_only
[SFX] Tiếng thở gấp gáp, nặng nề dội lại trong không gian hẹp.
Góc nhìn POV: Ánh đèn pin quét qua vách đất ẩm ướt, đầy rễ cây.

// SHOT 02 | 5s | EWS | narrator_only
Sâu dưới những khu rừng rậm rạp của miền Nam Việt Nam, 
một cuộc chiến bí mật đang diễn ra.

// SHOT 03 | 5s | MS | narrator_only
[SFX] Tiếng đất cát lạo xạo khi có người trườn qua.
Trong không gian chật hẹp này, không có chỗ cho xe tăng.
Đây là thế giới của những "Tunnel Rat".

// SHOT 04 | 4s | MCU | dialogue_dominant
[Soldier] Cố lên nào... đừng kẹt lại lúc này.
[SFX] Tiếng tim đập thình thịch tăng dần.

// SHOT 05 | 5s | ECU | narrator_only
Ánh đèn pin dừng lại ở một con bọ cạp đang bò trên trần hầm.
Chỉ một sơ suất nhỏ, bóng tối sẽ nuốt chửng lấy bạn mãi mãi.
```
