# Shot-Level Style Distillation

> Created: 2026-04-19
> Status: Implemented (Phase 1-3) — Phase 4 (Frontend) deferred

## Summary

Refactor pipeline tạo prompt ảnh/video để sử dụng kiến trúc **2-stage per-shot style distillation**:
- **Stage 1**: Batch distill `style_prompt` + `video_constraint` thành per-shot styles ngay sau khi storyboard tạo xong
- **Stage 2**: Sử dụng per-shot styles (thay vì full template prompts) khi generate image/video prompts

Giải quyết 3 vấn đề gốc:
1. Template prompts (dành cho LLM) bị ghép sai vào Image/Video API prompts — ảnh hưởng tất cả loại ảnh (character, scene, prop, shot)
2. Style chung bao quát nhiều trạng thái gây nhiễu khi inject vào prompt của 1 shot cụ thể
3. Character/Scene/Prop image generation bị double/triple style injection

## Problem Statement

### Vấn đề 1: Template prompt contamination

`style_prompt` và `video_constraint` là prompt dạng "Role: Art Director..." dành cho LLM, nhưng đang bị **ghép nguyên vẹn** vào prompt gửi trực tiếp cho AI tạo ảnh (Flux/GPT-Image) và AI tạo video (Kling/Doubao):

- [image_generation_service.go:280](file:///g:/VS-Project/huobao-drama/application/services/image_generation_service.go#L280): `prompt = stylePrompt + "\n\n" + prompt`
- [video_generation_service.go:409](file:///g:/VS-Project/huobao-drama/application/services/video_generation_service.go#L409): `prompt = constraintPrompt + "\n\n" + prompt`

### Vấn đề 2: Style bleeding across shots

`style_prompt` chung được inject toàn bộ vào mỗi shot qua `%s` trong image templates. Shot close-up mắt nhận cả phần không liên quan — gây LLM confuse.

### Vấn đề 3: video_constraint hoàn toàn bị bypass ở tầng LLM

`video_constraint` **KHÔNG BAO GIỜ** tham gia ở tầng LLM. Nó chỉ bị ghép sai ở tầng Video API. Nếu bỏ ghép đó thì `video_constraint` hoàn toàn không được sử dụng.

### Vấn đề 4: Character/Scene/Prop image bị contamination tương tự

Các extraction prompts (`character_extraction.txt`, `scene_extraction.txt`, `prop_extraction.txt`) đã inject `style_prompt` qua `%s` ở tầng LLM — LLM đã viết appearance/prompt **có style**. Nhưng khi generate image, `image_generation_service.go:280` lại ghép thêm `style_prompt` lần nữa.

| Entity | Extraction LLM có style? | Service riêng ghép style? | Image API ghép style (line 280)? |
|--------|---|---|---|
| Character | ✅ Có (qua `%s`) | Đã bỏ ([line 353](file:///g:/VS-Project/huobao-drama/application/services/character_library_service.go#L353)) | ❌ Ghép thừa |
| Scene | ✅ Có (qua `%s`) | — | ❌ Ghép thừa |
| Prop | ✅ Có (qua `%s`) | ❌ Hardcode "anime" ([line 175](file:///g:/VS-Project/huobao-drama/application/services/prop_service.go#L175)) | ❌ Ghép thừa |
| Shot | ✅ Có (qua `%s` khi AI-generated) | — | ❌ Ghép thừa |

**Kết luận**: Bỏ `image_generation_service.go:280` giải quyết contamination cho **tất cả** loại image. Không cần distill cho character/scene/prop — chỉ cho shots.

## Prior Solutions

- [custom-prompt-templates.md](file:///g:/VS-Project/huobao-drama/plans/custom-prompt-templates.md) — Plan gốc cho hệ thống template (Implemented). Plan này mở rộng hệ thống đó.

## Research Findings

### Codebase Patterns

**Luồng Image hiện tại (LLM-generated prompt):**
```
frame_prompt_service.go:253 → WithDramaFirstFramePrompt()
  → resolveEffectiveStyle() → style_prompt từ template (~1000c)
  → fmt.Sprintf(image_first_frame.txt, fullStylePrompt, ratio)
  → LLM viết image prompt
  → image_generation_service.go:280 → ghép THÊM style_prompt lần 2 (❌)
  → Gửi Flux/GPT-Image
```

**Luồng Video hiện tại:**
```
storyboard_service.go:1629 → generateVideoPrompt()
  → Ghép shot properties trực tiếp (KHÔNG qua LLM, KHÔNG dùng template)
  → video_generation_service.go:409 → ghép video_constraint (~3400c) (❌)
  → Gửi Kling/Doubao
```

**Key files:**
- [frame_prompt_service.go](file:///g:/VS-Project/huobao-drama/application/services/frame_prompt_service.go) — LLM image/video prompt generation
- [storyboard_service.go:1629](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L1629) — Auto video prompt composition
- [prompt_i18n.go:58-79](file:///g:/VS-Project/huobao-drama/application/services/prompt_i18n.go#L58-L79) — `resolveEffectiveStyle()`
- [prompt_template_service.go](file:///g:/VS-Project/huobao-drama/application/services/prompt_template_service.go) — Template resolver
- [image_generation_service.go:276-286](file:///g:/VS-Project/huobao-drama/application/services/image_generation_service.go#L276-L286) — Image API style injection
- [video_generation_service.go:407-414](file:///g:/VS-Project/huobao-drama/application/services/video_generation_service.go#L407-L414) — Video API constraint injection

**3 nguồn video_prompt:**

| Source | Code | Có LLM? |
|--------|------|---------|
| `auto` | `storyboard_service.go:1629` | Không |
| `ai` | `frame_prompt_service.go:573` | Có |
| `manual` | `storyboard_update_full.go:191` | Không |

## Proposed Solution

### Architecture: 2-Stage Per-Shot Style Distillation

```
User clicks "Split Shots"
│
├─ Step 1: AI tạo Storyboard (HIỆN TẠI — không đổi)
│  └─ Output: N shots với properties
│
└─ Step 2: Batch Distill (TỰ ĐỘNG sau step 1)
   ├─ LLM Call A: style_prompt + all shots → image_style per shot
   ├─ LLM Call B: video_constraint + all shots → video_style per shot
   │  (2 calls chạy SONG SONG)
   └─ Lưu image_style + video_style vào DB per shot

Sau đó khi generate:
  Image: image_keyframe.txt + shot.image_style → LLM → prompt → Image API (sạch)
  Video: shot properties + shot.video_style → compose → Video API (sạch)
```

### Distill Prompt Design

**Nguyên tắc**: Prompt distill phải **tổng quát** — hoạt động cho mọi template (Ghibli, pixel, CG5, The Well Studio, realistic...). Không chứa ví dụ của bất kỳ style cụ thể nào.

**Image Style Distill Prompt** (`image_style_distill.txt`):

```
You are a Shot-Level Visual Style Adapter. Your task is to read a general Art Direction Guide and produce a concise, shot-specific visual style description for each storyboard shot.

[Art Direction Guide — Reference Material]
--- BEGIN STYLE GUIDE ---
%s
--- END STYLE GUIDE ---

[Instructions]
For each shot listed below, analyze the Art Direction Guide and extract ONLY the visual style elements that are relevant to that specific shot's content, composition, and mood.

Your output for each shot should be a concise free-text description (50-150 words) written as direct visual descriptors suitable for AI image generation. Do NOT write instructions or role descriptions — write visual attributes.

Focus on:
- Rendering style and medium that apply to the shot's subject matter
- Color palette and lighting appropriate for the shot's atmosphere
- Texture and material qualities visible in the shot's framing
- Composition principles that enhance the shot's intent

Omit any style elements from the Guide that are irrelevant to the specific shot's content or framing.

[Storyboard Shots]
%s

[Output Format]
Return ONLY a valid JSON array. Do not include any explanation outside the JSON.

[
  {"shot_number": 1, "image_style": "<concise visual style for this shot>"},
  {"shot_number": 2, "image_style": "<concise visual style for this shot>"}
]

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response, including all JSON values, STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

**Video Style Distill Prompt** (`video_style_distill.txt`):

```
You are a Shot-Level Video Constraint Adapter. Your task is to read a general Video Production Guide and produce a concise, shot-specific video generation constraint for each storyboard shot.

[Video Production Guide — Reference Material]
--- BEGIN CONSTRAINT GUIDE ---
%s
--- END CONSTRAINT GUIDE ---

[Instructions]
For each shot listed below, analyze the Video Production Guide and extract ONLY the constraints and directives that are relevant to that specific shot's action, camera work, and pacing.

Your output for each shot should be a concise free-text description (50-150 words) written as direct video generation directives. Do NOT write role descriptions or system instructions — write actionable production constraints.

Focus on:
- Camera behavior appropriate for the shot's type, angle, and movement
- Motion characteristics relevant to the shot's action and subject
- Pacing and timing suitable for the shot's duration and intensity
- Visual continuity rules that apply to the shot's context

Omit any constraints from the Guide that are irrelevant to the specific shot's content or camera setup.

[Storyboard Shots]
%s

[Output Format]
Return ONLY a valid JSON array. Do not include any explanation outside the JSON.

[
  {"shot_number": 1, "video_style": "<concise video constraint for this shot>"},
  {"shot_number": 2, "video_style": "<concise video constraint for this shot>"}
]

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response, including all JSON values, STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

**Cách ghép `style_prompt` mượt mà**: `style_prompt` nằm giữa `--- BEGIN/END STYLE GUIDE ---` markers. Dù nó viết dạng "You are a top Art Director..." — LLM vẫn hiểu đây là tài liệu tham chiếu, không phải role của mình, vì role đã define ở đầu ("You are a Shot-Level Visual Style Adapter").

### DB Changes

Thêm 2 fields vào bảng `storyboards`:

```go
// domain/models/drama.go — Storyboard struct
ImageStyle *string `gorm:"type:text" json:"image_style,omitempty"`
VideoStyle *string `gorm:"type:text" json:"video_style,omitempty"`
```

GORM auto-migrate sẽ tự thêm columns. Nullable — các shots cũ không bị ảnh hưởng.

### Service Changes

#### 1. Tạo `StyleDistillService` (mới)

```go
// application/services/style_distill_service.go

type StyleDistillService struct {
    db         *gorm.DB
    aiService  *AIService
    promptI18n *PromptI18n
    log        *logger.Logger
}

// BatchDistillStyles chạy 2 LLM calls song song cho tất cả shots
func (s *StyleDistillService) BatchDistillStyles(episodeID uint, dramaID uint) error {
    // 1. Load all storyboards for episode
    // 2. Load drama → get style_prompt + video_constraint from template
    // 3. If no template or no style_prompt → skip image distill
    // 4. If no template or no video_constraint → skip video distill
    // 5. Build shot context JSON array
    // 6. Parallel:
    //    a. LLM call with image_style_distill.txt + style_prompt + shots
    //    b. LLM call with video_style_distill.txt + video_constraint + shots
    // 7. Parse JSON response → update storyboards.image_style / video_style
}
```

#### 2. Tích hợp vào Split Shots flow

| Flow | File | Method | Trigger distill |
|------|------|--------|-----------------|
| Standard split | `storyboard_service.go` | `saveStoryboards()` | ✅ |
| Voiceover director | `storyboard_service.go` | `SaveVoiceoverShots()` | ✅ |
| Nursery rhyme | `storyboard_nursery_service.go` | `SaveNurseryShots()` | ✅ |
| MV Maker | `storyboard_mv_service.go` | `SaveMVShots()` | ✅ |

Sau mỗi save thành công → gọi `styleDistillService.BatchDistillStyles()`.

#### 3. Tạo method mới cho per-shot style lookup (KHÔNG sửa `resolveEffectiveStyle()`)

`resolveEffectiveStyle()` có **11 callers** (extraction prompts, frame prompts, video extraction...). Extraction prompts (character/scene/prop) vẫn cần full `style_prompt` — chúng trích xuất từ kịch bản, không phải per-shot. **Chỉ frame prompts** (first/key/last/action) cần dùng `shot.image_style`.

```go
// Tạo method mới trong frame_prompt_service.go:
func resolveStyleForShot(storyboardID uint, dramaID uint, style string, customStyle string) string {
    // 1. Load storyboard → check image_style
    // 2. Nếu image_style != nil → return image_style (đã distill)
    // 3. Else → fallback: resolveEffectiveStyle(dramaID, style, customStyle)
}

// Chỉ 4 frame prompt methods gọi method mới này:
// - WithDramaFirstFramePrompt → resolveStyleForShot
// - WithDramaKeyFramePrompt → resolveStyleForShot  
// - WithDramaLastFramePrompt → resolveStyleForShot
// - WithDramaActionSequenceFramePrompt → resolveStyleForShot

// 7 callers khác giữ nguyên resolveEffectiveStyle():
// - extraction prompts (character/scene/prop) — vẫn cần full style_prompt
// - video_extraction — vẫn cần full style_prompt  
// - rapid_cut — vẫn cần full style_prompt
```

#### 4. Ghép video_style ở tầng Video API (thay vì sửa `generateVideoPrompt()`)

Vấn đề: `generateVideoPrompt()` được gọi trong `saveStoryboards()` — lúc đó shot chưa có ID, chưa có `video_style` (distill chạy sau). Giải pháp: **giữ `generateVideoPrompt()` như cũ**, thay đổi tại `video_generation_service.go`:

```go
// video_generation_service.go — thay vì prepend video_constraint:
// 1. Load storyboard → check video_style
// 2. Nếu video_style != nil → prepend video_style (đã distill, ngắn gọn)
// 3. Else → skip (không prepend gì)
```

Đơn giản hơn, ít breaking change.

#### 5. Bỏ API-level injection

```diff
// image_generation_service.go:276-286 — ảnh hưởng TẤT CẢ loại image (character/scene/prop/shot)
- if drama.Style != "" && drama.Style != "realistic" {
-     stylePrompt := s.promptI18n.WithDramaStylePrompt(...)
-     prompt = stylePrompt + "\n\n" + prompt
- }

// video_generation_service.go:407-414
- constraintPrompt := s.promptI18n.WithDramaVideoConstraintPrompt(...)
- prompt = constraintPrompt + "\n\n" + prompt
+ // Thay bằng: prepend shot.video_style nếu có
```

#### 6. Bỏ hardcode style trong Prop

```diff
// prop_service.go:175
- imageStyle := "Modern Japanese anime style"
```

Prop extraction đã có style qua `%s` trong `prop_extraction.txt` → prop.Prompt đã bao gồm style. Hardcode "anime" gây conflict khi dùng template khác (VD: CG5 3D CGI).

## Acceptance Criteria

### Distillation
- [ ] Tạo prompt files `image_style_distill.txt` và `video_style_distill.txt`
- [ ] Thêm fields `image_style`, `video_style` vào model `Storyboard` (auto-migrate)
- [ ] Tạo `StyleDistillService` với `BatchDistillStyles()`
- [ ] Batch distill chạy tự động sau storyboard tạo xong (tất cả flows)
- [ ] `image_style` / `video_style` lưu vào DB per shot
- [ ] Define rõ shot context format cho distill: `shot_number`, `action`, `result`, `location`, `atmosphere`, `shot_type`, `angle`, `movement`, `characters`

### Integration
- [ ] Tạo `resolveStyleForShot()` — chỉ 4 frame prompt methods dùng, 7 callers khác giữ nguyên
- [ ] `video_generation_service.go` prepend `shot.video_style` thay vì `video_constraint`

### Cleanup (ảnh hưởng tất cả image types: character/scene/prop/shot)
- [ ] Bỏ ghép `style_prompt` ở `image_generation_service.go:280`
- [ ] Bỏ ghép `video_constraint` ở `video_generation_service.go:409`
- [ ] Bỏ hardcode `"Modern Japanese anime style"` ở `prop_service.go:175`

### Backward Compatibility
- [ ] Drama không có template → skip distill, behavior cũ hoạt động
- [ ] Shots cũ không có image_style/video_style → fallback `resolveEffectiveStyle()`
- [ ] Character/Scene/Prop image gen hoạt động bình thường sau bỏ API injection (extraction đã có style)

### Frontend (Optional — có thể để sau)
- [ ] Hiển thị image_style/video_style trên UI
- [ ] Nút "Re-distill" khi user muốn regenerate styles

## Technical Considerations

### Dependencies
- Không cần package mới — sử dụng AIService hiện có cho LLM calls
- GORM auto-migrate xử lý DB migration

### Risks

1. **LLM output parsing**: Batch JSON response có thể malformed → Robust JSON extraction, fallback nếu parse fail
2. **Shot count lớn**: 50+ shots → JSON input/output dài → Chunk thành batches 20 shots nếu cần
3. **Tăng thời gian Split Shots**: Thêm 2 LLM calls → ~10-20s thêm → Chạy song song, user thấy progress
4. **Template không có style_prompt hoặc video_constraint**: Skip distill cho phần đó, fallback behavior cũ

### Alternatives Considered

1. **Inject full style_prompt vào %s (hiện tại)**: Bỏ vì gây nhiễu style — cùng 1 style prompt áp cho tất cả shots
2. **Bỏ style injection hoàn toàn**: Bỏ vì mất consistency — mỗi shot tự do phong cách
3. **Inject video_constraint vào video_extraction.txt**: Bỏ vì video_constraint quá dài và context-dependent

## Implementation Steps

### Phase 1: Cleanup — Bỏ API-level injection (có thể ship độc lập)
- Task 1: Bỏ ghép `style_prompt` ở `image_generation_service.go:276-286`
- Task 2: Bỏ ghép `video_constraint` ở `video_generation_service.go:407-414`
- Task 3: Bỏ hardcode `"Modern Japanese anime style"` ở `prop_service.go:175`
- Task 4: Test character/scene/prop/shot image gen vẫn hoạt động (extraction đã có style)

> **Lưu ý**: Phase 1 có thể ship trước, độc lập. Nó fix contamination ngay mà không cần distill. Distill (Phase 2-3) là enhancement tiếp theo.

### Phase 2: Foundation — Distill Service
- Task 5: Thêm `image_style`, `video_style` fields vào Storyboard model
- Task 6: Tạo prompt files `image_style_distill.txt`, `video_style_distill.txt`
- Task 7: Tạo `StyleDistillService` với `BatchDistillStyles()` — define shot context format: `{shot_number, action, result, location, atmosphere, shot_type, angle, movement, characters}`

### Phase 3: Integration — Per-Shot Style
- Task 8: Tích hợp distill vào `saveStoryboards()` (standard flow)
- Task 9: Tích hợp distill vào `SaveVoiceoverShots()`, `SaveNurseryShots()`, `SaveMVShots()`
- Task 10: Tạo `resolveStyleForShot()` — chỉ 4 frame prompt methods dùng (first/key/last/action). 7 callers khác (`resolveEffectiveStyle`) giữ nguyên
- Task 11: Sửa `video_generation_service.go` — prepend `shot.video_style` thay vì `video_constraint`
- Task 12: Test backward compatibility (drama không template, shots cũ, distill fail)

### Phase 4: Frontend (Optional — có thể để sau)
- Task 13: Hiển thị `image_style` / `video_style` trong shot detail panel
- Task 14: Nút "Re-distill" khi user muốn regenerate styles

## Edge Cases

| Case | Xử lý |
|------|--------|
| Drama không có template | Skip distill. Fallback behavior cũ |
| Template không có `style_prompt` | Skip image distill |
| Template không có `video_constraint` | Skip video distill |
| Shots cũ (trước feature) | `image_style`/`video_style` = NULL → `resolveStyleForShot()` fallback `resolveEffectiveStyle()` |
| User edit shot sau distill | Cần nút "Re-distill" hoặc chấp nhận style cũ |
| LLM trả về JSON malformed | Robust JSON extraction (tìm `[` đến `]`). Log warning, skip shots lỗi, các shots khác vẫn có style |
| Episode có 50+ shots | Chunk thành batches ~20 shots per LLM call |
| Drama dùng dropdown style (không template) | Có thể skip distill cho dropdown styles đơn giản (style text ngắn → distill output cũng ngắn, ít giá trị) |
| Distill LLM call fail (timeout/API error) | Log error, `image_style`/`video_style` = NULL, fallback behavior cũ. User thấy storyboard hoàn tất, có thể dùng "Re-distill" button sau |
| Character/Scene/Prop image gen sau bỏ API injection | Hoạt động bình thường — extraction đã inject style qua `%s`. Character `buildCharacterPrompt()` đã không ghép style |
| Prop image gen | Cần bỏ hardcode "anime" ở `prop_service.go:175` để tránh conflict với template style |

## References

- [custom-prompt-templates.md](file:///g:/VS-Project/huobao-drama/plans/custom-prompt-templates.md) — Plan gốc (Implemented)
- [prompt_i18n.go](file:///g:/VS-Project/huobao-drama/application/services/prompt_i18n.go) — Style resolution
- [frame_prompt_service.go](file:///g:/VS-Project/huobao-drama/application/services/frame_prompt_service.go) — LLM frame prompt generation
- [storyboard_service.go](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go) — Storyboard creation + auto prompts
