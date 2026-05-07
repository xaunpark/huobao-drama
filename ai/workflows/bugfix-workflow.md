# Workflow: Bug Fix

> Step-by-step procedure for fixing bugs.

## Prerequisites
- Load `ai/routing/debugging-router.md`
- Load `ai/memory/risks.md` — check if this is a known fragile area

## Steps

### 1. Reproduce
- Identify exact steps to reproduce the bug
- Check server logs for error messages
- Check browser console for frontend errors

### 2. Locate
- Use the Layer Identification table in `ai/routing/debugging-router.md`
- Trace the data flow from HTTP request to response
- Identify the specific file and function

### 3. Root Cause
- Read the relevant code section (not entire file)
- Check `docs/solutions/` for similar past issues
- Identify the root cause, not just symptoms

### 4. Fix
- Make minimal surgical change
- Don't refactor unrelated code
- Preserve existing behavior for unaffected paths

### 5. Verify
- Reproduce the original bug steps — confirm fix
- Check no regression in related features
- `go build ./...` compiles
- Frontend builds if touched: `cd web && npm run build:check`

### 6. Document
- If novel bug: create entry in `docs/solutions/`
- If in known fragile area: update `ai/memory/risks.md`
- Run `/compound` if significant learning
