package endpoints

import (
	"fiber-api/api/handlers/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/sushan531/auth-sqlc/generated"
)

func UserRouter(route fiber.Router, queries *generated.Queries) {
	route.Get("/profile", handlers.GetProfileHandler(queries))
}
