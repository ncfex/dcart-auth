package inbound

import (
	"context"
)

type CreateTokenParams struct {
	UserID string `json:"user_id" validate:"required"`
}

type ValidateTokenResponse struct {
	Subject string
}

type TokenService interface {
	CreateTokenPair(ctx context.Context, r CreateTokenParams) (*TokenPairDTO, error)

	// at
	CreateAccessToken(r CreateTokenParams) (*TokenDTO, error)
	ValidateAccessToken(r TokenRequest) (*ValidateTokenResponse, error)

	// rt
	CreateRefreshToken(ctx context.Context, r CreateTokenParams) (*TokenDTO, error)
	ValidateRefreshToken(ctx context.Context, r TokenRequest) (*ValidateTokenResponse, error)
	RevokeRefreshToken(ctx context.Context, r TokenRequest) error
}
