package services

import (
	"context"
	"fmt"

	"github.com/ncfex/dcart-auth/internal/application/ports/inbound"
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

func (as *authService) Register(ctx context.Context, req inbound.RegisterRequest) (*inbound.UserResponse, error) {
	user, err := as.userSvc.CreateUser(ctx, req.Username, req.Password)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &inbound.UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

func (as *authService) Login(ctx context.Context, req inbound.LoginRequest) (*inbound.TokenPairResponse, error) {
	user, err := as.userSvc.ValidateWithCreds(ctx, req.Username, req.Password)
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
	return &inbound.TokenPairResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (as *authService) Refresh(ctx context.Context, req inbound.RefreshRequest) (*inbound.TokenPairResponse, error) {
	refreshToken, err := as.tokenSvc.ValidateRefreshToken(ctx, req.TokenString)
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
	return &inbound.TokenPairResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (as *authService) Logout(ctx context.Context, req inbound.LogoutRequest) error {
	if err := as.tokenSvc.RevokeRefreshToken(ctx, req.TokenString); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

func (as *authService) Validate(ctx context.Context, req inbound.ValidateRequest) (*inbound.ValidateResponse, error) {
	userID, err := as.tokenSvc.ValidateAccessToken(req.TokenString)
	if err != nil {
		return nil, fmt.Errorf("validate access token: %w", err)
	}

	user, err := as.userSvc.ValidateWithID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("validate with id: %w", err)
	}
	return &inbound.ValidateResponse{
		Valid: true,
		User: inbound.UserResponse{
			ID:       user.ID,
			Username: user.Username,
		},
	}, nil
}
