package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sushan531/jwk-auth/service"
)

// JWT middleware for protecting routes
func JWTMiddleware(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Missing authorization header"})
		}

		// Check Bearer format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid authorization header format"})
		}

		token := parts[1]

		// Verify token
		claims, err := authService.VerifyToken(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		// Store claims in context for use in handlers
		c.Locals("claims", claims)
		c.Locals("user_id", claims["user_id"])
		c.Locals("user_email", claims["user_email"])

		return c.Next()
	}
}
