package inbound

import (
	"context"

	"github.com/ncfex/dcart-auth/internal/application/ports/types"
	"github.com/ncfex/dcart-auth/internal/application/queries"
)

type UserReadModel interface {
	GetUserByID(ctx context.Context, query queries.GetUserByIDQuery) (*types.UserResponse, error)
	GetUserByUsername(ctx context.Context, query queries.GetUserByUsernameQuery) (*types.UserResponse, error)
}
