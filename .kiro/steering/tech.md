---
inclusion: always
---

# Technology Stack

## Core Technologies
- **Language**: Go 1.25
- **Web Framework**: Fiber v2 (Fast HTTP framework)
- **Database**: PostgreSQL with `lib/pq` driver
- **Code Generation**: SQLC via `github.com/sushan531/hk_ims_sqlc`

## Key Dependencies
- `github.com/gofiber/fiber/v2` - Web framework
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/google/uuid` - UUID generation
- `github.com/shopspring/decimal` - Decimal handling
- `github.com/sushan531/hk_ims_sqlc/generated` - Generated database queries

## Common Commands

### Development
```bash
# Run the application
go run main.go

# Build the application
go build -o fiber-api

# Install dependencies
go mod tidy

# Update dependencies
go mod download
```

### Database
- Database URL format: `postgres://user:password@host:port/dbname?sslmode=disable`
- Uses SQLC for type-safe database queries
- Generated queries are imported from external package

### Server Configuration
- Default port: 3000
- API prefix: `/api`
- Database connection handled in main.go