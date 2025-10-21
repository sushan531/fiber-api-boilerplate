package presenter

import (
	"github.com/sushan531/auth-sqlc/generated"
)

// UserProfileResponse represents user profile response data
type UserProfileResponse struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	Role     string `json:"role,omitempty"`
	Email    string `json:"email"`
}

// UserProfileFetchResponse creates a standardized user profile response
func UserProfileFetchResponse(data generated.GetUserProfileRow) BaseResponse {
	return BaseResponse{
		Success: true,
		Data: UserProfileResponse{
			UserID:   data.UserProfileID.String(),
			FullName: data.FullName,
			Role:     data.UserRole.String,
			Email:    data.UserEmail,
		},
		Message: "Profile retrieved successfully",
	}
}
