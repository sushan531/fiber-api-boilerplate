# Fiber Auth API

A robust Go-based REST API boilerplate for authentication services built with the Fiber framework. This API provides secure user authentication, JWT token management with JWK (JSON Web Key) support, and device-aware session handling.

## 🚀 Features

- **User Authentication**: Registration, login, and profile management
- **JWT Token Management**: Access and refresh tokens with JWK-based signing
- **Device Detection**: Automatic device type detection (Web, Android, iOS)
- **Session Management**: Device-specific session keys for enhanced security
- **Input Validation**: Comprehensive request validation with detailed error messages
- **Structured Logging**: Emoji-prefixed logging for better development experience
- **Standardized Responses**: Consistent API response format across all endpoints
- **Environment Configuration**: Environment-based configuration management
- **PostgreSQL Integration**: SQLC-generated queries for type-safe database operations

## 🏗️ Architecture

The project follows a clean, layered architecture with clear separation of concerns:

```
fiber-api/
├── api/
│   ├── errors/           # Standardized error handling
│   ├── handlers/         # HTTP request handlers
│   │   └── helpers/      # Handler utility functions
│   ├── middleware/       # HTTP middleware (JWT, device detection)
│   ├── models/          # Request/response data structures
│   ├── presenter/       # Response formatting layer
│   ├── routes/          # Route definitions and grouping
│   ├── services/        # Business logic and service layer
│   └── validators/      # Input validation layer
├── config/              # Configuration management
├── pkg/
│   └── logger/          # Structured logging utilities
└── main.go             # Application entry point
```

## 🛠️ Technology Stack

- **Language**: Go 1.25.1
- **Web Framework**: [Fiber v2](https://github.com/gofiber/fiber) (Express-inspired Go web framework)
- **Database**: PostgreSQL with `lib/pq` driver
- **Authentication**: Custom JWK-based JWT authentication via `sushan531/jwk-auth`
- **Database Layer**: SQLC generated queries via `sushan531/auth-sqlc`
- **Password Hashing**: bcrypt for secure password storage
- **Device Detection**: User-Agent parsing for device type identification

## 📋 Prerequisites

- Go 1.25.1 or higher
- PostgreSQL database
- Git

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd fiber-api
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Environment Configuration

Copy the example environment file and configure your settings:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Server Configuration
PORT=3000
HOST=localhost

# Database Configuration
DATABASE_URL=postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable

# Environment
ENVIRONMENT=development
```

### 4. Database Setup

Ensure your PostgreSQL database is running and accessible with the connection string provided in your `.env` file.

### 5. Run the Application

```bash
# Development
go run main.go

# Build and run
go build -o fiber-api
./fiber-api
```

The server will start on `http://localhost:3000` by default.

## 📚 API Documentation

### Base URL
```
http://localhost:3000
```

### Authentication Endpoints

#### User Registration
```http
POST /api/signup
Content-Type: application/json

{
  "user_email": "user@example.com",
  "password": "securepassword123",
  "full_name": "John Doe",
  "user_role": "user",
  "address": "123 Main St"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user_id": "uuid-here",
    "email": "user@example.com",
    "full_name": "John Doe",
    "created_at": "2024-01-01T12:00:00Z"
  },
  "message": "User registered successfully"
}
```

#### User Login
```http
POST /api/login
Content-Type: application/json

{
  "user_email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer"
  },
  "message": "Authentication successful"
}
```

#### Token Refresh
```http
POST /api/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Protected Endpoints

#### Get User Profile
```http
GET /api/user/profile
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user_id": "uuid-here",
    "email": "user@example.com",
    "full_name": "John Doe",
    "role": "user",
    "address": "123 Main St"
  },
  "message": "Profile retrieved successfully"
}
```

### Device Detection

The API automatically detects device types based on User-Agent headers:

- **Web**: Desktop browsers, mobile web browsers
- **Android**: Native Android applications
- **iOS**: Native iOS applications

Device information is used for:
- Session key generation
- Security logging
- Device-specific token management

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `3000` |
| `HOST` | Server host | `localhost` |
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `ENVIRONMENT` | Application environment | `development` |

### Database Configuration

The application expects a PostgreSQL database with the following connection string format:
```
postgres://username:password@host:port/database?sslmode=disable
```

## 🧪 Development

### Project Structure Conventions

#### Handler Pattern
- Handlers are factory functions that return `fiber.Handler`
- Dependencies (queries, services) are injected as parameters
- Use structured logging with emoji prefixes (🚀, ❌, 🔒)

#### Service Layer
- Services encapsulate business logic and external dependencies
- Use dependency injection pattern with config structs
- Services manage their own lifecycle (Close() methods)

#### Error Handling
- Use standardized error responses from `api/errors` package
- Log errors with context (user email, operation)
- Return appropriate HTTP status codes
- Validate input at handler level using `api/validators`

#### Response Format
All API responses follow this structure:
```json
{
  "success": boolean,
  "data": object,
  "message": string,
  "error": {
    "code": string,
    "message": string,
    "details": object
  }
}
```

### Adding New Endpoints

1. **Create Model** (if needed): Add request/response models in `api/models/`
2. **Add Validation**: Create validators in `api/validators/`
3. **Create Handler**: Implement handler in `api/handlers/`
4. **Add Route**: Register route in `api/routes/`
5. **Update Presenter**: Add response formatting in `api/presenter/`

### Common Commands

```bash
# Run the application
go run main.go

# Build the application
go build -o fiber-api

# Install dependencies
go mod tidy

# Download dependencies
go mod download

# Run tests (when available)
go test ./...

# Format code
go fmt ./...

# Vet code
go vet ./...
```

## 🔒 Security Features

- **Password Hashing**: bcrypt with default cost
- **JWT Security**: JWK-based token signing and verification
- **Device-Specific Sessions**: Separate session keys per device type
- **Input Validation**: Comprehensive request validation
- **SQL Injection Prevention**: SQLC-generated type-safe queries
- **Structured Error Handling**: No sensitive information leakage

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fiber](https://github.com/gofiber/fiber) - Express-inspired web framework for Go
- [SQLC](https://sqlc.dev/) - Generate type-safe code from SQL
- [JWK Auth](https://github.com/sushan531/jwk-auth) - JWT/JWK authentication library

## 📞 Support

If you have any questions or need help, please open an issue in the repository.