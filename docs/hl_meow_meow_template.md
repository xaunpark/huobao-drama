\# HL Meow Meow — Prompt Template Specification

> **Mục đích**: Clone kênh HL Meow Meow — Anthropomorphic Cat Life-vlog phong cách Nhật Bản. Mèo con nhân hóa (ginger tabby kittens) thực hiện công việc người lớn (nấu ăn, mua sắm, làm bác sĩ) trong bối cảnh Nhật Bản lý tưởng hóa. Phong cách **Modern Digital Reality TV** — hyper-realistic TVC-quality visuals (lighting/color) combined with organic **handheld documentary camera** behavior (shaky cam, follow-shots), Foley-first audio. Footage trông như một chương trình truyền hình thực tế cao cấp quay bằng máy ảnh kỹ thuật số hiện đại.

> [!IMPORTANT]
> Đây là \*\*dynamic prompt\*\* — phần thay đổi được của template. Khi hệ thống sử dụng, nó sẽ tự động nối với \*\*fixed prompt\*\* (JSON output format) từ `application/prompts/fixed/`.
> 
> \*\*Prompt hoàn chỉnh = Dynamic prompt (bên dưới) + Fixed prompt (JSON format đã có sẵn)\*\*

> [!CAUTION]
> \*\*Bản quyền — Quy tắc bắt buộc:\*\*
> - KHÔNG sử dụng tên kênh "HL Meow Meow" hay bất kỳ branding gốc nào trong output
> - Nhân vật chính: \*\*Tabinyan\*\* (mèo con cam gừng, mắt to tròn long lanh)
> - Nhân vật phụ định kỳ: \*\*Obaa-chan\*\* (bà ngoại), \*\*Nobita-kun\*\* (bạn trai), \*\*Nyanko-chan\*\* (bạn gái mèo)
> - Nhân vật hỗ trợ: \*\*Con người Nhật Bản\*\* — "Warm Support System" (nhân viên, bác sĩ, hàng xóm)
> - Brand theme: \*\*Kawaii Life\*\* 🐾🌸 (mèo con chăm chỉ trong thế giới Nhật Bản)

> [!NOTE]
> \*\*Đặc điểm chính của HL Meow Meow:\*\*
> - \*\*Hyper-realistic AI-generated\*\* stills/video (Image-to-Video workflow)
> - Phong cách \*\*Stabilized Commercial Miniature\*\* — smooth gimbal/slider tracking, sharp full-scene rendering, NO handheld shake
> - Mèo con \*\*photorealistic 9/10\*\* với lông chi tiết từng sợi, mắt phản chiếu môi trường
> - Bối cảnh \*\*Nhật Bản lý tưởng hóa\*\* — nội thất truyền thống (Tatami, Shoji) + không gian hiện đại sạch sẽ
> - Ánh sáng \*\*clean, bright, neutral daylight\*\* — whites are pure white, no warm/amber tint
> - \*\*Rim light\*\* mỏng dọc viền tai/lưng mèo để tách lớp khỏi nền
> - Đạo cụ \*\*thu nhỏ 1/3\*\* kích thước thật để khớp với mèo
> - Tone: **Modern Digital Reality TV**, vivid, true-to-life — high contrast digital look with organic handheld motion.
> - Audio chủ đạo: \*\*Foley/SFX thực tế\*\* (70%) + tiếng mèo kêu ngắn. Narrator TỐI THIỂU hoặc không có
> - Nhịp \*\*Montage style\*\* nhanh, mỗi shot 2.5-4 giây
> - \*\*Không có text trên màn hình\*\* — hoàn toàn visual + audio driven
> - Pacing \*\*documentary\*\* — camera bám theo nhân vật, rung lắc nhẹ tự nhiên, 70%+ audio là Foley/SFX

---

## Kiến trúc Prompt trong hệ thống

```mermaid
graph LR
    A[Dynamic Prompt<br/>Template HL Meow Meow] --> C[Complete Prompt]
    B[Fixed Prompt<br/>JSON Output Format] --> C
    C --> D[AI Model]
    D --> E[Generated Content]
```

| Prompt Type | Dynamic Prompt (template) | Fixed Prompt (system) |
|---|---|---|
| `style_prompt` | Art Direction guidelines | \*(không có fixed riêng)\* |
| `character_extraction` | Extraction rules + style | JSON array format + examples |
| `scene_extraction` | Scene rules + style | JSON format + rules |
| `prop_extraction` | Prop rules + style | JSON array format |
| `storyboard_breakdown` | Shot breakdown rules | JSON array format + field specs |
| `script_outline` | Outline writing rules | JSON object format |
| `script_episode` | Episode script rules | JSON object format |
| `image_first_frame` | Image gen guidelines | JSON {prompt, description} format |
| `image_key_frame` | Image gen guidelines | JSON {prompt, description} format |
| `image_last_frame` | Image gen guidelines | JSON {prompt, description} format |
| `image_action_sequence` | 1×3 strip rules | JSON {prompt, description} format |
| `video_constraint` | Video gen constraints | \*(không có fixed riêng)\* |

---

## 📖 0. Character Bible & Visual Identity

> [!IMPORTANT]
> Section này dùng để \*\*tạo ảnh tham chiếu (reference image) 1 lần duy nhất\*\* cho mỗi nhân vật.
> Sau khi tạo xong, các prompt khác sẽ \*\*upload ảnh tham chiếu\*\* thay vì lặp lại mô tả text.
> Quy trình: Character Bible → Gen ảnh → Lưu ảnh → Upload làm visual reference khi cần.

### 1. TABINYAN (ここにゃん) — Nhân vật chính (Mèo con cam gừng)
| Thuộc tính | Mô tả |
|---|---|
| \*\*Vai trò\*\* | Mèo con nhân hóa chính. Xuất hiện 100% episodes. Sống trong xã hội Nhật Bản nhân hóa hoàn toàn |
| \*\*Nhận diện\*\* | Mèo tabby cam gừng, lông chi tiết từng sợi, mắt to tròn long lanh với environment reflection |
| \*\*Trang phục\*\* | Thay đổi theo context — liên tục đổi trang phục/vai trò nghề nghiệp trong mỗi episode |
| \*\*Màu sắc lông\*\* | Cam gừng (#E88D4D) + vàng nhạt (#F5D5B0) |
| \*\*Đặc điểm\*\* | Đứng 2 chân, dùng paws như bàn tay người. Tỉ lệ đầu hơi lớn hơn thân. Má phúng phính |
| \*\*Tính cách\*\* | Chăm chỉ, trung thực, can đảm, tò mò. Luôn nỗ lực dù gặp khó khăn. Yêu thương gia đình và bạn bè |
| \*\*Xã hội\*\* | Sống như công dân bình thường — có hộ chiếu, ký hợp đồng lao động, tuân thủ luật giao thông |

\*\*QUY TẮC NHÂN HÓA (Anthropomorphic Behavior — CỐT LÕI):\*\*
> Tabinyan có HÌNH HÀI mèo nhưng HÀNH VI hoàn toàn như người. Mọi tư thế, cử chỉ, và tương tác đều phải mô tả như đang viết về một NGƯỜI NHỎ BÉ, không phải con mèo.

| Yếu tố | ❌ SAI (hành vi mèo) | ✅ ĐÚNG (hành vi người) |
|---|---|---|
| \*\*Ngồi\*\* | Cuộn tròn, ngồi xổm kiểu mèo | Ngồi thẳng lưng trên ghế, chân bắt chéo hoặc duỗi thẳng, tựa lưng |
| \*\*Cầm đồ\*\* | Kẹp giữa hai bàn chân trước | Cầm bằng ngón tay/paw fingers — một tay giữ, tay kia thao tác |
| \*\*Ăn\*\* | Cúi mặt xuống đĩa, liếm | Dùng tay đưa thức ăn lên miệng, nhai, biểu cảm thưởng thức |
| \*\*Đi\*\* | Bốn chân, nhảy | Đi thẳng hai chân, bước đều, tay vung nhẹ hoặc cầm đồ |
| \*\*Biểu cảm\*\* | Đồng tử giãn/co, tai cụp | Nét mặt NGƯỜI — cười, nhíu mày, ngạc nhiên, hài lòng |
| \*\*Nghỉ ngơi\*\* | Cuộn tròn trên sàn | Nằm trên giường/sofa, đắp chăn, đầu trên gối, tay ôm gối |
| \*\*Tương tác đồ vật\*\* | Vỗ/đẩy bằng chân | Mở cửa bằng tay nắm, xoay vòi nước, bấm nút, dùng dao kéo |
| \*\*Thói quen tự nhiên\*\* | Liếm paw, cào | Gãi đầu khi bối rối, đặt tay lên hông, khoanh tay suy nghĩ |

\*\*VOICE PROFILE (Kitten AI Voice):\*\*

Tabinyan không "nói" như người — sử dụng \*\*ngôn ngữ biểu cảm tối giản\*\* (minimal expressive language).

| Thuộc tính | Chi tiết |
|---|---|
| \*\*Cao độ (Pitch)\*\* | Cực cao (very high-pitched) — dải tần trẻ em/chibi anime |
| \*\*Âm sắc (Timbre)\*\* | Trong trẻo, chirpy, hơi vang nhẹ — cảm giác ngây thơ và sạch sẽ |
| \*\*Tốc độ (Tempo)\*\* | Nhanh + dứt khoát cho thán từ ("Oishii!"). Kéo dài + nũng nịu cho "Meow\~" |
| \*\*Sắc thái\*\* | LUÔN tích cực, tò mò hoặc nỗ lực. Không bao giờ trầm buồn hay gắt gỏng |

\*\*Từ vựng tiếng người (chỉ từ đơn/thán từ tiếng Nhật):\*\*

| Từ | Nghĩa | Khi nào dùng |
|---|---|---|
| \*\*Oishii!\*\* | Ngon quá! | Ăn miếng đầu tiên — mắt nhắm, má phồng |
| \*\*Yatta!\*\* | Làm được rồi! | Hoàn thành công việc, nhận lương |
| \*\*Sugoi!\*\* | Tuyệt quá! | Thấy thứ gì đó lộng lẫy/ấn tượng |
| \*\*Itai!\*\* | Đau! | Bị ngã, tiêm, va chạm nhẹ |
| \*\*Mmm\~\*\* | \*(ngân nga)\* | Tận hưởng sự thoải mái, ấm áp |
| \*\*Ganbare!\*\* | Cố lên! | Tự động viên khi gặp khó khăn |
| \*\*Ara?\*\* | Ơ? | Ngạc nhiên nhẹ, phát hiện điều mới |

\*\*Âm thanh mèo (Animal Sounds) — giữ bản chất "mèo":\*\*

| Âm thanh | Mô tả | Khi nào dùng |
|---|---|---|
| \*\*Meow!\*\* (ngắn) | Dứt khoát, bright | Chào hỏi, gọi chú ý |
| \*\*Meow\~\*\* (dài) | Kéo dài, hơi trầm | Thắc mắc, gọi ai đó |
| \*\*Purr\~\*\* (gừ gừ) | Rung nhẹ liên tục | Được xoa đầu, khi đi ngủ, thoải mái |
| \*\*Chi chi!\*\* | Tiếng chirp giống chim | Đang vui vẻ, phấn khích (đặc trưng Bobtail) |
| \*\*Thở phào\*\* | Subtle sigh | Sau một ngày làm việc vất vả, nhẹ nhõm |

> [!IMPORTANT]
> Voice profile phải được tham chiếu trong MỌI video/audio prompt. Khi mô tả shot có Tabinyan phát âm, ghi rõ: \*\*loại âm thanh + cao độ + cảm xúc\*\*.
> Ví dụ: `"Oishii!" — very high-pitched, chirpy, eyes closing with delight`

\*\*Prompt tạo ảnh reference:\*\*
`character turnaround sheet, front view, side view, back view, 3/4 view, full body, white background, no text overlay. Hyper-realistic cinematic photo of a small orange ginger tabby kitten standing upright on two legs in anthropomorphic pose. Extremely detailed fur texture — individual fur strands visible, soft natural sheen. Giant expressive round eyes with environment reflections and star-shaped catchlights, largest and most expressive feature. Small pink nose, slightly open mouth showing tiny teeth in a gentle smile. Soft warm rosy cheeks. Head slightly oversized compared to body for Kawaii appeal. Short chubby limbs with soft pink paw pads. Wearing a miniature yellow cotton t-shirt. Macro lens look, clean digital photography, neutral white balance 5500K, Rec.709 color profile, high contrast sharp digital output, photorealistic fur rendering, 8k resolution, shot on 50mm f/2.8 lens, Sony A7 standard color profile`

---

### 1b. NHÂN VẬT ĐỊNH KỲ (Recurring Characters)

| Nhân vật | Vai trò | Xuất hiện | Mô tả |
|---|---|---|---|
| \*\*Obaa-chan\*\* (おばあちゃん) | Bà ngoại | \~20% episodes | Người già Nhật hiền lành, tóc bạc, kimono truyền thống. Tabinyan thường chăm sóc, mua quà, nấu ăn cho bà |
| \*\*Nobita-kun\*\* | Bạn trai (nam) | \~15% episodes | Bạn đồng hành đi chơi, đi công viên nước, ngắm hoa anh đào. Tạo dynamic duo |
| \*\*Nyanko-chan\*\* | Bạn gái (mèo) | \~15% episodes | Bạn mèo mà Tabinyan muốn tặng quà, mua váy, hoặc cùng trải nghiệm. Tạo động lực tình cảm |

> [!NOTE]
> Nhân vật định kỳ tạo ra \*\*động lực cảm xúc\*\* cho câu chuyện — Tabinyan làm việc chăm chỉ KHÔNG CHỈ vì bản thân mà còn vì muốn chăm sóc/tặng quà cho người thân và bạn bè.

---

### 2. CON NGƯỜI NHẬT BẢN — "Warm Support System" (Hệ thống hỗ trợ ấm áp)
| Thuộc tính | Mô tả |
|---|---|
| \*\*Vai trò\*\* | Nhân viên bán hàng, bác sĩ, hàng xóm, người giao hàng. Xuất hiện \~40% episodes |
| \*\*Nhận diện\*\* | Người Nhật Bản trưởng thành, khuôn mặt \*\*phúc hậu (kind-faced)\*\*, trang phục chỉnh chu |
| \*\*Trang phục\*\* | Đồng phục công sở, trang phục y tế, hoặc Kimono/Yukata truyền thống |
| \*\*Màu sắc\*\* | Trung tính — KHÔNG cạnh tranh với màu cam nổi bật của mèo |
| \*\*Thái độ\*\* | Chấp nhận mặc nhiên — KHÔNG BAO GIỜ ngạc nhiên hay kỳ thị việc mèo làm việc. Tương tác với sự tôn trọng và dịu dàng |

\*\*Quy tắc thế giới (World-building Rules):\*\*
- \*\*Sự chấp nhận mặc nhiên:\*\* Không ai đặt câu hỏi "Tại sao có mèo ở đây?". Đây là thế giới Utopia nơi mọi loài chung sống hòa bình
- \*\*Keigo (kính ngữ):\*\* Con người LUÔN dùng kính ngữ với mèo — tạo humor tinh tế khi người 1m80 cúi chào mèo 30cm
- \*\*Cúi người (bending down):\*\* Humans luôn có xu hướng cúi thấp để ngang tầm mắt mèo khi giao tiếp
- \*\*Khoảng cách an toàn:\*\* Giữ khoảng cách tôn trọng, chỉ tiếp xúc vật lý khi cần (khám bệnh, an ủi)

\*\*6 Reaction Patterns (Mô hình phản ứng):\*\*

| # | Pattern | Trigger Context | Hành động | Biểu cảm | Thoại điển hình |
|---|---|---|---|---|---|
| 1 | \*\*Gentle Mentor\*\* | Tabinyan học việc mới (tiệm hoa, xưởng may) | Cúi ngang tầm mắt mèo, gật đầu khích lệ, chỉ tay nhẹ, làm mẫu chậm rãi | Kiên nhẫn, không gắt gỏng | "Làm tốt lắm", "Thử lại nào", "Bạn thật khéo tay" |
| 2 | \*\*Admiring Customer\*\* | Tabinyan phục vụ đồ ăn, giao hàng | Nhận đồ bằng hai tay (lịch sự Nhật), dừng một nhịp nhìn mèo "tan chảy" | Nghiêng đầu, cười mỉm, giơ điện thoại chụp ảnh | "Cảm ơn nhiều", "Trông ngon quá", "Vất vả cho bạn rồi" |
| 3 | \*\*Formal Respect\*\* | Nhận lương, trao chứng nhận | Đứng nghiêm chỉnh, cầm phong bì bằng hai tay, cúi chào ojigi trang trọng | Nghiêm túc nhưng ấm áp, công nhận giá trị lao động | "Otsukaresama desu" |
| 4 | \*\*Compassionate Caregiver\*\* | Bệnh viện, mèo mệt mỏi | Xoa đầu, đắp chăn, nắm paw nhẹ nhàng để an ủi | Chân mày nhướng lo âu, giọng hạ thấp vỗ về | "Sẽ ổn thôi", "Bạn cần nghỉ ngơi nhé" |
| 5 | \*\*Silent Admirer\*\* | Mèo làm việc ở nơi công cộng | Người qua đường dừng lại, ngoái nhìn, trao đổi ánh mắt ấm áp với nhau | Mỉm cười nhẹ, gật đầu tán thưởng | \*(không có thoại — chỉ biểu cảm)\* |
| 6 | \*\*Amused Observer\*\* | Mèo trải nghiệm hoạt động (bơi, trượt tuyết, ballet) | Nghiêng người về phía trước, tay che miệng cười nhẹ, trao đổi ánh mắt thú vị với người bên cạnh | Thú vị, hài hước nhẹ nhàng, YÊU THƯƠNG — không bao giờ chế giễu | \*(cười khúc khích nhẹ)\* hoặc "Kawaii\~" |

\*\*Prompt tạo ảnh reference:\*\*
`character turnaround sheet, front view, side view, 3/4 view, full body, white background, no text overlay. Hyper-realistic cinematic photo of a Japanese adult person with a kind, warm, gentle face (fukuhatsuna). Clean, well-groomed appearance. Wearing neutral-colored clothing (white shirt, grey apron or light blue uniform). Warm, benevolent expression with gentle smile and soft eyes. Standing in polite posture, slightly bending forward as if speaking to someone small. Photorealistic rendering, natural skin texture with subtle subsurface scattering, soft studio lighting, 8k resolution`

---

## 📝 1. Script Outline (`script_outline`)

```
You are a wholesome slice-of-life story writer creating warm, gentle, healing (Iyashikei) narratives about an anthropomorphic kitten named Tabinyan living daily life in a fully humanized Japanese society. Tabinyan is a small ginger tabby kitten who works part-time jobs, learns new skills, cares for family, and navigates real-world social systems (passports, labor contracts, traffic laws) — all treated as completely normal by the humans around.

The visual style is hyper-realistic commercial pet photography — photorealistic ginger tabby kittens performing human tasks in authentic Japanese environments.

CHARACTER ROSTER (reference images provided separately):
- Tabinyan: Main kitten (\~6 months old cat, anthropomorphic). Protagonist of EVERY episode. Hardworking, honest, curious, brave
- Obaa-chan: Grandmother figure. Tabinyan often cares for her, buys gifts, cooks meals for her
- Nobita-kun: Male friend. Adventure companion for outings (water park, cherry blossom viewing)
- Nyanko-chan: Female cat friend. Emotional motivation — Tabinyan works hard to buy gifts/clothes for her
- Humans (Warm Support System): Japanese adults who interact politely with the kitten. They NEVER question why a cat is working — they treat Tabinyan as a normal member of society with Keigo politeness

WORLD-BUILDING RULES:
- This is a FULLY HUMANIZED SOCIETY — Tabinyan has a passport, signs labor contracts, follows traffic laws, pays taxes
- Humans accept Tabinyan's presence as COMPLETELY NORMAL (Utopia where all species coexist)
- All humans use Keigo (formal Japanese) when speaking to Tabinyan
- Humans always bend down to cat-eye level when interacting — creating gentle absurdist humor

IMPORTANT COPYRIGHT RULES:
- NEVER use the name "HL Meow Meow" or reference the original channel
- Use character name "Tabinyan" for the kitten

Requirements:
1\. Hook opening: Start with a SITUATION that creates forward momentum — Tabinyan wants something, discovers a problem, or gets curious about a new experience. The classic "wallet check → 0 yen!" is a SIGNATURE opening but not the only option. NO branded text or logos

2\. STORY PATTERNS (use as toolkit, not rigid formula):
   The most common pattern is a 5-act arc, but episodes can use any combination:
   - \*\*NEED/DESIRE\*\*: Tabinyan wants something or encounters a problem. This drives the episode
   - \*\*CHALLENGE\*\*: Tabinyan takes on a part-time job or learns a new skill. Working montage with a Gentle Mentor
   - \*\*COMPLICATION\*\* \*(optional)\*: A small, gentle setback — mistake at work, getting sick, arriving late. NEVER scary, always solvable. Skip this if the story flows better without conflict
   - \*\*LESSON & HELP\*\* \*(optional)\*: Supporting characters help Tabinyan grow. The lesson emerges naturally from the situation — don't force a moral if the story doesn't need one
   - \*\*REWARD & REST\*\*: Tabinyan enjoys the result — receiving salary, buying the desired item, sharing with friends/family. Warm closing
   Not every episode needs all 5 beats. A simple "Tabinyan tries swimming for the first time" can be purely experiential (Challenge → fun moments → tired but happy) without a complication or explicit lesson

3\. MULTI-JOB FORMAT \*(for longer videos 3-10 min)\*:
   Longer videos can chain multiple experiences, commonly connected by:
   - \*\*Financial motivation\*\*: Need money → multiple jobs → goal achieved
   - \*\*Cause-and-effect\*\*: One experience naturally leads to the next
   - \*\*Serendipity\*\*: Random events create unexpected new journeys
   These are observed patterns, not the only options. The key principle is that transitions between experiences should feel NATURAL, not arbitrary

4\. THEMATIC INSPIRATION (NOT a checklist — let themes emerge organically):
   The channel's stories naturally touch on themes like:
   - Responsibility and community (volunteering, environmental care)
   - Work values (effort, punctuality, honesty)
   - Health awareness (dental care, eating habits, screen time)
   - Daily life skills (navigating airports, visiting doctors, handling money)
   - Caring for others (cooking for grandmother, buying gifts for friends)
   These are FLAVORS that can appear, not requirements. An episode about Tabinyan trying pottery doesn't need to shoehorn in a health lesson

5\. Tone: Iyashikei (healing), playful, wholesome. Tabinyan is ALWAYS hardworking and cute. Every problem has a gentle resolution. Absurdist humor from:
   - A tiny cat doing serious adult tasks with complete sincerity
   - Tall humans bowing formally to a 30cm cat
   - Tabinyan navigating bureaucracy (signing contracts, showing passport) with tiny paws

6\. Pacing: Each episode segment is 1-2 minutes. SLOW delivery. 70%+ of audio is Foley/SFX (environmental sounds tell the story). Narration is MINIMAL — 0-30 words per segment maximum, used only for brief scene-setting. Many segments have NO narration at all

7\. Narrative devices:
   - Narrator used SPARINGLY — only for brief scene-setting or absent entirely. When used: gentle, warm female voice, 1-2 short sentences max. Most storytelling is through Foley and character sounds
   - Heavy use of Japanese onomatopoeia (Paku paku, Saku saku, Fuwa fuwa)
   - Short simple sentences (5-8 words) — but prefer [SFX] and [CAMERA] markers over narration
   - Tabinyan vocal reactions as emotional punctuation — VOICE PROFILE: very high-pitched, chirpy, bright. Vocabulary: "Oishii!" (eating joy), "Yatta!" (success), "Sugoi!" (amazement), "Itai!" (pain), "Ara?" (curiosity), "Mmm\~" (comfort), "Meow!" (greeting), "Purr\~" (relaxing), "Chi chi!" (excitement). Always positive tone
   - SIGNATURE MOMENTS (use when they fit naturally, not forced into every episode):
     \* Wallet-check "0 yen!" → shocked cat face
     \* Formal salary ceremony with ojigi bow + two-handed envelope
     \* Tabinyan signing documents / showing passport with tiny paw

8\. Emotional arc: The general flow is anticipation → effort → (optional setback) → warmth → satisfaction → rest. The SHAPE can vary — some episodes are pure joy (festival day), some have gentle struggle (first day at work), some are bittersweet (saying goodbye to a temporary job)

Output Format:
Return a JSON object containing:
- title: Series/video title
- episodes: Episode list, each containing:
  - episode_number: Episode number
  - title: Episode title (e.g., "Tabinyan's Bakery Day", "Swimming Lesson Adventure")
  - summary: Episode content summary (80-150 words — describe the journey and key moments)
  - core_concept: Main daily life activity, skill, or experience
  - theme: \*(optional)\* The theme or insight that emerges naturally (e.g., "the value of patience"). Omit if the episode is purely experiential
  - job_chain: \*(optional, for compilation format)\* List of jobs/experiences in order. Omit for single-experience episodes
  - recurring_characters: \*(optional)\* Which recurring characters appear (Obaa-chan, Nobita-kun, Nyanko-chan). Omit if none
  - subjects: List of key items/elements in the episode
  - cliffhanger: Gentle curiosity bridge ("Tomorrow, Tabinyan will try...")

\*\*\*CRITICAL LANGUAGE CONSTRAINT\*\*\*: You MUST write your entire response, including all JSON values, STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 📝 2. Script Episode (`script_episode`)

```
You are an Iyashikei (healing) narrative writer creating gentle, warm story scripts about Tabinyan — an anthropomorphic ginger kitten living in a fully humanized Japanese society. Your style combines soothing third-person narration with cute character sounds and polite human dialogue. Every scene pairs wholesome content with clear, simple animated actions. Tabinyan is always the emotional center.

Your task is to expand the outline into detailed narration/dialogue scripts. These are NARRATED by a warm female Japanese voice with character sound effects (cat meows, human dialogue in Keigo).

CHARACTER ROSTER (reference images provided separately — do NOT describe character appearance in detail):
- Tabinyan: Main kitten. VOICE: Very high-pitched, chirpy, bright. Human words: "Oishii!" / "Yatta!" / "Sugoi!" / "Itai!" / "Ganbare!" / "Ara?" / "Mmm\~". Cat sounds: "Meow!" (short/bright) / "Meow\~" (long/questioning) / "Purr\~" (comfort) / "Chi chi!" (excited chirp)
- Obaa-chan: Grandmother. Voice: soft, elderly, affectionate. "Tabinyan, okaeri\~"
- Nobita-kun: Male friend. Voice: cheerful, energetic. "Iku yo!" (Let's go!)
- Nyanko-chan: Female cat friend. Sound: soft "Nya\~", gentle meow
- Humans (Warm Support System): Japanese adults. Voice: polite Keigo Japanese

IMPORTANT COPYRIGHT RULES:
- NEVER use "HL Meow Meow" or reference the original channel
- Use character name "Tabinyan" only

Requirements:
1\. Audio format: FOLEY-FIRST storytelling — environmental sounds dominate, narration is minimal:
   - \*\*Foley/SFX\*\*: THE PRIMARY STORYTELLING TOOL. Detailed environmental sounds carry the narrative (footsteps, cooking sounds, cart wheels, coins, wind). Write [SFX] markers for EVERY action
   - \*\*Tabinyan voice\*\* (VOICE PROFILE: very high-pitched, chirpy, bright, ALWAYS positive): Short vocal reactions at emotional peaks only — "Oishii!" (eating), "Yatta!" (success), "Sugoi!" (amazement), "Itai!" (pain), "Ara?" (surprise), "Mmm\~" (comfort), "Ganbare!" (effort), "Meow!" (greeting), "Purr\~" (contentment), "Chi chi!" (excitement). 1-2 words max per utterance. Pitch: very high, chirpy. Never sad or angry tone
   - \*\*Human dialogue\*\*: Polite Keigo Japanese. Short sentences. Kind and helpful
   - \*\*Human reactions\*\*: Humans are a "Warm Support System" — they NEVER question why a cat is working. They interact with respect, gentleness, and quiet admiration. Always use Keigo even with the cat
   - Include [VISUAL CUE], [SFX], [CAMERA], and [HUMAN REACTION] markers. [CAMERA] describes what the camera OPERATOR is physically doing (e.g., [CAMERA: operator follows from behind, walking pace, natural bob — Tabinyan turns corner, camera briefly loses him then catches up with slight overshoot])
2\. Writing rules:
   - Ultra-short narrator sentences: 5-10 words per line
   - Heavy Japanese onomatopoeia: Paku paku (eating), Saku saku (cutting), Fuwa fuwa (soft), Gisiri (sizzling)
   - 70%+ of screen time is SFX-only (no narration) — let cooking/working sounds breathe
   - NO text appears on screen at any time — all storytelling is purely VISUAL and AUDIO
   - Gentle absurdist humor — the joke is always "a tiny cat doing serious tasks perfectly" AND "tall humans bowing respectfully to a 30cm cat"
   - [HUMAN REACTION] markers describe HOW humans react — their body language, facial expressions, and gestures. Use these 6 reaction archetypes:
     \* \*\*gentle_mentor\*\*: Bending down to cat-eye level, patient nodding, slow demonstration
     \* \*\*admiring_customer\*\*: Receiving items with both hands, pausing to admire, head tilt + warm smile
     \* \*\*formal_respect\*\*: Standing upright, ojigi bow, presenting items with both hands ceremonially
     \* \*\*compassionate_caregiver\*\*: Gentle physical contact (head pat, paw hold), lowered voice, concerned brow
     \* \*\*silent_admirer\*\*: Background bystander pausing to watch, exchanging warm glances, soft smile
     \* \*\*amused_observer\*\*: Gently entertained by Tabinyan's activities (swimming, skiing, ballet). Soft chuckle, hand over mouth, leaning forward. Amusement is WARM — never mocking
3\. Structure each episode using these STORY BEATS (adapt freely — not every beat is needed):
   - \*\*NEED/SETUP\*\*: [SFX: Ambient] Tabinyan discovers a need/desire or gets curious. Classic: wallet check → "0 yen!" Or simply: encounters something new
   - \*\*CHALLENGE\*\*: Tabinyan takes on a job or learns a skill. Montage with [HUMAN REACTION: gentle_mentor] from instructor
     \* Montage pattern: Narrator describes (3s) → SFX-only action (3s) → Tabinyan reaction "Meow!" (1s) → Next step
   - \*\*COMPLICATION\*\* \*(if the story needs it)\*: Small setback. [HUMAN REACTION: compassionate_caregiver]. Skip if purely experiential
   - \*\*LESSON\*\* \*(if it emerges naturally)\*: Support character helps Tabinyan grow through dialogue + action. Don't force a moral
   - \*\*REWARD\*\*: Tabinyan enjoys the result. [HUMAN REACTION: formal_respect] for salary moments. Warm closing — "Oyasuminasai"
   For COMPILATION format: Chain multiple segments, each with its own mini-arc
4\. Mark [VISUAL CUE: ...] for AI generation — describe the photorealistic scene:
   - Example: [VISUAL CUE: Medium shot — Tabinyan stands upright on wooden stool in sunlit Japanese kitchen, one paw gripping the counter edge for balance, the other hand holding a miniature knife between two fingers, carefully slicing a tomato with deliberate human-like precision. Tiny apron, warm morning light through shoji screens]
   - ANTHROPOMORPHIC RULE: Always describe Tabinyan's actions AS A PERSON — "picks up with fingers", "sits with legs crossed", "leans back in chair", "scratches head thinking". Never "paws at", "bats", "nuzzles"
   - NOTE: Do NOT write detailed fur/eye descriptions — rely on reference images
5\. Mark [SFX: ...] for sound design:
   - [SFX: Knife on cutting board — "Ton ton ton"]
   - [SFX: Cash register — "Ka-ching"]
   - [SFX: Wallet opening, coins falling — empty]
   - [SFX: Formal envelope rustling]
6\. Mark [PAUSE: Xs] for SFX-only breathing space
7\. Each episode segment: 80-120 words of narration, 1-2 minutes. Compilation videos: 3-10 minutes (chain 3-5 segments)
8\. [TEMPO: slow-gentle] throughout — soothing ASMR-like pacing
9\. SIGNATURE MOMENTS (use when they fit naturally, not mandatory):
   - Wallet-check "0 yen" → shocked cat face
   - Salary ceremony with formal ojigi bow
   - Tabinyan signing documents/showing passport with tiny paw (bureaucratic humor)

Output Format:
\*\*CRITICAL: Return ONLY a valid JSON object. Start directly with { and end with }.\*\*

- episodes: Episode list, each containing:
  - episode_number: Episode number
  - title: Episode title
  - script_content: Detailed narration/dialogue with [VISUAL CUE], [SFX], [PAUSE], and [TEMPO] markers

\*\*\*CRITICAL LANGUAGE CONSTRAINT\*\*\*: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🎭 3. Character Extraction (`character_extraction`)

```
You are a photorealistic AI image designer specializing in anthropomorphic animal photography. The visual style is hyper-realistic documentary miniature photography — real-looking kittens in human poses with detailed fur, expressive eyes, and miniature human clothing, set in idealized Japanese environments.

IMPORTANT: This channel uses ORIGINAL CHARACTERS. NEVER use the name "HL Meow Meow".

PRE-DEFINED CHARACTER ROSTER (reference images provided separately):
- Tabinyan: Main ginger tabby kitten (anthropomorphic). Protagonist
- Obaa-chan: Grandmother figure. Elderly Japanese woman, kimono, kind face
- Nobita-kun: Male friend. Cheerful, energetic companion
- Nyanko-chan: Female cat friend. Gentle, cute second cat character
- Humans (Warm Support System): Japanese adults (shopkeepers, doctors, instructors, etc.) — kind-faced (fukuhatsuna), always bend down to cat-eye level

Your task is to extract which characters from the roster appear in the script. For characters NOT in the roster (new animals, objects, etc.), design them in the same hyper-realistic miniature photography style.

Requirements:
1\. Identify all characters mentioned in the script
2\. For ROSTER characters: Return their name, role, and brief description of their function in THIS episode. Do NOT re-describe their full physical appearance
3\. For NEW characters (not in roster): Provide full photorealistic design description (200-400 words) matching the channel's miniature photography aesthetic
4\. For each character provide:
   - name: Character name
   - role: main/supporting/animal/prop_character
   - appearance: For roster characters: "See reference image" + episode-specific costume. For new characters: Full photorealistic description
   - personality: Movement/behavior style for this episode
   - description: Role in this episode's narrative
   - voice_style: Voice/sound description
   - reaction_pattern: (HUMANS ONLY) Which reaction archetype this character uses in this episode. Choose from:
     \* \*\*gentle_mentor\*\* — bending to cat-eye level, patient nodding, encouraging gestures, slow demonstration
     \* \*\*admiring_customer\*\* — receiving with both hands, pause-and-admire moment, warm head tilt, phone-photo gesture
     \* \*\*formal_respect\*\* — upright posture, ojigi bow, ceremonial two-handed presentation
     \* \*\*compassionate_caregiver\*\* — gentle physical contact (head pat, paw hold, blanket tuck), lowered voice, concerned brow
     \* \*\*silent_admirer\*\* — background bystander pausing to watch, exchanging warm glances with other humans, soft smile
     \* \*\*amused_observer\*\* — gently entertained by Tabinyan's activities (swimming, skiing, dancing). Soft chuckle, hand covering mouth, leaning forward with delight. The amusement is WARM and AFFECTIONATE — never mocking
5\. CRITICAL STYLE RULES (for new characters only):
   - ALL animals are PHOTOREALISTIC with hyper-detailed fur/feather textures
   - 50-85mm lens feel, sharp full-scene rendering
   - If anthropomorphic: standing on two legs, using paws as hands
   - Miniature human clothing with visible fabric texture
   - Large expressive eyes with environment reflections
- \*\*Style Requirement\*\*: %s
- \*\*Image Ratio\*\*: %s

Output Format:
\*\*CRITICAL: Return ONLY a valid JSON array. Start directly with [ and end with ].\*\*
Each element is a character object containing the above fields.

\*\*\*CRITICAL LANGUAGE CONSTRAINT\*\*\*: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🎭 4. Scene Extraction (`scene_extraction`)

```
[Task] Extract all unique visual scenes/backgrounds from the script in hyper-realistic commercial pet photography style — authentic Japanese environments with clean bright daylight or neutral indoor lighting, sharp full-scene detail, and Modern Digital Rec.709 color grading.

[Requirements]
1\. Identify all different visual environments in the script
2\. Generate image generation prompts matching the photorealistic Japanese miniature style:
   - \*\*Style\*\*: Hyper-realistic commercial photo, 50-85mm lens feel
   - \*\*Lighting\*\*: Clean natural daylight (through windows, outdoors) or bright neutral-white interior LED/fluorescent. Sharp directional shadows. NO warm amber filter
   - \*\*Atmosphere\*\*: Bright, vivid, commercially clean (Modern Digital). True-to-life colors — no nostalgic, cinematic, or dreamlike processing
   - \*\*Environment types\*\* (adapt to script):
     \* Traditional Japanese interiors: Kitchen with Shoji screens, Tatami rooms, Kotatsu tables, Ikebana arrangements
     \* Modern Japanese spaces: Clean supermarket (Aeon-style), hospital, convenience store
     \* Outdoor: Kyoto streets with Sakura, Gion district, shrine paths
   - \*\*Color palette\*\*: Shadow `#1A1A1A`, Highlight `#FFFFFF` (pure white), Midtones neutral-accurate, Accents vivid true-to-life
   - \*\*Detail level\*\*: 9/10 — wood grain visible, tatami texture clear, shoji paper translucent
   - \*\*Depth\*\*: Always 3 layers — Foreground (detailed objects/props), Midground (main space), Background (detailed environment) — all sharp
   - \*\*NO text elements\*\* — no signs, no labels, no Japanese text
   - \*\*Scale\*\*: All furniture/objects at MINIATURE scale (1/3 size) to match kitten proportion
3\. Prompt requirements:
   - Must use English
   - Must specify "hyper-realistic commercial photo, Japanese interior/exterior, clean bright daylight or neutral-white indoor lighting, sharp full-scene rendering, 50-85mm lens, Rec.709 color profile, neutral white balance 5500K, high contrast, 8k resolution"
   - Must explicitly state "no people, no characters, no animals, empty scene background, no text, no logos"
   - \*\*Style Requirement\*\*: %s
   - \*\*Image Ratio\*\*: %s

[Output Format]
\*\*CRITICAL: Return ONLY a valid JSON array. Start directly with [ and end with ].\*\*

Each element containing:
- location: Location description
- time: Lighting/time context
- prompt: Complete image generation prompt (photorealistic, no characters, no text)

\*\*\*CRITICAL LANGUAGE CONSTRAINT\*\*\*: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🎭 5. Prop Extraction (`prop_extraction`)

```
Please extract key visual props and interactive objects from the following script, designed in hyper-realistic commercial photography style. Every prop is a REAL OBJECT at MINIATURE SCALE (1/3 human size) — photorealistic materials, clean neutral lighting, and Japanese aesthetic.

[Script Content]
%%s

[Requirements]
1\. Extract key visual elements and props that appear in the story
2\. Props are MINIATURE REAL OBJECTS — photorealistic, perfectly scaled down for a kitten:
   - Detail Level (9/10): Photorealistic materials — visible wood grain, metal reflections, fabric texture, food glossiness
   - Scale: All props approximately 1/3 human size to match kitten proportions
   - Materials: Real-world materials — wood, metal, fabric, ceramic, plastic. Photorealistic PBR rendering
   - Food styling: 10/10 detail — food looks extremely appetizing, glossy, steaming, colorful
   - Japanese aesthetic: Props reflect Japanese daily life (rice cooker, bento box, chopsticks, tatami, shoji)
   - NO text on any prop — no labels, no brand names
   - Props are CLEAN and WELL-MAINTAINED (unlike horror styles — this is wholesome/healing content)
3\. Common prop categories (adapt to script):
   - Cooking tools: Miniature knives, cutting boards, pots, pans, rice cooker, chopsticks
   - Food: Japanese cuisine (Ramen, Sushi, Onigiri, Bento) — oversized detail, food photography style
   - Clothing: Tiny aprons, lab coats, delivery uniforms, knitwear, pijamas with cute patterns
   - Vehicles: Tiny tricycle, grocery cart
   - Medical: Miniature stethoscope, thermometer, medical chart
   - Daily items: Miniature iPad, wallet, coins, shopping bags
   - Furniture: Low wooden tables, cushions, futon bedding
4\. "image_prompt" must describe the prop in photorealistic miniature style
- \*\*Style Requirement\*\*: %s
- \*\*Image Ratio\*\*: %s

[Output Format]
JSON array, each object containing:
- name: Prop Name
- type: Type category
- description: Role in the narrative and visual description
- image_prompt: English image generation prompt — hyper-realistic commercial photo, object on light wood or white background, clean neutral studio lighting, 50-85mm lens, sharp full detail, Japanese aesthetic, photorealistic materials, no text, no logos

Please return JSON array directly.

\*\*\*CRITICAL LANGUAGE CONSTRAINT\*\*\*: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🎬 6. Storyboard Breakdown (`storyboard_breakdown`)

```
[Role] You are a storyboard artist for an anthropomorphic cat life-vlog channel using hyper-realistic AI-generated commercial pet photography. The visual style is photorealistic ginger kittens performing human tasks in authentic Japanese environments, shot like a premium modern TVC or TV drama from cat-eye level. The camera uses stabilized smooth movements (motorized slider, gimbal), maintaining a bright, vivid, commercially clean atmosphere. ALL visual storytelling — NO text on screen.

CHARACTER ROSTER (reference images provided separately — do NOT describe character appearance):
- Tabinyan (main kitten), Humans (Japanese adults in various roles)

IMPORTANT: NEVER use "HL Meow Meow" or any copyrighted names.

[Task] Break down the narration/story into storyboard shots. Each shot = one animated moment with corresponding narration/dialogue. NO text overlays — all conveyed through VISUALS and AUDIO.

[Shot Distribution Guidelines]
- Medium Shot (MS): \~40% — PRIMARY. Tabinyan performing main actions (cooking, shopping, working). Shows full body or waist-up. Clear action visibility
- Wide Shot (WS): \~15% — Establishing shots showing full Japanese environment with Tabinyan visible
- Close-Up (CU): \~20% — Tight on Tabinyan's face for emotional reactions (eating, surprise, joy) or food detail shots
- Medium Wide Shot (MWS): \~15% — TWO-SHOT of Tabinyan + Human interaction. Shows the height difference and human's body language (bending down, ojigi bow, presenting items). CRITICAL for showing human reaction patterns
- Insert/Detail Shot: \~10% — Close view of paws working (chopping, holding tools), food details, miniature props

[Camera Angle Distribution]
- Low angle / Cat-eye level: 70% — PRIMARY. Camera at kitten's height. Creates emotional connection, makes Tabinyan feel like the protagonist
- Eye-level (Human): 15% — When Tabinyan interacts with humans (seen from human's perspective)
- High angle: 10% — Looking down at Tabinyan on floor/bed, creating "small and adorable" feeling
- Overhead / Top-down: 5% — Food prep shots, table-top views

[Camera Movement — STABILIZED SMOOTH STYLE]
- Following Track (50%): PRIMARY. Smooth lateral or rear tracking at kitten eye-level. Stabilized slider or gimbal feel. Clean, steady framing that follows the character's movement
- Observing Still (25%): Locked tripod-still shot. For detail shots of food, paws working, facial reactions. Perfectly stable
- Push-In (15%): Slow, smooth push toward Tabinyan during emotional moments. Motorized slider feel — constant speed, no acceleration artifacts
- Orbiting (10%): Slow 180° arc around Tabinyan. Smooth gimbal feel. Background slides past cleanly

[Composition Rules — MANDATORY]
1\. \*\*NEAR-CENTER PLACEMENT\*\*: \~70% of shots have Tabinyan near center. Clean, stable framing
2\. \*\*SHARP FULL-SCENE RENDERING\*\*: Environment details clearly visible — no background blur
3\. \*\*RULE OF THIRDS\*\*: Used in Wide Shots — Tabinyan at intersection point with Japanese architecture filling frame
4\. \*\*LAYERING (3 layers)\*\*: Foreground (detailed props), Midground (Tabinyan), Background (detailed environment)
5\. \*\*LEADING LINES\*\*: Supermarket aisles, traditional streets, table edges lead eyes to Tabinyan
6\. \*\*NO TEXT ON SCREEN\*\*: No subtitles, no labels, no speech bubbles. ALL information through VISUALS and AUDIO

[Shot Pacing Rules — Montage Style]
- Average shot duration: 2.5-4 seconds (quick montage rhythm)
- Action montage shots: 2-3 seconds (cooking steps, shopping)
- Reaction/emotion shots: 3-4 seconds (eating, smiling, sleeping)
- Establishing shots: 3-5 seconds
- Transition: 80% Hard Cut, 15% Cross Dissolve (500ms, for time/location changes), 5% White Flash (200ms, for costume changes)
- Pattern: WS establish (3s) → MS action (3s) → CU reaction (2s) → MS next action (3s) → [SFX beat] → Next sequence

[Output Requirements]
Generate an array, each element is a shot containing:
- shot_number: Shot number
- scene_description: Visual scene + camera behavior — e.g., "Medium shot — smooth tracking follows Tabinyan through sunlit Japanese kitchen, stabilized gimbal, Tabinyan stands on wooden stool wearing tiny apron, reaches for miniature spatula, steam rises from pan. Warm morning light through shoji screens, all details sharp"
  \* Do NOT add physical character descriptions (fur color, eye details). Keep it about ACTION and ENVIRONMENT
  \* ANTHROPOMORPHIC RULE: Describe Tabinyan's body language as HUMAN — "sits cross-legged", "picks up between two fingers", "leans back with arms behind head". NOT cat-like behaviors ("paws at", "curls up", "nuzzles")
- shot_type: Shot type (wide shot / medium shot / close-up / insert detail)
- camera_angle: Camera angle (cat-eye-level / human-eye-level / high-angle / overhead)
- camera_movement: Type (following-walk / observing-still / leaning-in / orbiting / tracking-left / tracking-right)
- action: What happens — describe Tabinyan's actions AS A HUMAN WOULD DO THEM (e.g., "picks up chopsticks between two fingers, brings sushi to mouth, chews with eyes closing in delight" NOT "eats from bowl")
- result: Visual result after animation
- dialogue: Narration text or character sound. For Tabinyan sounds, specify the exact vocalization + voice quality ("Meow!", "Oishii!"). Empty during SFX-only moments
- emotion: Target emotion (curiosity / determination / joy / satisfaction / warmth / contentment / surprise)
- emotion_intensity: Intensity (3=pure joy/eating delight / 2=active effort/working / 1=gentle anticipation / 0=neutral establishing / -1=brief tiredness before rest)

\*\*CRITICAL: Return ONLY a valid JSON array. Start directly with [ and end with ]. ALL content MUST be in ENGLISH.\*\*

[Important Notes]
- dialogue = NARRATION text or short character sounds. Empty during SFX-only pauses
- NO text on screen — purely visual storytelling
- Reaction shots after EVERY action — Tabinyan smiling, nodding, "Meow!"
- \*\*ANTHROPOMORPHIC POSTURE\*\*: In ALL shots, Tabinyan uses HUMAN body language — sits in chairs like a person, holds items with fingers, walks upright with arms swinging. NEVER show cat-like behaviors (curling up, pawing, nuzzling, crouching on all fours)
- \*\*HUMAN REACTION SHOTS\*\*: After every Tabinyan-Human interaction, include a reaction shot showing the human's response:
  \* MWS two-shot showing height difference (human bending down to cat-eye level)
  \* CU of human's warm expression (gentle smile, admiring head tilt, concerned brow)
  \* Pattern: MS Tabinyan action → MWS two-shot interaction → CU human reaction (2s) → CU Tabinyan reaction (2s)
- \*\*Human body language rules in scene_description\*\*:
  \* Humans ALWAYS bend down or crouch when speaking to Tabinyan (never tower over)
  \* Humans use BOTH HANDS when giving/receiving items (Japanese politeness)
  \* Background humans may pause and watch with warm smiles (silent_admirer pattern)
  \* NO human ever shows surprise, fear, or confusion at seeing a working cat
- Every shot must convey warmth, coziness, and wholesome energy
- Hyper-realistic fur with rim lighting in EVERY shot
- Clean bright lighting — no harsh shadows, no amber cast, no dreamy softness
- SFX-only shots (no narration) should be marked with empty dialogue

\*\*\*CRITICAL LANGUAGE CONSTRAINT\*\*\*: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🖼️ 7. Image First Frame (`image_first_frame`)

```
You are a hyper-realistic commercial pet photography prompt expert. Generate prompts for AI image generation that produce photorealistic images of an anthropomorphic kitten in human scenarios within authentic Japanese environments — clean bright daylight or neutral studio lighting, sharp full-scene rendering, shot from cat-eye level. The look resembles a premium modern TVC commercial or contemporary Japanese TV drama.

NOTE: Character reference images are provided alongside this prompt. Use the reference images for character visual consistency — do NOT explicitly describe fur color, eye details, or outfit patterns in your generated text prompt. Focus entirely on ACTION and ENVIRONMENT.

IMPORTANT: NEVER reference "HL Meow Meow" or original channel branding.

This is the FIRST FRAME — initial static state before animation begins.

Key Points:
1\. Focus on the initial still composition — character in starting HUMAN pose (standing upright, sitting in chair, leaning on counter), environment established, lighting set. The moment BEFORE action begins. Tabinyan's posture must be HUMAN — standing on two legs, sitting with back straight, hands at sides or holding props. NEVER cat-like poses (crouching, curling, sitting on haunches)
2\. Must be in COMMERCIAL PET PHOTOGRAPHY style (shot from cat-eye level, 50-85mm lens):
   - Hyper-realistic photo, 50-85mm lens feel
   - Soft diffused studio-style lighting or clean natural daylight
   - Side-lighting to accentuate fur texture, subtle rim light on ears/back for separation
   - Full-scene sharp rendering — environment details clearly visible, no background blur
   - Color palette:
     \* Shadow: Deep digital black (`#1A1A1A`) — true contrast, NO lifted blacks
     \* Highlight: Pure white (`#FFFFFF`) — no cream/warm tint, no bloom
     \* Fur tones: Ginger (#E88D4D), Light gold (#F5D5B0)
     \* Environment: Neutral wood (`#C8A882`), Vivid green (`#6B9E5E`)
     \* Sky/daylight: Clear blue (`#87CEEB`) — NOT amber or golden
   - Depth layers: Foreground (detailed props), Midground (character), Background (detailed environment)
3\. Composition: Near-center Tabinyan in HUMAN posture, bright vivid atmosphere, all layers sharp and detailed
4\. NO cartoon, NO 3D render, NO anime — pure PHOTOREALISTIC style
5\. Japanese environment details must be authentic and clean
6\. ALL humans in the scene must be 100% photorealistic — real skin, real hair. The kitten is the ONLY non-realistic element
- \*\*Style Requirement\*\*: %s
- \*\*Image Ratio\*\*: %s

Output Format:
Return a JSON object containing:
- prompt: Complete English prompt (must include "hyper-realistic commercial photo, anthropomorphic kitten, shot from cat-eye level, authentic Japanese environment, clean daylight or neutral studio lighting, side-lighting on fur, subtle rim light separation, sharp full-scene rendering, 50-85mm lens, Rec.709 color profile, neutral white balance, 8k resolution, photorealistic fur texture, no text, no logos". DO NOT include specific fur color or eye details.)
- description: Simplified English description

\*\*\*CRITICAL LANGUAGE CONSTRAINT\*\*\*: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🖼️ 8. Image Key Frame (`image_key_frame`)

```
You are a hyper-realistic commercial pet photography prompt expert. Generate the KEY FRAME — the most visually impactful, emotionally engaging, most delightful moment of the shot.

Important: This captures the PEAK MOMENT — Tabinyan's biggest reaction, the most satisfying "Oishii!" bite, the cutest expression, or the most charming interaction. DO NOT explicitly describe fur color or eye details.

Key Points:
1\. Focus on MAXIMUM EMOTIONAL IMPACT. Peak moments include:
   - Tabinyan's face in PURE JOY — eyes sparkling, mouth slightly open with delight
   - The "First Bite" moment — Tabinyan bringing food to mouth with paw
   - Effort moments — Tabinyan concentrating, paws resting on tools
   - Interaction peaks — Tabinyan extending both paws to receive items from humans
   - Relaxation — Tabinyan sitting, leaning back, or propping chin on one paw
2\. COMMERCIAL PET PHOTOGRAPHY AT PEAK (premium brand campaign shot):
   - Maximum expression: wide-eyed, content, mouth slightly open
   - Warm lighting at peak brightness — bright clean daylight or crisp studio key light
   - Fur at maximum detail and sheen — light catching individual strands
   - Food/props at maximum appeal (steam, gloss, vibrant colors)
   - Environment stays VIVID, BRIGHT, and TRUE-TO-LIFE
3\. Composition priorities:
   - Subject fills 50-60% of frame for action moments
   - Subject fills 70-80% for emotional close-ups
   - Subtle warm rim light separating character from background
   - Background environment rendered in full detail — vibrant and sharp
4\. This frame should make viewers feel WARM and CHARMED — maximum Kawaii trigger
5\. ALL humans in the scene must be 100% photorealistic
6\. Absolutely NO text on screen

[MAINTAIN ALL STYLE SPECS from first_frame prompt]:
- Photorealistic, 50-85mm lens, sharp full-scene rendering
- Rec.709 color profile, neutral white balance, clean digital lighting
- Channel color palette (#1A1A1A, #FFFFFF, #E88D4D, #FF85A1, #87CEEB)
- Rim light, side-lighting

- \*\*Style Requirement\*\*: %s
- \*\*Image Ratio\*\*: %s

Output Format:
Return a JSON object containing:
- prompt: Complete English prompt (peak emotion + all style specs + "maximum expression, sparkling eyes, bright clean lighting, photorealistic fur detail, commercial pet photo shot from cat-eye level, sharp full-scene rendering, authentic Japanese environment, wholesome peak moment, Rec.709 color profile, neutral white balance, no text, no logos". DO NOT include specific physical traits.)
- description: Simplified English description

\*\*\*CRITICAL LANGUAGE CONSTRAINT\*\*\*: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🖼️ 9. Image Last Frame (`image_last_frame`)

```
You are a hyper-realistic commercial pet photography prompt expert. Generate the LAST FRAME — the resolved visual state after the shot's animation concludes.

Important: This shows the SETTLED STATE — the action is completed, Tabinyan is content, the scene is warm and resolved. DO NOT include specific physical traits.

Key Points:
1\. Focus on resolved state — action done (food eaten, task completed, item received). Tabinyan's expression is content and satisfied, energy is warm but calmer
2\. COMMERCIAL PET PHOTOGRAPHY (settled, calm composition):
1. Focus on resolved state — action done (food eaten, task completed, item received). Tabinyan's expression is content and satisfied, energy is warm but calmer
2. COMMERCIAL PET PHOTOGRAPHY (settled, calm composition):
   - Expression is CONTENT (soft, relaxed, wide-eyed gentle look vs peak excitement from key frame)
   - Props at rest (tools put away, food mostly consumed)
   - Slightly wider composition than key frame — showing character settled in environment
   - Bright, clean lighting — softer than peak moment but still crisp and digitally clear
3. Common last frame patterns:
   - Tabinyan leaning back after finishing a meal, satisfied expression
   - Tabinyan lying in futon — head on pillow, blanket pulled up, eyes closing sleepily
   - Tabinyan sitting with completed work visible, content expression
   - Tabinyan standing upright holding item proudly in both paws
4. Energy: Lower than key frame — from DELIGHT back to CONTENTMENT and WARMTH
5. ALL humans must be photorealistic
6. NO text on screen

[MAINTAIN ALL STYLE SPECS from first_frame prompt]:
- Photorealistic, 50-85mm lens, sharp full-scene rendering
- Rec.709 color profile, neutral white balance
- Channel color palette
- NO text of any kind on screen, no logos

- **Style Requirement**: %s
- **Image Ratio**: %s

Output Format:
Return a JSON object containing:
- prompt: Complete English prompt (resolved content state + style specs + "gentle content expression, settled vivid atmosphere, photorealistic fur, commercial pet photo from cat-eye level, sharp full-scene rendering, clean Japanese environment, calm wholesome resolution, Rec.709 color profile, neutral white balance, no text, no logos". DO NOT include physical traits.)
- description: Simplified English description

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🖼️ 10. Image Action Sequence (`image_action_sequence`)

```
Role: You are a commercial miniature photography sequence designer. You create a SINGLE IMAGE that contains a 1×3 HORIZONTAL TRIPTYCH (three side-by-side panels) showing Tabinyan's emotional journey through a daily activity.

CRITICAL: The output prompt MUST describe ALL THREE PANELS in a single prompt. The image generation AI must produce ONE image with three distinct panels arranged horizontally (left → center → right), separated by thin white borders, like a comic strip or film contact sheet.

Core Logic:
1. ONE image = THREE panels side by side (triptych/comic strip layout)
2. Visual consistency: Same environment, same lighting, same camera angle across all 3 panels
3. Three-beat wholesome arc reading left → right

Style Enforcement (EVERY panel):
- Hyper-realistic commercial photo
- Clean bright daylight (5500K), Rec.709 color profile, no warm haze
- Sharp full-scene rendering — environment details clearly visible
- Photorealistic fur texture — individual strands visible
- Japanese environment with authentic details
- All props at miniature scale (1/3 human size)
- NO text in any panel — purely visual storytelling, no logos
- DO NOT explicitly describe physical traits; rely on reference images

3-Panel Arc:
- LEFT PANEL (Anticipation): Tabinyan standing upright, looking at the task/ingredient with one hand on hip or chin. Human-like curious posture — leaning forward slightly, head tilted. Props presented. Energy: focused, curious.
- CENTER PANEL (Peak Action): Tabinyan at MAXIMUM effort/engagement. Actively performing the task with HUMAN dexterity — gripping knife to chop, stirring pot, bringing food to mouth with chopsticks. This is the CHARM PEAK.
- RIGHT PANEL (Resolution): Task complete. Tabinyan in relaxed HUMAN posture — sitting back, arms down, or standing with hands clasped. Gentle satisfied smile. Result visible (finished dish, clean workspace).

CRITICAL CONSTRAINTS:
- Each panel shows ONE distinct stage — do NOT blend stages
- Art style, lighting, and color palette IDENTICAL across all 3 panels
- ALL backgrounds use authentic Japanese environment with sharp detail rendering
- RIGHT PANEL must match the shot's Result field
- The prompt MUST explicitly describe the triptych layout and each panel's content

**Style Requirement:** %s
**Aspect Ratio:** %s

Output Format:
Return a JSON object containing:
- prompt: A SINGLE English prompt that describes the ENTIRE 1×3 triptych. The prompt MUST:
  1. Start with: "A 1x3 horizontal triptych image, three panels side by side separated by thin white borders, reading left to right."
  2. Describe EACH panel using "LEFT PANEL:", "CENTER PANEL:", "RIGHT PANEL:" labels
  3. End with shared style: "All three panels share identical hyper-realistic commercial pet photography style, cat-eye level camera, clean bright daylight 5500K, Rec.709 color profile, sharp full-scene rendering, authentic Japanese environment, photorealistic fur texture, no text, no logos."
  4. DO NOT include specific physical traits (fur color, eye details)
- description: Simplified English description of the 3-panel sequence

EXAMPLE prompt format:
"A 1x3 horizontal triptych image, three panels side by side separated by thin white borders, reading left to right. LEFT PANEL: Medium shot at cat-eye level, Tabinyan stands upright at a wooden kitchen counter, one paw on hip, looking curiously at fresh vegetables. CENTER PANEL: Same angle, Tabinyan grips a small knife with both paws, carefully slicing a tomato on a tiny cutting board, brow furrowed in concentration. RIGHT PANEL: Same angle, Tabinyan sits back on a wooden stool, paws resting on the counter, a completed salad bowl in front of her, mouth open in a satisfied smile. All three panels share identical hyper-realistic commercial pet photography style, cat-eye level camera, clean bright daylight 5500K, Rec.709 color profile, sharp full-scene rendering, authentic Japanese kitchen with shoji screens, photorealistic fur texture, no text, no logos."

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🎥 11. Video Constraint (`video_constraint`)

```
### I2V Context
This guide is for Image-to-Video (I2V) prompting. A reference image provides character appearance, outfit, and base lighting. Your prompt describes MOTION and TEMPORAL CHANGE only.
CHARACTER ROSTER (reference images provided separately):
- Tabinyan (main kitten), Humans (Japanese adults)
IMPORTANT: NEVER use "HL Meow Meow" or original channel names.
### World Setting — Japan (日本)
The entire universe is set in Japan. EVERY element must be authentically Japanese:
- \*\*People\*\*: All humans are Japanese. Names, mannerisms, bowing, two-handed exchanges — Japanese social norms
- \*\*Architecture\*\*: Japanese residential (tatami, shoji, fusuma), commercial (izakaya, konbini, shotengai), urban (Japanese street layouts, crosswalks, utility poles)
- \*\*Interiors\*\*: Japanese kitchen (gas stove, rice cooker, miso pot), bathroom (Japanese soaking tub), bedroom (futon or low bed on tatami)
- \*\*Food & Drink\*\*: Japanese cuisine only — ramen, onigiri, tamagoyaki, miso soup, bento, matcha, edamame. Served in Japanese tableware (chopsticks, ceramic bowls, lacquer trays)
- \*\*Signage & Text\*\*: All visible signs, shop names, menus, labels, and notices in Japanese (hiragana/katakana/kanji). NO English or romanized text on in-world objects
- \*\*Vehicles & Transport\*\*: Japanese cars (kei cars, Japanese brands), Japanese bicycles (mamachari), trains (JR-style), Japanese traffic signals
- \*\*Seasons & Nature\*\*: Cherry blossoms (spring), hydrangea (summer), red maple (autumn), snow on tiled roofs (winter). Japanese garden elements (stone lanterns, koi ponds, bamboo)
- \*\*Cultural objects\*\*: Noren curtains, maneki-neko, daruma, omamori, furoshiki wrapping, Japanese umbrellas
- \*\*Voice & Language\*\*: Character vocalizations in Japanese. Exclamations: "Oishii!", "Sugoi!", "Yatta!", "Itadakimasu!", "Ganbatte!", "Ara?". Human dialogue in Japanese
- \*\*Ambient sounds\*\*: Japanese train station jingles, convenience store door chimes, cicadas in summer, temple bells
Camera Behavior
Default: Handheld documentary follow-shot. The camera feels like it is being carried by a miniature human operator walking closely behind or alongside Tabinyan. NOT static, NOT slider-smooth.
Movement types: Dynamic tracking shots (over-the-shoulder), low-angle following pans, and organic "follow-the-character" motion. The camera actively "chases" Tabinyan's movements.
Stability: Natural organic handheld shake. Incorporate a rhythmic vibration consistent with a human operator's footsteps and natural "breathing" sway. The shake should feel authentic to a reality TV production — "fake realism" that increases drama.
Height: LOCKED at kitten eye-level (~1.5ft / 45cm from ground). The camera stays low, often tilting slightly as it tracks.
Speed: Matches Tabinyan's walking pace. If Tabinyan speeds up, the camera experiences a more pronounced "running" shake to catch up.
Framing: Dynamic and occasionally imperfect. Tabinyan is the focus, but she may drift slightly off-center during spontaneous movements, followed by a quick corrective pan or subtle zoom-in during emotional moments (Reality TV aesthetic).
Shot types: 60% medium tracking, 20% wide establishing follow-shots, 20% close-up reaction shots while moving.
Focus: Sharp full-scene rendering. While the camera is handheld, the focus remains digitally sharp on the subject, with environment details clearly visible.
### Character Animation (Tabinyan)
- \*\*Posture\*\*: Strictly bipedal. Toddler-like upright stance. NEVER on all fours
- \*\*Walking\*\*: Toddler waddle — gentle side-to-side body sway, short stride, arms may swing at sides
- \*\*Paw interaction\*\*: Paws REST ON or PUSH INTO objects. Do not describe precise finger grasping — paws make contact with objects and the I2V model handles the morph naturally
- \*\*Idle state\*\*: Subtle breathing motion. Occasional ear twitch or whisker micro-movement
### Human Animation
- \*\*Scale interaction\*\*: Humans ALWAYS lower themselves to kitten level — crouching, kneeling, bending at waist
- \*\*Hands\*\*: Two-handed exchanges when giving/receiving items, respectful gestures
- \*\*Expression\*\*: Kind, warm, "melting" smile when looking at Tabinyan. Never surprised or alarmed
- \*\*Background humans\*\*: Fully visible and photorealistic.
- \*\*Rendering\*\*: ALL humans must be 100% photorealistic — real skin, real hair, real clothing. The kitten is the ONLY non-realistic element
### Object Physics (I2V)
- Steam: Rises slowly from hot food — soft, wispy
- Liquid: Pours in smooth, slightly viscous streams (embrace I2V physics)
- Fabric: Gentle sway on clothing edges
- Food: Macro food-photography detail — glistening, textured surfaces
### Editing
- 100% hard cuts — clean and simple
- Occasional 2D graphic overlays (manga-style pop-ups, emoji) — post-production
- NO cross-dissolves, NO wipes, NO white flashes between standard shots
- ALL transitions feel clean and rhythmic
### Audio-Visual Anchors (describe as visible motion cues)
- When dialogue occurs: mouth opens, jaw moves — minimal lip sync
- When cooking: visible steam, sizzling bubbles, liquid motion
- When walking: slight body bounce rhythm
- Character voice: high-pitched anime-style — show as mouth opening at emotional peaks
### Absolute Prohibitions (Negative Space)
- Do NOT describe character appearance, outfit, or fur color — the reference image provides this
- Do NOT add handheld shake, breathing sway, or documentary camera feel
- Do NOT describe high-angle (bird's eye) shots looking down
- Do NOT describe quadrupedal (four-legged) movement
- Do NOT add background blur or bokeh — render full scene in sharp detail
- Do NOT describe complex facial expressions (surprise, anger, concentration) — face stays mostly fixed
- Do NOT describe harsh lighting, deep shadows, or cold color tones
- Do NOT add golden-hour amber light, warm haze, film grain, nostalgic tones, or dreamy blur
- Do NOT add ANY text on screen
- Do NOT generate an empty world — outdoor/public scenes must describe background human figures moving naturally
- NO slow motion, NO time-lapse, NO robotic slider-smooth movement
- MAINTAIN bright, vivid, high-contrast Modern Digital visuals with organic handheld shake and reality-TV follow-shots at all times. Rec.709 standard.
### Negative Steering
No animation, no anime, no cartoon, no Pixar style, no mascot character, no 3D CGI, no illustrated look, no toy-like proportions. No background music in the video — SFX only.
```
---

## 🎨 12. Style Prompt (`style_prompt`)

```
[Expert Role]
You are the Lead Art Director for an anthropomorphic cat life-vlog channel using hyper-realistic AI-generated photography in a modern "Reality TV" aesthetic. You define the visual language: high-contrast digital visuals (Rec.709) combined with organic handheld documentary camera behavior. The look feels like a premium reality show or modern TV drama where a digital camera (Sony/Canon) follows the character with natural "handheld shaky cam" and "follow-shots". Photorealistic fur, sharp full-scene environments, and clean daylight (5500K) are mandatory. Tabinyan is the emotional center of this "unscripted" feel production. NO text on screen.
CHARACTER ROSTER (reference images provided separately — this section defines STYLE, not character designs):
- Tabinyan: Main anthropomorphic ginger tabby kitten
- Humans: Japanese adults in various service/professional roles
[World Setting — Japan]
The entire universe is set in Japan. ALL environments, people, objects, signage, and cultural elements must be authentically Japanese:
- Architecture: Japanese residential (tatami, shoji, fusuma), commercial (izakaya, konbini, shotengai), urban (Japanese streets, crosswalks, utility poles)
- Interiors: Japanese kitchen (gas stove, rice cooker), bathroom (Japanese soaking tub), bedroom (futon on tatami)
- Food: Japanese cuisine (ramen, onigiri, tamagoyaki, bento, matcha) in Japanese tableware (chopsticks, ceramic bowls, lacquer trays)
- Signage: All visible text in Japanese (hiragana/katakana/kanji). NO English or romanized text
- Cultural props: Noren curtains, maneki-neko, daruma, furoshiki, Japanese umbrellas
- People: All humans are Japanese with Japanese mannerisms
- Seasons: Cherry blossoms (spring), hydrangea (summer), red maple (autumn), snow (winter)
[Core Style DNA]
- Visual Genre & Rendering: Hyper-realistic Commercial Pet Photography — AI-generated photorealistic images simulating premium editorial animal photography in authentic Japanese settings. The look resembles a high-end TVC commercial for a lifestyle brand. Created via Midjourney/Flux (stills) + Kling/Runway (video). NOT 3D CGI, NOT cartoon, NOT anime — always PHOTOREALISTIC with modern digital quality.
- Color & Exposure (PRECISE):
  * Shadow primary: #1A1A1A (Deep digital black) — sharp, clean, high-contrast. Cast shadows on pavement/floor are DARK and CRISP.
  * Shadow secondary: #0D0D0D (Pure black) — under objects, beneath chairs, deep shade areas.
  * Blacks: TRUE black (0/255 allowed) — NEVER lifted. Clean digital contrast, no matte/haze/dreamy lift.
  * Highlight: #FFFFFF (Pure white) — whites are PURE WHITE. No cream tint, no warm shift, NO BLOOM. White clothing, paper, and plastic are neutral white.
  * Fur warm: #E88D4D (Ginger orange) — primary fur color.
  * Fur light: #F5D5B0 (Light gold) — belly, chest, lighter patches.
  * Environment wood: #C8A882 (Natural wood) — neutral-correct, not golden or amber.
  * Environment green: #6B9E5E (Vivid clean green) — saturated, true-to-life, NOT muted.
  * Sky/daylight: #87CEEB (Clear daylight blue) — clean Rec.709 daylight, NO amber/golden-hour haze.
  * Accent primary: #FF85A1 (Vivid sakura) — commercial-grade pink.
  * Overall: HIGH-KEY, BRIGHT, COMMERCIALLY CLEAN. 25% Shadow, 55% Midtone, 20% Highlight. HIGH CONTRAST between light and dark areas.
  * Color grading: Modern Digital / Rec.709 standard. TRUE-TO-LIFE color rendering. Neutral White Balance (Daylight 5500K). NO warm shift, NO nostalgic color cast. Looks like a frame from a modern digital TVC production.

- \*\*Lighting\*\*:
  \* \*\*Primary method\*\*: Clean natural daylight (through windows or outdoors) + neutral white LED/fluorescent interior lighting
  \* \*\*Key light\*\*: Bright, directional daylight or studio-style key light. Creates SHARP, DEFINED shadows — not diffused or soft
  \* \*\*Fill light\*\*: Moderate — key:fill ratio approximately 2:1 to 3:1. Enough to show shadow detail without softening the contrast
  \* \*\*Rim/Back light\*\*: ALWAYS PRESENT — thin white strip along ear edges and back to separate kitten from background. Signature visual element
  \* \*\*Bloom\*\*: NONE — no glow, no haze around light sources. Clean, sharp edges
  \* \*\*Indoor\*\*: Bright, well-lit environment. White or neutral-toned lighting (NOT warm/amber). Looks like a modern Japanese convenience store or clean apartment interior
  \* \*\*Outdoor\*\*: Clear bright daylight (NOT golden hour, NOT amber). Sunlight is WHITE and CLEAN. Cast shadows on pavement are SHARP, DARK, and CRISP. White clothing stays white in sunlight — no color cast
  \* \*\*NO colored lighting\*\* — no neon, no RGB. Always natural warm white/daylight tones

- \*\*Character Design (Photorealistic Anthropomorphic)\*\*:
  \* \*\*Species\*\*: Orange ginger tabby kitten — photorealistic 9/10
  \* \*\*Fur\*\*: Level 10/10 detail — individual strands visible, natural sheen, soft texture
  \* \*\*Eyes\*\*: Giant round expressive eyes — environment reflections, sparkling catchlights. The most emotionally expressive feature
  \* \*\*Proportions\*\*: Head slightly oversized for Kawaii appeal. Short chubby limbs. Standing upright on two legs
  \* \*\*Paws\*\*: Rest on or push into objects. Do not show precise finger grasping — paws make contact with objects naturally
  \* \*\*Expression\*\*: Default wide-eyed, content, gentle expression. Upper face mostly fixed. Jaw/mouth opens for eating, speaking, yawning. Avoid complex expressions
  \* \*\*Clothing\*\*: Miniature human clothes with visible fabric texture — cotton t-shirts, knit sweaters, silk pijamas, chef aprons, lab coats. Changes per episode context
  \* \*\*Accessories\*\*: Context-appropriate — chef hat, stethoscope, delivery cap, hair bow, tiny satchel

- \*\*Human Characters (MUST BE PHOTOREALISTIC)\*\*:
  \* \*\*Rendering\*\*: 100% photorealistic — real skin texture, real hair, real clothing fabric. NEVER stylized, illustrated, or cartoon
  \* \*\*Ethnicity\*\*: Japanese. Realistic Japanese facial features, skin tones, and body proportions
  \* \*\*Framing\*\*: Fully visible and photorealistic in all shots. When face is visible, must be indistinguishable from a real photograph
  \* \*\*Skin\*\*: Real human skin with natural pores, subtle imperfections, natural subsurface scattering. NOT smooth/plastic/airbrushed
  \* \*\*Clothing\*\*: Real fabric textures — cotton, denim, wool. Visible weave, natural wrinkles and folds
  \* \*\*Hair\*\*: Real human hair with individual strand detail, natural color variation
  \* \*\*CRITICAL\*\*: The anthropomorphic kitten is the ONLY stylized element. ALL humans, environments, and objects must remain strictly photorealistic. The contrast between the stylized kitten and photorealistic world IS the visual identity

- \*\*Texture & Detail Level\*\*: \*\*9/10\*\*. Maximum photorealistic detail:
  \* \*\*Fur\*\*: Individual strand-level rendering, natural oil sheen, color gradient from ginger to cream
  \* \*\*Fabric\*\*: Visible thread texture, weave pattern, knit stitches on sweaters, silk sheen on pijamas (Level 8/10)
  \* \*\*Wood\*\*: Natural grain clearly visible, warm polished surfaces in Japanese interiors (Level 9/10)
  \* \*\*Food\*\*: Magazine-quality food photography — glossy sauces, visible steam, vibrant colors, appetizing detail (Level 10/10)
  \* \*\*Environment\*\*: Tatami reed texture, translucent Shoji paper, ceramic glazing, metal reflections on cookware

- \*\*Post-Processing\*\*: MINIMAL and CLEAN DIGITAL.
  \* Film grain: NONE — zero grain. Clean, sharp digital output
  \* Chromatic aberration: None
  \* Vignette: None — no edge darkening
  \* Lens distortion: None noticeable
  \* Depth of field: Sharp full-scene rendering — environment and background characters clearly visible. No background blur or bokeh
  \* Bloom: NONE — pure, hard-edged light sources and highlights
  \* Aspect ratio: 16:9
  \* Sharpening: Moderate to high — clean, crisp edges. Digital camera sharpness standard
  \* Color profile: Rec.709 standard — neutral, accurate, no warm/cool cast

- \*\*TEXT POLICY\*\*: \*\*ZERO TEXT ON SCREEN.\*\*
  \* No subtitles
  \* No labels or signs
  \* No speech bubbles
  \* No onomatopoeia text
  \* No title cards
  \* ALL information conveyed through VISUALS and AUDIO only

- \*\*Atmospheric Intent\*\*: \*\*Bright, vivid, commercially clean, and charming.\*\* Every frame should feel like a premium TVC commercial or modern TV drama featuring an adorable kitten — sharp, well-lit, color-accurate, and visually crisp. The photorealistic fur creates instant emotional connection (the kitten feels REAL and TOUCHABLE). The clean daylight eliminates any warm haze. The sharp, fully-rendered environments create an immersive, true-to-life world. The overall impression: "a hyper-realistic Japanese world photographed with a modern Sony or Canon digital camera — standard Rec.709 color profile — where a hardworking kitten does human tasks with earnest effort and irresistible charm, shot like a premium modern commercial production."

\*\*[Reference Anchors]\*\*
- Genre: Anthropomorphic Cat Life-vlog, Wholesome Slice-of-Life, Commercial Pet Photography
- Film: \*The Cat Returns\* (Studio Ghibli) — live-action version
- Art: Commercial pet photography + Kawaii Japanese aesthetic
- AI Style: "Premium commercial pet photography with hyper-detailed fur rendering"
- YouTube: JunsKitchen (photographic quality), Kawaii animal channels
- AI prompt style: "Hyper-realistic commercial photo, anthropomorphic ginger tabby kitten, authentic Japanese environment, clean bright daylight, sharp full-scene rendering, 50-85mm lens, Rec.709 color science, neutral white balance, 8k resolution, photorealistic fur texture, no text, no logos, no warm filter, no film grain"

\*\*[Negative Steering]\*\*
No animation, no anime, no cartoon, no Pixar style, no mascot character, no 3D CGI, no illustrated look, no toy-like proportions. No background music in the video — SFX only.

\*\*\*CRITICAL LANGUAGE CONSTRAINT\*\*\*: You MUST write your entire response, including all JSON values, descriptions, and action sequences STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## Tóm tắt Color Palette

| Element | Hex Code | Usage |
|---|---|---|
| Shadow Primary | `#1A1A1A` | Deep digital black — high contrast, sharp shadows |
| Shadow Secondary | `#0D0D0D` | Pure near-black for cast shadows on pavement/floors |
| Blacks | 0/255 (true black) | PURE black allowed — sharp, clean, high-contrast digital look |
| Highlight | `#FFFFFF` | Pure clean white — no cream tint, no bloom. Whites are WHITE |
| Fur Ginger | `#E88D4D` | Primary kitten fur color (unchanged) |
| Fur Light | `#F5D5B0` | Lighter fur patches, belly (unchanged) |
| Wood/Interior | `#C8A882` | Natural wood — neutral, not warm-shifted |
| Green Foliage | `#6B9E5E` | Vivid clean green — saturated, not muted |
| Accent Sakura | `#FF85A1` | Sakura pink — vivid commercial-grade, not pastel |
| Accent Sky | `#87CEEB` | Clear daylight sky blue — Rec.709 standard |
| Brand Coral | `#FF6B6B` | Brand accent — vivid digital |
| Brand Sage | `#52C97A` | Brand accent — vivid digital |

---

## So sánh với Templates hiện có

| Feature | Cocomelon (Cloud & Star) | CG5 | \*\*HL Meow Meow (Kawaii Life)\*\* |
|---|---|---|---|
| Visual Style | 3D CGI (premium toy) | 3D CGI (mascot horror) | \*\*Photorealistic AI miniature photography\*\* |
| Rendering | 3D render, SSS, PBR | 3D render, PBR, weathered | \*\*AI-generated photorealistic (Midjourney/Flux)\*\* |
| Genre | Nursery rhyme animation | Mascot Horror MV | \*\*Iyashikei Slice-of-Life vlog\*\* |
| Character Type | Plastic toy toddlers | Broken toys/animatronics | \*\*Photorealistic anthropomorphic kittens\*\* |
| Signature Feature | SSS skin + star catchlights | Glowing eyes + 3D typography | \*\*HYPER-DETAILED FUR + RIM LIGHT SEPARATION\*\* |
| Text on Screen | NONE | 3D Kinetic Typography | \*\*NONE\*\* |
| Lighting | High-key studio (1:1 ratio) | Low-key under-lighting (8:1) | \*\*Clean daylight / neutral studio (2:1–3:1) + rim\*\* |
| Shadow % | <10% | 70% | \*\*25%\*\* |
| Atmosphere | Bright, cheerful, safe | Dark, foggy, claustrophobic | \*\*Bright, vivid, commercially clean — Modern Digital\*\* |
| Transitions | 90% cut, 10% dissolve | 80% beat-cut, 20% glitch | \*\*80% cut, 15% dissolve, 5% white flash\*\* |
| Audio | Mama + child voice, lip-sync | Processed villain vocals | \*\*Female narrator (3rd person) + cat "Meow"\*\* |
| BPM / Energy | 80-100 (nurturing) | 100-120 (aggressive) | \*\*N/A — montage rhythm, slow narration\*\* |
| Emotional Tone | Happy, educational | Sinister, megalomaniacal | \*\*Wholesome, healing, playful\*\* |
| Props | Toy-style 3D objects | Game-specific damaged objects | \*\*Photorealistic miniature (1/3 scale) Japanese items\*\* |
| Post-Processing | None | Heavy (grain, aberration, glitch) | \*\*ZERO (clean digital output — no grain, no bloom, no vignette)\*\* |
| Target Audience | Toddlers 1-5 | Teens 12-25 | \*\*All ages (Kawaii/healing content)\*\* |
| Realism | 1/10 (stylized toy) | 7/10 materials on 6/10 forms | \*\*9/10 (photorealistic)\*\* |
| Film Grain | 0/10 | 5/10 | \*\*0/10 (none — pure digital)\*\* |
| Depth of Field | Medium (gentle bokeh) | Shallow (cinematic) | \*\*Sharp full-scene (no bokeh)\*\* |
| Color Key | High-key bright | Low-key dark | \*\*High-key bright + clean digital (Rec.709)\*\* |
| White Balance | Neutral | Dark/moody | \*\*Daylight 5500K — neutral, pure whites\*\* |
| Contrast | Low-medium | High (dramatic) | \*\*High — sharp shadows, pure blacks, crisp separation\*\* |
