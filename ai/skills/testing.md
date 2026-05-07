# Skill: Testing — Testing Approach

> This codebase has minimal automated tests. Here's how to verify changes.

## Current Test Coverage

Only 1 test file exists: `application/services/storyboard_parser_test.go`

```bash
# Run the only test
go test ./application/services/ -run TestParser -v
```

## Manual Testing Strategy

### Backend API Testing

```bash
# 1. Start the server
go run main.go

# 2. Test health endpoint
curl http://localhost:5678/health

# 3. Test API endpoints with curl or Postman
curl http://localhost:5678/api/v1/dramas
curl -X POST http://localhost:5678/api/v1/dramas -H "Content-Type: application/json" -d '{"title":"Test"}'
```

### Frontend Testing

1. Start dev server: `cd web && npm run dev`
2. Open `http://localhost:3012`
3. Navigate through features manually
4. Check browser console for errors
5. Check Network tab for API failures

### AI Integration Testing

1. Configure a test AI provider in Settings → AI Config
2. Create a test drama with simple script
3. Run storyboard generation — check JSON output
4. Run image generation — check image URLs resolve
5. Run video generation — check async completion

### FFmpeg Testing

```bash
# Verify FFmpeg works
ffmpeg -version

# Test a simple merge
ffmpeg -i input1.mp4 -i input2.mp4 -filter_complex concat=n=2:v=1:a=1 output.mp4
```

## Adding New Tests

If you need to add tests for new features:

```go
// application/services/xxx_test.go
package services

import "testing"

func TestXxx(t *testing.T) {
    // Test without database — use mock or in-memory SQLite
    // Test parsing logic, not full service flows
}
```

### What to Test (High Value)
- JSON parsing from AI responses
- Prompt template rendering
- Data transformation logic
- Edge cases in business rules

### What NOT to Test (Low Value Here)
- GORM CRUD operations
- Gin handler routing
- External API calls (mock these)

## Instrumentation

```bash
./scripts/log-skill.sh "testing" "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
```
