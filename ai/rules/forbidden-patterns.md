# Rules: Forbidden Patterns — Anti-Patterns to Avoid

> Violations of these rules have caused bugs in the past.

## Architecture Violations

### ❌ DO NOT add business logic to handlers
Handlers parse input and format output. All logic goes in services.
```go
// BAD - handler doing business logic
func (h *XxxHandler) CreateXxx(c *gin.Context) {
    // ... complex data transformation here ...
}

// GOOD - handler delegates to service
func (h *XxxHandler) CreateXxx(c *gin.Context) {
    result, err := h.service.Create(input)
    c.JSON(200, result)
}
```

### ❌ DO NOT add a repository abstraction layer
The project uses GORM directly in services. This is deliberate (Decision D-003).

### ❌ DO NOT add a DI framework
Manual DI in `routes.go` is intentional (Decision D-004).

### ❌ DO NOT create new UI component libraries
Element Plus is the sole UI library (Decision D-009). No Vuetify, Ant Design, etc.

## Code Quality

### ❌ DO NOT delete Chinese comments
Many backend files have Chinese (中文) comments. Preserve them as-is.

### ❌ DO NOT modify prompt output format without updating parsing
If you change a prompt template's output JSON structure, you MUST update the
corresponding service's parsing logic. These are tightly coupled:
- `storyboard_*.txt` ↔ `storyboard_service.go` JSON parsing
- `character_extraction.txt` ↔ `character_library_service.go`

### ❌ DO NOT use Options API in Vue components
All Vue components use `<script setup lang="ts">` (Composition API only).

### ❌ DO NOT mix TailwindCSS v3 and v4 patterns
This project uses TailwindCSS v4.1 with `@tailwindcss/postcss`.

## Data Safety

### ❌ DO NOT remove GORM model fields
Removing fields from GORM models breaks auto-migration and existing data.
Only ADD fields, never remove. Use nullable types for new optional fields.

### ❌ DO NOT change table names
GORM's `TableName()` methods are the source of truth. Changing them breaks existing databases.

### ❌ DO NOT modify migrations/init.sql without understanding
This file is a reference schema. GORM auto-migrate is the actual migration mechanism.

## API Design

### ❌ DO NOT create routes outside routes.go
All routes are registered in `api/routes/routes.go:SetupRouter()`. No distributed route registration.

### ❌ DO NOT break existing API contracts
Frontend depends on specific response shapes. Changing response format breaks the UI.

## Infrastructure

### ❌ DO NOT enable CGO
The Go binary must compile with `CGO_ENABLED=0` for Docker compatibility (Decision D-002).

### ❌ DO NOT hardcode file paths
Use `config.yaml` values for storage paths, base URLs, etc.

### ❌ DO NOT assume FFmpeg is available
FFmpeg is a runtime dependency. Check for errors and handle gracefully.
