package middleware

import (
	"crypto/sha256"
	"fiber-api/api/errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sushan531/jwk-auth/service"
	"github.com/ua-parser/uap-go/uaparser"
)

// JWT middleware for protecting endpoints
func JWTMiddleware(tokenService service.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return errors.AuthenticationError(c, "Missing authorization header")
		}

		// Check Bearer format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return errors.AuthenticationError(c, "Invalid authorization header format")
		}

		token := parts[1]

		// Verify token
		claims, err := tokenService.VerifyToken(token)
		if err != nil {
			return errors.AuthenticationError(c, "Invalid or expired token")
		}

		// Validate device fingerprint
		storedFingerprint, hasFingerprintClaim := claims["device_fingerprint"].(string)
		if !hasFingerprintClaim || storedFingerprint == "" {
			return errors.AuthenticationError(c, "Invalid token: missing device fingerprint")
		}

		// Get current User-Agent from request and validate against stored fingerprint
		currentUserAgent := c.Get("User-Agent")
		if currentUserAgent == "" {
			return errors.AuthenticationError(c, "Missing User-Agent header")
		}

		// Generate fingerprint from current User-Agent and compare
		currentFingerprint := generateDeviceFingerprintHash(currentUserAgent)
		if currentFingerprint != storedFingerprint {
			return errors.AuthenticationError(c, "Device fingerprint mismatch")
		}

		// Store claims in context for use in handlers
		c.Locals("claims", claims)
		c.Locals("user_id", claims["user_id"])
		c.Locals("user_email", claims["user_email"])

		return c.Next()
	}
}

// generateDeviceFingerprintHash creates a device fingerprint hash from User-Agent
// This is a simplified version for middleware use
func generateDeviceFingerprintHash(userAgent string) string {
	if userAgent == "" {
		return generateHash("unknown")
	}

	parser := uaparser.NewFromSaved()
	client := parser.Parse(userAgent)

	// Extract key components for fingerprinting
	platform := normalizeString(client.Os.Family)
	browser := normalizeString(client.UserAgent.Family)
	version := normalizeString(client.UserAgent.Major)

	// Create a composite string for hashing
	composite := fmt.Sprintf("%s|%s|%s", platform, browser, version)

	// Generate SHA256 hash of the composite string
	return generateHash(composite)
}

// generateHash creates a SHA256 hash of the input string
func generateHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", hash)
}

// normalizeString normalizes strings for consistent fingerprinting
func normalizeString(s string) string {
	if s == "" {
		return "unknown"
	}
	return strings.ToLower(strings.TrimSpace(s))
}
