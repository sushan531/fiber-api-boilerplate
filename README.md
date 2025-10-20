# Fiber Auth API

A Go-based REST API boilerplate for authentication services built with the Fiber framework. This application provides secure user authentication with JWT tokens, JWK-based signing, and PostgreSQL integration.

## Features

- 🔐 User registration and authentication
- 🔑 JWT token management with refresh tokens
- 🛡️ JWK (JSON Web Key) based authentication
- 🗄️ PostgreSQL database integration
- 👥 Role-based access control
- 🚀 High-performance Fiber web framework

## Tech Stack

- **Language**: Go 1.25.1
- **Web Framework**: Fiber v2
- **Database**: PostgreSQL
- **Authentication**: JWK-based JWT tokens
- **Password Hashing**: bcrypt
- **Database Layer**: SQLC generated queries

## Prerequisites

- Go 1.25.1 or higher
- PostgreSQL database
- Git

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd fiber-api
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up your PostgreSQL database and update the connection string in `main.go`:
```go
dbURL := "postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable"
```

## Running the Application

### Development
```bash
go run main.go
```

### Production Build
```bash
go build -o fiber-api
./fiber-api
```

The server will start on port 3000 by default.

## API Endpoints

### Authentication Routes (`/api`)

#### Register User
```http
POST /api/signup
Content-Type: application/json

{
  "user_email": "user@example.com",
  "password": "securepassword",
  "full_name": "John Doe",
  "user_role": "user"
}
```

#### Login
```http
POST /api/login
Content-Type: application/json

{
  "user_email": "user@example.com",
  "password": "securepassword"
}
```

#### Refresh Token
```http
POST /api/refresh
Content-Type: application/json

{
  "refresh_token": "your_refresh_token_here"
}
```

### Protected User Routes (`/api/user`)

All routes under `/api/user` require a valid JWT token in the Authorization header:

```http
Authorization: Bearer <your_jwt_token>
```

## Project Structure

```
api/
├── handlers/          # HTTP request handlers
│   ├── helpers/       # Handler utility functions
│   └── auth.go        # Authentication endpoints
├── middleware/        # HTTP middleware (JWT auth)
├── models/           # Request/response data structures
├── presenter/        # Response formatting layer
├── routes/           # Route definitions
└── services/         # Business logic and service layer
```

## Architecture

This project follows a layered architecture with clear separation of concerns:

- **Handlers**: Process HTTP requests and responses
- **Services**: Contain business logic and manage dependencies
- **Models**: Define data structures for requests/responses
- **Middleware**: Handle cross-cutting concerns like authentication
- **Presenters**: Format API responses consistently

## Key Features

### Security
- Password hashing with bcrypt
- JWT tokens with JWK-based signing
- Session key management
- Token refresh mechanism

### Database
- PostgreSQL integration with connection pooling
- SQLC generated queries for type safety
- Proper SQL injection prevention

### Error Handling
- Structured JSON error responses
- Comprehensive logging with context
- Appropriate HTTP status codes

## Development

### Adding New Routes
1. Define models in `api/models/`
2. Create handlers in `api/handlers/`
3. Register routes in `api/routes/`
4. Add to service registration in `api/services/server.go`

### Database Changes
This project uses SQLC for database operations. Update your SQL schemas and regenerate queries as needed.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.