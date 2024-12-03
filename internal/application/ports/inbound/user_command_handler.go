package inbound

import (
	"context"

	"github.com/ncfex/dcart-auth/internal/application/commands"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type UserCommandHandler interface {
	HandleRegisterUser(ctx context.Context, cmd commands.RegisterUserCommand) (*userDomain.User, error)
}
