# Integrate Flow Tool Method 1 API (Auto-Clustering & Thin Proxy)

> Created: 2026-04-08
> Status: Implemented

## Summary
Refactor the Flow-Tool API integration in the backend to adopt Method 1 (Single-Shot Request) and transform the Go backend into a performant "Thin Proxy". This eliminates isolated requests that scatter context across multiple Google accounts, drastically reduces RAM overhead by passing raw local paths instead of Base64 strings, enforcing `"quality": "fast_lp"` for video generations natively, and stripping out hardcoded logic that improperly auto-guesses generation models based on image counts.

## Problem Statement
The current implementation:
1. Calls `/v1/upload` locally in a `for` loop, extracting `media_id`s manually causing job clusters to fracture horizontally across random bot accounts (triggering `404` & `PUBLIC_ERROR_IP_INPUT_IMAGE`).
2. Tries to forcibly predict the API `mode` solely based on reference image count. For instance, single-image sequences natively intended for the `R2V` reference model get wrongly downshifted to `I2V_S` mode. 
3. Blindly reads image files via `convertImageToBase64` indiscriminately into memory. A batch of 8 large High-Res frames creates ~80MB payload BLOAT, severely blocking concurrent Go requests and causing unnecessary IO wait times.

## Research Findings

### API Improvements
- **API Spec (Method 1) Updates:** The API explicitly states that the `images` block array takes `"local_path"`, `"url"`, or `"base64"` directly without translation. Example: `{ "mode": "R2V", "images": ["C:\\path\\to\\local\\image1.png"], "quality": "fast_lp" }`. Flow-tool handles native uploading, scaling, and session persistence efficiently.

### Codebase Bottlenecks
- `application/services/video_generation_service.go:323`: Calls `s.convertImageToBase64(resolvedImagePath)`. Extremely fatal for RAM bloat on Batch sequences.
- `pkg/video/flowtool_video_client.go:224`: Counts URLs logic `len(options.ReferenceImageURLs) > 0` to decide if `mode = "R2V"` or `mode = "I2V_S"`. 
- Clients orchestrating their own manual cross-cluster uploads: `FlowToolImageClient.uploadImage(...)`.

## Proposed Solution

### Approach
1. **Refactor Client Request Structs:**
   - Erase old `ReferenceImageIDs`, `StartImageID`, and `EndImageID` properties. Add native `Images []string` struct tag mapping back to Flow-tool Spec.

2. **Delete Memory-Bloating Methods & File Translators:**
   - Rip out `convertImageToBase64()` completely from the backend's loop cycles entirely. 
   - Ensure the resolving logic prioritizes finding `LocalPath` out of the database and passing it raw (i.e. `"C:\gemini\..."`) to `Options.ReferenceImages`.

3. **Demolish Auto-Routing Logic:**
   - Remove the `if-else` loops from the Client (`pkg/video/...`) which infer models based on length. 
   - Trust and port the specific generation mode over from the `GenerationMode` config / User DB parameter (e.g., when the user specifies Batch R2V, respect `"R2V"`; when creating a video shot defaulting to `"R2V"`, respect it).

4. **Hardcode Global Settings:**
   - Enforce `"quality": "fast_lp"` inside `flowtool_video_client.go` statically without dynamic mapping hooks. Leave image generation unpinned (as image endpoints inherently just bind parameter passthrough without overriding constraints).

### Code Examples
```go
// INSIDE application/services/video_generation_service.go 
// Extract LocalPaths natively WITHOUT base64 conversion
if videoGen.ImageURL != nil {
    resolvedImagePath := *videoGen.ImageURL
    if videoGen.ImageGenID != nil {
        // DB Fetch LocalPath...
        resolvedImagePath = *imgGen.LocalPath 
    }
    opts = append(opts, video.WithReferenceImages([]string{resolvedImagePath}))
}

// INSIDE pkg/video/flowtool_video_client.go
type flowToolJobRequest struct {
	Prompt        string   `json:"prompt"`
	Mode          string   `json:"mode"`
	Images        []string `json:"images,omitempty"`
	Quality       string   `json:"quality"`
	Ratio         string   `json:"ratio"`
	WaitForResult bool     `json:"wait_for_result"`
}

// Client blindly obeys configuration!
mode := options.Mode // Passed directly from service GenerationMode, not calculated
// mode logic mapping (direct extraction):
// if mode == "shot_i2v" then mode = "I2V_S" (if user specifically asked for i2v)
// if mode == "shot_r2v" then mode = "R2V"
// if mode == "direct_r2v" then mode = "R2V"

jobReq := flowToolJobRequest{
	Prompt:  prompt,
	Mode:    mode, 
	Images:  options.ReferenceImages, // No uploads done. Payload is just strings.
	Quality: "fast_lp", // Permanently forced
	Ratio:   ratio,
}
// JSON POST payload executes almost instantly taking < 1KB
```

## Acceptance Criteria
- [x] RAM overhead metrics crash to `<10KB` per Flow Tool Video request sent.
- [x] `convertImageToBase64()` and `c.uploadImage()` routines are eradicated from Flow-Tool dependencies in image and video client blocks.
- [x] Shot Videos successfully default into `"R2V"` mode (using the Shot reference Image), not falling back to `"I2V_S"`. 
- [x] Endpoint `/v1/jobs` captures a single array in `"images"` seamlessly processing `fast_lp` generation.
- [x] Cross-account validation mismatch (`404 media not found`) is completely eliminated via Flow-tool auto-batching integration.

## Technical Considerations

### Risks
- Local Path Accessibility: We assume the target Flow Tool environment resides on a filesystem sharing the same accessible paths as the Golang local application context. Given the `C:\` references and architecture context, it works nicely for Local agent environments. If migrated to docker or remote networking later, a reverse volume binding for paths will be required.

### Future Work 
- If Flow-tool expands its model pipeline (e.g. standard `quality` features), the forced `"fast_lp"` hardcode logic can be exposed back to a configuration UI flag for users. For now, it stays locked per directive.

## Implementation Steps

**Approach:**
- **Task 1**: Sweep `video_generation_service.go` and `image_generation_service.go`. Remove `convertImageToBase64` functionality. Remap the extraction flows to simply capture strings natively (LocalPath priority).
- **Task 2**: Sweep `flowtool_video_client.go` and `flowtool_image_client.go`. Remove upload loop constructs. Delete magic byte checks. Add raw `images` slice to structural JSON mapping.
- **Task 3**: Update the mode mappings globally in `GenerateVideo` to enforce `fast_lp`, and ensure `R2V` becomes the unadulterated routing choice for valid default single-ref Shot videos or batch setups based purely on User Mode preference, not image counting rules.

## References
- `docs/API_MULTI_ACCOUNT_CLUSTER_GUIDE.md`
