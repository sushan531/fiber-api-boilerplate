# Fiber API

A REST API built with Go Fiber framework for user management and authentication, featuring PostgreSQL integration and type-safe database queries via SQLC.

## Technology Stack

- **Language**: Go 1.25
- **Web Framework**: [Fiber v2](https://github.com/gofiber/fiber) - Fast HTTP framework
- **Database**: PostgreSQL with `lib/pq` driver
- **Code Generation**: SQLC via `github.com/sushan531/hk_ims_sqlc`

## Features

- User signup with profile management
- PostgreSQL database integration
- RESTful API design
- Type-safe database queries with SQLC
- JSON request/response handling
- Layered architecture (Routes → Handlers → Presenters)

## Project Structure

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

## Getting Started

### Prerequisites

- Go 1.25 or higher
- PostgreSQL database
- Git

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd fiber-api
```

2. Install dependencies:
```bash
go mod download
```

3. Set up your database connection:
   - Update the database URL in `main.go`
   - Format: `postgres://user:password@host:port/dbname?sslmode=disable`

### Running the Application

```bash
# Run in development mode
go run main.go

# Build the application
go build -o fiber-api

# Run the built binary
./fiber-api
```

The server will start on `http://localhost:3000`

## API Endpoints

### Base
- `GET /` - Welcome message

### Authentication
- `POST /api/signup` - User registration

## Architecture

### Layered Architecture

The project follows a clean layered architecture:

1. **Routes**: Define HTTP endpoints and route to handlers
2. **Handlers**: Process requests, interact with database, return responses
3. **Presenters**: Format response data consistently

### Response Format

All API responses follow a consistent structure:
```json
{
  "status": "success|error",
  "data": {},
  "error": null
}
```

## Development

### Code Organization

- **Handlers**: Business logic and database operations
  - Use dependency injection for database queries
  - Return Fiber errors for HTTP status codes
  - Parse request bodies into structs with validation tags

- **Presenters**: Response formatting
  - Return `*fiber.Map` with consistent structure
  - Keep response formatting separate from business logic

- **Routes**: HTTP routing configuration
  - Group related routes (e.g., `/api` prefix)
  - Pass dependencies (queries) to handlers

### Naming Conventions

- Handlers: `{Entity}{Action}Handler` (e.g., `UserSignUpHandler`)
- Routes: `{Entity}Router` (e.g., `AuthRouter`)
- Presenters: `{Action}Response` (e.g., `SignUpSuccessResponse`)
- Structs: PascalCase with JSON tags for API models

## Dependencies

Key dependencies used in this project:

- `github.com/gofiber/fiber/v2` - Web framework
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/google/uuid` - UUID generation
- `github.com/shopspring/decimal` - Decimal handling
- `github.com/sushan531/hk_ims_sqlc/generated` - Generated database queries

## Database

- Uses SQLC for type-safe database queries
- Generated queries are imported from external package
- Handles SQL null types appropriately (`sql.NullString`, etc.)
- Uses context from Fiber for database operations

## License

[Add your license here]

## Contributing

[Add contribution guidelines here]
