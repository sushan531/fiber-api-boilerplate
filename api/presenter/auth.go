package presenter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sushan531/auth-sqlc/generated"
	"github.com/sushan531/jwk-auth/service"
)

func SignUpSuccessResponse(data generated.Auth) *fiber.Map {
	return &fiber.Map{
		"status": true,
		"data":   data,
		"error":  nil,
	}
}

func SignInSuccessResponse(data service.TokenPair) *fiber.Map {
	return &fiber.Map{
		"status":        true,
		"access_token":  data.AccessToken,
		"refresh_token": data.RefreshToken,
		"error":         nil,
	}
}
