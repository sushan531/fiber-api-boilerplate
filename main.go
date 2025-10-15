package main

import (
	"database/sql"
	"fiber-api/api/routes"

	_ "github.com/lib/pq"
	"github.com/sushan531/hk_ims_sqlc/generated"

	"log"

	"github.com/gofiber/fiber/v2"
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

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Rest api.")
	})
	api := app.Group("/api")
	routes.AuthRouter(api, queries)

	log.Fatal(app.Listen(":3000"))
}
