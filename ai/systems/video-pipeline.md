# System: Video Pipeline — Video Generation & Merge

> Deep documentation of the video generation and merge subsystem.

## Pipeline Flow

```
Storyboard (with image) → Video Generation → Video Review → Video Merge → Episode
```

## Video Generation

### Service: `application/services/video_generation_service.go` (39KB)

#### Capabilities
- Image-to-video generation (i2v)
- Multiple provider support (5 providers)
- Batch generation per episode
- Video upscaling
- Status polling (async)

#### Provider Factory: `pkg/video/video_client.go` (10KB)
```go
type VideoClient interface {
    GenerateVideo(imageURL, prompt string, opts ...Option) (*VideoResult, error)
    GetVideoStatus(taskID string) (*VideoStatus, error)
    UpscaleVideo(videoURL string) (*VideoResult, error)
}
```

#### Providers
| Provider | File | Async Model | Notes |
|----------|------|-------------|-------|
| Doubao (Volcengine Ark) | `volces_ark_client.go` | Task ID + polling | Default provider |
| MiniMax | `minimax_client.go` | Task ID + polling | Chinese AI provider |
| OpenAI Sora | `openai_sora_client.go` | Task ID + polling | OpenAI video |
| FlowTool | `flowtool_video_client.go` | Task ID + polling | API aggregation |
| Chatfire | `chatfire_client.go` | Task ID + polling | Custom API |

#### FlowTool Generation Mode Routing

FlowTool supports 4 video modes. The `generation_mode` field in the video generation request controls which mode is used:

| FlowTool Mode | Meaning | `generation_mode` values that map here | Images |
|:---:|:---:|:---:|:---:|
| **T2V** | Text-to-Video | `"t2v"` | None |
| **I2V_S** | Image-to-Video (Start frame) | `"i2v_s"` | 1 first-frame image |
| **I2V_SE** | Image-to-Video (Start+End) | `"first_last"`, `*first_last*` | 2 images (first + last) |
| **R2V** | Reference-to-Video | `"direct_r2v"`, `"shot_r2v"`, `"shot_i2v"` | 1-N reference images |

**Routing logic** (`flowtool_video_client.go:69-77`):
```
if known R2V aliases → R2V
else if contains "first_last" → I2V_SE
else if "t2v" → T2V
else → strings.ToUpper(generation_mode)   // catches "i2v_s" → "I2V_S"
```

> [!WARNING]
> If `generation_mode` is not explicitly sent from frontend, backend defaults to `"shot_i2v"` → which maps to **R2V**, not I2V_S. This was a bug that silently misrouted First Frame shots to R2V.

#### Batch Action Studio Generation Modes

The Batch Action Studio (`BatchGenerationDialog.vue`) exposes 3 generation modes:

| UI Label | `generationMode` value | Image `frame_type` | Video `generation_mode` | FlowTool Mode |
|:---:|:---:|:---:|:---:|:---:|
| First Frame | `'first'` | `'first'` | `'i2v_s'` | **I2V_S** |
| Keyframe | `'key'` | `'key'` | (default→`'direct_r2v'`) | **R2V** |
| Action Sequence | `'action'` | `'action'` | (default→`'direct_r2v'`) | **R2V** |

#### Manual Editor (ProfessionalEditor) Auto-Detection

When using single reference mode, the editor checks `selectedImage.frame_type`:
- `frame_type === 'first'` → sends `generation_mode: 'i2v_s'` → FlowTool **I2V_S**
- Other frame types → no explicit `generation_mode` → backend defaults → FlowTool **R2V**

#### Async Pattern
1. Submit generation request → receive task_id
2. Poll provider API periodically for status
3. On completion: download video URL to local storage
4. Update DB record with local path

## Video Merge (FFmpeg)

### Service: `application/services/video_merge_service.go` (21KB)

#### FFmpeg Engine: `infrastructure/external/ffmpeg/ffmpeg.go` (29KB)

#### Capabilities
- Concatenate multiple video clips
- Crossfade transitions between clips
- Audio extraction and mixing
- Grid composition (action sequences)
- Text overlay (subtitles, watermarks)

#### Merge Flow
```
1. Collect video files for all shots in order
2. Verify all files exist locally
3. Build FFmpeg filter complex
4. Execute FFmpeg subprocess
5. Save merged output to storage
6. Create VideoMerge record
```

## Video Review (AI-Powered)

### Service: `application/services/video_review_service.go` (13KB)
### Prompt: `application/prompts/video_review_rubric.txt`

- Uses multimodal AI to review generated videos
- Scores on predefined rubric criteria
- Stores review results in VideoReview model

## Video Prompt Construction

### Constraint Prompts
- `video_constraint_prefixes.txt` — Standard video constraints
- `video_constraint_rapid_cut.txt` — Rapid cut specific constraints
- `video_distill_combined.txt` — Combined distillation output
- `video_style_distill.txt` — Style-specific constraints

### Prompt Assembly
```
Base constraint + Shot description + Style override + Motion hints
→ Final video generation prompt
```

## Critical Warnings

- All video providers are **async** — no synchronous generation
- Resource transfer cron is critical for URL permanence
- FFmpeg subprocess can hang on large files — no explicit timeout
- Video upscale may not be supported by all providers
- Provider API rate limits can cause batch failures
