package main

import (
	"database/sql"
	"fiber-api/api/middleware"
	"fiber-api/api/routes"

	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
	"github.com/sushan531/auth-sqlc/generated"

	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/sushan531/jwk-auth/core/config"
	"github.com/sushan531/jwk-auth/core/database"
	"github.com/sushan531/jwk-auth/core/manager"
	"github.com/sushan531/jwk-auth/core/repository"
	"github.com/sushan531/jwk-auth/service"
)

func main() {
	dbURL := "postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	app := fiber.New()
	queries := generated.New(db)

	cfg := config.LoadConfig()

	db, err = database.NewConnection(cfg.Database)

	// Initialize repositories and managers
	userRepo := repository.NewUserAuthRepository(queries)
	jwkManager := manager.NewJwkManager(userRepo, cfg)
	jwtManager := manager.NewJwtManager(jwkManager)
	authService := service.NewAuthService(jwtManager, jwkManager, cfg)

	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Rest api.")
	})

	authRoute := app.Group("/api")
	routes.AuthRouter(authRoute, queries, jwkManager, authService)

	userRoute := app.Group("/api/user", middleware.JWTMiddleware(authService))
	routes.UserRouter(userRoute, queries)

	log.Fatal(app.Listen(":3000"))
}
