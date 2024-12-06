package inbound

import (
	"context"

	"github.com/ncfex/dcart-auth/internal/application/commands"
	"github.com/ncfex/dcart-auth/internal/application/ports/types"
)

type UserWriteModel interface {
	AuthenticateUser(ctx context.Context, cmd commands.AuthenticateUserCommand) (*types.UserResponse, error)
	RegisterUser(ctx context.Context, cmd commands.RegisterUserCommand) (*types.UserResponse, error)
	ChangePassword(ctx context.Context, cmd commands.ChangePasswordCommand) error
}
