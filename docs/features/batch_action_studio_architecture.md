# Batch Action Studio — Architecture & Data Flow

> **Mục đích**: Tài liệu hệ thống mô tả đầy đủ kiến trúc, luồng dữ liệu, và pipeline xử lý của **Batch Action Studio** — chế độ batch xử lý ảnh + video trong ProfessionalEditor.

---

## 1. Tổng quan

Batch Action Studio (`BatchGenerationDialog.vue`) là giao diện batch processing cho việc tạo ảnh và video từ storyboard. Nó xử lý 3 phase tuần tự:

```
[Extract Prompts] → [Generate Images] → [Generate Videos]
```

User có thể chạy tất cả phases ("Run All Phases") hoặc từng phase riêng lẻ.

### Vị trí trong codebase

| Component | File |
|---|---|
| **Frontend Dialog** | `web/src/components/editor/BatchGenerationDialog.vue` |
| **Storyboard Model** | `domain/models/drama.go` → `Storyboard` struct |
| **Image Prompt Gen (Pipeline 2)** | `application/services/frame_prompt_service.go` |
| **Mechanical Video Prompt (Pipeline 1)** | `application/services/storyboard_service.go` → `generateVideoPrompt()` |
| **Style Distillation (Pipeline 3)** | `application/services/style_distill_service.go` |
| **Video Generation Service** | `application/services/video_generation_service.go` |
| **Prompt Templates** | `application/prompts/` (embed files) |
| **Template Override System** | `application/services/prompt_template_service.go` |

---

## 2. Hai giai đoạn xử lý

### Giai đoạn A: Pre-distillation (tự động)

Chạy **tự động ngay sau khi storyboard được tạo** (không cần user trigger).
Gọi `BatchDistillStyles(episodeID, dramaID)` — chạy dưới dạng goroutine background.

```
BatchDistillStyles
├── Goroutine 1: distillImageStyles()
│   Input:  style_prompt (từ template) + shot contexts (batch 20 shots)
│   Output: image_style per shot (50-150 words)
│   Lưu:    storyboards.image_style
│
└── Goroutine 2: distillVideoCombined()
    Input:  video_constraint (từ template) + shot contexts (batch 20 shots)
    Output: video_style + video_prompt_distilled per shot
    Lưu:    storyboards.video_style, storyboards.video_prompt_distilled
```

**Đặc điểm:**
- Batch processing: 1 LLM call xử lý tối đa 20 shots
- 2 goroutines chạy song song → tổng 2 LLM calls cho toàn bộ episode
- Kết quả lưu trực tiếp vào DB columns của Storyboard

### Giai đoạn B: Batch Action Studio (user trigger)

User mở dialog → chọn model → bấm Run.

```
Phase 1: Extract Prompts  → processPrompt() per shot
Phase 2: Generate Images   → processImage() per shot
Phase 3: Generate Videos   → processVideo() per shot
```

---

## 3. Chi tiết từng Phase

### Phase 1: Extract Prompts (`processPrompt`)

```
Frontend:  generateFramePrompt(storyboard_id, { frame_type: 'action' | 'key' })
Backend:   frame_prompt_service.go → GenerateFramePrompt()
LLM Call:  1 call per shot
```

**Input cho LLM:**
- Shot properties (action, result, camera, location, atmosphere, dialogue...)
- Character descriptions + reference images
- Scene/background reference images
- Prop reference images
- Template prompt override: `video_extraction` (nếu template có) hoặc `video_extraction.txt` (default)

**Output:**
- `image_prompt` — prompt mô tả chi tiết cho image generation
- Lưu vào `storyboards.image_prompt` + sessionStorage

**Generation Mode (user chọn):**

| Mode | frame_type | Mô tả |
|---|---|---|
| **Action Sequence** | `action` | Grid 4-6 frames thể hiện chuyển động liên tục |
| **Keyframe** | `key` | 1 frame đại diện cho shot |

### Phase 2: Generate Images (`processImage`)

```
Frontend:  imageAPI.generateImage({ prompt, storyboard_id, frame_type, reference_images })
Backend:   image_generation_service.go
API Call:  1 call per shot (Flux/Stable Diffusion)
```

**Input:**
- `image_prompt` (từ Phase 1)
- `image_style` (đã pre-distill ở Giai đoạn A) — prepend vào prompt
- Reference images: character images + scene/background + props

**Output:**
- Grid image (action sequence) hoặc single keyframe
- Lưu vào `images` table + `storyboards.composed_image`

### Phase 3: Generate Videos (`processVideo`)

```
Frontend:  videoAPI.generateVideo({ prompt, storyboard_id, reference_image_urls, model, provider })
Backend:   video_generation_service.go → GenerateVideo()
API Call:  1 call per shot (Kling/Doubao/Runway/etc)
```

**Prompt assembly (backend, dòng ~370):**
```go
prompt := videoGen.Prompt  // = video_prompt_distilled || video_prompt || action
if storyboard.VideoStyle != nil {
    prompt = *storyboard.VideoStyle + "\n\n" + prompt
}
// Final: [video_style constraint] + "\n\n" + [video_prompt narrative]
```

**Prompt priority chain (frontend, dòng 549):**
```typescript
prompt: sb.video_prompt_distilled || sb.video_prompt || sb.action || "Cinematic video"
```

| Source | Chất lượng | Khi nào có |
|---|---|---|
| `video_prompt_distilled` | ⭐⭐⭐ AI narrative | Sau BatchDistillStyles (nếu có video_constraint template) |
| `video_prompt` | ⭐⭐ Mechanical | Auto-gen khi tạo/update storyboard |
| `action` | ⭐ Raw | Luôn có (từ storyboard breakdown) |

**Input:**
- Prompt (assembled ở trên)
- Reference image: grid image từ Phase 2 (`image.local_path || image.image_url`)
- Video model + provider (user chọn)
- Duration: 5s default
- Mode: `reference_mode: 'multiple'`

**Output:**
- Video file
- Lưu vào `video_generations` table

---

## 4. Storyboard Data Model — Các trường liên quan

```go
// domain/models/drama.go — Storyboard struct

// === Shot Properties (từ storyboard breakdown) ===
Action           *string  // Hành động chính trong shot
Result           *string  // Kết quả/trạng thái cuối shot
Location         *string  // Bối cảnh
Time             *string  // Thời điểm
ShotType         *string  // Close-up, Medium, Wide...
Angle            *string  // Eye level, Low angle, High angle...
Movement         *string  // Static, Pan, Dolly, Handheld...
Atmosphere       *string  // Mood/atmosphere

// === Audio Context ===
Dialogue         *string  // Thoại nhân vật
NarratorScript   *string  // Lời kể voice-over
SoundEffect      *string  // Hiệu ứng âm thanh
AudioMode        *string  // narrator_only | dialogue_dominant

// === Generated Prompts ===
ImagePrompt      *string  // AI-generated image prompt (Phase 1)
VideoPrompt      *string  // Mechanical video prompt (Pipeline 1 — auto)
VideoPromptSource string  // "auto" | "ai" | "manual"

// === Pre-distilled Fields (Giai đoạn A) ===
ImageStyle       *string  // Distilled từ style_prompt template
VideoStyle       *string  // Distilled từ video_constraint template (constraint)
VideoPromptDistilled *string // Distilled từ video_constraint template (narrative)

// === Generated Assets ===
ComposedImage    *string  // Grid image URL
VideoURL         *string  // Video URL
```

---

## 5. Template System — Prompt Resolution

### Fallback Chain

```
Template override (PromptTemplatePrompts.{field}) → Default embed file (prompts/{file}.txt)
```

Khi drama có `prompt_template_id` → hệ thống kiểm tra template có field tương ứng không.
Nếu có → dùng template content. Nếu không → dùng file mặc định.

### Các prompt types liên quan đến Batch Action

| Prompt Type | Default File | Dùng ở | Mục đích |
|---|---|---|---|
| `style_prompt` | `style_prompt.txt` | Distill → `image_style` | Visual DNA cho image gen |
| `video_constraint` | `video_constraint_prefixes.txt` | Distill → `video_style` + `video_prompt_distilled` | Camera/animation rules + narrative prompt |
| `video_extraction` | `video_extraction.txt` | Pipeline 2 (BatchReferenceVideoStudio) | AI viết video prompt (R2V mode) |
| `image_first_frame` | `image_first_frame.txt` | Phase 1 (keyframe mode) | Image prompt generation |
| `image_action_sequence` | `image_action_sequence.txt` | Phase 1 (action sequence mode) | Grid prompt generation |
| `image_key_frame` | `image_key_frame.txt` | Phase 1 (keyframe mode) | Keyframe prompt generation |

### Distillation Prompt Files (system-level, không override được)

| File | Output | Vai trò |
|---|---|---|
| `image_style_distill.txt` | `image_style` | Prompt hệ thống cho distill image style |
| `video_distill_combined.txt` | `video_style` + `video_prompt_distilled` | Prompt hệ thống cho distill cả constraint lẫn narrative |
| `video_style_distill.txt` | `video_style` | Legacy — chỉ dùng cho Narrative MV part-aware mode |

---

## 6. Chi phí LLM — Ví dụ 10 shots

| Bước | Loại | Requests | Thời điểm |
|---|---|---|---|
| Distill image_style | Batch LLM | **1** | Tự động sau tạo storyboard |
| Distill video_style + video_prompt | Batch LLM | **1** | Tự động sau tạo storyboard |
| Extract image prompts | Per-shot LLM | **10** | User trigger Phase 1 |
| Generate images | Per-shot Image API | **10** | User trigger Phase 2 |
| Generate videos | Per-shot Video API | **10** | User trigger Phase 3 |
| **Tổng** | | **32** | |

---

## 7. Luồng dữ liệu tổng hợp

```
┌─────────────────────────────────────────────────────────────────┐
│ TEMPLATE (user tạo, ví dụ: HL Meow Meow, Cocomelon, CG5)      │
│                                                                 │
│  style_prompt ─────────────┐                                    │
│  video_constraint ─────────┤                                    │
│  image_action_sequence ────┤                                    │
│  image_first_frame ────────┤                                    │
│  ...                       │                                    │
└────────────────────────────┼────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│ GIAI ĐOẠN A: Pre-distillation (auto, background)               │
│                                                                 │
│  style_prompt + shots ──LLM──→ image_style (per shot)           │
│  video_constraint + shots ──LLM──→ video_style (per shot)       │
│                                  + video_prompt_distilled       │
│                                                                 │
│  Lưu vào DB: storyboards.{image_style, video_style,            │
│               video_prompt_distilled}                           │
└─────────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│ GIAI ĐOẠN B: Batch Action Studio (user trigger)                 │
│                                                                 │
│  Phase 1 — Extract Prompts                                      │
│  ┌────────────────────────────────────────────────┐             │
│  │ shot + chars + scenes + props                  │             │
│  │ + template (image_action_sequence/first_frame) │             │
│  │ ──LLM──→ image_prompt                          │             │
│  └────────────────────────────────────────────────┘             │
│                             │                                   │
│  Phase 2 — Generate Images                                      │
│  ┌────────────────────────────────────────────────┐             │
│  │ [image_style] + image_prompt                   │             │
│  │ + reference images (chars, scene, props)       │             │
│  │ ──Image API──→ grid image / keyframe           │             │
│  └────────────────────────────────────────────────┘             │
│                             │                                   │
│  Phase 3 — Generate Videos                                      │
│  ┌────────────────────────────────────────────────┐             │
│  │ [video_style] + [video_prompt_distilled]       │             │
│  │ + reference image (grid từ Phase 2)            │             │
│  │ ──Video API──→ video                           │             │
│  └────────────────────────────────────────────────┘             │
└─────────────────────────────────────────────────────────────────┘
```

---

## 8. Mechanical Video Prompt (Pipeline 1) — Fallback

Khi `video_prompt_distilled` không tồn tại (ví dụ: drama không có template, hoặc distill fail), hệ thống fallback về `video_prompt` cơ học:

```go
// storyboard_service.go → generateVideoPrompt()
// Chạy auto khi tạo/update storyboard

Output format:
"Action: {action}. Result: {result}. Camera movement: {movement}. 
 Shot type: {shot_type}. Camera angle: {angle}. Scene: {location, time}. 
 Atmosphere: {atmosphere}. Dialogue/Narration: {dialogue/narrator}. 
 Sound effects: {sfx}. =VideoRatio: 16:9"
```

**Đặc điểm:**
- Ghép label:value thuần túy, không có narrative
- Tự động inject lip-sync/mouth-closed constraints dựa trên audio_mode
- Inject character voice styles cho dialogue shots
- Source: `video_prompt_source = "auto"`

---

## 9. So sánh với Batch Reference Video Studio

| Đặc điểm | Batch Action Studio | Batch Reference Video Studio |
|---|---|---|
| **Component** | `BatchGenerationDialog.vue` | `BatchReferenceVideoStudio.vue` |
| **Nút mở** | "Batch" button | "R2V" button |
| **Image gen** | ✅ Tạo grid/keyframe | ❌ Không tạo ảnh |
| **Video prompt source** | `video_prompt_distilled` (pre-distilled) | Pipeline 2: AI viết qua `video_extraction` |
| **Video gen mode** | I2V (image-to-video) từ grid | R2V (reference-to-video) từ character refs |
| **Use case** | Production — full pipeline | Quick iteration — chỉ video |

---

## 10. Concurrency & Rate Limiting

- Batch Action Studio sử dụng `runConcurrently(items, limit, worker)`
- `limit` = `maxConcurrentThreads` (cấu hình từ AI Settings)
- Mỗi phase chạy tuần tự (Phase 1 xong → Phase 2 → Phase 3)
- Trong mỗi phase, shots chạy concurrent tối đa `limit` threads
- User có thể bấm "Stop" để dừng giữa chừng (`shouldStop` flag)

---

## 11. Post-processing (Optional)

Sau khi gen video xong, user có thể:

| Action | Chức năng |
|---|---|
| **Upscale All Videos** | Nâng cấp chất lượng video (HD) |
| **Review Videos** | AI review chất lượng video → score + verdict |
| **Download All Videos** | Tải tất cả video thành ZIP |
| **Reset All Data** | Xóa toàn bộ prompt/ảnh/video, reset về trạng thái chờ |

---

## 12. Changelog

| Ngày | Thay đổi |
|---|---|
| 2026-04-27 | Thêm `video_prompt_distilled` — AI narrative prompt thay thế mechanical concatenation |
| 2026-04-27 | Gộp `distillVideoStyles` + `distillVideoPrompts` thành `distillVideoCombined` (tiết kiệm 1 API call) |
| 2026-04-27 | Mở rộng `shotContext` thêm dialogue, narrator, sound_effect, character_voices, audio_mode |
| 2026-04-27 | Frontend ưu tiên `video_prompt_distilled` > `video_prompt` > `action` |
