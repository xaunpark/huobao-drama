# Skill: Refactoring — Safe Refactoring Patterns

> Use when restructuring code without changing behavior.

## Golden Rule

Every refactoring must:
1. Leave existing behavior unchanged
2. Pass the same manual tests as before
3. Not increase the public API surface

## Safe Refactoring Patterns for This Codebase

### Pattern 1: Extract Service Method
When a service function is too large (common in `storyboard_service.go`):
1. Identify a self-contained block of logic
2. Extract to a private method on the same struct
3. Pass only required parameters
4. Keep in the SAME file (don't create new files without reason)

### Pattern 2: Split Handler Actions
When a handler has too many methods:
1. Create a new handler file for the extracted group
2. Register new handler in `routes.go`
3. Keep the original handler for existing routes

### Pattern 3: Extract Prompt Constants
When prompt text is inline in service code:
1. Create new `.txt` file in `application/prompts/`
2. Add Go embed directive in `prompts/prompts.go`
3. Reference via `prompt_i18n.go`

### Pattern 4: Consolidate Duplicate AI Calls
When similar AI patterns repeat:
1. Create helper function in the service
2. Parameterize the differences
3. Don't create a "generic AI caller" — keep it domain-specific

## Dangerous Refactorings (Avoid)

- ❌ **Don't add a repository layer** — project deliberately uses GORM directly (Decision D-003)
- ❌ **Don't introduce DI framework** — manual DI in routes.go is intentional (Decision D-004)
- ❌ **Don't split `routes.go`** — single route file is a feature, not a bug
- ❌ **Don't rename Chinese-commented code** — preserve existing comment language
- ❌ **Don't restructure `domain/models/`** — GORM auto-migrate depends on stable model structure

## Refactoring Checklist

- [ ] Identified specific, measurable improvement
- [ ] No behavior change (same inputs → same outputs)
- [ ] Existing imports still work (no broken references)
- [ ] GORM model changes checked against `migrations/init.sql`
- [ ] Changes compiled (`go build ./...`)
- [ ] Frontend types still match backend models

## Instrumentation

```bash
./scripts/log-skill.sh "refactor" "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
```
