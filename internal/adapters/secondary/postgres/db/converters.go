package db

import (
	"database/sql"
	"time"

	tokenDomain "github.com/ncfex/dcart-auth/internal/domain/token"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

func ToUserDomain(dbUser *User) *userDomain.User {
	return &userDomain.User{
		ID:           dbUser.ID,
		Username:     dbUser.Username,
		PasswordHash: dbUser.PasswordHash,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}
}

func ToUserDB(domainUser *userDomain.User) *User {
	return &User{
		ID:           domainUser.ID,
		Username:     domainUser.Username,
		PasswordHash: domainUser.PasswordHash,
		CreatedAt:    domainUser.CreatedAt,
		UpdatedAt:    domainUser.UpdatedAt,
	}
}

func ToRefreshTokenDomain(dbToken *RefreshToken) *tokenDomain.RefreshToken {
	var revokedAt *time.Time
	if dbToken.RevokedAt.Valid {
		revokedAt = &dbToken.RevokedAt.Time
	}

	return &tokenDomain.RefreshToken{
		Token:     dbToken.Token,
		CreatedAt: dbToken.CreatedAt,
		UpdatedAt: dbToken.UpdatedAt,
		UserID:    dbToken.UserID,
		ExpiresAt: dbToken.ExpiresAt,
		RevokedAt: revokedAt,
	}
}

func ToRefreshTokenDB(domainToken *tokenDomain.RefreshToken) *RefreshToken {
	var revokedAt sql.NullTime
	if domainToken.RevokedAt != nil {
		revokedAt = sql.NullTime{
			Time:  *domainToken.RevokedAt,
			Valid: true,
		}
	}

	return &RefreshToken{
		Token:     domainToken.Token,
		CreatedAt: domainToken.CreatedAt,
		UpdatedAt: domainToken.UpdatedAt,
		UserID:    domainToken.UserID,
		ExpiresAt: domainToken.ExpiresAt,
		RevokedAt: revokedAt,
	}
}
