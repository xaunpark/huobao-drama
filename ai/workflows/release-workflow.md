# Workflow: Release

> Step-by-step procedure for preparing a release.

## Prerequisites
- All features complete and tested
- No blocking bugs

## Steps

### 1. Pre-Release Checks
- [ ] `go build ./...` compiles without errors
- [ ] `cd web && npm run build:check` passes
- [ ] All planned features tested manually
- [ ] No critical bugs in `todos/`
- [ ] `docs/solutions/` updated with new learnings

### 2. Version Bump
- Update `configs/config.yaml` → `app.version`
- Update `README.md` changelog section

### 3. Build
```bash
# Frontend
cd web && npm run build && cd ..

# Backend
go build -ldflags="-w -s" -o huobao-drama .
```

### 4. Docker Build
```bash
docker compose build
docker compose up -d
# Verify health: curl http://localhost:5678/health
```

### 5. Post-Release
- Tag release in git: `git tag v1.0.x`
- Run `/housekeeping` to archive completed plans/todos
- Update `README.md` changelog

## Rollback Plan
- Keep previous binary/Docker image
- SQLite DB backup: `cp data/drama_generator.db data/drama_generator.bak.db`
- Rollback: restore previous binary + backup DB
