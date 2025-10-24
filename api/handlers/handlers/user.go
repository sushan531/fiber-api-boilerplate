package handlers

import (
	"fiber-api/api/errors"
	"fiber-api/api/handlers/presenter"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sushan531/auth-sqlc/generated"
)

func GetProfileHandler(queries *generated.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		// Extract user ID from JWT claims
		userID, ok := c.Locals("user_id").(string)
		userUuidID, _ := uuid.Parse(userID)
		if !ok {
			return errors.AuthenticationError(c, "Invalid user session")
		}

		// Fetch user profile
		userProfile, err := queries.GetUserProfile(ctx, userUuidID)
		if err != nil {
			return errors.NotFoundError(c, "User profile not found")
		}

		return c.JSON(presenter.UserProfileFetchResponse(userProfile))
	}
}
