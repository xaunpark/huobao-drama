# Pipeline Map — Data Flow Pipelines

> Traces how data moves through the system. Use for debugging and feature planning.

## Pipeline 1: Script → Storyboard

```
User Input (Script Text)
  │
  ├─ [Character Extraction Pipeline]
  │   User pastes script → POST /api/v1/episodes/:id/characters/extract
  │   → character_library_service.go:ExtractCharacters()
  │   → prompts/character_extraction.txt + AI call
  │   → Creates Character records in DB
  │
  ├─ [Scene Extraction Pipeline]
  │   Script text → (manual or via storyboard generation)
  │   → prompts/scene_extraction.txt + AI call
  │   → Creates Scene records in DB
  │
  └─ [Storyboard Generation Pipeline]
      POST /api/v1/episodes/:id/storyboards
      → storyboard_service.go (dispatcher)
        ├─ Standard: storyboard_story_breakdown.txt
        ├─ Visual Unit: storyboard_visual_unit.txt
        ├─ Nursery Rhyme: storyboard_nursery_service.go + storyboard_nursery_rhyme.txt
        ├─ MV Maker: storyboard_mv_service.go + storyboard_mv_gaming_horror.txt
        ├─ Narrative MV: storyboard_narrative_service.go + storyboard_narrative_director.txt
        └─ Voiceover: storyboard_composition_service.go
      → AI generates JSON shot list
      → Parsed into Storyboard records in DB
```

## Pipeline 2: Storyboard → Images

```
Storyboard Records (DB)
  │
  ├─ [Style Distillation]  (optional, per-shot)
  │   POST /api/v1/episodes/:id/distill-styles
  │   → style_distill_service.go
  │   → prompts/image_style_distill.txt + prompts/video_style_distill.txt
  │   → Updates Storyboard.ImageStyle, VideoStyle, VideoPromptDistilled
  │
  ├─ [Frame Prompt Generation]  (per-shot)
  │   POST /api/v1/storyboards/:id/frame-prompt
  │   → frame_prompt_service.go
  │   → prompts/image_first_frame.txt | image_key_frame.txt | image_last_frame.txt
  │   → Creates FramePrompt records
  │
  └─ [Image Generation]  (per-shot or batch)
      POST /api/v1/images/episode/:id/batch
      → image_generation_service.go
      → Constructs full prompt: style + character + scene + shot description
      → pkg/image/ client (OpenAI DALL-E / Gemini / Volcengine / FlowTool)
      → Creates ImageGeneration records
      → resource_transfer_service downloads URL → local storage
```

## Pipeline 3: Images → Videos

```
ImageGeneration Records (DB, with local images)
  │
  └─ [Video Generation]  (per-image or batch)
      POST /api/v1/videos/episode/:id/batch
      → video_generation_service.go
      → Constructs video prompt: motion + constraints + style
      → pkg/video/ client (Doubao / MiniMax / Sora / FlowTool / Chatfire)
      → Returns task_id (ASYNC!)
      → Polling loop checks provider API for completion
      → Creates VideoGeneration records
      → resource_transfer_service downloads URL → local storage
      │
      ├─ [Generation Mode Routing (FlowTool)]
      │   Frontend frame_type → generation_mode → FlowTool mode:
      │     'first'  + generation_mode='i2v_s'     → I2V_S  (start frame strict)
      │     'key'    + (default='shot_i2v')         → R2V    (reference)
      │     'action' + (default='direct_r2v')       → R2V    (reference)
      │     'first_last' + generation_mode='first_last' → I2V_SE (start+end)
      │   See: ai/systems/video-pipeline.md for full routing matrix
      │
      ├─ [Video Upscale]  (optional)
      │   POST /api/v1/videos/:id/upscale
      │   → video_generation_service.go:UpscaleVideo()
      │   → Provider-specific upscale API
      │
      └─ [Video Review]  (optional, AI-powered)
          POST /api/v1/videos/:id/review
          → video_review_service.go
          → prompts/video_review_rubric.txt + multimodal AI
          → Creates VideoReview record with scores
```

## Pipeline 4: Videos → Final Episode

```
VideoGeneration Records (DB, with local videos)
  │
  ├─ [Video Merge]
  │   POST /api/v1/video-merges
  │   → video_merge_service.go
  │   → ffmpeg.go: concat videos with transitions
  │   → Creates VideoMerge record with merged file
  │
  └─ [Episode Finalization]
      POST /api/v1/episodes/:id/finalize
      → drama_service.go:FinalizeEpisode()
      → ffmpeg.go: compose all shots in order
      → Adds transitions, audio, overlays
      → Creates final episode video file
      → Updates Episode.VideoURL
```

## Pipeline 5: Rapid Cut (Parallel Pipeline)

```
Standard Storyboard
  │
  └─ POST /api/v1/episodes/:id/rapid-cut
     → rapid_cut_service.go
     → prompts/rapid_cut_merge.txt
     → Merges multiple standard shots into fast-paced sequences
     → Creates new Storyboard records (is_production=true, pacing_mode="rapid_cut")
     → Uses prompts/image_action_sequence_rapid_cut.txt for 3x3 grids
     → Uses prompts/video_constraint_rapid_cut.txt for video constraints
```

## Data Lifecycle

```
External URL (provider CDN)
  → resource_transfer_scheduler (cron)
    → Downloads to ./data/storage/
      → Updates record.LocalPath
        → Served via /static endpoint
```

## Critical Async Points

1. **Video generation** — providers return task_id, need polling
2. **Resource transfer** — background cron downloads external URLs
3. **Batch operations** — image/video batch gen are sequential, not parallel
4. **FFmpeg** — subprocess execution, can timeout on large files
