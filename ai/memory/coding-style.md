# Coding Style — Code Style Conventions

> Inferred from actual codebase patterns. Follow these, don't introduce new styles.

## Go Backend

### Naming
- **Files**: `snake_case.go` (e.g., `drama_service.go`, `video_generation_service.go`)
- **Types**: `PascalCase` (e.g., `DramaHandler`, `ImageGenerationService`)
- **Constructors**: `NewXxxHandler()`, `NewXxxService()`
- **Methods**: `PascalCase` for exported, `camelCase` for private
- **Variables**: `camelCase` (e.g., `dramaHandler`, `aiService`)

### Formatting
- Standard `gofmt`
- Tabs for indentation
- Chinese comments throughout — preserved, not translated

### Error Handling
- Return errors up the call stack: `return nil, fmt.Errorf("context: %w", err)`
- `fmt.Printf` for debug logging in AI clients (not structured Zap)
- `log.Fatal` / `log.Fatalw` for unrecoverable errors at startup
- HTTP errors: `c.JSON(statusCode, gin.H{"error": message})`

### Import Style
- Standard lib → third-party → internal packages
- Aliased imports in routes.go: `handlers2`, `middlewares2`, `services2` (legacy pattern — match it)

### Comments
- Mix of Chinese (中文) and English comments
- Chinese comments for business logic explanations
- English for technical/API comments
- Preserve whatever language existing comments use

## Vue Frontend

### Component Style
- **SFC**: `<script setup lang="ts">` (Composition API only)
- **No Options API** — all new components use `<script setup>`
- **Element Plus components**: `el-` prefixed (e.g., `el-table`, `el-dialog`)

### File Naming
- **Components**: `PascalCase.vue` (e.g., `LanguageSwitcher.vue`)
- **API modules**: `kebab-case.ts` (e.g., `character-library.ts`)
- **Types**: `kebab-case.ts` matching domain (e.g., `drama.ts`, `video.ts`)
- **Stores**: `camelCase.ts` (e.g., `episode.ts`)

### TypeScript
- Strict mode enabled
- Interfaces for type definitions (not `type` aliases for objects)
- `export function` pattern in API modules, not classes

### CSS
- TailwindCSS v4 utility classes in templates
- `sass-embedded` for SCSS when needed
- No CSS-in-JS

## Prompt Templates

### File Format
- Plain `.txt` files in `application/prompts/`
- No markdown formatting inside prompts
- Go template syntax: `{{.VariableName}}` for dynamic values
- Some use `%s` sprintf-style placeholders

### Prompt Structure
- System prompt sets AI role and output format
- User prompt provides the specific data/content
- Most request JSON output format from AI
- Include example outputs in prompts for consistency

## JSON API Convention

- All responses wrapped in: `{"data": ...}` for success, `{"error": "..."}` for failure
- Pagination: `?page=1&page_size=20`
- IDs as path parameters: `/api/v1/resource/:id`
- Nested resources: `/api/v1/episodes/:episode_id/storyboards`
