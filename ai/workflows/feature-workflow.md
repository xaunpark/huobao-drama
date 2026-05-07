# Workflow: Feature Implementation

> Step-by-step procedure for implementing new features.

## Prerequisites
- Load `ai/routing/backend-router.md` or `ai/routing/frontend-router.md`
- Check `plans/` for existing plans
- Check `docs/solutions/` for prior knowledge

## Steps

### 1. Plan
- Create implementation plan in `plans/{feature-name}.md`
- Define scope, affected files, verification criteria
- Check `ai/memory/decisions.md` — don't re-debate settled decisions

### 2. Backend Implementation
1. **Model**: Add/modify domain model in `domain/models/`
2. **Service**: Implement business logic in `application/services/`
3. **Handler**: Create HTTP handler in `api/handlers/`
4. **Routes**: Register in `api/routes/routes.go`
5. **Prompts**: Add prompt templates in `application/prompts/` if AI-related
6. **Verify**: `go build ./...` compiles without errors

### 3. Frontend Implementation
1. **Types**: Add TypeScript types in `web/src/types/`
2. **API client**: Add API module in `web/src/api/`
3. **View**: Create page component in `web/src/views/`
4. **Route**: Register in `web/src/router/`
5. **Verify**: `npm run build:check` passes

### 4. Integration Test
1. Start backend: `go run main.go`
2. Start frontend: `cd web && npm run dev`
3. Test feature end-to-end via web UI
4. Check browser console for errors
5. Check server logs for errors

### 5. Document
- Update `plans/{feature-name}.md` with completion status
- Create `docs/solutions/` entry if solving a novel problem
- Update `ai/indexes/feature-map.md` if adding new feature
