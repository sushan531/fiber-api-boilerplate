package routes

import (
	"fiber-api/api/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/sushan531/auth-sqlc/generated"
	"github.com/sushan531/jwk-auth/core/manager"
	"github.com/sushan531/jwk-auth/service"
)

func AuthRouter(route fiber.Router, queries *generated.Queries, jwkManager manager.JwkManager, tokenService service.TokenService) {
	route.Post("/signup", handlers.UserSignUpHandler(queries))
	route.Post("/login", handlers.LoginHandler(queries, jwkManager, tokenService))
	route.Post("/refresh", handlers.RefreshTokenHandler(queries, tokenService))

	// Device detection example endpoints
	route.Get("/device-info", handlers.DeviceInfoHandler())
	route.Post("/login-example", handlers.LoginResponseHandler())
}
