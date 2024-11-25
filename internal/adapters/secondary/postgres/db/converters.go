package db

import (
	"time"

	tokenDomain "github.com/ncfex/dcart-auth/internal/domain/token"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

func ToUserDomain(dbUser *User) *userDomain.User {
	return &userDomain.User{
		ID:           dbUser.ID.String(),
		Username:     dbUser.Username,
		PasswordHash: dbUser.PasswordHash,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}
}

func ToRefreshTokenDomain(dbToken *RefreshToken) *tokenDomain.RefreshToken {
	var revokedAt *time.Time
	if dbToken.RevokedAt.Valid {
		revokedAt = &dbToken.RevokedAt.Time
	}

	return &tokenDomain.RefreshToken{
		Token:     dbToken.Token,
		UserID:    dbToken.UserID.String(),
		CreatedAt: dbToken.CreatedAt,
		UpdatedAt: dbToken.UpdatedAt,
		ExpiresAt: dbToken.ExpiresAt,
		RevokedAt: revokedAt,
	}
}
