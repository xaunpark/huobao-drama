---
id: "pending-p2-optimize-distill-nplus1"
priority: "p2"
status: "pending"
created: "2026-04-19"
title: "Optimize N+1 query in buildShotContexts"
---

# Optimize N+1 query in buildShotContexts

## Context
During `/review` of the style distillation integration (Phase 1-3), an N+1 query was identified in `style_distill_service.go:buildShotContexts`.
Character associations are loaded per-shot in a loop using `s.db.Model(&sb).Association("Characters").Find(...)`.

## Task
- [ ] Refactor `loadStoryboards` or `buildShotContexts` to Preload characters up front in a single query.
- [ ] Or collect all unique IDs and do an `IN` query.

## Impact
This is `P2` strictly because the batch limits are `20` maximum per chunk, and this runs asynchronously in a goroutine. However, fixing it is a best practice.
