package handlers

import (
	"fiber-api/api/errors"

	"github.com/gofiber/fiber/v2"
)

func ReturnErrorResponse(c *fiber.Ctx, err error) error {
	return errors.SendError(c, fiber.StatusBadRequest, errors.NewAPIError(
		errors.ErrCodeInternal,
		err.Error(),
		nil,
	))

}
