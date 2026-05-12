---
name: hl-meow-meow-script
description: Generate detailed cinematic scripts for HL Meow Meow anthropomorphic cat life-vlog episodes
---

# HL Meow Meow — Cinematic Script Generator

Tạo **kịch bản phân cảnh chi tiết** (shot-by-shot cinematic script) cho video mèo nhân hóa phong cách HL Meow Meow.
Mỗi SHOT trong output có thể dùng trực tiếp làm prompt cho AI image/video generation (Midjourney, Flux, Kling, Veo).

## Khi nào dùng

- User yêu cầu viết kịch bản cho một episode mèo mới (ý tưởng hoặc outline)
- User cung cấp video gốc cần phân tích thành script chi tiết
- User cần chuyển đổi một concept/outline thành storyboard shots

## Instrumentation

```bash
./scripts/log-skill.sh "hl-meow-meow-script" "manual" "$$"
```

## What do you want to do?

1. **Viết script từ ý tưởng/outline** → Read `workflows/write-from-idea.md`
2. **Phân tích video gốc thành script** → Read `workflows/analyze-video.md`
3. **Xem character & world reference** → Continue reading this file (Section: World Bible)
4. **Xem output format** → Read `templates/shot-format.md`
5. **Tạo Title + Description + Thumbnail Prompt** → Read `references/title-description.md`

---

## Nguyên tắc cốt lõi (LUÔN NẠP)

### 1. Thế giới (World-building)

- **Bối cảnh**: Xã hội Nhật Bản lý tưởng hóa (Utopia). Mèo và người chung sống bình đẳng.
- **Quy tắc chấp nhận mặc nhiên**: KHÔNG AI đặt câu hỏi "Tại sao mèo lại ở đây?". Mèo có hộ chiếu, ký hợp đồng lao động, trả thuế. Đây là bình thường.
- **Keigo (Kính ngữ)**: Người LUÔN dùng kính ngữ với mèo — tạo humor tinh tế khi người 1m80 cúi chào mèo 30cm.
- **Bending Down**: Người luôn cúi thấp ngang tầm mắt mèo khi giao tiếp trực tiếp.
- **Ngôn ngữ**: Tất cả signage, thoại nhân vật phụ (người) bằng tiếng Nhật. Tabinyan dùng thán từ Nhật ngắn.

### 2. Nhân vật (Character Roster)

| Nhân vật | Vai trò | Xuất hiện | Đặc điểm |
|---|---|---|---|
| **Tabinyan** | Mèo con gừng nhân hóa — nhân vật chính | 100% | Đứng 2 chân, dùng paw như tay. Chăm chỉ, tò mò, can đảm. Voice: cực cao, chirpy |
| **Obaa-chan** | Bà ngoại | ~20% | Người Nhật già, kimono, hiền lành. Động lực: Tabinyan chăm sóc bà |
| **Nobita-kun** | Bạn trai | ~15% | Đi chơi, phiêu lưu cùng Tabinyan |
| **Nyanko-chan** | Bạn gái mèo | ~15% | Động lực tình cảm: Tabinyan mua quà cho cô ấy |
| **Người Nhật** | Warm Support System | ~40% | Nhân viên, bác sĩ, hàng xóm. Luôn tử tế, dùng Keigo |

### 3. Tabinyan Voice Profile

Tabinyan KHÔNG "nói" như người. Chỉ dùng **thán từ tiếng Nhật cực ngắn** + **âm thanh mèo**:

| Âm thanh | Khi nào | Pitch |
|---|---|---|
| "Oishii!" | Ăn miếng đầu tiên | Cực cao, chirpy |
| "Yatta!" | Hoàn thành việc | Cực cao, phấn khích |
| "Sugoi!" | Ngạc nhiên | Cao, mở mắt to |
| "Itai!" | Đau nhẹ | Cao, ngắn |
| "Ara?" | Phát hiện điều mới | Tò mò, ngân nhẹ |
| "Meow!" | Chào hỏi | Bright, dứt khoát |
| "Purr~" | Thoải mái | Thấp, rung nhẹ |
| "Chi chi!" | Vui vẻ | Chirp giống chim |

### 4. 6 Human Reaction Patterns

Mỗi cảnh có người tương tác, PHẢI chỉ định reaction pattern:

| Pattern | Trigger | Hành động |
|---|---|---|
| **gentle_mentor** | Dạy việc | Cúi ngang mắt mèo, gật đầu, làm mẫu chậm |
| **admiring_customer** | Nhận đồ | Hai tay nhận, nghiêng đầu tan chảy, giơ phone chụp |
| **formal_respect** | Trao lương/chứng nhận | Đứng nghiêm, ojigi, hai tay trao |
| **compassionate_caregiver** | Mèo mệt/đau | Xoa đầu, nắm paw nhẹ, giọng thấp |
| **silent_admirer** | Nơi công cộng | Người qua đường dừng lại nhìn, mỉm cười |
| **amused_observer** | Mèo hoạt động vui | Nghiêng người, che miệng cười, "Kawaii~" |

### 5. Storytelling DNA — 3-Beat Arc

Mọi hành động đều tuân theo mô hình 3 nhịp:

1. **Anticipation (Chuẩn bị)**: Tabinyan đứng nghiêm túc, chống nạnh hoặc sờ cằm nhìn vào nhiệm vụ. Ánh sáng trong, sáng. Energy: tò mò.
2. **Peak Action (Nỗ lực tối đa)**: Tabinyan dùng paw lóng ngóng nhưng cực tập trung. Close-up để bắt biểu cảm. Energy: charm peak.
3. **Resolution (Thỏa mãn)**: Lau mồ hôi, cười hài lòng, "Yatta!" hoặc "Oishii!". Energy: ấm áp, contentment.

### 6. Story Beat Template — Công thức chuẩn "GOAL → WORK → FULFILL"

> **Quy tắc vàng**: Tabinyan KHÔNG BAO GIỜ đi làm vì thích làm. Luôn có **MỤC TIÊU CỤ THỂ** ở đầu (mua váy, đi công viên nước, tặng quà bà) → phát hiện **hết tiền** → đi làm kiếm tiền → **thỏa mãn mục tiêu ban đầu**.

#### A. Công thức SINGLE EPISODE (1-2 phút, 10-15 shots)

```
GOAL → 0 YEN! → WORK → [COMPLICATION] → EARN → FULFILL GOAL
```

1. **GOAL (Mục tiêu)**: Tabinyan muốn thực hiện điều gì đó cụ thể:
   - Mua váy tặng Nyanko-chan
   - Đi công viên nước với Nobita-kun
   - Mua quà sinh nhật cho Obaa-chan
   - Đi ngắm hoa anh đào
   - Mua kimono mới cho lễ hội
2. **0 YEN! (Trigger)**: Kiểm ví → trống rỗng. "Ara?" → lật ví ngược → không rơi gì. Đây là **signature moment** bắt buộc.
3. **WORK (Làm việc)**: Tabinyan tìm việc part-time liên quan. Human: gentle_mentor dạy việc. Montage SFX-heavy (3-5 shots).
4. **COMPLICATION** *(optional)*: Khó khăn nhỏ (làm đổ, bị bỏng nhẹ, lạc đường). Human: compassionate_caregiver.
5. **EARN (Nhận lương)**: Ceremony trang trọng. Human: formal_respect, ojigi, phong bì hai tay. "Otsukaresama desu". Tabinyan: "Yatta!"
6. **FULFILL GOAL (Thỏa mãn)**: Tabinyan dùng tiền thực hiện mục tiêu ban đầu. Đây là **climax cảm xúc** — Nyanko-chan mặc váy mới, Obaa-chan mở quà, cả nhóm ở công viên nước.
   - Kết thúc: cảnh ấm áp, "Oyasuminasai", hoặc hai nhân vật ngồi bên nhau.

#### B. Công thức COMPILATION (3-10 phút, 25-60 shots)

Video dài = chuỗi **3-5 mini-arcs** nối nhau bằng kết nối nhân quả:

```
[Mini-arc 1: GOAL→WORK→FULFILL] → [Transition] → [Mini-arc 2: GOAL→WORK→FULFILL] → ...
```

**3 loại kết nối giữa các mini-arcs:**

| Loại | Cách hoạt động | Ví dụ |
|---|---|---|
| **Tài chính** | Mục tiêu mới cần thêm tiền | "Muốn đi onsen nhưng hết tiền → làm thêm ở konbini" |
| **Hệ quả** | Hành động trước gây ra tình huống mới | "Ăn quá nhiều takoyaki → đau bụng → đi bệnh viện" |
| **Tình cờ** | Sự kiện ngẫu nhiên mở ra trải nghiệm | "Nhặt ví tiền → mang đến đồn cảnh sát → được khen thưởng" |

**Nén thời gian**: Video cảm giác "1 ngày bận rộn" nhưng thực tế là nhiều ngày. Dấu hiệu:
- Manager nói "Hẹn gặp 8 giờ sáng mai"
- Nhận "tiền lương tháng này" (今月分のお給料)
- Cảnh phục hồi bệnh (nằm viện) trước khi bắt đầu arc mới

#### C. Bài học đạo đức (Moral Lessons) — tự nhiên, KHÔNG ÉP

Bài học nảy sinh từ **hệ quả tự nhiên** của hành động, không phải thuyết giáo:

| Lĩnh vực | Trigger | Bài học |
|---|---|---|
| **Sức khỏe** | Dùng điện thoại trong tối → cận thị; ăn quá nhanh → đau bụng | Thói quen sinh hoạt lành mạnh |
| **Trung thực** | Nhặt ví → trả cảnh sát → được thưởng | Trung thực mang lại phần thưởng |
| **Môi trường** | Tình nguyện dọn bãi biển | Bảo vệ sinh thái |
| **Kỹ năng sống** | Đi khám răng, làm thủ tục sân bay | Kỹ năng thực tế cho trẻ |
| **Đúng giờ** | Đi trễ → bị cảnh cáo/mất việc | Tôn trọng thời gian |
| **Gọn gàng** | Không tìm thấy đồ → bài học sắp xếp | Ngăn nắp, tổ chức |
| **Hiếu thảo** | Chăm sóc Obaa-chan bệnh | Yêu thương gia đình |

### 7. Visual Style (KHÔNG MÔ TẢ TRONG SHOT)

Script **KHÔNG** mô tả style/aesthetic/lighting trong Visual. Chỉ mô tả **hành động vật lý** và **context/environment**.
Style (Rec.709, Modern Digital, Handheld camera) được xử lý riêng bởi template `style_prompt` và `video_constraint`.

> **CẤM trong Visual**: "warm lighting", "Kodak Portra", "cinematic look", "shallow DOF", "bokeh".
> **CHỈ viết**: camera setup + subject + action + context + audio.

### 8. Anthropomorphic Rules (BẮT BUỘC)

Tabinyan có HÌNH HÀI mèo nhưng HÀNH VI hoàn toàn như NGƯỜI:

| ❌ SAI (hành vi mèo) | ✅ ĐÚNG (hành vi người) |
|---|---|
| Cuộn tròn, ngồi xổm kiểu mèo | Ngồi thẳng lưng trên ghế, bắt chéo chân |
| Kẹp giữa hai bàn chân | Cầm bằng ngón tay/paw, một tay giữ tay kia thao tác |
| Cúi mặt xuống đĩa liếm | Dùng tay đưa thức ăn lên miệng, nhai, biểu cảm |
| Bốn chân, nhảy | Đi thẳng hai chân, bước đều, tay vung nhẹ |
| Cuộn tròn trên sàn | Nằm trên giường, đắp chăn, đầu trên gối |

### 9. Audio-First Storytelling

- **70%+ screen time là SFX-only** — để âm thanh môi trường kể chuyện
- **Foley là công cụ kể chuyện chính**: bước chân lạch bạch, dao thái cộc cộc, nước sôi sùng sục
- **Narrator tối thiểu**: 5-10 từ/câu, chỉ bổ sung khi hình ảnh không đủ
- **ZERO text on screen**: Không phụ đề, không label, không speech bubble

---

## References

- Template chi tiết: `docs/hl_meow_meow_template.md`
- Voiceover format guide: `docs/features/voiceover_script_format_guide.md`
- Pipeline: `ai/indexes/pipeline-map.md` (Pipeline 1: Script → Storyboard)
