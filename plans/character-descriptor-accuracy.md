# Character Descriptor Accuracy — Variant Prompt & Episode Descriptor Redesign

> Created: 2026-04-24
> Status: Implemented ✓ (2026-04-24)
> Parent: plans/episode-character-variant.md (Phase 2 extension)

## Summary

Hai vấn đề tồn tại sau khi hoàn thành Phase 1 của character variant system:
1. `variant_prompt` chỉ được tạo cho character có nhiều trạng thái (state variants), trong khi mọi character đều cần nó để hỗ trợ I2I generation từ ảnh tham chiếu của người dùng.
2. `episode_descriptor` được derive từ `character_prompt`/`appearance` (mô tả nhân vật trừu tượng do AI tưởng tượng), không phải từ `variant_prompt` (mô tả nhân vật thực tế trong episode). Khi người dùng dùng I2I với ảnh reference riêng, descriptor cũ trở nên sai → shot prompts gọi sai nhân vật.

## Problem Statement

### Problem 1: Variant Prompt thiếu cho nhiều character

**Hiện tại:** Extraction prompt hiểu `variant_prompt` là "delta thay đổi so với base". Logic này chỉ tạo variant_prompt cho state variant characters ("Anh Tôm (bị thương)"), không cho base characters hay characters 1 trạng thái.

**Kỳ vọng:** Mọi character đều có khả năng được generate từ ảnh tham chiếu (người dùng upload). `variant_prompt` phải tồn tại cho tất cả.

### Problem 2: Episode Descriptor không nhất quán sau I2I generation

**Chain hiện tại (bị lỗi):**
```
Extraction → episode_descriptor = "mô tả AI tưởng tượng"
                                          ↓
User upload ref image → I2I generate → Ảnh thực = "ngoại hình từ ref image của user"
                                          ↓
Shot prompts dùng episode_descriptor cũ → AI tạo nhân vật sai ngoại hình
```

**Root cause:** `episode_descriptor` được build từ `character_prompt` (T2I source of truth), nhưng sau khi user dùng I2I, source of truth chuyển sang ref image — descriptor không được update.

## Research Findings

### Codebase Patterns

#### Nơi episode_descriptor được dùng
- `storyboard_service.go` — `charDescMap` builder (4 nơi duplicate) — dùng trong shot image/video prompts
- `storyboard_nursery_service.go` — tương tự
- Priority hiện tại: `EpisodeDescriptor > Appearance > Description`

#### Nơi variant_prompt được tạo
- `application/prompts/fixed/character_extraction.txt` — fixed instructions
- `docs/cocomelon_template.md` — dynamic template (character_extraction section, ~line 490-539)
- `application/services/script_generation_service.go` — parse và save vào DB

#### Fields liên quan trong Character model
```go
CharacterPrompt   *string  // Full T2I standalone prompt
VariantPrompt     *string  // Delta I2I prompt (hiện chỉ có cho state variants)
EpisodeDescriptor *string  // Short identifier cho shot prompts (hiện derive từ T2I)
```

### Prior Solutions
- `plans/episode-character-variant.md` — Phase 1: thiết kế 3 fields, đã implement
- `todos/review-character-variants.md` — P3 note về charDescMap duplication

## Proposed Solution

### Nguyên lý cốt lõi

> **`variant_prompt` là nguồn truth về ngoại hình nhân vật trong episode.**
> `episode_descriptor` phải derive từ `variant_prompt`, không từ `character_prompt`.

Chain mới:
```
Extraction → variant_prompt = "Nhân vật trông như thế nào trong episode này, khi có ref image"
                     ↓
             episode_descriptor = compact(variant_prompt)  ← derive từ đây
                     ↓
Shot prompts gọi nhân vật → luôn đúng với ngoại hình episode-specific
```

### Phase 2A: Tái định nghĩa variant_prompt trong extraction prompt

**Thay đổi ở:** `character_extraction.txt` (fixed prompt) + `cocomelon_template.md` (dynamic)

> [!IMPORTANT]
> **Review note:** Đây là 2 thay đổi độc lập trong prompt, cần viết instructions riêng biệt:
> - **Change A:** Remove rule "empty string if no changes" → thay bằng "ALWAYS populate"
> - **Change B:** Add rule "episode_descriptor MUST derive from variant_prompt content"

**Định nghĩa mới của 3 fields:**

| Field | Định nghĩa mới | Thay đổi so với hiện tại |
|---|---|---|
| `character_prompt` | Full standalone T2I prompt. Đủ để generate nhân vật từ đầu không cần ảnh. ~200-400 words | Không thay đổi |
| `variant_prompt` | **ALWAYS required — NEVER empty string.** Prompt I2I hoàn chỉnh mô tả nhân vật AS THEY APPEAR IN THIS EPISODE. ~100-200 words | **Changed:** hiện tại cho phép `""` nếu không có costume change — sẽ bị xóa rule này |
| `episode_descriptor` | Compact visual identifier 20-30 words. **MUST be derived from variant_prompt content**, not from character_prompt | **Changed:** hiện tại AI tự generate, có thể derive từ character_prompt — sẽ yêu cầu explicit derive từ variant_prompt |

**Change A — Quy tắc xây dựng variant_prompt (ALWAYS required):**

- **Base character, 1 trạng thái, không có costume change:**
  ```
  Given a reference image of [Name], maintain the exact character design —
  same proportions, colors, and distinctive features. No modifications for this episode.
  Context: [brief episode context].
  ```

- **Base character có episode costume/state:**
  ```
  Given a reference image of [Name], maintain base design but apply:
  [episode-specific costume/accessory/state]. Keep all other features unchanged.
  ```

- **State variant** ("Anh Tôm (bị thương)"):
  ```
  Given a reference image of [Base Name], apply the following modifications:
  [specific physical changes]. All other base features unchanged.
  ```

**Change B — Quy tắc derive episode_descriptor từ variant_prompt:**
- Trích xuất các visual identifiers quan trọng nhất TỪ variant_prompt: trang phục, màu sắc đặc trưng, điểm nhận dạng
- 20-30 words, không dùng tên nhân vật
- Đủ cho AI tạo ảnh nhận ra nhân vật khi kết hợp với reference image

### Phase 2B: Cập nhật charDescMap priority logic

**Hiện tại (trong storyboard services):**
```go
charDescMap[name] = EpisodeDescriptor > Appearance > Description
```

**Sau Phase 2A,** không cần thay đổi logic — vì `EpisodeDescriptor` bây giờ đã được build đúng từ `variant_prompt`, tự nhiên nó chính xác hơn.

**Bonus: giải quyết P3 todo** — refactor charDescMap builder thành shared utility:
```go
// pkg/utils/character_utils.go
func BuildCharDescMap(characters []models.Character) map[string]string
```

### Phase 2C: UI hint sau I2I generation với custom reference

Khi user upload ảnh reference riêng và generate I2I:
1. `variant_prompt` (AI extraction) mô tả nhân vật theo text description trong MV input
2. Ảnh reference của user có thể show một ngoại hình khác với text đó
3. Sau khi generate xong, `episode_descriptor` trong DB vẫn là text cũ → shot prompts sẽ mô tả nhân vật theo text, không theo ảnh thực

**Implementation approach (Chosen: auto-derive từ variant_prompt):**

Sau I2I generation thành công, hiện alert trong Edit Prompt dialog:
```
⚠️ Reference image uploaded. Your Episode Descriptor may not match 
your reference image. [Sync from Variant Prompt] [Keep Current]
```

Nút **"Sync from Variant Prompt"**:
- Gọi `PATCH /api/v1/characters/:id` với `episode_descriptor = <derived from current variant_prompt>`
- Derivation logic: extract first 20-30 words describing visual appearance từ variant_prompt
- Thực hiện client-side (string truncation/extraction), không cần AI call thêm
- Sau khi sync: toast "Episode Descriptor updated. Shot prompts will now reference this description."

**Tại sao không dùng AI call để analyze ref image:**
- Expensive (extra API call per character)
- Không cần thiết — variant_prompt đã describe đủ appearance
- Auto-derive từ variant_prompt là deterministic và instant

## Acceptance Criteria

- [ ] Mọi character (base, 1 trạng thái, multiple states) đều có `variant_prompt` sau extraction
- [ ] `variant_prompt` cho base characters = "Given reference image, maintain exact design + episode context"
- [ ] `variant_prompt` cho state variants = "Given base reference, apply these specific modifications"
- [ ] `episode_descriptor` được AI generate từ nội dung `variant_prompt` (không phải từ `character_prompt`)
- [ ] Shot prompts dùng `episode_descriptor` mới → nhân vật nhất quán với episode appearance
- [ ] Sau I2I generation với custom ref, UI hiện gợi ý update `episode_descriptor`
- [ ] `charDescMap` builder được extract thành shared utility (bonus, P3)

## Technical Considerations

### Scope thay đổi

| Component | Thay đổi |
|---|---|
| `character_extraction.txt` (fixed) | Cập nhật instructions cho variant_prompt (ALWAYS required, new definition) + episode_descriptor phải derive từ variant_prompt |
| `cocomelon_template.md` (dynamic) | Cập nhật character_extraction section để align với định nghĩa mới |
| `EpisodeWorkflow.vue` | Thêm UI hint sau I2I generation |
| `storyboard_service.go` / `storyboard_nursery_service.go` | Optional: refactor charDescMap thành shared util |

### Không thay đổi

- Backend Go models (Character struct đã có đủ 3 fields)
- Database schema (không cần migration)
- Generation endpoints
- Priority logic (`EpisodeDescriptor > Appearance > Description`) vẫn đúng

### Risks

- **Extraction token count tăng:** variant_prompt cho mọi character sẽ tăng response length. Hiện đã tăng max_tokens lên 6000. Cần monitor.
- **Prompt regression:** Thay đổi extraction prompt có thể ảnh hưởng đến chất lượng extraction cho các content khác. Cần test với nhiều input types.
- **Backward compatibility:** Characters đã tồn tại trong DB có `variant_prompt = NULL`. Sau fix, re-extract sẽ populate. Không ảnh hưởng runtime vì logic vẫn fallback về Appearance.

### Alternatives Considered

**Alt A: Hai episode_descriptor (T2I và I2I riêng biệt)**
- Phức tạp hơn: shot prompts phải biết mỗi character dùng mode nào
- Rejected: thêm complexity không cần thiết

**Alt B: Auto-derive episode_descriptor từ ảnh sau khi generate (image captioning)**
- Chính xác nhất nhưng tốn kém (extra AI call per character)
- Rejected: not worth cost for this use case

**Alt C (chosen): variant_prompt là source of truth, episode_descriptor derive từ nó**
- Single source of truth
- Episode_descriptor luôn reflect ngoại hình episode-specific
- Manual update hint sau I2I với custom ref xử lý edge case còn lại

## Implementation Steps

### Step 1: Cập nhật extraction prompt (fixed)
- File: `application/prompts/fixed/character_extraction.txt`
- **Change A:** Xóa rule "Set to empty string if no changes" ở `variant_prompt` → thay bằng "ALWAYS required, never empty"
- **Change B:** Thêm rule explicit: `episode_descriptor` MUST be derived from `variant_prompt` content, not from `character_prompt`
- Cập nhật example trong file (Driver Wang: thêm variant_prompt, update episode_descriptor)

### Step 2: Cập nhật dynamic template (Cocomelon)
- File: `docs/cocomelon_template.md` — character_extraction section
- Align với fixed prompt: xóa rule empty string, thêm ví dụ variant_prompt cho base characters

### Step 3: UI — Sync Episode Descriptor sau I2I
- File: `web/src/views/drama/EpisodeWorkflow.vue`
- Sau khi `saveAndGenerate` với I2I mode thành công → hiện alert với 2 nút: "Sync from Variant Prompt" / "Keep Current"
- "Sync from Variant Prompt": derive episode_descriptor từ variant_prompt (client-side extract), gọi `PATCH /api/v1/characters/:id`

### Step 4 (Optional P3): Refactor charDescMap
- File: `pkg/utils/character_utils.go` (new)
- Extract shared builder, update callers in storyboard services

## Verification

Sau khi implement, thực hiện theo thứ tự:

1. **Re-extract 1 episode test** với input có mixed characters (base + variant)
2. **Kiểm tra DB:** Mọi character đều có `variant_prompt != NULL` và `episode_descriptor != NULL`
3. **Kiểm tra nội dung:** `episode_descriptor` của base characters (1 trạng thái) phải phản ánh episode context, không phải base appearance
4. **Generate storyboard:** Verify shot prompts dùng descriptor mới — mô tả đúng nhân vật trong episode
5. **Test I2I flow:** Upload ref image → generate → xác nhận alert hiện → bấm Sync → verify episode_descriptor trong DB được update

## References

- `plans/episode-character-variant.md` — Phase 1 original plan
- `todos/review-character-variants.md` — P3 charDescMap refactor note
- `todos/pending-p3-consolidate-character-extraction-services.md` — service consolidation
- `application/prompts/fixed/character_extraction.txt` — Change A & B target
- `docs/cocomelon_template.md` — dynamic template target
