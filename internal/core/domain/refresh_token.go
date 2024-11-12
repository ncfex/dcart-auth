package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
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
