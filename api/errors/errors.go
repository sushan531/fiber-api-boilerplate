package errors

import (
	"github.com/gofiber/fiber/v2"
)

// APIError represents a structured API error
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Success bool     `json:"success"`
	Error   APIError `json:"error"`
}

// Common error codes
const (
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeAuthentication = "AUTHENTICATION_ERROR"
	ErrCodeAuthorization  = "AUTHORIZATION_ERROR"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeInternal       = "INTERNAL_ERROR"
	ErrCodeDuplicate      = "DUPLICATE_ERROR"
)

// NewAPIError creates a new API error
func NewAPIError(code, message, details string) APIError {
	return APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// SendError sends a structured error response
func SendError(c *fiber.Ctx, status int, err APIError) error {
	return c.Status(status).JSON(ErrorResponse{
		Success: false,
		Error:   err,
	})
}

// Common error responses
func ValidationError(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusBadRequest, NewAPIError(
		ErrCodeValidation,
		message,
		"",
	))
}

func AuthenticationError(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusUnauthorized, NewAPIError(
		ErrCodeAuthentication,
		message,
		"",
	))
}

func NotFoundError(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusNotFound, NewAPIError(
		ErrCodeNotFound,
		message,
		"",
	))
}

func InternalError(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusInternalServerError, NewAPIError(
		ErrCodeInternal,
		message,
		"",
	))
}
