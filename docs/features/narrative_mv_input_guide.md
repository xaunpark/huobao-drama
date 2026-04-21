# Narrative MV Mode: Input Guide

The **Narrative MV (Cinematic Short Film)** mode is designed to generate a 3-part cinematic short film driven by a structured "Story Bible". Unlike standard lyrics-synced music video modes, the Narrative MV mode treats music as emotional atmosphere while following independent story logic.

## The 3-Part Structure

Every Narrative MV requires these three acts:

1. **PROLOGUE (No Music)**: Pure cinematic storytelling. Establishes the world, characters, and the inciting mood before any music plays.
2. **MUSIC FILM**: The song plays. The visual story continues, but instead of illustrating the lyrics directly, the visuals are guided by the *emotional tone* of the music segments.
3. **EPILOGUE (No Music)**: A final, silent cinematic sequence resolving the story after the music concludes.

---

## Story Bible Format

The input for the Narrative MV mode requires a specific, marker-based syntax called a **Story Bible**. Copy your entire Story Bible into the editor's Script content area. 

Here is the exact schema:

### 1. `[STORY_BIBLE]` (Required)
The overarching narrative context, world description, and core themes.

### 2. `[CHARACTERS]` (Required)
A pipe-delimited `|` list of characters formatted as `Name | Description | Role`. The AI will reference these to maintain casting consistency.

### 3. `[PROLOGUE]` (Required)
Describe what happens before the music starts.
**Important:** You must include a `duration: Xs` tag on the same line as the Prologue marker.

### 4. `[MUSIC_SEGMENTS]` (Required)
A breakdown of the song structure with exact timestamps in `(M:SS - M:SS)` format.
*   **Emotions**: Append `— emotion: <description>` to provide the emotional tone for the segment.
*   **Sync Points**: Use indented `[SYNC_POINT]` entries to map pivotal story beats to specific music moments. Sync points require a type:
    *   `convergent`: The plot action peaks simultaneously with the song's energy.
    *   `parallel`: The visuals run alongside the song's energy but don't strictly peak with it.
    *   `irony`: The visuals contrast sharply with the music (e.g., an upbeat hook during a tragic story moment).

### 5. `[LYRICS]` (Optional)
The raw lyrics of the song to be used as timeline anchors by the AI. This is helpful for post-production editing context. 

### 6. `[EPILOGUE]` (Required)
Describe what happens after the music stops. 
**Important:** You must include a `duration: Xs` tag on the same line as the Epilogue marker.

---

## Full Example: "The Empty Studio"

Below is a complete, valid Story Bible ready to be pasted into the Professional Editor.

```text
[STORY_BIBLE]
An aging ballet dancer visits her old, abandoned rehearsal studio one last time before it's demolished. The story is about letting go of past glory and accepting the passage of time.
Emotional core: Melancholy transitioning into a bittersweet sense of peace.

[CHARACTERS]
An | 60 years old, elegant, wearing an oversized coat over her old rehearsal clothes. | Protagonist

[PROLOGUE] duration: 45s
An arrives at the dark studio. Rain streaks the windows. She slowly takes off her coat to reveal her faded ballet practice clothes. She walks to the barre, running her hand along the wood, stirring up dust. She prepares her posture, closing her eyes as she remembers the past.

[MUSIC_SEGMENTS]
(0:00 - 0:30) INTRO — emotion: quiet, atmospheric, slow buildup
(0:30 - 1:15) VERSE 1 — emotion: nostalgic, slightly sad
    [SYNC_POINT] An touches the worn out mark on the floor where she used to practice pirouettes — parallel
(1:15 - 2:00) CHORUS — emotion: swelling, powerful, overwhelming memory
    [SYNC_POINT] An executes a perfect, full-body grand jeté, transcending her age for a brief second — convergent
(2:00 - 2:30) OUTRO — emotion: fading, peaceful resolution

[LYRICS] (đã lược bỏ một phần để ngắn gọn cho mục đích làm ví dụ)
(0:30 - 0:35) Dust in the light
(0:36 - 0:40) Where I used to fly
(0:41 - 0:45) The music has faded
(0:46 - 0:50) But I still remember
(0:51 - 0:55) I try to hold on
(0:56 - 1:00) But the floor is too cold
(1:01 - 1:10) And the mirror shows someone I don't know
(1:15 - 1:20) But when I was here
(1:22 - 1:28) I was queen of the air
(1:30 - 1:35) I touched the sky
(1:38 - 1:45) And I didn't care
(1:48 - 1:55) Now it's time to say goodbye

[EPILOGUE] duration: 30s
The music ends. An is out of breath, smiling softly. She picks up her overcoat, drapes it over her shoulders, and walks out the door. The final shot lingers on the empty studio bathed in the morning light, perfectly still.
```

## How the AI Processes This

When you generate shots using the **Narrative MV** split mode:
1.  **Phase 1 (Story Planning)**: The AI reads the `[STORY_BIBLE]` and analyzes the sync points, plotting out lighting arcs and motifs guaranteed to span all three parts. It logs the total music duration.
2.  **Phase 2 (Shot Generation)**: The AI acts as a film director, composing a list of shots. The visuals are completely driven by the story logic and the acting notes mapped to the emotional tone of the current music segment. The backend automatically calculates precise start/end timestamps based on the assigned shot durations, ensuring the Music Film part precisely covers the length of the uploaded song.
3.  **UI Grouping**: In the Professional Editor, the resulting shots will be clearly divided into `🎬 PROLOGUE`, `🎵 MUSIC FILM`, and `🌅 EPILOGUE` with distinct color-coding.
