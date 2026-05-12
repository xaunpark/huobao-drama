# HL Meow Meow — Title, Description & Thumbnail Generator

Tạo tiêu đề YouTube, mô tả video, và prompt tạo thumbnail theo đúng công thức kênh gốc.

> **Bối cảnh**: Đã có ảnh tham chiếu (reference image) của nhân vật chính là mèo trong trang phục trung tính. Thumbnail prompt chỉ cần mô tả thêm trang phục/ngoại hình theo video cụ thể + bố cục + biểu cảm.

---

## PHẦN 1: TITLE

### Cấu trúc bắt buộc

```
【Emotion Tag】CharacterName、[Setup/Problem]…[Cliffhanger/Result] + Emoji Cluster
```

### Thành phần

#### 1. Emotion Tag (mở đầu — LUÔN CÓ)

| Tag | Khi nào | Tần suất |
|---|---|---|
| 【感動】 | Kết quả cảm động, nỗ lực được đền đáp | Rất cao |
| 【まさかの展開】 | Plot twist, kết quả bất ngờ | Rất cao |
| 【衝撃】 | Sốc, ngoài mong đợi | Cao |
| 【まさか】 | Bất ngờ nhẹ hơn 衝撃 | Cao |
| 【号泣】 | Khóc, cảm động sâu sắc | Cao |
| 【大号泣】 | Cực kỳ cảm động | Trung bình |
| 【大逆転】 | Lật ngược tình thế | Trung bình |
| 【まさかの大成功】 | Thành công bất ngờ | Trung bình |
| 【できるかな？】 | Thử thách — liệu có làm được? | Thấp |
| 【まさかのご褒美】 | Phần thưởng bất ngờ | Thấp |
| 【まさかの挑戦】 | Thử thách bất ngờ | Thấp |
| 【まさかの高熱】 | Biến cố sức khỏe | Thấp |
| 【やらかした…】 | Gây chuyện, sai lầm | Thấp |

Biến thể: emoji trước tag (`💥【衝撃】`), hoặc emoji đơn thay tag (`😱`).

#### 2. Tên nhân vật + Setup + Cliffhanger

| Thành phần | Rule |
|---|---|
| **Tên** | Luôn có, theo sau bởi `、` |
| **Setup** | 1 mệnh đề mô tả vấn đề. Dùng con số cụ thể (`15万円`, `39.2℃`) |
| **`…`** | LUÔN CÓ giữa Setup và Cliffhanger |
| **Cliffhanger** | Gợi ý, KHÔNG spoil: `夢が叶った！？`, `結末がすごすぎた`, `に涙` |
| **Emoji** | 2-5 cuối title. 1 cảm xúc (😭🥹😱) + chủ đề (🐙🍩🚲) |

#### Chiều dài: 40-60 ký tự Nhật + emoji

### Title Examples

```
【感動】ここにゃん、初めてのたこ焼き作りで人生激変…屋台が大行列になった理由🥹🐙
【号泣】ここにゃん、お茶をこぼしてスマホが壊れた…15万円のiPhoneを買うために初バイト😭📱✨
【まさかの大成功】ここにゃん、ドーナツに夢中でお店を開店…初売上で夢のウォーターパークへ行けた物語🍩💰🌊
```

---

## PHẦN 2: DESCRIPTION

### Cấu trúc 5 phần

#### 1. Hook (2-3 dòng)

Nhắc lại premise + emoji. Shock statement → mô tả sự cố.

```
ここにゃんが大変なことに…！？😱🐝
はちみつパンを食べていたら、突然ハチの大群が家に侵入！？
顔を刺されてしまい、ぷっくり腫れてしまったここにゃん…。
```

#### 2. Story Recap (6-12 dòng)

Tóm tắt beats chính. Mỗi beat = 1-2 dòng, xen emoji, chi tiết cụ thể.

#### 3. Climax Teaser (2-3 dòng)

Gợi ý kết quả + "ぜひ最後まで見届けてください✨"

#### 4. CTA (2-3 dòng)

```
💖 チャンネル登録・高評価よろしくお願いします！
🔔 通知ONで最新動画もチェックしてね！
```

Hoặc dạng ✔ content tags (xem ví dụ đầy đủ ở cuối file).

#### 5. Hashtags

Core: `#ここにゃん #子猫 #猫動画 #感動 #kitten #cat #癒し #かわいい猫`
Thêm 3-5 tag chủ đề.

---

## PHẦN 3: THUMBNAIL PROMPT

### Visual DNA — Phân tích từ thumbnails gốc

Từ 3 thumbnail tham chiếu, rút ra các quy tắc sau:

#### Composition Formula

```
[NHÂN VẬT CHÍNH ~40-50% frame] + [KEY PROP / DESIRE OBJECT] + [BỐI CẢNH GỢI Ý]
```

- **Nhân vật**: Luôn chiếm 40-50% diện tích frame. KHÔNG BAO GIỜ nhỏ hay bị cắt quá nhiều.
- **Key Prop**: 1 đối tượng duy nhất kể câu chuyện (chiếc váy, nhiệt kế, máy may). Đặt cạnh nhân vật.
- **Bối cảnh**: Đủ để set context nhưng KHÔNG quá chi tiết. Bối cảnh phục vụ story, không phải decoration.

#### 3 Layout Patterns

| Pattern | Mô tả | Khi nào |
|---|---|---|
| **Subject + Desire** | Nhân vật BÊN TRÁI nhìn sang đối tượng mong muốn BÊN PHẢI | Mục tiêu/ước mơ: mua váy, muốn đồ ăn |
| **Subject + Situation** | Nhân vật Ở GIỮA, bối cảnh bao quanh kể tình huống | Đang làm việc: may vá, nấu ăn, giao hàng |
| **Subject + Crisis** | Close-up nhân vật, 1 element chỉ vấn đề (nhiệt kế, bandage) | Biến cố: ốm, bị thương, gặp nạn |

#### Expression Rules

| Emotion | Biểu cảm mèo | Khi nào |
|---|---|---|
| **Tò mò / Muốn** | Mắt mở to long lanh, đầu hơi nghiêng, miệng hé nhẹ | Nhìn đồ muốn mua, khám phá |
| **Hào hứng / Quyết tâm** | Miệng mở rộng (ha!), mắt sáng, tai dựng | Đang làm việc, thành công |
| **Buồn / Ốm** | Mắt nhắm hoặc mở nửa, tai cụp, miệng mím | Bệnh viện, thất bại, mệt |
| **Sốc / Hoảng** | Mắt tròn xoe, miệng O, paw giơ lên | Phát hiện hết tiền, tai nạn |
| **Hài lòng** | Mắt nhắm cười, miệng cong lên | Ăn ngon, nhận lương, hoàn thành |

#### Outfit Rules

> **Bối cảnh**: Đã có reference image nhân vật trong trang phục trung tính.
> Thumbnail prompt chỉ cần mô tả **lớp trang phục thêm** phù hợp video.

| Chủ đề | Outfit thêm |
|---|---|
| Nấu ăn | Miniature white chef apron, tiny chef hat |
| Bệnh viện (bệnh nhân) | Hospital patient gown (red/white pattern), cooling pad on forehead |
| Bệnh viện (bác sĩ) | White lab coat, tiny stethoscope around neck |
| May vá / Xưởng | Simple work apron |
| Giao hàng | Delivery cap, small delivery bag |
| Lễ hội | Miniature yukata/kimono |
| Bơi lội | Tiny swim goggles on forehead, arm floats |
| Trượt tuyết | Miniature winter jacket, ski goggles |
| Văn phòng | Tiny necktie, miniature suit jacket |
| Mua sắm (trung tính) | Default outfit — không cần thêm |

#### Style Specs (LUÔN ÁP DỤNG)

```
- Hyper-realistic commercial pet photography
- Clean bright lighting (5500K daylight hoặc bright indoor)
- Sharp rendering — NO motion blur
- Vivid, high saturation — must POP on YouTube feed
- 16:9 aspect ratio (1280×720 minimum)
- Background: contextual, slightly simplified (không quá chi tiết — focus vào nhân vật)
- Color: bright, high-contrast, commercially clean
- NO text in the generated image (text overlay thêm riêng)
```

#### Text Overlay Specs (THÊM RIÊNG, KHÔNG trong AI prompt)

Text được thêm riêng sau khi có ảnh, KHÔNG yêu cầu AI tạo text:

```
- Font: Bold Japanese sans-serif (Gothic/rounded)
- Color: White text + heavy black outline/shadow (readability tối đa)
- Vị trí: Phía trên hoặc dưới, KHÔNG che mặt nhân vật
- Nội dung: 1 câu ngắn gây tò mò (5-12 ký tự Nhật) + dấu ? hoặc !
- Watermark: "By [ChannelName]" bottom-right, nhỏ
```

**Text content** — chọn 1 trong các patterns:

| Pattern | Ví dụ | Khi nào |
|---|---|---|
| Câu hỏi gợi tò mò | `このドレス、入るかな？` | Mục tiêu/ước mơ |
| Statement + surprise | `雨に濡れただけなのに。！` | Biến cố bất ngờ |
| Action + surprise | `子猫が着物を縫っている。？` | Hoạt động lạ thường |
| Số cụ thể + câu hỏi | `39℃？` | Nhấn mạnh severity |
| Reaction | `まさか…！？` | Twist/sốc |

---

### Thumbnail Prompt Template

```
[COMPOSITION]: [layout pattern] — [vị trí nhân vật trong frame] + [key prop/object]
[SUBJECT]: Anthropomorphic ginger tabby kitten (Tabinyan), standing upright / [pose], [expression description]. [Outfit addition if any].
[KEY PROP]: [1 đối tượng chính kể câu chuyện — mô tả vật lý cụ thể]
[ENVIRONMENT]: [bối cảnh gợi ý — đủ context, không quá chi tiết]
[STYLE]: Hyper-realistic commercial pet photography, clean bright daylight, vivid high-saturation colors, sharp rendering, 16:9, no text, no logos.
```

### Ví dụ Prompt hoàn chỉnh

#### Ví dụ 1: "Muốn mua váy" (Subject + Desire layout)

```
Prompt: Hyper-realistic commercial pet photo, 16:9. An anthropomorphic ginger tabby kitten stands upright on the left side of the frame, looking through a clothing shop window display with wide curious eyes, mouth slightly open, head tilted. She wears a simple pink t-shirt (default neutral outfit). Inside the shop window on the right side, a small red polka-dot dress hangs on display, brightly lit. The background is a Japanese street storefront exterior with a glass window. Clean bright daylight, vivid saturated colors, sharp rendering, no text, no logos.

Text overlay (thêm riêng): "このドレス、入るかな？"
```

#### Ví dụ 2: "Bị ốm 39 độ" (Subject + Crisis layout)

```
Prompt: Hyper-realistic commercial pet photo, 16:9. Close-up shot of an anthropomorphic ginger tabby kitten lying in a hospital bed, wearing a red-and-white patterned hospital gown. A blue cooling pad sits on her forehead. Her eyes are half-closed with a tired, vulnerable expression. A realistic human hand enters from the right side of the frame holding a digital thermometer near her face. White hospital interior background with medical equipment softly visible. Clean bright indoor lighting, vivid colors, sharp rendering, no text, no logos.

Text overlay (thêm riêng): "雨に濡れただけなのに。！" + "39℃?"
```

#### Ví dụ 3: "May kimono" (Subject + Situation layout)

```
Prompt: Hyper-realistic commercial pet photo, 16:9. An anthropomorphic ginger tabby kitten sits at a wooden table operating a small pink-and-white sewing machine, both paws pressing fabric through the machine. Her mouth is wide open in an excited expression, ears perked up. Flowing pink floral kimono fabric extends from the sewing machine across the table. Thread spools in pink, red, and white are scattered on the table. A small green potted plant sits in the background. Bright warm indoor lighting, vivid saturated colors, sharp rendering, no text, no logos.

Text overlay (thêm riêng): "子猫が着物を縫っている。？"
```

---

## QUY TRÌNH TẠO BỘ 3 (Title + Description + Thumbnail)

### Input cần

Từ script/outline, xác định:
1. **GOAL**: Mục tiêu ban đầu
2. **TRIGGER**: Sự cố / vấn đề
3. **WORK**: Công việc / nỗ lực
4. **CLIMAX**: Kết quả bất ngờ nhất
5. **EMOTION**: Cảm xúc chính
6. **KEY VISUAL**: 1 hình ảnh đại diện nhất cho video (moment nào sẽ làm thumbnail?)

### Bước 1: Title (3-5 variants)

Chọn Emotion Tag → viết `【Tag】Name、[Setup]…[Cliffhanger] + Emoji`
Chọn variant clickbait nhất mà vẫn trung thực.

### Bước 2: Thumbnail Prompt

Xác định:
- **Layout pattern**: Subject+Desire / Subject+Situation / Subject+Crisis
- **Key moment**: Moment nào trong video gây tò mò nhất khi thấy thumbnail?
- **Expression**: Matching emotion tag (tò mò / hào hứng / buồn / sốc)
- **Outfit addition**: Cần thêm gì so với reference image?
- **Key prop**: 1 đối tượng kể câu chuyện

Viết prompt theo template. Thêm text overlay specs riêng.

### Bước 3: Description

Viết 5 phần: Hook → Story Recap → Climax Teaser → CTA → Hashtags.

### Output Format

```
## Title
【感動】たびにゃん、初めてのたこ焼き作りで人生激変…屋台が大行列になった理由🥹🐙

## Thumbnail Prompt
[full prompt — no text in image]

## Thumbnail Text Overlay
Line 1: "たこ焼き、大成功！？"
Watermark: "By たびにゃん"

## Description
[full 5-part description in Japanese]
```

---

## Localization

Khi dùng cho Tabinyan:
- Thay `ここにゃん` → `たびにゃん`
- Thay `#ここにゃん` → `#たびにゃん`
- Thay `By ここにゃん` → `By たびにゃん`
- Giữ nguyên toàn bộ cấu trúc, tone, và công thức

---

## Description Full Example

```
ここにゃんが大変なことに…！？😱🐝
はちみつパンを食べていたら、突然ハチの大群が家に侵入！？
顔を刺されてしまい、ぷっくり腫れてしまったここにゃん…。

でも泣きながらも、自分で服を着てバスに乗り、病院へ向かいます🥺🏥
受付、診察、お薬までも全部ひとりで頑張る姿に感動…。

最後には「もう二度とハチを呼ばないようにしよう！」と反省して、
しっかりお片付けするここにゃんでした✨🐱🍯

今回の動画は、
✔ 子猫のかわいい日常
✔ ハチ騒動＆病院ストーリー
✔ 日本の日常風景
✔ 感動ストーリー
✔ 癒し＆ほっこり動画
が好きな方におすすめです💛

ぜひ最後まで見て、ここにゃんを応援してください✨
コメント＆高評価もよろしくお願いします🐾

#ここにゃん #子猫 #猫動画 #感動 #ハチ #病院 #kitten #cat #癒し #かわいい猫 #日本の日常 #子猫成長記録 #猫アニメ #動物動画 #ほっこり
```
