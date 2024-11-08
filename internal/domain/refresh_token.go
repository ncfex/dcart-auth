package domain

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	database "github.com/ncfex/dcart-auth/internal/infrastructure/database/sqlc"
)

var (
	ErrTokenNotFound      = errors.New("token not found")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenRevoked       = errors.New("token revoked")
	ErrTokenInvalid       = errors.New("token invalid")
	ErrTokenInvalidIssuer = errors.New("token invalid issuer")
	ErrTokenInvalidClaims = errors.New("token invalid claims")
	ErrTokenSigningFailed = errors.New("token signing failed")
)

type RefreshToken struct {
	Token     string     `json:"token"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	UserID    uuid.UUID  `json:"user_id"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

func (rt *RefreshToken) FromDB(dbToken *database.RefreshToken) {
	rt.Token = dbToken.Token
	rt.CreatedAt = dbToken.CreatedAt
	rt.UpdatedAt = dbToken.UpdatedAt
	rt.UserID = dbToken.UserID
	rt.ExpiresAt = dbToken.ExpiresAt
	if dbToken.RevokedAt.Valid {
		rt.RevokedAt = &dbToken.RevokedAt.Time
	}
}

func (rt *RefreshToken) ToDB() *database.RefreshToken {
	var revokedAt sql.NullTime
	if rt.RevokedAt != nil {
		revokedAt = sql.NullTime{
			Time:  *rt.RevokedAt,
			Valid: true,
		}
	}

	return &database.RefreshToken{
		Token:     rt.Token,
		CreatedAt: rt.CreatedAt,
		UpdatedAt: rt.UpdatedAt,
		UserID:    rt.UserID,
		ExpiresAt: rt.ExpiresAt,
		RevokedAt: revokedAt,
	}
}

func NewRefreshTokenFromDB(dbToken *database.RefreshToken) *RefreshToken {
	token := &RefreshToken{}
	token.FromDB(dbToken)
	return token
}
