# Feature Map — Features → Files

> Maps every feature to its implementation files. Use for navigation and impact analysis.

## Drama Management

| Feature | Backend | Frontend |
|---------|---------|----------|
| CRUD dramas | `api/handlers/drama.go` → `application/services/drama_service.go` | `web/src/views/drama/` · `web/src/api/drama.ts` |
| Save outline | `handlers/drama.go:SaveOutline` | Editor views |
| Save episodes | `handlers/drama.go:SaveEpisodes` | Editor views |
| Save progress | `handlers/drama.go:SaveProgress` | Editor views |
| Drama stats | `handlers/drama.go:GetDramaStats` | Dashboard |

## Character System

| Feature | Backend | Frontend |
|---------|---------|----------|
| Character CRUD | `handlers/character_library.go` → `services/character_library_service.go` | `web/src/api/character-library.ts` |
| AI char extraction | `handlers/character_library.go:ExtractCharacters` → `services/character_library_service.go` | Episode view |
| Batch image gen | `handlers/character_library.go:BatchGenerateCharacterImages` | Character library UI |
| Character library | `services/character_library_service.go` (23KB) | Character library views |
| Char image upload | `handlers/upload.go` → `services/upload_service.go` | Upload dialog |

## Storyboard System

| Feature | Backend | Frontend |
|---------|---------|----------|
| Standard mode | `services/storyboard_service.go` | `web/src/views/storyboard/` |
| Rapid Cut mode | `handlers/rapid_cut.go` → `services/rapid_cut_service.go` | Storyboard view |
| Nursery Rhyme mode | `services/storyboard_nursery_service.go` | Storyboard view |
| MV Maker mode | `services/storyboard_mv_service.go` | Storyboard view |
| Narrative MV mode | `services/storyboard_narrative_service.go` | Storyboard view |
| Voiceover Director | `services/storyboard_composition_service.go` (24KB) | Storyboard view |
| Style distillation | `services/style_distill_service.go` (22KB) | Episode view |
| Shot preservation | Prompt: `prompts/storyboard_preserve_shots.txt` | Regenerate button |

## Image Generation

| Feature | Backend | Frontend |
|---------|---------|----------|
| Generate image | `handlers/image_generation.go` → `services/image_generation_service.go` | `web/src/views/generation/ImageGeneration.vue` |
| Batch per episode | `handlers/image_generation.go:BatchGenerateForEpisode` | Episode batch actions |
| Background extract | `handlers/image_generation.go:ExtractBackgroundsForEpisode` | Episode view |
| Upload image | `handlers/image_generation.go:UploadImage` | Upload dialog |

## Video Generation

| Feature | Backend | Frontend |
|---------|---------|----------|
| Generate video | `handlers/video_generation.go` → `services/video_generation_service.go` | `web/src/views/generation/VideoGeneration.vue` |
| Batch per episode | `handlers/video_generation.go:BatchGenerateForEpisode` | Episode batch actions |
| Video upscale | `handlers/video_generation.go:UpscaleVideo` | Video card actions |
| Reset status | `handlers/video_generation.go:ResetVideoStatus` | Video card actions |
| Video review (AI) | `handlers/video_review.go` → `services/video_review_service.go` | Video review UI |

## Video Merge

| Feature | Backend | Frontend |
|---------|---------|----------|
| Merge videos | `handlers/video_merge.go` → `services/video_merge_service.go` | Merge dialog |
| Episode finalize | `handlers/drama.go:FinalizeEpisode` → FFmpeg | Episode actions |
| Download episode | `handlers/drama.go:DownloadEpisodeVideo` | Download button |

## Scene Management

| Feature | Backend | Frontend |
|---------|---------|----------|
| Scene CRUD | `handlers/scene.go` | Scene panels |
| Scene image gen | `handlers/scene.go:GenerateSceneImage` | Scene card |
| Full prompt view | `handlers/scene.go:GetSceneFullPrompt` | Debug dialog |

## Props System

| Feature | Backend | Frontend |
|---------|---------|----------|
| Prop CRUD | `handlers/prop.go` → `services/prop_service.go` | Prop panels |
| AI prop extraction | `handlers/prop.go:ExtractProps` | Episode view |
| Prop image gen | `handlers/prop.go:GenerateImage` | Prop card |
| Associate to shots | `handlers/storyboard.go:AssociateProps` | Storyboard editor |

## Frame Prompts

| Feature | Backend | Frontend |
|---------|---------|----------|
| Generate frame prompt | `handlers/frame_prompt.go` → `services/frame_prompt_service.go` | Storyboard shot |
| List frame prompts | `handlers/storyboard.go:GetStoryboardFramePrompts` | Shot detail view |

## AI Configuration

| Feature | Backend | Frontend |
|---------|---------|----------|
| AI config CRUD | `handlers/ai_config.go` | `web/src/views/settings/` |
| Test connection | `handlers/ai_config.go:TestConnection` | Settings page |

## Audio

| Feature | Backend | Frontend |
|---------|---------|----------|
| Extract audio | `handlers/audio_extraction.go` | Audio tools |
| Batch extract | `handlers/audio_extraction.go:BatchExtractAudio` | Batch actions |

## Settings

| Feature | Backend | Frontend |
|---------|---------|----------|
| Language toggle | `handlers/settings.go` | `web/src/components/LanguageSwitcher.vue` |
| Prompt templates | `handlers/prompt_template_handler.go` → `services/prompt_template_service.go` | Settings/templates view |

## Script Generation

| Feature | Backend | Frontend |
|---------|---------|----------|
| Generate characters | `handlers/script_generation.go` → `services/script_generation_service.go` | `web/src/views/script/` |
| Script outline | Prompt: `script_outline_generation.txt` | Script editor |
| Episode generation | Prompt: `script_episode_generation.txt` | Script editor |

## Asset Management

| Feature | Backend | Frontend |
|---------|---------|----------|
| Asset CRUD | `handlers/asset.go` → `services/asset_service.go` | Asset library |
| Import from gen | `handlers/asset.go:ImportFromImageGen/ImportFromVideoGen` | Generation views |
| Duration update | `services/asset_duration_update.go` | Background |
