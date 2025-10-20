# Authentication Services

This package provides reusable authentication services that can be imported and used in other Go modules.

## Features

- Complete authentication system with signup, login, and token refresh
- JWT token management with JWK (JSON Web Key) support
- bcrypt password hashing
- PostgreSQL database integration
- Fiber web framework integration

## Usage

### Option 1: Full Server Service

Use the complete server service with all routes and middleware:

```go
package main

import (
    "your-module/api/services"
    "github.com/sushan531/jwk-auth/core/config"
)

func main() {
    cfg := config.LoadConfig()
    
    serverService, err := services.NewServerService(services.ServerConfig{
        DatabaseURL: "postgres://user:pass@localhost:5432/db?sslmode=disable",
        Port:        "3000",
        Config:      cfg,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer serverService.Close()

    // Register all routes (auth + user)
    serverService.RegisterAllRoutes()
    
    // Start server
    log.Fatal(serverService.Start())
}
```

### Option 2: Auth Manager Only

Use just the authentication manager in your existing application:

```go
package main

import (
    "your-module/api/services"
    "github.com/gofiber/fiber/v2"
    "github.com/sushan531/jwk-auth/core/config"
)

func main() {
    cfg := config.LoadConfig()
    
    authManager, err := services.NewAuthManager(services.AuthServiceConfig{
        DatabaseURL: "postgres://user:pass@localhost:5432/db?sslmode=disable",
        Config:      cfg,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer authManager.Close()

    // Create your own Fiber app
    app := fiber.New()
    
    // Register auth routes
    authGroup := app.Group("/api/auth")
    authManager.RegisterRoutes(authGroup)
    
    // Add your custom routes
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("My custom app with auth")
    })
    
    log.Fatal(app.Listen(":3000"))
}
```

## API Endpoints

When you register the auth routes, you get these endpoints:

- `POST /signup` - User registration
- `POST /login` - User authentication
- `POST /refresh` - Token refresh

## Dependencies

Make sure to import this module in your `go.mod`:

```go
require (
    your-module v0.1.0
    github.com/sushan531/jwk-auth latest
    github.com/sushan531/auth-sqlc latest
    github.com/gofiber/fiber/v2 latest
    // ... other dependencies
)
```

## Configuration

The services require a `config.Config` object from the `jwk-auth` package. Make sure your configuration includes:

- Database connection settings
- JWT signing keys
- Other auth-related configuration

## Database

The services expect a PostgreSQL database with the schema defined in the `auth-sqlc` package.

## File Structure

- `auth_manager.go` - Core authentication manager with database and JWT handling
- `server.go` - Complete server service with Fiber app and route registration
- `README.md` - Documentation and usage examples