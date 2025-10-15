package routes

import (
	"fiber-api/api/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/sushan531/hk_ims_sqlc/generated"
)

func AuthRouter(route fiber.Router, queries *generated.Queries) {
	route.Post("/signup", handlers.UserSignUpHandler(queries))
	route.Post("/login", handlers.LoginHandler(queries))
}
