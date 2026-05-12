# Sprunki Characters — Voice Profiles (Adult Version)

Tài liệu tham chiếu Voice Profile cho tất cả nhân vật Sprunki.

> **Mục đích**: Cung cấp giá trị `voice_style` để nhập vào Character DB.
> Pipeline tự động load `voice_style` → `shotContext.CharacterVoices` → gửi tới LLM distill → inject vào video prompt cuối cùng.

---

## Cách sử dụng

### 1. Khi tạo Drama, nhập `voice_style` cho mỗi Character

Copy giá trị từ cột **DB `voice_style`** bên dưới vào field `voice_style` của Character trong Drama.

### 2. Pipeline tự động xử lý

```
Character.VoiceStyle → buildShotContexts() → shotContext.CharacterVoices
→ "Oren: Male, mid-20s, laid-back...; Pinki: Female, mid-20s, warm..."
→ JSON gửi LLM distill → LLM viết video_prompt có tham chiếu voice/speaking behavior
→ getMouthConstraint() inject dialogue + lip-sync vào video_prompt_distilled
```

### 3. Kết quả trong video prompt cuối cùng

```
video_prompt_distilled = [LLM narrative with speaking cues] + [Dialogue constraint]
```

Ví dụ cuối cùng:
```
Camera slowly pushes in. Oren leans against kitchen counter, phone in one hand, eyes scanning the screen. 
He sets the phone down, rubs the back of his neck, exhales. His mouth opens — speaking with laid-back, 
slightly raspy delivery, casual rhythm...
Dialogue: Rent's due tomorrow, huh. The character is actively speaking, lip-syncing naturally to the dialog, mouth moving
```

---

## Voice Profile Registry

### Oren (Main Character)

| Field | Value |
|---|---|
| **DB `voice_style`** | `Male, mid-20s, laid-back American. Medium-low pitch, slightly raspy, casual cadence. West Coast chill vibe. Speaks in short phrases with occasional slang. Beat-boxes and hums between sentences.` |
| **Age impression** | 25 |
| **Vocal range** | Medium-low |
| **Speaking speed** | Slow-medium, unhurried |
| **Signature sounds** | Beat-boxing, finger tapping, humming melody fragments |

**Dialogue samples** (dùng trong script `[Oren]` tags):

| Line | Context | Delivery Note |
|---|---|---|
| `Yo.` | Greeting | Low, single word |
| `That's fire.` | Impressed | Genuine, slightly awed |
| `Nah, I'm good.` | Decline | Relaxed, not rude |
| `Bro, what?` | Confused/shocked | Flat, disbelief |
| `Let's get it.` | Motivated | Rising energy |
| `Aw man...` | Disappointed | Low, drawn out |
| `Pizza time, baby.` | Food excitement | Upbeat, sing-song |
| `Aight, bet.` | Agree | Quick, decisive |
| `Rent's due...` | Dread | Trailing off |
| `That was... actually kinda fun.` | Post-hustle satisfaction | Surprised, warm |

---

### Pinki (Girlfriend)

| Field | Value |
|---|---|
| **DB `voice_style`** | `Female, mid-20s, warm and bright. Medium-high pitch, clear articulation, playful but with moments of genuine concern. Coffee shop energy — friendly and engaging.` |
| **Age impression** | 24 |
| **Vocal range** | Medium-high |
| **Speaking speed** | Medium, expressive |
| **Signature sounds** | Small laugh, "hmm" when thinking, light sigh when worried |

**Dialogue samples**:

| Line | Context | Delivery Note |
|---|---|---|
| `Babe, did you pay rent?` | Checking on Oren | Concerned but loving |
| `You're the worst. But also the best.` | Teasing | Playful sarcasm, smile |
| `Okay but like... how?` | Skeptical of Oren's plan | Eyebrow raised |
| `I'm so proud of you!` | Oren succeeds | Genuine, bright |
| `Come here.` | Comfort/hug | Soft, warm |
| `You spent HOW much on Uber Eats?` | Budget shock | Rising pitch, incredulous |

---

### Simon (Best Friend)

| Field | Value |
|---|---|
| **DB `voice_style`** | `Male, mid-20s, high energy, extroverted. Medium-high pitch, fast-talking, enthusiastic. Always sounds like he just had 3 espressos. Marketing bro energy.` |
| **Age impression** | 26 |
| **Vocal range** | Medium-high |
| **Speaking speed** | Fast, excited |
| **Signature sounds** | "BRO", hand claps for emphasis, phone camera shutter |

**Dialogue samples**:

| Line | Context | Delivery Note |
|---|---|---|
| `BRO. Bro bro bro.` | Has an idea | Escalating excitement |
| `Trust me on this one.` | Pitching plan | Confident, salesman |
| `This is CONTENT.` | Sees something filmable | Hyped, pointing at phone |
| `We're literally gonna blow up.` | Optimistic | Over-the-top |
| `...okay that didn't work.` | Plan fails | Quick pivot, unbothered |
| `Networking, baby. It's all networking.` | Social event | Smooth, schmoozy |

---

### Durple (Close Friend)

| Field | Value |
|---|---|
| **DB `voice_style`** | `Male, late 20s, smooth and unhurried. Low-medium pitch, deliberate pacing, jazz musician cool. Thoughtful pauses between phrases. Barry White meets chill philosopher.` |
| **Age impression** | 28 |
| **Vocal range** | Low-medium |
| **Speaking speed** | Slow, deliberate |
| **Signature sounds** | Trumpet riffs, wine glass clink, long thoughtful "hmm..." |

**Dialogue samples**:

| Line | Context | Delivery Note |
|---|---|---|
| `Hmm. Let me think about that.` | Considering | Slow, deliberate |
| `Love isn't about grand gestures, man.` | Advice | Smooth, wise |
| `That's beautiful.` | Moved | Low, genuine |
| `I'll bring the wine.` | Social plan | Casual, refined |
| `Life's a solo, brother. Play it your way.` | Philosophy | Jazz-cool delivery |

---

### Gray (Quiet Friend)

| Field | Value |
|---|---|
| **DB `voice_style`** | `Male, late 20s, quiet and measured. Low pitch, soft volume, minimal words but each one counts. Introvert energy — comfortable with silence. Podcast narrator vibes.` |
| **Age impression** | 27 |
| **Vocal range** | Low |
| **Speaking speed** | Slow, minimal |
| **Signature sounds** | Keyboard typing, book page turning, long comfortable silences |

**Dialogue samples**:

| Line | Context | Delivery Note |
|---|---|---|
| `...` | Listening | *Nodding, fully present* |
| `You already know the answer.` | Wisdom | Calm, direct |
| `I'll pass.` | Decline social event | Gentle, no drama |
| `That's... a lot.` | Overwhelmed | Understated |
| `Here. I made you this.` | Gift | Quiet, thoughtful |
| `Have you considered just... not doing that?` | Reality check | Dry humor |

---

### Sky (Youngest)

| Field | Value |
|---|---|
| **DB `voice_style`** | `Male, early 20s, fresh and optimistic. Medium pitch, genuine curiosity in voice. New grad energy — eager but not annoying. Asks good questions.` |
| **Age impression** | 22 |
| **Vocal range** | Medium |
| **Speaking speed** | Medium, curious inflection |
| **Signature sounds** | "Oh!", "Wait—", typing on phone to look things up |

**Dialogue samples**:

| Line | Context | Delivery Note |
|---|---|---|
| `Wait, can I ask something?` | Curious | Genuine |
| `Both sides have a point though.` | Mediating | Balanced, diplomatic |
| `I googled it and actually...` | Research | Helpful, slightly nerdy |
| `Is it just me or is this kinda cool?` | Discovery | Wonder, fresh eyes |
| `I don't get it but I support you.` | Supporting friends | Earnest |

---

## Human Archetypes — Voice Profiles

Không cần nhập vào DB (xuất hiện tạm thời), nhưng dùng trong script dialogue:

| Archetype | Voice Description | Sample Lines |
|---|---|---|
| **Jake (Neighbor)** | Male, 30s, generic American friendly. Medium pitch, easygoing | "You good, bro?", "Hey, game's on tonight if you wanna come over" |
| **Manager/Boss** | Male/Female, 30-40s, professional-casual. Clear, authoritative but warm | "Nice work, man.", "Same time tomorrow?", "Here's your pay." |
| **Barista** | Any gender, 20s, friendly service. Upbeat, practiced | "The usual?", "That'll be $5.50." |
| **Bartender** | Male/Female, 30s, chill. Low-medium, conversational | "What can I get you?", "Rough day?", "That one's on the house." |

---

## Quick Reference — DB voice_style Values (Copy-Paste)

```
Oren:   Male, mid-20s, laid-back American. Medium-low pitch, slightly raspy, casual cadence. West Coast chill vibe. Speaks in short phrases with occasional slang. Beat-boxes and hums between sentences.

Pinki:  Female, mid-20s, warm and bright. Medium-high pitch, clear articulation, playful but with moments of genuine concern. Coffee shop energy — friendly and engaging.

Simon:  Male, mid-20s, high energy, extroverted. Medium-high pitch, fast-talking, enthusiastic. Always sounds like he just had 3 espressos. Marketing bro energy.

Durple: Male, late 20s, smooth and unhurried. Low-medium pitch, deliberate pacing, jazz musician cool. Thoughtful pauses between phrases. Barry White meets chill philosopher.

Gray:   Male, late 20s, quiet and measured. Low pitch, soft volume, minimal words but each one counts. Introvert energy — comfortable with silence. Podcast narrator vibes.

Sky:    Male, early 20s, fresh and optimistic. Medium pitch, genuine curiosity in voice. New grad energy — eager but not annoying. Asks good questions.
```
