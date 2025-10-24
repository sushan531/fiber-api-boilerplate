package validators

import (
	"fiber-api/api/models"
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResult holds validation results
type ValidationResult struct {
	IsValid bool              `json:"is_valid"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

// ValidateSignUp validates user registration input
func ValidateSignUp(input models.SignUp) ValidationResult {
	var errors []ValidationError

	// Email validation
	if input.UserEmail == "" {
		errors = append(errors, ValidationError{
			Field:   "user_email",
			Message: "Email is required",
		})
	} else if !isValidEmail(input.UserEmail) {
		errors = append(errors, ValidationError{
			Field:   "user_email",
			Message: "Invalid email format",
		})
	}

	// Password validation
	if input.Password == "" {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password is required",
		})
	} else if len(input.Password) < 8 {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must be at least 8 characters long",
		})
	}

	// Full name validation
	if input.FullName == "" {
		errors = append(errors, ValidationError{
			Field:   "full_name",
			Message: "Full name is required",
		})
	} else if len(strings.TrimSpace(input.FullName)) < 2 {
		errors = append(errors, ValidationError{
			Field:   "full_name",
			Message: "Full name must be at least 2 characters long",
		})
	}

	// User role validation (if provided)
	if input.UserRole != "" {
		validRoles := []string{"admin", "user", "moderator"}
		if !contains(validRoles, strings.ToLower(input.UserRole)) {
			errors = append(errors, ValidationError{
				Field:   "user_role",
				Message: fmt.Sprintf("Invalid role. Must be one of: %s", strings.Join(validRoles, ", ")),
			})
		}
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// ValidateLogin validates user login input
func ValidateLogin(input models.Login) ValidationResult {
	var errors []ValidationError

	// Email validation
	if input.UserEmail == "" {
		errors = append(errors, ValidationError{
			Field:   "user_email",
			Message: "Email is required",
		})
	} else if !isValidEmail(input.UserEmail) {
		errors = append(errors, ValidationError{
			Field:   "user_email",
			Message: "Invalid email format",
		})
	}

	// Password validation
	if input.Password == "" {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password is required",
		})
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// isValidEmail validates email format using regex
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
