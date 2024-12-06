package command

import (
	"context"

	"github.com/ncfex/dcart-auth/internal/application/ports/types"
)

type AuthenticateUserCommand struct {
	Username string
	Password string
}

type RegisterUserCommand struct {
	Username string
	Password string
}

type ChangePasswordCommand struct {
	UserID      string
	OldPassword string
	NewPassword string
}

type UserCommandPort interface {
	RegisterUser(ctx context.Context, cmd RegisterUserCommand) (*types.UserResponse, error)
	AuthenticateUser(ctx context.Context, cmd AuthenticateUserCommand) (*types.UserResponse, error)
	ChangePassword(ctx context.Context, cmd ChangePasswordCommand) error
}
