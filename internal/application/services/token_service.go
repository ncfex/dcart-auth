package services

import (
	"context"
	"fmt"

	"github.com/ncfex/dcart-auth/internal/application/ports/inbound"
	"github.com/ncfex/dcart-auth/internal/application/ports/outbound"
	"github.com/ncfex/dcart-auth/internal/domain/token"
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

func (ts *tokenService) CreateTokenPair(ctx context.Context, params inbound.CreateTokenParams) (tokenDomain.TokenPair, error) {
	accessTokenString, err := ts.CreateAccessToken(params)
	if err != nil {
		return token.TokenPair{}, fmt.Errorf("create access token: %w", err)
	}

	refreshToken, err := ts.CreateRefreshToken(ctx, params)
	if err != nil {
		return token.TokenPair{}, fmt.Errorf("create refresh token: %w", err)
	}
	return tokenDomain.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken.Token,
	}, nil
}

// at
func (ts *tokenService) CreateAccessToken(params inbound.CreateTokenParams) (string, error) {
	accessTokenString, err := ts.accessTokenGen.Generate(params.UserID)
	if err != nil {
		return "", fmt.Errorf("access token generate: %w", err)
	}
	return accessTokenString, nil
}

func (ts *tokenService) ValidateAccessToken(tokenString string) (string, error) {
	subjectString, err := ts.accessTokenGen.Validate(tokenString)
	if err != nil {
		return "", fmt.Errorf("access token validate: %w", err)
	}
	return subjectString, nil
}

// rt
func (ts *tokenService) CreateRefreshToken(ctx context.Context, params inbound.CreateTokenParams) (*tokenDomain.RefreshToken, error) {
	refreshTokenString, err := ts.refreshTokenGen.Generate("")
	if err != nil {
		return nil, fmt.Errorf("refresh token generate: %w", err)
	}

	refreshToken, err := tokenDomain.NewRefreshToken(refreshTokenString, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("new refresh token: %w", err)
	}

	if err := ts.tokenRepo.Add(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("store token: %w", err)
	}
	return refreshToken, nil
}

func (ts *tokenService) ValidateRefreshToken(ctx context.Context, tokenString string) (*token.RefreshToken, error) {
	refreshToken, err := ts.tokenRepo.GetByToken(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("get token string: %w", err)
	}

	if err := refreshToken.IsValid(); err != nil {
		return nil, fmt.Errorf("is valid: %w", err)
	}
	return refreshToken, nil
}

func (ts *tokenService) RevokeRefreshToken(ctx context.Context, tokenString string) error {
	token, err := ts.tokenRepo.GetByToken(ctx, tokenString)
	if err != nil {
		return fmt.Errorf("get token string: %w", err)
	}

	token.Revoke()
	if err := ts.tokenRepo.Save(ctx, token); err != nil {
		return fmt.Errorf("save token: %w", err)
	}
	return nil
}
