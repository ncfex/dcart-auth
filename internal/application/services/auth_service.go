package services

import (
	"context"
	"fmt"

	"github.com/ncfex/dcart-auth/internal/application/commands"
	"github.com/ncfex/dcart-auth/internal/application/ports/inbound"
	"github.com/ncfex/dcart-auth/internal/application/ports/types"
	"github.com/ncfex/dcart-auth/internal/application/queries"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type authService struct {
	userCommandHandler inbound.UserCommandHandler
	userQueryHandler   inbound.UserQueryHandler
	tokenSvc           inbound.TokenService
}

func NewAuthService(
	userCommandHandler inbound.UserCommandHandler,
	userQueryHandler inbound.UserQueryHandler,
	tokenSvc inbound.TokenService,
) inbound.AuthenticationService {
	return &authService{
		userCommandHandler: userCommandHandler,
		userQueryHandler:   userQueryHandler,
		tokenSvc:           tokenSvc,
	}
}

func (as *authService) Register(ctx context.Context, req types.RegisterRequest) (*types.UserResponse, error) {
	if _, err := as.userQueryHandler.HandleGetUserByUsername(ctx, queries.GetUserByUsernameQuery{
		Username: req.Username,
	}); err == nil {
		return nil, fmt.Errorf("get user username: %w", userDomain.ErrUserAlreadyExists)
	} else if err != userDomain.ErrUserNotFound {
		return nil, fmt.Errorf("get user username: %w", err)
	}

	registerCommand := commands.RegisterUserCommand{
		Username: req.Username,
		Password: req.Password,
	}
	user, err := as.userCommandHandler.HandleRegisterUser(ctx, registerCommand)
	if err != nil {
		return nil, fmt.Errorf("register user: %w", err)
	}

	return &types.UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

func (as *authService) Login(ctx context.Context, req types.LoginRequest) (*types.TokenPairResponse, error) {
	existingUser, err := as.userQueryHandler.HandleGetUserByUsername(ctx, queries.GetUserByUsernameQuery{
		Username: req.Username,
	})
	if err != nil {
		return nil, fmt.Errorf("get user username: %w", err)
	}
	if existingUser == nil {
		return nil, fmt.Errorf("user not found")
	}

	correct := existingUser.Authenticate(req.Password)
	if !correct {
		return nil, fmt.Errorf("wrong password")
	}

	createTokenParams := types.CreateTokenParams{
		UserID: existingUser.GetID(),
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

func (as *authService) Refresh(ctx context.Context, req types.TokenRequest) (*types.TokenResponse, error) {
	refreshToken, err := as.tokenSvc.ValidateRefreshToken(ctx, types.TokenRequest{Token: req.Token})
	if err != nil {
		return nil, fmt.Errorf("validate refresh token: %w", err)
	}

	getUserByIdQuery := queries.GetUserByIDQuery{
		UserID: refreshToken.Subject,
	}
	user, err := as.userQueryHandler.HandleGetUserById(ctx, getUserByIdQuery)
	if err != nil {
		return nil, fmt.Errorf("get existing user : %w", err)
	}

	params := types.CreateTokenParams{
		UserID: user.ID,
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

	getUserByIdQuery := queries.GetUserByIDQuery{
		UserID: validateResp.Subject,
	}
	user, err := as.userQueryHandler.HandleGetUserById(ctx, getUserByIdQuery)
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
