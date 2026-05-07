# Rules: API — API Design Rules

> Follow when adding or modifying API endpoints.

## Endpoint Design

1. **Base path**: All API endpoints under `/api/v1/`
2. **Resource naming**: Plural nouns (`/dramas`, `/videos`, `/storyboards`)
3. **Nesting**: One level deep max (`/episodes/:id/storyboards`)
4. **Actions**: `POST /resource/:id/verb` for non-CRUD operations

## Request/Response Format

### Success Response
```json
{
  "data": { ... },
  "total": 100,       // for lists
  "page": 1,          // for paginated
  "page_size": 20
}
```

### Error Response
```json
{
  "error": "Human-readable error message"
}
```

### Status Codes
- `200` — Success (GET, PUT, DELETE)
- `201` — Created (POST that creates a resource)
- `400` — Bad request (invalid input)
- `404` — Not found
- `500` — Server error

## Input Validation

- Validate in handler using Gin binding (`c.ShouldBindJSON()`)
- Return `400` with descriptive error for invalid input
- Don't validate in services — assume handlers validated

## Pagination

- Query params: `?page=1&page_size=20`
- Default page size: 20
- Response includes `total` count

## Long-Running Operations

- Return immediately with task ID
- Client polls `GET /api/v1/tasks/:task_id` for status
- Used by: video generation, batch operations

## CORS

- Origins configured in `config.yaml`: `server.cors_origins`
- Default: `http://localhost:3012` (frontend dev server)
- Middleware in `api/middlewares/`
