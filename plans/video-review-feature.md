# Video Review Feature — Batch Action Studio

> Created: 2026-04-20
> Status: Implemented ✓

## Summary

Thêm nút "Review Videos" vào Batch Action Studio, cho phép chấm điểm chất lượng video AI bằng Gemini Vision (qua Locally-AI). Backend extract frames → tạo Contact Sheet → gửi cho Gemini → parse điểm số → lưu DB. Điểm gắn với `video_gen_id` cụ thể, tự mất khi video thay đổi.

## Problem Statement

Sau khi batch generate video, không có cách tự động đánh giá chất lượng video. User phải mở từng video xem bằng mắt → tốn thời gian và dễ bỏ sót lỗi (character drift, reverse motion, brand logos...).

## Prior Solutions

Không tìm thấy solution nào liên quan trong `docs/solutions/`.

## Research Findings

### Codebase Patterns

- **Async Task pattern**: `application/services/task_service.go` — `CreateTask` → goroutine → `UpdateTaskResult`. Frontend poll via `taskAPI.getStatus()`.
- **Batch Operations UI**: `web/src/components/editor/BatchGenerationDialog.vue` — `startUpscaleAll` (line 659) là pattern mẫu: iterate shots, filter by state, call API, poll status.
- **AI Client**: `pkg/ai/openai_client.go` — `ChatMessage.Content` = `string`. Không thay đổi struct này.
- **Video Model**: `domain/models/video_generation.go` — `VideoGeneration` with `LocalPath`, `VideoURL`, `StoryboardID`.
- **Routes**: `api/routes/routes.go` (line 185-195) — video routes under `api.Group("/videos")`.
- **DB Migration**: `infrastructure/database/database.go` (line 68-98) — `AutoMigrate` list.

### External Reference — FlowKit Video Review

- Source: `docs/video_review_export/` (5 files)
- Rubric: 6 dimensions × weighted scores → overall 0-10
- Error catalog: 14 error types (5 CRITICAL, 5 HIGH, 4 MINOR)
- Contact sheet approach: ffmpeg `tile` filter, 320px thumbnails with timestamp overlay

### Locally-AI Multimodal API

- Source: `docs/API-Locally-image.md`
- Format: `content: [{type: "text", text: ...}, {type: "image_url", image_url: {url: "file:///..."}}]`
- Model: `gemini/auto`
- Image pasted via ClipboardEvent into Gemini UI → 1 image per request is reliable

## Proposed Solution

### Architecture

```
Frontend (Vue)                Backend (Go)                    External
─────────────                ──────────────                  ─────────
POST /videos/:id/review  →  CreateTask("video_review")  
                             goroutine:                    
                               1. Get video local_path
                               2. ffmpeg → Contact Sheet    ← FFmpeg subprocess
                               3. Build multimodal request
                               4. POST /v1/chat/completions → Locally-AI (Gemini)
                               5. Parse JSON response
                               6. INSERT video_reviews
                               7. UpdateTaskResult

GET /tasks/:id (poll)    →  Return task status + result
GET /videos/:id/review   →  Return latest VideoReview
```

### Key Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| AI Provider | Gemini via Locally-AI | User requirement; $0 cost |
| Image format | Contact Sheet (1 grid image) | Locally-AI pastes images into Gemini UI; 1 image = reliable, 40 images = crash |
| Concurrency | 4 parallel requests | User confirmed Locally supports multiple workers |
| Frame rate | 8fps fixed | User requirement; ~40 frames per 5s video |
| Score persistence | FK to `video_gen_id` | Score auto-invalidates when video changes |
| Struct approach | New `MultimodalChatMessage` | Avoid changing existing `ChatMessage` (string Content) |
| Re-review | Allowed | POST creates new review; GET returns latest |

## Acceptance Criteria

- [ ] Nút "Review Videos" hiển thị trong Batch Action Studio
- [ ] Click nút → duyệt qua shots có video Ready → trigger review
- [ ] Mỗi shot hiện điểm số (0-10) + verdict (excellent/good/acceptable/poor/unusable) trong cột "Score"
- [ ] Điểm được lưu vào DB gắn với `video_gen_id`
- [ ] Khi mở lại dialog → điểm cũ vẫn hiển thị (nếu video chưa thay đổi)
- [ ] Khi video đổi (xóa/tạo mới) → điểm hiển thị "—"
- [ ] Review chạy async (poll status), 4 shot song song
- [ ] Contact Sheet được tạo chính xác: 8fps, 8 cột, timestamp overlay

## Technical Considerations

### Dependencies
- `ffmpeg` + `ffprobe` (đã có trên server)
- Locally-AI server phải đang chạy + Gemini worker online

### Risks
- Locally-AI offline → review fails. Mitigation: check `/health` trước, fail fast với message rõ.
- Gemini trả text thay vì JSON → parse fail. Mitigation: robust JSON extraction (tìm `{...}` trong response).
- Contact Sheet quá lớn (video dài) → Gemini timeout. Mitigation: giới hạn max frames.

### Alternatives Considered
- **Gửi từng frame riêng**: Rejected — Locally-AI paste 40 ảnh vào UI sẽ crash.
- **Sửa `ChatMessage.Content` thành `interface{}`**: Rejected — breaking change cho toàn bộ text generation code.
- **Claude Vision (Anthropic)**: Rejected — User yêu cầu chỉ dùng Gemini.

---

## Implementation Steps

### Phase 1: GORM Model + Migration (~30m)

**Files to create/modify:**
- `domain/models/video_review.go` — NEW
- `infrastructure/database/database.go` — ADD to AutoMigrate

```go
// domain/models/video_review.go
type VideoReview struct {
    ID           uint           `gorm:"primarykey" json:"id"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

    VideoGenID   uint    `gorm:"not null;index" json:"video_gen_id"`
    StoryboardID *uint   `gorm:"index" json:"storyboard_id,omitempty"`

    OverallScore float64 `gorm:"not null" json:"overall_score"`
    Verdict      string  `gorm:"type:varchar(20);not null" json:"verdict"`
    Dimensions   string  `gorm:"type:text" json:"dimensions"`
    Errors       string  `gorm:"type:text" json:"errors"`
    FixGuide     string  `gorm:"type:text" json:"fix_guide"`

    FramesAnalyzed int     `json:"frames_analyzed"`
    FPSUsed        float64 `json:"fps_used"`
    HasCritical    bool    `gorm:"default:false" json:"has_critical_errors"`
}

func (VideoReview) TableName() string { return "video_reviews" }
```

---

### Phase 2: Contact Sheet Generator (~1h)

**Files to modify:**
- `infrastructure/external/ffmpeg/ffmpeg.go` — ADD method `CreateContactSheet`

> ⚠️ **Review Finding**: Đã có sẵn package `infrastructure/external/ffmpeg/ffmpeg.go` với `GetVideoDuration()`, `getVideoResolution()`, `downloadVideo()` etc. KHÔNG tạo file mới ở `services/`. Thêm method vào struct `FFmpeg` hiện tại để reuse code.

> ⚠️ **Review Finding**: Code ffmpeg hiện tại dùng `-hwaccel cuda` + `-c:v h264_nvenc`. Contact sheet generation chỉ xuất ảnh JPG → KHÔNG cần GPU flags. Nhưng cần lưu ý nếu ffmpeg trên server không có font cho `drawtext` filter → cần fallback.

Core logic:
```go
func (f *FFmpeg) CreateContactSheet(videoPath string, fps float64, outDir string) (string, int, error) {
    // 1. f.GetVideoDuration(videoPath) — reuse existing method
    // 2. totalFrames = duration * fps
    // 3. rows = ceil(totalFrames / 8)
    // 4. ffmpeg -i video -vf "fps=8,scale=320:-1,drawtext=timestamp,tile=8xN" -q:v 2 output.jpg
    //    (NO hwaccel/nvenc — this is image output only)
    // 5. return outputPath, totalFrames, nil
}
```

---

### Phase 3: Multimodal Locally-AI Client (~1h)

**Files to create:**
- `pkg/ai/multimodal.go` — NEW

Separate from existing `ChatMessage` to avoid breaking changes:

```go
// pkg/ai/multimodal.go
type ContentPart struct {
    Type     string    `json:"type"`
    Text     string    `json:"text,omitempty"`
    ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
    URL string `json:"url"`
}

type MultimodalMessage struct {
    Role    string        `json:"role"`
    Content []ContentPart `json:"content"`
}

type MultimodalRequest struct {
    Model    string              `json:"model"`
    Messages []MultimodalMessage `json:"messages"`
}

// SendMultimodal sends a multimodal request to an OpenAI-compatible endpoint
func SendMultimodal(baseURL, apiKey string, req *MultimodalRequest) (string, error) {
    // Direct HTTP POST, bypassing existing ChatCompletion()
    // Parse response.choices[0].message.content
}
```

---

### Phase 4: VideoReviewService + API Endpoints (~3-4h)

**Files to create:**
- `application/services/video_review_service.go` — NEW
- `api/handlers/video_review.go` — NEW
- `application/prompts/video_review_rubric.txt` — NEW (review prompt)

**Files to modify:**
- `api/routes/routes.go` — ADD 2 routes

#### Service Logic (pseudo-code):

```go
func (s *VideoReviewService) ReviewVideo(videoGenID uint, taskID string) {
    // 1. Fetch video + storyboard info
    // 2. Resolve video path (local_path or download from minio_url)
    // 3. Create temp dir
    // 4. CreateContactSheet(videoPath, 8.0, tmpDir)
    // 5. Build review prompt from rubric template
    // 6. Get Locally-AI config from AIService
    // 7. SendMultimodal(baseURL, apiKey, request with contact sheet file:// URL)
    // 8. Parse JSON from response
    // 9. Calculate overall_score = weighted sum of dimensions
    // 10. Determine verdict
    // 11. Save VideoReview to DB
    // 12. UpdateTaskResult(taskID, review)
}
```

#### API Handler:

```go
// POST /videos/:id/review — Trigger async review
func (h *VideoReviewHandler) ReviewVideo(c *gin.Context) {
    videoGenID := parseID(c)
    task := taskService.CreateTask("video_review", videoGenID)
    go reviewService.ReviewVideo(videoGenID, task.ID)
    response.Success(c, gin.H{"task_id": task.ID})
}

// GET /videos/:id/review — Get latest review result
func (h *VideoReviewHandler) GetVideoReview(c *gin.Context) {
    videoGenID := parseID(c)
    review := db.Where("video_gen_id = ?", videoGenID).Order("created_at DESC").First()
    response.Success(c, review)
}
```

#### Handler constructor needs dependencies:

```go
// In routes.go — VideoReviewHandler needs: db, taskService, aiService, ffmpeg, log
videoReviewHandler := handlers.NewVideoReviewHandler(db, log, aiService)
```

#### Routes (add to `routes.go`):

```go
videos.POST("/:id/review", videoReviewHandler.ReviewVideo)
videos.GET("/:id/review", videoReviewHandler.GetVideoReview)
```

#### Review Prompt (`application/prompts/video_review_rubric.txt`):

Adapted from FlowKit's battle-tested rubric:
- 6 scoring dimensions with weights
- 14 error types (CRITICAL/HIGH/MINOR)
- "You are viewing a contact sheet grid of N frames at 8fps with timestamps"
- JSON-only output requirement

---

### Phase 5: Frontend Changes (~2-3h)

**Files to modify:**
- `web/src/api/video.ts` — ADD 2 API methods
- `web/src/components/editor/BatchGenerationDialog.vue` — ADD button, column, logic

#### 5.1 API Methods

```typescript
// web/src/api/video.ts
reviewVideo(videoGenId: number) {
  return request.post<{ task_id: string }>(`/videos/${videoGenId}/review`)
},

getVideoReview(videoGenId: number) {
  return request.get<any>(`/videos/${videoGenId}/review`)
},
```

#### 5.2 BatchGenerationDialog Changes

**Template additions:**
1. Button "Review Videos" in `action-row` (after Upscale All)
2. Column "Score" in table (after Video column)

**Script additions:**
1. `reviewScores` reactive state
2. `startReviewAll()` — clone pattern from `startUpscaleAll`:
   - Filter shots with video state `done` or `hd`
   - Use `runConcurrently` with concurrency **4**
   - Skip if already has review for current `videoId`
   - Call `videoAPI.reviewVideo(videoId)`
   - Poll `taskAPI.getStatus(task_id)` every 3s
   - Store result in `reviewScores[sb.id]`
3. Verdict display helpers: `getVerdictTagType()`, `getReviewScore()`
4. In `initTaskStates`: fetch existing reviews per `activeVideoId`

---

### Phase 6: Integration Testing (~1h)

1. Start Locally-AI + Gemini worker
2. Generate video for 1 shot
3. Trigger single review via API: `POST /api/v1/videos/{id}/review`
4. Verify contact sheet created correctly
5. Verify Gemini response parsed to valid JSON
6. Verify score saved to DB
7. Test batch review from UI: click "Review Videos"
8. Verify scores persist on dialog reopen
9. Delete video → verify score shows "—"
10. Generate new video → verify old score gone

---

## File Inventory (New/Modified)

| File | Action | Description |
|------|--------|-------------|
| `domain/models/video_review.go` | **CREATE** | GORM model |
| `infrastructure/database/database.go` | MODIFY | Add to AutoMigrate |
| `infrastructure/external/ffmpeg/ffmpeg.go` | MODIFY | Add `CreateContactSheet` method |
| `pkg/ai/multimodal.go` | **CREATE** | Multimodal request types + sender |
| `application/services/video_review_service.go` | **CREATE** | Core review logic |
| `application/prompts/video_review_rubric.txt` | **CREATE** | Review prompt |
| `api/handlers/video_review.go` | **CREATE** | API handlers |
| `api/routes/routes.go` | MODIFY | Add 2 routes |
| `web/src/api/video.ts` | MODIFY | Add 2 API methods |
| `web/src/components/editor/BatchGenerationDialog.vue` | MODIFY | UI changes |

## References

- FlowKit review system: `docs/video_review_export/`
- Locally-AI multimodal API: `docs/API-Locally-image.md`
