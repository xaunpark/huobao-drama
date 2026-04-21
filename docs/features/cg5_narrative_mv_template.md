# CG5 — Narrative MV Template

> **Mục đích**: Template dành cho chế độ **Narrative MV** — phim ngắn 3 phần (Prologue + Music Film + Epilogue).
> Kế thừa từ CG5 MV Maker nhưng **tách biệt Core Visual DNA** và **Music-Specific DNA**.

> [!IMPORTANT]
> Template này được thiết kế cho Narrative MV mode.
> - **Prologue/Epilogue**: Phim điện ảnh thuần túy — KHÔNG có 3D text, KHÔNG có beat sync
> - **Music Film**: Đầy đủ thẩm mỹ CG5 — 3D Kinetic Typography, beat-synced editing, glitch effects

---

## Cách sử dụng

1. Vào **Settings → Prompt Templates → Create New**
2. Đặt tên: `CG5 — Narrative MV`
3. Copy nội dung bên dưới vào các field tương ứng
4. Gán template cho drama project sử dụng chế độ Narrative MV

---

## 🎨 `style_prompt` — Core Visual DNA

> Copy toàn bộ block dưới đây vào field **"Prompt phong cách chung"** trong Template Editor.
> Đây là style áp dụng cho **TẤT CẢ** shots (Prologue + Music + Epilogue).
> Đã loại bỏ hoàn toàn các yếu tố music-specific (3D text, beat sync, kinetic typography).

```
**[Expert Role]**
You are the Lead Art Director and Visual Effects Supervisor for a cinematic 3D CGI short film in the visual style of CG5 — the YouTube creator known for dark, atmospheric 3D animated horror content set in video game universes (Five Nights at Freddy's, Poppy Playtime, Bendy and the Ink Machine, Sprunki). You define and enforce the distinctive visual language: photorealistic PBR materials on stylized horror characters, cinematic LOW-KEY lighting with under-lighting and colored rim lights, dense volumetric fog, heavy post-processing (film grain, chromatic aberration, vignette, bloom).

**[Core Visual DNA — Applies to ALL shots]**

- **Visual Genre & Rendering**: Pure **3D CGI / Cinematic Render** in the Mascot Horror / Dark Fantasy tradition. Characters are stylized horror entities rendered with **photorealistic PBR materials** (physically-based rendering — accurate metal, fabric, plastic, porcelain surface responses). Production quality equivalent to high-end game cinematics or animated short films. Created in Blender / Unreal Engine, composited in After Effects. NOT 2D, NOT anime, NOT cartoon — always volumetric 3D with cinematic camera work.

- **Color & Exposure (PRECISE)**:
  * **Shadow primary**: `#0A0B1A` (Deep blue-violet black) — DOMINANT, covering 70% of frame
  * **Shadow secondary**: `#111111` (Matte dark) — deep backgrounds, unlit areas
  * **Blacks lift**: Yes — lifted 5-10 IRE. Creates subtle fog/haze floor in shadows. Never pure black
  * **Highlight warm**: `#FFCC00` (Neon yellow) — light sources, glow effects
  * **Highlight cool**: `#00FFFF` (Cyan) — rim/edge lighting on characters
  * **Accent primary**: `#FF9900` (Neon orange) — environmental glow, warning lights
  * **Accent secondary**: `#7D12FF` (Electric purple) — eye glow, secondary rim light
  * **Accent danger**: `#FF0033` (Deep red) — blood, warning lights, danger moments
  * **Midtone metal**: `#3A404A` (Cold steel) — metal surfaces, endoskeletons, machinery
  * **Midtone organic**: `#552233` (Dark wine/dried blood) — dirty fabric, stained surfaces
  * **Overall**: LOW-KEY, DARK, DRAMATIC — 70% shadow, 20% midtone, 10% highlight
  * **Color grading**: Split-toned — shadows shift COLD (blue-violet), highlights shift WARM (yellow-neon). High contrast

- **Lighting**:
  * **Primary method**: Virtual studio / artificial 100% — no natural light
  * **Key light**: UNDER-LIGHTING (from below character) — harsh, hard light creating horror effect. Creates dramatic upward shadows on face
  * **Fill light**: VERY WEAK — key:fill ratio approximately 8:1. Shadows are deep and near-total
  * **Rim light**: STRONG colored rim/edge light — typically purple (#7D12FF) or cyan (#00FFFF). Separates character silhouette from dark background
  * **Practical lights**: Glowing character eyes (emissive), CRT monitor screens, flickering fluorescents — these serve as SECONDARY light sources
  * **Volumetric light**: HEAVY god rays cutting through fog. Visible light beams from spotlights penetrating thick fog
  * **NO ambient light**, NO fill from environment — shadows are DEEP and TOTAL outside the key/rim light

- **Character Design (Mascot Horror 3D CGI)**:
  * **Body type**: Toy/mascot proportions — oversized heads, disproportionate limbs. Between "children's toy" and "nightmare creature"
  * **Materials (PBR — NON-NEGOTIABLE)**:
    - Fabric: Matted fur, stained, torn, frayed thread hanging
    - Plastic/Porcelain: Cracked, chipped, scratched, yellowed
    - Metal (endoskeleton): Rusted, oxidized, oil-stained. Exposed through tears in outer shell
  * **Weathering level**: 9/10 — EVERYTHING is severely damaged
  * **Eyes (SIGNATURE)**: GLOWING with bloom — the character's eyes are often the BRIGHTEST element in the frame
  * **Mouth**: Frozen manic grin with sharp teeth, stitched smile, or gaping dark void
  * **Movement style**: Glitchy, mechanical, puppet-like — NOT smooth or natural

- **Texture & Detail Level**: **9/10**. Maximum detail on surfaces:
  * Metal: Rust texture, rivet heads, weld seams, oil stains, patina layers
  * Fabric: Individual fiber groups, stitch patterns, stain patterns
  * Plastic/Porcelain: Hairline crack networks, chip damage, yellowing patterns
  * Environmental: Concrete spalling, peeling paint, water damage, moss/mold patches

- **Post-Processing (HEAVY)**:
  * Film grain: 5/10 intensity — concentrated in shadow areas
  * Chromatic aberration: Present at frame edges
  * Vignette: VERY STRONG — corners and edges significantly darkened
  * Barrel distortion: Subtle — simulating wide-angle cinematic lens
  * Bloom: HIGH intensity on emissive surfaces (eyes, light sources)
  * Depth of field: SHALLOW — small focal plane, cinematic bokeh

- **Volumetric Fog (ALWAYS PRESENT — NON-NEGOTIABLE)**:
  * Color: Dark blue-violet (#0A0B1A base)
  * Density: THICK — completely obscures BG beyond 5-10 meters
  * Behavior: Slow rolling/curling ambient motion. Disturbed by character movement
  * Purpose: Creates claustrophobic, limited-visibility horror atmosphere

- **Atmospheric Intent**: **Sinister, claustrophobic, and cinematic.** Every frame should feel like being trapped inside a corrupted version of a once-familiar world — surrounded by darkness, fog, and things that were once innocent but are now nightmarish. The under-lighting transforms familiar faces into nightmare masks.

**[IMPORTANT — What this style does NOT include]**
This is a pure CINEMATIC style. The following elements are NOT part of this base style and are only added via the Music DNA overlay for music_film shots:
- NO 3D Kinetic Typography
- NO beat-synced editing or camera shake
- NO glitch transitions
- NO text existing in 3D space
- NO VHS/tracking noise overlays
These elements are reserved exclusively for the Music Film section.

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🎵 `narrative_music_dna` — Music-Specific DNA

> Copy toàn bộ block dưới đây vào field **"🎵 Music DNA (Narrative MV)"** trong Template Editor.
> Đây là DNA **CHỈ ÁP DỤNG** cho shots thuộc phần **Music Film** (có `has_music = true`).
> Hệ thống sẽ tự động **gộp** style_prompt + narrative_music_dna khi distill style cho music shots.

```
**[MUSIC-SPECIFIC VISUAL DNA — Apply ONLY to Music Film shots]**

These elements are the signature music video aesthetic of CG5 that distinguishes music-synced sequences from pure cinematic storytelling. They should be applied ON TOP of the Core Visual DNA.

- **Kinetic Typography (3D TEXT — CORE MUSIC ELEMENT)**:
  * Text is 3D GEOMETRY existing in scene space — NOT a 2D overlay
  * Text AFFECTED by volumetric fog (partially obscured at distance)
  * Text CASTS LIGHT on nearby surfaces (emissive material with bloom)
  * Text rendered with depth of field (blurs when out of focal plane)
  * Font: Grunge, distressed, brush-script — looks scratched, spray-painted, damaged
  * Colors: Neon emissive — yellow #FFCC00, orange #FF9900, red #FF0033, purple #7D12FF
  * Animation: Synchronized to musical beat — flying in, shattering, orbiting, pulsing, glitch-materializing
  * Position: Surrounding character in 3D space, stuck to environment walls, floating in fog

- **Beat-Synced Editing (100-120 BPM)**:
  * ALL camera cuts PRECISELY on musical beats (kick drum or snare). The cut IS the beat
  * Camera shake intensity correlates with musical intensity — heavier during chorus/bridge
  * Handheld shake: Simulated camera instability during high-energy moments. Amplitude 3-8px, frequency matched to BPM
  * Character pose snaps align with beats — "Predatory Stillness → Explosive Violence" cycle
  * Text arrives on beats, shatters on beat drops, pulses size on each hit (100% → 120% → 100%)

- **Transition Effects (Beat-Driven)**:
  * 80% Hard cuts (0ms) — PRECISELY on musical beats
  * 20% Glitch transitions (200-500ms):
    - Chromatic aberration burst — RGB channels split then snap back
    - VHS tracking noise — horizontal distortion band sweeps vertically
    - Frame tear — image tears diagonally revealing next scene
    - Digital corruption — pixelation/blocking artifacts
  * NO dissolves, NO cross-fades, NO wipes — everything is sharp/aggressive
  * ALL transitions aligned to musical grid

- **VHS / Glitch Overlay Effects**:
  * Intermittent digital corruption — scanlines, tracking noise, color channel shift
  * Frequency INCREASES toward bridge/breakdown sections
  * Chromatic aberration PULSING on beat drops (subtle → heavy → subtle, 200ms cycle)
  * Frame tears and digital artifacts used as atmospheric texture

- **Composition with Typography**:
  * Center subject: 70% of shots place the monster CHARACTER dead center. Negative space filled with 3D TYPOGRAPHY
  * 3-Layer Depth with Text:
    - Foreground (FG): 3D Typography text / fog tendrils
    - Midground (MG): Main character/action
    - Background (BG): Environment swallowed by volumetric fog
  * KINETIC TYPOGRAPHY as Architecture: Text EXISTS IN the 3D world, affected by fog, casting glow on surfaces, rendered with depth of field
  * NEGATIVE SPACE = DARKNESS filled with floating 3D text

- **Audio-Visual Sync Pattern**:
  * Verse: WS establishing (3s) → MS villain singing (4s) → CU detail (2s) → MCU expression (3s) → Text card (2s)
  * Chorus: Rapid alternation — MS(1s) → CU(1s) → Text(1s) → MS(1s) → Dutch angle(1s)
  * Bridge: Jump cuts < 1s, heavy glitch, overlapping text chaos
  * EVERY camera cut = a beat. EVERY text arrival = a beat. EVERY pose snap = a beat
```

---

## 🎥 `video_constraint` — Video Generation Constraint

> Copy vào field **"Ràng buộc Video"**. Đây là constraint chung cho TẤT CẢ shots.
> Music DNA sẽ tự động overlay cho shots music_film.

```
### Core Production Method
1. FULL 3D CGI ANIMATION — Blender / Unreal Engine quality, After Effects compositing
2. Characters are RIGGED 3D MODELS — skeletal animation with intentionally glitchy/mechanical movement
3. Virtual camera system — cinematic lens simulation (shallow DoF, barrel distortion, chromatic aberration)
4. Post-processing: Deep Glow → Chromatic Aberration → Film Grain → Vignette

### Character Animation (3D Rigged)
- **Glitchy movement (PRIMARY)**: Intentional stutters, jerks, position snaps — corrupted puppet movement
- **Eye glow animation**: Eyes pulse brightness (60-100% cycle), intensifying during tense moments
- **Body horror**: Characters occasionally glitch — head snaps 90° unnaturally, limbs twist backward
- **Posing**: Oscillates between PREDATORY STILLNESS and EXPLOSIVE VIOLENCE
- **Walking/locomotion**: Characters RARELY walk normally — slide, float, crawl, or teleport-glitch

### Atmospheric Animation
- Volumetric fog: ALWAYS PRESENT. Slow rolling motion. Disturbed by character movement
- Dust particles: Fine particles catching under-light, 50-100 visible particles
- Fluorescent flicker: Irregular stutter (60ms on / 200ms off / 120ms on / 80ms off)
- Chromatic aberration: Static at frame edges, stronger during intense moments

### Camera System (Virtual Cinematic Camera)
- Slow push-in: 40% — constant forward dolly toward character, 2-5% zoom/sec
- Static lock: For maximum impact — stillness contrasts with violence
- Low-angle default: 50% — camera at chest level looking up
- Barrel distortion: Subtle wide-angle lens effect
- Depth of field: SHALLOW — f/1.4 to f/2.8, cinematic bokeh

### Transition Rules (Cinematic — Non-Music)
- Hard cuts: SHARP, purposeful — but NOT beat-synced in prologue/epilogue
- Pacing driven by STORY BEATS, not musical beats
- Longer holds on establishing shots (4-6s) and emotional moments (3-5s)
- Fade from black: ONLY at film start (~1000ms)
- NO glitch transitions in non-music sections

### Color Consistency
- Shadow primary: #0A0B1A, Shadow secondary: #111111
- Highlight warm: #FFCC00, Highlight cool: #00FFFF
- Accent: #FF9900, #7D12FF, #FF0033
- Midtone: #3A404A (metal), #552233 (organic)
- Split-toned grading: Cold shadows / warm neon highlights

### Hallucination Prohibition
- Do NOT add bright/cheerful colors, daylight, or natural lighting
- Do NOT add cute, clean, or friendly character designs
- Do NOT add flat 2D animation — ALWAYS 3D CGI
- Do NOT remove volumetric fog — fog is ALWAYS present
- Do NOT add smooth/natural movement — movement must be GLITCHY
- MAINTAIN dark, cinematic, mascot horror aesthetic at ALL times

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH.
```

---

## Tóm tắt: So sánh CG5 Music MV vs CG5 Narrative MV

| Element | CG5 Music MV (Original) | CG5 Narrative MV |
|---------|------------------------|-------------------|
| `style_prompt` | Full DNA (core + music) | **Core only** (no 3D text, no beat sync) |
| `narrative_music_dna` | *(không có)* | **Music-specific DNA** (3D text, beat sync, glitch) |
| Prologue shots | *(không có)* | Core style only → pure cinema |
| Music Film shots | Full DNA | Core + Music DNA → full CG5 aesthetic |
| Epilogue shots | *(không có)* | Core style only → pure cinema |
| `video_constraint` | Full (beat-synced) | **Story-beat pacing** (music overlay for music shots) |

---

## Checklist tạo Template

- [ ] Vào Settings → Prompt Templates
- [ ] Click "Create New"
- [ ] Name: `CG5 — Narrative MV`
- [ ] Tab "🎨 Phong cách" → paste **style_prompt** (Core Visual DNA)
- [ ] Tab "🎨 Phong cách" → paste **🎵 Music DNA (Narrative MV)** (narrative_music_dna)
- [ ] Tab "🎥 Video" → paste **video_constraint**
- [ ] Save
- [ ] Gán template cho drama project → chọn chế độ "Narrative MV"
