# Shot Format — Output Template

Mỗi SHOT trong output script phải tuân theo format chuẩn này.
Mục tiêu: **mỗi SHOT có thể dùng trực tiếp làm prompt cho AI image/video generation**.

---

## Định nghĩa SHOT

SHOT = một đoạn hình ảnh liên tục giữa hai lần CUT.

**Tạo SHOT mới khi:**
- Có CUT / chuyển cảnh
- Đổi góc camera (wide → close, POV → side…)
- Đổi hành động chính
- Đổi vị trí / không gian
- Nhảy thời gian

**KHÔNG tạo shot mới nếu:**
- Chỉ có chuyển động nhẹ trong cùng khung hình
- Camera pan/zoom nhưng vẫn cùng hành động

---

## Format mỗi SHOT

```
SHOT XX | duration | shot_type
[BGM: ...]
[SFX: ...]
[DIA: ...]
[CAM: ...]
[HUMAN REACTION: pattern_name]

Visual: ...
```

### Quy tắc từng field

#### `SHOT XX | duration | shot_type`
- `XX`: số thứ tự (01, 02, 03...)
- `duration`: thời lượng (3s, 4s, 5s...)
- `shot_type`: MS (Medium Shot), WS (Wide Shot), CU (Close-Up), MWS (Medium Wide Shot), INSERT (Detail Shot), ECU (Extreme Close-Up)

#### `[BGM]` — Background Music
- Chỉ ghi khi mood thay đổi hoặc ở SHOT đầu tiên
- Ví dụ: `[BGM] Gentle acoustic guitar, warm and playful`

#### `[SFX]` — Sound Effects
- **BẮT BUỘC cho MỌI shot** — Foley là công cụ kể chuyện chính
- Cụ thể, không chung chung. Dùng onomatopoeia Nhật khi phù hợp
- Ví dụ: `[SFX] Knife on cutting board — "Ton ton ton"; oil sizzling — "Juu~"]`

#### `[DIA]` — Dialogue
- Tabinyan: chỉ thán từ ngắn + pitch description
- Người Nhật: Keigo, ngắn gọn
- Narrator: 5-10 từ mỗi câu, tiếng Anh
- Ví dụ: `[DIA — Tabinyan, very high-pitched, chirpy] "Oishii!"`
- Ví dụ: `[DIA — Shopkeeper, polite Keigo] "Irasshaimase~"`
- Ví dụ: `[NARRATOR] Tabinyan begins her morning routine.`

#### `[CAM]` — Camera Direction
- Mô tả camera operator đang LÀM GÌ vật lý
- Handheld documentary style: rung nhẹ, theo sau, corrective pan
- Ví dụ: `[CAM] Handheld follow from behind, walking pace, natural bob — Tabinyan turns corner, camera overshoots slightly then corrects`
- Ví dụ: `[CAM] Static low angle, cat-eye level, slight breathing sway`

#### `[HUMAN REACTION]` — Reaction Pattern
- Chỉ khi có người tương tác
- Giá trị: `gentle_mentor`, `admiring_customer`, `formal_respect`, `compassionate_caregiver`, `silent_admirer`, `amused_observer`

#### `Visual:` — Visual Prompt (QUAN TRỌNG NHẤT)

**Công thức BẮT BUỘC:**
```
[CAMERA SETUP] + [SUBJECT + VỊ TRÍ] + [ACTION CHÍNH] + [CONTEXT/ENVIRONMENT]
```

**Luật viết Visual:**

1. **Camera LUÔN đứng đầu**: shot type + angle + movement
2. **Subject rõ ràng**: ai xuất hiện, vị trí trong frame (center, left, background)
3. **Action = 1 hành động chính**: dùng động từ vật lý (walks, reaches, lifts, chews...)
4. **Context**: môi trường (Japanese kitchen, supermarket aisle...), vật thể xung quanh
5. **Blocking & Staging**: nhân vật đứng ở đâu, di chuyển hướng nào, khoảng cách giữa các nhân vật

**CẤM TUYỆT ĐỐI trong Visual:**
- ❌ Mô tả style hình ảnh (lighting, aesthetic, color grading, lens...)
- ❌ Mô tả cảm xúc trừu tượng ("sad", "tense", "warm atmosphere")
- ❌ Multi-action trong 1 shot
- ❌ Thêm nội dung không có trong kịch bản
- ❌ Mô tả ngoại hình nhân vật chi tiết (fur color, eye shape) — dùng reference image
- ❌ Viết kiểu văn học, ẩn dụ

**✅ Ưu tiên: mô tả vật lý (physical description), front-load thông tin quan trọng đầu câu.**

---

## Shot Distribution (cho episode standard 1-2 phút)

| Shot Type | Tỷ lệ | Khi nào |
|---|---|---|
| Medium Shot (MS) | ~40% | Tabinyan thực hiện hành động chính |
| Close-Up (CU) | ~20% | Biểu cảm, chi tiết thức ăn |
| Wide Shot (WS) | ~15% | Establishing, full environment |
| Medium Wide (MWS) | ~15% | Two-shot Tabinyan + Human |
| Insert/Detail | ~10% | Paw làm việc, prop detail |

## Camera Angle Distribution

| Angle | Tỷ lệ | Mô tả |
|---|---|---|
| Cat-eye level (Low) | 70% | Camera ở tầm mắt mèo (~45cm) |
| Human eye-level | 15% | Khi tương tác với người |
| High angle | 10% | Nhìn xuống mèo trên sàn/giường |
| Overhead/Top-down | 5% | Food prep, table-top |

## Shot Pacing

| Loại | Duration | Ghi chú |
|---|---|---|
| Action montage | 2-3s | Cooking steps, working |
| Reaction/emotion | 3-4s | Eating, smiling, sleeping |
| Establishing | 3-5s | Wide shot environment |
| Transition | 100% hard cut | NO dissolve, NO fade |

## Interaction Pattern

Khi Tabinyan tương tác với người, tuân theo pattern:
```
MS Tabinyan action → MWS two-shot → CU human reaction (2s) → CU Tabinyan reaction (2s)
```

---

## Ví dụ hoàn chỉnh

### Episode: "Tabinyan làm đầu bếp ramen"

```
SHOT 01 | 4s | WS
[BGM] Gentle acoustic guitar, warm morning vibe
[SFX] Distant street sounds, birds chirping, bicycle bell
[CAM] Handheld follow, approaching from street level, slight bob

Visual: Wide shot at cat-eye level, camera approaching a small Japanese ramen shop from across a narrow shotengai street. Noren curtains sway gently in the doorway. A sandwich board sign sits outside. Morning pedestrians walk past in the background.

SHOT 02 | 3s | MS
[SFX] Fabric rustling, tiny footsteps on wooden floor — "Pata pata"
[CAM] Static, cat-eye level, slight breathing sway

Visual: Medium shot inside the kitchen. Tabinyan stands upright on a wooden step stool, center frame, tying a miniature white apron behind her back with both paws. A row of ramen bowls lines the shelf behind her.

SHOT 03 | 3s | INSERT
[SFX] Knife on cutting board — "Ton ton ton ton"
[CAM] Static overhead, slight zoom

Visual: Top-down insert shot. Two small paws grip a miniature knife, slicing green onions on a tiny wooden cutting board. Thin rings of scallion scatter across the board.

SHOT 04 | 3s | CU
[SFX] Oil sizzling — "Juu~"; gentle bubbling
[DIA — Tabinyan, very high-pitched, focused] "Mmm~"
[CAM] Static, eye-level, slight push-in

Visual: Close-up of Tabinyan's face from the front, slightly low angle. She looks down intently at a small pot, brow slightly furrowed in concentration. Steam rises from below frame.

SHOT 05 | 4s | MS
[SFX] Ceramic bowl placed on counter — "Koton"; ladle scooping broth — "Chapu"
[CAM] Handheld tracking left, following paw movement

Visual: Medium shot, side angle. Tabinyan stands at the counter, carefully ladling steaming broth from a pot into a small ceramic ramen bowl with both paws. Noodles already sit in the bowl. She places a single slice of chashu pork on top with precise, deliberate placement.

SHOT 06 | 3s | MWS
[SFX] Gentle footsteps approaching — adult shoes
[DIA — Customer, polite Keigo] "Sumimasen~ Ramen hitotsu onegaishimasu."
[HUMAN REACTION: admiring_customer]
[CAM] Static, slightly wider, cat-eye level

Visual: Medium wide two-shot. A Japanese woman in a beige cardigan sits at the counter, bending forward to be at Tabinyan's eye level. Tabinyan stands on the opposite side of the counter holding the completed ramen bowl with both paws, presenting it forward.

SHOT 07 | 3s | CU
[SFX] Soft gasp
[HUMAN REACTION: admiring_customer]
[CAM] Static, eye-level

Visual: Close-up of the customer's face. She receives the bowl with both hands, pauses, tilts her head slightly with a melting smile as she looks down at the ramen, then glances at Tabinyan.

SHOT 08 | 2s | CU
[DIA — Tabinyan, very high-pitched, proud chirp] "Meow!"
[SFX] Tiny satisfied exhale
[CAM] Static, cat-eye level

Visual: Close-up of Tabinyan's face, center frame. Wide sparkling eyes, mouth slightly open in a proud little chirp. She stands straight with one paw on her hip.

SHOT 09 | 4s | MS
[SFX] Chopsticks picking up noodles — "Zuru zuru"; satisfied eating sounds
[DIA — Customer, eyes closing] "Oishii~"
[CAM] Handheld, slight push-in during reaction

Visual: Medium shot from Tabinyan's POV across the counter. The customer lifts noodles with chopsticks, brings them to her mouth, chews. Her eyes close and she smiles deeply.

SHOT 10 | 3s | CU
[DIA — Tabinyan, very high-pitched, triumphant] "Yatta!"
[SFX] Tiny fist pump — fabric rustling
[CAM] Static, cat-eye level, slight breathing sway

Visual: Close-up of Tabinyan doing a small fist pump with one paw, the other paw resting on her apron. Eyes sparkling, mouth open in a wide, satisfied grin.
```
