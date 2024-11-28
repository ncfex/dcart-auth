package inbound

import (
	"context"
)

type AuthenticationService interface {
	Register(ctx context.Context, req RegisterRequest) (*UserDTO, error)
	Login(ctx context.Context, req LoginRequest) (*TokenPairDTO, error)
	Refresh(ctx context.Context, req TokenRequest) (*TokenDTO, error)
	Logout(ctx context.Context, req TokenRequest) error
	Validate(ctx context.Context, req TokenRequest) (*ValidateResponse, error)
}
