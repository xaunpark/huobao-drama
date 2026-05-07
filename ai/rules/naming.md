# Rules: Naming — Naming Conventions

> Follow these naming conventions for all new code.

## Go Backend

| Element | Convention | Example |
|---------|-----------|---------|
| File | `snake_case.go` | `video_generation_service.go` |
| Package | lowercase, single word | `handlers`, `services`, `models` |
| Exported Type | `PascalCase` | `DramaHandler`, `AIService` |
| Constructor | `NewXxxType()` | `NewDramaHandler()` |
| Exported Method | `PascalCase` verb | `CreateDrama()`, `GetStoryboards()` |
| Private method | `camelCase` | `buildPrompt()`, `parseResponse()` |
| Variable | `camelCase` | `dramaService`, `promptI18n` |
| Constant | `PascalCase` or `ALL_CAPS` | `DefaultTimeout`, `MAX_RETRIES` |
| Table name | Plural `snake_case` | `dramas`, `storyboards`, `video_generations` |
| Junction table | `{table1}_{table2}` | `storyboard_characters`, `episode_characters` |

## Vue Frontend

| Element | Convention | Example |
|---------|-----------|---------|
| Component file | `PascalCase.vue` | `ImageGeneration.vue`, `LanguageSwitcher.vue` |
| API client file | `kebab-case.ts` | `character-library.ts`, `prompt-template.ts` |
| Type file | `kebab-case.ts` | `drama.ts`, `video.ts` |
| Store file | `camelCase.ts` | `episode.ts` |
| Composable | `use{Name}.ts` | `useEpisode.ts` |
| View directory | `kebab-case/` | `drama/`, `storyboard/` |

## API Endpoints

| Pattern | Convention | Example |
|---------|-----------|---------|
| List | `GET /api/v1/{resources}` | `GET /api/v1/dramas` |
| Create | `POST /api/v1/{resources}` | `POST /api/v1/dramas` |
| Get | `GET /api/v1/{resources}/:id` | `GET /api/v1/dramas/1` |
| Update | `PUT /api/v1/{resources}/:id` | `PUT /api/v1/dramas/1` |
| Delete | `DELETE /api/v1/{resources}/:id` | `DELETE /api/v1/dramas/1` |
| Nested | `POST /api/v1/{parent}/:id/{child}` | `POST /api/v1/episodes/1/storyboards` |
| Action | `POST /api/v1/{resources}/:id/{verb}` | `POST /api/v1/videos/1/upscale` |

## Prompt Templates

| Pattern | Convention | Example |
|---------|-----------|---------|
| Main prompt | `{domain}_{action}.txt` | `storyboard_story_breakdown.txt` |
| Mode variant | `{domain}_{action}_{mode}.txt` | `video_constraint_rapid_cut.txt` |
| Format spec | `{domain}_{mode}_format.txt` | `storyboard_nursery_rhyme_format.txt` |

## Channel Templates (docs/)

| Pattern | Convention | Example |
|---------|-----------|---------|
| Template | `{channel_name}_template.md` | `cg5_template.md`, `cocomelon_template.md` |
| Variant | `{channel}_{variant}_template.md` | `cg5_poppy_playtime_template.md` |
