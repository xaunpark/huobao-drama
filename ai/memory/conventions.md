# Conventions — Naming & Structural Conventions

> Follow these conventions when adding new code. Derived from existing patterns.

## Adding a New API Resource

1. Create domain model in `domain/models/{resource}.go`
2. Create handler in `api/handlers/{resource}.go`
3. Create service in `application/services/{resource}_service.go`
4. Register routes in `api/routes/routes.go:SetupRouter()`
5. Create frontend API client in `web/src/api/{resource}.ts`
6. Create frontend types in `web/src/types/{resource}.ts`

## Adding a New AI Provider

1. Implement the interface in `pkg/ai/` (text), `pkg/image/` (image), or `pkg/video/` (video)
2. Register in the factory function of the respective client.go
3. Add provider name to `domain/models/ai_config.go` provider enum
4. Test via Settings → AI Config → Test Connection in web UI

## Adding a New Storyboard Mode

1. Create mode-specific service in `application/services/storyboard_{mode}_service.go`
2. Create prompt template in `application/prompts/storyboard_{mode}.txt`
3. Create format template in `application/prompts/storyboard_{mode}_format.txt`
4. Add dispatch case in `application/services/storyboard_service.go`
5. Add UI option in frontend storyboard view
6. Create implementation plan in `plans/{mode}-mode.md`

## Naming Conventions

### Database Tables
- Plural snake_case: `dramas`, `episodes`, `storyboards`, `characters`
- Junction tables: `storyboard_characters`, `storyboard_props`, `episode_characters`

### API Endpoints
- RESTful: `GET /api/v1/{resource}` (list), `POST /api/v1/{resource}` (create)
- Nested: `POST /api/v1/episodes/:episode_id/storyboards` (generate for episode)
- Action: `POST /api/v1/videos/:id/upscale` (verb suffix for non-CRUD)

### Go Packages
- Flat within each layer — no deep nesting
- `handlers`, `services`, `models` (not `handler`, `service`, `model`)

### Frontend Routes
- Kebab-case paths: `/drama-list`, `/storyboard-edit`
- Nested: `/drama/:id/episode/:episodeId`

## Structural Conventions

### Service Dependencies
- Services receive `*gorm.DB` directly — no repository abstraction
- Services receive `*logger.Logger` for logging
- Services receive `*config.Config` when needing configuration
- Services receive other services when needing cross-domain logic

### Handler Dependencies
- Handlers receive services, config, logger, db
- Handlers do NOT access other handlers
- Handlers do input validation → service call → response formatting

### Model Fields
- Use pointers for optional fields: `*string`, `*int`, `*bool`
- Use `datatypes.JSON` for flexible/nested data
- Always include `CreatedAt`, `UpdatedAt`, `DeletedAt`
- Runtime fields with `gorm:"-"` tag

### Prompt Template Naming
- `{domain}_{action}.txt` — e.g., `storyboard_story_breakdown.txt`
- Variant: `{domain}_{action}_{variant}.txt` — e.g., `video_constraint_rapid_cut.txt`
- Format: `{domain}_{mode}_format.txt` — output format instructions
