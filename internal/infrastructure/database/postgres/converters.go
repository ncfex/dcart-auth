package postgres

import (
	"database/sql"
	"time"

	tokenDomain "github.com/ncfex/dcart-auth/internal/core/domain/token"
	userDomain "github.com/ncfex/dcart-auth/internal/core/domain/user"
	db "github.com/ncfex/dcart-auth/internal/infrastructure/database/postgres/sqlc"
)

func ToUserDomain(dbUser *db.User) *userDomain.User {
	return &userDomain.User{
		ID:           dbUser.ID,
		Username:     dbUser.Username,
		PasswordHash: dbUser.PasswordHash,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}
}

func ToUserDB(domainUser *userDomain.User) *db.User {
	return &db.User{
		ID:           domainUser.ID,
		Username:     domainUser.Username,
		PasswordHash: domainUser.PasswordHash,
		CreatedAt:    domainUser.CreatedAt,
		UpdatedAt:    domainUser.UpdatedAt,
	}
}

func ToRefreshTokenDomain(dbToken *db.RefreshToken) *tokenDomain.RefreshToken {
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

func ToRefreshTokenDB(domainToken *tokenDomain.RefreshToken) *db.RefreshToken {
	var revokedAt sql.NullTime
	if domainToken.RevokedAt != nil {
		revokedAt = sql.NullTime{
			Time:  *domainToken.RevokedAt,
			Valid: true,
		}
	}

	return &db.RefreshToken{
		Token:     domainToken.Token,
		CreatedAt: domainToken.CreatedAt,
		UpdatedAt: domainToken.UpdatedAt,
		UserID:    domainToken.UserID,
		ExpiresAt: domainToken.ExpiresAt,
		RevokedAt: revokedAt,
	}
}
