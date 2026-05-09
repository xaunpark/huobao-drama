# Decisions — Key Technical Decisions

> Why things were done this way. Check here before re-debating decisions.

## D-001: SQLite as Primary Database
- **Date**: Project inception
- **Decision**: Use SQLite instead of PostgreSQL/MySQL
- **Rationale**: Single-binary deployment, zero configuration, easy backup (copy file)
- **Consequences**: No concurrent write scaling, must use WAL mode, file permissions critical
- **Status**: Active, no plans to change

## D-002: Pure Go SQLite Driver
- **Date**: v1.0.3 (2026-01-16)
- **Decision**: Switch from CGO SQLite to `modernc.org/sqlite`
- **Rationale**: Enable `CGO_ENABLED=0` for simple cross-compilation and Docker builds
- **Consequences**: Slightly slower queries, but eliminated all CGO build issues
- **Status**: Active

## D-003: DDD Layer Structure
- **Decision**: Use DDD-inspired 4-layer architecture without strict DDD patterns
- **Rationale**: Organize code by concern, not by DDD orthodoxy
- **Consequences**: No repository pattern (services use GORM directly), no value objects, no aggregates
- **Status**: Active — don't add repository pattern unless there's a concrete reason

## D-004: Manual Dependency Injection
- **Decision**: Wire all dependencies manually in `routes.go:SetupRouter()`
- **Rationale**: Simple, explicit, no magic, easy to trace dependencies
- **Consequences**: `SetupRouter()` is large and growing, but predictable
- **Status**: Active — don't add Wire, dig, or other DI frameworks

## D-005: AI Provider Configs in Database
- **Decision**: Store AI API keys and configs in database, not config files
- **Rationale**: Users manage multiple providers via web UI, can add/test/switch at runtime
- **Consequences**: Config not in version control, must backup DB to preserve configs
- **Status**: Active

## D-006: Prompt Template Layering
- **Decision**: Static prompts in `.txt` files → I18n layer → User override in DB
- **Rationale**: Balance between developer-maintained defaults and user customization
- **Consequences**: Complex resolution chain in `prompt_i18n.go`, but flexible
- **Status**: Active

## D-007: TailwindCSS v4 for Frontend
- **Decision**: Use TailwindCSS 4.1 with `@tailwindcss/postcss`
- **Rationale**: Latest version, utility-first CSS for rapid UI development
- **Consequences**: v4 has breaking changes from v3 (different plugin system, config format)
- **Status**: Active — do NOT mix v3 and v4 patterns

## D-008: Vue Composition API
- **Decision**: Use `<script setup>` Composition API exclusively
- **Rationale**: Better TypeScript support, cleaner component logic
- **Consequences**: All new components must use Composition API
- **Status**: Active

## D-009: Element Plus as UI Framework
- **Decision**: Use Element Plus as the sole UI component library
- **Rationale**: Comprehensive component set, Vue 3 native, Chinese community support
- **Consequences**: Don't introduce Vuetify, Ant Design, or other UI libraries
- **Status**: Active

## D-010: Resource Transfer Pattern
- **Decision**: Background cron downloads external URLs to local storage
- **Rationale**: AI provider image/video URLs expire (TTL 1-24 hours)
- **Consequences**: Requires scheduler, doubles storage, but ensures asset permanence
- **Status**: Active — critical for production reliability

## D-011: Storyboard Mode Dispatching
- **Decision**: Single `storyboard_service.go` dispatches to mode-specific services
- **Rationale**: Central entry point for all storyboard generation
- **Consequences**: Very large file (77KB), complex conditional logic
- **Status**: Active, but identified as tech debt for potential refactoring

## D-012: Channel Template Architecture
- **Decision**: Large markdown files in `docs/` as reference templates for AI style cloning
- **Rationale**: Detailed aesthetic specifications need rich formatting and examples
- **Consequences**: 50-80KB files that are reference docs, not runtime config
- **Status**: Active — new channel templates added regularly

## D-013: Explicit FlowTool Generation Mode Routing
- **Date**: 2026-05-08
- **Decision**: Frontend must explicitly send `generation_mode: 'i2v_s'` for First Frame video generation instead of relying on backend defaults
- **Rationale**: Backend defaults `generation_mode` to `"shot_i2v"` when not specified, which FlowTool maps to **R2V** (Reference-to-Video). For First Frame images, the correct FlowTool mode is **I2V_S** (Image-to-Video Start frame), which produces semantically different results (video starts exactly from the provided frame vs using it loosely as reference). The default mapping silently misrouted all First Frame shots to R2V.
- **Consequences**: 
  - Batch Action Studio conditionally sets `generation_mode` based on `generationMode` dropdown value
  - Manual Editor auto-detects `frame_type === 'first'` on selected image and sets `generation_mode: 'i2v_s'`
  - `video.ts` type union widened to include `'i2v_s'`
  - Backend Go code unchanged — `generation_mode` is a pass-through string field
- **Files changed**: `BatchGenerationDialog.vue`, `ProfessionalEditor.vue`, `web/src/types/video.ts`
- **Status**: Active
