# Project Structure

## Architecture Pattern
This project follows a layered architecture with clear separation of concerns:

```
api/
├── handlers/          # HTTP request handlers (controllers)
│   ├── helpers/       # Handler utility functions
│   ├── auth.go        # Authentication endpoints
│   └── user.go        # User management endpoints
├── middleware/        # HTTP middleware (JWT auth, etc.)
├── models/           # Request/response data structures
├── presenter/        # Response formatting layer
├── routes/           # Route definitions and grouping
└── services/         # Business logic and service layer
```

## Key Conventions

### Handler Pattern
- Handlers are factory functions that return `fiber.Handler`
- Dependencies (queries, services) are injected as parameters
- Use structured logging with emoji prefixes (🚀, ❌, 🔒)

### Service Layer
- Services encapsulate business logic and external dependencies
- Use dependency injection pattern with config structs
- Services manage their own lifecycle (Close() methods)

### Models & Data
- Request/response models in `api/models/`
- Use struct tags for JSON binding (`json:"field_name"`)
- Database operations use SQLC generated queries
- Password hashing with bcrypt before storage

### Authentication Flow
- JWT tokens with JWK-based signing
- Middleware extracts and validates Bearer tokens
- Claims stored in Fiber context as `c.Locals()`
- Refresh token mechanism for token renewal

### Error Handling
- Return structured JSON error responses
- Log errors with context (user email, operation)
- Use appropriate HTTP status codes
- Validate input at handler level

### Route Organization
- Group routes by functionality (`/api`, `/api/user`)
- Apply middleware at group level
- Separate public and protected routes