package handlers

import (
	"database/sql"
	"fiber-api/api/handlers/helpers"
	"fiber-api/api/models"
	"fiber-api/api/presenter"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sushan531/auth-sqlc/generated"
	"github.com/sushan531/jwk-auth/core/manager"
	"github.com/sushan531/jwk-auth/service"
	"golang.org/x/crypto/bcrypt"
)

func UserSignUpHandler(queries *generated.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		// Parse request body
		var input models.SignUp
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request payload",
			})
		}
		// Validate basic fields (optional but recommended)
		if input.UserEmail == "" || input.Password == "" || input.FullName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email, password, and full name are required",
			})
		}
		// 🔒 Hash the password before storing
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("❌ Failed to hash password for %s: %v", input.UserEmail, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process password",
			})
		}
		// Prepare user parameters
		userParams := generated.InsertUserProfileParams{
			UserEmail: input.UserEmail,
			Password:  string(hashedPassword),
			FullName:  input.FullName,
			UserRole:  sql.NullString{String: input.UserRole, Valid: input.UserRole != ""},
			Address:   sql.NullString{String: input.Address, Valid: input.Address != ""},
		}
		// Insert new user record
		user, err := queries.InsertUserProfile(ctx, userParams)
		if err != nil {
			log.Printf("❌ Failed to insert user %s: %v", input.UserEmail, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to create user — possibly duplicate email",
			})
		}
		// Return success response
		return c.Status(fiber.StatusCreated).JSON(presenter.SignUpSuccessResponse(user))
	}
}

func LoginHandler(queries *generated.Queries, jwkManager manager.JwkManager, tokenService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		// Parse request body
		var input models.Login
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request payload",
			})
		}

		// Fetch user auth record
		auth, err := queries.GetUserAuth(ctx, input.UserEmail)
		if err != nil {
			log.Printf("❌ Failed to fetch user for email %s: %v", input.UserEmail, err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		// Validate password using bcrypt
		if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(input.Password)); err != nil {
			log.Printf("❌ Invalid password attempt for user %s", input.UserEmail)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}

		// Create a new session key
		keyID, err := jwkManager.CreateSessionKey(auth.UserProfileID, "web")
		if err != nil {
			log.Printf("❌ Failed to create session key for user %s: %v", input.UserEmail, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create session key",
			})
		}

		// Create JWT claims
		claims, err := helpers.CreateJWTClaims(queries, ctx, auth.UserProfileID)
		if err != nil {
			log.Printf("❌ Failed to create JWT claims for user %s: %v", input.UserEmail, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create JWT claims",
			})
		}

		// Generate token pair
		tokenPair, err := tokenService.GenerateTokenPairWithKeyID(claims.ToMap(), keyID)
		if err != nil {
			log.Printf("❌ Failed to generate tokens for user %s: %v", input.UserEmail, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate tokens",
			})
		}

		// Return successful response
		return c.JSON(presenter.SignInSuccessResponse(*tokenPair))
	}
}

func RefreshTokenHandler(queries *generated.Queries, authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		// Parse request body
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request payload",
			})
		}
		// Verify the refresh token
		refreshClaims, err := authService.VerifyRefreshToken(req.RefreshToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired refresh token",
			})
		}
		// Parse user_id from claims
		userID, err := helpers.ExtractUserIdFromMapObj(refreshClaims)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		// Extract keyID from token
		keyID, err := authService.ExtractKeyIDFromToken(req.RefreshToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid refresh token",
			})
		}
		// Create new JWT claims
		claims, err := helpers.CreateJWTClaims(queries, ctx, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create JWT claims for token refresh",
			})
		}
		// Generate refreshed tokens
		tokenPair, err := authService.RefreshTokensWithKeyID(req.RefreshToken, claims.ToMap(), keyID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Failed to refresh tokens",
			})
		}
		return c.JSON(presenter.SignInSuccessResponse(*tokenPair))
	}
}
