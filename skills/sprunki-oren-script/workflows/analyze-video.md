# Workflow: Phân tích Video gốc thành Script (Adult Version)

## Đầu vào

1. **Video nguồn** (Sprunki hoặc tương tự)
2. **Phạm vi phân tích** (ví dụ: 00:00 — 02:30)

## Prerequisites

- Đọc `SKILL.md` (world rules, character roster — ADULT VERSION)
- Đọc `templates/shot-format.md` (output format)
- Đọc `references/voice-profiles.md` (để nhận diện voice cues)

## Quy trình

### Bước 1: Xem tổng thể

Ghi chú: tổng cuts, nhịp edit, nhân vật, locations, emotional arc.
**Đặc biệt chú ý (Adult Context):**
- Những props mang tính người lớn (coffee, beer, laptop, xe hơi)
- Bối cảnh người lớn (bar, apartment, văn phòng)
- Biểu cảm (chill, tự tin, mệt mỏi — KHÔNG childish)

### Bước 2: Chia SHOT theo CUT

SHOT mới khi: CUT, đổi camera, đổi action, đổi không gian, nhảy thời gian.

### Bước 3: Trích xuất từng SHOT

| Element | Trích xuất |
|---|---|
| Duration | Cut-to-cut |
| Shot type | WS/MS/CU/MWS/INSERT |
| Camera | Static? Tracking? Handheld? |
| Subject | Ai, ở đâu trong frame |
| Action | 1 hành động chính (chú ý body language: swagger, lean) |
| Environment | American context cụ thể (có adult cues) |
| Audio | SFX? Dialogue? BGM? |

### Bước 4: Viết Visual [NOTE]

```
[NOTE] [CAMERA] + [SUBJECT + VỊ TRÍ] + [ACTION] + [CONTEXT]
```
*Lưu ý:* Cố gắng sử dụng các từ ngữ trưởng thành (VD: "leans on the counter", "holds a coffee mug", "checks phone") để đảm bảo AI generate ra character với adult proportions.

### Bước 5: Audio Extraction

- `[SFX]` — BẮT BUỘC có cho mỗi shot (ambient của bar, apartment, traffic, coffee machine)
- `[DIA]` — Match tên nhân vật. Lời thoại nên được diễn đạt theo đúng Voice Profile (ví dụ: Oren thì ngắn gọn, chill, "Aight").
- `[BGM]` — Chỉ ghi khi có sự thay đổi.
- `[CAM]` — Hướng camera.

### Bước 6: Compile + Verify

- [ ] Mỗi Visual có thể dùng làm AI prompt.
- [ ] Không có style/lighting descriptions.
- [ ] Chỉ 1 action per shot.
- [ ] Camera đứng đầu Visual.
- [ ] Mọi shot đều có SFX.
- [ ] **ADULT CHECK**: Các prompt có đủ cues (prop, location, body language) để AI tạo ra nhân vật người trưởng thành không?
- [ ] **VOICE CHECK**: Dialogue match với Voice Profile không?
