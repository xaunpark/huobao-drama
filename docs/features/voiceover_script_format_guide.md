# 📝 Voice-over Script Format Guide
## Cú pháp kịch bản Structured Shot cho AI Director

---

## Tổng quan

Hệ thống AI Director nhận kịch bản voice-over đã được **phân shot thủ công** bằng marker `// SHOT XX`. AI **chỉ bổ sung** visual metadata — **không bao giờ tách lại** hoặc gộp shot.

**Quy tắc cốt lõi:** Mọi dòng nội dung PHẢI có tag. Không có dòng "tự hiểu" — tag tường minh đảm bảo trích xuất chính xác.

**Luồng hoạt động:**

```
AI viết kịch bản (ChatGPT/Claude/Gemini)
    │  output script với // SHOT markers + [TAG] tường minh
    ▼
Paste vào hệ thống → chọn "Visual Unit" mode
    │
    ▼
AI Director → giữ nguyên N shots → enrich metadata
    │  Trích xuất: NarratorScript, DialogueText, SFX, BGM, CAM...
    ▼
Shot list hoàn chỉnh → Video prompt có narrator context
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
| `ELS` | Extreme Long Shot | Panorama, toàn cảnh |
| `LS` | Long Shot | Toàn thân nhân vật |
| `WS` | Wide Shot | Cảnh rộng, nhiều nhân vật |
| `MWS` | Medium Wide Shot | Đầu gối trở lên |
| `MS` | Medium Shot | Nửa người |
| `MCU` | Medium Close-Up | Ngực trở lên |
| `CU` | Close-Up | Khuôn mặt, cảm xúc |
| `ECU` | Extreme Close-Up | Mắt, bàn tay, chi tiết |
| `title-card` | Title Card | Slide text |
| `tactical-map` | Tactical Map | Bản đồ, infographic |

---

## Hệ thống Tag

> **QUAN TRỌNG:** Mọi dòng nội dung bên trong shot PHẢI có tag. Dùng `[NARRATOR]` cho lời kể — KHÔNG để dòng trống không tag.

### Bảng tag đầy đủ

| Tag | Loại | Ý nghĩa | Maps to field | Ảnh hưởng audio_mode? |
|-----|------|---------|--------------|----------------------|
| `[NARRATOR]` | Audio | Lời kể / voice-over | `narrator_script` | → `narrator_only` |
| `[Tên NV]` | Audio | Dialogue nhân vật | `dialogue_text` | → `dialogue_dominant` |
| `[CROWD]` | Audio | Tiếng đám đông | `dialogue_text` + `dialogue_type=crowd` | → `dialogue_dominant` |
| `[SFX]` | Metadata | Hiệu ứng âm thanh | `sound_effect` | ❌ Không |
| `[BGM]` | Metadata | Nhạc nền / mood | `bgm_prompt` | ❌ Không |
| `[CAM]` | Metadata | Chỉ đạo camera | `movement` | ❌ Không |
| `[VFX]` | Metadata | Hiệu ứng hình ảnh | `visual_description` (append) | ❌ Không |
| `[NOTE]` | Metadata | Ghi chú đạo diễn | `reason_for_shot` (append) | ❌ Không |

### Cách trích xuất dữ liệu

Hệ thống trích xuất từ script thành 2 trường riêng biệt:

| Trường | Chứa gì | Dùng để |
|--------|---------|---------|
| **`script_segment`** | `[NARRATOR]` + `[Character]` + `[CROWD]` text | Hiển thị trong UI, script reference |
| **`narrator_script`** | CHỈ `[NARRATOR]` text (không tag) | Inject vào video prompt, TTS tương lai |

**Ví dụ trích xuất:**
```
// SHOT 05 | 5s | MS
[SFX] Tiếng đất cát lạo xạo
[NARRATOR] Đây là thế giới của những Tunnel Rat.
[CAM] Slow push-in
[Soldier] Cố lên nào!
```

→ `script_segment`: `"Đây là thế giới của những Tunnel Rat.\n[Soldier] Cố lên nào!"`
→ `narrator_script`: `"Đây là thế giới của những Tunnel Rat."`
→ `sound_effect`: `"Tiếng đất cát lạo xạo"`
→ `movement`: `"Slow push-in"`
→ `dialogue`: `"Cố lên nào!"`

> **[SFX], [BGM], [CAM], [VFX], [NOTE]** KHÔNG xuất hiện trong `script_segment` — chúng được trích xuất vào trường riêng.

---

### Chi tiết từng tag

#### `[NARRATOR]` — Voice-over Narration ⭐
Lời kể chuyện. Video minh họa hình ảnh, nhân vật **không mở miệng**.

```
[NARRATOR] Sâu dưới những khu rừng rậm rạp của miền Nam Việt Nam,
[NARRATOR] một cuộc chiến bí mật đang diễn ra.
```
→ `narrator_script: "Sâu dưới những khu rừng rậm rạp... một cuộc chiến bí mật đang diễn ra."`
→ Video prompt: `"Narration (voice-over): Sâu dưới... mouth closed, silent expression"`

#### `[Character]` — Dialogue
Thoại nhân vật cụ thể. Nhân vật **mở miệng nói**.

```
[Soldier] Cố lên nào... đừng kẹt lại lúc này.
```
→ `dialogue: "Cố lên nào... đừng kẹt lại lúc này."`
→ Video prompt: `"Dialogue: Cố lên nào... actively speaking, lip-syncing"`

#### `[SFX]` — Sound Effect
Nhiều `[SFX]` trong cùng shot → nối bằng `; `.

```
[SFX] Tiếng bom nổ rền, mảnh vỡ bay
[SFX] Tiếng la hét từ xa
```
→ `sound_effect: "Tiếng bom nổ rền, mảnh vỡ bay; Tiếng la hét từ xa"`

#### `[BGM]` — Background Music
```
[BGM] Tense military drums, low brass, building suspense
```

#### `[CAM]` — Camera Direction
```
[CAM] Slow Ken Burns zoom in, slight pan left to right
[CAM] Handheld shake, POV crawling through tunnel
```

#### `[VFX]` — Visual Effect
```
[VFX] Sepia tone, film grain, aged document look
```

#### `[NOTE]` — Director's Note
```
[NOTE] Shot này cần tạo cảm giác ngột ngạt, claustrophobic
[NOTE] Reference: Band of Brothers Episode 4
```

---

## Ví dụ hoàn chỉnh

### Ví dụ 1: Lịch sử quân sự (phong cách Simple History)

```
// SHOT 01 | 3s | title-card
[BGM] Dark ambient drone, low rumble
[NOTE] Title card cho phần mở đầu

// SHOT 02 | 6s | CU
[SFX] Tiếng thở gấp gáp, nặng nề dội lại trong không gian hẹp.
[CAM] Handheld shake, POV nhìn từ mắt lính
[NARRATOR] Góc nhìn POV: Ánh đèn pin quét qua vách đất ẩm ướt, đầy rễ cây.
[NOTE] Mở đầu ấn tượng — claustrophobic

// SHOT 03 | 5s | ELS
[BGM] Tense orchestral, war atmosphere
[NARRATOR] Sâu dưới những khu rừng rậm rạp của miền Nam Việt Nam,
[NARRATOR] một cuộc chiến bí mật đang diễn ra.
[VFX] Desaturated military palette, lifted blacks

// SHOT 04 | 5s | MS
[SFX] Tiếng đất cát lạo xạo khi có người trườn qua.
[NARRATOR] Trong không gian chật hẹp này, không có chỗ cho xe tăng hay máy bay ném bom.
[NARRATOR] Đây là thế giới của những "Tunnel Rat" — lính đánh chuột cống của quân đội Mỹ.
[CAM] Slow push-in, tightening frame

// SHOT 05 | 4s | MCU | dialogue_dominant
[Soldier] Cố lên nào... đừng kẹt lại lúc này.
[SFX] Tiếng tim đập thình thịch tăng dần.

// SHOT 06 | 5s | ECU
[CAM] Slow zoom into scorpion, hold on detail
[NARRATOR] Ánh đèn pin dừng lại ở một con bọ cạp đang bò trên trần hầm.
[NARRATOR] Chỉ một sơ suất nhỏ, bóng tối sẽ nuốt chửng lấy bạn mãi mãi.
[SFX] Single heartbeat thud, then silence

// SHOT 07 | 3s | title-card
[BGM] Transition to epic documentary strings
[NOTE] Title card cho phần 2

// SHOT 08 | 6s | ELS
[SFX] Tiếng bom B-52 nổ rền vang từ xa.
[VFX] Archival photo filter, slight vignette
[NARRATOR] Quân đội Bắc Việt và Việt Cộng đã xây dựng một mạng lưới địa đạo khổng lồ.
[CAM] Ken Burns slow zoom out, revealing scale

// SHOT 09 | 5s | tactical-map
[NOTE] Hiển thị bản đồ miền Nam Việt Nam, highlight hệ thống địa đạo
[NARRATOR] Hệ thống này trải dài hàng trăm dặm, vươn sâu vào các căn cứ của Mỹ.
[VFX] Red lines spreading across map like veins, pulsing animation

// SHOT 10 | 6s | MS
[SFX] Tiếng cuốc đất nện xuống đều đặn.
[NARRATOR] Địa đạo không chỉ là nơi trú ẩn, mà còn là bệnh viện, kho vũ khí và bếp ăn.

// SHOT 11 | 4s | WS
[NARRATOR] Bom B-52 nổ trên mặt đất gần như không thể chạm tới những tầng hầm sâu nhất.
[SFX] Muffled explosion from above, dirt falling

// SHOT 12 | 5s | MCU | dialogue_dominant
[NVA Officer] Các đồng chí, giữ im lặng. Quân Mỹ đang ở ngay trên đầu chúng ta.
[BGM] Music drops to near silence, tension

// SHOT 13 | 3s | WS | dialogue_dominant
[CROWD] Tiếng thì thầm xôn xao trong hầm tối.
[SFX] Faint footsteps from above

// SHOT 14 | 6s | MS
[CAM] Steady tracking shot along tunnel wall
[NARRATOR] Quân đội Mỹ đã thử bơm nước, dùng hơi cay và thậm chí là đánh thuốc nổ.
[NARRATOR] Nhưng cấu trúc chữ U và các vách ngăn kín đã vô hiệu hóa tất cả.
```

### Ví dụ 2: Câu chuyện Kinh Thánh (phong cách Minno)

```
// SHOT 01 | 5s | WS
[BGM] Gentle pastoral flute, warm and playful
[NARRATOR] Moses was a simple shepherd, taking care of his sheep in the desert.

// SHOT 02 | 3s | MS | dialogue_dominant
[Moses] Come on, little guys. Stay together!
[SFX] Sheep bleating

// SHOT 03 | 4s | MS
[NARRATOR] But Moses was also... a bit clumsy.

// SHOT 04 | 3s | MCU | dialogue_dominant
[Moses] Whoa! Oof!
[SFX] Tripping sound, dust cloud
[VFX] Cartoon dust poof effect

// SHOT 05 | 5s | WS
[NARRATOR] One day, while walking through the wilderness, he saw something incredible.
[NARRATOR] A bush was on fire — but it didn't burn up!
[SFX] Crackling fire, mystical hum

// SHOT 06 | 3s | CU | dialogue_dominant
[Moses] That's... odd. Why is it still green?
[CAM] Slow push-in on Moses' confused face

// SHOT 07 | 4s | MS
[NARRATOR] Then, a voice called out to him from the flames.
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
[NARRATOR] God told Moses that He had seen the suffering of His people in Egypt
[NARRATOR] and that He was sending Moses to set them free.
[BGM] Music builds with hope, strings rising

// SHOT 13 | 4s | CU | dialogue_dominant
[Moses] Me? But I'm nobody! I can't even talk to people without stuttering!

// SHOT 14 | 4s | MS | dialogue_dominant
[God] I will be with you. I will teach you what to say.
[NOTE] Moment of divine reassurance — visual should convey warmth and power

// SHOT 15 | 5s | WS
[NARRATOR] And so began one of the greatest rescue missions in all of history.
[CROWD] murmuring with hope and wonder
[BGM] Triumphant, uplifting finale
```

### Ví dụ 3: Khoa học / Vũ trụ (phong cách Kurzgesagt)

```
// SHOT 01 | 4s | ELS
[BGM] Ethereal synth, cosmic ambient
[NARRATOR] Imagine a star so massive, it makes our Sun look like a marble.
[VFX] Size comparison animation — Sun shrinking next to a hypergiant

// SHOT 02 | 5s | MS
[NARRATOR] UY Scuti, the largest known star, is 1,700 times wider than the Sun.
[NARRATOR] If placed at the center of our solar system, it would swallow Jupiter.
[CAM] Slow zoom out revealing scale

// SHOT 03 | 4s | CU
[SFX] Deep cosmic rumble, low frequency vibration
[NARRATOR] But here's the twist — this giant is dying.
[VFX] Star surface pulsating, unstable red glow

// SHOT 04 | 3s | ECU
[NOTE] Dramatic pause shot — build tension before the reveal
[SFX] Heartbeat-like pulsation from the star

// SHOT 05 | 6s | ELS
[BGM] Music swells with epic brass and strings
[NARRATOR] When UY Scuti finally collapses, the explosion will be so powerful
[NARRATOR] that it could outshine an entire galaxy for weeks.
[VFX] Supernova explosion animation, shockwave expanding outward
[SFX] Massive explosion reverberating through space

// SHOT 06 | 5s | WS
[NARRATOR] This is called a supernova — the most violent event in the universe.
[CAM] Slow pull-out from explosion debris field

// SHOT 07 | 4s | MS
[NARRATOR] But from that destruction, something beautiful emerges.
[SFX] Gentle chimes, rebirth sound
[BGM] Transition to hopeful, wonder-filled melody

// SHOT 08 | 5s | CU
[NARRATOR] The scattered elements — iron, gold, carbon —
[NARRATOR] become the building blocks of new planets and, eventually, life.
[VFX] Colorful nebula forming, particles coalescing into matter

// SHOT 09 | 4s | ECU
[CAM] Extreme close-up on human hand, then zoom out to reveal stars
[NARRATOR] Every atom in your body was forged inside a dying star.
[NOTE] Emotional peak — connection between cosmic and personal

// SHOT 10 | 5s | ELS
[NARRATOR] You are, quite literally, made of stardust.
[BGM] Final note — sustained, ethereal, fading out
[VFX] Camera pulls out through galaxy, human figure fades into starfield
```

### Ví dụ 4: Kinh dị / Bí ẩn (phong cách MrBallen / Dark History)

```
// SHOT 01 | 4s | ECU
[BGM] Single low drone, unsettling
[SFX] Static buzz, old radio tuning
[NARRATOR] On the night of February 2nd, 1959, something went terribly wrong.

// SHOT 02 | 5s | ELS
[NARRATOR] Nine experienced hikers entered the Ural Mountains of Russia.
[NARRATOR] None of them would make it out alive.
[VFX] Desaturated blue-gray palette, snow particles drifting
[CAM] Slow aerial drone shot over snow-covered mountain

// SHOT 03 | 3s | title-card
[NOTE] Title card: "The Dyatlov Pass Incident"
[BGM] Music drops to near silence

// SHOT 04 | 5s | MS
[NARRATOR] The group was led by Igor Dyatlov, 23, an engineering student
[NARRATOR] with years of trekking experience.
[SFX] Wind howling, boots crunching on snow

// SHOT 05 | 4s | CU
[CAM] Close-up on compass, then pan to snowy horizon
[NARRATOR] They had planned a Category III expedition —
[NARRATOR] the most difficult grade in the Soviet system.

// SHOT 06 | 5s | WS
[NARRATOR] On the evening of February 1st, they set up camp on the slope
[NARRATOR] of a mountain the Mansi people called "Kholat Syakhl."
[SFX] Tent fabric flapping in wind
[NOTE] "Kholat Syakhl" = "Dead Mountain" in Mansi language — emphasize irony

// SHOT 07 | 6s | MS
[BGM] Tension building — sparse piano, dissonant
[NARRATOR] Sometime between midnight and the early morning hours,
[NARRATOR] something caused all nine hikers to slash their tent open from the inside
[NARRATOR] and flee into -30°C temperatures... without shoes.
[SFX] Fabric ripping, panicked breathing

// SHOT 08 | 3s | CU | dialogue_dominant
[Rescuer] Это невозможно... кто бы сделал это?
[NOTE] Rescuer dialogue in Russian — translation: "This is impossible... who would do this?"

// SHOT 09 | 5s | ECU
[CAM] Slow zoom on footprints in snow leading away from torn tent
[NARRATOR] When search teams finally found the camp, the scene made no sense.
[SFX] Eerie wind, distant wolves howling

// SHOT 10 | 5s | tactical-map
[NARRATOR] The bodies were found scattered across a mile-wide area.
[NARRATOR] Some were barely clothed. Others had catastrophic internal injuries
[NARRATOR] — but no external wounds.
[VFX] Map overlay showing body locations, red markers with distances
[BGM] Dark crescendo, building dread
```

### Ví dụ 5: Văn minh cổ đại (phong cách Full Documentary)

```
## PHẦN 1: BÌNH MINH CỦA NỀN VĂN MINH

// SHOT 01 | 3s | title-card
[BGM] Ancient ethnic flute, mysterious and reverent
[NOTE] Opening title — golden text on dark background

// SHOT 02 | 6s | ELS
[NARRATOR] In the heart of North Africa, along the banks of the mighty Nile River,
[NARRATOR] one of humanity's greatest civilizations rose from the desert sands.
[SFX] Wind blowing across desert dunes
[CAM] Aerial drone shot, sweeping across Nile delta at sunrise
[VFX] Warm golden hour lighting, slight lens flare

// SHOT 03 | 5s | WS
[NARRATOR] Ancient Egypt began as scattered farming villages along the fertile floodplains.
[NARRATOR] Each year, the Nile would overflow its banks, depositing rich black soil.
[SFX] Gentle water flowing, birds calling

// SHOT 04 | 4s | MS | dialogue_dominant
[Farmer] Look at this soil — black as night! The gods have blessed us again!
[SFX] Digging sounds, tools clinking

// SHOT 05 | 3s | WS | dialogue_dominant
[CROWD] Excited chatter and celebration
[SFX] Clapping, joyful voices

// SHOT 06 | 5s | MS
[NARRATOR] The Egyptians called their country "Kemet" — the Black Land —
[NARRATOR] a name born from the dark, life-giving mud.
[CAM] Slow pan across fertile fields, farmers working

## PHẦN 2: SỰ THỐNG NHẤT

// SHOT 07 | 6s | LS
[BGM] Epic percussion, war drums building
[NARRATOR] Around 3100 BC, a powerful warrior king named Narmer
[NARRATOR] united Upper and Lower Egypt into a single kingdom.
[SFX] War drums, horns blowing in the distance

// SHOT 08 | 4s | MCU | dialogue_dominant
[Narmer] From this day forward, Upper and Lower Egypt are ONE!
[SFX] Sword raised, crowd roaring

// SHOT 09 | 3s | WS | dialogue_dominant
[CROWD] Cheering and chanting, thunderous applause

// SHOT 10 | 5s | CU
[NARRATOR] He wore the double crown — the white crown of the south
[NARRATOR] merged with the red crown of the north —
[NARRATOR] symbolizing the unification of two worlds.
[CAM] Slow push-in on crown detail, then tilt up to reveal face

## PHẦN 3: THỜI ĐẠI KIM TỰ THÁP

// SHOT 11 | 5s | ELS
[BGM] Majestic orchestral, awe-inspiring
[NARRATOR] The age of the pyramids followed.
[NARRATOR] At Giza, tens of thousands of workers hauled massive limestone blocks.
[SFX] Stones grinding, ropes creaking
[CAM] Extreme wide revealing pyramid under construction

// SHOT 12 | 4s | MS | dialogue_dominant
[Worker 1] Pull! Pull harder!

// SHOT 13 | 3s | MS | dialogue_dominant
[Worker 2] It's too heavy!

// SHOT 14 | 4s | MCU | dialogue_dominant
[Overseer] The Pharaoh demands perfection! Keep moving!
[SFX] Whip crack, stones scraping

// SHOT 15 | 6s | ELS
[NARRATOR] The Great Pyramid of Khufu rose to a staggering height of 146 meters,
[NARRATOR] remaining the tallest structure on Earth for nearly four thousand years.
[CAM] Slow vertical tilt from base to apex
[BGM] Music reaches triumphant peak
[VFX] Golden sunlight catching the limestone cap, slight glow effect

// SHOT 16 | 4s | WS
[NARRATOR] But the pyramids were more than just tombs.
[NARRATOR] They were a statement — a declaration that Egypt would last forever.
[NOTE] Final shot of this section — sense of timelessness and permanence
```

---

## Quy tắc viết script

### 1. Mọi dòng PHẢI có tag
```
✅ [NARRATOR] The kingdom prospered for centuries.
✅ [Soldier] Move forward!
✅ [SFX] Explosion

❌ The kingdom prospered for centuries.     ← THIẾU TAG
```

### 2. Tên nhân vật nhất quán
```
✅ [Moses] ... (xuyên suốt script)
❌ Lúc [Moses] lúc [Mô-sê]
```

### 3. Một dòng = một người nói
```
✅ [Moses] Let my people go!
✅ [God] I will be with you.

❌ [Moses] Let my people go! [God] No!
```

### 4. Dialogue nên ngắn gọn

| Loại | Độ dài | Ví dụ |
|------|--------|-------|
| **Reaction** | 1-3 từ | `[Soldier] Incredible!` |
| **Soft line** | 4-8 từ | `[Farmer] The harvest looks good!` |
| **Quote** | 1-2 câu | `[Caesar] I came, I saw, I conquered.` |
| **Full dialogue** | Max ~10 giây | Nhiều dòng `[A]` `[B]` liên tiếp |
| **Crowd** | Bất kỳ | `[CROWD] Long live the king!` |

### 5. Mỗi shot = 1 hình ảnh cụ thể
Nếu cảnh thay đổi → tạo shot mới.

### 6. Yêu cầu tối thiểu
Cần **≥ 3 marker `// SHOT`** để hệ thống nhận diện chế độ Structured.

### 7. Markdown headers được phép
`#`, `##`, `---` đều bị bỏ qua khi parse — dùng thoải mái để tổ chức script.

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
| **Mở đầu** (20%) | Chủ yếu `[NARRATOR]` | Thiết lập bối cảnh |
| **Phát triển** (40%) | Cân bằng | Câu chuyện + dialogue tăng dần |
| **Cao trào** (30%) | Nhiều `[Character]` | Kịch tính, trao đổi nhanh |
| **Kết** (10%) | Quay lại `[NARRATOR]` | Tổng kết, suy ngẫm |

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
- Phân script thành từng shot bằng marker `// SHOT XX`
- Mỗi shot = 1 hình ảnh cụ thể (~4-6 giây)
- MỌI DÒNG NỘI DUNG PHẢI CÓ TAG — không có dòng nào thiếu tag

**CÚ PHÁP SHOT MARKER:**
```
// SHOT [số] | [duration]s | [shot_type] | [audio_mode]
```
- Chỉ `// SHOT [số]` là bắt buộc. Các tham số khác tùy chọn.
- shot_type: ELS, LS, WS, MWS, MS, MCU, CU, ECU, title-card, tactical-map
- audio_mode: narrator_only (mặc định) hoặc dialogue_dominant

**HỆ THỐNG TAG (bắt buộc cho MỌI dòng nội dung):**

Audio tags (xác định loại shot):
- `[NARRATOR]` = lời kể / voice-over → nhân vật KHÔNG nói, miệng đóng
- `[Tên Nhân Vật]` = dialogue → nhân vật NÓI, miệng mở
- `[CROWD]` = tiếng đám đông

Metadata tags (bổ sung thông tin):
- `[SFX]` = hiệu ứng âm thanh (ví dụ: `[SFX] Tiếng bom nổ rền`)
- `[BGM]` = nhạc nền (ví dụ: `[BGM] Tense military drums, building suspense`)
- `[CAM]` = chỉ đạo camera (ví dụ: `[CAM] Slow Ken Burns zoom in`)
- `[VFX]` = hiệu ứng hình ảnh (ví dụ: `[VFX] Sepia tone, film grain`)
- `[NOTE]` = ghi chú đạo diễn (ví dụ: `[NOTE] Tạo cảm giác ngột ngạt`)

**QUY TẮC QUAN TRỌNG:**
1. KHÔNG để dòng nào thiếu tag — dùng `[NARRATOR]` cho mọi lời kể
2. Tên nhân vật phải nhất quán trong suốt script
3. Dialogue nên ngắn: reaction 1-3 từ, quote 1-2 câu, max ~10 giây
4. Mỗi shot nên có 1-3 câu `[NARRATOR]` (~4-6 giây)
5. Cứ 3-4 shot narrator → chèn 1 shot dialogue để tạo nhịp
6. `[BGM]` chỉ cần đặt ở shot đầu tiên khi mood thay đổi
7. Dùng ## Headers để tổ chức script (sẽ bị bỏ qua khi parse)

**VÍ DỤ:**
```
// SHOT 01 | 3s | title-card
[BGM] Dark ambient drone, low rumble
[NOTE] Title card mở đầu

// SHOT 02 | 6s | ELS
[SFX] Tiếng gió sa mạc thổi qua
[NARRATOR] Giữa vùng sa mạc khô cằn, một đoàn quân đang hành quân dưới nắng gắt.
[VFX] Desaturated palette, heat haze effect

// SHOT 03 | 5s | MS
[CAM] Slow tracking shot following soldiers
[NARRATOR] Họ đã đi suốt 3 ngày đêm không nghỉ.

// SHOT 04 | 4s | MCU | dialogue_dominant
[Officer] Tiến lên! Không được dừng lại!
[SFX] Tiếng bước chân nặng nề trên cát

// SHOT 05 | 5s | WS
[NOTE] Tạo cảm giác rộng lớn và cô đơn
[NARRATOR] Nhưng phía trước, một thử thách lớn hơn đang chờ đợi.
[BGM] Music builds with ominous brass
```

Bây giờ hãy viết script cho: [CHỦ ĐỀ CỤ THỂ]
````

---

## Voice Profile (Tương lai)

Hệ thống hỗ trợ đặt voice profile để đảm bảo đồng nhất giọng xuyên suốt video:

| Cấp | Trường | Ví dụ |
|-----|--------|-------|
| **Drama** (project) | `narrator_voice_profile` | `"Male, 30s, British accent, deep baritone, calm documentary tone"` |
| **Character** | `voice_style` | `"Female, child-like, cheerful, high-pitched"` |

Voice profile được inject vào video prompt và sẽ được dùng cho TTS integration trong tương lai.

---

## Lưu ý kỹ thuật

1. **Ngôn ngữ:** Script viết bằng **bất kỳ ngôn ngữ nào**. Hệ thống auto output trường kỹ thuật bằng tiếng Anh.

2. **Auto-detect:** ≥ 3 marker `// SHOT` → kích hoạt Structured mode.

3. **Backward compat:** Dòng không tag vẫn hoạt động (tự động = narrator), nhưng **khuyến cáo KHÔNG dùng** — luôn dùng `[NARRATOR]`.

4. **Tag alias:** `[CAM]` = `[CAMERA]`, `[NOTE]` = `[DIR]`.

5. **Nhiều SFX:** `[SFX]` × nhiều → nối bằng `; `. Các tag metadata khác lấy giá trị cuối.

6. **Narrator nhiều dòng:**
   ```
   [NARRATOR] Take off your sandals, for you are standing on holy ground.
   [NARRATOR] I am the God of your father.
   ```
   → `narrator_script: "Take off your sandals... \n I am the God..."`

7. **Duration tự động:** ~16 từ tiếng Anh ≈ 7 giây (~140 WPM).

8. **Thứ tự tag:** Không quan trọng. `[SFX]` đầu, `[BGM]` cuối, xen kẽ — đều OK.

9. **Video prompt:** Narrator text được inject vào video prompt để AI tạo hình ảnh minh họa phù hợp với lời kể. Shot narrator → nhân vật đóng miệng. Shot dialogue → nhân vật mở miệng.
