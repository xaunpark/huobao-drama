---
date: "2026-04-19"
problem_type: "logic-error"
component: "prompt_resolution"
severity: "high"
symptoms:
  - "Distilled shot-specific style (e.g., ImageStyle) is completely ignored during frame generation."
  - "The final prompt uses the original massive project-level template style instead."
root_cause: "incorrect_override_precedence"
tags:
  - "style-distillation"
  - "templates"
  - "prompt-generation"
  - "resolveEffectiveStyle"
---

# Fix: Full Template Style Overriding Distilled Per-Shot Style

## 1. Problem Statement
When running the `StyleDistillService` and successfully generating distilled per-shot `ImageStyle` descriptions, these specific styles were not being applied to the actual generated frames. Instead, the final prompts always reverted to the project's massive, general style template, rendering the entire distillation process useless and re-introducing style bleeding across shots.

## 2. Symptoms Observed
- AI successfully extracts/distills concise styles per shot (e.g., "3D CGI, volumetric fog, dark palette").
- During single-shot generation (`generateFirstFrame`, etc.), the actual LLM prompt prepends the full massive template style prompt instead of the concise per-shot string.

## 3. Investigation Steps

| Step | Action | Finding |
|------|--------|---------|
| 1 | Check DB `Storyboard.ImageStyle` | Correctly populated with concise distilled strings. |
| 2 | Trace `generateFirstFrame` | It reads `ImageStyle` and calls `WithDramaFirstFramePrompt(dramaID, shotStyle)` |
| 3 | Trace `WithDramaFirstFramePrompt` | Calls `resolveEffectiveStyle(dramaID, shotStyle, "")` |
| 4 | Analyze `resolveEffectiveStyle` | Prioritizes: 1) `style == "custom"`, 2) Template existence, 3) Dropdown key. Since `shotStyle` is a string (not "custom") and a template exists, it always returns the full template, overriding our specific string. |

## 4. Root Cause Analysis
The `resolveEffectiveStyle` function was designed early on to force templates if they exist. It assumes any incoming `style` argument is either the word `"custom"` (where it uses the `customStyle` field) or a dropdown key (like `"ghibli"`). When we pass it a fully-formed distilled string, it fails the `"custom"` check. It then successfully finds the project's overall style template and completely discards our passed-in string, returning the template.

## 5. Working Solution
Instead of trying to hack `resolveEffectiveStyle` or change every caller method signature to include `customStyle`, we created a dedicated formatting method `FormatFramePromptWithStyle` that **bypasses** the resolution chain entirely if we already know we have a distilled style.

```go
// In prompt_i18n.go
func (p *PromptI18n) FormatFramePromptWithStyle(dramaID uint, promptKey string, style string) string {
	imageRatio := "16:9"
	template := p.resolvePrompt(dramaID, promptKey)
	return formatPromptWithVars(template, map[string]string{
		"{{STYLE}}": style,
		"{{RATIO}}": imageRatio,
	})
}

// In frame_prompt_service.go
func resolveStyleForShot(sb models.Storyboard) (string, bool) {
	if sb.ImageStyle != nil && *sb.ImageStyle != "" {
		return *sb.ImageStyle, true
	}
	return "", false
}

// In generation methods:
if shotStyle, isDistilled := resolveStyleForShot(sb); isDistilled {
    // Distilled style: format template directly, bypassing resolveEffectiveStyle
    dynamicPrompt = s.promptI18n.FormatFramePromptWithStyle(dramaID, "image_first_frame", shotStyle)
} else {
    dynamicPrompt = s.promptI18n.WithDramaFirstFramePrompt(dramaID, dramaStyle)
}
```

## 6. Prevention Strategies
- **Beware `resolveEffectiveStyle`**: Never pass a fully-formed prompt text as the primary `style` argument to functions that rely on `resolveEffectiveStyle` (unless wrapped as `style="custom"` and passed in `customStyle`).
- **Separation of Concerns**: Separate the resolution of "which template to use" from the formatting of that template.

## 7. Cross-References
- Distillation Plan: `plans/shot-style-distill.md`
- Commit: `934b7f1` (Fix: bypass resolveEffectiveStyle when using distilled per-shot styles)
