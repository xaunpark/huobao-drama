# Dependency Map — External Dependencies

> Lists all external dependencies and their roles. Use for upgrade planning and vulnerability assessment.

## Backend (Go 1.23)

### Core Dependencies

| Dependency | Version | Purpose | Impact |
|-----------|---------|---------|--------|
| `gin-gonic/gin` | 1.9.1 | HTTP web framework | Critical — all API routes |
| `gorm.io/gorm` | 1.30.0 | ORM framework | Critical — all DB operations |
| `gorm.io/driver/sqlite` | 1.6.0 | SQLite GORM driver | Critical — DB connection |
| `modernc.org/sqlite` | 1.34.4 | Pure Go SQLite (no CGO) | Critical — enables cross-compilation |
| `spf13/viper` | 1.17.0 | Configuration management | High — config loading |
| `go.uber.org/zap` | 1.26.0 | Structured logging | High — all logging |
| `google/uuid` | 1.6.0 | UUID generation | Medium — entity IDs |
| `robfig/cron/v3` | 3.0.1 | Cron scheduling | Medium — resource transfer scheduler |
| `gorm.io/datatypes` | 1.2.0 | JSON column type for GORM | Medium — flexible data storage |

### Also Available (indirect/unused potential)

| Dependency | Purpose |
|-----------|---------|
| `gorm.io/driver/mysql` | MySQL support (available but not used) |
| `gorm.io/driver/postgres` | PostgreSQL support (available but not used) |

### Runtime Dependencies

| Tool | Purpose | Required |
|------|---------|----------|
| **FFmpeg** | Video processing, merging, audio extraction | YES — runtime |
| **SQLite** | Database (via pure Go driver) | Built-in |

## Frontend (Node 18+)

### Production Dependencies

| Dependency | Version | Purpose |
|-----------|---------|---------|
| `vue` | 3.4+ | SPA framework |
| `vue-router` | 4.2+ | Client-side routing |
| `pinia` | 2.1+ | State management |
| `element-plus` | 2.5+ | UI component library |
| `@element-plus/icons-vue` | 2.3+ | Icon set |
| `axios` | 1.6+ | HTTP client |
| `vue-i18n` | 9.14+ | Internationalization |
| `dayjs` | 1.11+ | Date formatting |
| `lodash-es` | 4.17+ | Utility functions |
| `jszip` | 3.10+ | ZIP file handling |
| `cropperjs` | 2.1+ | Image cropping UI |
| `@ffmpeg/ffmpeg` | 0.12+ | Browser-side FFmpeg (WASM) |
| `@ffmpeg/util` | 0.12+ | FFmpeg WASM utilities |

### Dev Dependencies

| Dependency | Version | Purpose |
|-----------|---------|---------|
| `vite` | 5+ | Build tool + dev server |
| `typescript` | 5.3+ | Type checking |
| `@vitejs/plugin-vue` | 5+ | Vue SFC compilation |
| `tailwindcss` | 4.1+ | CSS framework |
| `@tailwindcss/postcss` | 4.1+ | TailwindCSS PostCSS plugin (v4 style) |
| `postcss` | 8.4+ | CSS processing |
| `autoprefixer` | 10.4+ | CSS vendor prefixes |
| `sass-embedded` | 1.97+ | SCSS compilation |
| `vue-tsc` | 2.2+ | Vue TypeScript checking |

## Infrastructure Dependencies

| Tool | Version | Purpose |
|------|---------|---------|
| Docker | Any | Containerization |
| Docker Compose | v2 | Service orchestration |
| Nginx | Any | Reverse proxy (production) |
| Node.js | 20 | Frontend build (in Docker) |
| Alpine Linux | Latest | Runtime base image |

## External AI Service Dependencies

| Provider | Services Used | Config Location |
|----------|--------------|-----------------|
| OpenAI | Text (GPT), Image (DALL-E), Video (Sora) | DB `ai_configs` table |
| Google Gemini | Text, Image (Imagen) | DB `ai_configs` table |
| Doubao (Volcengine Ark) | Video generation | DB `ai_configs` table |
| MiniMax | Video generation | DB `ai_configs` table |
| Volcengine | Image generation | DB `ai_configs` table |
| FlowTool | Image + Video (aggregation) | DB `ai_configs` table |
| Chatfire | Video generation | DB `ai_configs` table |
| Ollama (local) | Text via OpenAI-compatible API | DB `ai_configs` table |

> **Note**: All AI provider configs are user-managed via the web UI and stored in the database, NOT in config files.
