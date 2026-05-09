# Risks — Known Fragile Areas & Technical Debt

> Check this before touching these areas. Known sources of bugs and complexity.

## 🔴 Critical Risks

### 1. storyboard_service.go — God File (77KB)
- **Risk**: Too large, 6 modes with interleaved logic, hard to reason about
- **Impact**: Any change risks breaking other modes
- **Workaround**: Read specific functions only, test all modes after changes
- **Debt**: Should be split into mode-specific dispatchers (partially done with nursery/mv/narrative services)

### 2. SQLite Concurrent Write Contention
- **Risk**: `database is locked` errors under write-heavy load
- **Impact**: API requests fail intermittently
- **Workaround**: WAL mode, minimize write transactions, retry logic
- **Evidence**: `data/drama_generator.db.corrupt` file exists — corruption has occurred

### 3. External URL Expiration
- **Risk**: AI provider image/video URLs expire (1-24 hour TTL)
- **Impact**: Lost assets if resource transfer fails or is delayed
- **Workaround**: Resource transfer scheduler runs periodically
- **Fragile point**: `infrastructure/scheduler/resource_transfer_scheduler.go`

### 4. No Automated Tests
- **Risk**: Only 1 test file (`storyboard_parser_test.go`), no CI/CD testing
- **Impact**: Regressions go undetected until manual testing
- **Mitigation**: Manual testing, careful code review

## 🟠 High Risks

### 5. prompt_i18n.go — Complex Resolution Chain (27KB)
- **Risk**: 3-layer prompt resolution (static → i18n → custom template override)
- **Impact**: Wrong prompt used if resolution chain breaks
- **Fragile**: Template service dependency injection order matters

### 6. AI Response Parsing
- **Risk**: AI models return unpredictable formats (non-JSON, markdown-wrapped JSON, truncated)
- **Impact**: Storyboard generation fails silently or with cryptic errors
- **Locations**: `storyboard_service.go`, `character_library_service.go`
- **Pattern**: Always clean AI response before JSON parsing (strip markdown fences, fix trailing commas)

### 7. FFmpeg Command Construction
- **Risk**: Special characters in filenames break FFmpeg commands
- **Impact**: Video merge fails
- **Evidence**: Multiple `fix_*.py` and `fix_*.go` scripts at root level suggest historical data fixing
- **Location**: `infrastructure/external/ffmpeg/ffmpeg.go` (29KB)

### 8. Frontend-Backend Type Sync
- **Risk**: TypeScript types manually mirror Go models, no code generation
- **Impact**: Type mismatches cause runtime errors
- **Locations**: `web/src/types/` vs `domain/models/`

### 9. Video Generation Mode Silent Misrouting
- **Risk**: If frontend doesn't explicitly send `generation_mode`, backend defaults to `"shot_i2v"` which maps to FlowTool's **R2V** mode. First Frame images should use **I2V_S** instead.
- **Impact**: Videos generated with wrong FlowTool endpoint — R2V treats the image as loose reference (style/character guide), I2V_S treats it as the exact starting frame. Wrong mode = wrong video output, no error reported.
- **Pattern**: Any new frame type or generation mode must verify the full chain: frontend `generation_mode` → backend `GenerationMode` field → `flowtool_video_client.go` mode resolution → FlowTool API `mode` parameter.
- **Fix applied**: 2026-05-08, D-013 in `decisions.md`
- **Watch for**: Adding Last Frame or Panel batch modes would need similar explicit routing.

## 🟡 Medium Risks

### 10. Storyboard Model Field Explosion
- **Risk**: `Storyboard` model has 50+ fields (voice-over, nursery, MV, rapid cut...)
- **Impact**: Large DB rows, complex queries, migration risk
- **Evidence**: Fields added incrementally per mode — `domain/models/drama.go`

### 11. No Authentication
- **Risk**: No user authentication or authorization
- **Impact**: Anyone with network access can control the system
- **Design**: Intended for single-user/local use, but exposed on 0.0.0.0

### 12. Debug Print Statements
- **Risk**: `fmt.Printf` used extensively in AI clients instead of structured logging
- **Impact**: Noisy logs, potential credential leakage in log output
- **Locations**: `pkg/ai/openai_client.go`, `pkg/ai/gemini_client.go`

### 13. Root-Level Utility Scripts
- **Risk**: 10+ Go and Python utility scripts at project root (`check_db.go`, `fix_dates.py`, etc.)
- **Impact**: Confusion about what's part of the application vs one-off scripts
- **Evidence**: `check_db.go`, `check_db2.go`, `dump.go`, `fix_dates.py`, `fix_empty_dates.py`, etc.

## Known Data Issues

- `data/drama_generator.db.corrupt` — SQLite corruption has occurred before
- `data/drama_generator.bak.db` (64MB) — Manual backup exists
- `storyboards.json` (466KB) at root — appears to be a data dump, not application code

## Technical Debt Backlog

From `todos/`:
- `deprecate-nursery-rhyme-mode.md` — Nursery mode may be deprecated
- `pending-p2-optimize-distill-nplus1.md` — N+1 query in style distillation
- `pending-p3-remove-dead-style-prompt.md` — Unused style prompt code
- `mv-maker-genre-profiles.md` — Genre profiles not fully implemented
