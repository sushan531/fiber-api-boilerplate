package endpoints

import (
	"database/sql"
	"fiber-api/api/handlers/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/sushan531/auth-sqlc/generated"
	"github.com/sushan531/jwk-auth/core/manager"
	"github.com/sushan531/jwk-auth/service"
)

func AuthRouter(route fiber.Router, db *sql.DB, queries *generated.Queries, jwkManager manager.JwkManager, tokenService service.TokenService) {
	route.Post("/signup", handlers.UserSignUpHandler(queries, db))
	route.Post("/login", handlers.LoginHandler(queries, jwkManager, tokenService))
	route.Post("/refresh", handlers.RefreshTokenHandler(queries, tokenService))
}
