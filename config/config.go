package config

import (
	"os"

	"github.com/sushan531/jwk-auth/core/config"
)

// AppConfig holds all application configuration
type AppConfig struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWK      *config.Config
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	URL string
}

// LoadAppConfig loads configuration from environment variables with defaults
func LoadAppConfig() *AppConfig {
	return &AppConfig{
		Server: ServerConfig{
			Port: getEnv("PORT", "3000"),
			Host: getEnv("HOST", "localhost"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable"),
		},
		JWK: config.LoadConfig(),
	}
}

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
