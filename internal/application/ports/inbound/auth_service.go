package inbound

import (
	"context"
)

type AuthenticationService interface {
	Register(ctx context.Context, req RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req LoginRequest) (*TokenPairResponse, error)
	Refresh(ctx context.Context, req RefreshRequest) (*TokenPairResponse, error)
	Logout(ctx context.Context, req LogoutRequest) error
	Validate(ctx context.Context, req ValidateRequest) (*ValidateResponse, error)
}
