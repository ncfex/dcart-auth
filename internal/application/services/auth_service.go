package services

import (
	"context"
	"fmt"

	"github.com/ncfex/dcart-auth/internal/application/ports/primary/command"
	"github.com/ncfex/dcart-auth/internal/application/ports/primary/query"
	"github.com/ncfex/dcart-auth/internal/application/ports/primary/services"
	"github.com/ncfex/dcart-auth/internal/application/ports/types"
)

type authService struct {
	userCommandHandler command.UserCommandPort
	userQueryHandler   query.UserQueryPort
	tokenSvc           services.TokenService
}

func NewAuthService(
	userCommandHandler command.UserCommandPort,
	userQueryHandler query.UserQueryPort,
	tokenSvc services.TokenService,
) services.AuthenticationService {
	return &authService{
		userCommandHandler: userCommandHandler,
		userQueryHandler:   userQueryHandler,
		tokenSvc:           tokenSvc,
	}
}

func (as *authService) Register(ctx context.Context, req types.RegisterRequest) (*types.UserResponse, error) {
	registerCommand := command.RegisterUserCommand{
		Username: req.Username,
		Password: req.Password,
	}
	return as.userCommandHandler.RegisterUser(ctx, registerCommand)
}

func (as *authService) Login(ctx context.Context, req types.LoginRequest) (*types.TokenPairResponse, error) {
	authenticateCmd := command.AuthenticateUserCommand{
		Username: req.Username,
		Password: req.Password,
	}
	authenticatedUser, err := as.userCommandHandler.AuthenticateUser(ctx, authenticateCmd)
	if err != nil {
		return nil, fmt.Errorf("create token pair: %w", err)
	}

	createTokenParams := types.CreateTokenParams{
		UserID: authenticatedUser.ID,
	}
	tokenPair, err := as.tokenSvc.CreateTokenPair(ctx, createTokenParams)
	if err != nil {
		return nil, fmt.Errorf("create token pair: %w", err)
	}
	return &types.TokenPairResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (as *authService) ChangePassword(ctx context.Context, req types.ChangePasswordRequest) error {
	changePasswordCmd := command.ChangePasswordCommand{
		UserID:      req.UserID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}
	if err := as.userCommandHandler.ChangePassword(ctx, changePasswordCmd); err != nil {
		return fmt.Errorf("change password: %w", err)
	}
	return nil
}

func (as *authService) Refresh(ctx context.Context, req types.TokenRequest) (*types.TokenResponse, error) {
	refreshToken, err := as.tokenSvc.ValidateRefreshToken(ctx, types.TokenRequest{Token: req.Token})
	if err != nil {
		return nil, fmt.Errorf("validate refresh token: %w", err)
	}

	// just to check user exists
	getUserByIdQuery := query.GetUserByIDQuery{
		UserID: refreshToken.Subject,
	}
	_, err = as.userQueryHandler.GetUserByID(ctx, getUserByIdQuery)
	if err != nil {
		return nil, fmt.Errorf("get existing user : %w", err)
	}

	params := types.CreateTokenParams{
		UserID: refreshToken.Subject,
	}
	accessToken, err := as.tokenSvc.CreateAccessToken(params)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}
	return &types.TokenResponse{
		Token: accessToken.Token,
	}, nil
}

func (as *authService) Logout(ctx context.Context, req types.TokenRequest) error {
	if err := as.tokenSvc.RevokeRefreshToken(ctx, types.TokenRequest{Token: req.Token}); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

func (as *authService) Validate(ctx context.Context, req types.TokenRequest) (*types.ValidateResponse, error) {
	validateResp, err := as.tokenSvc.ValidateAccessToken(req)
	if err != nil {
		return nil, fmt.Errorf("validate access token: %w", err)
	}

	getUserByIdQuery := query.GetUserByIDQuery{
		UserID: validateResp.Subject,
	}
	user, err := as.userQueryHandler.GetUserByID(ctx, getUserByIdQuery)
	if err != nil {
		return nil, fmt.Errorf("get existing user : %w", err)
	}
	return &types.ValidateResponse{
		Valid: true,
		User: types.UserResponse{
			ID:       user.ID,
			Username: user.Username,
		},
	}, nil
}
