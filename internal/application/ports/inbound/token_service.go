package inbound

import (
	"context"

	"github.com/ncfex/dcart-auth/internal/application/ports/types"
)

type TokenService interface {
	CreateTokenPair(ctx context.Context, r types.CreateTokenParams) (*types.TokenPairResponse, error)

	// at
	CreateAccessToken(r types.CreateTokenParams) (*types.TokenResponse, error)
	ValidateAccessToken(r types.TokenRequest) (*types.ValidateTokenResponse, error)

	// rt
	CreateRefreshToken(ctx context.Context, r types.CreateTokenParams) (*types.TokenResponse, error)
	ValidateRefreshToken(ctx context.Context, r types.TokenRequest) (*types.ValidateTokenResponse, error)
	RevokeRefreshToken(ctx context.Context, r types.TokenRequest) error
}
