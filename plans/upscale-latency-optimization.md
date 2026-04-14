# Optimize Video Upscale Polling and Validation

> Created: 2026-04-12
> Status: Implemented ✓

## Summary
Optimize the Flow-Tool API polling and validation logic for video upscale tasks to eliminate phantom delays after the AI server completes the background work. The solution removes the synchronous FFmpeg duration check for upscales, replaces it with a fast file-size validation (< 1MB), and implements a smart polling strategy with a 30s initial delay and 3s interval so that backend UI refreshes immediately upon download.

## Problem Statement
When Flow-Tool completes an Upscale job, the client UI takes several minutes to ultimately update. The root cause is a linear processing pipeline: the 10-second polling interval combined with synchronous downloading of massive 1080p files and heavy FFmpeg duration probing block the database state update. Furthermore, optimizing polling requires taking care of exact Timeout boundaries (10 minutes) safely.

## Proposed Solution

### Approach
1. **Smart Polling Logic**:
   - Limit timeout correctly without side effects: maintain 10-minute timeout by adopting `interval = 3s` and `maxAttempts = 200`.
   - Specifically for Upscale jobs (`status == models.VideoStatusUpscaling`), execute a heavy initial sleep (30 seconds) right before the very first check, avoiding useless API calls to the Flow-Tool API.
2. **Remove FFmpeg for Upscale**:
   - Bypass `s.ffmpeg.GetVideoDuration` for Upscale jobs because the duration remains completely unchanged from the original video.
3. **File Size Validation**:
   - Replace FFmpeg's error-checking duties with a lightweight OS stat check to verify downloaded file size. If `< 1MB (1,048,576 bytes)`, reject the file as corrupted and enforce a `FAILED` state.

### Code Examples
```go
// Inside pollTaskStatus
if videoGen.Status == models.VideoStatusUpscaling && attempt == 0 {
    time.Sleep(30 * time.Second)
}

// Inside completeVideoGeneration
shouldProbe := localVideoPath != nil && s.ffmpeg != nil && (duration == nil || *duration == 0) && !isUpscale

fileInfo, err := os.Stat(absPath)
if err == nil && fileInfo.Size() < 1048576 {
    // Return failed state
}
```

## Acceptance Criteria
- [ ] Upscale polling waits 30 seconds initially, then polls every 3s.
- [ ] Timeout bounding for polling remains correctly at exactly 600s (10 minutes).
- [ ] FFmpeg duration probe is intentionally avoided if the task is an upscale.
- [ ] Any downloaded upscale video < 1MB causes an immediate `models.VideoStatusUpscaleFailed` state logging a "file corrupted/empty" message.

## Technical Considerations

### Risks
- Normal video generation (`status == models.VideoStatusProcessing`) must continue to poll natively without the 30s penalty, as short T2V generations can finish inside of 10s.
- The 1MB File Size limit: Generative 1080p upscales almost always significantly exceed 1MB, keeping this logic highly robust.

## Implementation Steps

**Approach:**
- Task 1: Refactor `pollTaskStatus` to handle variable intervals safely (interval=3s, maxAttempts=200).
- Task 2: Refactor `pollTaskStatus` to sleep 30s only once during the 0th index if the context implies it's an Upscaling Task context.
- Task 3: Patch `completeVideoGeneration` block for FFmpeg duration checking to safely bypass if it's an upscale. Note: A boolean signature param `isUpscale` or similar might need tracking.
- Task 4: Follow up the `DownloadFromURLWithPath` inside `completeVideoGeneration` with the lightweight `os.Stat(absPath)` checkpoint filtering out <1MB files.

## References
- Client UI issue reporting synchronization delays
- Flow-Tool Localhost Upscale execution constraints
