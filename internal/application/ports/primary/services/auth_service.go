package services

import (
	"context"

	"github.com/ncfex/dcart-auth/internal/application/ports/types"
)

type AuthenticationService interface {
	Register(ctx context.Context, req types.RegisterRequest) (*types.UserResponse, error)
	Login(ctx context.Context, req types.LoginRequest) (*types.TokenPairResponse, error)
	ChangePassword(ctx context.Context, req types.ChangePasswordRequest) error
	Refresh(ctx context.Context, req types.TokenRequest) (*types.TokenResponse, error)
	Logout(ctx context.Context, req types.TokenRequest) error
	Validate(ctx context.Context, req types.TokenRequest) (*types.ValidateResponse, error)
}
