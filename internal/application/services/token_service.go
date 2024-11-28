package services

import (
	"context"
	"fmt"

	"github.com/ncfex/dcart-auth/internal/application/ports/inbound"
	"github.com/ncfex/dcart-auth/internal/application/ports/outbound"
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

func (ts *tokenService) CreateTokenPair(ctx context.Context, r inbound.CreateTokenParams) (*inbound.TokenPairDTO, error) {
	accessToken, err := ts.CreateAccessToken(r)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	refreshToken, err := ts.CreateRefreshToken(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}
	return &inbound.TokenPairDTO{
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
	}, nil
}

// at
func (ts *tokenService) CreateAccessToken(r inbound.CreateTokenParams) (*inbound.TokenDTO, error) {
	accessTokenString, err := ts.accessTokenGen.Generate(r.UserID)
	if err != nil {
		return nil, fmt.Errorf("access token generate: %w", err)
	}
	return &inbound.TokenDTO{
		Token: accessTokenString,
	}, nil
}

func (ts *tokenService) ValidateAccessToken(r inbound.TokenRequest) (*inbound.ValidateTokenResponse, error) {
	subjectString, err := ts.accessTokenGen.Validate(r.Token)
	if err != nil {
		return nil, fmt.Errorf("access token validate: %w", err)
	}
	return &inbound.ValidateTokenResponse{
		Subject: subjectString,
	}, nil
}

// rt
func (ts *tokenService) CreateRefreshToken(ctx context.Context, r inbound.CreateTokenParams) (*inbound.TokenDTO, error) {
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
	return &inbound.TokenDTO{
		Token: refreshToken.Token,
	}, nil
}

func (ts *tokenService) ValidateRefreshToken(ctx context.Context, r inbound.TokenRequest) (*inbound.ValidateTokenResponse, error) {
	refreshToken, err := ts.tokenRepo.GetByToken(ctx, r.Token)
	if err != nil {
		return nil, fmt.Errorf("get token string: %w", err)
	}

	if err := refreshToken.IsValid(); err != nil {
		return nil, fmt.Errorf("is valid: %w", err)
	}
	return &inbound.ValidateTokenResponse{
		Subject: refreshToken.UserID,
	}, nil
}

func (ts *tokenService) RevokeRefreshToken(ctx context.Context, r inbound.TokenRequest) error {
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
