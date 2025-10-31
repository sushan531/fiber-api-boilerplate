package endpoints

import (
	"database/sql"
	"fiber-api/api/handlers/handlers"
	"fiber-api/api/handlers/parsers"

	"github.com/gofiber/fiber/v2"
	"github.com/sushan531/auth-sqlc/generated"
	"github.com/sushan531/jwk-auth/core/manager"
	"github.com/sushan531/jwk-auth/service"
)

func AuthRouter(route fiber.Router, db *sql.DB, queries *generated.Queries, jwkManager manager.JwkManager, tokenService service.TokenService) {
	route.Post("/signup", parsers.SignupInputParser(db), handlers.UserSignUpHandler(queries, db))
	route.Post("/login", handlers.LoginHandler(queries, jwkManager, tokenService))
	route.Post("/refresh", handlers.RefreshTokenHandler(queries, tokenService))
}
