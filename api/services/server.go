package services

import (
	"fiber-api/api/middleware"
	"fiber-api/api/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sushan531/jwk-auth/core/config"
)

// ServerConfig holds the configuration for the server
type ServerConfig struct {
	DatabaseURL string
	Port        string
	Config      *config.Config
}

// ServerService encapsulates the entire server functionality
type ServerService struct {
	App            *fiber.App
	AuthAPIService *AuthAPIService
	Config         ServerConfig
}

// NewAPIServerService creates a new server service with all dependencies
func NewAPIServerService(cfg ServerConfig) (*ServerService, error) {
	// Create Fiber app
	app := fiber.New()

	// Initialize auth manager
	authService, err := NewAuthAPIService(AuthAPIServiceConfig{
		DatabaseURL: cfg.DatabaseURL,
		Config:      cfg.Config,
	})
	if err != nil {
		return nil, err
	}

	// Setup welcome route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Auth BoilerPlate Rest API.")
	})

	return &ServerService{
		App:            app,
		AuthAPIService: authService,
		Config:         cfg,
	}, nil
}

// RegisterAuthRoutes registers authentication routes
func (ss *ServerService) RegisterAuthRoutes() {
	authRoute := ss.App.Group("/api")
	ss.AuthAPIService.RegisterRoutes(authRoute)
}

// RegisterUserRoutes registers user routes with JWT middleware
func (ss *ServerService) RegisterUserRoutes() {
	userRoute := ss.App.Group("/api/user", middleware.JWTMiddleware(ss.AuthAPIService.GetAuthService()))
	routes.UserRouter(userRoute, ss.AuthAPIService.GetQueries())
}

// RegisterAllRoutes registers both auth and user routes
func (ss *ServerService) RegisterAllRoutes() {
	ss.RegisterAuthRoutes()
	ss.RegisterUserRoutes()
}

// Start starts the server on the configured port
func (ss *ServerService) Start() error {
	port := ss.Config.Port
	if port == "" {
		port = ":3000"
	}
	if port[0] != ':' {
		port = ":" + port
	}

	log.Printf("🚀 Server starting on port %s", port)
	return ss.App.Listen(port)
}

// Close closes all resources
func (ss *ServerService) Close() error {
	return ss.AuthAPIService.Close()
}

// GetApp returns the fiber app instance for custom configuration
func (ss *ServerService) GetApp() *fiber.App {
	return ss.App
}

// GetAuthManager returns the auth manager for external access
func (ss *ServerService) GetAuthManager() *AuthAPIService {
	return ss.AuthAPIService
}
