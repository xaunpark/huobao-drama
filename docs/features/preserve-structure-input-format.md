# Preserve Structure — Proper Input Format

> Tài liệu này mô tả cách viết script đầu vào đúng chuẩn cho chế độ **Preserve Structure** (Split Mode) trong Storyboard Generator.

---

## Tổng quan

Chế độ **Preserve Structure** giữ nguyên số lượng và ranh giới shot đã được xác định sẵn trong script. Hệ thống **KHÔNG tự chia** — nó chỉ **thêm metadata** (shot_type, angle, movement, emotion, dialogue, bgm...) vào từng shot đã có.

Có **2 luồng xử lý** khác nhau tùy vào format shot header trong script:

| Luồng | Trigger | Prompt sử dụng | Parser |
|---|---|---|---|
| **Preserve (Generic)** | `SHOT 01`, `Scene 1`, timestamps `(0:00 – 0:06)` | `storyboard_preserve_shots.txt` | AI tự parse (không có Go parser) |
| **Visual Unit Structured** | `// SHOT 01` (có prefix `//`) | `storyboard_visual_unit_structured.txt` | Go parser `parseStructuredShots()` → trích xuất tags trước khi gửi AI |

> [!IMPORTANT]
> **Visual Unit Structured** là luồng chính xác nhất vì Go parser trích xuất dialogue, SFX, BGM, CAM trước khi AI xử lý. Luồng **Preserve Generic** gửi raw text cho AI → AI có thể bỏ sót hoặc map sai dialogue.

---

## Khuyến nghị: Dùng Visual Unit Structured

### Shot Header Format

```
// SHOT {number}
// SHOT {number} | {duration}s
// SHOT {number} | {duration}s | {shot_type}
// SHOT {number} | {duration}s | {shot_type} | {audio_mode}
```

**Regex pattern** (source: `storyboard_service.go` line 484-489):
```regex
^//\s*SHOT\s+(\d+)(?:\s*\|\s*(\d+)s)?(?:\s*\|\s*([\w-]+))?(?:\s*\|\s*(narrator_only|dialogue_dominant))?
```

**Ví dụ hợp lệ:**
```
// SHOT 01
// SHOT 02 | 3s
// SHOT 03 | 4s | MS
// SHOT 04 | 3s | CU | dialogue_dominant
// SHOT 05 | 5s | WS | narrator_only
```

**Detection**: Cần >= 3 markers `// SHOT` để kích hoạt structured mode (`detectStructuredShots()`).

### Tag System

Giữa 2 shot headers, mỗi dòng được parse bởi regex `^\s*\[([^\]]+)\]\s*(.*)$`:

| Tag | Loại | Ảnh hưởng audio_mode? | Map vào field |
|---|---|---|---|
| `[NARRATOR]` text | Audio | Có → `narrator_only` | `narrator_text` |
| `[Character Name]` text | Audio | Có → `dialogue_dominant` | `dialogue_text` |
| `[CROWD]` text | Audio | Có → `dialogue_dominant` | `dialogue_type="crowd"` |
| `[SFX]` text | Metadata | Không | `sound_effect` |
| `[BGM]` text | Metadata | Không | `bgm_prompt` |
| `[CAM]` text | Metadata | Không | `movement` (camera direction) |
| `[VFX]` text | Metadata | Không | Append vào `visual_description` |
| `[NOTE]` / `[DIR]` text | Metadata | Không | `reason_for_shot` / `visual_description` |
| (không có tag) | Audio | Có → `narrator_only` | `narrator_text` |

> [!WARNING]
> **Tag `[Character Name]` phải ĐÚNG là tên nhân vật**, không thêm annotation.
> - ✅ `[Simon] Yo, Friday is bowling night!`
> - ❌ `[DIA — Simon, excited] "Yo, Friday is bowling night!"`
>
> Format `[DIA — Simon, excited]` sẽ bị parse thành character name = `DIA — Simon, excited` → **sai**.

### Dòng không có tag

Bất kỳ dòng nào **không match** pattern `[...]` sẽ được coi là **narrator** line:

```
// SHOT 05 | 4s | WS
So they waited. And waited. And... yeah.
```

→ `narrator_text = "So they waited. And waited. And... yeah."`, `audio_mode = "narrator_only"`

### Visual Description

Dòng bắt đầu bằng `Visual:` **không** phải là tag — nó sẽ được coi là narrator line bởi parser. Nên mô tả visual bằng tag `[NOTE]` hoặc để AI tự infer từ context:

```
// SHOT 01 | 4s | WS
[SFX] Distant basketball bouncing; cicadas buzzing
[CAM] Handheld approach from sidewalk, slow walking pace
[NOTE] Wide shot approaching suburban house, basketball hoop over garage, five pairs of sneakers on porch
```

---

## Ví dụ hoàn chỉnh — Đúng format

```
// SHOT 01 | 4s | WS
[BGM] Chill lo-fi beat, lazy weekend afternoon vibe
[SFX] Distant basketball bouncing on pavement; cicadas buzzing; sprinkler sounds
[CAM] Handheld approach from sidewalk, slow walking pace
[NOTE] Wide shot approaching typical American suburban house. Basketball hoop over garage, garden hose on lawn, five pairs of colored sneakers on porch steps.

// SHOT 02 | 3s | MS | dialogue_dominant
[SFX] Phone screen tapping; group chatter ambience
[Simon] Yo, Friday is bowling night at Galaxy Lanes! Who's in?
[CAM] Static, eye-level with characters on bedroom floor
[NOTE] Inside Oren's bedroom. Five Sprunki characters sit in circle on floor. Simon holds up phone showing bowling flyer.

// SHOT 03 | 3s | CU | dialogue_dominant
[SFX] Wallet snapping open — empty leather sound
[Oren] Bruh.
[CAM] Static close-up, slight push-in
[NOTE] Close-up of Oren's hands pulling wallet open. Completely empty. Flips upside down — nothing falls out.

// SHOT 04 | 3s | INSERT
[SFX] Four wallets opening in rapid succession; collective sigh
[Durple] ...this is tragic.
[CAM] Overhead, static top-down
[NOTE] Top-down shot of five wallets spread open in circle on carpet. All empty. One crumpled $1 bill in center.

// SHOT 05 | 4s | WS | narrator_only
[SFX] Wind blowing; distant dog barking; Oren yawning
[CAM] Static wide, time-lapse feel
So they waited. And waited. And... yeah.
[NOTE] Wide shot of garage sale from across street. Table fully set. Five characters waiting. Sidewalk empty — zero customers.
```

---

## So sánh: Script SAI vs ĐÚNG

### ❌ SAI — Format hiện tại trong `ep03-garage-sale-flip.md`

```
SHOT 02 | 3s | MS
[SFX] Phone screen tapping — "tap tap tap"; group chatter ambience
[DIA — Simon, excited] "Yo, Friday is bowling night at Galaxy Lanes! Who's in?"
[CAM] Static, eye-level with characters sitting on bedroom floor

Visual: Medium shot inside Oren's bedroom. All five Sprunki characters sit in a loose circle on the floor.
```

**Vấn đề:**
1. `SHOT 02` — thiếu prefix `//` → Go parser không nhận dạng shot boundary
2. `[DIA — Simon, excited]` — parser đọc tag name = `DIA — Simon, excited`, không phải `Simon`
3. `"Yo, Friday..."` — text trong quotes bị parser cắt quotes nhưng tag name sai nên dialogue bị mất
4. `Visual: ...` — parser coi là narrator line (không phải tag), nội dung visual bị nhồi vào narrator_text

### ✅ ĐÚNG — Format chuẩn

```
// SHOT 02 | 3s | MS | dialogue_dominant
[SFX] Phone screen tapping; group chatter ambience
[Simon] Yo, Friday is bowling night at Galaxy Lanes! Who's in?
[CAM] Static, eye-level with characters on bedroom floor
[NOTE] Inside Oren's bedroom. Five Sprunki characters sit in circle on floor. Simon holds up phone showing bowling flyer.
```

**Sửa lại:**
1. `// SHOT 02` — có prefix `//` → parser match
2. `[Simon]` — tag name = `Simon` → match character ID
3. Text không có quotes → clean dialogue text
4. `[NOTE]` — visual description đúng tag, không bị nhồi vào narrator

---

## Tham chiếu nhanh: Tất cả Audio Mode

| audio_mode | Khi nào | Trigger |
|---|---|---|
| `narrator_only` | Voice-over, không có nhân vật nói | `[NARRATOR]` hoặc dòng không tag |
| `dialogue_dominant` | Nhân vật nói trực tiếp | `[Character Name]` hoặc `[CROWD]` |

Có thể force audio_mode trong header:
```
// SHOT 04 | 3s | CU | narrator_only
// SHOT 05 | 3s | MS | dialogue_dominant
```

Nếu không chỉ định, AI tự infer từ tags trong shot.

---

## Tham chiếu nhanh: Shot Type codes

Dùng trong header field thứ 3:

| Code | Tên đầy đủ |
|---|---|
| `ELS` | Extreme Long Shot |
| `LS` | Long Shot |
| `WS` | Wide Shot |
| `MWS` | Medium Wide Shot |
| `MS` | Medium Shot |
| `MCU` | Medium Close-Up |
| `CU` | Close-Up |
| `ECU` | Extreme Close-Up |
| `INSERT` | Insert/Detail Shot |

---

## Source code references

| File | Line | Mục đích |
|---|---|---|
| [storyboard_service.go](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L483-L489) | 483-489 | `shotHeaderPattern` regex — parse shot headers |
| [storyboard_service.go](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L491-L496) | 491-496 | `detectStructuredShots()` — detection (>= 3 markers) |
| [storyboard_service.go](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L498-L604) | 498-604 | `parseStructuredShots()` — Go parser logic |
| [storyboard_service.go](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L606-L699) | 606-699 | `buildStructuredAnalysis()` — build AI context |
| [storyboard_service.go](file:///g:/VS-Project/huobao-drama/application/services/storyboard_service.go#L911-L942) | 911-942 | `processVisualUnitGeneration()` — routing logic |
| [storyboard_preserve_shots.txt](file:///g:/VS-Project/huobao-drama/application/prompts/storyboard_preserve_shots.txt) | — | Preserve generic prompt |
| [storyboard_visual_unit_structured.txt](file:///g:/VS-Project/huobao-drama/application/prompts/storyboard_visual_unit_structured.txt) | — | Visual Unit Structured prompt |
