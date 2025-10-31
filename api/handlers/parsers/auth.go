package parsers

import (
	"database/sql"
	"fiber-api/api/errors"
	"fiber-api/api/handlers/validators"
	"fiber-api/api/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func SignupInputParser(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tx, _ := db.Begin()
		defer tx.Rollback()
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
		c.Locals("user_email", input.UserEmail)
		c.Locals("full_name", input.FullName)
		c.Locals("user_role", input.UserRole)
		c.Locals("address", input.Address)
		c.Locals("hashed_password", hashedPassword)
		return c.Next()
	}
}
