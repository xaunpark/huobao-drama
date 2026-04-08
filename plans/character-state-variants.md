# Character State Variants — Flatten Approach

> Created: 2026-04-08
> Status: Implemented

## Summary

Nâng cấp hệ thống trích xuất nhân vật để tự động nhận diện và tạo các **biến thể trạng thái** (state variants) cho mỗi nhân vật khi có sự thay đổi rõ nét về mặt hình ảnh thị giác. Sử dụng Hướng B (Flatten) — tạo Character riêng cho mỗi state bằng naming convention, **không thay đổi schema DB**.

## Problem Statement

Hiện tại mỗi `Character` chỉ có 1 `Appearance` và 1 `ImageURL`. Khi downstream (frame prompt, video prompt) sử dụng character, nó inject cùng 1 mô tả ngoại hình cho tất cả shot — bất kể nhân vật đang ở trạng thái hình ảnh nào.

**Ví dụ thực tế:** "Xe tải không" vs "Xe tải chở đầy gỗ" vs "Xe tải đang được kéo bởi xích từ xe cứu hộ" — đây là những hình ảnh hoàn toàn khác nhau cần reference image riêng.

## Prior Solutions

Không tìm thấy giải pháp hiện có trong `docs/solutions/` hay `docs/explorations/`.

## Research Findings

### Codebase Architecture (Luồng hiện tại)

1. **Character Extraction** (2 entry points parallel):
   - `script_generation_service.go:66-217` — `processCharacterGeneration()` via `GenerateCharacters()`
   - `character_library_service.go:578-687` — `processCharacterExtraction()` via `ExtractCharactersFromScript()`
   
2. **Prompt Chain**:
   - Dynamic: `character_extraction.txt` (9 dòng, rất đơn giản)
   - Fixed: `fixed/character_extraction.txt` (41 dòng, chứa format + examples)

3. **Character Model** (`drama.go:40-64`):
   - `Name`, `Role`, `Appearance`, `Description`, `Personality`, `VoiceStyle`, `ImageURL`
   - Không có `parent_id` hay `state` field

4. **Shot Creation — AI assigns characters by ID**:
   - `storyboard_service.go:142-156`: Truyền `characterList` = `[{"id": 5, "name": "Truck"}, ...]`
   - AI output: `"characters": [5, 12]` — AI tự chọn ID
   - Format constraint: "characters array must contain only valid IDs from the Available Character List"

5. **Downstream Usage** (`frame_prompt_service.go:492-501`):
   - `buildStoryboardContext()` inject `char.Name + (char.Appearance)` vào prompt
   - Cơ chế này sẽ tự động hưởng lợi khi variant có name + appearance riêng

### Tiêu chí tạo State Variant (đã thống nhất)

#### TẠO trạng thái mới khi — có vật thể được thêm/bớt/biến đổi gắn liền:

| Thay đổi | Ví dụ |
|----------|-------|
| Thêm vật thể lớn lên/vào | Xe tải không → xe tải chở đầy gỗ |
| Kết nối với entity khác | Xe tải → xe tải được kéo bởi xích từ xe cứu hộ |
| Trang phục/phụ kiện lớn thay đổi | Chiến binh mặc đồ thường → chiến binh mặc giáp |
| Biến hình rõ ràng | Sâu bướm → bướm |

#### KHÔNG tạo khi — chỉ thay đổi tư thế/vị trí/hướng/biểu cảm:

| Thay đổi | Ví dụ |
|----------|-------|
| Di chuyển / dừng lại | Xe tải đứng → xe tải di chuyển |
| Lật / nghiêng / xoay | Xe tải đứng → xe tải bị lật |
| Biểu cảm | Nhân vật cười → nhân vật khóc |
| Hành động thoáng qua | Nhân vật giơ tay → hạ tay |

**Nguyên tắc vàng:** Nếu cần thêm/bớt pixel đáng kể vào reference image để mô tả đúng → tạo variant. Nếu chỉ xoay/di chuyển/biểu cảm → video prompt xử lý.

## Proposed Solution — Hướng B (Flatten)

### Approach

Tạo Character riêng cho mỗi state, sử dụng naming convention `"BaseName (State Description)"`. Không cần thêm `base_name` field vào output — tên trước ngoặc đã là base name.

**Zero schema change** — chỉ sửa prompt extraction + thêm guidance nhỏ vào shot assignment.

### Downstream tự động hưởng lợi

Khi `buildStoryboardContext()` inject `char.Name + (char.Appearance)`:
- Trước: `"Truck (A large industrial truck with red cab...)"` — 1 mô tả cho mọi shot
- Sau: `"Truck (Loaded with Logs) (A large red truck with timber logs stacked high...)"` — mô tả chính xác cho từng shot

## Acceptance Criteria

- [ ] AI trích xuất character tạo ra các state variant khi có sự thay đổi thị giác rõ rệt
- [ ] Mỗi variant có `name`, `appearance`, `description` riêng biệt
- [ ] AI khi tạo shot chọn đúng variant ID phù hợp với trạng thái trong shot đó
- [ ] Không có thay đổi schema DB
- [ ] Không ảnh hưởng đến các tính năng khác (image generation, storyboard split quality)
- [ ] Prompt có ví dụ cụ thể đa domain (phương tiện, con người, động vật, đồ vật)

## Technical Considerations

### Dependencies
- Không cần thêm package/library mới
- Không cần migration DB

### Risks
- **Token budget risk**: Prompt extraction dài hơn → có thể tốn thêm token. Mitigation: chỉ thêm vừa đủ, tập trung ví dụ
- **Over-splitting risk**: AI có thể tạo quá nhiều variant vụn vặt. Mitigation: prompt có ví dụ rõ ràng về khi nào KHÔNG tạo variant
- **Shot assignment dilution risk**: Thêm guidance vào shot prompt có thể làm loãng các split rules. Mitigation: chỉ thêm 5-8 dòng vào section Entity Requirements hiện có, không tạo section mới

### Alternatives Considered
- **Hướng A (Character States table)**: Schema change + migration + sửa nhiều code downstream. Rejected vì quá phức tạp cho mức độ cần thiết hiện tại.
- **Thêm `base_name` field**: Không cần vì naming convention đã đủ rõ. Nếu cần group sau này, parse tên bằng regex là đủ.

## Implementation Steps

### Task 1: Sửa prompt trích xuất nhân vật (SỬA LỚN)

**File**: `application/prompts/fixed/character_extraction.txt`

**Thay đổi:**
- Thêm section `[State Variant Rules]` với tiêu chí rõ ràng
- Thêm ví dụ cụ thể multi-domain (phương tiện, con người, động vật)
- Thêm ví dụ output hoàn chỉnh cho character có variant và không có variant
- Sửa naming convention: `"BaseName (State)"` pattern

### Task 2: Sửa dynamic prompt (SỬA NHỎ)

**File**: `application/prompts/character_extraction.txt`

**Thay đổi:** Thêm 1-2 dòng nhắc AI xem xét state variants.

### Task 3: Sửa prompt tạo shot — Entity Requirements (SỬA NHỎ)

**Files** (3 files cần sửa cùng nội dung):
- `application/prompts/storyboard_visual_unit_format.txt` (dòng 61-64, Entity Requirements)
- `application/prompts/storyboard_format_instructions.txt` (dòng 54-57, Entity Requirements)
- `application/prompts/storyboard_visual_unit_structured.txt` (dòng 46, characters matching)

**Thay đổi:** Thêm 5-8 dòng vào phần Entity Requirements hiện có.

**Không** tạo section mới, **không** thêm split rules, **không** ảnh hưởng logic chia shot.

### Task 4: Kiểm tra — không cần sửa service code

Xác nhận rằng **không cần sửa code Go** vì:
- `processCharacterGeneration()` đã parse JSON array → create Character records → không cần biết state variant
- `buildStoryboardContext()` đã dùng `char.Name + char.Appearance` → variant tự inject đúng
- Character image generation (`buildCharacterPrompt()`) dùng `Appearance` field → variant có appearance riêng → sinh ảnh đúng

## References

- Character model: `domain/models/drama.go:40-64`
- Fixed extraction prompt: `application/prompts/fixed/character_extraction.txt`
- Dynamic extraction prompt: `application/prompts/character_extraction.txt`
- Shot format (visual unit): `application/prompts/storyboard_visual_unit_format.txt`
- Shot format (breakdown): `application/prompts/storyboard_format_instructions.txt`
- Shot format (structured): `application/prompts/storyboard_visual_unit_structured.txt`
- Character extraction service: `application/services/script_generation_service.go`
- Frame prompt (downstream): `application/services/frame_prompt_service.go:492-501`
