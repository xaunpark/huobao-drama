# Infrastructure — Infrastructure Patterns & Configuration

> Infrastructure details for database, Docker, deployment, and runtime.

## Configuration

### Config File: `configs/config.yaml`
```yaml
app:
  name: "Huobao Drama API"
  version: "1.0.0"
  debug: true                    # true=development, false=production
  language: "zh"                 # zh|en — affects prompt language

server:
  port: 5678
  host: "0.0.0.0"
  cors_origins: ["http://localhost:3012"]
  read_timeout: 600              # 10 minutes (for long AI operations)
  write_timeout: 600

database:
  type: "sqlite"                 # Only sqlite supported currently
  path: "./data/drama_generator.db"
  max_idle: 10
  max_open: 100

storage:
  type: "local"
  local_path: "./data/storage"
  base_url: "http://localhost:5678/static"

ai:
  default_text_provider: "openai"
  default_image_provider: "openai"
  default_video_provider: "doubao"
```

### Config Loading
- `pkg/config/` — Uses Viper
- Reads from `configs/config.yaml`
- No environment variable override (config file only)

## Database

### SQLite Configuration
- **Driver**: `modernc.org/sqlite` (pure Go, no CGO)
- **Mode**: WAL (Write-Ahead Logging) for concurrent reads
- **File**: `./data/drama_generator.db`
- **Auto-migrate**: GORM auto-migration on startup
- **Soft deletes**: All models use `gorm.DeletedAt`

### Connection Pool
```
max_idle: 10
max_open: 100
```

### Schema Source of Truth
- Primary: `domain/models/*.go` (GORM auto-migrate)
- Reference: `migrations/init.sql` (20KB DDL)
- Migration: `migrations/20260126_add_local_path.sql`

## Docker

### Multi-Stage Build
```
Stage 1: node:20-alpine     → Build Vue frontend (npm install + build)
Stage 2: golang:1.23-alpine → Compile Go binary (CGO_ENABLED=0)
Stage 3: alpine:latest      → Runtime (binary + ffmpeg + frontend dist)
```

### Runtime Image Contents
```
/app/huobao-drama           → Go binary
/app/migrate                 → Migration binary
/app/web/dist/               → Frontend static files
/app/configs/config.yaml     → Configuration
/app/migrations/             → SQL migrations
/app/data/                   → Data volume mount point
```

### Networking
- Port: 5678
- Health check: `GET /health`
- Host access: `host.docker.internal` (for Ollama, etc.)

## File Storage

### Layout
```
./data/
  drama_generator.db         → SQLite database
  storage/                   → Generated assets
    images/                  → AI-generated images (downloaded from providers)
    videos/                  → AI-generated videos
    uploads/                 → User-uploaded files
    merged/                  → FFmpeg merged outputs
```

### Static Serving
- Route: `/static/*` → serves `./data/storage/`
- Base URL: `storage.base_url` in config
- Files persist across restarts (Docker volume or local disk)

## Serving Architecture

### Development Mode
```
Browser → localhost:3012 (Vite dev server)
          ├── /api/* → proxy to localhost:5678 (Go backend)
          └── /* → Vue SPA hot reload

Go backend on :5678
          ├── /api/v1/* → API endpoints
          ├── /static/* → Local storage files
          └── /health → Health check
```

### Production Mode (Single Binary)
```
Browser → localhost:5678 (Go backend)
          ├── /api/v1/* → API endpoints
          ├── /static/* → Local storage files
          ├── /assets/* → Frontend JS/CSS bundles
          ├── /health → Health check
          └── /* → SPA fallback (index.html)
```

## Timeout Configuration

| Component | Timeout | Reason |
|-----------|---------|--------|
| HTTP server read/write | 10 min | Long AI generation operations |
| OpenAI HTTP client | 30 min | Very long chat completions |
| FFmpeg operations | None (process) | Large video processing |
| Resource transfer | Cron interval | Background download |
