package main

import (
	"fiber-api/api/services"
	"log"

	_ "github.com/lib/pq"
	"github.com/sushan531/jwk-auth/core/config"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Database URL - you can move this to config later
	dbURL := "postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable"

	// Create server service
	serverService, err := services.NewAPIServerService(services.ServerConfig{
		DatabaseURL: dbURL,
		Port:        "3000",
		Config:      cfg,
	})
	if err != nil {
		log.Fatal("Failed to create server service:", err)
	}
	defer serverService.Close()

	// Register all routes
	serverService.RegisterAllRoutes()

	// Start the server
	log.Fatal(serverService.Start())
}
