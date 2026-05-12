# Workflow: Phân tích Video gốc thành Script

## Đầu vào

User cung cấp:
1. **Video nguồn** (HL Meow Meow hoặc video mèo tương tự)
2. **Phạm vi phân tích** (ví dụ: 00:00 — 02:45)

## Prerequisites

- Đọc `SKILL.md` (world rules, character roster)
- Đọc `templates/shot-format.md` (output format)
- Có khả năng xem/phân tích video

## Quy trình

### Bước 1: Xem toàn bộ video trước

Xem hết phạm vi được chỉ định 1 lần. Ghi chú:
- Tổng số cuts (chuyển cảnh)
- Nhịp edit chung (nhanh/chậm/montage)
- Nhân vật xuất hiện
- Locations
- Emotional arc tổng thể

### Bước 2: Chia SHOT theo CUT

Xem lại frame-by-frame. Tạo SHOT mới khi:
- Có CUT / chuyển cảnh
- Đổi góc camera
- Đổi hành động chính
- Đổi không gian
- Nhảy thời gian

**Ưu tiên chia shot giống nhịp edit của video gốc.**

### Bước 3: Trích xuất từng SHOT

Cho mỗi shot, ghi nhận:

| Element | Cách trích xuất |
|---|---|
| Duration | Đếm giây từ cut-to-cut |
| Shot type | Xác định: WS/MS/CU/MWS/INSERT |
| Camera angle | Cat-eye level? High? Human-eye? |
| Camera movement | Static? Tracking? Push-in? Handheld? |
| Subject | Ai xuất hiện, ở đâu trong frame |
| Action | 1 hành động chính (động từ vật lý) |
| Environment | Bối cảnh cụ thể |
| Audio | SFX nghe được? Dialogue? BGM? |
| Human Reaction | Pattern nào? (nếu có người) |

### Bước 4: Viết Visual cho mỗi SHOT

Áp dụng công thức:
```
[CAMERA SETUP] + [SUBJECT + VỊ TRÍ] + [ACTION CHÍNH] + [CONTEXT]
```

**Luật viết:**
1. Camera LUÔN đứng đầu
2. Subject + vị trí trong frame rõ ràng
3. 1 action chính duy nhất
4. Context/environment cụ thể
5. Blocking: vị trí, hướng di chuyển, khoảng cách

**CẤM:**
- ❌ Mô tả style (lighting, color, aesthetic)
- ❌ Cảm xúc trừu tượng
- ❌ Nhiều action trong 1 shot
- ❌ Thêm nội dung không có trong video
- ❌ Chi tiết ngoại hình nhân vật
- ❌ Văn phong ẩn dụ

### Bước 5: Audio Extraction

Cho mỗi shot, trích xuất riêng:
- `[SFX]` — âm thanh cụ thể (onomatopoeia nếu có)
- `[DIA]` — thoại rõ ràng + chỉ định giọng
- `[BGM]` — chỉ khi mood thay đổi
- `[CAM]` — mô tả chuyển động camera

### Bước 6: Compile Output

Ghép tất cả theo format:

```
PHẦN X: [Tên phần]
(timestamp range | số shot)

SHOT XX | timestamp | shot_type
[BGM: ...]
[SFX: ...]
[DIA: ...]
[CAM: ...]
[HUMAN REACTION: pattern]

Visual: ...
```

### Bước 7: Verification

Kiểm tra cuối:
- [ ] Mỗi Visual có thể dùng trực tiếp làm prompt
- [ ] Không có mô tả style/lighting/aesthetic
- [ ] Mỗi shot chỉ có 1 action chính
- [ ] Camera LUÔN đứng đầu Visual
- [ ] Blocking cụ thể cho mọi shot
- [ ] SFX cho mọi shot
- [ ] Anthropomorphic behavior đúng (hành vi NGƯỜI)

## Expected Output

Full script shot-by-shot, chia theo PHẦN, mỗi shot theo format chuẩn.
Không giải thích, không commentary — chỉ output script.
