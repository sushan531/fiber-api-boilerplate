---
inclusion: always
---

# Project Structure

## Directory Organization

```
fiber-api/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
└── api/                    # API layer
    ├── handlers/           # Request handlers (business logic)
    │   └── auth.go        # Authentication handlers
    ├── presenter/          # Response formatters
    │   └── auth.go        # Authentication response formatting
    └── routes/             # Route definitions
        └── auth.go        # Authentication routes
```

## Architecture Patterns

### Layered Architecture
- **Routes**: Define HTTP endpoints and route to handlers
- **Handlers**: Process requests, interact with database, return responses
- **Presenters**: Format response data consistently

### Code Organization Rules
1. **Handlers**: Business logic and database operations
   - Use dependency injection for database queries
   - Return Fiber errors for HTTP status codes
   - Parse request bodies into structs with validation tags

2. **Presenters**: Response formatting
   - Return `*fiber.Map` with consistent structure: `{status, data, error}`
   - Keep response formatting separate from business logic

3. **Routes**: HTTP routing configuration
   - Group related routes (e.g., `/api` prefix)
   - Pass dependencies (queries) to handlers
   - Use descriptive route names

### Naming Conventions
- Handlers: `{Entity}{Action}Handler` (e.g., `UserSignUpHandler`)
- Routes: `{Entity}Router` (e.g., `AuthRouter`)
- Presenters: `{Action}Response` (e.g., `SignUpSuccessResponse`)
- Structs: PascalCase with JSON tags for API models

### Database Integration
- Use SQLC generated queries for type safety
- Pass `*generated.Queries` to handlers via dependency injection
- Handle SQL null types appropriately (`sql.NullString`, etc.)
- Use context from Fiber for database operations