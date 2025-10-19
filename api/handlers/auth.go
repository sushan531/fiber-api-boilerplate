package handlers

import (
	"database/sql"
	"fiber-api/api/presenter"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sushan531/auth-sqlc/generated"
	"github.com/sushan531/jwk-auth/core/manager"
	"github.com/sushan531/jwk-auth/service"
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

func LoginHandler(queries *generated.Queries, jwkManager manager.JwkManager, tokenService service.AuthService) fiber.Handler {
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
		// Create session key
		keyID, err := jwkManager.CreateSessionKey(auth.UserProfileID, "web")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create session key"})
		}
		// Prepare claims
		claims := make(map[string]interface{})
		claims["user_id"] = auth.UserProfileID.String()
		claims["user_email"] = auth.UserEmail
		claims["role"] = "admin"

		tokenPair, err := tokenService.GenerateTokenPairWithKeyID(claims, keyID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to generate tokens"})
		}
		return c.JSON(presenter.SignInSuccessResponse(*tokenPair))
	}
}

//func VerifyRefreshToken(authService service.AuthService) fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		var req struct {
//			RefreshToken string `json:"refresh_token" binding:"required"`
//		}
//		if err := c.BodyParser(&req); err != nil {
//			return fiber.NewError(fiber.StatusBadRequest, err.Error())
//		}
//		claims, err := authService.VerifyRefreshToken(req.RefreshToken)
//		if err != nil {
//			return fiber.NewError(fiber.StatusBadRequest, err.Error())
//		}
//		authService.GenerateTokenPairWithKeyID()
//	}
//}
//
//func RefreshTokenHandler(authService service.AuthService) fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		var req struct {
//			RefreshToken string                 `json:"refresh_token"`
//			Claims       map[string]interface{} `json:"claims,omitempty"`
//		}
//
//		if err := c.BodyParser(&req); err != nil {
//			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
//		}
//
//		// Extract keyID from refresh token
//		keyID, err := authService.ExtractKeyIDFromToken(req.RefreshToken)
//		if err != nil {
//			return c.Status(400).JSON(fiber.Map{"error": "Invalid refresh token"})
//		}
//
//		// Refresh tokens
//		tokenPair, err := authService.RefreshTokensWithKeyID(req.RefreshToken, req.Claims, keyID)
//		if err != nil {
//			return c.Status(401).JSON(fiber.Map{"error": "Failed to refresh tokens"})
//		}
//
//		return c.JSON(tokenPair)
//	}
//}
