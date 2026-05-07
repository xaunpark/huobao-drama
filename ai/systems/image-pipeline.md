# System: Image Pipeline — Image Generation

> Deep documentation of the image generation subsystem.

## Pipeline Flow

```
Storyboard shot → Frame Prompt → Full Image Prompt → Image Generation → Local Storage
```

## Image Generation Service

### Service: `application/services/image_generation_service.go` (43KB)

#### Capabilities
- Text-to-image generation
- Batch generation per episode
- Reference image support (base64 embedding)
- Background extraction
- Upload (user-provided images)

#### Provider Factory: `pkg/image/image_client.go`
```go
type ImageClient interface {
    GenerateImage(prompt string, opts ...Option) ([]string, error)
}
```

#### Providers
| Provider | File | Notes |
|----------|------|-------|
| OpenAI DALL-E | `openai_image_client.go` | DALL-E 3, standard |
| Google Gemini/Imagen | `gemini_image_client.go` | Google image gen |
| Volcengine | `volcengine_image_client.go` | Bytedance |
| FlowTool | `flowtool_image_client.go` | API aggregation |

## Image Prompt Assembly

### Full Prompt Construction
```
1. Global style (from Drama.Style or Drama.CustomStyle)
2. Character appearances (from Character.Appearance + Character.ImageURL ref)
3. Scene/background description (from Scene.Prompt)
4. Shot-specific image prompt (from Storyboard.ImagePrompt)
5. Per-shot distilled style (from Storyboard.ImageStyle, if populated)
```

### Frame Types
| Type | Prompt Template | Purpose |
|------|----------------|---------|
| First Frame | `image_first_frame.txt` | Opening shot of scene |
| Key Frame | `image_key_frame.txt` | Peak action moment |
| Last Frame | `image_last_frame.txt` | Scene ending |
| Action Sequence | `image_action_sequence.txt` | 3x3 grid of motion |

## Reference Image Handling

- Character reference images stored as URLs or local paths
- Converted to base64 for API calls that support reference images
- Base64 encoding done in `image_generation_service.go`
- Memory-intensive for many characters with references

## Background Extraction

```
POST /api/v1/images/episode/:episode_id/backgrounds/extract
→ Scans storyboard locations
→ Groups unique backgrounds
→ Generates background images
→ Links to Scene records
```

## Resource Transfer

After generation:
1. Provider returns image URL (CDN, with TTL)
2. Resource transfer service downloads to `./data/storage/`
3. Updates `ImageGeneration.LocalPath` field
4. `/static/` endpoint serves local copy

## Critical Warnings

- Provider image URLs expire (1-24 hours) — transfer must complete
- Base64 reference images consume significant memory
- Batch generation is sequential (not parallel) per episode
- Some providers have size/resolution constraints
- DALL-E 3 may refuse prompts with certain content
