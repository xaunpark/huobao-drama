# Voice-over AI Director Mode — Chế Độ Đạo Diễn Video Voice-over

> Created: 2026-04-05
> Status: Implemented (Phase 1-3 Core)

## Summary

Thêm chế độ `splitMode = "visual_unit"` vào pipeline storyboard hiện tại. Chế độ này nhận script narrator và tạo shot list theo **visual unit** (thay vì action unit hoặc sentence boundary), kèm theo **audio strategy per-shot** (narrator_only / mixed / dialogue_dominant). Tương thích hoàn toàn với hệ thống Prompt Template và toàn bộ downstream pipeline (FramePrompt → Image → Video → Merge).

## Problem Statement

Pipeline hiện tại có 2 split modes:
- `breakdown`: Chia shot theo **action unit** (McKee theory) — phù hợp cho drama có nhân vật, thoại, hành động
- `preserve`: Giữ shot boundaries từ script có timestamp — phù hợp cho script đã chia sẵn

**Thiếu**: Chế độ cho **video voice-over** (narrator-driven) — nơi:
- Input là script narrator thuần (không có đối thoại nhân vật)
- Shot cần chia theo **thay đổi hình ảnh** (visual change), không theo grammar
- 1 câu có thể → nhiều shot, nhiều câu có thể → 1 shot
- Mỗi shot cần gắn **đoạn narrator tương ứng** cho post-production (TTS, timeline)
- Mỗi shot cần **quyết định audio** (narrator on/off, có dialogue accent không, ambience/music)

## Prior Solutions

Không tìm thấy solution hiện có trong `docs/solutions/`. Tuy nhiên:
- [custom-prompt-templates.md](file:///g:/VS-Project/huobao-drama/plans/custom-prompt-templates.md) — Pattern tạo prompt type mới + resolver fallback (**reuse pattern**)
- [rapid-cut-mode.md](file:///g:/VS-Project/huobao-drama/plans/rapid-cut-mode.md) — Pattern thêm mode mới vào Storyboard model với nullable fields (**reuse pattern**)

## Research Findings

### Codebase Patterns

**Split Mode mechanism** — [storyboard_service.go:144-154](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L144-L154):
- `splitMode` đã có auto-detect logic + router giữa `breakdown` / `preserve`
- Thêm `case "visual_unit"` là mở rộng tự nhiên của switch hiện tại

**Prompt Template system** — [prompt_template.go:29-43](file:///g:/VS-Project/huobao-drama/domain/models/prompt_template.go#L29-L43):
- `PromptTemplatePrompts` struct có 13 fields, mỗi field map sang 1 prompt type
- `PromptTypeToDefaultFile` map prompt type key → default .txt file
- `getPromptFromStruct()` switch statement resolve custom prompt

**Storyboard Model** — [drama.go:94-131](file:///g:/VS-Project/huobao-drama/domain/models/drama.go#L94-L131):
- Đã có pattern thêm nullable fields cho Rapid Cut (`IsProduction`, `PacingMode`, `SourceShotIDs`)
- GORM AutoMigrate tự thêm cột → migration dễ

**PromptI18n resolver** — [prompt_i18n.go:409-411](file:///g:/VS-Project/huobao-drama/application/services/prompt_i18n.go#L409-L411):
- `WithDramaStoryboardSystemPrompt()` → `resolvePrompt(dramaID, "storyboard_breakdown")`
- Pattern: mỗi mode có `WithDrama{Mode}SystemPrompt()` method riêng

### Reference Documents

- [v1 Analysis](file:///C:/Users/dinht/.gemini/antigravity/brain/d96c2f1f-f50c-4f4f-8b7c-4575189c0447/voiceover_director_integration_analysis.md) — Pipeline mapping, field-by-field comparison
- [v2 Analysis](file:///C:/Users/dinht/.gemini/antigravity/brain/d96c2f1f-f50c-4f4f-8b7c-4575189c0447/voiceover_director_v2_deep_analysis.md) — ScriptSegment, Prompt Template compatibility, Audio Strategy

## Proposed Solution

### Approach: Extend Existing Pipeline

Thêm `visual_unit` như splitMode thứ 3, theo đúng pattern đã dùng cho `preserve` và `rapid_cut`:

```
StoryboardService.GenerateStoryboard()
  splitMode: "breakdown"      ← existing (McKee action units)
  splitMode: "preserve"       ← existing (keep timestamps)
  splitMode: "visual_unit"    ← NEW (AI Director visual change rules + audio strategy)
```

### Data Model Changes

**15 nullable fields mới** trong `models.Storyboard`:

```go
// === Voice-over Director fields ===
ScriptSegment    *string        `gorm:"type:text" json:"script_segment"`
ScriptStartChar  *int           `json:"script_start_char"`
ScriptEndChar    *int           `json:"script_end_char"`
ShotReason       *string        `gorm:"type:text" json:"shot_reason"`
SplitRules       datatypes.JSON `gorm:"type:json" json:"split_rules"`
VisualType       *string        `gorm:"size:20" json:"visual_type"`
ShotRole         *string        `gorm:"size:30" json:"shot_role"`

// === Audio Strategy fields ===
AudioMode        *string        `gorm:"size:30" json:"audio_mode"`
NarratorEnabled  *bool          `json:"narrator_enabled"`
NarratorDucking  *bool          `json:"narrator_ducking"`
DialogueType     *string        `gorm:"size:30" json:"dialogue_type"`
AmbienceType     *string        `gorm:"size:50" json:"ambience_type"`
AmbienceLevel    *string        `gorm:"size:10" json:"ambience_level"`
MusicMood        *string        `gorm:"size:50" json:"music_mood"`
MusicLevel       *string        `gorm:"size:10" json:"music_level"`
```

**1 field mới** trong `models.PromptTemplatePrompts`:

```go
VisualUnitBreakdown string `json:"visual_unit_breakdown,omitempty"`
```

**1 entry mới** trong `PromptTypeToDefaultFile`:

```go
"visual_unit_breakdown": "storyboard_visual_unit.txt",
```

### AI Output Parse Struct

```go
// VoiceoverShot — AI output struct cho visual_unit mode
type VoiceoverShot struct {
    ShotID             int      `json:"shot_id"`
    ScriptSegment      string   `json:"script_segment"`
    ScriptStartChar    int      `json:"script_start_char"`
    ScriptEndChar      int      `json:"script_end_char"`
    EstimatedDuration  int      `json:"estimated_duration_sec"`
    VisualType         string   `json:"visual_type"`          // literal | symbolic
    ShotRole           string   `json:"shot_role"`            // establishing | action | detail | emotional | symbolic | transition | closing
    VisualDescription  string   `json:"visual_description"`
    ReasonForShot      string   `json:"reason_for_shot"`
    TriggeredRules     []string `json:"triggered_rules"`      // ["scene_change", "action_change", ...]

    // Shot cinematography (map to existing Storyboard fields)
    Title       string `json:"title"`
    ShotType    string `json:"shot_type"`
    Angle       string `json:"angle"`
    Movement    string `json:"movement"`
    Location    string `json:"location"`
    Time        string `json:"time"`
    Atmosphere  string `json:"atmosphere"`

    // Audio Strategy
    AudioMode       string `json:"audio_mode"`        // narrator_only | dialogue_dominant
    NarratorEnabled bool   `json:"narrator_enabled"`
    NarratorDucking bool   `json:"narrator_ducking"`
    DialogueType    string `json:"dialogue_type"`     // none | reaction | soft_line | quote
    DialogueText    string `json:"dialogue_text"`
    AmbienceType    string `json:"ambience_type"`
    AmbienceLevel   string `json:"ambience_level"`    // low | medium | high
    MusicMood       string `json:"music_mood"`
    MusicLevel      string `json:"music_level"`       // low | medium | high
    SoundEffect     string `json:"sound_effect"`
    BgmPrompt       string `json:"bgm_prompt"`

    // Existing fields pass-through
    Characters []uint      `json:"characters"`
    Props      []uint      `json:"props"`
    SceneID    *uint       `json:"scene_id"`
}
```

### Mapping VoiceoverShot → models.Storyboard

```
VoiceoverShot.ScriptSegment    → Storyboard.ScriptSegment
VoiceoverShot.VisualDescription → Storyboard.Description
VoiceoverShot.ReasonForShot    → Storyboard.ShotReason
VoiceoverShot.TriggeredRules   → Storyboard.SplitRules (JSON)
VoiceoverShot.VisualType       → Storyboard.VisualType
VoiceoverShot.ShotRole         → Storyboard.ShotRole
VoiceoverShot.Title            → Storyboard.Title
VoiceoverShot.ShotType         → Storyboard.ShotType
VoiceoverShot.DialogueText     → Storyboard.Dialogue (reuse existing field)
VoiceoverShot.AudioMode        → Storyboard.AudioMode
VoiceoverShot.*                → Storyboard.* (matching fields)
```

### Prompt Template Integration

```
Resolver chain (same pattern as existing):

1. User chọn visual_unit mode
2. storyboard_service.go calls: s.promptI18n.WithDramaVisualUnitSystemPrompt(dramaID)
3. prompt_i18n.go calls: p.resolvePrompt(dramaID, "visual_unit_breakdown")
4. prompt_template_service.go checks: template.VisualUnitBreakdown (custom?)
   → YES: return custom prompt
   → NO: return prompts.Get("storyboard_visual_unit.txt")
```

Frontend Settings page thêm 1 textarea trong tab "🎬 Phân cảnh":
- Label: "Visual Unit Breakdown (Voice-over)"
- Placeholder: default prompt content

### Pipeline Flow

```
Script Input
    ├── Extract Characters (reuse 100%)
    ├── Extract Scenes (reuse 100%)
    ├── Extract Props (reuse 100%)
    └── Storyboard Generation
         └── splitMode = "visual_unit"
              ├── Prompt: storyboard_visual_unit.txt (or template custom)
              ├── AI output: VoiceoverShot[] (shot + audio per shot)
              ├── Parse → map → models.Storyboard (15 new fields populated)
              └── Save to DB
                   └── Downstream (reuse 100%):
                        ├── FramePrompt Service
                        ├── Image Generation
                        ├── Video Generation
                        └── Video Merge / Timeline
```

## Acceptance Criteria

### Core
- [ ] `splitMode = "visual_unit"` tạo shot list theo visual unit rules
- [ ] Mỗi shot có `ScriptSegment` = đoạn narrator tương ứng
- [ ] Mỗi shot có `ShotReason` + `SplitRules` giải thích lý do chia
- [ ] Mỗi shot có `AudioMode` (narrator_only / dialogue_dominant)
- [ ] Không shot nào vượt 8 giây (hoặc `max_shot_duration_sec` nếu custom)
- [ ] 1 câu có thể → nhiều shot khi có nhiều visual unit
- [ ] Nhiều câu có thể → 1 shot khi cùng visual unit

### Prompt Template
- [ ] `visual_unit_breakdown` hiển thị trong Prompt Template Settings
- [ ] Template custom `visual_unit_breakdown` được dùng thay default khi gán vào Drama
- [ ] Template KHÔNG có custom → fallback về `storyboard_visual_unit.txt`
- [ ] Custom `storyboard_breakdown` KHÔNG ảnh hưởng `visual_unit` mode và ngược lại

### Audio Strategy
- [ ] ~80-90% shots là `narrator_only` chuộc từ visual
- [ ] AI tự quyết dialogue type per-shot (user không cần chỉ định)
- [ ] `Dialogue` field populated khi `dialogue_type != "none"` (reaction/quote/...)
- [ ] `NarratorEnabled = false` khi `audio_mode = "dialogue_dominant"`

### Backward Compatibility
- [ ] `breakdown` mode hoạt động 100% y hệt (zero regression)
- [ ] `preserve` mode hoạt động 100% y hệt
- [ ] Projects cũ không có visual_unit data → 15 fields mới = null → UI không hiện gì thêm
- [ ] Downstream pipeline (Image/Video/Merge) hoạt động bình thường với visual_unit shots

## Technical Considerations

### Dependencies
- Không thêm package Go mới (GORM `datatypes.JSON` đã có)
- Frontend: Element Plus đã có tabs, textarea, select — không cần thêm

### Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| AI output JSON format khác struct | 🟡 Medium | Parse VoiceoverShot riêng rồi map → Storyboard. SafeParseAIJSON đã xử lý nhiều format |
| AI chia shot quá nhiều/quá ít | 🟡 Medium | Duration limit trong prompt + lý do bắt buộc (self-check) |
| Prompt quá phức tạp (shot + audio + rules) | 🟡 Medium | Tách rõ sections trong prompt: [Visual Rules], [Audio Rules], [Output Schema] |
| Audio fields chưa dùng ngay downstream | 🟢 Low | Fields nullable, downstream ignore. Dùng khi có TTS + audio renderer integration |
| Format vars (%s) trong custom prompt | 🟢 Low | Dùng `formatPromptWithVars()` pattern từ prompt_i18n.go (đã có) |

### Alternatives Considered

| Alternative | Reason Rejected |
|-------------|----------------|
| Pipeline riêng (VoiceoverDirectorService) | Code duplication 70%+, maintenance gấp đôi, downstream vẫn phải dùng chung |
| Reuse `storyboard_breakdown` prompt type | Conflict khi user custom prompt cho drama mode — 2 mode phải tách biệt |
| Audio strategy as separate AI call | Mất semantic context từ shot planning step, thêm latency |
| Separate AudioPlannerService (bắt buộc) | Over-engineering cho Phase 1. Có thể thêm sau như optional re-plan |

## Implementation Steps

### Phase 1: Backend — Model + Prompt (Day 1)

- **Task 1.1**: Thêm 15 fields vào `models.Storyboard` trong [drama.go](file:///g:/VS-Project/huobao-drama/domain/models/drama.go#L94-L131)
  - 7 Voice-over Director fields + 8 Audio Strategy fields
  - Tất cả nullable, GORM AutoMigrate tự thêm cột
  
- **Task 1.2**: Thêm `VisualUnitBreakdown` vào `PromptTemplatePrompts` trong [prompt_template.go](file:///g:/VS-Project/huobao-drama/domain/models/prompt_template.go#L29-L43)
  - Thêm entry trong `PromptTypeToDefaultFile` map

- **Task 1.3**: Tạo prompt file `application/prompts/storyboard_visual_unit.txt`
  - AI Director role + Visual unit rules + Audio strategy rules + Output JSON schema
  - Chia rõ sections: [Role], [Visual Rules], [Audio Rules], [Output Format]

- **Task 1.4**: Tạo `application/prompts/fixed/storyboard_visual_unit_format.txt` (nếu cần)
  - JSON schema cứng cho VoiceoverShot output (hoặc embed trong prompt chính)

### Phase 2: Backend — Service Integration (Day 1-2)

- **Task 2.1**: Thêm `WithDramaVisualUnitSystemPrompt()` vào [prompt_i18n.go](file:///g:/VS-Project/huobao-drama/application/services/prompt_i18n.go)
  - `return p.resolvePrompt(dramaID, "visual_unit_breakdown")`

- **Task 2.2**: Thêm case `"visual_unit_breakdown"` vào `getPromptFromStruct()` trong [prompt_template_service.go](file:///g:/VS-Project/huobao-drama/application/services/prompt_template_service.go#L253-L283)
  - Thêm `VisualUnitBreakdown` vào `GetDefaultPrompts()` return

- **Task 2.3**: Thêm `case "visual_unit"` vào `processStoryboardGeneration()` trong [storyboard_service.go](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L206-L224)
  - Dùng `WithDramaVisualUnitSystemPrompt(dramaID)` cho system prompt
  - Tạo `VoiceoverShot` parse struct
  - Parse AI output → map sang `Storyboard` struct → populate 15 fields mới

- **Task 2.4**: Cập nhật `saveStoryboards()` hoặc tạo `saveVoiceoverShots()` variant
  - Populate `ScriptSegment`, `ShotReason`, `SplitRules`, `VisualType`, `ShotRole`
  - Populate `AudioMode`, `NarratorEnabled`, `NarratorDucking`, `DialogueType`
  - Populate `AmbienceType`, `AmbienceLevel`, `MusicMood`, `MusicLevel`
  - Map `DialogueText` → `Dialogue` field (reuse)
  - Map `VisualDescription` → `Description` field (reuse)
  - Generate `ImagePrompt` + `VideoPrompt` (reuse existing helper methods)

### Phase 3: Frontend — Split Mode + Display (Day 2-3)

- **Task 3.1**: Thêm "Visual Unit (Voice-over)" option vào split mode selector
  - Nơi user chọn mode trước khi generate storyboard
  - Label: "Voice-over / Narrator" hoặc tương tự

- **Task 3.2**: Cập nhật Prompt Template Settings page
  - Thêm textarea cho `visual_unit_breakdown` trong tab "🎬 Phân cảnh"
  - Label: "Visual Unit Breakdown (Voice-over)"
  - Load default button: gọi `/api/prompt-templates/defaults` → field `visual_unit_breakdown`

- **Task 3.3**: Cập nhật ProfessionalEditor.vue để hiển thị voice-over fields
  - Nếu shot có `ScriptSegment` ≠ null → hiển thị panel "📝 Narrator Script"
  - Nếu shot có `AudioMode` ≠ null → hiển thị badge audio mode
  - Nếu shot có `ShotReason` ≠ null → hiển thị collapsible reason/rules section
  - Conditional render: chỉ show khi fields có data (backward compatible)

- **Task 3.4**: Localization
  - Thêm labels cho tất cả fields mới (en + zh-CN + vi nếu có)

### Phase 4: Testing & Polish (Day 3-4)

- **Task 4.1**: Test visual_unit mode end-to-end
  - Input script narrator → generate storyboard → verify shot list quality
  - Verify ScriptSegment covers toàn bộ script (no gaps)
  - Verify no shot > 8 seconds

- **Task 4.2**: Test prompt template integration
  - Drama KHÔNG có template → default prompt
  - Drama CÓ template nhưng không custom visual_unit → default prompt  
  - Drama CÓ template VÀ custom visual_unit → custom prompt

- **Task 4.3**: Test backward compatibility
  - `breakdown` mode unchanged
  - `preserve` mode unchanged
  - Old projects with null fields render correctly

- **Task 4.4**: Test audio strategy output
  - Verify ~80-90% shots are narrator_only
  - Verify dialogue_dominant shots have narrator_enabled = false
  - Verify mixed shots have short dialogue (≤3 words for reaction)

## Rollback Strategy

Tất cả thay đổi là **additive** (nullable fields, new prompt file, new case branch). Rollback:
1. Remove `case "visual_unit"` từ storyboard_service.go → mode không khả dụng
2. Frontend: hide option từ selector
3. DB columns remain (nullable, harmless) — hoặc drop nếu cần

## References

- [v1 Analysis: Pipeline Mapping](file:///C:/Users/dinht/.gemini/antigravity/brain/d96c2f1f-f50c-4f4f-8b7c-4575189c0447/voiceover_director_integration_analysis.md)
- [v2 Analysis: ScriptSegment + Audio + Template](file:///C:/Users/dinht/.gemini/antigravity/brain/d96c2f1f-f50c-4f4f-8b7c-4575189c0447/voiceover_director_v2_deep_analysis.md)
- [custom-prompt-templates.md](file:///g:/VS-Project/huobao-drama/plans/custom-prompt-templates.md) — Pattern reference
- [rapid-cut-mode.md](file:///g:/VS-Project/huobao-drama/plans/rapid-cut-mode.md) — Pattern reference
- [storyboard_service.go](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go) — Primary integration point
- [prompt_template.go](file:///g:/VS-Project/huobao-drama/domain/models/prompt_template.go) — Template model
- [drama.go](file:///g:/VS-Project/huobao-drama/domain/models/drama.go) — Storyboard model
