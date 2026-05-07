# System: Frontend — Vue SPA Architecture

> Deep documentation of the frontend subsystem.

## Tech Stack
- Vue 3.4+ (Composition API with `<script setup>`)
- TypeScript 5+
- Vite 5 (build + dev server)
- Element Plus 2.5+ (UI components)
- TailwindCSS 4.1 (utility CSS)
- Pinia 2.1+ (state management)
- Vue Router 4 (routing)
- Axios (HTTP client)
- vue-i18n (internationalization)

## Page Structure

```
web/src/views/
  dashboard/       → Main dashboard with drama stats
  drama/           → Drama CRUD (list, create, edit, detail)
  editor/          → Visual drama editor (timeline, storyboard)
  generation/      → Image & video generation views
    ImageGeneration.vue (12KB)
    VideoGeneration.vue (12KB)
    components/    → Generation-specific components
  script/          → Script editor and generation
  settings/        → AI config, language, prompt templates
  storyboard/      → Storyboard editor
    StoryboardEdit.vue
  workflow/        → Workflow management
```

## State Management (Pinia)

Single store: `web/src/stores/episode.ts` (7KB)
- Manages current episode state
- Tracks storyboard generation status
- Handles batch operation progress

## API Client Layer

14 client modules in `web/src/api/`:
```
ai.ts              → AI config operations
asset.ts           → Asset library operations
audio.ts           → Audio extraction
character-library.ts → Character library CRUD
drama.ts           → Drama CRUD
frame.ts           → Frame prompt operations
generation.ts      → Generation triggers
image.ts           → Image generation
prompt-template.ts → Prompt template management
prop.ts            → Prop management
settings.ts        → Settings
task.ts            → Task status polling
video.ts           → Video generation
videoMerge.ts      → Video merge operations
```

## Type Definitions

10 type files in `web/src/types/`:
```
ai.ts, asset.ts, drama.ts, generation.ts, image.ts,
prompt-template.ts, prop.ts, timeline.ts, user.ts, video.ts
```

Types manually mirror Go domain models. No code generation.

## Routing

Vue Router in `web/src/router/`:
- Hash or History mode
- Lazy-loaded route components
- Nested routes for drama → episode → storyboard hierarchy

## i18n

`web/src/locales/` — Chinese (zh) and English (en) translation files
Language switcher: `web/src/components/LanguageSwitcher.vue`

## Critical Patterns

1. **Element Plus everywhere** — all UI uses `el-*` components
2. **Composition API only** — no Options API in new code
3. **Axios with proxy** — dev mode proxies `/api` to `:5678`
4. **No unit tests** — visual testing only
5. **TailwindCSS v4** — uses `@tailwindcss/postcss`, NOT v3 plugin

## Build

```bash
npm run dev         # Dev server on :3012
npm run build       # Production build → web/dist/
npm run build:check # TypeScript check + build
npm run build:skip  # Build without type check
```
