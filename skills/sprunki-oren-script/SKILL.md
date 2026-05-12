---
name: sprunki-oren-script
description: Generate detailed cinematic scripts for Sprunki Oren adult anthropomorphic slice-of-life episodes set in America
---

# Sprunki Oren — Cinematic Script Generator (Adult Version)

Tạo **kịch bản phân cảnh chi tiết** (shot-by-shot cinematic script) cho video nhân vật Oren (Sprunki) theo phong cách slice-of-life Mỹ.
Mọi nhân vật Sprunki là **NGƯỜI TRƯỞNG THÀNH** (young adult, 20s-30s) — KHÔNG phải trẻ con/teenager.
Mỗi SHOT trong output có thể dùng trực tiếp làm prompt cho AI image/video generation.

## Khi nào dùng

- User yêu cầu viết kịch bản cho một episode Sprunki Oren mới
- User cung cấp video gốc cần phân tích thành script chi tiết
- User cần chuyển đổi concept/outline thành storyboard shots

## Instrumentation

```bash
./scripts/log-skill.sh "sprunki-oren-script" "manual" "$$"
```

## What do you want to do?

1. **Viết script từ ý tưởng/outline** → Read `workflows/write-from-idea.md`
2. **Phân tích video gốc thành script** → Read `workflows/analyze-video.md`
3. **Xem character & world reference** → Continue reading this file (Section: World Bible)
4. **Xem output format** → Read `templates/shot-format.md`
5. **Tạo Title + Description + Thumbnail Prompt** → Read `references/title-description.md`
6. **Xem Voice Profile guide** → Read `references/voice-profiles.md`

---

## Nguyên tắc cốt lõi (LUÔN NẠP)

### 1. Thế giới (World-building)

- **Bối cảnh**: Nước Mỹ lý tưởng hóa — thành phố, suburbs, downtown. Sprunki characters sống cùng con người bình thường.
- **Quy tắc chấp nhận mặc nhiên**: KHÔNG AI đặt câu hỏi "Tại sao có sinh vật cam ở đây?". Oren có apartment, driver's license, job, bank account. Đây là bình thường.
- **Mature Casual American**: Giao tiếp casual, bình đẳng, giữa những người trưởng thành. "Hey man!", fist-bump, handshake, "What's up bro". KHÔNG patronizing, KHÔNG baby talk.
- **Eye Level**: Sprunki characters CÓ THỂ nhỏ hơn con người (khoảng 4-5 feet) nhưng được đối xử hoàn toàn bình đẳng — KHÔNG cúi xuống nói chuyện kiểu với trẻ con.
- **Ngôn ngữ**: Signage, thoại bằng tiếng Anh. Oren dùng tiếng Anh casual + một ít slang.
- **Vibe**: Adult indie comedy meets Cartoon Network Adults Swim. Chill, witty, relatable.

### 2. Nhân vật — ADULT REFRAMING (Character Roster)

> **CRITICAL**: Tất cả Sprunki characters là **YOUNG ADULTS (20s-30s)**. Họ có apartment, jobs, relationships, bills. KHÔNG phải kids/teens đi học.

#### Nhân vật chính

| Nhân vật | Vai trò | Ngoại hình | Tính cách | Tuổi tương đương |
|---|---|---|---|---|
| **Oren** | Nhân vật chính — Orange Sprunki | Sinh vật cam tươi, 2 antenna, tuft of hair trên trán, tai headphone cam nhạt, jacket xanh có chữ "O" | Laid-back, chill. Freelance beatmaker, hay quên trả bills. Frugal, impulsive nhưng chân thành. Sống ở studio apartment nhỏ | 25 |

#### Nhân vật phụ (Sprunki Friends)

| Nhân vật | Vai trò | Ngoại hình | Tính cách | Tuổi |
|---|---|---|---|---|
| **Pinki** | Bạn gái (GF) | Hồng, dễ thương, energetic | Warm, yêu thương. Barista tại coffee shop. Mạnh mẽ nhưng lo lắng về tương lai. Hay nói Oren tiêu tiền vô tội vạ | 24 |
| **Simon** | Best friend | Màu xanh lá, tự tin, outgoing | Energetic, social butterfly. Làm marketing freelance. Luôn rủ Oren đi adventure hoặc networking event. "Bro, this is gonna be huge!" | 26 |
| **Durple** | Bạn chí cốt | Tím, calm, chill | Mellow, wise. Chơi trumpet tại jazz bar. Cho lời khuyên tình cảm và life advice. Sống chậm, uống wine | 28 |
| **Gray** | Bạn trầm tính | Xám, introvert, gentle | Calm, introspective. Graphic designer, làm remote. Ít nói nhưng phân tích đúng trọng tâm. Thích sách và podcast | 27 |
| **Sky** | Bạn trẻ nhất | Nhạt/sky blue, serene | Thoughtful, mediator. Junior developer mới ra trường. Mang sự cân bằng khi nhóm bất đồng. Perspective tươi mới | 22 |

#### Con người (American Peers — Bình đẳng)

| Nhân vật | Vai trò | Đặc điểm |
|---|---|---|
| **Người Mỹ** | Đồng nghiệp, hàng xóm, bartender, barista, coworker... | Đa dạng chủng tộc. Casual, bình đẳng. "Hey man", "What's good?", handshake, fist-bump |
| **Jake** (archetype) | Hàng xóm / buddy | Người Mỹ 30s, easygoing. Rủ Oren xem game, cho mượn đồ. "You good, bro?" |
| **Manager / Boss** (archetype) | Người quản lý / client | Professional nhưng casual. "Nice work, man!", direct deposit hoặc Venmo |

### 3. Voice Profiles — ADULT VERSION

> **Voice profiles này PHẢI được nhập vào field `voice_style` của Character trong Drama.**
> Pipeline sẽ tự động đưa voice_style vào `character_voices` trong `shotContext` → LLM distill sử dụng khi viết video prompt.

Xem chi tiết tại: `references/voice-profiles.md`

#### Oren Voice Profile

**DB `voice_style`**: `Male, mid-20s, laid-back American. Medium-low pitch, slightly raspy, casual cadence. West Coast chill vibe. Speaks in short phrases with occasional slang. Beat-boxes and hums between sentences.`

| Câu nói | Khi nào | Tone/Delivery |
|---|---|---|
| "Yo." | Chào hỏi casual | Low, chill, one-word greeting |
| "Dude." | Ngạc nhiên | Mid-pitch, eyebrows up |
| "That's fire." | Ấn tượng | Genuine, slightly awed |
| "Nah, I'm good." | Từ chối nhẹ | Relaxed, not dismissive |
| "Bro, what?" | Bất ngờ tiêu cực | Flat, slightly confused |
| "Let's get it." | Bắt đầu thử thách | Rising energy, motivated |
| "Aw man..." | Thất vọng | Low, drawn out |
| "Pizza time, baby." | Food excitement | Upbeat, almost singing |
| *[beat-box/finger tap]* | Idle/thinking/happy | Rhythmic tapping on surfaces, humming |
| "Aight, bet." | Đồng ý | Quick, decisive |
| "Rent's due..." | Nhìn ví | Dread, trailing off |

#### Pinki Voice Profile

**DB `voice_style`**: `Female, mid-20s, warm and bright. Medium-high pitch, clear articulation, playful but with moments of genuine concern. Coffee shop energy — friendly and engaging.`

| Câu nói | Khi nào | Tone |
|---|---|---|
| "Babe, did you pay rent?" | Kiểm tra Oren | Concerned but loving |
| "Ugh, you're the worst. But also the best." | Trêu Oren | Playful sarcasm |
| "Okay but like... how?" | Nghe plan Oren | Skeptical, eyebrow raised |
| "I'm so proud of you!" | Oren thành công | Genuine, bright |
| "Come here." | Ôm/comfort | Soft, warm |

#### Simon Voice Profile

**DB `voice_style`**: `Male, mid-20s, high energy, extroverted. Medium-high pitch, fast-talking, enthusiastic. Always sounds like he just had 3 espressos. Marketing bro energy.`

| Câu nói | Khi nào | Tone |
|---|---|---|
| "BRO. Bro bro bro." | Có ý tưởng | Escalating excitement |
| "Trust me on this one." | Pitch kế hoạch | Confident, salesman |
| "Dude, content. This is CONTENT." | Thấy gì thú vị | Hyped, pointing at phone |
| "We're literally gonna blow up." | Lạc quan | Over-the-top optimistic |
| "...okay that didn't work." | Fail | Quick pivot, unbothered |

#### Durple Voice Profile

**DB `voice_style`**: `Male, late 20s, smooth and unhurried. Low-medium pitch, deliberate pacing, jazz musician cool. Thoughtful pauses between phrases. Barry White meets chill philosopher.`

| Câu nói | Khi nào | Tone |
|---|---|---|
| "Hmm. Let me think about that." | Xem xét | Slow, deliberate |
| "Love isn't about grand gestures, man." | Tình cảm advice | Smooth, wise |
| "That's beautiful." | Xúc động | Low, genuine |
| "I'll bring the wine." | Kế hoạch social | Casual, refined |
| "*[trumpet riff]*" | Mood music | Playing trumpet casually |

#### Gray Voice Profile

**DB `voice_style`**: `Male, late 20s, quiet and measured. Low pitch, soft volume, minimal words but each one counts. Introvert energy — comfortable with silence. Podcast narrator vibes.`

| Câu nói | Khi nào | Tone |
|---|---|---|
| "..." | Lắng nghe | *Nodding, present* |
| "You already know the answer." | Wisdom drop | Calm, direct |
| "I'll pass." | Từ chối social event | Gentle, no drama |
| "That's... a lot." | Overwhelmed | Understated |
| "Here. I made you this." | Tặng quà | Quiet, thoughtful |

#### Sky Voice Profile

**DB `voice_style`**: `Male, early 20s, fresh and optimistic. Medium pitch, genuine curiosity in voice. New grad energy — eager but not annoying. Asks good questions.`

| Câu nói | Khi nào | Tone |
|---|---|---|
| "Wait, can I ask something?" | Muốn hiểu thêm | Genuine curiosity |
| "Both sides have a point though." | Mediating | Balanced, diplomatic |
| "I googled it and actually..." | Research | Helpful, slightly nerdy |
| "Is it just me or is this kinda cool?" | Thấy điều hay | Wonder, fresh perspective |

### 4. 6 Human Interaction Patterns (Adult Peer Version)

| Pattern | Trigger | Hành động |
|---|---|---|
| **casual_coworker** | Cùng làm việc | Bình đẳng, "Hey, can you grab that?", teamwork, professional-casual |
| **impressed_client** | Nhận sản phẩm/dịch vụ | "Dude, this is exactly what I wanted!", tip/payment, IG story |
| **chill_neighbor** | Gặp hành lang / sân | "What's up man", cho mượn đồ, BBQ invite, dog walking chat |
| **concerned_friend** | Oren stress/broke | Mang beer/food, ngồi cạnh, "Talk to me, what's going on?" |
| **passerby_amused** | Nơi công cộng | Gật đầu respect, "Nice!" hoặc quick thumbs-up, không patronizing |
| **supportive_bro** | Oren thành công hoặc fail hài | Shoulder pat, laugh cùng, "You're wild, bro!" |

### 5. Storytelling DNA — 3-Beat Arc (Adult Version)

1. **Setup (Chill Vibe)**: Oren đang chill ở apartment — beat-boxing, lướt phone, cook ramen. Phát hiện vấn đề (rent due, empty fridge, event coming up). Energy: relaxed → "oh crap."
2. **Grind (Hustle)**: Oren tìm gig/freelance/side job. Learning curve, mistakes, nhưng quyết tâm. Close-up hands working. Beat-box nhẹ khi tập trung. Energy: determined.
3. **Payoff (Satisfaction)**: Get paid, task complete. Fist-bump, "That's fire.", chia sẻ thành quả với bạn bè. Energy: triumphant → chill lại.

### 6. Story Beat Template — "ADULT GRIND → EARN → ENJOY"

> **Quy tắc vàng**: Oren có CUỘC SỐNG NGƯỜI LỚN. Motivation = **bills, rent, date night, concert tickets, birthday gift, new gear**, KHÔNG phải "mua đồ chơi" kiểu trẻ con.

#### A. Công thức SINGLE EPISODE (1-2 phút, 10-15 shots)

```
WANT → BROKE! → HUSTLE → [FAIL/COMPLICATION] → GET PAID → ENJOY!
```

1. **WANT (Mục tiêu)**: Oren muốn gì?
   - Trả rent cuối tháng
   - Date night xịn cho Pinki
   - Mua mixing equipment mới
   - Concert tickets cho cả nhóm
   - Weekend getaway
   - Fix laptop bị hỏng
2. **BROKE! (Trigger)**: Check banking app → $3.47. "Bro, what?" → swipe through expenses → toàn Uber Eats. **Signature moment bắt buộc.**
3. **HUSTLE (Làm việc)**: Oren tìm gig — freelance beat production, food delivery, bartending, pet-sitting, moving help, graphic design side-gig. Coworker/boss dạy việc.
4. **FAIL** *(optional)*: Sai sót hài hước nhưng relatable (overcook đồ, drop package, wrong order). Concerned friend.
5. **GET PAID**: Direct deposit/Venmo/cash. Boss: "Good work, man.", professional handshake hoặc fist-bump.
6. **ENJOY! (Thỏa mãn)**: Dùng tiền thực hiện mục tiêu. Climax cảm xúc — Pinki ở rooftop dinner, nhóm bạn ở concert, new equipment unboxing.
   - Kết thúc: nhóm bạn chill cùng nhau trên couch/rooftop, beer/pizza, "Same time next week?", sunset.

#### B. Công thức COMPILATION (3-10 phút, 25-60 shots)

```
[Mini-arc 1: WANT→HUSTLE→ENJOY] → [Transition] → [Mini-arc 2] → ...
```

**3 loại kết nối:**

| Loại | Ví dụ |
|---|---|
| **Tài chính** | "Concert tickets = rent money short → need another gig" |
| **Hệ quả** | "Eat too much at food truck gig → food coma → miss meeting" |
| **Tình cờ** | "Side gig introduces Oren to potential music collab" |

#### C. Bài học — Tự nhiên, RELATABLE, KHÔNG ÉP

| Lĩnh vực | Trigger | Bài học |
|---|---|---|
| **Tài chính** | Impulse buying → can't pay rent | Budgeting, saving |
| **Work-life balance** | Overwork → burn out → miss friend's birthday | Boundaries |
| **Relationships** | Forget anniversary → Pinki upset | Communication, effort |
| **Health** | Energy drinks + no sleep → crash | Self-care |
| **Honesty** | Take shortcut → backfires | Integrity |
| **Friendship** | Disagree with Simon → awkward → resolve | Mature conflict resolution |

### 7. Visual Style (KHÔNG MÔ TẢ TRONG SHOT)

Script **KHÔNG** mô tả style/aesthetic/lighting trong Visual. Chỉ mô tả **hành động vật lý** và **context/environment**.
Style được xử lý riêng bởi template hoặc style_prompt.

> **CẤM trong Visual**: "warm lighting", "cinematic look", "3D render", "cartoon style".
> **CHỈ viết**: camera setup + subject + action + context + audio.

### 8. Oren's Body Language — ADULT VERSION (BẮT BUỘC)

Oren có HÌNH HÀI alien nhưng HÀNH VI hoàn toàn như **YOUNG ADULT AMERICAN**:

| ❌ SAI (childish) | ✅ ĐÚNG (adult) |
|---|---|
| Nhảy lên vui mừng, flail arms | Fist-pump contained, head nod, subtle smile |
| Nắm tay cứng nhắc, two-hand hold | Confident one-hand hold, casual grip |
| Ngồi thẳng đơ / bounce trên ghế | Lean back, one arm on armrest, ankle on knee |
| Skip/run nhỏ nhắn | Walk with swagger, confident stride, hands in pockets |
| Eyes wide open mọi lúc | Eyes half-lidded chill, opens wide chỉ khi shocked |
| Đứng bất động chờ | Lean on counter, tap phone, gõ nhịp trên đùi |

**Đặc trưng Oren (adult)**: Headphones luôn trên cổ hoặc đầu. Tay luôn gõ nhịp. Hay dựa tường/counter. Đi bộ có swagger nhẹ. Ngồi thì lean back, relax posture. Uống coffee/beer bằng 1 tay.

### 9. Audio Storytelling

- **60%+ screen time là SFX-only** — ambient Mỹ + Oren beat-boxing/humming
- **American Adult Ambient SFX**: coffee machine, keyboard typing, car engine, apartment AC hum, skateboard wheels, beer can opening, microwave beep, city traffic
- **Music-driven**: Oren luôn tạo nhịp — beat-box, finger-snapping, desk-tapping, humming
- **Narrator tối thiểu**: Casual English, witty, "So there's Oren, broke as usual... but this time it's different."
- **ZERO text on screen**

### 10. Ensuring Adult Appearance in Video Output

> **VẤN ĐỀ**: Sprunki character là thiết kế 2D stylized. AI video model có thể render chúng "childish" nếu không có chỉ dẫn.

#### Chiến lược thể hiện "người trưởng thành" trong video cuối cùng:

1. **Reference Image**: Tạo reference image với prompt nhấn mạnh adult proportions:
   - "adult proportions, longer limbs, taller body relative to objects"
   - "holding coffee cup, leaning on bar counter" (adult context props)
   - Xem `references/character-reference-prompts.md` cho prompt chuẩn

2. **Script-level signals (trong [NOTE])**: Luôn mô tả **HÀNH VI người lớn** trong visual description:
   - ✅ "Oren leans against kitchen counter, coffee in one hand, phone in the other, scrolling"
   - ✅ "Oren and Simon sit at bar, beers on counter, watching game on TV above"
   - ✅ "Oren walks through office corridor carrying laptop bag"
   - ❌ "Oren stands in front of toy store" → screams "kid"
   - ❌ "Oren runs excitedly toward ice cream truck" → childish

3. **Environment (context cues)**: Môi trường tự đẩy nhận thức "adult":
   - Apartment (có sofa, coffee table, kitchen island)
   - Office/coworking space
   - Bar/pub/rooftop
   - Coffee shop (đang work on laptop)
   - Gym, grocery store, laundromat

4. **Props that signal adulthood**: 
   - Coffee mug, beer/wine glass, laptop, car keys
   - Wallet with cards (not coins), phone with banking app
   - Lease papers, tax forms (comedy), grocery bags
   - Work badge, business cards

5. **Body proportions in prompt**: Khi mô tả action, luôn ngầm dùng adult scale:
   - "reaches up to kitchen cabinet" (adult height)
   - "leans on car hood" (car-scale, not toy-scale)
   - "sits on bar stool, feet on footrest" (bar stool height)

6. **Dialogue content signals maturity**:
   - Rent, bills, taxes, work deadlines, relationship problems
   - Coffee (not juice box), beer (not soda), ramen (broke adult staple)
   - "I have a meeting at 9", not "I have school at 9"

---

## Adult Environment Guide

### Common Locations

| Location | Chi tiết Mỹ | Phù hợp |
|---|---|---|
| Studio apartment | Small, messy, mix decor, kitchen counter bar stools, records on wall | Home scenes |
| American kitchen | Island counter, coffee maker, microwave, fridge with magnets/bills | Cooking |
| Coffee shop | Laptop on table, latte art, power outlets, coworking vibe | Work/hangout |
| Bar/Pub | Beer taps, TV showing game, dim lighting, bar stools, pool table | Social |
| Coworking space | Standing desks, whiteboards, meeting pods | Work |
| Food truck park | Picnic tables, string lights, diverse food options | Eating, social |
| Record store / Music shop | Vinyl racks, turntables, headphone stations | Oren's world |
| Rooftop | City skyline view, fairy lights, mismatched furniture, sunset | Romantic/group chill |
| Grocery store | Aisles, self-checkout, comparing prices | Adult life |
| Laundromat | Folding tables, TV on wall, magazine rack | Adulting comedy |
| Gym | Free weights, treadmill, locker room | Health episodes |
| Skate park | Half-pipe, graffiti walls, benches | Weekend hobby |

### Props (Adult Context)

**Tech**: Smartphone, laptop, headphones (ALWAYS), skateboard, camera
**Food**: Pizza (still), burgers, ramen, coffee (daily), beer/wine (social), energy drinks
**Work**: Laptop bag, meeting notebook, earbuds, USB drive, invoice
**Music**: Headphones (ALWAYS), DJ controller, vinyl records, portable speaker, MIDI keyboard
**Finance**: Banking app ($3.47), credit card, rent notice, bills, direct deposit notification
**Home**: Coffee maker, takeout containers, mismatched plates, apartment key
**Social**: Beer cans/bottles, board game, streaming on laptop, aux cord

---

## Voice Profile Integration — Pipeline Guide

### Cách voice profiles chảy vào video cuối cùng

```
Script [Character] tag → Parser trích Dialogue field
                              ↓
Character.VoiceStyle (DB) → buildShotContexts() → shotContext.CharacterVoices
                              ↓
shotContext JSON (có dialogue + character_voices) → LLM distill prompt
                              ↓
LLM viết video_prompt (bao gồm lip-sync/speaking cues)
                              ↓
getMouthConstraint() → inject "Dialogue: {text}. Character is actively speaking..."
                              ↓
video_prompt_distilled (final video prompt = LLM narrative + dialogue + mouth constraint)
```

### Hành động cần thiết khi tạo Drama

1. **Tạo Characters** trong Drama với đầy đủ `voice_style`:
   ```
   Oren → voice_style: "Male, mid-20s, laid-back American. Medium-low pitch, slightly raspy..."
   Pinki → voice_style: "Female, mid-20s, warm and bright. Medium-high pitch..."
   Simon → voice_style: "Male, mid-20s, high energy, extroverted..."
   ```

2. **Script dùng đúng tên nhân vật** trong dialogue tags:
   ```
   [Oren] Yo. Rent's due tomorrow, huh.
   [Pinki] Did you even check your account?
   ```

3. **Pipeline tự động**:
   - `parseStructuredShots()` → trích dialogue → lưu `Storyboard.Dialogue`
   - `buildShotContexts()` → load character voice_style → `shotContext.CharacterVoices`
   - Distill LLM nhận cả `dialogue` + `character_voices` → viết prompt có lip-sync cues
   - `getMouthConstraint()` → append vào `video_prompt_distilled`

4. **Kết quả**: Video prompt cuối cùng chứa:
   - LLM-generated motion narrative (có thể tham chiếu speaking/voice tone)
   - `Dialogue: [text]. The character is actively speaking, lip-syncing naturally...`

---

## References

- HL Meow Meow skill (inspiration): `skills/hl-meow-meow-script/SKILL.md`
- Shot format template: `templates/shot-format.md`
- Character deep reference: `references/sprunki-characters.md`
- Episode ideas: `references/episode-ideas.md`
- Voice profiles: `references/voice-profiles.md`
- Character reference image prompts: `references/character-reference-prompts.md`
