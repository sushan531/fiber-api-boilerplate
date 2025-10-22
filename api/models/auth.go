package models

// SignUp represents the request body for user registration
type SignUp struct {
	UserEmail string `json:"user_email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FullName  string `json:"full_name" binding:"required"`
	UserRole  string `json:"user_role" binding:"required"`
	Address   string `json:"address"`
}

// Login represents the request body for user authentication
type Login struct {
	UserEmail string `json:"user_email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// JWTClaims represents the claims structure for JWT tokens
type JWTClaims struct {
	UserID            string `json:"user_id"`
	UserEmail         string `json:"user_email"`
	Role              string `json:"role"`
	DeviceFingerprint string `json:"device_fingerprint"`
}

// ToMap converts JWTClaims struct to map[string]interface{} for JWT token generation
func (j *JWTClaims) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"user_id":            j.UserID,
		"user_email":         j.UserEmail,
		"role":               j.Role,
		"device_fingerprint": j.DeviceFingerprint,
	}
}
