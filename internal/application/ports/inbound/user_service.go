package inbound

import (
	"context"
)

type UserService interface {
	CreateUser(ctx context.Context, req RegisterRequest) (*UserDTO, error)
	VerifyWithCreds(ctx context.Context, req LoginRequest) (*UserDTO, error)
	GetExistingUser(ctx context.Context, userID string) (*UserDTO, error)
}
