package query

import (
	"context"

	"github.com/ncfex/dcart-auth/internal/application/ports/types"
)

type GetUserByIDQuery struct {
	UserID string
}

type GetUserByUsernameQuery struct {
	Username string
}

type UserQueryPort interface {
	GetUserByID(ctx context.Context, query GetUserByIDQuery) (*types.UserResponse, error)
	GetUserByUsername(ctx context.Context, query GetUserByUsernameQuery) (*types.UserResponse, error)
}
