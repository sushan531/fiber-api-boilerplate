package helpers

import (
	"context"
	"fiber-api/api/models"
	"fmt"

	"github.com/google/uuid"
	"github.com/sushan531/auth-sqlc/generated"
)

func CreateJWTClaims(queries *generated.Queries, context context.Context, userId uuid.UUID, deviceFingerprint string) (*models.JWTClaims, error) {
	profile, err := queries.GetUserProfile(context, userId)
	if err != nil {
		return nil, err
	}
	claims := &models.JWTClaims{
		UserID:            profile.UserProfileID.String(),
		UserEmail:         profile.UserEmail,
		Role:              "admin",
		DeviceFingerprint: deviceFingerprint,
	}
	return claims, nil
}

// extractUserIDFromClaims safely extracts and parses user_id from token claims
func ExtractUserIdFromMapObj(claims map[string]interface{}) (uuid.UUID, error) {
	raw, exists := claims["user_id"]
	if !exists {
		return uuid.Nil, fmt.Errorf("user_id missing in claims")
	}

	userIDStr, ok := raw.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("user_id is not a string")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user_id format")
	}

	return userID, nil
}
