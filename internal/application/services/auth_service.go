package services

import (
	"context"
	"fmt"

	"github.com/ncfex/dcart-auth/internal/application/ports/inbound"
)

type authService struct {
	userSvc  inbound.UserService
	tokenSvc inbound.TokenService
}

func NewAuthService(
	userSvc inbound.UserService,
	tokenSvc inbound.TokenService,
) inbound.AuthenticationService {
	return &authService{
		userSvc:  userSvc,
		tokenSvc: tokenSvc,
	}
}

func (as *authService) Register(ctx context.Context, req inbound.RegisterRequest) (*inbound.UserDTO, error) {
	user, err := as.userSvc.CreateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (as *authService) Login(ctx context.Context, req inbound.LoginRequest) (*inbound.TokenPairDTO, error) {
	user, err := as.userSvc.VerifyWithCreds(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("verify with creds: %w", err)
	}

	createTokenParams := inbound.CreateTokenParams{
		UserID: user.ID,
	}
	tokenPair, err := as.tokenSvc.CreateTokenPair(ctx, createTokenParams)
	if err != nil {
		return nil, fmt.Errorf("create token pair: %w", err)
	}
	return &inbound.TokenPairDTO{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (as *authService) Refresh(ctx context.Context, req inbound.TokenRequest) (*inbound.TokenDTO, error) {
	refreshToken, err := as.tokenSvc.ValidateRefreshToken(ctx, inbound.TokenRequest{Token: req.Token})
	if err != nil {
		return nil, fmt.Errorf("validate refresh token: %w", err)
	}

	user, err := as.userSvc.GetExistingUser(ctx, refreshToken.Subject)
	if err != nil {
		return nil, fmt.Errorf("get existing user : %w", err)
	}

	params := inbound.CreateTokenParams{
		UserID: user.ID,
	}
	accessToken, err := as.tokenSvc.CreateAccessToken(params)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}
	return &inbound.TokenDTO{
		Token: accessToken.Token,
	}, nil
}

func (as *authService) Logout(ctx context.Context, req inbound.TokenRequest) error {
	if err := as.tokenSvc.RevokeRefreshToken(ctx, inbound.TokenRequest{Token: req.Token}); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

func (as *authService) Validate(ctx context.Context, req inbound.TokenRequest) (*inbound.ValidateResponse, error) {
	validateResp, err := as.tokenSvc.ValidateAccessToken(req)
	if err != nil {
		return nil, fmt.Errorf("validate access token: %w", err)
	}

	user, err := as.userSvc.GetExistingUser(ctx, validateResp.Subject)
	if err != nil {
		return nil, fmt.Errorf("get existing user : %w", err)
	}
	return &inbound.ValidateResponse{
		Valid: true,
		User: inbound.UserDTO{
			ID:       user.ID,
			Username: user.Username,
		},
	}, nil
}
