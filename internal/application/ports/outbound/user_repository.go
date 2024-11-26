package outbound

import (
	"context"

	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *userDomain.User) error
	GetUserByID(ctx context.Context, userID string) (*userDomain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*userDomain.User, error)
}
