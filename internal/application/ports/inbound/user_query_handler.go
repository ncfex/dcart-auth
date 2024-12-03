package inbound

import (
	"context"

	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type UserQueryHandler interface {
	HandleGetUserById(ctx context.Context, cmd userDomain.GetUserByIDQuery) (*userDomain.User, error)
	HandleGetUserByUsername(ctx context.Context, cmd userDomain.GetUserByUsernameQuery) (*userDomain.User, error)
}
