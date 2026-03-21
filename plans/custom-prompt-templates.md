> Created: 2026-03-21
> Status: Implemented

## Summary

Tách hệ thống prompt hiện tại thành 2 lớp: **Fixed Requirement** (cấu trúc JSON output cố định cho Backend) và **Dynamic Template** (phần nội dung sáng tạo mà người dùng tùy biến theo từng dự án). Người dùng có thể tạo, lưu, và tái sử dụng các Prompt Template (Mẫu Prompt) trên nhiều project khác nhau, biến hệ thống từ "Công cụ làm phim" thành "Cỗ máy sinh Video đa mục đích".

## Problem Statement

Hiện tại, toàn bộ prompt của hệ thống được hardcode trong thư mục `application/prompts/*.txt` và được nhúng trực tiếp vào binary Go thông qua `embed.FS`. Mọi project đều dùng chung một bộ prompt duy nhất hướng đến phong cách "Làm phim điện ảnh". Người dùng không thể:
1. Thay đổi "Vai trò" (Role) của AI cho các loại video khác nhau (Explainer, Vlog, MV...).
2. Tùy chỉnh "Nguyên tắc" (Principles) phân rã kịch bản theo nhu cầu riêng.
3. Tái sử dụng một bộ cấu hình prompt đã tạo cho nhiều project.

## Research Findings

### Codebase Patterns

**Prompt Loading Layer:**
- [prompts.go](file:///g:/VS-Project/huobao-drama/application/prompts/prompts.go): Sử dụng `embed.FS` để load tất cả file `.txt` tại compile-time. Hàm `prompts.Get(name)` trả về string.
- [prompt_i18n.go](file:///g:/VS-Project/huobao-drama/application/services/prompt_i18n.go): Lớp trung gian (middleware) gọi `prompts.Get()` và format với các biến runtime (`style`, `imageRatio`). Đây là **điểm can thiệp (injection point)** chính.

**Prompt Files hiện tại (14 files):**

| File | Loại | Có Format Vars | Ai gọi |
|------|------|----------------|--------|
| `storyboard_story_breakdown.txt` | Dynamic (Role + Rules) | Không | `GetStoryboardSystemPrompt()` |
| `storyboard_format_instructions.txt` | **Fixed** (JSON Schema) | Không | `storyboard_service.go:168` |
| `character_extraction.txt` | Mixed | `%s` (style, ratio) | `GetCharacterExtractionPrompt()` |
| `prop_extraction.txt` | Mixed | `%s` (style, ratio) | `GetPropExtractionPrompt()` |
| `scene_extraction.txt` | Mixed | `%s` (style, ratio) | `GetSceneExtractionPrompt()` |
| `script_outline_generation.txt` | Dynamic (Role + Rules) | Không | `GetOutlineGenerationPrompt()` |
| `script_episode_generation.txt` | Dynamic (Role + Rules) | Không | `GetEpisodeScriptPrompt()` |
| `image_first_frame.txt` | Mixed | `%s` (style, ratio) | `GetFirstFramePrompt()` |
| `image_key_frame.txt` | Mixed | `%s` (style, ratio) | `GetKeyFramePrompt()` |
| `image_last_frame.txt` | Mixed | `%s` (style, ratio) | `GetLastFramePrompt()` |
| `image_action_sequence.txt` | Mixed | `%s` (style, ratio) | `GetActionSequenceFramePrompt()` |
| `video_constraint_prefixes.txt` | Mixed | Không | `GetVideoConstraintPrompt()` |
| `style_prompt.txt` | Dynamic | Không | `GetStylePrompt()` |

**Data Model hiện tại:**
- [drama.go](file:///g:/VS-Project/huobao-drama/domain/models/drama.go): Model `Drama` đã có `Style` và `CustomStyle` nhưng chưa có `PromptTemplateID`.

### Phân tách Fixed vs Dynamic cho từng file

Mỗi file `.txt` hiện tại cần được phân tách thành 2 phần:

**Phần Fixed (Backend khóa cứng, User không thấy):**
- Cấu trúc JSON output bắt buộc (key names, types, format)
- Ràng buộc ngôn ngữ (`CRITICAL LANGUAGE CONSTRAINT`)
- Ràng buộc dữ liệu liên kết (`character IDs`, `scene_id`)

**Phần Dynamic (User tùy biến):**
- Role/Persona của AI
- Principles/Rules cho AI tuân theo
- Hướng dẫn mô tả, ví dụ, gợi ý phong cách

## Proposed Solution

### Architecture Overview

```
┌─────────────────────────────────────────────────────┐
│                   Frontend (Vue.js)                  │
│  ┌─────────────┐  ┌──────────────────────────────┐  │
│  │ Drama Create │  │ Prompt Template Manager      │  │
│  │ (Dropdown    │  │ (CRUD + Preview + Tab UI)    │  │
│  │  chọn tpl)   │  │                              │  │
│  └──────┬──────┘  └──────────────┬───────────────┘  │
└─────────┼────────────────────────┼──────────────────┘
          │                        │
          ▼                        ▼
┌─────────────────────────────────────────────────────┐
│                Backend (Golang)                      │
│                                                      │
│  ┌──────────────────────────────────────────────┐   │
│  │           PromptResolver (MỚI)                │   │
│  │  GetPrompt(dramaID, promptType) string        │   │
│  │                                                │   │
│  │  1. Drama có template_id? ──NO──► prompts.Get()│   │
│  │               │YES                             │   │
│  │  2. Template có override? ──NO──► prompts.Get()│   │
│  │               │YES                             │   │
│  │  3. Return: DynamicTemplate + FixedRequirement │   │
│  └──────────────────────────────────────────────┘   │
│                                                      │
│  ┌────────────────┐  ┌────────────────────────┐     │
│  │ Fixed Prompts   │  │ prompt_templates (DB)  │     │
│  │ (embed.FS)      │  │ id, name, prompts JSON │     │
│  │ Không thay đổi  │  │ User CRUD              │     │
│  └────────────────┘  └────────────────────────┘     │
└─────────────────────────────────────────────────────┘
```

### Approach

#### Phase 1: Database & Model (Backend)

**1.1. Tạo Model `PromptTemplate`:**
```go
// domain/models/prompt_template.go
type PromptTemplate struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Name        string         `gorm:"type:varchar(200);not null" json:"name"`
    Description *string        `gorm:"type:text" json:"description"`
    Prompts     datatypes.JSON `gorm:"type:json" json:"prompts"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**1.2. Cấu trúc JSON của trường `Prompts`:**
```json
{
  "storyboard_breakdown": "You are a Motion Graphics Director...",
  "character_extraction": "You are a character concept designer...",
  "scene_extraction": "...",
  "prop_extraction": "...",
  "script_outline": "...",
  "script_episode": "...",
  "image_first_frame": "...",
  "image_key_frame": "...",
  "image_last_frame": "...",
  "image_action_sequence": "...",
  "video_constraint": "...",
  "style_prompt": "..."
}
```
> Mỗi key là optional. Nếu key không tồn tại hoặc value rỗng → Fallback về file `.txt` mặc định.

**1.3. Thêm cột `prompt_template_id` vào Drama:**
```go
// Thêm vào Drama struct
PromptTemplateID *uint          `gorm:"index" json:"prompt_template_id"`
PromptTemplate   PromptTemplate `gorm:"foreignKey:PromptTemplateID" json:"prompt_template,omitempty"`
```

#### Phase 2: Tách Fixed Requirements

**2.1. Tạo thư mục mới `application/prompts/fixed/`:**

Chứa các file chỉ có phần Fixed (JSON schema + constraints). Các file này **không bao giờ** được expose ra Frontend hoặc cho User sửa.

Danh sách file Fixed cần tách:
- `fixed/storyboard_format.txt` ← Phần JSON schema từ `storyboard_format_instructions.txt`
- `fixed/character_format.txt` ← Phần JSON output format từ `character_extraction.txt`
- `fixed/scene_format.txt` ← Phần JSON output format từ `scene_extraction.txt`
- `fixed/prop_format.txt` ← Phần JSON output format từ `prop_extraction.txt`
- `fixed/outline_format.txt` ← Phần JSON output format từ `script_outline_generation.txt`
- `fixed/episode_format.txt` ← Phần JSON output format từ `script_episode_generation.txt`
- `fixed/image_format.txt` ← Phần JSON output format chung cho image prompts
- `fixed/video_format.txt` ← Constraints cho video

**2.2. Sửa các file `.txt` gốc chỉ giữ phần Dynamic:**
Các file gốc (`storyboard_story_breakdown.txt`, `character_extraction.txt`...) sẽ chỉ giữ lại phần Role + Rules + Guidelines. Phần JSON schema bị cắt sang thư mục `fixed/`.

#### Phase 3: PromptResolver Service (Backend)

**3.1. Tạo `application/services/prompt_resolver.go`:**
```go
type PromptResolver struct {
    db         *gorm.DB
    promptI18n *PromptI18n
    log        *logger.Logger
}

// GetPrompt trả về prompt hoàn chỉnh = DynamicPart + FixedPart
func (r *PromptResolver) GetPrompt(dramaID uint, promptType string) string {
    // 1. Tìm Drama → lấy PromptTemplateID
    // 2. Nếu nil → Fallback: prompts.Get(promptType + ".txt")
    // 3. Nếu có Template → check template.Prompts[promptType]
    //    - Có value → dùng value đó (Dynamic phần User)
    //    - Không có → Fallback: prompts.Get(promptType + ".txt")
    // 4. Nối: DynamicPart + "\n\n" + prompts.GetFixed(promptType)
    // 5. Return chuỗi hoàn chỉnh
}
```

**3.2. Sửa `prompt_i18n.go`:**
Mỗi hàm `Get...Prompt()` sẽ nhận thêm tham số `dramaID` và delegate sang `PromptResolver`. Logic format `%s` (style, imageRatio) vẫn giữ nguyên.

#### Phase 4: CRUD API (Backend)

**4.1. Tạo `application/services/prompt_template_service.go`:**
- `ListTemplates()` → GET `/api/prompt-templates`
- `GetTemplate(id)` → GET `/api/prompt-templates/:id`
- `CreateTemplate(req)` → POST `/api/prompt-templates`
- `UpdateTemplate(id, req)` → PUT `/api/prompt-templates/:id`
- `DeleteTemplate(id)` → DELETE `/api/prompt-templates/:id`
- `GetDefaultPrompts()` → GET `/api/prompt-templates/defaults` (trả về nội dung Dynamic mặc định từ file `.txt` để User tham khảo/copy)

**4.2. Sửa `drama_service.go`:**
- `CreateDramaRequest` / `UpdateDramaRequest` thêm `PromptTemplateID *uint`

#### Phase 5: Frontend UI

**5.1. Trang quản lý Prompt Templates (Menu mới):**

Giao diện gồm:
- **Danh sách Templates** (bảng có Name, Description, ngày tạo, nút Sửa/Xóa/Nhân bản)
- **Dialog Tạo/Sửa Template** với cấu trúc Tab:

| Tab | Prompt Types trong Tab |
|-----|----------------------|
| 📝 Kịch bản | `script_outline`, `script_episode` |
| 🎭 Trích xuất | `character_extraction`, `scene_extraction`, `prop_extraction` |
| 🎬 Phân cảnh | `storyboard_breakdown` |
| 🖼️ Hình ảnh | `image_first_frame`, `image_key_frame`, `image_last_frame`, `image_action_sequence` |
| 🎥 Video | `video_constraint` |
| 🎨 Phong cách | `style_prompt` |

Mỗi tab hiển thị:
- **Textarea** để User nhập nội dung Dynamic
- **Nút "Tải Prompt Mặc Định"** gọi API `/defaults` để fill nội dung gốc vào Textarea
- **Placeholder ẩn** hiển thị nội dung gốc mờ bên dưới để User tham khảo
- **Label ghi chú**: "Để trống = Hệ thống sẽ dùng Prompt mặc định"

**5.2. Sửa Dialog Tạo/Sửa Drama:**
- Thêm Dropdown: **"Mẫu Prompt"** với options: `[Mặc định hệ thống]` + danh sách từ API

**5.3. Frontend Types:**
```typescript
// types/prompt_template.ts
interface PromptTemplate {
  id: number
  name: string
  description?: string
  prompts: Record<string, string> // key → dynamic prompt content
  created_at: string
  updated_at: string
}
```

## Acceptance Criteria

- [ ] Tạo được bảng `prompt_templates` trong Database (Auto-migrate)
- [ ] CRUD API cho Prompt Templates hoạt động đầy đủ
- [ ] API `/defaults` trả về nội dung prompt gốc từ file `.txt`
- [ ] Drama model có thể gán `prompt_template_id`
- [ ] PromptResolver hoạt động theo cơ chế Fallback: Template → Default
- [ ] Fixed Requirements (JSON Schema) luôn được append cuối prompt, User không sửa được
- [ ] Frontend có trang quản lý Templates với Tab UI
- [ ] Frontend Dialog tạo Drama có Dropdown chọn Template
- [ ] Project không gán Template → hệ thống chạy y như cũ (zero regression)
- [ ] Project gán Template → sử dụng Dynamic prompt từ Template
- [ ] Template để trống 1 prompt type → Fallback về file mặc định cho type đó

## Technical Considerations

### Dependencies
- Không cần thêm package Go mới (GORM đã hỗ trợ `datatypes.JSON`)
- Frontend: Không cần thêm thư viện mới (Element Plus đã có Tabs, Textarea, Select)

### Risks
1. **Migration cũ/mới**: Các project cũ không có `prompt_template_id` → cột nullable, Fallback tự động xử lý → **Risk: Thấp**
2. **Prompt format vars (`%s`)**: Một số prompt có `%s` cho `style`, `imageRatio`. Nếu User custom nhưng quên `%s` → `fmt.Sprintf` lỗi → **Giải pháp**: Dùng `strings.Replace()` thay vì `Sprintf` cho custom prompts, hoặc append style/ratio riêng biệt
3. **Prompt quá dài/ngắn**: User có thể nhập prompt bất kỳ → **Giải pháp**: Validate length (min 50, max 10000 chars)

### Alternatives Considered
1. **Lưu prompt trực tiếp trong Drama model**: Bỏ vì không tái sử dụng được
2. **Cho User sửa cả phần Fixed**: Bỏ vì phá vỡ Backend parsing
3. **Dùng file system thay Database**: Bỏ vì khó CRUD và khó scale

## Implementation Steps

### Phase 1: Backend Foundation (Ưu tiên cao)
- Task 1: Tạo model `PromptTemplate` + auto-migrate
- Task 2: Thêm `prompt_template_id` vào model `Drama`
- Task 3: Tách nội dung Fixed vs Dynamic trong các file `.txt` hiện tại
- Task 4: Tạo `PromptResolver` service với cơ chế Fallback
- Task 5: Sửa `prompt_i18n.go` để delegate sang `PromptResolver`

### Phase 2: Backend API
- Task 6: Tạo `PromptTemplateService` (CRUD)
- Task 7: Tạo API routes (`/api/prompt-templates/*`)
- Task 8: Sửa `DramaService` để xử lý `prompt_template_id`
- Task 9: API endpoint `/defaults` trả về nội dung Dynamic mặc định

### Phase 3: Frontend
- Task 10: Tạo TypeScript types cho `PromptTemplate`
- Task 11: Tạo API client cho prompt templates
- Task 12: Tạo trang `PromptTemplateManager.vue` (danh sách + CRUD dialog với Tab UI)
- Task 13: Sửa `CreateDramaDialog.vue` + `DramaList.vue` thêm Dropdown chọn Template
- Task 14: Thêm route + menu navigation

### Phase 4: Testing & Polish
- Task 15: Test Fallback mechanism (Template có/không/partial override)
- Task 16: Test backward compatibility (project cũ không gán Template)
- Task 17: Đảm bảo `%s` format vars hoạt động đúng với custom prompts

## References
- [prompt_i18n.go](file:///g:/VS-Project/huobao-drama/application/services/prompt_i18n.go) - Điểm can thiệp chính
- [prompts.go](file:///g:/VS-Project/huobao-drama/application/prompts/prompts.go) - Hệ thống embed hiện tại
- [drama.go](file:///g:/VS-Project/huobao-drama/domain/models/drama.go) - Model cần mở rộng
- [storyboard_service.go](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go) - Consumer chính của prompts
