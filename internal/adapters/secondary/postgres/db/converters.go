package db

import (
	"time"

	tokenDomain "github.com/ncfex/dcart-auth/internal/domain/token"
)

// todo improve
func ToRefreshTokenDomain(dbToken *RefreshToken) *tokenDomain.RefreshToken {
	var revokedAt time.Time
	if dbToken.RevokedAt.Valid {
		revokedAt = dbToken.RevokedAt.Time
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
