---
id: "pending-p3-remove-dead-style-prompt"
priority: "p3"
status: "pending"
created: "2026-04-19"
title: "Remove dead WithDramaStylePrompt function"
---

# Remove dead WithDramaStylePrompt function

## Context
During `/review` of the style distillation integration, it was found that `WithDramaStylePrompt` in `prompt_i18n.go:589` has no callers left after the `style_prompt` injection was removed from `image_generation_service.go`.

## Task
- [ ] Remove `WithDramaStylePrompt` from `application/services/prompt_i18n.go`.
- [ ] Check if `ResolveEffectiveStylePublic` (if exposed for similar reasons) is still needed.
- [ ] Update any references or tests.

## Impact
`P3` nice-to-have cleanup for code hygiene.
