# Implementation Plan: R2V (Reference to Video) Workflow Integration

> Created: 2026-04-01
> Status: Draft
> Related Spec: N/A

## Summary
Integration of a new "Batch Reference Video Studio" (BRVS) panel within the Professional Production Board. This workflow allows direct video generation from project assets (Scene, Characters, Props) bypassing the "Shot Image" creation phase. It includes an AI-powered Video Prompt extraction step that incorporates visual style from either the chosen Template or the Project's default settings.

## Problem Statement
The current workflow requires generating a intermediate "Shot Image" (I2V) to maintain style consistency. The new R2V (Reference to Video) capabilities of models like Flow-Tool allow direct generation from reference images. Users need a streamlined batch process to:
1. Extract video-specific prompts (motion-focused) with style already incorporated.
2. Automatically select relevant reference images (Scene/Character/Prop) within a 3-slot limit.
3. Generate videos in bulk without manual shot image creation.

## Research Findings

### Codebase Patterns
- **Frontend Batch Logic:** `web/src/components/editor/BatchGenerationDialog.vue` provides the template for multi-step batch processing (Prompt -> Image -> Video).
- **Backend Prompting:** `application/services/frame_prompt_service.go` handles LLM-based extraction. It uses `PromptI18n` to resolve style priorities: `Custom > Template > Project`.
- **Backend Video Generation:** `application/services/video_generation_service.go` already supports `multiple` reference mode and `ReferenceImageURLs`.
- **Flow-Tool Integration:** `pkg/video/flowtool_video_client.go` correctly implements R2V mode when multiple reference IDs are provided.

### Best Practices
- **Style Dissolution:** Incorporating style requirements into the LLM system prompt for prompt extraction prevents redundant "double-prompting" at the generation stage.
- **Reference Priority:** Shot stability is best maintained by prioritizing the environment (Scene) followed by key subjects (Characters).

---

## Proposed Solution

### 1. Backend: Video Prompt Extraction
- **New Prompt Template:** Create `application/prompts/video_extraction.txt`.
- **Logic:** The AI will be instructed to describe motion vectors, action flow, and temporal changes.
- **Style Priority:** Use `PromptI18n.resolveEffectiveStyle` to inject the correct style into the extraction instructions.
- **Composition:** The final prompt sent to model = `[Template Constraint]` + `[AI Styled Video Prompt]`.

### 2. Frontend: Batch Reference Video Studio (BRVS)
- **Component:** `BatchReferenceVideoStudio.vue` (cloned and modified from `BatchGenerationDialog.vue`).
- **Steps:** 
    1. **Prompt Phase:** Triggers `generateFramePrompt` with `frame_type: 'video'`.
    2. **Video Phase:** 
       - Automates reference collection (max 3 slots).
       - Priority: `1 Scene Image` -> `Up to 2 Character Images` -> `Prop Images (as fillers)`.
       - Calls `videoAPI.generateVideo` with `reference_mode: 'multiple'`.

### Code Reference (Draft)

#### Backend: FrameType Constant
```go
// application/services/frame_prompt_service.go
const (
    // ...
    FrameTypeVideo FrameType = "video" // New type for R2V prompt extraction
)
```

#### Backend: video_extraction.txt Structure
```text
[Instruction Header]
Dựa vào style dưới đây, hãy mô tả hành động trong kịch bản dưới dạng prompt video tập trung vào chuyển động và sự thay đổi theo thời gian:

Style Requirement:
%s

[Shot Context Variables]
...
```

---

## Acceptance Criteria
- [ ] New `BatchReferenceVideoStudio.vue` component exists and functions correctly.
- [ ] AI successfully extracts styled video prompts when `frame_type: 'video'` is requested.
- [ ] R2V generation correctly utilizes max 3 reference images (Scene + Characters/Props).
- [ ] Project/Template style is correctly incorporated into the extracted prompt (not prepended manually in the final step).
- [ ] Existing S2V/I2V workflows remain unaffected.

## Implementation Steps

### Phase 1: Backend Infrastructure
1. Create `application/prompts/video_extraction.txt`.
2. Update `domain/models/prompt_template.go`:
   - Add `VideoExtraction` to `PromptTemplatePrompts` struct.
   - Add `"video_extraction": "video_extraction.txt"` to `PromptTypeToDefaultFile` map.
3. Update `application/services/frame_prompt_service.go`:
   - Add `FrameTypeVideo` constant.
   - Implement `generateVideoPrompt` logic (similar to `generateFirstFrame` but with specific video templates).
   - Update `processFramePromptGeneration` switch case.
4. Update `application/services/prompt_template_service.go`:
   - Add `VideoExtraction` to `getPromptFromStruct` and `GetDefaultPrompts`.

### Phase 2: Frontend Implementation
1. Create `web/src/components/editor/BatchReferenceVideoStudio.vue`.
2. Implement step-by-step logic:
   - **Step 1 (Prompt):** Sequential extraction of video prompts for selected shots.
   - **Step 2 (Generation):** Asset collection logic (Scene/Char/Prop) + R2V API call.
3. Integrate into `web/src/views/drama/ProfessionalEditor.vue`:
   - Add "Batch R2V Studio" button/icon.
   - Add translation mappings in `locales/`.

## Risks & Considerations
- **Token Limits:** Detailed style descriptions + script context might hit LLM context limits.
- **Flow-Tool R2V Performance:** Generating video from 3 diverse references requires high coherence; prompt quality is critical.
- **Reference Availability:** If a shot has no Scene or characters, handle fallback (pure T2V or error message).
