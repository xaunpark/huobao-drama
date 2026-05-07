# Skill: Prompt Engineering — Writing AI Prompts for This System

> Use when creating or modifying prompt templates in `application/prompts/`.

## Prompt Architecture

```
application/prompts/{domain}_{action}.txt     → Static template
application/services/prompt_i18n.go           → I18n resolution layer
domain/models/prompt_template.go              → User override (DB)
application/services/{domain}_service.go      → Prompt consumer (parsing)
```

## Writing Effective Prompts

### Structure
```
1. Role definition ("You are a professional storyboard director...")
2. Task description ("Analyze the following script and break it into shots...")
3. Input specification ("Script content:\n{{.ScriptContent}}")
4. Output format ("Return a JSON array with the following structure:")
5. Format example (concrete JSON example)
6. Constraints ("Rules to follow: ...")
7. Anti-examples ("Do NOT do: ...")
```

### Output Format Rules
- **Always request JSON** for structured data (storyboard, characters, etc.)
- **Provide example output** — AI follows examples more reliably than descriptions
- **Specify field types** — "duration: integer in seconds", not just "duration"
- **Handle edge cases** — "If no dialogue exists, set dialogue to null"

### Placeholder Patterns
```
Go template syntax:    {{.VariableName}}
Sprintf syntax:        %s (legacy, some prompts still use this)
```

### Multi-Language Support
- Prompts are primarily English (even for Chinese-language output)
- I18n layer in `prompt_i18n.go` handles Chinese variants
- Custom user templates stored in DB override both

## Common Failure Modes

### 1. AI Returns Markdown-Wrapped JSON
```
```json
{"shots": [...]}
```
```
**Fix**: Service code must strip markdown fences before JSON parsing.
Already handled in `storyboard_service.go` — check before adding new parsing.

### 2. AI Truncates Long Output
**Cause**: `max_tokens` too low for the requested output
**Fix**: Increase max_tokens, or break prompt into smaller chunks

### 3. AI Ignores Format Instructions
**Fix**: Add concrete examples, move format instructions to the END of the prompt,
use "IMPORTANT:" prefix for critical rules

### 4. Provider-Specific Behavior
- GPT-4 follows JSON format well
- Gemini may add extra text around JSON
- Local models (Ollama) may have lower quality adherence
**Fix**: Test with target provider, add more examples

## Validation Procedure

1. Test prompt with OpenAI GPT-4 (most reliable baseline)
2. Test with Gemini (different formatting tendencies)
3. Verify JSON output parses correctly in Go service
4. Check all expected fields are populated
5. Test with edge cases (short input, very long input, non-English input)

## Existing Prompt Templates Reference

See `ai/routing/prompt-router.md` for full listing of 34 prompt templates.

## Instrumentation

```bash
./scripts/log-skill.sh "prompt-engineering" "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
```
