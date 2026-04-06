# 🎙️ Prompt Phân Tích: Narrator - Dialogue Interleaving trong Video Voice-over

> **Mục đích:** Phân tích CHÍNH XÁC cách video voice-over xử lý việc đan xen giữa narrator (giọng kể) và dialogue (đối thoại/lời nhân vật). Output dùng để tinh chỉnh hệ thống AI Director cho voice-over video generation.

> **Input:** Upload 1-2 video clips voice-over/documentary (2-5 phút mỗi clip, ưu tiên đoạn CÓ CẢ narrator và dialogue/quote/reaction). Hoặc upload 30-50 screenshots liên tục kèm transcript.

---

## Prompt sử dụng:

```
Bạn là chuyên gia phân tích audio-visual design cho video voice-over / documentary. Tôi cung cấp video/screenshots từ "[TÊN KÊNH / VIDEO]."

Nhiệm vụ: Phân tích CỰC KỲ CHI TIẾT cách video này đan xen giữa NARRATOR (giọng kể chuyện) và DIALOGUE (lời nhân vật, quote, phản ứng). Trả lời bằng TIẾNG VIỆT.

---

## 1. PHÂN LOẠI AUDIO MODE TỪNG SHOT

Xem xét TOÀN BỘ video, phân loại MỖI shot/cảnh vào một trong các loại sau:

| # | Timestamp | Duration | Audio Mode | Narrator? | Dialogue? | Dialogue Content | Nhận xét |
|---|-----------|----------|------------|-----------|-----------|------------------|----------|
| 1 | 00:00-00:04 | 4s | ? | ? | ? | ? | ? |
| 2 | ... | ... | ... | ... | ... | ... | ... |

**Audio Mode categories (chọn 1 trong 2):**
- **narrator_only**: Chỉ có giọng narrator, KHÔNG có lời nhân vật
- **dialogue_dominant**: Dialogue là chủ đạo hoặc có đan xen, narrator có thể im lặng, ducking, hoặc xuất hiện ngắn

Sau bảng, tổng hợp:
- Tổng số shot narrator_only: ? / ? (? %)
- Tổng số shot dialogue_dominant: ? / ? (? %)

---

## 2. PHÂN TÍCH DIALOGUE TYPES

Cho MỖI shot có dialogue (dialogue_dominant), phân loại dialogue vào một TYPE cụ thể:

| Type | Mô tả | Ví dụ từ video | Số lần | % |
|------|--------|----------------|--------|---|
| **reaction** | Âm thanh phản ứng ngắn, 1-3 từ (reo hò, tiếng thở dài, "wow", "incredible!") | ? | ? | ? |
| **soft_line** | Câu nói ngắn 4-8 từ, narrator vẫn chạy ở background | ? | ? | ? |
| **quote** | Trích dẫn lời nhân vật lịch sử/thật, narrator tạm dừng | ? | ? | ? |
| **narrated_quote** | Narrator tự đọc lời quote (không có voice actor riêng) | ? | ? | ? |
| **full_dialogue** | Đoạn đối thoại dài 2+ câu giữa các nhân vật | ? | ? | ? |
| **inner_voice** | Suy nghĩ nội tâm của nhân vật | ? | ? | ? |
| **crowd** | Tiếng đám đông, la hét, hô vang | ? | ? | ? |
| **ambient_voice** | Giọng nói xung quanh không rõ lời, tạo atmosphere | ? | ? | ? |

**Câu hỏi bổ sung:**
- Dialogue có dùng voice actor riêng hay narrator thay đổi giọng?
- Dialogue có được xử lý âm thanh khác narrator không? (reverb, EQ, filter, echo?)
- Có subtitle/text hiển thị kèm dialogue không?

---

## 3. PHÂN TÍCH PATTERN ĐÁNG CHÚ Ý: KHI NÀO DIALOGUE XUẤT HIỆN?

### 3a. Trigger conditions — Dialogue xuất hiện khi nào?

Xác định CÁC TÌNH HUỐNG nội dung dẫn đến dialogue (check tất cả áp dụng):

- [ ] Khi nhân vật lịch sử được giới thiệu lần đầu (quote nổi tiếng)
- [ ] Khi có emotional peak / climax trong story
- [ ] Khi có dramatic reveal / twist / bất ngờ
- [ ] Khi thể hiện conflict / đối đầu giữa hai bên
- [ ] Khi mô tả hành động cụ thể của nhân vật
- [ ] Khi cần "phá" nhịp narrator dài (visual/audio variety)
- [ ] Khi chuyển từ mô tả chung sang specific moment
- [ ] Khi quote là bằng chứng / evidence cho lập luận
- [ ] Random — không có pattern rõ ràng
- [ ] Khác: ?

### 3b. Vị trí dialogue trong dòng kể chuyện:
- Dialogue xuất hiện ở đầu video? Giữa? Cuối? Đều khắp?
- Có cluster (nhiều shots dialogue liên tiếp) không? Hay rải đều?
- Khoảng cách trung bình giữa 2 dialogue shots (bao nhiêu narrator-only shots ở giữa)?

### 3c. Dialogue density theo phút:
- Phút 0-1: ? shot dialogue / ? tổng shots
- Phút 1-2: ? / ?
- Phút 2-3: ? / ?
- ...
- Có trend tăng/giảm dialogue theo thời gian không?

---

## 4. PHÂN TÍCH KỸ THUẬT CHUYỂN ĐỔI NARRATOR ↔ DIALOGUE

### 4a. Transition VÀO dialogue:
- Narrator dừng đột ngột hay fade out?
- Có silence gap giữa narrator và dialogue không? Bao lâu?
- Có SFX transition (whoosh, reverb tail) không?
- Có visual cue (zoom in, cắt cảnh mới, đổi màu) khi dialogue bắt đầu?
- Narrator có "cue" cho dialogue không? (VD: "He said...", "She proclaimed...")
- Music: giữ nguyên, giảm, dừng, hay đổi khi dialogue bắt đầu?

### 4b. Transition RA KHỎI dialogue:
- Dialogue kết thúc bằng gì? (câu nói trọn vẹn? bị cắt? fade?)
- Narrator quay lại ngay lập tức hay có gap?
- Có echo/reverb tail trên dialogue sau khi narrator đã quay lại?
- Visual: giữ cảnh hay cắt khi quay lại narrator?

### 4c. Xử lý audio trong mode "dialogue_dominant" (nếu narrator vẫn đang nói):
- Narrator volume giảm (ducking) bao nhiêu % khi dialogue xuất hiện?
- Có crossfade giữa narrator và dialogue không?
- Dialogue volume so với narrator: louder / same / softer?
- Dialogue pan (stereo positioning): center? slight left/right? full stereo?

---

## 5. PHÂN TÍCH HÌNH ẢNH KHI CÓ DIALOGUE

### 5a. Visual treatment cho dialogue shots vs narrator shots:
- Shot size có thay đổi không? (VD: narrator → wide ; dialogue → close-up?)
- Camera movement có khác không?
- Color grading có khác không? (VD: flashback tint, sepia tone?)
- Có text overlay (quote text on screen) không?
- Nhân vật nói có xuất hiện trực tiếp không hay chỉ có voice-over?

### 5b. Với từng dialogue type, hình ảnh tương ứng là gì?
| Dialogue Type | Hình ảnh đi kèm | Ví dụ cụ thể |
|---------------|-----------------|---------------|
| reaction | ? | ? |
| soft_line | ? | ? |
| quote | ? | ? |
| narrated_quote | ? | ? |
| full_dialogue | ? | ? |

---

## 6. SO SÁNH MẬT ĐỘ DIALOGUE GIỮA CÁC VIDEO/KÊNH

Nếu có nhiều video, so sánh:
| Metric | Video 1 | Video 2 | Video 3 |
|--------|---------|---------|---------|
| Tổng shots | ? | ? | ? |
| % narrator_only | ? | ? | ? |
| % dialogue_dominant | ? | ? | ? |
| Dialogue types chủ yếu | ? | ? | ? |
| Khoảng cách TB giữa dialogue shots | ? | ? | ? |

---

## 7. KEY INSIGHTS & RULES EXTRACTION

Tổng hợp thành RULES có thể dùng trực tiếp cho AI Director system:

### 7a. Audio Mode Distribution Rule:
```
narrator_only: ~?% (range: ?%-?%)
dialogue_dominant: ~?% (range: ?%-?%)
```

### 7b. Dialogue Trigger Rules:
Liệt kê 5-8 rules dạng:
```
IF [condition] THEN audio_mode = [mode], dialogue_type = [type]
```

VD:
```
IF nhân vật lịch sử xuất hiện lần đầu THEN audio_mode = "dialogue_dominant", dialogue_type = "quote"
IF action scene + nhân vật la hét THEN audio_mode = "dialogue_dominant", dialogue_type = "reaction"
```

### 7c. Transition Rules:
```
Narrator → Dialogue: [mô tả kỹ thuật chuyển đổi chuẩn]
Dialogue → Narrator: [mô tả kỹ thuật chuyển đổi chuẩn]
Max dialogue duration before narrator returns: ? seconds
Min narrator gap between two dialogue moments: ? shots / ? seconds
```

### 7d. Dialogue Content Rules:
```
reaction: max ? words, tone: ?
soft_line: max ? words, narrator: ducked/off
quote: tone: ?, visual: ?, narrator: off
full_dialogue: max ? seconds, narrator: off, visual treatment: ?
```

### 7e. Visual-Audio Sync Rules:
```
Khi audio_mode thay đổi → visual [thay đổi / giữ nguyên]
Dialogue shot size preference: [CU / MS / varies]
Dialogue camera movement: [static / slow zoom / varies]
```

---

## 8. SAMPLE SHOT LIST WITH DIALOGUE (Output Mẫu)

Viết 10-15 shots mẫu (narrator xen dialogue) GIỐNG PHONG CÁCH VIDEO ĐÃ PHÂN TÍCH, về một chủ đề KHÁC. Format:

```json
[
  {
    "shot_id": 1,
    "script_segment": "...",
    "audio_mode": "narrator_only",
    "dialogue_type": "none",
    "dialogue_text": "",
    "visual_description": "...",
    "duration": 4,
    "note": "establishing shot, pure narration"
  },
  {
    "shot_id": 5,
    "script_segment": "...",
    "audio_mode": "dialogue_dominant",
    "dialogue_type": "reaction",
    "dialogue_text": "Incredible!",
    "visual_description": "...",
    "duration": 3,
    "note": "crowd reaction mixed with narrator"
  }
]
```

Trả lời bằng TIẾNG VIỆT, sample shots viết bằng TIẾNG ANH.
```
