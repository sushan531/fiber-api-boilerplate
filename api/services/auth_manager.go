package services

import (
	"database/sql"

	"github.com/sushan531/auth-sqlc/generated"
	"github.com/sushan531/jwk-auth/core/config"
	"github.com/sushan531/jwk-auth/core/manager"
	"github.com/sushan531/jwk-auth/core/repository"
	"github.com/sushan531/jwk-auth/service"
)

// AuthAPIServiceConfig holds the configuration for the auth service
type AuthAPIServiceConfig struct {
	DatabaseURL string
	Config      *config.Config
}

// AuthAPIService encapsulates all auth-related dependencies and functionality
type AuthAPIService struct {
	DB           *sql.DB
	Queries      *generated.Queries
	JWKManager   manager.JwkManager
	TokenService service.TokenService
	Config       *config.Config
}

// NewAuthAPIService creates a new auth manager with all dependencies initialized
func NewAuthAPIService(cfg AuthAPIServiceConfig) (*AuthAPIService, error) {
	// Initialize database connection
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Create queries instance
	queries := generated.New(db)

	// Initialize repositories and managers
	userRepo := repository.NewUserAuthRepository(queries)
	jwkManager := manager.NewJwkManager(userRepo, cfg.Config)
	jwtManager := manager.NewJwtManager(jwkManager)
	tokenService := service.NewTokenService(jwtManager, jwkManager, cfg.Config)

	return &AuthAPIService{
		DB:           db,
		Queries:      queries,
		JWKManager:   jwkManager,
		TokenService: tokenService,
		Config:       cfg.Config,
	}, nil
}

// Close closes the database connection
func (am *AuthAPIService) Close() error {
	return am.DB.Close()
}

// GetQueries returns the queries instance for external use
func (am *AuthAPIService) GetQueries() *generated.Queries {
	return am.Queries
}

// GetQueries returns the queries instance for external use
func (am *AuthAPIService) GetDB() *sql.DB {
	return am.DB
}

// GetJWKManager returns the JWK manager for external use
func (am *AuthAPIService) GetJWKManager() manager.JwkManager {
	return am.JWKManager
}

// GetAuthService returns the auth service for external use
func (am *AuthAPIService) GetAuthService() service.TokenService {
	return am.TokenService
}
