package handlers

import (
	"fiber-api/api/presenter"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sushan531/auth-sqlc/generated"
)

func GetProfileHandler(queries *generated.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		context := c.Context()
		userEmail := c.Locals("user_id").(uuid.UUID)
		userProfile, err := queries.GetUserProfile(context, userEmail)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return c.JSON(presenter.UserProfileFetchResponse(userProfile))
	}
}
