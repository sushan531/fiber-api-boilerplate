package presenter

import (
	"github.com/sushan531/auth-sqlc/generated"
	"github.com/sushan531/jwk-auth/service"
)

// BaseResponse represents the standard API response structure
type BaseResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// SignUpResponse represents user registration response data
type SignUpResponse struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	CreatedAt string `json:"created_at"`
}

// SignInResponse represents authentication response data
type SignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
}

// SignUpSuccessResponse creates a standardized signup success response
func SignUpSuccessResponse(data generated.Auth) BaseResponse {
	return BaseResponse{
		Success: true,
		Data: SignUpResponse{
			UserID: data.UserProfileID.String(),
			Email:  data.UserEmail,
			//FullName:  data.FullName,
			//CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
		Message: "User registered successfully",
	}
}

// SignInSuccessResponse creates a standardized signin success response
func SignInSuccessResponse(data service.TokenPair) BaseResponse {
	return BaseResponse{
		Success: true,
		Data: SignInResponse{
			AccessToken:  data.AccessToken,
			RefreshToken: data.RefreshToken,
			TokenType:    "Bearer",
		},
		Message: "Authentication successful",
	}
}
