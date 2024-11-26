package inbound

import (
	"context"

	"github.com/ncfex/dcart-auth/internal/domain/token"
)

type CreateTokenParams struct {
	UserID string
}

type TokenService interface {
	CreateTokenPair(ctx context.Context, params CreateTokenParams) (tokenPair token.TokenPair, err error)

	// at
	CreateAccessToken(params CreateTokenParams) (tokenString string, err error)
	ValidateAccessToken(tokenString string) (subjectString string, err error)

	// rt
	CreateRefreshToken(ctx context.Context, params CreateTokenParams) (refreshToken *token.RefreshToken, err error)
	ValidateRefreshToken(ctx context.Context, tokenString string) (refreshToken *token.RefreshToken, err error)
	RevokeRefreshToken(ctx context.Context, tokenString string) error
}
