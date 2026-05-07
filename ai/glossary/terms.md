# Glossary — Domain Terminology

> Key terms used throughout the codebase. Understand these before diving into code.

## Core Domain Terms

| Term | Definition |
|------|-----------|
| **Drama** | A video project/series containing episodes, characters, scenes, and props |
| **Episode** | A single video within a drama, containing script content and storyboards |
| **Storyboard** | A single shot/frame in a video — the atomic unit of production |
| **Scene** | A background/location used across multiple storyboards |
| **Character** | A named entity with appearance, personality, and optional reference image |
| **Prop** | An object or item used in scenes (weapons, vehicles, etc.) |
| **Shot** | Synonymous with storyboard — one visual unit in the timeline |

## Production Modes

| Term | Definition |
|------|-----------|
| **Standard** | Default storyboard generation — story breakdown into sequential shots |
| **Visual Unit** | Splits script by visual changes, not narrative beats |
| **Nursery Rhyme** | Mode for children's music videos — lyrics-synchronized shots |
| **MV Maker** | Mode for fan-made music videos — gaming horror aesthetic (CG5-style) |
| **Narrative MV** | Three-part structure: prologue → music film → epilogue |
| **Voiceover Director** | Mode for narrator-driven content with audio strategy per shot |
| **Rapid Cut** | Post-processing mode — merges standard shots into fast-paced sequences |

## AI Pipeline Terms

| Term | Definition |
|------|-----------|
| **Frame Prompt** | Detailed image generation prompt for a specific storyboard shot |
| **Frame Type** | First frame, key frame, last frame, or action sequence |
| **Action Sequence** | A 3x3 grid of images showing motion progression |
| **Style Distillation** | AI-generated per-shot visual style from channel template |
| **Video Prompt** | Text description of motion/animation for video generation |
| **Video Constraint** | Rules for video generation (camera movement, pacing, etc.) |

## Infrastructure Terms

| Term | Definition |
|------|-----------|
| **Resource Transfer** | Background download of external URLs to local storage |
| **AI Config** | User-managed API provider configuration stored in database |
| **Prompt Template** | User-customizable override for built-in AI prompts |
| **Channel Template** | Large reference document describing a YouTube channel's aesthetic DNA |

## UI Terms

| Term | Definition |
|------|-----------|
| **Composed Image** | Image assigned to a storyboard shot (from generation or upload) |
| **Batch Generate** | Generate images/videos for all shots in an episode at once |
| **Finalize** | Merge all shot videos into final episode video via FFmpeg |

## Abbreviations

| Abbrev | Full Form |
|--------|-----------|
| **DDD** | Domain-Driven Design |
| **i2v** | Image-to-Video generation |
| **t2i** | Text-to-Image generation |
| **MV** | Music Video |
| **SPA** | Single-Page Application |
| **SFC** | Single-File Component (Vue `.vue` file) |
| **WAL** | Write-Ahead Logging (SQLite mode) |
| **CGO** | C-Go interop (disabled in this project) |
