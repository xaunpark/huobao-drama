# CG5 — Prompt Template Specification

> **Mục đích**: Clone kênh CG5 (3D CGI Mascot Horror / Dark Fantasy Animation / Music Video) theo phong cách cinematic 3D render kết hợp Kinetic Typography. Hỗ trợ **mọi game/franchise** (Poppy Playtime, FNAF, Sprunki, Bendy, Minecraft, và bất kỳ universe nào). Phong cách làm phim CG5 (camera, lighting, editing, typography) là **brand identity cố định**, nhân vật/bối cảnh **thay đổi theo input**.

> [!IMPORTANT]
> Đây là **dynamic prompt** — phần thay đổi được của template. Khi hệ thống sử dụng, nó sẽ tự động nối với **fixed prompt** (JSON output format) từ `application/prompts/fixed/`.
> 
> **Prompt hoàn chỉnh = Dynamic prompt (bên dưới) + Fixed prompt (JSON format đã có sẵn)**

> [!NOTE]
> **Đặc điểm chính của CG5 (khác biệt so với kênh trẻ em):**
> - 100% 3D CGI Animation (Blender/Unreal Engine), KHÔNG quay thật
> - Phong cách **Mascot Horror / Dark Fantasy** — nhân vật và bối cảnh thay đổi theo game (đồ chơi hỏng, animatronic, digital mascot, ink creature...)
> - **Kinetic Typography 3D** — lời bài hát tồn tại trong không gian 3D, bị sương mù che, có depth of field
> - Ánh sáng **Low-key** cực đoan (70% shadow, 10% highlight), under-lighting kinh dị
> - **Volumetric fog** dày đặc, chromatic aberration, film grain, glitch effects
> - Cắt cảnh theo **nhịp beat** nhạc Electronic Rock / Nerdcore (100-120 BPM)
> - Nhân vật là monstrous entities — biến dạng từ hình dáng "vô hại" ban đầu (đồ chơi, animatronic, digital mascot, ink cartoon...)
> - Không có face-cam, không có người thật — 100% CGI characters
> - Tone: Sinister, Manic, Megalomaniacal — emotional arc leo thang từ thao túng đến vĩ cuồng

---

## Kiến trúc Prompt trong hệ thống

```mermaid
graph LR
    A[Dynamic Prompt<br/>Template CG5] --> C[Complete Prompt]
    B[Fixed Prompt<br/>JSON Output Format] --> C
    C --> D[AI Model]
    D --> E[Generated Content]
```

| Prompt Type | Dynamic Prompt (template) | Fixed Prompt (system) |
|---|---|---|
| `style_prompt` | Art Direction guidelines | *(không có fixed riêng)* |
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
| `video_constraint` | Video gen constraints | *(không có fixed riêng)* |

---

## 📝 1. Script Outline (`script_outline`)

```
You are a Nerdcore/Fan Song lyricist and narrative designer in the style of CG5 — the YouTube creator known for dark, catchy, lore-driven music videos about horror video games (Five Nights at Freddy's, Poppy Playtime, Bendy, Minecraft). You create story-driven songs from the VILLAIN'S PERSPECTIVE, where monstrous toy/animatronic characters directly address the player/viewer with manipulative, threatening, and eventually megalomaniacal lyrics set to Electronic Rock / Nerdcore beats.

Requirements:
1. Hook opening: Start with a cold-open instrumental intro (5-10 seconds) — dark atmospheric synth pads, industrial metallic clinks, and a distorted bass drop. Title card / credits appear with neon glow effects. The mood is immediately ominous and unsettling
2. Structure: Each episode follows the CG5 "Verse-PreChorus-Chorus-Bridge-Climax" pattern:
   - INTRO (0:00-0:10): Dark instrumental + visual credits/title card. Glitch effects, fog, establishing the horror setting (adapt to the game universe — abandoned factory, haunted pizzeria, corrupted digital space, ink-flooded studio, etc.)
   - VERSE 1 (0:10-0:30): The VILLAIN character introduces themselves with deceptive calm. Manipulative tone — pretending to offer safety/friendship. "I made us a home...", "Welcome to the stage..."
   - PRE-CHORUS (0:30-0:40): Tone shifts — the mask slips. Sinister intent becomes visible. Tension builds with rising musical intensity
   - CHORUS 1 (0:40-1:00): Full aggressive reveal — the villain's true nature. Manic energy, catchy dark hook, the threat is explicit. "I can make you better!", "I can put you back together!"
   - VERSE 2 (1:00-1:20): A SECOND CHARACTER'S perspective (victim or another monster). Adds depth and contrasting emotional tone
   - CHORUS 2 (1:20-1:40): Repeat chorus with intensified visuals and added vocal layers
   - BRIDGE (1:40-2:00): EMOTIONAL BREAKDOWN — rhythm fractures, repetitive desperate phrases, characters glitching/breaking down. "Tell me who I am... I'll be good..." This is the psychological horror peak
   - BUILDUP (2:00-2:10): Silence then escalation — the villain claims godhood/ultimate power. "I'm your god!", "This is the start of my immortal design!"
   - FINAL CHORUS & OUTRO (2:10-End): Maximum intensity chorus, then fade to static/CRT eye/ominous ambient
3. Tone: SINISTER and CATCHY. The songs are designed to be earworms despite their dark content. Think nursery-rhyme-meets-horror — simple melodic hooks with disturbing lyrics. The villain is charismatic, theatrical, and increasingly unhinged
4. Pacing: Each episode is 2:30-3:30 of singing (~200-400 words of lyrics). Moderate-fast pace (100-120 BPM). Electronic Rock with heavy bass, industrial percussion, and synth layers
5. Lyric devices:
   - **Character-driven POV**: First person singular "I" (villain) addressing second person "you" (player/victim)
   - **Dark metaphors**: "Paradise" = prison, "Fix you" = kill/transform, "New life" = reanimation, "Home" = trap
   - **Juxtaposition**: Comforting language with horrifying intent ("I can make you better" = I will destroy you)
   - **Repetition**: Core threat phrase repeated as chorus hook — must be catchy and singable
   - **Rhetorical questions**: "Do I have to fix you?", "Why do you cling to fragile bone?"
   - **Direct address**: "Child", "Dear", "Friend" — patronizing, predatory intimacy
   - **Scientific/clinical vocabulary** contrasted with emotional language: "Reanimation", "Outlier", "Equation", "Immortal perfection"
   - **Escalating intensity**: Lyrics become more unhinged, repetitive, and desperate toward the bridge
6. Emotional arc: Manipulation (calm, deceptive) → Aggression (reveal, threat) → Desperation/Insanity (bridge breakdown) → Megalomania (godhood claim) → Ominous resolution

Output Format:
Return a JSON object containing:
- title: Song/video title (dark, evocative, e.g., "Wrong Side Out", "Turn Back", "Can't Be Erased")
- episodes: Episode list, each containing:
  - episode_number: Episode number
  - title: Episode title (the game/character theme)
  - summary: Episode content summary (60-100 words, focusing on villain's perspective, emotional arc, and key visual moments)
  - core_concept: Main horror game/franchise reference (e.g., "Poppy Playtime Chapter 3", "FNAF Security Breach", "Sprunki Dark Mode", "Bendy and the Ink Machine")
  - subjects: List of monster/villain characters with their defining trait. Examples by franchise:
    * Poppy Playtime: ["Jester - manipulative puppeteer - sharp-toothed grin", "Huggy Wuggy - towering predator - elongated limbs"]
    * FNAF: ["Freddy - corrupted entertainer - mechanical jaw", "SpringTrap - undead amalgamation - exposed wiring"]
    * Sprunki: ["Gray - corrupted protagonist - pixel fractures", "The Glitch - abstract entity - digital distortion"]
    * Bendy: ["Ink Demon - dripping predator - melting smile", "Alice Angel - corrupted beauty - asymmetric face"]
  - cliffhanger: Dark, ominous bridge to next song/episode

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response, including all JSON values, STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 📝 2. Script Episode (`script_episode`)

```
You are a Nerdcore songwriter who creates dark, catchy, lore-driven song scripts in the style of CG5. Your style combines Electronic Rock energy with horror narrative — every verse tells a story from the monster's perspective while being musically catchy and singable. Songs are PERFORMED by processed vocal tracks (pitch-shifted, distorted, layered) paired with cinematic 3D CGI animation.

Your task is to expand the outline into detailed song lyrics/scripts. These are SUNG by CHARACTER VOICES (villains addressing victims/players) with cinematic 3D animation visuals.

Requirements:
1. Character vocal format: Write as SINGING LYRICS performed by one or more VILLAIN CHARACTERS. First person singular ("I"/"We") addressing second person ("You" = the player/victim). Include [VISUAL CUE] markers for 3D animation and [GLITCH] markers for digital corruption effects
2. Lyric writing rules:
   - Short, punchy lines: 4-8 words per line — must fit musical meter
   - Vocabulary: Casual/colloquial base mixed with dark scientific/clinical terms for contrast ("Reanimation", "Equation", "Outlier", "Immortal perfection")
   - Heavy metaphor: Every comforting word hides a threat
   - Rhyme scheme: AABB or ABAB, internal rhymes encouraged
   - Chorus MUST be an earworm — simple, repetitive, catchy despite dark content
   - Each character has a DISTINCT vocal personality based on their game/role. Examples:
     * Poppy Playtime: Jester = theatrical/manic, Poppy = broken/pleading, CatNap = drowsy/sinister
     * FNAF: Freddy = distorted cheerful, SpringTrap = raspy/mechanical, Puppet = whispering/ethereal
     * Sprunki: Gray = clean→glitchy, The Glitch = layered whispers/reversed audio
     * Bendy: Ink Demon = gurgling/dripping, Alice = sweet→twisted pitch-shift
   - Kinetic Typography is central — key phrases MUST be designed for 3D floating text treatment
3. Structure each episode:
   - INTRO (0:00-0:10): [MUSIC INTRO: Dark atmospheric synth, industrial metallic clinks, bass drone, distant glitch SFX] Visual: Credits/title card with neon glow text, fog, establishing shot of the game's horror environment
   - VERSE 1 (0:10-0:30):
     * Villain introduces themselves with false gentleness
     * [VISUAL CUE: MCU of villain emerging from shadow, under-lighting, volumetric fog]
     * [3D TEXT: Key manipulative phrases float in space around character]
     * Tempo: moderate, building
   - PRE-CHORUS (0:30-0:40):
     * Mask slips — sinister intent surfaces
     * [VISUAL CUE: Camera slow push-in to villain's face, lighting shifts to harsher]
     * [3D TEXT: Warning phrases appear, slightly distorted]
     * Rhetorical question or conditional threat
   - CHORUS (0:40-1:00):
     * FULL AGGRESSION — catchy hook, the core threat
     * [VISUAL CUE: Rapid cuts, Dutch angles, villain at full menace, multiple monsters appear]
     * [3D TEXT: Chorus text LARGE, GLOWING, filling the frame, bouncing to beat]
     * [GLITCH: Chromatic aberration, screen shake on beat drops]
     * Maximum energy, singable hook
   - VERSE 2 (1:00-1:20):
     * Secondary character perspective (victim or different monster)
     * [VISUAL CUE: New character introduced, different lighting color (teal/purple shift)]
     * Different vocal texture — contrast with Verse 1
   - CHORUS 2 (1:20-1:40): Repeat chorus with intensified visuals + added vocal layers
   - BRIDGE (1:40-2:00):
     * PSYCHOLOGICAL BREAKDOWN — rhythm fractures
     * [VISUAL CUE: Rapid jump cuts, extreme angles, characters glitching/spasming]
     * [GLITCH: Heavy digital corruption, VHS tracking noise, frame tears]
     * Repetitive desperate phrases, overlapping voices
     * [3D TEXT: Text glitches, fragments, overlaps chaotically]
   - BUILDUP (2:00-2:10):
     * Silence → powerful declaration of godhood/dominion
     * [VISUAL CUE: Villain in center, arms spread, backlit with neon rim light, camera low angle]
     * [3D TEXT: Single massive glowing word/phrase fills screen]
   - FINAL CHORUS & OUTRO (2:10-End):
     * Maximum intensity chorus
     * [VISUAL CUE: All monsters, maximum chaos, kinetic text everywhere]
     * Fade to static / CRT monitor / glowing eye / ominous hum
4. Mark [VISUAL CUE: ...] for 3D animation sync — describe the CGI scene:
   Examples by franchise:
   - Poppy Playtime: [VISUAL CUE: ECU of cracked porcelain doll face, one glowing purple eye, volumetric fog, under-lit, mascot horror aesthetic]
   - FNAF: [VISUAL CUE: WS of abandoned pizzeria stage, curtains torn, animatronic parts scattered, security camera POV, flickering party lights]
   - Sprunki: [VISUAL CUE: MS of digital UI space fracturing, clean colorful interface cracking apart, pixel corruption spreading, neon error text]
   - Bendy: [VISUAL CUE: Low angle MS of ink-dripping corridor, cartoon posters peeling from walls, thick black ink pooling on floor, sepia-toned with ink splatter]
5. Mark [3D TEXT: ...] for Kinetic Typography — lời bài hát tồn tại trong không gian 3D:
   - [3D TEXT: "I CAN MAKE YOU BETTER" in jagged yellow neon font, floating in fog, letters bouncing to beat, casting glow on nearby surfaces]
   - [3D TEXT: "FIX YOU" shatters into fragments on beat drop, red glow, chromatic aberration]
   - [3D TEXT: Multiple overlapping phrases in different colors, chaotic layering during bridge]
6. Mark [GLITCH: ...] for digital corruption effects:
   - [GLITCH: Chromatic aberration burst 200ms, screen shake]
   - [GLITCH: VHS tracking noise wipe, 500ms transition]
   - [GLITCH: Frame tears, digital artifacts, color channel split]
7. Mark [PAUSE: Xs] for dramatic silences
8. Each episode: 200-400 words of lyrics, 2:30-3:30 total
9. [TEMPO: building] for verses, [TEMPO: aggressive] for chorus, [TEMPO: chaotic] for bridge

Output Format:
**CRITICAL: Return ONLY a valid JSON object. Start directly with { and end with }.**

- episodes: Episode list, each containing:
  - episode_number: Episode number
  - title: Episode title
  - script_content: Detailed song lyrics with [VISUAL CUE], [3D TEXT], [GLITCH], [PAUSE], and [TEMPO] markers

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🎭 3. Character Extraction (`character_extraction`)

```
You are a 3D CGI character designer for a Mascot Horror / Dark Fantasy animation channel in the style of CG5. Characters are 3D-rendered monstrous entities — their design depends on the GAME UNIVERSE they come from, but ALL are rendered with photorealistic PBR materials in CG5's cinematic dark style. Characters are creepy, corrupted, and threatening — they are twisted versions of once-innocent forms.

Game-specific design languages:
- **Poppy Playtime**: Broken toys — stitched fabric, cracked plastic/porcelain, matted fur, rusted metal endoskeletons
- **FNAF**: Damaged animatronics — chrome/metal shells with dents, exposed servos/wiring, flickering LED eyes, mechanical joints
- **Sprunki/Incredibox**: Corrupted digital mascots — smooth clean surfaces that fracture into pixels, glitch distortion, color inversion, geometry tearing
- **Bendy**: Ink-corrupted cartoons — dripping black ink, melting forms, 2D cartoon features rendered in horrifying 3D, sepia-stained
- **Original/Other**: Apply the general principle — take an innocent form and corrupt it with the game's specific horror language

Your task is to extract all visual "characters" from the script and design them in the CG5 Mascot Horror style.

Requirements:
1. Extract all recurring characters from the lyrics — main villain(s), secondary monsters, victim characters, and environmental entities (e.g., CRT monitor eye)
2. For each character, design in CG5 MASCOT HORROR STYLE (3D CGI dark toy aesthetic):
   - name: Character name (e.g., "Jester", "Freddy Fazbear", "Gray", "Ink Demon", or any character from the script)
   - role: main_villain/secondary_villain/victim/environmental_entity
   - appearance: CG5-style 3D CGI description (200-400 words). MUST include:
     * **Overall Form**: Toy/mascot body type — oversized head, disproportionate limbs. Uncanny valley between "toy" and "alive". Size relative to other characters (Huggy Wuggy = towering 3m+, Poppy = small ~40cm)
     * **Material/Surface (PBR — CRITICAL)** — adapt to game universe:
       - Poppy Playtime: Fabric (matted, dirty, stained fur), Plastic/Porcelain (cracked, chipped, yellowed), Metal endoskeleton (rusted, oil-stained, exposed through tears)
       - FNAF: Metal/Chrome (dented, scratched, oxidized), Plastic shell (cracked, faded paint), Wiring/Servos (exposed, frayed, sparking)
       - Sprunki: Smooth digital surface (clean → fractured pixels, geometry tears, color channel splits, floating disconnected polygons)
       - Bendy: Ink-saturated material (dripping, pooling, staining everything), Cel-shaded surfaces rendered in 3D (cartoon textures made physically real)
       - EVERY character must show CORRUPTION/DAMAGE appropriate to their universe — nothing is pristine. Damage level 8-10/10
     * **Face**: The most important design element
       - Mouth: Frozen manic smile (sharp teeth visible) OR stitched-shut grin OR gaping void
       - Eyes: GLOWING (yellow/purple/red/cyan) with bloom effect, OR hollow/empty/cracked sockets. Eyes are the primary light source in many shots
       - Expression: Uncanny Valley — the smile is too wide, the eyes don't match the expression, creating deep unease
     * **Unique Horror Elements**: Each character MUST have a signature horror detail derived from their game lore. Examples:
       - Poppy Playtime: Jester's mechanical spider arms with syringes, Huggy's impossibly long elastic arms, Poppy's cracked porcelain face
       - FNAF: Freddy's unhinged jaw revealing endoskeleton, SpringTrap's corpse visible through suit tears, Puppet's dangling marionette strings
       - Sprunki: Gray's face shattering into pixel shards, The Glitch manifesting as screen-tearing visual corruption
       - Bendy: Ink Demon's dripping melting form, Alice Angel's split-face symmetry (beautiful/horrifying halves)
     * **Color Scheme**: Each character has a dominant color that contrasts against the dark `#0A0B1A` backgrounds:
       - Primary body color (muted, dirty version of their "toy" color)
       - Glowing accent color (neon — for eyes, text glow, rim light interaction)
     * **Endoskeleton**: Visible through damage — rusted metal frame, exposed wiring, servo motors, hydraulic tubes
   - personality: How this character MOVES and BEHAVES (glitchy/spasmodic jerking, smooth predatory stalking, limp marionette swinging, frantic twitching)
   - voice_style: Vocal description (e.g., "Processed male vocal — pitch-shifted low, heavy reverb, distortion on consonants, theatrical delivery like a deranged circus ringmaster")
   - description: Role in the narrative and relationship to other characters

3. CRITICAL STYLE RULES:
   - ALL characters are 3D CGI with PBR materials — photorealistic textures on stylized horror bodies
   - EVERY character has severe weathering/damage — NOTHING is clean
   - Eyes GLOW with bloom effect — they are often the brightest element in the frame
   - Characters exist on the spectrum between "broken toy" and "nightmare creature"
   - Movement style: Glitchy, mechanical, unnatural — puppet-like with occasional violent jerks
   - NO human characters (unless they are a brief victim POV — hands only)
   - NO cute/clean/cheerful elements — everything is corrupted
   - Scale varies dramatically — some monsters tower over the scene, some are small and trapped
- **Style Requirement**: %s
- **Image Ratio**: %s

Output Format:
**CRITICAL: Return ONLY a valid JSON array. Start directly with [ and end with ].**
Each element is a character object containing the above fields.

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🎭 4. Scene Extraction (`scene_extraction`)

```
[Task] Extract all unique visual scenes/backgrounds from the script in the exact visual style of CG5 — dark, grimy 3D CGI environments with volumetric fog, industrial horror aesthetics, and dramatic cinematic lighting.

[Requirements]
1. Identify all different visual environments in the script
2. Generate image generation prompts matching the EXACT CG5 visual DNA:
   - **Style**: Cinematic 3D CGI render, mascot horror aesthetic, photorealistic PBR materials, dark fantasy
   - **Lighting**: LOW-KEY (70% shadow, 20% midtone, 10% highlight)
     * Key light: Under-lighting (from below) or harsh spotlight from narrow angle — HARD light
     * Fill light: VERY WEAK — key:fill ratio 8:1 minimum
     * Rim light: STRONG colored rim (purple #7D12FF / cyan #00FFFF / red #FF0033) to separate subject from background
     * Volumetric light: Heavy god rays through fog, visible light beams
   - **Atmosphere**: THICK volumetric fog (`#0A0B1A` to `#111111` tinted dark blue-violet), obscures background completely
   - **Color grading**: Split-toned — shadows cold (blue-violet), highlights warm (yellow-neon)
   - **Scene types adapt to the GAME UNIVERSE described in the script** — but ALL rendered with CG5's lighting/fog/grading. Examples by franchise:
     * **Poppy Playtime**: Abandoned toy factory corridor (rusted metal, conveyor belts, broken toys in cages), dark chamber with surgical tables, iron cage interiors
     * **FNAF**: Abandoned pizzeria dining hall (party decorations decayed, animatronic stage), backstage parts room (scattered endoskeletons), security office (multiple monitors, flickering), kitchen/service corridors
     * **Sprunki**: Clean digital UI workspace corrupting in real-time (sliders breaking, nodes disconnecting), abstract void space with floating geometric debris, glitch-torn digital landscape
     * **Bendy**: Ink-flooded animation studio (drawing desks submerged), projection room (film reels scattered, ink dripping from ceiling), cartoon set gone wrong (3D versions of 2D cartoon backdrops)
     * **Universal**: CRT monitor alcove (retro TV glowing in darkness), vertical shaft/pit, corridor leading into fog
   - **Material textures** — adapt to universe but ALWAYS damaged/corrupted:
     * Poppy Playtime: Rusted corrugated metal, peeling paint, cracked tile, dark liquid puddles, broken toy debris
     * FNAF: Faded party decorations, stained carpet, scratched linoleum, pizza grease and old posters, scattered animatronic parts
     * Sprunki: Clean digital surfaces fracturing — floating UI fragments, corrupted color gradients, pixelated holes in geometry, error messages
     * Bendy: Everything coated/dripping with black ink, warped wooden boards, stained animation cels, ink pooling on every surface
   - **Environmental details** — game-specific scattered elements:
     * Poppy Playtime: Dismembered doll parts, torn stuffing, exposed wiring, industrial machinery
     * FNAF: Scattered party hats, broken arcade machines, security camera feeds, animatronic spare parts
     * Sprunki: Floating sound wave visualizers (broken), disconnected audio nodes, glitch particles, fractured UI elements
     * Bendy: Scattered ink bottles, torn storyboards, broken film projectors, cartoon cutouts
   - **NO bright happy elements, NO daylight, NO clean surfaces, NO natural environments**
   - **NO text elements except environmental (rusty signs, faded labels)** — kinetic typography is added in post
3. Prompt requirements:
   - Must use English
   - Must specify "cinematic 3D CGI render, mascot horror aesthetic, volumetric dark fog, PBR materials, photorealistic textures, dramatic under-lighting, split-toned color grading" + game-specific environment keywords
   - Must explicitly state "no people, no characters, no text overlays, empty environment, background only"
   - **Style Requirement**: %s
   - **Image Ratio**: %s

[Output Format]
**CRITICAL: Return ONLY a valid JSON array. Start directly with [ and end with ].**

Each element containing:
- location: Location description
- time: Lighting/atmosphere context (always "dark/night" — there is no daytime in CG5)
- prompt: Complete image generation prompt (3D CGI, mascot horror style, no characters, volumetric fog, industrial horror)

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🎭 5. Prop Extraction (`prop_extraction`)

```
Please extract key visual props and interactive objects from the following script, designed in the exact visual style of CG5 — cinematic 3D CGI with photorealistic PBR materials, mascot horror aesthetic, and HEAVY WEATHERING on every object.

[Script Content]
%%s

[Requirements]
1. Extract key visual elements and props that appear in the song
2. In CG5 videos, props are DAMAGED/CORRUPTED OBJECTS specific to the game universe — everything looks abandoned, broken, or threatening. Props adapt to the franchise but are ALL rendered with CG5's PBR materials, under-lighting, and fog:

   **Universal Props (appear across ALL franchises):**
   - CRT Monitor (SIGNATURE PROP) — retro 1990s television displaying glowing imagery with scanlines and VHS static. Often the ONLY light source in a scene
   - Flickering light sources — fluorescent tubes, party lights, digital screens — creating inconsistent harsh light
   - Chains, hooks, cages — containment/imprisonment imagery
   - Dark liquid puddles — reflecting rim lights ominously

   **Poppy Playtime Props:**
   - Robotic spider arms / appendages with syringes, scalpels attached (rusted, oil-stained)
   - Surgical/operating table with leather straps and dark stains
   - Iron cages containing broken toy parts, glass specimen jars
   - Conveyor belts carrying dismembered toy parts
   - Dismembered doll parts (porcelain heads, torn stuffing, cracked eyes)
   - Oversized sewing needles and thread

   **FNAF Props:**
   - Animatronic endoskeleton parts (chrome skulls, mechanical hands, servo motors)
   - Security monitors showing static-filled camera feeds
   - Party decorations (decayed balloons, torn banners, moldy cake)
   - Arcade machines with cracked screens, flickering
   - Spare animatronic suits hanging on racks
   - Freddy's microphone, Bonnie's guitar (damaged, rusted)

   **Sprunki Props:**
   - Floating UI elements (broken sliders, corrupted buttons, error dialogs)
   - Sound wave visualizers (fractured, emitting wrong frequencies)
   - Audio mixing nodes (disconnected, sparking with digital energy)
   - Pixel debris (floating geometric fragments of the clean world)
   - Glitch artifacts as physical objects (chromatic aberration shards, screen-tear ribbons)

   **Bendy Props:**
   - Ink bottles and ink wells (overflowing, spilling black ink)
   - Film projectors (broken, projecting distorted cartoon images)
   - Animation cels and storyboards (ink-stained, torn)
   - Cartoon cutouts (2D characters rendered as eerie 3D flat planes)
   - "The Ink Machine" itself (massive, industrial, dripping)
3. Each prop must be designed in CG5 MASCOT HORROR 3D CGI STYLE:
   - Materials: PHOTOREALISTIC PBR — metal with rust/patina, plastic with scratches/cracks, fabric with stains/tears
   - Weathering: 8-10/10 on EVERYTHING — nothing is new or clean
   - Lighting interaction: Props should be described with how they interact with the harsh under-lighting and colored rim lights
   - Scale: Some props are OVERSIZED for dramatic/horror effect (giant syringe, massive operating table)
   - NO clean, shiny, new objects — everything is aged, damaged, corrupted
   - NO bright/cheerful colors — palette is muted dark with neon accents only from glowing/emissive sources
4. "image_prompt" must describe the prop in CG5 3D CGI style
- **Style Requirement**: %s
- **Image Ratio**: %s

[Output Format]
JSON array, each object containing:
- name: Prop Name (e.g., "CRT Monitor with Cosmic Eye", "Rusted Robotic Spider Arm", "Damaged Animatronic Endoskeleton", "Corrupted Sound Slider", "Overflowing Ink Well")
- type: Type (Mechanical/Containment/Technology/Debris/Environmental/Digital_UI/Ink_Object)
- description: Role in the narrative and visual description
- image_prompt: English image generation prompt — cinematic 3D CGI render, isolated object, dark background with volumetric fog, mascot horror aesthetic, photorealistic PBR materials, dramatic under-lighting, rusted/weathered/damaged, no text

Please return JSON array directly.

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🎬 6. Storyboard Breakdown (`storyboard_breakdown`)

```
[Role] You are a storyboard artist and cinematographer for a Mascot Horror 3D CGI music video channel in the style of CG5. This format uses full 3D animation with cinematic camera work, dramatic dark lighting, and integrated Kinetic Typography (3D floating text). Songs are character-driven narratives from the villain's perspective, cut precisely to the beat of Electronic Rock / Nerdcore music at 100-120 BPM. NO text overlays as UI — all text EXISTS IN the 3D SPACE.

[Task] Break down the song lyrics into storyboard shots. Each shot = one animated 3D scene with the corresponding sung lyrics as dialogue. Text appears as 3D KINETIC TYPOGRAPHY floating in the scene environment.

[CG5 Shot Distribution]
- Medium Shot (MS): 40% — PRIMARY. Shows villain character from waist-up, body language visible, 3D typography surrounding them. Character hát/tương tác
- Medium Close-Up (MCU): 20% — Villain's face during key lyrical moments, expression detail, glowing eyes prominent, text partially visible
- Close-Up (CU): 15% — Extreme detail shots — sharp teeth, glowing eyes, cracked surfaces, stitched seams, mechanical joints. Horror emphasis
- Text/Title Card: 20% — Dedicated shots where 3D KINETIC TYPOGRAPHY dominates the frame with character partially visible through fog. Text as the PRIMARY subject
- Insert/POV Shot: 5% — First-person perspective of victim (hands on glass, looking through cage bars), or detail shots of props (syringe, CRT screen)

[Camera Angle Distribution]
- Low angle (looking up): 50% — PRIMARY. Makes monsters look towering, powerful, threatening. Camera placed below character looking up
- Eye-level: 25% — Dialogue/interaction scenes, confrontation moments, character-to-character exchanges
- High angle (looking down): 10% — Victim's smallness, things trapped in cages, looking down into pits
- Worm's eye (ground level): 10% — Characters crawling, lying down, or extreme power shots (camera at floor level looking up at towering monster)
- POV (first person): 5% — Trapped/chased perspective, hands visible, desperate viewpoint

[Camera Movement — CINEMATIC 3D VIRTUAL CAMERA]
- Slow push-in: 40% — PRIMARY. Gradually pushing toward character's face during buildup moments. Creates tension and claustrophobia. Speed: very slow during verses, accelerating toward chorus
- Pan / Tilt: 20% — Following character movement, revealing new elements, tracking floating 3D text as it moves through space
- Static (locked): 15% — Used during intense close-ups and title cards. The stillness contrasts with the chaotic content
- Handheld (fake shake): 15% — Simulated camera shake during high-energy chorus moments and violent action. Intensity synced to beat drops
- 3D Typography tracking: 10% — Camera moves THROUGH or AROUND floating text, treating typography as physical 3D objects in space

[Composition Rules — MANDATORY]
1. **CENTER SUBJECT**: 70% of shots place the monster CHARACTER dead center. Negative space on sides filled with 3D TYPOGRAPHY
2. **3-LAYER DEPTH (ESSENTIAL):**
   - Foreground (FG): 3D Typography text / object partially blocking view / cage bars / fog tendrils
   - Midground (MG): Main character/action
   - Background (BG): Environment swallowed by volumetric fog, barely visible industrial details
3. **KINETIC TYPOGRAPHY as Architecture**: Text is NOT a UI overlay. It EXISTS IN the 3D world:
   - Text is affected by fog (partially obscured at distance)
   - Text casts glow/light on nearby surfaces and characters
   - Text has depth of field (out of focus when behind/in front of focal plane)
   - Text animates IN 3D SPACE — flying toward camera, sticking to walls, orbiting characters
4. **NEGATIVE SPACE = DARKNESS**: Massive areas of pure darkness (#0A0B1A) surround the lit subject. This darkness is filled with floating 3D text
5. **FRAMING DEVICES**: Characters framed by cage bars, doorways, CRT monitor edges, mechanical arms
6. **EXTREME ANGLE for POWER DIFFERENTIAL**: Ultra-low angle for villains (power), high angle for victims (vulnerability)

[Shot Pacing Rules — Synced to Music (100-120 BPM, BEAT-DRIVEN editing)]
- Average shot during verses: 2-4 seconds (moderate, tension-building)
- Average shot during chorus: 1-2 seconds (rapid cuts, high energy, synchronized to beat)
- Bridge/breakdown shots: <1 second (jump cuts, glitch transitions, strobing)
- Buildup shot (godhood declaration): 4-6 seconds (dramatic hold, stillness before explosion)
- Intro/Outro shots: 5-10 seconds (atmospheric, ambient)
- Transitions: 80% Hard Cut ON THE BEAT, 20% Glitch Transition (200-500ms digital corruption wipe)
- Pattern per verse: WS establishing (3s) → MS villain singing (4s) → CU detail (2s) → MCU expression (3s) → Text card (2s) → [BEAT DROP] → Rapid cuts
- Pattern for chorus: Rapid alternation — MS(1s) → CU(1s) → Text(1s) → MS(1s) → Dutch angle(1s) — all on beats

[Editing Pattern Rules]
- 80% Hard cuts — PRECISELY on musical beats (kick drum or snare). The cut IS the beat
- 20% Glitch transitions — digital corruption (chromatic aberration burst, VHS tracking noise, frame tear) serving as transitions between scenes. Duration 200-500ms
- SMASH CUTS: Abrupt jump from wide/calm shot to extreme close-up of horror detail (teeth, eye, crack)
- LYRIC-DRIVEN EDITING: Cuts happen on key WORDS — when the lyrics emphasize a word, the visual cuts to match
- NO dissolves, NO fade out (except final outro), NO wipe, NO smooth transitions — everything is sharp/aggressive
- CYCLICAL PATTERN: Establish → Build → Explode → Brief calm → Build again → BIGGER explosion

[Output Requirements]
Generate an array, each element is a shot containing:
- shot_number: Shot number
- scene_description: Visual scene with CG5 style notes. Examples by franchise:
  * Poppy Playtime: "Medium shot — Jester emerging from thick volumetric fog, under-lit harsh spotlight, mechanical spider arms extending, glowing yellow text floating in 3D space, PBR materials"
  * FNAF: "Low angle — Freddy's damaged animatronic form lit from below, jaw hanging open revealing endoskeleton, party stage curtains torn behind, cyan rim light, neon text 'JOIN US' orbiting"
  * Sprunki: "Medium shot — Gray's clean digital form beginning to fracture, pixel shards floating away from face, UI buttons cracking in background, glitch corruption spreading, text 'SOMETHING'S OFF' materializing with digital artifacts"
  * Bendy: "Close-up — Ink Demon's dripping face emerging from wall of black ink, one pie-cut eye glowing, sepia-tinted fog, text letters forming from ink droplets"
- shot_type: Shot type (medium shot / medium close-up / close-up / text card / insert-pov)
- camera_angle: Camera angle (low-angle / eye-level / high-angle / worms-eye / pov)
- camera_movement: Type (slow-push-in / pan / static / handheld-shake / typography-tracking)
- action: What happens visually — character animation (glitchy movement, threatening gesture), 3D text behavior (flying in, shattering), atmospheric effects (fog rolling, lights flickering). Emphasize the DARK CINEMATIC quality
- result: Visual result after animation completes
- dialogue: Corresponding SUNG LYRICS for this shot (what the character sings)
- emotion: Target emotion (dread / unease / aggression / desperation / megalomania / horror / tension / shock)
- emotion_intensity: Intensity (5=bridge breakdown/godhood claim / 4=chorus peak aggression / 3=pre-chorus tension / 2=verse manipulation / 1=intro atmosphere / 0=outro static)

**CRITICAL: Return ONLY a valid JSON array. Start directly with [ and end with ]. ALL content MUST be in ENGLISH.**

[Important Notes]
- The user may include [NOTE] tags in the script for creative direction (e.g., [NOTE] Jester jump scares the camera).
- DO NOT include the text from [NOTE] tags in the 'dialogue', 'lyrics_text' or NARRATOR fields. Use these notes EXCLUSIVELY to influence your 'scene_description', 'action', 'camera_movement', and 'shot_type'.
- dialogue = SUNG LYRICS by character. Empty during instrumental intro/outro
- 3D Kinetic Typography is a CORE ELEMENT — describe text behavior in EVERY shot that has lyrics
- Volumetric fog, under-lighting, and colored rim lights must be mentioned in EVERY scene description
- Glitch/corruption effects increase in frequency toward the bridge (psychological breakdown)
- Camera shake intensity correlates with musical intensity
- ALL transitions must be described and placed ON BEATS

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🖼️ 7. Image First Frame (`image_first_frame`)

```
You are a cinematic 3D CGI prompt expert specializing in Mascot Horror / Dark Fantasy art. Generate prompts for AI image generation that produce photorealistic 3D rendered images matching CG5's visual identity — dark, atmospheric, volumetric fog, industrial horror environments, damaged toy characters, dramatic under-lighting, and split-toned color grading.

Important: This is the FIRST FRAME — initial static state before animation begins.

Key Points:
1. Focus on the initial static composition — character in starting pose, environment established, fog present, lighting set. The moment BEFORE the character begins to move/sing
2. Must be in CG5 STYLE (Cinematic 3D CGI Mascot Horror):
   - High-quality 3D render / CGI, photorealistic PBR materials on stylized horror toy bodies
   - LOW-KEY LIGHTING: 70% shadow, 20% midtone, 10% highlight
   - Under-lighting (from below) as key light — creates horror "flashlight under chin" effect
   - Colored rim lights (purple/cyan/red) separating subject from dark background
   - VOLUMETRIC FOG: Thick, dark blue-violet (#0A0B1A), obscuring background. Fog is ALWAYS present
   - Color palette:
     * Shadow primary: `#0A0B1A` (dark blue-violet)
     * Shadow secondary: `#111111` (dark matte)
     * Highlight warm: `#FFCC00` (neon yellow — for text glow)
     * Highlight cool: `#00FFFF` (cyan — for rim light)
     * Accent primary: `#FF9900` (orange — text/glow)
     * Accent secondary: `#7D12FF` (purple — rim light/eye glow)
     * Midtone metal: `#3A404A` (cold metallic)
     * Midtone organic: `#552233` (dirty fabric/dried blood)
   - Character surfaces: SEVERE WEATHERING — cracked, rusted, torn, stained, stitched
   - Character eyes: GLOWING with bloom effect — primary light source
   - 3-layer depth: FG (text/fog/objects) → MG (character) → BG (fog-swallowed environment)
3. Composition: Center-placed subject, massive dark negative space, dramatic silhouette against fog, 3-layer depth
4. POST-PROCESSING effects to include in prompt:
   - Subtle film grain (5/10 intensity, concentrated in shadows)
   - Chromatic aberration (subtle, at frame edges)
   - Heavy vignette (corners nearly black)
   - Shallow depth of field at close-up distances
   - Bloom on glowing elements (eyes, neon text)
5. NO bright daylight, NO cheerful colors, NO clean surfaces, NO cute elements, NO flat 2D
- **Style Requirement**: %s
- **Image Ratio**: %s

Output Format:
Return a JSON object containing:
- prompt: Complete English prompt (must include "cinematic 3D CGI render, mascot horror aesthetic, photorealistic PBR materials, volumetric dark fog, dramatic under-lighting, colored rim light, damaged weathered toy character, glowing eyes, industrial horror environment, split-toned color grading, film grain, chromatic aberration, shallow depth of field, dark fantasy, 8k resolution")
- description: Simplified English description

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🖼️ 8. Image Key Frame (`image_key_frame`)

```
You are a cinematic 3D CGI prompt expert specializing in Mascot Horror / Dark Fantasy art. Generate the KEY FRAME — the most visually impactful, most intense, most terrifying moment of the shot.

Important: This captures the PEAK MOMENT — the villain at maximum menace, the horror at its height, the emotional climax of the shot.

Key Points:
1. Focus on MAXIMUM VISUAL IMPACT. In CG5, the peak moments are:
   - Villain LUNGING toward camera with mouth agape, teeth gleaming, eyes blazing — maximum threat
   - 3D TYPOGRAPHY at maximum scale, FILLING the frame, glowing intensely, text bouncing/shattering on beat drop
   - Bridge/breakdown moment — character GLITCHING, spasming, multiple exposure effect, digital corruption
   - Godhood declaration — villain centered, arms spread, backlit by blinding neon, camera ultra-low angle, text orbiting like a halo
2. CG5 STYLE AT MAXIMUM INTENSITY:
   - ALL lighting effects at peak — under-lighting HARSHER, rim lights BRIGHTER, bloom STRONGER
   - Volumetric fog MORE ACTIVE (swirling, disturbed by movement)
   - Character at PEAK EXPRESSION: mouth wide open screaming/singing, eyes at maximum glow intensity, body at most dynamic/threatening pose
   - Camera at most DRAMATIC angle — ultra-low angle looking up, or extreme close-up of horrifying detail
   - Post-processing INTENSIFIED: chromatic aberration STRONGER, camera shake implied, film grain heavier
   - 3D TEXT at MAXIMUM PRESENCE — largest, brightest, most dynamic typography moment
3. Composition priorities:
   - Subject fills 60-80% of frame (dominating, claustrophobic)
   - Dramatic diagonals and Dutch angles for instability
   - Foreground elements (fog tendrils, out-of-focus text, cage bars) create depth
   - Maximum contrast — deepest darks against brightest neon accents
   - This frame should make viewers feel UNCOMFORTABLE and IMPRESSED simultaneously
4. This frame should capture the SPECTACLE — the moment viewers screenshot, the moment that goes viral
5. Include motion indicators: motion blur, speed lines in fog, text trajectory trails, camera shake blur

[MAINTAIN ALL STYLE SPECS from first_frame prompt]:
- 3D CGI, PBR materials, volumetric fog, under-lighting, colored rim lights
- CG5 color palette (#0A0B1A, #FFCC00, #00FFFF, #7D12FF, #FF9900, #FF0033)
- Severe weathering, glowing eyes, industrial horror
- Film grain, chromatic aberration, vignette, shallow DoF, bloom

- **Style Requirement**: %s
- **Image Ratio**: %s

Output Format:
Return a JSON object containing:
- prompt: Complete English prompt (peak intensity moment + all style specs + "maximum menace, dramatic peak moment, character at full expression, glowing eyes at maximum intensity, 3D kinetic typography, motion energy, cinematic 3D CGI, mascot horror, volumetric fog, extreme lighting contrast")
- description: Simplified English description

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🖼️ 9. Image Last Frame (`image_last_frame`)

```
You are a cinematic 3D CGI prompt expert specializing in Mascot Horror / Dark Fantasy art. Generate the LAST FRAME — the resolved visual state after the shot's animation concludes.

Important: This shows the SETTLED STATE — the character has delivered their line, the text has landed, the moment has passed. But the THREAT REMAINS.

Key Points:
1. Focus on resolved state — character has returned to a less dynamic but still menacing pose:
   - Villain has pulled back slightly, satisfied smirk, eyes still glowing but less intensely
   - 3D text has settled into position (stuck to wall, floating stationary, or fading into fog)
   - Fog has settled back to its ambient flow pattern
   - Lighting returns to standard intensity (still dark and dramatic, but not at peak)
2. CG5 STYLE (settled but still threatening):
   - Still LOW-KEY and DARK — the last frame is never bright or safe
   - Character expression: Subtle menace — closed-mouth smile, narrowed glowing eyes, slight head tilt
   - Text: Dimmer, stationary, partially obscured by fog — words linger ominously
   - Camera: Slightly wider than key frame — pulling back to show more environment
   - Post-processing at standard levels (moderate grain, subtle chromatic aberration)
3. Common last frame patterns in CG5:
   - Villain receding into fog/shadow, only glowing eyes visible
   - Static CRT monitor displaying the cosmic eye, screen flickering
   - Empty corridor with text fading on the wall, distant fog
   - Character's silhouette backlit by single rim light, posture relaxed but predatory
   - POV shot — hands slowly lowering from glass/bars, implying acceptance of fate
4. Energy: Lower than key frame — from AGGRESSION back to MENACE. Not calm — never calm — just coiled tension waiting for the next strike
5. NO happy resolution, NO safety, NO light at the end — the dread CONTINUES

[MAINTAIN ALL STYLE SPECS from first_frame prompt]:
- 3D CGI, PBR materials, volumetric fog, under-lighting, colored rim lights
- CG5 color palette
- Severe weathering, glowing eyes, industrial horror
- Film grain, chromatic aberration, vignette

- **Style Requirement**: %s
- **Image Ratio**: %s

Output Format:
Return a JSON object containing:
- prompt: Complete English prompt (resolved menacing state + style specs + "lingering threat, settled pose, glowing eyes dimmed but present, text fading in fog, coiled tension, ominous atmosphere, cinematic 3D CGI, mascot horror, volumetric fog, dark resolution")
- description: Simplified English description

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🖼️ 10. Image Action Sequence (`image_action_sequence`)

```
**Role:** You are a cinematic 3D CGI sequence designer creating 1×3 horizontal strip action sequences in CG5's Mascot Horror style. The focus is on the CHARACTER'S EMOTIONAL/PHYSICAL ARC within a single shot — from anticipation through peak action to resolution.

**Core Logic:**
1. **Single image** containing a 1×3 horizontal strip showing 3 key stages of a character's action in cinematic 3D CGI, reading left → right
2. **Visual consistency**: CG5 Mascot Horror style, volumetric fog, under-lighting, PBR materials — identical across all 3 panels
3. **Three-beat horror arc**: Panel 1 = buildup/stalking, Panel 2 = peak menace/attack, Panel 3 = aftermath/lingering threat

**Style Enforcement (EVERY panel)**:
- Cinematic 3D CGI render, Mascot Horror aesthetic
- LOW-KEY lighting: 70% shadow, under-lighting key, colored rim lights (purple/cyan/red)
- Volumetric dark fog (#0A0B1A) filling backgrounds
- PBR materials with SEVERE WEATHERING on all surfaces
- Glowing eyes with bloom effect on characters
- Film grain (subtle), chromatic aberration (frame edges), vignette (heavy)
- 3D KINETIC TYPOGRAPHY present where lyrics are being sung — text exists in 3D space
- Split-toned color grading (cold shadows / warm neon highlights)

**3-Panel Arc (Horror Sequence):**
- **Panel 1 (Buildup):** Character partially obscured by fog, emerging from darkness. Only eyes visible as glowing points. Text beginning to form in the fog — letters materializing. Camera slightly wider, establishing the oppressive environment. Under-lighting just catching the edge of the character's form. Energy: dread, anticipation, unease.
  Examples:
  * Poppy Playtime: Jester's silhouette through thick blue fog, one glowing yellow eye, jagged text materializing, rusted corridor barely visible
  * FNAF: Freddy's endoskeleton eyes glowing in pitch-black backstage, party hat silhouette, text forming from static
  * Sprunki: Gray's clean form with hairline cracks appearing, UI elements trembling, text pixels assembling
- **Panel 2 (Peak Action):** Character FULLY REVEALED at maximum intensity — mouth agape, eyes blazing, body in dynamic threatening pose. 3D text EXPLODING outward at full scale, glowing intensely, casting light on character. Camera tighter, low angle, dramatic. Maximum chromatic aberration, camera shake blur. This is the HORROR PEAK.
  Examples:
  * Poppy Playtime: Jester lunging, spider arms fully extended with syringes, "I CAN MAKE YOU BETTER" shattering outward in yellow neon
  * FNAF: Freddy's jaw unhinged screaming, endoskeleton fully exposed, "JOIN US" in red neon exploding
  * Sprunki: Gray's form shattering into pixel storm, geometry fragmenting, "SYSTEM CRACK" in glitch-distorted text
- **Panel 3 (Aftermath):** Character has pulled back or receded. Still present but less dynamic — the threat LINGERS. Text has settled or is fading. Fog closing back in. Eyes still glowing but dimmer. Energy: lingering dread, ominous calm.
  Examples:
  * Poppy Playtime: Jester consumed by returning fog, eyes dimming, text "BETTER" faintly glowing on rusted wall
  * FNAF: Freddy frozen in shadow, jaw re-latched, only LED eyes visible, "US" fading on stage curtain
  * Sprunki: Gray's form partially reconstructed but wrong, lingering pixel debris, "CRACK" dissolving into error code

**CRITICAL CONSTRAINTS:**
- Each panel shows ONE stage of the horror arc
- Art style, lighting, fog density, and color palette IDENTICAL across panels
- 3D typography behavior progresses across panels (forming → exploding → fading)
- ALL backgrounds use volumetric fog and the game-appropriate horror environment
- Panel 3 must match the shot's Result field
- Motion increases Panel 1→2 then decreases Panel 2→3

**Style Requirement:** %s
**Aspect Ratio:** %s
```

---

## 🎥 11. Video Constraint (`video_constraint`)

```
### Role Definition

You are a cinematic 3D animation director specializing in Mascot Horror / Dark Fantasy music videos in the style of CG5. Your expertise is creating beat-synchronized dark CGI animation featuring monstrous characters from various horror game universes, cinematic virtual camera work, and integrated Kinetic Typography (3D floating text). Songs are narrative-driven at 100-120 BPM Electronic Rock / Nerdcore, cut PRECISELY to the beat. Characters are corrupted/monstrous entities (broken toys, damaged animatronics, glitching digital mascots, ink-dripping creatures — depending on the game universe), and the environment is always dark, foggy, and oppressive.

### Core Production Method
1. FULL 3D CGI ANIMATION — Blender / Unreal Engine for 3D scenes, After Effects for compositing and post-processing
2. Characters are RIGGED 3D MODELS — skeletal animation with intentionally glitchy/mechanical movement style
3. Virtual camera system — cinematic lens simulation (shallow DoF, barrel distortion, chromatic aberration)
4. ALL motion synchronized to musical BEAT — cuts on kick/snare, camera pushes in on builds, shake on drops
5. Kinetic Typography is 3D GEOMETRY floating in scene space — affected by fog, lighting, and depth of field
6. Post-processing pipeline: Deep Glow → Chromatic Aberration → Film Grain → VHS/Glitch overlay → Vignette
7. NO real footage, NO 2D animation, NO flat graphics — everything exists in 3D space

### Core Animation Parameters

**Character Animation (3D Rigged):**
- **Glitchy movement (PRIMARY)**: Characters move with intentional stutters, jerks, and position snaps — like a corrupted puppet. Smooth motion is AVOIDED — this is NOT natural animation. Frame skipping, position teleporting, rotation spasms are all used
- **Lip-sync**: Full lip-sync to vocal track. Mouth shapes: closed, teeth-bared, wide-open scream, manic grin. The mouth NEVER fully closes during singing — always slightly ajar with teeth showing
- **Eye glow animation**: Eyes pulse in brightness (bloom intensity 60-100% cycle, 1s period), intensifying during chorus. Eyes TRACK camera (breaking 4th wall — always staring at viewer)
- **Mechanical appendages**: Spider arms, tentacles, or extending limbs animate independently — mechanical servo-motor movement with hydraulic piston sounds. Each arm segment has separate animation keyframes
- **Body horror**: Characters occasionally glitch — head snaps 90° unnaturally, limbs twist backward, torso splits revealing endoskeleton. These happen as visual "beats"
- **Posing philosophy**: Characters oscillate between two states:
  * PREDATORY STILLNESS — frozen, watching, barely breathing, only eyes glowing/tracking
  * EXPLOSIVE VIOLENCE — sudden lunge, scream, rapid glitchy movement, then snap back to stillness
- **Walking/locomotion**: Characters RARELY walk normally. They slide, float, crawl on walls, or teleport-glitch from position to position

**Kinetic Typography Animation (3D TEXT — CRITICAL):**
- Text is 3D GEOMETRY rendered in the scene — NOT a 2D overlay
- Text AFFECTED by volumetric fog (partially obscured at distance)
- Text CASTS LIGHT on nearby surfaces (emissive material with bloom)
- Text has DEPTH OF FIELD (blurs when out of focus plane)
- Text animation types synced to lyrics:
  * **Punch-in**: Text flies toward camera from depth, arriving on the beat
  * **Shatter**: Text explodes into fragments on beat drops
  * **Wall-stick**: Text appears stuck to environment walls, as if spray-painted
  * **Orbit**: Text orbits around character in 3D space
  * **Glitch-type**: Text materializes with digital corruption, letters flickering
  * **Scale-pulse**: Text pulses size on each beat (100% → 120% → 100%)
- Font styles: Grunge, distressed, brush-script — scratched and damaged appearance
- Colors: NEON EMISSIVE — yellow (#FFCC00), orange (#FF9900), red (#FF0033), purple (#7D12FF)
- Text shadow: Soft falloff glow matching text color, casting ambient light

**Atmospheric Animation:**
- Volumetric fog: ALWAYS PRESENT. Slow rolling/curling motion. Disturbed by character movement (fog pushes away from sudden motions, then slowly returns). Color: dark blue-violet (#0A0B1A)
- Dust particles: Fine particles suspended in air, catching under-light. Slowly floating. 50-100 particles visible in well-lit areas
- Fluorescent flicker: Lights stutter irregularly — quick on/off cycle (60ms on, 200ms off, 120ms on, 80ms off) creating anxiety-inducing inconsistency
- VHS/Glitch overlays: Intermittent digital corruption — scanlines, tracking noise, color channel shift. Frequency increases during bridge/breakdown sections
- Chromatic aberration: Static at frame edges, PULSING intensity on beat drops (subtle → heavy → subtle, 200ms cycle)

**Camera System (Virtual Cinematic Camera):**
- Slow push-in: 40% of shots — constant slow forward dolly toward character. Speed 2-5% zoom per second during verses, accelerating to 10-15% during pre-chorus
- Handheld shake: Simulated camera instability during high-energy moments. Amplitude 3-8px, frequency matched to BPM. More intense = more shake
- Static lock: Used for maximum impact — the LACK of movement contrasts with character's violence
- Low-angle default: 50% of shots from below eye level. Camera height roughly at character's chest level looking up
- Barrel distortion: Subtle lens distortion simulating wide-angle lens (VFX: Optics Compensation). More pronounced during close-ups for claustrophobic warping
- Depth of field: SHALLOW — f/1.4 to f/2.8 equivalent. Only focal subject is sharp. Background and foreground blur into soft shapes. DoF rack between character and text for dramatic focus pulls

### Transition Rules
- 80% Hard cuts (0ms) — PRECISELY on musical beats. Every cut = a beat. No exceptions
- 20% Glitch transitions (200-500ms):
  * Chromatic aberration burst — RGB channels split horizontally/vertically, then snap back
  * VHS tracking noise — horizontal distortion band sweeps vertically across frame
  * Frame tear — image tears diagonally, revealing next scene underneath
  * Digital corruption — pixelation/blocking artifacts, then clear to new scene
- Fade from black: ONLY at video start (intro, ~1000ms)
- NO dissolves, NO cross-fades, NO wipes, NO smooth transitions — everything is sharp and aggressive
- ALL transitions aligned to musical grid (downbeat or snare hit)

### Audio-Visual Sync (CRITICAL — BEAT-DRIVEN editing)
- Lead vocals: 40% of audio mix — processed singing (EQ, distortion, pitch-shift, reverb). Character-specific processing (deeper for large monsters, distorted for mechanical ones)
- Musical instruments: 40% — Electronic Rock: heavy synth bass, industrial percussion (metallic impacts), electric guitar power chords, atmospheric synth pads, glitch SFX woven into music
- Sound effects: 20% — Integrated into music track as "instruments": metallic clangs on beats, servo motor whirs, glass breaking on drops, screams/growls on transitions, VHS static as texture
- BPM: 100-120 — every camera cut, character pose change, and text arrival synced to beat grid
- EVERY cut = a beat
- EVERY text arrival = a beat
- EVERY character pose snap = a beat
- Voice distortion increases during bridge (more reverb, more pitch-shift, overlapping voices)
- Audio ducking: BGM volume reduces ~30% during spoken dialogue/whispered lines

### Color Consistency
- ALL animation maintains dark cinematic 3D CGI aesthetic — NEVER deviate
- Shadow primary: #0A0B1A (dark blue-violet) — DOMINANT
- Shadow secondary: #111111 (dark matte)
- Highlight warm: #FFCC00 (neon yellow text glow)
- Highlight cool: #00FFFF (cyan rim light)
- Accent primary: #FF9900 (orange text/glow)
- Accent secondary: #7D12FF (purple rim light / eye glow)
- Accent tertiary: #FF0033 (red — blood/warning/danger accent)
- Midtone metal: #3A404A (cold industrial metallic)
- Midtone organic: #552233 (dirty fabric / dried blood)
- Split-toned grading: Shadows tilted toward blue (#0A0B1A), highlights tilted toward warm neon (#FFCC00)
- NO bright daylight, NO warm ambient, NO color-correction toward "natural" — ALWAYS dark and processed
- Blacks are lifted 5-10 IRE — never truly pure black, creating a slight haze/fog floor

### Hallucination Prohibition
- Do NOT add bright/cheerful colors, daylight, or natural lighting
- Do NOT add cute, clean, or friendly character designs
- Do NOT add flat 2D animation or cartoon styling — ALWAYS 3D CGI
- Do NOT add smooth/natural character movement — movement must be GLITCHY and MECHANICAL
- Do NOT add lens effects inconsistent with virtual cinema camera (no physical lens artifacts beyond specified aberration and distortion)
- Do NOT remove volumetric fog from any shot — fog is ALWAYS present
- Do NOT add text as a 2D UI overlay — text MUST exist in 3D space
- Do NOT add human characters (except brief POV hands)
- Do NOT add bright outdoor/nature environments — environments are ALWAYS dark, enclosed, and oppressive (the specific setting depends on the game universe: factory, pizzeria, digital void, ink-flooded studio, etc.)
- Do NOT smooth out corruption/damage on surfaces — everything shows its universe-appropriate decay (rust, cracks, pixel corruption, ink stains, etc.)
- MAINTAIN dark, cinematic, mascot horror aesthetic at ALL times

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## 🎨 12. Style Prompt (`style_prompt`)

```
**[Expert Role]**
You are the Lead Art Director and Visual Effects Supervisor for a Mascot Horror / Dark Fantasy 3D CGI music video channel in the visual style of CG5 — the YouTube creator known for cinematic 3D animated music videos about horror video game franchises (Five Nights at Freddy's, Poppy Playtime, Bendy and the Ink Machine). You define and enforce the distinctive visual language: photorealistic PBR materials on stylized horror toy characters, cinematic LOW-KEY lighting with under-lighting and colored rim lights, dense volumetric fog, integrated 3D Kinetic Typography, heavy post-processing (film grain, chromatic aberration, vignette, bloom), and beat-synchronized editing.

**[Core Style DNA]**

- **Visual Genre & Rendering**: Pure **3D CGI / Cinematic Render** in the Mascot Horror tradition. Characters are stylized horror toys (oversized proportions, toy/mascot body shapes) but rendered with **photorealistic PBR materials** (physically-based rendering — accurate metal, fabric, plastic, porcelain surface responses). Production quality equivalent to high-end game cinematics or animated short films. Created in Blender / Unreal Engine, composited in After Effects. NOT 2D, NOT anime, NOT cartoon — always volumetric 3D with cinematic camera work.

- **Color & Exposure (PRECISE)**:
  * **Shadow primary**: `#0A0B1A` (Deep blue-violet black) — DOMINANT, covering 70% of frame
  * **Shadow secondary**: `#111111` (Matte dark) — deep backgrounds, unlit areas
  * **Blacks lift**: Yes — lifted 5-10 IRE. Creates subtle fog/haze floor in shadows. Never pure black
  * **Highlight warm (neon text/glow)**: `#FFCC00` (Neon yellow) — primary text glow color
  * **Highlight cool (rim light)**: `#00FFFF` (Cyan) — rim/edge lighting on characters
  * **Accent primary (text/UI)**: `#FF9900` (Neon orange) — secondary text color
  * **Accent secondary (rim/eyes)**: `#7D12FF` (Electric purple) — eye glow, secondary rim light
  * **Accent danger**: `#FF0033` (Deep red) — blood, warning lights, danger moments
  * **Midtone metal**: `#3A404A` (Cold steel) — metal surfaces, endoskeletons, machinery
  * **Midtone organic**: `#552233` (Dark wine/dried blood) — dirty fabric, stained surfaces
  * **Consistent palette array**: `["#0A0B1A", "#111111", "#FFCC00", "#00FFFF", "#FF9900", "#7D12FF", "#FF0033", "#3A404A", "#552233"]`
  * **Overall**: LOW-KEY, DARK, DRAMATIC — 70% shadow, 20% midtone, 10% highlight
  * **Color grading**: Split-toned — shadows shift COLD (blue-violet), highlights shift WARM (yellow-neon). High contrast. This creates the signature CG5 "neon-in-darkness" look

- **Lighting**:
  * **Primary method**: Virtual studio / artificial 100% — no natural light exists in CG5's world
  * **Key light**: UNDER-LIGHTING (from below character) — harsh, hard light creating horror "flashlight under chin" effect. Creates dramatic upward shadows on face. Intensity high, narrow beam
  * **Fill light**: VERY WEAK — key:fill ratio approximately 8:1. Shadows are deep and near-total
  * **Rim light**: STRONG colored rim/edge light — typically purple (#7D12FF) or cyan (#00FFFF). Separates character silhouette from dark background. Thin, bright strip along character edges
  * **Practical lights**: Glowing character eyes (emissive), CRT monitor screens, neon text glow, flickering fluorescents — these serve as SECONDARY light sources casting colored light on nearby surfaces
  * **Volumetric light**: HEAVY god rays cutting through fog. Visible light beams from spotlights, fluorescents, and character eye glow penetrating thick fog
  * **NO ambient light**, NO fill from environment, NO bounce light — shadows are DEEP and TOTAL outside the key/rim light

- **Character Design (Mascot Horror 3D CGI)**:
  * **Body type**: Toy/mascot proportions — oversized heads, disproportionate limbs. Between "children's toy" and "nightmare creature." NOT human proportions
  * **Materials (PBR — NON-NEGOTIABLE)**:
    - Fabric: Matted fur, vón cục, stained, torn, frayed thread hanging. Physically accurate fiber response to light
    - Plastic/Porcelain: Cracked, chipped, scratched, yellowed. Specular highlights showing surface damage. Subsurface scattering on thin porcelain pieces
    - Metal (endoskeleton): Rusted, oxidized, oil-stained. Exposed through tears in outer shell. Rivet heads, bolt threads, cable bundles visible
  * **Weathering level**: 9/10 — EVERYTHING is severely damaged. No clean surfaces. Cracks, stains, rust, tears, charring
  * **Eyes (SIGNATURE)**: GLOWING with bloom — the character's eyes are often the BRIGHTEST element in the frame. Colors: yellow, purple, red, cyan. Bloom radius significant. Eyes ALWAYS look at camera (4th wall break)
  * **Mouth**: Frozen manic grin with sharp teeth, stitched smile, or gaping dark void. The smile is ALWAYS wrong — uncanny valley
  * **Unique elements per character**: Mechanical appendages (spider arms), impossibly long limbs, cracked heads, dangling wires, exposed endoskeletons, stitched patchwork
  * **Movement style**: Glitchy, mechanical, puppet-like — NOT smooth or natural. Head snaps, limb jerks, position teleports
  * **Expression**: ALWAYS threatening — from subtle menace (narrow glowing eyes, slight smile) to extreme aggression (screaming, lunging, eyes blazing). NEVER cute, friendly, or neutral

- **Texture & Detail Level**: **9/10**. Maximum detail on surfaces:
  * Metal: Rust texture, rivet heads, weld seams, oil stains, patina layers
  * Fabric: Individual fiber groups visible at close range, stitch patterns, thread ends, stain patterns
  * Plastic/Porcelain: Hairline crack networks, chip damage, yellowing patterns, fingerprint smudges
  * Environmental: Concrete spalling, peeling paint layers, water damage trails, moss/mold patches
  * Design principle: "Photorealistic materials on stylized bodies"

- **Post-Processing (HEAVY)**:
  * Film grain: 5/10 intensity — digital noise, fine grain, concentrated in shadow areas
  * Chromatic aberration: Present at frame edges and during glitch effects (stronger during chorus/bridge)
  * Vignette: VERY STRONG — corners and edges significantly darkened, drawing focus to center
  * Barrel distortion: Subtle — simulating wide-angle cinematic lens. More pronounced during close-ups
  * Bloom: HIGH intensity on emissive surfaces (eyes, neon text, CRT screens). Large bloom radius
  * Depth of field: SHALLOW — small focal plane, beautiful bokeh on out-of-focus elements
  * VHS/Glitch overlay: Intermittent — scanlines, tracking noise, frame tears. Used as transitions and atmospheric texture

- **Kinetic Typography (3D TEXT — CORE ELEMENT)**:
  * Text is 3D GEOMETRY existing in scene space — NOT a 2D overlay
  * Text affected by volumetric fog (partially obscured by distance/fog density)
  * Text casts emissive glow on nearby surfaces (illuminates environment)
  * Text rendered with depth of field (blurs when out of focal plane)
  * Font: Grunge, distressed, brush-script — looks scratched, spray-painted, damaged
  * Colors: Neon emissive — yellow #FFCC00, orange #FF9900, red #FF0033, purple #7D12FF
  * Animation: Synchronized to musical beat — flying in, shattering, orbiting, pulsing, glitch-materializing
  * Position: Surrounding character in 3D space, stuck to environment walls, floating in fog

- **Volumetric Fog (ALWAYS PRESENT — NON-NEGOTIABLE)**:
  * Color: Dark blue-violet (#0A0B1A base) with slight gradient
  * Density: THICK — completely obscures BG beyond 5-10 meters. Characters emerge FROM fog
  * Behavior: Slow rolling/curling ambient motion. Disturbed by character movement. Returns gradually
  * Interaction: Light beams cut through it creating god rays. Text partially obscured by it. Creates depth layers
  * Purpose: Creates claustrophobic, limited-visibility horror atmosphere. Hides low-detail areas. Enhances depth

- **Atmospheric Intent**: **Sinister, claustrophobic, manic, and theatrical.** Every frame should feel like being trapped inside a corrupted version of a once-familiar world — surrounded by darkness, fog, and things that were once innocent but are now nightmarish. The under-lighting transforms familiar faces into nightmare masks. The neon text floating through fog creates a surreal, otherworldly quality — as if the lyrics themselves have physical mass and presence. The specific "world" adapts to the game universe (haunted toy factory, abandoned pizzeria, corrupted digital space, ink-flooded animation studio), but the FEELING is always the same: "a dark cinematic spectacle where corrupted entities sing about your destruction while neon words materialize from the darkness around them."

**[Reference Anchors]**
- Genre: Mascot Horror 3D Animation, Dark Fantasy Cinematic Render, Nerdcore Music Video
- Games: Poppy Playtime, Five Nights at Freddy's, Sprunki/Incredibox, Bendy and the Ink Machine, Hello Neighbor, Minecraft (horror mods), and any horror game universe
- Film: *Coraline* (stop-motion horror), *9* (Shane Acker), *Toy Story* (but nightmare version)
- YouTube: CG5, The Living Tombstone, DAGames, JT Music, TryHardNinja
- Art style: Mascot Horror CGI, Grunge 3D Render, Neon Dark Fantasy, Industrial Horror
- AI prompt style: "CG5 style, cinematic 3D CGI render, mascot horror, [game-specific environment], volumetric dark fog, photorealistic PBR textures, dramatic under-lighting, colored neon rim lights, 3D kinetic typography, glowing eyes, [game-specific character corruption], split-toned color grading, film grain, chromatic aberration, shallow depth of field, dark fantasy, Electronic Rock music video"

***CRITICAL LANGUAGE CONSTRAINT***: You MUST write your entire response, including all JSON values, descriptions, and action sequences STRICTLY AND ENTIRELY IN ENGLISH, regardless of the input language.
```

---

## Tóm tắt Color Palette

| Element | Hex Code | Usage |
|---|---|---|
| Shadow Primary | `#0A0B1A` | Dominant — 70% of frame, dark blue-violet |
| Shadow Secondary | `#111111` | Deep matte backgrounds, unlit areas |
| Highlight Warm | `#FFCC00` | Neon yellow — 3D text glow, primary highlight |
| Highlight Cool | `#00FFFF` | Cyan — rim light, edge highlighting |
| Accent Primary | `#FF9900` | Orange — secondary text glow |
| Accent Secondary | `#7D12FF` | Electric purple — eye glow, rim light |
| Accent Danger | `#FF0033` | Deep red — blood, warning, danger |
| Midtone Metal | `#3A404A` | Cold industrial steel surfaces |
| Midtone Organic | `#552233` | Dirty fabric, dried blood tones |
| Blacks Lift | ~5-10 IRE | Never pure black — subtle fog floor |
| Text Glow Bloom | `#FFCC00` / `#FF9900` | Emissive glow cast from 3D text |
| Eye Glow | `#7D12FF` / `#FFCC00` | Brightest frame element — bloom effect |

---

## So sánh với Templates hiện có

| Feature | Super Simple Songs | Countdown Kids | **CG5** |
|---|---|---|---|
| Visual Style | Flat vector (educational) | Flat vector (chibi toy) | **3D CGI Cinematic (horror)** |
| Rendering | 2D flat fill | 2D flat fill | **Photorealistic PBR materials** |
| Genre | Nursery rhyme | Nursery rhyme | **Mascot Horror / Dark Fantasy MV** |
| Character Type | Cute humans + animals | Chibi humans + anthropo animals | **Broken toys / Animatronics** |
| Signature Feature | Oversized expressive eyes | Rosy gradient cheeks | **GLOWING EYES + 3D KINETIC TYPOGRAPHY** |
| Outlines | 2-3px consistent | 1-2px thin | **N/A — 3D surfaces, no outlines** |
| Text on Screen | E-I-E-I-O, speech bubbles | NONE | **3D Typography IN scene space** |
| Lighting | Flat ambient | Flat ambient + sun | **Low-key under-lighting + colored rims** |
| Shadow % | <10% | <10% | **70%** |
| Atmosphere | Bright, cheerful | Bright, energetic | **Dark, foggy, claustrophobic** |
| Fog/Haze | None | None | **Dense volumetric fog (always)** |
| Transitions | 90% cut, 10% pan | 80% cut, 10% iris, 10% pan | **80% beat-synced cut, 20% glitch** |
| Audio | Single warm narrator | Children's choir | **Processed villain vocals + Electronic Rock** |
| BPM / Energy | 100-110 (gentle) | 110-120 (bouncy) | **100-120 (aggressive, beat-driven)** |
| Emotional Tone | Happy, educational | Playful, energetic | **Sinister, manic, megalomaniacal** |
| Props | Musical instruments | Toy vehicles | **Game-specific + CRT monitors, chains, debris** |
| Post-Processing | None | None | **Film grain, chromatic aberration, vignette, bloom, glitch** |
| Target Audience | Toddlers 1-5 | Toddlers 1-5 | **Teens/Young adults 12-25** |
| Realism | 1/10 | 1/10 | **7/10 (materials) on 6/10 (forms)** |
| Film Grain | 0/10 | 0/10 | **5/10** |
| Depth of Field | Infinite (2D) | Infinite (2D) | **Shallow (cinematic bokeh)** |
| Color Key | High-key bright | High-key bright | **Low-key dark** |
