package inbound

import (
	"context"

	"github.com/ncfex/dcart-auth/internal/application/queries"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type UserQueryHandler interface {
	HandleGetUserById(ctx context.Context, cmd queries.GetUserByIDQuery) (*userDomain.User, error)
	HandleGetUserByUsername(ctx context.Context, cmd queries.GetUserByUsernameQuery) (*userDomain.User, error)
}
