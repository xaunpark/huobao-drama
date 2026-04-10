# Manual AI Response Injection in Split Shots

> Created: 2026-04-09
> Status: Completed ✓

## Summary
Add a manual fallback capability to the "Split Shots" functionality. When users attempt to split a script via AI and face timeout issues or prefer bypassing the built-in AI calls, they can paste a raw JSON AI response directly into a new text area in the waiting UI. The system will then use this manual input to generate the storyboard, stopping any active waiting or polling.

## Problem Statement
The AI proxy gateway often times out at the hard 10-minute limit because analyzing thick scripts for Split Shots takes a very long time. When a timeout occurs, users are stuck. Permitting manual injection acts as an ultimate fallback, giving users control to retrieve responses manually via ChatGPT/Gemini web interfaces and injecting them straight into the workflow without experiencing blockers.

## Prior Solutions
N/A - This is a new fault-tolerance pattern in our UI.

## Research Findings

### Codebase Patterns
- `web/src/views/drama/EpisodeWorkflow.vue`: The Split shots action triggers `generationAPI.generateStoryboard` and usually just sets `generatingShots = true;` locking the UI until completion or error.
- `api/handlers/storyboard.go`: Has `GenerateStoryboard` which creates an async background task.
- `application/services/storyboard_service.go`: Method `processVisualUnitGeneration` combines prompt generation, AI calling, JSON parsing, and DB saving.
- **Critical Pattern 10 (Atomic State Transitions)**: Manual task execution must supersede or cancel the background polling gracefully. 

### Proposed Solution

### Approach
1. **Frontend Refactoring (`EpisodeWorkflow.vue`)**:
   - Create a `v-dialog` for the "Split Shots" loading state instead of a simple overlaid spinner.
   - The dialog will feature a `textarea` for manual JSON input and a "Xử lý thủ công" (Manual Process) button.
   - Include a "Copy Prompt" button so users can copy the payload and prompt to paste into their preferred AI web UI.
   - If user inputs JSON, hit a new endpoint: `POST /api/v1/episodes/:id/storyboards/manual-parse`.

2. **Backend Refactoring (`storyboard_service.go` & `storyboard.go`)**:
   - Extract the logic inside `processVisualUnitGeneration` that runs *after* `GenerateText(prompt)` into a distinct function: `ParseAndSaveVisualUnits(episodeID string, jsonResponse string)`.
   - Add an endpoint `POST /episodes/:episode_id/storyboards/manual-parse` taking `{ split_mode, ai_response }`.
   - Call `ParseAndSaveVisualUnits` directly to bypass `GenerateText`.

### Code Examples
```go
// api/handlers/storyboard.go
func (h *StoryboardHandler) ManualParseStoryboard(c *gin.Context) {
    episodeID := c.Param("episode_id")
    var req struct {
        SplitMode  string `json:"split_mode"`
        AiResponse string `json:"ai_response"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    // Synchronously parse and save
    err := h.storyboardService.ManualParseAndSave(episodeID, req.SplitMode, req.AiResponse)
    // ...
}
```

## Acceptance Criteria
- [ ] In the Split Shots state, users see a Dialog allowing manual input.
- [ ] Users can click a "Copy Prompt" button to grab the system+user prompt for the episode context.
- [ ] Users can paste AI JSON responses into the text area.
- [ ] When the manual response is submitted, standard storyboard mapping runs correctly.
- [ ] The workflow successfully transitions to the next phase upon manual completion.

## Technical Considerations

### Dependencies
- No new external libraries needed.

### Risks
- **Race Condition / Overwrites**: The originally triggered background timeout task could finish *after* manual completion and overwrite the user's manual fix.
  *Mitigation*: Update the backend endpoint to cleanly kill or cancel the `task_id` associated with the background generation.
- **Malformed Input**: User input might contain markdown formatting (``json ... ``) or syntax errors. 
  *Mitigation*: Rely on `utils.SafeParseAIJSON` which trims code blocks, and return clear HTTP 400 errors detailing syntax issues back to the UI.

## Implementation Steps

Tasks tracked in 03-tasks.md (or typical tracker).

**Approach:**
- Task 1: Create `POST /episodes/:episode_id/storyboards/manual-parse` endpoint.
- Task 2: Refactor `processVisualUnitGeneration` JSON-parsing and DB-saving code into a modular handler function.
- Task 3: In `EpisodeWorkflow.vue`, design the Split Shots loading Dialog with `<el-input type="textarea">`, "Cancel Auto", "Copy Prompt", and "Submit Manual" buttons.
- Task 4: Link Frontend Form to the new Backend endpoint and handle validation errors.

## References
- N/A
