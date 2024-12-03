package services

import (
	"context"
	"fmt"

	"github.com/ncfex/dcart-auth/internal/application/ports/inbound"
	"github.com/ncfex/dcart-auth/internal/application/ports/outbound"
	"github.com/ncfex/dcart-auth/internal/application/ports/types"
	tokenDomain "github.com/ncfex/dcart-auth/internal/domain/token"
)

type tokenService struct {
	accessTokenGen  outbound.TokenGeneratorValidator
	refreshTokenGen outbound.TokenGenerator
	tokenRepo       outbound.TokenRepository
}

func NewTokenService(
	accessTokenGen outbound.TokenGeneratorValidator,
	refreshTokenGen outbound.TokenGenerator,
	tokenRepo outbound.TokenRepository,
) inbound.TokenService {
	return &tokenService{
		accessTokenGen:  accessTokenGen,
		refreshTokenGen: refreshTokenGen,
		tokenRepo:       tokenRepo,
	}
}

func (ts *tokenService) CreateTokenPair(ctx context.Context, r types.CreateTokenParams) (*types.TokenPairResponse, error) {
	accessToken, err := ts.CreateAccessToken(r)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	refreshToken, err := ts.CreateRefreshToken(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}
	return &types.TokenPairResponse{
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
	}, nil
}

// at
func (ts *tokenService) CreateAccessToken(r types.CreateTokenParams) (*types.TokenResponse, error) {
	accessTokenString, err := ts.accessTokenGen.Generate(r.UserID)
	if err != nil {
		return nil, fmt.Errorf("access token generate: %w", err)
	}
	return &types.TokenResponse{
		Token: accessTokenString,
	}, nil
}

func (ts *tokenService) ValidateAccessToken(r types.TokenRequest) (*types.ValidateTokenResponse, error) {
	subjectString, err := ts.accessTokenGen.Validate(r.Token)
	if err != nil {
		return nil, fmt.Errorf("access token validate: %w", err)
	}
	return &types.ValidateTokenResponse{
		Subject: subjectString,
	}, nil
}

// rt
func (ts *tokenService) CreateRefreshToken(ctx context.Context, r types.CreateTokenParams) (*types.TokenResponse, error) {
	refreshTokenString, err := ts.refreshTokenGen.Generate("")
	if err != nil {
		return nil, fmt.Errorf("refresh token generate: %w", err)
	}

	refreshToken, err := tokenDomain.NewRefreshToken(refreshTokenString, r.UserID)
	if err != nil {
		return nil, fmt.Errorf("new refresh token: %w", err)
	}

	if err := ts.tokenRepo.Add(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("store token: %w", err)
	}
	return &types.TokenResponse{
		Token: refreshToken.Token,
	}, nil
}

func (ts *tokenService) ValidateRefreshToken(ctx context.Context, r types.TokenRequest) (*types.ValidateTokenResponse, error) {
	refreshToken, err := ts.tokenRepo.GetByToken(ctx, r.Token)
	if err != nil {
		return nil, fmt.Errorf("get token string: %w", err)
	}

	if err := refreshToken.IsValid(); err != nil {
		return nil, fmt.Errorf("is valid: %w", err)
	}
	return &types.ValidateTokenResponse{
		Subject: refreshToken.UserID,
	}, nil
}

func (ts *tokenService) RevokeRefreshToken(ctx context.Context, r types.TokenRequest) error {
	token, err := ts.tokenRepo.GetByToken(ctx, r.Token)
	if err != nil {
		return fmt.Errorf("get token string: %w", err)
	}

	token.Revoke()
	if err := ts.tokenRepo.Save(ctx, token); err != nil {
		return fmt.Errorf("save token: %w", err)
	}
	return nil
}
