# AI Pipeline Router — Generation Pipeline Context

> Load this when working with image generation, video generation, or AI text processing.

## Retrieval Order

1. **This file** (you're reading it)
2. `ai/systems/image-pipeline.md` — image generation deep dive
3. `ai/systems/video-pipeline.md` — video generation deep dive
4. `ai/memory/risks.md` — known API issues

## Pipeline Overview

```
┌──────────────┐    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  AI Text     │    │  Image       │    │  Video       │    │  FFmpeg      │
│  Generation  │───▶│  Generation  │───▶│  Generation  │───▶│  Merge       │
│  (prompts)   │    │  (t2i/i2i)   │    │  (i2v)       │    │  (compose)   │
└──────────────┘    └──────────────┘    └──────────────┘    └──────────────┘
```

## AI Text Generation

### Providers
- **OpenAI-compatible** (`pkg/ai/openai_client.go`) — GPT-4, local models via Ollama, any OpenAI-compatible API
- **Gemini** (`pkg/ai/gemini_client.go`) — Google Gemini
- **Multimodal** (`pkg/ai/multimodal.go`) — Vision/multimodal capabilities

### Key Service: `application/services/ai_service.go` (13KB)
- `NewAIService(db, log)` — creates client from DB-stored AI config
- `GetTextClient()` / `GetImageClient()` — factory methods
- AI configs stored in DB (`domain/models/ai_config.go`) — user manages via web UI

### Prompt System
- **Templates**: `application/prompts/*.txt` — 34 template files
- **I18n**: `application/services/prompt_i18n.go` (27KB) — bilingual (zh/en)
- **Custom templates**: User-defined via `domain/models/prompt_template.go` stored in DB

## Image Generation

### Providers (in `pkg/image/`)
| Provider | File | Notes |
|----------|------|-------|
| OpenAI (DALL-E) | `openai_image_client.go` | Standard DALL-E 3 |
| Gemini | `gemini_image_client.go` | Google Imagen |
| Volcengine | `volcengine_image_client.go` | Bytedance |
| FlowTool | `flowtool_image_client.go` | Aggregation API |

### Key Service: `application/services/image_generation_service.go` (43KB)
- Batch generation per episode
- Background extraction
- Reference image handling (base64 encoding)
- Local storage with transfer service

## Video Generation

### Providers (in `pkg/video/`)
| Provider | File | Notes |
|----------|------|-------|
| Doubao (Volcengine Ark) | `volces_ark_client.go` | Default, Bytedance |
| MiniMax | `minimax_client.go` | Chinese AI |
| OpenAI Sora | `openai_sora_client.go` | OpenAI video |
| FlowTool | `flowtool_video_client.go` | Aggregation API |
| Chatfire | `chatfire_client.go` | Custom API |

### Key Service: `application/services/video_generation_service.go` (39KB)
- Image-to-video generation
- Video upscaling
- Batch processing per episode
- Async status polling (video generation is asynchronous)

## FFmpeg Processing

### Key File: `infrastructure/external/ffmpeg/ffmpeg.go` (29KB)
- Video merging with transitions
- Audio extraction
- Subtitle/text overlay
- Grid composition (2x2, 2x3, 3x3 layouts)
- Episode finalization

### Key Service: `application/services/video_merge_service.go` (21KB)
- Video concatenation with transitions
- Audio mixing
- Progress tracking

## Critical Warnings

- Video generation is **async** — APIs return task IDs, need polling for completion
- Image URLs from providers **expire** — resource transfer service downloads them to local storage
- FFmpeg must be installed system-wide — not bundled
- `image_generation_service.go` handles base64 reference images — memory intensive
- Provider API keys stored in DB, not config file — user manages via web UI
- Max HTTP timeout is 30 minutes for OpenAI client (long-running generations)

## Escalation

If you need deeper understanding:
- Storyboard composition: Load `ai/systems/storyboard-system.md`
- Video review: Load `application/services/video_review_service.go`
- Style distillation: Load `application/services/style_distill_service.go`
