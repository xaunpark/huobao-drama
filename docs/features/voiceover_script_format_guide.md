# 📝 Voice-over Script Format Guide
## Cú pháp kịch bản Structured Shot cho AI Director

---

## Tổng quan

Hệ thống AI Director nhận kịch bản voice-over đã được **phân shot thủ công** bằng marker `// SHOT XX`. AI **chỉ bổ sung** visual metadata (visual_description, camera, atmosphere...) — **không bao giờ tách lại** hoặc gộp shot.

**Luồng hoạt động:**

```
AI viết kịch bản (ChatGPT/Claude/Gemini)
    │  output script với // SHOT markers
    ▼
Paste vào hệ thống → chọn "Visual Unit" mode
    │
    ▼
AI Director → giữ nguyên N shots → chỉ enrich metadata
    │
    ▼
Shot list hoàn chỉnh (visual_description, camera, audio strategy...)
```

---

## Cú pháp Shot Marker

```
// SHOT [số] | [duration] | [shot_type] | [audio_mode]
       ↑ bắt buộc    ↑ tùy chọn     ↑ tùy chọn    ↑ tùy chọn
```

### Các mức chi tiết

```
// SHOT 01                                      ← Chỉ ranh giới
// SHOT 02 | 5s                                 ← + duration cố định
// SHOT 03 | 4s | CU                            ← + shot type
// SHOT 04 | 6s | MCU | dialogue_dominant        ← Full control
```

### Giá trị hợp lệ

| Tham số | Giá trị | Ghi chú |
|---------|---------|---------|
| **duration** | `2s` – `12s` | Không chỉ định → AI tự ước tính |
| **shot_type** | `ELS` `LS` `WS` `MWS` `MS` `MCU` `CU` `ECU` `title-card` `tactical-map` | Không chỉ định → AI tự chọn |
| **audio_mode** | `narrator_only` `dialogue_dominant` | Không chỉ định → AI infer từ [Tag] |

### Bảng shot_type chi tiết

| Viết tắt | Tên đầy đủ | Dùng khi nào |
|----------|-----------|-------------|
| `ELS` | Extreme Long Shot | Panorama, toàn cảnh môi trường |
| `LS` | Long Shot | Toàn thân nhân vật, bối cảnh rộng |
| `WS` | Wide Shot | Cảnh rộng, nhiều nhân vật |
| `MWS` | Medium Wide Shot | Nhân vật từ đầu gối trở lên |
| `MS` | Medium Shot | Nửa người, tương tác |
| `MCU` | Medium Close-Up | Ngực trở lên, biểu cảm + bối cảnh |
| `CU` | Close-Up | Khuôn mặt, chi tiết, cảm xúc |
| `ECU` | Extreme Close-Up | Mắt, bàn tay, vật thể cận cảnh |
| `title-card` | Title Card | Slide text, tiêu đề phần |
| `tactical-map` | Tactical Map | Bản đồ chiến thuật, infographic |

---

## Hệ thống Tag

Bên trong mỗi shot, sử dụng **tag** `[TAG]` để chỉ định nội dung. Dòng không có tag = **narrator** (lời kể).

### Bảng tag đầy đủ

| Tag | Loại | Ý nghĩa | Maps to field | Ảnh hưởng audio_mode? |
|-----|------|---------|--------------|----------------------|
| *(không tag)* | Audio | Lời narrator / kể chuyện | `script_segment` | → `narrator_only` |
| `[Tên NV]` | Audio | Dialogue nhân vật | `dialogue_text` | → `dialogue_dominant` |
| `[CROWD]` | Audio | Tiếng đám đông | `dialogue_text` + `dialogue_type=crowd` | → `dialogue_dominant` |
| `[SFX]` | Metadata | Hiệu ứng âm thanh | `sound_effect` | ❌ Không |
| `[BGM]` | Metadata | Nhạc nền / mood | `bgm_prompt` | ❌ Không |
| `[CAM]` | Metadata | Chỉ đạo camera | `movement` | ❌ Không |
| `[VFX]` | Metadata | Hiệu ứng hình ảnh | `visual_description` (append) | ❌ Không |
| `[NOTE]` | Metadata | Ghi chú đạo diễn | `reason_for_shot` (append) | ❌ Không |

> **Quy tắc quan trọng:** Chỉ có **Audio tags** (`[Character]`, `[CROWD]`) thay đổi `audio_mode`. Tất cả **Metadata tags** (`[SFX]`, `[BGM]`, `[CAM]`, `[VFX]`, `[NOTE]`) **KHÔNG** ảnh hưởng audio_mode.

### Chi tiết từng tag

#### `[SFX]` — Sound Effect
Hiệu ứng âm thanh cụ thể. Nhiều `[SFX]` trong cùng shot sẽ được nối bằng `; `.

```
[SFX] Tiếng bom nổ rền, mảnh vỡ bay khắp nơi
[SFX] Tiếng la hét hoảng loạn từ xa
```
→ `sound_effect: "Tiếng bom nổ rền, mảnh vỡ bay khắp nơi; Tiếng la hét hoảng loạn từ xa"`

#### `[BGM]` — Background Music
Mood/phong cách nhạc nền. Mỗi shot chỉ nên có 1 `[BGM]` (cái cuối cùng thắng).

```
[BGM] Tense military drums, low brass, building suspense
```
→ `bgm_prompt: "Tense military drums, low brass, building suspense"`

#### `[CAM]` — Camera Direction
Chỉ đạo camera cụ thể. Override trường `movement` mà AI tự tạo.

```
[CAM] Slow Ken Burns zoom in, slight pan left to right
[CAM] Handheld shake, POV crawling through tunnel
[CAM] Drone shot pulling up from ground level
```
→ `movement: "Slow Ken Burns zoom in, slight pan left to right"`

#### `[VFX]` — Visual Effect
Hiệu ứng hình ảnh đặc biệt, append vào `visual_description`.

```
[VFX] Sepia tone, film grain, aged document look
[VFX] Red highlight on tactical map, pulsing glow
[VFX] Slow motion, particles floating
```
→ Appended: `visual_description += " [VFX: Sepia tone, film grain, aged document look]"`

#### `[NOTE]` — Director's Note
Ghi chú đạo diễn về ý đồ, không hiển thị trên video. Append vào `reason_for_shot`.

```
[NOTE] Shot này cần tạo cảm giác ngột ngạt, claustrophobic
[NOTE] Tương phản với shot trước — từ bóng tối ra ánh sáng
[NOTE] Reference: Band of Brothers Episode 4
```
→ `reason_for_shot += " | Director note: Shot này cần tạo cảm giác ngột ngạt, claustrophobic"`

---

## Ví dụ hoàn chỉnh

### Ví dụ 1: Lịch sử quân sự (phong cách Simple History)

```
// SHOT 01 | 3s | title-card
## PHẦN 1: CƠN ÁC MỘNG DƯỚI LÒNG ĐẤT
[BGM] Dark ambient drone, low rumble

// SHOT 02 | 6s | CU
[SFX] Tiếng thở gấp gáp, nặng nề dội lại trong không gian hẹp.
[CAM] Handheld shake, POV nhìn từ mắt lính
Góc nhìn POV: Ánh đèn pin quét qua vách đất ẩm ướt, đầy rễ cây.
[NOTE] Mở đầu ấn tượng — claustrophobic, khán giả cảm nhận ngay sự ngột ngạt

// SHOT 03 | 5s | ELS
[BGM] Tense orchestral, war atmosphere
Sâu dưới những khu rừng rậm rạp của miền Nam Việt Nam,
một cuộc chiến bí mật đang diễn ra.
[VFX] Desaturated military palette, lifted blacks

// SHOT 04 | 5s | MS
[SFX] Tiếng đất cát lạo xạo khi có người trườn qua.
Trong không gian chật hẹp này, không có chỗ cho xe tăng hay máy bay ném bom.
Đây là thế giới của những "Tunnel Rat" — lính đánh chuột cống của quân đội Mỹ.
[CAM] Slow push-in, tightening frame

// SHOT 05 | 4s | MCU | dialogue_dominant
[Soldier] Cố lên nào... đừng kẹt lại lúc này.
[SFX] Tiếng tim đập thình thịch tăng dần.

// SHOT 06 | 5s | ECU
[CAM] Slow zoom into scorpion, hold on detail
Ánh đèn pin dừng lại ở một con bọ cạp đang bò trên trần hầm.
Chỉ một sơ suất nhỏ, bóng tối sẽ nuốt chửng lấy bạn mãi mãi.
[SFX] Single heartbeat thud, then silence

// SHOT 07 | 3s | title-card
## PHẦN 2: BỐI CẢNH LỊCH SỬ
[BGM] Transition to epic documentary strings

// SHOT 08 | 6s | ELS
[SFX] Tiếng bom B-52 nổ rền vang từ xa.
[VFX] Archival photo filter, slight vignette
Quân đội Bắc Việt và Việt Cộng đã xây dựng một mạng lưới địa đạo khổng lồ.
[CAM] Ken Burns slow zoom out, revealing scale

// SHOT 09 | 5s | tactical-map
[NOTE] Hiển thị bản đồ miền Nam Việt Nam, highlight hệ thống địa đạo
Hệ thống này trải dài hàng trăm dặm, vươn sâu vào các căn cứ của Mỹ.
[VFX] Red lines spreading across map like veins, pulsing animation

// SHOT 10 | 6s | MS
[SFX] Tiếng cuốc đất nện xuống đều đặn.
Địa đạo không chỉ là nơi trú ẩn, mà còn là bệnh viện, kho vũ khí và bếp ăn.

// SHOT 11 | 4s | WS
Bom B-52 nổ trên mặt đất gần như không thể chạm tới những tầng hầm sâu nhất.
[SFX] Muffled explosion from above, dirt falling

// SHOT 12 | 5s | MCU | dialogue_dominant
[NVA Officer] Các đồng chí, giữ im lặng. Quân Mỹ đang ở ngay trên đầu chúng ta.
[BGM] Music drops to near silence, tension

// SHOT 13 | 3s | WS | dialogue_dominant
[CROWD] Tiếng thì thầm xôn xao trong hầm tối.
[SFX] Faint footsteps from above

// SHOT 14 | 6s | MS
[CAM] Steady tracking shot along tunnel wall
Quân đội Mỹ đã thử bơm nước, dùng hơi cay và thậm chí là đánh thuốc nổ.
Nhưng cấu trúc chữ U và các vách ngăn kín đã vô hiệu hóa tất cả.
```

### Ví dụ 2: Câu chuyện Kinh Thánh (phong cách Minno)

```
// SHOT 01 | 5s | WS
[BGM] Gentle pastoral flute, warm and playful
Moses was a simple shepherd, taking care of his sheep in the desert.

// SHOT 02 | 3s | MS | dialogue_dominant
[Moses] Come on, little guys. Stay together!
[SFX] Sheep bleating

// SHOT 03 | 4s | MS
But Moses was also... a bit clumsy.

// SHOT 04 | 3s | MCU | dialogue_dominant
[Moses] Whoa! Oof!
[SFX] Tripping sound, dust cloud
[VFX] Cartoon dust poof effect

// SHOT 05 | 5s | WS
One day, while walking through the wilderness, he saw something incredible.
A bush was on fire — but it didn't burn up!
[SFX] Crackling fire, mystical hum

// SHOT 06 | 3s | CU | dialogue_dominant
[Moses] That's... odd. Why is it still green?
[CAM] Slow push-in on Moses' confused face

// SHOT 07 | 4s | MS
Then, a voice called out to him from the flames.
[BGM] Music shifts to reverent, ethereal choir
[VFX] Soft golden glow emanating from bush

// SHOT 08 | 3s | CU | dialogue_dominant
[God] Moses! Moses!

// SHOT 09 | 2s | MCU | dialogue_dominant
[Moses] Here I am!

// SHOT 10 | 5s | MS | dialogue_dominant
[God] Take off your sandals, for you are standing on holy ground.

// SHOT 11 | 3s | MCU | dialogue_dominant
[Moses] Oh! Sandals off! Right away!
[SFX] Sandals thrown, landing with a thud

// SHOT 12 | 5s | WS
God told Moses that He had seen the suffering of His people in Egypt
and that He was sending Moses to set them free.
[BGM] Music builds with hope, strings rising

// SHOT 13 | 4s | CU | dialogue_dominant
[Moses] Me? But I'm nobody! I can't even talk to people without stuttering!

// SHOT 14 | 4s | MS | dialogue_dominant
[God] I will be with you. I will teach you what to say.
[NOTE] Moment of divine reassurance — visual should convey warmth and power

// SHOT 15 | 5s | WS
And so began one of the greatest rescue missions in all of history.
[CROWD] murmuring with hope and wonder
[BGM] Triumphant, uplifting finale
```

---

## Quy tắc viết script

### 1. Tên nhân vật nhất quán
```
✅ [Moses] ... (xuyên suốt script)
❌ Lúc [Moses] lúc [Mô-sê]
```

### 2. Một dòng = một người nói
```
✅ [Moses] Let my people go!
✅ [God] I will be with you.

❌ [Moses] Let my people go! [God] No!
```

### 3. Dialogue nên ngắn gọn

| Loại | Độ dài | Ví dụ |
|------|--------|-------|
| **Reaction** | 1-3 từ | `[Soldier] Incredible!` |
| **Soft line** | 4-8 từ | `[Farmer] The harvest looks good!` |
| **Quote** | 1-2 câu | `[Caesar] I came, I saw, I conquered.` |
| **Full dialogue** | Max ~10 giây | Nhiều dòng `[A]` `[B]` liên tiếp |
| **Crowd** | Bất kỳ | `[CROWD] Long live the king!` |

### 4. Mỗi shot = 1 hình ảnh cụ thể
Nếu cảnh thay đổi → tạo shot mới. Nếu cùng 1 hình ảnh → giữ trong cùng shot.

### 5. Markdown headers được phép
`#`, `##`, `---`, `**bold**`, bảng `|...|` đều bị bỏ qua khi parse — dùng thoải mái để tổ chức script.

### 6. Yêu cầu tối thiểu
Cần **≥ 3 marker `// SHOT`** để hệ thống nhận diện chế độ Structured.

---

## Nhịp điệu khuyến nghị

### Theo phong cách kênh

| Phong cách | Narrator : Dialogue | Shot TB | Shots/phút |
|-----------|---------------------|---------|------------|
| **Simple History** (quân sự) | 75-85% : 15-25% | 4-5s | 10-14 |
| **Minno** (trẻ em) | 50-60% : 40-50% | 3-5s | 10-15 |
| **General documentary** | 60-70% : 30-40% | 4-6s | 10-12 |

### Theo giai đoạn câu chuyện

| Giai đoạn | % script | Đặc điểm |
|-----------|---------|----------|
| **Mở đầu** (20% đầu) | Chủ yếu narrator | Thiết lập bối cảnh, ít dialogue |
| **Phát triển** (20-60%) | Cân bằng | Câu chuyện chính, dialogue tăng dần |
| **Cao trào** (60-90%) | Nhiều dialogue | Trao đổi nhanh, kịch tính |
| **Kết** (10% cuối) | Quay lại narrator | Tổng kết, suy ngẫm |

---

## Thông số kỹ thuật

| Metric | Giá trị |
|--------|---------|
| Duration mỗi shot | 2-8 giây (lý tưởng 4-5s) |
| Max dialogue liên tục | 10-15 giây rồi quay lại narrator |
| Min gap narrator giữa 2 cụm dialogue | 2 shots |
| Từ narrator / phút | ~140 từ tiếng Anh |

### Bảng quy mô script

| Thời lượng | Số shots | Narrator | Dialogue | Tổng ~words |
|-----------|----------|----------|----------|-------------|
| 1 phút | 10-14 | ~120 từ | ~30 từ | ~150 |
| 2 phút | 20-28 | ~250 từ | ~60 từ | ~310 |
| 3 phút | 30-42 | ~380 từ | ~90 từ | ~470 |
| 5 phút | 50-70 | ~630 từ | ~150 từ | ~780 |
| 7 phút | 70-98 | ~880 từ | ~210 từ | ~1090 |

---

## Mẫu prompt cho AI viết kịch bản

> Dán prompt sau vào AI bên ngoài (ChatGPT, Claude, Gemini...) để nó output script đúng format:

````
Viết kịch bản voice-over cho video [chủ đề] dài khoảng [N] phút.

**FORMAT BẮT BUỘC — Structured Shot:**
Phân script thành từng shot bằng marker `// SHOT XX`. Mỗi shot = 1 hình ảnh cụ thể.

**CÚ PHÁP SHOT MARKER:**
```
// SHOT [số] | [duration]s | [shot_type] | [audio_mode]
```
- Chỉ `// SHOT [số]` là bắt buộc. Các tham số khác tùy chọn.
- shot_type: ELS, LS, WS, MWS, MS, MCU, CU, ECU, title-card, tactical-map
- audio_mode: narrator_only (mặc định) hoặc dialogue_dominant

**HỆ THỐNG TAG (viết bên trong mỗi shot):**

Audio tags (thay đổi audio_mode):
- Dòng không tag = lời narrator → audio_mode = narrator_only
- `[Tên Nhân Vật]` = dialogue → audio_mode = dialogue_dominant
- `[CROWD]` = tiếng đám đông → audio_mode = dialogue_dominant

Metadata tags (KHÔNG thay đổi audio_mode):
- `[SFX]` = hiệu ứng âm thanh (ví dụ: `[SFX] Tiếng bom nổ rền`)
- `[BGM]` = nhạc nền / mood (ví dụ: `[BGM] Tense military drums, building suspense`)
- `[CAM]` = chỉ đạo camera (ví dụ: `[CAM] Slow Ken Burns zoom in`)
- `[VFX]` = hiệu ứng hình ảnh (ví dụ: `[VFX] Sepia tone, film grain`)
- `[NOTE]` = ghi chú đạo diễn (ví dụ: `[NOTE] Tạo cảm giác ngột ngạt`)

**QUY TẮC:**
1. Tên nhân vật phải nhất quán trong suốt script
2. Dialogue nên ngắn: reaction 1-3 từ, quote 1-2 câu, max ~10 giây
3. Mỗi shot nên có 1-3 câu narrator (~4-6 giây)
4. Cứ 3-4 shot narrator → chèn 1 shot dialogue để tạo nhịp
5. Dùng ## Headers để tổ chức script (sẽ bị bỏ qua khi parse)
6. Mỗi shot = 1 hình ảnh/cảnh cụ thể — nếu hình ảnh thay đổi → shot mới
7. [BGM] chỉ cần đặt ở shot đầu tiên của mỗi đoạn mood mới (không cần lặp lại)
8. [CAM] chỉ dùng khi muốn camera đặc biệt — mặc định AI sẽ tự chọn

**VÍ DỤ:**
```
// SHOT 01 | 3s | title-card
## PHẦN 1: TIÊU ĐỀ
[BGM] Dark ambient drone, low rumble

// SHOT 02 | 6s | ELS
[SFX] Tiếng gió sa mạc thổi qua
Giữa vùng sa mạc khô cằn, một đoàn quân đang hành quân dưới nắng gắt.
[VFX] Desaturated palette, heat haze effect

// SHOT 03 | 5s | MS
[CAM] Slow tracking shot following soldiers
Họ đã đi suốt 3 ngày đêm không nghỉ.

// SHOT 04 | 4s | MCU | dialogue_dominant
[Officer] Tiến lên! Không được dừng lại!
[SFX] Tiếng bước chân nặng nề trên cát

// SHOT 05 | 5s | WS
[NOTE] Tạo cảm giác về sự rộng lớn và cô đơn
Nhưng phía trước, một thử thách lớn hơn đang chờ đợi.
[BGM] Music builds with ominous brass
```

Bây giờ hãy viết script cho: [CHỦ ĐỀ CỤ THỂ]
````

---

## Lưu ý kỹ thuật

1. **Ngôn ngữ:** Script viết bằng **bất kỳ ngôn ngữ nào**. Hệ thống auto output trường kỹ thuật bằng tiếng Anh.

2. **Auto-detect:** Hệ thống tự nhận diện khi có ≥ 3 marker `// SHOT` → kích hoạt Structured mode.

3. **Nhân vật:** Tên trong `[Tag]` nên trùng với Character đã tạo trong project để auto-map ID.

4. **Tag alias:** `[CAM]` = `[CAMERA]`, `[NOTE]` = `[DIR]` — đều hoạt động.

5. **Nhiều SFX:** Nhiều `[SFX]` trong cùng shot sẽ được nối bằng `; `. Các tag khác chỉ lấy giá trị cuối cùng.

6. **Dialogue nhiều dòng:**
   ```
   [God] Take off your sandals, for you are standing on holy ground.
   [God] I am the God of your father.
   ```
   Sẽ được nối thành: `dialogue_text: "Take off your sandals... \n I am the God..."`

7. **Duration tự động:** Không chỉ định → AI ước tính ~16 từ tiếng Anh ≈ 7 giây (~140 WPM).

8. **Thứ tự tag trong shot:** Không quan trọng. Có thể đặt `[SFX]` đầu, `[BGM]` cuối, hoặc xen kẽ với narrator.
