# Workflow: Prompt Template Creation/Modification

> Step-by-step procedure for creating or modifying AI prompt templates.

## Prerequisites
- Load `ai/routing/prompt-router.md`
- Load `ai/skills/prompt-engineering.md`

## Steps

### 1. Understand Requirements
- What AI output is needed? (JSON structure, text format?)
- Which service will consume this prompt's output?
- Which AI providers will use this prompt?

### 2. Create Prompt File
```
application/prompts/{domain}_{action}.txt
```
Follow the structure in `ai/skills/prompt-engineering.md`.

### 3. Update Prompt I18n
If adding a new prompt key:
1. Add to `application/services/prompt_i18n.go`
2. Add both English and Chinese variants if applicable
3. Register the prompt key for custom template override

### 4. Update Service Consumer
The service that calls this prompt must:
1. Load the prompt via `prompt_i18n.GetPrompt(key)`
2. Fill template placeholders
3. Parse AI response correctly
4. Handle parsing failures gracefully

### 5. Test
1. Test with OpenAI GPT-4 (baseline)
2. Test with Gemini (different formatting)
3. Verify JSON parsing in service
4. Test edge cases (short input, long input)

### 6. Document
- Add to `ai/routing/prompt-router.md` prompt listing
- If novel pattern: add to `docs/solutions/`
