# Sprunki Characters — Hyper-Realistic Reference Image Prompts (Adult Version)

Chuyển ảnh 2D game art thành ảnh tham chiếu hyper-realistic **ADULT proportions**.

> **Input**: Ảnh 2D gốc nhân vật (làm image reference)
> **Output**: Character sheet hyper-realistic giữ đúng design gốc, nhưng với **adult body proportions**

> [!IMPORTANT]
> Reference image CỰC KỲ QUAN TRỌNG cho việc thể hiện "adult" trong video output.
> Nếu reference image trông childish → video sẽ render childish. 
> Phải đảm bảo reference có adult proportions + adult context props (coffee, laptop).

---

## Universal Prompt (dùng cho TẤT CẢ nhân vật — ADULT VERSION)

Upload ảnh 2D + paste prompt này:

```
Hyper-realistic character reference sheet. Convert this 2D character design into a photorealistic living creature with ADULT PROPORTIONS — taller, longer limbs, mature posture. Maintain the EXACT original design colors, outfit details, and distinguishing features from the reference image. Do not alter the character's identity or color scheme.

CRITICAL ADULT REFRAMING:
- Body proportions: adult human ratio (7-8 head heights), NOT child proportions (4-5 head heights)
- Posture: relaxed adult stance — weight on one hip, one hand in pocket, confident
- Height relative to environment: can reach kitchen counter, sit on bar stool comfortably, lean on car hood
- Expression: mature, calm, self-assured — NOT wide-eyed childish wonder

Style conversion rules:
- Skin/body surface: real organic texture with visible pores, subtle subsurface scattering, natural sheen — like a real living creature
- Eyes: glossy wet corneas with realistic light reflections and depth. Expression: calm, alert, adult awareness
- Clothing/accessories: real fabric and material textures — visible stitching, thread weave, natural wrinkles, proper draping, realistic zippers/buttons
- Any appendages (antennae, horns, tails): real biological tissue with organic curves and subtle vein-like details

Format: Full-body T-pose, character turnaround sheet showing front view, 3/4 view, and side view. Clean white background. Studio photography lighting: bright key light from upper-left, soft fill from right, rim light for subject separation. Shot on 85mm lens, 5500K daylight, Rec.709 color profile, sharp rendering. 2:3 portrait ratio. No text, no logos, no watermarks, no background elements.

Negative: No 3D CGI render, no Pixar style, no cartoon, no illustration, no anime, no toy-like look, no child proportions, no cute/kawaii aesthetic. This must look like a REAL adult creature photographed in a studio.
```

### Thêm dòng outfit theo nhân vật

Thêm **1 dòng** dưới đây vào cuối prompt (trước "Negative:"), tùy nhân vật:

| Nhân vật | Dòng thêm |
|---|---|
| **Oren** | `Outfit: blue zip-up jacket with letter "O" on chest, dark jeans, sneakers, light orange over-ear headphones around neck or on head. Adult male proportions, relaxed confident stance.` |
| **Pinki** | `Outfit: casual chic — pink and white top, comfortable jeans, white sneakers, small earrings. Adult female proportions, warm confident posture.` |
| **Simon** | `Outfit: green sporty hoodie, joggers, high-top sneakers, smartwatch on wrist. Adult male proportions, energetic forward-leaning stance.` |
| **Durple** | `Outfit: purple knit cardigan over button-up shirt, dark slacks, loafers. Adult male proportions, relaxed unhurried posture, one hand in pocket.` |
| **Gray** | `Outfit: plain gray oversized hoodie, dark joggers, plain white sneakers, wired earbuds around neck. Adult male proportions, quiet reserved stance.` |
| **Sky** | `Outfit: light blue casual button-up, khaki chinos, clean white canvas shoes, small backpack. Young adult male proportions, open curious posture.` |

## Biến thể

### A. Adult Context Shot (cho video reference — KHUYẾN KHÍCH)

Thay vì chỉ T-pose, tạo reference trong **adult environment**:

```
Hyper-realistic photo. Convert this 2D character into a photorealistic living creature with ADULT proportions. Maintain EXACT original design, colors, and outfit from reference.

The character is standing in a modern coffee shop, leaning casually against the counter, holding a coffee cup in one hand, phone in the other. Relaxed adult posture, weight on one hip. Ambient cafe background softly blurred.

Adult body proportions (7-8 head heights). Organic skin, glossy eyes, real fabric textures. Professional photography. 85mm, 5500K. No text, no overlays.

Negative: No child proportions, no cartoon, no CGI, no kawaii aesthetic.
```

### B. Expression Sheet (4 biểu cảm — ADULT)

```
Hyper-realistic expression sheet. Convert this 2D character into a photorealistic adult creature. Maintain EXACT original design from reference.

2x2 grid showing 4 ADULT expressions: neutral (calm, self-assured), amused (subtle smirk, one eyebrow raised), tired (half-lidded eyes, slight slouch), surprised (controlled, raised eyebrows — NOT cartoonish shock). Each expression is a head-and-shoulders close-up.

Adult proportions. Organic skin, glossy eyes, real textures. Clean white background. Studio lighting. No text.

Negative: No child expressions, no exaggerated cartoon reactions, no kawaii.
```

### C. Outfit Variant (đổi trang phục)

```
Hyper-realistic character reference. Convert this 2D character into a photorealistic adult creature. Maintain EXACT original body/face design from reference.

OUTFIT CHANGE: [mô tả trang phục mới — ví dụ: "wearing a bartender's black vest over white rolled-up sleeves, bar apron, name tag that reads 'OREN'"]

Full-body T-pose, front and 3/4 view. Adult proportions. Clean white background. Studio lighting. 85mm, 5500K, Rec.709. No text.

Negative: No child proportions, no cartoon, no CGI.
```

### D. Work/Gig Outfit Variants

| Gig | Outfit Description |
|---|---|
| Bartending | Black vest over rolled-up white shirt, bar apron, name tag |
| Food delivery | Branded cap, delivery bag slung over shoulder, bike helmet dangling |
| DJ set | Headphones ON (full coverage), casual dark outfit, behind decks |
| Pet-sitting | Casual clothes + leash in hand, slightly overwhelmed expression |
| Moving helper | Plain t-shirt, work gloves, slightly dusty |
| Office temp | Button-up (slightly too big), lanyard with temp badge |

---

## Checklist trước khi generate

- [ ] Ảnh 2D gốc đã upload làm reference
- [ ] Prompt KHÔNG mô tả ngoại hình (màu da, kiểu antenna, kiểu mắt...)
- [ ] Chỉ mô tả: style (hyper-realistic) + format (T-pose, turnaround, white BG)
- [ ] **ADULT proportions specified** (7-8 head heights, NOT child/chibi)
- [ ] **Negative steering includes**: "no child proportions, no kawaii, no cartoon"
- [ ] Có negative steering (no CGI, no cartoon, no Pixar)
