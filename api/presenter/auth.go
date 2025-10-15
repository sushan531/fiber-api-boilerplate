package presenter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sushan531/hk_ims_sqlc/generated"
)

func SignUpSuccessResponse(data generated.Auth) *fiber.Map {
	return &fiber.Map{
		"status": true,
		"data":   data,
		"error":  nil,
	}
}

func SignInSuccessResponse(data string) *fiber.Map {
	return &fiber.Map{
		"status": true,
		"data":   data,
		"error":  nil,
	}
}
