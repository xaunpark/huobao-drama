# Implementation Plan: Batch Generation System

> Created: 2026-03-18
> Status: Draft

## Summary
Add a "Batch Actions" capability to the Professional Editor to allow users to process all storyboards in an episode simultaneously. This includes extracting AI prompts, generating reference images (Action Sequences), and rendering videos (R2V) in a single workflow.

## Problem Statement
Currently, users must manually navigate to each shot, extract prompts, generate images, and then generate videos. For a production with dozens of shots, this is highly repetitive and time-consuming. A batch system will significantly improve productivity.

## Research Findings

### Codebase Patterns
- **Prompt Extraction**: `extractFramePrompt` in `ProfessionalEditor.vue` uses `generateFramePrompt` API and polls `taskAPI`.
- **Image Generation**: `generateFrameImage` uses `imageAPI.generateImage`.
- **Video Generation**: `generateVideo` uses `videoAPI.generateVideo`.
- **Persistence**: Most fields are saved via `dramaAPI.updateStoryboard`.
- **State Management**: Uses `sessionStorage` for temporary prompt storage, which should be shifted to DB for batch stability.

### Key Logic Locations
- `web/src/views/drama/ProfessionalEditor.vue`: Main UI and orchestration.
- `web/src/api/frame.ts`: Prompt extraction API.
- `web/src/api/image.ts`: Image generation API.
- `web/src/api/video.ts`: Video generation API.

## Proposed Solution

### Approach
1.  **Component**: Create a new component `BatchGenerationDialog.vue` to encapsulate the batch logic and UI to avoid further bloating `ProfessionalEditor.vue`.
2.  **Workflow**:
    - **Step 1: Batch Prompt Extraction**: Iterate through all `storyboards`, call `generateFramePrompt` for each. Track progress via a state map. Save results to both `sessionStorage` (for UI compatibility) and DB (for persistence).
    - **Step 2: Batch Image Generation**: Use the extracted prompts to call `imageAPI.generateImage`. Defaults to `frame_type: "action"`.
    - **Step 3: Batch Video Generation**: Once images are completed, trigger `videoAPI.generateVideo` using `reference_mode: "multiple"`.
3.  **UI**:
    - Add a "Batch" button to the storyboard panel header.
    - Show a dialog with three main sections (Prompts, Images, Videos).
    - Provide "Run All" and individual "Run Step" buttons.
    - Visual progress indicators (progress bars or status icons) for each storyboard.

### Code Examples

#### Batch Logic Skeleton
```typescript
const runBatchExtraction = async () => {
  for (const shot of storyboards.value) {
    try {
      updateBatchStatus(shot.id, 'extracting');
      const { task_id } = await generateFramePrompt(shot.id, { frame_type: 'action' });
      // Poll and save...
      await dramaAPI.updateStoryboard(shot.id, { image_prompt: result });
      updateBatchStatus(shot.id, 'prompt_ready');
    } catch (e) {
      updateBatchStatus(shot.id, 'error');
    }
  }
}
```

## Acceptance Criteria
- [ ] New "Batch Actions" button visible in the storyboard list header.
- [ ] Dialog opens with options to process all shots.
- [ ] Batch prompt extraction works and saves to DB.
- [ ] Batch image generation creates "Action Sequence" images for all shots.
- [ ] Batch video generation initiates R2V tasks using the generated images.
- [ ] Progress is visible to the user during processing.
- [ ] Individual shot view reflects the changes made during batch processing.

## Technical Considerations

### Dependencies
- Existing `dramaAPI`, `imageAPI`, `videoAPI`, `taskAPI`.
- `element-plus` for UI components.

### Risks
- **Rate Limiting**: Sending dozens of requests simultaneously might trigger API limits. Implementation should use a concurrency-limited queue (e.g., `p-limit`).
- **Long Polling**: Keeping many pollings active might be resource-heavy.
- **State Sync**: Ensuring the main editor state refreshes after batch operations.

### Alternatives Considered
- **Backend Batch**: Move batch logic to Go backend. *Rejected* because the current frontend already has well-tested logic for task orchestration and UI feedback. Frontend batching is faster to implement and provides better real-time feedback.

## Implementation Steps

1. **Phase 1: UI Foundation**
   - Create `BatchGenerationDialog.vue`.
   - Add trigger button to `ProfessionalEditor.vue`.
   - Implement basic dialog state and list of storyboards.

2. **Phase 2: Prompt Extraction**
   - Implement concurrent extraction logic.
   - Add DB persistence for the extracted prompts.

3. **Phase 3: Image & Video Generation**
   - Implement batch image generation (Action Sequence).
   - Implement batch video generation (R2V).
   - Add status polling for batch tasks.

4. **Phase 4: Polish**
   - Add "Run All" workflow.
   - Refine progress reporting.
   - Ensure state synchronization with the main editor.
