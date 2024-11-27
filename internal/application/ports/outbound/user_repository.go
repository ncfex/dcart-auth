package outbound

import (
	"context"

	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type UserRepository interface {
	Add(ctx context.Context, user *userDomain.User) error
	GetByID(ctx context.Context, id string) (*userDomain.User, error)
	GetByUsername(ctx context.Context, username string) (*userDomain.User, error)
}
