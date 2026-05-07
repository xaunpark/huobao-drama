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
