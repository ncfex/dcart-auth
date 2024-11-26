package services

import (
	"context"
	"fmt"

	"github.com/ncfex/dcart-auth/internal/application/ports/inbound"
	tokenDomain "github.com/ncfex/dcart-auth/internal/domain/token"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type service struct {
	userService  inbound.UserSevice
	tokenService inbound.TokenService
}

func NewAuthService(
	userService inbound.UserSevice,
	tokenService inbound.TokenService,
) *service {
	return &service{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (s *service) Register(ctx context.Context, username, password string) (*userDomain.User, error) {
	user, err := s.userService.CreateUser(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (s *service) Login(ctx context.Context, username, password string) (*tokenDomain.TokenPair, error) {
	user, err := s.userService.ValidateWithCreds(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("validate with creds: %w", err)
	}

	createTokenParams := inbound.CreateTokenParams{
		UserID: user.ID,
	}
	tokenPair, err := s.tokenService.CreateTokenPair(ctx, createTokenParams)
	if err != nil {
		return nil, fmt.Errorf("create token pair: %w", err)
	}
	return &tokenDomain.TokenPair{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (s *service) Refresh(ctx context.Context, tokenString string) (*tokenDomain.TokenPair, error) {
	refreshToken, err := s.tokenService.ValidateRefreshToken(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("validate refresh token: %w", err)
	}

	user, err := s.userService.ValidateWithID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("validate with id: %w", err)
	}

	params := inbound.CreateTokenParams{
		UserID: user.ID,
	}
	accessTokenString, err := s.tokenService.CreateAccessToken(params)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}
	return &tokenDomain.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (s *service) Logout(ctx context.Context, tokenString string) error {
	fmt.Println(tokenString)
	if err := s.tokenService.RevokeRefreshToken(ctx, tokenString); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

func (s *service) Validate(ctx context.Context, token string) (*userDomain.User, error) {
	userID, err := s.tokenService.ValidateAccessToken(token)
	if err != nil {
		return nil, fmt.Errorf("validate access token: %w", err)
	}

	user, err := s.userService.ValidateWithID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("validate with id: %w", err)
	}
	return user, nil
}
