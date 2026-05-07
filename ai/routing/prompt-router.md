# Prompt Router — AI Prompt Template Context

> Load this when creating or modifying AI prompt templates.

## Retrieval Order

1. **This file** (you're reading it)
2. `ai/skills/prompt-engineering.md` — prompt writing patterns
3. `ai/systems/storyboard-system.md` — if working on storyboard prompts

## Prompt System Architecture

### Location: `application/prompts/`

34 template files organized by function:

#### Character & Scene Extraction
- `character_extraction.txt` — Extract characters from scripts
- `scene_extraction.txt` — Extract scenes/locations
- `prop_extraction.txt` — Extract props/objects

#### Image Generation Prompts
- `image_first_frame.txt` — First frame generation
- `image_key_frame.txt` — Key frame generation
- `image_last_frame.txt` — Last frame generation
- `image_action_sequence.txt` — 3x3 grid action sequence
- `image_action_sequence_rapid_cut.txt` — Rapid cut variant
- `image_style_distill.txt` — Style distillation for images

#### Storyboard Generation (Multiple Modes)
- `storyboard_story_breakdown.txt` — Standard mode
- `storyboard_visual_unit.txt` — Visual unit split
- `storyboard_visual_unit_format.txt` — Format instructions
- `storyboard_visual_unit_structured.txt` — Structured output
- `storyboard_nursery_rhyme.txt` — Nursery rhyme mode
- `storyboard_nursery_rhyme_format.txt` — Nursery format
- `storyboard_mv_gaming_horror.txt` — MV gaming horror (CG5-style)
- `storyboard_mv_cinematic_movie.txt` — MV cinematic mode
- `storyboard_narrative_director.txt` — Narrative MV director
- `storyboard_narrative_format.txt` — Narrative format
- `storyboard_narrative_planner.txt` — Narrative planning
- `storyboard_preserve_shots.txt` — Regenerate without losing data
- `storyboard_format_instructions.txt` — General format rules

#### Video Generation
- `video_constraint_prefixes.txt` — Video prompt prefixes
- `video_constraint_rapid_cut.txt` — Rapid cut constraints
- `video_distill_combined.txt` — Combined distillation
- `video_extraction.txt` — Video data extraction
- `video_style_distill.txt` — Style distillation
- `video_review_rubric.txt` — AI video review criteria

#### Script Generation
- `script_episode_generation.txt` — Episode script generation
- `script_outline_generation.txt` — Outline generation

#### Other
- `style_prompt.txt` — Global style prompt
- `rapid_cut_merge.txt` — Rapid cut merge logic
- `prompts.go` — Go embed directives

### I18n Layer: `application/services/prompt_i18n.go` (27KB)
- Bilingual support (zh/en) based on `app.language` config
- `GetPrompt(key string)` — returns localized prompt
- Falls back to English if Chinese not available
- Custom templates (user-defined in DB) override defaults

### Custom Templates: `domain/models/prompt_template.go`
- Users can create custom prompt templates via web UI
- Stored in database, override built-in prompts
- CRUD operations via `application/services/prompt_template_service.go`

### Channel Templates (docs/)
Large reference templates for specific YouTube channel aesthetics:
- `docs/cg5_template.md` (74KB) — CG5 gaming horror
- `docs/cg5_poppy_playtime_template.md` (57KB) — Poppy Playtime variant
- `docs/cocomelon_template.md` (66KB) — Cocomelon style
- `docs/super_simple_songs_template.md` (50KB) — SSS style
- `docs/hl_meow_meow_template.md` (79KB) — Meow Meow style
- And 8+ more specialized templates

## Prompt Writing Conventions

1. **Language**: Prompts are primarily in English (even for Chinese users)
2. **Format**: Plain text `.txt` files, no markdown
3. **Placeholders**: Use `{{.VariableName}}` Go template syntax (some use `%s` sprintf)
4. **Output format**: Most prompts request JSON output from AI
5. **System prompt vs User prompt**: System prompts set role/context, user prompts contain data

## Critical Rules

- **DO NOT** change prompt output format without updating the parsing logic in corresponding service
- Storyboard prompts output JSON → parsed by `storyboard_service.go`
- Image prompts output text → used directly as image generation input
- Video prompts output text → used as video generation motion/scene description
- Always test with multiple AI providers — different models interpret prompts differently

## When to Load Channel Templates

Only load docs templates when:
- User is working on a specific channel's aesthetic
- Creating a new channel template
- Debugging visual inconsistency issues
- These are 50-80KB each — never load all at once
