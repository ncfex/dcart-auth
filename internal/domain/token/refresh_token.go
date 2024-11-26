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

type RefreshToken struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
	RevokedAt time.Time `json:"revoked_at,omitempty"`
}

func NewRefreshToken(tokenString string, userID string) (*RefreshToken, error) {
	if tokenString == "" || userID == "" {
		return nil, ErrTokenInvalid
	}

	now := time.Now()
	return &RefreshToken{
		Token:     tokenString,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now,
	}, nil
}

func (rt *RefreshToken) Revoke() {
	now := time.Now()
	rt.RevokedAt = now
	rt.UpdatedAt = now
}

func (rt *RefreshToken) Expire() {
	now := time.Now()
	rt.ExpiresAt = now
	rt.UpdatedAt = now
}

func (rt *RefreshToken) IsValid() error {
	if !rt.RevokedAt.IsZero() {
		return ErrTokenRevoked
	}

	if rt.ExpiresAt.Before(time.Now()) {
		return ErrTokenExpired
	}

	return nil
}
