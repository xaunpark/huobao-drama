---
priority: p4
status: ready
created: 2026-04-16
source: plans/mv-maker-mode.md
---

# Deprecate Standalone Nursery Rhyme → MV Maker Genre Profile

After MV Maker mode is stable and proven, consider migrating `nursery_rhyme` standalone mode into `mv_maker` as a genre profile (`mv_maker` + genre=`nursery`).

## Steps

- [ ] Verify mv_maker + nursery genre produces identical output to standalone nursery_rhyme mode
- [ ] Add `nursery` to MVGenrePromptMap pointing to existing `storyboard_nursery_rhyme.txt`
- [ ] Add nursery option to genre dropdown
- [ ] Mark standalone nursery_rhyme radio button as deprecated (visual indicator)
- [ ] Eventually remove standalone nursery_rhyme radio button after user migration

## Context

See [plans/mv-maker-mode.md](../plans/mv-maker-mode.md) — Option 1 (Safe) migration strategy.
