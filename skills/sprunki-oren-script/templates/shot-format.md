# Shot Format — Output Template (Sprunki Oren — Adult Version)

Mỗi SHOT trong output script phải tuân theo format **Visual Unit Structured** để tương thích với Go parser (`parseStructuredShots()`).

> [!IMPORTANT]
> Format này dùng prefix `// SHOT` để Go parser tự động trích xuất tags (SFX, dialogue, CAM, BGM...) trước khi gửi AI.
> Tham chiếu đầy đủ: `docs/features/preserve-structure-input-format.md`

---

## Format mỗi SHOT

```
// SHOT {XX} | {duration}s | {shot_type} | {audio_mode}
[BGM] description
[SFX] description
[Character Name] dialogue text
[CAM] camera direction
[NOTE] visual description — camera setup + subject + action + environment
```

### Quy tắc header

```
// SHOT {number} | {duration}s | {shot_type} | {audio_mode}
```

| Field | Bắt buộc | Giá trị |
|---|---|---|
| `// SHOT {number}` | ✅ | `// SHOT 01`, `// SHOT 02`... Prefix `//` BẮT BUỘC |
| `{duration}s` | Nên có | `2s`, `3s`, `4s`, `5s` |
| `{shot_type}` | Nên có | `WS`, `MWS`, `MS`, `MCU`, `CU`, `ECU`, `INSERT` |
| `{audio_mode}` | Tùy chọn | `narrator_only` hoặc `dialogue_dominant` (tự infer nếu không ghi) |

### Quy tắc tags

| Tag | Loại | Ảnh hưởng audio_mode | Ghi chú |
|---|---|---|---|
| `[BGM]` text | Metadata | Không | Chỉ khi mood thay đổi hoặc shot đầu |
| `[SFX]` text | Metadata | Không | **BẮT BUỘC mọi shot** — American adult ambient SFX |
| `[Character]` text | Audio | → `dialogue_dominant` | Tag = TÊN NHÂN VẬT, không thêm gì |
| `[NARRATOR]` text | Audio | → `narrator_only` | Hoặc dòng không tag = narrator |
| `[CAM]` text | Metadata | Không | Handheld documentary style |
| `[NOTE]` text | Metadata | Không | Visual description → `visual_description` |
| `[VFX]` text | Metadata | Không | Append vào `visual_description` |

> [!WARNING]
> **CRITICAL — Tag dialogue phải ĐÚNG tên nhân vật:**
> - ✅ `[Oren] Aight, bet.`
> - ✅ `[Simon] BRO. I have an idea.`
> - ✅ `[Pinki] Did you pay rent?`
> - ✅ `[Manager] Nice work, man.`
> - ❌ ~~`[DIA — Oren, chill] "Bruh."`~~ → Parser đọc tag = `DIA — Oren, chill` → SAI
> - ❌ ~~`[DIA — Simon, excited] "Dude!"`~~ → Parser không match character

> [!WARNING]
> **CRITICAL — Visual description dùng `[NOTE]`, KHÔNG dùng `Visual:`:**
> - ✅ `[NOTE] Wide shot of studio apartment. Coffee maker on counter, records on wall.`
> - ❌ ~~`Visual: Wide shot of studio apartment...`~~ → Parser coi là narrator line

### Dòng không có tag

Dòng không match pattern `[...]` → tự động thành narrator:

```
// SHOT 05 | 4s | WS | narrator_only
So they waited. And waited. And... rent was still due.
```

→ `narrator_text = "So they waited. And waited. And... rent was still due."`

---

## Quy tắc viết `[NOTE]` (Visual Description)

**Công thức:**
```
[NOTE] {Shot type + camera position}. {Subject + vị trí + action}. {Environment context}.
```

**Luật:**
1. Camera/shot type mở đầu
2. Subject rõ ràng (ai, ở đâu trong frame)
3. 1 hành động chính duy nhất
4. Context: **ADULT American environment** cụ thể
5. Oren body language: lean on things, swagger, tay trong túi, headphones, coffee

**Adult Context Cues trong [NOTE]** — luôn include ít nhất 1:
- Adult-scale furniture (bar stool, office chair, couch)
- Adult beverages (coffee, beer)
- Adult tech (laptop, phone with apps)
- Adult environment (apartment, bar, office, coffee shop)
- Adult posture (lean, slouch, confident stride)

**CẤM:**
- ❌ Style descriptions (lighting, color, render style)
- ❌ Cảm xúc trừu tượng ("tense", "warm")
- ❌ Multi-action trong 1 shot
- ❌ Ngoại hình chi tiết (da cam, antenna shape) — dùng reference image
- ❌ Childish environments (playground, toy store, school)

---

## Shot Distribution

| Shot Type | Tỷ lệ | Khi nào |
|---|---|---|
| Medium Shot (MS) | ~40% | Oren thực hiện hành động |
| Close-Up (CU) | ~20% | Biểu cảm, detail food/tech/phone |
| Wide Shot (WS) | ~15% | Establishing (apartment, street, bar) |
| Medium Wide (MWS) | ~15% | Two-shot Oren + Human/Friend |
| Insert/Detail | ~10% | Phone screen ($3.47), coffee cup, laptop screen |

## Camera Angle

| Angle | Tỷ lệ |
|---|---|
| Oren eye-level | 50% |
| Standard eye-level | 25% |
| High angle (overview) | 10% |
| Low angle (power/confidence) | 10% |
| Dutch angle (hài hước/chaos) | 5% |

## Pacing

| Loại | Duration |
|---|---|
| Action montage | 2-3s |
| Reaction | 2-3s |
| Establishing | 3-5s |
| Transition | 80% hard cut, 20% whip-pan |

---

## Ví dụ hoàn chỉnh

### Episode: "Oren's Rent Is Due"

```
// SHOT 01 | 4s | WS
[BGM] Chill lo-fi beat, lazy afternoon vibe
[SFX] City traffic hum; distant skateboard wheels; apartment AC unit
[CAM] Handheld follow, approaching from street level
[NOTE] Wide shot at eye-level approaching a brick apartment building with fire escapes. A skateboard leans against the stoop. Mailbox stuffed with envelopes. Urban neighborhood with parked cars and a bodega across the street.

// SHOT 02 | 3s | MS
[SFX] Coffee maker gurgling; phone notification buzz
[CAM] Static, eye-level, slight sway
[NOTE] Medium shot inside a small studio apartment. Oren stands at kitchen counter, one hand around a coffee mug, other hand holding phone, scrolling. Messy counter with takeout containers. Records on the wall behind. MIDI keyboard on desk in background.

// SHOT 03 | 3s | INSERT | dialogue_dominant
[SFX] Phone screen tap; banking app loading sound
[Oren] Bro, what?
[CAM] Static close-up, slight push-in
[NOTE] Insert shot of Oren's hand holding smartphone. Screen displays banking app showing "$3.47" balance. Uber Eats transaction history visible. Other hand sets coffee mug on counter.

// SHOT 04 | 2s | CU
[SFX] Wallet flipping open — cards, no cash; sigh
[CAM] Static, eye-level
[NOTE] Close-up of Oren center frame. Eyes half-lidded, jaw clenches slightly. Opens wallet — credit cards but no cash. Rubs the back of his neck. Antenna droop slightly.

// SHOT 05 | 3s | WS
[SFX] Door closing; footsteps on sidewalk; city ambient
[BGM] Beat picks up — upbeat lo-fi
[CAM] Handheld follow from behind, walking pace
[NOTE] Wide shot from behind. Oren walks down urban sidewalk, hands in jacket pockets, headphones on, slight swagger. Passes storefronts. Stops in front of a pub with "BARTENDER NEEDED" sign in window.

// SHOT 06 | 3s | MWS | dialogue_dominant
[SFX] Bar door opening; glasses clinking; ambient bar chatter
[Manager] Hey! You here about the bartending gig?
[Oren] Yo. Yeah, that's me.
[CAM] Static, slightly wider
[NOTE] Medium wide two-shot inside pub. Stocky man in black vest leans on bar counter, wiping a glass, looking at Oren with casual nod. Oren stands on other side, one hand on counter, other in pocket. Beer taps and bottles line the shelf behind.

// SHOT 07 | 3s | MS
[SFX] Cocktail shaker rattling; ice cubes clinking; liquid pouring
[CAM] Static, eye-level, slight push-in
[NOTE] Medium shot from side. Oren behind the bar, sleeves rolled up, shaking a cocktail shaker with focused expression. Tongue slightly out in concentration. Bottles and glasses arranged on counter.

// SHOT 08 | 2s | INSERT
[SFX] Liquid pouring; glass set down; garnish drop
[CAM] Overhead, static
[NOTE] Top-down insert. Two small orange hands pour cocktail from shaker into glass, then drop a lime wedge in. Cocktail tools and napkins frame the shot.

// SHOT 09 | 3s | CU | dialogue_dominant
[SFX] Cash register beep; customer applause
[Oren] That's fire.
[CAM] Static, eye-level
[NOTE] Close-up of Oren doing a contained fist-pump, small satisfied smirk. Bar counter in foreground. Customer in background giving thumbs-up with their drink.

// SHOT 10 | 3s | MWS | dialogue_dominant
[SFX] Venmo notification sound; phone buzz
[Manager] Good work tonight, man. Same time tomorrow?
[CAM] Static, wider
[NOTE] Medium wide shot. Manager extends hand for a handshake toward Oren, who takes it firmly. Manager's other hand holds phone showing Venmo. Oren glances at phone, eyes light up subtly.

// SHOT 11 | 3s | MS | dialogue_dominant
[SFX] Phone tapping; payment confirmation sound
[Oren] Aight. Rent's handled.
[CAM] Handheld, slight sway
[NOTE] Medium shot. Oren sits on apartment couch, laptop open on coffee table, phone in hand showing rent payment confirmation. Leans back, puts one foot on table, small exhale of relief. Headphones around neck.

// SHOT 12 | 4s | MS | dialogue_dominant
[SFX] Beer can opening; couch settling; lo-fi playlist starting
[Oren] ...same time next week.
[BGM] Lo-fi beat fades back in, mellow
[CAM] Static, eye-level, slow zoom out
[NOTE] Medium shot. Oren sunk into couch, beer in one hand, phone in other scrolling playlist. Headphones on. City lights through apartment window behind. Small satisfied smile, eyes half-closed. Takeout pizza box on coffee table next to laptop.
```
