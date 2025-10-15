package handlers

import (
	"database/sql"
	"fiber-api/api/presenter"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sushan531/hk_ims_sqlc/generated"
)

type SignUp struct {
	UserEmail string `json:"user_email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FullName  string `json:"full_name" binding:"required"`
	UserRole  string `json:"user_role" binding:"required"`
	Address   string `json:"address"`
}

type Login struct {
	UserEmail string `json:"user_email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func UserSignUpHandler(queries *generated.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		inputBody := new(SignUp)
		context := c.Context()
		if err := c.BodyParser(&inputBody); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		userParams := generated.InsertUserProfileParams{
			UserEmail: inputBody.UserEmail,
			Password:  inputBody.Password,
			FullName:  inputBody.FullName,
			UserRole:  sql.NullString{String: inputBody.UserRole, Valid: true},
			Address:   sql.NullString{String: inputBody.Address, Valid: true},
		}
		user, err := queries.InsertUserProfile(context, userParams)
		if err != nil {
			log.Printf("❌ Failed to insert user: %v", err)
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return c.JSON(presenter.SignUpSuccessResponse(user))
	}
}

func LoginHandler(queries *generated.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		inputBody := new(Login)
		context := c.Context()
		if err := c.BodyParser(&inputBody); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		auth, err := queries.GetUserAuth(context, inputBody.UserEmail)
		if err != nil {
			log.Printf("❌ Failed to Fetch user: %v", err)
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		if auth.Password != inputBody.Password {
			return fiber.ErrUnauthorized
		}
		return c.JSON(presenter.SignInSuccessResponse(auth.UserEmail))
	}
}
