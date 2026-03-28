# Rapid Cut Mode — Chế Độ Chia Cảnh Nhanh

> Created: 2026-03-28
> Status: Implemented (Full)

## Summary

Thêm chế độ "Rapid Cut" cho phép gộp 2-3 shot liên tiếp (từ Editorial storyboard gốc) thành 1 "production shot" duy nhất. Production shot này có **cấu trúc dữ liệu tương tự shot gốc** nhưng action field chứa multi-beat description, và ảnh 1x3 grid đại diện cho 3 micro-shots thay vì 3 phases của 1 action. Mục tiêu: tạo video 8s chứa 2-3 cảnh nhịp nhanh (~2.5-4s mỗi cảnh) thay vì 1 action kéo dài 8s.

## Problem Statement

- AI Video tools luôn tạo video 8s cố định, nhưng hầu hết action chỉ cần 2-4s
- Trong phim thực tế, mỗi cut cảnh chỉ ~3-5s, không phải 8s
- Test đã chứng minh AI video xử lý hoàn hảo multi-scene transitions trong 1 video 8s nếu prompt + ảnh 1x3 grid đủ tốt (xem `docs/thao-luan-chatgpt.txt`)

## Prior Solutions

Không tìm thấy solution liên quan trong `docs/solutions/`.

## Research Findings

### Codebase Patterns

#### Data Model (giữ nguyên)
- [Storyboard struct](file:///g:/VS-Project/huobao-drama/domain/models/drama.go#L94-L130): 25 trường, bao gồm relationships `Characters`, `Props`, `Background`
- Sử dụng GORM AutoMigrate — thêm trường mới tự động tạo cột

#### Pipeline hiện tại (giữ nguyên cho Standard mode)
1. **Storyboard Generation** ([storyboard_service.go:64](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L64)): AI chia script → 30 shots
2. **Image Generation** ([image_generation_service.go:749](file:///g:/VS-Project/huobao-drama/application/services/image_generation_service.go#L749)): Batch tạo ảnh cho mỗi shot
3. **Video Generation** ([video_generation_service.go:729](file:///g:/VS-Project/huobao-drama/application/services/video_generation_service.go#L729)): Batch tạo video từ ảnh

#### Prompt Files
- `storyboard_story_breakdown.txt` — System prompt cho AI chia shot (1 action = 1 shot)
- `storyboard_format_instructions.txt` — JSON output format cho shots
- `image_action_sequence.txt` — Prompt tạo ảnh 1x3 strip (Start→Peak→End **cho 1 action**)
- `video_constraint_prefixes.txt` — Constraint cho AI video (maintain SINGLE shot consistency)

#### API Endpoints liên quan
- `POST /api/v1/episodes/:episode_id/storyboards` — Generate storyboard
- `POST /api/v1/images/episode/:episode_id/batch` — Batch generate images
- `POST /api/v1/videos/episode/:episode_id/batch` — Batch generate videos
- `POST /api/v1/storyboards` — Create single storyboard
- `PUT /api/v1/storyboards/:id` — Update storyboard
- `DELETE /api/v1/storyboards/:id` — Delete storyboard

## Proposed Solution

### Approach: "Gộp shot — giữ nguyên cấu trúc"

Khi user chọn Rapid Cut:
1. AI nhận danh sách 30 shots gốc → quyết định gộp shots nào
2. Tạo ~10-15 **production shots mới** trong bảng `storyboards`
3. Production shots có **CÙNG cấu trúc dữ liệu** với shot gốc
4. Nhưng `action` field chứa multi-beat format, `characters`/`props`/`scene_id` được merge
5. 30 shots gốc **giữ nguyên**, chỉ được đánh dấu `is_production = false`
6. User làm việc trên Production Board với shots mới — **pipeline ảnh/video hoạt động y hệt**

### Thay đổi Data Model

Thêm 3 cột vào bảng `storyboards` (GORM AutoMigrate tự xử lý):

```go
type Storyboard struct {
    // ... existing fields ...
    
    // Rapid Cut fields (mới)
    IsProduction   bool            `gorm:"default:false" json:"is_production"`           // true = production shot (rapid cut result)
    PacingMode     *string         `gorm:"size:20" json:"pacing_mode"`                   // "standard" | "rapid_cut" 
    SourceShotIDs  datatypes.JSON  `gorm:"type:json" json:"source_shot_ids"`             // [1,2,3] — IDs of original shots merged into this
}
```

### Luồng xử lý

```
┌──────────────────────────────────────────────────────────┐
│ Editorial Phase (KHÔNG ĐỔI)                              │
│ User tạo storyboard → AI chia 30 shots                   │
│ Shots lưu DB: is_production=false, pacing_mode=null      │
└──────────────────────────────────────────────────────────┘
                        │
                        ▼
┌──────────────────────────────────────────────────────────┐
│ User chọn Rapid Cut Mode (MỚI)                           │
│                                                          │
│ 1. API call: POST /episodes/:id/rapid-cut                │
│ 2. Backend lấy 30 shots gốc → gửi cho AI                │
│ 3. AI quyết định gộp:                                    │
│    - Shots 1,2,3 → Production Shot 1                     │
│    - Shots 4,5 → Production Shot 2 (2 beats + padding)   │
│    - Shot 6 → Production Shot 3 (giữ solo — complex)     │
│    - ...                                                  │
│ 4. Tạo 10-15 production shots mới trong DB:              │
│    - is_production=true                                   │
│    - pacing_mode="rapid_cut"                              │
│    - source_shot_ids=[1,2,3]                              │
│    - action = multi-beat merged description               │
│    - characters = union(shot1.chars, shot2.chars, ...)    │
│    - props = union(shot1.props, shot2.props, ...)         │
│    - scene_id = primary scene hoặc null                   │
│ 5. Trả về danh sách production shots                     │
└──────────────────────────────────────────────────────────┘
                        │
                        ▼
┌──────────────────────────────────────────────────────────┐
│ Production Board (UI hiển thị)                           │
│ User thấy 10-15 shots thay vì 30                         │
│ Mỗi shot hiển thị "3 beats" description                  │
│                                                          │
│ Từ đây: pipeline HOÀN TOÀN GIỐNG Standard mode:         │
│ - Generate ảnh 1x3 (nhưng prompt khác — multi-scene)     │
│ - Generate video 8s từ ảnh                               │
│ - Batch generate cũng hoạt động bình thường              │
└──────────────────────────────────────────────────────────┘
```

### API Queries

**1. Khi UI query storyboards cho 1 episode:**
```
GET /api/v1/episodes/:episode_id/storyboards?view=production
```
- `view=editorial` (default): trả về tất cả shots có `is_production=false` hoặc `pacing_mode IS NULL`
- `view=production`: trả về shots có `is_production=true` (nếu có), fallback về editorial

**2. Khi user bấm "Rapid Cut":**
```
POST /api/v1/episodes/:episode_id/rapid-cut
Body: { "model": "optional-model-name" }
Response: { "task_id": "..." } (async như storyboard generation)
```

**3. Khi user muốn quay về Standard:**
- Xóa các production shots (`is_production=true`)
- UI query lại với `view=editorial`

### Prompt thay đổi

#### Prompt mới: `rapid_cut_merge.txt` (cho bước gộp shot)

```
[Role] You are a senior film editor specializing in fast-paced cinematic editing.

[Task] Given a list of storyboard shots, merge adjacent short shots into "rapid cut" 
production units. Each production unit should contain 2-3 sequential beats that fit 
within an 8-second video.

[Rules]
1. Group 2-3 adjacent shots into one production unit
2. Only merge shots that are SHORT (duration <= 5s each)
3. Keep LONG shots solo (duration > 6s, complex dialogue, emotional peaks)
4. Total duration of merged beats should be ~8 seconds
5. Output the merged action as multi-beat format:
   "BEAT 1: [action from shot 1]. BEAT 2: [action from shot 2]. BEAT 3: [action from shot 3]"
6. Merge characters, props, and scenes from all source shots
7. Combine shot_type, angle, movement for each beat

[Output Format]
[
  {
    "source_shot_ids": [1, 2, 3],
    "title": "From Solitude to Running Free",
    "action": "BEAT 1: Mai sits alone in café, staring at rain. BEAT 2: Bird's-eye view of busy city street. BEAT 3: Mai runs through rain, laughing.",
    "result": "Mai transforms from isolation to joyful freedom in the rain.",
    "shot_type": "Medium Shot → Bird's Eye → Tracking Shot",
    "angle": "Eye-level → Top-down → Low-angle",
    "movement": "Fixed → Fixed → Follow",
    "location": "Café interior → City street → Street level",
    "time": "Rainy afternoon",
    "atmosphere": "Melancholic → Urban bustle → Energetic liberation",
    "emotion": "Loneliness↓ → Neutrality→ → Joy↑↑",
    "duration": 8,
    "dialogue": "",
    "bgm_prompt": "Transition from soft piano to upbeat orchestral",
    "sound_effect": "Rain on window, city traffic, splashing footsteps",
    "characters": [1],
    "scene_id": null,
    "is_primary": true
  }
]
```

#### Prompt variant: `image_action_sequence_rapid_cut.txt` (cho tạo ảnh)

Thay đổi so với `image_action_sequence.txt`:
- Panel 1/2/3 = **3 micro-shots khác nhau** (thay vì Start/Peak/End của 1 action)
- Cho phép thay đổi nhân vật, bối cảnh, góc quay giữa các panel
- Nhấn mạnh "strong visual contrast between panels" và "smooth transition flow"

#### Prompt variant: `video_constraint_rapid_cut.txt` (cho tạo video)

Thay đổi so với `video_constraint_prefixes.txt`:
- Thay "SINGLE shot" → "3 scene transitions within 8 seconds"
- Cho phép "noticeable but smooth transitions between scenes"
- Cho phép thay đổi character, environment, camera giữa các beat

### Cách pipeline ảnh/video hoạt động không cần sửa

Vì production shots có **cùng cấu trúc dữ liệu** với shot gốc:

1. **Image Generation**: `BatchGenerateImagesForEpisode` query `storyboards` → dùng `image_prompt` → **hoạt động y hệt**
   - Chỉ cần đảm bảo query đúng shots (production nếu có rapid cut)
   
2. **Video Generation**: `BatchGenerateVideosForEpisode` query storyboards → tìm completed images → tạo video → **hoạt động y hệt**

3. **Frame Prompt Generation**: `GenerateFramePrompt` nhận storyboard ID → **hoạt động y hệt**
   - Nhưng cần detect `pacing_mode="rapid_cut"` để dùng prompt variant

### Detect rapid cut trong prompt generation

Tại `FramePromptService` và `ImageGenerationService`, khi tạo prompt cho storyboard:
- Kiểm tra `storyboard.PacingMode`
- Nếu `"rapid_cut"` → dùng `image_action_sequence_rapid_cut.txt`
- Nếu `null` hoặc `"standard"` → dùng `image_action_sequence.txt` (hiện tại)

Tương tự cho video constraint:
- Kiểm tra image's storyboard `PacingMode`  
- Nếu `"rapid_cut"` → dùng `video_constraint_rapid_cut.txt`
- Nếu `null` hoặc `"standard"` → dùng `video_constraint_prefixes.txt` (hiện tại)

## Implementation Steps

### Phase 1: Backend — Data Model + Rapid Cut Service (core)

- **Task 1.1**: Thêm 3 trường mới vào `Storyboard` model trong [drama.go](file:///g:/VS-Project/huobao-drama/domain/models/drama.go#L94)
  - `IsProduction`, `PacingMode`, `SourceShotIDs`
  - GORM AutoMigrate sẽ tự thêm cột

- **Task 1.2**: Tạo prompt file `rapid_cut_merge.txt` trong `application/prompts/`
  - AI prompt hướng dẫn gộp shots thành production units

- **Task 1.3**: Tạo `RapidCutService` trong `application/services/rapid_cut_service.go`
  - Method `GenerateRapidCut(episodeID, model)` — async, trả về task_id
  - Logic: query 30 shots gốc → tạo prompt → AI gộp → parse JSON → tạo production shots
  - Merge characters, props, scenes từ source shots
  - Tạo `image_prompt` và `video_prompt` cho production shots

- **Task 1.4**: Tạo handler + route cho Rapid Cut API
  - `POST /api/v1/episodes/:episode_id/rapid-cut` → generate rapid cut
  - `DELETE /api/v1/episodes/:episode_id/rapid-cut` → xóa production shots, quay về standard
  - Cập nhật `GET /api/v1/episodes/:episode_id/storyboards` → hỗ trợ query param `view=production|editorial`

### Phase 2: Backend — Prompt Variants

- **Task 2.1**: Tạo `image_action_sequence_rapid_cut.txt`
  - 3 panel = 3 micro-shots (khác nhân vật/bối cảnh/góc quay OK)

- **Task 2.2**: Tạo `video_constraint_rapid_cut.txt`
  - Cho phép scene transitions thay vì enforce single-shot consistency

- **Task 2.3**: Cập nhật `FramePromptService` để detect `pacing_mode` và chọn đúng prompt variant

- **Task 2.4**: Cập nhật `VideoGenerationService` để detect rapid cut mode và dùng đúng video constraint

- **Task 2.5**: Register prompt types mới trong `PromptI18n` và `PromptTemplate` system

### Phase 3: Frontend — UI

- **Task 3.1**: Thêm nút "Rapid Cut" sau bước Editorial trong `EpisodeWorkflow.vue` hoặc `ProfessionalEditor.vue`
  - Nút xuất hiện sau khi storyboard đã generate
  - Khi bấm → gọi API → hiện loading → hiện Production Board

- **Task 3.2**: Cập nhật storyboard list query để hỗ trợ `view` parameter
  - Khi ở Rapid Cut mode → query `view=production`
  - Hiển thị badge "Rapid Cut" trên mỗi production shot
  - Hiển thị "3 beats" breakdown bên trong mỗi shot card

- **Task 3.3**: Thêm nút "Back to Standard" để xóa production shots

### Phase 4: Testing & Polish

- **Task 4.1**: Test end-to-end: Standard mode vẫn hoạt động 100%
- **Task 4.2**: Test Rapid Cut: generate → review → batch image → batch video
- **Task 4.3**: Test chuyển đổi: Standard → Rapid Cut → Standard
- **Task 4.4**: Test edge cases: episode không có shots, shots rất dài, episode đã có video

## Acceptance Criteria

- [ ] Standard mode hoạt động 100% không thay đổi
- [ ] User có thể chọn Rapid Cut sau khi editorial storyboard xong
- [ ] AI gộp 30 shots thành 10-15 production shots hợp lý
- [ ] Production shots có đầy đủ characters, props, scenes merged
- [ ] Batch generate images cho production shots dùng đúng rapid cut prompt
- [ ] Batch generate videos cho production shots dùng đúng rapid cut constraint
- [ ] User có thể quay về Standard mode (xóa production shots)
- [ ] Không mất dữ liệu editorial shots trong bất kỳ trường hợp nào

## Technical Considerations

### Dependencies
- Không thêm dependency mới — chỉ dùng GORM, existing AI clients

### Risks
- **AI merge quality**: AI có thể gộp shots không hợp lý → cần review & cho phép user manual adjust
- **Prompt engineering**: Rapid cut prompt variants cần tuning nhiều lần
- **Data integrity**: Production shots phải reference đúng source shots

### Alternatives Considered

1. **Separate VideoUnit table**: Rejected — thêm entity mới tăng complexity, pipeline hiện tại phải sửa nhiều
2. **Direct micro-shot generation**: Rejected — mất granularity, không thể chuyển mode
3. **Two-pass merge**: Rejected — 2-3 AI calls, merge logic phức tạp, kết quả không tối ưu

## References

- ChatGPT discussion proving multi-scene AI video works: [thao-luan-chatgpt.txt](file:///g:/VS-Project/huobao-drama/docs/thao-luan-chatgpt.txt)
- Previous analysis: Conversation `488ccd78` artifacts `shot_splitting_analysis.md`, `rapid_cut_architecture.md`
