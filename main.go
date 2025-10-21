package main

import (
	"fiber-api/api/services"
	"fiber-api/config"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Load application configuration
	appConfig := config.LoadAppConfig()

	// Create server service
	serverService, err := services.NewAPIServerService(services.ServerConfig{
		DatabaseURL: appConfig.Database.URL,
		Port:        appConfig.Server.Port,
		Config:      appConfig.JWK,
	})
	if err != nil {
		log.Fatal("Failed to create server service:", err)
	}
	defer serverService.Close()

	// Register all routes
	serverService.RegisterAllRoutes()

	// Start the server
	log.Printf("🚀 Starting server on %s:%s", appConfig.Server.Host, appConfig.Server.Port)
	log.Fatal(serverService.Start())
}
