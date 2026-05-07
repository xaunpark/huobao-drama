# Infrastructure Router — Database, Docker, FFmpeg, Deployment

> Load this when working with database, Docker, deployment, or infrastructure.

## Retrieval Order

1. **This file** (you're reading it)
2. `ai/memory/infrastructure.md` — patterns & config details
3. `ai/memory/decisions.md` — why things were done this way

## Database (SQLite)

### Key Files
- `infrastructure/database/database.go` — Connection setup, WAL mode, auto-migrate
- `domain/models/*.go` — GORM model definitions (12 files)
- `migrations/init.sql` — Schema DDL (20KB, source of truth)
- `configs/config.yaml` — DB path: `./data/drama_generator.db`

### Patterns
- **ORM**: GORM v1.30 with SQLite driver (`modernc.org/sqlite` — pure Go, no CGO)
- **Mode**: WAL (Write-Ahead Logging) for concurrent reads
- **Migration**: Auto-migrate on startup via GORM
- **Soft deletes**: `gorm.DeletedAt` on all models
- **JSON fields**: `gorm.io/datatypes.JSON` for flexible data

### Common Issues
- `database is locked` — too many concurrent writes (WAL helps but doesn't eliminate)
- File permissions — SQLite needs write access to DB file AND directory (for `-wal`, `-journal`)
- Corrupt DB files — existing `drama_generator.db.corrupt` suggests this has happened

## Docker

### Key Files
- `Dockerfile` — Multi-stage build (frontend → backend → runtime)
- `docker-compose.yml` — Service definition with health checks
- `.env.example` — Mirror configuration for China

### Architecture
- **Stage 1**: Node 20 Alpine — build Vue frontend
- **Stage 2**: Go 1.23 Alpine — compile Go binary (CGO_ENABLED=0)
- **Stage 3**: Alpine latest — runtime with FFmpeg
- Port: 5678
- Data volume: `/app/data` (named volume `huobao-data`)
- Health check: `wget http://localhost:5678/health`

### Host Access
- `host.docker.internal` for accessing host services (Ollama, etc.)
- `extra_hosts: "host.docker.internal:host-gateway"` in compose

## FFmpeg

### Key File
- `infrastructure/external/ffmpeg/ffmpeg.go` (29KB)
- Called from `video_merge_service.go` and other services
- Not a library — shells out to `ffmpeg` binary

### Capabilities
- Video concatenation with crossfade transitions
- Audio extraction from video
- Grid image composition (2x2, 3x3)
- Text overlay (subtitles, watermarks)
- Video format conversion

## Storage

### Key File
- `infrastructure/storage/local_storage.go`
- Local filesystem storage at `./data/storage/`
- Served via static file handler at `/static`
- Base URL configurable: `storage.base_url` in config

### Resource Transfer
- `infrastructure/scheduler/resource_transfer_scheduler.go` — Cron job
- `application/services/resource_transfer_service.go` — Downloads external URLs to local
- Critical for preventing expired image/video URLs from providers

## Deployment Options

1. **Dev mode**: `go run main.go` + `cd web && npm run dev`
2. **Single binary**: Build frontend → embed in Go binary
3. **Docker Compose**: Recommended for production
4. **Traditional**: Binary + systemd service + Nginx reverse proxy

## Critical Warnings

- No database backups configured — manual only
- SQLite doesn't scale for high-write concurrency
- FFmpeg must be in PATH — runtime dependency
- Data directory must be writable by the running user
- Docker volumes for persistence — data lost without named volumes
