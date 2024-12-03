package inbound

import (
	"context"

	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type UserCommandHandler interface {
	HandleRegisterUser(ctx context.Context, cmd userDomain.RegisterUserCommand) (*userDomain.User, error)
}
