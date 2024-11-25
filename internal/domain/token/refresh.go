package token

import (
	"errors"
	"time"
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

type Token string

type TokenPair struct {
	AccessToken  Token `json:"access_token"`
	RefreshToken Token `json:"refresh_token"`
}

type RefreshToken struct {
	Token     string     `json:"token"`
	UserID    string     `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}
