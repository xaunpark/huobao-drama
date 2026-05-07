# Frontend Router — Vue 3 SPA Context

> Load this when modifying any Vue frontend code.

## Retrieval Order

1. **This file** (you're reading it)
2. `ai/rules/naming.md` — component naming
3. `ai/memory/coding-style.md` — frontend conventions

## Architecture Quick Ref

```
web/src/
  api/              → 14 Axios API client modules (one per backend resource)
  views/            → 8 view groups (page-level components)
    dashboard/      → Main dashboard
    drama/          → Drama CRUD views
    editor/         → Visual editor views
    generation/     → Image & video generation UI
    script/         → Script editing
    settings/       → App settings
    storyboard/     → Storyboard editor
    workflow/       → Workflow management
  components/       → Shared reusable components
    common/         → Generic UI components
    editor/         → Editor-specific components
    LanguageSwitcher.vue → i18n switcher
  composables/      → Vue composables (hooks)
  stores/           → Pinia stores (episode.ts is the main store)
  types/            → 10 TypeScript type definition files
  locales/          → i18n translation files
  router/           → Vue Router config
  utils/            → Utility functions
  assets/           → Static assets
```

## Tech Stack

| Tool | Version | Purpose |
|------|---------|---------|
| Vue | 3.4+ | SPA framework (Composition API) |
| TypeScript | 5+ | Type safety |
| Vite | 5 | Build tool + dev server |
| Element Plus | 2.5+ | UI component library |
| TailwindCSS | 4.1 | Utility-first CSS |
| Pinia | 2.1+ | State management |
| Vue Router | 4.2+ | Client-side routing |
| Axios | 1.6+ | HTTP client |
| vue-i18n | 9.14+ | Internationalization |
| dayjs | 1.11+ | Date formatting |
| jszip | 3.10+ | ZIP file handling |
| cropperjs | 2.1+ | Image cropping |

## Key Patterns

### API Client Pattern
- One file per resource in `web/src/api/`
- Uses Axios with base URL proxy to `:5678/api/v1`
- Pattern: `export function getXxx()`, `export function createXxx(data)`, etc.

### Type Pattern
- Types defined in `web/src/types/` — match Go domain models
- One file per domain area: `drama.ts`, `video.ts`, `image.ts`, etc.

### View Pattern
- Page components in `web/src/views/{area}/`
- Use Composition API (`<script setup lang="ts">`)
- Heavy use of Element Plus components (`el-table`, `el-dialog`, `el-form`, etc.)

### Store Pattern
- Pinia stores in `web/src/stores/`
- Currently only `episode.ts` — handles episode state management

## Build & Dev

```bash
# Dev server (proxies API to :5678)
cd web && npm run dev       # Starts on :3012

# Production build
cd web && npm run build     # Output: web/dist/

# Type check + build
cd web && npm run build:check
```

## Vite Config

- Proxy: `/api` → `http://localhost:5678`
- Dev port: `3012`
- Config: `web/vite.config.ts`

## When You Need Deeper Context

| Task | Load |
|------|------|
| Adding new page/view | Existing view as template + `web/src/router/` |
| Adding API client | `web/src/api/` existing file as template |
| Adding types | `web/src/types/` matching Go models |
| Working with state | `web/src/stores/episode.ts` |
| i18n changes | `web/src/locales/` |

## Critical Warnings

- TailwindCSS v4.1 uses `@tailwindcss/postcss` — NOT the v3 plugin
- Element Plus is the ONLY UI library — don't introduce alternatives
- Frontend has NO unit tests — visual testing only
- Some `.vue` files are very large — read targeted sections
