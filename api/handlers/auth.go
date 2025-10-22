package handlers

import (
	"database/sql"
	"fiber-api/api/errors"
	"fiber-api/api/handlers/helpers"
	"fiber-api/api/middleware"
	"fiber-api/api/models"
	"fiber-api/api/presenter"
	"fiber-api/api/validators"
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
			return errors.ValidationError(c, "Invalid request payload")
		}

		// Validate input
		validation := validators.ValidateSignUp(input)
		if !validation.IsValid {
			return errors.ValidationErrorWithDetails(c, "Validation failed", validation.Errors)
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("❌ Failed to hash password for %s: %v", input.UserEmail, err)
			return errors.InternalError(c, "Failed to process password")
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
			return errors.SendError(c, fiber.StatusConflict, errors.NewAPIError(
				errors.ErrCodeDuplicate,
				"Email already exists",
				nil,
			))
		}

		// Return success response
		return c.Status(fiber.StatusCreated).JSON(presenter.SignUpSuccessResponse(user))
	}
}

func LoginHandler(queries *generated.Queries, jwkManager manager.JwkManager, tokenService service.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		// Parse request body
		var input models.Login
		if err := c.BodyParser(&input); err != nil {
			return errors.ValidationError(c, "Invalid request payload")
		}

		// Validate input
		validation := validators.ValidateLogin(input)
		if !validation.IsValid {
			return errors.ValidationErrorWithDetails(c, "Validation failed", validation.Errors)
		}

		// Fetch user auth record
		auth, err := queries.GetUserAuth(ctx, input.UserEmail)
		if err != nil {
			log.Printf("❌ Failed to fetch user for email %s: %v", input.UserEmail, err)
			return errors.AuthenticationError(c, "Invalid email or password")
		}

		// Validate password using bcrypt
		if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(input.Password)); err != nil {
			log.Printf("❌ Invalid password attempt for user %s", input.UserEmail)
			return errors.AuthenticationError(c, "Invalid email or password")
		}

		// Get device type from middleware
		deviceType := middleware.GetDeviceType(c)

		// Generate device fingerprint from User-Agent from request
		userAgent := c.Get("User-Agent")
		deviceFingerprint := helpers.GenerateDeviceFingerprint(userAgent)

		// Create a new session key with device type
		keyID, err := jwkManager.CreateSessionKey(auth.UserProfileID.String(), string(deviceType))
		if err != nil {
			log.Printf("❌ Failed to create session key for user %s on device %s: %v", input.UserEmail, deviceType, err)
			return errors.InternalError(c, "Failed to create session key")
		}

		// Create JWT claims with device fingerprint
		claims, err := helpers.CreateJWTClaims(queries, ctx, auth.UserProfileID, deviceFingerprint.Hash)
		if err != nil {
			log.Printf("❌ Failed to create JWT claims for user %s: %v", input.UserEmail, err)
			return errors.InternalError(c, "Failed to create JWT claims")
		}

		// Generate token pair
		tokenPair, err := tokenService.GenerateTokenPairWithKeyID(claims.ToMap(), keyID)
		if err != nil {
			log.Printf("❌ Failed to generate tokens for user %s: %v", input.UserEmail, err)
			return errors.InternalError(c, "Failed to generate tokens")
		}

		// Return successful response
		log.Printf("🚀 User %s logged in successfully from %s device", input.UserEmail, deviceType)
		return c.JSON(presenter.SignInSuccessResponse(*tokenPair))
	}
}

func RefreshTokenHandler(queries *generated.Queries, tokenService service.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		// Parse request body
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.BodyParser(&req); err != nil {
			return errors.ValidationError(c, "Invalid request payload")
		}
		// Verify the refresh token
		refreshClaims, err := tokenService.VerifyRefreshToken(req.RefreshToken)
		if err != nil {
			return errors.AuthenticationError(c, "Invalid or expired refresh token")
		}
		// Parse user_id from claims
		userID, err := helpers.ExtractUserIdFromMapObj(refreshClaims)
		if err != nil {
			return errors.ValidationError(c, err.Error())
		}
		// Extract keyID from token
		keyID, err := tokenService.ExtractKeyIDFromToken(req.RefreshToken)
		if err != nil {
			return errors.AuthenticationError(c, "Invalid refresh token")
		}
		// Extract device fingerprint from refresh token claims
		storedFingerprint, hasFingerprintClaim := helpers.GetFingerprintFromClaims(refreshClaims)
		if !hasFingerprintClaim {
			return errors.AuthenticationError(c, "Invalid refresh token: missing device fingerprint")
		}

		// Validate current device fingerprint against stored one
		currentUserAgent := c.Get("User-Agent")
		if !helpers.ValidateDeviceFingerprint(currentUserAgent, storedFingerprint) {
			log.Printf("❌ Device fingerprint mismatch for user %s during token refresh", userID.String())
			return errors.AuthenticationError(c, "Device fingerprint mismatch")
		}

		// Create new JWT claims with same device fingerprint
		claims, err := helpers.CreateJWTClaims(queries, ctx, userID, storedFingerprint)
		if err != nil {
			return errors.InternalError(c, "Failed to create JWT claims for token refresh")
		}
		// Generate refreshed tokens
		tokenPair, err := tokenService.RefreshTokensWithKeyID(req.RefreshToken, claims.ToMap(), keyID)
		if err != nil {
			return errors.AuthenticationError(c, "Failed to refresh tokens")
		}
		return c.JSON(presenter.SignInSuccessResponse(*tokenPair))
	}
}
