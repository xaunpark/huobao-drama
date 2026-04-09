# Reset All Generated Data (Batch Action Studio)

> Created: 2026-04-09
> Status: Implemented ✓

## Summary
Add a "Reset All Generated Data" feature to the Batch Action Studio to allow users to clear all AI-generated properties (image prompts, images, videos) for an episode's shots before batch processing them again. This guarantees a clean slate, removing issues where "ghost" data from prior failures causes `initTaskStates` to mistakenly show tasks as "Ready". 

## Problem Statement
When running the Batch Action Studio or switching between Generation Modes (e.g. from "Keyframe" to "Action Sequence"), some shots might already be marked as "Ready" because they have legacy completed images/videos from past runs. If a user forces a retry and the task fails mid-way, shutting down and restarting the dialog will re-fetch the *old* successful database entries, masking the failure and marking the shot as "Ready" again. There is currently no way to explicitly wipe this data to start over freshly.

## Prior Solutions
- No direct exact prior implementations for Batch Studio, but `DeleteRapidCut` acts somewhat similarly for the Rapid Cut module (wiping out generated data for an episode).

## Research Findings

### Codebase Patterns
- `BatchGenerationDialog.vue:209`: Checks for image availability using `imageAPI.listImages({ storyboard_id, frame_type })`. So ghost data stems directly from `ImageGeneration` records living in the database.
- `api/routes/routes.go:142`: Routes under `episodes.POST("/:episode_id/...")` map to bulk episode actions (e.g., `characters/extract`, `finalize`). 
- `application/services/storyboard_service.go`: Controls storyboard data manipulation.

### Best Practices
- **Atomic Operation:** Doing a reset via 100 API calls (one per storyboard) from the UI is inefficient and failure-prone. A single backend endpoint `POST /api/v1/episodes/:episode_id/clear-generated-data` ensures atomicity.
- **Foreign Key Detachment:** Rather than deeply cascading deletes which might orphan files, breaking the foreign-key link (setting `storyboard_id = null`) or applying a GoRM soft-delete to `image_generations` / `video_generations` hides them cleanly without violating referential constraints.

## Proposed Solution

### Approach
1. **Backend**: Provide a new REST endpoint `POST /api/v1/episodes/:episode_id/clear-generated-data` in `StoryboardHandler`.
   - The Service will:
     1. Nullify `image_prompt`, `image_url`, `composed_image`, `video_url`, `video_prompt` for all Storyboards where `episode_id = X`.
     2. Soft-delete `video_generations` where `storyboard_id` belongs to this episode.
     3. Soft-delete `image_generations` where `storyboard_id` belongs to this episode.
2. **Frontend UI**:
   - Add a Danger button `<el-button type="danger" plain :disabled="isBatching">` labeled "Reset All Data" to the `BatchGenerationDialog.vue` config-row.
   - Attach a prompt dialog using ElMessageBox confirming destructive intent.
   - On confirm: Call the new API endpoint, then reload `localStoryboards` using `refreshStoryboardsFromDB()` and recalculate `initTaskStates()` -> UI goes completely to 0% "Pending".

### Code Examples
```go
// application/services/storyboard_service.go
func (s *StoryboardService) ClearGeneratedData(episodeID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get all storyboard IDs for this episode
		var sbIDs []uint
		if err := tx.Model(&domain.Storyboard{}).Where("episode_id = ?", episodeID).Pluck("id", &sbIDs).Error; err != nil {
			return err
		}
		if len(sbIDs) == 0 {
			return nil
		}

		// 2. Soft-delete generation records associated with these storyboards
		tx.Where("storyboard_id IN ?", sbIDs).Delete(&domain.ImageGeneration{})
		tx.Where("storyboard_id IN ?", sbIDs).Delete(&domain.VideoGeneration{})

		// 3. Nullify storyboard references
		if err := tx.Model(&domain.Storyboard{}).Where("id IN ?", sbIDs).Updates(map[string]interface{}{
			"image_prompt":   "",
			"video_prompt":   "",
			"image_url":      "",
			"composed_image": "",
			"video_url":      "",
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
```

## Acceptance Criteria
- [ ] User can click a "Reset All Data" button in Batch Action Studio.
- [ ] A confirmation dialog is required to prevent accidental erasure.
- [ ] Upon confirmation, all shots visually revert to "Pending" and 0% progress within 2 seconds.
- [ ] Changing Generation Modes afterwards guarantees no old outputs load.
- [ ] Re-running the batch process from 0 ensures entirely fresh backend tasks.

## Technical Considerations

### Dependencies
- None beyond existing GORM DB transaction models.

### Risks
- Users accidentally wiping 2 hours of generated video. (Mitigated tightly by the bold red Confirm Dialog).
- Soft delete implies we might accumulate disk space over time. (Existing file cleanup chron jobs or housekeeping should sweep orphaned `local_path` values).

### Alternatives Considered
- *Force Retry checkbox (Frontend only):* Discarded. Force retry complicates the state machine inside `BatchGenerationDialog.vue` without cleaning up the database representation. Restarting the dialog masks eventual errors with ghost records.

## Implementation Steps

Tasks tracked in 03-tasks.md (Phase 4).

**Approach:**
- Task 1: Create `ClearGeneratedData` inside `StoryboardService` utilizing a strict DB transaction.
- Task 2: Expose via `StoryboardHandler.ClearGeneratedData` on `POST /api/v1/episodes/:episode_id/clear-generated-data`.
- Task 3: Bind endpoint in `api/routes/routes.go`.
- Task 4: Add `clearBatchData` definition into `web/src/api/drama.ts`.
- Task 5: Add `<el-button type="danger" plain>` + ElMessageBox execution block to `web/src/components/editor/BatchGenerationDialog.vue`. Add `i18n` translations in zh-CN / en-US.

## References
- `g:\VS-Project\huobao-drama\web\src\components\editor\BatchGenerationDialog.vue`
- `g:\VS-Project\huobao-drama\api\routes\routes.go`
