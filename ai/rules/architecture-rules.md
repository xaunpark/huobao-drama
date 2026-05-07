# Rules: Architecture — Architecture Constraints

> These constraints maintain system integrity. Do not violate without an ADR.

## Layer Dependencies

```
✅ Allowed:
  api → application → domain
  api → infrastructure (for DI wiring only, in routes.go)
  application → domain
  application → pkg
  application → infrastructure (for external services)
  infrastructure → domain
  pkg → (no internal dependencies)

❌ Forbidden:
  domain → application (domain is pure)
  domain → infrastructure
  domain → api
  handlers → handlers (handlers don't call each other)
  pkg → application (pkg is generic, not domain-specific)
```

## Model Rules

1. All models must have `CreatedAt`, `UpdatedAt`, `DeletedAt` fields
2. All models must have a `TableName()` method
3. Use `*string`, `*int`, `*bool` for optional fields (nullable)
4. Use `datatypes.JSON` for dynamic/nested data
5. New fields must be nullable (don't break existing rows)
6. Never remove model fields — only add

## Service Rules

1. Services receive `*gorm.DB` directly — no repository abstraction
2. Services receive `*logger.Logger` for logging
3. Services must not import from `api/` package
4. Business logic lives ONLY in services, not handlers

## Handler Rules

1. Handlers parse HTTP input and format output — no business logic
2. One handler struct per resource domain
3. Constructors wire all dependencies via parameters
4. HTTP errors: `c.JSON(statusCode, gin.H{"error": message})`

## Route Rules

1. All routes defined in `api/routes/routes.go`
2. Use `api.Group("/resource")` for grouping
3. RESTful verbs: GET (read), POST (create/action), PUT (update), DELETE (remove)
4. Rate limiting applied at API group level

## Frontend Rules

1. Composition API only (`<script setup lang="ts">`)
2. Element Plus for all UI components
3. Pinia for state management
4. TypeScript types manually mirror Go models
5. API calls through `web/src/api/` modules (not inline Axios)
