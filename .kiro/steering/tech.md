# Technology Stack

## Core Technologies
- **Language**: Go 1.25.1
- **Web Framework**: Fiber v2 (Express-inspired Go web framework)
- **Database**: PostgreSQL with `lib/pq` driver
- **Authentication**: Custom JWK-based JWT authentication via `sushan531/jwk-auth`
- **Database Layer**: SQLC generated queries via `sushan531/auth-sqlc`

## Key Dependencies
- `github.com/gofiber/fiber/v2` - Web framework
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/sushan531/auth-sqlc` - Generated database queries
- `github.com/sushan531/jwk-auth` - JWT/JWK authentication library
- `golang.org/x/crypto/bcrypt` - Password hashing

## Common Commands
```bash
# Run the application
go run main.go

# Build the application
go build -o fiber-api

# Install dependencies
go mod tidy

# Download dependencies
go mod download
```

## Database Setup
The application expects a PostgreSQL database with connection string format:
```
postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable
```

Default server port is 3000.