package services

import (
	"context"
	"fmt"

	"github.com/ncfex/dcart-auth/internal/application/ports/inbound"
	tokenDomain "github.com/ncfex/dcart-auth/internal/domain/token"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type authService struct {
	userSvc  inbound.UserSevice
	tokenSvc inbound.TokenService
}

func NewAuthService(
	userSvc inbound.UserSevice,
	tokenSvc inbound.TokenService,
) inbound.AuthenticationService {
	return &authService{
		userSvc:  userSvc,
		tokenSvc: tokenSvc,
	}
}

func (as *authService) Register(ctx context.Context, username, password string) (*userDomain.User, error) {
	user, err := as.userSvc.CreateUser(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (as *authService) Login(ctx context.Context, username, password string) (*tokenDomain.TokenPair, error) {
	user, err := as.userSvc.ValidateWithCreds(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("validate with creds: %w", err)
	}

	createTokenParams := inbound.CreateTokenParams{
		UserID: user.ID,
	}
	tokenPair, err := as.tokenSvc.CreateTokenPair(ctx, createTokenParams)
	if err != nil {
		return nil, fmt.Errorf("create token pair: %w", err)
	}
	return &tokenDomain.TokenPair{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (as *authService) Refresh(ctx context.Context, tokenString string) (*tokenDomain.TokenPair, error) {
	refreshToken, err := as.tokenSvc.ValidateRefreshToken(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("validate refresh token: %w", err)
	}

	user, err := as.userSvc.ValidateWithID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("validate with id: %w", err)
	}

	params := inbound.CreateTokenParams{
		UserID: user.ID,
	}
	accessTokenString, err := as.tokenSvc.CreateAccessToken(params)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}
	return &tokenDomain.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (as *authService) Logout(ctx context.Context, tokenString string) error {
	if err := as.tokenSvc.RevokeRefreshToken(ctx, tokenString); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

func (as *authService) Validate(ctx context.Context, token string) (*userDomain.User, error) {
	userID, err := as.tokenSvc.ValidateAccessToken(token)
	if err != nil {
		return nil, fmt.Errorf("validate access token: %w", err)
	}

	user, err := as.userSvc.ValidateWithID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("validate with id: %w", err)
	}
	return user, nil
}
