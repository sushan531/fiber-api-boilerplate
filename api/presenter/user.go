package presenter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sushan531/auth-sqlc/generated"
)

func UserProfileFetchResponse(data generated.GetUserProfileRow) *fiber.Map {
	return &fiber.Map{
		"status": true,
		"data": &fiber.Map{
			"full_name": data.FullName,
			"role":      data.UserRole.String,
			"email":     data.UserEmail,
		},
		"error": nil,
	}
}
