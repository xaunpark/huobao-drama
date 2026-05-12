# Workflow: Viết Script từ Ý tưởng / Outline

## Đầu vào

User cung cấp một trong các dạng:
- **Ý tưởng 1 câu**: "Tabinyan đi làm shipper giao hàng"
- **Outline sơ bộ**: Danh sách các cảnh chính
- **Tên episode**: "Tabinyan's First Day at the Sushi Shop"

## Prerequisites

- Đọc `SKILL.md` (đã nạp tự động — world rules, character roster, storytelling DNA)
- Đọc `templates/shot-format.md` để nắm output format

## Quy trình

### Bước 1: Xác định Episode Structure

Từ ý tưởng, xây dựng **Story Beats** theo công thức chuẩn:

**Single Episode:**
```
GOAL → 0 YEN! → WORK → [COMPLICATION] → EARN → FULFILL GOAL
```

**Compilation (video dài):**
```
[Mini-arc 1: GOAL→WORK→FULFILL] → [Kết nối nhân quả] → [Mini-arc 2] → ...
```

**Xác định cho mỗi arc:**
- **GOAL**: Tabinyan muốn gì? (mua váy, đi chơi, tặng quà — phải CỤ THỂ)
- **0 YEN trigger**: Cảnh kiểm ví → trống rỗng (signature moment bắt buộc)
- **WORK**: Làm việc gì? Ở đâu? (liên quan đến goal nếu có thể)
- **EARN**: Salary ceremony (ojigi, phong bì, "Otsukaresama desu")
- **FULFILL**: Thỏa mãn mục tiêu ban đầu — đây là climax cảm xúc
- Locations: 2-4 địa điểm Nhật Bản
- Humans: vai trò + reaction pattern
- Moral lesson: tự nhiên, không ép (xem `references/episode-ideas.md`)
- Duration: 1-2 phút (10-15 shots) hoặc 3-10 phút compilation (25-60 shots)

### Bước 2: Xây dựng Shot List

Lập danh sách shots theo flow câu chuyện. Áp dụng các quy tắc:

**Shot Distribution:**
- WS establishing → MS action → CU reaction → MS next action (lặp lại)
- 40% MS, 20% CU, 15% WS, 15% MWS (two-shot), 10% Insert

**Pacing:**
- Action montage: 2-3s/shot
- Emotion/reaction: 3-4s/shot  
- Establishing: 3-5s/shot

**3-Beat Arc cho mỗi hành động:**
1. Setup (anticipation shot)
2. Action (peak effort shot)  
3. Payoff (reaction/satisfaction shot)

**Human Interaction Pattern:**
```
MS Tabinyan action → MWS two-shot → CU human reaction → CU Tabinyan reaction
```

### Bước 3: Viết từng SHOT

Cho mỗi shot, viết theo format trong `templates/shot-format.md`:

```
SHOT XX | duration | shot_type
[SFX: cụ thể, chi tiết]
[DIA: nếu có — Tabinyan chỉ thán từ ngắn]
[CAM: handheld behavior cụ thể]
[HUMAN REACTION: pattern_name — nếu có người]

Visual: [CAMERA] + [SUBJECT + VỊ TRÍ] + [ACTION] + [CONTEXT]
```

**Checklist cho mỗi Visual:**
- [ ] Camera đứng đầu câu (shot type, angle)
- [ ] Subject rõ ràng (ai, ở đâu trong frame)
- [ ] Chỉ 1 hành động chính (động từ vật lý)
- [ ] Context/environment mô tả rõ
- [ ] Blocking cụ thể (vị trí, hướng di chuyển, khoảng cách)
- [ ] KHÔNG mô tả style/lighting/aesthetic
- [ ] KHÔNG mô tả ngoại hình (fur color, eye shape)
- [ ] Anthropomorphic: hành vi NGƯỜI, không hành vi mèo

### Bước 4: Audio Layer

Xuyên suốt script:
- **[SFX]** cho MỌI shot — Foley là storytelling tool chính
- **[DIA]** chỉ khi cần: Tabinyan thán từ ngắn, Người dùng Keigo
- **[BGM]** chỉ ở shot đầu hoặc khi mood thay đổi lớn
- **[NARRATOR]** tối thiểu: 5-10 từ, chỉ khi visual không đủ
- **[PAUSE: Xs]** cho breathing space giữa các sequence
- **70%+ shot là SFX-only** — không narration

### Bước 5: Review & Polish

**Checklist cuối:**
- [ ] Story beats đầy đủ (NEED → CHALLENGE → REWARD)
- [ ] 3-Beat Arc cho mỗi chuỗi hành động
- [ ] Shot distribution đúng tỷ lệ (~40% MS, ~20% CU...)
- [ ] Camera angle: 70%+ là cat-eye level
- [ ] Human interactions có reaction pattern
- [ ] Tabinyan luôn ĐÚNG anthropomorphic behavior
- [ ] 70%+ shots là SFX-only
- [ ] ZERO text on screen
- [ ] Mỗi Visual có thể dùng trực tiếp làm AI prompt
- [ ] Blocking/staging cụ thể cho mọi shot

## Expected Output

Script hoàn chỉnh 10-40 shots, mỗi shot theo format chuẩn.
Có thể paste trực tiếp vào hệ thống hoặc dùng từng Visual làm image/video prompt.
